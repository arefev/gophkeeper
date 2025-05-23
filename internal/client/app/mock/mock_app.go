// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/client/app/app.go

// Package mock_app is a generated GoMock package.
package mock_app

import (
	context "context"
	reflect "reflect"

	model "github.com/arefev/gophkeeper/internal/client/tui/model"
	tea "github.com/charmbracelet/bubbletea"
	gomock "github.com/golang/mock/gomock"
)

// MockConnection is a mock of Connection interface.
type MockConnection struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionMockRecorder
}

// MockConnectionMockRecorder is the mock recorder for MockConnection.
type MockConnectionMockRecorder struct {
	mock *MockConnection
}

// NewMockConnection creates a new mock instance.
func NewMockConnection(ctrl *gomock.Controller) *MockConnection {
	mock := &MockConnection{ctrl: ctrl}
	mock.recorder = &MockConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnection) EXPECT() *MockConnectionMockRecorder {
	return m.recorder
}

// CheckTokenCmd mocks base method.
func (m *MockConnection) CheckTokenCmd() tea.Msg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTokenCmd")
	ret0, _ := ret[0].(tea.Msg)
	return ret0
}

// CheckTokenCmd indicates an expected call of CheckTokenCmd.
func (mr *MockConnectionMockRecorder) CheckTokenCmd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTokenCmd", reflect.TypeOf((*MockConnection)(nil).CheckTokenCmd))
}

// Delete mocks base method.
func (m *MockConnection) Delete(ctx context.Context, uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockConnectionMockRecorder) Delete(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockConnection)(nil).Delete), ctx, uuid)
}

// FileDownload mocks base method.
func (m *MockConnection) FileDownload(ctx context.Context, uuid string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileDownload", ctx, uuid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FileDownload indicates an expected call of FileDownload.
func (mr *MockConnectionMockRecorder) FileDownload(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileDownload", reflect.TypeOf((*MockConnection)(nil).FileDownload), ctx, uuid)
}

// FileUpload mocks base method.
func (m *MockConnection) FileUpload(ctx context.Context, path, metaName, metaType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileUpload", ctx, path, metaName, metaType)
	ret0, _ := ret[0].(error)
	return ret0
}

// FileUpload indicates an expected call of FileUpload.
func (mr *MockConnectionMockRecorder) FileUpload(ctx, path, metaName, metaType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileUpload", reflect.TypeOf((*MockConnection)(nil).FileUpload), ctx, path, metaName, metaType)
}

// GetList mocks base method.
func (m *MockConnection) GetList(ctx context.Context) (*[]model.MetaListData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx)
	ret0, _ := ret[0].(*[]model.MetaListData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockConnectionMockRecorder) GetList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockConnection)(nil).GetList), ctx)
}

// Login mocks base method.
func (m *MockConnection) Login(ctx context.Context, login, pwd string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, login, pwd)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockConnectionMockRecorder) Login(ctx, login, pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockConnection)(nil).Login), ctx, login, pwd)
}

// Register mocks base method.
func (m *MockConnection) Register(ctx context.Context, login, pwd string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, login, pwd)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockConnectionMockRecorder) Register(ctx, login, pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockConnection)(nil).Register), ctx, login, pwd)
}

// SetToken mocks base method.
func (m *MockConnection) SetToken(t string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetToken", t)
}

// SetToken indicates an expected call of SetToken.
func (mr *MockConnectionMockRecorder) SetToken(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetToken", reflect.TypeOf((*MockConnection)(nil).SetToken), t)
}

// TextUpload mocks base method.
func (m *MockConnection) TextUpload(ctx context.Context, txt []byte, metaName, metaType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TextUpload", ctx, txt, metaName, metaType)
	ret0, _ := ret[0].(error)
	return ret0
}

// TextUpload indicates an expected call of TextUpload.
func (mr *MockConnectionMockRecorder) TextUpload(ctx, txt, metaName, metaType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TextUpload", reflect.TypeOf((*MockConnection)(nil).TextUpload), ctx, txt, metaName, metaType)
}
