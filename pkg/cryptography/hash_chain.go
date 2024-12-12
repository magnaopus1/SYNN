package cryptography

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"golang.org/x/crypto/ripemd160"
)

var hashChainLock sync.Mutex

// HashChainGenerate: Generates a hash chain starting with a seed
func HashChainGenerate(seed []byte, length int) ([][]byte, error) {
    hashChain := make([][]byte, length)
    currentHash := sha256.Sum256(seed)
    hashChain[0] = currentHash[:]
    
    for i := 1; i < length; i++ {
        currentHash = sha256.Sum256(currentHash[:])
        hashChain[i] = currentHash[:]
    }
    
    LogHashOperation("HashChainGenerate", fmt.Sprintf("Generated hash chain of length %d", length))
    return hashChain, nil
}

// HashChainVerify: Verifies a hash chain by recalculating and matching hashes
func HashChainVerify(seed []byte, hashChain [][]byte) (bool, error) {
    currentHash := sha256.Sum256(seed)
    if !equal(currentHash[:], hashChain[0]) {
        LogHashOperation("HashChainVerify", "Initial seed mismatch")
        return false, errors.New("initial seed mismatch")
    }
    
    for i := 1; i < len(hashChain); i++ {
        currentHash = sha256.Sum256(currentHash[:])
        if !equal(currentHash[:], hashChain[i]) {
            LogHashOperation("HashChainVerify", "Hash chain verification failed")
            return false, errors.New("hash chain verification failed")
        }
    }
    
    LogHashOperation("HashChainVerify", "Hash chain verified successfully")
    return true, nil
}

// HashConcatenate: Concatenates two hashes and returns the SHA-256 hash of the result
func HashConcatenate(hash1, hash2 []byte) ([]byte, error) {
    concatenated := append(hash1, hash2...)
    result := sha256.Sum256(concatenated)
    LogHashOperation("HashConcatenate", "Concatenated two hashes")
    return result[:], nil
}

// HashTreeRoot: Calculates the root hash of a Merkle tree from a set of hashes
func HashTreeRoot(hashes [][]byte) ([]byte, error) {
    if len(hashes) == 0 {
        LogHashOperation("HashTreeRoot", "Empty hash list")
        return nil, errors.New("empty hash list")
    }
    
    for len(hashes) > 1 {
        var newLevel [][]byte
        for i := 0; i < len(hashes); i += 2 {
            if i+1 < len(hashes) {
                concatenated, _ := HashConcatenate(hashes[i], hashes[i+1])
                newLevel = append(newLevel, concatenated)
            } else {
                newLevel = append(newLevel, hashes[i])
            }
        }
        hashes = newLevel
    }
    
    LogHashOperation("HashTreeRoot", "Calculated Merkle tree root hash")
    return hashes[0], nil
}

// HashRandomize: Generates a randomized hash from input and a random nonce
func HashRandomize(input []byte) ([]byte, error) {
    nonce := make([]byte, 16)
    _, err := rand.Read(nonce)
    if err != nil {
        return nil, err
    }
    
    combined := append(input, nonce...)
    result := sha256.Sum256(combined)
    LogHashOperation("HashRandomize", "Randomized hash generated")
    return result[:], nil
}

// BinaryHashReduction: Performs binary reduction on a list of hashes
func BinaryHashReduction(hashes [][]byte) ([]byte, error) {
    if len(hashes) == 0 {
        LogHashOperation("BinaryHashReduction", "Empty hash list for reduction")
        return nil, errors.New("empty hash list")
    }
    
    reducedHash := hashes[0]
    for i := 1; i < len(hashes); i++ {
        reducedHash, _ = HashConcatenate(reducedHash, hashes[i])
    }
    
    LogHashOperation("BinaryHashReduction", "Binary hash reduction completed")
    return reducedHash, nil
}

// CompositeHashCalculate: Calculates a composite hash from multiple inputs
func CompositeHashCalculate(inputs [][]byte) ([]byte, error) {
    composite := make([]byte, 0)
    for _, input := range inputs {
        composite = append(composite, input...)
    }
    
    result := sha256.Sum256(composite)
    LogHashOperation("CompositeHashCalculate", "Composite hash calculated")
    return result[:], nil
}


// TokenizeHash: Tokenizes data by hashing and encoding to hexadecimal
func TokenizeHash(data []byte) (string, error) {
    hash := sha256.Sum256(data)
    token := hex.EncodeToString(hash[:])
    LogHashOperation("TokenizeHash", "Data tokenized")
    return token, nil
}

// RIPEMD128Hash: Hashes data using RIPEMD-128
func RIPEMD128Hash(data []byte) ([]byte, error) {
    hasher := ripemd160.New()
    hasher.Write(data)
    hash := hasher.Sum(nil)[:16] // Use the first 16 bytes for RIPEMD-128 equivalent
    LogHashOperation("RIPEMD128Hash", "Data hashed with RIPEMD-128")
    return hash, nil
}

