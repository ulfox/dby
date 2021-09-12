package v2

// Query hosts results from SQL methods
type Query struct {
	keys []string
}

// NewQueryFactory for creating a new Query
func NewQueryFactory() Query {
	Query := Query{
		keys: make([]string, 0),
	}
	return Query
}

// Clear for clearing the Query
func (c *Query) Clear() {
	c.keys = make([]string, 0)
}

// AddKey for appending a new Key to Query
func (c *Query) AddKey(k string) {
	c.keys = append(c.keys, k)
}

// GetKeys returns Query.KeysFound
func (c *Query) GetKeys() []string {
	return c.keys
}
