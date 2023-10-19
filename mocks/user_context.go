package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
)

// Add user to context
func AddUserToCtx(user *models.User) context.Context {
	c := context.Background()

	if user != nil {
		return utils.WithUser(c, user)
	}

	_user := models.User{
		ID:   uuid.MustParse("a577273e-12de-4acc-a4f8-de7fb5b86e37"),
		Type: models.MARKET_MAKER,
	}

	return utils.WithUser(c, &_user)
}
