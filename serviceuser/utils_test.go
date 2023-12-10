package serviceuser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAPIKey(t *testing.T) {
	key, err := GenerateAPIKey()
	assert.NoError(t, err)
	assert.Len(t, key, 44)
}
