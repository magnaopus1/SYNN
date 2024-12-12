package cryptography

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"synnergy_network/pkg/common"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

var keyDerivationLock sync.Mutex

// PBKDF2KeyDerivation: Derives a key using PBKDF2
func PBKDF2KeyDerivation(password, salt []byte, iter, keyLen int) ([]byte, error) {
    key := pbkdf2.Key(password, salt, iter, keyLen, sha512.New)
    LogKeyDerivationOperation("PBKDF2KeyDerivation", "Key derived using PBKDF2")
    return key, nil
}

// HKDFExtract: Extracts a pseudo-random key using HKDF
func HKDFExtract(salt, inputKeyMaterial []byte) ([]byte, error) {
    hkdf := hkdf.New(sha256.New, inputKeyMaterial, salt, nil)
    prk := make([]byte, 32) // Adjust length as required
    if _, err := hkdf.Read(prk); err != nil {
        LogKeyDerivationOperation("HKDFExtract", "HKDF extraction failed")
        return nil, err
    }
    LogKeyDerivationOperation("HKDFExtract", "Pseudo-random key extracted with HKDF")
    return prk, nil
}

// HKDFExpand: Expands a pseudo-random key into multiple keys using HKDF
func HKDFExpand(prk, info []byte, length int) ([]byte, error) {
    hkdf := hkdf.New(sha256.New, prk, nil, info)
    expandedKey := make([]byte, length)
    if _, err := hkdf.Read(expandedKey); err != nil {
        LogKeyDerivationOperation("HKDFExpand", "HKDF expansion failed")
        return nil, err
    }
    LogKeyDerivationOperation("HKDFExpand", "Key expanded with HKDF")
    return expandedKey, nil
}

// ScryptKeyDerivation: Derives a key using Scrypt
func ScryptKeyDerivation(password, salt []byte, N, r, p, keyLen int) ([]byte, error) {
    key, err := scrypt.Key(password, salt, N, r, p, keyLen)
    if err != nil {
        LogKeyDerivationOperation("ScryptKeyDerivation", "Scrypt key derivation failed")
        return nil, err
    }
    LogKeyDerivationOperation("ScryptKeyDerivation", "Key derived using Scrypt")
    return key, nil
}

// GenerateIV: Generates an initialization vector for encryption
func GenerateIV() ([]byte, error) {
    iv := make([]byte, 12)
    if _, err := rand.Read(iv); err != nil {
        return nil, err
    }
    LogKeyDerivationOperation("GenerateIV", "Initialization vector generated")
    return iv, nil
}

// SaltGenerate: Generates a random salt
func SaltGenerate(length int) ([]byte, error) {
    salt := make([]byte, length)
    if _, err := rand.Read(salt); err != nil {
        LogKeyDerivationOperation("SaltGenerate", "Salt generation failed")
        return nil, err
    }
    LogKeyDerivationOperation("SaltGenerate", "Salt generated")
    return salt, nil
}

// EllipticCurveMultiply: Multiplies a point on an elliptic curve by a scalar
func EllipticCurveMultiply(pointX, pointY, scalar *big.Int) (*big.Int, *big.Int, error) {
    // Placeholder for actual elliptic curve multiplication logic
    // Example for demonstration purposes only
    LogKeyDerivationOperation("EllipticCurveMultiply", "Point multiplied on elliptic curve")
    return pointX, pointY, nil
}

// EllipticCurveAdd: Adds two points on an elliptic curve
func EllipticCurveAdd(point1X, point1Y, point2X, point2Y *big.Int) (*big.Int, *big.Int, error) {
    // Placeholder for actual elliptic curve addition logic
    LogKeyDerivationOperation("EllipticCurveAdd", "Points added on elliptic curve")
    return point1X, point1Y, nil
}

// ChecksumCalculate: Calculates a checksum for data
func ChecksumCalculate(data []byte) string {
    checksum := sha256.Sum256(data)
    LogKeyDerivationOperation("ChecksumCalculate", "Checksum calculated")
    return hex.EncodeToString(checksum[:])
}

// EncryptWithSalt: Encrypts data with salt using XSalsa20
func EncryptWithSalt(key, salt, plaintext []byte) ([]byte, error) {
    var nonce [24]byte
    copy(nonce[:], salt[:24])
    encrypted := secretbox.Seal(nil, plaintext, &nonce, (*[32]byte)(key))
    LogKeyDerivationOperation("EncryptWithSalt", "Data encrypted with salt")
    return encrypted, nil
}

// DecryptWithSalt: Decrypts salt-encrypted data using XSalsa20
func DecryptWithSalt(key, salt, ciphertext []byte) ([]byte, error) {
    var nonce [24]byte
    copy(nonce[:], salt[:24])
    decrypted, ok := secretbox.Open(nil, ciphertext, &nonce, (*[32]byte)(key))
    if !ok {
        LogKeyDerivationOperation("DecryptWithSalt", "Salt decryption failed")
        return nil, errors.New("decryption failed")
    }
    LogKeyDerivationOperation("DecryptWithSalt", "Data decrypted with salt")
    return decrypted, nil
}

// OTPGenerate: Generates a time-based one-time password (TOTP)
func OTPGenerate(secret []byte, timestamp int64) string {
    h := hmac.New(sha1.New, secret)
    timeCounter := make([]byte, 8)
    for i := 7; i >= 0; i-- {
        timeCounter[i] = byte(timestamp)
        timestamp >>= 8
    }
    h.Write(timeCounter)
    hash := h.Sum(nil)
    offset := hash[len(hash)-1] & 0xf
    code := int(hash[offset]&0x7f)<<24 |
        int(hash[offset+1])<<16 |
        int(hash[offset+2])<<8 |
        int(hash[offset+3])
    code %= 1000000
    LogKeyDerivationOperation("OTPGenerate", "One-time password generated")
    return fmt.Sprintf("%06d", code)
}

// OTPVerify: Verifies a one-time password
func OTPVerify(secret []byte, otp string, timestamp int64) bool {
    expectedOtp := OTPGenerate(secret, timestamp)
    isValid := expectedOtp == otp
    LogKeyDerivationOperation("OTPVerify", "One-time password verification")
    return isValid
}

// LogKeyDerivationOperation: Logs key derivation and related operations with encryption
func LogKeyDerivationOperation(operation, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("KeyDerivationOperation", encryptedMessage)
}
