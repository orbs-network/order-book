package abi

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const expectedAbi = `0x0d7a16c3000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000002e00000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000022b1c8c1227a0000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000896d9b9eee18f6c88c5575b7824783402937557500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d500b1d8e8ef31e21c99d1db9a6444d3adf12700000000000000000000000000000000000000000000000022b1c8c1227a0000000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000002ee46d8d20020520d5266f3cacc7c41e1aadd4c600000000000000000000000092fdb5485c5eacaa3c93d68196070820e2253d3100000000000000000000000000000000000000000000000000000000000003e80000000000000000000000000000000000000000000000000000000065de5b60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000003c499c542cef5e3811e1192ce70d8cc03d5c3359000000000000000000000000000000000000000000000001e02be4ae6c8400000000000000000000000000004a9d6b0b19cbffcb0255550661ecb7014283c60e000000000000000000000000000000000000000000000000000000000000008430783434313635643331366465666266373734646631643135303137343632323031383235373364663536646534343931646534626263303465386236333539346434383065663134363438373261363162353231613331383239393538636337626539363631353632313830663737306433303539353362356235383134616232316300000000000000000000000000000000000000000000000000000000`

const reactor = "0x2ee46d8d20020520d5266f3cacc7c41e1aadd4c6"
const swapper = "0x92fdb5485c5eacaa3c93d68196070820e2253d31"
const filler = "0x896d9b9eee18f6c88c5575b78247834029375575"
const inToken = "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"
const outToken = "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359"

func TestOrder_EncodeFragData(t *testing.T) {

	// Example usage
	info := Info{
		Reactor:                      common.HexToAddress(reactor),
		Swapper:                      common.HexToAddress(swapper),
		Nonce:                        big.NewInt(1000),
		Deadline:                     big.NewInt(1709071200),
		AdditionalValidationContract: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		AdditionalValidationData:     []byte{},
	}

	inAmount := big.NewInt(0)
	inAmount.SetString("40000000000000000000", 10)

	outAmount := big.NewInt(0)
	outAmount.SetString("34600000000000000000", 10)

	input := Input{
		Token:  common.HexToAddress(inToken),
		Amount: inAmount,
	}

	output := Output{
		Token:     common.HexToAddress(outToken),
		Amount:    outAmount,
		Recipient: common.HexToAddress("0x4A9D6b0b19CBFfCB0255550661eCB7014283c60E"),
	}

	order := Order{
		Info:                   info,
		ExclusiveFiller:        common.HexToAddress(filler),
		ExclusivityOverrideBps: big.NewInt(0),
		Input:                  input,
		Outputs:                []Output{output},
	}

	orderWithAmount := OrderWithAmount{
		Order:  order,
		Amount: inAmount,
	}

	sigSample := "0x44165d316defbf774df1d1501746220182573df56de4491de4bbc04e8b63594d480ef1464872a61b521a31829958cc7be9661562180f770d305953b5b5814ab21c"
	signedOrder := SignedOrder{
		OrderWithAmount: orderWithAmount,
		Signature:       []byte(sigSample),
	}

	signedOrders := []SignedOrder{
		signedOrder,
		//signedOrder,
	}

	packedData, err := PackSignedOrders(context.Background(), signedOrders)
	assert.NoError(t, err)
	abi := fmt.Sprintf("0x%x", packedData)
	assert.Equal(t, expectedAbi, abi)

}
