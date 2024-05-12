package appcli

import (
	"bytes"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
)

func TestAskOne(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      string
		expectedError bool
	}{
		{
			name:          "basic input",
			input:         "example input\n",
			expected:      "example input",
			expectedError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := bytes.NewBufferString(tc.input)
			ds := NewDumbSurvey(input, nil)
			var response string
			err := ds.AskOne(&survey.Input{Message: "Please input:"}, &response)
			assert.Equal(t, tc.expected, response)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetStringFromPrompt(t *testing.T) {
	tests := []struct {
		name     string
		prompt   survey.Prompt
		expected string
	}{
		{
			name:     "Input prompt",
			prompt:   &survey.Input{Message: "What is your name?"},
			expected: "What is your name?",
		},
		{
			name:     "Confirm prompt",
			prompt:   &survey.Confirm{Message: "Do you agree?"},
			expected: "Do you agree?",
		},
		{
			name:     "Select prompt",
			prompt:   &survey.Select{Message: "Choose your option:"},
			expected: "Choose your option:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringFromPrompt(tt.prompt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAsk(t *testing.T) {
	cases := []struct {
		name          string
		questions     []*survey.Question
		input         string
		expected      map[string]interface{}
		expectedError bool
	}{
		{
			name: "single question",
			questions: []*survey.Question{
				{
					Name:   "name",
					Prompt: &survey.Input{Message: "What is your name?"},
				},
			},
			input:         "John\n",
			expected:      map[string]interface{}{"name": "John"},
			expectedError: false,
		},
		{
			name: "multiple questions",
			questions: []*survey.Question{
				{
					Name:   "name",
					Prompt: &survey.Input{Message: "What is your name?"},
				},
				{
					Name:   "age",
					Prompt: &survey.Input{Message: "What is your age?"},
				},
			},
			input:         "John\n25\n",
			expected:      map[string]interface{}{"name": "John", "age": "25"},
			expectedError: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := bytes.NewBufferString(tc.input)
			ds := NewDumbSurvey(input, nil)
			response := make(map[string]interface{})
			err := ds.Ask(tc.questions, &response)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, response)
			}
		})
	}
}
