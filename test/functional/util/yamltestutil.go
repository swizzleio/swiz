package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"testing"
)

// YamlTestUtil provides a set of test functions that parses yaml and allows for string based indexing
type YamlTestUtil struct {
	y map[string]interface{}
}

// NewYamlTestUtil creates a new YAML test utility helper
func NewYamlTestUtil(yamlData string) (*YamlTestUtil, error) {
	y, err := parseYAML(yamlData)
	if err != nil {
		return nil, err
	}

	return &YamlTestUtil{
		y: y,
	}, nil
}

// AssertEqual validates the value in the YAML is equal
func (y YamlTestUtil) AssertEqual(t *testing.T, expected interface{}, path string) {
	val, err := fetchValue(y.y, path)
	assert.NoError(t, err)
	assert.Equal(t, expected, val)

}

// parseYAML takes a YAML string and parses it into a nested map
func parseYAML(yamlData string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// fetchValue traverses a nested map according to a dot-separated path
func fetchValue(data map[string]interface{}, path string) (interface{}, error) {
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
