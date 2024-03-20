package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

const DEFAULT_PAGE_SIZE int = 10

type paginationResponse[T any] struct {
	Data       T   `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func NewPaginationResponse[T any](ctx context.Context, data T, total int) paginationResponse[T] {

	paginator := utils.NewPaginator(ctx)

	return paginationResponse[T]{
		Data:       data,
		Page:       paginator.Page,
		PageSize:   paginator.PageSize,
		Total:      total,
		TotalPages: (total + paginator.PageSize - 1) / paginator.PageSize,
	}
}

func PaginationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}
		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil {
			pageSize = DEFAULT_PAGE_SIZE
		}

		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = DEFAULT_PAGE_SIZE
		}

		ctx := utils.WithPaginationCtx(r.Context(), &utils.Paginator{Page: page, PageSize: pageSize})

		logctx.Debug(ctx, "pagination middleware", logger.String("url", r.URL.Path), logger.Int("page", page), logger.Int("pageSize", pageSize))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
