package types

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/beacon"
)

type MessageID []byte

func MessageIDForValidatorPKAndRole(pk []byte, role beacon.RoleType) MessageID {
	panic("implement")
}

func (msgID MessageID) String() string {
	return hex.EncodeToString(msgID)
}

type Type uint64

const (
	// Consensus are all QBFT consensus related messages
	Consensus Type = iota
	// Sync are all QBFT sync messages
	Sync
	// PostConsensusSignature are all partial signatures sent after consensus
	PostConsensusSignature
)

type MessageEncoder interface {
	// Encode returns a msg encoded bytes or error
	Encode() ([]byte, error)
	// Decode returns error if decoding failed
	Decode(data []byte) error
}

type MessageDigest interface {
	// Digest returns a digest of the msg including any msg specific data, used for signing and verification
	Digest() []byte
}

type MessageSignature interface {
	GetSignature() []byte
	GetSigners() []NodeID
	// IsValidSignature returns true if signature is valid (against message and signers)
	IsValidSignature(nodes []Node) bool
	// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
	MatchedSigners(ids []NodeID) bool
	// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
	Aggregate(signedMsg MessageSignature) error
}

// SSVMessage is the main message passed within the SSV network, it can contain different types of messages (QBTF, Sync, etc.)
type SSVMessage interface {
	MessageEncoder
	MessageDigest

	GetType() Type
	// GetID returns a unique msg ID that is used to identify to which validator should the message be sent for processing
	GetID() MessageID
	// GetData returns message data as byte slice
	GetData() []byte
}
