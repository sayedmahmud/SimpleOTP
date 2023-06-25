package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
)

// Marshals the JSON, encrypts it, and returns a base64 encoded string
func Encrypt(e *Entry) string {
	plaintext, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key.Key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	encryptedData := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(encryptedData)
}

// Decrypts a base64 encoded and encrypted string, and unmarshals it
func Decrypt(entry string) (*Entry, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(entry)
	if err != nil {
		return nil, errors.New("failed to decode base64")
	}

	block, err := aes.NewCipher(key.Key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt")
	}

	var entryData Entry
	err = json.Unmarshal(plaintext, &entryData)
	if err != nil {
		panic(err)
	}

	return &entryData, nil
}

func hash(str string) string {
	hashedName := sha512.Sum512([]byte(str))
	// Hex
	return hex.EncodeToString(hashedName[:])
}