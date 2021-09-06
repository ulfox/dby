package db

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// AssertData is used to for converting interface objects to
// map of interfaces or array of interfaces
type AssertData struct {
	D0    map[string]string
	A0    []string
	Error error
	Cache Cache
}

// NewConvertFactory for initializing AssertData
func NewConvertFactory() *AssertData {
	assertData := &AssertData{
		D0:    make(map[string]string),
		A0:    make([]string, 0),
		Cache: NewCacheFactory(),
	}
	return assertData
}

// Clear for resetting AssertData
func (a *AssertData) Clear() {
	a.Cache.Clear()
	a.D0 = make(map[string]string)
	a.A0 = make([]string, 0)
}

// Input sets a data source that can be used for assertion
func (a *AssertData) Input(o interface{}) *AssertData {
	a.Clear()
	a.Cache.V1 = o
	return a
}

func (a *AssertData) toBytes() {
	a.Cache.B, a.Cache.E = yaml.Marshal(a.Cache.V1)
	if a.Cache.E != nil {
		a.Error = a.Cache.E
	}
}

// GetMap for converting a map[interface{}]interface{} into a map[string]string
func (a *AssertData) GetMap() map[string]string {
	if a.Cache.E != nil {
		a.Error = a.Cache.E
		return nil
	}

	a.toBytes()
	if a.Cache.E != nil {
		return nil
	}

	a.Cache.E = yaml.Unmarshal(a.Cache.B, &a.D0)
	if a.Cache.E != nil {
		a.Error = a.Cache.E
		return nil
	}
	return a.D0
}

// GetArray for converting a []interface{} to []string
func (a *AssertData) GetArray() []string {
	if a.Cache.E != nil {
		a.Error = a.Cache.E
		return nil
	}

	_, isArray := a.Cache.V1.([]interface{})
	if !isArray {
		a.Cache.E = wrapErr(fmt.Errorf(notArrayObj), getFn())
		a.Error = a.Cache.E
		return nil
	}

	a.toBytes()
	if a.Cache.E != nil {
		return nil
	}

	a.Cache.E = yaml.Unmarshal(a.Cache.B, &a.A0)
	if a.Cache.E != nil {
		a.Error = a.Cache.E
		return nil
	}

	return a.A0
}

// Key copies initial interface object and returns a map of interfaces{}
// Used to easily pipe interfaces
func (a *AssertData) Key(k string) *AssertData {
	if a.Cache.E != nil {
		a.Error = a.Cache.E
		return a
	}

	_, isMap := a.Cache.V1.(map[interface{}]interface{})
	if !isMap {
		a.Cache.E = wrapErr(fmt.Errorf(notAMap), getFn())
		a.Error = a.Cache.E
		return a
	}

	a.Cache.V1 = a.Cache.V1.(map[interface{}]interface{})[k]

	return a
}

// Index getting an interface{} from a []interface{}
func (a *AssertData) Index(i int) *AssertData {
	if a.Cache.E != nil {
		a.Error = a.Cache.E
		return a
	}

	_, isArray := a.Cache.V1.([]interface{})
	if !isArray {
		a.Cache.E = wrapErr(fmt.Errorf(notArrayObj), getFn())
		a.Error = a.Cache.E
		return a
	}
	a.Cache.V1 = a.Cache.V1.([]interface{})[i]

	return a
}
