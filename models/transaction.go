package models

import (
	"time"
)

type Status string

const (
	TX_SUCCESS Status = "success"
	TX_FAILURE Status = "failure"
	TX_PENDING Status = "pending"
)

type Tx struct {
	Status Status
	TxHash string
	// When tx still pending, nil block and timestamp
	Block     *int64
	Timestamp *time.Time
}
