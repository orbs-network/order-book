package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type GetUserByApiKeyFunc func(ctx context.Context, apiKey string) (*models.User, error)

// ValidateUserMiddleware validates the user by the API key in the request header
func ValidateUserMiddleware(getUserByApiKey GetUserByApiKeyFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			key, err := BearerToken(r, "X-API-KEY")

			if err != nil {
				logctx.Warn(r.Context(), "incorrect API key format", logger.Error(err))
				restutils.WriteJSONError(r.Context(), w, http.StatusBadRequest, "Invalid API key (ensure the format is 'Bearer YOUR-API-KEY')")
				return
			}

			user, err := getUserByApiKey(r.Context(), key)
			if err != nil {
				if err == models.ErrNotFound {
					logctx.Warn(r.Context(), "user not found by api key")
					restutils.WriteJSONError(r.Context(), w, http.StatusUnauthorized, "Unauthorized")
				} else {
					logctx.Error(r.Context(), "unexpected error getting user by api key", logger.Error(err))
					restutils.WriteJSONError(r.Context(), w, http.StatusInternalServerError, "Internal server error")
				}
				return
			}

			ctx := utils.WithUserCtx(r.Context(), user)

			logctx.Debug(ctx, "found user by api key", logger.String("userId", user.Id.String()))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BearerToken(r *http.Request, header string) (string, error) {
	rawToken := r.Header.Get(header)
	pieces := strings.SplitN(rawToken, " ", 2)

	if len(pieces) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}

	token := strings.TrimSpace(pieces[1])

	return token, nil
}
