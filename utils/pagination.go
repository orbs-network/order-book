package utils

import (
	"context"
)

type Paginator struct {
	Page     int
	PageSize int
}

func (p *Paginator) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Paginator) Limit() int {
	return p.PageSize
}

// NewPaginator returns a new paginator from the context
func NewPaginator(ctx context.Context) *Paginator {
	pg := GetPaginationCtx(ctx)

	if pg == nil {
		return &Paginator{Page: 1, PageSize: 10}
	}

	return &Paginator{Page: pg.Page, PageSize: pg.PageSize}
}

// PaginationBounds returns the start and stop values for a paginated query
func PaginationBounds(ctx context.Context) (start, stop int64) {
	paginator := NewPaginator(ctx)
	start = int64(paginator.Offset())
	stop = start + int64(paginator.Limit())
	return
}
