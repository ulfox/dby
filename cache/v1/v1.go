package v1

// Cache for easily sharing state between map operations and methods
// V{1,2} are common interface{} placeholders, C{1,2,3} are counters used
// for condition, while []Keys is used by path discovery methods
// to keep track and derive the right path.
type Cache struct {
	v1, v2     interface{}
	v3         []interface{}
	b          []byte
	e          error
	c1, c2, c3 int
	keys       []string
}

// NewCacheFactory for creating a new Cache
func NewCacheFactory() Cache {
	cache := Cache{
		v3:   make([]interface{}, 0),
		b:    make([]byte, 0),
		keys: make([]string, 0),
	}
	return cache
}

// E can be used to set or get e. If no argument is passed
// the method returns e value. If an argument is passed, the method
// sets e's value to that of the arguments
func (c *Cache) E(e ...error) error {
	if len(e) > 0 {
		c.e = e[0]
	}
	return c.e
}

// Clear for clearing the cache and setting all
// content to nil
func (c *Cache) Clear() {
	c.v1, c.v2 = nil, nil
	c.v3 = make([]interface{}, 0)
	c.c1, c.c2, c.c3 = 0, 0, 0
	c.b = make([]byte, 0)
	c.e = nil
	c.keys = make([]string, 0)
}

// DropLastKey removes the last key added from the list
func (c *Cache) DropLastKey() {
	if len(c.keys) > 0 {
		c.keys = c.keys[:len(c.keys)-1]
	}
}

// DropKeys removes all keys from the list
func (c *Cache) DropKeys() {
	c.keys = make([]string, 0)
}

// AddKey for appending a new Key to cache
func (c *Cache) AddKey(k string) {
	c.keys = append(c.keys, k)
}

// GetKeys returns Cache.Keys
func (c *Cache) GetKeys() []string {
	return c.keys
}

// DropV3 removes all keys from the list
func (c *Cache) DropV3() {
	c.v3 = make([]interface{}, 0)
}

// V3 returns Cache.V3
func (c *Cache) V3(v ...interface{}) []interface{} {
	if len(v) > 0 {
		c.v3 = append(c.v3, v[0])
	}
	return c.v3
}

// C1 can be used to set or get c1. If no argument is passed
// the method returns c1 value. If an argument is passed, the method
// sets c1's value to that of the arguments
func (c *Cache) C1(i ...int) int {
	if len(i) > 0 {
		c.c1 = i[0]
	}
	return c.c1
}

// C2 can be used to set or get c2. If no argument is passed
// the method returns c2 value. If an argument is passed, the method
// sets c2's value to that of the arguments
func (c *Cache) C2(i ...int) int {
	if len(i) > 0 {
		c.c2 = i[0]
	}
	return c.c2
}

// C3 can be used to set or get c3. If no argument is passed
// the method returns c3 value. If an argument is passed, the method
// sets c3's value to that of the arguments
func (c *Cache) C3(i ...int) int {
	if len(i) > 0 {
		c.c3 = i[0]
	}
	return c.c3
}

// V1 can be used to set or get v1. If no argument is passed
// the method returns v1 value. If an argument is passed, the method
// sets v1's value to that of the arguments
func (c *Cache) V1(v ...interface{}) interface{} {
	if len(v) > 0 {
		c.v1 = v[0]
	}
	return c.v1
}

// V2 can be used to set or get v2. If no argument is passed
// the method returns v2 value. If an argument is passed, the method
// sets v2's value to that of the arguments
func (c *Cache) V2(v ...interface{}) interface{} {
	if len(v) > 0 {
		c.v2 = v[0]
	}

	return c.v2
}

// B can be used to set or get b. If no argument is passed
// the method returns b value. If an argument is passed, the method
// sets b's value to that of the arguments
func (c *Cache) B(b ...[]byte) []byte {
	if len(b) > 0 {
		c.b = b[0]
	}

	return c.b
}

// BE expects 2 inputs, a byte array and an error.
// If byte array is not nil, b will be set to the
// value of the new byte array. Error is returned
// with the new/old byte array at the end always
func (c *Cache) BE(b []byte, e error) error {
	if b != nil {
		c.b = b
	}

	return e
}

// V1E expects 2 inputs, a byte array and an error.
// If byte array is not nil, b will be set to the
// value of the new byte array. Error is returned
// with the new/old byte array at the end always
func (c *Cache) V1E(v interface{}, e error) error {
	if v != nil {
		c.v1 = v
	}

	return e
}
