package utils_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils_Pagination(t *testing.T) {

	t.Run("NewPaginator", func(t *testing.T) {

		ctx := context.Background()

		// Test case 1: no pagination context in context
		paginator := utils.NewPaginator(ctx)
		assert.Equal(t, 1, paginator.Page)
		assert.Equal(t, 10, paginator.PageSize)

		// Test case 2: pagination context in context
		pg := &utils.Paginator{Page: 2, PageSize: 20}
		ctx = utils.WithPaginationCtx(ctx, pg)

		paginator = utils.NewPaginator(ctx)
		assert.Equal(t, 2, paginator.Page)
		assert.Equal(t, 20, paginator.PageSize)
	})

	t.Run("PaginationBounds", func(t *testing.T) {
		ctx := context.Background()

		// Test case 1: no pagination context in context
		start, stop := utils.PaginationBounds(ctx)
		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(10), stop)

		// Test case 2: pagination context in context
		pg := &utils.Paginator{Page: 2, PageSize: 20}
		ctx = utils.WithPaginationCtx(ctx, pg)

		start, stop = utils.PaginationBounds(ctx)
		assert.Equal(t, int64(20), start)
		assert.Equal(t, int64(40), stop)
	})

	t.Run("Offset", func(t *testing.T) {
		pg := &utils.Paginator{Page: 2, PageSize: 20}
		assert.Equal(t, 20, pg.Offset())
	})

	t.Run("Limit", func(t *testing.T) {
		pg := &utils.Paginator{Page: 2, PageSize: 20}
		assert.Equal(t, 20, pg.Limit())
	})
}
