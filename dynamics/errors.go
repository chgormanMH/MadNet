package dynamics

import (
	"errors"
)

var (
	// ErrRawStorageNilPointer is an error which results from a
	// RawStorage struct which has not been initialized.
	ErrRawStorageNilPointer = errors.New("Invalid RawStorage: nil pointer")

	// ErrZeroEpoch is an error which is raised whenever the epoch is given
	// as zero; there is no zero epoch.
	ErrZeroEpoch = errors.New("invalid epoch: no zero epoch")
)
