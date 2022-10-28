// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/calculator_service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	decimal "github.com/shopspring/decimal"
)

// MockTransactionStore is a mock of TransactionStore interface.
type MockTransactionStore struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionStoreMockRecorder
}

// MockTransactionStoreMockRecorder is the mock recorder for MockTransactionStore.
type MockTransactionStoreMockRecorder struct {
	mock *MockTransactionStore
}

// NewMockTransactionStore creates a new mock instance.
func NewMockTransactionStore(ctrl *gomock.Controller) *MockTransactionStore {
	mock := &MockTransactionStore{ctrl: ctrl}
	mock.recorder = &MockTransactionStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionStore) EXPECT() *MockTransactionStoreMockRecorder {
	return m.recorder
}

// CalcAmountByPeriod mocks base method.
func (m *MockTransactionStore) CalcAmountByPeriod(ctx context.Context, userID int64, moment time.Time, currencyID string) (map[string]decimal.Decimal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalcAmountByPeriod", ctx, userID, moment, currencyID)
	ret0, _ := ret[0].(map[string]decimal.Decimal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalcAmountByPeriod indicates an expected call of CalcAmountByPeriod.
func (mr *MockTransactionStoreMockRecorder) CalcAmountByPeriod(ctx, userID, moment, currencyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalcAmountByPeriod", reflect.TypeOf((*MockTransactionStore)(nil).CalcAmountByPeriod), ctx, userID, moment, currencyID)
}

// MockCurrencyExchanger is a mock of CurrencyExchanger interface.
type MockCurrencyExchanger struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyExchangerMockRecorder
}

// MockCurrencyExchangerMockRecorder is the mock recorder for MockCurrencyExchanger.
type MockCurrencyExchangerMockRecorder struct {
	mock *MockCurrencyExchanger
}

// NewMockCurrencyExchanger creates a new mock instance.
func NewMockCurrencyExchanger(ctrl *gomock.Controller) *MockCurrencyExchanger {
	mock := &MockCurrencyExchanger{ctrl: ctrl}
	mock.recorder = &MockCurrencyExchangerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCurrencyExchanger) EXPECT() *MockCurrencyExchangerMockRecorder {
	return m.recorder
}

// GetMultiplier mocks base method.
func (m *MockCurrencyExchanger) GetMultiplier(ctx context.Context, currency string, date time.Time) (decimal.Decimal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMultiplier", ctx, currency, date)
	ret0, _ := ret[0].(decimal.Decimal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMultiplier indicates an expected call of GetMultiplier.
func (mr *MockCurrencyExchangerMockRecorder) GetMultiplier(ctx, currency, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMultiplier", reflect.TypeOf((*MockCurrencyExchanger)(nil).GetMultiplier), ctx, currency, date)
}
