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

func initializeStorageCE(currentEpoch uint32) *Storage {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database
	rs := &RawStorage{}
	rs.standardParameters()

	// Prepare LinkedList
	node, ll, err := CreateLinkedList(currentEpoch, rs)
	if err != nil {
		panic(err)
	}
	err = database.SetNode(node)
	if err != nil {
		panic(err)
	}
	err = database.SetLinkedList(ll)
	if err != nil {
		panic(err)
	}

	s := &Storage{}
	err = s.Init(database, storageLogger)
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
	s.Start()

	// Check currentEpoch == 1 (in the database)
	currentEpoch, err := s.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != 1 {
		t.Fatal("invalid currentEpoch: does not match current value")
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

// Test Storage Init with database initialized
func TestStorageInit2(t *testing.T) {
	storageLogger := newLogger()
	database := initializeDB()

	// Initialize database
	epoch := uint32(25519)
	rsTrue := &RawStorage{}
	rsTrue.standardParameters()

	// Prepare LinkedList
	node, ll, err := CreateLinkedList(epoch, rsTrue)
	if err != nil {
		t.Fatal(err)
	}
	err = database.SetNode(node)
	if err != nil {
		t.Fatal(err)
	}
	err = database.SetLinkedList(ll)
	if err != nil {
		t.Fatal(err)
	}

	s := &Storage{}
	err = s.Init(database, storageLogger)
	if err != nil {
		t.Fatal(err)
	}
	s.Start()

	// Ensure currentEpoch matches value from database
	currentEpoch, err := s.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if currentEpoch != epoch {
		t.Fatal("invalid currentEpoch: does not match current value")
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

	// Incorrectly initialize database;
	// this should raise an error during initialization
	// when running loadStorage.
	epoch := uint32(1)
	ll := &LinkedList{
		epochLastUpdated: epoch,
		currentEpoch:     epoch,
	}
	err := database.SetLinkedList(ll)
	if err != nil {
		t.Fatal(err)
	}

	s := &Storage{}
	err = s.Init(database, storageLogger)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageStartGood(t *testing.T) {
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

func TestStorageCheckUpdate(t *testing.T) {
	fieldBad := "invalid"
	valueBad := "invalid"
	epochGood := uint32(25519)
	err := checkUpdate(fieldBad, valueBad, epochGood)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	fieldGood := "maxBytes"
	valueGood := "1234567890"
	err = checkUpdate(fieldGood, valueGood, epochGood)
	if err != nil {
		t.Fatal(err)
	}

	epochBad := uint32(0)
	err = checkUpdate(fieldGood, valueGood, epochBad)
	if !errors.Is(err, ErrInvalidUpdateValue) {
		t.Fatal("Should have raised error (2)")
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
	epoch := uint32(25519)
	s := initializeStorageCE(epoch)
	// We attempt to load an epoch for which we do not have data for;
	// this should raise an error.
	err := s.LoadStorage(1)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// Test success of UpdateStorageInstance again
func TestStorageLoadStorage3(t *testing.T) {
	epoch := uint32(1)
	s := initializeStorageCE(epoch)
	rs := &RawStorage{}
	rs.standardParameters()
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	newMaxBytes := uint32(12345)
	rsNew := &RawStorage{}
	rsNew.standardParameters()
	rsNew.MaxBytes = newMaxBytes
	newEpoch := uint32(10)
	newNode := &Node{
		thisEpoch:  newEpoch,
		rawStorage: rsNew,
	}
	err = s.addNode(newNode)
	if err != nil {
		t.Fatal(err)
	}

	newMaxBytes2 := uint32(123456)
	rsNew2 := &RawStorage{}
	rsNew2.standardParameters()
	rsNew2.MaxBytes = newMaxBytes2
	newEpoch2 := uint32(100)
	newNode2 := &Node{
		thisEpoch:  newEpoch2,
		rawStorage: rsNew2,
	}
	err = s.addNode(newNode2)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadStorage(epoch)
	if err != nil {
		t.Fatal(err)
	}
	retRS, err := s.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err := retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsBytes) {
		t.Fatal("invalid rawStorage")
	}
}

func TestStorageSetGetCurrentEpoch1(t *testing.T) {
	epoch := uint32(1)
	s := initializeStorageCE(epoch)
	retCE, err := s.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if retCE != epoch {
		t.Fatal("Invalid current epoch (1)")
	}

	newEpoch := uint32(25519)
	err = s.SetCurrentEpoch(newEpoch)
	if err != nil {
		t.Fatal(err)
	}
	retCE, err = s.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if retCE != newEpoch {
		t.Fatal("Invalid current epoch (2)")
	}
}

func TestStorageSetGetCurrentEpoch2(t *testing.T) {
	s := initializeStorage()
	// Should raise error for attempting to set current epoch to 0
	badEpoch := uint32(0)
	err := s.SetCurrentEpoch(badEpoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeGood1(t *testing.T) {
	origEpoch := uint32(1)
	s := initializeStorageCE(origEpoch)
	rs := &RawStorage{}
	rs.standardParameters()
	rsStandardBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	newMaxBytes := uint32(12345)
	rs.MaxBytes = newMaxBytes
	epoch := uint32(10)
	newNode := &Node{
		prevEpoch:  0,
		thisEpoch:  epoch,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err = s.addNode(newNode)
	if err != nil {
		t.Fatal(err)
	}
	rsNewBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Check everything
	origNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if origNode.prevEpoch != origEpoch {
		t.Fatal("origNode.prevEpoch is invalid")
	}
	if origNode.thisEpoch != origEpoch {
		t.Fatal("origNode.thisEpoch is invalid")
	}
	if origNode.nextEpoch != epoch {
		t.Fatal("origNode.nextEpoch is invalid")
	}
	retRS, err := origNode.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err := retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsStandardBytes) {
		t.Fatal("invalid RawStorage")
	}

	addedNode, err := s.database.GetNode(epoch)
	if err != nil {
		t.Fatal(err)
	}
	if addedNode.prevEpoch != origEpoch {
		t.Fatal("addedNode.prevEpoch is invalid")
	}
	if addedNode.thisEpoch != epoch {
		t.Fatal("addedNode.thisEpoch is invalid")
	}
	if addedNode.nextEpoch != epoch {
		t.Fatal("addedNode.nextEpoch is invalid")
	}
	retRS, err = addedNode.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err = retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsNewBytes) {
		t.Fatal("invalid RawStorage (2)")
	}
}

func TestStorageAddNodeGood2(t *testing.T) {
	origEpoch := uint32(10)
	s := initializeStorageCE(origEpoch)
	rs := &RawStorage{}
	rs.standardParameters()
	rsStandardBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	newMaxBytes := uint32(12345)
	rs.MaxBytes = newMaxBytes
	rsNewBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	epoch := uint32(1)
	newNode := &Node{
		prevEpoch:  0,
		thisEpoch:  epoch,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err = s.addNode(newNode)
	if err != nil {
		t.Fatal(err)
	}

	// Check everything
	addedNode, err := s.database.GetNode(epoch)
	if err != nil {
		t.Fatal(err)
	}
	if addedNode.prevEpoch != epoch {
		t.Fatal("addedNode.prevEpoch is invalid")
	}
	if addedNode.thisEpoch != epoch {
		t.Fatal("addedNode.thisEpoch is invalid")
	}
	if addedNode.nextEpoch != origEpoch {
		t.Fatal("addedNode.nextEpoch is invalid")
	}
	retRS, err := addedNode.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err := retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsNewBytes) {
		t.Fatal("invalid RawStorage")
	}

	origNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if origNode.prevEpoch != epoch {
		t.Fatal("origNode.prevEpoch is invalid")
	}
	if origNode.thisEpoch != origEpoch {
		t.Fatal("origNode.thisEpoch is invalid")
	}
	if origNode.nextEpoch != origEpoch {
		t.Fatal("origNode.nextEpoch is invalid")
	}
	retRS, err = origNode.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err = retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsStandardBytes) {
		t.Fatal("invalid RawStorage")
	}
}

func TestStorageAddNodeGood3(t *testing.T) {
	origEpoch := uint32(100)
	s := initializeStorageCE(origEpoch)
	rs := &RawStorage{}
	rs.standardParameters()
	rsStandardBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	newMaxBytes := uint32(12345)
	rs.MaxBytes = newMaxBytes
	rsNewBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	addedEpoch := uint32(10)
	newNode := &Node{
		prevEpoch:  0,
		thisEpoch:  addedEpoch,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err = s.addNode(newNode)
	if err != nil {
		t.Fatal(err)
	}

	rs.standardParameters()
	addedEpoch2 := uint32(1)
	newNode2 := &Node{
		prevEpoch:  0,
		thisEpoch:  addedEpoch2,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err = s.addNode(newNode2)
	if err != nil {
		t.Fatal(err)
	}

	// Check everything
	origNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if origNode.prevEpoch != addedEpoch {
		t.Fatal("origNode.prevEpoch is invalid")
	}
	if origNode.thisEpoch != origEpoch {
		t.Fatal("origNode.thisEpoch is invalid")
	}
	if origNode.nextEpoch != origEpoch {
		t.Fatal("origNode.nextEpoch is invalid")
	}
	retRS, err := origNode.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err := retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsStandardBytes) {
		t.Fatal("invalid RawStorage")
	}

	addedNode, err := s.database.GetNode(addedEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if addedNode.prevEpoch != addedEpoch2 {
		t.Fatal("addedNode.prevEpoch is invalid")
	}
	if addedNode.thisEpoch != addedEpoch {
		t.Fatal("addedNode.thisEpoch is invalid")
	}
	if addedNode.nextEpoch != origEpoch {
		t.Fatal("addedNode.nextEpoch is invalid")
	}
	retRS, err = addedNode.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err = retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsNewBytes) {
		t.Fatal("invalid RawStorage")
	}

	addedNode2, err := s.database.GetNode(addedEpoch2)
	if err != nil {
		t.Fatal(err)
	}
	if addedNode2.prevEpoch != addedEpoch2 {
		t.Fatal("addedNode2.prevEpoch is invalid")
	}
	if addedNode2.thisEpoch != addedEpoch2 {
		t.Fatal("addedNode2.thisEpoch is invalid")
	}
	if addedNode2.nextEpoch != addedEpoch {
		t.Fatal("addedNode2.nextEpoch is invalid")
	}
	retRS, err = addedNode2.rawStorage.Copy()
	if err != nil {
		t.Fatal(err)
	}
	retRSBytes, err = retRS.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retRSBytes, rsStandardBytes) {
		t.Fatal("invalid RawStorage")
	}
}

func TestStorageAddNodeBad1(t *testing.T) {
	s := initializeStorage()
	rs := &RawStorage{}
	newNode := &Node{
		prevEpoch:  0,
		thisEpoch:  0,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err := s.addNode(newNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeBad2(t *testing.T) {
	s := initializeStorage()
	rs := &RawStorage{}
	newNode := &Node{
		prevEpoch:  0,
		thisEpoch:  1,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err := s.addNode(newNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

/*
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
*/

/*
// Test failure of UpdateStorage
func TestStorageUpdateStorageBad(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	field := "invalid"
	value := ""
	err := s.UpdateStorage(field, value, epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
*/

/*
// Test success of UpdateStorage;
func TestStorageUpdateStorageValueGood(t *testing.T) {
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
	err = s.UpdateStorage(field, value, epoch)
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
*/

/*
// Test failure of UpdateStorageValue
func TestStorageUpdateStorageValueBad1(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	field := "invalid"
	value := ""
	err := s.updateStorageValue(field, value, epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
*/

/*
// Test failure of UpdateStorageValue
func TestStorageUpdateStorageValueBad2(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	field := "maxBytes"
	value := "1"

	// Create invalid data at CurrentEpochKey
	ceKey := s.database.makeCurrentEpochKey()
	invalidBytes := []byte("Invalid Bytes")
	err := s.database.rawDB.SetValue(ceKey, invalidBytes)
	if err != nil {
		t.Fatal(err)
	}

	err = s.updateStorageValue(field, value, epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
*/

/*
// Test failure of UpdateStorageValue
func TestStorageUpdateStorageValueBad3(t *testing.T) {
	s := initializeStorage()
	epoch := uint32(25519)
	field := "maxBytes"
	value := "1"

	// Create invalid data at HighestEpochKey
	heKey := s.database.makeHighestEpochKey()
	invalidBytes := []byte("Invalid Bytes")
	err := s.database.rawDB.SetValue(heKey, invalidBytes)
	if err != nil {
		t.Fatal(err)
	}

	err = s.updateStorageValue(field, value, epoch)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
*/

/*
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
	err = s.updateStorageValue(field, value, epoch)
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
*/

/*
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
	err = s.updateStorageValue(field, value, epoch)
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
*/

/*
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
	err := s.updateStorageValue(field, value, epoch)
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

	// Check that current RawStorage is valid;
	// this should match the updated value.
	rsCurrBytes, err := s.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsCurrBytes, rsNewBytes) {
		t.Fatal("RawStorage values do not match (3)")
	}
}
*/

/*
func TestStorageGetSetCurrentEpoch(t *testing.T) {
	currentEpoch := uint32(1)
	highestEpoch := uint32(10)
	// Initialize epochs to standardParameters
	s := initializeStorageCE(currentEpoch, highestEpoch)
	retCE, err := s.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if retCE != currentEpoch {
		t.Fatal("currentEpochs are not equal (1)")
	}

	epoch := uint32(0)
	err = s.SetCurrentEpoch(epoch)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	// Set currentEpoch in database and then check (another time)
	epoch = 4294967295
	err = s.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	retCE, err = s.GetCurrentEpoch()
	if err != nil {
		t.Fatal(err)
	}
	if retCE != epoch {
		t.Fatal("currentEpochs are not equal (2)")
	}
}
*/
