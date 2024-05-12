//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/swizzleio/swiz/cmd/cmds"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGenerateSimple(t *testing.T) {
	cmd := []string{"swiz", "config", "generate"}
	expect := []ExpectResponse{
		{
			Output: "Scanning for AWS accounts...\n",
			Action: NoAction,
		},
		{
			Output:   "Provide the name of your AWS account \n",
			Response: "dev-sandbox",
		},
		{
			Output:   "Enter the AWS account id \n",
			Response: "012345678901\n",
		},
		{
			Output:   "What region do you want to use for this account \n",
			Response: "us-east-2",
		},
		{
			Output:   "What domain name do you want to use for this environment \n",
			Response: "example.com",
		},
		{
			Output:   "Specify a comma seperated list for any global parameters (i.e. LogLevel,VpcId) \n",
			Response: "LogLevel,Apm",
		},
		{
			Output:   "Provide a value for the global parameter LogLevel \n",
			Response: "Debug",
		},
		{
			Output:   "Provide a value for the global parameter Apm \n",
			Response: "Enabled",
		},
		{
			Output:   "Enter the filename of the IaC template for your stack (leave blank to exit): \n",
			Response: "",
		},
		{
			Output:   "Name the enclave that AWS account dev-sandbox will be part of (leave blank to ignore) \n",
			Response: "",
		},
	}
	RunTestWithMocks(t, cmd, expect, 2000, func(t *testing.T, mocks cmds.FixtureMocks, resp string) {
		assert.Equal(t, "Version is dev(n/a)\n", resp)
		fmt.Println("Captured Output:", resp)
	})
}
