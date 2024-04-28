package appcli

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
)

// This provides a more testable version of survey, instead of smart prompting, it just uses a scanner

type DumbSurvey struct {
	scanner *bufio.Scanner
}

// AskOne is a wrapper for the survey.AskOne() func
func (d *DumbSurvey) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	fmt.Printf("%v ", getStringFromPrompt(p))

	if d.scanner == nil {
		d.scanner = bufio.NewScanner(os.Stdin)
	}

	d.scanner.Scan()
	if d.scanner.Err() != nil {
		return fmt.Errorf("error reading %v - %v", d.scanner.Err(), d.scanner.Text())
	}
	fmt.Println(d.scanner.Text())
	//fmt.Println(d.scanner.Text())
	//fmt.Println(d.scanner.Bytes())
	response = d.scanner.Text()

	return nil
}

// Ask is a wrapper for the survey.Ask() func
func (d DumbSurvey) Ask(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
	answers := &map[string]string{}
	for _, q := range qs {
		resp := ""
		err := d.AskOne(q.Prompt, &resp, opts...)
		if err != nil {
			return err
		}

		(*answers)[q.Name] = resp
	}

	response = answers
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
