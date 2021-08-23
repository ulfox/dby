package db

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type cache struct {
	v1   interface{}
	v2   map[interface{}]interface{}
	keys []string
}

type Query struct {
	KeysFound []string
	Results   []interface{}
}

type SQL struct {
	Query Query
	cache cache
}

func (d *SQL) Clear() *SQL {
	d.Query = Query{}
	d.cache = cache{}

	return d
}

func (d *SQL) clearCache() *SQL {
	d.cache = cache{}
	return d
}

func (d *SQL) dropLastKey() {
	if len(d.cache.keys) > 0 {
		d.cache.keys = d.cache.keys[:len(d.cache.keys)-1]
	}
}

func (d *SQL) dropKeys() {
	d.cache.keys = []string{}
}

func (d *SQL) getObj(k string, o interface{}) (interface{}, bool) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		return d.getArrayObject(k, o)
	}

	for thisKey, thisObj := range obj {
		d.cache.keys = append(d.cache.keys, thisKey.(string))
		if thisKey == k {
			return thisObj, true
		}

		if thisObjMap, isMap := thisObj.(map[interface{}]interface{}); isMap {
			if objFinal, found := d.getObj(k, thisObjMap); found {
				return objFinal, found
			}
			d.dropLastKey()
			continue
		}

		if arrayObj, found := d.getArrayObject(k, thisObj); found {
			return arrayObj, found
		}
		d.dropLastKey()
	}
	return nil, false
}

func (d *SQL) getArrayObject(k string, o interface{}) (interface{}, bool) {
	if arrayObj, isArray := o.([]interface{}); isArray {
		for _, thisArrayObj := range arrayObj {
			if arrayObjFinal, found := d.getObj(k, thisArrayObj); found {
				return arrayObjFinal, found
			}
		}
	}
	return nil, false
}

func (d *SQL) getFromIndex(k []string, o interface{}) (interface{}, error) {
	if getObjectType(o) != ArrayObj {
		return nil, errors.New(NotArrayObj)
	}

	i, err := getIndex(k[0])
	if err != nil {
		return nil, err
	}

	if i > len(o.([]interface{}))-1 {
		return nil, errors.New(ArrayOutOfRange)
	}

	if len(k) > 1 {
		return d.getPath(k[1:], o.([]interface{})[i])
	}

	return o.([]interface{})[i], nil
}

func (d *SQL) getPath(k []string, o interface{}) (interface{}, error) {
	obj, err := interfaceToMap(o)
	if err != nil {
		return d.getFromIndex(k, o)
	}

	for thisKey, thisObj := range obj {
		if thisKey == k[0] {
			d.cache.keys = append(d.cache.keys, k[0])
			if len(k) > 1 {
				objFinal, err := d.getPath(k[1:], thisObj)
				if err != nil {
					return objFinal, err
				}
				return objFinal, nil
			} else {
				return thisObj, nil
			}
		}
	}
	return nil, errors.New(KeyDoesNotExist)
}

func (d *SQL) delPath(k string, o interface{}) error {
	keys := strings.Split(k, ".")

	if len(keys) == 1 {
		for kn := range o.(map[interface{}]interface{}) {
			if kn.(string) == keys[0] {
				delete(o.(map[interface{}]interface{}), kn)
				break
			}
		}
	} else if len(keys) > 1 {
		d.dropKeys()
		obj, err := d.getPath(keys[:len(keys)-1], o)
		if err != nil {
			return err
		}

		d.dropKeys()
		for kn := range obj.(map[interface{}]interface{}) {
			if kn.(string) == keys[len(keys)-1] {
				delete(obj.(map[interface{}]interface{}), kn)
				break
			}
		}
	}

	return nil
}

func (d *SQL) get(k string, o interface{}) ([]string, error) {
	var err error
	var key string

	d.clearCache()
	d.cache.v1, err = copyMap(o)
	if err != nil {
		return nil, err
	}

	for {
		_, found := d.getObj(k, d.cache.v1)
		if found {
			key = strings.Join(d.cache.keys, ".")
			d.Query.KeysFound = append(d.Query.KeysFound, key)
			err := d.delPath(key, d.cache.v1)
			if err != nil {
				return d.Query.KeysFound, err
			}
			d.dropKeys()
		} else {
			break
		}
	}

	return d.Query.KeysFound, nil
}

func (d *SQL) getFirst(k string, o interface{}) (interface{}, error) {
	d.clearCache()
	obj, found := d.getObj(k, o)
	if found {
		return obj, nil
	}

	return nil, errors.New(KeyDoesNotExist)
}

func (d *SQL) upsertRecursive(k []string, o, v interface{}) error {
	var exists bool

	obj, err := interfaceToMap(o)
	if err != nil {
		return err
	}

	for thisKey, thisObj := range obj {
		if thisKey == k[0] {
			if len(k) > 1 {
				err := d.upsertRecursive(k[1:], thisObj, v)
				if err != nil {
					return err
				}
			} else {
				t := getObjectType(thisObj)
				switch t {
				case MapObj:
					for kn := range thisObj.(map[interface{}]interface{}) {
						delete(thisObj.(map[interface{}]interface{}), kn)
					}
				case ArrayObj:
					thisObj = nil
				}

				exists = false
				break
			}
			exists = true
			break
		}
	}

	if !exists {
		obj[k[0]] = make(map[interface{}]interface{})
		if len(k) > 1 {
			err := d.upsertRecursive(k[1:], obj[k[0]], v)
			if err != nil {
				return err
			}
		} else {
			obj[k[0]] = v
		}
	}

	return nil
}

func (d *SQL) upsert(k string, i, o interface{}) error {
	keys := strings.Split(k, ".")
	return d.upsertRecursive(keys, o, i)
}

func (d *SQL) mergeDBs(path string, o interface{}) error {
	var dataNew interface{}

	ok, err := fileExists(path)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New(FileNotExist)
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	yaml.Unmarshal(f, &dataNew)

	obj, err := interfaceToMap(dataNew)
	if err != nil {
		return err
	}

	for kn, vn := range obj {
		err = d.upsert(kn.(string), vn, o)
		if err != nil {
			return err
		}
	}
	return nil
}
