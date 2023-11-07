// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// NetStorage is an autogenerated mock type for the netStorage type
type NetStorage struct {
	mock.Mock
}

type NetStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *NetStorage) EXPECT() *NetStorage_Expecter {
	return &NetStorage_Expecter{mock: &_m.Mock}
}

// AddNet provides a mock function with given fields: ctx, ip, list
func (_m *NetStorage) AddNet(ctx context.Context, ip string, list string) error {
	ret := _m.Called(ctx, ip, list)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, ip, list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NetStorage_AddNet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddNet'
type NetStorage_AddNet_Call struct {
	*mock.Call
}

// AddNet is a helper method to define mock.On call
//   - ctx context.Context
//   - ip string
//   - list string
func (_e *NetStorage_Expecter) AddNet(ctx interface{}, ip interface{}, list interface{}) *NetStorage_AddNet_Call {
	return &NetStorage_AddNet_Call{Call: _e.mock.On("AddNet", ctx, ip, list)}
}

func (_c *NetStorage_AddNet_Call) Run(run func(ctx context.Context, ip string, list string)) *NetStorage_AddNet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *NetStorage_AddNet_Call) Return(_a0 error) *NetStorage_AddNet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NetStorage_AddNet_Call) RunAndReturn(run func(context.Context, string, string) error) *NetStorage_AddNet_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteNet provides a mock function with given fields: ctx, ip
func (_m *NetStorage) DeleteNet(ctx context.Context, ip string) error {
	ret := _m.Called(ctx, ip)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, ip)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NetStorage_DeleteNet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteNet'
type NetStorage_DeleteNet_Call struct {
	*mock.Call
}

// DeleteNet is a helper method to define mock.On call
//   - ctx context.Context
//   - ip string
func (_e *NetStorage_Expecter) DeleteNet(ctx interface{}, ip interface{}) *NetStorage_DeleteNet_Call {
	return &NetStorage_DeleteNet_Call{Call: _e.mock.On("DeleteNet", ctx, ip)}
}

func (_c *NetStorage_DeleteNet_Call) Run(run func(ctx context.Context, ip string)) *NetStorage_DeleteNet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *NetStorage_DeleteNet_Call) Return(_a0 error) *NetStorage_DeleteNet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NetStorage_DeleteNet_Call) RunAndReturn(run func(context.Context, string) error) *NetStorage_DeleteNet_Call {
	_c.Call.Return(run)
	return _c
}

// Net provides a mock function with given fields: ctx, ip
func (_m *NetStorage) Net(ctx context.Context, ip string) (string, error) {
	ret := _m.Called(ctx, ip)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, ip)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, ip)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetStorage_Net_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Net'
type NetStorage_Net_Call struct {
	*mock.Call
}

// Net is a helper method to define mock.On call
//   - ctx context.Context
//   - ip string
func (_e *NetStorage_Expecter) Net(ctx interface{}, ip interface{}) *NetStorage_Net_Call {
	return &NetStorage_Net_Call{Call: _e.mock.On("Net", ctx, ip)}
}

func (_c *NetStorage_Net_Call) Run(run func(ctx context.Context, ip string)) *NetStorage_Net_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *NetStorage_Net_Call) Return(_a0 string, _a1 error) *NetStorage_Net_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NetStorage_Net_Call) RunAndReturn(run func(context.Context, string) (string, error)) *NetStorage_Net_Call {
	_c.Call.Return(run)
	return _c
}

// NewNetStorage creates a new instance of NetStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNetStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *NetStorage {
	mock := &NetStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}