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
	ProposalCheckInterval      = 12 * time.Hour // Interval for checking governance proposal submissions
	MaxProposalsPerUser        = 5              // Maximum number of proposals allowed per user
	MinimumStakeForProposal    = 1000.0         // Minimum stake required to submit a proposal
)

// GovernanceProposalRestrictionAutomation monitors and restricts the submission of governance proposals across the network
type GovernanceProposalRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	userProposalCount      map[string]int // Tracks the number of proposals submitted per user
}

// NewGovernanceProposalRestrictionAutomation initializes and returns an instance of GovernanceProposalRestrictionAutomation
func NewGovernanceProposalRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GovernanceProposalRestrictionAutomation {
	return &GovernanceProposalRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		userProposalCount:     make(map[string]int),
	}
}

// StartProposalMonitoring starts continuous monitoring of governance proposal submissions
func (automation *GovernanceProposalRestrictionAutomation) StartProposalMonitoring() {
	ticker := time.NewTicker(ProposalCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorProposals()
		}
	}()
}

// monitorProposals checks recent governance proposals and enforces submission restrictions
func (automation *GovernanceProposalRestrictionAutomation) monitorProposals() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent proposals from Synnergy Consensus
	recentProposals := automation.consensusSystem.GetRecentProposals()

	for _, proposal := range recentProposals {
		// Validate proposal submission limits and minimum stake requirements
		if !automation.validateProposalLimit(proposal) {
			automation.flagProposalViolation(proposal, "Exceeded maximum number of proposals per user")
		} else if !automation.validateStakeRequirement(proposal) {
			automation.flagProposalViolation(proposal, "Insufficient stake to submit the proposal")
		}
	}
}

// validateProposalLimit checks if a user has exceeded the maximum number of allowed proposals
func (automation *GovernanceProposalRestrictionAutomation) validateProposalLimit(proposal common.GovernanceProposal) bool {
	currentProposalCount := automation.userProposalCount[proposal.UserID]
	if currentProposalCount+1 > MaxProposalsPerUser {
		return false
	}

	// Update the proposal count for the user
	automation.userProposalCount[proposal.UserID]++
	return true
}

// validateStakeRequirement checks if the user meets the minimum stake requirement for submitting a proposal
func (automation *GovernanceProposalRestrictionAutomation) validateStakeRequirement(proposal common.GovernanceProposal) bool {
	return proposal.Stake >= MinimumStakeForProposal
}

// flagProposalViolation flags a proposal that violates system rules and logs it in the ledger
func (automation *GovernanceProposalRestrictionAutomation) flagProposalViolation(proposal common.GovernanceProposal, reason string) {
	fmt.Printf("Governance proposal violation: User %s, Reason: %s\n", proposal.UserID, reason)

	// Log the violation into the ledger
	automation.logProposalViolation(proposal, reason)
}

// logProposalViolation logs the flagged proposal violation into the ledger with full details
func (automation *GovernanceProposalRestrictionAutomation) logProposalViolation(proposal common.GovernanceProposal, violationReason string) {
	// Encrypt the proposal data before logging
	encryptedData := automation.encryptProposalData(proposal)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("proposal-violation-%s-%d", proposal.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Governance Proposal Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for governance proposal violation. Reason: %s. Encrypted Data: %s", proposal.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance proposal violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Governance proposal violation logged for user: %s\n", proposal.UserID)
	}
}

// encryptProposalData encrypts proposal data before logging for security
func (automation *GovernanceProposalRestrictionAutomation) encryptProposalData(proposal common.GovernanceProposal) string {
	data := fmt.Sprintf("User ID: %s, Proposal ID: %s, Stake: %.2f, Timestamp: %d", proposal.UserID, proposal.ProposalID, proposal.Stake, proposal.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting proposal data:", err)
		return data
	}
	return string(encryptedData)
}
