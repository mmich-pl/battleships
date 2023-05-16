// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	models "battleships/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// BattleshipClient is an autogenerated mock type for the BattleshipClient type
type BattleshipClient struct {
	mock.Mock
}

// Board provides a mock function with given fields: endpoint
func (_m *BattleshipClient) Board(endpoint string) ([]string, error) {
	ret := _m.Called(endpoint)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(endpoint)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(endpoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(endpoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Description provides a mock function with given fields: endpoint
func (_m *BattleshipClient) Description(endpoint string) (*models.DescriptionResponse, error) {
	ret := _m.Called(endpoint)

	var r0 *models.DescriptionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.DescriptionResponse, error)); ok {
		return rf(endpoint)
	}
	if rf, ok := ret.Get(0).(func(string) *models.DescriptionResponse); ok {
		r0 = rf(endpoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DescriptionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(endpoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GameStatus provides a mock function with given fields: endpoint
func (_m *BattleshipClient) GameStatus(endpoint string) (*models.StatusResponse, error) {
	ret := _m.Called(endpoint)

	var r0 *models.StatusResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.StatusResponse, error)); ok {
		return rf(endpoint)
	}
	if rf, ok := ret.Get(0).(func(string) *models.StatusResponse); ok {
		r0 = rf(endpoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.StatusResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(endpoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitGame provides a mock function with given fields: endpoint, nick, desc, targetNick, wpbot
func (_m *BattleshipClient) InitGame(endpoint string, nick string, desc string, targetNick string, wpbot bool) error {
	ret := _m.Called(endpoint, nick, desc, targetNick, wpbot)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, bool) error); ok {
		r0 = rf(endpoint, nick, desc, targetNick, wpbot)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBattleshipClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewBattleshipClient creates a new instance of BattleshipClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBattleshipClient(t mockConstructorTestingTNewBattleshipClient) *BattleshipClient {
	mock := &BattleshipClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}