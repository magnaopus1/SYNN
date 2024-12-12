package quantum_cryptography

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

type PostQuantumCertificate struct {
    CertID       string
    Owner        string
    PublicKey    []byte
    IssuedAt     time.Time
    ExpiresAt    time.Time
    Revoked      bool
}

var certificateStore = make(map[string]PostQuantumCertificate)
var certLock sync.Mutex

// RegisterPostQuantumCertificate: Registers a new post-quantum certificate on the blockchain
func RegisterPostQuantumCertificate(certID, owner string, publicKey []byte, validityPeriod time.Duration) error {
    certLock.Lock()
    defer certLock.Unlock()

    if _, exists := certificateStore[certID]; exists {
        LogCertificateOperation("RegisterPostQuantumCertificate", "Certificate already exists: "+certID)
        return errors.New("certificate already exists")
    }

    certificate := PostQuantumCertificate{
        CertID:    certID,
        Owner:     owner,
        PublicKey: publicKey,
        IssuedAt:  time.Now(),
        ExpiresAt: time.Now().Add(validityPeriod),
        Revoked:   false,
    }
    certificateStore[certID] = certificate

    LogCertificateOperation("RegisterPostQuantumCertificate", fmt.Sprintf("Registered certificate %s for owner %s", certID, owner))
    return nil
}

// RevokePostQuantumCertificate: Revokes an existing post-quantum certificate
func RevokePostQuantumCertificate(certID string) error {
    certLock.Lock()
    defer certLock.Unlock()

    cert, exists := certificateStore[certID]
    if !exists {
        LogCertificateOperation("RevokePostQuantumCertificate", "Certificate not found: "+certID)
        return errors.New("certificate not found")
    }

    if cert.Revoked {
        LogCertificateOperation("RevokePostQuantumCertificate", "Certificate already revoked: "+certID)
        return errors.New("certificate already revoked")
    }

    cert.Revoked = true
    certificateStore[certID] = cert

    LogCertificateOperation("RevokePostQuantumCertificate", "Revoked certificate "+certID)
    return nil
}

// VerifyPostQuantumCertificate: Verifies the authenticity and validity of a post-quantum certificate
func VerifyPostQuantumCertificate(certID string, publicKey []byte) (bool, error) {
    certLock.Lock()
    defer certLock.Unlock()

    cert, exists := certificateStore[certID]
    if !exists {
        LogCertificateOperation("VerifyPostQuantumCertificate", "Certificate not found: "+certID)
        return false, errors.New("certificate not found")
    }

    if cert.Revoked {
        LogCertificateOperation("VerifyPostQuantumCertificate", "Certificate is revoked: "+certID)
        return false, errors.New("certificate is revoked")
    }

    if time.Now().After(cert.ExpiresAt) {
        LogCertificateOperation("VerifyPostQuantumCertificate", "Certificate expired: "+certID)
        return false, errors.New("certificate expired")
    }

    publicKeyHash := sha256.Sum256(publicKey)
    certKeyHash := sha256.Sum256(cert.PublicKey)
    valid := publicKeyHash == certKeyHash

    LogCertificateOperation("VerifyPostQuantumCertificate", fmt.Sprintf("Verification result for certificate %s: %t", certID, valid))
    return valid, nil
}

// Helper Functions

// LogCertificateOperation: Logs certificate operations with encryption
func LogCertificateOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("CertificateOperation", encryptedMessage)
}
