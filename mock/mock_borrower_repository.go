// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository/borrower.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/frencius/loan-service/model"
	gomock "github.com/golang/mock/gomock"
)

// MockIBorrowerRepository is a mock of IBorrowerRepository interface.
type MockIBorrowerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIBorrowerRepositoryMockRecorder
}

// MockIBorrowerRepositoryMockRecorder is the mock recorder for MockIBorrowerRepository.
type MockIBorrowerRepositoryMockRecorder struct {
	mock *MockIBorrowerRepository
}

// NewMockIBorrowerRepository creates a new mock instance.
func NewMockIBorrowerRepository(ctrl *gomock.Controller) *MockIBorrowerRepository {
	mock := &MockIBorrowerRepository{ctrl: ctrl}
	mock.recorder = &MockIBorrowerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBorrowerRepository) EXPECT() *MockIBorrowerRepositoryMockRecorder {
	return m.recorder
}

// GetBorrowerByID mocks base method.
func (m *MockIBorrowerRepository) GetBorrowerByID(ctx context.Context, id string) (*model.Borrower, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBorrowerByID", ctx, id)
	ret0, _ := ret[0].(*model.Borrower)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBorrowerByID indicates an expected call of GetBorrowerByID.
func (mr *MockIBorrowerRepositoryMockRecorder) GetBorrowerByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBorrowerByID", reflect.TypeOf((*MockIBorrowerRepository)(nil).GetBorrowerByID), ctx, id)
}
