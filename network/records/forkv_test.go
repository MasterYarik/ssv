package records

import (
	crand "crypto/rand"
	"github.com/bloxapp/ssv/network/commons"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ForkVersionEntry(t *testing.T) {
	priv, _, err := crypto.GenerateSecp256k1Key(crand.Reader)
	require.NoError(t, err)
	sk := fromInterfacePrivKey(priv)
	ip, err := commons.IPAddr()
	require.NoError(t, err)
	node, err := CreateLocalNode(sk, "", ip, commons.DefaultUDP, commons.DefaultTCP)
	require.NoError(t, err)

	require.NoError(t, SetForkVersionEntry(node, "v0"))
	t.Log("ENR with fork version:", node.Node().String())

	fv, err := GetForkVersionEntry(node.Node().Record())
	require.NoError(t, err)
	require.Equal(t, "v0", fv)
}
