package objs

import (
	"time"

	"github.com/MadBase/MadNet/errorz"

	mdefs "github.com/MadBase/MadNet/consensus/objs/capn"
	"github.com/MadBase/MadNet/consensus/objs/ovstate"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/dynamics"
	gUtils "github.com/MadBase/MadNet/utils"
	capnp "zombiezen.com/go/capnproto2"
)

// OwnValidatingState ...
type OwnValidatingState struct {
	VAddr                []byte
	GroupKey             []byte
	RoundStarted         int64
	PreVoteStepStarted   int64
	PreCommitStepStarted int64
	ValidValue           *Proposal
	LockedValue          *Proposal
}

// UnmarshalBinary takes a byte slice and returns the corresponding
// OwnValidatingState object
func (b *OwnValidatingState) UnmarshalBinary(data []byte) error {
	bh, err := ovstate.Unmarshal(data)
	if err != nil {
		return err
	}
	return b.UnmarshalCapn(bh)
}

// UnmarshalCapn unmarshals the capnproto definition of the object
func (b *OwnValidatingState) UnmarshalCapn(bh mdefs.OwnValidatingState) error {
	err := ovstate.Validate(bh)
	if err != nil {
		return err
	}
	b.VAddr = gUtils.CopySlice(bh.VAddr())
	b.GroupKey = gUtils.CopySlice(bh.GroupKey())
	b.RoundStarted = bh.RoundStarted()
	b.PreVoteStepStarted = bh.PreVoteStepStarted()
	b.PreCommitStepStarted = bh.PreCommitStepStarted()
	if bh.HasLockedValue() {
		b.LockedValue = &Proposal{}
		err := b.LockedValue.UnmarshalCapn(bh.LockedValue())
		if err != nil {
			return err
		}
	}
	if bh.HasValidValue() {
		b.ValidValue = &Proposal{}
		err := b.ValidValue.UnmarshalCapn(bh.ValidValue())
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalBinary takes the OwnValidatingState object and returns the canonical
// byte slice
func (b *OwnValidatingState) MarshalBinary() ([]byte, error) {
	if b == nil {
		return nil, errorz.ErrInvalid{}.New("not initialized")
	}
	bh, err := b.MarshalCapn(nil)
	if err != nil {
		return nil, err
	}
	return ovstate.Marshal(bh)
}

// MarshalCapn marshals the object into its capnproto definition
func (b *OwnValidatingState) MarshalCapn(seg *capnp.Segment) (mdefs.OwnValidatingState, error) {
	if b == nil {
		return mdefs.OwnValidatingState{}, errorz.ErrInvalid{}.New("not initialized")
	}
	var bh mdefs.OwnValidatingState
	if seg == nil {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			return bh, err
		}
		tmp, err := mdefs.NewRootOwnValidatingState(seg)
		if err != nil {
			return bh, err
		}
		bh = tmp
	} else {
		tmp, err := mdefs.NewOwnValidatingState(seg)
		if err != nil {
			return bh, err
		}
		bh = tmp
	}
	err := bh.SetVAddr(b.VAddr)
	if err != nil {
		return mdefs.OwnValidatingState{}, err
	}
	err = bh.SetGroupKey(b.GroupKey)
	if err != nil {
		return mdefs.OwnValidatingState{}, err
	}
	bh.SetRoundStarted(b.RoundStarted)
	bh.SetPreVoteStepStarted(b.PreVoteStepStarted)
	bh.SetPreCommitStepStarted(b.PreCommitStepStarted)
	if b.LockedValue != nil {
		LockedValue, err := b.LockedValue.MarshalCapn(seg)
		if err != nil {
			return bh, err
		}
		err = bh.SetLockedValue(LockedValue)
		if err != nil {
			return mdefs.OwnValidatingState{}, err
		}
	}
	if b.ValidValue != nil {
		ValidValue, err := b.ValidValue.MarshalCapn(seg)
		if err != nil {
			return bh, err
		}
		err = bh.SetValidValue(ValidValue)
		if err != nil {
			return mdefs.OwnValidatingState{}, err
		}
	}
	return bh, nil
}

func (b *OwnValidatingState) PTOExpired(height uint32) bool {
	rs := b.RoundStarted
	proposalStepTO := dynamics.GetProposalStepTimeout(height)
	return rs+int64(proposalStepTO)/constants.OneBillion < time.Now().Unix()
}

func (b *OwnValidatingState) PVTOExpired(height uint32) bool {
	rs := b.PreVoteStepStarted
	prevoteStepTO := dynamics.GetPreVoteStepTimeout(height)
	return rs+int64(prevoteStepTO)/constants.OneBillion < time.Now().Unix()
}

func (b *OwnValidatingState) PCTOExpired(height uint32) bool {
	rs := b.PreCommitStepStarted
	precommitStepTO := dynamics.GetPreCommitStepTimeout(height)
	return rs+int64(precommitStepTO)/constants.OneBillion < time.Now().Unix()
}

func (b *OwnValidatingState) DBRNRExpired(height uint32) bool {
	rs := b.PreCommitStepStarted
	dbrnrTO := dynamics.GetDeadBlockRoundNextRoundTimeout(height)
	return rs+int64(dbrnrTO)/constants.OneBillion < time.Now().Unix()
}

func (b *OwnValidatingState) SetRoundStarted() {
	now := time.Now()
	b.RoundStarted = now.Unix()
	b.PreVoteStepStarted = 0
	b.PreCommitStepStarted = 0
}

func (b *OwnValidatingState) SetPreVoteStepStarted() {
	now := time.Now()
	b.PreVoteStepStarted = now.Unix()
	b.PreCommitStepStarted = 0
}

func (b *OwnValidatingState) SetPreCommitStepStarted() {
	now := time.Now()
	b.PreCommitStepStarted = now.Unix()
}
