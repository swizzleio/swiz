package appcli

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	mocksurvey "github.com/swizzleio/swiz/mocks/pkg/cli"
	"io"
	"reflect"
	"strings"
	"testing"
)

func getMocks(inputs []string) (mockWrite *bytes.Buffer, mockErr *bytes.Buffer, l *SwizCli) {
	mockWrite = bytes.NewBufferString("")
	if inputs == nil {
		inputs = []string{}
	}
	mockRead := strings.NewReader(strings.Join(inputs, "\n") + "\n")
	mockErr = bytes.NewBufferString("")

	l = &SwizCli{
		output: mockWrite,
		input:  mockRead,
		err:    mockErr,
		survey: &mocksurvey.SurveyWrapper{},
	}
	return
}

func TestSwizCli_Info(t *testing.T) {
	writeStr, _, l := getMocks(nil)

	l.Info("All your base are belong to us %v %v", 42, "foo")

	assert.Equal(t, "All your base are belong to us 42 foo", writeStr.String())
}

func TestSwizCli_Infoln(t *testing.T) {
	writeStr, _, l := getMocks(nil)

	l.Infoln("All your base are belong to us")

	assert.Equal(t, "All your base are belong to us\n", writeStr.String())
}

func TestSwizCli_Ask(t *testing.T) {
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	type args struct {
		prompt   string
		required bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
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
			got, err := l.Ask(tt.args.prompt, tt.args.required)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ask() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwizCli_AskAutocomplete(t *testing.T) {
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	type args struct {
		prompt   string
		required bool
		complete func(toComplete string) []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
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
			got, err := l.AskAutocomplete(tt.args.prompt, tt.args.required, tt.args.complete)
			if (err != nil) != tt.wantErr {
				t.Errorf("AskAutocomplete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AskAutocomplete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwizCli_AskConfirm(t *testing.T) {
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	type args struct {
		prompt string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
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
			got, err := l.AskConfirm(tt.args.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("AskConfirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AskConfirm() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwizCli_AskOptions(t *testing.T) {
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	type args struct {
		prompt  string
		options []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
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
			got, err := l.AskOptions(tt.args.prompt, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("AskOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AskOptions() got = %v, want %v", got, tt.want)
			}
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
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   io.Writer
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
			if got := l.getOutput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwizCli_getInput(t *testing.T) {
	type fields struct {
		output io.Writer
		input  io.Reader
		err    io.Writer
		reader *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   *bufio.Reader
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
			if got := l.getInput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
