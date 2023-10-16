package models

import "errors"

type Side string

const (
	BUY  Side = "buy"
	SELL Side = "sell"
)

var ErrInvalidSide = errors.New("invalid side")

func StrToSide(s string) (Side, error) {
	switch s {
	case "buy":
		return BUY, nil
	case "sell":
		return SELL, nil
	default:
		return "", ErrInvalidSide
	}
}

func (s Side) String() string {
	return string(s)
}
