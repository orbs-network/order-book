// Pub/Sub events

package models

import (
	"fmt"

	"github.com/google/uuid"
)

func CreateUserOrdersEventKey(userId uuid.UUID) string {
	return fmt.Sprintf("user_orders:%s", userId)
}
