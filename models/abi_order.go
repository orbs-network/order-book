package models

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type OrderInfo struct {
	Reactor                      common.Address
	Swapper                      common.Address
	Nonce                        *big.Int
	Deadline                     *big.Int
	AdditionalValidationContract common.Address
	AdditionalValidationData     []byte
}

type PartialInput struct {
	Token  common.Address
	Amount *big.Int
}

type PartialOutput struct {
	Token     common.Address
	Amount    *big.Int
	Recipient common.Address
}

type AbiFragment struct {
	Info                   OrderInfo
	ExclusiveFiller        common.Address
	ExclusivityOverrideBps *big.Int
	Input                  PartialInput
	Outputs                []PartialOutput
}

// Custom JSON serialization for *big.Int fields
// func (o *OrderInfo) MarshalJSON() ([]byte, error) {
// 	type Alias OrderInfo
// 	return json.Marshal(&struct {
// 		Nonce    string `json:"nonce"`
// 		Deadline string `json:"deadline"`
// 		*Alias
// 	}{
// 		Nonce:    o.Nonce.String(),
// 		Deadline: o.Deadline.String(),
// 		Alias:    (*Alias)(o),
// 	})
// }

// func (p *PartialInput) MarshalJSON() ([]byte, error) {
// 	type Alias PartialInput
// 	return json.Marshal(&struct {
// 		Amount string `json:"Amount"`
// 		*Alias
// 	}{
// 		Amount: p.Amount.String(),
// 		Alias:  (*Alias)(p),
// 	})
// }

// func (p *PartialOutput) MarshalJSON() ([]byte, error) {
// 	type Alias PartialOutput
// 	return json.Marshal(&struct {
// 		Amount string `json:"Amount"`
// 		*Alias
// 	}{
// 		Amount: p.Amount.String(),
// 		Alias:  (*Alias)(p),
// 	})
// }

func getAbiArguments() (abi.Arguments, error) {
	abiFrag, err := abi.NewType("tuple", "AbiFragment", []abi.ArgumentMarshaling{
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
		},
		}},
	)
	if err != nil {
		return abi.Arguments{}, err
	}

	// Define the ABI arguments
	arguments := abi.Arguments{
		{
			Type: abiFrag,
			Name: "abiFrag",
		},
	}
	return arguments, nil
}
func EncodeFragData(ctx context.Context, frag AbiFragment) (string, error) {
	args, err := getAbiArguments()
	if err != nil {
		logctx.Error(ctx, "args.Pack failed %s", logger.Error(err))
		return "", err
	}

	//Encode the data
	encodedData, err := args.Pack(frag)
	if err != nil {
		logctx.Error(ctx, "args.Pack failed %s", logger.Error(err))
		return "", err
	}

	hexRes := fmt.Sprintf("%x", encodedData)
	logctx.Debug(ctx, "EncodeFragData", logger.String("buffer", hexRes))
	return hexRes, nil
}
