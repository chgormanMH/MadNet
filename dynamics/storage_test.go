package dynamics

import (
	"bytes"
	"errors"
	"testing"
)

func initializeStorage() *Storage {
	storageLogger := newLogger()
	database := initializeDB()

	s := &Storage{}
	err := s.Init(database, storageLogger)
	if err != nil {
		panic(err)
	}
	s.Start()
	return s
}

func initializeStorageCE(currentEpoch, highestEpoch uint32) *Storage {
	if currentEpoch > highestEpoch {
		panic("Cannot have currentEpoch > highestEpoch")
	}
	s := initializeStorage()

	// Initialize RawStorage to standardParameters between currentEpoch and highestEpoch
	rsStandard := &RawStorage{}
	rsStandard.standardParameters()
	for epochCurr := currentEpoch; epochCurr <= highestEpoch; epochCurr++ {
		err := s.database.SetRawStorage(epochCurr, rsStandard)
		if err != nil {
			panic(err)
		}
	}
	err := s.database.SetCurrentEpoch(currentEpoch)
	if err != nil {
		panic(err)
	}
	s.currentEpoch = currentEpoch
	err = s.database.SetHighestEpoch(highestEpoch)
	if err != nil {
		panic(err)
	}
	rsStandardBytes, err := rsStandard.Marshal()
	if err != nil {
		panic(err)
	}

	// Now to check
	for epochCurr := currentEpoch; epochCurr <= highestEpoch; epochCurr++ {
		rs, err := s.database.GetRawStorage(epochCurr)
		if err != nil {
			panic(err)
		}
		rsBytes, err := rs.Marshal()
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(rsBytes, rsStandardBytes) {
			panic("RawStorage values do not match")
		}
	}
	return s
}

// Test Storage Init with nothing initialized
func TestStorageInit1(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	s := &Storage{}
	err := s.Init(database, storageLogger)
	if err != nil {
		t.Fatal(err)
	}

	// Nothing was initialized, so we should have currentEpoch == 1
	if s.currentEpoch != 1 {
		t.Fatal("invalid currentEpoch: should be 1")
	}

	// Check currentEpoch == 1 (in the database)
	currentEpoch, err := s.database.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != s.currentEpoch {
		t.Fatal("invalid currentEpoch: does not match current value")
	}

	// Check highestEpoch == 1 (in database)
	highestEpoch, err := s.database.GetHighestEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if highestEpoch != 1 {
		t.Fatal("invalid highestEpoch: should be 1")
	}

	rs := &RawStorage{}
	rs.standardParameters()
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Check rawStorage == standardParameters
	storageRSBytes, err := s.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, storageRSBytes) {
		t.Fatal("rawStorage values do not match")
	}
}

// Test Storage Init with database initialized for currentEpoch
// and associated rawStorage, but not highestEpoch
func TestStorageInit2(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database
	epoch := uint32(25519)
	err := database.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	err = database.SetRawStorage(epoch, rsTrue)
	if err != nil {
		t.Fatal(err)
	}

	s := &Storage{}
	err = s.Init(database, storageLogger)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure currentEpoch == epoch
	if s.currentEpoch != epoch {
		t.Fatal("invalid currentEpoch")
	}

	// Ensure currentEpoch matches value from database
	currentEpoch, err := s.database.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != s.currentEpoch {
		t.Fatal("invalid currentEpoch: does not match current value")
	}

	// Ensure highestEpoch is now set to epoch
	highestEpoch, err := s.database.GetHighestEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if highestEpoch != epoch {
		t.Fatal("invalid highestEpoch: should match epoch variable")
	}

	rs := &RawStorage{}
	rs.standardParameters()
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Check rawStorage == standardParameters
	storageRSBytes, err := s.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, storageRSBytes) {
		t.Fatal("rawStorage values do not match")
	}
}

// Test Storage Init with database initialized incorrectly:
// currentEpoch is set but no associated rawStorage.
func TestStorageInit3(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database with current epoch but no rawStorage;
	// this should raise an error during initialization
	// when running GetRawStorage.
	epoch := uint32(25519)
	err := database.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}

	s := &Storage{}
	err = s.Init(database, storageLogger)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// Test Storage Init with database initialized again;
// currentEpoch is set with associated rawStorage, but highestEpoch
// is invalid (highestEpoch < currentEpoch).
func TestStorageInit4(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database
	epoch := uint32(25519)
	err := database.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	err = database.SetHighestEpoch(1)
	if err != nil {
		t.Fatal(err)
	}
	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	err = database.SetRawStorage(epoch, rsTrue)
	if err != nil {
		t.Fatal(err)
	}

	s := &Storage{}
	err = s.Init(database, storageLogger)
	if err != nil {
		t.Fatal(err)
	}

	if s.currentEpoch != epoch {
		t.Fatal("invalid currentEpoch")
	}

	currentEpoch, err := s.database.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != s.currentEpoch {
		t.Fatal("invalid currentEpoch: does not match current value")
	}

	highestEpoch, err := s.database.GetHighestEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if highestEpoch != epoch {
		t.Fatal("invalid highestEpoch: should match epoch variable")
	}

	rs := &RawStorage{}
	rs.standardParameters()
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	storageRSBytes, err := s.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, storageRSBytes) {
		t.Fatal("rawStorage values do not match")
	}
}
func TestStorageStart(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	s := &Storage{}
	err := s.Init(database, storageLogger)
	if err != nil {
		t.Fatal(err)
	}
	s.Start()
}

// Test ensures we panic when running Start before Init.
// This happens from attempting to close a closed channel.
func TestStorageStartFail(t *testing.T) {
	s := &Storage{}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should panic")
		}
	}()
	s.Start()
}

// Test ensures storage has is initialized to the correct values.
func TestStorageInitialized(t *testing.T) {
	s := initializeStorage()

	maxBytesReturned := s.GetMaxBytes()
	if maxBytesReturned != maxBytes {
		t.Fatal("Incorrect MaxBytes")
	}

	maxProposalSizeReturned := s.GetMaxProposalSize()
	if maxProposalSizeReturned != maxProposalSize {
		t.Fatal("Incorrect MaxProposalSize")
	}

	srvrMsgTimeoutReturned := s.GetSrvrMsgTimeout()
	if srvrMsgTimeoutReturned != srvrMsgTimeout {
		t.Fatal("Incorrect srvrMsgTimeout")
	}

	msgTimeoutReturned := s.GetMsgTimeout()
	if msgTimeoutReturned != msgTimeout {
		t.Fatal("Incorrect msgTimeout")
	}

	proposalStepTimeoutReturned := s.GetProposalStepTimeout()
	if proposalStepTimeoutReturned != proposalStepTO {
		t.Fatal("Incorrect proposalStepTO")
	}

	preVoteStepTimeoutReturned := s.GetPreVoteStepTimeout()
	if preVoteStepTimeoutReturned != preVoteStepTO {
		t.Fatal("Incorrect preVoteStepTO")
	}

	preCommitStepTimeoutReturned := s.GetPreCommitStepTimeout()
	if preCommitStepTimeoutReturned != preCommitStepTO {
		t.Fatal("Incorrect preCommitStepTO")
	}

	deadBlockRoundNextRoundTimeoutReturned := s.GetDeadBlockRoundNextRoundTimeout()
	if deadBlockRoundNextRoundTimeoutReturned != dBRNRTO {
		t.Fatal("Incorrect deadBlockRoundNextRoundTimeout")
	}

	downloadTimeoutReturned := s.GetDownloadTimeout()
	if downloadTimeoutReturned != downloadTO {
		t.Fatal("Incorrect downloadTimeout")
	}

	minTxBurnedFee := s.GetMinTxBurnedFee()
	if minTxBurnedFee.Sign() != 0 {
		t.Fatal("Incorrect minTxBurnedFee")
	}

	txValidVersion := s.GetTxValidVersion()
	if txValidVersion != 0 {
		t.Fatal("Incorrect txValidVersion")
	}

	minVSBurnedFee := s.GetMinValueStoreBurnedFee()
	if minVSBurnedFee.Sign() != 0 {
		t.Fatal("Incorrect minValueStoreBurnedFee")
	}

	vsTxValidVersion := s.GetValueStoreTxValidVersion()
	if vsTxValidVersion != 0 {
		t.Fatal("Incorrect valueStoreTxValidVersion")
	}

	minASBurnedFee := s.GetMinAtomicSwapBurnedFee()
	if minASBurnedFee.Sign() != 0 {
		t.Fatal("Incorrect minAtomicSwapBurnedFee")
	}

	asStopEpoch := s.GetAtomicSwapValidStopEpoch()
	if asStopEpoch != 0 {
		t.Fatal("Incorrect atomicSwapStopValidStopEpoch")
	}

	dsTxValidVersion := s.GetDataStoreTxValidVersion()
	if dsTxValidVersion != 0 {
		t.Fatal("Incorrect dataStoreTxValidVersion")
	}
}

// Test success of UpdateStorageInstance
func TestStorageLoadStorage1(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)

	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	rsTrueBytes, err := rsTrue.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadStorage(epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsBytes, err := s.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, rsTrueBytes) {
		t.Fatal("rawStorage values do not match")
	}
}

// Test success of UpdateStorageInstance again
func TestStorageLoadStorage2(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	err := s.database.SetRawStorage(epoch, s.rawStorage)
	if err != nil {
		t.Fatal(err)
	}

	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	rsTrueBytes, err := rsTrue.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadStorage(epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsBytes, err := s.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, rsTrueBytes) {
		t.Fatal("rawStorage values do not match")
	}
}

// Test failure of UpdateStorageInstance
func TestStorageLoadStorage3(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	s.rawStorage = nil
	// This should fail because we have an invalid rawStorage value
	// and we are assuming we are able to use the current rawStorage value
	err := s.LoadStorage(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// Test failure of UpdateStorageValue
func TestStorageUpdateStorageValueBad(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	field := "invalid"
	value := ""
	err := s.UpdateStorageValue(field, value, epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// Test success of UpdateStorageValue;
// Here, we have currentEpoch < epoch < highestEpoch.
func TestStorageUpdateStorageValueGood1(t *testing.T) {
	currentEpoch := uint32(1)
	highestEpoch := uint32(10)
	// Initialize epochs to standardParameters
	s := initializeStorageCE(currentEpoch, highestEpoch)

	rsStandard := &RawStorage{}
	rsStandard.standardParameters()
	// Check to confirm RawStorage values are valid
	rsStandardBytes, err := rsStandard.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Now to update MaxBytes
	field := "maxBytes"
	value := "1234567890"
	epoch := uint32(5)
	err = s.UpdateStorageValue(field, value, epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsNew, err := rsStandard.Copy()
	if err != nil {
		t.Fatal(err)
	}
	err = rsNew.UpdateValue(field, value)
	if err != nil {
		t.Fatal(err)
	}
	rsNewBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Now to check updated value
	for epochCurr := currentEpoch; epochCurr < epoch; epochCurr++ {
		rs, err := s.database.GetRawStorage(epochCurr)
		if err != nil {
			t.Fatal(err)
		}
		rsBytes, err := rs.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(rsBytes, rsStandardBytes) {
			t.Fatal("RawStorage values do not match (2)")
		}
	}
	for epochCurr := epoch; epochCurr <= highestEpoch; epochCurr++ {
		rs, err := s.database.GetRawStorage(epochCurr)
		if err != nil {
			t.Fatal(err)
		}
		rsBytes, err := rs.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(rsBytes, rsNewBytes) {
			t.Fatal("RawStorage values do not match (3)")
		}
	}
}

// Test success of UpdateStorageValue;
// Here, we have currentEpoch < highestEpoch < epoch.
func TestStorageUpdateStorageValueGood2(t *testing.T) {
	currentEpoch := uint32(1)
	highestEpoch := uint32(10)
	// Initialize epochs to standardParameters
	s := initializeStorageCE(currentEpoch, highestEpoch)

	rsStandard := &RawStorage{}
	rsStandard.standardParameters()
	// Check to confirm RawStorage values are valid
	rsStandardBytes, err := rsStandard.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Now to update MaxBytes
	field := "maxBytes"
	value := "1234567890"
	epoch := uint32(15)
	err = s.UpdateStorageValue(field, value, epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsNew, err := rsStandard.Copy()
	if err != nil {
		t.Fatal(err)
	}
	err = rsNew.UpdateValue(field, value)
	if err != nil {
		t.Fatal(err)
	}
	rsNewBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Now to check updated value
	for epochCurr := currentEpoch; epochCurr < epoch; epochCurr++ {
		rs, err := s.database.GetRawStorage(epochCurr)
		if err != nil {
			t.Fatal(err)
		}
		rsBytes, err := rs.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(rsBytes, rsStandardBytes) {
			t.Fatal("RawStorage values do not match (2)")
		}
	}
	rsEpoch, err := s.database.GetRawStorage(epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsEpochBytes, err := rsEpoch.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsEpochBytes, rsNewBytes) {
		t.Fatal("RawStorage values do not match (3)")
	}

	// Check that RawStorage after epoch has not been set
	_, err = s.database.GetRawStorage(epoch + 1)
	if !errors.Is(err, ErrKeyNotPresent) {
		t.Fatal("Should have raised ErrKeyNotPresent error")
	}
}

// Test success of UpdateStorageValue;
// Here, we have epoch < currentEpoch < highestEpoch.
func TestStorageUpdateStorageValueGood3(t *testing.T) {
	currentEpoch := uint32(5)
	highestEpoch := uint32(10)
	// Initialize epochs to standardParameters
	s := initializeStorageCE(currentEpoch, highestEpoch)

	rsStandard := &RawStorage{}
	rsStandard.standardParameters()

	// Now to update MaxBytes
	field := "maxBytes"
	value := "1234567890"
	epoch := uint32(1)
	err := s.UpdateStorageValue(field, value, epoch)
	if err != nil {
		t.Fatal(err)
	}
	rsNew, err := rsStandard.Copy()
	if err != nil {
		t.Fatal(err)
	}
	err = rsNew.UpdateValue(field, value)
	if err != nil {
		t.Fatal(err)
	}
	rsNewBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Now to check updated value
	for epochCurr := currentEpoch; epochCurr <= highestEpoch; epochCurr++ {
		rs, err := s.database.GetRawStorage(epochCurr)
		if err != nil {
			t.Fatal(err)
		}
		rsBytes, err := rs.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(rsBytes, rsNewBytes) {
			t.Fatal("RawStorage values do not match (2)")
		}
	}

	// Check that RawStorage after epoch has not been set
	_, err = s.database.GetRawStorage(highestEpoch + 1)
	if !errors.Is(err, ErrKeyNotPresent) {
		t.Fatal("Should have raised ErrKeyNotPresent error")
	}
}
