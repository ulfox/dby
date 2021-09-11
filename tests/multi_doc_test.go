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

	state, err := db.NewStorageFactory()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 1)

	err = state.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 2)

	assert.Equal(t, state.AD, 1)
	err = state.Switch(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, state.AD, 0)

	err = state.DeleteDoc(1)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 1)

	err = state.DeleteDoc(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 0)

	err = state.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 1)

	err = state.DeleteAll(true).
		ImportDocs("../docs/examples/manifests/deployment.yaml", true)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(state.Data), 8)

	err = state.SetNames("kind", "metadata.name")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.ListDocs()), 8)

	for _, j := range []string{
		"spec.selector.matchLabels.version",
		"metadata.labels.version",
		"spec.selector.version",
		"spec.template.selector.matchLabels.version",
		"spec.template.metadata.labels.version",
	} {
		err = state.UpdateGlobal(
			j,
			"v0.3.0",
		)
		assert.Equal(t, err, nil)
	}

	for _, j := range state.ListDocs() {
		if strings.HasPrefix(j, "horizontalpodautoscaler/") {
			continue
		}
		err = state.SwitchDoc(j)
		assert.Equal(t, err, nil)

		val, err := state.GetPath("metadata.labels.version")
		assert.Equal(t, err, nil)
		assert.Equal(t, val, "v0.3.0")
	}

	err = state.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 9)
	err = state.SetNames("kind", "metadata.name")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.ListDocs()), 8)
	assert.Equal(t, len(state.Data), 9)
	for i, j := range state.Lib {
		err = state.Switch(j)
		assert.Equal(t, err, nil)

		kind, err := state.GetPath("kind")
		assert.Equal(t, err, nil)

		name, err := state.GetPath("metadata.name")
		assert.Equal(t, err, nil)

		sKind, ok := kind.(string)
		assert.Equal(t, ok, true)

		sName, ok := name.(string)
		assert.Equal(t, ok, true)
		assert.Equal(t, i, strings.ToLower(sKind)+"/"+strings.ToLower(sName))
	}

	c0 := 0
	err = state.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 10)
	err = state.AddDoc()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(state.Data), 11)
	err = state.SetNames("kind", "metadata.name")
	assert.Equal(t, err, nil)

	for i := range state.Data {
		err = state.Switch(i)
		assert.Equal(t, err, nil)

		kind, err := state.GetPath("kind")
		if err != nil {
			c0++
			continue
		}

		name, err := state.GetPath("metadata.name")
		assert.Equal(t, err, nil)

		sKind, ok := kind.(string)
		assert.Equal(t, ok, true)

		sName, ok := name.(string)
		assert.Equal(t, ok, true)

		assert.Equal(t, state.Lib[strings.ToLower(sKind)+"/"+strings.ToLower(sName)], i)
	}
	assert.Equal(t, c0, 3)

	data := state.GetPathGlobal("metadata.name")
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

	for i, j := range state.FindKeysGlobal("name") {
		err = state.Switch(i)
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
