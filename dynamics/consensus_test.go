package dynamics

import "testing"

func TestGetMaxBytes(t *testing.T) {
	height := uint32(1)
	maxBytesReturned := GetMaxBytes(height)
	if maxBytes != maxBytesReturned {
		t.Fatal("Invalid MaxBytes returned")
	}
}

func TestGetMaxProposalSize(t *testing.T) {
	height := uint32(1)
	maxProposalSizeReturned := GetMaxProposalSize(height)
	if maxProposalSize != maxProposalSizeReturned {
		t.Fatal("Invalid MaxProposalSize returned")
	}
}

func TestGetSrvrMsgTimeout(t *testing.T) {
	height := uint32(1)
	srvrMsgTimeoutReturned := GetSrvrMsgTimeout(height)
	if srvrMsgTimeout != srvrMsgTimeoutReturned {
		t.Fatal("Invalid SrvrMsgTimeout returned")
	}
}

func TestGetMsgTimeout(t *testing.T) {
	height := uint32(1)
	msgTimeoutReturned := GetMsgTimeout(height)
	if msgTimeout != msgTimeoutReturned {
		t.Fatal("Invalid MsgTimeout returned")
	}
}

func TestGetProposalStepTimeout(t *testing.T) {
	height := uint32(1)
	proposalStepTimeoutReturned := GetProposalStepTimeout(height)
	if proposalStepTO != proposalStepTimeoutReturned {
		t.Fatal("Invalid ProposalStepTimeout returned")
	}
}

func TestGetPreVoteStepTimeout(t *testing.T) {
	height := uint32(1)
	prevoteStepTimeoutReturned := GetPreVoteStepTimeout(height)
	if preVoteStepTO != prevoteStepTimeoutReturned {
		t.Fatal("Invalid PreVoteStepTimeout returned")
	}
}

func TestGetPreCommitStepTimeout(t *testing.T) {
	height := uint32(1)
	precommitStepTimeoutReturned := GetPreCommitStepTimeout(height)
	if preCommitStepTO != precommitStepTimeoutReturned {
		t.Fatal("Invalid PreCommitStepTimeout returned")
	}
}

func TestGetDeadBlockRoundNextRoundTimeout(t *testing.T) {
	height := uint32(1)
	dbrnrTimeoutReturned := GetDeadBlockRoundNextRoundTimeout(height)
	if dBRNRTO != dbrnrTimeoutReturned {
		t.Fatal("Invalid DeadBlockRoundNextRoundTimeout returned")
	}
}

func TestDownloadTimeout(t *testing.T) {
	height := uint32(1)
	downloadTimeoutReturned := GetDownloadTimeout(height)
	if downloadTO != downloadTimeoutReturned {
		t.Fatal("Invalid DownloadTimeout returned")
	}
}
