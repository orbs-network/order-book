package models

import (
	"errors"

	"github.com/google/uuid"
)

type UserType string

const (
	MARKET_MAKER  UserType = "MARKET_MAKER"
	LIQUIDITY_HUB UserType = "LIQUIDITY_HUB"
	ADMIN         UserType = "ADMIN"
)

func (u UserType) String() string {
	return string(u)
}

type User struct {
	Id uuid.UUID
	// Public key
	Pk   string
	Type UserType
}

var ErrInvalidUserType = errors.New("invalid user type")

func StrToUserType(str string) (UserType, error) {
	switch str {
	case "MARKET_MAKER":
		return MARKET_MAKER, nil
	case "LIQUIDITY_HUB":
		return LIQUIDITY_HUB, nil
	case "ADMIN":
		return ADMIN, nil
	default:
		return "", ErrInvalidUserType
	}
}
