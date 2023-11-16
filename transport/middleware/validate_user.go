package middleware

import (
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func ValidateUserMiddleware(s service.OrderBookService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			publicKey := r.Header.Get("X-Public-Key")
			if publicKey == "" {
				logctx.Warn(r.Context(), "missing public key header")
				http.Error(w, "Missing public key", http.StatusBadRequest)
				return
			}

			user, err := s.GetUserByPublicKey(r.Context(), publicKey)
			if err != nil {
				if err == models.ErrUserNotFound {
					logctx.Warn(r.Context(), "user not found", logger.String("publicKey", publicKey))
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				} else {
					logctx.Error(r.Context(), "unexpected error getting user by public key", logger.Error(err), logger.String("publicKey", publicKey))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}

			ctx := utils.WithUserCtx(r.Context(), user)

			logctx.Info(ctx, "found user by public key", logger.String("publicKey", publicKey), logger.String("userId", user.Id.String()))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
