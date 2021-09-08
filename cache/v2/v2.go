package v1

// Cache hosts results from SQL methods
type Cache struct {
	KeysFound []string
	Results   interface{}
}

// NewCacheFactory for creating a new Cache
func NewCacheFactory() Cache {
	cache := Cache{
		KeysFound: make([]string, 0),
	}
	return cache
}

// Clear for clearing the cache
func (c *Cache) Clear() {
	c.Results = nil
	c.KeysFound = make([]string, 0)
}

// AddKey for appending a new Key to cache
func (c *Cache) AddKey(k string) {
	c.KeysFound = append(c.KeysFound, k)
}

// GetKeys returns Cache.KeysFound
func (c *Cache) GetKeys() []string {
	return c.KeysFound
}
