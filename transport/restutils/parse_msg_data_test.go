package restutils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestConvertToAbiFragment(t *testing.T) {
	// Example test data
	testData := map[string]interface{}{
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

	result, err := ConvertToAbiFragment(testData)

	assert.NoError(t, err)

	// Info
	assert.Equal(t, common.HexToAddress("0x0B94c1A3E11F8aaA25D27cAf8DD05818e6f2Ad97"), result.Info.Reactor)
	assert.Equal(t, common.HexToAddress("0x8fd379246834eac74B8419FfdA202CF8051F7A03"), result.Info.Swapper)
	assert.Equal(t, big.NewInt(1000), result.Info.Nonce)
	assert.Equal(t, big.NewInt(1709071200), result.Info.Deadline)
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000000"), result.Info.AdditionalValidationContract)
	assert.Equal(t, []byte("0x"), result.Info.AdditionalValidationData)

	// ExclusiveFiller
	assert.Equal(t, common.HexToAddress("0x1a08D64Fb4a7D0b6DA5606A1e4619c147C3fB95e"), result.ExclusiveFiller)

	// ExclusivityOverrideBps
	assert.Equal(t, big.NewInt(0), result.ExclusivityOverrideBps)

	// Input
	assert.Equal(t, common.HexToAddress("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"), result.Input.Token)

	expectedAmount := new(big.Int)
	expectedAmount.SetString("40000000000000000000", 10)
	assert.Equal(t, expectedAmount, result.Input.Amount)

	// Outputs
	assert.Len(t, result.Outputs, 1)
	assert.Equal(t, common.HexToAddress("0x3c499c542cef5e3811e1192ce70d8cc03d5c3359"), result.Outputs[0].Token)
	assert.Equal(t, big.NewInt(34600000), result.Outputs[0].Amount)

}
