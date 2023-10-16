package models

import "errors"

type Status string

const (
	STATUS_OPEN     Status = "OPEN"
	STATUS_PENDING  Status = "PENDING"
	STATUS_FILLED   Status = "FILLED"
	STATUS_CANCELED Status = "CANCELED"
)

func (s Status) String() string {
	return string(s)
}

var ErrInvalidStatus = errors.New("invalid status")

func StrToStatus(s string) (Status, error) {
	switch s {
	case "OPEN":
		return STATUS_OPEN, nil
	case "PENDING":
		return STATUS_PENDING, nil
	case "FILLED":
		return STATUS_FILLED, nil
	case "CANCELED":
		return STATUS_CANCELED, nil
	default:
		return "", ErrInvalidStatus
	}
}
