package daikin

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func encryptAESCFB(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid IV size: must be %d bytes", aes.BlockSize)
	}
	var (
		mode       = cipher.NewCFBEncrypter(block, iv)
		ciphertext = make([]byte, len(plaintext))
	)
	mode.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func decryptAESCFB(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid IV size: must be %d bytes", aes.BlockSize)
	}
	var (
		mode      = cipher.NewCFBDecrypter(block, iv)
		plaintext = make([]byte, len(ciphertext))
	)
	mode.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
