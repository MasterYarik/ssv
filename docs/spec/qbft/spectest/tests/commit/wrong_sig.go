package commit

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongSignature tests a single commit received with a wrong signature
func WrongSignature() *tests.SpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		}),
		testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(testingutils.TestingSK2, types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(testingutils.TestingSK3, types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(testingutils.TestingSK1, types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		}),
	}
	return &tests.SpecTest{
		Name:          "wrong commit signature",
		Pre:           pre,
		PostRoot:      "a35c30905e199f647a8702210fd04a69a6049527cae31e256b0b5b75395a5750",
		Messages:      msgs,
		ExpectedError: "commit msg invalid: commit msg signature invalid: failed to verify signature",
	}
}