package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/orbs-network/order-book/models"
)

var BigInt = big.Int{}

var AbiFragment = models.AbiFragment{
	Info: models.OrderInfo{
		Reactor:                      common.Address{},
		Swapper:                      common.Address{},
		Nonce:                        &BigInt,
		Deadline:                     &BigInt,
		AdditionalValidationContract: common.Address{},
	},
	ExclusiveFiller:        common.Address{},
	ExclusivityOverrideBps: &BigInt,
	Input: models.PartialInput{
		Token:  common.Address{},
		Amount: &BigInt,
	},
	Outputs: []models.PartialOutput{
		{
			Token:     common.Address{},
			Amount:    &BigInt,
			Recipient: common.Address{},
		},
	},
}
