package appcli

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"os"
)

// This provides a more testable version of survey, instead of smart prompting, it just uses a scanner

type SurveyResponse struct {
	Prompt           string
	ResponseList     []string
	NextResponse     int
	FuzzyPromptMatch bool
	ActualPrompts    []string
}

type DumbSurvey struct {
	scanner      *bufio.Reader
	responseList []SurveyResponse
}

func NewDumbSurvey() *DumbSurvey {
	return &DumbSurvey{
		scanner:      bufio.NewReader(os.Stdin),
		responseList: []SurveyResponse{},
	}
}

func (d *DumbSurvey) AddResponses(responses []SurveyResponse) {
	for _, resp := range responses {
		_, err := d.getResponse(resp.Prompt)
		if err != nil {
			// Only add if the prompt is not there
			d.responseList = append(d.responseList, resp)
		}
	}
}

func (d *DumbSurvey) AddResponse(response SurveyResponse) {
	d.AddResponses([]SurveyResponse{response})
}

func (d DumbSurvey) GetPrompts(prompt string) ([]string, error) {
	resp, err := d.getResponseObj(prompt)
	if err != nil {
		return []string{}, err
	}

	return resp.ActualPrompts, nil
}

// AskOne is a wrapper for the survey.AskOne() func
func (d *DumbSurvey) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	prompt := getStringFromPrompt(p)
	//fmt.Printf("%v ", prompt)

	resp, err := d.getResponse(prompt)
	if err != nil {
		return err
	}

	pResponse := response.(*string)
	*pResponse = resp

	//strReader := strings.NewReader("dev-sandbox\n012345678901\nus-east-2\nexample.com\nLogLevel,Apm\nDebug\nEnabled\n\n\n")
	//fmt.Println(resp)

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

func (d *DumbSurvey) getResponseObj(prompt string) (*SurveyResponse, error) {
	var resp *SurveyResponse
	var err error
	found := false
	for _, r := range d.responseList {
		if r.FuzzyPromptMatch && fuzzy.Match(prompt, r.Prompt) {
			resp = &r
			found = true
		} else if prompt == r.Prompt {
			resp = &r
			found = true
		}

		if found {
			break
		}
	}

	if !found {
		err = fmt.Errorf("could not find response")
	}

	return resp, err
}

func (d *DumbSurvey) getResponse(prompt string) (string, error) {
	resp, err := d.getResponseObj(prompt)
	if err != nil {
		return "", err
	}

	respStr := resp.ResponseList[resp.NextResponse]
	resp.NextResponse = (resp.NextResponse + 1) % len(resp.ResponseList)

	return respStr, err
}
