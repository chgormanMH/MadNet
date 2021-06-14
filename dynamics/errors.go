package dynamics

import (
	"errors"
)

var (
	// ErrStorageInstanceNilPointer is an error which results from a
	// StorageInstance struct which has not been initialized.
	ErrStorageInstanceNilPointer = errors.New("Invalid StorageInstance: nil pointer")

	// ErrZeroEpoch is an error which is raised whenever the epoch is given
	// as zero; there is no zero epoch.
	ErrZeroEpoch = errors.New("invalid epoch: no zero epoch")
)
