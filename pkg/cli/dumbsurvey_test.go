package appcli

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
)

func TestAddResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    []SurveyResponse
		expected int
	}{
		{
			name: "Add single response",
			input: []SurveyResponse{
				{Prompt: "Test prompt 1", ResponseList: []string{"Response 1"}, NextResponse: 0},
			},
			expected: 1,
		},
		{
			name: "Add multiple responses",
			input: []SurveyResponse{
				{Prompt: "Test prompt 1", ResponseList: []string{"Response 1"}, NextResponse: 0},
				{Prompt: "Test prompt 2", ResponseList: []string{"Response 2"}, NextResponse: 0},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDumbSurvey()
			ds.AddResponses(tt.input)
			assert.Equal(t, tt.expected, len(ds.responseList))
		})
	}
}

func TestGetPrompts(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *DumbSurvey
		prompt   string
		expected []string
		err      string
	}{
		{
			name: "Prompt found",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				ds.AddResponse(SurveyResponse{Prompt: "Test", ActualPrompts: []string{"Test 1", "Test 2"}})
				return ds
			},
			prompt:   "Test",
			expected: []string{"Test 1", "Test 2"},
			err:      "",
		},
		{
			name: "Prompt not found",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				return ds
			},
			prompt:   "Non-existent",
			expected: []string{},
			err:      "could not find response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := tt.setup()
			actual, err := ds.GetPrompts(tt.prompt)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestAskOne(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *DumbSurvey
		prompt  survey.Prompt
		wantErr bool
	}{
		{
			name: "Successful ask",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				ds.AddResponse(SurveyResponse{Prompt: "Enter your name:", ResponseList: []string{"John Doe"}, NextResponse: 0})
				return ds
			},
			prompt:  &survey.Input{Message: "Enter your name:"},
			wantErr: false,
		},
		{
			name: "Prompt not found",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				return ds
			},
			prompt:  &survey.Input{Message: "Enter your age:"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := tt.setup()
			var resp string
			err := ds.AskOne(tt.prompt, &resp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "John Doe", resp)
			}
		})
	}
}

func TestAsk(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *DumbSurvey
		questions []*survey.Question
		expected  map[string]interface{}
		wantErr   bool
	}{
		{
			name: "Successful multiple asks",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				ds.AddResponses([]SurveyResponse{
					{Prompt: "Name", ResponseList: []string{"Alice"}, NextResponse: 0},
					{Prompt: "Age", ResponseList: []string{"30"}, NextResponse: 0},
				})
				return ds
			},
			questions: []*survey.Question{
				{Name: "name", Prompt: &survey.Input{Message: "Name"}},
				{Name: "age", Prompt: &survey.Input{Message: "Age"}},
			},
			expected: map[string]interface{}{"name": "Alice", "age": "30"},
			wantErr:  false,
		},
		{
			name: "One question prompt not found",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				ds.AddResponses([]SurveyResponse{
					{Prompt: "Name", ResponseList: []string{"Bob"}, NextResponse: 0},
				})
				return ds
			},
			questions: []*survey.Question{
				{Name: "name", Prompt: &survey.Input{Message: "Name"}},
				{Name: "age", Prompt: &survey.Input{Message: "Age"}},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := tt.setup()
			response := make(map[string]interface{})
			err := ds.Ask(tt.questions, &response)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, response)
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

func TestGetResponseObj(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *DumbSurvey
		prompt   string
		expected *SurveyResponse
		wantErr  bool
	}{
		{
			name: "Find exact match",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				ds.AddResponse(SurveyResponse{Prompt: "Location?", ResponseList: []string{"New York"}, NextResponse: 0})
				return ds
			},
			prompt:   "Location?",
			expected: &SurveyResponse{Prompt: "Location?", ResponseList: []string{"New York"}, NextResponse: 0},
			wantErr:  false,
		},
		{
			name: "Prompt not found",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				return ds
			},
			prompt:   "Location?",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := tt.setup()
			result, err := ds.getResponseObj(tt.prompt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetResponse(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *DumbSurvey
		prompt   string
		expected string
		wantErr  bool
	}{
		{
			name: "Get valid response",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				ds.AddResponse(SurveyResponse{Prompt: "Favorite color?", ResponseList: []string{"Blue", "Red"}, NextResponse: 0})
				return ds
			},
			prompt:   "Favorite color?",
			expected: "Blue",
			wantErr:  false,
		},
		{
			name: "Prompt not found",
			setup: func() *DumbSurvey {
				ds := NewDumbSurvey()
				return ds
			},
			prompt:   "Favorite color?",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := tt.setup()
			result, err := ds.getResponse(tt.prompt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
