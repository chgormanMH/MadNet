package dynamics

import (
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/utils"
)

// LinkedList is a doubly linked list which will store nodes corresponding
// to changes to dynamic parameters.
// We store the largest epoch which has been updated.
type LinkedList struct {
	currentEpoch     uint32
	epochLastUpdated uint32
}

func makeLinkedListKey() *NodeKey {
	nk := &NodeKey{
		prefix: dbprefix.PrefixStorageNodeKey(),
		epoch:  0,
	}
	return nk
}

// GetEpochLastUpdated returns highest epoch with changes
func (ll *LinkedList) GetEpochLastUpdated() uint32 {
	return ll.epochLastUpdated
}

// SetEpochLastUpdated returns highest epoch with changes
func (ll *LinkedList) SetEpochLastUpdated(epoch uint32) error {
	if epoch == 0 {
		return ErrZeroEpoch
	}
	ll.epochLastUpdated = epoch
	return nil
}

// GetCurrentEpoch returns current epoch
func (ll *LinkedList) GetCurrentEpoch() uint32 {
	return ll.currentEpoch
}

// SetCurrentEpoch sets the current epoch
func (ll *LinkedList) SetCurrentEpoch(epoch uint32) error {
	if epoch == 0 {
		return ErrZeroEpoch
	}
	ll.currentEpoch = epoch
	return nil
}

// Marshal marshals LinkedList
func (ll *LinkedList) Marshal() ([]byte, error) {
	if !ll.IsValid() {
		return nil, ErrInvalid
	}
	eluBytes := utils.MarshalUint32(ll.epochLastUpdated)
	ceBytes := utils.MarshalUint32(ll.currentEpoch)
	v := append(eluBytes, ceBytes...)
	return v, nil
}

// Unmarshal unmarshals LinkedList
func (ll *LinkedList) Unmarshal(v []byte) error {
	if len(v) != 8 {
		return ErrInvalidNode
	}
	elu, _ := utils.UnmarshalUint32(v[0:4])
	ce, _ := utils.UnmarshalUint32(v[4:8])
	ll.epochLastUpdated = elu
	ll.currentEpoch = ce
	if !ll.IsValid() {
		return ErrInvalid
	}
	return nil
}

// IsValid returns true if LinkedList is valid
func (ll *LinkedList) IsValid() bool {
	if ll.epochLastUpdated == 0 || ll.currentEpoch == 0 {
		return false
	}
	return true
}

// CreateLinkedList creates the first node in a LinkedList.
// These values can then be stored in the database.
func CreateLinkedList(epoch uint32, rs *RawStorage) (*Node, *LinkedList, error) {
	if epoch == 0 {
		return nil, nil, ErrZeroEpoch
	}
	rsCopy, err := rs.Copy()
	if err != nil {
		return nil, nil, err
	}
	node := &Node{
		thisEpoch:  epoch,
		prevEpoch:  epoch,
		nextEpoch:  epoch,
		rawStorage: rsCopy,
	}
	linkedList := &LinkedList{
		epochLastUpdated: epoch,
		currentEpoch:     epoch,
	}
	return node, linkedList, nil
}
