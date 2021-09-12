package tests

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
	storage, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"path-1.sub-path-1",
		map[string][]string{
			"sub-path-2": {"value-1", "value-2"},
			"sub-path-3": {"value-3", "value-4"},
		},
	)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
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

	err = storage.Upsert(
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

	err = storage.Upsert(
		"path-4",
		map[string]int{
			"sub-path-1": 0,
			"sub-path-2": 1,
		},
	)
	assert.Equal(t, err, nil)

	f, err := ioutil.ReadFile(path)
	assert.Equal(t, err, nil)

	v := storage.State.GetData()
	yaml.Unmarshal(f, &v)

	testUpsert := []struct {
		Key   string
		Value string
	}{
		{"key-1", "value-1"},
		{"key-2", "value-2"},
	}

	data, ok := storage.State.GetData().(map[interface{}]interface{})
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
