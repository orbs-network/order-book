package models

import "errors"

var ErrOrderAlreadyExists = errors.New("order already exists")
var ErrOrderNotFound = errors.New("order not found")
var ErrUnexpectedError = errors.New("unexpected error")
var ErrMarshalError = errors.New("marshal error")
var ErrNoUserInContext = errors.New("no user in context")
var ErrUnauthorized = errors.New("user not allowed to perform this action")
var ErrOrderNotOpen = errors.New("order must be status open to perform this action")
var ErrTransactionFailed = errors.New("transaction failed")
var ErrInsufficientLiquity = errors.New("not enough liquidity in book to satisfy amountIn")
