package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orbs-network/order-book/utils"
	"github.com/stretchr/testify/assert"
)

func TestExtractPubKeyMiddleware(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicKey := utils.GetPubKeyCtx(r.Context())
		_, err := w.Write([]byte(publicKey))
		assert.NoError(t, err)
	})

	pubKey := "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
	middleware := ExtractPubKeyMiddleware(handlerFunc)

	t.Run("should extract public key from header and add to context", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/some-endpoint", nil)
		assert.NoError(t, err)

		req.Header.Set("X-Public-Key", pubKey)

		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, pubKey, rr.Body.String())
	})

	t.Run("should return error if public key header is missing", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Missing public key\n", rr.Body.String())
	})
}
