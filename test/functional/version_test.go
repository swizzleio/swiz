//go:build functionaltest

package functional

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	cmd := []string{"swiz", "version"}
	_, resp := RunCommandWithMocks(t, cmd, nil, DefaultCmdTimeoutSec)

	assert.Equal(t, "Version is dev(n/a)\n", resp)
	fmt.Println("Captured Output:", resp)
}
