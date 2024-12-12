package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/contracts"
)

const (
	ContractFinalizationCheckInterval = 10 * time.Minute // Interval for checking contract finalization status
)

// ContractFinalizationAutomation automates the finalization of smart contracts, legal smart contracts, and Ricardian contracts.
type ContractFinalizationAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validating contract finalization
	ledgerInstance    *ledger.Ledger                        // Ledger for logging finalization events
	contractManager   *contracts.ContractManager            // Manages smart contracts and their statuses
	stateMutex        *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewContractFinalizationAutomation initializes the contract finalization automation
func NewContractFinalizationAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, contractManager *contracts.ContractManager, stateMutex *sync.RWMutex) *ContractFinalizationAutomation {
	return &ContractFinalizationAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		contractManager:  contractManager,
		stateMutex:       stateMutex,
	}
}

// StartContractFinalizationMonitor starts monitoring the contract finalization process
func (automation *ContractFinalizationAutomation) StartContractFinalizationMonitor() {
	ticker := time.NewTicker(ContractFinalizationCheckInterval)

	go func() {
		for range ticker.C {
			automation.finalizePendingContracts()
		}
	}()
}

// finalizePendingContracts checks for contracts ready for finalization and processes them
func (automation *ContractFinalizationAutomation) finalizePendingContracts() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch pending contracts that are ready for finalization
	pendingContracts, err := automation.contractManager.FetchPendingContracts()
	if err != nil {
		fmt.Println("Error fetching pending contracts:", err)
		return
	}

	for _, contract := range pendingContracts {
		if automation.consensusEngine.ValidateContract(contract) {
			automation.finalizeContract(contract)
		} else {
			fmt.Printf("Contract %s failed consensus validation.\n", contract.ID)
		}
	}
}

// finalizeContract completes the finalization process for a given contract and logs it into the ledger
func (automation *ContractFinalizationAutomation) finalizeContract(contract contracts.Contract) {
	// Finalize the contract within the contract manager
	err := automation.contractManager.FinalizeContract(contract)
	if err != nil {
		fmt.Printf("Error finalizing contract %s: %v\n", contract.ID, err)
		return
	}

	// Log contract finalization into the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("contract-finalization-%s", contract.ID),
		Timestamp: time.Now().Unix(),
		Type:      "Contract Finalization",
		Status:    "Completed",
		Details:   fmt.Sprintf("Contract %s finalized successfully.", contract.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log contract finalization for contract %s: %v\n", contract.ID, err)
	} else {
		fmt.Printf("Contract %s finalized and logged successfully.\n", contract.ID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ContractFinalizationAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

