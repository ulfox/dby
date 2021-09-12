package tests

import (
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestGet run unit tests for searching for keys
func TestGet(t *testing.T) {
	t.Parallel()

	storage, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)

	err = storage.Upsert(
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
	assertData.Input(storage.State.GetData())

	assertData.
		Key("path-1").
		Index(0).
		Key("subpath-01").
		Index(0)

	assert.Equal(t, assertData.GetError(), nil)
	m, err := assertData.GetMap()
	assert.Equal(t, assertData.GetError(), nil)
	assert.Equal(t, err, nil)

	assert.Equal(t, m["k01"], "v01")

	obj, err := storage.GetPath("path-1.[1].subpath-11.[0]")
	assert.Equal(t, err, nil)
	assertData.Input(obj)
	m, err = assertData.GetMap()
	assert.Equal(t, assertData.GetError(), nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, m["k11"], "v11")

	assert.Equal(t, assertData.GetError(), nil)
}
