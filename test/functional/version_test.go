//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/swizzleio/swiz/test/functional/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	cmd := []string{"swiz", "version"}
	_, resp := util.RunCommandWithMocks(t, cmd, nil, util.DefaultCmdTimeoutSec, false, nil)

	assert.Equal(t, "Version is dev(n/a)\n", resp)
	fmt.Println("Captured Output:", resp)
}
