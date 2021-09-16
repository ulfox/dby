package tests

import (
	"os"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestCacheV2 run unit tests for cachev2
func TestCacheV2(t *testing.T) {
	t.Parallel()

	path := ".test/db-cache-v2.yaml"
	storage, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"test",
		1,
	)
	assert.Equal(t, err, nil)

	v, ok := storage.SQL.GetQCache("test")
	assert.Equal(t, v, (*map[string]*interface{})(nil))
	assert.Equal(t, ok, false)

	val, err := storage.FindKeys("test")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(val), 1)
	assert.Equal(t, val[0], "test")

	v, ok = storage.SQL.GetQCache("test")
	assert.Equal(t, ok, true)
	assert.Equal(t, (*(*v)["test"]), 1)

	err = storage.Delete("test")
	assert.Equal(t, err, nil)

	v, ok = storage.SQL.GetQCache("test")
	assert.Equal(t, v, (*map[string]*interface{})(nil))
	assert.Equal(t, ok, false)

	val, err = storage.FindKeys("test")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(val), 0)

	err = storage.Upsert(
		"key-1",
		"value-1",
	)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"path-1",
		map[string]string{
			"key-1": "value-1",
		},
	)
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"path-2",
		map[string]string{
			"key-1": "value-1",
		},
	)
	assert.Equal(t, err, nil)

	_, ok = storage.SQL.GetQCache("key-1")
	assert.Equal(t, ok, false)

	val, err = storage.FindKeys("key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(val), 3)

	v, ok = storage.SQL.GetQCache("key-1")
	assert.Equal(t, ok, true)
	assert.Equal(t, (*(*v)["key-1"]), "value-1")
	assert.Equal(t, len((*v)), 3)

	err = storage.Delete("path-1.key-1")
	assert.Equal(t, err, nil)

	v, ok = storage.SQL.GetQCache("key-1")
	assert.Equal(t, ok, true)
	assert.Equal(t, (*(*v)["path-2.key-1"]), "value-1")
	assert.Equal(t, len((*v)), 2)

	err = storage.Delete("key-1")
	assert.Equal(t, err, nil)

	_, ok = storage.SQL.GetQCache("key-1")
	assert.Equal(t, ok, false)

	val, err = storage.FindKeys("key-1")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(val), 1)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}

// TestCacheV2Mem run unit tests for cachev2
func TestCacheV2Mem(t *testing.T) {
	t.Parallel()

	storage, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)

	err = storage.Upsert(
		"test",
		1,
	)
	assert.Equal(t, err, nil)

	v, ok := storage.SQL.GetQCache("test")
	assert.Equal(t, v, (*map[string]*interface{})(nil))
	assert.Equal(t, ok, false)

	val, err := storage.FindKeys("test")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(val), 1)
	assert.Equal(t, val[0], "test")

	v, ok = storage.SQL.GetQCache("test")
	assert.Equal(t, ok, true)
	assert.Equal(t, (*(*v)["test"]), 1)

	err = storage.Delete("test")
	assert.Equal(t, err, nil)

	v, ok = storage.SQL.GetQCache("test")
	assert.Equal(t, v, (*map[string]*interface{})(nil))
	assert.Equal(t, ok, false)

}
