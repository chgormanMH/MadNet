package dynamics

import (
	"bytes"
	"testing"
)

func TestSU32IsValid(t *testing.T) {
	su := &SortedUint32{}
	su.size = upperBound
	ok := su.IsValid()
	if ok {
		t.Fatal("Should not be ok (2)")
	}

	su.size = 1
	su.blkNums = []uint32{1, 2}
	su.values = []uint32{0, 0}
	ok = su.IsValid()
	if ok {
		t.Fatal("Should not be ok (3)")
	}

	su.size = 2
	su.blkNums = []uint32{1, 2}
	su.values = []uint32{0, 0, 0}
	ok = su.IsValid()
	if ok {
		t.Fatal("Should not be ok (4)")
	}

	su.size = 2
	su.blkNums = []uint32{2, 2}
	su.values = []uint32{0, 0}
	ok = su.IsValid()
	if ok {
		t.Fatal("Should not be ok (5)")
	}

	su.size = 0
	su.blkNums = []uint32{}
	su.values = []uint32{}
	ok = su.IsValid()
	if !ok {
		t.Fatal("Should be ok (1)")
	}

	su.size = 2
	su.blkNums = []uint32{1, 2}
	su.values = []uint32{0, 1}
	ok = su.IsValid()
	if !ok {
		t.Fatal("Should be ok (2)")
	}
}

func TestSU32Marshal(t *testing.T) {
	su := &SortedUint32{}
	ret, err := su.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	retTrue := []byte{byte(typeUint32), 0, 0}
	if !bytes.Equal(ret, retTrue) {
		t.Fatal("Invalid marshalled byte slice (1)")
	}

	su.size = 2
	su.blkNums = []uint32{1, 2}
	su.values = []uint32{0, 1}
	ret, err = su.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	retTrue = []byte{byte(typeUint32), 0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1}
	if !bytes.Equal(ret, retTrue) {
		t.Fatal("Invalid marshalled byte slice (2)")
	}

	su.size = 2
	su.blkNums = []uint32{1, 2, 3}
	su.values = []uint32{0, 1}
	_, err = su.Marshal()
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestSU32Unmarshal(t *testing.T) {
	v := []byte{}
	su := &SortedUint32{}
	err := su.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	v = []byte{0, 0, 0}
	err = su.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	v = []byte{byte(typeUint32), 0, 0}
	err = su.Unmarshal(v)
	if err != nil {
		t.Fatal(err)
	}

	v = []byte{byte(typeUint32), 0, 0, 0}
	err = su.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}

	v = []byte{byte(typeUint32), 0, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1}
	err = su.Unmarshal(v)
	if err != nil {
		t.Fatal(err)
	}
	suTrue := &SortedUint32{}
	suTrue.size = 2
	suTrue.blkNums = []uint32{1, 2}
	suTrue.values = []uint32{0, 1}

	v = []byte{byte(typeUint32), 0, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1}
	err = su.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (4)")
	}
}

func TestSU32AddElement(t *testing.T) {
	su := &SortedUint32{}
	ok := su.IsValid()
	if !ok {
		t.Fatal("Should be ok (1)")
	}

	blkNum := uint32(1)
	value := uint32(7)
	err := su.AddElement(blkNum, value)
	if err != nil {
		t.Fatal(err)
	}
	ok = su.IsValid()
	if !ok {
		t.Fatal("Should be ok (2)")
	}

	blkNum = 2
	value = 5
	err = su.AddElement(blkNum, value)
	if err != nil {
		t.Fatal(err)
	}
	ok = su.IsValid()
	if !ok {
		t.Fatal("Should be ok (3)")
	}

	blkNum = 1
	value = 10
	err = su.AddElement(blkNum, value)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	su.size = 0
	su.blkNums = []uint32{1, 2}
	su.values = []uint32{1}
	err = su.AddElement(blkNum, value)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	su.size = 0
	su.blkNums = []uint32{}
	su.values = []uint32{}
	for k := 0; k < int(upperBound-1); k++ {
		blkNum = uint32(k + 1)
		value = uint32(3*k + 1)
		err = su.AddElement(blkNum, value)
		if err != nil {
			t.Fatal(err)
		}
	}
	blkNum = upperBound + 1
	value = 0
	err = su.AddElement(blkNum, value)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}
}

func TestSU32SeekValue(t *testing.T) {
	maxBlock := 10
	su := &SortedUint32{}
	for k := 0; k < maxBlock; k++ {
		err := su.AddElement(uint32(k), uint32(k)*uint32(k))
		if err != nil {
			t.Fatal(err)
		}
	}
	blkNum := uint32(7)
	value, err := su.SeekValue(blkNum)
	if err != nil {
		t.Fatal(err)
	}
	if value != (blkNum * blkNum) {
		t.Fatal("Invalid value (1)")
	}

	blkNum = uint32(maxBlock) + 1
	value, err = su.SeekValue(blkNum)
	if err != nil {
		t.Fatal(err)
	}
	if value != uint32((maxBlock-1)*(maxBlock-1)) {
		t.Fatal("Invalid value (2)")
	}
}

func TestSU32SeekValueFail(t *testing.T) {
	su := &SortedUint32{}
	su.size = 0
	su.blkNums = []uint32{1, 2}
	su.values = []uint32{1}
	_, err := su.SeekValue(0)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	su.size = 0
	su.blkNums = []uint32{}
	su.values = []uint32{}
	_, err = su.SeekValue(0)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}

	err = su.AddElement(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = su.SeekValue(0)
	if err == nil {
		t.Fatal("Should have raised error (3)")
	}
}
