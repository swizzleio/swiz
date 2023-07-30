package appcli

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"strings"
)

type SwizCli struct {
	output io.Writer
	input  io.Reader
	err    io.Writer
	reader *bufio.Reader
	survey SurveyWrapper
}

const (
	TransformModeNone = iota
	TransformModeTrimSpace
	TransformModeCamelCase
)

type TransformModeType int

type AskManyOpts struct {
	Key           string
	Message       string
	Required      bool
	TransformMode TransformModeType
}

type SwizClier interface {
	Info(format string, i ...interface{})
	Infoln(i ...interface{})
	Ask(prompt string, required bool) (string, error)
	AskAutocomplete(prompt string, required bool, complete func(toComplete string) []string) (string, error)
	AskConfirm(prompt string) (bool, error)
	AskOptions(prompt string, options []string) (string, error)
	AskMany(prompts []AskManyOpts) (map[string]string, error)
}

// NewCli creates a new cli
func NewCli(o io.Writer, i io.Reader) SwizClier {
	return &SwizCli{
		output: o,
		input:  i,
		err:    o,
		survey: &SurveyWrap{},
	}
}

// Info outputs an informational message
func (l *SwizCli) Info(format string, i ...interface{}) {
	_, _ = fmt.Fprintf(l.getOutput(), format, i...)
}

// Infoln outputs an informational message with a newline
func (l *SwizCli) Infoln(i ...interface{}) {
	l.Info(fmt.Sprintln(i...))
}

// Ask asks a question
func (l *SwizCli) Ask(prompt string, required bool) (string, error) {
	return l.AskAutocomplete(prompt, required, nil)
}

// AskAutocomplete asks a question with autocomplete
func (l *SwizCli) AskAutocomplete(prompt string, required bool, complete func(toComplete string) []string) (string, error) {
	var resp string
	question := &survey.Input{
		Message: prompt,
		Suggest: complete,
	}

	var err error
	if required {
		err = l.survey.AskOne(question, &resp, survey.WithValidator(survey.Required))
	} else {
		err = l.survey.AskOne(question, &resp)
	}
	if err != nil {
		return "", err
	}

	return resp, nil
}

// AskConfirm asks the user to confirm with Y/n
func (l *SwizCli) AskConfirm(prompt string) (bool, error) {
	var resp bool
	question := &survey.Confirm{
		Message: prompt,
	}

	err := l.survey.AskOne(question, &resp)
	if err != nil {
		return false, err
	}

	return resp, nil
}

// AskOptions asks with a list of options
func (l *SwizCli) AskOptions(prompt string, options []string) (string, error) {
	var resp string
	question := &survey.Select{
		Message: prompt,
		Options: options,
	}

	err := l.survey.AskOne(question, &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// AskMany asks many questions
func (l *SwizCli) AskMany(prompts []AskManyOpts) (map[string]string, error) {
	// TODO: Update the other Ask* calls to use this or an internal method
	answers := map[string]interface{}{}
	qs := []*survey.Question{}
	strAnswers := map[string]string{}

	for _, prompt := range prompts {
		q := &survey.Question{
			Name:   prompt.Key,
			Prompt: &survey.Input{Message: prompt.Message},
		}

		if prompt.Required {
			q.Validate = survey.Required
		}

		switch prompt.TransformMode {
		case TransformModeTrimSpace:
			q.Transform = survey.TransformString(strings.TrimSpace)
		case TransformModeCamelCase:
			q.Transform = survey.TransformString(convertToCamelCase)
		}

		answers[prompt.Key] = ""
		qs = append(qs, q)
	}

	err := l.survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	for k, v := range answers {
		strAnswers[k] = fmt.Sprintf("%v", v)
	}

	return strAnswers, nil
}

// convertToCamelCase converts the string to camel case
func convertToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	words := strings.Split(s, " ")

	// Convert the first word to lowercase
	words[0] = strings.ToLower(words[0])

	// Convert subsequent words to title case
	for i := 1; i < len(words); i++ {
		words[i] = cases.Title(language.Und, cases.NoLower).String(words[i])
	}

	// Join the words together
	return strings.Join(words, "")
}

// getOutput fetches the output writer or if not set, uses stdout
func (l *SwizCli) getOutput() io.Writer {
	if l.output != nil {
		return l.output
	}

	return os.Stdout
}

// getInput fetches the input reader or if not set, uses stdin
func (l *SwizCli) getInput() *bufio.Reader {
	if l.reader != nil {
		return l.reader
	}

	input := l.input
	if input == nil {
		input = os.Stdin
	}

	l.reader = bufio.NewReader(input)

	return l.reader
}
