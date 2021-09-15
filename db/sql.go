package db

import (
	"io/ioutil"
	"strconv"
	"strings"

	v1 "github.com/ulfox/dby/cache/v1"
	v2 "github.com/ulfox/dby/cache/v2"
	"gopkg.in/yaml.v2"
)

// SQL is the core struct for working with maps.
type SQL struct {
	v1.Cache
	v2.Query
}

// NewSQLFactory creates a new empty SQL
func NewSQLFactory() *SQL {
	sql := &SQL{
		Query: v2.NewQueryFactory(),
		Cache: v1.NewCacheFactory(),
	}
	return sql
}

// Clear deletes all objects from Query and Cache structures
func (s *SQL) Clear() *SQL {
	s.Query.Clear()
	s.Cache.Clear()

	return s
}

func (s *SQL) getObj(k string, o *interface{}) (*interface{}, bool) {
	_, isMap := (*o).(map[interface{}]interface{})
	if !isMap {
		return s.getArrayObject(k, o)
	}

	for thisKey, thisObj := range (*o).(map[interface{}]interface{}) {
		s.Cache.AddKey(thisKey.(string))
		if thisKey == k {
			return &thisObj, true
		}

		if objFinal, found := s.getObj(k, &thisObj); found {
			return objFinal, found
		}
		s.Cache.DropLastKey()
	}

	return nil, false
}

func (s *SQL) getArrayObject(k string, o *interface{}) (*interface{}, bool) {
	if o == nil {
		return nil, false
	}
	_, isArray := (*o).([]interface{})
	if !isArray {
		return nil, false
	}

	for i, thisArrayObj := range (*o).([]interface{}) {
		s.Cache.AddKey("[" + strconv.Itoa(i) + "]")
		arrayObjFinal, found := s.getObj(k, &thisArrayObj)
		if found {
			return arrayObjFinal, found
		}

		s.Cache.DropLastKey()
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

func (s *SQL) getFromIndex(k []string, o *interface{}) (*interface{}, error) {
	_, isArray := (*o).([]interface{})
	if !isArray {
		return nil, wrapErr(notArrayObj)
	}
	v := (*o).([]interface{})

	i, err := s.getIndex(k[0])
	if err != nil {
		return nil, wrapErr(err)
	}

	if i > len((*o).([]interface{}))-1 {
		return nil, wrapErr(
			arrayOutOfRange,
			strconv.Itoa(i),
			strconv.Itoa(len((*o).([]interface{}))-1),
		)
	}

	if len(k) > 1 {
		return s.getPath(k[1:], &v[i])
	}

	return &v[i], nil
}

func (s *SQL) getPath(k []string, o *interface{}) (*interface{}, error) {
	if err := checkKeyPath(k); err != nil {
		return nil, wrapErr(err)
	}

	_, ok := (*o).(map[interface{}]interface{})
	if !ok {
		return s.getFromIndex(k, o)
	}

	if len(k) == 0 {
		return nil, wrapErr(keyDoesNotExist, k[0])
	}

	for thisKey, thisObj := range (*o).(map[interface{}]interface{}) {
		if thisKey != k[0] {
			continue
		}
		s.Cache.AddKey(k[0])
		if len(k) == 1 {
			return &thisObj, nil
		}

		objFinal, err := s.getPath(k[1:], &thisObj)
		if err != nil {
			return nil, wrapErr(err)
		}
		return objFinal, nil
	}

	return nil, wrapErr(keyDoesNotExist, k[0])
}

func (s *SQL) deleteArrayItem(k string, o *interface{}) error {
	if o == nil {
		return wrapErr(notArrayObj)
	}

	i, err := s.getIndex(k)
	if err != nil {
		return wrapErr(err)
	}

	(*o).([]interface{})[i] = (*o).([]interface{})[len((*o).([]interface{}))-1]
	(*o).([]interface{})[len((*o).([]interface{}))-1] = ""
	*o = (*o).([]interface{})[:len((*o).([]interface{}))-1]

	return nil
}

func (s *SQL) deleteItem(k string, o *interface{}) error {
	_, ok := (*o).(map[interface{}]interface{})
	if !ok {
		return s.deleteArrayItem(k, o)
	}

	for kn := range (*o).(map[interface{}]interface{}) {
		if kn.(string) == k {
			delete((*o).(map[interface{}]interface{}), kn)
			return nil
		}
	}
	return wrapErr(keyDoesNotExist, k)
}

func (s *SQL) delPath(k string, o *interface{}) error {
	keys := strings.Split(k, ".")
	if err := checkKeyPath(keys); err != nil {
		return wrapErr(err)
	}

	if len(keys) == 0 {
		return wrapErr(invalidKeyPath, k)
	}

	if len(keys) == 1 {
		if err := s.deleteItem(keys[0], o); err != nil {
			return wrapErr(keyDoesNotExist, k)
		}
		return nil
	}

	s.Cache.DropKeys()
	obj, err := s.getPath(keys[:len(keys)-1], o)
	if err != nil {
		return wrapErr(err)
	}

	s.Cache.DropKeys()
	if err := s.deleteItem(keys[len(keys)-1], obj); err != nil {
		return wrapErr(keyDoesNotExist, k)
	}

	return nil
}

func (s *SQL) findKeys(k string, o *interface{}) ([]string, error) {
	var err error
	var key string

	s.Clear()
	err = s.V1E(copyMap(*o))
	if err != nil {
		return nil, wrapErr(err)
	}

	for {
		obj := s.V1()
		if _, found := s.getObj(k, &obj); !found {
			break
		}

		key = strings.Join(s.Cache.GetKeys(), ".")
		s.Query.AddKey(key)

		if err := s.delPath(key, &obj); err != nil {
			return s.Query.GetKeys(), wrapErr(err)
		}
		s.Cache.DropKeys()
	}

	return s.Query.GetKeys(), nil
}

func (s *SQL) getFirst(k string, o *interface{}) (*interface{}, error) {
	s.Clear()

	keys, err := s.findKeys(k, o)
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

	s.Cache.C1(len(keySlice))
	if len(keys) == 1 {
		path, err := s.getPath(keySlice, o)
		return path, wrapErr(err)
	}

	for i, key := range keys[1:] {
		if len(strings.Split(key, ".")) < s.Cache.C1() {
			s.Cache.C1(len(strings.Split(key, ".")))
			s.Cache.C2(i + 1)
		}
	}

	path, err := s.getPath(strings.Split(keys[s.Cache.C2()], "."), o)
	return path, wrapErr(err)
}

func (s *SQL) toInterfaceMap(v interface{}) (interface{}, error) {
	var dataNew interface{}

	if v == nil {
		return make(map[interface{}]interface{}), nil
	}

	dataBytes, err := yaml.Marshal(&v)
	if err != nil {
		return nil, wrapErr(err)
	}
	err = yaml.Unmarshal(dataBytes, &dataNew)
	if err != nil {
		return nil, wrapErr(err)
	}

	return dataNew, nil
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

	obj[k[0]] = emptyMap()

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
