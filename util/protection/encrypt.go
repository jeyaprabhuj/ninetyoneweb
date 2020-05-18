package protection

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	var block cipher.Block
	block, _ = aes.NewCipher(key)

	var aead cipher.AEAD
	aead, _ = cipher.NewGCM(block)

	var nonce []byte
	nonce = make([]byte, aead.NonceSize())

	io.ReadFull(rand.Reader, nonce)

	return aead.Seal(nonce, nonce, data, nil), nil
}
