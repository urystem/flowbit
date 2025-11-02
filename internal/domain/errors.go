package domain

import "errors"

var (
	ErrSymbolNotFound = errors.New("symbol not found")
	ErrInternal       = errors.New("internal error")
	// ErrInvalidPrice   = errors.New("invalid price")
)
