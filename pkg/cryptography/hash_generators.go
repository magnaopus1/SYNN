package cryptography

import (
	"crypto/sha256"
	"errors"
	"sync"
)

var hashGeneratorLock sync.Mutex

// HashCombine: Combines two hashes using XOR and hashes the result with SHA-256
func HashCombine(hash1, hash2 []byte) ([]byte, error) {
    if len(hash1) != len(hash2) {
        LogHashOperation("HashCombine", "Hash length mismatch for combination")
        return nil, errors.New("hashes must be of equal length to combine")
    }
    
    combined := make([]byte, len(hash1))
    for i := range hash1 {
        combined[i] = hash1[i] ^ hash2[i]
    }
    result := sha256.Sum256(combined)
    LogHashOperation("HashCombine", "Hashes combined with XOR and SHA-256")
    return result[:], nil
}

// HashGenerateVector: Generates a vector of hashes from a seed
func HashGenerateVector(seed []byte, length int) ([][]byte, error) {
    if length <= 0 {
        LogHashOperation("HashGenerateVector", "Invalid vector length")
        return nil, errors.New("vector length must be greater than zero")
    }
    
    vector := make([][]byte, length)
    currentHash := sha256.Sum256(seed)
    vector[0] = currentHash[:]
    
    for i := 1; i < length; i++ {
        currentHash = sha256.Sum256(currentHash[:])
        vector[i] = currentHash[:]
    }
    
    LogHashOperation("HashGenerateVector", "Hash vector generated successfully")
    return vector, nil
}

// TruncateHash: Truncates a hash to a specified length (in bytes)
func TruncateHash(hash []byte, length int) ([]byte, error) {
    if length <= 0 || length > len(hash) {
        LogHashOperation("TruncateHash", "Invalid truncation length")
        return nil, errors.New("invalid truncation length")
    }
    truncatedHash := hash[:length]
    LogHashOperation("TruncateHash", "Hash truncated successfully")
    return truncatedHash, nil
}

// BlockchainHashRoot: Calculates a hash root for a blockchain segment
func BlockchainHashRoot(hashes [][]byte) ([]byte, error) {
    if len(hashes) == 0 {
        LogHashOperation("BlockchainHashRoot", "Empty hash list for root calculation")
        return nil, errors.New("empty hash list for blockchain root calculation")
    }
    
    for len(hashes) > 1 {
        var newLevel [][]byte
        for i := 0; i < len(hashes); i += 2 {
            if i+1 < len(hashes) {
                combined, _ := HashCombine(hashes[i], hashes[i+1])
                newLevel = append(newLevel, combined)
            } else {
                newLevel = append(newLevel, hashes[i])
            }
        }
        hashes = newLevel
    }
    
    LogHashOperation("BlockchainHashRoot", "Blockchain hash root calculated")
    return hashes[0], nil
}

// NodeIdentityHash: Generates a unique node identity hash
func NodeIdentityHash(nodeID string, timestamp int64) ([]byte, error) {
    data := []byte(nodeID + string(timestamp))
    hash := sha256.Sum256(data)
    LogHashOperation("NodeIdentityHash", "Node identity hash generated")
    return hash[:], nil
}


