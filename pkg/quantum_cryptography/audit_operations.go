package quantum_cryptography

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

var qrEncryptionAuditEnabled bool
var forwardSecrecyEnabled bool
var postQuantumTrustChecksEnabled bool
var isogenyAuditLogEnabled bool
var qrKeyUsageHistory = make(map[string]int) // Tracks quantum-resistant key usage by ID
var auditLock sync.Mutex

// EnableQREncryptionAudit: Activates auditing for quantum-resistant encryption
func EnableQREncryptionAudit() {
    auditLock.Lock()
    defer auditLock.Unlock()
    qrEncryptionAuditEnabled = true
    LogAuditOperation("EnableQREncryptionAudit", "Quantum-resistant encryption audit enabled")
}

// DisableQREncryptionAudit: Deactivates auditing for quantum-resistant encryption
func DisableQREncryptionAudit() {
    auditLock.Lock()
    defer auditLock.Unlock()
    qrEncryptionAuditEnabled = false
    LogAuditOperation("DisableQREncryptionAudit", "Quantum-resistant encryption audit disabled")
}

// TrackQRKeyUsage: Logs the usage count of a specific quantum-resistant encryption key
func TrackQRKeyUsage(keyID string) {
    auditLock.Lock()
    defer auditLock.Unlock()
    if qrEncryptionAuditEnabled {
        qrKeyUsageHistory[keyID]++
        LogAuditOperation("TrackQRKeyUsage", fmt.Sprintf("Usage count for key %s incremented to %d", keyID, qrKeyUsageHistory[keyID]))
    }
}

// ClearQRKeyHistory: Clears the history of quantum-resistant key usage
func ClearQRKeyHistory() {
    auditLock.Lock()
    defer auditLock.Unlock()
    qrKeyUsageHistory = make(map[string]int)
    LogAuditOperation("ClearQRKeyHistory", "Quantum-resistant key usage history cleared")
}

// EnableQRForwardSecrecy: Enables forward secrecy for quantum-resistant encryption
func EnableQRForwardSecrecy() {
    auditLock.Lock()
    defer auditLock.Unlock()
    forwardSecrecyEnabled = true
    LogAuditOperation("EnableQRForwardSecrecy", "Quantum-resistant forward secrecy enabled")
}

// DisableQRForwardSecrecy: Disables forward secrecy for quantum-resistant encryption
func DisableQRForwardSecrecy() {
    auditLock.Lock()
    defer auditLock.Unlock()
    forwardSecrecyEnabled = false
    LogAuditOperation("DisableQRForwardSecrecy", "Quantum-resistant forward secrecy disabled")
}

// EnablePostQuantumTrustChecks: Activates post-quantum trust checks
func EnablePostQuantumTrustChecks() {
    auditLock.Lock()
    defer auditLock.Unlock()
    postQuantumTrustChecksEnabled = true
    LogAuditOperation("EnablePostQuantumTrustChecks", "Post-quantum trust checks enabled")
}

// DisablePostQuantumTrustChecks: Deactivates post-quantum trust checks
func DisablePostQuantumTrustChecks() {
    auditLock.Lock()
    defer auditLock.Unlock()
    postQuantumTrustChecksEnabled = false
    LogAuditOperation("DisablePostQuantumTrustChecks", "Post-quantum trust checks disabled")
}

// EnableIsogenyAuditLog: Activates audit logging for isogeny-based cryptographic operations
func EnableIsogenyAuditLog() {
    auditLock.Lock()
    defer auditLock.Unlock()
    isogenyAuditLogEnabled = true
    LogAuditOperation("EnableIsogenyAuditLog", "Isogeny-based audit logging enabled")
}

// DisableIsogenyAuditLog: Deactivates audit logging for isogeny-based cryptographic operations
func DisableIsogenyAuditLog() {
    auditLock.Lock()
    defer auditLock.Unlock()
    isogenyAuditLogEnabled = false
    LogAuditOperation("DisableIsogenyAuditLog", "Isogeny-based audit logging disabled")
}

// Helper Functions

// LogAuditOperation: Logs audit operations with encryption
func LogAuditOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("AuditOperation", encryptedMessage)
}
