package dynamics

import (
	"errors"
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
func (s *Storage) Init(database *Database, logger *logrus.Logger) error {
	// initialize channel
	s.startChan = make(chan struct{})

	// initialize database
	s.database = database

	// initialize logger
	s.logger = logger

	s.rawStorage = &RawStorage{}
	// Get currentEpoch and associated rawStorage values
	currentEpoch, err := s.database.GetCurrentEpoch()
	if err != nil {
		if !errors.Is(err, ErrKeyNotPresent) {
			utils.DebugTrace(s.logger, err)
			return err
		}
		// currentEpoch is not set;
		// we load standard parameters
		currentEpoch = 1
		s.currentEpoch = currentEpoch
		err = s.database.SetCurrentEpoch(s.currentEpoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
		s.rawStorage.standardParameters()
		err = s.database.SetRawStorage(s.currentEpoch, s.rawStorage)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	} else {
		// No error
		s.currentEpoch = currentEpoch
		rs, err := s.database.GetRawStorage(s.currentEpoch)
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

	// Look for highestEpoch and set its value if necessary
	highestEpoch, err := s.database.GetHighestEpoch()
	if err != nil {
		if !errors.Is(err, ErrKeyNotPresent) {
			utils.DebugTrace(s.logger, err)
			return err
		}
		highestEpoch = currentEpoch
		err = s.database.SetHighestEpoch(highestEpoch)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	}
	// Ensure highestEpoch >= currentEpoch
	if highestEpoch < currentEpoch {
		highestEpoch = currentEpoch
		err = s.database.SetHighestEpoch(highestEpoch)
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

// CheckForUpdates looks for updates to system parameters
func (s *Storage) CheckForUpdates() error {
	select {
	case <-s.startChan:
	}
	return nil
}

// UpdateStorage updates the database to include changes that must be made
// to the database
func (s *Storage) UpdateStorage() error {
	select {
	case <-s.startChan:
	}
	return nil
}

// UpdateStorageValue ...
func (s *Storage) UpdateStorageValue(field, value string, epoch uint32) error {
	select {
	case <-s.startChan:
	}
	lowestEpoch := epoch
	if lowestEpoch < s.currentEpoch {
		lowestEpoch = s.currentEpoch
	}
	highestEpoch, err := s.database.GetHighestEpoch()
	if err != nil {
		return err
	}
	if epoch > highestEpoch {
		// We now need to update highestEpoch to reflect this change
		highestEpoch = epoch
	}
	for epochCurr := lowestEpoch; epochCurr <= highestEpoch; epochCurr++ {
		rsCurr, err := s.database.GetRawStorage(epochCurr)
		if err != nil {
			return err
		}
		err = rsCurr.UpdateValue(field, value)
		if err != nil {
			return err
		}
		err = s.database.SetRawStorage(epochCurr, rsCurr)
		if err != nil {
			return err
		}
	}
	err = s.database.SetHighestEpoch(highestEpoch)
	if err != nil {
		return err
	}
	return nil
}

// LoadStorage updates RawStorage to the correct value
// defined by the epoch.
func (s *Storage) LoadStorage(epoch uint32) error {
	select {
	case <-s.startChan:
	}
	s.Lock()
	defer s.Unlock()

	// Check for any updates to parameters that have not been called;
	// if some exist, perform those updates.
	err := s.CheckForUpdates()
	if err != nil {
		utils.DebugTrace(s.logger, err)
		return err
	}

	// Search for RawStorage for epoch at correct location
	rs, err := s.database.GetRawStorage(epoch)
	if err != nil {
		if !errors.Is(err, ErrKeyNotPresent) {
			utils.DebugTrace(s.logger, err)
			return err
		}
		// There is no current rawStorage for the specified epoch;
		// continue to use current rawStorage
		err = s.database.SetRawStorage(epoch, s.rawStorage)
		if err != nil {
			utils.DebugTrace(s.logger, err)
			return err
		}
	} else {
		s.rawStorage, err = rs.Copy()
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
