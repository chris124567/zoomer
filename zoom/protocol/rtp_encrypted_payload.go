package protocol

import (
	"encoding/binary"
)

var (
	PREFIX = []byte{0, 0}
	SUFFIX = []byte{0}
)

const (
	LEN_PREFIX = 2 // Should match len(PREFIX)
	LEN_SUFFIX = 1 // Should match len(SUFFIX)
	// LEN_IV     = 12
	LEN_TAG = 16
)

type RtpEncryptedPayload struct {
	IV         []byte
	Ciphertext []byte
	Tag        []byte
}

func NewRtpEncryptedPayload(IV []byte, CiphertextWithTag []byte) *RtpEncryptedPayload {
	// TODO: assert it has at least length of tag or more
	Ciphertext := CiphertextWithTag[:len(CiphertextWithTag)-LEN_TAG]
	Tag := CiphertextWithTag[len(CiphertextWithTag)-LEN_TAG:]
	return &RtpEncryptedPayload{
		IV,
		Ciphertext,
		Tag,
	}
}

func (encryptedPayload *RtpEncryptedPayload) Unmarshal(payload []byte) {
	lenActualHeader := 15
	lenHeader := LEN_PREFIX /* 00 prefix */ + lenActualHeader /* actual header */ + LEN_SUFFIX /* 0 suffix */

	// TODO: assert payload has size of header
	header := payload[LEN_PREFIX:lenActualHeader]
	lenIV := int(header[2])
	// TODO: assert length of IV with 12
	IV := header[3 : 3+lenIV]

	lenBody := int(binary.BigEndian.Uint16(header[0:2]))
	// TODO: assert payload has size of header + body
	ciphertext := payload[lenHeader : lenHeader+lenBody]
	tag := payload[lenHeader+lenBody:]

	encryptedPayload.IV = IV
	encryptedPayload.Ciphertext = ciphertext
	encryptedPayload.Tag = tag
}

func (encryptedPayload *RtpEncryptedPayload) Marshal() []byte {
	// Parse the length of the ciphertext to a big endian uint16 bytes
	lenBody := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBody, uint16(len(encryptedPayload.Ciphertext)))

	// Parse the length of the IV to a big endian uint8 byte
	lenIV := []byte{uint8(len(encryptedPayload.IV))}
	// TODO: might need range checks to ensure validity

	/*
	 Finalize encrypted payload
	 [ 00 00 PREFIX ] [ length body ] [ length iv ] [ iv ] [ 00 SUFFIX ] [ ciphertext ] [ tag ]
	*/
	header := append(PREFIX, lenBody...)
	header = append(header, lenIV...)
	header = append(header, encryptedPayload.IV...)
	header = append(header, SUFFIX...)

	encodedPayload := append(header, encryptedPayload.Ciphertext...)
	encodedPayload = append(encodedPayload, encryptedPayload.Tag...)
	return encodedPayload
}
