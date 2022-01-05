package qbft

type MessageType int

const (
	ProposalType MessageType = iota
	PrepareType
	CommitType
	RoundChangeType
)

type MessageEncoder interface {
	// Encode returns a msg encoded bytes or error
	Encode() ([]byte, error)
	// Decode returns error if decoding failed
	Decode(data []byte) error
	// Digest returns a digest of the msg including any msg specific data, used for signing and verification
	Digest() []byte
}

type Message interface {
	MessageEncoder

	// GetType returns the msg type
	GetType() MessageType
	// GetHeight returns the QBFT instance height (similar to block number)
	GetHeight() uint64
	// GetRound returns the round for which the msg was sent
	GetRound() Round
	// GetInstanceIdentifier returns the instance identifier this msg belongs to
	GetInstanceIdentifier() []byte

	// GetProposalData returns proposal specific data
	GetProposalData() ProposalData
	// GetPrepareData returns prepare specific data
	GetPrepareData() PrepareData
	// GetCommitData returns commit specific data
	GetCommitData() PrepareData
	// GetRoundChangeData returns round change specific data
	GetRoundChangeData() RoundChangeData
}

type SignedMessage interface {
	MessageEncoder

	// GetMessage returns the message for which this signature is for
	GetMessage() Message
	// GetSignature returns the signature
	GetSignature() []byte
	// GetSignerIds returns the ids of signers (according to node configuration)
	GetSignerIds() []NodeID
	// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
	MatchedSigners(ids []NodeID) bool
	// IsValidSignature returns true if signature is valid (against message and signers)
	IsValidSignature(nodes []Node) bool
	// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
	Aggregate(signedMsg SignedMessage) error
}

type ProposalData interface {
	// GetData returns the data for which this QBFT instance tries to decide, can be any arbitrary data
	GetData() []byte
	// GetRoundChangeJustification returns a signed message with quorum as justification for the round change
	GetRoundChangeJustification() []SignedMessage
	// GetPrepareJustification returns a signed message with quorum as justification for a prepared round change
	GetPrepareJustification() []SignedMessage
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
	// Should NOT be included in digest
	GetNextProposalData() []byte
	// GetRoundChangeJustification returns an aggregated signed prepare message for the last prepared state
	// Should NOT be included in digest
	GetRoundChangeJustification() SignedMessage
}

// ValueCheck is an interface which validates the pre-prepare value passed to the node.
// It's kept minimal to allow the implementation to have all the check logic.
type ValueCheck interface {
	Check(value []byte) error
}
