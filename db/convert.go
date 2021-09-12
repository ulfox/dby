package db

import (
	v1 "github.com/ulfox/dby/cache/v1"
	"gopkg.in/yaml.v2"
)

// AssertData is used to for converting interface objects to
// map of interfaces or array of interfaces
type AssertData struct {
	d0    map[string]string
	s0    *string
	s1    []string
	i0    *int
	i1    []int
	cache v1.Cache
}

// NewConvertFactory for initializing AssertData
func NewConvertFactory() *AssertData {
	ad := &AssertData{
		d0:    make(map[string]string),
		s1:    make([]string, 0),
		cache: v1.NewCacheFactory(),
	}
	return ad
}

// Clear for resetting AssertData
func (a *AssertData) Clear() {
	a.cache.Clear()
	a.d0 = make(map[string]string)
	a.s1 = make([]string, 0)
	a.i1 = make([]int, 0)
	a.i0 = nil
	a.s0 = nil
}

// GetError returns the any error set to AssertData
func (a *AssertData) GetError() error {
	return a.cache.E()
}

func (a *AssertData) setErr(e ...error) *AssertData {
	if len(e) > 0 {
		a.cache.E(e[0])
	}
	return a
}

// Input sets a data source that can be used for assertion
func (a *AssertData) Input(o interface{}) *AssertData {
	a.Clear()
	a.cache.V1(o)
	return a
}

// GetString asserts the input as string
func (a *AssertData) GetString() (string, error) {
	if a.GetError() != nil {
		return "", a.GetError()
	}

	s, isString := a.cache.V1().(string)
	if !isString {
		a.setErr(wrapErr(notAType, "string"))
		return "", a.GetError()
	}

	return s, nil
}

// GetInt asserts the input as int
func (a *AssertData) GetInt() (int, error) {
	if a.GetError() != nil {
		return 0, a.GetError()
	}

	i, isInt := a.cache.V1().(int)
	if !isInt {
		a.setErr(wrapErr(notAType, "int"))
		return 0, a.GetError()
	}

	return i, nil
}

// GetMap for converting a map[interface{}]interface{} into a map[string]string
func (a *AssertData) GetMap() (map[string]string, error) {
	if a.GetError() != nil {
		return nil, a.GetError()
	}

	a.cache.E(a.cache.BE(yaml.Marshal(a.cache.V1())))
	if a.GetError() != nil {
		return nil, a.GetError()
	}

	a.cache.E(yaml.Unmarshal(a.cache.B(), &a.d0))
	if a.GetError() != nil {
		return nil, a.GetError()
	}
	return a.d0, nil
}

// GetArray for converting a []interface{} to []string
func (a *AssertData) GetArray() ([]string, error) {
	if a.GetError() != nil {
		return nil, a.GetError()
	}

	_, isArray := a.cache.V1().([]interface{})
	if !isArray {
		a.setErr(wrapErr(notArrayObj))
		return nil, a.GetError()
	}

	a.cache.E(a.cache.BE(yaml.Marshal(a.cache.V1())))
	if a.GetError() != nil {
		return nil, a.GetError()
	}

	a.cache.E(yaml.Unmarshal(a.cache.B(), &a.s1))
	if a.GetError() != nil {
		return nil, a.GetError()
	}

	return a.s1, nil
}

// Key copies initial interface object and returns a map of interfaces{}
// Used to easily pipe interfaces
func (a *AssertData) Key(k string) *AssertData {
	if a.GetError() != nil {
		return a
	}

	_, isMap := a.cache.V1().(map[interface{}]interface{})
	if !isMap {
		return a.setErr(wrapErr(notAMap))
	}

	a.cache.V1(a.cache.V1().(map[interface{}]interface{})[k])

	return a
}

// Index getting an interface{} from a []interface{}
func (a *AssertData) Index(i int) *AssertData {
	if a.GetError() != nil {
		return a
	}

	_, isArray := a.cache.V1().([]interface{})
	if !isArray {
		return a.setErr(wrapErr(notArrayObj))
	}
	a.cache.V1(a.cache.V1().([]interface{})[i])

	return a
}
