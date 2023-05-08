package query

import "errors"

var (
	ErrTableNotSet        = errors.New("table name not provided")
	ErrValidation         = errors.New("invalid")
	ErrUpdateValuesNotSet = errors.New("value for update not set")
	ErrNo                 = errors.New("value for update not set")
)
