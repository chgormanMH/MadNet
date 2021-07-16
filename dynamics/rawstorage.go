package dynamics

import (
	"encoding/json"
	"math/big"
	"strconv"
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

/*
struct {
	Type
	Value
	Name
}

type struct interface {
	GetName() string
	GetType() enum
	GetValue() interface{}/string (parse using type for enum)
}

*/

// Marshal performs json.Marshal on the RawStorage struct.
func (rs *RawStorage) Marshal() ([]byte, error) {
	if rs == nil {
		return nil, ErrRawStorageNilPointer
	}
	return json.Marshal(rs)
}

// Unmarshal performs json.Unmarshal on the RawStorage struct.
func (rs *RawStorage) Unmarshal(v []byte) error {
	if rs == nil {
		return ErrRawStorageNilPointer
	}
	if len(v) == 0 {
		return ErrUnmarshalEmpty
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

// IsValid returns true if we can successfully make a copy
func (rs *RawStorage) IsValid() bool {
	_, err := rs.Copy()
	if err != nil {
		return false
	}
	return true
}

// ParseInt64/Uint32/BigInt; fix this
// *DRY*
//
// Do enum for field (look at iota/opcode )

// UpdateValue updates the field with the appropriate value.
func (rs *RawStorage) UpdateValue(field, value string) error {
	switch field {
	case "maxBytes":
		// uint32
		v64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		v := uint32(v64)
		rs.SetMaxBytes(v)
	case "proposalStepTimeout":
		// time.Duration
		v64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if v64 < 0 {
			return ErrInvalidUpdateValue
		}
		v := time.Duration(v64)
		rs.SetProposalStepTimeout(v)
	case "preVoteStepTimeout":
		// time.Duration
		v64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if v64 < 0 {
			return ErrInvalidUpdateValue
		}
		v := time.Duration(v64)
		rs.SetPreVoteStepTimeout(v)
	case "preCommitStepTimeout":
		// time.Duration
		v64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if v64 < 0 {
			return ErrInvalidUpdateValue
		}
		v := time.Duration(v64)
		rs.SetPreCommitStepTimeout(v)
	case "msgTimeout":
		// time.Duration
		v64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if v64 < 0 {
			return ErrInvalidUpdateValue
		}
		v := time.Duration(v64)
		rs.SetMsgTimeout(v)
	case "minTxBurnedFee":
		// *big.Int
		v, valid := new(big.Int).SetString(value, 10)
		if !valid {
			return ErrInvalidUpdateValue
		}
		if v.Sign() < 0 {
			return ErrInvalidUpdateValue
		}
		rs.SetMinTxBurnedFee(v)
	case "txValidVersion":
		// uint32
		v64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		v := uint32(v64)
		rs.SetTxValidVersion(v)
	case "minValueStoreBurnedFee":
		// *big.Int
		v, valid := new(big.Int).SetString(value, 10)
		if !valid {
			return ErrInvalidUpdateValue
		}
		if v.Sign() < 0 {
			return ErrInvalidUpdateValue
		}
		rs.SetMinValueStoreBurnedFee(v)
	case "valueStoreTxValidVersion":
		// uint32
		v64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		v := uint32(v64)
		rs.SetValueStoreTxValidVersion(v)
	case "minAtomicSwapBurnedFee":
		// *big.Int
		v, valid := new(big.Int).SetString(value, 10)
		if !valid {
			return ErrInvalidUpdateValue
		}
		if v.Sign() < 0 {
			return ErrInvalidUpdateValue
		}
		rs.SetMinAtomicSwapBurnedFee(v)
	case "atomicSwapValidStopEpoch":
		// uint32
		v64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		v := uint32(v64)
		rs.SetAtomicSwapValidStopEpoch(v)
	case "dataStoreTxValidVersion":
		// uint32
		v64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		v := uint32(v64)
		rs.SetDataStoreTxValidVersion(v)
	default:
		return ErrInvalidUpdateValue
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
	rs.MaxProposalSize = value
}

// GetMaxProposalSize returns the maximum size of bytes allowed in a proposal
func (rs *RawStorage) GetMaxProposalSize() uint32 {
	return rs.MaxProposalSize
}

// GetSrvrMsgTimeout returns the time before timeout of server message
func (rs *RawStorage) GetSrvrMsgTimeout() time.Duration {
	return rs.SrvrMsgTimeout
}

// GetMsgTimeout returns the timeout to receive a message
func (rs *RawStorage) GetMsgTimeout() time.Duration {
	return rs.MsgTimeout
}

// SetMsgTimeout sets the timeout to receive a message
func (rs *RawStorage) SetMsgTimeout(value time.Duration) {
	rs.MsgTimeout = value
	rs.SrvrMsgTimeout = (3 * value) / 4
}

// GetProposalStepTimeout returns the proposal step timeout
func (rs *RawStorage) GetProposalStepTimeout() time.Duration {
	return rs.ProposalStepTimeout
}

// SetProposalStepTimeout sets the proposal step timeout
func (rs *RawStorage) SetProposalStepTimeout(value time.Duration) {
	rs.ProposalStepTimeout = value
	sum := rs.ProposalStepTimeout + rs.PreVoteStepTimeout + rs.PreCommitStepTimout
	rs.DownloadTimeout = sum
	rs.DeadBlockRoundNextRoundTimeout = (5 * sum) / 2
}

// GetPreVoteStepTimeout returns the prevote step timeout
func (rs *RawStorage) GetPreVoteStepTimeout() time.Duration {
	return rs.PreVoteStepTimeout
}

// SetPreVoteStepTimeout sets the prevote step timeout
func (rs *RawStorage) SetPreVoteStepTimeout(value time.Duration) {
	rs.PreVoteStepTimeout = value
	sum := rs.ProposalStepTimeout + rs.PreVoteStepTimeout + rs.PreCommitStepTimout
	rs.DownloadTimeout = sum
	rs.DeadBlockRoundNextRoundTimeout = (5 * sum) / 2
}

// GetPreCommitStepTimeout returns the precommit step timeout
func (rs *RawStorage) GetPreCommitStepTimeout() time.Duration {
	return rs.PreCommitStepTimout
}

// SetPreCommitStepTimeout sets the precommit step timeout
func (rs *RawStorage) SetPreCommitStepTimeout(value time.Duration) {
	rs.PreCommitStepTimout = value
	sum := rs.ProposalStepTimeout + rs.PreVoteStepTimeout + rs.PreCommitStepTimout
	rs.DownloadTimeout = sum
	rs.DeadBlockRoundNextRoundTimeout = (5 * sum) / 2
}

// GetDeadBlockRoundNextRoundTimeout returns the timeout required before
// moving into the DeadBlockRound
func (rs *RawStorage) GetDeadBlockRoundNextRoundTimeout() time.Duration {
	return rs.DeadBlockRoundNextRoundTimeout
}

// GetDownloadTimeout returns the timeout for downloads
func (rs *RawStorage) GetDownloadTimeout() time.Duration {
	return rs.DownloadTimeout
}

// GetMinTxBurnedFee returns the minimun tx burned fee
func (rs *RawStorage) GetMinTxBurnedFee() *big.Int {
	if rs.MinTxBurnedFee == nil {
		rs.MinTxBurnedFee = new(big.Int)
	}
	return rs.MinTxBurnedFee
}

// SetMinTxBurnedFee sets the minimun tx burned fee
func (rs *RawStorage) SetMinTxBurnedFee(value *big.Int) error {
	if value == nil {
		return ErrInvalidValue
	}
	if rs.MinTxBurnedFee == nil {
		rs.MinTxBurnedFee = new(big.Int)
	}
	if value.Sign() < 0 {
		return ErrInvalidValue
	}
	rs.MinTxBurnedFee.Set(value)
	return nil
}

// GetTxValidVersion returns the valid version of tx
func (rs *RawStorage) GetTxValidVersion() uint32 {
	return rs.TxValidVersion
}

// SetTxValidVersion sets the minimun tx burned fee
func (rs *RawStorage) SetTxValidVersion(value uint32) {
	rs.TxValidVersion = value
}

// GetMinValueStoreBurnedFee returns the minimun ValueStore burned fee
func (rs *RawStorage) GetMinValueStoreBurnedFee() *big.Int {
	if rs.MinValueStoreBurnedFee == nil {
		rs.MinValueStoreBurnedFee = new(big.Int)
	}
	return rs.MinValueStoreBurnedFee
}

// SetMinValueStoreBurnedFee sets the minimun ValueStore burned fee
func (rs *RawStorage) SetMinValueStoreBurnedFee(value *big.Int) error {
	if value == nil {
		return ErrInvalidValue
	}
	if rs.MinValueStoreBurnedFee == nil {
		rs.MinValueStoreBurnedFee = new(big.Int)
	}
	if value.Sign() < 0 {
		return ErrInvalidValue
	}
	rs.MinValueStoreBurnedFee.Set(value)
	return nil
}

// GetValueStoreTxValidVersion returns the valid version of ValueStore
func (rs *RawStorage) GetValueStoreTxValidVersion() uint32 {
	return rs.ValueStoreTxValidVersion
}

// SetValueStoreTxValidVersion sets the valid version of ValueStore
func (rs *RawStorage) SetValueStoreTxValidVersion(value uint32) {
	rs.ValueStoreTxValidVersion = value
}

// GetMinAtomicSwapBurnedFee returns the minimun AtomicSwap burned fee
func (rs *RawStorage) GetMinAtomicSwapBurnedFee() *big.Int {
	if rs.MinAtomicSwapBurnedFee == nil {
		rs.MinAtomicSwapBurnedFee = new(big.Int)
	}
	return rs.MinAtomicSwapBurnedFee
}

// SetMinAtomicSwapBurnedFee sets the minimun AtomicSwap burned fee
func (rs *RawStorage) SetMinAtomicSwapBurnedFee(value *big.Int) error {
	if value == nil {
		return ErrInvalidValue
	}
	if rs.MinAtomicSwapBurnedFee == nil {
		rs.MinAtomicSwapBurnedFee = new(big.Int)
	}
	if value.Sign() < 0 {
		return ErrInvalidValue
	}
	rs.MinAtomicSwapBurnedFee.Set(value)
	return nil
}

// GetAtomicSwapValidStopEpoch returns the valid version of AtomicSwap
func (rs *RawStorage) GetAtomicSwapValidStopEpoch() uint32 {
	return rs.AtomicSwapValidStopEpoch
}

// SetAtomicSwapValidStopEpoch sets the valid version of AtomicSwap
func (rs *RawStorage) SetAtomicSwapValidStopEpoch(value uint32) {
	rs.AtomicSwapValidStopEpoch = value
}

// GetDataStoreTxValidVersion returns the valid version of DataStore
func (rs *RawStorage) GetDataStoreTxValidVersion() uint32 {
	return rs.DataStoreTxValidVersion
}

// SetDataStoreTxValidVersion sets the valid version of AtomicSwap
func (rs *RawStorage) SetDataStoreTxValidVersion(value uint32) {
	rs.DataStoreTxValidVersion = value
}
