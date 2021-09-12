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

	err = storage.Delete("test.path.key-1")
	assert.Equal(t, err, nil)

	val, err := storage.GetPath("test.path.key-1")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, val, nil)

	err = storage.Delete("test.path.key-12")
	assert.NotEqual(t, err, nil)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}
