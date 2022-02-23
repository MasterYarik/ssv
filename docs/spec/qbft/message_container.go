package qbft

type MsgContainer interface {
	// AllMessagedForHeight returns all messages for the specific height
	AllMessagedForHeight(height uint64) []*SignedMessage
	// MessagesForHeightAndRound returns all msgs for height and round, empty slice otherwise
	MessagesForHeightAndRound(height uint64, round Round) []*SignedMessage
	// AddIfDoesntExist will add a msg with the following unique params: (round, signer)
	// Will prevent byzantine signers to send multiple msgs with different values for the same round
	AddIfDoesntExist(msg *SignedMessage) bool
}