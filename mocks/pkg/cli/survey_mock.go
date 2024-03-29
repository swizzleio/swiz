// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocksurvey

import (
	survey "github.com/AlecAivazis/survey/v2"
	mock "github.com/stretchr/testify/mock"
)

// SurveyWrapper is an autogenerated mock type for the SurveyWrapper type
type SurveyWrapper struct {
	mock.Mock
}

// Ask provides a mock function with given fields: qs, response, opts
func (_m *SurveyWrapper) Ask(qs []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, qs, response)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*survey.Question, interface{}, ...survey.AskOpt) error); ok {
		r0 = rf(qs, response, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AskOne provides a mock function with given fields: p, response, opts
func (_m *SurveyWrapper) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, p, response)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(survey.Prompt, interface{}, ...survey.AskOpt) error); ok {
		r0 = rf(p, response, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSurveyWrapper interface {
	mock.TestingT
	Cleanup(func())
}

// NewSurveyWrapper creates a new instance of SurveyWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSurveyWrapper(t mockConstructorTestingTNewSurveyWrapper) *SurveyWrapper {
	mock := &SurveyWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
