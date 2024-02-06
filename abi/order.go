package abi

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const executeBatchJson = `[{
	"inputs": [
		{
			"components": [
				{
					"internalType": "bytes",
					"name": "order",
					"type": "bytes"
				},
				{
					"internalType": "bytes",
					"name": "sig",
					"type": "bytes"
				}
			],
			"name": "orders",
			"type": "tuple[]"
		}
	],
	"name": "executeBatch",
	"outputs": [],
	"stateMutability": "payable",
	"type": "function"
}]`

type Info struct {
	Reactor                      common.Address
	Swapper                      common.Address
	Nonce                        *big.Int
	Deadline                     *big.Int
	AdditionalValidationContract common.Address
	AdditionalValidationData     []byte
}

type Input struct {
	Token  common.Address
	Amount *big.Int
}

type Output struct {
	Token     common.Address
	Amount    *big.Int
	Recipient common.Address
}

type Order struct {
	Info                   Info
	ExclusiveFiller        common.Address
	ExclusivityOverrideBps *big.Int
	Input                  Input
	Outputs                []Output
}
type OrderWithAmount struct {
	Order  Order
	Amount *big.Int
}

type SignedOrder struct {
	OrderWithAmount OrderWithAmount
	Signature       []byte
}

type ExecuteBatchTuple struct {
	Order []byte
	Sig   []byte
}

func PackSignedOrders(signedOrders []SignedOrder) ([]byte, error) {
	orderType, err := abi.NewType(
		"tuple", "OrderWithAmount", []abi.ArgumentMarshaling{
			{Name: "Order", Type: "tuple", Components: []abi.ArgumentMarshaling{
				{Name: "Info", Type: "tuple", Components: []abi.ArgumentMarshaling{
					{Name: "Reactor", Type: "address"},
					{Name: "Swapper", Type: "address"},
					{Name: "Nonce", Type: "uint256"},
					{Name: "Deadline", Type: "uint256"},
					{Name: "AdditionalValidationContract", Type: "address"},
					{Name: "AdditionalValidationData", Type: "bytes"},
				}},
				{Name: "ExclusiveFiller", Type: "address"},
				{Name: "ExclusivityOverrideBps", Type: "uint256"},
				{Name: "Input", Type: "tuple", Components: []abi.ArgumentMarshaling{
					{Name: "Token", Type: "address"},
					{Name: "Amount", Type: "uint256"},
				}},
				{Name: "Outputs", Type: "tuple[]", Components: []abi.ArgumentMarshaling{
					{Name: "Token", Type: "address"},
					{Name: "Amount", Type: "uint256"},
					{Name: "Recipient", Type: "address"},
				}},
			}},
			{Name: "Amount", Type: "uint256"},
		})

	if err != nil {
		return nil, err
	}

	orderArgs := abi.Arguments{
		{
			Type: orderType,
			Name: "OrderWithAmount",
		},
	}

	tuples := []ExecuteBatchTuple{}
	for _, order := range signedOrders {

		bytesOrderWithAmount, err := orderArgs.Pack(order.OrderWithAmount)
		if err != nil {
			return nil, err
		}

		tuples = append(tuples, ExecuteBatchTuple{
			Order: bytesOrderWithAmount,
			Sig:   order.Signature,
		})
	}

	execBatchAbi, err := abi.JSON(strings.NewReader(executeBatchJson))
	if err != nil {
		return nil, err
	}

	data, err := execBatchAbi.Pack("executeBatch", tuples)

	if err != nil {
		fmt.Println("Error packing data for executeBatch function:", err)
		return nil, err
	}

	return data, nil
}
