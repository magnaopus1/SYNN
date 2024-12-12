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
	VotingRightCheckInterval = 12 * time.Hour // Interval for checking governance voting rights
	MinVotingTokens          = 1.0          // Minimum tokens required to have voting rights
	MaxVotingTokens          = 10000000000.0       // Maximum tokens allowed for voting to prevent centralization
)

// GovernanceVotingRightRestrictionAutomation monitors and restricts governance voting rights across the network
type GovernanceVotingRightRestrictionAutomation struct {
	consensusSystem      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	stateMutex           *sync.RWMutex
	userVotingRightCount map[string]float64 // Tracks governance tokens per user for voting rights
}

// NewGovernanceVotingRightRestrictionAutomation initializes and returns an instance of GovernanceVotingRightRestrictionAutomation
func NewGovernanceVotingRightRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GovernanceVotingRightRestrictionAutomation {
	return &GovernanceVotingRightRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		userVotingRightCount: make(map[string]float64),
	}
}

// StartVotingRightMonitoring starts continuous monitoring of governance voting rights
func (automation *GovernanceVotingRightRestrictionAutomation) StartVotingRightMonitoring() {
	ticker := time.NewTicker(VotingRightCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorVotingRights()
		}
	}()
}

// monitorVotingRights checks recent governance token holdings and enforces voting right restrictions
func (automation *GovernanceVotingRightRestrictionAutomation) monitorVotingRights() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent governance token holdings from Synnergy Consensus
	recentVotingTokens := automation.consensusSystem.GetUserVotingTokenHoldings()

	for userID, tokenCount := range recentVotingTokens {
		// Validate voting right based on token holding limits
		if !automation.validateVotingRight(userID, tokenCount) {
			automation.flagVotingRightViolation(userID, tokenCount, "User's voting rights do not meet the required token limits")
		}
	}
}

// validateVotingRight checks if a user's governance token holdings meet the voting right requirements
func (automation *GovernanceVotingRightRestrictionAutomation) validateVotingRight(userID string, tokenCount float64) bool {
	if tokenCount < MinVotingTokens || tokenCount > MaxVotingTokens {
		return false
	}
	return true
}

// flagVotingRightViolation flags a governance voting right violation and logs it in the ledger
func (automation *GovernanceVotingRightRestrictionAutomation) flagVotingRightViolation(userID string, tokenCount float64, reason string) {
	fmt.Printf("Governance voting right violation: User %s, Reason: %s, Token Count: %.2f\n", userID, reason, tokenCount)

	// Log the violation into the ledger
	automation.logVotingRightViolation(userID, tokenCount, reason)
}

// logVotingRightViolation logs the flagged governance voting right violation into the ledger with full details
func (automation *GovernanceVotingRightRestrictionAutomation) logVotingRightViolation(userID string, tokenCount float64, violationReason string) {
	// Encrypt the governance token holding data before logging
	encryptedData := automation.encryptVotingRightData(userID, tokenCount)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("voting-right-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Governance Voting Right Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for governance voting right violation. Reason: %s. Encrypted Data: %s", userID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance voting right violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Governance voting right violation logged for user: %s\n", userID)
	}
}

// encryptVotingRightData encrypts governance token holding data before logging for security
func (automation *GovernanceVotingRightRestrictionAutomation) encryptVotingRightData(userID string, tokenCount float64) string {
	data := fmt.Sprintf("User ID: %s, Voting Token Holdings: %.2f", userID, tokenCount)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting voting right data:", err)
		return data
	}
	return string(encryptedData)
}
