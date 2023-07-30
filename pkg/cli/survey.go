package appcli

import "github.com/AlecAivazis/survey/v2"

// SurveyWrapper is a wrapper around the survey package. This makes the cli package unit testable.
//
//go:generate mockery --name SurveyWrapper --filename survey_mock.go --output ../../mocks/pkg/cli --outpkg mocksurvey
type SurveyWrapper interface {
	AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
	Ask(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error
}

// SurveyWrap wraps the survey package
type SurveyWrap struct {
}

// AskOne is a wrapper for the survey.AskOne() func
func (SurveyWrap) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, opts...)
}

// Ask is a wrapper for the survey.Ask() func
func (SurveyWrap) Ask(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
	return survey.Ask(qs, response, opts...)
}
