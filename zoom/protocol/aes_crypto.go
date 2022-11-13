package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
)

const (
	AES_GCM_TYPE       = 0x02
	AES_GCM_TAG_LENGTH = 16
)

type AesGcmCrypto struct {
	cipher cipher.AEAD
}

func NewAesGcmCrypto(sharedMeetingKey, secretNonce []byte) (*AesGcmCrypto, error) {
	derivedKey := DeriveEncryptionKey(sharedMeetingKey, secretNonce)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCMWithTagSize(block, AES_GCM_TAG_LENGTH)
	if err != nil {
		return nil, err
	}

	return &AesGcmCrypto{
		cipher: aesgcm,
	}, nil
}

func (c *AesGcmCrypto) Encrypt(nonce, plaintext []byte) ([]byte, error) {
	ciphertext := c.cipher.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nil
}

func (c *AesGcmCrypto) Decrypt(nonce, cipherText []byte) ([]byte, error) {
	plainText, err := c.cipher.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func DeriveEncryptionKey(sharedMeetingKey, secretNonce []byte) []byte {
	// TODO: assert key length & secret nonce length
	message := make([]byte, 0)
	message = append(message, secretNonce...)
	message = append(message, AES_GCM_TYPE)

	mac := hmac.New(sha256.New, sharedMeetingKey)
	mac.Write(message)
	return mac.Sum(nil)
}
