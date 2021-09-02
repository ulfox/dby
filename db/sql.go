package db

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Cache for easily sharing state between map operations and methods
// v1 & v2 are common interface{} placeholders while keys is used by
// path discovery methods to keep track and derive the right path.
type Cache struct {
	V1   interface{}
	V2   int
	Keys []string
}

// Query hosts results from SQL methods
type Query struct {
	KeysFound []string
	Results   interface{}
}

// SQL is the core struct for working with maps.
type SQL struct {
	Query Query
	Cache Cache
}

// NewSQLFactory creates a new empty SQL
func NewSQLFactory() *SQL {
	sql := &SQL{
		Query: Query{},
		Cache: Cache{},
	}
	return sql
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
	if len(d.Cache.Keys) > 0 {
		d.Cache.Keys = d.Cache.Keys[:len(d.Cache.Keys)-1]
	}
}

func (d *SQL) dropKeys() {
	d.Cache.Keys = []string{}
}

func (d *SQL) getObj(k string, o interface{}) (interface{}, bool) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		return d.getArrayObject(k, o)
	}

	for thisKey, thisObj := range obj {
		d.Cache.Keys = append(d.Cache.Keys, thisKey.(string))
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
		return d.loopArray(k, arrayObj)
	}
	return nil, false
}

func (d *SQL) loopArray(k string, o []interface{}) (interface{}, bool) {
	for _, thisArrayObj := range o {
		if arrayObjFinal, found := d.getObj(k, thisArrayObj); found {
			return arrayObjFinal, found
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
		return nil, errors.New(
			fmt.Sprintf(
				arrayOutOfRange,
				strconv.Itoa(i),
				strconv.Itoa(len(o.([]interface{}))-1),
			),
		)
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

	if len(k) == 0 {
		return nil, errors.New(fmt.Sprintf(keyDoesNotExist, k[0]))
	}

	for thisKey, thisObj := range obj {
		if thisKey == k[0] {
			d.Cache.Keys = append(d.Cache.Keys, k[0])
			if len(k) == 1 {
				return thisObj, nil
			}

			objFinal, err := d.getPath(k[1:], thisObj)
			if err != nil {
				return nil, errors.Wrap(err, "getPath")
			}
			return objFinal, nil
		}
	}

	return nil, errors.New(fmt.Sprintf(keyDoesNotExist, k[0]))
}

func (d *SQL) deleteItem(k string, o interface{}) bool {
	for kn := range o.(map[interface{}]interface{}) {
		if kn.(string) == k {
			delete(o.(map[interface{}]interface{}), kn)
			return true
		}
	}
	return false
}

func (d *SQL) delPath(k string, o interface{}) error {
	keys := strings.Split(k, ".")

	if len(keys) == 0 {
		return errors.New(fmt.Sprintf(invalidKeyPath, k))
	}

	if len(keys) == 1 {
		if !d.deleteItem(keys[0], o) {
			return errors.New(fmt.Sprintf(keyDoesNotExist, k))
		}
		return nil
	}

	d.dropKeys()
	obj, err := d.getPath(keys[:len(keys)-1], o)
	if err != nil {
		return errors.Wrap(err, "delPath")
	}

	d.dropKeys()
	if !d.deleteItem(keys[len(keys)-1], obj) {
		return errors.New(fmt.Sprintf(keyDoesNotExist, k))
	}

	return nil
}

func (d *SQL) get(k string, o interface{}) ([]string, error) {
	var err error
	var key string

	d.Clear()
	d.Cache.V1, err = copyMap(o)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}

	for {
		if _, found := d.getObj(k, d.Cache.V1); !found {
			break
		}

		key = strings.Join(d.Cache.Keys, ".")
		d.Query.KeysFound = append(d.Query.KeysFound, key)

		if err := d.delPath(key, d.Cache.V1); err != nil {
			return d.Query.KeysFound, errors.Wrap(err, "get")
		}
		d.dropKeys()
	}

	return d.Query.KeysFound, nil
}

func (d *SQL) getFirst(k string, o interface{}) (interface{}, error) {
	d.Clear()

	obj, err := d.get(k, o)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(keyDoesNotExist, k))
	}

	if len(obj) == 0 {
		return nil, errors.New(fmt.Sprintf(keyDoesNotExist, k))
	}

	if len(obj) == 1 {
		return d.getPath(strings.Split(obj[0], "."), o)
	}

	for i, key := range obj[1:] {
		d.dropKeys()
		if len(strings.Split(key, ".")) < len(strings.Split(obj[0], ".")) {
			d.dropKeys()
			d.Cache.V2 = i
		}
	}

	return d.getPath(strings.Split(obj[d.Cache.V2], "."), o)
}

func (d *SQL) upsertRecursive(k []string, o, v interface{}) error {
	d.Clear()

	obj, err := interfaceToMap(o)
	if err != nil {
		return errors.Wrap(err, "upsertRecursive")
	}

	for thisKey, thisObj := range obj {
		if thisKey == k[0] {
			if len(k) > 1 {
				return errors.Wrap(
					d.upsertRecursive(k[1:], thisObj, v),
					"upsertRecursive",
				)
			}

			switch getObjectType(thisObj) {
			case mapObj:
				for kn := range thisObj.(map[interface{}]interface{}) {
					delete(thisObj.(map[interface{}]interface{}), kn)
				}
			case arrayObj:
				thisObj = nil
			}

			break
		}
	}

	obj[k[0]] = make(map[interface{}]interface{})

	if len(k) > 1 {
		return errors.Wrap(
			d.upsertRecursive(k[1:], obj[k[0]], v),
			"upsertRecursive",
		)
	}

	obj[k[0]] = v

	return nil
}

func (d *SQL) mergeDBs(path string, o interface{}) error {
	var dataNew interface{}

	ok, err := fileExists(path)
	if err != nil {
		return errors.Wrap(err, "mergeDBs")
	}

	if !ok {
		return errors.New(fmt.Sprintf(fileNotExist, path))
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
		err = d.upsertRecursive(strings.Split(kn.(string), "."), o, vn)
		if err != nil {
			return errors.Wrap(err, "mergeDBs")
		}
	}
	return nil
}
