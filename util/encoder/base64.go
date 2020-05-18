package encoder

import (
	"encoding/base64"
)

func Base64EncodeToString(data []byte) string {
	base64Cipher := make([]byte, base64.RawStdEncoding.EncodedLen(len(data)))
	base64.RawStdEncoding.Encode(base64Cipher, data)
	return string(base64Cipher)
}

func Base64DecodeToByte(data string) ([]byte, error) {
	bData := []byte(data)
	cipher := make([]byte, base64.RawStdEncoding.DecodedLen(len(bData)))
	_, err := base64.RawStdEncoding.Decode(cipher, bData)
	if err != nil {
		return nil, err
	} else {
		return cipher, nil
	}
}

func Base64DecodeToString(data string) (string, error) {
	b, err := Base64DecodeToByte(data)
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
