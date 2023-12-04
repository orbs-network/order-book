// A quick script to test EIP712 signature verification
// Usage: go run scripts/signing/main.go

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/orbs-network/order-book/service"
)

var sig = "0xf577460acb728646ceb2f561a96db6c84aeb99334bcc938e7d21fb1fa83c90f75efb5e8c64c53b98b4d1d97b693a89dbd7606225f68cddffbc075f402494b7331b"

var pubKey = "0x6a04ab98d9e4774ad806e302dddeb63bea16b5cb5f223ee77478e861bb583eb336b6fbcb60b5b3d4f1551ac45e5ffc4936466e7d98f6c7c0ec736539f74691a6"

var message = `{"permitted": {"token": "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359", "amount": 10000000}, "spender": "0x21Da9737764527e75C17F1AB26Cb668b66dEE0a0", "nonce": 2576040678, "deadline": 1709382304, "witness": {"info": {"reactor": "0x21Da9737764527e75C17F1AB26Cb668b66dEE0a0", "swapper": "0xE3682CCecefBb3C3fe524BbFF1598B2BBaC0d6E3", "nonce": 2576040678, "deadline": 1709382304, "additionalValidationContract": "0x1a08D64Fb4a7D0b6DA5606A1e4619c147C3fB95e", "additionalValidationData": "0x"}, "decayStartTime": 1709382304, "decayEndTime": 1709382304, "exclusiveFiller": "0x1a08D64Fb4a7D0b6DA5606A1e4619c147C3fB95e", "exclusivityOverrideBps": "0", "inputToken": "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359", "inputStartAmount": 10000000, "inputEndAmount": 10000000}}`

func main() {

	client := service.EthereumClient{}

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
