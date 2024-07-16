package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testABComb(t *testing.T, m *PairMngr, tkn1, tkn2 string) {
	// find side 1
	p := m.Resolve(tkn1, tkn2)
	assert.NotNil(t, p)
	// find oposit side
	p = m.Resolve(tkn2, tkn1)
	assert.NotNil(t, p)
	// fale same token2
	p = m.Resolve(tkn2, tkn2)
	assert.Nil(t, p)
	// fale same token1
	p = m.Resolve(tkn1, tkn1)
	assert.Nil(t, p)
}

func TestOrder_PairMngr(t *testing.T) {
	m := NewPairMngr()

	t.Run("when all data is provided", func(t *testing.T) {
		p1 := m.Resolve("ETH", "USD")
		p2 := m.Resolve("USD", "ETH")
		p3 := m.Resolve("USD", "XXX")
		assert.Equal(t, p1, p2)
		assert.Nil(t, p3)
	})

	t.Run("bi direction explicit", func(t *testing.T) {
		p := m.Resolve("MATIC", "USDT")
		assert.NotNil(t, p)
		p = m.Resolve("USDT", "MATIC")
		assert.NotNil(t, p)
		p = m.Resolve("USDT", "USDT")
		assert.Nil(t, p)
		p = m.Resolve("MATIC", "MATIC")
		assert.Nil(t, p)

	})

	t.Run("all pair should work", func(t *testing.T) {
		symbolPairs := GetAllSymbols()
		for _, sp := range symbolPairs {
			arr := strings.Split(sp.String(), "-")
			assert.Equal(t, len(arr), 2)
			// resolve both sides
			testABComb(t, m, arr[0], arr[1])
		}
	})
}
