package postconsensus

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// WrongBeaconChainSig tests an invalid PostConsensusMessage PartialSignature
func WrongBeaconChainSig() *tests.SpecTest {
	dr := testingutils.DecidedRunner()

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsgWithWrongSig(testingutils.TestingSK1, 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    "Invalid PostConsensusMessage PartialSignature",
		DutyRunner:              dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "cbcefe579470d914c3c230bd45cee06e9c5723460044b278a0c629a742551b02",
		ExpectedError:           "partial post consensus sig invalid: could not verify beacon partial Signature: could not verify Signature from iBFT member 1",
	}
}
