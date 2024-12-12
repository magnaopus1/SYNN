package cryptography

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var digitalSignLock sync.Mutex

// SignData: Signs data using a specified private key
func SignData(privateKey []byte, data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    signature, err := blscrypto.Sign(privateKey, hash[:])
    if err != nil {
        LogSignatureOperation("SignData", "Data signing failed")
        return nil, errors.New("failed to sign data")
    }
    LogSignatureOperation("SignData", "Data signed successfully")
    return signature, nil
}

// VerifySignature: Verifies a signature for given data and public key
func VerifySignature(publicKey []byte, data []byte, signature []byte) (bool, error) {
    hash := sha256.Sum256(data)
    valid, err := blscrypto.Verify(publicKey, hash[:], signature)
    if err != nil || !valid {
        LogSignatureOperation("VerifySignature", "Signature verification failed")
        return false, errors.New("failed to verify signature")
    }
    LogSignatureOperation("VerifySignature", "Signature verified successfully")
    return true, nil
}

// ECDSASignature: Signs data using ECDSA and a specified private key
func ECDSASignature(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, []byte, error) {
    hash := sha256.Sum256(data)
    r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
    if err != nil {
        LogSignatureOperation("ECDSASignature", "ECDSA signing failed")
        return nil, nil, errors.New("failed to sign data using ECDSA")
    }
    LogSignatureOperation("ECDSASignature", "Data signed with ECDSA")
    return r.Bytes(), s.Bytes(), nil
}

// ECDSAVerify: Verifies an ECDSA signature for given data and public key
func ECDSAVerify(publicKey *ecdsa.PublicKey, data []byte, rBytes []byte, sBytes []byte) (bool, error) {
    hash := sha256.Sum256(data)
    r := new(big.Int).SetBytes(rBytes)
    s := new(big.Int).SetBytes(sBytes)
    valid := ecdsa.Verify(publicKey, hash[:], r, s)
    if !valid {
        LogSignatureOperation("ECDSAVerify", "ECDSA verification failed")
        return false, errors.New("failed to verify ECDSA signature")
    }
    LogSignatureOperation("ECDSAVerify", "ECDSA signature verified successfully")
    return true, nil
}

// BLSSign: Signs data using BLS and a specified private key
func BLSSign(privateKey []byte, data []byte) ([]byte, error) {
    signature, err := blscrypto.Sign(privateKey, data)
    if err != nil {
        LogSignatureOperation("BLSSign", "BLS signing failed")
        return nil, errors.New("failed to sign data using BLS")
    }
    LogSignatureOperation("BLSSign", "Data signed with BLS")
    return signature, nil
}

// BLSVerify: Verifies a BLS signature for given data and public key
func BLSVerify(publicKey []byte, data []byte, signature []byte) (bool, error) {
    valid, err := blscrypto.Verify(publicKey, data, signature)
    if err != nil || !valid {
        LogSignatureOperation("BLSVerify", "BLS verification failed")
        return false, errors.New("failed to verify BLS signature")
    }
    LogSignatureOperation("BLSVerify", "BLS signature verified successfully")
    return true, nil
}

// EDDSASign: Signs data using EDDSA and a specified private key
func EDDSASign(privateKey []byte, data []byte) ([]byte, error) {
    signature, err := eddsa.Sign(privateKey, data)
    if err != nil {
        LogSignatureOperation("EDDSASign", "EDDSA signing failed")
        return nil, errors.New("failed to sign data using EDDSA")
    }
    LogSignatureOperation("EDDSASign", "Data signed with EDDSA")
    return signature, nil
}

// EDDSAVerify: Verifies an EDDSA signature for given data and public key
func EDDSAVerify(publicKey []byte, data []byte, signature []byte) (bool, error) {
    valid, err := eddsa.Verify(publicKey, data, signature)
    if err != nil || !valid {
        LogSignatureOperation("EDDSAVerify", "EDDSA verification failed")
        return false, errors.New("failed to verify EDDSA signature")
    }
    LogSignatureOperation("EDDSAVerify", "EDDSA signature verified successfully")
    return true, nil
}

// LogSignatureOperation: Logs signature operations with encryption
func LogSignatureOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SignatureOperation", encryptedMessage)
}
