package controller

import (
	"github.com/bloxapp/ssv/protocol/v1/message"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ProcessChangeRound check basic pipeline validation than check if height or round is higher than the last one. if so, update
func (c *Controller) ProcessChangeRound(msg *message.SignedMessage) error {
	if err := c.ValidateChangeRoundMsg(msg); err != nil {
		return err
	}
	lastMsg, err := c.changeRoundStorage.GetLastChangeRoundMsg(c.Identifier)
	if err != nil {
		return errors.Wrap(err, "failed to get last change round msg")
	}

	if lastMsg == nil {
		// no last changeRound msg exist, save the first one
		c.logger.Debug("no last change round exist. saving first one", zap.Int64("NewHeight", int64(msg.Message.Height)), zap.Int64("NewRound", int64(msg.Message.Round)))
		return c.changeRoundStorage.SaveLastChangeRoundMsg(msg)
	}

	logger := c.logger.With(
		zap.Int64("lastHeight", int64(lastMsg.Message.Height)),
		zap.Int64("NewHeight", int64(msg.Message.Height)),
		zap.Int64("lastRound", int64(lastMsg.Message.Round)),
		zap.Int64("NewRound", int64(msg.Message.Round)))

	if msg.Message.Height < lastMsg.Message.Height {
		// height is lower than the last known
		logger.Debug("new changeRoundMsg height is lower than last changeRoundMsg")
		return nil
	} else if msg.Message.Height == lastMsg.Message.Height {
		if msg.Message.Round <= lastMsg.Message.Round {
			// round is not higher than last known
			logger.Debug("new changeRoundMsg round is lower than last changeRoundMsg")
			return nil
		}
	}

	// new msg is higher than last one, save.
	logger.Debug("last change round updated")
	return c.changeRoundStorage.SaveLastChangeRoundMsg(msg)
}

// ValidateChangeRoundMsg - validation for read mode change round msg
// validating -
// basic validation, signature, changeRound data
func (c *Controller) ValidateChangeRoundMsg(msg *message.SignedMessage) error {
	return c.fork.ValidateChangeRoundMsg(c.ValidatorShare, c.Identifier).Run(msg)
}
