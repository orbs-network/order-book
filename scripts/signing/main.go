// A quick script to test EIP712 signature verification
// Usage: go run scripts/signing/main.go

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/orbs-network/order-book/service"
)

var sig = "0x977939edfd70b5e0a802b0874fa681749ac3669516f8250fc11426bd3d02e3f353a874aec1d1394965507c7404a24a53cc4427917c2b214a695ea432aa79ccba1b"

var pubKey = "0x6a04ab98d9e4774ad806e302dddeb63bea16b5cb5f223ee77478e861bb583eb336b6fbcb60b5b3d4f1551ac45e5ffc4936466e7d98f6c7c0ec736539f74691a6"

var message = map[string]interface{}{
	"permitted": map[string]interface{}{
		"token":  "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
		"amount": "20000000000",
	},
	"spender":  "0x21Da9737764527e75C17F1AB26Cb668b66dEE0a0",
	"nonce":    "845753781",
	"deadline": "1709657651",
	"witness": map[string]interface{}{
		"info": map[string]interface{}{
			"reactor":                      "0x21Da9737764527e75C17F1AB26Cb668b66dEE0a0",
			"swapper":                      "0xE3682CCecefBb3C3fe524BbFF1598B2BBaC0d6E3",
			"nonce":                        "845753781",
			"deadline":                     "1709657651",
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
				"recipient": "0x8fd379246834eac74B8419FfdA202CF8051F7A03",
			},
		},
	},
}

func main() {

	client := service.EvmClient{}

	ctx := context.Background()

	verified, err := client.VerifySignature(ctx, service.VerifySignatureInput{
		MessageData: message,
		Signature:   sig,
		PublicKey:   pubKey,
	})

	if err != nil {
		log.Fatalf("Error verifying signature: %v", err)
	}

	fmt.Printf("Signature verified: %v\n", verified)
}
