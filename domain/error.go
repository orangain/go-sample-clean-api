package domain

import "errors"

var (
	ErrNotFound = errors.New("Item not found")
	ErrConflict = errors.New("Conflict with the current state")
)
