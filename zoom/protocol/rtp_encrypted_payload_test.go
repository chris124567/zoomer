package protocol

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestUnmarshalRtpEncryptedPayload(t *testing.T) {
	encryptedPayload, err := hex.DecodeString("000000860c18000000000000000000000000916b4946439355dd851c3e8fcb0c9ca360e2fa595aaee6917c22177d66bd11863b6c830eed9a767e15193adba82a0e2295053ad6196daffd70377ca6b9be9301472d2ae9ea0320fc69d73d4a0e486db9665945bb5c09fd8db175a01b6bc000bec59a5bfdfece8a4363a40ef06092061bea0762e8ed19b5118837d599bf503f928c621500c4c6c81b9e57f31b768124eecb85922d0331")
	if err != nil {
		t.Error(err)
		return
	}

	zoomEncryptedPayload := &RtpEncryptedPayload{}
	zoomEncryptedPayload.Unmarshal(encryptedPayload)

	expectedIv, err := hex.DecodeString("180000000000000000000000")
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(zoomEncryptedPayload.IV, expectedIv) != 0 {
		t.Error("IV did not match expected")
	}

	expectedCiphertext, err := hex.DecodeString("916b4946439355dd851c3e8fcb0c9ca360e2fa595aaee6917c22177d66bd11863b6c830eed9a767e15193adba82a0e2295053ad6196daffd70377ca6b9be9301472d2ae9ea0320fc69d73d4a0e486db9665945bb5c09fd8db175a01b6bc000bec59a5bfdfece8a4363a40ef06092061bea0762e8ed19b5118837d599bf503f928c621500c4c6")
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(zoomEncryptedPayload.Ciphertext, expectedCiphertext) != 0 {
		t.Error("Tag did not match expected")
	}

	expectedTag, err := hex.DecodeString("c81b9e57f31b768124eecb85922d0331")
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(zoomEncryptedPayload.Tag, expectedTag) != 0 {
		t.Error("Tag did not match expected")
	}
}

func TestMarshalRtpEncryptedPayload(t *testing.T) {
	IV, err := hex.DecodeString("180000000000000000000000")
	if err != nil {
		t.Error(err)
		return
	}

	ciphertextWithTag, err := hex.DecodeString("916b4946439355dd851c3e8fcb0c9ca360e2fa595aaee6917c22177d66bd11863b6c830eed9a767e15193adba82a0e2295053ad6196daffd70377ca6b9be9301472d2ae9ea0320fc69d73d4a0e486db9665945bb5c09fd8db175a01b6bc000bec59a5bfdfece8a4363a40ef06092061bea0762e8ed19b5118837d599bf503f928c621500c4c6c81b9e57f31b768124eecb85922d0331")
	if err != nil {
		t.Error(err)
		return
	}

	zoomEncryptedPayload := NewRtpEncryptedPayload(IV, ciphertextWithTag)
	encryptedPayload := zoomEncryptedPayload.Marshal()

	expectedEncryptedPayload, err := hex.DecodeString("000000860c18000000000000000000000000916b4946439355dd851c3e8fcb0c9ca360e2fa595aaee6917c22177d66bd11863b6c830eed9a767e15193adba82a0e2295053ad6196daffd70377ca6b9be9301472d2ae9ea0320fc69d73d4a0e486db9665945bb5c09fd8db175a01b6bc000bec59a5bfdfece8a4363a40ef06092061bea0762e8ed19b5118837d599bf503f928c621500c4c6c81b9e57f31b768124eecb85922d0331")
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(encryptedPayload, expectedEncryptedPayload) != 0 {
		t.Error("Plaintext did not match expected")
	}
}
