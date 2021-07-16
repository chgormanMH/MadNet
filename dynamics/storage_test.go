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

func TestStorageAddNodeHeadGood(t *testing.T) {
	origEpoch := uint32(1)
	s := initializeStorageCE(origEpoch)
	headNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if !headNode.IsValid() {
		t.Fatal("headNode should be valid")
	}

	newEpoch := origEpoch + 1
	rsNew := &RawStorage{}
	rsNew.standardParameters()
	rsBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567
	node := &Node{
		thisEpoch:  newEpoch,
		rawStorage: rsNew,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be prevalid")
	}
	rsNewBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	err = s.addNodeHead(node, headNode)
	if err != nil {
		t.Fatal(err)
	}

	origNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if origNode.prevEpoch != origEpoch || origNode.thisEpoch != origEpoch || origNode.nextEpoch != newEpoch {
		t.Fatal("origNode invalid (1)")
	}
	rsOrigBytes, err := origNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsOrigBytes, rsBytes) {
		t.Fatal("origNode invalid (2)")
	}

	retNode, err := s.database.GetNode(newEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if retNode.prevEpoch != origEpoch || retNode.thisEpoch != newEpoch || retNode.nextEpoch != newEpoch {
		t.Fatal("retNode invalid (1)")
	}
	retBytes, err := retNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retBytes, rsNewBytes) {
		t.Fatal("retNode invalid (2)")
	}
}

func TestStorageAddNodeHeadBad1(t *testing.T) {
	origEpoch := uint32(1)
	s := initializeStorageCE(origEpoch)
	headNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if !headNode.IsValid() {
		t.Fatal("headNode should be valid")
	}

	node := &Node{}
	if node.IsPreValid() {
		t.Fatal("node should not be prevalid")
	}

	err = s.addNodeHead(node, headNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeHeadBad2(t *testing.T) {
	origEpoch := uint32(1)
	s := initializeStorageCE(origEpoch)
	headNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if !headNode.IsValid() {
		t.Fatal("headNode should be valid")
	}

	rs := &RawStorage{}
	node := &Node{
		thisEpoch:  1,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be prevalid")
	}

	err = s.addNodeHead(node, headNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeTailGood(t *testing.T) {
	origEpoch := uint32(10)
	s := initializeStorageCE(origEpoch)
	tailNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if !tailNode.IsValid() {
		t.Fatal("tailNode should be valid")
	}

	newEpoch := origEpoch - 1
	rsNew := &RawStorage{}
	rsNew.standardParameters()
	rsBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567
	node := &Node{
		thisEpoch:  newEpoch,
		rawStorage: rsNew,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be prevalid")
	}
	rsNewBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	err = s.addNodeTail(node, tailNode)
	if err != nil {
		t.Fatal(err)
	}

	origNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if origNode.prevEpoch != newEpoch || origNode.thisEpoch != origEpoch || origNode.nextEpoch != origEpoch {
		t.Fatal("origNode invalid (1)")
	}
	rsOrigBytes, err := origNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsOrigBytes, rsBytes) {
		t.Fatal("origNode invalid (2)")
	}

	retNode, err := s.database.GetNode(newEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if retNode.prevEpoch != newEpoch || retNode.thisEpoch != newEpoch || retNode.nextEpoch != origEpoch {
		t.Fatal("retNode invalid (1)")
	}
	retBytes, err := retNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retBytes, rsNewBytes) {
		t.Fatal("retNode invalid (2)")
	}
}

func TestStorageAddNodeTailBad1(t *testing.T) {
	origEpoch := uint32(10)
	s := initializeStorageCE(origEpoch)
	tailNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if !tailNode.IsValid() {
		t.Fatal("tailNode should be valid")
	}

	node := &Node{}
	if node.IsPreValid() {
		t.Fatal("node should not be prevalid")
	}

	err = s.addNodeTail(node, tailNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeTailBad2(t *testing.T) {
	origEpoch := uint32(10)
	s := initializeStorageCE(origEpoch)
	tailNode, err := s.database.GetNode(origEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if !tailNode.IsValid() {
		t.Fatal("tailNode should be valid")
	}

	rs := &RawStorage{}
	node := &Node{
		thisEpoch:  origEpoch,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should not be prevalid")
	}

	err = s.addNodeTail(node, tailNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeSplitGood(t *testing.T) {
	first := uint32(1)
	s := initializeStorageCE(first)
	prevNode, err := s.database.GetNode(first)
	if err != nil {
		t.Fatal(err)
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be valid")
	}
	last := uint32(10)
	prevNode.nextEpoch = last
	err = s.database.SetNode(prevNode)
	if err != nil {
		t.Fatal(err)
	}

	// Set up nextnode
	rs := &RawStorage{}
	rs.standardParameters()
	rsBytes, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	nextNode := &Node{
		prevEpoch:  first,
		thisEpoch:  last,
		nextEpoch:  last,
		rawStorage: rs,
	}
	if !nextNode.IsValid() {
		t.Fatal("nextNode should be valid")
	}
	err = s.database.SetNode(nextNode)
	if err != nil {
		t.Fatal(err)
	}

	// Set up node
	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 123456
	rsNewBytes, err := rsNew.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	newEpoch := uint32(5)
	node := &Node{
		thisEpoch:  newEpoch,
		rawStorage: rsNew,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}

	err = s.addNodeSplit(node, prevNode, nextNode)
	if err != nil {
		t.Fatal(err)
	}

	// Check everything
	firstNode, err := s.database.GetNode(first)
	if err != nil {
		t.Fatal(err)
	}
	if firstNode.prevEpoch != first || firstNode.thisEpoch != first || firstNode.nextEpoch != newEpoch {
		t.Fatal("firstNode invalid (1)")
	}
	retBytes, err := firstNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retBytes, rsBytes) {
		t.Fatal("firstNode invalid (2)")
	}

	middleNode, err := s.database.GetNode(newEpoch)
	if err != nil {
		t.Fatal(err)
	}
	if middleNode.prevEpoch != first || middleNode.thisEpoch != newEpoch || middleNode.nextEpoch != last {
		t.Fatal("middleNode invalid (1)")
	}
	retBytes, err = middleNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retBytes, rsNewBytes) {
		t.Fatal("middleNode invalid (2)")
	}

	lastNode, err := s.database.GetNode(last)
	if err != nil {
		t.Fatal(err)
	}
	if lastNode.prevEpoch != newEpoch || lastNode.thisEpoch != last || lastNode.nextEpoch != last {
		t.Fatal("lastNode invalid (1)")
	}
	retBytes, err = lastNode.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retBytes, rsBytes) {
		t.Fatal("lastNode invalid (2)")
	}
}

func TestStorageAddNodeSplitBad1(t *testing.T) {
	first := uint32(1)
	s := initializeStorageCE(first)
	prevNode, err := s.database.GetNode(first)
	if err != nil {
		t.Fatal(err)
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be valid")
	}
	last := uint32(10)
	prevNode.nextEpoch = last
	err = s.database.SetNode(prevNode)
	if err != nil {
		t.Fatal(err)
	}

	// Set up nodes
	rs := &RawStorage{}
	rs.standardParameters()
	nextNode := &Node{
		prevEpoch:  first,
		thisEpoch:  last,
		nextEpoch:  last,
		rawStorage: rs,
	}
	if !nextNode.IsValid() {
		t.Fatal("nextNode should be valid")
	}
	err = s.database.SetNode(nextNode)
	if err != nil {
		t.Fatal(err)
	}

	node := &Node{}

	err = s.addNodeSplit(node, prevNode, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestStorageAddNodeSplitBad2(t *testing.T) {
	first := uint32(1)
	s := initializeStorageCE(first)
	prevNode, err := s.database.GetNode(first)
	if err != nil {
		t.Fatal(err)
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be valid")
	}
	last := uint32(10)
	prevNode.nextEpoch = last
	err = s.database.SetNode(prevNode)
	if err != nil {
		t.Fatal(err)
	}

	// Set up nodes
	rs := &RawStorage{}
	rs.standardParameters()
	nextNode := &Node{
		prevEpoch:  first,
		thisEpoch:  last,
		nextEpoch:  last,
		rawStorage: rs,
	}
	if !nextNode.IsValid() {
		t.Fatal("nextNode should be valid")
	}
	err = s.database.SetNode(nextNode)
	if err != nil {
		t.Fatal(err)
	}

	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 123456
	newEpoch := uint32(100)
	node := &Node{
		thisEpoch:  newEpoch,
		rawStorage: rsNew,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}

	err = s.addNodeSplit(node, prevNode, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// Test addNode when adding to Head
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

// Test addNode when adding to behind Head == Tail
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

// Test addNode when adding to behind Head == Tail, and then behind new Tail
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

// Test addNode when adding to behind Head == Tail, and then
// in between Tail and Head
func TestStorageAddNodeGood4(t *testing.T) {
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
	addedEpoch := uint32(1)
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
	addedEpoch2 := uint32(10)
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
	if origNode.prevEpoch != addedEpoch2 {
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
	if addedNode.prevEpoch != addedEpoch {
		t.Fatal("addedNode.prevEpoch is invalid")
	}
	if addedNode.thisEpoch != addedEpoch {
		t.Fatal("addedNode.thisEpoch is invalid")
	}
	if addedNode.nextEpoch != addedEpoch2 {
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
	if addedNode2.prevEpoch != addedEpoch {
		t.Fatal("addedNode2.prevEpoch is invalid")
	}
	if addedNode2.thisEpoch != addedEpoch2 {
		t.Fatal("addedNode2.thisEpoch is invalid")
	}
	if addedNode2.nextEpoch != origEpoch {
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

func TestStorageAddNodeGood5(t *testing.T) {
	origEpoch := uint32(1000)
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
	addedEpoch := uint32(100)
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

	rs.standardParameters()
	addedEpoch3 := uint32(10)
	newNode3 := &Node{
		prevEpoch:  0,
		thisEpoch:  addedEpoch3,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err = s.addNode(newNode3)
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
	if addedNode.prevEpoch != addedEpoch3 {
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
	if addedNode2.nextEpoch != addedEpoch3 {
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
	if !bytes.Equal(retRSBytes, rsNewBytes) {
		t.Fatal("invalid RawStorage")
	}

	addedNode3, err := s.database.GetNode(addedEpoch3)
	if err != nil {
		t.Fatal(err)
	}
	if addedNode3.prevEpoch != addedEpoch2 {
		t.Fatal("addedNode3.prevEpoch is invalid")
	}
	if addedNode3.thisEpoch != addedEpoch3 {
		t.Fatal("addedNode3.thisEpoch is invalid")
	}
	if addedNode3.nextEpoch != addedEpoch {
		t.Fatal("addedNode3.nextEpoch is invalid")
	}
	retRS, err = addedNode3.rawStorage.Copy()
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

func TestStorageAddNodeBad3(t *testing.T) {
	origEpoch := uint32(10)
	s := initializeStorageCE(origEpoch)
	rs := &RawStorage{}
	newNode := &Node{
		prevEpoch:  0,
		thisEpoch:  1,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err := s.addNode(newNode)
	if err != nil {
		t.Fatal(err)
	}

	// Add same epoch again; should raise error
	newNode2 := &Node{
		prevEpoch:  0,
		thisEpoch:  1,
		nextEpoch:  0,
		rawStorage: rs,
	}
	err = s.addNode(newNode2)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

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
