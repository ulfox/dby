package db

import (
	"io/ioutil"
	"strconv"
	"strings"

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
		Cache: NewCacheFactory(),
	}
	return sql
}

// Clear deletes all objects from Query and Cache structures
func (s *SQL) Clear() *SQL {
	s.Query.Clear()
	s.Cache.Clear()

	return s
}

func (s *SQL) getObj(k string, o interface{}) (interface{}, bool) {
	// The object is either a map or an array.
	// If isMap returns false then check the array

	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		return s.getArrayObject(k, o)
	}

	for thisKey, thisObj := range obj {
		s.Cache.Keys = append(s.Cache.Keys, thisKey.(string))
		if thisKey == k {
			return thisObj, true
		}

		// Call self again
		if objFinal, found := s.getObj(k, thisObj); found {
			return objFinal, found
		}
		s.Cache.dropLastKey()
	}

	return nil, false
}

func (s *SQL) getArrayObject(k string, o interface{}) (interface{}, bool) {
	// This is always called after object has been
	// checked if it is a map. If isArray is false then
	// the object is neither and we should return false
	// since we do not support the required operation
	if o == nil {
		return nil, false
	}
	arrayObj, isArray := o.([]interface{})
	if !isArray {
		return nil, false
	}

	for i, thisArrayObj := range arrayObj {
		s.Cache.Keys = append(s.Cache.Keys, "["+strconv.Itoa(i)+"]")
		arrayObjFinal, found := s.getObj(k, thisArrayObj)
		if found {
			return arrayObjFinal, found
		}

		s.Cache.dropLastKey()
	}

	return nil, false
}

func (s *SQL) getIndex(k string) (int, error) {
	if !strings.HasPrefix(k, "[") || !strings.HasSuffix(k, "]") {
		return 0, wrapErr(notAnIndex, k)
	}

	intVar, err := strconv.Atoi(k[1 : len(k)-1])
	if err != nil {
		return 0, wrapErr(err)
	}
	return intVar, nil
}

func (s *SQL) getFromIndex(k []string, o interface{}) (interface{}, error) {
	if getObjectType(o) != arrayObj {
		return nil, wrapErr(notArrayObj)
	}

	i, err := s.getIndex(k[0])
	if err != nil {
		return nil, wrapErr(err)
	}

	if i > len(o.([]interface{}))-1 {
		return nil, wrapErr(
			arrayOutOfRange,
			strconv.Itoa(i),
			strconv.Itoa(len(o.([]interface{}))-1),
		)
	}

	if len(k) > 1 {
		return s.getPath(k[1:], o.([]interface{})[i])
	}

	return o.([]interface{})[i], nil
}

func (s *SQL) getPath(k []string, o interface{}) (interface{}, error) {
	if err := checkKeyPath(k); err != nil {
		return nil, wrapErr(err)
	}

	obj, err := interfaceToMap(o)
	if err != nil {
		return s.getFromIndex(k, o)
	}

	if len(k) == 0 {
		return nil, wrapErr(keyDoesNotExist, k[0])
	}

	for thisKey, thisObj := range obj {
		if thisKey != k[0] {
			continue
		}
		s.Cache.Keys = append(s.Cache.Keys, k[0])
		if len(k) == 1 {
			return thisObj, nil
		}

		objFinal, err := s.getPath(k[1:], thisObj)
		if err != nil {
			return nil, wrapErr(err)
		}
		return objFinal, nil
	}

	return nil, wrapErr(keyDoesNotExist, k[0])
}

func (s *SQL) deleteArrayItem(k string, o interface{}) bool {
	if o == nil {
		return false
	}
	for ki, kn := range o.([]interface{}) {
		if kn.(map[interface{}]interface{})[k] != nil {
			o.([]interface{})[ki] = make(map[interface{}]interface{})
			return true
		}
	}
	return false
}

func (s *SQL) deleteItem(k string, o interface{}) bool {
	_, ok := o.(map[interface{}]interface{})
	if !ok {
		return s.deleteArrayItem(k, o)
	}

	for kn := range o.(map[interface{}]interface{}) {
		if kn.(string) == k {
			delete(o.(map[interface{}]interface{}), kn)
			return true
		}
	}
	return false
}

func (s *SQL) delPath(k string, o interface{}) error {
	keys := strings.Split(k, ".")
	if err := checkKeyPath(keys); err != nil {
		return wrapErr(err)
	}

	if len(keys) == 0 {
		return wrapErr(invalidKeyPath, k)
	}

	if len(keys) == 1 {
		if !s.deleteItem(keys[0], o) {
			return wrapErr(keyDoesNotExist, k)
		}
		return nil
	}

	s.Cache.dropKeys()
	obj, err := s.getPath(keys[:len(keys)-1], o)
	if err != nil {
		return wrapErr(err)
	}

	s.Cache.dropKeys()
	if !s.deleteItem(keys[len(keys)-1], obj) {
		return wrapErr(keyDoesNotExist, k)
	}

	return nil
}

func (s *SQL) get(k string, o interface{}) ([]string, error) {
	var err error
	var key string

	s.Clear()
	s.Cache.V1, err = copyMap(o)
	if err != nil {
		return nil, wrapErr(err)
	}

	for {
		if _, found := s.getObj(k, s.Cache.V1); !found {
			break
		}

		key = strings.Join(s.Cache.Keys, ".")
		s.Query.KeysFound = append(s.Query.KeysFound, key)

		if err := s.delPath(key, s.Cache.V1); err != nil {
			return s.Query.KeysFound, wrapErr(err)
		}
		s.Cache.dropKeys()
	}

	return s.Query.KeysFound, nil
}

func (s *SQL) getFirst(k string, o interface{}) (interface{}, error) {
	s.Clear()

	keys, err := s.get(k, o)
	if err != nil {
		return nil, wrapErr(err)
	}

	if len(keys) == 0 {
		return nil, wrapErr(keyDoesNotExist, k)
	}

	keySlice := strings.Split(keys[0], ".")
	if err := checkKeyPath(keySlice); err != nil {
		return nil, wrapErr(err)
	}

	s.Cache.C1 = len(keySlice)
	if len(keys) == 1 {
		path, err := s.getPath(keySlice, o)
		return path, wrapErr(err)
	}

	for i, key := range keys[1:] {
		if len(strings.Split(key, ".")) < s.Cache.C1 {
			s.Cache.C1 = len(strings.Split(key, "."))
			s.Cache.C2 = i + 1
		}
	}

	path, err := s.getPath(strings.Split(keys[s.Cache.C2], "."), o)
	return path, wrapErr(err)
}

func (s *SQL) upsertRecursive(k []string, o, v interface{}) error {
	s.Clear()

	if err := checkKeyPath(k); err != nil {
		return wrapErr(err)
	}

	obj, err := interfaceToMap(o)
	if err != nil {
		return wrapErr(err)
	}

	for thisKey, thisObj := range obj {
		if thisKey != k[0] {
			continue
		}

		if len(k) > 1 {
			return wrapErr(s.upsertRecursive(k[1:], thisObj, v))
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
		return wrapErr(s.upsertRecursive(k[1:], obj[k[0]], v))
	}

	obj[k[0]] = v

	return nil
}

func (s *SQL) mergeDBs(path string, o interface{}) error {
	var dataNew interface{}

	ok, err := fileExists(path)
	if err != nil {
		return wrapErr(err)
	}

	if !ok {
		return wrapErr(fileNotExist, path)
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return wrapErr(err)
	}

	yaml.Unmarshal(f, &dataNew)

	obj, err := interfaceToMap(dataNew)
	if err != nil {
		return wrapErr(err)
	}

	for kn, vn := range obj {
		err = s.upsertRecursive(strings.Split(kn.(string), "."), o, vn)
		if err != nil {
			return wrapErr(err)
		}
	}
	return nil
}
