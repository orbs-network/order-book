package models

import (
	"math/big"

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
