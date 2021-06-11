package dynamics

import (
	"sync"
)

// Database is an abstraction of the header trie and the object storage
type Database struct {
	sync.Mutex
	//rawDB  *rawDataBase
	//logger *logrus.Logger
}

// GetCurrentStorageInstance returns the current StorageInstance
// from the database
func (db *Database) GetCurrentStorageInstance() (*StorageInstance, error) {
	// Look up currentEpoch
	currentEpoch, err := db.GetCurrentEpoch()
	if err != nil {
		return nil, err
	}
	// Look up corresponding StorageInstance
	si, err := db.GetStorageInstance(currentEpoch)
	if err != nil {
		return nil, err
	}
	return si, nil
}

// GetStorageInstance returns the StorageInstance for epoch from the database
func (db *Database) GetStorageInstance(epoch uint32) (*StorageInstance, error) {
	panic("not implemented")
	// Look up currentEpoch
	// Look up corresponding StorageInstance
}

// SetStorageInstance sets the StorageInstance for epoch in the database
func (db *Database) SetStorageInstance(epoch uint32, si *StorageInstance) error {
	panic("not implemented")
	// Store StorageInstance at correct location
}

// GetCurrentEpoch returns the current epoch from the database
func (db *Database) GetCurrentEpoch() (uint32, error) {
	panic("not implemented")
	// Look up currentEpoch
}

// SetCurrentEpoch sets the current epoch in the database
func (db *Database) SetCurrentEpoch(epoch uint32) error {
	panic("not implemented")
	// Store value at correct location
}

// GetHighestEpoch returns the highest epoch from the database
// which has a non-nil StorageInstance value
func (db *Database) GetHighestEpoch() (uint32, error) {
	panic("not implemented")
	// Look up highestEpoch from the correct location
}

// SetHighestEpoch sets the highest epoch in the database
// which has a non-nil StorageInstance value
func (db *Database) SetHighestEpoch(epoch uint32) error {
	panic("not implemented")
	// Set highestEpoch at the correct location
}
