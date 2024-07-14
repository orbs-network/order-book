package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_PairMngr(t *testing.T) {
	m := NewPairMngr()

	t.Run("when all data is provided", func(t *testing.T) {
		p1 := m.Resolve("ETH", "USD")
		p2 := m.Resolve("USD", "ETH")
		p3 := m.Resolve("USD", "XXX")
		assert.Equal(t, p1, p2)
		assert.Nil(t, p3)
	})

	t.Run("all pair should work", func(t *testing.T) {
		symbolPairs := GetAllSymbols()
		for _, sp := range symbolPairs {
			arr := strings.Split(sp.String(), "-")
			assert.Equal(t, len(arr), 2)
			pair := m.Resolve(arr[0], arr[1])
			assert.NotNil(t, pair)
		}
	})
}
