package dynamics

import (
	"errors"
)

var (
	// ErrRawStorageNilPointer is an error which results from a
	// RawStorage struct which has not been initialized.
	ErrRawStorageNilPointer = errors.New("invalid RawStorage: nil pointer")

	// ErrZeroEpoch is an error which is raised whenever the epoch is given
	// as zero; there is no zero epoch.
	ErrZeroEpoch = errors.New("invalid epoch: no zero epoch")

	// ErrUnmarshalEmpty is an error which is raised whenever attempting
	// to unmarshal an empty byte slice.
	ErrUnmarshalEmpty = errors.New("invalid: attempting to unmarshal empty byte slice")

	// ErrKeyNotPresent is an error which is raised when a key is not present
	// in the database.
	ErrKeyNotPresent = errors.New("invalid: Key not found")

	// ErrInvalidUpdateValue is an error which is returned when the data
	// for updating rawStorage is invalid.
	ErrInvalidUpdateValue = errors.New("invalid update value for storage")

	// ErrInvalidValue is an error which is returned when the value is invalid.
	ErrInvalidValue = errors.New("invalid value")
)
