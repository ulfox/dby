package tests

import (
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestGetFirst run unit tests on GetSingle key
func TestGetFirst(t *testing.T) {
	t.Parallel()

	storage, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.Equal(t, err, nil)

	val, err := storage.GetFirst("key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-1")

	val, err = storage.GetFirst("key-2")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-2")

	err = storage.Upsert(
		"path-1",
		map[string][]string{
			"key-3": {"value-3"},
			"key-4": {"value-4"},
		},
	)

	assert.Equal(t, err, nil)

	val, err = storage.GetFirst("key-3")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, []interface{}{"value-3"})

	val, err = storage.GetFirst("key-4")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, []interface{}{"value-4"})
	err = storage.Upsert(
		"key-3",
		map[string][]string{
			"key-5": {"value-5"},
			"key-6": {"value-6"},
		},
	)

	assert.Equal(t, err, nil)

	val, err = storage.GetFirst("key-3")
	assert.Equal(t, err, nil)
	assert.Equal(t, val.(map[interface{}]interface{})["key-5"].([]interface{})[0], "value-5")
	assert.Equal(t, val.(map[interface{}]interface{})["key-6"].([]interface{})[0], "value-6")

	err = storage.Upsert(
		"test",
		map[string]string{},
	)
	assert.Equal(t, err, nil)
	err = storage.Upsert(
		"key-3",
		map[string][]string{},
	)
	assert.Equal(t, err, nil)
	err = storage.Upsert(
		"path-1",
		map[string][]string{},
	)
	assert.Equal(t, err, nil)
	err = storage.Upsert(
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

	val, err = storage.GetFirst("key-30")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, 30)
}
