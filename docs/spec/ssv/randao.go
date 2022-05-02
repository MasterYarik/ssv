package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func (dr *Runner) SignRandaoPreConsensus(epoch spec.Epoch, slot spec.Slot, signer types.KeyManager) (*PartialSignatureMessage, error) {
	sig, r, err := signer.SignRandaoReveal(epoch, dr.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign partial randao reveal")
	}

	dr.State.RandaoPartialSig.SigRoot = ensureRoot(r)

	// generate partial sig for randao
	msg := &PartialSignatureMessage{
		Type:             RandaoPartialSig,
		Slot:             slot,
		PartialSignature: sig,
		SigningRoot:      ensureRoot(r),
		Signers:          []types.OperatorID{dr.Share.OperatorID},
	}

	return msg, nil
}

// ProcessRandaoMessage process randao msg, returns true if it has quorum for partial signatures.
// returns true only once (first time quorum achieved)
func (dr *Runner) ProcessRandaoMessage(msg *SignedPartialSignatureMessage) (bool, error) {
	if err := dr.canProcessRandaoMsg(msg); err != nil {
		return false, errors.Wrap(err, "can't process randao message")
	}

	prevQuorum := dr.State.RandaoPartialSig.HasQuorum()

	if err := dr.State.RandaoPartialSig.AddSignature(msg.Message); err != nil {
		return false, errors.Wrap(err, "could not add partial randao signature")
	}

	if prevQuorum {
		return false, nil
	}

	return dr.State.RandaoPartialSig.HasQuorum(), nil
}

// canProcessRandaoMsg returns true if it can process randao message, false if not
func (dr *Runner) canProcessRandaoMsg(msg *SignedPartialSignatureMessage) error {
	if err := dr.validatePartialSigMsg(msg, dr.State.RandaoPartialSig, dr.CurrentDuty.Slot); err != nil {
		return errors.Wrap(err, "randao msg invalid")
	}

	if dr.randaoSigTimeout(dr.BeaconNetwork.EstimatedCurrentSlot()) {
		return errors.New("randao sig collection timeout")
	}

	return nil
}

// randaoSigTimeout returns true if collecting randao sigs timed out
func (dr *Runner) randaoSigTimeout(currentSlot spec.Slot) bool {
	return dr.partialSigCollectionTimeout(dr.State.RandaoPartialSig, currentSlot)
}