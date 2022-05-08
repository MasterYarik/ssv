package messages

import (
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
)

// MultipleMessageSigners tests >1 PostConsensusMessage Signers
func MultipleMessageSigners() *tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.DecidedRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsgWithMsgMultiSigners(ks.Shares[1], 1, qbft.FirstHeight)),
	}

	return &tests.SpecTest{
		Name:                    ">1 PostConsensusMessage Signers",
		Runner:                  dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "cbcefe579470d914c3c230bd45cee06e9c5723460044b278a0c629a742551b02",
		ExpectedError:           "partial post valcheck sig invalid: SignedPartialSignatureMessage invalid: invalid PartialSignatureMessage signers",
	}
}
