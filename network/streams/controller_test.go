package streams

import (
	"bytes"
	"context"
	forksv1 "github.com/bloxapp/ssv/network/forks/v1"
	ssv_protocol "github.com/bloxapp/ssv/protocol/v1/message"
	libp2pnetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestStreamCtrl(t *testing.T) {
	hosts := testHosts(t, 3)

	prot := protocol.ID("/test/protocol")

	//logger := zaptest.NewLogger(t)
	logger := zap.L()
	ctrl0 := NewStreamController(context.Background(), logger.With(zap.String("who", "node-0")),
		hosts[0], forksv1.New(), time.Second)
	ctrl1 := NewStreamController(context.Background(), logger.With(zap.String("who", "node-0")),
		hosts[1], forksv1.New(), time.Second)

	t.Run("handle request", func(t *testing.T) {
		hosts[0].SetStreamHandler(prot, func(stream libp2pnetwork.Stream) {
			msg, res, done, err := ctrl0.HandleStream(stream)
			defer done()
			require.NoError(t, err)
			require.NotNil(t, msg)
			resp, err := dummyMsg().MarshalJSON()
			require.NoError(t, err)
			require.NoError(t, res(resp))
		})
		d, err := dummyMsg().MarshalJSON()
		require.NoError(t, err)
		res, err := ctrl1.Request(hosts[0].ID(), prot, d)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.True(t, bytes.Equal(res, d))
	})

	t.Run("request deadline", func(t *testing.T) {
		timeout := time.Millisecond * 10
		ctrl0.(*streamCtrl).requestTimeout = timeout
		hosts[1].SetStreamHandler(prot, func(stream libp2pnetwork.Stream) {
			msg, s, done, err := ctrl0.HandleStream(stream)
			done()
			require.NoError(t, err)
			require.NotNil(t, msg)
			require.NotNil(t, s)
			<-time.After(timeout + time.Millisecond)
		})
		d, err := dummyMsg().MarshalJSON()
		require.NoError(t, err)
		res, err := ctrl0.Request(hosts[0].ID(), prot, d)
		require.Error(t, err)
		require.Nil(t, res)
	})

}

func dummyMsg() *ssv_protocol.SSVMessage {
	return &ssv_protocol.SSVMessage{Data: []byte("dummy")}
}
