package service

import "github.com/orbs-network/order-book/models"

type OrderIter interface {
	HasNext() bool
	Next() *models.Order
}
