package dynamics

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// Database is an abstraction for object storage
type Database struct {
	sync.Mutex
	rawDB  rawDataBase
	logger *logrus.Logger
}

func (db *Database) SetNode(node *Node) error {
	if !node.IsValid() {
		return ErrInvalidNode
	}
	nodeKey, err := makeNodeKey(node.thisEpoch)
	if err != nil {
		return err
	}
	key, err := nodeKey.Marshal()
	if err != nil {
		return err
	}
	nodeBytes, err := node.Marshal()
	if err != nil {
		return err
	}
	err = db.rawDB.SetValue(key, nodeBytes)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetNode(epoch uint32) (*Node, error) {
	nodeKey, err := makeNodeKey(epoch)
	if err != nil {
		return nil, err
	}
	key, err := nodeKey.Marshal()
	if err != nil {
		return nil, err
	}
	v, err := db.rawDB.GetValue(key)
	if err != nil {
		return nil, err
	}
	node := &Node{}
	err = node.Unmarshal(v)
	if err != nil {
		return nil, err
	}
	if !node.IsValid() {
		return nil, ErrInvalidNode
	}
	return node, nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// SetLinkedList saves LinkedList to the database
func (db *Database) SetLinkedList(ll *LinkedList) error {
	value, err := ll.Marshal()
	if err != nil {
		return err
	}
	llKey := makeLinkedListKey()
	key, err := llKey.Marshal()
	if err != nil {
		return err
	}
	err = db.rawDB.SetValue(key, value)
	if err != nil {
		return err
	}
	return nil
}

// GetLinkedList retrieves LinkedList from the database
func (db *Database) GetLinkedList() (*LinkedList, error) {
	llKey := makeLinkedListKey()
	key, err := llKey.Marshal()
	if err != nil {
		return nil, err
	}
	v, err := db.rawDB.GetValue(key)
	if err != nil {
		return nil, err
	}
	ll := &LinkedList{}
	err = ll.Unmarshal(v)
	if err != nil {
		return nil, err
	}
	return ll, nil
}
