package mocks

import (
	"context"

	"github.com/orbs-network/order-book/utils"
)

func AddPaginationToCtx(page, pageSize int) context.Context {
	c := context.Background()

	pg := &utils.Paginator{Page: page, PageSize: pageSize}

	return utils.WithPaginationCtx(c, pg)
}
