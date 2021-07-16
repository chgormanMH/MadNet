package dynamics

import (
	"bytes"
	"errors"
	"testing"

	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/constants/dbprefix"
)

func TestNodeKeyMarshal(t *testing.T) {
	nk := &NodeKey{}
	_, err := nk.Marshal()
	if err == nil {
		t.Fatal("Should have raised error")
	}
	epoch := uint32(1)
	nk, err = makeNodeKey(epoch)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(nk.prefix, dbprefix.PrefixStorageNodeKey()) {
		t.Fatal("prefixes do not match")
	}
	if nk.epoch != 1 {
		t.Fatal("epochs do not match")
	}
}

func TestNodeLinkedListMakeKeys(t *testing.T) {
	epoch := uint32(0)
	_, err := makeNodeKey(epoch)
	if !errors.Is(err, ErrZeroEpoch) {
		t.Fatal("Should have returned error for zero epoch")
	}

	epoch = 1
	nk, err := makeNodeKey(epoch)
	if err != nil {
		t.Fatal(err)
	}
	if nk.epoch != epoch {
		t.Fatal("epochs do not match")
	}
	if !bytes.Equal(nk.prefix, dbprefix.PrefixStorageNodeKey()) {
		t.Fatal("prefixes do not match (1)")
	}

	llk := makeLinkedListKey()
	if llk.epoch != 0 {
		t.Fatal("epoch should be 0")
	}
	if !bytes.Equal(nk.prefix, dbprefix.PrefixStorageNodeKey()) {
		t.Fatal("prefixes do not match (2)")
	}
}

func TestLinkedListMarshal(t *testing.T) {
	ll := &LinkedList{}
	if ll.IsValid() {
		t.Fatal("Should not have valid LinkedList")
	}
	_, err := ll.Marshal()
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	invalidBytes := []byte{0, 1, 2, 3, 4}
	err = ll.Unmarshal(invalidBytes)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	invalidBytes2 := make([]byte, 8)
	err = ll.Unmarshal(invalidBytes2)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	v := []byte{255, 255, 255, 255, 0, 0, 0, 1}
	err = ll.Unmarshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if ll.epochLastUpdated != constants.MaxUint32 {
		t.Fatal("Invalid LinkedList (1)")
	}
	if ll.currentEpoch != 1 {
		t.Fatal("Invalid LinkedList (2)")
	}
}

func TestLinkedListGetSet(t *testing.T) {
	ll := &LinkedList{}
	err := ll.SetEpochLastUpdated(0)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}
	err = ll.SetCurrentEpoch(0)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	elu := uint32(123456)
	err = ll.SetEpochLastUpdated(elu)
	if err != nil {
		t.Fatal(err)
	}
	retElu := ll.GetEpochLastUpdated()
	if retElu != elu {
		t.Fatal("Invalid EpochLastUpdated")
	}

	ce := uint32(25519)
	err = ll.SetCurrentEpoch(ce)
	if err != nil {
		t.Fatal(err)
	}
	retCe := ll.GetCurrentEpoch()
	if retCe != ce {
		t.Fatal("Invalid CurrentEpoch")
	}

	if !ll.IsValid() {
		t.Fatal("LinkedList should be valid")
	}
}

func TestCreateLinkedList(t *testing.T) {
	epoch := uint32(0)
	_, _, err := CreateLinkedList(epoch, nil)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	epoch = 1
	_, _, err = CreateLinkedList(epoch, nil)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	rs := &RawStorage{}
	rs.standardParameters()
	node, linkedlist, err := CreateLinkedList(epoch, rs)
	if err != nil {
		t.Fatal(err)
	}
	if node.thisEpoch != epoch {
		t.Fatal("invalid thisEpoch")
	}
	if node.prevEpoch != epoch {
		t.Fatal("invalid prevEpoch")
	}
	if node.nextEpoch != epoch {
		t.Fatal("invalid nextEpoch")
	}
	if linkedlist.epochLastUpdated != epoch {
		t.Fatal("invalid epochLastUpdated")
	}

	err = linkedlist.SetCurrentEpoch(0)
	if err == nil {
		t.Fatal("Should have raised error")
	}

	epoch = 2
	err = linkedlist.SetCurrentEpoch(epoch)
	if err != nil {
		t.Fatal(err)
	}
	if linkedlist.currentEpoch != epoch {
		t.Fatal("invalid currentEpoch")
	}
}

func TestNodeMarshal(t *testing.T) {
	node := &Node{}
	_, err := node.Marshal()
	if err == nil {
		t.Fatal("Should have raied error (1)")
	}

	epoch := uint32(1)
	rs := &RawStorage{}
	rs.standardParameters()
	node, _, err = CreateLinkedList(epoch, rs)
	if err != nil {
		t.Fatal(err)
	}

	nodeBytes, err := node.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	node2 := &Node{}
	err = node2.Unmarshal(nodeBytes)
	if err != nil {
		t.Fatal(err)
	}
	if node.thisEpoch != node2.thisEpoch {
		t.Fatal("invalid thisEpoch")
	}
	if node.prevEpoch != node2.prevEpoch {
		t.Fatal("invalid prevEpoch")
	}
	if node.nextEpoch != node2.nextEpoch {
		t.Fatal("invalid nextEpoch")
	}
	rsBytes, err := node.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rs2Bytes, err := node2.rawStorage.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rsBytes, rs2Bytes) {
		t.Fatal("invalid RawStroage")
	}

	v := []byte{}
	err = node.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	v = make([]byte, 13)
	err = node.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}
}

type wNode struct {
	node *Node
}

func TestNodeIsValid(t *testing.T) {
	wNode := &wNode{}
	if wNode.node.IsValid() {
		t.Fatal("Node should not be valid (0)")
	}

	node := &Node{}
	if node.IsValid() {
		t.Fatal("Node should not be valid (1)")
	}

	node.prevEpoch = 3
	node.thisEpoch = 2
	node.nextEpoch = 3
	if node.IsValid() {
		t.Fatal("Node should not be valid (2)")
	}

	node.prevEpoch = 1
	node.thisEpoch = 3
	node.nextEpoch = 2
	if node.IsValid() {
		t.Fatal("Node should not be valid (3)")
	}

	node.prevEpoch = 1
	node.thisEpoch = 2
	node.nextEpoch = 3
	if node.IsValid() {
		t.Fatal("Node should not be valid (4)")
	}

	node.rawStorage = &RawStorage{}
	if !node.IsValid() {
		t.Fatal("Node should be valid")
	}
}

func TestNodeIsPreValid(t *testing.T) {
	wNode := &wNode{}
	if wNode.node.IsPreValid() {
		t.Fatal("Node should not be prevalid (0)")
	}

	node := &Node{}
	if node.IsPreValid() {
		t.Fatal("Node should not be prevalid (1)")
	}

	node.prevEpoch = 0
	node.thisEpoch = 0
	node.nextEpoch = 0
	if node.IsPreValid() {
		t.Fatal("Node should not be prevalid (2)")
	}

	node.prevEpoch = 1
	node.thisEpoch = 1
	node.nextEpoch = 0
	if node.IsPreValid() {
		t.Fatal("Node should not be prevalid (3)")
	}

	node.prevEpoch = 0
	node.thisEpoch = 1
	node.nextEpoch = 1
	if node.IsPreValid() {
		t.Fatal("Node should not be prevalid (4)")
	}

	node.prevEpoch = 0
	node.thisEpoch = 1
	node.nextEpoch = 0
	if node.IsPreValid() {
		t.Fatal("Node should not be prevalid (5)")
	}

	node.rawStorage = &RawStorage{}
	if !node.IsPreValid() {
		t.Fatal("Node should be prevalid")
	}
}

// SetNode with prevNode at Head
func TestNodeSetEpochsGood1(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	nodeEpoch := uint32(25519)
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  nodeEpoch,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567890
	first := uint32(1)
	last := uint32(123456789)
	prevEpoch := uint32(257)
	prevNode := &Node{
		prevEpoch:  first,
		thisEpoch:  prevEpoch,
		nextEpoch:  last,
		rawStorage: rsNew,
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be Valid")
	}
	if prevNode.thisEpoch >= node.thisEpoch {
		t.Fatal("Should have prevNode.thisEpoch < node.thisEpoch")
	}
	err = node.SetEpochs(prevNode, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Now need to confirm all epochs are good.
	if prevNode.prevEpoch != first {
		t.Fatal("prevNode.prevEpoch is incorrect")
	}
	if prevNode.thisEpoch != prevEpoch {
		t.Fatal("prevNode.thisEpoch is incorrect")
	}
	if prevNode.nextEpoch != nodeEpoch {
		t.Fatal("prevNode.nextEpoch is incorrect; it does not point to new nodeEpoch")
	}
	if node.prevEpoch != prevEpoch {
		t.Fatal("prevNode.prevEpoch is incorrect; it does not equal prevEpoch")
	}
	if node.thisEpoch != nodeEpoch {
		t.Fatal("node.thisEpoch is incorrect")
	}
	if node.nextEpoch != last {
		t.Fatal("node.nextEpoch is incorrect; it does not point to last")
	}
}

// SetNode with prevNode at not at Head
func TestNodeSetEpochsGood2(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	nodeEpoch := uint32(25519)
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  nodeEpoch,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567890
	first := uint32(1)
	last := uint32(1)
	prevEpoch := uint32(1)
	prevNode := &Node{
		prevEpoch:  first,
		thisEpoch:  prevEpoch,
		nextEpoch:  last,
		rawStorage: rsNew,
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be Valid")
	}
	if prevNode.thisEpoch >= node.thisEpoch {
		t.Fatal("Should have prevNode.thisEpoch < node.thisEpoch")
	}
	err = node.SetEpochs(prevNode, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Now need to confirm all epochs are good.
	if prevNode.prevEpoch != first {
		t.Fatal("prevNode.prevEpoch is incorrect")
	}
	if prevNode.thisEpoch != prevEpoch {
		t.Fatal("prevNode.thisEpoch is incorrect")
	}
	if prevNode.nextEpoch != nodeEpoch {
		t.Fatal("prevNode.nextEpoch is incorrect; it does not point to new nodeEpoch")
	}
	if node.prevEpoch != prevEpoch {
		t.Fatal("prevNode.prevEpoch is incorrect; it does not equal prevEpoch")
	}
	if node.thisEpoch != nodeEpoch {
		t.Fatal("node.thisEpoch is incorrect")
	}
	if node.nextEpoch != nodeEpoch {
		t.Fatal("node.nextEpoch is incorrect; it does not point to last")
	}
}

// SetNode with nextNode at Tail
func TestNodeSetEpochsGood3(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	nodeEpoch := uint32(1)
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  nodeEpoch,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567890
	first := uint32(2)
	last := uint32(3)
	nextEpoch := uint32(2)
	nextNode := &Node{
		prevEpoch:  first,
		thisEpoch:  nextEpoch,
		nextEpoch:  last,
		rawStorage: rsNew,
	}
	if !nextNode.IsValid() {
		t.Fatal("nextNode should be Valid")
	}
	if node.thisEpoch >= nextNode.thisEpoch {
		t.Fatal("Should have node.thisEpoch < nextNode.thisEpoch")
	}
	err = node.SetEpochs(nil, nextNode)
	if err != nil {
		t.Fatal(err)
	}

	// Now need to confirm all epochs are good.
	if node.prevEpoch != nodeEpoch {
		t.Fatal("node.prevEpoch is incorrect; it does not equal nodeEpoch")
	}
	if node.thisEpoch != nodeEpoch {
		t.Fatal("node.thisEpoch is incorrect")
	}
	if node.nextEpoch != nextEpoch {
		t.Fatal("node.nextEpoch is incorrect; it does not point to nextEpoch")
	}
	if nextNode.prevEpoch != nodeEpoch {
		t.Fatal("nextNode.prevEpoch is incorrect")
	}
	if nextNode.thisEpoch != nextEpoch {
		t.Fatal("nextNode.thisEpoch is incorrect")
	}
	if nextNode.nextEpoch != last {
		t.Fatal("nextNode.nextEpoch is incorrect; it does not point to last")
	}
}

// We should raise an error when having node not PreValid
func TestNodeSetEpochsBad1(t *testing.T) {
	node := &Node{}
	err := node.SetEpochs(nil, nil)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for prevNode being invalid
func TestNodeSetEpochsBad2(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	prevNode := &Node{}
	if prevNode.IsValid() {
		t.Fatal("prevNode should not be valid")
	}
	err := node.SetEpochs(prevNode, nil)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for nextNode not nil
func TestNodeSetEpochsBad3(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	prevNode := &Node{
		prevEpoch:  1,
		thisEpoch:  257,
		nextEpoch:  123456789,
		rawStorage: rs,
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be Valid")
	}
	nextNode := &Node{}
	if nextNode == nil {
		t.Fatal("nextNode should not be nil")
	}
	err := node.SetEpochs(prevNode, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for prevNode.thisEpoch >= node.thisEpoch
func TestNodeSetEpochsBad4(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	prevNode := &Node{
		prevEpoch:  1,
		thisEpoch:  25519,
		nextEpoch:  123456789,
		rawStorage: rs,
	}
	if !prevNode.IsValid() {
		t.Fatal("prevNode should be Valid")
	}
	if prevNode.thisEpoch < node.thisEpoch {
		t.Fatal("Should have prevNode.thisEpoch >= node.thisEpoch")
	}
	err := node.SetEpochs(prevNode, nil)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for prevNode not nil
func TestNodeSetEpochsBad5(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	if !node.IsPreValid() {
		t.Fatal("node should be preValid")
	}
	prevNode := &Node{}
	if prevNode.IsValid() {
		t.Fatal("prevNode should not be valid")
	}
	if prevNode == nil {
		t.Fatal("prevNode should not be nil")
	}
	nextNode := &Node{}
	err := node.SetEpochs(prevNode, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for prevNode not nil
func TestNodeSetEpochsBad6(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	nextNode := &Node{}
	if nextNode.IsValid() {
		t.Fatal("nextNode should not be valid")
	}
	err := node.SetEpochs(nil, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for prevNode not nil
func TestNodeSetEpochsBad7(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567890
	nextNode := &Node{
		prevEpoch:  1,
		thisEpoch:  257,
		nextEpoch:  123456789,
		rawStorage: rsNew,
	}
	if !nextNode.IsValid() {
		t.Fatal("nextNode should not be valid")
	}
	if node.thisEpoch < nextNode.thisEpoch {
		t.Fatal("We should not have node.thisEpoch >= nextNode.thisEpoch to raise error")
	}
	err = node.SetEpochs(nil, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// We should raise error for prevNode not nil
func TestNodeSetEpochsBad8(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()
	node := &Node{
		prevEpoch:  0,
		thisEpoch:  25519,
		nextEpoch:  0,
		rawStorage: rs,
	}
	rsNew, err := rs.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rsNew.MaxBytes = 1234567890
	nextNode := &Node{
		prevEpoch:  1,
		thisEpoch:  123456,
		nextEpoch:  123456789,
		rawStorage: rsNew,
	}
	if !nextNode.IsValid() {
		t.Fatal("nextNode should not be valid")
	}
	if node.thisEpoch >= nextNode.thisEpoch {
		t.Fatal("We should have node.thisEpoch >= nextNode.thisEpoch")
	}
	if nextNode.IsTail() {
		t.Fatal("We should not have nextNode is tail to raise error")
	}
	err = node.SetEpochs(nil, nextNode)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}
