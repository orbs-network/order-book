package models

import "context"

type OrderIter interface {
	HasNext() bool
	Next(ctx context.Context) *Order
}
