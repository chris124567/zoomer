package protocol

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestDeriveEncryptionKey(t *testing.T) {
	sharedKey, err := hex.DecodeString("5e832c8b269df39b8d2ec046d97e2b1119d83bf479fc3ad41f6f8ff5dd128ac6")
	if err != nil {
		t.Error(err)
		return
	}
	secretNonce, err := hex.DecodeString("33cd55ac374712fcf3f2b41bad4203b8")
	if err != nil {
		t.Error(err)
		return
	}

	streamKey := DeriveEncryptionKey(sharedKey, secretNonce)

	expectedStreamKey, err := hex.DecodeString("aa1a283465d5194fd4e57ac16ee3299f831d26b6745557db90ebcf0f9e3a9ce1")
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(streamKey, expectedStreamKey) != 0 {
		t.Error("Stream key did not match expected")
	}
}

func TestAesGcmDecrypt(t *testing.T) {
	sharedKey, err := hex.DecodeString("5e832c8b269df39b8d2ec046d97e2b1119d83bf479fc3ad41f6f8ff5dd128ac6")
	if err != nil {
		t.Error(err)
		return
	}
	secretNonce, err := hex.DecodeString("33cd55ac374712fcf3f2b41bad4203b8")
	if err != nil {
		t.Error(err)
		return
	}

	aesGcmCrypto, err := NewAesGcmCrypto(sharedKey, secretNonce)
	if err != nil {
		t.Error(err)
		return
	}

	pkt, err := hex.DecodeString("040000000000000000000000146ae168b1cb41172f40d2d157dd464c3810046831f281e80719cfdf67e4a35e5c8740af62f0ccf298b144db8089d0b2772946c0fe7088a64cadf71db033845588986db14e16343a8637b8f8e4ad6f6995f4c620c72ca3b0aeb1f9fa5ecceae54a7f02b0721998c47668f5700787a6b9d09a2895a4cf2ba0300d09a8f59a13e431ba3d953763f72d49685295089eb3bece20fe25f90abc3b468f90af472c67767792add95eb9bc43b03d3f2429b08b6873b7d2bc7be97c72018d")
	if err != nil {
		t.Error(err)
		return
	}

	nonce := pkt[:12]
	cipherText := pkt[12:]
	plainText, err := aesGcmCrypto.Decrypt(nonce, cipherText)
	if err != nil {
		t.Error(err)
		return
	}

	expectedPlaintext, err := hex.DecodeString("00000141cc0006000776da724a28e4c47f0080f13319e021b63b4e6421ccf3408a17067fa4d09220002a3faa3f455597fb2f9aa2fd0b99510e5eefb712c7e234ea2521c662e927f50df99be14640ea5d62570000000141007fb300018001ddb69c928a39311f0000030000030047e00000000141003decc0006000776da724a28e4c470000030000030047e000000001410016fb300018001ddb69c928a39311ff0000030000030047e0")
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(plainText, expectedPlaintext) != 0 {
		t.Error("Plaintext did not match expected")
	}
}

func TestAesGcmEncrypt(t *testing.T) {
	sharedKey, err := hex.DecodeString("5e832c8b269df39b8d2ec046d97e2b1119d83bf479fc3ad41f6f8ff5dd128ac6")
	if err != nil {
		t.Error(err)
		return
	}
	secretNonce, err := hex.DecodeString("33cd55ac374712fcf3f2b41bad4203b8")
	if err != nil {
		t.Error(err)
		return
	}

	aesGcmCrypto, err := NewAesGcmCrypto(sharedKey, secretNonce)
	if err != nil {
		t.Error(err)
		return
	}

	nonce, err := hex.DecodeString("040000000000000000000000")
	if err != nil {
		t.Error(err)
		return
	}

	plaintext, err := hex.DecodeString("00000141cc0006000776da724a28e4c47f0080f13319e021b63b4e6421ccf3408a17067fa4d09220002a3faa3f455597fb2f9aa2fd0b99510e5eefb712c7e234ea2521c662e927f50df99be14640ea5d62570000000141007fb300018001ddb69c928a39311f0000030000030047e00000000141003decc0006000776da724a28e4c470000030000030047e000000001410016fb300018001ddb69c928a39311ff0000030000030047e0")
	if err != nil {
		t.Error(err)
		return
	}

	ciphertext, err := aesGcmCrypto.Encrypt(nonce, plaintext)
	if err != nil {
		t.Error(err)
		return
	}

	expectedCiphertext, err := hex.DecodeString("146ae168b1cb41172f40d2d157dd464c3810046831f281e80719cfdf67e4a35e5c8740af62f0ccf298b144db8089d0b2772946c0fe7088a64cadf71db033845588986db14e16343a8637b8f8e4ad6f6995f4c620c72ca3b0aeb1f9fa5ecceae54a7f02b0721998c47668f5700787a6b9d09a2895a4cf2ba0300d09a8f59a13e431ba3d953763f72d49685295089eb3bece20fe25f90abc3b468f90af472c67767792add95eb9bc43b03d3f2429b08b6873b7d2bc7be97c72018d")
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(ciphertext, expectedCiphertext) != 0 {
		t.Error("Plaintext did not match expected")
	}
}
