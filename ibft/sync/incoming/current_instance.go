package incoming

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/network"
	"github.com/bloxapp/ssv/storage/kv"
	"go.uber.org/zap"
)

func (s *ReqHandler) handleGetCurrentInstanceReq(msg *network.SyncChanObj) {
	retMsg := &network.SyncMessage{
		Lambda: s.identifier,
		Type:   network.Sync_GetCurrentInstance,
	}

	if s.currentInstanceMsg != nil {
		retMsg.SignedMessages = []*proto.SignedMessage{s.currentInstanceMsg}
	} else {
		retMsg.Error = kv.EntryNotFoundError
	}

	if err := s.network.RespondToGetCurrentInstance(msg.Stream, retMsg); err != nil {
		s.logger.Error("failed to send current instance req", zap.Error(err))
	}
}