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
	ContractModificationCheckInterval = 10 * time.Second // Interval for checking contract modifications
	MaxModificationsPerContract       = 100              // Maximum number of modifications allowed per contract
	ModificationWindow                = 30 * 24 * time.Hour // Time window to track modifications (30 days)
)

// ContractModificationRestrictionAutomation monitors and restricts contract modifications across the network
type ContractModificationRestrictionAutomation struct {
	consensusSystem             *consensus.SynnergyConsensus
	ledgerInstance              *ledger.Ledger
	stateMutex                  *sync.RWMutex
	contractModificationCount   map[string]int // Tracks the number of modifications per contract
}

// NewContractModificationRestrictionAutomation initializes and returns an instance of ContractModificationRestrictionAutomation
func NewContractModificationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractModificationRestrictionAutomation {
	return &ContractModificationRestrictionAutomation{
		consensusSystem:           consensusSystem,
		ledgerInstance:            ledgerInstance,
		stateMutex:                stateMutex,
		contractModificationCount: make(map[string]int),
	}
}

// StartContractModificationMonitoring begins continuous monitoring of contract modifications
func (automation *ContractModificationRestrictionAutomation) StartContractModificationMonitoring() {
	ticker := time.NewTicker(ContractModificationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorContractModifications()
		}
	}()
}

// monitorContractModifications checks recent contract modifications and enforces modification limits
func (automation *ContractModificationRestrictionAutomation) monitorContractModifications() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent contract modifications from Synnergy Consensus
	recentModifications := automation.consensusSystem.GetRecentContractModifications()

	for _, modification := range recentModifications {
		// Validate modification limits
		if !automation.validateModificationLimit(modification) {
			automation.flagModificationViolation(modification, "Exceeded maximum number of modifications for this contract")
		}
	}
}

// validateModificationLimit checks if a contract has exceeded the modification limit within the given time window
func (automation *ContractModificationRestrictionAutomation) validateModificationLimit(modification common.ContractModification) bool {
	currentModCount := automation.contractModificationCount[modification.ContractID]
	if currentModCount+1 > MaxModificationsPerContract {
		return false
	}

	// Update the modification count for the contract
	automation.contractModificationCount[modification.ContractID]++
	return true
}

// flagModificationViolation flags a contract modification that violates the system's rules and logs it in the ledger
func (automation *ContractModificationRestrictionAutomation) flagModificationViolation(modification common.ContractModification, reason string) {
	fmt.Printf("Contract modification violation: Contract %s, Reason: %s\n", modification.ContractID, reason)

	// Log the violation into the ledger
	automation.logModificationViolation(modification, reason)
}

// logModificationViolation logs the flagged contract modification violation into the ledger with full details
func (automation *ContractModificationRestrictionAutomation) logModificationViolation(modification common.ContractModification, violationReason string) {
	// Encrypt the contract modification data
	encryptedData := automation.encryptModificationData(modification)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-modification-violation-%s-%d", modification.ContractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Contract Modification Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Contract %s flagged for modification violation. Reason: %s. Encrypted Data: %s", modification.ContractID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log contract modification violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Contract modification violation logged for contract: %s\n", modification.ContractID)
	}
}

// encryptModificationData encrypts contract modification data before logging for security
func (automation *ContractModificationRestrictionAutomation) encryptModificationData(modification common.ContractModification) string {
	data := fmt.Sprintf("Contract ID: %s, Modification Timestamp: %d", modification.ContractID, modification.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting contract modification data:", err)
		return data
	}
	return string(encryptedData)
}
