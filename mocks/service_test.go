package mocks

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestService_MockOrderBookStore(t *testing.T) {
	ctx := context.Background()

	t.Run("MockOrderBookStore- AddVal2Set", func(t *testing.T) {
		store := MockOrderBookStore{
			Error: nil,
			Sets:  make(map[string]map[string]struct{}),
		}
		err := store.AddVal2Set(ctx, "setA", "val1")
		assert.NoError(t, err)
		err = store.AddVal2Set(ctx, "setA", "val2")
		assert.NoError(t, err)
		err = store.AddVal2Set(ctx, "setB", "val1")
		assert.NoError(t, err)
		err = store.AddVal2Set(ctx, "setA", "val1")
		assert.Equal(t, err, models.ErrValAlreadyInSet)

	})
}
