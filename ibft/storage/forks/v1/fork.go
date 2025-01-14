package v1

import (
	"github.com/bloxapp/ssv/protocol/v1/message"
)

// ForkV1 is the genesis fork for controller
type ForkV1 struct {
}

// Identifier returns the identifier of the fork
func (f ForkV1) Identifier(pk []byte, role message.RoleType) []byte {
	return message.NewIdentifier(pk, role)
}

// EncodeSignedMsg encodes signed message
func (f ForkV1) EncodeSignedMsg(msg *message.SignedMessage) ([]byte, error) {
	return msg.Encode()
}

// DecodeSignedMsg decodes signed message
func (f ForkV1) DecodeSignedMsg(data []byte) (*message.SignedMessage, error) {
	msgV1 := &message.SignedMessage{}
	err := msgV1.Decode(data)
	if err != nil {
		return nil, err
	}
	return msgV1, nil
}
