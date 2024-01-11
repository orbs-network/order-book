package mocks

import "github.com/orbs-network/order-book/service"

var SupportedTokens = map[string]service.SupportedToken{
	"matic": {
		Address:  "0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0",
		Decimals: 18,
	},
	"usdc": {
		Address:  "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Decimals: 6,
	},
}
