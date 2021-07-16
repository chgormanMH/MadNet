package dynamics

import (
	"bytes"

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
	if !node.IsValid() {
		return nil, nil, ErrInvalidNode
	}
	linkedList := &LinkedList{
		epochLastUpdated: epoch,
		currentEpoch:     epoch,
	}
	if !linkedList.IsValid() {
		return nil, nil, ErrInvalid
	}
	return node, linkedList, nil
}

// NodeKey stores the necessary information to load a Node
type NodeKey struct {
	prefix []byte
	epoch  uint32
}

func makeNodeKey(epoch uint32) (*NodeKey, error) {
	if epoch == 0 {
		return nil, ErrZeroEpoch
	}
	nk := &NodeKey{
		prefix: dbprefix.PrefixStorageNodeKey(),
		epoch:  epoch,
	}
	return nk, nil
}

// Marshal converts NodeKey into the byte slice
func (nk *NodeKey) Marshal() ([]byte, error) {
	if !bytes.Equal(nk.prefix, dbprefix.PrefixStorageNodeKey()) {
		return nil, ErrInvalidNodeKey
	}
	epochBytes := utils.MarshalUint32(nk.epoch)
	key := []byte{}
	key = append(key, nk.prefix...)
	key = append(key, epochBytes...)
	return key, nil
}

// Node contains the necessary information about RawStorage
type Node struct {
	thisEpoch  uint32
	prevEpoch  uint32
	nextEpoch  uint32
	rawStorage *RawStorage
}

// Marshal marshals a Node
func (n *Node) Marshal() ([]byte, error) {
	rsBytes, err := n.rawStorage.Marshal()
	if err != nil {
		return nil, err
	}
	teBytes := utils.MarshalUint32(n.thisEpoch)
	peBytes := utils.MarshalUint32(n.prevEpoch)
	neBytes := utils.MarshalUint32(n.nextEpoch)
	v := []byte{}
	v = append(v, teBytes...)
	v = append(v, peBytes...)
	v = append(v, neBytes...)
	v = append(v, rsBytes...)
	return v, nil
}

// Unmarshal unmarshals a Node
func (n *Node) Unmarshal(v []byte) error {
	if len(v) < 12 {
		return ErrInvalidNode
	}
	thisEpoch, _ := utils.UnmarshalUint32(v[0:4])
	prevEpoch, _ := utils.UnmarshalUint32(v[4:8])
	nextEpoch, _ := utils.UnmarshalUint32(v[8:12])
	rs := &RawStorage{}
	err := rs.Unmarshal(v[12:])
	if err != nil {
		return err
	}
	n.thisEpoch = thisEpoch
	n.prevEpoch = prevEpoch
	n.nextEpoch = nextEpoch
	n.rawStorage, err = rs.Copy()
	if err != nil {
		return err
	}
	return nil
}

// IsValid returns true if Node is valid
func (n *Node) IsValid() bool {
	if n == nil {
		return false
	}
	if n.thisEpoch == 0 || n.prevEpoch == 0 || n.nextEpoch == 0 {
		// node has not set values; invalid
		return false
	}
	if n.prevEpoch > n.thisEpoch || n.thisEpoch > n.nextEpoch {
		// node has not been correctly defined; invalid
		return false
	}
	_, err := n.rawStorage.Copy()
	if err != nil {
		// invalid RawStorage; invalid
		return false
	}
	return true
}

// IsPreValid returns true if Node is ready to be added to database
func (n *Node) IsPreValid() bool {
	if n == nil {
		return false
	}
	if n.thisEpoch == 0 || n.prevEpoch != 0 || n.nextEpoch != 0 {
		// Only thisEpoch should be set; invalid
		return false
	}
	_, err := n.rawStorage.Copy()
	if err != nil {
		// invalid RawStorage; invalid
		return false
	}
	return true
}

func (n *Node) Copy() (*Node, error) {
	if !n.IsValid() {
		return nil, ErrInvalidNode
	}
	nodeBytes, err := n.Marshal()
	if err != nil {
		return nil, err
	}
	nodeCopy := &Node{}
	err = nodeCopy.Unmarshal(nodeBytes)
	if err != nil {
		return nil, err
	}
	return nodeCopy, nil
}

/*
Need to have doubly linked list

First Node will have prevEpoch point to self.
Final Node will have nextEpoch point to self		.
*/

// SetEpochs sets n.prevEpoch and n.nextEpoch.
func (n *Node) SetEpochs(prevNode *Node, nextNode *Node) error {
	if !n.IsPreValid() {
		return ErrInvalid
	}
	if prevNode.IsValid() && nextNode == nil && prevNode.thisEpoch < n.thisEpoch {
		if prevNode.IsHead() {
			// n is the new Head
			// Update prevNode.nextEpoch
			prevNode.nextEpoch = n.thisEpoch
			// Update epochs for n;
			// must point backward to prevNode and forward to self
			n.prevEpoch = prevNode.thisEpoch
			n.nextEpoch = n.thisEpoch
			return nil
		}
		// prevNode is not head, so prevNode.thisEpoch < prevNode.nextEpoch
		// Update n
		n.nextEpoch = prevNode.nextEpoch
		n.prevEpoch = prevNode.thisEpoch
		// Update prevNode
		prevNode.nextEpoch = n.thisEpoch
		return nil
	}
	if prevNode == nil && nextNode.IsValid() && n.thisEpoch < nextNode.thisEpoch && nextNode.IsTail() {
		// n is the new Tail
		// Update nextNode.prevEpoch
		nextNode.prevEpoch = n.thisEpoch
		// Update n;
		// must point backward to self and forward to nextNode
		n.prevEpoch = n.thisEpoch
		n.nextEpoch = nextNode.thisEpoch
		return nil
	}
	return ErrInvalid
}

func (n *Node) AddNext(epoch uint32, rs *RawStorage) error {
	return nil
}

func (n *Node) SplitNodes(epoch uint32, rs *RawStorage) error {
	return nil
}

// IsHead returns true if Node is end of linked list
func (n *Node) IsHead() bool {
	if n.thisEpoch == n.nextEpoch {
		return true
	}
	return false
}

// IsTail returns true if Node is beginning of linked list
func (n *Node) IsTail() bool {
	if n.thisEpoch == n.prevEpoch {
		return true
	}
	return false
}
