package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Signature2Bytes(t *testing.T) {
	sig1 := "123456789ABCDEF"
	sig2 := "0x" + sig1

	bt1 := Signature2Bytes(sig1)
	bt2 := Signature2Bytes(sig2)

	assert.Equal(t, bt1, bt2)
}
