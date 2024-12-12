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
	ContractTerminationCheckInterval = 10 * time.Second // Interval for checking contract terminations
	MaxTerminationsPerContract       = 1                // Maximum number of terminations allowed per contract
	TerminationTimeWindow            = 24 * 7 * time.Hour // Time window to monitor contract termination attempts (1 week)
)

// ContractTerminationRestrictionAutomation monitors and restricts contract terminations across the network
type ContractTerminationRestrictionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	contractTerminationCount map[string]int // Tracks contract termination attempts per contract
}

// NewContractTerminationRestrictionAutomation initializes and returns an instance of ContractTerminationRestrictionAutomation
func NewContractTerminationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractTerminationRestrictionAutomation {
	return &ContractTerminationRestrictionAutomation{
		consensusSystem:          consensusSystem,
		ledgerInstance:           ledgerInstance,
		stateMutex:               stateMutex,
		contractTerminationCount: make(map[string]int),
	}
}

// StartContractTerminationMonitoring begins continuous monitoring of contract terminations
func (automation *ContractTerminationRestrictionAutomation) StartContractTerminationMonitoring() {
	ticker := time.NewTicker(ContractTerminationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorContractTerminations()
		}
	}()
}

// monitorContractTerminations checks recent contract terminations and enforces termination limits
func (automation *ContractTerminationRestrictionAutomation) monitorContractTerminations() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent contract terminations from Synnergy Consensus
	recentTerminations := automation.consensusSystem.GetRecentContractTerminations()

	for _, termination := range recentTerminations {
		// Validate termination limits
		if !automation.validateTerminationLimit(termination) {
			automation.flagTerminationViolation(termination, "Exceeded maximum number of terminations for this contract")
		}
	}
}

// validateTerminationLimit checks if a contract has exceeded the termination limit within the time window
func (automation *ContractTerminationRestrictionAutomation) validateTerminationLimit(termination common.ContractTermination) bool {
	currentTermCount := automation.contractTerminationCount[termination.ContractID]
	if currentTermCount+1 > MaxTerminationsPerContract {
		return false
	}

	// Update the termination count for the contract
	automation.contractTerminationCount[termination.ContractID]++
	return true
}

// flagTerminationViolation flags a contract termination that violates system rules and logs it in the ledger
func (automation *ContractTerminationRestrictionAutomation) flagTerminationViolation(termination common.ContractTermination, reason string) {
	fmt.Printf("Contract termination violation: Contract %s, Reason: %s\n", termination.ContractID, reason)

	// Log the violation into the ledger
	automation.logTerminationViolation(termination, reason)
}

// logTerminationViolation logs the flagged contract termination violation into the ledger with full details
func (automation *ContractTerminationRestrictionAutomation) logTerminationViolation(termination common.ContractTermination, violationReason string) {
	// Encrypt the contract termination data
	encryptedData := automation.encryptTerminationData(termination)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-termination-violation-%s-%d", termination.ContractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Contract Termination Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Contract %s flagged for termination violation. Reason: %s. Encrypted Data: %s", termination.ContractID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log contract termination violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Contract termination violation logged for contract: %s\n", termination.ContractID)
	}
}

// encryptTerminationData encrypts contract termination data before logging for security
func (automation *ContractTerminationRestrictionAutomation) encryptTerminationData(termination common.ContractTermination) string {
	data := fmt.Sprintf("Contract ID: %s, Termination Timestamp: %d", termination.ContractID, termination.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting contract termination data:", err)
		return data
	}
	return string(encryptedData)
}
