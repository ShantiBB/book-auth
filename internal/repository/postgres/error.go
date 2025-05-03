package postgres

import "errors"

var (
	DuplicateError = errors.New("duplicate key error")
)
