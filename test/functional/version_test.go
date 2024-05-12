//go:build functionaltest

package functional

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	cmd := []string{"swiz", "version"}
	expect := map[string]string{}
	RunTestWithMocks(t, cmd, expect, func(t *testing.T, appFs afero.Fs, resp string) {
		assert.Equal(t, "Version is dev(n/a)", resp)
		fmt.Println("Captured Output:", resp)
	})
}
