package appcli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocksurvey "github.com/swizzleio/swiz/mocks/pkg/cli"
)

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
	// TODO
}

func TestSwizCli_AskAutocomplete(t *testing.T) {
	// TODO
}

func TestSwizCli_AskConfirm(t *testing.T) {
	// TODO
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

			got, err := l.AskOptions("Hey", []string{"one", "two"})
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
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	type args struct {
		prompts []AskManyOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &SwizCli{
				output: tt.fields.output,
				input:  tt.fields.input,
				err:    tt.fields.err,
				reader: tt.fields.reader,
			}
			got, err := l.AskMany(tt.args.prompts)
			if (err != nil) != tt.wantErr {
				t.Errorf("AskMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AskMany() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToCamelCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToCamelCase(tt.args.s); got != tt.want {
				t.Errorf("convertToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwizCli_getOutput(t *testing.T) {
	// TODO
}

func TestSwizCli_getInput(t *testing.T) {
	// TODO
}
