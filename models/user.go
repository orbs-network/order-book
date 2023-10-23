package models

import (
	"github.com/google/uuid"
)

type UserType string

const (
	MARKET_MAKER  UserType = "MARKET_MAKER"
	LIQUIDITY_HUB UserType = "LIQUIDITY_HUB"
	ADMIN         UserType = "ADMIN"
)

type User struct {
	ID   uuid.UUID
	Type UserType
}
