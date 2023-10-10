// Code generated by mockery v2.20.0. DO NOT EDIT.

package mockaws

import (
	context "context"

	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation"

	mock "github.com/stretchr/testify/mock"
)

// CfDescribeStacksPaginator is an autogenerated mock type for the CfDescribeStacksPaginator type
type CfDescribeStacksPaginator struct {
	mock.Mock
}

// HasMorePages provides a mock function with given fields:
func (_m *CfDescribeStacksPaginator) HasMorePages() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NextPage provides a mock function with given fields: ctx, optFns
func (_m *CfDescribeStacksPaginator) NextPage(ctx context.Context, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *cloudformation.DescribeStacksOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)); ok {
		return rf(ctx, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...func(*cloudformation.Options)) *cloudformation.DescribeStacksOutput); ok {
		r0 = rf(ctx, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudformation.DescribeStacksOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...func(*cloudformation.Options)) error); ok {
		r1 = rf(ctx, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCfDescribeStacksPaginator interface {
	mock.TestingT
	Cleanup(func())
}

// NewCfDescribeStacksPaginator creates a new instance of CfDescribeStacksPaginator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCfDescribeStacksPaginator(t mockConstructorTestingTNewCfDescribeStacksPaginator) *CfDescribeStacksPaginator {
	mock := &CfDescribeStacksPaginator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
