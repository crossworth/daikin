package daikin

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func decryptAESCFB(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	var (
		mode      = cipher.NewCFBDecrypter(block, iv)
		plaintext = make([]byte, len(ciphertext))
	)
	mode.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
