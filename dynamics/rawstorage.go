package dynamics

import (
	"encoding/json"
	"math/big"
	"time"
)

// RawStorage is the struct which actually stores everything
type RawStorage struct {
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

// Marshal performs json.Marshal on the RawStorage struct.
func (rs *RawStorage) Marshal() ([]byte, error) {
	return json.Marshal(rs)
}

// Unmarshal performs json.Unmarshal on the RawStorage struct.
func (rs *RawStorage) Unmarshal(v []byte) error {
	if rs == nil {
		return ErrRawStorageNilPointer
	}
	return json.Unmarshal(v, rs)
}

// Copy makes a complete copy of RawStorage struct.
func (rs *RawStorage) Copy() (*RawStorage, error) {
	rsBytes, err := rs.Marshal()
	if err != nil {
		return nil, err
	}
	c := &RawStorage{}
	err = c.Unmarshal(rsBytes)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Overwrite replaces the current StorageGetInterface contents with the copy
// taken as an argument
func (rs *RawStorage) Overwrite(c *RawStorage) error {
	cBytes, err := c.Marshal()
	if err != nil {
		return err
	}
	err = rs.Unmarshal(cBytes)
	if err != nil {
		return err
	}
	return nil
}

// standardParameters initializes RawStorage with the standard (original)
// parameters for the system.
func (rs *RawStorage) standardParameters() {
	rs.MaxBytes = maxBytes
	rs.MaxProposalSize = maxProposalSize
	rs.ProposalStepTimeout = proposalStepTO
	rs.PreVoteStepTimeout = preVoteStepTO
	rs.PreCommitStepTimout = preCommitStepTO
	rs.DeadBlockRoundNextRoundTimeout = dBRNRTO
	rs.DownloadTimeout = downloadTO
	rs.SrvrMsgTimeout = srvrMsgTimeout
	rs.MsgTimeout = msgTimeout
}

// GetMaxBytes returns the maximum allowed bytes
func (rs *RawStorage) GetMaxBytes() uint32 {
	return rs.MaxBytes
}

// SetMaxBytes sets the maximum allowed bytes
func (rs *RawStorage) SetMaxBytes(value uint32) {
	rs.MaxBytes = value
}

// GetMaxProposalSize returns the maximum size of bytes allowed in a proposal
func (rs *RawStorage) GetMaxProposalSize() uint32 {
	return rs.MaxProposalSize
}

// SetMaxProposalSize sets the maximum size of bytes allowed in a proposal
func (rs *RawStorage) SetMaxProposalSize(value uint32) {
	rs.MaxProposalSize = value
}

// GetSrvrMsgTimeout returns the time before timeout of server message
func (rs *RawStorage) GetSrvrMsgTimeout() time.Duration {
	return rs.SrvrMsgTimeout
}

// SetSrvrMsgTimeout sets the time before timeout of server message
func (rs *RawStorage) SetSrvrMsgTimeout(value time.Duration) {
	rs.SrvrMsgTimeout = value
}

// GetMsgTimeout returns the timeout to receive a message
func (rs *RawStorage) GetMsgTimeout() time.Duration {
	return rs.MsgTimeout
}

// SetMsgTimeout sets the timeout to receive a message
// TODO: enforce correct relationship!
func (rs *RawStorage) SetMsgTimeout(value time.Duration) {
	rs.MsgTimeout = value
}

// GetProposalStepTimeout returns the proposal step timeout
func (rs *RawStorage) GetProposalStepTimeout() time.Duration {
	return rs.ProposalStepTimeout
}

// SetProposalStepTimeout sets the proposal step timeout
// TODO: enforce correct relationship!
func (rs *RawStorage) SetProposalStepTimeout(value time.Duration) {
	rs.ProposalStepTimeout = value
}

// GetPreVoteStepTimeout returns the prevote step timeout
func (rs *RawStorage) GetPreVoteStepTimeout() time.Duration {
	return rs.PreVoteStepTimeout
}

// SetPreVoteStepTimeout sets the prevote step timeout
// TODO: enforce correct relationship!
func (rs *RawStorage) SetPreVoteStepTimeout(value time.Duration) {
	rs.PreVoteStepTimeout = value
}

// GetPreCommitStepTimeout returns the precommit step timeout
func (rs *RawStorage) GetPreCommitStepTimeout() time.Duration {
	return rs.PreCommitStepTimout
}

// SetPreCommitStepTimeout sets the precommit step timeout
// TODO: enforce correct relationship!
func (rs *RawStorage) SetPreCommitStepTimeout(value time.Duration) {
	rs.PreCommitStepTimout = value
}

// GetDeadBlockRoundNextRoundTimeout returns the timeout required before
// moving into the DeadBlockRound
func (rs *RawStorage) GetDeadBlockRoundNextRoundTimeout() time.Duration {
	return rs.DeadBlockRoundNextRoundTimeout
}

// SetDeadBlockRoundNextRoundTimeout sets the timeout required before
// moving into the DeadBlockRound
// TODO: enforce correct relationship!
func (rs *RawStorage) SetDeadBlockRoundNextRoundTimeout(value time.Duration) {
	rs.DeadBlockRoundNextRoundTimeout = value
}

// GetDownloadTimeout returns the timeout for downloads
func (rs *RawStorage) GetDownloadTimeout() time.Duration {
	return rs.DownloadTimeout
}

// SetDownloadTimeout sets the timeout for downloads
// TODO: enforce correct relationship!
func (rs *RawStorage) SetDownloadTimeout(value time.Duration) {
	rs.DownloadTimeout = value
}
