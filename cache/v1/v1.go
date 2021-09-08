package v1

// Cache for easily sharing state between map operations and methods
// V{1,2} are common interface{} placeholders, C{1,2,3} are counters used
// for condition, while []Keys is used by path discovery methods
// to keep track and derive the right path.
type Cache struct {
	V1, V2     interface{}
	V3         []interface{}
	B          []byte
	E          error
	C1, C2, C3 int
	Keys       []string
}

// NewCacheFactory for creating a new Cache
func NewCacheFactory() Cache {
	cache := Cache{
		V3:   make([]interface{}, 0),
		B:    make([]byte, 0),
		Keys: make([]string, 0),
	}
	return cache
}

// GetError returns cache error field
func (c *Cache) GetError() error {
	return c.E
}

// Clear for clearing the cache and setting all
// content to nil
func (c *Cache) Clear() {
	c.V1, c.V2 = nil, nil
	c.V3 = make([]interface{}, 0)
	c.C1, c.C2, c.C3 = 0, 0, 0
	c.B = make([]byte, 0)
	c.E = nil
	c.Keys = make([]string, 0)
}

// DropLastKey removes the last key added from the list
func (c *Cache) DropLastKey() {
	if len(c.Keys) > 0 {
		c.Keys = c.Keys[:len(c.Keys)-1]
	}
}

// DropKeys removes all keys from the list
func (c *Cache) DropKeys() {
	c.Keys = make([]string, 0)
}

// AddKey for appending a new Key to cache
func (c *Cache) AddKey(k string) {
	c.Keys = append(c.Keys, k)
}

// GetKeys returns Cache.Keys
func (c *Cache) GetKeys() []string {
	return c.Keys
}