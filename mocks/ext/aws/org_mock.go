// Code generated by mockery v2.20.0. DO NOT EDIT.

package mockaws

import (
	context "context"

	organizations "github.com/aws/aws-sdk-go-v2/service/organizations"
	mock "github.com/stretchr/testify/mock"
)

// Orger is an autogenerated mock type for the Orger type
type Orger struct {
	mock.Mock
}

// ListAccounts provides a mock function with given fields: _a0, _a1, _a2
func (_m *Orger) ListAccounts(_a0 context.Context, _a1 *organizations.ListAccountsInput, _a2 ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *organizations.ListAccountsOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *organizations.ListAccountsInput, ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *organizations.ListAccountsInput, ...func(*organizations.Options)) *organizations.ListAccountsOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*organizations.ListAccountsOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *organizations.ListAccountsInput, ...func(*organizations.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewOrger interface {
	mock.TestingT
	Cleanup(func())
}

// NewOrger creates a new instance of Orger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOrger(t mockConstructorTestingTNewOrger) *Orger {
	mock := &Orger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
