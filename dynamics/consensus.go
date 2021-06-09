package dynamics

import (
	"sync"
	"time"
)

// Copied from constants/consensus.go
const (
	//dEADBLOCKROUND   uint32 = 5
	//dEADBLOCKROUNDNR        = dEADBLOCKROUND - 1
	maxBytes        = 3000000
	maxProposalSize = maxBytes
	srvrMsgTimeout  = 3 * time.Second  // Do not go lower than 2 seconds!
	msgTimeout      = 4 * time.Second  // Do not go lower than 2 seconds!
	proposalStepTO  = 4 * time.Second  // 4 * time.Second
	preVoteStepTO   = 3 * time.Second  // 4 * time.Second
	preCommitStepTO = 3 * time.Second  // 4 * time.Second
	dBRNRTO         = 24 * time.Second // Should be a significant wait before entering
	downloadTO      = proposalStepTO + preVoteStepTO + preCommitStepTO
)

// Dynamics contains the list of "constants" which may be changed
// dynamically to reflect protocol updates.
// The point is that these values are essentially constant but may be changed
// in future.

// TODO: Will need a data structure to store more than just the most recent
//		 parameter value.

type storage struct {
	sync.RWMutex
	MaxBytes                       uint32 // Done
	MaxProposalSize                uint32 // Done
	SrvrMsgTimeout                 time.Duration
	MsgTimeout                     time.Duration
	ProposalStepTimeout            time.Duration // Done
	PreVoteStepTimeout             time.Duration // Done
	PreCommitStepTimout            time.Duration // Done
	DeadBlockRoundNextRoundTimeout time.Duration // Done
	DownloadTimeout                time.Duration
}

var d storage

func init() {
	firstBlock := uint32(1)
	err := SetMaxBytes(maxBytes, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetMaxProposalSize(maxProposalSize, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetSrvrMsgTimeout(srvrMsgTimeout, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetMsgTimeout(msgTimeout, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetProposalStepTimeout(proposalStepTO, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetPreVoteStepTimeout(preVoteStepTO, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetPreCommitStepTimeout(preCommitStepTO, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetDeadBlockRoundNextRoundTimeout(dBRNRTO, firstBlock)
	if err != nil {
		panic(err)
	}
	err = SetDownloadTimeout(downloadTO, firstBlock)
	if err != nil {
		panic(err)
	}
}

// GetMaxBytes returns the maximum allowed bytes
func GetMaxBytes(height uint32) uint32 {
	d.RLock()
	defer d.RUnlock()
	return d.MaxBytes
}

// SetMaxBytes sets the maximum allowed bytes
func SetMaxBytes(value uint32, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.MaxBytes = value
	return nil
}

// GetMaxProposalSize returns the maximum size of bytes allowed in a proposal
func GetMaxProposalSize(height uint32) uint32 {
	d.RLock()
	defer d.RUnlock()
	return d.MaxProposalSize
}

// SetMaxProposalSize sets the maximum size of bytes allowed in a proposal
func SetMaxProposalSize(value uint32, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.MaxProposalSize = value
	return nil
}

// GetSrvrMsgTimeout returns the time before timeout of server message
func GetSrvrMsgTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.SrvrMsgTimeout
}

// SetSrvrMsgTimeout sets the time before timeout of server message
func SetSrvrMsgTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.SrvrMsgTimeout = value
	return nil
}

// GetMsgTimeout returns the timeout to receive a message
func GetMsgTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.MsgTimeout
}

// SetMsgTimeout sets the timeout to receive a message
func SetMsgTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.MsgTimeout = value
	return nil
}

// GetProposalStepTimeout returns the proposal step timeout
func GetProposalStepTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.ProposalStepTimeout
}

// SetProposalStepTimeout sets the proposal step timeout
func SetProposalStepTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.ProposalStepTimeout = value
	return nil
}

// GetPreVoteStepTimeout returns the prevote step timeout
func GetPreVoteStepTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.PreVoteStepTimeout
}

// SetPreVoteStepTimeout sets the prevote step timeout
func SetPreVoteStepTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.PreVoteStepTimeout = value
	return nil
}

// GetPreCommitStepTimeout returns the precommit step timeout
func GetPreCommitStepTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.PreCommitStepTimout
}

// SetPreCommitStepTimeout sets the precommit step timeout
func SetPreCommitStepTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.PreCommitStepTimout = value
	return nil
}

// GetDeadBlockRoundNextRoundTimeout returns the timeout required before
// moving into the DeadBlockRound
func GetDeadBlockRoundNextRoundTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.DeadBlockRoundNextRoundTimeout
}

// SetDeadBlockRoundNextRoundTimeout sets the timeout required before
// moving into the DeadBlockRound
func SetDeadBlockRoundNextRoundTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.DeadBlockRoundNextRoundTimeout = value
	return nil
}

// GetDownloadTimeout returns the timeout for downloads
func GetDownloadTimeout(height uint32) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.DownloadTimeout
}

// SetDownloadTimeout sets the timeout for downloads
func SetDownloadTimeout(value time.Duration, height uint32) error {
	d.RLock()
	defer d.RUnlock()
	d.DownloadTimeout = value
	return nil
}
