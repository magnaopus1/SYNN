package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/governance"
)

// Configuration for DAO governance rule enforcement
const (
	GovernanceCheckInterval      = 20 * time.Second // Interval to check for rule compliance in DAO governance
	MaxProposalViolations        = 2                // Maximum violations before restricting proposal submission
	MaxVotingViolations          = 2                // Maximum violations before restricting voting rights
	MinimumQuorumPercentage      = 60               // Minimum quorum for proposal approval
	MinimumApprovalPercentage    = 50               // Minimum approval percentage for proposal passing
)

// DAOGovernanceRuleEnforcementAutomation monitors and enforces DAO governance compliance
type DAOGovernanceRuleEnforcementAutomation struct {
	governanceManager *governance.GovernanceManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	proposalViolations map[string]int // Tracks violations per proposal
	votingViolations   map[string]int // Tracks voting rule violations per voter
}

// NewDAOGovernanceRuleEnforcementAutomation initializes the DAO governance rule enforcement automation
func NewDAOGovernanceRuleEnforcementAutomation(governanceManager *governance.GovernanceManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DAOGovernanceRuleEnforcementAutomation {
	return &DAOGovernanceRuleEnforcementAutomation{
		governanceManager: governanceManager,
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		enforcementMutex:  enforcementMutex,
		proposalViolations: make(map[string]int),
		votingViolations:   make(map[string]int),
	}
}

// StartGovernanceEnforcement begins continuous monitoring and enforcement of DAO governance rules
func (automation *DAOGovernanceRuleEnforcementAutomation) StartGovernanceEnforcement() {
	ticker := time.NewTicker(GovernanceCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkGovernanceCompliance()
		}
	}()
}

// checkGovernanceCompliance monitors DAO proposals and votes for compliance with governance rules
func (automation *DAOGovernanceRuleEnforcementAutomation) checkGovernanceCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Check compliance for ongoing proposals
	for _, proposalID := range automation.governanceManager.GetOngoingProposals() {
		automation.validateProposal(proposalID)
		automation.validateVotingCompliance(proposalID)
	}
}

// validateProposal checks if a proposal meets quorum and approval thresholds
func (automation *DAOGovernanceRuleEnforcementAutomation) validateProposal(proposalID string) {
	quorumMet := automation.governanceManager.GetQuorumPercentage(proposalID) >= MinimumQuorumPercentage
	approvalMet := automation.governanceManager.GetApprovalPercentage(proposalID) >= MinimumApprovalPercentage

	if !quorumMet || !approvalMet {
		fmt.Printf("Proposal %s does not meet quorum or approval thresholds.\n", proposalID)
		automation.handleProposalViolation(proposalID)
	}
}

// validateVotingCompliance ensures voters are following voting rules (e.g., no double voting)
func (automation *DAOGovernanceRuleEnforcementAutomation) validateVotingCompliance(proposalID string) {
	for _, voterID := range automation.governanceManager.GetVoters(proposalID) {
		if automation.governanceManager.HasVotingViolations(voterID, proposalID) {
			fmt.Printf("Voting rule violation detected for voter %s on proposal %s.\n", voterID, proposalID)
			automation.handleVotingViolation(voterID)
		}
	}
}

// handleProposalViolation restricts proposal submissions from entities with repeated violations
func (automation *DAOGovernanceRuleEnforcementAutomation) handleProposalViolation(proposalID string) {
	automation.proposalViolations[proposalID]++

	if automation.proposalViolations[proposalID] >= MaxProposalViolations {
		err := automation.governanceManager.RestrictProposalSubmission(proposalID)
		if err != nil {
			fmt.Printf("Failed to restrict proposal %s due to rule violations: %v\n", proposalID, err)
			automation.logGovernanceAction(proposalID, "Failed Proposal Restriction")
		} else {
			fmt.Printf("Proposal %s restricted due to repeated rule violations.\n", proposalID)
			automation.logGovernanceAction(proposalID, "Proposal Restricted for Rule Violations")
			automation.proposalViolations[proposalID] = 0
		}
	} else {
		automation.logGovernanceAction(proposalID, "Proposal Compliance Violation Detected")
	}
}

// handleVotingViolation restricts voting rights for voters with repeated rule violations
func (automation *DAOGovernanceRuleEnforcementAutomation) handleVotingViolation(voterID string) {
	automation.votingViolations[voterID]++

	if automation.votingViolations[voterID] >= MaxVotingViolations {
		err := automation.governanceManager.RestrictVotingRights(voterID)
		if err != nil {
			fmt.Printf("Failed to restrict voting rights for voter %s: %v\n", voterID, err)
			automation.logGovernanceAction(voterID, "Failed Voting Restriction")
		} else {
			fmt.Printf("Voting rights restricted for voter %s due to repeated rule violations.\n", voterID)
			automation.logGovernanceAction(voterID, "Voting Rights Restricted for Rule Violations")
			automation.votingViolations[voterID] = 0
		}
	} else {
		automation.logGovernanceAction(voterID, "Voting Compliance Violation Detected")
	}
}

// logGovernanceAction securely logs actions related to DAO governance rule enforcement
func (automation *DAOGovernanceRuleEnforcementAutomation) logGovernanceAction(entityID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Entity ID: %s", action, entityID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("governance-enforcement-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DAO Governance Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance enforcement action for entity %s: %v\n", entityID, err)
	} else {
		fmt.Println("Governance enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DAOGovernanceRuleEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualGovernanceCheck allows administrators to manually check compliance for a specific proposal
func (automation *DAOGovernanceRuleEnforcementAutomation) TriggerManualGovernanceCheck(proposalID string) {
	fmt.Printf("Manually triggering governance compliance check for proposal: %s\n", proposalID)

	automation.validateProposal(proposalID)
	automation.validateVotingCompliance(proposalID)
}
