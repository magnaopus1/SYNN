package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    KeyManagementInterval         = 30 * time.Second // Interval for checking key usage and security
    MaxKeyRotationRetries         = 3                // Maximum retries for key rotation
    SubBlocksPerBlock             = 1000             // Number of sub-blocks in a block
    MaxKeyUsageBeforeRotation     = 1000             // Maximum key usage count before rotation is triggered
    EncryptionKeyExpirationPeriod = 90 * 24 * time.Hour // Period before encryption key expiration
)

// KeyManagementSecurityProtocol monitors and enforces key security, including key rotation
type KeyManagementSecurityProtocol struct {
    consensusSystem        *consensus.SynnergyConsensus  // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger                // Ledger for logging key management-related events
    stateMutex             *sync.RWMutex                 // Mutex for thread-safe access
    keyUsageCount          map[string]int                // Tracks how many times a key is used
    keyRotationRetryCount  map[string]int                // Tracks retries for failed key rotations
    lastKeyRotation        map[string]time.Time          // Tracks the last time a key was rotated
    rotationCycleCount     int                           // Counter for key rotation cycles
}

// NewKeyManagementSecurityProtocol initializes the automation for key management
func NewKeyManagementSecurityProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *KeyManagementSecurityProtocol {
    return &KeyManagementSecurityProtocol{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        keyUsageCount:         make(map[string]int),
        keyRotationRetryCount: make(map[string]int),
        lastKeyRotation:       make(map[string]time.Time),
        rotationCycleCount:    0,
    }
}

// StartKeyManagement starts the continuous loop for monitoring and enforcing key management policies
func (protocol *KeyManagementSecurityProtocol) StartKeyManagement() {
    ticker := time.NewTicker(KeyManagementInterval)

    go func() {
        for range ticker.C {
            protocol.monitorAndEnforceKeyPolicies()
        }
    }()
}

// monitorAndEnforceKeyPolicies checks key usage, triggers key rotations, and enforces security policies
func (protocol *KeyManagementSecurityProtocol) monitorAndEnforceKeyPolicies() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch the list of keys used within the system
    activeKeys := protocol.consensusSystem.FetchActiveKeys()

    for _, key := range activeKeys {
        if protocol.shouldRotateKey(key) {
            fmt.Printf("Key rotation required for key: %s\n", key.ID)
            protocol.rotateKey(key)
        } else {
            fmt.Printf("Key %s is within usage limits.\n", key.ID)
        }
    }

    protocol.rotationCycleCount++
    fmt.Printf("Key management cycle #%d completed.\n", protocol.rotationCycleCount)

    if protocol.rotationCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeKeyRotationCycle()
    }
}

// shouldRotateKey checks if a key has exceeded its usage or expiration period
func (protocol *KeyManagementSecurityProtocol) shouldRotateKey(key common.EncryptionKey) bool {
    usageCount, exists := protocol.keyUsageCount[key.ID]
    if !exists {
        protocol.keyUsageCount[key.ID] = 0
    }

    lastRotation, exists := protocol.lastKeyRotation[key.ID]
    if !exists {
        protocol.lastKeyRotation[key.ID] = time.Now()
    }

    // Check if the key has exceeded its maximum usage or expiration period
    if usageCount >= MaxKeyUsageBeforeRotation || time.Since(lastRotation) >= EncryptionKeyExpirationPeriod {
        return true
    }
    return false
}

// rotateKey attempts to rotate the encryption key and update the consensus and ledger
func (protocol *KeyManagementSecurityProtocol) rotateKey(key common.EncryptionKey) {
    encryptedKeyData := protocol.encryptKeyData(key)

    // Attempt to rotate the key through the consensus system
    rotationSuccess := protocol.consensusSystem.RotateEncryptionKey(key, encryptedKeyData)

    if rotationSuccess {
        fmt.Printf("Key rotation successful for key: %s.\n", key.ID)
        protocol.logKeyRotationEvent(key, "Rotated")
        protocol.resetKeyRotationRetry(key.ID)
        protocol.lastKeyRotation[key.ID] = time.Now()
    } else {
        fmt.Printf("Error rotating key: %s. Retrying...\n", key.ID)
        protocol.retryKeyRotation(key)
    }
}

// retryKeyRotation retries key rotation in case of failure
func (protocol *KeyManagementSecurityProtocol) retryKeyRotation(key common.EncryptionKey) {
    protocol.keyRotationRetryCount[key.ID]++
    if protocol.keyRotationRetryCount[key.ID] < MaxKeyRotationRetries {
        protocol.rotateKey(key)
    } else {
        fmt.Printf("Max retries reached for rotating key: %s. Key rotation failed.\n", key.ID)
        protocol.logKeyRotationFailure(key)
    }
}

// resetKeyRotationRetry resets the retry count for a key rotation
func (protocol *KeyManagementSecurityProtocol) resetKeyRotationRetry(keyID string) {
    protocol.keyRotationRetryCount[keyID] = 0
}

// finalizeKeyRotationCycle finalizes the key rotation cycle and logs the result in the ledger
func (protocol *KeyManagementSecurityProtocol) finalizeKeyRotationCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeKeyRotationCycle()
    if success {
        fmt.Println("Key rotation cycle finalized successfully.")
        protocol.logKeyRotationCycleFinalization()
    } else {
        fmt.Println("Error finalizing key rotation cycle.")
    }
}

// logKeyRotationEvent logs a successful key rotation into the ledger
func (protocol *KeyManagementSecurityProtocol) logKeyRotationEvent(key common.EncryptionKey, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("key-rotation-%s-%s", key.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Key Rotation Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Encryption key %s was %s.", key.ID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with key rotation event for key: %s.\n", key.ID)
}

// logKeyRotationFailure logs the failure of a key rotation into the ledger
func (protocol *KeyManagementSecurityProtocol) logKeyRotationFailure(key common.EncryptionKey) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("key-rotation-failure-%s", key.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Key Rotation Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to rotate encryption key: %s after maximum retries.", key.ID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with key rotation failure for key: %s.\n", key.ID)
}

// logKeyRotationCycleFinalization logs the finalization of a key rotation cycle into the ledger
func (protocol *KeyManagementSecurityProtocol) logKeyRotationCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("key-rotation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Key Rotation Cycle Finalization",
        Status:    "Finalized",
        Details:   "Key rotation cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with key rotation cycle finalization.")
}

// encryptKeyData encrypts key data before performing key rotation
func (protocol *KeyManagementSecurityProtocol) encryptKeyData(key common.EncryptionKey) common.EncryptionKey {
    encryptedData, err := encryption.EncryptData(key.Data)
    if err != nil {
        fmt.Println("Error encrypting key data:", err)
        return key
    }

    key.EncryptedData = encryptedData
    fmt.Println("Encryption key data successfully encrypted.")
    return key
}

// triggerEmergencyKeyRevocation triggers the revocation of keys in case of a security breach
func (protocol *KeyManagementSecurityProtocol) triggerEmergencyKeyRevocation(key common.EncryptionKey) {
    fmt.Printf("Emergency key revocation triggered for key: %s.\n", key.ID)
    success := protocol.consensusSystem.TriggerEmergencyKeyRevocation(key)

    if success {
        protocol.logKeyRotationEvent(key, "Revoked")
        fmt.Println("Emergency key revocation executed successfully.")
    } else {
        fmt.Println("Emergency key revocation failed.")
    }
}
