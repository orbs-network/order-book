package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type EthereumClient struct{}

type VerifySignatureInput struct {
	// The message data map that was signed
	MessageData map[string]interface{}
	// The public key of the user that purportedly signed the message
	PublicKey string
	// The signature of the message
	Signature string
}

// Returns true if the signature is valid (the `PublicKey` matches the recovered one from the `Signature`), false otherwise
//
// https://blog.hook.xyz/validate-eip-712/
func (e *EthereumClient) VerifySignature(ctx context.Context, input VerifySignatureInput) (bool, error) {
	// Prepend "04" to the public key to ensure it's in the uncompressed format
	fullPubKey := "04" + strings.TrimPrefix(input.PublicKey, "0x")

	// Decode the hex-encoded public key
	pubKeyBytes, err := hex.DecodeString(fullPubKey)
	if err != nil {
		logctx.Error(ctx, "error decoding hex public key", logger.Error(err), logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("error decoding hex public key: %v", err)
	}

	// Validate the format and length of the public key
	if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 4 {
		logctx.Error(ctx, "invalid public key format", logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("invalid public key format")
	}

	// Convert the byte representation of the public key to an ECDSA public key
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		logctx.Error(ctx, "failed to unmarshal public key", logger.Error(err), logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	// Decode the hex-encoded signature
	signatureBytes, err := hex.DecodeString(strings.TrimPrefix(input.Signature, "0x"))
	if err != nil {
		logctx.Error(ctx, "error decoding hex signature", logger.Error(err), logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Normalize the `v` value in the signature (adjust for Ethereum's signature format)
	v := signatureBytes[64]
	if v == 27 || v == 28 {
		logctx.Info(ctx, "signature v value is normalized", logger.String("publicKey", fullPubKey))
		v -= 27
		signatureBytes[64] = v
	}

	// EIP712 domain
	domain := apitypes.TypedDataDomain{
		Name:              "Permit2",
		ChainId:           math.NewHexOrDecimal256(137),
		VerifyingContract: "0x000000000022d473030f116ddee9f6b43ac78ba3",
	}

	// TODO: confirm with Zlotin rePermit payload
	// EIP712 message types
	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"PermitWitnessTransferFrom": {
			{Name: "permitted", Type: "TokenPermissions"},
			{Name: "spender", Type: "address"},
			{Name: "nonce", Type: "uint256"},
			{Name: "deadline", Type: "uint256"},
			{Name: "witness", Type: "PartialOrder"},
		},
		"TokenPermissions": {
			{Name: "token", Type: "address"},
			{Name: "amount", Type: "uint256"},
		},
		"PartialOrder": {
			{Name: "info", Type: "OrderInfo"},
			{Name: "exclusiveFiller", Type: "address"},
			{Name: "exclusivityOverrideBps", Type: "uint256"},
			{Name: "input", Type: "PartialInput"},
			{Name: "outputs", Type: "PartialOutput[]"},
		},
		"OrderInfo": {
			{Name: "reactor", Type: "address"},
			{Name: "swapper", Type: "address"},
			{Name: "nonce", Type: "uint256"},
			{Name: "deadline", Type: "uint256"},
			{Name: "additionalValidationContract", Type: "address"},
			{Name: "additionalValidationData", Type: "bytes"},
		},
		"PartialInput": {
			{Name: "token", Type: "address"},
			{Name: "amount", Type: "uint256"},
		},
		"PartialOutput": {
			{Name: "token", Type: "address"},
			{Name: "amount", Type: "uint256"},
			{Name: "recipient", Type: "address"},
		},
	}

	// Create the TypedData object
	typedData := apitypes.TypedData{
		PrimaryType: "PermitWitnessTransferFrom",
		Types:       types,
		Domain:      domain,
		Message:     input.MessageData,
	}

	// Hash the message data
	dataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		logctx.Error(ctx, "failed to hash structured data", logger.Error(err), logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("failed to hash structured data: %v", err)
	}

	// Hash the domain separator
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		logctx.Error(ctx, "failed to hash domain separator", logger.Error(err), logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("failed to hash domain separator: %v", err)
	}

	// Reconstruct the exact message that was signed - concatenate the domain separator and the hash of the data
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(dataHash)))

	// Compute the Keccak256 hash of the final message
	hashBytes := crypto.Keccak256(rawData)
	hash := common.BytesToHash(hashBytes)

	// Recover the public key from the signature
	recoveredPub, err := crypto.SigToPub(hash.Bytes(), signatureBytes)
	if err != nil {
		logctx.Error(ctx, "failed to recover public key from signature", logger.Error(err), logger.String("publicKey", fullPubKey))
		return false, fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	// Convert the recovered public key to bytes
	recoveredPubBytes := crypto.FromECDSAPub(recoveredPub)

	// Convert the original public key to bytes
	originalPubBytes := crypto.FromECDSAPub(pubKey)

	// Compare the recovered public key with the original public key
	if !bytes.Equal(recoveredPubBytes, originalPubBytes) {
		logctx.Warn(ctx, "signature does not match", logger.String("recoveredPub", hex.EncodeToString(recoveredPubBytes)), logger.String("originalPub", hex.EncodeToString(originalPubBytes)))
		return false, fmt.Errorf("signature does not match")
	}

	logctx.Info(ctx, "signature is valid", logger.String("recoveredPub", hex.EncodeToString(recoveredPubBytes)), logger.String("originalPub", hex.EncodeToString(originalPubBytes)))
	// Signature is valid
	return true, nil
}
