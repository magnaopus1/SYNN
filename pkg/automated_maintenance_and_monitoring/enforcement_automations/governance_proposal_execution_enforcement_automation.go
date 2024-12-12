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

// Configuration for governance proposal execution enforcement automation
const (
	ProposalExecutionCheckInterval = 15 * time.Second // Interval to check for execution-ready proposals
)

// GovernanceProposalExecutionEnforcementAutomation monitors and enforces the execution of approved governance proposals
type GovernanceProposalExecutionEnforcementAutomation struct {
	governanceManager   *governance.GovernanceManager
	consensusEngine     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	enforcementMutex    *sync.RWMutex
	proposalExecutionMap map[string]string // Tracks execution status of each proposal
}

// NewGovernanceProposalExecutionEnforcementAutomation initializes the governance proposal execution enforcement automation
func NewGovernanceProposalExecutionEnforcementAutomation(governanceManager *governance.GovernanceManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *GovernanceProposalExecutionEnforcementAutomation {
	return &GovernanceProposalExecutionEnforcementAutomation{
		governanceManager:    governanceManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		proposalExecutionMap: make(map[string]string),
	}
}

// StartProposalExecutionEnforcement begins continuous monitoring and enforcement of proposal executions
func (automation *GovernanceProposalExecutionEnforcementAutomation) StartProposalExecutionEnforcement() {
	ticker := time.NewTicker(ProposalExecutionCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkProposalExecutionCompliance()
		}
	}()
}

// checkProposalExecutionCompliance monitors each approved proposal and triggers execution if ready
func (automation *GovernanceProposalExecutionEnforcementAutomation) checkProposalExecutionCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, proposalID := range automation.governanceManager.GetExecutionReadyProposals() {
		if automation.proposalExecutionMap[proposalID] != "Executed" {
			automation.executeProposal(proposalID)
		}
	}
}

// executeProposal triggers the execution of a governance proposal and logs the action
func (automation *GovernanceProposalExecutionEnforcementAutomation) executeProposal(proposalID string) {
	err := automation.governanceManager.ExecuteProposal(proposalID)
	if err != nil {
		fmt.Printf("Failed to execute proposal %s: %v\n", proposalID, err)
		automation.logProposalAction(proposalID, "Execution Failed")
	} else {
		fmt.Printf("Proposal %s executed successfully.\n", proposalID)
		automation.proposalExecutionMap[proposalID] = "Executed"
		automation.logProposalAction(proposalID, "Executed")
	}
}

// logProposalAction securely logs actions related to governance proposal execution enforcement
func (automation *GovernanceProposalExecutionEnforcementAutomation) logProposalAction(proposalID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Proposal ID: %s", action, proposalID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("proposal-execution-enforcement-%s-%d", proposalID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Proposal Execution Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log proposal execution enforcement action for proposal %s: %v\n", proposalID, err)
	} else {
		fmt.Println("Proposal execution enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *GovernanceProposalExecutionEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualProposalExecution allows administrators to manually enforce the execution of a specific proposal
func (automation *GovernanceProposalExecutionEnforcementAutomation) TriggerManualProposalExecution(proposalID string) {
	fmt.Printf("Manually triggering execution for proposal: %s\n", proposalID)

	if automation.governanceManager.GetProposalStatus(proposalID) == "Execution Ready" {
		automation.executeProposal(proposalID)
	} else {
		fmt.Printf("Proposal %s is not ready for execution.\n", proposalID)
		automation.logProposalAction(proposalID, "Manual Execution Attempt Failed")
	}
}
