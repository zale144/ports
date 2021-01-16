// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/zale144/ports/portdomainservice/internal/model"
)

// MockStore is an autogenerated mock type for the store type
type MockStore struct {
	mock.Mock
}

// GetPorts provides a mock function with given fields: ctx
func (_m *MockStore) GetPorts(ctx context.Context) ([]model.Port, error) {
	ret := _m.Called(ctx)

	var r0 []model.Port
	if rf, ok := ret.Get(0).(func(context.Context) []model.Port); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Port)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SavePorts provides a mock function with given fields: ctx, ports
func (_m *MockStore) SavePorts(ctx context.Context, ports []model.Port) error {
	ret := _m.Called(ctx, ports)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []model.Port) error); ok {
		r0 = rf(ctx, ports)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}