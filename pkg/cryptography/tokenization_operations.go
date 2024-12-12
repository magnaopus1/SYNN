package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

// TokenizeData creates a unique, encrypted token for the given data
func TokenizeData(data []byte, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error creating cipher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %v", err)
	}

	token := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(token), nil
}

// DetokenizeData decrypts a token back into the original data
func DetokenizeData(token string, secretKey string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("error creating cipher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %v", err)
	}

	tokenData, err := hex.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("error decoding token: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(tokenData) < nonceSize {
		return nil, errors.New("token data is too short")
	}

	nonce, ciphertext := tokenData[:nonceSize], tokenData[nonceSize:]
	data, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %v", err)
	}

	return data, nil
}

// TokenizeWithExpiration creates a token with an expiration time
func TokenizeWithExpiration(data []byte, secretKey string, expiration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiration).Unix()
	dataWithExpiration := append(data, []byte(fmt.Sprintf("|%d", expirationTime))...)

	token, err := TokenizeData(dataWithExpiration, secretKey)
	if err != nil {
		return "", fmt.Errorf("error creating token with expiration: %v", err)
	}

	return token, nil
}

// ValidateTokenExpiration checks if the token is expired
func ValidateTokenExpiration(token string, secretKey string) (bool, error) {
	data, err := DetokenizeData(token, secretKey)
	if err != nil {
		return false, fmt.Errorf("error detokenizing data: %v", err)
	}

	// Separate expiration time from data
	dataParts := bytes.Split(data, []byte("|"))
	if len(dataParts) < 2 {
		return false, errors.New("invalid token format")
	}

	expirationTime, err := strconv.ParseInt(string(dataParts[1]), 10, 64)
	if err != nil {
		return false, fmt.Errorf("error parsing expiration time: %v", err)
	}

	if time.Now().Unix() > expirationTime {
		return false, errors.New("token has expired")
	}

	return true, nil
}
