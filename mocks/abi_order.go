package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/orbs-network/order-book/abi"
)

var BigInt = big.Int{}

var AbiFragment = abi.Order{
	Info: abi.Info{
		Reactor:                      common.Address{},
		Swapper:                      common.Address{},
		Nonce:                        &BigInt,
		Deadline:                     &BigInt,
		AdditionalValidationContract: common.Address{},
	},
	ExclusiveFiller:        common.Address{},
	ExclusivityOverrideBps: &BigInt,
	Input: abi.Input{
		Token:  common.Address{},
		Amount: &BigInt,
	},
	Outputs: []abi.Output{
		{
			Token:     common.Address{},
			Amount:    &BigInt,
			Recipient: common.Address{},
		},
	},
}

var MsgData = map[string]interface{}{
	"witness": map[string]interface{}{
		"info": map[string]interface{}{
			"reactor":                      "0x0B94c1A3E11F8aaA25D27cAf8DD05818e6f2Ad97",
			"swapper":                      "0x8fd379246834eac74B8419FfdA202CF8051F7A03",
			"nonce":                        "1000",
			"deadline":                     "1709071200",
			"additionalValidationContract": "0x0000000000000000000000000000000000000000",
			"additionalValidationData":     "0x",
		},
		"exclusiveFiller":        "0x1a08D64Fb4a7D0b6DA5606A1e4619c147C3fB95e",
		"exclusivityOverrideBps": "0",
		"input": map[string]interface{}{
			"token":  "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270",
			"amount": "40000000000000000000",
		},
		"outputs": []interface{}{
			map[string]interface{}{
				"token":  "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
				"amount": "34600000",
			},
		},
	},
}
