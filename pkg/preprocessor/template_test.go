package preprocessor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanTemplateParam(t *testing.T) {
	param := CleanTemplateParam("{{paramName}}")
	assert.Equal(t, "paramName", param, "Expected CleanTemplateParam to return paramName")

	param = CleanTemplateParam("paramName")
	assert.Equal(t, "paramName", param, "Expected CleanTemplateParam to return paramName")
}

func TestIsTemplateReplaceParam(t *testing.T) {
	isParam := IsTemplateReplaceParam("{{paramName}}")
	assert.True(t, isParam, "Expected IsTemplateReplaceParam to return true")

	isParam = IsTemplateReplaceParam("paramName")
	assert.False(t, isParam, "Expected IsTemplateReplaceParam to return false")
}

func TestParseTemplateTokens(t *testing.T) {
	testCases := []struct {
		desc        string
		template    string
		replaceIdx  map[string]string
		expectedRes string
	}{
		{
			desc:     "Basic replacement",
			template: "Hello, {{name:5}}, welcome to {{location:3}}",
			replaceIdx: map[string]string{
				"name":     "Maggie",
				"location": "San Francisco",
			},
			expectedRes: "Hello, Maggi, welcome to San",
		},
		{
			desc:     "No truncation necessary",
			template: "My name is {{name:10}}",
			replaceIdx: map[string]string{
				"name": "John",
			},
			expectedRes: "My name is John",
		},
		{
			desc:     "Variable not found",
			template: "Hello, {{name:5}}",
			replaceIdx: map[string]string{
				"age": "25",
			},
			expectedRes: "Hello, ",
		},
		{
			desc:     "Multiple replacements with the same variable",
			template: "Hi, {{name:2}}. Bye, {{name:2}}.",
			replaceIdx: map[string]string{
				"name": "Alice",
			},
			expectedRes: "Hi, Al. Bye, Al.",
		},
		// Add more test cases as necessary...
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := ParseTemplateTokens(tC.template, tC.replaceIdx)
			assert.Equal(t, tC.expectedRes, result, "Expected ParseTemplateTokens to return the replaced string")
		})
	}
}
