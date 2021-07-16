package dynamics

import (
	"github.com/MadBase/MadNet/utils"
)

// Node contains the necessary information about RawStorage
type Node struct {
	thisEpoch  uint32
	prevEpoch  uint32
	nextEpoch  uint32
	rawStorage *RawStorage
}

// Marshal marshals a Node
func (n *Node) Marshal() ([]byte, error) {
	if !n.IsValid() {
		return nil, ErrInvalidNode
	}
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
		return ErrInvalid
	}
	thisEpoch, _ := utils.UnmarshalUint32(v[0:4])
	prevEpoch, _ := utils.UnmarshalUint32(v[4:8])
	nextEpoch, _ := utils.UnmarshalUint32(v[8:12])
	n.thisEpoch = thisEpoch
	n.prevEpoch = prevEpoch
	n.nextEpoch = nextEpoch
	n.rawStorage = &RawStorage{}
	err := n.rawStorage.Unmarshal(v[12:])
	if err != nil {
		return err
	}
	if !n.IsValid() {
		return ErrInvalidNode
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
	if prevNode.IsValid() && nextNode.IsValid() && prevNode.thisEpoch < n.thisEpoch && n.thisEpoch < nextNode.thisEpoch {
		// In this setting, we want to add a new node in between prevNode and nextNode
		//
		// Update prevNode;
		// must point forward to n
		prevNode.nextEpoch = n.thisEpoch
		// Update epochs for n;
		// must point backward to prevNode and forward to nextNode
		n.prevEpoch = prevNode.thisEpoch
		n.nextEpoch = nextNode.thisEpoch
		// Update  nextNode;
		// must point backward to n
		nextNode.prevEpoch = n.thisEpoch
		return nil
	}
	if prevNode.IsValid() && nextNode == nil && prevNode.thisEpoch < n.thisEpoch && prevNode.IsHead() {
		// n is the new Head
		// Update prevNode.nextEpoch
		prevNode.nextEpoch = n.thisEpoch
		// Update epochs for n;
		// must point backward to prevNode and forward to self
		n.prevEpoch = prevNode.thisEpoch
		n.nextEpoch = n.thisEpoch
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

// IsHead returns true if Node is end of linked list
func (n *Node) IsHead() bool {
	if !n.IsValid() {
		return false
	}
	if n.thisEpoch == n.nextEpoch {
		return true
	}
	return false
}

// IsTail returns true if Node is beginning of linked list
func (n *Node) IsTail() bool {
	if !n.IsValid() {
		return false
	}
	if n.thisEpoch == n.prevEpoch {
		return true
	}
	return false
}
