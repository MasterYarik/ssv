package testingutils

import (
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
)

var BaseValidator = func(keySet *TestKeySet) *ssv.Validator {
	ret := ssv.NewValidator(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		NewTestingStorage(),
		testShare(keySet),
		NewTestingKeyManager(),
	)
	ret.DutyRunners[types.BNRoleAttester] = AttesterRunner(keySet)
	ret.DutyRunners[types.BNRoleProposer] = ProposerRunner(keySet)
	ret.DutyRunners[types.BNRoleAggregator] = AggregatorRunner(keySet)
	ret.DutyRunners[types.BNRoleSyncCommittee] = SyncCommitteeRunner(keySet)
	ret.DutyRunners[types.BNRoleSyncCommitteeContribution] = SyncCommitteeContributionRunner(keySet)
	return ret
}
