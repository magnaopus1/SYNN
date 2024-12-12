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

// Configuration for governance decision enforcement automation
const (
	GovernanceCheckInterval = 15 * time.Second // Interval to check governance decision compliance
)

// GovernanceDecisionEnforcementAutomation monitors and enforces execution of approved governance decisions
type GovernanceDecisionEnforcementAutomation struct {
	governanceManager    *governance.GovernanceManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	proposalStatusMap    map[string]string // Tracks status of each governance proposal
}

// NewGovernanceDecisionEnforcementAutomation initializes the governance decision enforcement automation
func NewGovernanceDecisionEnforcementAutomation(governanceManager *governance.GovernanceManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *GovernanceDecisionEnforcementAutomation {
	return &GovernanceDecisionEnforcementAutomation{
		governanceManager:    governanceManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		proposalStatusMap:    make(map[string]string),
	}
}

// StartGovernanceEnforcement begins continuous monitoring and enforcement of governance decisions
func (automation *GovernanceDecisionEnforcementAutomation) StartGovernanceEnforcement() {
	ticker := time.NewTicker(GovernanceCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkGovernanceCompliance()
		}
	}()
}

// checkGovernanceCompliance verifies the status of governance proposals and enforces decisions as needed
func (automation *GovernanceDecisionEnforcementAutomation) checkGovernanceCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, proposalID := range automation.governanceManager.GetApprovedProposals() {
		status := automation.governanceManager.GetProposalStatus(proposalID)

		if status == "Approved" {
			automation.executeGovernanceDecision(proposalID)
		}
	}
}

// executeGovernanceDecision executes the decision for an approved proposal and logs the action
func (automation *GovernanceDecisionEnforcementAutomation) executeGovernanceDecision(proposalID string) {
	err := automation.governanceManager.ExecuteProposal(proposalID)
	if err != nil {
		fmt.Printf("Failed to execute governance decision for proposal %s: %v\n", proposalID, err)
		automation.logGovernanceAction(proposalID, "Execution Failed")
	} else {
		fmt.Printf("Governance decision executed for proposal %s.\n", proposalID)
		automation.updateProposalStatus(proposalID, "Executed")
		automation.logGovernanceAction(proposalID, "Decision Executed")
	}
}

// updateProposalStatus updates the proposal status in the tracking map
func (automation *GovernanceDecisionEnforcementAutomation) updateProposalStatus(proposalID string, status string) {
	automation.proposalStatusMap[proposalID] = status
}

// logGovernanceAction securely logs actions related to governance decision enforcement
func (automation *GovernanceDecisionEnforcementAutomation) logGovernanceAction(proposalID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Proposal ID: %s", action, proposalID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("governance-enforcement-%s-%d", proposalID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Governance Decision Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance decision enforcement action for proposal %s: %v\n", proposalID, err)
	} else {
		fmt.Println("Governance decision enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *GovernanceDecisionEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualProposalExecution allows administrators to manually execute a specific governance proposal
func (automation *GovernanceDecisionEnforcementAutomation) TriggerManualProposalExecution(proposalID string) {
	fmt.Printf("Manually triggering governance decision execution for proposal: %s\n", proposalID)

	status := automation.governanceManager.GetProposalStatus(proposalID)
	if status == "Approved" {
		automation.executeGovernanceDecision(proposalID)
	} else {
		fmt.Printf("Proposal %s is not approved and cannot be executed.\n", proposalID)
		automation.logGovernanceAction(proposalID, "Manual Execution Attempt Failed")
	}
}
