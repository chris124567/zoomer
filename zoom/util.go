package zoom

import (
	"encoding/base64"
	"strings"
)

func ZoomEscapedBase64Decode(encoded string) ([]byte, error) {
	escaped := strings.ReplaceAll(encoded, "_", "/")
	escaped = strings.ReplaceAll(escaped, "-", "+")
	return base64.RawStdEncoding.DecodeString(escaped)
}
