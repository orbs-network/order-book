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
			InSize:  decimal.Zero,
		})
	}

	t.Run("happy path", func(t *testing.T) {
		var expected = `[{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"1000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"2000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"3000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"4000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"5000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"6000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"7000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"8000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"9000"},{"inSize":"0","orderId":"97522a1a-7648-4b4f-97b2-38d90a5e2bd0","outSize":"10000"}]`

		res, err := MarshalOrderFrags(frags)
		assert.NoError(t, err)
		assert.Equal(t, string(res), expected)
	})
}
