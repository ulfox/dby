package db

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Cache for easily sharing state between map operations and methods
// v1 & v2 are common interface{} placeholders while keys is used by
// path discovery methods to keep track and derive the right path.
type Cache struct {
	v1   interface{}
	v2   map[interface{}]interface{}
	keys []string
}

// Query hosts results from SQL methods
type Query struct {
	KeysFound []string
	Results   []interface{}
}

// SQL is the core struct for working with maps.
type SQL struct {
	Query Query
	Cache Cache
}

// Clear deletes all objects from Query and Cache structures
func (d *SQL) Clear() *SQL {
	d.Query = Query{}
	d.Cache = Cache{}

	return d
}

func (d *SQL) clearCache() *SQL {
	d.Cache = Cache{}
	return d
}

func (d *SQL) dropLastKey() {
	if len(d.Cache.keys) > 0 {
		d.Cache.keys = d.Cache.keys[:len(d.Cache.keys)-1]
	}
}

func (d *SQL) dropKeys() {
	d.Cache.keys = []string{}
}

func (d *SQL) getObj(k string, o interface{}) (interface{}, bool) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		return d.getArrayObject(k, o)
	}

	for thisKey, thisObj := range obj {
		d.Cache.keys = append(d.Cache.keys, thisKey.(string))
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
	if getObjectType(o) != arrayObj {
		return nil, errors.Wrap(errors.New(notArrayObj), "getFromIndex")
	}

	i, err := getIndex(k[0])
	if err != nil {
		return nil, errors.Wrap(err, "getFromIndex")
	}

	if i > len(o.([]interface{}))-1 {
		return nil, errors.Wrap(errors.New(arrayOutOfRange), "getFromIndex")
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
			d.Cache.keys = append(d.Cache.keys, k[0])
			if len(k) > 1 {
				objFinal, err := d.getPath(k[1:], thisObj)
				if err != nil {
					return objFinal, errors.Wrap(err, "getPath")
				}
				return objFinal, nil
			}
			return thisObj, nil
		}
	}
	return nil, errors.Wrap(errors.New(keyDoesNotExist), "getPath")
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
			return errors.Wrap(err, "delPath")
		}

		d.dropKeys()
		deleted := false
		for kn := range obj.(map[interface{}]interface{}) {
			if kn.(string) == keys[len(keys)-1] {
				delete(obj.(map[interface{}]interface{}), kn)
				deleted = true
				break
			}
		}

		if !deleted {
			return errors.Wrap(errors.New(keyDoesNotExist), "delPath")
		}
	}

	return nil
}

func (d *SQL) get(k string, o interface{}) ([]string, error) {
	var err error
	var key string

	d.clearCache()
	d.Cache.v1, err = copyMap(o)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}

	for {
		_, found := d.getObj(k, d.Cache.v1)
		if found {
			key = strings.Join(d.Cache.keys, ".")
			d.Query.KeysFound = append(d.Query.KeysFound, key)
			err := d.delPath(key, d.Cache.v1)
			if err != nil {
				return d.Query.KeysFound, errors.Wrap(err, "get")
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

	return nil, errors.Wrap(errors.New(keyDoesNotExist), "getFirst")
}

func (d *SQL) upsertRecursive(k []string, o, v interface{}) error {
	var exists bool

	obj, err := interfaceToMap(o)
	if err != nil {
		return errors.Wrap(err, "upsertRecursive")
	}

	for thisKey, thisObj := range obj {
		if thisKey == k[0] {
			if len(k) > 1 {
				err := d.upsertRecursive(k[1:], thisObj, v)
				if err != nil {
					return errors.Wrap(err, "upsertRecursive")
				}
			} else {
				t := getObjectType(thisObj)
				switch t {
				case mapObj:
					for kn := range thisObj.(map[interface{}]interface{}) {
						delete(thisObj.(map[interface{}]interface{}), kn)
					}
				case arrayObj:
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
				return errors.Wrap(err, "upsertRecursive")
			}
		} else {
			obj[k[0]] = v
		}
	}

	return nil
}

func (d *SQL) upsert(k string, i, o interface{}) error {
	keys := strings.Split(k, ".")
	return errors.Wrap(d.upsertRecursive(keys, o, i), "upsert")
}

func (d *SQL) mergeDBs(path string, o interface{}) error {
	var dataNew interface{}

	ok, err := fileExists(path)
	if err != nil {
		return errors.Wrap(err, "mergeDBs")
	}

	if !ok {
		return errors.Wrap(errors.New(fileNotExist), "mergeDBs")
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "mergeDBs")
	}

	yaml.Unmarshal(f, &dataNew)

	obj, err := interfaceToMap(dataNew)
	if err != nil {
		return errors.Wrap(err, "mergeDBs")
	}

	for kn, vn := range obj {
		err = d.upsert(kn.(string), vn, o)
		if err != nil {
			return errors.Wrap(err, "mergeDBs")
		}
	}
	return nil
}
