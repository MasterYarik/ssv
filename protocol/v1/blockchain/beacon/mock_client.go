// Code generated by MockGen. DO NOT EDIT.
// Source: ./client.go

// Package beacon is a generated GoMock package.
package beacon

import (
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	phase0 "github.com/attestantio/go-eth2-client/spec/phase0"
	message "github.com/bloxapp/ssv/protocol/v1/message"
	gomock "github.com/golang/mock/gomock"
	bls "github.com/herumi/bls-eth-go-binary/bls"
	reflect "reflect"
)

// MockBeacon is a mock of Beacon interface
type MockBeacon struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconMockRecorder
}

// MockBeaconMockRecorder is the mock recorder for MockBeacon
type MockBeaconMockRecorder struct {
	mock *MockBeacon
}

// NewMockBeacon creates a new mock instance
func NewMockBeacon(ctrl *gomock.Controller) *MockBeacon {
	mock := &MockBeacon{ctrl: ctrl}
	mock.recorder = &MockBeaconMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeacon) EXPECT() *MockBeaconMockRecorder {
	return m.recorder
}

// SignIBFTMessage mocks base method
func (m *MockBeacon) SignIBFTMessage(message *message.ConsensusMessage, pk []byte, forkVersion string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIBFTMessage", message, pk)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIBFTMessage indicates an expected call of SignIBFTMessage
func (mr *MockBeaconMockRecorder) SignIBFTMessage(message, pk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIBFTMessage", reflect.TypeOf((*MockBeacon)(nil).SignIBFTMessage), message, pk)
}

// SignAttestation mocks base method
func (m *MockBeacon) SignAttestation(data *phase0.AttestationData, duty *Duty, pk []byte) (*phase0.Attestation, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignAttestation", data, duty, pk)
	ret0, _ := ret[0].(*phase0.Attestation)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignAttestation indicates an expected call of SignAttestation
func (mr *MockBeaconMockRecorder) SignAttestation(data, duty, pk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignAttestation", reflect.TypeOf((*MockBeacon)(nil).SignAttestation), data, duty, pk)
}

// AddShare mocks base method
func (m *MockBeacon) AddShare(shareKey *bls.SecretKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddShare", shareKey)
	ret0, _ := ret[0].(error)
	return ret0
}

func (m *MockBeacon) RemoveShare(pubKey string) error {
	//TODO implement me
	panic("implement me")
}

// AddShare indicates an expected call of AddShare
func (mr *MockBeaconMockRecorder) AddShare(shareKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddShare", reflect.TypeOf((*MockBeacon)(nil).AddShare), shareKey)
}

// GetDomain mocks base method
func (m *MockBeacon) GetDomain(data *phase0.AttestationData) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomain", data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomain indicates an expected call of GetDomain
func (mr *MockBeaconMockRecorder) GetDomain(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomain", reflect.TypeOf((*MockBeacon)(nil).GetDomain), data)
}

// ComputeSigningRoot mocks base method
func (m *MockBeacon) ComputeSigningRoot(object interface{}, domain []byte) ([32]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComputeSigningRoot", object, domain)
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ComputeSigningRoot indicates an expected call of ComputeSigningRoot
func (mr *MockBeaconMockRecorder) ComputeSigningRoot(object, domain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComputeSigningRoot", reflect.TypeOf((*MockBeacon)(nil).ComputeSigningRoot), object, domain)
}

// GetDuties mocks base method
func (m *MockBeacon) GetDuties(epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*Duty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDuties", epoch, validatorIndices)
	ret0, _ := ret[0].([]*Duty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDuties indicates an expected call of GetDuties
func (mr *MockBeaconMockRecorder) GetDuties(epoch, validatorIndices interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDuties", reflect.TypeOf((*MockBeacon)(nil).GetDuties), epoch, validatorIndices)
}

// GetValidatorData mocks base method
func (m *MockBeacon) GetValidatorData(validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]*v1.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorData", validatorPubKeys)
	ret0, _ := ret[0].(map[phase0.ValidatorIndex]*v1.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorData indicates an expected call of GetValidatorData
func (mr *MockBeaconMockRecorder) GetValidatorData(validatorPubKeys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorData", reflect.TypeOf((*MockBeacon)(nil).GetValidatorData), validatorPubKeys)
}

// GetAttestationData mocks base method
func (m *MockBeacon) GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttestationData", slot, committeeIndex)
	ret0, _ := ret[0].(*phase0.AttestationData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttestationData indicates an expected call of GetAttestationData
func (mr *MockBeaconMockRecorder) GetAttestationData(slot, committeeIndex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttestationData", reflect.TypeOf((*MockBeacon)(nil).GetAttestationData), slot, committeeIndex)
}

// SubmitAttestation mocks base method
func (m *MockBeacon) SubmitAttestation(attestation *phase0.Attestation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitAttestation", attestation)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubmitAttestation indicates an expected call of SubmitAttestation
func (mr *MockBeaconMockRecorder) SubmitAttestation(attestation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitAttestation", reflect.TypeOf((*MockBeacon)(nil).SubmitAttestation), attestation)
}

// SubscribeToCommitteeSubnet mocks base method
func (m *MockBeacon) SubscribeToCommitteeSubnet(subscription []*v1.BeaconCommitteeSubscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToCommitteeSubnet", subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeToCommitteeSubnet indicates an expected call of SubscribeToCommitteeSubnet
func (mr *MockBeaconMockRecorder) SubscribeToCommitteeSubnet(subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToCommitteeSubnet", reflect.TypeOf((*MockBeacon)(nil).SubscribeToCommitteeSubnet), subscription)
}

// MockKeyManager is a mock of KeyManager interface
type MockKeyManager struct {
	ctrl     *gomock.Controller
	recorder *MockKeyManagerMockRecorder
}

// MockKeyManagerMockRecorder is the mock recorder for MockKeyManager
type MockKeyManagerMockRecorder struct {
	mock *MockKeyManager
}

// NewMockKeyManager creates a new mock instance
func NewMockKeyManager(ctrl *gomock.Controller) *MockKeyManager {
	mock := &MockKeyManager{ctrl: ctrl}
	mock.recorder = &MockKeyManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyManager) EXPECT() *MockKeyManagerMockRecorder {
	return m.recorder
}

// SignIBFTMessage mocks base method
func (m *MockKeyManager) SignIBFTMessage(message *message.ConsensusMessage, pk []byte, forkVersion string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIBFTMessage", message, pk)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIBFTMessage indicates an expected call of SignIBFTMessage
func (mr *MockKeyManagerMockRecorder) SignIBFTMessage(message, pk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIBFTMessage", reflect.TypeOf((*MockKeyManager)(nil).SignIBFTMessage), message, pk)
}

// SignAttestation mocks base method
func (m *MockKeyManager) SignAttestation(data *phase0.AttestationData, duty *Duty, pk []byte) (*phase0.Attestation, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignAttestation", data, duty, pk)
	ret0, _ := ret[0].(*phase0.Attestation)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignAttestation indicates an expected call of SignAttestation
func (mr *MockKeyManagerMockRecorder) SignAttestation(data, duty, pk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignAttestation", reflect.TypeOf((*MockKeyManager)(nil).SignAttestation), data, duty, pk)
}

// AddShare mocks base method
func (m *MockKeyManager) AddShare(shareKey *bls.SecretKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddShare", shareKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddShare indicates an expected call of AddShare
func (mr *MockKeyManagerMockRecorder) AddShare(shareKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddShare", reflect.TypeOf((*MockKeyManager)(nil).AddShare), shareKey)
}

// MockSigner is a mock of Signer interface
type MockSigner struct {
	ctrl     *gomock.Controller
	recorder *MockSignerMockRecorder
}

// MockSignerMockRecorder is the mock recorder for MockSigner
type MockSignerMockRecorder struct {
	mock *MockSigner
}

// NewMockSigner creates a new mock instance
func NewMockSigner(ctrl *gomock.Controller) *MockSigner {
	mock := &MockSigner{ctrl: ctrl}
	mock.recorder = &MockSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSigner) EXPECT() *MockSignerMockRecorder {
	return m.recorder
}

// SignIBFTMessage mocks base method
func (m *MockSigner) SignIBFTMessage(message *message.ConsensusMessage, pk []byte, forkVersion string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIBFTMessage", message, pk)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIBFTMessage indicates an expected call of SignIBFTMessage
func (mr *MockSignerMockRecorder) SignIBFTMessage(message, pk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIBFTMessage", reflect.TypeOf((*MockSigner)(nil).SignIBFTMessage), message, pk)
}

// SignAttestation mocks base method
func (m *MockSigner) SignAttestation(data *phase0.AttestationData, duty *Duty, pk []byte) (*phase0.Attestation, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignAttestation", data, duty, pk)
	ret0, _ := ret[0].(*phase0.Attestation)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignAttestation indicates an expected call of SignAttestation
func (mr *MockSignerMockRecorder) SignAttestation(data, duty, pk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignAttestation", reflect.TypeOf((*MockSigner)(nil).SignAttestation), data, duty, pk)
}

// MockSigningUtil is a mock of SigningUtil interface
type MockSigningUtil struct {
	ctrl     *gomock.Controller
	recorder *MockSigningUtilMockRecorder
}

// MockSigningUtilMockRecorder is the mock recorder for MockSigningUtil
type MockSigningUtilMockRecorder struct {
	mock *MockSigningUtil
}

// NewMockSigningUtil creates a new mock instance
func NewMockSigningUtil(ctrl *gomock.Controller) *MockSigningUtil {
	mock := &MockSigningUtil{ctrl: ctrl}
	mock.recorder = &MockSigningUtilMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSigningUtil) EXPECT() *MockSigningUtilMockRecorder {
	return m.recorder
}

// GetDomain mocks base method
func (m *MockSigningUtil) GetDomain(data *phase0.AttestationData) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomain", data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomain indicates an expected call of GetDomain
func (mr *MockSigningUtilMockRecorder) GetDomain(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomain", reflect.TypeOf((*MockSigningUtil)(nil).GetDomain), data)
}

// ComputeSigningRoot mocks base method
func (m *MockSigningUtil) ComputeSigningRoot(object interface{}, domain []byte) ([32]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComputeSigningRoot", object, domain)
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ComputeSigningRoot indicates an expected call of ComputeSigningRoot
func (mr *MockSigningUtilMockRecorder) ComputeSigningRoot(object, domain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComputeSigningRoot", reflect.TypeOf((*MockSigningUtil)(nil).ComputeSigningRoot), object, domain)
}
