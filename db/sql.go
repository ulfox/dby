package db

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Query hosts results from SQL methods
type Query struct {
	KeysFound []string
	Results   interface{}
}

// Clear deletes all objects from Query
func (q *Query) Clear() *Query {
	q.KeysFound = nil
	q.Results = nil

	return q
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
	d.Query.Clear()
	d.Cache.Clear()

	return d
}

func (d *SQL) getObj(k string, o interface{}) (interface{}, bool) {
	// The object is either a map or an array.
	// If isMap returns false then check the array
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap && o != nil {
		return d.getArrayObject(k, o)
	}

	for thisKey, thisObj := range obj {
		d.Cache.Keys = append(d.Cache.Keys, thisKey.(string))
		if thisKey == k {
			return thisObj, true
		}

		// Call self again
		if objFinal, found := d.getObj(k, thisObj); found {
			return objFinal, found
		}
		d.Cache.dropLastKey()
	}
	return nil, false
}

func (d *SQL) getArrayObject(k string, o interface{}) (interface{}, bool) {
	// This is always called after object has been
	// checked if it is a map. If isArray is false then
	// the object is neither and we should return false
	// since we do not support the required operation
	arrayObj, isArray := o.([]interface{})
	if !isArray {
		return nil, false
	}
	for _, thisArrayObj := range arrayObj {
		if arrayObjFinal, found := d.getObj(k, thisArrayObj); found {
			return arrayObjFinal, found
		}
	}
	return nil, false
}

func (d *SQL) getIndex(k string) (int, error) {
	if !strings.HasPrefix(k, "[") || !strings.HasSuffix(k, "]") {
		return 0, wrapErr(fmt.Errorf(notAnIndex, k), getFn())
	}

	intVar, err := strconv.Atoi(k[1 : len(k)-1])
	if err != nil {
		return 0, wrapErr(err, getFn())
	}
	return intVar, nil
}

func (d *SQL) getFromIndex(k []string, o interface{}) (interface{}, error) {
	if getObjectType(o) != arrayObj {
		return nil, wrapErr(errors.New(notArrayObj), getFn())
	}

	i, err := d.getIndex(k[0])
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	if i > len(o.([]interface{}))-1 {
		return nil, wrapErr(
			fmt.Errorf(
				arrayOutOfRange,
				strconv.Itoa(i),
				strconv.Itoa(len(o.([]interface{}))-1),
			),
			getFn(),
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
		return nil, wrapErr(fmt.Errorf(keyDoesNotExist, k[0]), getFn())
	}

	for thisKey, thisObj := range obj {
		if thisKey != k[0] {
			continue
		}
		d.Cache.Keys = append(d.Cache.Keys, k[0])
		if len(k) == 1 {
			return thisObj, nil
		}

		objFinal, err := d.getPath(k[1:], thisObj)
		if err != nil {
			return nil, wrapErr(err, getFn())
		}
		return objFinal, nil
	}

	return nil, wrapErr(fmt.Errorf(keyDoesNotExist, k[0]), getFn())
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
		return wrapErr(fmt.Errorf(invalidKeyPath, k), getFn())
	}

	if len(keys) == 1 {
		if !d.deleteItem(keys[0], o) {
			return wrapErr(fmt.Errorf(keyDoesNotExist, k), getFn())
		}
		return nil
	}

	d.Cache.dropKeys()
	obj, err := d.getPath(keys[:len(keys)-1], o)
	if err != nil {
		return wrapErr(err, getFn())
	}

	d.Cache.dropKeys()
	if !d.deleteItem(keys[len(keys)-1], obj) {
		return wrapErr(fmt.Errorf(keyDoesNotExist, k), getFn())
	}

	return nil
}

func (d *SQL) get(k string, o interface{}) ([]string, error) {
	var err error
	var key string

	d.Clear()
	d.Cache.V1, err = copyMap(o)
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	for {
		if _, found := d.getObj(k, d.Cache.V1); !found {
			break
		}

		key = strings.Join(d.Cache.Keys, ".")
		d.Query.KeysFound = append(d.Query.KeysFound, key)

		if err := d.delPath(key, d.Cache.V1); err != nil {
			return d.Query.KeysFound, wrapErr(err, getFn())
		}
		d.Cache.dropKeys()
	}

	return d.Query.KeysFound, nil
}

func (d *SQL) getFirst(k string, o interface{}) (interface{}, error) {
	d.Clear()

	obj, err := d.get(k, o)
	if err != nil {
		return nil, wrapErr(fmt.Errorf(keyDoesNotExist, k), getFn())
	}

	if len(obj) == 0 {
		return nil, wrapErr(fmt.Errorf(keyDoesNotExist, k), getFn())
	}

	d.Cache.C1 = len(strings.Split(obj[0], "."))
	if len(obj) == 1 {
		return d.getPath(strings.Split(obj[0], "."), o)
	}
	for i, key := range obj {
		if len(strings.Split(key, ".")) < d.Cache.C1 {
			d.Cache.C1 = len(strings.Split(key, "."))
			d.Cache.C2 = i
		}
	}

	return d.getPath(strings.Split(obj[d.Cache.C2], "."), o)
}

func (d *SQL) upsertRecursive(k []string, o, v interface{}) error {
	d.Clear()

	obj, err := interfaceToMap(o)
	if err != nil {
		return wrapErr(err, getFn())
	}

	for thisKey, thisObj := range obj {
		if thisKey != k[0] {
			continue
		}

		if len(k) > 1 {
			return wrapErr(d.upsertRecursive(k[1:], thisObj, v), getFn())
		}

		switch getObjectType(thisObj) {
		case mapObj:
			deleteMap(thisObj)
		case arrayObj:
			thisObj = nil
		}

		break
	}

	obj[k[0]] = make(map[interface{}]interface{})

	if len(k) > 1 {
		return wrapErr(d.upsertRecursive(k[1:], obj[k[0]], v), getFn())
	}

	obj[k[0]] = v

	return nil
}

func (d *SQL) mergeDBs(path string, o interface{}) error {
	var dataNew interface{}

	ok, err := fileExists(path)
	if err != nil {
		return wrapErr(err, getFn())
	}

	if !ok {
		return wrapErr(fmt.Errorf(fileNotExist, path), getFn())
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return wrapErr(err, getFn())
	}

	yaml.Unmarshal(f, &dataNew)

	obj, err := interfaceToMap(dataNew)
	if err != nil {
		return wrapErr(err, getFn())
	}

	for kn, vn := range obj {
		err = d.upsertRecursive(strings.Split(kn.(string), "."), o, vn)
		if err != nil {
			return wrapErr(err, getFn())
		}
	}
	return nil
}
