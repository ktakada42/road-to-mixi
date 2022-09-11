// Code generated by MockGen. DO NOT EDIT.
// Source: friend_list_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	model "problem1/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
)

// MockFriendListService is a mock of FriendListService interface.
type MockFriendListService struct {
	ctrl     *gomock.Controller
	recorder *MockFriendListServiceMockRecorder
}

// MockFriendListServiceMockRecorder is the mock recorder for MockFriendListService.
type MockFriendListServiceMockRecorder struct {
	mock *MockFriendListService
}

// NewMockFriendListService creates a new mock instance.
func NewMockFriendListService(ctrl *gomock.Controller) *MockFriendListService {
	mock := &MockFriendListService{ctrl: ctrl}
	mock.recorder = &MockFriendListServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriendListService) EXPECT() *MockFriendListServiceMockRecorder {
	return m.recorder
}

// CheckUserExist mocks base method.
func (m *MockFriendListService) CheckUserExist(userId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserExist", userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserExist indicates an expected call of CheckUserExist.
func (mr *MockFriendListServiceMockRecorder) CheckUserExist(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserExist", reflect.TypeOf((*MockFriendListService)(nil).CheckUserExist), userId)
}

// GetFriendListByUserId mocks base method.
func (m *MockFriendListService) GetFriendListByUserId(c echo.Context) (*model.FriendList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendListByUserId", c)
	ret0, _ := ret[0].(*model.FriendList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendListByUserId indicates an expected call of GetFriendListByUserId.
func (mr *MockFriendListServiceMockRecorder) GetFriendListByUserId(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendListByUserId", reflect.TypeOf((*MockFriendListService)(nil).GetFriendListByUserId), c)
}

// GetFriendListOfFriendsByUserId mocks base method.
func (m *MockFriendListService) GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendListOfFriendsByUserId", c)
	ret0, _ := ret[0].(*model.FriendList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendListOfFriendsByUserId indicates an expected call of GetFriendListOfFriendsByUserId.
func (mr *MockFriendListServiceMockRecorder) GetFriendListOfFriendsByUserId(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendListOfFriendsByUserId", reflect.TypeOf((*MockFriendListService)(nil).GetFriendListOfFriendsByUserId), c)
}

// GetFriendListOfFriendsByUserIdWithPaging mocks base method.
func (m *MockFriendListService) GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) (*model.FriendList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendListOfFriendsByUserIdWithPaging", c)
	ret0, _ := ret[0].(*model.FriendList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendListOfFriendsByUserIdWithPaging indicates an expected call of GetFriendListOfFriendsByUserIdWithPaging.
func (mr *MockFriendListServiceMockRecorder) GetFriendListOfFriendsByUserIdWithPaging(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendListOfFriendsByUserIdWithPaging", reflect.TypeOf((*MockFriendListService)(nil).GetFriendListOfFriendsByUserIdWithPaging), c)
}

// InsertUserLink mocks base method.
func (m *MockFriendListService) InsertUserLink(ulfr *model.UserLinkForRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUserLink", ulfr)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertUserLink indicates an expected call of InsertUserLink.
func (mr *MockFriendListServiceMockRecorder) InsertUserLink(ulfr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUserLink", reflect.TypeOf((*MockFriendListService)(nil).InsertUserLink), ulfr)
}
