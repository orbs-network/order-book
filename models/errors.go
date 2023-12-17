package models

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrClashingOrderId = errors.New("order already exists")
var ErrClashingClientOrderId = errors.New("order already exists")
var ErrUnexpectedError = errors.New("unexpected error")
var ErrMarshalError = errors.New("marshal error")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrNoUserInContext = errors.New("no user in context")
var ErrUnauthorized = errors.New("user not allowed to perform this action")
var ErrOrderNotOpen = errors.New("order must be status open to perform this action")
var ErrInsufficientLiquity = errors.New("not enough liquidity in book to satisfy amountIn")
var ErrInAmount = errors.New("amountIn should be positive")
var ErrSwapInvalid = errors.New("orders in the swap can not fill any longer")
var ErrOrderPending = errors.New("order is pending")
var ErrOrderFilled = errors.New("order is filled")
var ErrInvalidInput = errors.New("invalid input")
var ErrSignatureVerificationError = errors.New("signature verification error")
var ErrSignatureVerificationFailed = errors.New("signature verification failed")
var ErrInvalidSize = errors.New("updated sizeFilled is greater than size")

// store generic errors
var ErrValAlreadyInSet = errors.New("the value is already a member of the set")
