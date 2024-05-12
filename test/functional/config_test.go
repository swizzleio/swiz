//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/swizzleio/swiz/cmd/cmds"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"os"
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
			Response: "012345678901",
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
		// Get the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		fh := fileutil.NewFileUrlHelper(mocks.Fs)
		buf, err := fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
		assert.NoError(t, err)

		yaml := `version: 1
default_enclave: NameMe
naming_scheme: '{{env_name:32}}-{{stack_name:32}}'
enclave_def:
    - name: NameMe
      default_provider: dev-sandbox
      default_iac: Cloudformation
      providers:
        - name: dev-sandbox
          provider_id: AWS
          account_id: "012345678901"
          region: us-east-2
      env_behavior:
        deploy_all_stacks: true
      domain_name: example.com
      params:
        Apm: Enabled
        LogLevel: Debug
stack_cfg: []
`
		assert.Equal(t, yaml, string(buf))
	})
}
