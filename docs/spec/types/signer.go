package types

import (
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// DomainType is a unique identifier for signatures, 2 identical pieces of data signed with different domains will result in different sigs
type DomainType []byte

var (
	PrimusTestnet = DomainType("primus_testnet")
)

type SignatureType []byte

var (
	QBFTSigType          = []byte{1, 0, 0, 0}
	PostConsensusSigType = []byte{2, 0, 0, 0}
)

type BeaconSigner interface {
	// SignAttestation signs the given attestation
	SignAttestation(data *spec.AttestationData, duty *beacon.Duty, pk []byte) (*spec.Attestation, []byte, error)
	// IsAttestationSlashable returns error if attestation data is slashable
	IsAttestationSlashable(data *spec.AttestationData) error
}

// SSVSigner used for all SSV specific signing
type SSVSigner interface {
	SignRoot(root []byte, sigType SignatureType, pk []byte) (Signature, error)
}

// KeyManager is an interface responsible for all key manager functions
type KeyManager interface {
	BeaconSigner
	SSVSigner
	// AddShare saves a share key
	AddShare(shareKey *bls.SecretKey) error
}

// SSVKeyManager implements the KeyManager interface with all of its funcs
type SSVKeyManager struct {
	keys   map[string]*bls.SecretKey // holds pub keys as key and secret key as value
	domain DomainType
}

func NewSSVKeyManager(domain DomainType) KeyManager {
	return &SSVKeyManager{
		keys:   make(map[string]*bls.SecretKey),
		domain: domain,
	}
}

// SignAttestation signs the given attestation
func (s *SSVKeyManager) SignAttestation(data *spec.AttestationData, duty *beacon.Duty, pk []byte) (*spec.Attestation, []byte, error) {
	panic("implement from beacon ")
}

func (s *SSVKeyManager) IsAttestationSlashable(data *spec.AttestationData) error {
	panic("implement")
}

func (s *SSVKeyManager) SignRoot(root []byte, sigType SignatureType, pk []byte) (Signature, error) {
	if k, found := s.keys[hex.EncodeToString(pk)]; found {
		computedRoot := ComputeSigningRoot(root, ComputeSignatureDomain(s.domain, sigType))
		return k.SignByte(computedRoot).Serialize(), nil
	}
	return nil, errors.New("pk not found")
}

// AddShare saves a share key
func (s *SSVKeyManager) AddShare(sk *bls.SecretKey) error {
	s.keys[sk.GetPublicKey().GetHexString()] = sk
	return nil
}
