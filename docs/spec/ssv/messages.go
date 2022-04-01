package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type PostConsensusMessage struct {
	Height          qbft.Height
	DutySignature   []byte // The beacon chain partial Signature for a duty
	DutySigningRoot []byte // the root signed in DutySignature
	Signers         []types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (pcsm *PostConsensusMessage) Encode() ([]byte, error) {
	return json.Marshal(pcsm)
}

// Decode returns error if decoding failed
func (pcsm *PostConsensusMessage) Decode(data []byte) error {
	return json.Unmarshal(data, pcsm)
}

func (pcsm *PostConsensusMessage) GetRoot() ([]byte, error) {
	marshaledRoot, err := pcsm.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PostConsensusMessage")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

func (pcsm *PostConsensusMessage) Validate() error {
	if len(pcsm.DutySignature) != 96 {
		return errors.New("PostConsensusMessage sig invalid")
	}
	if len(pcsm.DutySigningRoot) != 32 {
		return errors.New("DutySigningRoot invalid")
	}
	if len(pcsm.Signers) != 1 {
		return errors.New("invalid PostConsensusMessage signers")
	}
	return nil
}

type SignedPostConsensusMessage struct {
	Message   *PostConsensusMessage
	Signature types.Signature
	Signers   []types.OperatorID
}

// Encode returns a msg encoded bytes or error
func (spcsm *SignedPostConsensusMessage) Encode() ([]byte, error) {
	return json.Marshal(spcsm)
}

// Decode returns error if decoding failed
func (spcsm *SignedPostConsensusMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &spcsm)
}

func (spcsm *SignedPostConsensusMessage) GetSignature() types.Signature {
	return spcsm.Signature
}

func (spcsm *SignedPostConsensusMessage) GetSigners() []types.OperatorID {
	return spcsm.Signers
}

func (spcsm *SignedPostConsensusMessage) GetRoot() ([]byte, error) {
	return spcsm.Message.GetRoot()
}

func (spcsm *SignedPostConsensusMessage) Aggregate(signedMsg types.MessageSignature) error {
	//if !bytes.Equal(spcsm.GetRoot(), signedMsg.GetRoot()) {
	//	return errors.New("can't aggregate msgs with different roots")
	//}
	//
	//// verify no matching Signers
	//for _, signerID := range spcsm.Signers {
	//	for _, toMatchID := range signedMsg.GetSigners() {
	//		if signerID == toMatchID {
	//			return errors.New("signer IDs partially/ fully match")
	//		}
	//	}
	//}
	//
	//allSigners := append(spcsm.Signers, signedMsg.GetSigners()...)
	//
	//// verify and aggregate
	//sig1, err := blsSig(spcsm.Signature)
	//if err != nil {
	//	return errors.Wrap(err, "could not parse DutySignature")
	//}
	//
	//sig2, err := blsSig(signedMsg.GetSignature())
	//if err != nil {
	//	return errors.Wrap(err, "could not parse DutySignature")
	//}
	//
	//sig1.Add(sig2)
	//spcsm.Signature = sig1.Serialize()
	//spcsm.Signers = allSigners
	//return nil
	panic("implement")
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (spcsm *SignedPostConsensusMessage) MatchedSigners(ids []types.OperatorID) bool {
	toMatchCnt := make(map[types.OperatorID]int)
	for _, id := range ids {
		toMatchCnt[id]++
	}

	foundCnt := make(map[types.OperatorID]int)
	for _, id := range spcsm.GetSigners() {
		foundCnt[id]++
	}

	for id, cnt := range toMatchCnt {
		if cnt != foundCnt[id] {
			return false
		}
	}
	return true
}

func blsSig(sig []byte) (*bls.Sign, error) {
	ret := &bls.Sign{}
	if err := ret.Deserialize(sig); err != nil {
		return nil, errors.Wrap(err, "could not covert DutySignature byts to bls.sign")
	}
	return ret, nil
}

func (spcsm *SignedPostConsensusMessage) Validate() error {
	if len(spcsm.Signature) != 96 {
		return errors.New("SignedPostConsensusMessage sig invalid")
	}
	if len(spcsm.Signers) != 1 {
		return errors.New("no SignedPostConsensusMessage signers")
	}
	return spcsm.Message.Validate()
}
