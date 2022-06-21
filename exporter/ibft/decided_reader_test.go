package ibft

import (
	"github.com/bloxapp/ssv/ibft/proto"
	ssvstorage "github.com/bloxapp/ssv/storage"
	"github.com/bloxapp/ssv/storage/basedb"
	"github.com/bloxapp/ssv/storage/collections"
	"github.com/bloxapp/ssv/validator/storage"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestDecidedReader_checkDecided(t *testing.T) {
	reader := setupDecidedReaderForTest(t)
	d := reader.(*decidedReader)
	require.NoError(t, d.storage.SaveDecided(&proto.SignedMessage{
		Message: &proto.Message{
			Type:      0,
			Round:     0,
			Lambda:    d.identifier,
			SeqNumber: 0,
			Value:     nil,
		},
		Signature: nil,
		SignerIds: []uint64{1, 2, 3},
	}))

	t.Run("same signers, ignore", func(t *testing.T) {
		known, updatedMsg, err := d.checkDecided(&proto.SignedMessage{
			Message: &proto.Message{
				Type:      0,
				Round:     0,
				Lambda:    d.identifier,
				SeqNumber: 0,
				Value:     nil,
			},
			Signature: nil,
			SignerIds: []uint64{1, 2, 3},
		})
		require.NoError(t, err)
		require.True(t, known)
		require.Equal(t, 3, len(updatedMsg.SignerIds))
	})

	t.Run("update to 4 signers", func(t *testing.T) {
		known, updatedMsg, err := d.checkDecided(&proto.SignedMessage{
			Message: &proto.Message{
				Type:      0,
				Round:     0,
				Lambda:    d.identifier,
				SeqNumber: 0,
				Value:     nil,
			},
			Signature: nil,
			SignerIds: []uint64{2, 3, 4},
		})
		require.NoError(t, err)
		require.False(t, known)
		require.Equal(t, 4, len(updatedMsg.SignerIds))
		require.NoError(t, d.storage.SaveDecided(updatedMsg))
	})

	t.Run("same 4 signers, ignore", func(t *testing.T) {
		known, updatedMsg, err := d.checkDecided(&proto.SignedMessage{
			Message: &proto.Message{
				Type:      0,
				Round:     0,
				Lambda:    d.identifier,
				SeqNumber: 0,
				Value:     nil,
			},
			Signature: nil,
			SignerIds: []uint64{2, 3, 4, 1},
		})
		require.NoError(t, err)
		require.True(t, known)
		require.Equal(t, 4, len(updatedMsg.SignerIds))
	})
}

func setupDecidedReaderForTest(t *testing.T) Reader {
	logger := zap.L()
	db, err := ssvstorage.GetStorageFactory(basedb.Options{
		Type:   "badger-memory",
		Logger: logger,
		Path:   "",
	})
	require.NoError(t, err)
	ibftStorage := collections.NewIbft(db, logger, "attestation")
	_ = bls.Init(bls.BLS12_381)

	pubKey := &bls.PublicKey{}
	bls.GetGeneratorOfPublicKey(pubKey)

	cr := NewDecidedReader(DecidedReaderOptions{
		Logger:  logger,
		Storage: &ibftStorage,
		Network: nil,
		Config:  nil,
		ValidatorShare: &storage.Share{
			NodeID:       1,
			PublicKey:    pubKey,
			Committee:    nil,
			Metadata:     nil,
			OwnerAddress: "",
			Operators:    nil,
		},
		Out: nil,
	})

	return cr
}
