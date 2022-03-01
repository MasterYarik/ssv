package qbft

import (
	"github.com/bloxapp/ssv/docs/spec/utils"
	"github.com/bloxapp/ssv/utils/threadsafe"
	"github.com/pkg/errors"
	"sync"
)

type proposedValueCheck func(data []byte) error

type IInstance interface {
	// Start will start the new instance with a specific value and height
	Start(value []byte, height uint64)
	// ProcessMsg will process a signed msg
	ProcessMsg(msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error)
	// IsDecided will return true and a non-empty byte slice if instance decided.
	IsDecided() (bool, []byte)
	// GetHeight returns the instance's height
	GetHeight() uint64
}

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	state      State
	config     Config
	valueCheck proposedValueCheck

	proposeContainer     MsgContainer
	prepareContainer     MsgContainer
	commitContainer      MsgContainer
	roundChangeContainer MsgContainer

	decided      *threadsafe.SafeBool
	decidedValue *threadsafe.SafeBytes
	processMsgF  *utils.ThreadSafeF
	startOnce    sync.Once
	startValue   []byte
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height uint64) {
	i.startOnce.Do(func() {
		i.startValue = value
		i.state.Round = FirstRound
		i.state.Height = height

		// propose if this node is the proposer
		if proposer(i.state, FirstRound) == i.config.GetID() {
			proposal, err := createProposal(i.state, i.config, i.startValue, nil, nil)
			if err != nil {
				// TODO log
			}
			if err := i.config.GetNetwork().Broadcast(proposal); err != nil {
				// TODO - log
			}
		}
	})
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {
	res := i.processMsgF.Run(func() interface{} {
		switch msg.Message.MsgType {
		case ProposalMsgType:
			return uponProposal(i.state, i.config, msg, i.proposeContainer)
		case PrepareMsgType:
			return uponPrepare(i.state, i.config, msg, i.prepareContainer, i.commitContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = uponCommit(i.state, i.config, msg, i.commitContainer)
			i.decided.Set(decided)
			if decided {
				i.decidedValue.Set(decidedValue)
			}

			// TODO - Roberto comment: we should send a decided msg here
			return err
		case RoundChangeMsgType:
			return uponRoundChange(i.state, i.config, msg, i.roundChangeContainer, i.valueCheck)
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.decided.Get(), i.decidedValue.Get(), aggregatedCommit, nil
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	return i.decided.Get(), i.decidedValue.Get()
}

// GetHeight interface implementation
func (i *Instance) GetHeight() uint64 {
	return i.state.Height
}
