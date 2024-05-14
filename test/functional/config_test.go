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
	expect := []*ExpectResponse{
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
	RunTestWithMocks(t, cmd, expect, 2, func(t *testing.T, mocks cmds.FixtureMocks, resp string) {
		// Get the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		fh := fileutil.NewFileUrlHelper(mocks.Fs)
		buf, err := fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
		assert.NoError(t, err)

		y, err := NewYamlTestUtil(string(buf))
		y.AssertEqual(t, 1, "version")
		y.AssertEqual(t, "NameMe", "default_enclave")
		y.AssertEqual(t, "{{env_name:32}}-{{stack_name:32}}", "naming_scheme")
		y.AssertEqual(t, "NameMe", "enclave_def[name=NameMe].name")
		y.AssertEqual(t, "dev-sandbox", "enclave_def[name=NameMe].default_provider")
		y.AssertEqual(t, "Cloudformation", "enclave_def[name=NameMe].default_iac")
		y.AssertEqual(t, "example.com", "enclave_def[name=NameMe].domain_name")
		y.AssertEqual(t, "dev-sandbox", "enclave_def[name=NameMe].providers[name=dev-sandbox].name")
		y.AssertEqual(t, "AWS", "enclave_def[name=NameMe].providers[name=dev-sandbox].provider_id")
		y.AssertEqual(t, "012345678901", "enclave_def[name=NameMe].providers[name=dev-sandbox].account_id")
		y.AssertEqual(t, "us-east-2", "enclave_def[name=NameMe].providers[name=dev-sandbox].region")
		y.AssertEqual(t, true, "enclave_def[name=NameMe].env_behavior.deploy_all_stacks")
		y.AssertEqual(t, "Enabled", "enclave_def[name=NameMe].params.Apm")
		y.AssertEqual(t, "Debug", "enclave_def[name=NameMe].params.LogLevel")
		y.AssertEqual(t, []interface{}{}, "stack_cfg")
	})
}

func TestConfigGenerateFull(t *testing.T) {
	cmd := []string{"swiz", "config", "generate"}
	expect := []*ExpectResponse{
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
			Response: "foobar.yaml",
			Match:    ExactOnce,
		},
		{
			Output:   "Enter the name of the stack: \n",
			Response: "Awesomesauce",
			Match:    ExactOnce,
		},
		{
			Output:   "Enter the filename of the IaC template for your stack (leave blank to exit): \n",
			Response: "anotherone.yaml",
			Match:    ExactOnce,
		},
		{
			Output:   "Enter the name of the stack: \n",
			Response: "Notsoawesome",
			Match:    ExactOnce,
		},
		{
			Output:   "Enter the filename of the IaC template for your stack (leave blank to exit): \n",
			Response: "",
			Match:    ExactOnce,
		},
		{
			Output:   "Name the enclave that AWS account dev-sandbox will be part of (leave blank to ignore) \n",
			Response: "partycity",
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

		y, err := NewYamlTestUtil(string(buf))
		y.AssertEqual(t, 1, "version")
		y.AssertEqual(t, "partycity", "default_enclave")
		y.AssertEqual(t, "{{env_name:32}}-{{stack_name:32}}", "naming_scheme")
		y.AssertEqual(t, "partycity", "enclave_def[name=partycity].name")
		y.AssertEqual(t, "dev-sandbox", "enclave_def[name=partycity].default_provider")
		y.AssertEqual(t, "Cloudformation", "enclave_def[name=partycity].default_iac")
		y.AssertEqual(t, "example.com", "enclave_def[name=partycity].domain_name")
		y.AssertEqual(t, "dev-sandbox", "enclave_def[name=partycity].providers[name=dev-sandbox].name")
		y.AssertEqual(t, "AWS", "enclave_def[name=partycity].providers[name=dev-sandbox].provider_id")
		y.AssertEqual(t, "012345678901", "enclave_def[name=partycity].providers[name=dev-sandbox].account_id")
		y.AssertEqual(t, "us-east-2", "enclave_def[name=partycity].providers[name=dev-sandbox].region")
		y.AssertEqual(t, true, "enclave_def[name=partycity].env_behavior.deploy_all_stacks")
		y.AssertEqual(t, "Enabled", "enclave_def[name=partycity].params.Apm")
		y.AssertEqual(t, "Debug", "enclave_def[name=partycity].params.LogLevel")
		y.AssertEqual(t, "Awesomesauce", "stack_cfg[name=Awesomesauce].name")
		y.AssertEqual(t, "file://./out/Awesomesauce-cfg.yaml", "stack_cfg[name=Awesomesauce].config_file")
		y.AssertEqual(t, 1, "stack_cfg[name=Awesomesauce].order")
		y.AssertEqual(t, "Notsoawesome", "stack_cfg[name=Notsoawesome].name")
		y.AssertEqual(t, "file://./out/Notsoawesome-cfg.yaml", "stack_cfg[name=Notsoawesome].config_file")
		y.AssertEqual(t, 1, "stack_cfg[name=Notsoawesome].order")
	})
}

func TestConfigGenerateMinimal(t *testing.T) {
	cmd := []string{"swiz", "config", "generate"}
	expect := []*ExpectResponse{
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
			Response: "",
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
	RunTestWithMocks(t, cmd, expect, 2, func(t *testing.T, mocks cmds.FixtureMocks, resp string) {
		// Get the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		fh := fileutil.NewFileUrlHelper(mocks.Fs)
		buf, err := fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
		assert.NoError(t, err)

		y, err := NewYamlTestUtil(string(buf))
		y.AssertEqual(t, 1, "version")
		y.AssertEqual(t, "NameMe", "default_enclave")
		y.AssertEqual(t, "{{env_name:32}}-{{stack_name:32}}", "naming_scheme")
		y.AssertEqual(t, "NameMe", "enclave_def[name=NameMe].name")
		y.AssertEqual(t, "dev-sandbox", "enclave_def[name=NameMe].default_provider")
		y.AssertEqual(t, "Cloudformation", "enclave_def[name=NameMe].default_iac")
		y.AssertEqual(t, "example.com", "enclave_def[name=NameMe].domain_name")
		y.AssertEqual(t, "dev-sandbox", "enclave_def[name=NameMe].providers[name=dev-sandbox].name")
		y.AssertEqual(t, "AWS", "enclave_def[name=NameMe].providers[name=dev-sandbox].provider_id")
		y.AssertEqual(t, "012345678901", "enclave_def[name=NameMe].providers[name=dev-sandbox].account_id")
		y.AssertEqual(t, "us-east-2", "enclave_def[name=NameMe].providers[name=dev-sandbox].region")
		y.AssertEqual(t, true, "enclave_def[name=NameMe].env_behavior.deploy_all_stacks")
		y.AssertEqual(t, map[string]interface{}{}, "enclave_def[name=NameMe].params")
		y.AssertEqual(t, []interface{}{}, "stack_cfg")
	})
}

func TestConfigGenerateCustomLocation(t *testing.T) {
	cmd := []string{"swiz", "config", "generate", "--output", "file://blah/foobar.yaml"}
	expect := []*ExpectResponse{
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
			Response: "",
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
	RunTestWithMocks(t, cmd, expect, 2, func(t *testing.T, mocks cmds.FixtureMocks, resp string) {
		fh := fileutil.NewFileUrlHelper(mocks.Fs)
		buf, err := fh.OpenUrl("file://blah/foobar.yaml")
		assert.NoError(t, err)

		y, err := NewYamlTestUtil(string(buf))
		y.AssertEqual(t, "NameMe", "default_enclave")
	})
}
