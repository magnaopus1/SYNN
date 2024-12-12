package ledger

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// RecordContractStorage stores data in a contract's storage
func (l *SmartContractLedger ) RecordContractStorage(contractAddress, key, value string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.ContractStorage[contractAddress]; !exists {
		l.ContractStorage[contractAddress] = ContractStorage{
			ContractAddress: contractAddress,
			Data:            make(map[string]string),
			Timestamp:       time.Now(),
		}
	}

	l.ContractStorage[contractAddress].Data[key] = value
	return nil
}

// RecordContractInteraction logs an interaction with a smart contract.
func (l *SmartContractLedger) RecordContractInteraction(initiatingContractID, receivingContractID, caller, functionName string, parameters map[string]interface{}, executionResult, interactionData, encryptedInteractionData string) error {
	l.Lock()
	defer l.Unlock()

	// Create a unique Interaction ID (you may use a more robust method to generate this)
	interactionID := fmt.Sprintf("interaction-%d", time.Now().UnixNano())

	interaction := ContractInteraction{
		InteractionID:            interactionID,
		InitiatingContractID:      initiatingContractID,
		ReceivingContractID:       receivingContractID,
		Caller:                    caller,
		FunctionName:              functionName,
		Parameters:                parameters,
		Timestamp:                 time.Now(),
		ExecutionResult:           executionResult,
		InteractionData:           interactionData,
		EncryptedInteractionData:  encryptedInteractionData,
	}

	// Initialize if not present
	if _, exists := l.ContractInteractions[initiatingContractID]; !exists {
		l.ContractInteractions[initiatingContractID] = []ContractInteraction{}
	}

	// Append the interaction to the contract's interaction log
	l.ContractInteractions[initiatingContractID] = append(l.ContractInteractions[initiatingContractID], interaction)

	return nil
}


// RetrieveContractInteraction retrieves interactions for a given contract address
func (l *SmartContractLedger) RetrieveContractInteraction(contractAddress string) ([]ContractInteraction, error) {
	l.Lock()
	defer l.Unlock()

	interactions, exists := l.ContractInteractions[contractAddress]
	if !exists {
		return nil, errors.New("no interactions found for the contract")
	}

	return interactions, nil
}

// LogDeployment logs the deployment of a cross-chain smart contract.
func LogDeployment(contractID, owner string, connectedChains string) {
	timestamp := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf(
		"[%s] Contract Deployment: ID=%s, Owner=%s, ConnectedChains=%s",
		timestamp, contractID, owner, connectedChains,
	)
	log.Println(logMessage)
}

// RemoveContractStorage removes key-value data from the contract's storage
func (l *SmartContractLedger) RemoveContractStorage(contractAddress, key string) error {
	l.Lock()
	defer l.Unlock()

	if contractStorage, exists := l.ContractStorage[contractAddress]; exists {
		delete(contractStorage.Data, key)
		return nil
	}

	return errors.New("contract storage not found")
}

// RetrieveContractStorage retrieves the stored data for a given contract
func (l *SmartContractLedger) RetrieveContractStorage(contractAddress, key string) (string, error) {
	l.Lock()
	defer l.Unlock()

	if contractStorage, exists := l.ContractStorage[contractAddress]; exists {
		value, exists := contractStorage.Data[key]
		if !exists {
			return "", errors.New("key not found in contract storage")
		}
		return value, nil
	}

	return "", errors.New("contract storage not found")
}

// LogExecutionFailure logs contract execution failures
func LogExecutionFailure(contractID, reason string) {
	timestamp := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf("[%s] Execution Failure: ContractID=%s, Reason=%s", timestamp, contractID, reason)
	log.Println(logMessage)
}

// RetrieveContractDeployment retrieves the deployment details for a contract.
func (l *SmartContractLedger) RetrieveContractDeployment(contractID string) (ContractDeployment, error) {
	l.Lock()
	defer l.Unlock()

	deployment, exists := l.ContractDeployments[contractID]
	if !exists {
		return ContractDeployment{}, errors.New("contract deployment not found")
	}

	return deployment, nil
}

// LogExecution logs the successful execution of a contract.
func LogExecution(executionID, contractID string) {
	timestamp := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf(
		"[%s] Contract Execution: ExecutionID=%s, ContractID=%s",
		timestamp, executionID, contractID,
	)
	log.Println(logMessage)
}

// RecordContractDeployment logs the deployment of a smart contract
func (l *SmartContractLedger) RecordContractDeployment(contractAddress, deployerAddress, code string) error {
    l.Lock()
    defer l.Unlock()

    deployment := ContractDeployment{
        ContractID:     contractAddress,    // Matching the struct's field name
        Deployer:       deployerAddress,    // Matching the struct's field name
        ContractCode:   code,               // Matching the struct's field name
        DeployedAt:     time.Now(),         // Matching the struct's field name
        Status:         "open",             // Set initial status to "open"
    }

    l.ContractDeployments[contractAddress] = deployment
    return nil
}

// RetrieveContract retrieves the deployment details of a smart contract
func (l *SmartContractLedger) RetrieveContract(contractAddress string) (*ContractDeployment, error) {
	l.Lock()
	defer l.Unlock()

	deployment, exists := l.ContractDeployments[contractAddress]
	if !exists {
		return nil, errors.New("contract not found")
	}

	return &deployment, nil
}

// RecordMigration logs the migration of a smart contract
func (l *SmartContractLedger) RecordMigration(oldContractAddress, newContractAddress, migratorAddress string) error {
	l.Lock()
	defer l.Unlock()

	migration := ContractMigration{
		OldContractAddress: oldContractAddress,
		NewContractAddress: newContractAddress,
		MigratorAddress:    migratorAddress,
		Timestamp:          time.Now(),
	}

	l.ContractMigrations[oldContractAddress] = migration
	return nil
}

// RecordContractSignature logs a contract signature
func (l *SmartContractLedger) RecordContractSignature(contractAddress, signerAddress, signature string) error {
	l.Lock()
	defer l.Unlock()

	signatureRecord := ContractSignature{
		ContractAddress: contractAddress,
		SignerAddress:   signerAddress,
		Signature:       signature,
		Timestamp:       time.Now(),
	}

	l.ContractSignatures[contractAddress] = append(l.ContractSignatures[contractAddress], signatureRecord)
	return nil
}



// RetrieveMigrationRecord retrieves the MigrationID of a migration record for a given original and new contract ID.
func (l *SmartContractLedger) RetrieveMigrationRecord(originalContractID, newContractID string) (string, error) {
    key := fmt.Sprintf("%s-%s", originalContractID, newContractID)

    l.Lock()
    defer l.Unlock()

    // Check if the record exists in the migration records map
    if record, exists := l.MigrationRecords[key]; exists {
        return record.MigrationID, nil // Use MigrationID as the unique identifier or "hash"
    }

    return "", fmt.Errorf("migration record not found for original contract ID: %s, new contract ID: %s", originalContractID, newContractID)
}



// RecordContractDeployment logs the deployment of a smart contract in a state channel.
func (l *SmartContractLedger) RecordStateChannelContractDeployment(contractID, deployer, contractCode string) error {
	l.Lock()
	defer l.Unlock()

	// Log contract deployment
	deployment := ContractDeployment{
		ContractID:   contractID,
		Deployer:     deployer,
		ContractCode: contractCode,
		DeployedAt:   time.Now(),
	}

	// Store the deployment in the ledger
	l.ContractDeployments[contractID] = deployment
	return nil
}

// RecordContractClosure logs the closure of a smart contract.
func (l *SmartContractLedger) RecordContractClosure(contractID, closer string) error {
	l.Lock()
	defer l.Unlock()

	contract, exists := l.ContractDeployments[contractID]
	if !exists {
		return errors.New("contract not found")
	}

	// Mark the contract as closed
	contract.Status = "closed"
	contract.ClosedBy = closer
	contract.ClosedAt = time.Now()

	// Update the contract in the ledger
	l.ContractDeployments[contractID] = contract
	return nil
}
