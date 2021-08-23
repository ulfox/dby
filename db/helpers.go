package db

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Some common objects
func getObjectType(o interface{}) ObjectType {
	_, isMap := o.(map[interface{}]interface{})
	if isMap {
		return 1
	}

	_, isArray := o.([]interface{})
	if isArray {
		return 2
	}

	_, isArrayMap := o.([]map[interface{}]interface{})
	if isArrayMap {
		return 3
	}

	_, isMapStringString := o.(map[string]string)
	if isMapStringString {
		return 4
	}

	_, isMapStringInterface := o.(map[string]interface{})
	if isMapStringInterface {
		return 5
	}

	_, isMapStringArrayString := o.(map[string][]string)
	if isMapStringArrayString {
		return 6
	}

	_, isArrayMapStringArrayString := o.([]map[string][]string)
	if isArrayMapStringArrayString {
		return 7
	}

	_, isArrayMapStringString := o.([]map[string]string)
	if isArrayMapStringString {
		return 8
	}

	_, isArrayMapStringInterface := o.([]map[string]interface{})
	if isArrayMapStringInterface {
		return 9
	}

	_, isArrayMapStringArrayInterface := o.([]map[string][]interface{})
	if isArrayMapStringArrayInterface {
		return 10
	}

	return UnknownObj
}

func copyMap(o interface{}) (interface{}, error) {
	obj, err := interfaceToMap(o)
	if err != nil {
		return nil, err
	}

	var cache interface{}

	data, err := yaml.Marshal(&obj)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func interfaceToMap(o interface{}) (map[interface{}]interface{}, error) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		if o != nil {
			return nil, errors.New(NotAMap)
		}
		obj = make(map[interface{}]interface{})
	}
	return obj, nil
}

// makeDirs create directories if they do not exist
func makeDirs(p string, m os.FileMode) error {
	if p == "/" {
		// Nothing to do here
		return nil
	}

	p = strings.TrimSuffix(p, "/")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, m)
		if err != nil {
			return err
		}
	}
	return nil
}

// fileExists for checking if a file exists
func fileExists(filepath string) (bool, error) {
	f, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false, nil
	} else if f.IsDir() {
		return false, errors.New(DictNotFile)
	}

	return true, nil
}
