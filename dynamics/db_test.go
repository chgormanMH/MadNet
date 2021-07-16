package dynamics

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

type mockRawDB struct {
	rawDB map[string]string
}

func (m *mockRawDB) GetValue(key []byte) ([]byte, error) {
	strKey := string(key)
	strValue, ok := m.rawDB[strKey]
	if !ok {
		return nil, ErrKeyNotPresent
	}
	value := []byte(strValue)
	return value, nil
}

func (m *mockRawDB) SetValue(key []byte, value []byte) error {
	strKey := string(key)
	strValue := string(value)
	m.rawDB[strKey] = strValue
	return nil
}

func (m *mockRawDB) DeleteValue(key []byte) error {
	strKey := string(key)
	_, ok := m.rawDB[strKey]
	if !ok {
		return ErrKeyNotPresent
	}
	delete(m.rawDB, strKey)
	return nil
}

func TestMock(t *testing.T) {
	key := []byte("Key")
	value := []byte("Key")

	m := &mockRawDB{}
	m.rawDB = make(map[string]string)

	_, err := m.GetValue(key)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	err = m.SetValue(key, value)
	if err != nil {
		t.Fatal(err)
	}

	retValue, err := m.GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(retValue, value) {
		t.Fatal("values do not match")
	}

	err = m.DeleteValue(key)
	if err != nil {
		t.Fatal(err)
	}
	_, err = m.GetValue(key)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	return logger
}

func initializeDB() *Database {
	logger := newLogger()
	db := &Database{}
	db.logger = logger
	mock := &mockRawDB{}
	mock.rawDB = make(map[string]string)
	db.rawDB = mock
	return db
}

func TestGetSetNode(t *testing.T) {
	db := initializeDB()

	node := &Node{}
	err := db.SetNode(node)
	if err == nil {
		t.Fatal("Should have raised error")
	}

	node.prevEpoch = 1
	node.thisEpoch = 1
	node.nextEpoch = 1
	node.rawStorage = &RawStorage{}
	err = db.SetNode(node)
	if err != nil {
		t.Fatal(err)
	}
	nodeBytes, err := node.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	node2, err := db.GetNode(1)
	if err != nil {
		t.Fatal(err)
	}
	node2Bytes, err := node2.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(nodeBytes, node2Bytes) {
		t.Fatal("nodes do not match")
	}
}

func TestGetSetLinkedList(t *testing.T) {
	db := initializeDB()

	ll := &LinkedList{}
	err := db.SetLinkedList(ll)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	ll.epochLastUpdated = 1
	ll.currentEpoch = 1
	err = db.SetLinkedList(ll)
	if err != nil {
		t.Fatal(err)
	}
	llBytes := ll.Marshal()

	ll2, err := db.GetLinkedList()
	if err != nil {
		t.Fatal(err)
	}
	ll2Bytes := ll2.Marshal()
	if !bytes.Equal(llBytes, ll2Bytes) {
		t.Fatal("LinkedLists do not match")
	}
}
