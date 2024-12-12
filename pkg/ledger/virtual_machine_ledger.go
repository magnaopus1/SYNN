package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordBytecodeDeployment records the deployment of bytecode to the ledger.
func (l *VirtualMachineLedger) RecordBytecodeDeployment(contractID string, code []byte) error {
	l.Lock()
	defer l.Unlock()

	// Check if the bytecode for the contract is already deployed
	if _, exists := l.BytecodeStore[contractID]; exists {
		return errors.New("bytecode already deployed for this contract")
	}

	bytecode := Bytecode{
		ContractID: contractID,
		Code:       code,
		DeployedAt: time.Now(),
	}

	// Store the deployed bytecode
	l.BytecodeStore[contractID] = bytecode
	return nil
}




// GetBytecode retrieves the bytecode of a deployed contract.
func (l *VirtualMachineLedger) GetBytecode(contractID string) ([]byte, error) {
	l.Lock()
	defer l.Unlock()

	bytecode, exists := l.BytecodeStore[contractID]
	if !exists {
		return nil, errors.New("bytecode not found for the contract")
	}

	return bytecode.Code, nil
}


// GetRecord retrieves the details of a specific contract record.
func (l *VirtualMachineLedger) GetRecord(contractID string) (ContractState, error) {
	l.Lock()
	defer l.Unlock()

	contractState, exists := l.ContractState[contractID]
	if !exists {
		return ContractState{}, errors.New("contract state not found")
	}

	return contractState, nil
}

// RecordContractExecution logs the execution of a contract and updates its state.
func (l *VirtualMachineLedger) RecordContractExecution(contractID string, stateUpdate map[string]interface{}) error {
	l.Lock()
	defer l.Unlock()

	contractState, exists := l.ContractState[contractID]
	if !exists {
		// Initialize a new state if the contract does not exist yet
		contractState = ContractState{
			ContractID: contractID,
			StateData:  make(map[string]interface{}),
			UpdatedAt:  time.Now(),
		}
	}

	// Update the contract's state
	for key, value := range stateUpdate {
		contractState.StateData[key] = value
	}
	contractState.UpdatedAt = time.Now()

	l.ContractState[contractID] = contractState
	return nil
}


// UpdateBalance updates the balance of an account after contract execution.
func (l *VirtualMachineLedger) UpdateBalance(A *AccountsWalletLedger, accountID string, newBalance uint64) error {
	l.Lock()
	defer l.Unlock()

	account, exists := A.AccountsWalletLedgerState.Accounts[accountID]
	if !exists {
		return errors.New("account not found")
	}

	// Convert uint64 newBalance to float64 before assignment
	account.Balance = float64(newBalance)
	A.AccountsWalletLedgerState.Accounts[accountID] = account
	return nil
}



// StoreContract stores the state of a contract in the ledger.
func (l *VirtualMachineLedger) StoreContract(contractID string, state map[string]interface{}) error {
	l.Lock()
	defer l.Unlock()

	contractState := ContractState{
		ContractID: contractID,
		StateData:  state,
		UpdatedAt:  time.Now(),
	}

	// Store or update the contract state
	l.ContractState[contractID] = contractState
	return nil
}

// GetContractState retrieves the state of a contract.
func (l *VirtualMachineLedger) GetContractState(contractID string) (map[string]interface{}, error) {
	l.Lock()
	defer l.Unlock()

	contractState, exists := l.ContractState[contractID]
	if !exists {
		return nil, errors.New("contract state not found")
	}

	return contractState.StateData, nil
}

func (l *VirtualMachineLedger) LogEntry(logType, description string, vmID string, severity string) error {
	l.Lock()
	defer l.Unlock()

	entry := VMLogEntry{
		LogID:       GenerateUniqueID(), // Generate unique ID for the log
		Timestamp:   time.Now(),
		VMID:        vmID,
		Event:       logType,
		Severity:    severity,
		Details:     description,
	}

	// Add log entry to VMLogEntries slice
	l.VMLogEntries = append(l.VMLogEntries, entry)
	return nil
}




// LogExecution logs the execution of a contract, recording the transaction.
func (l *VirtualMachineLedger) LogExecution(contractID, txID, logDetails string) error {
	return l.LogEntry(
		"contract_execution",                             // logType
		"Contract "+contractID+" executed in transaction "+txID+": "+logDetails, // description
		contractID,                                      // vmID (use contractID as a reference VM ID here)
		"info",                                          // severity
	)
}

// LogBlock logs an event related to a block in the ledger.
func (l *VirtualMachineLedger) LogBlock(blockID, logDetails string) error {
	return l.LogEntry(
		"block",                      // logType
		"Block "+blockID+": "+logDetails, // description
		blockID,                      // vmID (use blockID as a reference VM ID here)
		"info",                       // severity
	)
}


// Utility function for generating unique log entry IDs.
func GenerateUniqueID() string {
	return time.Now().Format("20060102150405") // Simple timestamp-based ID generation
}

// RecordScheduledEvent logs a scheduled event in the ledger with additional details.
func (l *VirtualMachineLedger) RecordScheduledEvent(eventID, vmID, initiator string, executionTime time.Time, opcode byte) {
    event := VMEventLog{
        EventID:     eventID,
        VMID:        vmID,
        EventType:   "Scheduled",
        Timestamp:   time.Now(),
        Initiator:   initiator,
        Description: fmt.Sprintf("Execution Time: %s, Opcode: %d", executionTime.Format(time.RFC3339), opcode),
        Impact:      "low", // Default impact level for scheduled events
    }
    l.VMEventLog = append(l.VMEventLog, event)
    fmt.Printf("Scheduled Event Recorded: %v\n", event)
}

// RecordFailedEvent logs a failed event with additional details in the ledger.
func (l *VirtualMachineLedger) RecordFailedEvent(eventID, vmID, initiator, errorMsg string) {
    event := VMEventLog{
        EventID:     eventID,
        VMID:        vmID,
        EventType:   "Failed",
        Timestamp:   time.Now(),
        Initiator:   initiator,
        Description: fmt.Sprintf("Error: %s", errorMsg),
        Impact:      "high", // Default impact level for failed events
    }
    l.VMEventLog = append(l.VMEventLog, event)
    fmt.Printf("Failed Event Recorded: %v\n", event)
}


// RecordExecutedEvent logs an executed event in the ledger.
func (l *VirtualMachineLedger) RecordExecutedEvent(vmID, initiator, description string) {
    event := VMEventLog{
        EventID:     GenerateUniqueID(), // Assuming a function to generate a unique ID
        VMID:        vmID,
        EventType:   "Executed",
        Timestamp:   time.Now(),
        Initiator:   initiator,
        Description: description,
        Impact:      "medium", // Default impact level for executed events
    }
    l.VMEventLog = append(l.VMEventLog, event)
    fmt.Printf("Executed Event Recorded: %v\n", event)
}

// GetEventLog retrieves the ledger's event log.
func (l *VirtualMachineLedger) GetEventLog() []VMEventLog {
    return l.VMEventLog
}





// Modify RecordEvent method to log event details as a string
func (l *Ledger) RecordLedgerEvent(eventType, details string) {
	logMessage := fmt.Sprintf("[%s] %s: %s", time.Now().Format(time.RFC3339), eventType, details)
	// Log the event or store it as required; for demonstration, we'll just print
	fmt.Println(logMessage)
}