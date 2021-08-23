package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
	"gopkg.in/yaml.v2"
)

// TestUpsert run unit tests on Upsert
func TestUpsert(t *testing.T) {
	t.Parallel()

	path := "local/db-upsert.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.Equal(t, err, nil)

	f, err := ioutil.ReadFile(path)
	assert.Equal(t, err, nil)

	yaml.Unmarshal(f, &state.Data)

	testUpsert := []struct {
		Key   string
		Value string
	}{
		{"key-1", "value-1"},
		{"key-2", "value-2"},
	}

	for _, testCase := range testUpsert {
		data, ok := state.Data.(map[interface{}]interface{})
		assert.Equal(t, ok, true)

		assert.Equal(
			t,
			data["test"].(map[interface{}]interface{})["path"].(map[interface{}]interface{})[testCase.Key],
			testCase.Value,
			fmt.Sprintf("Expected: %v", testCase.Value),
		)
	}
}

// TestGetSingle run unit tests on GetSingle key
func TestGetSingle(t *testing.T) {
	t.Parallel()

	path := "local/db-query.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.Equal(t, err, nil)

	val, err := state.GetFirst("key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-1")
}

// TestGetPath run unit tests on Get object from path
func TestGetPath(t *testing.T) {
	t.Parallel()

	path := "local/db-get-path.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.Equal(t, err, nil)

	val, err := state.GetPath("test.path.key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-1")

	err = state.Upsert(
		"some",
		[]map[string][]string{
			{
				"array": {
					"value-1",
					"value-2",
				},
			},
		},
	)

	assert.Equal(t, err, nil)

	val, err = state.GetPath("some.[0].array")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, []interface{}{"value-1", "value-2"})

	err = state.Upsert(
		"array",
		[]string{
			"value-1",
			"value-2",
		},
	)

	assert.Equal(t, err, nil)

	val, err = state.GetPath("array.[0]")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-1")

	val, err = state.GetPath("array.[1]")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-2")
}

// TestDelete run unit tests for deleting objects
// from a given path
func TestDelete(t *testing.T) {
	t.Parallel()

	path := "local/db-delete-key.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.Equal(t, err, nil)

	err = state.Delete("test.path.key-1")
	assert.Equal(t, err, nil)

	val, err := state.GetPath("test.path.key-1")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, val, nil)
}

// TestGet run unit tests for searching for keys
func TestGet(t *testing.T) {
	t.Parallel()

	path := "local/db-get-keys.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"test.path",
		map[string]map[string]string{
			"key-1": {
				"v": "1",
			},
			"key-2": {
				"key-1": "1",
			},
		},
	)

	assert.Equal(t, err, nil)

	keys, err := state.Get("key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(keys), 2)

	for _, j := range keys {
		_, err = state.GetPath(j)
		assert.Equal(t, err, nil)
	}

}
