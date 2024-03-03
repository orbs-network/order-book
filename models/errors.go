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
var ErrInsufficientLiquity = errors.New("not enough liquidity in book to satisfy inAmount")
var ErrInAmount = errors.New("inAmount should be positive")
var ErrSwapInvalid = errors.New("orders in the swap can not fill any longer")
var ErrOrderPending = errors.New("order is pending")
var ErrOrderNotPending = errors.New("order is not pending")
var ErrOrderFilled = errors.New("order is filled")
var ErrOrderNotUnfilled = errors.New("order should be completely unfilled")
var ErrOrderNotPartialFilled = errors.New("order should be partially filled")
var ErrOrderCancelled = errors.New("order is cancelled")
var ErrInvalidInput = errors.New("invalid input")
var ErrSignatureVerificationError = errors.New("signature verification error")
var ErrSignatureVerificationFailed = errors.New("signature verification failed")
var ErrUnexpectedSizeFilled = errors.New("unexpected sizeFilled")
var ErrUnexpectedSizePending = errors.New("unexpected sizePending")
var ErrIterFail = errors.New("failed to get bid/ask iterator from store")
var ErrTokenNotsupported = errors.New("token is not supported")
var ErrMinOutAmount = errors.New("OutAmount is less than MinOutAmount")

// store generic errors
var ErrValAlreadyInSet = errors.New("the value is already a member of the set")
