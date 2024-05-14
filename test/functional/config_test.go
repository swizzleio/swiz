//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/swizzleio/swiz/cmd/cmds"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// parseYAML takes a YAML string and parses it into a nested map
func parseYAML2(yamlData string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// fetchValue traverses a nested map according to a dot-separated path
func fetchValue2(data map[string]interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		if strings.Contains(part, "[") {
			// Split the part into the key and the condition
			idxStart := strings.Index(part, "[")
			idxEnd := strings.Index(part, "]")
			if idxStart == -1 || idxEnd == -1 {
				return nil, fmt.Errorf("syntax error in array access")
			}

			key := part[:idxStart]
			condition := part[idxStart+1 : idxEnd]

			// Navigate into the map
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("attempted to navigate a non-map element")
			}

			// Get the slice from the map
			slice, ok := currentMap[key].([]interface{})
			if !ok {
				return nil, fmt.Errorf("element at key '%s' is not a slice", key)
			}

			// Parse the condition, expected format: name=Value
			eqIndex := strings.Index(condition, "=")
			if eqIndex == -1 {
				// Handle index-based access
				index, err := strconv.Atoi(condition)
				if err != nil {
					return nil, fmt.Errorf("invalid index '%s': %v", condition, err)
				}
				if index < 0 || index >= len(slice) {
					return nil, fmt.Errorf("index out of range: %d", index)
				}
				current = slice[index]
			} else {
				// Handle conditional access
				attr := condition[:eqIndex]
				value := condition[eqIndex+1:]
				found := false

				for _, item := range slice {
					asMap, ok := item.(map[string]interface{})
					if !ok {
						continue
					}
					if asMap[attr] == value {
						current = item
						found = true
						break
					}
				}

				if !found {
					return nil, fmt.Errorf("no element matches the condition: %s", condition)
				}
			}
		} else {
			currentMap, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("attempted to navigate a non-map element at '%s'", part)
			}
			current = currentMap[part]
		}
	}

	return current, nil
}

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
