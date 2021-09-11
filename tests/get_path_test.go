package tests

import (
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestGetPath run unit tests on Get object from path
func TestGetPath(t *testing.T) {
	t.Parallel()

	state, err := db.NewStorageFactory()
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
}
