package cryptography

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"sync"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

var hashOperationsLock sync.Mutex

// BLAKE2BHash: Hashes data using BLAKE2b-256
func BLAKE2BHash(data []byte) ([]byte, error) {
    hash := blake2b.Sum256(data)
    LogHashOperation("BLAKE2BHash", "Data hashed with BLAKE2b-256")
    return hash[:], nil
}


// HMACSHA256: Generates HMAC-SHA256 for the given data and key
func HMACSHA256(key, data []byte) ([]byte, error) {
    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    result := mac.Sum(nil)
    LogHashOperation("HMACSHA256", "HMAC-SHA256 generated")
    return result, nil
}

// HMACSHA512: Generates HMAC-SHA512 for the given data and key
func HMACSHA512(key, data []byte) ([]byte, error) {
    mac := hmac.New(sha512.New, key)
    mac.Write(data)
    result := mac.Sum(nil)
    LogHashOperation("HMACSHA512", "HMAC-SHA512 generated")
    return result, nil
}

// HMACBlake2b: Generates HMAC-BLAKE2b for the given data and key
func HMACBlake2b(key, data []byte) ([]byte, error) {
    mac := hmac.New(blake2b.New256, key)
    mac.Write(data)
    result := mac.Sum(nil)
    LogHashOperation("HMACBlake2b", "HMAC-BLAKE2b generated")
    return result, nil
}


// MD5Hash: Hashes data using MD5
func MD5Hash(data []byte) ([]byte, error) {
    hash := md5.Sum(data)
    LogHashOperation("MD5Hash", "Data hashed with MD5")
    return hash[:], nil
}



// SHAKE128Hash: Hashes data using SHAKE128 (variable output)
func SHAKE128Hash(data []byte, length int) ([]byte, error) {
    hash := make([]byte, length)
    shake := sha3.NewShake128()
    shake.Write(data)
    shake.Read(hash)
    LogHashOperation("SHAKE128Hash", "Data hashed with SHAKE128")
    return hash, nil
}

// SHAKE256Hash: Hashes data using SHAKE256 (variable output)
func SHAKE256Hash(data []byte, length int) ([]byte, error) {
    hash := make([]byte, length)
    shake := sha3.NewShake256()
    shake.Write(data)
    shake.Read(hash)
    LogHashOperation("SHAKE256Hash", "Data hashed with SHAKE256")
    return hash, nil
}

// HashSHA256: Hashes data using SHA-256
func HashSHA256(data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    LogHashOperation("HashSHA256", "Data hashed with SHA-256")
    return hash[:], nil
}

// HashSHA512: Hashes data using SHA-512
func HashSHA512(data []byte) ([]byte, error) {
    hash := sha512.Sum512(data)
    LogHashOperation("HashSHA512", "Data hashed with SHA-512")
    return hash[:], nil
}

// HashRIPEMD160: Hashes data using RIPEMD-160
func HashRIPEMD160(data []byte) ([]byte, error) {
    hasher := ripemd160.New()
    hasher.Write(data)
    hash := hasher.Sum(nil)
    LogHashOperation("HashRIPEMD160", "Data hashed with RIPEMD-160")
    return hash, nil
}

// HashSHA1: Hashes data using SHA-1
func HashSHA1(data []byte) ([]byte, error) {
    hasher := sha1.New()
    hasher.Write(data)
    hash := hasher.Sum(nil)
    LogHashOperation("HashSHA1", "Data hashed with SHA-1")
    return hash, nil
}

// HashMurmur: Hypothetical function for hashing using MurmurHash (placeholder for integration)
func HashMurmur(data []byte) ([]byte, error) {
    // Assuming a hypothetical MurmurHash function integration
    hash := murmur.Sum128(data) // Hypothetical library for MurmurHash
    LogHashOperation("HashMurmur", "Data hashed with MurmurHash")
    return hash, nil
}


