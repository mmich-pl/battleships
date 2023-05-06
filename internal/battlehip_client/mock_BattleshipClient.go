// Code generated by mockery v2.23.1. DO NOT EDIT.

package battlehip_client

import (
	models "battleships/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// MockBattleshipClient is an autogenerated mock type for the BattleshipClient type
type MockBattleshipClient struct {
	mock.Mock
}

type MockBattleshipClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBattleshipClient) EXPECT() *MockBattleshipClient_Expecter {
	return &MockBattleshipClient_Expecter{mock: &_m.Mock}
}

// Board provides a mock function with given fields: endpoint
func (_m *MockBattleshipClient) Board(endpoint string) ([]string, error) {
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

// MockBattleshipClient_Board_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Board'
type MockBattleshipClient_Board_Call struct {
	*mock.Call
}

// Board is a helper method to define mock.On call
//   - endpoint string
func (_e *MockBattleshipClient_Expecter) Board(endpoint interface{}) *MockBattleshipClient_Board_Call {
	return &MockBattleshipClient_Board_Call{Call: _e.mock.On("Board", endpoint)}
}

func (_c *MockBattleshipClient_Board_Call) Run(run func(endpoint string)) *MockBattleshipClient_Board_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockBattleshipClient_Board_Call) Return(_a0 []string, _a1 error) *MockBattleshipClient_Board_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBattleshipClient_Board_Call) RunAndReturn(run func(string) ([]string, error)) *MockBattleshipClient_Board_Call {
	_c.Call.Return(run)
	return _c
}

// Description provides a mock function with given fields: endpoint
func (_m *MockBattleshipClient) Description(endpoint string) (*models.DescriptionResponse, error) {
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

// MockBattleshipClient_Description_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Description'
type MockBattleshipClient_Description_Call struct {
	*mock.Call
}

// Description is a helper method to define mock.On call
//   - endpoint string
func (_e *MockBattleshipClient_Expecter) Description(endpoint interface{}) *MockBattleshipClient_Description_Call {
	return &MockBattleshipClient_Description_Call{Call: _e.mock.On("Description", endpoint)}
}

func (_c *MockBattleshipClient_Description_Call) Run(run func(endpoint string)) *MockBattleshipClient_Description_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockBattleshipClient_Description_Call) Return(_a0 *models.DescriptionResponse, _a1 error) *MockBattleshipClient_Description_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBattleshipClient_Description_Call) RunAndReturn(run func(string) (*models.DescriptionResponse, error)) *MockBattleshipClient_Description_Call {
	_c.Call.Return(run)
	return _c
}

// Fire provides a mock function with given fields: endpoint, coords
func (_m *MockBattleshipClient) Fire(endpoint string, coords string) (*models.ShootResult, error) {
	ret := _m.Called(endpoint, coords)

	var r0 *models.ShootResult
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*models.ShootResult, error)); ok {
		return rf(endpoint, coords)
	}
	if rf, ok := ret.Get(0).(func(string, string) *models.ShootResult); ok {
		r0 = rf(endpoint, coords)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ShootResult)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(endpoint, coords)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBattleshipClient_Fire_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fire'
type MockBattleshipClient_Fire_Call struct {
	*mock.Call
}

// Fire is a helper method to define mock.On call
//   - endpoint string
//   - coords string
func (_e *MockBattleshipClient_Expecter) Fire(endpoint interface{}, coords interface{}) *MockBattleshipClient_Fire_Call {
	return &MockBattleshipClient_Fire_Call{Call: _e.mock.On("Fire", endpoint, coords)}
}

func (_c *MockBattleshipClient_Fire_Call) Run(run func(endpoint string, coords string)) *MockBattleshipClient_Fire_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockBattleshipClient_Fire_Call) Return(_a0 *models.ShootResult, _a1 error) *MockBattleshipClient_Fire_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBattleshipClient_Fire_Call) RunAndReturn(run func(string, string) (*models.ShootResult, error)) *MockBattleshipClient_Fire_Call {
	_c.Call.Return(run)
	return _c
}

// GameStatus provides a mock function with given fields: endpoint
func (_m *MockBattleshipClient) GameStatus(endpoint string) (*models.StatusResponse, error) {
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

// MockBattleshipClient_GameStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GameStatus'
type MockBattleshipClient_GameStatus_Call struct {
	*mock.Call
}

// GameStatus is a helper method to define mock.On call
//   - endpoint string
func (_e *MockBattleshipClient_Expecter) GameStatus(endpoint interface{}) *MockBattleshipClient_GameStatus_Call {
	return &MockBattleshipClient_GameStatus_Call{Call: _e.mock.On("GameStatus", endpoint)}
}

func (_c *MockBattleshipClient_GameStatus_Call) Run(run func(endpoint string)) *MockBattleshipClient_GameStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockBattleshipClient_GameStatus_Call) Return(_a0 *models.StatusResponse, _a1 error) *MockBattleshipClient_GameStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBattleshipClient_GameStatus_Call) RunAndReturn(run func(string) (*models.StatusResponse, error)) *MockBattleshipClient_GameStatus_Call {
	_c.Call.Return(run)
	return _c
}

// InitGame provides a mock function with given fields: endpoint, nick, desc, targetNick, wpbot
func (_m *MockBattleshipClient) InitGame(endpoint string, nick string, desc string, targetNick string, wpbot bool) error {
	ret := _m.Called(endpoint, nick, desc, targetNick, wpbot)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, bool) error); ok {
		r0 = rf(endpoint, nick, desc, targetNick, wpbot)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBattleshipClient_InitGame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitGame'
type MockBattleshipClient_InitGame_Call struct {
	*mock.Call
}

// InitGame is a helper method to define mock.On call
//   - endpoint string
//   - nick string
//   - desc string
//   - targetNick string
//   - wpbot bool
func (_e *MockBattleshipClient_Expecter) InitGame(endpoint interface{}, nick interface{}, desc interface{}, targetNick interface{}, wpbot interface{}) *MockBattleshipClient_InitGame_Call {
	return &MockBattleshipClient_InitGame_Call{Call: _e.mock.On("InitGame", endpoint, nick, desc, targetNick, wpbot)}
}

func (_c *MockBattleshipClient_InitGame_Call) Run(run func(endpoint string, nick string, desc string, targetNick string, wpbot bool)) *MockBattleshipClient_InitGame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(string), args[4].(bool))
	})
	return _c
}

func (_c *MockBattleshipClient_InitGame_Call) Return(_a0 error) *MockBattleshipClient_InitGame_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBattleshipClient_InitGame_Call) RunAndReturn(run func(string, string, string, string, bool) error) *MockBattleshipClient_InitGame_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockBattleshipClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockBattleshipClient creates a new instance of MockBattleshipClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockBattleshipClient(t mockConstructorTestingTNewMockBattleshipClient) *MockBattleshipClient {
	mock := &MockBattleshipClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
