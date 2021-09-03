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