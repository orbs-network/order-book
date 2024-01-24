package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Swap(t *testing.T) {

	orderId := uuid.Must(uuid.Parse("97522a1a-7648-4b4f-97b2-38d90a5e2bd0"))
	var frags []OrderFrag
	for i := 0; i < 10; i++ {
		sz := int32(1000 * (i + 1))
		frags = append(frags, OrderFrag{
			OrderId: orderId,
			OutSize: decimal.NewFromInt32(sz),
		})
	}

	t.Run("happy path", func(t *testing.T) {
		var expected = `[{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"1000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"2000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"3000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"4000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"5000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"6000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"7000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"8000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"9000"},{"orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","size":"10000"}]`
		res, err := MarshalOrderFrags(frags)
		assert.NoError(t, err)
		assert.Equal(t, string(res), expected)
	})
}
