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

	statefulEmpty, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(statefulEmpty.Data), 1)

	assertData := db.NewConvertFactory()
	assertData.Input(statefulEmpty.Data[statefulEmpty.AD])

	emptyMap, err := assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(emptyMap), 0)

	err = statefulEmpty.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)
	assert.Equal(t, err, nil)

	statefulEmpty = nil

	statefulData, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(statefulData.Data), 1)
	assertData.Input(statefulData.Data[statefulData.AD])

	assertData.Key("test")
	assert.Equal(t, assertData.GetError(), nil)
	assertData.Key("path")
	assert.Equal(t, assertData.GetError(), nil)

	dataMap, err := assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dataMap), 2)

	statefulData.InMem(true)
	statefulData.DeleteAll(true)
	assert.Equal(t, len(statefulData.Data), 0)

	statefulData, dataMap = nil, nil

	statefulDataDue, err := db.NewStorageFactory(path)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(statefulDataDue.Data), 1)
	assertData.Input(statefulDataDue.Data[statefulDataDue.AD])

	assertData.Key("test")
	assert.Equal(t, assertData.GetError(), nil)
	assertData.Key("path")
	assert.Equal(t, assertData.GetError(), nil)

	dataMap, err = assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(dataMap), 2)

	statelessData, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(statelessData.Data), 1)

	assertData.Input(statelessData.Data[statelessData.AD])

	emptyMap, err = assertData.GetMap()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(emptyMap), 0)

	err = statelessData.Upsert(
		"test.path",
		map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	)
	assert.Equal(t, err, nil)
	assertData.Input(statelessData.Data[statelessData.AD])

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
