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
	AutoSettlementCheckInterval   = 5 * time.Second  // Interval for checking Layer 2 auto settlements
	MaxSettlementFailures         = 3                // Maximum number of settlement failures before restriction
	MaxSettlementDuration         = 10 * time.Minute // Maximum duration allowed for a Layer 2 settlement
)

// Layer2AutoSettlementRestrictionAutomation monitors and restricts Layer 2 auto-settlement processes across the network
type Layer2AutoSettlementRestrictionAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	settlementFailureCount map[string]int // Tracks settlement failures per process
}

// NewLayer2AutoSettlementRestrictionAutomation initializes and returns an instance of Layer2AutoSettlementRestrictionAutomation
func NewLayer2AutoSettlementRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *Layer2AutoSettlementRestrictionAutomation {
	return &Layer2AutoSettlementRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		settlementFailureCount: make(map[string]int),
	}
}

// StartSettlementMonitoring starts continuous monitoring of Layer 2 auto-settlement processes
func (automation *Layer2AutoSettlementRestrictionAutomation) StartSettlementMonitoring() {
	ticker := time.NewTicker(AutoSettlementCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorSettlementProcesses()
		}
	}()
}

// monitorSettlementProcesses checks recent Layer 2 auto-settlement attempts and enforces restrictions
func (automation *Layer2AutoSettlementRestrictionAutomation) monitorSettlementProcesses() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent Layer 2 settlements from Synnergy Consensus
	recentSettlements := automation.consensusSystem.GetLayer2Settlements()

	for processID, settlement := range recentSettlements {
		// Validate settlement time and failure count
		if !automation.validateSettlement(processID, settlement) {
			automation.flagSettlementViolation(processID, settlement, "Exceeded maximum settlement duration or failure threshold")
		}
	}
}

// validateSettlement checks if a Layer 2 auto-settlement process meets time and failure requirements
func (automation *Layer2AutoSettlementRestrictionAutomation) validateSettlement(processID string, settlement common.Layer2Settlement) bool {
	// Check if the settlement duration exceeds the maximum allowed time
	if settlement.Duration > MaxSettlementDuration {
		automation.settlementFailureCount[processID]++
		return false
	}

	// Check if the process has exceeded the maximum allowed settlement failures
	if automation.settlementFailureCount[processID] >= MaxSettlementFailures {
		return false
	}

	// Reset failure count for successful settlements
	automation.settlementFailureCount[processID] = 0
	return true
}

// flagSettlementViolation flags a settlement violation and logs it in the ledger
func (automation *Layer2AutoSettlementRestrictionAutomation) flagSettlementViolation(processID string, settlement common.Layer2Settlement, reason string) {
	fmt.Printf("Layer 2 auto-settlement violation: Process %s, Reason: %s\n", processID, reason)

	// Log the violation into the ledger
	automation.logSettlementViolation(processID, settlement, reason)
}

// logSettlementViolation logs the flagged settlement violation into the ledger with full details
func (automation *Layer2AutoSettlementRestrictionAutomation) logSettlementViolation(processID string, settlement common.Layer2Settlement, violationReason string) {
	// Encrypt the settlement data before logging
	encryptedData := automation.encryptSettlementData(processID, settlement)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("settlement-violation-%s-%d", processID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Layer 2 Auto Settlement Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Process %s flagged for auto-settlement violation. Reason: %s. Encrypted Data: %s", processID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log auto-settlement violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Auto-settlement violation logged for process: %s\n", processID)
	}
}

// encryptSettlementData encrypts settlement data before logging for security
func (automation *Layer2AutoSettlementRestrictionAutomation) encryptSettlementData(processID string, settlement common.Layer2Settlement) string {
	data := fmt.Sprintf("Process ID: %s, Settlement Duration: %s, Timestamp: %d", processID, settlement.Duration.String(), settlement.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting settlement data:", err)
		return data
	}
	return string(encryptedData)
}
