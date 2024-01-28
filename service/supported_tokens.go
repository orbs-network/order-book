package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type Token struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
}

type TokenMap map[string]*Token

type SupportedTokens struct {
	Name2Token TokenMap
	Adrs2Token TokenMap
}

type GetTokensRes struct {
	Tokens TokenMap `json:"tokens"`
}

func (s *SupportedTokens) ByName(name string) *Token {
	if token, ok := s.Name2Token[strings.ToUpper(name)]; ok {
		return token
	}
	return nil
}
func (s *SupportedTokens) ByAddress(adrs string) *Token {
	if token, ok := s.Adrs2Token[strings.ToUpper(adrs)]; ok {
		return token
	}
	return nil
}

func (s *SupportedTokens) AsJson() ([]byte, error) {
	res := GetTokensRes{
		Tokens: s.Name2Token,
	}

	return json.Marshal(res)
}

func NewSupportedTokens(ctx context.Context, filePath string) *SupportedTokens {
	name2Token, err := loadSupportedTokens(ctx, filePath)
	if err != nil {
		return nil
	}
	adrs2Token := TokenMap{}
	// set name fields in token, and assign to address map
	for name, token := range name2Token {
		// set name field
		token.Name = strings.ToUpper(name)
		// add address entry
		adrs2Token[strings.ToUpper(token.Address)] = token
	}

	return &SupportedTokens{
		Name2Token: name2Token,
		Adrs2Token: adrs2Token,
	}
}
func loadSupportedTokens(ctx context.Context, filePath string) (TokenMap, error) {
	file, err := os.ReadFile(filePath)

	if err != nil {
		logctx.Error(ctx, "failed to read supported tokens file", logger.Error(err), logger.String("file-path", filePath))
		return nil, fmt.Errorf("failed to read supported tokens file: %s", err)
	}

	var st TokenMap
	err = json.Unmarshal(file, &st)
	if err != nil {
		logctx.Error(ctx, "failed to unmarshal supported tokens file", logger.Error(err), logger.String("file-path", filePath))
		return nil, fmt.Errorf("failed to unmarshal supported tokens file: %s", err)
	}

	return st, nil
}
