package db

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func checkKeyPath(k []string) error {
	for _, j := range k {
		if j == "" {
			return fmt.Errorf(emptyKey, strings.Join(k, "."))
		}
	}
	return nil
}

func copyMap(o interface{}) (interface{}, error) {
	obj, err := interfaceToMap(o)
	if err != nil {
		return nil, wrapErr(err)
	}

	var cache interface{}

	data, err := yaml.Marshal(&obj)
	if err != nil {
		return nil, wrapErr(err)
	}

	err = yaml.Unmarshal(data, &cache)
	if err != nil {
		return nil, wrapErr(err)
	}

	return cache, nil
}

func emptyMap() map[interface{}]interface{} {
	return make(map[interface{}]interface{})
}

func interfaceToMap(o interface{}) (map[interface{}]interface{}, error) {
	obj, isMap := o.(map[interface{}]interface{})
	if !isMap {
		if o != nil {
			return nil, wrapErr(notAMap)
		}
		obj = emptyMap()
	}
	return obj, nil
}

func deleteMap(o interface{}) {
	for kn := range o.(map[interface{}]interface{}) {
		delete(o.(map[interface{}]interface{}), kn)
	}
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
			return wrapErr(err)
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
		return false, wrapErr(dictNotFile, filepath)
	}

	return true, nil
}
