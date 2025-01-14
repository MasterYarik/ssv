// Code generated by MockGen. DO NOT EDIT.
// Source: ./controller.go

// Package mocks is a generated GoMock package.
package mocks

import (
	phase0 "github.com/attestantio/go-eth2-client/spec/phase0"
	eth1 "github.com/bloxapp/ssv/eth1"
	validator "github.com/bloxapp/ssv/operator/validator"
	message "github.com/bloxapp/ssv/protocol/v1/blockchain/beacon"
	validator0 "github.com/bloxapp/ssv/protocol/v1/validator"
	gomock "github.com/golang/mock/gomock"
	event "github.com/prysmaticlabs/prysm/async/event"
	reflect "reflect"
)

// MockController is a mock of Controller interface
type MockController struct {
	ctrl     *gomock.Controller
	recorder *MockControllerMockRecorder
}

// MockControllerMockRecorder is the mock recorder for MockController
type MockControllerMockRecorder struct {
	mock *MockController
}

// NewMockController creates a new mock instance
func NewMockController(ctrl *gomock.Controller) *MockController {
	mock := &MockController{ctrl: ctrl}
	mock.recorder = &MockControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockController) EXPECT() *MockControllerMockRecorder {
	return m.recorder
}

// ListenToEth1Events mocks base method
func (m *MockController) ListenToEth1Events(feed *event.Feed) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ListenToEth1Events", feed)
}

// ListenToEth1Events indicates an expected call of ListenToEth1Events
func (mr *MockControllerMockRecorder) ListenToEth1Events(feed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenToEth1Events", reflect.TypeOf((*MockController)(nil).ListenToEth1Events), feed)
}

// StartValidators mocks base method
func (m *MockController) StartValidators() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartValidators")
}

// StartValidators indicates an expected call of StartValidators
func (mr *MockControllerMockRecorder) StartValidators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartValidators", reflect.TypeOf((*MockController)(nil).StartValidators))
}

// GetValidatorsIndices mocks base method
func (m *MockController) GetValidatorsIndices() []phase0.ValidatorIndex {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorsIndices")
	ret0, _ := ret[0].([]phase0.ValidatorIndex)
	return ret0
}

// GetValidatorsIndices indicates an expected call of GetValidatorsIndices
func (mr *MockControllerMockRecorder) GetValidatorsIndices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorsIndices", reflect.TypeOf((*MockController)(nil).GetValidatorsIndices))
}

// GetValidator mocks base method
func (m *MockController) GetValidator(pubKey string) (validator0.IValidator, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", pubKey)
	ret0, _ := ret[0].(validator0.IValidator)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator
func (mr *MockControllerMockRecorder) GetValidator(pubKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockController)(nil).GetValidator), pubKey)
}

// UpdateValidatorMetaDataLoop mocks base method
func (m *MockController) UpdateValidatorMetaDataLoop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateValidatorMetaDataLoop")
}

// UpdateValidatorMetaDataLoop indicates an expected call of UpdateValidatorMetaDataLoop
func (mr *MockControllerMockRecorder) UpdateValidatorMetaDataLoop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateValidatorMetaDataLoop", reflect.TypeOf((*MockController)(nil).UpdateValidatorMetaDataLoop))
}

// StartNetworkMediators mocks base method
func (m *MockController) StartNetworkMediators() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartNetworkHandlers")
}

// StartNetworkMediators indicates an expected call of StartNetworkMediators
func (mr *MockControllerMockRecorder) StartNetworkMediators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartNetworkHandlers", reflect.TypeOf((*MockController)(nil).StartNetworkMediators))
}

// Eth1EventHandler mocks base method
func (m *MockController) Eth1EventHandler(handlers ...validator.ShareEventHandlerFunc) eth1.SyncEventHandler {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range handlers {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Eth1EventHandler", varargs...)
	ret0, _ := ret[0].(eth1.SyncEventHandler)
	return ret0
}

// Eth1EventHandler indicates an expected call of Eth1EventHandler
func (mr *MockControllerMockRecorder) Eth1EventHandler(handlers ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Eth1EventHandler", reflect.TypeOf((*MockController)(nil).Eth1EventHandler), handlers...)
}

// GetAllValidatorShares mocks base method
func (m *MockController) GetAllValidatorShares() ([]*message.Share, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllValidatorShares")
	ret0, _ := ret[0].([]*message.Share)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllValidatorShares indicates an expected call of GetAllValidatorShares
func (mr *MockControllerMockRecorder) GetAllValidatorShares() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllValidatorShares", reflect.TypeOf((*MockController)(nil).GetAllValidatorShares))
}
