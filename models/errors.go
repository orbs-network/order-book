package models

import "errors"

var ErrOrderAlreadyExists = errors.New("order already exists")
var ErrOrderNotFound = errors.New("order not found")
var ErrUnexpectedError = errors.New("unexpected error")
var ErrNoUserInContext = errors.New("no user in context")
var ErrUnauthorized = errors.New("user not allowed to perform this action")
