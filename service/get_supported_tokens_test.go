package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestGetSupportedTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("ReadFileError", func(t *testing.T) {
		filePath := "nonexistent-file.json"
		tokens, err := service.GetSupportedTokens(ctx, filePath)
		assert.Error(t, err)
		assert.Nil(t, tokens)
	})

	t.Run("UnmarshalError", func(t *testing.T) {
		filePath := "invalid-file.json"
		createInvalidFile(filePath)
		defer deleteFile(filePath)

		tokens, err := service.GetSupportedTokens(ctx, filePath)
		assert.Error(t, err)
		assert.Nil(t, tokens)
	})

	t.Run("Success", func(t *testing.T) {
		filePath := "valid-file.json"
		createValidFile(filePath, []byte(`{
			"0XBTC": {
					"address": "0x71b821aa52a49f32eed535fca6eb5aa130085978",
					"decimals": 8
			}}`))
		defer deleteFile(filePath)

		tokens, err := service.GetSupportedTokens(ctx, filePath)
		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		expected := service.SupportedTokens{
			"0XBTC": {
				Address:  "0x71b821aa52a49f32eed535fca6eb5aa130085978",
				Decimals: 8,
			},
		}
		assert.Equal(t, expected, tokens)
	})

}

func createInvalidFile(filePath string) {
	os.WriteFile(filePath, []byte("invalid-json"), 0644)
}

func createValidFile(filePath string, data []byte) {
	os.WriteFile(filePath, data, 0644)
}

func deleteFile(filePath string) {
	os.Remove(filePath)
}
