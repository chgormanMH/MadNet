package constants

import "time"

// Consensus params
const (
	DEADBLOCKROUND   uint32 = 5
	DEADBLOCKROUNDNR        = DEADBLOCKROUND - 1
	MaxBytes                = 3000000
	MaxProposalSize         = MaxBytes
	SrvrMsgTimeout          = 3 * time.Second // Do not go lower than 2 seconds!
	MsgTimeout              = 4 * time.Second // Do not go lower than 2 seconds!
	ProposalStepTO          = 4 * time.Second //4 * time.Second
	PreVoteStepTO           = 3 * time.Second //4 * time.Second
	PreCommitStepTO         = 3 * time.Second //4 * time.Second
	DBRNRTO                 = 24 * time.Second
	DownloadTO              = ProposalStepTO + PreVoteStepTO + PreCommitStepTO
)

// AdminHandlerKid returns a constant byte slice to be used as Key ID
func AdminHandlerKid() []byte {
	return []byte("constant")
}

// StorageCurrentEpoch returns a constant byte slice to be used for storing
// current epoch
func StorageCurrentEpoch() []byte {
	return []byte("CurrentEpoch")
}

// StorageHighestEpoch returns a constant byte slice to be used for storing
// highest epoch
func StorageHighestEpoch() []byte {
	return []byte("HighestEpoch")
}
