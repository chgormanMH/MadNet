package dynamics

import (
	"errors"
)

var (
	ErrInvalidSortedArray = errors.New("Invalid Sorted Array")
	ErrIncorrectData      = errors.New("Invalid: incorrect data")
	ErrInvalidFormat      = errors.New("Invalid: incorrect format")
	ErrIncorrectAdd       = errors.New("Invalid: invalid additional element")
	ErrEmptyArray         = errors.New("Invalid entry: empty array")
	ErrNotPresent         = errors.New("Invalid: entry not present")
)
