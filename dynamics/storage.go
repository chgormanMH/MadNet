package dynamics

import (
	"encoding/json"
	"math/big"
	"sync"
	"time"
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
	database         *Database
	startChan        chan struct{}
	startOnce        sync.Once
	*StorageInstance // change this out entirely on epoch boundaries
	currentEpoch     uint32
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
		return err
	}
	if currentEpoch == 0 {
		// currentEpoch has not been initialized.
		// We are starting from the beginning
		currentEpoch = 1
		err := s.database.SetCurrentEpoch(currentEpoch)
		if err != nil {
			return err
		}
		s.currentEpoch = currentEpoch
	} else {
		s.currentEpoch = currentEpoch
	}

	s.StorageInstance = &StorageInstance{}
	si, err := s.database.GetCurrentStorageInstance()
	if err != nil {
		return err
	}
	// ^^^ TODO:
	//	   Should this take in epoch as argument?
	//	   If so, how would be know what the current epoch actually is?
	if si == nil {
		// No StorageInstance present; set standard parameters
		s.StorageInstance.standardParameters()
		err := s.database.SetStorageInstance(currentEpoch, s.StorageInstance)
		if err != nil {
			return err
		}
	} else {
		s.StorageInstance.Overwrite(si)
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

// UpdateStorageInstance updates StorageInstance to the correct value
// defined by the epoch.
//
// This is not yet implemented and will probably be involved.
func (s *Storage) UpdateStorageInstance(epoch uint32) error {
	s.Lock()
	defer s.Unlock()

	// Check for any updates to parameters that have not been called;
	//		if some exist, perform those updates.
	s.CheckForUpdates()

	// Search for StorageInstance for epoch at correct location.
	si, err := s.database.GetStorageInstance(epoch)
	if err != nil {
		return err
	}
	if si == nil {
		// Not present; continue using current one and store it
		err := s.database.SetStorageInstance(epoch, si)
		if err != nil {
			return err
		}
	} else {
		err := s.StorageInstance.Overwrite(si)
		if err != nil {
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
	return s.StorageInstance.GetMaxBytes()
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
	return s.StorageInstance.GetMaxProposalSize()
}

/*
// SetMaxProposalSize sets the maximum size of bytes allowed in a proposal
func (s *Storage) SetMaxProposalSize(value uint32, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.MaxProposalSize = value
	return nil
}
*/

// GetSrvrMsgTimeout returns the time before timeout of server message
func (s *Storage) GetSrvrMsgTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.StorageInstance.GetSrvrMsgTimeout()
}

/*
// SetSrvrMsgTimeout sets the time before timeout of server message
func (s *Storage) SetSrvrMsgTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.SrvrMsgTimeout = value
	return nil
}
*/

// GetMsgTimeout returns the timeout to receive a message
func (s *Storage) GetMsgTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.StorageInstance.GetMsgTimeout()
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
	return s.StorageInstance.GetProposalStepTimeout()
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
	return s.StorageInstance.GetPreVoteStepTimeout()
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
	return s.StorageInstance.GetPreCommitStepTimeout()
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
	return s.StorageInstance.GetDeadBlockRoundNextRoundTimeout()
}

/*
// SetDeadBlockRoundNextRoundTimeout sets the timeout required before
// moving into the DeadBlockRound
func (s *Storage) SetDeadBlockRoundNextRoundTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.DeadBlockRoundNextRoundTimeout = value
	return nil
}
*/

// GetDownloadTimeout returns the timeout for downloads
func (s *Storage) GetDownloadTimeout() time.Duration {
	select {
	case <-s.startChan:
	}
	s.RLock()
	defer s.RUnlock()
	return s.StorageInstance.GetDownloadTimeout()
}

/*
// SetDownloadTimeout sets the timeout for downloads
func (s *Storage) SetDownloadTimeout(value time.Duration, epoch uint32) error {
	panic("not implemented")
	s.Lock()
	defer s.Unlock()
	s.DownloadTimeout = value
	return nil
}
*/

// StorageInstance is the struct which actually
type StorageInstance struct {
	MaxBytes                       uint32        `json:"maxBytes,omitempty"`
	MaxProposalSize                uint32        `json:"maxProposalSize,omitempty"`
	ProposalStepTimeout            time.Duration `json:"proposalStepTimeout,omitempty"`
	PreVoteStepTimeout             time.Duration `json:"preVoteStepTimeout,omitempty"`
	PreCommitStepTimout            time.Duration `json:"preCommitStepTimeout,omitempty"`
	DeadBlockRoundNextRoundTimeout time.Duration `json:"deadBlockRoundNextRoundTimeout,omitempty"`
	DownloadTimeout                time.Duration `json:"downloadTimeout,omitempty"`
	SrvrMsgTimeout                 time.Duration `json:"srvrMsgTimeout,omitempty"`
	MsgTimeout                     time.Duration `json:"msgTimeout,omitempty"`

	MinTxBurnedFee *big.Int `json:"minTxBurnedFee,omitempty"`
	TxValidVersion uint32   `json:"txValidVersion,omitempty"`

	MinValueStoreBurnedFee   *big.Int `json:"minValueStoreBurnedFee,omitempty"` // TODO: Do we need to worry about storing variable-size values?
	ValueStoreTxValidVersion uint32   `json:"valueStoreTxValidVersion,omitempty"`

	MinAtomicSwapBurnedFee   *big.Int `json:"minAtomicSwapBurnedFee,omitempty"` // TODO: Do we need to worry about storing variable-size values?
	AtomicSwapValidStopEpoch uint32   `json:"atomicSwapValidStopEpoch,omitempty"`

	DataStoreTxValidVersion uint32 `json:"dataStoreTxValidVersion,omitempty"`
}

// Marshal performs json.Marshal on the StorageInstance struct.
func (si *StorageInstance) Marshal() ([]byte, error) {
	return json.Marshal(si)
}

// Unmarshal performs json.Unmarshal on the StorageInstance struct.
func (si *StorageInstance) Unmarshal(v []byte) error {
	if si == nil {
		return ErrStorageInstanceNilPointer
	}
	return json.Unmarshal(v, si)
}

// Copy makes a complete copy of StorageInstance struct.
func (si *StorageInstance) Copy() (*StorageInstance, error) {
	siBytes, err := si.Marshal()
	if err != nil {
		return nil, err
	}
	c := &StorageInstance{}
	err = c.Unmarshal(siBytes)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Overwrite replaces the current StorageGetInterface contents with the copy
// taken as an argument
func (si *StorageInstance) Overwrite(c *StorageInstance) error {
	cBytes, err := c.Marshal()
	if err != nil {
		return err
	}
	err = si.Unmarshal(cBytes)
	if err != nil {
		return err
	}
	return nil
}

// standardParameters initializes StorageInstance with the standard (original)
// parameters for the system.
func (si *StorageInstance) standardParameters() {
	si.MaxBytes = maxBytes
	si.MaxProposalSize = maxProposalSize
	si.ProposalStepTimeout = proposalStepTO
	si.PreVoteStepTimeout = preVoteStepTO
	si.PreCommitStepTimout = preCommitStepTO
	si.DeadBlockRoundNextRoundTimeout = dBRNRTO
	si.DownloadTimeout = downloadTO
	si.SrvrMsgTimeout = srvrMsgTimeout
	si.MsgTimeout = msgTimeout
}

// GetMaxBytes returns the maximum allowed bytes
func (si *StorageInstance) GetMaxBytes() uint32 {
	return si.MaxBytes
}

/*
// SetMaxBytes sets the maximum allowed bytes
func (d *StorageInstance) SetMaxBytes(value uint32, epoch uint32) error {
	d.MaxBytes = value
	return nil
}
*/

// GetMaxProposalSize returns the maximum size of bytes allowed in a proposal
func (si *StorageInstance) GetMaxProposalSize() uint32 {
	return si.MaxProposalSize
}

/*
// SetMaxProposalSize sets the maximum size of bytes allowed in a proposal
func (d *StorageInstance) SetMaxProposalSize(value uint32, epoch uint32) error {
	d.MaxProposalSize = value
	return nil
}
*/

// GetSrvrMsgTimeout returns the time before timeout of server message
func (si *StorageInstance) GetSrvrMsgTimeout() time.Duration {
	return si.SrvrMsgTimeout
}

/*
// SetSrvrMsgTimeout sets the time before timeout of server message
func (d *StorageInstance) SetSrvrMsgTimeout(value time.Duration, epoch uint32) error {
	d.SrvrMsgTimeout = value
	return nil
}
*/

// GetMsgTimeout returns the timeout to receive a message
func (si *StorageInstance) GetMsgTimeout() time.Duration {
	return si.MsgTimeout
}

/*
// SetMsgTimeout sets the timeout to receive a message
func (d *StorageInstance) SetMsgTimeout(value time.Duration, epoch uint32) error {
	d.MsgTimeout = value
	return nil
}
*/

// GetProposalStepTimeout returns the proposal step timeout
func (si *StorageInstance) GetProposalStepTimeout() time.Duration {
	return si.ProposalStepTimeout
}

/*
// SetProposalStepTimeout sets the proposal step timeout
func (d *StorageInstance) SetProposalStepTimeout(value time.Duration, epoch uint32) error {
	d.ProposalStepTimeout = value
	return nil
}
*/

// GetPreVoteStepTimeout returns the prevote step timeout
func (si *StorageInstance) GetPreVoteStepTimeout() time.Duration {
	return si.PreVoteStepTimeout
}

/*
// SetPreVoteStepTimeout sets the prevote step timeout
func (d *StorageInstance) SetPreVoteStepTimeout(value time.Duration, epoch uint32) error {
	d.PreVoteStepTimeout = value
	return nil
}
*/

// GetPreCommitStepTimeout returns the precommit step timeout
func (si *StorageInstance) GetPreCommitStepTimeout() time.Duration {
	return si.PreCommitStepTimout
}

/*
// SetPreCommitStepTimeout sets the precommit step timeout
func (d *StorageInstance) SetPreCommitStepTimeout(value time.Duration, epoch uint32) error {
	d.PreCommitStepTimout = value
	return nil
}
*/

// GetDeadBlockRoundNextRoundTimeout returns the timeout required before
// moving into the DeadBlockRound
func (si *StorageInstance) GetDeadBlockRoundNextRoundTimeout() time.Duration {
	return si.DeadBlockRoundNextRoundTimeout
}

/*
// SetDeadBlockRoundNextRoundTimeout sets the timeout required before
// moving into the DeadBlockRound
func (d *StorageInstance) SetDeadBlockRoundNextRoundTimeout(value time.Duration, epoch uint32) error {
	d.DeadBlockRoundNextRoundTimeout = value
	return nil
}
*/

// GetDownloadTimeout returns the timeout for downloads
func (si *StorageInstance) GetDownloadTimeout() time.Duration {
	return si.DownloadTimeout
}

/*
// SetDownloadTimeout sets the timeout for downloads
func (d *StorageInstance) SetDownloadTimeout(value time.Duration, epoch uint32) error {
	d.DownloadTimeout = value
	return nil
}
*/
