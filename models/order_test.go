package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var id = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var userId = uuid.MustParse("00000000-0000-0000-0000-000000000002")

func TestOrder_MapToOrder(t *testing.T) {
	order := Order{}

	t.Run("when all data is provided", func(t *testing.T) {
		data := map[string]string{
			"id":        id.String(),
			"userId":    userId.String(),
			"price":     "10.99",
			"symbol":    "USDC-ETH",
			"size":      "1000",
			"signature": "signature",
			"status":    "OPEN",
			"side":      "buy",
		}

		err := order.MapToOrder(data)

		priceDec, _ := decimal.NewFromString("10.99")
		sizeDec, _ := decimal.NewFromString("1000")

		assert.NoError(t, err)
		assert.Equal(t, data["id"], order.Id.String())
		assert.Equal(t, data["userId"], order.UserId.String())
		assert.Equal(t, priceDec, order.Price)
		assert.Equal(t, "USDC-ETH", order.Symbol.String())
		assert.Equal(t, sizeDec, order.Size)
		assert.Equal(t, "signature", order.Signature)
		assert.Equal(t, "OPEN", order.Status.String())
		assert.Equal(t, "buy", order.Side.String())
	})

	t.Run("when some data is missing", func(t *testing.T) {
		data := map[string]string{
			"id":        uuid.New().String(),
			"userId":    uuid.New().String(),
			"price":     "10.0",
			"size":      "42343324",
			"signature": "signature",
			"status":    "OPEN",
			"side":      "buy",
		}

		err := order.MapToOrder(data)
		assert.Error(t, err)
	})

	t.Run("when some data is invalid", func(t *testing.T) {
		data := map[string]string{
			"id":        "invalid-uuid",
			"userId":    "invalid-uuid",
			"price":     "invalid-decimal",
			"symbol":    "invalid-symbol",
			"size":      "invalid-decimal",
			"signature": "signature",
			"status":    "invalid-status",
			"side":      "invalid-side",
		}

		err := order.MapToOrder(data)
		assert.Error(t, err)
	})
}
