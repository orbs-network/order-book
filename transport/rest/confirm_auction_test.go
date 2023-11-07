package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/transport/rest"
	"github.com/stretchr/testify/assert"
)

func TestHandler_confirmAuction(t *testing.T) {

	res := createServer(t)
	assert.True(t, res)

	t.Run("Happy Path - Cobfirm Auction", func(t *testing.T) {

		auctionId := uuid.New().String()

		url := fmt.Sprintf("http://localhost:8080/lh/v1/confirm_auction/%s", auctionId)

		response, err := http.Get(url)
		assert.NoError(t, err)

		// Decode the response body into the struct
		var actualRes rest.ConfirmAuctionRes
		err = json.NewDecoder(response.Body).Decode(&actualRes)
		assert.NoError(t, err)
		assert.Equal(t, len(actualRes.Fragments), 3)
		assert.Equal(t, actualRes.Fragments[0].AmountOut, "1")
		assert.Equal(t, actualRes.Fragments[1].AmountOut, "2")
		assert.Equal(t, actualRes.Fragments[2].AmountOut, "1.5")

		//assert.Equal(t, expectedRes, actualRes)
	})

}
