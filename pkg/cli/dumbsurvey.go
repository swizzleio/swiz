package appcli

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io"
	"io/fs"
	"os"
	"strings"
)

type DumbSurvey struct {
	input  *bufio.Reader
	output io.Writer
}

func NewDumbSurvey(input io.Reader, output io.Writer) *DumbSurvey {
	if input == nil {
		input = os.Stdin
	}
	if output == nil {
		output = os.Stdout
	}

	return &DumbSurvey{
		input:  bufio.NewReader(input),
		output: output,
	}
}

// AskOne is a wrapper for the survey.AskOne() func
func (d *DumbSurvey) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	prompt := getStringFromPrompt(p)
	_, oErr := fmt.Fprintf(d.output, "%v \n", prompt)
	if oErr != nil {
		return oErr
	}

	resp, err := d.input.ReadString('\n')
	if err != nil {
		if _, ok := err.(*fs.PathError); !ok {
			// Only return if it's not a path error
			return err
		}
	}

	resp = strings.TrimSpace(resp)

	switch v := response.(type) {
	case *string:
		*v = resp
	case *bool:
		*v = strings.ToLower(resp) == "t"
	// Add more case statements if you need to handle more types
	default:
		return fmt.Errorf("unsupported response type: %T", response)
	}

	return nil
}

// Ask is a wrapper for the survey.Ask() func
func (d DumbSurvey) Ask(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
	answers := map[string]string{}
	for _, q := range qs {
		resp := ""
		err := d.AskOne(q.Prompt, &resp, opts...)
		if err != nil {
			return err
		}

		answers[q.Name] = resp
	}

	pResponse := response.(*map[string]interface{})
	for k, v := range answers {
		var iface interface{} = v
		(*pResponse)[k] = iface
	}

	return nil
}

func getStringFromPrompt(prompt survey.Prompt) string {
	resp := ""
	switch p := prompt.(type) {
	case *survey.Input:
		resp = p.Message
	case *survey.Confirm:
		resp = p.Message
	case *survey.Select:
		resp = p.Message
	}

	return resp
}
