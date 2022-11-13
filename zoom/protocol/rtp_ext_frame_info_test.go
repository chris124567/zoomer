package protocol

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"
)

func TestUnmarshalRtpExtensionFrameInfo(t *testing.T) {
	payload, err := hex.DecodeString("02c007030702055f")
	if err != nil {
		t.Error(err)
		return
	}

	extFrameInfo := &RtpExtFrameInfo{}
	err = extFrameInfo.Unmarshal(payload)
	if err != nil {
		t.Error(err)
		return
	}

	expectedExtFrameInfo := &RtpExtFrameInfo{
		Version:       uint8(2),
		Start:         true,
		End:           true,
		Independent:   false,
		Required:      false,
		Base:          false,
		TemporalID:    uint8(0),
		CurrentFrame:  uint16(1795),
		PreviousFrame: uint16(1794),
		BaseFrame:     uint16(1375),
	}

	if *extFrameInfo != *expectedExtFrameInfo {
		t.Error("Does not equal")
		return
	}

}

func TestMarshalRtpExtensionFrameInfo(t *testing.T) {
	extFrameInfo := &RtpExtFrameInfo{
		Version:       uint8(2),
		Start:         true,
		End:           true,
		Independent:   false,
		Required:      false,
		Base:          false,
		TemporalID:    uint8(0),
		CurrentFrame:  uint16(1795),
		PreviousFrame: uint16(1794),
		BaseFrame:     uint16(1375),
	}

	payload, err := extFrameInfo.Marshal()
	if err != nil {
		t.Error(err)
		return
	}

	log.Printf("payload = %v", hex.EncodeToString(payload))
	expectedPayload, err := hex.DecodeString("02c007030702055f")
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(payload, expectedPayload) != 0 {
		t.Error("payload did not match expected")
	}
}
