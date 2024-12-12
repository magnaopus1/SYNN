package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	KeyRotationInterval        = 24 * time.Hour // Interval for automatic key rotation
	MaxKeyAge                  = 30 * 24 * time.Hour // Max age before key rotation is enforced
	KeyRotationRetryLimit      = 3              // Number of retry attempts for key rotation
	KeyRotationLedgerEntryType = "Key Rotation"  // Ledger entry type for key rotation events
)

// KeyRotationExecutionAutomation manages encryption key rotation across the network
type KeyRotationExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
	ledgerInstance    *ledger.Ledger                        // Ledger instance for logging key rotations
	encryptionManager *encryption.Manager                   // Encryption manager responsible for key operations
	keyRotationMutex  *sync.RWMutex                         // Mutex for thread-safe key operations
}

// NewKeyRotationExecutionAutomation initializes the key rotation automation
func NewKeyRotationExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, encryptionManager *encryption.Manager, keyRotationMutex *sync.RWMutex) *KeyRotationExecutionAutomation {
	return &KeyRotationExecutionAutomation{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		encryptionManager: encryptionManager,
		keyRotationMutex:  keyRotationMutex,
	}
}

// StartKeyRotationMonitoring starts the automatic key rotation process
func (automation *KeyRotationExecutionAutomation) StartKeyRotationMonitoring() {
	ticker := time.NewTicker(KeyRotationInterval)

	go func() {
		for range ticker.C {
			automation.checkAndRotateKeys()
		}
	}()
}

// checkAndRotateKeys checks the current encryption keys and rotates them if necessary
func (automation *KeyRotationExecutionAutomation) checkAndRotateKeys() {
	automation.keyRotationMutex.Lock()
	defer automation.keyRotationMutex.Unlock()

	keys := automation.encryptionManager.GetAllKeys()

	for _, key := range keys {
		if time.Since(key.CreatedAt) > MaxKeyAge {
			automation.rotateKey(key)
		}
	}
}

// rotateKey handles the rotation of a specific encryption key
func (automation *KeyRotationExecutionAutomation) rotateKey(key *encryption.Key) {
	for i := 0; i < KeyRotationRetryLimit; i++ {
		fmt.Printf("Attempting to rotate key: %s (Attempt %d)\n", key.ID, i+1)
		
		newKey, err := automation.encryptionManager.RotateKey(key)
		if err != nil {
			fmt.Printf("Key rotation failed for key %s: %v\n", key.ID, err)
			continue
		}

		// Log the successful key rotation into the ledger
		automation.logKeyRotationInLedger(key.ID, newKey.ID)

		// Notify the consensus engine about the key rotation
		automation.consensusEngine.NotifyKeyRotation(key.ID, newKey.ID)

		fmt.Printf("Successfully rotated key: %s -> %s\n", key.ID, newKey.ID)
		return
	}

	fmt.Printf("Failed to rotate key %s after %d attempts.\n", key.ID, KeyRotationRetryLimit)
}

// logKeyRotationInLedger securely logs the key rotation event in the ledger
func (automation *KeyRotationExecutionAutomation) logKeyRotationInLedger(oldKeyID string, newKeyID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("key-rotation-%s-%d", oldKeyID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      KeyRotationLedgerEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("Rotated key from %s to %s", oldKeyID, newKeyID),
	}

	// Encrypt the details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log key rotation in the ledger for key %s: %v\n", oldKeyID, err)
	} else {
		fmt.Printf("Key rotation successfully logged in the ledger for key %s.\n", oldKeyID)
	}
}

// encryptData encrypts sensitive data before logging in the ledger
func (automation *KeyRotationExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualKeyRotation allows administrators to manually rotate a specific key
func (automation *KeyRotationExecutionAutomation) TriggerManualKeyRotation(keyID string) {
	fmt.Printf("Manually triggering key rotation for key %s...\n", keyID)

	key := automation.encryptionManager.GetKeyByID(keyID)
	if key != nil {
		automation.rotateKey(key)
	} else {
		fmt.Printf("Key %s not found.\n", keyID)
	}
}
