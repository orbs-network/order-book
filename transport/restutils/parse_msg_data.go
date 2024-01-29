package restutils

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/orbs-network/order-book/models"
)

func ConvertToAbiFragment(data map[string]interface{}) (models.AbiFragment, error) {
	var af models.AbiFragment

	witness, ok := data["witness"].(map[string]interface{})
	if !ok {
		return af, fmt.Errorf("witness is not a valid object")
	}

	info, ok := witness["info"].(map[string]interface{})
	if !ok {
		return af, fmt.Errorf("info is not a valid object")
	}

	orderInfo, err := parseOrderInfo(info)
	if err != nil {
		return af, err
	}

	af.Info = orderInfo

	exclusiveFiller, ok := witness["exclusiveFiller"].(string)
	if !ok {
		return af, fmt.Errorf("exclusiveFiller is not a valid string")
	}
	af.ExclusiveFiller = common.HexToAddress(exclusiveFiller)

	exclusivityOverrideBps, ok := witness["exclusivityOverrideBps"].(string)
	if !ok {
		return af, fmt.Errorf("exclusivityOverrideBps is not a valid string")
	}
	exclusivityOverrideBpsBI, err := hexToBigInt(exclusivityOverrideBps)
	if err != nil {
		return af, err
	}
	af.ExclusivityOverrideBps = exclusivityOverrideBpsBI

	input, ok := witness["input"].(map[string]interface{})
	if !ok {
		return af, fmt.Errorf("input is not a valid object")
	}
	af.Input, err = parsePartialInput(input)
	if err != nil {
		return af, err
	}

	outputs, ok := witness["outputs"].([]interface{})
	if !ok {
		return af, fmt.Errorf("outputs is not a valid array")
	}

	for _, o := range outputs {
		outputMap, ok := o.(map[string]interface{})
		if !ok {
			return af, fmt.Errorf("output is not a valid object")
		}
		output, err := parsePartialOutput(outputMap)
		if err != nil {
			return af, err
		}
		af.Outputs = append(af.Outputs, output)
	}

	return af, nil
}

func parseOrderInfo(data map[string]interface{}) (models.OrderInfo, error) {
	var oi models.OrderInfo
	// Similar validation and parsing for OrderInfo fields...
	// Example:
	reactor, ok := data["reactor"].(string)
	if !ok {
		return oi, fmt.Errorf("reactor is not a valid string")
	}
	oi.Reactor = common.HexToAddress(reactor)

	swapper, ok := data["swapper"].(string)
	if !ok {
		return oi, fmt.Errorf("swapper is not a valid string")
	}
	oi.Swapper = common.HexToAddress(swapper)

	nonce, ok := data["nonce"].(string)
	if !ok {
		return oi, fmt.Errorf("nonce is not a valid string")
	}

	nonceBI, err := stringToBigInt(nonce)
	if err != nil {
		return oi, err
	}
	oi.Nonce = nonceBI

	deadline, ok := data["deadline"].(string)
	if !ok {
		return oi, fmt.Errorf("deadline is not a valid string")
	}
	deadlineBI, err := stringToBigInt(deadline)
	if err != nil {
		return oi, err
	}
	oi.Deadline = deadlineBI

	additionalValidationContract, ok := data["additionalValidationContract"].(string)
	if !ok {
		return oi, fmt.Errorf("additionalValidationContract is not a valid string")
	}
	oi.AdditionalValidationContract = common.HexToAddress(additionalValidationContract)

	additionalValidationData, ok := data["additionalValidationData"].(string)
	if !ok {
		return oi, fmt.Errorf("additionalValidationData is not a valid string")
	}
	oi.AdditionalValidationData = []byte(additionalValidationData)

	return oi, nil
}

func parsePartialInput(data map[string]interface{}) (models.PartialInput, error) {
	var pi models.PartialInput

	token, ok := data["token"].(string)
	if !ok {
		return pi, fmt.Errorf("token is not a valid string")
	}
	pi.Token = common.HexToAddress(token)

	amount, ok := data["amount"].(string)
	if !ok {
		return pi, fmt.Errorf("amount is not a valid string")
	}

	amountBI, err := stringToBigInt(amount)
	if err != nil {
		return pi, err
	}

	pi.Amount = amountBI

	return pi, nil
}

func parsePartialOutput(data map[string]interface{}) (models.PartialOutput, error) {
	var po models.PartialOutput

	token, ok := data["token"].(string)
	if !ok {
		return po, fmt.Errorf("token is not a valid string")
	}
	po.Token = common.HexToAddress(token)

	amount, ok := data["amount"].(string)
	if !ok {
		return po, fmt.Errorf("amount is not a valid string")
	}

	amountBI, err := stringToBigInt(amount)
	if err != nil {
		return po, err
	}

	po.Amount = amountBI

	return po, nil
}

func hexToBigInt(hexStr string) (*big.Int, error) {
	value := new(big.Int)
	_, ok := value.SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert hex string to big.Int")
	}
	return value, nil
}

func stringToBigInt(str string) (*big.Int, error) {
	value := new(big.Int)
	if _, ok := value.SetString(str, 10); !ok {
		return nil, fmt.Errorf("failed to convert string to big.Int")
	}
	return value, nil
}
