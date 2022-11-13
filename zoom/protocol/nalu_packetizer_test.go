package protocol

import (
	"bytes"
	"log"
	"math/rand"
	"testing"
)

func TestNaluPacketizer(t *testing.T) {

	// Generate a random payload
	payload := make([]byte, 1800)
	rand.Read(payload)

	packetizer := NewNaluPacketizer()

	// Packetize the payload
	pkts, err := packetizer.Marshal(payload)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Depacketize the payload
	var decoded []byte
	for _, pkt := range pkts {
		log.Printf("%v", pkt[:32])
		decoded, err = packetizer.Unmarshal(pkt)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if bytes.Compare(payload, decoded) != 0 {
		log.Fatal("Decoding failed")
	}
}
