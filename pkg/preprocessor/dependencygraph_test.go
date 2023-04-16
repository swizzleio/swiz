package preprocessor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDepTreeTestCircDependencies(t *testing.T) {
	depTree := NewDepTree()

	err := depTree.AddNodes(map[string][]string{
		"a": {},
		"b": {"a", "e"},
		"c": {"d"},
		"d": {"a", "c"},
		"e": {"a"},
	})
	assert.NoError(t, err)

	err = depTree.Build()
	assert.NoError(t, err)

	circDepends := depTree.GetCircularDependencies()
	assert.ElementsMatch(t, []string{"c", "d"}, circDepends)
}
