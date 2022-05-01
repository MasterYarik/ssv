package valcheck

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongDutyPubKey tests decided value with duty validator pubkey != the duty runner's pubkey
func WrongDutyPubKey() *tests.SpecTest {
	dr := testingutils.AttesterRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes(testingutils.TestConsensusWrongDutyPKDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusWrongDutyPKDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusWrongDutyPKDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusWrongDutyPKDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestConsensusWrongDutyPKDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestConsensusWrongDutyPKDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestConsensusWrongDutyPKDataByts),
		}), nil),
	}

	return &tests.SpecTest{
		Name:                    "wrong decided value's pubkey",
		Runner:                  dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "3f82ba9763ce97791e62f6daff599692f82608dbff222e8f6562a48a34f08272",
		ExpectedError:           "decided value is invalid: decided value's validator pk is wrong",
	}
}
