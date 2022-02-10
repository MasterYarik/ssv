package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/types"
)

type MessageType int

const (
	ProposalType MessageType = iota
	PrepareType
	CommitType
	RoundChangeType
)

type ProposalData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
	// GetRoundChangeJustification returns a signed message with quorum as justification for the round change
	GetRoundChangeJustification() []*SignedMessage
	// GetPrepareJustification returns a signed message with quorum as justification for a prepared round change
	GetPrepareJustification() []*SignedMessage
}

type PrepareData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
}

type CommitData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
}

type RoundChangeData interface {
	GetPreparedValue() []byte
	GetPreparedRound() Round
	// GetNextProposalData returns NOT nil byte array if the signer is the next round's proposal.
	GetNextProposalData() []byte
	// GetRoundChangeJustification returns signed prepare messages for the last prepared state
	GetRoundChangeJustification() []*SignedMessage
}

type Message struct {
	MsgType    MessageType
	Height     uint64 // QBFT instance height
	Round      Round  // QBFT round for which the msg is for
	Identifier []byte // instance identifier this msg belongs to
	Data       []byte
}

// GetProposalData returns proposal specific data
func (msg *Message) GetProposalData() ProposalData {
	panic("implement")
}

// GetPrepareData returns prepare specific data
func (msg *Message) GetPrepareData() PrepareData {
	panic("implement")
}

// GetCommitData returns commit specific data
func (msg *Message) GetCommitData() PrepareData {
	panic("implement")
}

// GetRoundChangeData returns round change specific data
func (msg *Message) GetRoundChangeData() RoundChangeData {
	panic("implement")
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	panic("implement")
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	panic("implement")
}

// GetRoot returns the root used for signing and verification
func (msg *Message) GetRoot() []byte {
	panic("implement")
}

type SignedMessage struct {
	Signature types.Signature
	Signers   []types.NodeID
	Message   *Message // message for which this signature is for
}

func (signedMsg *SignedMessage) GetSignature() []byte {
	return signedMsg.Signature
}
func (signedMsg *SignedMessage) GetSigners() []types.NodeID {
	return signedMsg.Signers
}

// IsValidSignature returns true if signature is valid (against message and signers)
func (signedMsg *SignedMessage) IsValidSignature(domain types.DomainType, nodes []*types.Node) error {
	pks := make([][]byte, 0)
	for _, id := range signedMsg.Signers {
		for _, n := range nodes {
			if id == n.GetID() {
				pks = append(pks, n.GetPublicKey())
			}
		}
	}

	return signedMsg.Signature.VerifyMultiPubKey(
		signedMsg.Message.GetRoot(),
		domain,
		types.QBFTSigType,
		pks,
	)
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (signedMsg *SignedMessage) MatchedSigners(ids []types.NodeID) bool {
	panic("implement")
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (signedMsg *SignedMessage) Aggregate(sig types.MessageSignature) error {
	panic("implement")
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessage) Encode() ([]byte, error) {
	panic("implement")
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	panic("implement")
}

// GetRoot returns the root used for signing and verification
func (signedMsg *SignedMessage) GetRoot() []byte {
	return signedMsg.Message.GetRoot()
}
