package tests

import (
	"strings"
	"testing"

	"github.com/likexian/gokit/assert"
	"github.com/ulfox/dby/db"
)

// TestMultiDoc run tests for docs
func TestMultiDoc(t *testing.T) {
	t.Parallel()

	storage, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 1)

	err = storage.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 2)

	assert.Equal(t, storage.GetAD(), 1)
	err = storage.Switch(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, storage.GetAD(), 0)

	err = storage.DeleteDoc(1)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 1)

	err = storage.DeleteDoc(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 0)

	err = storage.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 1)

	err = storage.DeleteAll(true).
		ImportDocs("../docs/examples/manifests/deployment.yaml", true)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(storage.GetAllData()), 8)

	err = storage.SetNames("kind", "metadata.name")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.ListDocs()), 8)

	for _, j := range []string{
		"spec.selector.matchLabels.version",
		"metadata.labels.version",
		"spec.selector.version",
		"spec.template.selector.matchLabels.version",
		"spec.template.metadata.labels.version",
	} {
		err = storage.UpdateGlobal(
			j,
			"v0.3.0",
		)
		assert.Equal(t, err, nil)
	}

	for _, j := range storage.ListDocs() {
		if strings.HasPrefix(j, "horizontalpodautoscaler/") {
			continue
		}
		err = storage.SwitchDoc(j)
		assert.Equal(t, err, nil)

		val, err := storage.GetPath("metadata.labels.version")
		assert.Equal(t, err, nil)
		assert.Equal(t, val, "v0.3.0")
	}

	err = storage.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 9)
	err = storage.SetNames("kind", "metadata.name")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.ListDocs()), 8)
	assert.Equal(t, len(storage.GetAllData()), 9)
	for i, j := range storage.Lib() {
		err = storage.Switch(j)
		assert.Equal(t, err, nil)

		kind, err := storage.GetPath("kind")
		assert.Equal(t, err, nil)

		name, err := storage.GetPath("metadata.name")
		assert.Equal(t, err, nil)

		sKind, ok := kind.(string)
		assert.Equal(t, ok, true)

		sName, ok := name.(string)
		assert.Equal(t, ok, true)
		assert.Equal(t, i, strings.ToLower(sKind)+"/"+strings.ToLower(sName))
	}

	c0 := 0
	err = storage.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 10)
	err = storage.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(storage.GetAllData()), 11)
	err = storage.SetNames("kind", "metadata.name")
	assert.Equal(t, err, nil)

	for i := range storage.GetAllData() {
		err = storage.Switch(i)
		assert.Equal(t, err, nil)

		kind, err := storage.GetPath("kind")
		if err != nil {
			c0++
			continue
		}

		name, err := storage.GetPath("metadata.name")
		assert.Equal(t, err, nil)

		sKind, ok := kind.(string)
		assert.Equal(t, ok, true)

		sName, ok := name.(string)
		assert.Equal(t, ok, true)

		doc, ok := storage.LibIndex(strings.ToLower(sKind) + "/" + strings.ToLower(sName))
		assert.Equal(t, ok, true)
		assert.Equal(t, doc, i)
	}
	assert.Equal(t, c0, 3)

	data := storage.GetPathGlobal("metadata.name")
	assert.Equal(t, len(data), 8)

	dataMap := map[int][]string{
		0: {"metadata.name", "spec.metrics.[0].resource.name", "spec.scaleTargetRef.name"},
		1: {"metadata.name", "spec.template.spec.containers.[0].name"},
		2: {"spec.ports.[0].name", "metadata.name"},
		3: {"metadata.name"},
		4: {"metadata.name", "spec.metrics.[0].resource.name", "spec.scaleTargetRef.name"},
		5: {"metadata.name", "spec.template.spec.containers.[0].name"},
		6: {"metadata.name", "spec.ports.[0].name"},
		7: {"metadata.name"},
	}

	for i, j := range storage.FindKeysGlobal("name") {
		err = storage.Switch(i)
		assert.Equal(t, err, nil)

		assert.Equal(t, len(j), len(dataMap[i]))

		for _, m := range dataMap[i] {
			assert.Equal(t, checkArrays(m, j), true)
		}

	}
}

func checkArrays(v string, l []string) bool {
	for _, j := range l {
		if j == v {
			return true
		}
	}
	return false
}
