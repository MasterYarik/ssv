package consensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// UnknownDuty tests decided value with an unknown duty type
func UnknownDuty() *tests.SpecTest {
	dr := testingutils.BaseRunner()
	if err := dr.StartNewInstance([]byte{1, 2, 3, 4}); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK2, 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts),
		}), nil),
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK3, 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes(testingutils.TestConsensusUnkownDutyTypeDataByts),
		}), nil),
	}

	return &tests.SpecTest{
		Name:                    "unknown decided value's duty type",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "a56b370679d010d217e2c87873e6820f011d9ace64fbfed91c2b185da3a79bba",
		ExpectedError:           "decided value is invalid: decided value's duty has wrong beacon role type",
	}
}
