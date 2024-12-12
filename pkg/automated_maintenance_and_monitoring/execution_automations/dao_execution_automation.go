package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/dao"
)

const (
	DAOProposalCheckInterval      = 10 * time.Minute // Interval for checking new DAO proposals
	DAOExecutionTimeout           = 24 * time.Hour   // Timeout for proposal execution after approval
	DAOVoteThresholdPercentage    = 0.6              // Approval threshold (60%) for DAO proposal execution
)

// DAOExecutionAutomation manages decentralized autonomous organization (DAO) proposal executions
type DAOExecutionAutomation struct {
	consensusEngine    *synnergy_consensus.SynnergyConsensus // Synnergy Consensus for proposal validation
	ledgerInstance     *ledger.Ledger                        // Ledger for logging DAO proposal executions
	daoManager         *dao.DAOManager                       // DAO manager for handling proposals and votes
	executionMutex     *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewDAOExecutionAutomation initializes DAO execution automation
func NewDAOExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, daoManager *dao.DAOManager, executionMutex *sync.RWMutex) *DAOExecutionAutomation {
	return &DAOExecutionAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		daoManager:       daoManager,
		executionMutex:   executionMutex,
	}
}

// StartDAOExecutionMonitor begins monitoring DAO proposals and automates execution
func (automation *DAOExecutionAutomation) StartDAOExecutionMonitor() {
	ticker := time.NewTicker(DAOProposalCheckInterval)

	go func() {
		for range ticker.C {
			automation.processPendingDAOProposals()
		}
	}()
}

// processPendingDAOProposals retrieves and processes any pending DAO proposals that meet the criteria for execution
func (automation *DAOExecutionAutomation) processPendingDAOProposals() {
	automation.executionMutex.Lock()
	defer automation.executionMutex.Unlock()

	// Fetch pending proposals from the DAO manager
	pendingProposals, err := automation.daoManager.FetchPendingProposals()
	if err != nil {
		fmt.Println("Error fetching pending DAO proposals:", err)
		return
	}

	for _, proposal := range pendingProposals {
		// Validate the proposal with Synnergy Consensus
		if automation.validateDAOProposal(proposal) {
			automation.executeAndLogProposal(proposal)
		} else {
			fmt.Printf("DAO proposal %s failed consensus validation or did not meet vote threshold.\n", proposal.ID)
		}
	}
}

// validateDAOProposal validates the DAO proposal using Synnergy Consensus and checks voting thresholds
func (automation *DAOExecutionAutomation) validateDAOProposal(proposal dao.DAOProposal) bool {
	// Validate with Synnergy Consensus
	isValid := automation.consensusEngine.ValidateDAOProposal(proposal)

	// Check if proposal has reached the required voting threshold
	if !isValid || proposal.GetVotePercentage() < DAOVoteThresholdPercentage {
		return false
	}

	return true
}

// executeAndLogProposal executes the validated DAO proposal and logs the event in the ledger
func (automation *DAOExecutionAutomation) executeAndLogProposal(proposal dao.DAOProposal) {
	// Execute the DAO proposal
	err := automation.daoManager.ExecuteProposal(proposal)
	if err != nil {
		fmt.Printf("Error executing DAO proposal %s: %v\n", proposal.ID, err)
		return
	}

	// Log the DAO execution in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("dao-execution-%s", proposal.ID),
		Timestamp: time.Now().Unix(),
		Type:      "DAO Execution",
		Status:    "Completed",
		Details:   fmt.Sprintf("DAO proposal %s successfully executed.", proposal.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log DAO proposal %s execution: %v\n", proposal.ID, err)
	} else {
		fmt.Printf("DAO proposal %s executed and logged successfully.\n", proposal.ID)
	}
}

// encryptData encrypts sensitive data before logging it in the ledger
func (automation *DAOExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
