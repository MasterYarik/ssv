package testingutils

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var TestingConfig = &qbft.Config{
	Signer:     NewTestingKeyManager(),
	SigningPK:  TestingSK1.GetPublicKey().Serialize(),
	Domain:     types.PrimusTestnet,
	ValueCheck: ssv.BeaconAttestationValueCheck(ssv.NowTestNetwork),
	Storage:    NewTestingStorage(),
	Network:    NewTestingNetwork(),
}

var TestingShare = &types.Share{
	OperatorID:      1,
	ValidatorPubKey: TestingValidatorPubKey[:],
	SharePubKey:     TestingSK1.GetPublicKey().Serialize(),
	DomainType:      types.PrimusTestnet,
	Quorum:          3,
	PartialQuorum:   2,
	Committee: []*types.Operator{
		{
			OperatorID: 1,
			PubKey:     TestingSK1.GetPublicKey().Serialize(),
		},
		{
			OperatorID: 2,
			PubKey:     TestingSK2.GetPublicKey().Serialize(),
		},
		{
			OperatorID: 3,
			PubKey:     TestingSK3.GetPublicKey().Serialize(),
		},
		{
			OperatorID: 4,
			PubKey:     TestingSK4.GetPublicKey().Serialize(),
		},
	},
}
var BaseInstance = func() *qbft.Instance {
	ret := qbft.NewInstance(TestingConfig, nil, nil)
	ret.State = &qbft.State{
		Share:                           TestingShare,
		ID:                              []byte{1, 2, 3, 4},
		Round:                           qbft.FirstRound,
		Height:                          qbft.FirstHeight,
		LastPreparedRound:               qbft.NoRound,
		LastPreparedValue:               nil,
		ProposalAcceptedForCurrentRound: nil,
	}
	ret.State.ProposeContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.State.PrepareContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.State.CommitContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	ret.State.RoundChangeContainer = &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	return ret
}

func NewTestingQBFTController(identifier []byte) *qbft.Controller {
	ret := qbft.NewController(
		[]byte{1, 2, 3, 4},
		TestingShare,
		types.PrimusTestnet,
		NewTestingKeyManager(),
		ssv.BeaconAttestationValueCheck(ssv.NowTestNetwork),
		NewTestingStorage(),
		NewTestingNetwork(),
	)
	ret.Identifier = identifier
	ret.Domain = types.PrimusTestnet
	return ret
}