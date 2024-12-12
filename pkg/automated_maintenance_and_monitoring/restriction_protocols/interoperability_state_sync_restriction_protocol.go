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
	StateSyncCheckInterval      = 10 * time.Second // Interval for checking state synchronization
	MaxStateSyncFailures        = 5                // Maximum number of allowed state sync failures before restriction
	MaxStateSyncTime            = 1 * time.Minute  // Maximum time allowed for a state sync operation
)

// InteroperabilityStateSyncRestrictionAutomation monitors and restricts state synchronization processes across interoperable systems
type InteroperabilityStateSyncRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	stateSyncFailureCount  map[string]int // Tracks state sync failures per system
}

// NewInteroperabilityStateSyncRestrictionAutomation initializes and returns an instance of InteroperabilityStateSyncRestrictionAutomation
func NewInteroperabilityStateSyncRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *InteroperabilityStateSyncRestrictionAutomation {
	return &InteroperabilityStateSyncRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		stateSyncFailureCount: make(map[string]int),
	}
}

// StartStateSyncMonitoring starts continuous monitoring of state synchronization processes
func (automation *InteroperabilityStateSyncRestrictionAutomation) StartStateSyncMonitoring() {
	ticker := time.NewTicker(StateSyncCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorStateSync()
		}
	}()
}

// monitorStateSync checks recent state synchronization attempts and enforces restrictions
func (automation *InteroperabilityStateSyncRestrictionAutomation) monitorStateSync() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent state synchronization data from Synnergy Consensus
	recentStateSyncs := automation.consensusSystem.GetStateSyncData()

	for systemID, stateSync := range recentStateSyncs {
		// Validate state synchronization based on time and failure count
		if !automation.validateStateSync(systemID, stateSync) {
			automation.flagStateSyncViolation(systemID, stateSync, "Exceeded time limit or failure threshold for state synchronization")
		}
	}
}

// validateStateSync checks if a state synchronization operation meets time and failure requirements
func (automation *InteroperabilityStateSyncRestrictionAutomation) validateStateSync(systemID string, stateSync common.StateSync) bool {
	// Check if the state sync time exceeds the maximum allowed
	if stateSync.Duration > MaxStateSyncTime {
		automation.stateSyncFailureCount[systemID]++
		return false
	}

	// Check if the system has exceeded the maximum allowed failures
	if automation.stateSyncFailureCount[systemID] >= MaxStateSyncFailures {
		return false
	}

	// Reset failure count for successful state syncs
	automation.stateSyncFailureCount[systemID] = 0
	return true
}

// flagStateSyncViolation flags a state sync violation and logs it in the ledger
func (automation *InteroperabilityStateSyncRestrictionAutomation) flagStateSyncViolation(systemID string, stateSync common.StateSync, reason string) {
	fmt.Printf("State sync violation: System %s, Reason: %s\n", systemID, reason)

	// Log the violation into the ledger
	automation.logStateSyncViolation(systemID, stateSync, reason)
}

// logStateSyncViolation logs the flagged state sync violation into the ledger with full details
func (automation *InteroperabilityStateSyncRestrictionAutomation) logStateSyncViolation(systemID string, stateSync common.StateSync, violationReason string) {
	// Encrypt the state sync data before logging
	encryptedData := automation.encryptStateSyncData(systemID, stateSync)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("state-sync-violation-%s-%d", systemID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "State Sync Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("System %s flagged for state sync violation. Reason: %s. Encrypted Data: %s", systemID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log state sync violation into ledger: %v\n", err)
	} else {
		fmt.Printf("State sync violation logged for system: %s\n", systemID)
	}
}

// encryptStateSyncData encrypts state sync data before logging for security
func (automation *InteroperabilityStateSyncRestrictionAutomation) encryptStateSyncData(systemID string, stateSync common.StateSync) string {
	data := fmt.Sprintf("System ID: %s, Sync Duration: %s, Timestamp: %d", systemID, stateSync.Duration.String(), stateSync.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting state sync data:", err)
		return data
	}
	return string(encryptedData)
}
