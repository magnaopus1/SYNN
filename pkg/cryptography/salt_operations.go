package cryptography

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var saltOpsLock sync.Mutex

// GenerateSalt: Generates a random salt of the specified length in bytes
func GenerateSalt(length int) ([]byte, error) {
    if length <= 0 {
        LogSaltOperation("GenerateSalt", "Invalid salt length specified")
        return nil, errors.New("salt length must be greater than zero")
    }

    salt := make([]byte, length)
    _, err := rand.Read(salt)
    if err != nil {
        LogSaltOperation("GenerateSalt", "Salt generation failed")
        return nil, err
    }

    LogSaltOperation("GenerateSalt", "Salt generated successfully")
    return salt, nil
}

// VerifySalt: Verifies that a given salt matches the expected salt
func VerifySalt(givenSalt, expectedSalt []byte) (bool, error) {
    if len(givenSalt) != len(expectedSalt) {
        LogSaltOperation("VerifySalt", "Salt length mismatch")
        return false, errors.New("salt length mismatch")
    }

    for i := range givenSalt {
        if givenSalt[i] != expectedSalt[i] {
            LogSaltOperation("VerifySalt", "Salt verification failed")
            return false, nil
        }
    }

    LogSaltOperation("VerifySalt", "Salt verified successfully")
    return true, nil
}

// SaltHashCombine: Combines salt with a hash of the data using SHA-256
func SaltHashCombine(data, salt []byte) (string, error) {
    hash := sha256.New()
    hash.Write(salt)
    hash.Write(data)
    combinedHash := hash.Sum(nil)

    LogSaltOperation("SaltHashCombine", "Data combined with salt and hashed")
    return hex.EncodeToString(combinedHash), nil
}

// LogSaltOperation: Logs salt-related operations with encryption
func LogSaltOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SaltOperation", encryptedMessage)
}
