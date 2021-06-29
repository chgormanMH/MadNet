package dynamics

import (
	"bytes"
	"testing"

	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/utils"
	"github.com/sirupsen/logrus"
)

type mockRawDB struct {
	rawDB map[string]string
}

func (m *mockRawDB) GetValue(key []byte) ([]byte, error) {
	strKey := string(key)
	strValue, ok := m.rawDB[strKey]
	if !ok {
		return nil, nil
	}
	value := []byte(strValue)
	return value, nil
}

func (m *mockRawDB) SetValue(key []byte, value []byte) error {
	strKey := string(key)
	strValue := string(value)
	m.rawDB[strKey] = strValue
	return nil
}

func initializeDB() *Database {
	db := &Database{}
	logger := logrus.New()
	db.logger = logger
	mock := &mockRawDB{}
	mock.rawDB = make(map[string]string)
	db.rawDB = mock
	return db
}

func TestGetCurrentRawStorageKey(t *testing.T) {
	db := initializeDB()
	rs, currentEpoch, err := db.GetCurrentRawStorage()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != 0 {
		t.Fatal("currentEpoch should be 0")
	}
	if rs != nil {
		t.Fatal("rawStorage should be nil")
	}

	epochTrue := uint32(1)
	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	err = db.SetCurrentEpoch(epochTrue)
	if err != nil {
		t.Fatal(err)
	}
	err = db.SetRawStorage(epochTrue, rsTrue)
	if err != nil {
		t.Fatal(err)
	}

	rs, currentEpoch, err = db.GetCurrentRawStorage()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != epochTrue {
		t.Fatal("currentEpochs do not agree")
	}
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rsTrueBytes, err := rsTrue.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, rsTrueBytes) {
		t.Fatal("rawStorage values do not agree")
	}
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
		t.Fatal("Incorrect RawStorageKey (1)")
	}

	epoch = 4294967295
	key, err = db.makeRawStorageKey(epoch)
	if err != nil {
		t.Fatal(err)
	}
	prefix = dbprefix.PrefixRawStorageKey()
	epochBytes = utils.MarshalUint32(epoch)
	keyTrue = []byte{}
	keyTrue = append(keyTrue, prefix...)
	keyTrue = append(keyTrue, epochBytes...)
	if !bytes.Equal(key, keyTrue) {
		t.Fatal("Incorrect RawStorageKey (2)")
	}
}

func TestGetSetRawStorage(t *testing.T) {
	db := initializeDB()
	epoch := uint32(0)
	_, err := db.GetRawStorage(epoch)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	rs := &RawStorage{}
	err = db.SetRawStorage(epoch, rs)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	epoch = uint32(1)
	rs.standardParameters()
	err = db.SetRawStorage(epoch, rs)
	if err != nil {
		t.Fatal(err)
	}

	rsRcvd, err := db.GetRawStorage(epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rsRcvdBytes, err := rsRcvd.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, rsRcvdBytes) {
		t.Fatal("rawStorage are not equal")
	}

	err = db.SetRawStorage(epoch, nil)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestMakeCurrentEpochKey(t *testing.T) {
	db := initializeDB()
	key := db.makeCurrentEpochKey()
	prefix := dbprefix.PrefixRawStorageKey()
	currentEpoch := constants.StorageCurrentEpoch()
	keyTrue := []byte{}
	keyTrue = append(keyTrue, prefix...)
	keyTrue = append(keyTrue, currentEpoch...)
	if !bytes.Equal(key, keyTrue) {
		t.Fatal("Incorrect CurrentEpochKey")
	}
}

func TestGetSetCurrentEpoch(t *testing.T) {
	db := initializeDB()
	epoch := uint32(0)
	err := db.SetCurrentEpoch(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}

	// No CurrentEpoch present in database; should return 0
	curEpoch, err := db.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if curEpoch != 0 {
		t.Fatal("currentEpoch should be 0")
	}

	// Set currentEpoch in database and then check
	epoch = uint32(1)
	err = db.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	curEpoch, err = db.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if curEpoch != epoch {
		t.Fatal("currentEpochs are not equal (1)")
	}

	// Set currentEpoch in database and then check (again)
	epoch = uint32(25519)
	err = db.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	curEpoch, err = db.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if curEpoch != epoch {
		t.Fatal("currentEpochs are not equal (2)")
	}
}

func TestMakeHighestEpochKey(t *testing.T) {
	db := initializeDB()
	key := db.makeHighestEpochKey()
	prefix := dbprefix.PrefixRawStorageKey()
	highestEpoch := constants.StorageHighestEpoch()
	keyTrue := []byte{}
	keyTrue = append(keyTrue, prefix...)
	keyTrue = append(keyTrue, highestEpoch...)

	if !bytes.Equal(key, keyTrue) {
		t.Fatal("Incorrect HighestEpochKey")
	}
}

func TestGetSetHighestEpoch(t *testing.T) {
	db := initializeDB()
	epoch := uint32(0)
	err := db.SetHighestEpoch(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}

	// No HighestEpoch present in database; should return 0
	highestEpoch, err := db.GetHighestEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if highestEpoch != 0 {
		t.Fatal("highestEpoch should be 0")
	}

	// Set highestEpoch in database and check
	epoch = uint32(1)
	err = db.SetHighestEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	highestEpoch, err = db.GetHighestEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if highestEpoch != epoch {
		t.Fatal("highestEpochs are not equal (1)")
	}

	// Set highestEpoch in database and check (again)
	epoch = uint32(25519)
	err = db.SetHighestEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	highestEpoch, err = db.GetHighestEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if highestEpoch != epoch {
		t.Fatal("highestEpochs are not equal (2)")
	}
}
