package tests

import (
	"os"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestGeneric run generic tests for all scenarios
func TestGeneric(t *testing.T) {
	t.Parallel()

	path := ".test/db-generic.yaml"
	storage, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)
	storage.InMem(true)

	err = storage.Upsert(
		".someKey",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.NotEqual(t, err, nil)

	err = storage.Upsert(
		".",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)

	assert.NotEqual(t, err, nil)

	err = storage.Upsert(
		"k01",
		nil,
	)

	assert.Equal(t, err, nil)

	err = storage.Upsert(
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

	val, err := storage.GetFirst("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, val, "v03")
	assertData := db.NewConvertFactory()

	val, err = storage.GetFirst("03")
	assert.Equal(t, err, nil)

	assertData.Input(val)
	s, err := assertData.GetString()
	assert.Equal(t, assertData.GetError(), nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, s, "v05")

	keys, err := storage.Get("1")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(keys), 1)

	err = storage.Upsert(
		"i",
		[]int{
			1,
			2,
			3,
		},
	)
	assert.Equal(t, err, nil)

	val, err = storage.GetFirst("i")
	assert.Equal(t, err, nil)
	assertData.Input(val)

	i, err := assertData.GetArray()
	assert.Equal(t, assertData.GetError(), nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(i), 3)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}
