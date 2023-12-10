package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var id = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var userId = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000003")

func TestOrder_OrderToMap(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	order := Order{
		Id:          id,
		ClientOId:   clientOId,
		UserId:      userId,
		Price:       decimal.NewFromFloat(10.99),
		Symbol:      "USDC-ETH",
		Size:        decimal.NewFromInt(1000),
		SizeFilled:  decimal.NewFromInt(600),
		SizePending: decimal.NewFromInt(400),
		Signature: Signature{
			Eip712Sig: "signature",
			Eip712MsgData: map[string]interface{}{
				"message": "data",
			},
		},
		Side:      BUY,
		Timestamp: timestamp,
	}

	eip712MsgDataStr := "{\"message\":\"data\"}"

	expectedMap := map[string]string{
		"id":            order.Id.String(),
		"clientOId":     order.ClientOId.String(),
		"userId":        order.UserId.String(),
		"price":         order.Price.String(),
		"symbol":        order.Symbol.String(),
		"size":          order.Size.String(),
		"sizePending":   order.SizePending.String(),
		"sizeFilled":    order.SizeFilled.String(),
		"side":          order.Side.String(),
		"timestamp":     order.Timestamp.Format(time.RFC3339),
		"eip712Sig":     order.Signature.Eip712Sig,
		"eip712MsgData": eip712MsgDataStr,
	}

	actualMap := order.OrderToMap()

	assert.Equal(t, expectedMap, actualMap)
}

func TestOrder_MapToOrder(t *testing.T) {
	order := Order{}

	t.Run("when all data is provided", func(t *testing.T) {
		data := map[string]string{
			"id":            id.String(),
			"clientOId":     clientOId.String(),
			"userId":        userId.String(),
			"price":         "10.99",
			"symbol":        "USDC-ETH",
			"size":          "1000",
			"sizePending":   "0",
			"sizeFilled":    "0",
			"side":          "buy",
			"timestamp":     "2021-01-01T00:00:00Z",
			"clientOrderId": id.String(),
			"eip712Sig":     "signature",
			"eip712MsgData": "{\"message\":\"data\"}",
		}

		err := order.MapToOrder(data)
		assert.NoError(t, err)

		priceDec, _ := decimal.NewFromString("10.99")
		sizeDec, _ := decimal.NewFromString("1000")

		assert.NoError(t, err)
		assert.Equal(t, data["id"], order.Id.String())
		assert.Equal(t, data["clientOId"], order.ClientOId.String())
		assert.Equal(t, data["userId"], order.UserId.String())
		assert.Equal(t, priceDec, order.Price)
		assert.Equal(t, "USDC-ETH", order.Symbol.String())
		assert.Equal(t, sizeDec, order.Size)
		assert.Equal(t, Signature{Eip712Sig: "signature", Eip712MsgData: map[string]interface{}{
			"message": "data",
		}}, order.Signature)
		assert.Equal(t, "buy", order.Side.String())
		assert.Equal(t, "2021-01-01 00:00:00 +0000 UTC", order.Timestamp.String())
	})

	t.Run("when some data is missing", func(t *testing.T) {
		data := map[string]string{
			"id":        uuid.New().String(),
			"userId":    uuid.New().String(),
			"price":     "10.0",
			"size":      "42343324",
			"signature": "signature",
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
			"side":      "invalid-side",
		}

		err := order.MapToOrder(data)
		assert.Error(t, err)
	})
}
