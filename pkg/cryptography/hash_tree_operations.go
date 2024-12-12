package cryptography

import (
	"crypto/sha256"
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var hashTreeOpsLock sync.Mutex

// HashTreeGenerate: Generates a Merkle tree from a list of data hashes and returns the root hash
func HashTreeGenerate(data [][]byte) ([]byte, [][]byte, error) {
    if len(data) == 0 {
        LogHashTreeOperation("HashTreeGenerate", "Empty data list provided for Merkle tree generation")
        return nil, nil, errors.New("data list cannot be empty")
    }

    currentLevel := data
    var allLevels [][]byte

    for len(currentLevel) > 1 {
        var nextLevel [][]byte
        allLevels = append(allLevels, currentLevel...)

        for i := 0; i < len(currentLevel); i += 2 {
            if i+1 < len(currentLevel) {
                combinedHash := hashPair(currentLevel[i], currentLevel[i+1])
                nextLevel = append(nextLevel, combinedHash)
            } else {
                // If odd number of hashes, duplicate the last hash to ensure pairs
                combinedHash := hashPair(currentLevel[i], currentLevel[i])
                nextLevel = append(nextLevel, combinedHash)
            }
        }
        currentLevel = nextLevel
    }

    rootHash := currentLevel[0]
    LogHashTreeOperation("HashTreeGenerate", "Merkle tree generated with root hash")
    return rootHash, allLevels, nil
}

// HashTreeVerify: Verifies a leaf node in a Merkle tree using its Merkle proof and root hash
func HashTreeVerify(leaf []byte, proof [][]byte, rootHash []byte) (bool, error) {
    computedHash := leaf

    for _, siblingHash := range proof {
        computedHash = hashPair(computedHash, siblingHash)
    }

    if !equal(computedHash, rootHash) {
        LogHashTreeOperation("HashTreeVerify", "Merkle proof verification failed")
        return false, errors.New("merkle proof verification failed")
    }

    LogHashTreeOperation("HashTreeVerify", "Merkle proof verified successfully")
    return true, nil
}

// hashPair: Helper function to hash two byte slices together using SHA-256
func hashPair(left, right []byte) []byte {
    combined := append(left, right...)
    hash := sha256.Sum256(combined)
    return hash[:]
}

// equal: Helper function to check byte slice equality
func equal(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
    }
    }
    return true
}

// LogHashTreeOperation: Logs hash tree generation and verification operations with encryption
func LogHashTreeOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("HashTreeOperation", encryptedMessage)
}
