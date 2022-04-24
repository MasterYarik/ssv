package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/ssv"
)

var BaseValidator = func() *ssv.Validator {
	ret := ssv.NewValidator(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		NewTestingStorage(),
		TestingShare,
		NewTestingKeyManager(),
		ssv.BeaconAttestationValueCheck(ssv.NowTestNetwork),
	)
	ret.DutyRunners[beacon.RoleTypeAttester] = BaseRunner()
	return ret
}