package cryptography

import (
	"encoding/hex"
	"errors"
	"sync"
	"synnergy_network/pkg/common"

	"golang.org/x/crypto/sha3"
)

var shakeOpsLock sync.Mutex

// Shake128Hash: Generates a SHAKE128 hash with variable output length
func Shake128Hash(data []byte, outputLength int) (string, error) {
    if outputLength <= 0 {
        LogShakeOperation("Shake128Hash", "Invalid output length specified")
        return "", errors.New("output length must be greater than zero")
    }

    hasher := sha3.NewShake128()
    hasher.Write(data)
    hash := make([]byte, outputLength)
    hasher.Read(hash)

    LogShakeOperation("Shake128Hash", "SHAKE128 hash generated")
    return hex.EncodeToString(hash), nil
}

// Shake256Hash: Generates a SHAKE256 hash with variable output length
func Shake256Hash(data []byte, outputLength int) (string, error) {
    if outputLength <= 0 {
        LogShakeOperation("Shake256Hash", "Invalid output length specified")
        return "", errors.New("output length must be greater than zero")
    }

    hasher := sha3.NewShake256()
    hasher.Write(data)
    hash := make([]byte, outputLength)
    hasher.Read(hash)

    LogShakeOperation("Shake256Hash", "SHAKE256 hash generated")
    return hex.EncodeToString(hash), nil
}

// LogShakeOperation: Logs SHAKE operations with encryption
func LogShakeOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("ShakeOperation", encryptedMessage)
}
