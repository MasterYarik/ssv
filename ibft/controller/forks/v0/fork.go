package v0

import (
	"github.com/bloxapp/ssv/ibft"
	"github.com/bloxapp/ssv/ibft/controller"
	"github.com/bloxapp/ssv/ibft/instance/forks"
	v02 "github.com/bloxapp/ssv/ibft/instance/forks/v0"
	"github.com/bloxapp/ssv/ibft/pipeline"
)

// ForkV0 is the genesis fork for controller
type ForkV0 struct {
	ctrl         *controller.Controller
	instanceFork forks.Fork
}

// New returns new ForkV0
func New() *ForkV0 {
	return &ForkV0{
		instanceFork: v02.New(),
	}
}

// SlotTick implementation
func (v0 *ForkV0) SlotTick(slot uint64) {

}

// Apply fork on controller
func (v0 *ForkV0) Apply(ctrl ibft.Controller) {
	v0.ctrl = ctrl.(*controller.Controller)
}

// InstanceFork returns instance fork
func (v0 *ForkV0) InstanceFork() forks.Fork {
	return v0.instanceFork
}

// ValidateDecidedMsg impl
func (v0 *ForkV0) ValidateDecidedMsg() pipeline.Pipeline {
	return v0.ctrl.ValidateDecidedMsgV0()
}