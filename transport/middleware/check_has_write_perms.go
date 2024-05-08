package middleware

import (
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// CheckUserHasPermsMiddleware checks if the user has the correct permissions to access an endpoint.
//
// `allowedUserTypes“ is a list of user types that are allowed to access the endpoint.
func CheckUserHasPermsMiddleware(allowedUserTypes []models.UserType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := utils.GetUserCtx(r.Context())
			if user == nil {
				logctx.Error(r.Context(), "user should be in context")
				restutils.WriteJSONError(r.Context(), w, http.StatusUnauthorized, "User not found")
				return
			}

			found := false
			for _, u := range allowedUserTypes {
				if u == user.Type {
					found = true
					break
				}
			}

			if !found {
				logctx.Debug(r.Context(), "user does not have correct permissions", logger.String("userType", user.Type.String()), logger.String("userId", user.Id.String()))
				restutils.WriteJSONError(r.Context(), w, http.StatusForbidden, "user does not have correct permissions")
				return
			}

			logctx.Debug(r.Context(), "user has correct permissions", logger.String("userType", user.Type.String()), logger.String("userId", user.Id.String()))
			next.ServeHTTP(w, r)
		})
	}
}
