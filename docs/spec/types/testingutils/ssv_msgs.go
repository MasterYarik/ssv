package testingutils

import (
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	"github.com/bloxapp/ssv/docs/spec/ssv"
	"github.com/bloxapp/ssv/docs/spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestAttesterConsensusData = &types.ConsensusData{
	Duty:            TestingAttesterDuty,
	AttestationData: TestingAttestationData,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

var TestAggregatorConsensusData = &types.ConsensusData{
	Duty:            TestingAggregatorDuty,
	AttestationData: TestingAttestationData,
}
var TestAggregatorConsensusDataByts, _ = TestAggregatorConsensusData.Encode()

var TestConsensusUnkownDutyTypeData = &types.ConsensusData{
	Duty:            TestingUnknownDutyType,
	AttestationData: TestingAttestationData,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var SSVMsg = func(qbftMsg *qbft.SignedMessage, postMsg *ssv.SignedPostConsensusMessage) *types.SSVMessage {
	var msgType types.MsgType
	var data []byte
	if qbftMsg != nil {
		msgType = types.SSVConsensusMsgType
		data, _ = qbftMsg.Encode()
	} else if postMsg != nil {
		msgType = types.SSVPostConsensusMsgType
		data, _ = postMsg.Encode()
	} else {
		panic("msg type undefined")
	}

	msgID := types.MessageIDForValidatorPKAndRole(TestingValidatorPubKey[:], beacon.RoleTypeAttester)

	return &types.SSVMessage{
		MsgType: msgType,
		MsgID:   msgID,
		Data:    data,
	}
}

var PostConsensusAttestationMsgWithMsgMultiSigners = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, false, false, true, false)
}

var PostConsensusAttestationMsgWithNoMsgSigners = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, false, false, true, false)
}

var PostConsensusAttestationMsgWithWrongSig = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, false, true, false, false)
}

var PostConsensusAttestationMsgWithWrongRoot = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, true, false, false, false)
}

var PostConsensusAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPostConsensusMessage {
	return postConsensusAttestationMsg(sk, id, height, false, false, false, false)
}

var postConsensusAttestationMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	height qbft.Height,
	wrongRoot bool,
	wrongBeaconSig bool,
	noMsgSigners bool,
	multiMsgSigners bool,
) *ssv.SignedPostConsensusMessage {
	signer := NewTestingKeyManager()
	signedAtt, root, _ := signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, sk.GetPublicKey().Serialize())

	if wrongBeaconSig {
		signedAtt, _, _ = signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, TestingWrongSK.GetPublicKey().Serialize())
	}

	if wrongRoot {
		root = []byte{1, 2, 3, 4}
	}

	postConsensusMsg := &ssv.PostConsensusMessage{
		Height:          height,
		DutySignature:   signedAtt.Signature[:],
		DutySigningRoot: root,
		Signers:         []types.OperatorID{id},
	}

	if noMsgSigners {
		postConsensusMsg.Signers = []types.OperatorID{}
	}
	if multiMsgSigners {
		postConsensusMsg.Signers = []types.OperatorID{id, 5}
	}

	sig, _ := signer.SignRoot(postConsensusMsg, types.PostConsensusSigType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPostConsensusMessage{
		Message:   postConsensusMsg,
		Signature: sig,
		Signers:   []types.OperatorID{id},
	}
}
