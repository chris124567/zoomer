package protocol

import (
	"fmt"
	"log"
)

/*

NALU packetizer will chop packets up that are longer than the
allowed size of the Maximum Transmission Unit (called MTU).

This packetizer does not yet strictly conform to what Zoom is
doing. This is because Zoom seems to not be using a different algorithm.

The algorithm we use here is to chop pieces into multiple MTU's, the max size
we then have a little shit packet at the end that could be just a few bytes long.

Zoom instead seems to be using something more advanced will prefers a consistent
packet size across all packets. Maybe they just divvy up the size until it's lower
than the MTU?


*/

const (
	// TODO: figure out actual max MTU size
	MTU_MAX = 900

	MASK_NALU_HEADER_FORBIDDEN_BIT = 0x80
	MASK_NALU_HEADER_NRI_BITS      = 0x60
	MASK_NALU_HEADER_TYPE          = 0x1F
	FU_A                           = 28
	FU_B                           = 29
	BYTE_SINGLE                    = 0x00

	MASK_FU_HEADER_START_BIT    = 0x80
	MASK_FU_HEADER_END_BIT      = 0x40
	MASK_FU_HEADER_RESERVED_BIT = 0x20
	MASK_FU_HEADER_TYPE_BITS    = 0x1F
)

func IsSingle(b byte) bool {
	return b == BYTE_SINGLE
}
func IsFragmented(b byte) bool {
	return b&MASK_NALU_HEADER_TYPE == FU_A
}
func IsFragmentedEnd(b byte) bool {
	return b&MASK_FU_HEADER_END_BIT == MASK_FU_HEADER_END_BIT
}

type NaluPacketizer struct {
	buffer []byte
}

func NewNaluPacketizer() *NaluPacketizer {
	return &NaluPacketizer{
		buffer: make([]byte, 0),
	}
}

func (parser *NaluPacketizer) Unmarshal(data []byte) ([]byte, error) {
	var complete []byte

	// There should at least be one byte to define FU-A or Single
	if len(data) < 1 {
		return nil, ErrInvalidLength
	}

	if IsFragmented(data[0]) {
		// Fragmented should have at least 2 bytes of header and 1 byte of content
		if len(data) < 3 {
			return nil, ErrInvalidLength
		}

		parser.buffer = append(parser.buffer, data[2:]...)

		if IsFragmentedEnd(data[1]) {
			// FU-A end fragments probably have the tag
			complete = parser.buffer
			defer func() {
				parser.buffer = make([]byte, 0)
			}()
		} else {
			// Nothing to do yet, start or continuation
			return nil, nil
		}
	} else if IsSingle(data[0]) {
		// Single should have at least 1 byte of content
		if len(data) < 2 {
			return nil, ErrInvalidLength
		}
		complete = data[1:]
	} else {
		return nil, fmt.Errorf("unknown NALU format %v", data[0])
	}
	return complete, nil
}

func (parser *NaluPacketizer) Marshal(payload []byte) ([][]byte, error) {
	pkts := make([][]byte, 0)
	if len(payload) > MTU_MAX {
		// TODO: fragment
		bytesWritten := 0
		for bytesWritten < len(payload) {
			cursor := bytesWritten
			end := cursor + MTU_MAX
			// Final chunk may be smaller than the MTU
			if len(payload)-bytesWritten < MTU_MAX {
				end = len(payload)
			}

			// Default to continuation bit
			header := byte(0)
			if cursor == 0 {
				// Set start bit
				header |= MASK_FU_HEADER_START_BIT
			} else if end == len(payload) {
				// Set end bit
				header |= MASK_FU_HEADER_END_BIT
			}
			log.Printf("header=%v", header)

			fragment := payload[cursor:end]
			encoded := append([]byte{FU_A, header}, fragment...)
			pkts = append(pkts, encoded)
			bytesWritten += len(fragment)
		}
	} else {
		encoded := append([]byte{BYTE_SINGLE}, payload...)
		pkts = append(pkts, encoded)
	}

	return pkts, nil
}
