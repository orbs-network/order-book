package mocks

import (
	"github.com/orbs-network/order-book/models"
)

var Pk = "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
var UserType = models.MARKET_MAKER

var User = models.User{
	Id:   UserId,
	Pk:   Pk,
	Type: UserType,
}
