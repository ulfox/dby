package tests

import (
	"os"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// fileExists for checking if a file exists
func fileExists(filepath string) bool {
	f, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	} else if f.IsDir() {
		return true
	}

	return true
}

// TestStorage run unit tests on storage
func TestStorage(t *testing.T) {
	t.Parallel()

	path := ".test/db-storage.yaml"
	assert.Equal(t, fileExists(path), false)

	empty, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(empty.State.GetAllData()), 1)

	assertData := db.NewConvertFactory()
	assertData.Input(empty.State.GetData())

	emptyMap, err := assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(emptyMap), 0)

	err = empty.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)
	assert.Equal(t, err, nil)

	empty = nil

	data, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(data.State.GetAllData()), 1)
	assertData.Input(data.State.GetData())

	assertData.Key("test")
	assert.Equal(t, assertData.GetError(), nil)
	assertData.Key("path")
	assert.Equal(t, assertData.GetError(), nil)

	dataMap, err := assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dataMap), 2)

	data.InMem(true)
	data.DeleteAll(true)
	assert.Equal(t, len(data.State.GetAllData()), 0)

	data, dataMap = nil, nil

	dataDue, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dataDue.State.GetAllData()), 1)
	assertData.Input(dataDue.State.GetData())

	assertData.Key("test")
	assert.Equal(t, assertData.GetError(), nil)
	assertData.Key("path")
	assert.Equal(t, assertData.GetError(), nil)

	dataMap, err = assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dataMap), 2)

	memData, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(memData.State.GetAllData()), 1)

	assertData.Input(memData.State.GetData())

	emptyMap, err = assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(emptyMap), 0)

	err = memData.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)
	assert.Equal(t, err, nil)
	assertData.Input(memData.State.GetData())

	assertData.Key("test")
	assert.Equal(t, assertData.GetError(), nil)
	assertData.Key("path")
	assert.Equal(t, assertData.GetError(), nil)

	emptyMap, err = assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(emptyMap), 2)

	err = os.Remove(path)
	assert.Equal(t, err, nil)
}
