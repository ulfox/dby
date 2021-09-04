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
func (s *AssertData) Clear() {
	s.Cache.Clear()
	s.D0 = make(map[string]string)
	s.A0 = make([]string, 0)
}

// Input sets a data source that can be used for assertion
func (s *AssertData) Input(o interface{}) *AssertData {
	s.Clear()
	s.Cache.V1 = o
	return s
}

func (s *AssertData) toBytes() {
	s.Cache.B, s.Cache.E = yaml.Marshal(s.Cache.V1)
	if s.Cache.E != nil {
		s.Error = s.Cache.E
	}
}

// GetMap for converting a map[interface{}]interface{} into a map[string]string
func (s *AssertData) GetMap() map[string]string {
	if s.Cache.E != nil {
		s.Error = s.Cache.E
		return nil
	}

	s.toBytes()
	if s.Cache.E != nil {
		return nil
	}

	s.Cache.E = yaml.Unmarshal(s.Cache.B, &s.D0)
	if s.Cache.E != nil {
		s.Error = s.Cache.E
		return nil
	}
	return s.D0
}

// GetArray for converting a []interface{} to []string
func (s *AssertData) GetArray() []string {
	if s.Cache.E != nil {
		s.Error = s.Cache.E
		return nil
	}

	_, isArray := s.Cache.V1.([]interface{})
	if !isArray {
		s.Cache.E = wrapErr(fmt.Errorf(notArrayObj), getFn())
		s.Error = s.Cache.E
		return nil
	}

	s.toBytes()
	if s.Cache.E != nil {
		return nil
	}

	s.Cache.E = yaml.Unmarshal(s.Cache.B, &s.A0)
	if s.Cache.E != nil {
		s.Error = s.Cache.E
		return nil
	}

	return s.A0
}

// Key copies initial interface object and returns a map of interfaces{}
// Used to easily pipe interfaces
func (s *AssertData) Key(k string) *AssertData {
	if s.Cache.E != nil {
		s.Error = s.Cache.E
		return s
	}

	_, isMap := s.Cache.V1.(map[interface{}]interface{})
	if !isMap {
		s.Cache.E = wrapErr(fmt.Errorf(notAMap), getFn())
		s.Error = s.Cache.E
		return s
	}

	s.Cache.V1 = s.Cache.V1.(map[interface{}]interface{})[k]

	return s
}

// Index getting an interface{} from a []interface{}
func (s *AssertData) Index(i int) *AssertData {
	if s.Cache.E != nil {
		s.Error = s.Cache.E
		return s
	}

	_, isArray := s.Cache.V1.([]interface{})
	if !isArray {
		s.Cache.E = wrapErr(fmt.Errorf(notArrayObj), getFn())
		s.Error = s.Cache.E
		return s
	}
	s.Cache.V1 = s.Cache.V1.([]interface{})[i]

	return s
}
