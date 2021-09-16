package tests

import (
	"os"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestDelete run unit tests for deleting objects
// from a given path
func TestDelete(t *testing.T) {
	t.Parallel()

	path := ".test/db-delete-key.yaml"
	storage, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
			"key-3": "value-3",
			"key-4": "value-4",
		},
	)

	assert.Equal(t, err, nil)
	val, err := storage.GetFirst("key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "value-1")

	err = storage.Delete("test.path.key-1")
	assert.Equal(t, err, nil)

	val, err = storage.GetPath("test.path.key-1")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, val, nil)

	err = storage.Delete("test.path.key-12")
	assert.NotEqual(t, err, nil)

	err = storage.Upsert(
		"key-33",
		map[string][]string{
			"key-56": {"value-5", "value-6"},
			"key-78": {"value-7", "value-8"},
		},
	)

	assert.Equal(t, err, nil)
	err = storage.Delete("key-33.key-56")
	assert.Equal(t, err, nil)
	val, err = storage.GetPath("key-33")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(val.(map[interface{}]interface{})), 1)

	err = storage.Delete("key-33.key-78.[0]")
	assert.Equal(t, err, nil)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}
