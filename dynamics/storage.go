package dynamics

import (
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

// TODO: Will need a data structure to store more than just the most recent
//		 parameter value.

// StorageGetInterface is the interface that all Storage structs must match
// to be valid. These will be used to store the constants which may change
// each epoch as governance determines.
type StorageGetInterface interface {
	GetMaxBytes() uint32
	GetMaxProposalSize() uint32
	GetProposalStepTimeout() time.Duration
	GetPreVoteStepTimeout() time.Duration
	GetPreCommitStepTimout() time.Duration
	DeadBlockRoundNextRoundTimeout() time.Duration
	DownloadTimeout() time.Duration
	SrvrMsgTimeout() time.Duration
	MsgTimeout() time.Duration
}

// Storage is the struct which will implement the StorageGetInterface interface.
type Storage struct {
	sync.RWMutex
	database     *Database
	startChan    chan struct{}
	startOnce    sync.Once
	rawStorage   *RawStorage // change this out entirely on epoch boundaries
	currentEpoch uint32
	logger       *logrus.Logger
}

// Init initializes the Storage structure.
func (s *Storage) Init(database *Database) error {
	// initialize channel
	s.startChan = make(chan struct{})

	// initialize database
	s.database = database

	// Update currentEpoch and highest written to reflect this.
	currentEpoch, err := s.database.GetCurrentEpoch()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	if currentEpoch == 0 {
		// currentEpoch has not been initialized.
		// We are starting from the beginning
		currentEpoch = 1
		err := s.database.SetCurrentEpoch(currentEpoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		s.currentEpoch = currentEpoch
	} else {
		s.currentEpoch = currentEpoch
	}

	s.rawStorage = &RawStorage{}
	rs, err := s.database.GetCurrentRawStorage()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	// ^^^ TODO:
	//	   Should this take in epoch as argument?
	//	   If so, how would be know what the current epoch actually is?
	if rs == nil {
		// No RawStorage present; set standard parameters
		s.rawStorage.standardParameters()
		err := s.database.SetRawStorage(currentEpoch, s.rawStorage)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	} else {
		s.rawStorage.Overwrite(rs)
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

// CheckForUpdates looks for updates to system parameters
func (s *Storage) CheckForUpdates() error {
	return nil
}

// UpdateStorageInstance updates RawStorage to the correct value
// defined by the epoch.
func (s *Storage) UpdateStorageInstance(epoch uint32) error {
	s.Lock()
	defer s.Unlock()

	// Check for any updates to parameters that have not been called;
	//		if some exist, perform those updates.
	s.CheckForUpdates()

	// Search for RawStorage for epoch at correct location.
	rs, err := s.database.GetRawStorage(epoch)
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}
	if rs == nil {
		// Not present; continue using current one and store it
		err := s.database.SetRawStorage(epoch, rs)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	} else {
		err := s.rawStorage.Overwrite(rs)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
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

// GetMaxBytes returns the maximum allowed bytes
func (s *Storage) GetMaxBytes() uint32 {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetMaxBytes()
}

/*
// SetMaxBytes sets the maximum allowed bytes
func (s *Storage) SetMaxBytes(value uint32, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.MaxBytes = value
	return nil
}
*/

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

/*
// SetMsgTimeout sets the timeout to receive a message
func (s *Storage) SetMsgTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.MsgTimeout = value
	return nil
}
*/

// GetProposalStepTimeout returns the proposal step timeout
func (s *Storage) GetProposalStepTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetProposalStepTimeout()
}

/*
// SetProposalStepTimeout sets the proposal step timeout
func (s *Storage) SetProposalStepTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.ProposalStepTimeout = value
	return nil
}
*/

// GetPreVoteStepTimeout returns the prevote step timeout
func (s *Storage) GetPreVoteStepTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetPreVoteStepTimeout()
}

/*
// SetPreVoteStepTimeout sets the prevote step timeout
func (s *Storage) SetPreVoteStepTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.PreVoteStepTimeout = value
	return nil
}
*/

// GetPreCommitStepTimeout returns the precommit step timeout
func (s *Storage) GetPreCommitStepTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.rawStorage.GetPreCommitStepTimeout()
}

/*
// SetPreCommitStepTimeout sets the precommit step timeout
func (s *Storage) SetPreCommitStepTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.PreCommitStepTimout = value
	return nil
}
*/

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
