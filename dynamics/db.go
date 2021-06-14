package dynamics

import (
	"sync"

	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/utils"
	"github.com/sirupsen/logrus"
)

// Database is an abstraction for object storage
type Database struct {
	sync.Mutex
	rawDB  *rawDataBase
	logger *logrus.Logger
}

// GetCurrentStorageInstance returns the current RawStorage
// from the database
func (db *Database) GetCurrentStorageInstance() (*RawStorage, error) {
	// Look up currentEpoch
	currentEpoch, err := db.GetCurrentEpoch()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	// Look up corresponding RawStorage
	si, err := db.GetStorageInstance(currentEpoch)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	return si, nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (db *Database) makeStorageInstanceKey(epoch uint32) ([]byte, error) {
	if epoch == 0 {
		return nil, ErrZeroEpoch
	}
	prefix := dbprefix.PrefixStorageInstanceKey()
	epochBytes := utils.MarshalUint32(epoch)
	key := []byte{}
	key = append(key, prefix...)
	key = append(key, epochBytes...)
	return key, nil
}

// GetStorageInstance returns the RawStorage for epoch from the database
func (db *Database) GetStorageInstance(epoch uint32) (*RawStorage, error) {
	key, err := db.makeStorageInstanceKey(epoch)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	v, err := db.rawDB.GetValue(key)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	si := &RawStorage{}
	err = si.Unmarshal(v)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	return si, nil
	// Look up currentEpoch
	// Look up corresponding RawStorage
}

// SetStorageInstance sets the RawStorage for epoch in the database
func (db *Database) SetStorageInstance(epoch uint32, si *RawStorage) error {
	key, err := db.makeStorageInstanceKey(epoch)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	value, err := si.Marshal()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	err = db.rawDB.SetValue(key, value)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	return nil
	// Store RawStorage at correct location
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (db *Database) makeCurrentEpochKey() ([]byte, error) {
	prefix := dbprefix.PrefixStorageInstanceKey()
	currentEpoch := constants.StorageCurrentEpoch()
	key := []byte{}
	key = append(key, prefix...)
	key = append(key, currentEpoch...)
	return key, nil
}

// GetCurrentEpoch returns the current epoch from the database
func (db *Database) GetCurrentEpoch() (uint32, error) {
	key, err := db.makeCurrentEpochKey()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return 0, err
	}
	v, err := db.rawDB.GetValue(key)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return 0, err
	}
	value, err := utils.UnmarshalUint32(v)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return 0, err
	}
	return value, nil
	// Look up currentEpoch
}

// SetCurrentEpoch sets the current epoch in the database
func (db *Database) SetCurrentEpoch(epoch uint32) error {
	if epoch == 0 {
		return ErrZeroEpoch
	}
	key, err := db.makeCurrentEpochKey()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	value := utils.MarshalUint32(epoch)
	err = db.rawDB.SetValue(key, value)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	return nil
	// Store value at correct location
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (db *Database) makeHighestEpochKey() ([]byte, error) {
	prefix := dbprefix.PrefixStorageInstanceKey()
	highestEpoch := constants.StorageHighestEpoch()
	key := []byte{}
	key = append(key, prefix...)
	key = append(key, highestEpoch...)
	return key, nil
}

// GetHighestEpoch returns the highest epoch from the database
// which has a non-nil RawStorage value
func (db *Database) GetHighestEpoch() (uint32, error) {
	key, err := db.makeHighestEpochKey()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return 0, err
	}
	v, err := db.rawDB.GetValue(key)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return 0, err
	}
	value, err := utils.UnmarshalUint32(v)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return 0, err
	}
	return value, nil
}

// SetHighestEpoch sets the highest epoch in the database
// which has a non-nil RawStorage value
func (db *Database) SetHighestEpoch(epoch uint32) error {
	if epoch == 0 {
		return ErrZeroEpoch
	}
	key, err := db.makeHighestEpochKey()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	value := utils.MarshalUint32(epoch)
	err = db.rawDB.SetValue(key, value)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	return nil
	// Set highestEpoch at the correct location
}
