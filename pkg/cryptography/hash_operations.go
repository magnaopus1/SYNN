package cryptography

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

var hashOpsLock sync.Mutex

// SHA256Hash: Hashes data using SHA-256
func SHA256Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    if err := LogHashOperation(ledger, "SHA256Hash", "Data hashed with SHA-256"); err != nil {
        return nil, err
    }
    return hash[:], nil
}

// SHA512Hash: Hashes data using SHA-512
func SHA512Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hash := sha512.Sum512(data)
    if err := LogHashOperation(ledger, "SHA512Hash", "Data hashed with SHA-512"); err != nil {
        return nil, err
    }
    return hash[:], nil
}

// Blake2Hash: Hashes data using BLAKE2b-256
func Blake2Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hash := blake2b.Sum256(data)
    if err := LogHashOperation(ledger, "Blake2Hash", "Data hashed with BLAKE2b-256"); err != nil {
        return nil, err
    }
    return hash[:], nil
}

// SHA3256Hash: Hashes data using SHA3-256
func SHA3256Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hash := sha3.Sum256(data)
    if err := LogHashOperation(ledger, "SHA3256Hash", "Data hashed with SHA3-256"); err != nil {
        return nil, err
    }
    return hash[:], nil
}

// SHA3512Hash: Hashes data using SHA3-512
func SHA3512Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hash := sha3.Sum512(data)
    if err := LogHashOperation(ledger, "SHA3512Hash", "Data hashed with SHA3-512"); err != nil {
        return nil, err
    }
    return hash[:], nil
}

// RIPEMD160Hash: Hashes data using RIPEMD-160
func RIPEMD160Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hasher := ripemd160.New()
    hasher.Write(data)
    hash := hasher.Sum(nil)
    if err := LogHashOperation(ledger, "RIPEMD160Hash", "Data hashed with RIPEMD-160"); err != nil {
        return nil, err
    }
    return hash, nil
}

// SHA1Hash: Hashes data using SHA-1
func SHA1Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hasher := sha1.New()
    hasher.Write(data)
    hash := hasher.Sum(nil)
    if err := LogHashOperation(ledger, "SHA1Hash", "Data hashed with SHA-1"); err != nil {
        return nil, err
    }
    return hash, nil
}

// CRC32Hash: Calculates the CRC32 checksum of data
func CRC32Hash(ledger *ledger.Ledger, data []byte) uint32 {
    checksum := crc32.ChecksumIEEE(data)
    _ = LogHashOperation(ledger, "CRC32Hash", "Data hashed with CRC32") // Ignoring error for checksum
    return checksum
}

// FNVHash: Hashes data using FNV-1a
func FNVHash(ledger *ledger.Ledger, data []byte) (uint64, error) {
    hasher := fnv.New64a()
    hasher.Write(data)
    hash := hasher.Sum64()
    if err := LogHashOperation(ledger, "FNVHash", "Data hashed with FNV-1a"); err != nil {
        return 0, err
    }
    return hash, nil
}

// KECCAK256Hash: Hashes data using KECCAK-256 (common in Ethereum)
func KECCAK256Hash(ledger *ledger.Ledger, data []byte) ([]byte, error) {
    hash := sha3.NewLegacyKeccak256()
    hash.Write(data)
    result := hash.Sum(nil)
    if err := LogHashOperation(ledger, "KECCAK256Hash", "Data hashed with KECCAK-256"); err != nil {
        return nil, err
    }
    return result, nil
}

// XORCipher: Encrypts/Decrypts data using XOR cipher with a given key
func XORCipher(ledger *ledger.Ledger, data, key []byte) ([]byte, error) {
    if len(key) == 0 {
        return nil, errors.New("key length must be greater than zero")
    }
    result := make([]byte, len(data))
    for i := range data {
        result[i] = data[i] ^ key[i%len(key)]
    }
    if err := LogEncryptionOperation(ledger, "XORCipher", "Data XOR encrypted"); err != nil {
        return nil, err
    }
    return result, nil
}

// ChaCha20Encrypt encrypts data using ChaCha20 with a nonce.
func ChaCha20Encrypt(ledger *ledger.Ledger, key, nonce, plaintext []byte) ([]byte, error) {
    if len(key) != chacha20.KeySize || len(nonce) != chacha20.NonceSize {
        return nil, fmt.Errorf("invalid key or nonce size for ChaCha20")
    }

    ciphertext := make([]byte, len(plaintext))
    stream, err := chacha20.NewUnauthenticatedCipher(key, nonce)
    if err != nil {
        LogEncryptionOperation(ledger, "ChaCha20Encrypt", "ChaCha20 encryption failed")
        return nil, fmt.Errorf("chacha20 encryption failed: %w", err)
    }

    stream.XORKeyStream(ciphertext, plaintext)
    LogEncryptionOperation(ledger, "ChaCha20Encrypt", "Data encrypted with ChaCha20")
    return ciphertext, nil
}

// ChaCha20Decrypt decrypts data encrypted with ChaCha20.
func ChaCha20Decrypt(ledger *ledger.Ledger, key, nonce, ciphertext []byte) ([]byte, error) {
    if len(key) != chacha20.KeySize || len(nonce) != chacha20.NonceSize {
        return nil, fmt.Errorf("invalid key or nonce size for ChaCha20")
    }

    plaintext := make([]byte, len(ciphertext))
    stream, err := chacha20.NewUnauthenticatedCipher(key, nonce)
    if err != nil {
        LogEncryptionOperation(ledger, "ChaCha20Decrypt", "ChaCha20 decryption failed")
        return nil, fmt.Errorf("chacha20 decryption failed: %w", err)
    }

    stream.XORKeyStream(plaintext, ciphertext)
    LogEncryptionOperation(ledger, "ChaCha20Decrypt", "Data decrypted with ChaCha20")
    return plaintext, nil
}

// XSalsa20Encrypt encrypts data using XSalsa20 with a nonce.
func XSalsa20Encrypt(ledger *ledger.Ledger, key, nonce, plaintext []byte) ([]byte, error) {
    if len(key) != 32 || len(nonce) != 24 {
        return nil, fmt.Errorf("invalid key or nonce length for XSalsa20")
    }

    var keyArr [32]byte
    var nonceArr [24]byte
    copy(keyArr[:], key)
    copy(nonceArr[:], nonce)

    encrypted := secretbox.Seal(nil, plaintext, &nonceArr, &keyArr)
    LogEncryptionOperation(ledger, "XSalsa20Encrypt", "Data encrypted with XSalsa20")
    return encrypted, nil
}

// XSalsa20Decrypt decrypts XSalsa20 encrypted data.
func XSalsa20Decrypt(ledger *ledger.Ledger, key, nonce, ciphertext []byte) ([]byte, error) {
    if len(key) != 32 || len(nonce) != 24 {
        return nil, fmt.Errorf("invalid key or nonce length for XSalsa20")
    }

    var keyArr [32]byte
    var nonceArr [24]byte
    copy(keyArr[:], key)
    copy(nonceArr[:], nonce)

    decrypted, ok := secretbox.Open(nil, ciphertext, &nonceArr, &keyArr)
    if !ok {
        LogEncryptionOperation(ledger, "XSalsa20Decrypt", "XSalsa20 decryption failed")
        return nil, fmt.Errorf("decryption failed: secretbox open returned false")
    }

    LogEncryptionOperation(ledger, "XSalsa20Decrypt", "Data decrypted with XSalsa20")
    return decrypted, nil
}

// BlowfishEncrypt encrypts data using Blowfish in CBC mode.
func BlowfishEncrypt(ledger *ledger.Ledger, key, plaintext []byte) ([]byte, error) {
    block, err := blowfish.NewCipher(key)
    if err != nil {
        LogEncryptionOperation(ledger, "BlowfishEncrypt", "Blowfish encryption failed")
        return nil, fmt.Errorf("failed to create Blowfish cipher: %w", err)
    }

    padding := block.BlockSize() - len(plaintext)%block.BlockSize()
    padText := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

    ciphertext := make([]byte, len(padText))
    iv := make([]byte, block.BlockSize())
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %w", err)
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext, padText)

    LogEncryptionOperation(ledger, "BlowfishEncrypt", "Data encrypted with Blowfish")
    return append(iv, ciphertext...), nil
}

// BlowfishDecrypt decrypts data encrypted with Blowfish in CBC mode.
func BlowfishDecrypt(ledger *ledger.Ledger, key, ciphertext []byte) ([]byte, error) {
    block, err := blowfish.NewCipher(key)
    if err != nil {
        LogEncryptionOperation(ledger, "BlowfishDecrypt", "Failed to create Blowfish cipher")
        return nil, fmt.Errorf("blowfish decryption failed: %w", err)
    }

    if len(ciphertext) < block.BlockSize() {
        LogEncryptionOperation(ledger, "BlowfishDecrypt", "Ciphertext too short")
        return nil, fmt.Errorf("ciphertext too short")
    }

    iv := ciphertext[:block.BlockSize()]
    ciphertext = ciphertext[block.BlockSize():]

    if len(ciphertext)%block.BlockSize() != 0 {
        LogEncryptionOperation(ledger, "BlowfishDecrypt", "Invalid ciphertext length")
        return nil, fmt.Errorf("invalid ciphertext length")
    }

    mode := cipher.NewCBCDecrypter(block, iv)
    decrypted := make([]byte, len(ciphertext))
    mode.CryptBlocks(decrypted, ciphertext)

    padding := int(decrypted[len(decrypted)-1])
    if padding < 1 || padding > block.BlockSize() {
        LogEncryptionOperation(ledger, "BlowfishDecrypt", "Invalid padding")
        return nil, fmt.Errorf("invalid padding")
    }

    decrypted = decrypted[:len(decrypted)-padding]
    LogEncryptionOperation(ledger, "BlowfishDecrypt", "Data decrypted with Blowfish")
    return decrypted, nil
}


// LogHashOperation logs hashing and encryption operations with encryption.
func LogHashOperation(ledger *ledger.Ledger, operation string, details string) error {
    encryption := &common.Encryption{} // Create an instance of the Encryption struct
    encryptedMessage, err := encryption.EncryptData("AES", []byte("Operation: "+operation+" - Details: "+details), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt hash operation log: %v", err)
    }
    return ledger.CryptographyLedger.LogHashEvent("HashOperation", encryptedMessage)
}

// LogEncryptionOperation logs encryption and decryption operations with encryption.
func LogEncryptionOperation(ledger *ledger.Ledger, operation string, details string) error {
    encryption := &common.Encryption{} // Create an instance of the Encryption struct
    encryptedMessage, err := encryption.EncryptData("AES", []byte("Operation: "+operation+" - Details: "+details), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt encryption operation log: %v", err)
    }
    return ledger.CryptographyLedger.LogEncryptionEvent("EncryptionOperation", encryptedMessage)
}
