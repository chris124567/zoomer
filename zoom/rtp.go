package zoom

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/chris124567/zoomer/zoom/api/protocol"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

/*type ZoomHeader struct {
	prefix  uint16
	bodyLen uint16
	ivLen   uint8
	IV      [12]byte
}
*/
const (
	RTP_EXTENSION_ID_UUID       = 7
	RTP_EXTENSION_ID_RESOLUTION = 6
	RTP_EXTENSION_ID_FRAME_INFO = 4
)

type ZoomRtpDecoder struct {
	NaluPacketizer    *protocol.NaluPacketizer
	ParticipantRoster *ZoomParticipantRoster
}

func NewZoomRtpDecoder(participantRoster *ZoomParticipantRoster) *ZoomRtpDecoder {
	return &ZoomRtpDecoder{
		// TODO: NALU packetizer thinks its one big stream, so doesnt multiplex fragmented units
		// very well here, maybe needs to map ssrc -> packetizer
		NaluPacketizer:    protocol.NewNaluPacketizer(),
		ParticipantRoster: participantRoster,
	}
}

func (parser *ZoomRtpDecoder) SetSharedMeetingKey(k []byte) {
	parser.ParticipantRoster.SetSharedMeetingKey(k)
}

func (parser *ZoomRtpDecoder) Decode(rawPkt []byte) (decoded []byte, err error) {
	// 1. Decode the RTP packet
	rtpPacket := &rtp.Packet{}
	err = rtpPacket.Unmarshal(rawPkt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	id := rtpPacket.GetExtension(RTP_EXTENSION_ID_UUID)

	resolutionMeta := &protocol.RtpExtResolution{}
	err = resolutionMeta.Unmarshal(rtpPacket.GetExtension(RTP_EXTENSION_ID_RESOLUTION))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	svcMeta := &protocol.RtpExtFrameInfo{}
	err = svcMeta.Unmarshal(rtpPacket.GetExtension(RTP_EXTENSION_ID_FRAME_INFO))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("rtp header [M = %v] [PT type=%v] [SN seq=%v] [TS timestamp=%v] [P padding=%v size=%v] [ssrc=%v csrc=%v]", rtpPacket.Marker, rtpPacket.PayloadType, rtpPacket.SequenceNumber, rtpPacket.Timestamp, rtpPacket.Padding, rtpPacket.PaddingSize, rtpPacket.SSRC, rtpPacket.CSRC)
	log.Printf("rtp extensions [RtpId id=%v] [meta=%v] [%v]", id, svcMeta, resolutionMeta)
	log.Printf("rtp payload [PYLD size=%v]", len(rtpPacket.Payload))

	payload := rtpPacket.Payload
	// TODO: header length
	if len(payload) < 35 {
		return nil, errors.New("payload does not have enough bytes")
	}

	// 2. Check Fragmented Unit
	complete, err := parser.NaluPacketizer.Unmarshal(payload)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 3. Exit if we don't yet have a complete packet in our NALU packetizer.
	// TODO: ideally this would be more tightly coupled with the loop that it's running in.
	if complete == nil {
		return nil, nil
	}

	// 4. Decode the inner encrypted payload
	decodedPayload := &protocol.RtpEncryptedPayload{}
	decodedPayload.Unmarshal(complete)

	// 5. Decrypt the ciphertext
	secretNonce, err := parser.ParticipantRoster.GetSecretNonceForSSRC((int)(rtpPacket.SSRC))
	if err != nil {
		return nil, fmt.Errorf("cannot decode packet for ssrc=%v", rtpPacket.SSRC)
	}

	sharedMeetingKey := parser.ParticipantRoster.GetSharedMeetingKey()
	if true {
		log.Printf("body key=%v sn=%v iv=%v body=%v tag=%v", hex.EncodeToString(sharedMeetingKey), hex.EncodeToString(secretNonce), hex.EncodeToString(decodedPayload.IV), hex.EncodeToString(decodedPayload.Ciphertext), hex.EncodeToString(decodedPayload.Tag))
	}

	decryptor, err := protocol.NewAesGcmCrypto(sharedMeetingKey, secretNonce)
	if err != nil {
		return nil, err
	}

	ciphertextWithTag := append(decodedPayload.Ciphertext, decodedPayload.Tag...)
	plaintext, err := decryptor.Decrypt(decodedPayload.IV, ciphertextWithTag)
	if err != nil {
		return nil, err
	}
	log.Printf("rtp: decrypted=%v", hex.EncodeToString(plaintext))

	decoded = plaintext
	return
}

type ZoomRtpEncoder struct {
	id             int
	ssrc           int
	resolution     *protocol.RtpExtResolution
	messageCounter int
	timestamp      int

	// Frame information
	baseFrameCounter    int
	currentFrameCounter int

	// Encryption information
	ParticipantRoster *ZoomParticipantRoster
}

func NewZoomRtpEncoder(roster *ZoomParticipantRoster, ssrc int, id int, width, height int) *ZoomRtpEncoder {
	return &ZoomRtpEncoder{
		id:   id,
		ssrc: ssrc,
		resolution: &protocol.RtpExtResolution{
			Width:  uint16(width),
			Height: uint16(height),
		},

		messageCounter: 0,
		timestamp:      2746202358,

		ParticipantRoster: roster,
	}
}

func (parser *ZoomRtpEncoder) Encode(payload []byte) ([]byte, error) {

	encryptedPayload, err := encryptPayloadWithRoster(parser.ParticipantRoster, parser.ssrc, parser.messageCounter, payload)
	if err != nil {
		return nil, err
	}

	// TODO: properly NALU packetize the payload if too long
	packetizedPayload := append(make([]byte, 1), encryptedPayload...)

	// Wrap encoded payload in RTP packets
	// TODO: multiple packets, currently only single packets
	parser.timestamp = parser.timestamp + (parser.messageCounter * 15000)
	p := &rtp.Packet{
		Header: rtp.Header{
			Version:          2,
			Padding:          false,
			Extension:        false,
			Marker:           false,
			PayloadType:      99,
			SequenceNumber:   uint16(parser.messageCounter),
			Timestamp:        uint32(parser.timestamp),
			SSRC:             uint32(parser.ssrc),
			CSRC:             []uint32{},
			ExtensionProfile: 0,
			Extensions:       []rtp.Extension{},
		},
		Payload:     packetizedPayload,
		PaddingSize: 0,
	}

	// TODO: increment UUID whenever reconnecting, prob big endian but single byte?
	p.Header.SetExtension(RTP_EXTENSION_ID_UUID, []byte{0x01})
	if parser.resolution != nil {
		resolution, err := parser.resolution.Marshal()
		if err != nil {
			return nil, err
		}
		p.Header.SetExtension(RTP_EXTENSION_ID_RESOLUTION, resolution)
	}

	// Set extension screen size
	rtpResolution, err := parser.resolution.Marshal()
	if err != nil {
		return nil, err
	}
	p.Header.SetExtension(RTP_EXTENSION_ID_RESOLUTION, rtpResolution)

	// TODO: extension frame info should probably be more advanced
	rtpFrameInfo := &protocol.RtpExtFrameInfo{
		Version:       2,
		Start:         true,  // TODO: NALU packetizer
		End:           true,  // TODO: NALU packetize
		Independent:   false, // TODO: passed along with payload
		Required:      false, // TODO: passed along with payload
		Base:          false, // TODO: passed along with payload
		TemporalID:    0,     // TODO: ????
		CurrentFrame:  uint16(parser.currentFrameCounter + 1),
		PreviousFrame: uint16(parser.currentFrameCounter),
		BaseFrame:     uint16(parser.baseFrameCounter),
	}
	rtpFrameInfoBytes, err := rtpFrameInfo.Marshal()
	if err != nil {
		return nil, err
	}
	p.Header.SetExtension(RTP_EXTENSION_ID_FRAME_INFO, rtpFrameInfoBytes)

	rawPkt, err := p.Marshal()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return rawPkt, nil
}

func encryptPayloadWithRoster(roster *ZoomParticipantRoster, ssrc int, messageCounter int, plaintext []byte) ([]byte, error) {
	secretNonce, err := roster.GetSecretNonceForSSRC(ssrc)
	if err != nil {
		return nil, ErrParticipantExists
	}

	sharedMeetingKey := roster.GetSharedMeetingKey()

	// Build 12 byte IV
	IV := make([]byte, 12)
	binary.BigEndian.PutUint16(IV, uint16(messageCounter))
	IV = IV[:12]

	encryptor, err := protocol.NewAesGcmCrypto(sharedMeetingKey, secretNonce)
	if err != nil {
		return nil, err
	}

	ciphertextWithTag, err := encryptor.Encrypt(IV, plaintext)
	if err != nil {
		return nil, err
	}

	if true {
		log.Printf("encrypted body key=%v sn=%v iv=%v ciphertextWithTag=%v", hex.EncodeToString(sharedMeetingKey), hex.EncodeToString(secretNonce), hex.EncodeToString(IV), hex.EncodeToString(ciphertextWithTag))
	}

	encodedPayload := protocol.NewRtpEncryptedPayload(IV, ciphertextWithTag)
	encodedPayloadInBytes := encodedPayload.Marshal()
	return encodedPayloadInBytes, nil
}

func RtcpProcess(rawPkt []byte) ([]rtcp.Packet, error) {
	p, err := rtcp.Unmarshal(rawPkt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// log.Printf("%v", p[0])
	return p, nil
}
