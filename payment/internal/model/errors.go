package model

import "errors"

var (
	ErrInvalidArgument       = errors.New("order_uuid and user_uuid must be set")
	ErrPaymentInternalServer = errors.New("internal error while processing payment")
)
