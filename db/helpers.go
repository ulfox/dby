package db

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func getFn() string {
	pc, _, no, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return fmt.Sprintf("%s#%s\n", details.Name(), strconv.Itoa(no))
	}
	return ""
}

func getIndex(k string) (int, error) {
	if strings.HasPrefix(k, "[") && strings.HasSuffix(k, "]") {
		intVar, err := strconv.Atoi(k[1 : len(k)-1])
		if err != nil {
			return 0, wrapErr(err, getFn())
		}
		return intVar, nil
	}
	return 0, wrapErr(fmt.Errorf(notAnIndex, k), getFn())
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
		return nil, wrapErr(err, getFn())
	}

	var cache interface{}

	data, err := yaml.Marshal(&obj)
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	err = yaml.Unmarshal(data, &cache)
	if err != nil {
		return nil, wrapErr(err, getFn())
	}

	return cache, nil
}

func interfaceToMap(o interface{}) (map[interface{}]interface{}, error) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		if o != nil {
			return nil, wrapErr(errors.New(notAMap), getFn())
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
			return wrapErr(err, getFn())
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
		return false, wrapErr(fmt.Errorf(dictNotFile, filepath), getFn())
	}

	return true, nil
}
