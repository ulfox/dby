package v2

import (
	"strings"
)

// Query hosts results from SQL methods
type Query struct {
	cache map[int]map[string]map[string]*interface{}
	first map[int]map[string]string
	index int
}

// NewQueryFactory for creating a new Query
func NewQueryFactory() Query {
	Query := Query{
		cache: make(map[int]map[string]map[string]*interface{}),
		first: make(map[int]map[string]string),
	}
	return Query
}

// UpdateQCache for storing paths for a given key and the related
// path objects
func (c *Query) UpdateQCache(path string, o *interface{}) {
	slicedPath := strings.Split(path, ".")
	key := slicedPath[len(slicedPath)-1]

	if _, ok := c.cache[c.index][key]; !ok {
		c.cache[c.index][key] = map[string]*interface{}{
			path: o,
		}
		c.first[c.index][key] = path
		return
	}

	c.cache[c.index][key][path] = o

	first := c.first[c.index][key]
	if len(strings.Split(first, ".")) >= len(slicedPath) {
		c.first[c.index][key] = path
	}
}

// DeleteQCachePath for removing a path from the Query cache
func (c *Query) DeleteQCachePath(path string) {
	slicedPath := strings.Split(path, ".")
	key := slicedPath[len(slicedPath)-1]

	if _, ok := c.cache[c.index][key]; !ok {
		return
	}

	delete(c.cache[c.index][key], path)

	first, ok := c.first[c.index][key]
	if !ok {
		return
	}

	if first == path {
		delete(c.cache[c.index], key)
	}
}

// SetQIndex for setting cache index
func (c *Query) SetQIndex(i int) {
	c.index = i
	if _, ok := c.cache[c.index]; !ok {
		c.cache[c.index] = make(map[string]map[string]*interface{})
	}
	if _, ok := c.first[c.index]; !ok {
		c.first[c.index] = make(map[string]string)
	}
}

// GetQCache for returning the Query cache
func (c *Query) GetQCache(k string) (*map[string]*interface{}, bool) {
	keys, ok := c.cache[c.index][k]
	if !ok {
		return nil, false
	}
	return &keys, ok
}

// GetQCachePath for returning a path for a given key from the Query cache
func (c *Query) GetQCachePath(path string) (*interface{}, bool) {
	slicedPath := strings.Split(path, ".")
	key := slicedPath[len(slicedPath)-1]

	_, ok := c.cache[c.index][key]
	if !ok {
		return nil, false
	}

	obj, ok := c.cache[c.index][key][path]
	return obj, ok
}

// KeyExists check if a path exists for a given key
func (c *Query) KeyExists(path string) bool {
	slicedPath := strings.Split(path, ".")
	key := slicedPath[len(slicedPath)-1]

	p, ok := c.cache[c.index][key]
	if !ok {
		return false
	}

	_, ok = p[path]
	return ok
}

// FirstKey get the highest order path for a given key
func (c *Query) FirstKey(k string) (string, *interface{}, bool) {
	first, ok := c.first[c.index][k]
	if !ok {
		return "", nil, false
	}
	return first, c.cache[c.index][k][first], true
}

// QCacheKeys get all keys from Query cache
func (c *Query) QCacheKeys(k string) []string {
	var keys []string
	for i := range c.cache[c.index][k] {
		keys = append(keys, i)
	}
	return keys
}
