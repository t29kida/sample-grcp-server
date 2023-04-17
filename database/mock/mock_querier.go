// Code generated by MockGen. DO NOT EDIT.
// Source: querier.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"
	database "sample-grpc-server/database"

	gomock "github.com/golang/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreateArticle mocks base method.
func (m *MockQuerier) CreateArticle(arg0 context.Context, arg1 database.CreateArticleParams) (*database.CreateArticleResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateArticle", arg0, arg1)
	ret0, _ := ret[0].(*database.CreateArticleResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateArticle indicates an expected call of CreateArticle.
func (mr *MockQuerierMockRecorder) CreateArticle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateArticle", reflect.TypeOf((*MockQuerier)(nil).CreateArticle), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockQuerier) CreateSession(arg0 context.Context, arg1 database.CreateSessionParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockQuerierMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockQuerier)(nil).CreateSession), arg0, arg1)
}

// DeleteArticle mocks base method.
func (m *MockQuerier) DeleteArticle(arg0 context.Context, arg1 database.DeleteArticleParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteArticle", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteArticle indicates an expected call of DeleteArticle.
func (mr *MockQuerierMockRecorder) DeleteArticle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteArticle", reflect.TypeOf((*MockQuerier)(nil).DeleteArticle), arg0, arg1)
}

// GetArticle mocks base method.
func (m *MockQuerier) GetArticle(arg0 context.Context, arg1 database.GetArticleParams) (*database.GetArticleResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArticle", arg0, arg1)
	ret0, _ := ret[0].(*database.GetArticleResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArticle indicates an expected call of GetArticle.
func (mr *MockQuerierMockRecorder) GetArticle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArticle", reflect.TypeOf((*MockQuerier)(nil).GetArticle), arg0, arg1)
}

// GetArticles mocks base method.
func (m *MockQuerier) GetArticles(arg0 context.Context, arg1 database.GetArticlesParams) (*database.GetArticlesResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArticles", arg0, arg1)
	ret0, _ := ret[0].(*database.GetArticlesResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArticles indicates an expected call of GetArticles.
func (mr *MockQuerierMockRecorder) GetArticles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArticles", reflect.TypeOf((*MockQuerier)(nil).GetArticles), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockQuerier) GetSession(arg0 context.Context, arg1 database.GetSessionParams) (*database.GetSessionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(*database.GetSessionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockQuerierMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockQuerier)(nil).GetSession), arg0, arg1)
}

// Login mocks base method.
func (m *MockQuerier) Login(arg0 context.Context, arg1 database.LoginParams) (*database.LoginResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*database.LoginResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockQuerierMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockQuerier)(nil).Login), arg0, arg1)
}

// SignUp mocks base method.
func (m *MockQuerier) SignUp(arg0 context.Context, arg1 database.SignUpParams) (*database.SignUpResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0, arg1)
	ret0, _ := ret[0].(*database.SignUpResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockQuerierMockRecorder) SignUp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockQuerier)(nil).SignUp), arg0, arg1)
}

// UpdateArticle mocks base method.
func (m *MockQuerier) UpdateArticle(arg0 context.Context, arg1 database.UpdateArticleParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateArticle", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateArticle indicates an expected call of UpdateArticle.
func (mr *MockQuerierMockRecorder) UpdateArticle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateArticle", reflect.TypeOf((*MockQuerier)(nil).UpdateArticle), arg0, arg1)
}