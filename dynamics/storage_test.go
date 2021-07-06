package dynamics

import (
	"bytes"
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
