package signedmsg

import (
	"errors"
	"github.com/bloxapp/ssv/protocol/v1/message"
	"github.com/bloxapp/ssv/protocol/v1/qbft/pipelines"
)

// MsgTypeCheck is the pipeline to check message type
func MsgTypeCheck(msgType message.ConsensusMessageType) pipelines.SignedMessagePipeline {
	return pipelines.WrapFunc("type check", func(signedMessage *message.SignedMessage) error {
		if signedMessage.Message.MsgType != msgType {
			return errors.New("message type is wrong")
		}
		return nil
	})
}
