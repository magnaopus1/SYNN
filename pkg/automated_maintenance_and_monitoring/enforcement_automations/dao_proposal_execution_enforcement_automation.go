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

// Configuration for DAO proposal execution enforcement
const (
	ExecutionCheckInterval     = 15 * time.Second // Interval to check for proposal execution readiness
	MinimumApprovalPercentage  = 51              // Minimum approval percentage required for execution
)

// DAOProposalExecutionEnforcementAutomation monitors and enforces execution of DAO-approved proposals
type DAOProposalExecutionEnforcementAutomation struct {
	governanceManager *governance.GovernanceManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
}

// NewDAOProposalExecutionEnforcementAutomation initializes the proposal execution enforcement automation
func NewDAOProposalExecutionEnforcementAutomation(governanceManager *governance.GovernanceManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DAOProposalExecutionEnforcementAutomation {
	return &DAOProposalExecutionEnforcementAutomation{
		governanceManager: governanceManager,
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		enforcementMutex:  enforcementMutex,
	}
}

// StartExecutionEnforcement begins continuous monitoring and enforcement for DAO proposal execution
func (automation *DAOProposalExecutionEnforcementAutomation) StartExecutionEnforcement() {
	ticker := time.NewTicker(ExecutionCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkProposalExecution()
		}
	}()
}

// checkProposalExecution identifies approved proposals ready for execution and enforces their processing
func (automation *DAOProposalExecutionEnforcementAutomation) checkProposalExecution() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, proposalID := range automation.governanceManager.GetApprovedProposals() {
		if automation.validateProposalApproval(proposalID) {
			automation.executeProposal(proposalID)
		}
	}
}

// validateProposalApproval ensures the proposal meets the minimum approval percentage before execution
func (automation *DAOProposalExecutionEnforcementAutomation) validateProposalApproval(proposalID string) bool {
	approvalMet := automation.governanceManager.GetApprovalPercentage(proposalID) >= MinimumApprovalPercentage

	if !approvalMet {
		fmt.Printf("Proposal %s does not meet the approval threshold for execution.\n", proposalID)
		automation.logExecutionAction(proposalID, "Execution Denied: Insufficient Approval")
		return false
	}
	return true
}

// executeProposal carries out the execution of a DAO-approved proposal if all conditions are met
func (automation *DAOProposalExecutionEnforcementAutomation) executeProposal(proposalID string) {
	err := automation.consensusEngine.ValidateExecution(proposalID)
	if err != nil {
		fmt.Printf("Consensus validation failed for proposal %s execution: %v\n", proposalID, err)
		automation.logExecutionAction(proposalID, "Execution Failed: Consensus Validation Error")
		return
	}

	err = automation.governanceManager.ExecuteProposal(proposalID)
	if err != nil {
		fmt.Printf("Failed to execute proposal %s: %v\n", proposalID, err)
		automation.logExecutionAction(proposalID, "Execution Failed")
	} else {
		fmt.Printf("Successfully executed proposal %s.\n", proposalID)
		automation.logExecutionAction(proposalID, "Proposal Executed Successfully")
	}
}

// logExecutionAction securely logs actions related to DAO proposal execution enforcement
func (automation *DAOProposalExecutionEnforcementAutomation) logExecutionAction(proposalID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Proposal ID: %s", action, proposalID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("proposal-execution-%s-%d", proposalID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DAO Proposal Execution Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log proposal execution action for proposal %s: %v\n", proposalID, err)
	} else {
		fmt.Println("Proposal execution action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DAOProposalExecutionEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualExecution allows administrators to manually trigger the execution of a proposal if all requirements are met
func (automation *DAOProposalExecutionEnforcementAutomation) TriggerManualExecution(proposalID string) {
	fmt.Printf("Manually triggering execution for proposal: %s\n", proposalID)

	if automation.validateProposalApproval(proposalID) {
		automation.executeProposal(proposalID)
	} else {
		fmt.Printf("Manual execution denied for proposal %s due to insufficient approval.\n", proposalID)
		automation.logExecutionAction(proposalID, "Manual Execution Denied")
	}
}
