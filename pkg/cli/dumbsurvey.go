package appcli

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
)

// This provides a more testable version of survey, instead of smart prompting, it just uses a scanner

type DumbSurvey struct {
	scanner *bufio.Reader
}

// AskOne is a wrapper for the survey.AskOne() func
func (d *DumbSurvey) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	fmt.Printf("%v ", getStringFromPrompt(p))

	if d.scanner == nil {
		//strReader := strings.NewReader("dev-sandbox\n012345678901\nus-east-2\nexample.com\nLogLevel,Apm\nDebug\nEnabled\n\n\n")
		d.scanner = bufio.NewReader(os.Stdin)
		//d.scanner = bufio.NewReader(strReader)
	}

	line, err := d.scanner.ReadString('\n')
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	//fmt.Println(line)

	pResponse := response.(*string)
	*pResponse = line

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
