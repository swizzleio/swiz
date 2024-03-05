// Code generated by mockery v2.20.0. DO NOT EDIT.

package mockaws

import (
	context "context"

	sts "github.com/aws/aws-sdk-go-v2/service/sts"
	mock "github.com/stretchr/testify/mock"
)

// Stser is an autogenerated mock type for the Stser type
type Stser struct {
	mock.Mock
}

// GetCallerIdentity provides a mock function with given fields: ctx, params, optFns
func (_m *Stser) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *sts.GetCallerIdentityOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.Options)) *sts.GetCallerIdentityOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sts.GetCallerIdentityOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *sts.GetCallerIdentityInput, ...func(*sts.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewStser interface {
	mock.TestingT
	Cleanup(func())
}

// NewStser creates a new instance of Stser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStser(t mockConstructorTestingTNewStser) *Stser {
	mock := &Stser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
