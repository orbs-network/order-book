package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestSupportedTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("ReadFileError", func(t *testing.T) {
		filePath := "nonexistent-file.json"
		tokens := service.NewSupportedTokens(ctx, filePath)
		assert.Nil(t, tokens)

	})

	t.Run("UnmarshalError", func(t *testing.T) {
		filePath := "invalid-file.json"
		createInvalidFile(filePath)
		defer deleteFile(filePath)

		tokens := service.NewSupportedTokens(ctx, filePath)
		assert.Nil(t, tokens)
	})

	t.Run("Success", func(t *testing.T) {
		filePath := "valid-file.json"

		inText := `{"0XBTC":{"address":"0x71b821aa52a49f32eed535fca6eb5aa130085978","decimals":8}}`
		expectedOutText := fmt.Sprintf(`{"tokens":%s}`, inText)
		createValidFile(filePath, []byte(inText))
		defer deleteFile(filePath)

		st := service.NewSupportedTokens(ctx, filePath)

		assert.NotNil(t, st)
		json, err := st.AsJson()
		assert.NoError(t, err)

		outText := string(json[:])
		t.Logf("out text:  %s", outText)
		t.Logf("expected:  %s", expectedOutText)
		assert.Equal(t, expectedOutText, outText)
	})

}

func createInvalidFile(filePath string) {
	_ = os.WriteFile(filePath, []byte("invalid-json"), 0644)
}

func createValidFile(filePath string, data []byte) {
	_ = os.WriteFile(filePath, data, 0644)
}

func deleteFile(filePath string) {
	_ = os.Remove(filePath)
}
