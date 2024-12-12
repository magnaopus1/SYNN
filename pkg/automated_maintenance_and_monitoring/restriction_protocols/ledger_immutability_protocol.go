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
	ImmutabilityCheckInterval = 5 * time.Minute // Interval for checking ledger immutability
)

// LedgerImmutabilityProtocolAutomation ensures that ledger entries remain immutable and are protected from tampering
type LedgerImmutabilityProtocolAutomation struct {
	consensusSystem    *consensus.SynnergyConsensus
	ledgerInstance     *ledger.Ledger
	stateMutex         *sync.RWMutex
	immutableLedgerHashes map[string]string // Tracks hashes of immutable ledger entries
}

// NewLedgerImmutabilityProtocolAutomation initializes and returns an instance of LedgerImmutabilityProtocolAutomation
func NewLedgerImmutabilityProtocolAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *LedgerImmutabilityProtocolAutomation {
	return &LedgerImmutabilityProtocolAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		immutableLedgerHashes: make(map[string]string),
	}
}

// StartImmutabilityMonitoring starts continuous monitoring of ledger immutability
func (automation *LedgerImmutabilityProtocolAutomation) StartImmutabilityMonitoring() {
	ticker := time.NewTicker(ImmutabilityCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkLedgerImmutability()
		}
	}()
}

// checkLedgerImmutability verifies that all ledger entries remain immutable by comparing their stored hashes
func (automation *LedgerImmutabilityProtocolAutomation) checkLedgerImmutability() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch all ledger entries for verification
	ledgerEntries := automation.ledgerInstance.GetAllEntries()

	for _, entry := range ledgerEntries {
		// Get the stored hash for this ledger entry
		storedHash, exists := automation.immutableLedgerHashes[entry.ID]
		if !exists {
			// If no hash is stored, generate and store the hash for future verification
			newHash := automation.generateHash(entry)
			automation.immutableLedgerHashes[entry.ID] = newHash
			continue
		}

		// Generate a new hash of the current ledger entry for comparison
		currentHash := automation.generateHash(entry)

		// Compare the hashes to ensure immutability
		if currentHash != storedHash {
			// Trigger violation if the hashes do not match
			automation.flagImmutabilityViolation(entry)
		}
	}
}

// generateHash generates a hash of a ledger entry's data for immutability verification
func (automation *LedgerImmutabilityProtocolAutomation) generateHash(entry common.LedgerEntry) string {
	data := fmt.Sprintf("%s-%s-%d-%s", entry.ID, entry.Type, entry.Timestamp, entry.Details)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error generating hash for ledger entry:", err)
		return ""
	}
	return string(encryptedData)
}

// flagImmutabilityViolation flags a violation of ledger immutability and logs it in the ledger
func (automation *LedgerImmutabilityProtocolAutomation) flagImmutabilityViolation(entry common.LedgerEntry) {
	fmt.Printf("Immutability violation detected: Ledger Entry ID %s\n", entry.ID)

	// Log the violation in the ledger
	automation.logImmutabilityViolation(entry)
}

// logImmutabilityViolation logs the flagged immutability violation into the ledger with full details
func (automation *LedgerImmutabilityProtocolAutomation) logImmutabilityViolation(entry common.LedgerEntry) {
	violationDetails := fmt.Sprintf("Immutability violation detected for ledger entry %s.", entry.ID)

	// Create a ledger entry for immutability violation
	violationEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("immutability-violation-%s-%d", entry.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Ledger Immutability Violation",
		Status:    "Flagged",
		Details:   violationDetails,
	}

	// Add the immutability violation to the ledger
	err := automation.ledgerInstance.AddEntry(violationEntry)
	if err != nil {
		fmt.Println("Failed to log immutability violation:", err)
	} else {
		fmt.Println("Immutability violation logged.")
	}
}

