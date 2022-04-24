package attestations

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// DutyCommitteeIndexNotMatchingAttestations tests that a duty committee index == attestation committee index
func DutyCommitteeIndexNotMatchingAttestations() *tests.SpecTest {
	dr := testingutils.BaseRunner()

	consensusData := &types.ConsensusData{
		Duty: &beacon.Duty{
			Type:                    beacon.RoleTypeAttester,
			PubKey:                  testingutils.TestingValidatorPubKey,
			Slot:                    12,
			ValidatorIndex:          1,
			CommitteeIndex:          5,
			CommitteesAtSlot:        36,
			CommitteeLength:         128,
			ValidatorCommitteeIndex: 11,
		},
		AttestationData: testingutils.TestingAttestationData,
	}
	startingValue, _ := consensusData.Encode()

	// the starting value is not the same as the actual proposal!
	if err := dr.StartNewInstance(testingutils.TestAttesterConsensusDataByts); err != nil {
		panic(err.Error())
	}

	msgs := []*types.SSVMessage{
		testingutils.SSVMsg(testingutils.SignQBFTMsg(testingutils.TestingSK1, 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes(startingValue, nil, nil),
		}), nil),
	}

	return &tests.SpecTest{
		Name:                    "duty committee index matches attestation committee index",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "039e927a1858548cc411afe6442ee7222d285661058621bdb1e02405c0f344d4",
		ExpectedError:           "failed to process consensus msg: could not process msg: proposal invalid: proposal not justified: proposal value invalid: attestation data CommitteeIndex != duty CommitteeIndex",
	}
}