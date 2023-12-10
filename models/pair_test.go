package models

import (
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
}
