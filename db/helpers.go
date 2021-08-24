package db

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func getIndex(k string) (int, error) {
	if strings.HasPrefix(k, "[") && strings.HasSuffix(k, "]") {
		intVar, err := strconv.Atoi(k[1 : len(k)-1])
		if err != nil {
			return 0, errors.Wrap(err, "getIndex")
		}
		return intVar, nil
	}
	return 0, errors.Wrap(errors.New(notAnIndex), "getIndex")
}

// Some common objects
func getObjectType(o interface{}) objectType {
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

	return unknownObj
}

func copyMap(o interface{}) (interface{}, error) {
	obj, err := interfaceToMap(o)
	if err != nil {
		return nil, errors.Wrap(err, "copyMap")
	}

	var cache interface{}

	data, err := yaml.Marshal(&obj)
	if err != nil {
		return nil, errors.Wrap(err, "copyMap")
	}

	err = yaml.Unmarshal(data, &cache)
	if err != nil {
		return nil, errors.Wrap(err, "copyMap")
	}

	return cache, nil
}

func interfaceToMap(o interface{}) (map[interface{}]interface{}, error) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		if o != nil {
			return nil, errors.Wrap(errors.New(notAMap), "interfaceToMap")
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
			return errors.Wrap(err, "makeDirs")
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
		return false, errors.Wrap(errors.New(dictNotFile), "fileExists")
	}

	return true, nil
}
