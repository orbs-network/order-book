package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type SupportedToken struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
}

type SupportedTokens map[string]SupportedToken

func LoadSupportedTokens(ctx context.Context, filePath string) (SupportedTokens, error) {
	file, err := os.ReadFile(filePath)

	if err != nil {
		logctx.Error(ctx, "failed to read supported tokens file", logger.Error(err), logger.String("file-path", filePath))
		return nil, fmt.Errorf("failed to read supported tokens file: %s", err)
	}

	var tokens SupportedTokens
	err = json.Unmarshal(file, &tokens)
	if err != nil {
		logctx.Error(ctx, "failed to unmarshal supported tokens file", logger.Error(err), logger.String("file-path", filePath))
		return nil, fmt.Errorf("failed to unmarshal supported tokens file: %s", err)
	}

	return tokens, nil
}
