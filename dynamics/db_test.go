package dynamics

import (
	"bytes"
	"testing"

	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/utils"
	"github.com/sirupsen/logrus"
)

func initializeDB() *Database {
	db := &Database{}
	logger := logrus.New()
	db.logger = logger
	return db
}

func TestMakeRawStorageKey(t *testing.T) {
	db := initializeDB()
	epoch := uint32(0)
	_, err := db.makeRawStorageKey(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}

	epoch = 1
	key, err := db.makeRawStorageKey(epoch)
	if err != nil {
		t.Fatal(err)
	}
	prefix := dbprefix.PrefixRawStorageKey()
	epochBytes := utils.MarshalUint32(epoch)
	keyTrue := []byte{}
	keyTrue = append(keyTrue, prefix...)
	keyTrue = append(keyTrue, epochBytes...)
	if !bytes.Equal(key, keyTrue) {
		t.Fatal("Incorrect RawStorageKey")
	}
}

func TestGetRawStorage(t *testing.T) {
	db := initializeDB()
	epoch := uint32(0)
	_, err := db.GetRawStorage(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestSetRawStorage(t *testing.T) {
	db := initializeDB()
	epoch := uint32(0)
	rs := &RawStorage{}
	err := db.SetRawStorage(epoch, rs)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestMakeCurrentEpochKey(t *testing.T) {
	db := initializeDB()
	key, err := db.makeCurrentEpochKey()
	if err != nil {
		t.Fatal(err)
	}
	prefix := dbprefix.PrefixRawStorageKey()
	currentEpoch := constants.StorageCurrentEpoch()
	keyTrue := []byte{}
	keyTrue = append(keyTrue, prefix...)
	keyTrue = append(keyTrue, currentEpoch...)
	if !bytes.Equal(key, keyTrue) {
		t.Fatal("Incorrect CurrentEpochKey")
	}
}

func TestSetCurrentEpoch(t *testing.T) {
	epoch := uint32(0)
	db := initializeDB()
	err := db.SetCurrentEpoch(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestMakeHighestEpochKey(t *testing.T) {
	db := initializeDB()
	key, err := db.makeHighestEpochKey()
	if err != nil {
		t.Fatal(err)
	}
	prefix := dbprefix.PrefixRawStorageKey()
	highestEpoch := constants.StorageHighestEpoch()
	keyTrue := []byte{}
	keyTrue = append(keyTrue, prefix...)
	keyTrue = append(keyTrue, highestEpoch...)

	if !bytes.Equal(key, keyTrue) {
		t.Fatal("Incorrect HighestEpochKey")
	}
}

func TestSetHighestEpoch(t *testing.T) {
	epoch := uint32(0)
	db := initializeDB()
	err := db.SetHighestEpoch(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
