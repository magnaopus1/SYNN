package quantum_cryptography

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

var qrSecurityLevel int
var qrKeys = make(map[string]QRKeyPair)
var securityLock sync.Mutex

// SetQRSecurityLevel: Sets the quantum-resistant security level for the system
func SetQRSecurityLevel(level int) error {
    securityLock.Lock()
    defer securityLock.Unlock()

    if level < 1 || level > 5 { // Assuming security level ranges from 1 to 5
        LogSecurityOperation("SetQRSecurityLevel", "Invalid security level: "+fmt.Sprint(level))
        return errors.New("invalid security level")
    }

    qrSecurityLevel = level
    LogSecurityOperation("SetQRSecurityLevel", fmt.Sprintf("Quantum-resistant security level set to %d", level))
    return nil
}

// GetQRSecurityLevel: Retrieves the current quantum-resistant security level
func GetQRSecurityLevel() (int, error) {
    securityLock.Lock()
    defer securityLock.Unlock()

    LogSecurityOperation("GetQRSecurityLevel", fmt.Sprintf("Current security level: %d", qrSecurityLevel))
    return qrSecurityLevel, nil
}

// QRBulkKeyRotation: Performs bulk rotation of all quantum-resistant keys to maintain security
func QRBulkKeyRotation() error {
    securityLock.Lock()
    defer securityLock.Unlock()

    for keyID, keyPair := range qrKeys {
        newPublicKey, newPrivateKey, err := qrcrypto.GenerateKeyPair()
        if err != nil {
            LogSecurityOperation("QRBulkKeyRotation", fmt.Sprintf("Key rotation failed for keyID %s", keyID))
            return errors.New("bulk key rotation failed")
        }

        // Update the key pair in the map
        qrKeys[keyID] = QRKeyPair{
            PublicKey:  newPublicKey,
            PrivateKey: newPrivateKey,
        }
        LogSecurityOperation("QRBulkKeyRotation", fmt.Sprintf("Key rotated for keyID %s", keyID))
    }

    LogSecurityOperation("QRBulkKeyRotation", "Bulk key rotation completed for all quantum-resistant keys")
    return nil
}

// Helper Functions

// LogSecurityOperation: Logs security operations with encryption
func LogSecurityOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SecurityOperation", encryptedMessage)
}
