package mocks

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type MockBcBackend struct {
	backend *backends.SimulatedBackend
	auth    *bind.TransactOpts
}

func (m *MockBcBackend) Backend() *backends.SimulatedBackend {
	return m.backend
}

func (m *MockBcBackend) Auth() *bind.TransactOpts {
	return m.auth
}

func (m *MockBcBackend) Commit() {
	m.backend.Commit()
}

func (m *MockBcBackend) CreateTx(tx *types.Transaction, shouldMine bool) (*types.Transaction, error) {
	signedTx, err := m.auth.Signer(m.auth.From, tx)
	if err != nil {
		return nil, err
	}

	if err := m.backend.SendTransaction(context.Background(), signedTx); err != nil {
		return nil, err
	}

	if shouldMine {
		m.backend.Commit()
	}

	return signedTx, nil
}

func (m *MockBcBackend) BalanceOf(ctx context.Context, token, adrs string) (*big.Int, error) {
	return nil, nil
}

func NewMockBcBackend() *MockBcBackend {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	blockGasLimit := uint64(4712388)

	return &MockBcBackend{
		backend: backends.NewSimulatedBackend(genesisAlloc, blockGasLimit),
		auth:    auth,
	}
}
