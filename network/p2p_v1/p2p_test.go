package p2pv1

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/protocol"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestP2pNetwork_SubscribeBroadcast(t *testing.T) {
	n := 4
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loggerFactory := func(who string) *zap.Logger {
		//logger := zaptest.NewLogger(t).With(zap.String("who", who))
		logger := zap.L().With(zap.String("who", who))
		return logger
	}
	ln, err := CreateAndStartLocalNet(ctx, loggerFactory, n, n/2, false)
	require.NoError(t, err)
	require.Len(t, ln.Nodes, n)

	pks := []string{"b768cdc2b2e0a859052bf04d1cd66383c96d95096a5287d08151494ce709556ba39c1300fbb902a0e2ebb7c31dc4e400",
		"824b9024767a01b56790a72afb5f18bb0f97d5bddb946a7bd8dd35cc607c35a4d76be21f24f484d0d478b99dc63ed170"}

	routers := make([]*dummyRouter, n)
	for i, node := range ln.Nodes {
		routers[i] = &dummyRouter{i: i, logger: zaptest.NewLogger(t).With(zap.Int("i", i))}
		node.UseMessageRouter(routers[i])
	}

	for _, pk := range pks {
		vpk, err := hex.DecodeString(pk)
		require.NoError(t, err)
		for _, node := range ln.Nodes {
			require.NoError(t, node.Subscribe(vpk))
		}
	}
	// let the nodes subscribe
	// TODO: remove timeout
	<-time.After(time.Second * 3)

	msg1, err := dummyMsg(pks[0], 1)
	require.NoError(t, err)
	msg2, err := dummyMsg(pks[1], 2)
	require.NoError(t, err)
	msg3, err := dummyMsg(pks[0], 3)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, ln.Nodes[1].Broadcast(*msg1))
		<-time.After(time.Millisecond * 10)
		require.NoError(t, ln.Nodes[2].Broadcast(*msg3))
		<-time.After(time.Millisecond * 2)
		require.NoError(t, ln.Nodes[1].Broadcast(*msg1))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-time.After(time.Millisecond * 10)
		require.NoError(t, ln.Nodes[1].Broadcast(*msg2))
		<-time.After(time.Millisecond * 2)
		require.NoError(t, ln.Nodes[2].Broadcast(*msg1))
		require.NoError(t, ln.Nodes[1].Broadcast(*msg3))
	}()

	wg.Wait()

	// waiting for messages
	wg.Add(1)
	go func() {
		ct, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		defer wg.Done()
		for _, r := range routers {
			for ct.Err() == nil && atomic.LoadUint64(&r.count) < uint64(2) {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	wg.Wait()

	for _, r := range routers {
		require.GreaterOrEqual(t, atomic.LoadUint64(&r.count), uint64(2), "router", r.i)
	}

	for _, node := range ln.Nodes {
		require.NoError(t, node.(*p2pNetwork).Close())
	}
}

type dummyRouter struct {
	logger *zap.Logger
	count  uint64
	i      int
}

func (r *dummyRouter) Route(message protocol.SSVMessage) {
	c := atomic.AddUint64(&r.count, 1)
	r.logger.Debug("got message",
		zap.String("identifier", hex.EncodeToString(message.GetID())),
		zap.Uint64("count", c))
}

func dummyMsg(pkHex string, height int) (*protocol.SSVMessage, error) {
	pk, err := hex.DecodeString(pkHex)
	if err != nil {
		return nil, err
	}
	id := protocol.NewIdentifier(pk, beacon.RoleTypeAttester)
	msgData := fmt.Sprintf(`{
	  "message": {
		"type": 3,
		"round": 2,
		"identifier": "%s",
		"height": %d,
		"value": "bk0iAAAAAAACAAAAAAAAAAbYXFSt2H7SQd5q5u+N0bp6PbbPTQjU25H1QnkbzTECahIBAAAAAADmi+NJfvXZ3iXp2cfs0vYVW+EgGD7DTTvr5EkLtiWq8WsSAQAAAAAAIC8dZTEdD3EvE38B9kDVWkSLy40j0T+TtSrrrBqVjo4="
	  },
	  "signature": "sVV0fsvqQlqliKv/ussGIatxpe8LDWhc9uoaM5WpjbiYvvxUr1eCpz0ja7UT1PGNDdmoGi6xbMC1g/ozhAt4uCdpy0Xdfqbv2hMf2iRL5ZPKOSmMifHbd8yg4PeeceyN",
	  "signers": [1,3,4]
	}`, id, height)
	return &protocol.SSVMessage{
		MsgType: protocol.SSVConsensusMsgType,
		ID:      id,
		Data:    []byte(msgData),
	}, nil
}