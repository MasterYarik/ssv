package types

import (
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
)

// ConsensusData holds all relevant duty and data Decided on by consensus
type ConsensusData struct {
	Duty            *beacon.Duty
	AttestationData *phase0.AttestationData
	BlockData       *phase0.BeaconBlock
}

func (cid *ConsensusData) Encode() ([]byte, error) {
	//m := make(map[string]string)
	//if cid.Duty != nil {
	//	d, err := json.Marshal(cid.Duty)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "duty marshaling failed")
	//	}
	//	m["duty"] = hex.EncodeToString(d)
	//}
	//
	//if cid.AttestationData != nil {
	//	d, err := ssz.Marshal(cid.AttestationData)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "attestation data marshaling failed")
	//	}
	//	m["att_data"] = hex.EncodeToString(d)
	//}
	//
	//if cid.BlockData != nil {
	//	d, err := ssz.Marshal(cid.BlockData)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "block data marshaling failed")
	//	}
	//	m["block_data"] = hex.EncodeToString(d)
	//}
	return json.Marshal(cid)
}

func (cid *ConsensusData) Decode(data []byte) error {
	//m := make(map[string]string)
	//if err := json.Unmarshal(data, &m); err != nil {
	//	return errors.Wrap(err, "could not unmarshal ConsensusData")
	//}
	//
	//if val, ok := m["duty"]; ok {
	//	cid.Duty = &beacon.Duty{}
	//	d, err := hex.DecodeString(val)
	//	if err != nil {
	//		return errors.Wrap(err, "Duty decode string failed")
	//	}
	//	if err := json.Unmarshal(d, cid.Duty); err != nil {
	//		cid.Duty = nil
	//		return errors.Wrap(err, "could not unmarshal duty")
	//	}
	//}
	//
	//if val, ok := m["att_data"]; ok {
	//	cid.AttestationData = &phase0.AttestationData{}
	//	d, err := hex.DecodeString(val)
	//	if err != nil {
	//		return errors.Wrap(err, "AttestationData decode string failed")
	//	}
	//	if err := ssz.Unmarshal(d, cid.AttestationData); err != nil {
	//		cid.AttestationData = nil
	//		return errors.Wrap(err, "could not unmarshal AttestationData")
	//	}
	//}
	//
	//if val, ok := m["block_data"]; ok {
	//	cid.BlockData = &phase0.BeaconBlock{}
	//	d, err := hex.DecodeString(val)
	//	if err != nil {
	//		return errors.Wrap(err, "BlockData decode string failed")
	//	}
	//	if err := ssz.Unmarshal(d, cid.BlockData); err != nil {
	//		cid.BlockData = nil
	//		return errors.Wrap(err, "could not unmarshal BeaconBlock")
	//	}
	//}
	return json.Unmarshal(data, &cid)
}