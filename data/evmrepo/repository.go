package evmrepo

import "github.com/ethereum/go-ethereum/ethclient"

type evmRepository struct {
	client ethclient.Client
}

func NewEvmRepository(client ethclient.Client) (*evmRepository, error) {
	return &evmRepository{
		client: client,
	}, nil
}
