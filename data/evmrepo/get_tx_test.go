package evmrepo

import (
	"context"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func TestGetTx(t *testing.T) {
	ethClient, err := ethclient.Dial("https://nd-629-499-152.p2pify.com/9d54c0800de991110a4e8e5dc6300b3a")
	if err != nil {
		log.Fatalf("error creating eth client: %v", err)
	}
	defer ethClient.Close()
	client, err := NewEvmRepository(*ethClient)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	successfulTx := "0xefd319bb86b954a8e8cd7d9396546db8d3251910209cd8b1b9a674ef8585f226"
	// failedTx := "0x5dcbfe934287c50363e5c82502739aadd4d535a1f7c0ccd7a8088fb4dfd800da"

	tx, err := client.GetTx(context.TODO(), successfulTx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, successfulTx, tx.TxHash)

}
