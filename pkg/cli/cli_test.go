package appcli

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocksurvey "github.com/swizzleio/swiz/mocks/pkg/cli"
)

const mockArgMissing = "(Missing)"

func getMocks(inputs []string) (mockWrite *bytes.Buffer, mockErr *bytes.Buffer,
	mockSurvey *mocksurvey.SurveyWrapper, l *SwizCli) {
	mockWrite = bytes.NewBufferString("")
	if inputs == nil {
		inputs = []string{}
	}
	mockRead := strings.NewReader(strings.Join(inputs, "\n") + "\n")
	mockErr = bytes.NewBufferString("")
	mockSurvey = &mocksurvey.SurveyWrapper{}

	l = &SwizCli{
		output: mockWrite,
		input:  mockRead,
		err:    mockErr,
		survey: mockSurvey,
	}
	return
}

func TestSwizCli_Info(t *testing.T) {
	writeStr, _, _, l := getMocks(nil)

	l.Info("All your base are belong to us %v %v", 42, "foo")

	assert.Equal(t, "All your base are belong to us 42 foo", writeStr.String())
}

func TestSwizCli_Infoln(t *testing.T) {
	writeStr, _, _, l := getMocks(nil)

	l.Infoln("All your base are belong to us")

	assert.Equal(t, "All your base are belong to us\n", writeStr.String())
}

func TestSwizCli_Ask(t *testing.T) {
	_, _, surv, l := getMocks(nil)
	question := &survey.Input{
		Message: "Heya",
	}

	surv.On("AskOne", question, mock.AnythingOfType("*string")).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*string)
		*arg = "Foobar"
	}).Return(nil)

	got, err := l.Ask("Heya", false)
	assert.Equal(t, "Foobar", got)
	assert.NoError(t, err)
}

func TestSwizCli_AskAutocomplete(t *testing.T) {
	tests := []struct {
		name        string
		argPrompt   string
		argRequired bool
		argComplete AutocompleteFunc
		want        string
		wantErr     bool
	}{
		{
			name:        "required no func",
			argPrompt:   "Hey",
			argRequired: true,
			argComplete: nil,
			want:        "Sup",
			wantErr:     false,
		},
		{
			name:        "not required",
			argPrompt:   "Hey",
			argRequired: false,
			argComplete: nil,
			want:        "Sup",
			wantErr:     false,
		},
		{
			name:        "required",
			argPrompt:   "Hey",
			argRequired: true,
			argComplete: func(toComplete string) []string { return []string{"hi", "there"} },
			want:        "Sup",
			wantErr:     false,
		},
		{
			name:        "error",
			argPrompt:   "Hey",
			argRequired: false,
			argComplete: nil,
			want:        "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, surv, l := getMocks(nil)
			question := &survey.Input{
				Message: tt.argPrompt,
				Suggest: tt.argComplete,
			}
			var retErr error
			if tt.wantErr {
				retErr = fmt.Errorf("Ahhhhhhh!")
			}
			var mockAskOne *mock.Call

			// This is needed because mocks blows up when matching functions
			questionMatch := mock.MatchedBy(func(arg *survey.Input) bool {
				msgMatch := arg.Message == tt.argPrompt
				suggestMatch := true
				if tt.argComplete != nil {
					suggestMatch = arg.Suggest != nil
				}

				return msgMatch && suggestMatch
			})
			if tt.argRequired {
				// Required arg needs an autocomplete function
				mockAskOne = surv.On("AskOne", questionMatch, mock.AnythingOfType("*string"), mock.AnythingOfType("survey.AskOpt"))
			} else {
				mockAskOne = surv.On("AskOne", question, mock.AnythingOfType("*string"))
			}
			mockAskOne.Run(func(args mock.Arguments) {
				arg := args.Get(1).(*string)
				*arg = tt.want
			}).Return(retErr)

			got, err := l.AskAutocomplete(tt.argPrompt, tt.argRequired, tt.argComplete)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSwizCli_AskConfirm(t *testing.T) {
	tests := []struct {
		name    string
		argMsg  string
		argOpt  []string
		want    bool
		wantErr bool
	}{
		{
			name:    "happy case",
			argMsg:  "Hey",
			want:    true,
			wantErr: false,
		},
		{
			name:    "error case",
			argMsg:  "Hey",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, surv, l := getMocks(nil)
			question := &survey.Confirm{
				Message: tt.argMsg,
			}
			var retErr error
			if tt.wantErr {
				retErr = fmt.Errorf("Ahhhhhhh!")
			}
			surv.On("AskOne", question, mock.AnythingOfType("*bool")).Run(func(args mock.Arguments) {
				arg := args.Get(1).(*bool)
				*arg = tt.want
			}).Return(retErr)

			got, err := l.AskConfirm(tt.argMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSwizCli_AskOptions(t *testing.T) {
	tests := []struct {
		name    string
		argMsg  string
		argOpt  []string
		want    string
		wantErr bool
	}{
		{
			name:    "happy case",
			argMsg:  "Hey",
			argOpt:  []string{"one", "two"},
			want:    "Sup",
			wantErr: false,
		},
		{
			name:    "error case",
			argMsg:  "Hey",
			argOpt:  []string{"one", "two"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, surv, l := getMocks(nil)
			question := &survey.Select{
				Message: tt.argMsg,
				Options: tt.argOpt,
			}
			var retErr error
			if tt.wantErr {
				retErr = fmt.Errorf("Ahhhhhhh!")
			}
			surv.On("AskOne", question, mock.AnythingOfType("*string")).Run(func(args mock.Arguments) {
				arg := args.Get(1).(*string)
				*arg = tt.want
			}).Return(retErr)

			got, err := l.AskOptions(tt.argMsg, tt.argOpt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSwizCli_AskMany(t *testing.T) {
	tests := []struct {
		name    string
		argOpt  []AskManyOpts
		want    map[string]string
		wantErr bool
	}{
		{
			name: "happy case",
			argOpt: []AskManyOpts{
				{
					Key:           "Foo",
					Message:       "What is foo?",
					Required:      false,
					TransformMode: TransformModeNone,
				},
				{
					Key:           "Bar",
					Message:       "What is love?",
					Required:      false,
					TransformMode: TransformModeNone,
				},
			},
			want: map[string]string{
				"Foo": "What is foo?",
				"Bar": "What is love?",
			},
			wantErr: false,
		},
		{
			name: "mixed case",
			argOpt: []AskManyOpts{
				{
					Key:           "Foo",
					Message:       "What is foo?",
					Required:      false,
					TransformMode: TransformModeNone,
				},
				{
					Key:           "Bar",
					Message:       "What is love?",
					Required:      true,
					TransformMode: TransformModeNone,
				},
				{
					Key:           "Stuff",
					Message:       "Lots of stuff!",
					Required:      true,
					TransformMode: TransformModeTrimSpace,
				},
				{
					Key:           "Blah",
					Message:       "Waaah",
					Required:      false,
					TransformMode: TransformModeCamelCase,
				},
			},
			want: map[string]string{
				"Foo":   "What is foo?",
				"Bar":   "What is love?",
				"Stuff": "Lots of stuff!",
				"Blah":  "Waaah",
			},
			wantErr: false,
		},
			{
				name: "fail case",
				argOpt: []AskManyOpts{
					{
						Key:           "Foo",
						Message:       "What is foo?",
						Required:      false,
						TransformMode: TransformModeNone,
					},
				},
				want:    nil,
				wantErr: true,
			},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, surv, l := getMocks(nil)

			var retErr error
			if tt.wantErr {
				retErr = fmt.Errorf("Ahhhhhhh!")
			}
			surv.On("Ask", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				arg := args.Get(1).(*map[string]interface{})
				for _, p := range tt.argOpt {
					_, ok := (*arg)[p.Key]
					assert.True(t, ok)
					(*arg)[p.Key] = p.Message
				}
			}).Return(retErr)

			got, err := l.AskMany(tt.argOpt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_convertToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "normal",
			want: "normal",
		},
		{
			name: "All Your Base Are Belong To Us",
			want: "allYourBaseAreBelongToUs",
		},
		{
			name: "Foo-Bar Stuff",
			want: "foo-barStuff",
		},
		{
			name: "La la la",
			want: "laLaLa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToCamelCase(tt.name)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSwizCli_getOutput(t *testing.T) {
	w, _, _, l := getMocks(nil)
	write := l.getOutput()
	assert.Equal(t, w, write)
	l.output = nil
	write = l.getOutput()
	assert.Equal(t, os.Stdout, write)
}
