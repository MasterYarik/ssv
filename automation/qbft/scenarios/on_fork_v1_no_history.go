package scenarios

import (
	"fmt"
	"sync"
	"time"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/automation/commons"
	"github.com/bloxapp/ssv/automation/qbft/runner"
	"github.com/bloxapp/ssv/network"
	forksprotocol "github.com/bloxapp/ssv/protocol/forks"
	"github.com/bloxapp/ssv/protocol/v1/blockchain/beacon"
	"github.com/bloxapp/ssv/protocol/v1/message"
	qbftstorage "github.com/bloxapp/ssv/protocol/v1/qbft/storage"
	"github.com/bloxapp/ssv/protocol/v1/validator"
)

// OnForkV1NoHistoryScenario is the scenario name for OnForkV1NoHistory
const OnForkV1NoHistoryScenario = "OnForkV1NoHistory"

type onForkV1NoHistory struct {
	logger     *zap.Logger
	share      *beacon.Share
	sks        map[uint64]*bls.SecretKey
	validators []validator.IValidator
}

// newOnForkV1NoHistory creates 'on fork v1' scenario
func newOnForkV1NoHistory(logger *zap.Logger) runner.Scenario {
	return &onForkV1NoHistory{logger: logger}
}

func (f *onForkV1NoHistory) NumOfOperators() int {
	return 4
}

func (f *onForkV1NoHistory) NumOfBootnodes() int {
	return 0
}

func (f *onForkV1NoHistory) NumOfFullNodes() int {
	return 0
}

func (f *onForkV1NoHistory) Name() string {
	return OnForkV1NoHistoryScenario
}

// PreExecution will create messages in v0 format
func (f *onForkV1NoHistory) PreExecution(ctx *runner.ScenarioContext) error {
	share, sks, validators, err := commons.CreateShareAndValidators(ctx.Ctx, f.logger, ctx.LocalNet, ctx.KeyManagers, ctx.Stores)
	if err != nil {
		return errors.Wrap(err, "could not create share")
	}
	// save all references
	f.validators = validators
	f.sks = sks
	f.share = share

	// setting up routers
	routers := make([]*runner.Router, f.NumOfOperators())
	loggerFactory := func(who string) *zap.Logger {
		logger := zap.L().With(zap.String("who", who))
		return logger
	}

	for i, node := range ctx.LocalNet.Nodes {
		routers[i] = &runner.Router{
			Logger:      loggerFactory(fmt.Sprintf("msgRouter-%d", i)),
			Controllers: f.validators[i].(*validator.Validator).Ibfts(),
		}
		node.UseMessageRouter(routers[i])
	}

	return nil
}

func (f *onForkV1NoHistory) Execute(ctx *runner.ScenarioContext) error {
	if len(f.sks) == 0 || f.share == nil {
		return errors.New("pre-execution failed")
	}

	var wg sync.WaitGroup
	var startErr error
	for _, val := range f.validators {
		wg.Add(1)
		go func(val validator.IValidator) {
			defer wg.Done()
			if err := val.Start(); err != nil {
				startErr = errors.Wrap(err, "could not start validator")
			}
			<-time.After(time.Second * 3)
		}(val)
	}
	wg.Wait()

	if startErr != nil {
		return errors.Wrap(startErr, "could not start validators")
	}

	// running instances pre-fork
	if err := f.startInstances(message.Height(1), message.Height(6)); err != nil {
		return errors.Wrap(err, "could not start instances")
	}

	// forking
	for i := 0; i < f.NumOfOperators(); i++ {
		wg.Add(3)
		go func(node network.P2PNetwork) {
			defer wg.Done()
			if err := node.(forksprotocol.ForkHandler).OnFork(forksprotocol.V1ForkVersion); err != nil {
				f.logger.Panic("could not fork network to v1", zap.Error(err))
			}
		}(ctx.LocalNet.Nodes[i])
		go func(val validator.IValidator) {
			defer wg.Done()
			<-time.After(time.Second)
			if err := val.OnFork(forksprotocol.V1ForkVersion); err != nil {
				f.logger.Panic("could not fork to v1", zap.Error(err))
			}
		}(f.validators[i])
		go func(store qbftstorage.QBFTStore) {
			defer wg.Done()
			<-time.After(time.Second)
			if err := store.(forksprotocol.ForkHandler).OnFork(forksprotocol.V1ForkVersion); err != nil {
				f.logger.Panic("could not fork qbft store to v1", zap.Error(err))
			}
		}(ctx.Stores[i])
	}
	wg.Wait()

	f.logger.Debug("------ after fork, waiting 10 seconds...")
	// waiting 10 sec after fork
	<-time.After(time.Second * 10)
	f.logger.Debug("------ starting instances")

	for i := 0; i < f.NumOfOperators(); i++ {
		peers, err := ctx.LocalNet.Nodes[i].Peers(f.share.PublicKey.Serialize())
		if err != nil {
			return errors.Wrap(err, "could not check peers of topic")
		}
		if len(peers) < f.NumOfOperators()/2 {
			return errors.Errorf("node %d could not find enough peers after fork: %d", i, len(peers))
		}
	}

	// running instances post-fork
	if err := f.startInstances(message.Height(7), message.Height(9)); err != nil {
		return errors.Wrap(err, "could not start instance after fork")
	}

	return nil
}

func (f *onForkV1NoHistory) PostExecution(ctx *runner.ScenarioContext) error {
	expectedMsgCount := 9
	msgs, err := ctx.Stores[0].GetDecided(message.NewIdentifier(f.share.PublicKey.Serialize(), message.RoleTypeAttester), message.Height(0), message.Height(expectedMsgCount))
	if err != nil {
		return err
	}
	f.logger.Debug("msgs count", zap.Int("len", len(msgs)))
	if len(msgs) < expectedMsgCount {
		return errors.New("node-0 didn't sync all messages")
	}

	msg, err := ctx.Stores[0].GetLastDecided(message.NewIdentifier(f.share.PublicKey.Serialize(), message.RoleTypeAttester))
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("could not find last decided")
	}
	if msg.Message.Height != message.Height(expectedMsgCount) {
		return errors.Errorf("wrong msg height: %d", msg.Message.Height)
	}

	return nil
}

func (f *onForkV1NoHistory) startInstances(from, to message.Height) error {
	var wg sync.WaitGroup

	h := from

	for h <= to {
		for i := uint64(1); i < uint64(f.NumOfOperators()); i++ {
			wg.Add(1)
			go func(node validator.IValidator, index uint64, seqNumber message.Height) {
				if err := startNode(node, seqNumber, []byte("value"), f.logger); err != nil {
					f.logger.Error("could not start node", zap.Uint64("node", index-1), zap.Error(err))
				}
				wg.Done()
			}(f.validators[i-1], i, h)
		}
		wg.Wait()
		h++
	}
	return nil
}
