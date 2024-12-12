package state_channels

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewSmartContractStateChannel initializes a new state channel with smart contract functionality
func NewSmartContractStateChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.SmartContractStateChannel {
	return &common.SmartContractStateChannel{
		ChannelID:      channelID,
		Participants:   participants,
		State:          make(map[string]interface{}),
		SmartContracts: make(map[string]*common.Contract),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
	}
}

// DeployContract deploys a smart contract into the state channel
func (sc *common.SmartContractStateChannel) DeployContract(contractID string, contractCode string, params map[string]interface{}) (*common.Contract, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return nil, errors.New("state channel is closed")
	}

	// Check if the contract already exists
	if _, exists := sc.SmartContracts[contractID]; exists {
		return nil, errors.New("contract already exists in the channel")
	}

	// Deploy the smart contract
	contract := &smart_contracts.Contract{
		ContractID: contractID,
		Code:       contractCode,
		Params:     params,
		CreatedAt:  time.Now(),
	}

	// Store the contract in the channel
	sc.SmartContracts[contractID] = contract

	// Log the contract deployment in the ledger
	err := sc.Ledger.RecordContractDeployment(sc.ChannelID, contractID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log contract deployment: %v", err)
	}

	fmt.Printf("Smart contract %s deployed in state channel %s\n", contractID, sc.ChannelID)
	return contract, nil
}

// ExecuteContract executes a deployed smart contract within the state channel
func (sc *SmartContractStateChannel) ExecuteContract(contractID string, function string, args map[string]interface{}) (map[string]interface{}, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Check if the contract exists
	contract, exists := sc.SmartContracts[contractID]
	if !exists {
		return nil, fmt.Errorf("contract %s not found in channel", contractID)
	}

	// Execute the contract's function
	result, err := contract.ExecuteFunction(function, args)
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	// Log the contract execution in the ledger
	err = sc.Ledger.RecordContractExecution(sc.ChannelID, contractID, function, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log contract execution: %v", err)
	}

	fmt.Printf("Smart contract %s executed function %s in state channel %s\n", contractID, function, sc.ChannelID)
	return result, nil
}

// CloseContract finalizes and closes a contract within the state channel
func (sc *SmartContractStateChannel) CloseContract(contractID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Check if the contract exists
	contract, exists := sc.SmartContracts[contractID]
	if !exists {
		return fmt.Errorf("contract %s not found in channel", contractID)
	}

	// Finalize the contract
	contract.Finalized = true

	// Log the contract closure in the ledger
	err := sc.Ledger.RecordContractClosure(sc.ChannelID, contractID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract closure: %v", err)
	}

	fmt.Printf("Smart contract %s closed in state channel %s\n", contractID, sc.ChannelID)
	return nil
}

// UpdateContractState securely updates the internal state of the contract
func (sc *SmartContractStateChannel) UpdateContractState(contractID string, key string, value interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Check if the contract exists
	contract, exists := sc.SmartContracts[contractID]
	if !exists {
		return fmt.Errorf("contract %s not found in channel", contractID)
	}

	// Update the state of the contract
	contract.State[key] = value

	// Log the state update in the ledger
	err := sc.Ledger.RecordStateUpdate(sc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log contract state update: %v", err)
	}

	fmt.Printf("State of contract %s updated: %s = %v in channel %s\n", contractID, key, value, sc.ChannelID)
	return nil
}

// RetrieveContractState retrieves the state of a smart contract in the channel
func (sc *SmartContractStateChannel) RetrieveContractState(contractID string, key string) (interface{}, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Check if the contract exists
	contract, exists := sc.SmartContracts[contractID]
	if !exists {
		return nil, fmt.Errorf("contract %s not found in channel", contractID)
	}

	value, exists := contract.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found in contract %s", key, contractID)
	}

	fmt.Printf("Retrieved state from contract %s: %s = %v\n", contractID, key, value)
	return value, nil
}

// RetrieveContract retrieves a deployed contract by its ID
func (sc *SmartContractStateChannel) RetrieveContract(contractID string) (*smart_contracts.Contract, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	contract, exists := sc.SmartContracts[contractID]
	if !exists {
		return nil, fmt.Errorf("contract %s not found in channel", contractID)
	}

	fmt.Printf("Retrieved contract %s from state channel %s\n", contractID, sc.ChannelID)
	return contract, nil
}

// CloseChannel closes the state channel and all associated contracts
func (sc *SmartContractStateChannel) CloseChannel() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if !sc.IsOpen {
		return errors.New("state channel is already closed")
	}

	// Finalize all contracts before closing the channel
	for contractID := range sc.SmartContracts {
		err := sc.CloseContract(contractID)
		if err != nil {
			return fmt.Errorf("failed to close contract %s: %v", contractID, err)
		}
	}

	sc.IsOpen = false

	// Log the channel closure in the ledger
	err := sc.Ledger.RecordChannelClosure(sc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("State channel %s closed\n", sc.ChannelID)
	return nil
}
