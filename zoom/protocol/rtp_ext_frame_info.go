package protocol

import (
	"encoding/binary"
	"fmt"
)

/* H264 SVC layer / or frame information

Zoom uses a mangled form of "Frame Marking RTP Header Extension"
https://www.ietf.org/id/draft-ietf-avtext-framemarking-13.html

+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  ID=? |  L=1  |S|E|I|D|B| TID |   LID         | (TL0PICIDX omitted)
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

The most obvious distinction between the spec and Zoom's implementation is that theirs
contains additional data that contains 3 frame counters. While the exact names and uses
for these frame counters is not known, they do behave very consistently which allows us
to make a guess at what the might be doing.

I took the liberty to name these frame counters, giving a bit of meaning:
1. Current Frame: increases on every frame (fragmented packets contain the same counter)
2. Previous Frame: seems to always indicate the previous frame, but this may indicate the temporal layer
3. Base Frame: seems to only change when serious changes occur such as a lot of movement or screen size changes

Eitherway, the exact meaning of these frame counters is currently unknown.

Q: Why is there a need to include this fragmentation info in an extension header?
A?: Maybe in an e2ee chat we can't distinguish keyframes from other frames so
    they are marked in order that the SFU can figure which frames to send

[2 192 0 114 0 113 0 1] =>
[	2 		= version?
	192 	= pkt info?
	0 114 	= current frame/layer counter
	0 113 	= previous frame/layer counter
	0 1		= base reference frame counter
]

Fragmented:
	10  11 1000 184 -> start IDR frame
	00  11 1000 56 -> continuation IDR bits (probably last bits)
	01  11 1000 120 -> end IDR frame

	// This can't only be for pred frames, also gets used for IDR
	10  10 0100 164 -> start PRED frame
	00	10 0100 36 -> continuation PRED frame
	01	10 0100 100 -> end PRED frame

	// Med intensity changes
	10 00 0000  128 -> start
	00 00 0000	0 	-> continuation
	01 00 0000  64  -> end

	// High intensity changes
	10	10 0000 160 -> start
	00	10 0000 32 ->  continue single (probably predicted frames)
	01	10 0000 96 ->   end
Unfragmented:
	11  00 0000 192 -> single (probably predicted frames)

	pkt_info bytes [S, E, I, R, B, T, ??, ??]
	S = start
	E = end
	I = independent frame
	R = required (seems to be inverted)
	B = Base Layer Sync = 1 if only depends on base layer (resolution changes SPS PPS, I frames)
	T = Temporal Layer Sync = 1 if P frame


baseFrame only changes if there have been spatial changes (screen resolution)
	rtp id=[3] version=2 pktInfo=192 curFrame=2568 prevFrame=2567 baseFrame=229 width=1920 height=924
	rtp id=[3] version=2 pktInfo=184 curFrame=2569 prevFrame=2569 baseFrame=2569 width=1150 height=914
*/
type RtpExtFrameInfo struct {
	Version       uint8
	Start         bool
	End           bool
	Independent   bool
	Required      bool
	Base          bool
	TemporalID    uint8
	CurrentFrame  uint16
	PreviousFrame uint16
	BaseFrame     uint16
}

func (ext *RtpExtFrameInfo) Unmarshal(data []byte) error {
	if data == nil {
		return ErrNoData
	}

	ext.Version = data[0]

	pktInfo := data[1]
	ext.Start = pktInfo&0x80 == 0x80
	ext.End = pktInfo&0x40 == 0x40
	ext.Independent = pktInfo&0x20 == 0x20
	ext.Required = pktInfo&0x10 == 0x10
	ext.Base = pktInfo&0x08 == 0x08
	ext.TemporalID = pktInfo & 0x07

	ext.CurrentFrame = binary.BigEndian.Uint16(data[2:4])
	ext.PreviousFrame = binary.BigEndian.Uint16(data[4:6])
	ext.BaseFrame = binary.BigEndian.Uint16(data[6:8])
	return nil
}

func (ext *RtpExtFrameInfo) Marshal() ([]byte, error) {

	data := make([]byte, 0)
	data = append(data, ext.Version)

	pktInfo := ext.TemporalID
	if ext.Start {
		pktInfo |= 0x80
	}
	if ext.End {
		pktInfo |= 0x40
	}
	if ext.Independent {
		pktInfo |= 0x40
	}
	if ext.Required {
		pktInfo |= 0x10
	}
	if ext.Base {
		pktInfo |= 0x08
	}
	data = append(data, pktInfo)

	currentFrame := make([]byte, 2)
	binary.BigEndian.PutUint16(currentFrame, ext.CurrentFrame)
	data = append(data, currentFrame...)

	prevFrame := make([]byte, 2)
	binary.BigEndian.PutUint16(prevFrame, ext.PreviousFrame)
	data = append(data, prevFrame...)

	baseFrame := make([]byte, 2)
	binary.BigEndian.PutUint16(baseFrame, ext.BaseFrame)
	data = append(data, baseFrame...)

	return data, nil
}

func (ext *RtpExtFrameInfo) String() string {
	return fmt.Sprintf("[RtpExtFrameInfo start=%v end=%v indep=%v req=%v, base=%v temp=%v cur=%v prev=%v base=%v]", ext.Start, ext.End, ext.Independent, ext.Required, ext.Base, ext.TemporalID, ext.CurrentFrame, ext.PreviousFrame, ext.BaseFrame)
}
