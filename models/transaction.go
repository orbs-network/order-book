package models

type Status string

const (
	TX_SUCCESS Status = "success"
	TX_FAILURE Status = "failure"
	TX_PENDING Status = "pending"
)

type Tx struct {
	Status Status
	TxHash string
}
