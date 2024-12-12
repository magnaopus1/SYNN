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
	"synnergy_network_demo/cross_chain"
)

const (
	CrossChainCheckInterval = 5 * time.Minute // Interval for checking cross-chain contract status
	CrossChainFinalizationTimeout = 30 * time.Minute // Timeout duration for contract finalization
)

// CrossChainContractExecutionAutomation handles the execution and finalization of smart contracts across different chains.
type CrossChainContractExecutionAutomation struct {
	consensusEngine        *synnergy_consensus.SynnergyConsensus // Consensus engine for cross-chain validation
	ledgerInstance         *ledger.Ledger                        // Ledger for logging cross-chain contract events
	contractManager        *contracts.ContractManager            // Manages cross-chain smart contracts
	crossChainManager      *cross_chain.CrossChainManager        // Cross-chain manager to facilitate cross-chain interaction
	stateMutex             *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewCrossChainContractExecutionAutomation initializes cross-chain contract execution automation
func NewCrossChainContractExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, contractManager *contracts.ContractManager, crossChainManager *cross_chain.CrossChainManager, stateMutex *sync.RWMutex) *CrossChainContractExecutionAutomation {
	return &CrossChainContractExecutionAutomation{
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		contractManager:    contractManager,
		crossChainManager:  crossChainManager,
		stateMutex:         stateMutex,
	}
}

// StartCrossChainExecutionMonitor begins continuous monitoring and execution of cross-chain smart contracts
func (automation *CrossChainContractExecutionAutomation) StartCrossChainExecutionMonitor() {
	ticker := time.NewTicker(CrossChainCheckInterval)

	go func() {
		for range ticker.C {
			automation.processPendingCrossChainContracts()
		}
	}()
}

// processPendingCrossChainContracts checks for cross-chain contracts that need execution or finalization
func (automation *CrossChainContractExecutionAutomation) processPendingCrossChainContracts() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch pending cross-chain contracts
	pendingContracts, err := automation.contractManager.FetchPendingCrossChainContracts()
	if err != nil {
		fmt.Println("Error fetching pending cross-chain contracts:", err)
		return
	}

	for _, contract := range pendingContracts {
		// Validate the contract across chains
		if automation.validateCrossChainContract(contract) {
			automation.executeAndFinalizeContract(contract)
		} else {
			fmt.Printf("Cross-chain contract %s failed consensus validation.\n", contract.ID)
		}
	}
}

// validateCrossChainContract ensures the contract passes the consensus validation across multiple chains
func (automation *CrossChainContractExecutionAutomation) validateCrossChainContract(contract contracts.Contract) bool {
	// Validate contract within Synnergy Consensus
	isValidInPrimaryChain := automation.consensusEngine.ValidateContract(contract)

	// Check with the cross-chain manager for validation on other chains
	isValidInSecondaryChains := automation.crossChainManager.ValidateContractAcrossChains(contract)

	return isValidInPrimaryChain && isValidInSecondaryChains
}

// executeAndFinalizeContract completes the execution and finalization of the contract and logs it into the ledger
func (automation *CrossChainContractExecutionAutomation) executeAndFinalizeContract(contract contracts.Contract) {
	// Execute contract on all chains involved
	err := automation.crossChainManager.ExecuteCrossChainContract(contract)
	if err != nil {
		fmt.Printf("Error executing cross-chain contract %s: %v\n", contract.ID, err)
		return
	}

	// Finalize the contract once executed
	err = automation.contractManager.FinalizeContract(contract)
	if err != nil {
		fmt.Printf("Error finalizing cross-chain contract %s: %v\n", contract.ID, err)
		return
	}

	// Log contract finalization into the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("cross-chain-contract-finalization-%s", contract.ID),
		Timestamp: time.Now().Unix(),
		Type:      "Cross-Chain Contract Finalization",
		Status:    "Completed",
		Details:   fmt.Sprintf("Cross-chain contract %s successfully executed and finalized.", contract.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log finalization for cross-chain contract %s: %v\n", contract.ID, err)
	} else {
		fmt.Printf("Cross-chain contract %s finalized and logged successfully.\n", contract.ID)
	}
}

// encryptData encrypts sensitive contract data before logging it into the ledger
func (automation *CrossChainContractExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

