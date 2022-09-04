// Code generated by MockGen. DO NOT EDIT.
// Source: friend_list_usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	model "problem1/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
)

// MockFriendListUseCase is a mock of FriendListUseCase interface.
type MockFriendListUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockFriendListUseCaseMockRecorder
}

// MockFriendListUseCaseMockRecorder is the mock recorder for MockFriendListUseCase.
type MockFriendListUseCaseMockRecorder struct {
	mock *MockFriendListUseCase
}

// NewMockFriendListUseCase creates a new mock instance.
func NewMockFriendListUseCase(ctrl *gomock.Controller) *MockFriendListUseCase {
	mock := &MockFriendListUseCase{ctrl: ctrl}
	mock.recorder = &MockFriendListUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriendListUseCase) EXPECT() *MockFriendListUseCaseMockRecorder {
	return m.recorder
}

// GetFriendListByUserId mocks base method.
func (m *MockFriendListUseCase) GetFriendListByUserId(c echo.Context) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendListByUserId", c)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendListByUserId indicates an expected call of GetFriendListByUserId.
func (mr *MockFriendListUseCaseMockRecorder) GetFriendListByUserId(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendListByUserId", reflect.TypeOf((*MockFriendListUseCase)(nil).GetFriendListByUserId), c)
}