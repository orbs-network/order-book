package models

type OrderIter interface {
	HasNext() bool
	Next() *Order
}
