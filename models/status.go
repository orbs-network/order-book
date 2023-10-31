package models

import "errors"

type Status string

const (
	STATUS_OPEN   Status = "OPEN"
	STATUS_FILLED Status = "FILLED"
)

func (s Status) String() string {
	return string(s)
}

var ErrInvalidStatus = errors.New("invalid status")

func StrToStatus(s string) (Status, error) {
	switch s {
	case "OPEN":
		return STATUS_OPEN, nil
	case "FILLED":
		return STATUS_FILLED, nil
	default:
		return "", ErrInvalidStatus
	}
}
