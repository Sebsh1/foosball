// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package player is a generated GoMock package.
package player

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreatePlayer mocks base method.
func (m *MockService) CreatePlayer(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePlayer", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePlayer indicates an expected call of CreatePlayer.
func (mr *MockServiceMockRecorder) CreatePlayer(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlayer", reflect.TypeOf((*MockService)(nil).CreatePlayer), ctx, name)
}

// DeletePlayer mocks base method.
func (m *MockService) DeletePlayer(ctx context.Context, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlayer", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlayer indicates an expected call of DeletePlayer.
func (mr *MockServiceMockRecorder) DeletePlayer(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlayer", reflect.TypeOf((*MockService)(nil).DeletePlayer), ctx, id)
}

// GetPlayer mocks base method.
func (m *MockService) GetPlayer(ctx context.Context, id uint) (*Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlayer", ctx, id)
	ret0, _ := ret[0].(*Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlayer indicates an expected call of GetPlayer.
func (mr *MockServiceMockRecorder) GetPlayer(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlayer", reflect.TypeOf((*MockService)(nil).GetPlayer), ctx, id)
}

// GetPlayers mocks base method.
func (m *MockService) GetPlayers(ctx context.Context, ids []uint) ([]*Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlayers", ctx, ids)
	ret0, _ := ret[0].([]*Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlayers indicates an expected call of GetPlayers.
func (mr *MockServiceMockRecorder) GetPlayers(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlayers", reflect.TypeOf((*MockService)(nil).GetPlayers), ctx, ids)
}

// GetTopPlayersByRating mocks base method.
func (m *MockService) GetTopPlayersByRating(ctx context.Context, top int) ([]*Player, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopPlayersByRating", ctx, top)
	ret0, _ := ret[0].([]*Player)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopPlayersByRating indicates an expected call of GetTopPlayersByRating.
func (mr *MockServiceMockRecorder) GetTopPlayersByRating(ctx, top interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopPlayersByRating", reflect.TypeOf((*MockService)(nil).GetTopPlayersByRating), ctx, top)
}

// UpdatePlayers mocks base method.
func (m *MockService) UpdatePlayers(ctx context.Context, players []*Player, ratingChange []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePlayers", ctx, players, ratingChange)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePlayers indicates an expected call of UpdatePlayers.
func (mr *MockServiceMockRecorder) UpdatePlayers(ctx, players, ratingChange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePlayers", reflect.TypeOf((*MockService)(nil).UpdatePlayers), ctx, players, ratingChange)
}
