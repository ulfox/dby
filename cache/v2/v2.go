package v1

// Cache hosts results from SQL methods
type Cache struct {
	v1 []string
}

// NewCacheFactory for creating a new Cache
func NewCacheFactory() Cache {
	cache := Cache{
		v1: make([]string, 0),
	}
	return cache
}

// Clear for clearing the cache
func (c *Cache) Clear() {
	c.v1 = make([]string, 0)
}

// AddKey for appending a new Key to cache
func (c *Cache) AddKey(k string) {
	c.v1 = append(c.v1, k)
}

// GetKeys returns Cache.KeysFound
func (c *Cache) GetKeys() []string {
	return c.v1
}
