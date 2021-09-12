package db

import (
	"fmt"
	"strings"
)

const (
	ire = "index error"
)

// state struct used by dby storage
type state struct {
	data   []interface{}
	buffer []*interface{}
	lib    map[string]int
	ad     int
}

// newStateFactory for creating a new v3 State
func newStateFactory() *state {
	s := state{
		data:   make([]interface{}, 0),
		buffer: make([]*interface{}, 0),
		lib:    make(map[string]int),
	}
	return &s
}

// Clear for clearing the v3 state
func (c *state) Clear() {
	c.data, c.buffer, c.lib = nil, nil, nil

	c.data = make([]interface{}, 0)
	c.buffer = make([]*interface{}, 0)
	c.lib = make(map[string]int)
}

// SetAD for setting new Active Document index
func (c *state) SetAD(i int) error {
	if err := c.IndexInRange(i); err != nil {
		return wrapErr(err)
	}
	c.ad = i
	return nil
}

// GetAD returns the current active document index
func (c *state) GetAD() int {
	return c.ad
}

// PushData for appending data to the data array
func (c *state) PushData(d interface{}) {
	c.data = append(c.data, d)
}

// PushBuffer for appending data to the buffer array
func (c *state) PushBuffer(d interface{}) {
	c.buffer = append(c.buffer, &d)
}

// GetAllData returns the data array
func (c *state) GetAllData() []interface{} {
	return c.data
}

// GetAllBuffer returns the buffer array
func (c *state) GetAllBuffer() []*interface{} {
	return c.buffer
}

// GetData returns the data in the c.ad index from the data array
func (c *state) GetData() interface{} {
	data, _ := c.GetDataFromIndex(c.GetAD())
	return data
}

// GetDataFromIndex returns the i'th element from the data array
func (c *state) GetDataFromIndex(i int) (interface{}, error) {
	if err := c.IndexInRange(i); err != nil {
		return nil, wrapErr(err)
	}
	return c.data[i], nil
}

// SetData sets to input value the data in the c.ad index from the data array
func (c *state) SetData(v interface{}) error {
	return c.SetDataFromIndex(v, c.GetAD())
}

// SetDataFromIndex sets to input value the i'th element from the data array
func (c *state) SetDataFromIndex(v interface{}, i int) error {
	if err := c.IndexInRange(i); err != nil {
		return wrapErr(err)
	}
	c.data[i] = v
	return nil
}

// GetBufferFromIndex returns the i'th element from the buffer array
func (c *state) GetBufferFromIndex(i int) (*interface{}, error) {
	if len(c.buffer)-1 >= i {
		return c.buffer[i], nil
	}
	return nil, fmt.Errorf(ire)
}

// SetDataFromIndex sets to input value the i'th element from the data array
func (c *state) SetBufferFromIndex(v interface{}, i int) error {
	if len(c.buffer)-1 >= i {
		c.buffer[i] = &v
		return nil
	}
	return fmt.Errorf(ire)

}

// IndexInRange check if index is within data array range
func (c *state) IndexInRange(i int) error {
	if len(c.data)-1 >= i {
		return nil
	}
	return fmt.Errorf(ire)
}

// Lib returns the lib map
func (c *state) Lib() map[string]int {
	return c.lib
}

// AddDoc for adding a document to the lib map
func (c *state) AddDoc(k string, i int) error {
	if err := c.IndexInRange(i); err != nil {
		return wrapErr(err)
	}
	c.lib[k] = i
	return nil
}

// LibIndex returns the index for a given doc name
func (c *state) LibIndex(doc string) (int, bool) {
	i, exists := c.lib[strings.ToLower(doc)]
	return i, exists
}

// RemoveDocName removes a doc from the lib
func (c *state) RemoveDocName(i int) error {
	if err := c.IndexInRange(i); err != nil {
		return wrapErr(err)
	}
	for k, v := range c.lib {
		if v == i {
			delete(c.lib, k)
		}
	}

	return nil
}

// DeleteData for deleting the i'th element from the data array
func (c *state) DeleteData(i int) error {
	if err := c.IndexInRange(i); err != nil {
		return wrapErr(err)
	}

	if err := c.RemoveDocName(i); err != nil {
		return wrapErr(err)
	}

	c.data[i] = nil
	c.data = append(c.data[:i], c.data[i+1:]...)

	if c.ad == i {
		if c.ad > 0 {
			c.ad = c.ad - 1
		} else {
			c.ad = 0
		}
	}

	return nil
}

// // CopyBufferToData for copying buffer array over data array
// func (c *state) CopyBufferToData() {
// 	copy(c.data, c.buffer)
// }

// UnsetDataArray for deleting all data. This sets data = nil
func (c *state) UnsetDataArray() {
	c.data = nil
}

// DeleteAllData calls PurgeAllData first and then creates a new empty array
func (c *state) DeleteAllData() {
	c.UnsetDataArray()
	c.data = make([]interface{}, 0)
}

// UnsetBufferArray This sets buffer = nil
func (c *state) UnsetBufferArray() {
	c.buffer = nil
}

// DeleteBuffer deletes the data from the buffer array
func (c *state) DeleteBuffer() {
	c.UnsetBufferArray()
	c.buffer = make([]*interface{}, 0)
}

// ClearLib removes all keys from the lib map
func (c *state) ClearLib() {
	c.lib = nil
	c.lib = make(map[string]int)
}
