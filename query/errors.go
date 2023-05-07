package query

import "errors"

var (
	ErrTableNotSet = errors.New("table name not provided")
	ErrValidation  = errors.New("invalid")
)
