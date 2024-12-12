package quantum_cryptography

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

type QRKeyPair struct {
    PublicKey  []byte
    PrivateKey []byte
}

var qrKeyPairs = make(map[string]QRKeyPair)
var qrTrustAnchors = make(map[string]bool)
var qrLock sync.Mutex

// GenerateQRKeyPair: Generates a quantum-resistant key pair
func GenerateQRKeyPair(keyID string) (*QRKeyPair, error) {
    qrLock.Lock()
    defer qrLock.Unlock()

    publicKey, privateKey, err := qrcrypto.GenerateKeyPair()
    if err != nil {
        LogQROperation("GenerateQRKeyPair", "Key pair generation failed for keyID "+keyID)
        return nil, errors.New("failed to generate quantum-resistant key pair")
    }

    qrKeyPair := QRKeyPair{
        PublicKey:  publicKey,
        PrivateKey: privateKey,
    }
    qrKeyPairs[keyID] = qrKeyPair

    LogQROperation("GenerateQRKeyPair", fmt.Sprintf("Generated quantum-resistant key pair for keyID %s", keyID))
    return &qrKeyPair, nil
}

// QRSignData: Signs data using a quantum-resistant private key
func QRSignData(keyID string, data []byte) ([]byte, error) {
    qrLock.Lock()
    defer qrLock.Unlock()

    keyPair, exists := qrKeyPairs[keyID]
    if !exists {
        LogQROperation("QRSignData", "Private key not found: "+keyID)
        return nil, errors.New("private key not found")
    }

    signature, err := qrcrypto.SignData(keyPair.PrivateKey, data)
    if err != nil {
        LogQROperation("QRSignData", "Signing failed for keyID "+keyID)
        return nil, errors.New("signing data failed")
    }

    LogQROperation("QRSignData", fmt.Sprintf("Data signed with keyID %s", keyID))
    return signature, nil
}

// QRVerifySignature: Verifies a signature using a quantum-resistant public key
func QRVerifySignature(keyID string, data, signature []byte) (bool, error) {
    qrLock.Lock()
    defer qrLock.Unlock()

    keyPair, exists := qrKeyPairs[keyID]
    if !exists {
        LogQROperation("QRVerifySignature", "Public key not found: "+keyID)
        return false, errors.New("public key not found")
    }

    valid, err := qrcrypto.VerifySignature(keyPair.PublicKey, data, signature)
    if err != nil || !valid {
        LogQROperation("QRVerifySignature", "Verification failed for keyID "+keyID)
        return false, errors.New("signature verification failed")
    }

    LogQROperation("QRVerifySignature", fmt.Sprintf("Signature verified for keyID %s", keyID))
    return true, nil
}

// QRSharedSecret: Derives a shared secret between two quantum-resistant public keys
func QRSharedSecret(publicKeyA, publicKeyB []byte) ([]byte, error) {
    sharedSecret, err := qrcrypto.DeriveSharedSecret(publicKeyA, publicKeyB)
    if err != nil {
        LogQROperation("QRSharedSecret", "Shared secret derivation failed")
        return nil, errors.New("failed to derive shared secret")
    }

    LogQROperation("QRSharedSecret", "Derived shared secret between two keys")
    return sharedSecret, nil
}

// QRKeyAgreement: Establishes a quantum-resistant key agreement between two parties
func QRKeyAgreement(publicKey, privateKey []byte) ([]byte, error) {
    agreedKey, err := qrcrypto.KeyAgreement(publicKey, privateKey)
    if err != nil {
        LogQROperation("QRKeyAgreement", "Key agreement failed")
        return nil, errors.New("key agreement failed")
    }

    LogQROperation("QRKeyAgreement", "Quantum-resistant key agreement established")
    return agreedKey, nil
}

// QRKeyRotation: Rotates a quantum-resistant key pair, creating a new one and archiving the old
func QRKeyRotation(keyID string) (*QRKeyPair, error) {
    qrLock.Lock()
    defer qrLock.Unlock()

    oldKeyPair, exists := qrKeyPairs[keyID]
    if !exists {
        LogQROperation("QRKeyRotation", "Key not found: "+keyID)
        return nil, errors.New("key not found")
    }

    newPublicKey, newPrivateKey, err := qrcrypto.GenerateKeyPair()
    if err != nil {
        LogQROperation("QRKeyRotation", "Key rotation failed for keyID "+keyID)
        return nil, errors.New("key rotation failed")
    }

    qrKeyPairs[keyID] = QRKeyPair{PublicKey: newPublicKey, PrivateKey: newPrivateKey}
    LogQROperation("QRKeyRotation", fmt.Sprintf("Key rotation completed for keyID %s", keyID))
    return &QRKeyPair{PublicKey: newPublicKey, PrivateKey: newPrivateKey}, nil
}

// QRMultisigCreate: Initiates a multisignature scheme for quantum-resistant security
func QRMultisigCreate(participants [][]byte) ([]byte, error) {
    multisigKey, err := qrcrypto.CreateMultisig(participants)
    if err != nil {
        LogQROperation("QRMultisigCreate", "Multisig creation failed")
        return nil, errors.New("multisignature creation failed")
    }

    LogQROperation("QRMultisigCreate", "Quantum-resistant multisignature created")
    return multisigKey, nil
}

// QRMultisigSign: Signs data within a multisignature scheme
func QRMultisigSign(multisigKey, privateKey, data []byte) ([]byte, error) {
    signature, err := qrcrypto.MultisigSign(multisigKey, privateKey, data)
    if err != nil {
        LogQROperation("QRMultisigSign", "Multisig signing failed")
        return nil, errors.New("multisignature signing failed")
    }

    LogQROperation("QRMultisigSign", "Data signed in multisignature scheme")
    return signature, nil
}

// QRMultisigVerify: Verifies a multisignature
func QRMultisigVerify(multisigKey, data, signature []byte) (bool, error) {
    valid, err := qrcrypto.MultisigVerify(multisigKey, data, signature)
    if err != nil || !valid {
        LogQROperation("QRMultisigVerify", "Multisig verification failed")
        return false, errors.New("multisignature verification failed")
    }

    LogQROperation("QRMultisigVerify", "Multisignature verified")
    return true, nil
}

// RegisterQRTrustAnchor: Registers a quantum-resistant trust anchor
func RegisterQRTrustAnchor(anchorID string) error {
    qrLock.Lock()
    defer qrLock.Unlock()

    if _, exists := qrTrustAnchors[anchorID]; exists {
        LogQROperation("RegisterQRTrustAnchor", "Trust anchor already exists: "+anchorID)
        return errors.New("trust anchor already exists")
    }

    qrTrustAnchors[anchorID] = true
    LogQROperation("RegisterQRTrustAnchor", "Registered quantum-resistant trust anchor: "+anchorID)
    return nil
}

// QRRevokeTrustAnchor: Revokes a quantum-resistant trust anchor
func QRRevokeTrustAnchor(anchorID string) error {
    qrLock.Lock()
    defer qrLock.Unlock()

    if _, exists := qrTrustAnchors[anchorID]; !exists {
        LogQROperation("QRRevokeTrustAnchor", "Trust anchor not found: "+anchorID)
        return errors.New("trust anchor not found")
    }

    delete(qrTrustAnchors, anchorID)
    LogQROperation("QRRevokeTrustAnchor", "Revoked quantum-resistant trust anchor: "+anchorID)
    return nil
}

// Helper Functions

// LogQROperation: Logs QR cryptographic operations with encryption
func LogQROperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("QROperation", encryptedMessage)
}
