package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
	"gopkg.in/yaml.v2"
)

// TestUpsert run unit tests on Upsert
func TestUpsert(t *testing.T) {
	t.Parallel()

	path := ".test/db-upsert.yaml"
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

	err = state.Upsert(
		"path-1.sub-path-1",
		map[string][]string{
			"sub-path-2": {"value-1", "value-2"},
			"sub-path-3": {"value-3", "value-4"},
		},
	)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"path-2",
		[]map[string][]string{
			{
				"sub-path-1": {"value-1", "value-2"},
			},
			{
				"sub-path-2": {"value-3", "value-4"},
			},
		},
	)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"path-3",
		[]map[string]string{
			{
				"sub-path-1": "value-1",
			},
			{
				"sub-path-2": "value-2",
			},
		},
	)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"path-4",
		map[string]int{
			"sub-path-1": 0,
			"sub-path-2": 1,
		},
	)
	assert.Equal(t, err, nil)

	f, err := ioutil.ReadFile(path)
	assert.Equal(t, err, nil)

	yaml.Unmarshal(f, &state.Data[state.AD])

	testUpsert := []struct {
		Key   string
		Value string
	}{
		{"key-1", "value-1"},
		{"key-2", "value-2"},
	}

	data, ok := state.Data[state.AD].(map[interface{}]interface{})
	assert.Equal(t, ok, true)
	for _, testCase := range testUpsert {

		assert.Equal(
			t,
			data["test"].(map[interface{}]interface{})["path"].(map[interface{}]interface{})[testCase.Key],
			testCase.Value,
			fmt.Sprintf("Expected: %v", testCase.Value),
		)
	}

	assert.Equal(
		t,
		data["path-2"].([]interface{})[0].(map[interface{}]interface{})["sub-path-1"].([]interface{})[0],
		"value-1",
	)
	assert.Equal(
		t,
		data["path-2"].([]interface{})[0].(map[interface{}]interface{})["sub-path-1"].([]interface{})[1],
		"value-2",
	)
	assert.Equal(
		t,
		data["path-2"].([]interface{})[1].(map[interface{}]interface{})["sub-path-2"].([]interface{})[0],
		"value-3",
	)
	assert.Equal(
		t,
		data["path-2"].([]interface{})[1].(map[interface{}]interface{})["sub-path-2"].([]interface{})[1],
		"value-4",
	)

	assert.Equal(
		t,
		data["path-3"].([]interface{})[0].(map[interface{}]interface{})["sub-path-1"],
		"value-1",
	)
	assert.Equal(
		t,
		data["path-3"].([]interface{})[1].(map[interface{}]interface{})["sub-path-2"],
		"value-2",
	)

	assert.Equal(t, data["path-4"].(map[interface{}]interface{})["sub-path-1"], 0)
	assert.Equal(t, data["path-4"].(map[interface{}]interface{})["sub-path-2"], 1)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}

// TestGetFirst run unit tests on GetSingle key
func TestGetFirst(t *testing.T) {
	t.Parallel()

	path := ".test/db-query.yaml"
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

	val, err = state.GetFirst("key-2")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-2")

	err = state.Upsert(
		"path-1",
		map[string][]string{
			"key-3": {"value-3"},
			"key-4": {"value-4"},
		},
	)

	assert.Equal(t, err, nil)

	val, err = state.GetFirst("key-3")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, []interface{}{"value-3"})

	val, err = state.GetFirst("key-4")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, []interface{}{"value-4"})
	err = state.Upsert(
		"key-3",
		map[string][]string{
			"key-5": {"value-5"},
			"key-6": {"value-6"},
		},
	)

	assert.Equal(t, err, nil)

	val, err = state.GetFirst("key-3")
	assert.Equal(t, err, nil)
	assert.Equal(t, val.(map[interface{}]interface{})["key-5"].([]interface{})[0], "value-5")
	assert.Equal(t, val.(map[interface{}]interface{})["key-6"].([]interface{})[0], "value-6")

	err = state.Upsert(
		"test",
		map[string]string{},
	)
	assert.Equal(t, err, nil)
	err = state.Upsert(
		"key-3",
		map[string][]string{},
	)
	assert.Equal(t, err, nil)
	err = state.Upsert(
		"path-1",
		map[string][]string{},
	)
	assert.Equal(t, err, nil)
	err = state.Upsert(
		"to.array-10",
		map[string][]map[string]int{
			"key-10": {
				{"key-20": 20},
				{"key-30": 30},
				{"key-40": 40},
			},
		},
	)
	assert.Equal(t, err, nil)

	val, err = state.GetFirst("key-30")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, 30)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}

// TestGetPath run unit tests on Get object from path
func TestGetPath(t *testing.T) {
	t.Parallel()

	path := ".test/db-get-path.yaml"
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
					"value-3",
					"value-4",
				},
			},
		},
	)

	assert.Equal(t, err, nil)

	val, err = state.GetPath("some.[0].array")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, []interface{}{"value-3", "value-4"})

	err = state.Upsert(
		"array",
		[]string{
			"value-5",
			"value-6",
		},
	)

	assert.Equal(t, err, nil)

	val, err = state.GetPath("array.[0]")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-5")

	val, err = state.GetPath("array.[1]")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-6")

	val, err = state.GetPath("array.[2]")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, val, nil)

	val, err = state.GetPath("some.[2].array")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, val, nil)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}

// TestDelete run unit tests for deleting objects
// from a given path
func TestDelete(t *testing.T) {
	t.Parallel()

	path := ".test/db-delete-key.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
			"key-3": "value-3",
			"key-4": "value-4",
		},
	)

	assert.Equal(t, err, nil)

	err = state.Delete("test.path.key-1")
	assert.Equal(t, err, nil)

	val, err := state.GetPath("test.path.key-1")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, val, nil)

	err = state.Delete("test.path.key-12")
	assert.NotEqual(t, err, nil)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}

// TestGet run unit tests for searching for keys
func TestGet(t *testing.T) {
	t.Parallel()

	path := ".test/db-get-keys.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		"path-1",
		[]map[string][]map[string]string{
			{
				"subpath-01": {
					{"k01": "v01"},
				},
				"subpath-02": {
					{"k02": "v02"},
				},
			},
			{
				"subpath-11": {
					{"k11": "v11"},
				},
				"subpath-12": {
					{"k12": "v12"},
				},
			},
		},
	)
	assert.Equal(t, err, nil)

	assertData := db.NewConvertFactory()
	assertData.Input(state.Data[state.AD])

	assertData.
		Key("path-1").
		Index(0).
		Key("subpath-01").
		Index(0)

	assert.Equal(t, assertData.Error, nil)
	assert.Equal(t, assertData.GetMap()["k01"], "v01")

	obj, err := state.GetPath("path-1.[1].subpath-11.[0]")
	assert.Equal(t, err, nil)
	assertData.Input(obj)
	assert.Equal(t, assertData.GetMap()["k11"], "v11")

	assert.Equal(t, assertData.Error, nil)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}

// TestGeneric run generic tests for all scenarios
func TestGeneric(t *testing.T) {
	t.Parallel()

	path := ".test/db-with-empties.yaml"
	state, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = state.Upsert(
		".someKey",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.NotEqual(t, err, nil)

	err = state.Upsert(
		".",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.NotEqual(t, err, nil)

	err = state.Upsert(
		"k01",
		nil,
	)

	assert.Equal(t, err, nil)

	err = state.Upsert(
		"k",
		[]map[string][]map[string]string{
			{
				"0": {
					{"1": "v03"},
				},
				"2": {
					{"03": "v05"},
				},
			},
			{
				"3": {
					{"2": "v11"},
				},
				"4": {
					{"3": "v12"},
				},
			},
		},
	)
	assert.Equal(t, err, nil)
	val, err := state.GetFirst("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "v03")
	assertData := db.NewConvertFactory()

	val, err = state.GetFirst("03")
	assert.Equal(t, err, nil)

	assertData.Input(val)
	assert.Equal(t, assertData.Cache.V1, "v05")

	keys, err := state.Get("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(keys), 1)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}
