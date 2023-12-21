package mocks

import (
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
)

var SwapTx = models.SwapTx{
	SwapId: uuid.MustParse("b3f9e3a0-5b7a-4b7a-8b0a-9b9b9b9b9b9b"),
	TxHash: "0x5dcbfe934287c50363e5c82502739aadd4d535a1f7c0ccd7a8088fb4dfd800da",
}
