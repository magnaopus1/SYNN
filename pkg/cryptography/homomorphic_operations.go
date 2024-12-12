package cryptography

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
)

// Lock for concurrent access protection
var homomorphicOpsLock sync.Mutex

// HomomorphicEncrypt: Encrypts data using a homomorphic encryption scheme (e.g., Paillier)
func HomomorphicEncrypt(publicKey *big.Int, plaintext *big.Int) (*big.Int, error) {
    if publicKey == nil || plaintext == nil {
        LogHomomorphicOperation("HomomorphicEncrypt", "Invalid public key or plaintext provided")
        return nil, errors.New("public key and plaintext must not be nil")
    }

    // Example: Paillier encryption steps (simplified for illustration)
    nSquared := new(big.Int).Mul(publicKey, publicKey)
    r, err := randomInRange(new(big.Int).SetInt64(1), publicKey)
    if err != nil {
        LogHomomorphicOperation("HomomorphicEncrypt", "Random number generation failed")
        return nil, err
    }

    g := new(big.Int).Add(publicKey, big.NewInt(1))
    c1 := new(big.Int).Exp(g, plaintext, nSquared)    // g^plaintext mod n^2
    c2 := new(big.Int).Exp(r, publicKey, nSquared)    // r^n mod n^2

    ciphertext := new(big.Int).Mul(c1, c2)            // g^plaintext * r^n mod n^2
    ciphertext.Mod(ciphertext, nSquared)

    LogHomomorphicOperation("HomomorphicEncrypt", "Data encrypted with homomorphic scheme")
    return ciphertext, nil
}

// HomomorphicDecrypt: Decrypts data encrypted with a homomorphic encryption scheme (e.g., Paillier)
func HomomorphicDecrypt(privateKey *big.Int, publicKey *big.Int, ciphertext *big.Int) (*big.Int, error) {
    if privateKey == nil || publicKey == nil || ciphertext == nil {
        LogHomomorphicOperation("HomomorphicDecrypt", "Invalid private key, public key, or ciphertext provided")
        return nil, errors.New("keys and ciphertext must not be nil")
    }

    // Example: Paillier decryption steps (simplified for illustration)
    nSquared := new(big.Int).Mul(publicKey, publicKey)
    x := new(big.Int).Exp(ciphertext, privateKey, nSquared)
    x.Sub(x, big.NewInt(1))
    x.Div(x, publicKey)

    plaintext := new(big.Int).Mod(x, publicKey)       // Final result mod publicKey

    LogHomomorphicOperation("HomomorphicDecrypt", "Data decrypted with homomorphic scheme")
    return plaintext, nil
}

// Helper function to generate a random number in the range [min, max)
func randomInRange(min, max *big.Int) (*big.Int, error) {
    rangeSize := new(big.Int).Sub(max, min)
    r, err := rand.Int(rand.Reader, rangeSize)
    if err != nil {
        return nil, err
    }
    return r.Add(r, min), nil
}

// LogHomomorphicOperation: Logs homomorphic encryption and decryption operations with encryption
func LogHomomorphicOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return ledger.LogEvent("HomomorphicOperation", encryptedMessage)
}
