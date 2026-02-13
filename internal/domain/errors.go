package domain

import "errors"

var (
	ErrOrderSizePositive = errors.New("order size must be greater than 0")
	ErrNoPackSizes       = errors.New("no pack sizes available")
	ErrEmptyPackSizes    = errors.New("pack sizes cannot be empty")
	ErrInvalidPackSize   = errors.New("invalid pack size")
	ErrTooManyPackSizes  = errors.New("too many pack sizes")
)
