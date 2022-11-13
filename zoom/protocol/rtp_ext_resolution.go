package protocol

import (
	"encoding/binary"
	"fmt"
)

type RtpExtResolution struct {
	Width  uint16
	Height uint16
}

func (ext *RtpExtResolution) Unmarshal(data []byte) error {
	if data == nil {
		return ErrNoData
	}

	ext.Width = binary.BigEndian.Uint16(data[0:2])
	ext.Height = binary.BigEndian.Uint16(data[2:4])

	return nil
}

func (ext *RtpExtResolution) Marshal() ([]byte, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint16(data[0:], ext.Width)
	binary.BigEndian.PutUint16(data[2:], ext.Height)
	return data, nil
}

func (ext *RtpExtResolution) String() string {
	return fmt.Sprintf("[RtpExtResolution width=%v height=%v]", ext.Width, ext.Height)
}
