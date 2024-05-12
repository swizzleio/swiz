//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/swizzleio/swiz/cmd/cmds"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	cmd := []string{"swiz", "version"}
	RunTestWithMocks(t, cmd, nil, 2, func(t *testing.T, mocks cmds.FixtureMocks, resp string) {
		assert.Equal(t, "Version is dev(n/a)\n", resp)
		fmt.Println("Captured Output:", resp)
	})
}
