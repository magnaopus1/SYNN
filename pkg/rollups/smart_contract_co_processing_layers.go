package rollups

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewSmartContractCoProcessingLayer initializes a new Smart Contract Co-Processing Layer (SCCL)
func NewSmartContractCoProcessingLayer(layerID, rollupID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.SmartContractCoProcessingLayer {
	return &common.SmartContractCoProcessingLayer{
		LayerID:        layerID,
		RollupID:       rollupID,
		Contracts:      []*common.SmartContract{},
		Results:        make(map[string]interface{}),
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		Consensus:      consensus,
	}
}

// AddContract adds a new smart contract to the co-processing layer
func (sccl *common.SmartContractCoProcessingLayer) AddContract(contract *common.SmartContract) error {
	sccl.mu.Lock()
	defer sccl.mu.Unlock()

	if sccl.IsFinalized {
		return errors.New("co-processing layer is already finalized, no new contracts can be added")
	}

	// Encrypt the contract data
	encryptedContract, err := sccl.Encryption.EncryptData([]byte(contract.ContractID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt contract: %v", err)
	}
	contract.ContractID = string(encryptedContract)

	// Add the contract to the layer
	sccl.Contracts = append(sccl.Contracts, contract)

	// Log the contract addition in the ledger
	err = sccl.Ledger.RecordContractAddition(sccl.LayerID, contract.ContractID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract addition: %v", err)
	}

	fmt.Printf("Smart contract %s added to co-processing layer %s\n", contract.ContractID, sccl.LayerID)
	return nil
}

// FinalizeLayer finalizes the co-processing layer, marking it ready for result syncing
func (sccl *common.SmartContractCoProcessingLayer) FinalizeLayer() error {
	sccl.mu.Lock()
	defer sccl.mu.Unlock()

	if sccl.IsFinalized {
		return errors.New("co-processing layer is already finalized")
	}

	// Process each contract off-chain
	for _, contract := range sccl.Contracts {
		// Execute the smart contract and store the result
		result, err := contract.ExecuteOffChain()
		if err != nil {
			return fmt.Errorf("failed to execute contract %s: %v", contract.ContractID, err)
		}

		// Store the result, assuming contract.ContractID as key
		sccl.Results[contract.ContractID] = result

		// Log the result in the ledger
		err = sccl.Ledger.RecordContractExecution(sccl.LayerID, contract.ContractID, result, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log contract execution: %v", err)
		}
	}

	sccl.IsFinalized = true
	fmt.Printf("Smart contract co-processing layer %s finalized\n", sccl.LayerID)
	return nil
}

// SyncResults synchronizes the co-processed smart contract results back to the rollup
func (sccl *common.SmartContractCoProcessingLayer) SyncResults() error {
	sccl.mu.Lock()
	defer sccl.mu.Unlock()

	if !sccl.IsFinalized {
		return errors.New("co-processing layer is not finalized, cannot sync results")
	}

	// Sync results to the rollup via the network
	for contractID, result := range sccl.Results {
		err := sccl.NetworkManager.BroadcastData(sccl.RollupID, []byte(fmt.Sprintf("ContractID: %s, Result: %v", contractID, result)))
		if err != nil {
			return fmt.Errorf("failed to sync contract result %s: %v", contractID, err)
		}

		// Log the sync in the ledger
		err = sccl.Ledger.RecordResultSync(sccl.LayerID, contractID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log result sync: %v", err)
		}
	}

	fmt.Printf("Results synced from co-processing layer %s to rollup %s\n", sccl.LayerID, sccl.RollupID)
	return nil
}

// VerifyResults uses Synnergy Consensus to verify the co-processed results
func (sccl *common.SmartContractCoProcessingLayer) VerifyResults() (bool, error) {
	sccl.mu.Lock()
	defer sccl.mu.Unlock()

	// Verify the co-processed results using Synnergy Consensus
	valid, err := sccl.Consensus.VerifyResults(sccl.LayerID, sccl.Results)
	if err != nil {
		return false, fmt.Errorf("failed to verify co-processing layer results: %v", err)
	}

	// Log the verification in the ledger
	err = sccl.Ledger.RecordResultVerification(sccl.LayerID, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to log result verification: %v", err)
	}

	fmt.Printf("Co-processing layer %s results verified successfully\n", sccl.LayerID)
	return valid, nil
}

// RetrieveResult retrieves the result of a specific contract execution from the co-processing layer
func (sccl *common.SmartContractCoProcessingLayer) RetrieveResult(contractID string) (interface{}, error) {
	sccl.mu.Lock()
	defer sccl.mu.Unlock()

	result, exists := sccl.Results[contractID]
	if !exists {
		return nil, fmt.Errorf("result for contract %s not found in co-processing layer %s", contractID, sccl.LayerID)
	}

	fmt.Printf("Retrieved result for contract %s from co-processing layer %s\n", contractID, sccl.LayerID)
	return result, nil
}
