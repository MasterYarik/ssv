package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/docs/spec/types"
	"time"
)

// DutyRunners is a map of duty runners mapped by msg id hex.
type DutyRunners map[beacon.RoleType]*Runner

// DutyRunnerForMsgID returns a Runner from the provided msg ID, or nil if not found
func (ci DutyRunners) DutyRunnerForMsgID(msgID types.MessageID) *Runner {
	role := msgID.GetRoleType()
	return ci[role]
}

type Network interface {
	Broadcast(message types.Encoder) error
}

// Storage is a persistent storage for the SSV
type Storage interface {
}

type BeaconNode interface {
	// GetBeaconNetwork returns the beacon network the node is on
	GetBeaconNetwork() BeaconNetwork
	// GetAttestationData returns attestation data by the given slot and committee index
	GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error)
	// SubmitAttestation submit the attestation to the node
	SubmitAttestation(attestation *phase0.Attestation) error
	// GetBeaconBlock returns beacon block by the given slot and committee index
	GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (*altair.BeaconBlock, error)
	// SubmitBeaconBlock submit the block to the node
	SubmitBeaconBlock(block *altair.SignedBeaconBlock) error
	// SubmitAggregateSelectionProof returns an AggregateAndProof object
	SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, slotSig []byte) (*phase0.AggregateAndProof, error)
	// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
	SubmitSignedAggregateSelectionProof(msg *phase0.SignedAggregateAndProof) error
}

// Available networks.
const (
	// PraterNetwork represents the Prater test network.
	PraterNetwork BeaconNetwork = "prater"

	// MainNetwork represents the main network.
	MainNetwork BeaconNetwork = "mainnet"

	// NowTestNetwork is a simple test network with genesis time always equal to now, meaning now is slot 0
	NowTestNetwork BeaconNetwork = "now_test_network"
)

// BeaconNetwork represents the network.
type BeaconNetwork string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) BeaconNetwork {
	switch n {
	case string(PraterNetwork):
		return PraterNetwork
	case string(MainNetwork):
		return MainNetwork
	case string(NowTestNetwork):
		return NowTestNetwork
	default:
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n BeaconNetwork) ForkVersion() []byte {
	switch n {
	case PraterNetwork:
		return []byte{0x00, 0x00, 0x10, 0x20}
	case MainNetwork:
		return []byte{0, 0, 0, 0}
	case NowTestNetwork:
		return []byte{0x99, 0x99, 0x99, 0x99}
	default:
		return nil
	}
}

// MinGenesisTime returns min genesis time value
func (n BeaconNetwork) MinGenesisTime() uint64 {
	switch n {
	case PraterNetwork:
		return 1616508000
	case MainNetwork:
		return 1606824023
	case NowTestNetwork:
		return uint64(time.Now().Unix())
	default:
		return 0
	}
}

// SlotDurationSec returns slot duration
func (n BeaconNetwork) SlotDurationSec() time.Duration {
	return 12 * time.Second
}

// SlotsPerEpoch returns number of slots per one epoch
func (n BeaconNetwork) SlotsPerEpoch() uint64 {
	return 32
}

// EstimatedCurrentSlot returns the estimation of the current slot
func (n BeaconNetwork) EstimatedCurrentSlot() phase0.Slot {
	return n.EstimatedSlotAtTime(time.Now().Unix())
}

// EstimatedSlotAtTime estimates slot at the given time
func (n BeaconNetwork) EstimatedSlotAtTime(time int64) phase0.Slot {
	genesis := int64(n.MinGenesisTime())
	if time < genesis {
		return 0
	}
	return phase0.Slot(uint64(time-genesis) / uint64(n.SlotDurationSec().Seconds()))
}

// EstimatedCurrentEpoch estimates the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#compute_start_slot_at_epoch
func (n BeaconNetwork) EstimatedCurrentEpoch() phase0.Epoch {
	return n.EstimatedEpochAtSlot(n.EstimatedCurrentSlot())
}

// EstimatedEpochAtSlot estimates epoch at the given slot
func (n BeaconNetwork) EstimatedEpochAtSlot(slot phase0.Slot) phase0.Epoch {
	return phase0.Epoch(slot / phase0.Slot(n.SlotsPerEpoch()))
}
