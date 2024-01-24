package rest

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

func getArguments() (abi.Arguments, error) {
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
		{Name: "ExclusivityOverride", Type: "uint256"},
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
		log.Fatalf("Failed to create ABI type for PartialOrder: %v", err)
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
func encodeFragData(frag AbiFragment) string {
	args, _ := getArguments()

	//Encode the data
	encodedData, err := args.Pack(frag)
	if err != nil {
		log.Fatalf("Failed to ABI encode frag: %v", err)
	}

	fmt.Printf("Encoded Data (hex): %x\n", encodedData)
	fmt.Printf("Length: %d\n", len(encodedData))
	return fmt.Sprintf("%x", encodedData)
}
