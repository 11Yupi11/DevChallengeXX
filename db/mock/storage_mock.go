// Code generated by MockGen. DO NOT EDIT.
// Source: ./storage.go

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	sql "database/sql"
	db "dev-challenge/db"
	models "dev-challenge/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddCellInput mocks base method.
func (m *MockStorage) AddCellInput(ctx context.Context, tx *sql.Tx, data db.Input) (*models.Data, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCellInput", ctx, tx, data)
	ret0, _ := ret[0].(*models.Data)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddCellInput indicates an expected call of AddCellInput.
func (mr *MockStorageMockRecorder) AddCellInput(ctx, tx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCellInput", reflect.TypeOf((*MockStorage)(nil).AddCellInput), ctx, tx, data)
}

// BeginTransaction mocks base method.
func (m *MockStorage) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTransaction", ctx)
	ret0, _ := ret[0].(*sql.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTransaction indicates an expected call of BeginTransaction.
func (mr *MockStorageMockRecorder) BeginTransaction(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTransaction", reflect.TypeOf((*MockStorage)(nil).BeginTransaction), ctx)
}

// GetCellInput mocks base method.
func (m *MockStorage) GetCellInput(ctx context.Context, sheetID, cellID string) (*models.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCellInput", ctx, sheetID, cellID)
	ret0, _ := ret[0].(*models.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCellInput indicates an expected call of GetCellInput.
func (mr *MockStorageMockRecorder) GetCellInput(ctx, sheetID, cellID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCellInput", reflect.TypeOf((*MockStorage)(nil).GetCellInput), ctx, sheetID, cellID)
}

// GetCellInputBatch mocks base method.
func (m *MockStorage) GetCellInputBatch(ctx context.Context, tx *sql.Tx, sheetID string, cells []string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCellInputBatch", ctx, tx, sheetID, cells)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCellInputBatch indicates an expected call of GetCellInputBatch.
func (mr *MockStorageMockRecorder) GetCellInputBatch(ctx, tx, sheetID, cells interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCellInputBatch", reflect.TypeOf((*MockStorage)(nil).GetCellInputBatch), ctx, tx, sheetID, cells)
}

// GetIDList mocks base method.
func (m *MockStorage) GetIDList(ctx context.Context, tx *sql.Tx, cellID string) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIDList", ctx, tx, cellID)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIDList indicates an expected call of GetIDList.
func (mr *MockStorageMockRecorder) GetIDList(ctx, tx, cellID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIDList", reflect.TypeOf((*MockStorage)(nil).GetIDList), ctx, tx, cellID)
}

// GetInputBatchByIDs mocks base method.
func (m *MockStorage) GetInputBatchByIDs(ctx context.Context, tx *sql.Tx, IDs []int) (*[]db.Input, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInputBatchByIDs", ctx, tx, IDs)
	ret0, _ := ret[0].(*[]db.Input)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInputBatchByIDs indicates an expected call of GetInputBatchByIDs.
func (mr *MockStorageMockRecorder) GetInputBatchByIDs(ctx, tx, IDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInputBatchByIDs", reflect.TypeOf((*MockStorage)(nil).GetInputBatchByIDs), ctx, tx, IDs)
}

// GetSheetInput mocks base method.
func (m *MockStorage) GetSheetInput(ctx context.Context, sheetID string) (map[string]models.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSheetInput", ctx, sheetID)
	ret0, _ := ret[0].(map[string]models.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSheetInput indicates an expected call of GetSheetInput.
func (mr *MockStorageMockRecorder) GetSheetInput(ctx, sheetID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSheetInput", reflect.TypeOf((*MockStorage)(nil).GetSheetInput), ctx, sheetID)
}
