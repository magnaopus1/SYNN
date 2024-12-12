package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	GovernanceCheckInterval       = 5 * time.Minute // Interval for checking new proposals
	ProposalApprovalThreshold     = 0.6             // Threshold of votes required for approval (60%)
	ProposalTimeout               = 24 * time.Hour  // Timeout for proposals to be voted on
)

// GovernanceProposalAutomation manages the execution and processing of governance proposals
type GovernanceProposalAutomation struct {
	consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for proposal voting
	ledgerInstance  *ledger.Ledger                        // Ledger to track proposal execution
	proposalMutex   *sync.RWMutex                         // Mutex for thread-safe proposal execution
}

// NewGovernanceProposalAutomation initializes the governance proposal automation
func NewGovernanceProposalAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, proposalMutex *sync.RWMutex) *GovernanceProposalAutomation {
	return &GovernanceProposalAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		proposalMutex:   proposalMutex,
	}
}

// StartGovernanceProposalMonitoring begins the governance proposal monitoring process
func (automation *GovernanceProposalAutomation) StartGovernanceProposalMonitoring() {
	ticker := time.NewTicker(GovernanceCheckInterval)
	go func() {
		for range ticker.C {
			automation.evaluateGovernanceProposals()
		}
	}()
}

// evaluateGovernanceProposals checks for pending proposals and validates them for execution
func (automation *GovernanceProposalAutomation) evaluateGovernanceProposals() {
	automation.proposalMutex.Lock()
	defer automation.proposalMutex.Unlock()

	proposals := automation.consensusEngine.GetPendingProposals()

	for _, proposal := range proposals {
		if proposal.IsExpired(ProposalTimeout) {
			automation.rejectProposal(proposal.ID, "Proposal expired")
		} else if proposal.VoteCount() >= ProposalApprovalThreshold {
			automation.approveProposal(proposal.ID)
		}
	}
}

// approveProposal approves the governance proposal and executes the changes
func (automation *GovernanceProposalAutomation) approveProposal(proposalID string) {
	fmt.Printf("Approving and executing governance proposal %s...\n", proposalID)

	err := automation.consensusEngine.ExecuteProposal(proposalID)
	if err != nil {
		fmt.Printf("Failed to execute governance proposal %s: %v\n", proposalID, err)
		return
	}

	// Log the successful proposal execution in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("proposal-execution-%s-%d", proposalID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Governance Proposal Execution",
		Status:    "Success",
		Details:   fmt.Sprintf("Governance proposal %s successfully executed.", proposalID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance proposal execution for proposal %s: %v\n", proposalID, err)
	} else {
		fmt.Printf("Governance proposal %s execution successfully logged in the ledger.\n", proposalID)
	}
}

// rejectProposal rejects a governance proposal and logs the reason for rejection
func (automation *GovernanceProposalAutomation) rejectProposal(proposalID string, reason string) {
	fmt.Printf("Rejecting governance proposal %s. Reason: %s\n", proposalID, reason)

	// Log the proposal rejection in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("proposal-rejection-%s-%d", proposalID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Governance Proposal Rejection",
		Status:    "Rejected",
		Details:   fmt.Sprintf("Governance proposal %s was rejected. Reason: %s", proposalID, reason),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log governance proposal rejection for proposal %s: %v\n", proposalID, err)
	} else {
		fmt.Printf("Governance proposal %s rejection successfully logged in the ledger.\n", proposalID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *GovernanceProposalAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualProposalExecution allows administrators to manually trigger governance proposal execution
func (automation *GovernanceProposalAutomation) TriggerManualProposalExecution(proposalID string) {
	fmt.Printf("Manually triggering execution for governance proposal %s...\n", proposalID)
	automation.approveProposal(proposalID)
}

// TriggerManualProposalRejection allows administrators to manually reject a governance proposal
func (automation *GovernanceProposalAutomation) TriggerManualProposalRejection(proposalID string, reason string) {
	fmt.Printf("Manually triggering rejection for governance proposal %s. Reason: %s\n", proposalID, reason)
	automation.rejectProposal(proposalID, reason)
}
