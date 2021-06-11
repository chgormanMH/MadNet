package dynamics

import (
	"errors"
)

var (
	// ErrStorageInstanceNilPointer is an error which results from a
	// StorageInstance struct which has not been initialized.
	ErrStorageInstanceNilPointer = errors.New("Invalid StorageInstance: nil pointer")
)
