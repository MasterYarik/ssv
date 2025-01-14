package controller

import (
	"context"
	"testing"
	"time"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	forksprotocol "github.com/bloxapp/ssv/protocol/forks"
	beaconprotocol "github.com/bloxapp/ssv/protocol/v1/blockchain/beacon"
	"github.com/bloxapp/ssv/protocol/v1/message"
	protocolp2p "github.com/bloxapp/ssv/protocol/v1/p2p"
	"github.com/bloxapp/ssv/protocol/v1/qbft"
	"github.com/bloxapp/ssv/protocol/v1/qbft/instance"
	qbftstorage "github.com/bloxapp/ssv/protocol/v1/qbft/storage"
	testingprotocol "github.com/bloxapp/ssv/protocol/v1/testing"
	"github.com/bloxapp/ssv/storage/basedb"
	"github.com/bloxapp/ssv/storage/kv"
)

//
//type testStorage struct {
//	highestDecided *message.SignedMessage
//	msgs           map[string]*message.SignedMessage
//	lock           sync.Mutex
//}
//
//func newTestStorage(highestDecided *message.SignedMessage) qbftstorage.QBFTStore {
//	return &testStorage{
//		highestDecided: highestDecided,
//		msgs:           map[string]*message.SignedMessage{},
//		lock:           sync.Mutex{},
//	}
//}
//
//func msgKey(identifier []byte, Height message.Height) string {
//	return fmt.Sprintf("%s_%d", string(identifier), Height)
//}
//
//func (s *testStorage) GetLastDecided(identifier message.Identifier) (*message.SignedMessage, error) {
//	return s.highestDecided, nil
//}
//
//// SaveLastDecided saves the given decided message, after checking that it is indeed the highest
//func (s *testStorage) SaveLastDecided(signedMsg ...*message.SignedMessage) error {
//	return nil
//}
//
//// GetDecided returns historical decided messages in the given range
//func (s *testStorage) GetDecided(identifier message.Identifier, from message.Height, to message.Height) ([]*message.SignedMessage, error) {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//
//	var msgs []*message.SignedMessage
//	for i := from; i <= to; i++ {
//		k := msgKey(identifier, i)
//		if msg, ok := s.msgs[k]; ok {
//			msgs = append(msgs, msg)
//		}
//	}
//	return msgs, nil
//}
//
//// SaveDecided saves historical decided messages
//func (s *testStorage) SaveDecided(signedMsg ...*message.SignedMessage) error {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//
//	for _, msg := range signedMsg {
//		if msg == nil || msg.Message == nil {
//			continue
//		}
//		k := msgKey(msg.Message.Identifier, msg.Message.Height)
//		s.msgs[k] = msg
//	}
//	return nil
//}
//
//// SaveCurrentInstance saves the state for the current running (not yet decided) instance
//func (s *testStorage) SaveCurrentInstance(identifier message.Identifier, state *qbft.State) error {
//	return nil
//}
//
//// GetCurrentInstance returns the state for the current running (not yet decided) instance
//func (s *testStorage) GetCurrentInstance(identifier message.Identifier) (*qbft.State, bool, error) {
//	return nil, false, nil
//}
//
//// GetLastChangeRoundMsg returns the latest broadcasted msg from the instance
//func (s *testStorage) GetLastChangeRoundMsg(identifier message.Identifier) (*message.SignedMessage, error) {
//	return nil, nil
//}
//
//func (s *testStorage) SaveLastChangeRoundMsg(msg *message.SignedMessage) error {
//	return nil
//}
//
//func (s *testStorage) CleanLastChangeRound(identifier message.Identifier) {}

//func TestDecidedRequiresSync(t *testing.T) {
//	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
//	secretKeys, _ := testingprotocol.GenerateBLSKeys(uids...)
//
//	height0 := atomic.Value{}
//	height0.Store(message.Height(0))
//
//	height3 := atomic.Value{}
//	height3.Store(message.Height(3))
//
//	tests := []struct {
//		name            string
//		currentInstance instance.Instancer
//		highestDecided  *message.SignedMessage
//		msg             *message.SignedMessage
//		expectedRes     bool
//		expectedErr     string
//		initState       uint32
//	}{
//		{
//			"decided from future, requires sync.",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height3,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  4,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			true,
//			"",
//			Ready,
//		},
//		{
//			"decided from future, requires sync. current is nil",
//			nil,
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  4,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			true,
//			"",
//			Ready,
//		},
//		{
//			"decided when init failed to sync",
//			nil,
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			true,
//			"",
//			NotStarted,
//		},
//		{
//			"decided from far future, requires sync.",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height3,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  10,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			true,
//			"",
//			Ready,
//		},
//		{
//			"decided from past, doesn't requires sync.",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height3,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//			"",
//			Ready,
//		},
//		{
//			"decided for current",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height3,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//			"",
//			Ready,
//		},
//		{
//			"decided for seq 0",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height0,
//			}),
//			nil,
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  0,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//			"",
//			Ready,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			storage := newTestStorage(test.highestDecided)
//			currentInstanceLock := &sync.RWMutex{}
//			ctrl := Controller{
//				currentInstance:     test.currentInstance,
//				instanceStorage:     storage,
//				changeRoundStorage:  storage,
//				state:               test.initState,
//				currentInstanceLock: currentInstanceLock,
//				forkLock:            &sync.Mutex{},
//			}
//
//			ctrl.fork = forksfactory.NewFork(forksprotocol.V0ForkVersion)
//			ctrl.decidedFactory = factory.NewDecidedFactory(zap.L(), ctrl.getNodeMode(), storage, nil)
//			ctrl.decidedStrategy = ctrl.decidedFactory.GetStrategy()
//
//			res, err := ctrl.decidedRequiresSync(test.msg)
//			require.EqualValues(t, test.expectedRes, res)
//			if len(test.expectedErr) > 0 {
//				require.EqualError(t, err, test.expectedErr)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}

//func TestDecideIsCurrentInstance(t *testing.T) {
//	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
//	secretKeys, _ := testingprotocol.GenerateBLSKeys(uids...)
//
//	height1 := atomic.Value{}
//	height1.Store(message.Height(1))
//
//	height4 := atomic.Value{}
//	height4.Store(message.Height(4))
//
//	tests := []struct {
//		name            string
//		currentInstance instance.Instancer
//		msg             *message.SignedMessage
//		expectedRes     bool
//	}{
//		{
//			"current instance",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height1,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			true,
//		},
//		{
//			"current instance nil",
//			nil,
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//		},
//		{
//			"current instance empty",
//			&instance.Instance{},
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  1,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//		},
//		{
//			"current instance seq lower",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height1,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//		},
//		{
//			"current instance seq higher",
//			instance.NewInstanceWithState(&qbft.State{
//				Height: height4,
//			}),
//			testingprotocol.SignMsg(t, secretKeys, []message.OperatorID{message.OperatorID(1)}, &message.ConsensusMessage{
//				MsgType: message.CommitMsgType,
//				Height:  2,
//				Data:    commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//			}),
//			false,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			currentInstanceLock := &sync.RWMutex{}
//			ibft := Controller{
//				currentInstance:     test.currentInstance,
//				currentInstanceLock: currentInstanceLock,
//				forkLock:            &sync.Mutex{},
//			}
//			require.EqualValues(t, test.expectedRes, ibft.decidedForCurrentInstance(test.msg))
//		})
//	}
//}

func TestForceDecided(t *testing.T) {
	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
	sks, nodes := testingprotocol.GenerateBLSKeys(uids...)
	pi, err := protocolp2p.GenPeerID()
	require.NoError(t, err)

	network := protocolp2p.NewMockNetwork(zap.L(), pi, 10)

	identifier := []byte("Identifier_11")
	s1 := testingprotocol.PopulatedStorage(t, sks, 3, 3)
	i1 := populatedIbft(1, identifier, network, s1, sks, nodes, newTestSigner())
	// test before sync
	highest, err := i1.(*Controller).decidedStrategy.GetLastDecided(identifier)
	require.NotNil(t, highest)
	require.NoError(t, err)
	require.EqualValues(t, 3, highest.Message.Height)

	time.Sleep(time.Second * 1) // wait for sync to complete

	go func() {
		time.Sleep(time.Millisecond * 500) // wait for instance to start

		signers := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}

		encodedCommit, err := (&message.CommitData{Data: []byte("value")}).Encode()
		require.NoError(t, err)
		decidedMsg := testingprotocol.AggregateSign(t, sks, signers, &message.ConsensusMessage{
			MsgType:    message.CommitMsgType,
			Height:     message.Height(4),
			Round:      message.Round(1),
			Identifier: identifier,
			Data:       encodedCommit,
		})

		require.NoError(t, i1.(*Controller).processDecidedMessage(decidedMsg))
	}()

	res, err := i1.StartInstance(instance.ControllerStartInstanceOptions{
		Logger:    zap.L(),
		SeqNumber: 4,
		Value:     commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
	})
	require.NoError(t, err)
	require.True(t, res.Decided)

	highest, err = i1.(*Controller).decidedStrategy.GetLastDecided(identifier)
	require.NotNil(t, highest)
	require.NoError(t, err)
	require.EqualValues(t, 4, highest.Message.Height)
}

func TestSyncAfterDecided(t *testing.T) {
	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
	sks, nodes := testingprotocol.GenerateBLSKeys(uids...)
	pi, err := protocolp2p.GenPeerID()
	require.NoError(t, err)

	identifier := []byte("Identifier_11")

	decidedMsg := testingprotocol.AggregateSign(t, sks, uids, &message.ConsensusMessage{
		MsgType:    message.CommitMsgType,
		Height:     message.Height(10),
		Round:      message.Round(3),
		Identifier: identifier,
		Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
	})

	network := protocolp2p.NewMockNetwork(zap.L(), pi, 10)
	network.SetLastDecidedHandler(generateLastDecidedHandler(t, identifier, decidedMsg))
	network.SetGetHistoryHandler(generateGetHistoryHandler(t, sks, uids, identifier, 4, 10))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	network.Start(ctx)
	network.AddPeers(message.Identifier(identifier).GetValidatorPK(), network)

	s1 := testingprotocol.PopulatedStorage(t, sks, 3, 4)
	i1 := populatedIbft(1, identifier, network, s1, sks, nodes, newTestSigner())

	_ = populatedIbft(2, identifier, network, testingprotocol.PopulatedStorage(t, sks, 3, 10), sks, nodes, newTestSigner())

	// test before sync
	highest, err := i1.(*Controller).decidedStrategy.GetLastDecided(identifier)
	require.NotNil(t, highest)
	require.NoError(t, err)
	require.EqualValues(t, 4, highest.Message.Height)

	require.NoError(t, i1.(*Controller).processDecidedMessage(decidedMsg))

	time.Sleep(time.Millisecond * 500) // wait for sync to complete
	highest, err = i1.(*Controller).decidedStrategy.GetLastDecided(identifier)
	require.NotNil(t, highest)
	require.NoError(t, err)
	require.EqualValues(t, message.Height(10), highest.Message.Height)
}

func TestSyncFromScratchAfterDecided(t *testing.T) {
	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
	sks, nodes := testingprotocol.GenerateBLSKeys(uids...)
	db, _ := kv.New(basedb.Options{
		Type:   "badger-memory",
		Path:   "",
		Logger: zap.L(),
	})
	pi, err := protocolp2p.GenPeerID()
	require.NoError(t, err)

	identifier := []byte("Identifier_11")
	decidedMsg := testingprotocol.AggregateSign(t, sks, uids, &message.ConsensusMessage{
		MsgType:    message.CommitMsgType,
		Height:     message.Height(10),
		Round:      message.Round(3),
		Identifier: identifier,
		Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
	})

	network := protocolp2p.NewMockNetwork(zap.L(), pi, 10)
	network.SetLastDecidedHandler(generateLastDecidedHandler(t, identifier, decidedMsg))
	network.SetGetHistoryHandler(generateGetHistoryHandler(t, sks, uids, identifier, 0, 10))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	network.Start(ctx)
	network.AddPeers(message.Identifier(identifier).GetValidatorPK(), network)

	s1 := qbftstorage.NewQBFTStore(db, zap.L(), "attestations")
	i1 := populatedIbft(1, identifier, network, s1, sks, nodes, newTestSigner())

	_ = populatedIbft(2, identifier, network, testingprotocol.PopulatedStorage(t, sks, 3, 10), sks, nodes, newTestSigner())

	require.NoError(t, i1.(*Controller).processDecidedMessage(decidedMsg))

	time.Sleep(time.Millisecond * 500) // wait for sync to complete
	highest, err := i1.(*Controller).decidedStrategy.GetLastDecided(identifier)
	require.NotNil(t, highest)
	require.NoError(t, err)
	require.EqualValues(t, 10, highest.Message.Height)
}

func TestValidateDecidedMsg(t *testing.T) {
	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
	sks, nodes := testingprotocol.GenerateBLSKeys(uids...)
	pi, err := protocolp2p.GenPeerID()
	require.NoError(t, err)

	network := protocolp2p.NewMockNetwork(zap.L(), pi, 10)
	identifier := []byte("Identifier_11")
	ibft := populatedIbft(1, identifier, network, testingprotocol.PopulatedStorage(t, sks, 3, 10), sks, nodes, newTestSigner())

	tests := []struct {
		name          string
		msg           *message.SignedMessage
		expectedError error
	}{
		{
			"valid",
			testingprotocol.AggregateSign(t, sks, uids, &message.ConsensusMessage{
				MsgType:    message.CommitMsgType,
				Height:     message.Height(11),
				Round:      message.Round(3),
				Identifier: identifier,
				Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
			}),
			nil,
		},
		{
			"invalid msg stage",
			testingprotocol.AggregateSign(t, sks, uids, &message.ConsensusMessage{
				MsgType:    message.PrepareMsgType,
				Height:     message.Height(11),
				Round:      message.Round(3),
				Identifier: identifier,
				Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
			}),
			errors.New("message type is wrong"),
		},
		{
			"invalid msg sig",
			testingprotocol.AggregateInvalidSign(t, sks, &message.ConsensusMessage{
				MsgType:    message.CommitMsgType,
				Height:     message.Height(11),
				Round:      message.Round(3),
				Identifier: identifier,
				Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
			}),
			errors.New("failed to verify signature"),
		},
		{
			"valid first decided",
			testingprotocol.AggregateSign(t, sks, uids, &message.ConsensusMessage{
				MsgType:    message.CommitMsgType,
				Height:     message.Height(0),
				Round:      message.Round(3),
				Identifier: identifier,
				Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
			}),
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedError != nil {
				err := ibft.(*Controller).ValidateDecidedMsg(test.msg)
				require.EqualError(t, err, test.expectedError.Error())
			} else {
				require.NoError(t, ibft.(*Controller).ValidateDecidedMsg(test.msg))
			}
		})
	}
}

//func TestController_checkDecidedMessageSigners(t *testing.T) {
//	uids := []message.OperatorID{message.OperatorID(1), message.OperatorID(2), message.OperatorID(3), message.OperatorID(4)}
//	secretKeys, nodes := testingprotocol.GenerateBLSKeys(uids...)
//	skQuorum := map[message.OperatorID]*bls.SecretKey{}
//	for i, sk := range secretKeys {
//		skQuorum[i] = sk
//	}
//	delete(skQuorum, 4)
//	identifier := []byte("Identifier_2")
//
//	incompleteDecided := testingprotocol.AggregateSign(t, skQuorum, uids[:3], &message.ConsensusMessage{
//		MsgType:    message.CommitMsgType,
//		Height:     message.Height(2),
//		Identifier: identifier[:],
//		Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//	})
//
//	completeDecided := testingprotocol.AggregateSign(t, secretKeys, uids, &message.ConsensusMessage{
//		MsgType:    message.CommitMsgType,
//		Height:     message.Height(2),
//		Identifier: identifier[:],
//		Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
//	})
//
//	share := &beaconprotocol.Share{
//		NodeID:    1,
//		PublicKey: secretKeys[1].GetPublicKey(),
//		Committee: nodes,
//	}
//
//	id := atomic.Value{}
//	id.Store(message.Identifier(identifier))
//
//	height := atomic.Value{}
//	height.Store(message.Height(2))
//
//	storage := newTestStorage(nil)
//	currentInstanceLock := &sync.RWMutex{}
//	ctrl := Controller{
//		ValidatorShare: share,
//		currentInstance: instance.NewInstanceWithState(&qbft.State{
//			Identifier: id,
//			Height:     height,
//		}),
//		instanceStorage:     storage,
//		changeRoundStorage:  storage,
//		currentInstanceLock: currentInstanceLock,
//		forkLock:            &sync.Mutex{},
//	}
//
//	ctrl.fork = forksfactory.NewFork(forksprotocol.V0ForkVersion)
//	ctrl.decidedFactory = factory.NewDecidedFactory(zap.L(), ctrl.getNodeMode(), storage, nil)
//	ctrl.decidedStrategy = ctrl.decidedFactory.GetStrategy()
//
//	_, err := ctrl.decidedStrategy.SaveDecided(incompleteDecided)
//	require.NoError(t, err)
//
//	// check message with similar number of signers
//	require.True(t, ctrl.checkDecidedMessageSigners(incompleteDecided, incompleteDecided))
//	// check message with more signers
//	require.False(t, ctrl.checkDecidedMessageSigners(incompleteDecided, completeDecided))
//}

// TODO: (lint) fix test
//nolint
func populatedIbft(
	nodeID message.OperatorID,
	identifier []byte,
	network protocolp2p.MockNetwork,
	ibftStorage qbftstorage.QBFTStore,
	sks map[message.OperatorID]*bls.SecretKey,
	nodes map[message.OperatorID]*beaconprotocol.Node,
	signer beaconprotocol.Signer,
) IController {
	share := &beaconprotocol.Share{
		NodeID:    nodeID,
		PublicKey: sks[1].GetPublicKey(),
		Committee: nodes,
	}

	opts := Options{
		Context:        context.Background(),
		Role:           message.RoleTypeAttester,
		Identifier:     identifier,
		Logger:         zap.L(),
		Storage:        ibftStorage,
		Network:        network,
		InstanceConfig: qbft.DefaultConsensusParams(),
		ValidatorShare: share,
		Version:        forksprotocol.V0ForkVersion, // TODO need to check v1 fork too? (:Niv)
		Beacon:         nil,                         // ?
		Signer:         signer,
		SyncRateLimit:  time.Millisecond * 100,
		SigTimeout:     time.Second * 5,
		ReadMode:       false,
	}
	ret := New(opts)

	ret.(*Controller).state = Ready // as if they are already synced
	return ret
}

type testSigner struct {
}

func newTestSigner() beaconprotocol.Signer {
	return &testSigner{}
}

func (s *testSigner) AddShare(shareKey *bls.SecretKey) error {
	return nil
}

func (s *testSigner) SignIBFTMessage(message *message.ConsensusMessage, pk []byte, forkVersion string) ([]byte, error) {
	return nil, nil
}

func (s *testSigner) SignAttestation(data *spec.AttestationData, duty *beaconprotocol.Duty, pk []byte) (*spec.Attestation, []byte, error) {
	return nil, nil, nil
}

func commitDataToBytes(t *testing.T, input *message.CommitData) []byte {
	ret, err := input.Encode()
	require.NoError(t, err)
	return ret
}

func generateGetHistoryHandler(t *testing.T, sks map[message.OperatorID]*bls.SecretKey, uids []message.OperatorID, identifier []byte, from, to int) protocolp2p.EventHandler {
	return func(e protocolp2p.MockMessageEvent) *message.SSVMessage {
		decidedMsgs := make([]*message.SignedMessage, 0)
		heights := make([]message.Height, 0)
		for i := from; i <= to; i++ {
			decidedMsgs = append(decidedMsgs, testingprotocol.AggregateSign(t, sks, uids, &message.ConsensusMessage{
				MsgType:    message.CommitMsgType,
				Height:     message.Height(i),
				Round:      message.Round(3),
				Identifier: identifier,
				Data:       commitDataToBytes(t, &message.CommitData{Data: []byte("value")}),
			}))
			heights = append(heights, message.Height(i))
		}

		sm := &message.SyncMessage{
			Protocol: message.LastDecidedType,
			Params: &message.SyncParams{
				Height:     heights,
				Identifier: identifier,
			},
			Data:   decidedMsgs,
			Status: message.StatusSuccess,
		}
		em, err := sm.Encode()
		require.NoError(t, err)

		msg := &message.SSVMessage{
			MsgType: message.SSVDecidedMsgType,
			ID:      identifier,
			Data:    em,
		}

		return msg
	}
}

func generateLastDecidedHandler(t *testing.T, identifier []byte, decidedMsg *message.SignedMessage) protocolp2p.EventHandler {
	return func(e protocolp2p.MockMessageEvent) *message.SSVMessage {
		sm := &message.SyncMessage{
			Protocol: message.LastDecidedType,
			Params: &message.SyncParams{
				Height:     []message.Height{message.Height(9), message.Height(10)},
				Identifier: identifier,
			},
			Data:   []*message.SignedMessage{decidedMsg},
			Status: message.StatusSuccess,
		}
		em, err := sm.Encode()
		require.NoError(t, err)

		msg := &message.SSVMessage{
			MsgType: message.SSVDecidedMsgType,
			ID:      identifier,
			Data:    em,
		}

		return msg
	}
}
