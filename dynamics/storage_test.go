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

	if s.currentEpoch != 1 {
		t.Fatal("invalid currentEpoch: should be 1")
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

// Test Storage Init with database initialized
func TestStorageInit2(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database
	epoch := uint32(25519)
	database.SetCurrentEpoch(epoch)
	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	err := database.SetRawStorage(epoch, rsTrue)
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

// Test Storage Init with database initialized incorrectly
func TestStorageInit3(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database with current epoch but no rawStorage;
	// this should lead to an error
	epoch := uint32(25519)
	database.SetCurrentEpoch(epoch)

	s := &Storage{}
	err := s.Init(database, storageLogger)
	if err == nil {
		t.Fatal("Should have raised error")
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

func TestStorageStartFail(t *testing.T) {
	s := &Storage{}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should panic")
		}
	}()
	s.Start()
}

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
func TestStorageUpdateInstance1(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)

	rsTrue := &RawStorage{}
	rsTrue.standardParameters()
	rsTrueBytes, err := rsTrue.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	err = s.UpdateStorageInstance(epoch)
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
func TestStorageUpdateInstance2(t *testing.T) {
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

	err = s.UpdateStorageInstance(epoch)
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
func TestStorageUpdateInstance3(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	s.rawStorage = nil
	// This should fail because we have an invalid rawStorage value
	// and we are assuming we are able to use the current rawStorage value
	err := s.UpdateStorageInstance(epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
