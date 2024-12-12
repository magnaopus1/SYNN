package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRollupContract initializes a new rollup contract
func NewRollupContract(contractID, contractOwner string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.RollupContract {
	return &common.RollupContract{
		ContractID:    contractID,
		ContractOwner: contractOwner,
		ContractState: make(map[string]interface{}),
		IsDeployed:    false,
		Ledger:        ledgerInstance,
		Encryption:    encryptionService,
		Consensus:     consensus,
	}
}

// DeployContract deploys a smart contract onto the rollup
func (rc *common.RollupContract) DeployContract(contractCode string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.IsDeployed {
		return errors.New("contract is already deployed")
	}

	// Encrypt contract code before deployment
	encryptedContractCode, err := rc.Encryption.EncryptData([]byte(contractCode), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt contract code: %v", err)
	}

	// Mark the contract as deployed
	rc.IsDeployed = true

	// Log the deployment event in the ledger
	err = rc.Ledger.RecordContractDeployment(rc.ContractID, rc.ContractOwner, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract deployment: %v", err)
	}

	fmt.Printf("Smart contract %s deployed by owner %s on rollup\n", rc.ContractID, rc.ContractOwner)
	return nil
}

// ExecuteContract handles the execution of the smart contract on the rollup
func (rc *common.RollupContract) ExecuteContract(transaction *common.Transaction) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if !rc.IsDeployed {
		return errors.New("contract is not deployed")
	}

	// Encrypt the transaction before execution
	encryptedTx, err := rc.Encryption.EncryptData([]byte(transaction.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	transaction.TxID = string(encryptedTx)

	// Append transaction to contract's transaction list
	rc.Transactions = append(rc.Transactions, transaction)

	// Log contract execution event in the ledger
	err = rc.Ledger.RecordContractExecution(rc.ContractID, transaction.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract execution: %v", err)
	}

	// Example: Update contract state (this would vary based on the smart contract logic)
	rc.ContractState["lastExecutedTx"] = transaction.TxID

	fmt.Printf("Smart contract %s executed with transaction %s\n", rc.ContractID, transaction.TxID)
	return nil
}

// ValidateContract ensures that the smart contract is valid using Synnergy Consensus
func (rc *common.RollupContract) ValidateContract() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if !rc.IsDeployed {
		return errors.New("contract is not deployed")
	}

	// Use consensus to validate the contract's integrity
	err := rc.Consensus.ValidateContract(rc.ContractID, rc.ContractState)
	if err != nil {
		return fmt.Errorf("contract validation failed: %v", err)
	}

	// Log the validation in the ledger
	err = rc.Ledger.RecordContractValidation(rc.ContractID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract validation: %v", err)
	}

	fmt.Printf("Smart contract %s validated successfully\n", rc.ContractID)
	return nil
}

// RetrieveContract retrieves the details of the deployed contract
func (rc *common.RollupContract) RetrieveContract() (*common.RollupContract, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if !rc.IsDeployed {
		return nil, errors.New("contract is not deployed")
	}

	fmt.Printf("Retrieved smart contract %s\n", rc.ContractID)
	return rc, nil
}

// UpdateContractState securely updates the state of the smart contract
func (rc *common.RollupContract) UpdateContractState(key string, value interface{}) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if !rc.IsDeployed {
		return errors.New("contract is not deployed")
	}

	// Update contract state
	rc.ContractState[key] = value

	// Log the state update in the ledger
	err := rc.Ledger.RecordContractStateUpdate(rc.ContractID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract state update: %v", err)
	}

	fmt.Printf("Contract %s state updated: %s = %v\n", rc.ContractID, key, value)
	return nil
}

// RetrieveContractState retrieves the current state of the contract
func (rc *common.RollupContract) RetrieveContractState(key string) (interface{}, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	value, exists := rc.ContractState[key]
	if !exists {
		return nil, fmt.Errorf("contract state key %s not found", key)
	}

	fmt.Printf("Retrieved contract state: %s = %v\n", key, value)
	return value, nil
}
