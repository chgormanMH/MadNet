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

// GetCurrentRawStorage returns the current RawStorage
// from the database
func (db *Database) GetCurrentRawStorage() (*RawStorage, error) {
	// Look up currentEpoch
	currentEpoch, err := db.GetCurrentEpoch()
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	if currentEpoch == 0 {
		// TODO: Need to do something specific if currentEpoch == 0.
		// Load standard parameters or return error?
		panic("not implemented")
	}
	// Look up corresponding RawStorage
	rs, err := db.GetRawStorage(currentEpoch)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	return rs, nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (db *Database) makeRawStorageKey(epoch uint32) ([]byte, error) {
	if epoch == 0 {
		return nil, ErrZeroEpoch
	}
	prefix := dbprefix.PrefixRawStorageKey()
	epochBytes := utils.MarshalUint32(epoch)
	key := []byte{}
	key = append(key, prefix...)
	key = append(key, epochBytes...)
	return key, nil
}

// GetRawStorage returns the RawStorage for epoch from the database
func (db *Database) GetRawStorage(epoch uint32) (*RawStorage, error) {
	key, err := db.makeRawStorageKey(epoch)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	v, err := db.rawDB.GetValue(key)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	rs := &RawStorage{}
	err = rs.Unmarshal(v)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return nil, err
	}
	return rs, nil
}

// SetRawStorage sets the RawStorage for epoch in the database
func (db *Database) SetRawStorage(epoch uint32, rs *RawStorage) error {
	key, err := db.makeRawStorageKey(epoch)
	if err != nil {
		utils.DebugTrace(db.logger, err)
		return err
	}
	value, err := rs.Marshal()
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
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (db *Database) makeCurrentEpochKey() ([]byte, error) {
	prefix := dbprefix.PrefixRawStorageKey()
	currentEpoch := constants.StorageCurrentEpoch()
	key := []byte{}
	key = append(key, prefix...)
	key = append(key, currentEpoch...)
	return key, nil
}

// GetCurrentEpoch returns the current epoch from the database
// TODO: What should happen if value is 0 or does not exist?
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
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (db *Database) makeHighestEpochKey() ([]byte, error) {
	prefix := dbprefix.PrefixRawStorageKey()
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
}
