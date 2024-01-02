package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sig = "0xe040ee7d67a021e8112dd16362f57132d86157e75e72c4aee23c875ae5ee1534065437043aeecfb9da5e44ec96ad042405ec7182d8d5bb2bf0b1a57e4244a8b71c"

var pubKey = "0xc9421bf7f3625d35b517b6af2fd0049f661209437ad216d681a5801739a71d784b2b0751c6951c5f412242a5c610022dfcbbe635f6002a362f8c4c1eb0bb1383"

var message = map[string]interface{}{
	"permitted": map[string]interface{}{
		"token":  "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
		"amount": "20000000000",
	},
	"spender":  "0x21Da9737764527e75C17F1AB26Cb668b66dEE0a0",
	"nonce":    "1438370532",
	"deadline": "1711965900",
	"witness": map[string]interface{}{
		"info": map[string]interface{}{
			"reactor":                      "0x21Da9737764527e75C17F1AB26Cb668b66dEE0a0",
			"swapper":                      "0xE3682CCecefBb3C3fe524BbFF1598B2BBaC0d6E3",
			"nonce":                        "1438370532",
			"deadline":                     "1711965900",
			"additionalValidationContract": "0x1a08D64Fb4a7D0b6DA5606A1e4619c147C3fB95e",
			"additionalValidationData":     "0x",
		},
		"exclusiveFiller":        "0x1a08D64Fb4a7D0b6DA5606A1e4619c147C3fB95e",
		"exclusivityOverrideBps": "0",
		"input": map[string]interface{}{
			"token":  "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
			"amount": "20000000000",
		},
		"outputs": []interface{}{
			map[string]interface{}{
				"token":     "0x11cd37bb86f65419713f30673a480ea33c826872",
				"amount":    "10000000000000000000",
				"recipient": "0x0A84F9A73d9cF507E2028B5e9D6a842F9632BC29",
			},
		},
	},
}

func TestEvmClient_VerifySignature(t *testing.T) {
	ctx := context.TODO()
	evmClient := &EvmClient{}

	t.Run("successfully verify signature - should return true", func(t *testing.T) {
		input := VerifySignatureInput{
			PublicKey:   pubKey,
			Signature:   sig,
			MessageData: message,
		}

		result, err := evmClient.VerifySignature(ctx, input)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("invalid public key  - should return false", func(t *testing.T) {
		input := VerifySignatureInput{
			PublicKey:   "0x123",
			Signature:   sig,
			MessageData: message,
		}

		result, err := evmClient.VerifySignature(ctx, input)

		assert.ErrorContains(t, err, "error decoding hex public key")
		assert.False(t, result)
	})

}
