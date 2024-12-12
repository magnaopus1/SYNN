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
	TokenHoldingCheckInterval     = 12 * time.Hour // Interval for checking governance token holdings
	MaxGovernanceTokenHolding     = 100000.0       // Maximum allowed governance token holding per user
	MinGovernanceTokenHolding     = 100.0          // Minimum required governance token holding per user
)

// GovernanceTokenHoldingRestrictionAutomation monitors and restricts governance token holdings across the network
type GovernanceTokenHoldingRestrictionAutomation struct {
	consensusSystem           *consensus.SynnergyConsensus
	ledgerInstance            *ledger.Ledger
	stateMutex                *sync.RWMutex
	userGovernanceTokenHoldings map[string]float64 // Tracks governance token holdings per user
}

// NewGovernanceTokenHoldingRestrictionAutomation initializes and returns an instance of GovernanceTokenHoldingRestrictionAutomation
func NewGovernanceTokenHoldingRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GovernanceTokenHoldingRestrictionAutomation {
	return &GovernanceTokenHoldingRestrictionAutomation{
		consensusSystem:          consensusSystem,
		ledgerInstance:           ledgerInstance,
		stateMutex:               stateMutex,
		userGovernanceTokenHoldings: make(map[string]float64),
	}
}

// StartTokenHoldingMonitoring starts continuous monitoring of governance token holdings
func (automation *GovernanceTokenHoldingRestrictionAutomation) StartTokenHoldingMonitoring() {
	ticker := time.NewTicker(TokenHoldingCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorTokenHoldings()
		}
	}()
}

// monitorTokenHoldings checks recent governance token holdings and enforces holding restrictions
func (automation *GovernanceTokenHoldingRestrictionAutomation) monitorTokenHoldings() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent governance token holdings from Synnergy Consensus
	recentTokenHoldings := automation.consensusSystem.GetUserGovernanceTokenHoldings()

	for userID, holdings := range recentTokenHoldings {
		// Validate governance token holding limits
		if !automation.validateTokenHoldingLimit(userID, holdings) {
			automation.flagTokenHoldingViolation(userID, holdings, "Exceeded maximum allowed governance token holdings")
		} else if !automation.validateTokenHoldingMinimum(userID, holdings) {
			automation.flagTokenHoldingViolation(userID, holdings, "Governance token holdings below the minimum required")
		}
	}
}

// validateTokenHoldingLimit checks if a user has exceeded the maximum governance token holding limit
func (automation *GovernanceTokenHoldingRestrictionAutomation) validateTokenHoldingLimit(userID string, holdings float64) bool {
	return holdings <= MaxGovernanceTokenHolding
}

// validateTokenHoldingMinimum checks if a user meets the minimum governance token holding requirement
func (automation *GovernanceTokenHoldingRestrictionAutomation) validateTokenHoldingMinimum(userID string, holdings float64) bool {
	return holdings >= MinGovernanceTokenHolding
}

// flagTokenHoldingViolation flags a governance token holding violation and logs it in the ledger
func (automation *GovernanceTokenHoldingRestrictionAutomation) flagTokenHoldingViolation(userID string, holdings float64, reason string) {
	fmt.Printf("Governance token holding violation: User %s, Reason: %s, Holdings: %.2f\n", userID, reason, holdings)

	// Log the violation into the ledger
	automation.logTokenHoldingViolation(userID, holdings, reason)
}

// logTokenHoldingViolation logs the flagged governance token holding violation into the ledger with full details
func (automation *GovernanceTokenHoldingRestrictionAutomation) logTokenHoldingViolation(userID string, holdings float64, violationReason string) {
	// Encrypt the governance token holding data before logging
	encryptedData := automation.encryptTokenHoldingData(userID, holdings)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("token-holding-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Governance Token Holding Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for governance token holding violation. Reason: %s. Encrypted Data: %s", userID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance token holding violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Governance token holding violation logged for user: %s\n", userID)
	}
}

// encryptTokenHoldingData encrypts governance token holding data before logging for security
func (automation *GovernanceTokenHoldingRestrictionAutomation) encryptTokenHoldingData(userID string, holdings float64) string {
	data := fmt.Sprintf("User ID: %s, Governance Token Holdings: %.2f", userID, holdings)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting token holding data:", err)
		return data
	}
	return string(encryptedData)
}
