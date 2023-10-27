package rest

import (
	"net/http"

	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// ExtractPkMiddleware extracts the public key from the X-Public-Key header and adds it to the context
func ExtractPkMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		publicKey := r.Header.Get("X-Public-Key")
		if publicKey == "" {
			logctx.Warn(r.Context(), "missing public key header")
			http.Error(w, "Missing public key", http.StatusBadRequest)
			return
		}

		ctx := utils.WithPkCtx(r.Context(), publicKey)

		logctx.Info(ctx, "found public key header", logger.String("publicKey", publicKey))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}