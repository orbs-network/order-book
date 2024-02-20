package evmrepo

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
)

func setup() (*evmRepository, *mocks.MockBcBackend) {
	mockBcBackend := mocks.NewMockBcBackend()

	client, err := NewEvmRepository(mockBcBackend.Backend())
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	return client, mockBcBackend
}

// TODO: figure out how to test all flows
func TestEvmRepo_GetTx(t *testing.T) {

	address := common.HexToAddress("0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e")

	transaction := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		GasPrice: big.NewInt(1000000000),
		Gas:      21000,
		To:       &address,
		Value:    big.NewInt(1000000000000000000),
		Data:     []byte{},
	})

	t.Run("GetTx returns TX_PENDING when transaction is pending", func(t *testing.T) {
		client, mockBcBackend := setup()
		pendingTx, _ := mockBcBackend.CreateTx(transaction, false)
		tx, err := client.GetTx(context.Background(), pendingTx.Hash().Hex())
		assert.Equal(t, models.Tx{
			Status:    models.TX_PENDING,
			TxHash:    pendingTx.Hash().Hex(),
			Block:     nil,
			Timestamp: nil,
		}, *tx)
		assert.NoError(t, err)
	})

	t.Run("GetTx returns TX_SUCCESS when transaction is successful", func(t *testing.T) {
		client, mockBcBackend := setup()
		successfulTx, _ := mockBcBackend.CreateTx(transaction, true)
		tx, err := client.GetTx(context.Background(), successfulTx.Hash().Hex())

		expectedBlockNumber := int64(1)
		expectedBlockTs := time.Unix(10, 0)

		assert.Equal(t, models.Tx{
			Status:    models.TX_SUCCESS,
			TxHash:    successfulTx.Hash().Hex(),
			Block:     &expectedBlockNumber,
			Timestamp: &expectedBlockTs,
		}, *tx)
		assert.NoError(t, err)
	})

	t.Run("`ErrNotFound` err and nil when transaction is not found", func(t *testing.T) {
		client, _ := setup()
		tx, err := client.GetTx(context.Background(), "0x023r2323r")
		assert.Nil(t, tx)
		assert.Equal(t, models.ErrNotFound, err)
	})

}
