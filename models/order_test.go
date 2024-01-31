package models

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var id = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var userId = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000003")

var abiFragment = AbiFragment{}

func TestOrder_OrderToMap(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	order := Order{
		Id:          id,
		ClientOId:   clientOId,
		UserId:      userId,
		Price:       decimal.NewFromFloat(10.99),
		Symbol:      "MATIC-USDC",
		Size:        decimal.NewFromInt(1000),
		SizeFilled:  decimal.NewFromInt(600),
		SizePending: decimal.NewFromInt(400),
		Signature: Signature{
			Eip712Sig:   "signature",
			AbiFragment: abiFragment,
		},
		Side:      BUY,
		Timestamp: timestamp,
	}

	expectedMap := map[string]string{
		"id":          order.Id.String(),
		"clientOId":   order.ClientOId.String(),
		"userId":      order.UserId.String(),
		"price":       order.Price.String(),
		"symbol":      order.Symbol.String(),
		"size":        order.Size.String(),
		"sizePending": order.SizePending.String(),
		"sizeFilled":  order.SizeFilled.String(),
		"side":        order.Side.String(),
		"timestamp":   order.Timestamp.Format(time.RFC3339),
		"eip712Sig":   order.Signature.Eip712Sig,
		"abiFragment": "{\"Info\":{\"Reactor\":\"0x0000000000000000000000000000000000000000\",\"Swapper\":\"0x0000000000000000000000000000000000000000\",\"Nonce\":null,\"Deadline\":null,\"AdditionalValidationContract\":\"0x0000000000000000000000000000000000000000\",\"AdditionalValidationData\":null},\"ExclusiveFiller\":\"0x0000000000000000000000000000000000000000\",\"ExclusivityOverrideBps\":null,\"Input\":{\"Token\":\"0x0000000000000000000000000000000000000000\",\"Amount\":null},\"Outputs\":null}",
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
			"symbol":        "MATIC-USDC",
			"size":          "1000",
			"sizePending":   "0",
			"sizeFilled":    "0",
			"side":          "buy",
			"timestamp":     "2021-01-01T00:00:00Z",
			"clientOrderId": id.String(),
			"eip712Sig":     "signature",
			"abiFragment":   "{\"Info\":{\"Reactor\":\"0x0000000000000000000000000000000000000000\",\"Swapper\":\"0x0000000000000000000000000000000000000000\",\"Nonce\":null,\"Deadline\":null,\"AdditionalValidationContract\":\"0x0000000000000000000000000000000000000000\",\"AdditionalValidationData\":null},\"ExclusiveFiller\":\"0x0000000000000000000000000000000000000000\",\"ExclusivityOverrideBps\":null,\"Input\":{\"Token\":\"0x0000000000000000000000000000000000000000\",\"Amount\":null},\"Outputs\":null}",
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
		assert.Equal(t, "MATIC-USDC", order.Symbol.String())
		assert.Equal(t, sizeDec, order.Size)
		assert.Equal(t, Signature{Eip712Sig: "signature", AbiFragment: abiFragment}, order.Signature)
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

func TestOrder_Fill(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name     string
		order    Order
		expected Order
		fillSize decimal.Decimal
		isFilled bool
		error    error
	}{
		{
			name: "size 1000, sizeFilled 0, sizePending 1000",
			order: Order{
				Size:        decimal.NewFromInt(1000),
				SizeFilled:  decimal.NewFromInt(0),
				SizePending: decimal.NewFromFloat(1000),
			},
			expected: Order{
				Size:        decimal.NewFromInt(1000),
				SizeFilled:  decimal.NewFromInt(1000),
				SizePending: decimal.Zero,
			},
			fillSize: decimal.NewFromInt(1000),
			isFilled: true,
		},
		{
			name: "size 1000, sizeFilled 500, sizePending 500",
			order: Order{
				Size:        decimal.NewFromInt(1000),
				SizeFilled:  decimal.NewFromInt(500),
				SizePending: decimal.NewFromInt(500),
			},
			expected: Order{
				Size:        decimal.NewFromInt(1000),
				SizeFilled:  decimal.NewFromInt(1000),
				SizePending: decimal.Zero,
			},
			fillSize: decimal.NewFromInt(500),
			isFilled: true,
		},
		{
			name: "partial fill",
			order: Order{
				Size:        decimal.NewFromFloat(23782378.50),
				SizeFilled:  decimal.NewFromFloat(2.38),
				SizePending: decimal.NewFromFloat(1238.12),
			},
			expected: Order{
				Size:        decimal.NewFromFloat(23782378.50),
				SizeFilled:  decimal.NewFromFloat(2.38).Add(decimal.NewFromFloat(10)),
				SizePending: decimal.NewFromFloat(1238.12).Sub(decimal.NewFromFloat(10)),
			},
			fillSize: decimal.NewFromFloat(10),
			isFilled: false,
		},
		{
			name: "total size is less than requested fill size",
			order: Order{
				Size:        decimal.NewFromFloat(10.00),
				SizeFilled:  decimal.NewFromFloat(9.00),
				SizePending: decimal.NewFromFloat(2.89),
			},
			expected: Order{
				Size:        decimal.NewFromFloat(10.00),
				SizeFilled:  decimal.NewFromFloat(9.00),
				SizePending: decimal.NewFromFloat(2.89),
			},
			fillSize: decimal.NewFromFloat(2.00),
			isFilled: false,
			error:    ErrUnexpectedSizeFilled,
		},
		{
			name: "size to be filled is greater than size pending",
			order: Order{
				Size:        decimal.NewFromFloat(10.00),
				SizeFilled:  decimal.NewFromFloat(2.00),
				SizePending: decimal.NewFromFloat(2.00),
			},
			expected: Order{
				Size:        decimal.NewFromFloat(10.00),
				SizeFilled:  decimal.NewFromFloat(2.00),
				SizePending: decimal.NewFromFloat(2.00),
			},
			fillSize: decimal.NewFromFloat(4.00),
			isFilled: false,
			error:    ErrUnexpectedSizePending,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isFilled, err := test.order.Fill(ctx, test.fillSize)
			assert.Equal(t, test.expected.Size.String(), test.order.Size.String(), "size should be equal")
			assert.Equal(t, test.expected.SizeFilled.String(), test.order.SizeFilled.String(), "sizeFilled should be equal")
			assert.Equal(t, test.expected.SizePending.String(), test.order.SizePending.String(), "sizePending should be equal")
			assert.Equal(t, test.isFilled, isFilled, "isFilled should be equal")
			assert.Equal(t, test.error, err, "error should be equal")
		})
	}
}

func TestOrder_Unlock(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name     string
		order    Order
		expected Order
		lockSize decimal.Decimal
		error    error
	}{
		{
			name: "sizePending 1000, lockSize 1000",
			order: Order{
				SizePending: decimal.NewFromFloat(1000),
			},
			expected: Order{
				SizePending: decimal.Zero,
			},
			lockSize: decimal.NewFromInt(1000),
		},
		{
			name: "sizePending 500, lockSize 1000",
			order: Order{
				SizePending: decimal.NewFromInt(500),
			},
			expected: Order{
				SizePending: decimal.NewFromInt(500),
			},
			lockSize: decimal.NewFromInt(1000),
			error:    ErrUnexpectedSizeFilled,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.order.Unlock(ctx, test.lockSize)
			assert.Equal(t, test.expected.SizePending.String(), test.order.SizePending.String(), "sizePending should be equal")
			assert.Equal(t, test.error, err, "error should be equal")
		})
	}
}
