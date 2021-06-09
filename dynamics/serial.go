package dynamics

import (
	"errors"

	"github.com/MadBase/MadNet/utils"
)

const (
	typeUint32 uint32 = 1
	typeInt32  uint32 = 2
	typeUint64 uint32 = 3
	typeInt64  uint32 = 4
	upperBound uint32 = 65536
)

// The types here are sorted arrays which hold values.
// blkNum only increases in value;
// the particular value type depends on the function.

type SortedUint32 struct {
	size    uint32
	blkNums []uint32
	values  []uint32
}

func (su *SortedUint32) Marshal() ([]byte, error) {
	if !su.IsValid() {
		return nil, ErrInvalidSortedArray
	}
	fmt := byte(typeUint32)
	sizeBytes := utils.MarshalUint16(uint16(su.size))
	ret := make([]byte, 0, 3+su.size*8)
	ret = append(ret, fmt)
	ret = append(ret, sizeBytes...)
	for k := 0; k < int(su.size); k++ {
		blkBytes := utils.MarshalUint32(su.blkNums[k])
		valueBytes := utils.MarshalUint32(su.values[k])
		ret = append(ret, blkBytes...)
		ret = append(ret, valueBytes...)
	}
	return ret, nil
}

func (su *SortedUint32) Unmarshal(v []byte) error {
	vLen := len(v)
	if vLen < 3 {
		return ErrIncorrectData
	}
	fmt := uint32(v[0])
	if fmt != typeUint32 {
		return ErrInvalidFormat
	}
	size16, _ := utils.UnmarshalUint16(v[1:3])
	// No error checking because length is correct
	size := uint32(size16)

	// Check size
	q := (vLen - 3) / 8
	r := (vLen - 3) % 8
	if (q != int(size)) || (r != 0) {
		return ErrIncorrectData
	}
	su.size = size
	su.blkNums = make([]uint32, int(size))
	su.values = make([]uint32, int(size))
	for k := 0; k < int(su.size); k++ {
		data := utils.CopySlice(v[3+k*8 : 3+(k+1)*8])
		blkNum, _ := utils.UnmarshalUint32(data[0:4])
		su.blkNums[k] = blkNum
		value, _ := utils.UnmarshalUint32(data[4:8])
		su.values[k] = value
	}
	if !su.IsValid() {
		return ErrInvalidSortedArray
	}
	return nil
}

// IsValid returns if the sorted list is valid.
// To be valid, the list must be sorted and must be strictly positive
// in length less than upperBound elements.
func (su *SortedUint32) IsValid() bool {
	size := int(su.size)
	if size < 0 || size >= int(upperBound) {
		// Must have correct size
		return false
	}
	if (size != len(su.blkNums)) || (size != len(su.values)) {
		// Lengths must agree
		return false
	}
	for k := 0; k < size-1; k++ {
		if su.blkNums[k+1] <= su.blkNums[k] {
			// blkNums must be strictly increasing in value
			return false
		}
	}
	return true
}

func (su *SortedUint32) AddElement(blkNum uint32, value uint32) error {
	ok := su.IsValid()
	if !ok {
		return ErrInvalidSortedArray
	}
	if su.size >= upperBound-1 {
		return errors.New("Too large of array")
	}
	if su.size != 0 {
		if su.blkNums[su.size-1] >= blkNum {
			return ErrIncorrectAdd
		}
	}
	su.size++
	su.blkNums = append(su.blkNums, blkNum)
	su.values = append(su.values, value)
	return nil
}

func (su *SortedUint32) SeekValue(blkNum uint32) (uint32, error) {
	// Check array is valid
	ok := su.IsValid()
	if !ok {
		return 0, ErrInvalidSortedArray
	}
	// If empty, then return no value
	if su.size == 0 {
		return 0, ErrEmptyArray
	}
	// Check if block number is larger than last
	if su.blkNums[su.size-1] <= blkNum {
		return su.values[su.size-1], nil
	}
	// If before first element, return error
	if su.blkNums[0] > blkNum {
		return 0, ErrNotPresent
	}

	// Now must perform binary search to find correct location
	//
	// We let B := blkNum, B[i] := blkNums[i]
	//
	// From the above calculations, we have shown
	//		B[0] <= B < B[size-1]
	//
	// At each stage in our search, we have
	//		B[lower] <= B < B[upper]
	//
	// where lower == lowerIdx and upper == upperIdx.
	// We have lower < upper.
	// Because
	//		middle = lower + (upper-lower)/2
	//
	// we always have
	//		lower <= middle < upper
	//
	// It is impossible for middle == upper.
	// This is because
	//		floor(alpha/2) <= alpha/2 < alpha
	//
	// Here, alpha >= 0 is an integer.
	// In our case, take alpha = upper - lower,
	// and the above proposition follows because lower < upper.
	//
	// If we ever have lower == middle, then we are finished and stop.
	lowerIdx := 0
	upperIdx := int(su.size - 1)
	middleIdx := lowerIdx + (upperIdx-lowerIdx)/2
	idx := int(su.size)
	for {
		if lowerIdx == middleIdx {
			idx = middleIdx
			break
		}
		if blkNum >= su.blkNums[middleIdx] {
			lowerIdx = middleIdx
		} else {
			upperIdx = middleIdx
		}
		middleIdx = lowerIdx + (upperIdx-lowerIdx)/2
		continue
	}
	if (su.blkNums[idx] > blkNum) || (su.blkNums[idx+1] < blkNum) {
		// This should never evaluate as true because we are starting
		// with a valid array.
		return 0, errors.New("Invalid: Something went very wrong")
	}
	return su.values[idx], nil
}
