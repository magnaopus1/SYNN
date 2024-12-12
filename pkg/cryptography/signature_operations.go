package cryptography

import (

	"crypto/ed25519"
	"time"
	"fmt"
	"encoding/hex"
)


// SignWithTimestamp creates a timestamped signature for verification at a specific time
func SignWithTimestamp(data []byte, privateKey ed25519.PrivateKey) (string, error) {
	timestamp := time.Now().Unix()
	dataWithTimestamp := fmt.Sprintf("%s|%d", data, timestamp)
	signature, err := SignData([]byte(dataWithTimestamp), privateKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

// VerifyTimestampSign verifies a signature that includes a timestamp
func VerifyTimestampSign(data []byte, timestamp int64, signature []byte, publicKey ed25519.PublicKey) (bool, error) {
	dataWithTimestamp := fmt.Sprintf("%s|%d", data, timestamp)
	return VerifySignature([]byte(dataWithTimestamp), signature, publicKey)
}
