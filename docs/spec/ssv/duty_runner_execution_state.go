package ssv

import (
	"crypto/sha256"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

// DutyExecutionState holds all the relevant progress the duty execution progress
type DutyExecutionState struct {
	RunningInstance *qbft.Instance

	DecidedValue *types.ConsensusData

	SignedAttestation *spec.Attestation
	SignedProposal    *spec.SignedBeaconBlock

	PostConsensusSignatures map[types.OperatorID][]byte
	PostConsensusSigRoot    []byte
	// Quorum is the number of min signatures needed for quorum
	Quorum uint64

	Finished bool
}

func NewDutyExecutionState(quorum uint64) *DutyExecutionState {
	return &DutyExecutionState{
		Quorum:                  quorum,
		PostConsensusSignatures: make(map[types.OperatorID][]byte),
		Finished:                false,
	}
}

func (pcs *DutyExecutionState) AddPostConsensusPartialSig(sigMsg *PartialSignatureMessage) error {
	if len(sigMsg.Signers) != 1 {
		return errors.New("PartialSignatureMessage has != 1 Signers")
	}

	if pcs.PostConsensusSignatures[sigMsg.Signers[0]] == nil {
		pcs.PostConsensusSignatures[sigMsg.Signers[0]] = sigMsg.PartialSignature
	}
	return nil
}

// ReconstructAttestationSig aggregates collected partial sigs, reconstructs a valid sig and returns an attestation obj with the reconstructed sig
func (pcs *DutyExecutionState) ReconstructAttestationSig(validatorPubKey []byte) (*spec.Attestation, error) {
	// Reconstruct signatures
	signature, err := types.ReconstructSignatures(pcs.PostConsensusSignatures)
	if err != nil {
		return nil, errors.Wrap(err, "failed to reconstruct signatures")
	}
	if err := types.VerifyReconstructedSignature(signature, validatorPubKey, pcs.PostConsensusSigRoot); err != nil {
		return nil, errors.Wrap(err, "failed to verify reconstruct signature")
	}

	blsSig := spec.BLSSignature{}
	copy(blsSig[:], signature.Serialize())
	pcs.SignedAttestation.Signature = blsSig
	return pcs.SignedAttestation, nil
}

func (pcs *DutyExecutionState) HasPostConsensusSigQuorum() bool {
	return uint64(len(pcs.PostConsensusSignatures)) >= pcs.Quorum
}

// SetFinished will mark this execution state as finished
func (pcs *DutyExecutionState) SetFinished() {
	pcs.Finished = true
}

// GetRoot returns the root used for signing and verification
func (pcs *DutyExecutionState) GetRoot() ([]byte, error) {
	marshaledRoot, err := pcs.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyExecutionState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// IsFinished returns true if this execution state is finished
func (pcs *DutyExecutionState) IsFinished() bool {
	return pcs.Finished
}

// Encode returns the encoded struct in bytes or error
func (pcs *DutyExecutionState) Encode() ([]byte, error) {
	return json.Marshal(pcs)
}

// Decode returns error if decoding failed
func (pcs *DutyExecutionState) Decode(data []byte) error {
	return json.Unmarshal(data, &pcs)
}
