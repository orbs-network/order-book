package service

import (
	"errors"

	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/data/storeblockchain"
)

type EvmClient struct {
	orderBookStore  store.OrderBookStore
	blockchainStore storeblockchain.BlockchainStore
}

func NewEvmSvc(obStore store.OrderBookStore, bcStore storeblockchain.BlockchainStore) (*EvmClient, error) {

	if obStore == nil {
		return nil, errors.New("obStore is nil")
	}

	if bcStore == nil {
		return nil, errors.New("bcStore is nil")
	}

	return &EvmClient{orderBookStore: obStore, blockchainStore: bcStore}, nil
}
