package qbft

import (
	"bytes"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/pkg/errors"
)

func uponPrepare(state State, config IConfig, signedPrepare *SignedMessage, prepareMsgContainer, commitMsgContainer *MsgContainer) error {
	// TODO - if we receive a prepare before a proposal and return an error we will never process the prepare msg, we still need to add it to the container
	if state.ProposalAcceptedForCurrentRound == nil {
		return errors.New("not proposal accepted for prepare")
	}

	if err := validSignedPrepareForHeightRoundAndValue(
		state,
		config,
		signedPrepare,
		state.Height,
		state.Round,
		state.ProposalAcceptedForCurrentRound.Message.GetProposalData().GetData(),
		state.Share.Committee,
	); err != nil {
		return errors.Wrap(err, "invalid prepare msg")
	}

	addedMsg, err := prepareMsgContainer.AddIfDoesntExist(signedPrepare)
	if err != nil {
		return errors.Wrap(err, "could not add prepare msg to container")
	}
	if !addedMsg {
		return nil // uponPrepare was already called
	}

	if !state.Share.HasQuorum(len(prepareMsgContainer.MessagesForRound(state.Round))) {
		return nil // no quorum yet
	}

	if didSendCommitForHeightAndRound(state, commitMsgContainer) {
		return nil // already moved to commit stage
	}

	proposedValue := state.ProposalAcceptedForCurrentRound.Message.GetProposalData().GetData()

	state.LastPreparedValue = proposedValue
	state.LastPreparedRound = state.Round

	commitMsg := createCommit(state, proposedValue)
	if err := config.GetNetwork().Broadcast(commitMsg); err != nil {
		return errors.Wrap(err, "failed to broadcast commit message")
	}

	return nil
}

func getRoundChangeJustification(state State, config IConfig, prepareMsgContainer MsgContainer) *SignedMessage {
	if state.LastPreparedValue == nil {
		return nil
	}

	prepareMsgs := prepareMsgContainer.MessagesForRound(state.LastPreparedRound)
	validPrepares := validPreparesForHeightRoundAndDigest(
		state,
		config,
		prepareMsgs,
		state.Height,
		state.LastPreparedRound,
		state.LastPreparedValue,
		state.Share.Committee,
	)
	if state.Share.HasQuorum(len(prepareMsgs)) {
		return validPrepares
	}
	return nil
}

// validPreparesForHeightRoundAndDigest returns an aggregated prepare msg for a specific Height and round
func validPreparesForHeightRoundAndDigest(
	state State,
	config IConfig,
	prepareMessages []*SignedMessage,
	height uint64,
	round Round,
	value []byte,
	operators []*types.Operator) *SignedMessage {
	var aggregatedPrepareMsg *SignedMessage
	for _, signedMsg := range prepareMessages {
		if err := validSignedPrepareForHeightRoundAndValue(state, config, signedMsg, height, round, value, operators); err == nil {
			if aggregatedPrepareMsg == nil {
				aggregatedPrepareMsg = signedMsg
			} else {
				aggregatedPrepareMsg.Aggregate(signedMsg)
			}
		}
	}
	return aggregatedPrepareMsg
}

// validSignedPrepareForHeightRoundAndValue known in dafny spec as validSignedPrepareForHeightRoundAndDigest
// https://entethalliance.github.io/client-spec/qbft_spec.html#dfn-qbftspecification
func validSignedPrepareForHeightRoundAndValue(
	state State,
	config IConfig,
	signedPrepare *SignedMessage,
	height uint64,
	round Round,
	value []byte,
	operators []*types.Operator) error {
	if signedPrepare.Message.MsgType != PrepareMsgType {
		return errors.New("prepare msg type is wrong")
	}
	if signedPrepare.Message.Height != height {
		return errors.New("msg Height wrong")
	}
	if signedPrepare.Message.Round != round {
		return errors.New("msg round wrong")
	}
	if bytes.Compare(signedPrepare.Message.GetPrepareData().GetData(), value) != 0 {
		return errors.New("msg Identifier wrong")
	}

	if len(signedPrepare.GetSigners()) != 1 {
		return errors.New("prepare msg allows 1 signer")
	}

	if err := signedPrepare.Signature.VerifyByOperators(signedPrepare, config.GetSignatureDomainType(), types.QBFTSigType, operators); err != nil {
		return errors.Wrap(err, "prepare msg signature invalid")
	}
	return nil
}

func createPrepare(state State, newRound Round, value []byte) *SignedMessage {
	/**
	Prepare(
	                    signPrepare(
	                        UnsignedPrepare(
	                            |current.blockchain|,
	                            newRound,
	                            digest(m.proposedBlock)),
	                        current.id
	                        )
	                );
	*/
	panic("implement")
}
