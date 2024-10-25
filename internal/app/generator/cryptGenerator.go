// Package generator содержит набор инструментов для генерации и шифрования данных.
package generator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// CryptGenerator Структура выполняющая шифрование строк с использованием ключа шифрования.
type CryptGenerator struct {
	// key - секретное слово шифрования и дешифрования.
	key []byte
}

// NewCryptGenerator Конструктор шифрователя key - ключ шифрования.
func NewCryptGenerator(key string) *CryptGenerator {
	return &CryptGenerator{key: []byte(key)}
}

// EncodeValue производит декодирование строки, value - принимает защифровонную строку.
func (c *CryptGenerator) EncodeValue(value string) (string, error) {
	byteMsg := []byte(value)
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("could not encrypt: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecodeValue - производит декодирование строки переднноый в value.
func (c *CryptGenerator) DecodeValue(value string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("could not base64 decode: %v", err)
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("could not create new cipher: %v", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
