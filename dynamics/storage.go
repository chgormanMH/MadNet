package dynamics

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/MadBase/MadNet/utils"
	"github.com/sirupsen/logrus"
)

/*
PROPOSAL ON CHAIN
PROPOSAL GETS VOTED ON
IF PROPOSAL PASSES IT BECOMES ACTIVE IN FUTURE ( EPOCH OF ACTIVE > EPOCH OF FINAL VOTE + 1 )
WHEN PROPOSAL PASSES AN EVENT IS EMITTED FROM THE GOVERNANCE CONTRACT
THIS EVENT IS OBSERVED BY THE NODES
THE NODES FETCH THE NEW VALUES AND STORE IN THE DATABASE FOR FUTURE USE
ON THE EPOCH BOUNDARY OF NOT ACTIVE TO ACTIVE, THE STORAGE STRUCT MUST BE UPDATED IN MEMORY FROM
 THE VALUES STORED IN THE DB
*/

// Dynamics contains the list of "constants" which may be changed
// dynamically to reflect protocol updates.
// The point is that these values are essentially constant but may be changed
// in future.

// StorageGetInterface is the interface that all Storage structs must match
// to be valid. These will be used to store the constants which may change
// each epoch as governance determines.
type StorageGetInterface interface {
	GetMaxBytes() uint32
	GetMaxProposalSize() uint32
	GetProposalStepTimeout() time.Duration
	GetPreVoteStepTimeout() time.Duration
	GetPreCommitStepTimout() time.Duration
	GetDeadBlockRoundNextRoundTimeout() time.Duration
	GetDownloadTimeout() time.Duration
	GetSrvrMsgTimeout() time.Duration
	GetMsgTimeout() time.Duration
}

// Storage is the struct which will implement the StorageGetInterface interface.
type Storage struct {
	sync.RWMutex
	database   *Database
	startChan  chan struct{}
	startOnce  sync.Once
	rawStorage *RawStorage // change this out entirely on epoch boundaries
	logger     *logrus.Logger
}

// checkUpdate confirms the specified update is valid.
func checkUpdate(field, value string, epoch uint32) error {
	if epoch == 0 {
		return ErrInvalidUpdateValue
	}
	rs := &RawStorage{}
	err := rs.UpdateValue(field, value)
	if err != nil {
		return err
	}
	return nil
}

// Init initializes the Storage structure.
// TODO: we will need to worry about initialization when not starting
// 		 from the beginning. May need to ability to call someone else.
func (s *Storage) Init(database *Database, logger *logrus.Logger) error {
	s.Lock()
	defer s.Unlock()
	// initialize channel
	s.startChan = make(chan struct{})

	// initialize database
	s.database = database

	// initialize logger
	s.logger = logger

	// Get LinkedList
	var currentEpoch uint32
	linkedList, err := s.database.GetLinkedList()
	if err != nil {
		if !errors.Is(err, ErrKeyNotPresent) {
			utils.DebugTrace(s.logger, err)
			return err
		}
		// We assume we are at the beginning
		// We need to set currentEpoch,
		// begin a new LinkedList and Node,
		// and store this information
		currentEpoch = 1
		rs := &RawStorage{}
		rs.standardParameters()
		node, linkedList, err := CreateLinkedList(currentEpoch, rs)
		if err != nil {
			return err
		}
		if !node.IsHead() || !node.IsTail() {
			// Something is very wrong; initial node should be head and tail
			utils.DebugTrace(s.logger, ErrInvalidNode)
			return ErrInvalidNode
		}
		err = s.database.SetLinkedList(linkedList)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		err = s.database.SetNode(node)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		s.rawStorage, err = node.rawStorage.Copy()
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	} else {
		// No error
		elu := linkedList.GetEpochLastUpdated()
		eluNode, err := s.database.GetNode(elu)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		if !eluNode.IsHead() {
			// Something is very wrong; eluNode should be head
			utils.DebugTrace(s.logger, ErrInvalidNode)
			return ErrInvalidNode
		}
		currentEpoch = linkedList.GetCurrentEpoch()
		rs, err := s.loadStorage(currentEpoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		s.rawStorage, err = rs.Copy()
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	}
	return nil
}

// Start allows normal operations to begin. This MUST be called after Init
// and can only be called once.
func (s *Storage) Start() {
	s.startOnce.Do(func() {
		close(s.startChan)
	})
}

/*
// UpdateStorage updates the database to include changes that must be made
// to the database
func (s *Storage) UpdateStorage(field, value string, epoch uint32) error {
	select {
	case <-s.startChan:
	}

	err := checkUpdate(field, value, epoch)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}

	err = s.updateStorageValue(field, value, epoch)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	return nil
}

// updateStorageValue updates the stored RawStorage values.
//
// There are a few cases that must be handled.
// Throughout, we let E == epoch (from the function argument),
// C == currentEpoch, and H == highestEpoch.
// We have three possibilities:
//
//          E           C                               H
//		|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
//
//                      C               E               H
//		|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
//
//                      C                               H           E
//		|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
//
// Naturally, we also allow for E == C or E == H (and even C == H).
//
// When E <= H, then we need to update epochs starting with max(C, E).
// In this case, we stop updating epochs at H.
// Thus, we have minEpoch = max(C, E) and maxEpoch = H.
// To update, we load RawStorage, update the value, and write RawStorage.
// This begins at minEpoch and ends at maxEpoch.
//
// Otherwise, we have E > H. In this case, we set minEpoch = H+1
// and maxEpoch = E. When updating epochs, no value will have been set,
// so we load RawStorage from the previous epoch and save it
// to the current value. When we reach maxEpoch (== epoch), we load
// the previous epoch, update the value, and then write it.
func (s *Storage) updateStorageValue(field, value string, epoch uint32) error {
	select {
	case <-s.startChan:
	}

	// We now set the lowest epoch which we must change
	var minEpoch uint32
	var maxEpoch uint32
	currentEpoch, err := s.database.GetCurrentEpoch()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	highestEpoch, err := s.database.GetHighestEpoch()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	updateHighestEpoch := false
	if epoch <= highestEpoch {
		// minEpoch == max(epoch, currentEpoch)
		minEpoch = epoch
		if currentEpoch > epoch {
			minEpoch = currentEpoch
		}
		maxEpoch = highestEpoch
	} else {
		// The epoch we need to update is beyond what we stored.
		// We now need to update highestEpoch to reflect this change.
		updateHighestEpoch = true
		minEpoch = highestEpoch + 1
		maxEpoch = epoch
	}

	for epochIter := minEpoch; epochIter <= maxEpoch; epochIter++ {
		rsCurr := &RawStorage{}
		updateValue := true

		if epochIter > epoch {
			// We load RawStorage, update value, and write RawStorage
			rsCurr, err = s.database.GetRawStorage(epochIter)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
		} else if epochIter < epoch {
			// Load RawStorage from previous epoch and write RawStorage
			// to current epoch
			rsCurr, err = s.database.GetRawStorage(epochIter - 1)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
			updateValue = false

		} else if updateHighestEpoch {
			// epochIter == epoch and epoch > highestEpoch
			rsCurr, err = s.database.GetRawStorage(epochIter - 1)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
		} else {
			// epochIter == epoch and epoch <= highestEpoch
			rsCurr, err = s.database.GetRawStorage(epochIter)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
		}

		if updateValue {
			err = rsCurr.UpdateValue(field, value)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
			// DO NOT ACT ON LOCAL VERSION!
			if epochIter == currentEpoch {
				s.rawStorage, err = rsCurr.Copy()
				if err != nil {
					utils.DebugTrace(s.logger, err)
					return err
				}
			}
		}
		err = s.database.SetRawStorage(epochIter, rsCurr)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	}

	if updateHighestEpoch {
		err = s.database.SetHighestEpoch(epoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	}
	return nil
}
*/

// LoadStorage updates RawStorage to the correct value defined by the epoch.
func (s *Storage) LoadStorage(epoch uint32) error {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	rs, err := s.loadStorage(epoch)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	s.rawStorage, err = rs.Copy()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	return nil
}

// loadStorage looks for the appropriate RawStorage value in the database
// and returns that value.
//
// We start at the most updated epoch and proceed backwards until we arrive
// at the node with
//		epoch >= node.thisEpoch
func (s *Storage) loadStorage(epoch uint32) (*RawStorage, error) {
	ll, err := s.database.GetLinkedList()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return nil, err
	}
	elu := ll.GetEpochLastUpdated()
	currentNode, err := s.database.GetNode(elu)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return nil, err
	}

	// Loop backwards through the LinkedList
	for {
		if epoch >= currentNode.thisEpoch {
			rs, err := currentNode.rawStorage.Copy()
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return nil, err
			}
			return rs, nil
		}
		// If we have reached the tail node, then we do not have a node
		// for this specific epoch; we raise an error.
		if currentNode.IsTail() {
			utils.DebugTrace(s.logger, err)
			return nil, ErrInvalid
		}
		// We proceed backward in the linked list of nodes
		prevEpoch := currentNode.prevEpoch
		currentNode, err = s.database.GetNode(prevEpoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return nil, err
		}
	}
}

// addNode adds an additional node to the databae.
// This node can be added anywhere.
// If the node is added at the head, then LinkedList must be updated
// to reflect this change.
func (s *Storage) addNode(node *Node) error {
	select {
	case <-s.startChan:
	}
	s.Lock()
	defer s.Unlock()

	// Ensure node.rawStorage and node.thisEpoch are valid;
	// other parameters should not be set
	if !node.IsPreValid() {
		return ErrInvalid
	}

	// Get LinkedList and Head
	ll, err := s.database.GetLinkedList()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	elu := ll.GetEpochLastUpdated()
	currentNode, err := s.database.GetNode(elu)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}

	if node.thisEpoch > currentNode.thisEpoch {
		// node to be added is strictly ahead of ELU
		err = s.addNodeHead(node, currentNode)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		return nil
	}

	if node.thisEpoch == currentNode.thisEpoch {
		// Node is already present; raise error
		return ErrInvalid
	}

	if currentNode.IsTail() {
		// We are at the end of the LinkedList
		// We need to add node before currentNode
		err = s.addNodeTail(node, currentNode)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		return nil
	}

	prevNode := &Node{}

	// Loop backwards through the LinkedList
	for {
		// Get previous node
		prevNode, err = s.database.GetNode(currentNode.prevEpoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		if prevNode.thisEpoch < node.thisEpoch && node.thisEpoch < currentNode.thisEpoch {
			// We need to add node in between prevNode and currentNode
			err = s.addNodeSplit(node, prevNode, currentNode)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
			return nil
		}
		if node.thisEpoch == prevNode.thisEpoch {
			// Node is already present; raise error
			return ErrInvalid
		}
		if prevNode.IsTail() {
			err = s.addNodeTail(node, prevNode)
			if err != nil {
				utils.DebugTrace(s.logger, err)
				return err
			}
			return nil
		}
		currentNode, err = prevNode.Copy()
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	}
}

func (s *Storage) addNodeHead(node, headNode *Node) error {
	if !node.IsPreValid() || !headNode.IsValid() {
		return ErrInvalid
	}
	if !headNode.IsHead() || node.thisEpoch <= headNode.thisEpoch {
		// We require headNode to be head and node.thisEpoch < headNode.thisEpoch
		return ErrInvalid
	}
	err := node.SetEpochs(headNode, nil)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	// Store the nodes after changes have been made
	err = s.database.SetNode(headNode)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	err = s.database.SetNode(node)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}

	// Update EpochLastUpdated
	ll, err := s.database.GetLinkedList()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	// We need to update EpochLastUpdated
	err = ll.SetEpochLastUpdated(node.thisEpoch)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	err = s.database.SetLinkedList(ll)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	return nil
}

func (s *Storage) addNodeSplit(node, prevNode, nextNode *Node) error {
	if !node.IsPreValid() || !prevNode.IsValid() || !nextNode.IsValid() {
		return ErrInvalid
	}
	if (prevNode.thisEpoch >= node.thisEpoch) || (node.thisEpoch >= nextNode.thisEpoch) {
		return ErrInvalid
	}
	err := node.SetEpochs(prevNode, nextNode)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	// Store the nodes after changes have been made
	err = s.database.SetNode(prevNode)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	err = s.database.SetNode(nextNode)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	err = s.database.SetNode(node)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	return nil
}

func (s *Storage) addNodeTail(node, tailNode *Node) error {
	if !node.IsPreValid() || !tailNode.IsValid() {
		return ErrInvalid
	}
	if !tailNode.IsTail() || node.thisEpoch >= tailNode.thisEpoch {
		// We require tailNode to be tail and node.thisEpoch < tailNode.thisEpoch
		return ErrInvalid
	}
	err := node.SetEpochs(nil, tailNode)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	err = s.database.SetNode(tailNode)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	err = s.database.SetNode(node)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	return nil
}

/*
Methodology for performing Setter:

First, need to know current epoch
	Should probably check to ensure that epoch is in the future.
	TODO: Need to figure out what to do if epoch is in the past!

Second, need to know highest epoch written.
	We will always have currentEpoch <= highestEpoch (highest epoch written)

As noted above, we MUST have
		currentEpoch < changeEpoch
	Unsure what to do if
		changeEpoch <= currentEpoch

if changeEpoch < highestEpoch {
	for epoch in [changeEpoch, ..., highestEpoch]:
		load epoch data
		change appropriate value
		write epoch data
} else {
	copy highestEpoch data
	for epoch in [highestEpoch+1, ..., changeEpoch-1] {

		// We just copy the highestEpoch data to each epoch until right before
		// changeEpoch
		write epoch data
	}
	change appropriate value
	write epoch data for changeEpoch
}

Note on above mentioned TODO:
If a valid Setter has been called and we missed that update because we were
offline for a significant amount of time, then we should make that change
immediately and update current to highest written.
How this could happen in practice (missing stated update) I am not sure,
but this is what could be done to make those changes.
*/

// SetCurrentEpoch sets the current epoch
func (s *Storage) SetCurrentEpoch(epoch uint32) error {
	select {
	case <-s.startChan:
	}
	s.Lock()
	defer s.Unlock()
	ll, err := s.database.GetLinkedList()
	if err != nil {
		return err
	}
	err = ll.SetCurrentEpoch(epoch)
	if err != nil {
		return err
	}
	err = s.database.SetLinkedList(ll)
	if err != nil {
		return err
	}
	return nil
}

// GetCurrentEpoch returns the current epoch
func (s *Storage) GetCurrentEpoch() (uint32, error) {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	ll, err := s.database.GetLinkedList()
	if err != nil {
		return 0, err
	}
	currentEpoch := ll.GetCurrentEpoch()
	return currentEpoch, nil
}

// GetMaxBytes returns the maximum allowed bytes
func (s *Storage) GetMaxBytes() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMaxBytes()
}

// GetMaxProposalSize returns the maximum size of bytes allowed in a proposal
func (s *Storage) GetMaxProposalSize() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMaxProposalSize()
}

// GetSrvrMsgTimeout returns the time before timeout of server message
func (s *Storage) GetSrvrMsgTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetSrvrMsgTimeout()
}

// GetMsgTimeout returns the timeout to receive a message
func (s *Storage) GetMsgTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMsgTimeout()
}

// GetProposalStepTimeout returns the proposal step timeout
func (s *Storage) GetProposalStepTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetProposalStepTimeout()
}

// GetPreVoteStepTimeout returns the prevote step timeout
func (s *Storage) GetPreVoteStepTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetPreVoteStepTimeout()
}

// GetPreCommitStepTimeout returns the precommit step timeout
func (s *Storage) GetPreCommitStepTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetPreCommitStepTimeout()
}

// GetDeadBlockRoundNextRoundTimeout returns the timeout required before
// moving into the DeadBlockRound
func (s *Storage) GetDeadBlockRoundNextRoundTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetDeadBlockRoundNextRoundTimeout()
}

// GetDownloadTimeout returns the timeout for downloads
func (s *Storage) GetDownloadTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetDownloadTimeout()
}

// GetMinTxBurnedFee returns the minimum transaction fee
func (s *Storage) GetMinTxBurnedFee() *big.Int {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMinTxBurnedFee()
}

// GetTxValidVersion returns the transaction valid version
func (s *Storage) GetTxValidVersion() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetTxValidVersion()
}

// GetMinValueStoreBurnedFee returns the minimum transaction fee for ValueStore
func (s *Storage) GetMinValueStoreBurnedFee() *big.Int {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMinValueStoreBurnedFee()
}

// GetValueStoreTxValidVersion returns the ValueStore valid version
func (s *Storage) GetValueStoreTxValidVersion() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetValueStoreTxValidVersion()
}

// GetMinAtomicSwapBurnedFee returns the minimum transaction fee for AtomicSwap
func (s *Storage) GetMinAtomicSwapBurnedFee() *big.Int {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMinAtomicSwapBurnedFee()
}

// GetAtomicSwapValidStopEpoch returns the last epoch at which AtomicSwap is valid
func (s *Storage) GetAtomicSwapValidStopEpoch() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetAtomicSwapValidStopEpoch()
}

// GetDataStoreTxValidVersion returns the DataStore valid version
func (s *Storage) GetDataStoreTxValidVersion() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetDataStoreTxValidVersion()
}
