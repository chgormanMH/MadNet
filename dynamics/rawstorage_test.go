package dynamics

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func TestRawStorageMarshal(t *testing.T) {
	rs := &RawStorage{}
	_, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	s := &Storage{}
	_, err = s.rawStorage.Marshal()
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestRawStorageUnmarshal(t *testing.T) {
	rs := &RawStorage{}
	v, err := rs.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rs2 := &RawStorage{}
	err = rs2.Unmarshal(v)
	if err != nil {
		t.Fatal(err)
	}

	v = []byte{}
	rs3 := &RawStorage{}
	err = rs3.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (1)")
	}

	s := &Storage{}
	err = s.rawStorage.Unmarshal(v)
	if err == nil {
		t.Fatal("Should have raised error (2)")
	}
}

func TestRawStorageCopy(t *testing.T) {
	// Copy empty RawStorage
	rs1 := &RawStorage{}
	rs2, err := rs1.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rs1Bytes, err := rs1.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rs2Bytes, err := rs2.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rs1Bytes, rs2Bytes) {
		t.Fatal("Should have equal bytes (1)")
	}

	// Copy RawStorage with parameters
	rs1.standardParameters()
	rs2, err = rs1.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rs1Bytes, err = rs1.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rs2Bytes, err = rs2.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rs1Bytes, rs2Bytes) {
		t.Fatal("Should have equal bytes (2)")
	}

	// Copy RawStorage with some parameters set to zero
	rs1.MaxBytes = 0
	rs2, err = rs1.Copy()
	if err != nil {
		t.Fatal(err)
	}
	rs1Bytes, err = rs1.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	rs2Bytes, err = rs2.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rs1Bytes, rs2Bytes) {
		t.Fatal("Should have equal bytes (3)")
	}

	s := &Storage{}
	_, err = s.rawStorage.Copy()
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

func TestRawStorageStandardParameters(t *testing.T) {
	rs := &RawStorage{}
	rs.standardParameters()

	retMaxBytes := rs.GetMaxBytes()
	if retMaxBytes != maxBytes {
		t.Fatal("Should be equal (1)")
	}

	retMaxProposalSize := rs.GetMaxProposalSize()
	if retMaxProposalSize != maxProposalSize {
		t.Fatal("Should be equal (2)")
	}

	retSrvrMsgTimeout := rs.GetSrvrMsgTimeout()
	if retSrvrMsgTimeout != srvrMsgTimeout {
		t.Fatal("Should be equal (3)")
	}

	retMsgTimeout := rs.GetMsgTimeout()
	if retMsgTimeout != msgTimeout {
		t.Fatal("Should be equal (4)")
	}

	retProposalTimeout := rs.GetProposalStepTimeout()
	if retProposalTimeout != proposalStepTO {
		t.Fatal("Should be equal (5)")
	}

	retPreVoteTimeout := rs.GetPreVoteStepTimeout()
	if retPreVoteTimeout != preVoteStepTO {
		t.Fatal("Should be equal (6)")
	}

	retPreCommitTimeout := rs.GetPreCommitStepTimeout()
	if retPreCommitTimeout != preCommitStepTO {
		t.Fatal("Should be equal (7)")
	}

	retDBRNRTO := rs.GetDeadBlockRoundNextRoundTimeout()
	if retDBRNRTO != dBRNRTO {
		t.Fatal("Should be equal (8)")
	}

	retDownloadTO := rs.GetDownloadTimeout()
	if retDownloadTO != downloadTO {
		t.Fatal("Should be equal (9)")
	}
}

func TestRawStorageMaxBytes(t *testing.T) {
	rs := &RawStorage{}
	retMaxBytes0 := rs.GetMaxBytes()
	if retMaxBytes0 != 0 {
		t.Fatal("Should be zero")
	}

	value := uint32(10000)
	rs.SetMaxBytes(value)
	retMaxBytes := rs.GetMaxBytes()
	if retMaxBytes != value {
		t.Fatal("Should be equal (1)")
	}

	retMaxProposalSize := rs.GetMaxProposalSize()
	if retMaxProposalSize != value {
		t.Fatal("Should be equal (2)")
	}
}

func TestRawStorageMaxProposalSize(t *testing.T) {
	rs := &RawStorage{}
	retMaxProposalSize0 := rs.GetMaxProposalSize()
	if retMaxProposalSize0 != 0 {
		t.Fatal("Should be zero (2)")
	}

	value := uint32(10000)
	rs.SetMaxBytes(value)
	retMaxProposalSize := rs.GetMaxProposalSize()
	if retMaxProposalSize != value {
		t.Fatal("Should be equal (2)")
	}
}

func TestRawStorageMsgTimeout(t *testing.T) {
	rs := &RawStorage{}
	retMsgTimeout0 := rs.GetMsgTimeout()
	if retMsgTimeout0 != 0 {
		t.Fatal("Should be zero")
	}

	value := time.Second
	rs.SetMsgTimeout(value)
	retMsgTimeout := rs.GetMsgTimeout()
	if retMsgTimeout != value {
		t.Fatal("Should be equal (1)")
	}

	valueSrvrMsg := (3 * value) / 4
	retSrvrMsgTimeout := rs.GetSrvrMsgTimeout()
	if retSrvrMsgTimeout != valueSrvrMsg {
		t.Fatal("Should be equal (2)")
	}
}

func TestRawStorageSrvrMsgTimeout(t *testing.T) {
	rs := &RawStorage{}
	retSrvrMsgTimeout0 := rs.GetSrvrMsgTimeout()
	if retSrvrMsgTimeout0 != 0 {
		t.Fatal("Should be zero")
	}

	value := time.Second
	rs.SetMsgTimeout(value)
	valueSrvrMsg := (3 * value) / 4
	retSrvrMsgTimeout := rs.GetSrvrMsgTimeout()
	if retSrvrMsgTimeout != valueSrvrMsg {
		t.Fatal("Should be equal")
	}
}

func TestRawStorageConsensusTimeouts(t *testing.T) {
	rs := &RawStorage{}

	retPropTOv0 := rs.GetProposalStepTimeout()
	if retPropTOv0 != 0 {
		t.Fatal("Should be zero (1)")
	}
	retPreVoteTOv0 := rs.GetPreVoteStepTimeout()
	if retPreVoteTOv0 != 0 {
		t.Fatal("Should be zero (2)")
	}
	retPreCommitTOv0 := rs.GetPreCommitStepTimeout()
	if retPreCommitTOv0 != 0 {
		t.Fatal("Should be zero (3)")
	}
	retDownloadTOv0 := rs.GetDownloadTimeout()
	if retDownloadTOv0 != 0 {
		t.Fatal("Should be zero (4)")
	}
	retDBRNRTOv0 := rs.GetDeadBlockRoundNextRoundTimeout()
	if retDBRNRTOv0 != 0 {
		t.Fatal("Should be zero (5)")
	}

	propValue := 10 * time.Second
	rs.SetProposalStepTimeout(propValue)

	retPropTOv1 := rs.GetProposalStepTimeout()
	if retPropTOv1 != propValue {
		t.Fatal("Should be equal (1)")
	}
	retPreVoteTOv1 := rs.GetPreVoteStepTimeout()
	if retPreVoteTOv1 != 0 {
		t.Fatal("Should be zero (6)")
	}
	retPreCommitTOv1 := rs.GetPreCommitStepTimeout()
	if retPreCommitTOv1 != 0 {
		t.Fatal("Should be zero (7)")
	}
	retDownloadTOv1 := rs.GetDownloadTimeout()
	if retDownloadTOv1 != propValue {
		t.Fatal("Should be equal (2)")
	}
	retDBRNRTOv1 := rs.GetDeadBlockRoundNextRoundTimeout()
	if retDBRNRTOv1 != ((5 * propValue) / 2) {
		t.Fatal("Should be equal (3)")
	}

	preVoteValue := 20 * time.Second
	rs.SetPreVoteStepTimeout(preVoteValue)

	retPropTOv2 := rs.GetProposalStepTimeout()
	if retPropTOv2 != propValue {
		t.Fatal("Should be equal (4)")
	}
	retPreVoteTOv2 := rs.GetPreVoteStepTimeout()
	if retPreVoteTOv2 != preVoteValue {
		t.Fatal("Should be equal (5)")
	}
	retPreCommitTOv2 := rs.GetPreCommitStepTimeout()
	if retPreCommitTOv2 != 0 {
		t.Fatal("Should be zero (8)")
	}
	retDownloadTOv2 := rs.GetDownloadTimeout()
	if retDownloadTOv2 != (propValue + preVoteValue) {
		t.Fatal("Should be equal (6)")
	}
	retDBRNRTOv2 := rs.GetDeadBlockRoundNextRoundTimeout()
	if retDBRNRTOv2 != ((5 * (propValue + preVoteValue)) / 2) {
		t.Fatal("Should be equal (7)")
	}

	preCommitValue := 30 * time.Second
	rs.SetPreCommitStepTimeout(preCommitValue)

	retPropTOv3 := rs.GetProposalStepTimeout()
	if retPropTOv3 != propValue {
		t.Fatal("Should be equal (8)")
	}
	retPreVoteTOv3 := rs.GetPreVoteStepTimeout()
	if retPreVoteTOv3 != preVoteValue {
		t.Fatal("Should be equal (9)")
	}
	retPreCommitTOv3 := rs.GetPreCommitStepTimeout()
	if retPreCommitTOv3 != preCommitValue {
		t.Fatal("Should be equal (10)")
	}
	retDownloadTOv3 := rs.GetDownloadTimeout()
	if retDownloadTOv3 != (propValue + preVoteValue + preCommitValue) {
		t.Fatal("Should be equal (11)")
	}
	retDBRNRTOv3 := rs.GetDeadBlockRoundNextRoundTimeout()
	if retDBRNRTOv3 != ((5 * (propValue + preVoteValue + preCommitValue)) / 2) {
		t.Fatal("Should be equal (12)")
	}
}

/*
func TestRawStorageUpdateValueGood(t *testing.T) {
	rs := &RawStorage{}

	retMaxBytes := rs.GetMaxBytes()
	if retMaxBytes != 0 {
		t.Fatal("value should be zero")
	}

	field := "maxBytes"
	value := uint32(1000)
	valueStr := strconv.Itoa(int(value))
	err := rs.UpdateValue(field, valueStr)
	if err != nil {
		t.Fatal(err)
	}
	retMaxBytesNew := rs.GetMaxBytes()
	if retMaxBytesNew != value {
		t.Fatal("incorrect value: should match submitted value")
	}
}

func TestRawStorageUpdateBad(t *testing.T) {
	rs := &RawStorage{}

	field := "thisShouldFail"
	value := uint32(1000)
	valueStr := strconv.Itoa(int(value))
	err := rs.UpdateValue(field, valueStr)
	if !errors.Is(err, ErrInvalidUpdateValue) {
		t.Fatal("Did not raise an error or the correct error")
	}
}
*/

func TestMakeJSONBytes(t *testing.T) {
	field := "field"
	value := "value"
	jsonBytesTrue := []byte("{\"" + field + "\":" + value + "}")
	jsonBytes := makeJSONBytes(field, value)
	if !bytes.Equal(jsonBytes, jsonBytesTrue) {
		t.Fatal("jsonBytes do not agree")
	}
}

// Should produce valid update value
func TestCheckUpdateValueGood(t *testing.T) {
	fieldGood := "maxBytes"
	valueGood := "100000"
	jsonBytes, err := checkUpdateValue(fieldGood, valueGood)
	if err != nil {
		t.Fatal(err)
	}
	jsonBytesTrue := makeJSONBytes(fieldGood, valueGood)
	if !bytes.Equal(jsonBytes, jsonBytesTrue) {
		t.Fatal("invalid jsonBytes returned")
	}
}

// Should produce error for invalid field.
func TestCheckUpdateValueBad1(t *testing.T) {
	fieldBad1 := "field"
	valueBad1 := "1000"
	_, err := checkUpdateValue(fieldBad1, valueBad1)
	if !errors.Is(err, ErrInvalidUpdateValue) {
		t.Fatal("Should have raised error for invalid update value")
	}
}

// Should produce an error for submitting correct field but invalid value
func TestCheckUpdateValueBad2(t *testing.T) {
	fieldBad2 := "maxBytes"
	valueBad2 := "\"value\""
	_, err := checkUpdateValue(fieldBad2, valueBad2)
	if err == nil {
		t.Fatal("Should have raised error")
	}
}

// Should produce an error for submitting correct field but invalid value
func TestCheckUpdateValueBad3(t *testing.T) {
	fieldBad3 := "maxBytes"
	valueBad3 := "value"
	_, err := checkUpdateValue(fieldBad3, valueBad3)
	if !errors.Is(err, ErrInvalidUpdateValue) {
		t.Fatal("Should have raised error")
	}
}
