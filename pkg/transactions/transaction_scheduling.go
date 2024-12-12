package transactions

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// TransactionScheduler manages the scheduling of blockchain transactions.
type TransactionScheduler struct {
	mutex          sync.Mutex                       // For thread safety
	ScheduleList   map[string]*TransactionSchedule  // Scheduled transactions list
	Ledger         *ledger.Ledger                   // Reference to the ledger for transaction logging
	ContractManager *common.SmartContractManager    // Reference to the SmartContractManager
}

// TransactionSchedule defines a scheduled transaction with its conditions.
type TransactionSchedule struct {
    TransactionID  string              // Unique ID of the transaction
    ScheduledTime  time.Time           // Time when the transaction is scheduled to be executed
    BlockHeight    uint64              // Block height at which the transaction will be executed (optional)
    Condition      string              // Condition to trigger the transaction (e.g., contract state change)
    Recurring      bool                // Whether this transaction is recurring
    NextExecution  time.Time           // Next execution time for recurring transactions
    EncryptedTx    string              // Encrypted transaction data
    ValidatorID    string              // Validator responsible for the transaction execution
    Executed       bool                // Whether the transaction has already been executed
}

// Condition defines the structure for the transaction condition
type Condition struct {
	ContractID    string      `json:"contractID"`
	ConditionKey  string      `json:"conditionKey"`
	ExpectedValue interface{} `json:"expectedValue"`
}


// NewTransactionScheduler initializes a new transaction scheduler.
func NewTransactionScheduler(ledgerInstance *ledger.Ledger) *TransactionScheduler {
	return &TransactionScheduler{
		ScheduleList: make(map[string]*TransactionSchedule),
		Ledger:       ledgerInstance,
	}
}

// ScheduleTransaction schedules a transaction based on time, block height, or condition.
func (ts *TransactionScheduler) ScheduleTransaction(txID string, scheduledTime time.Time, blockHeight uint64, condition string, encryptedTx string, validatorID string, recurring bool) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if _, exists := ts.ScheduleList[txID]; exists {
		return errors.New("transaction already scheduled")
	}

	ts.ScheduleList[txID] = &TransactionSchedule{
		TransactionID:  txID,
		ScheduledTime:  scheduledTime,
		BlockHeight:    blockHeight,
		Condition:      condition,
		Recurring:      recurring,
		NextExecution:  scheduledTime,
		EncryptedTx:    encryptedTx,
		ValidatorID:    validatorID,
		Executed:       false,
	}
	return nil
}

// CheckAndExecuteScheduledTransactions checks if any scheduled transactions meet their execution criteria.
func (ts *TransactionScheduler) CheckAndExecuteScheduledTransactions(currentTime time.Time, currentBlockHeight uint64) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	for _, schedule := range ts.ScheduleList {
		if !schedule.Executed && ts.shouldExecute(schedule, currentTime, currentBlockHeight) {
			ts.executeTransaction(schedule)
			if !schedule.Recurring {
				schedule.Executed = true
			} else {
				// Reschedule the next execution for recurring transactions
				schedule.NextExecution = schedule.NextExecution.Add(time.Hour * 24) // Example of recurring daily transactions
			}
		}
	}
}

// shouldExecute checks whether the transaction meets its scheduled condition.
func (ts *TransactionScheduler) shouldExecute(schedule *TransactionSchedule, currentTime time.Time, currentBlockHeight uint64) bool {
	// Check if it's time-based
	if !schedule.ScheduledTime.IsZero() && currentTime.After(schedule.ScheduledTime) {
		return true
	}

	// Check if it's block-height-based
	if schedule.BlockHeight > 0 && currentBlockHeight >= schedule.BlockHeight {
		return true
	}

	// Check if it's condition-based
	if schedule.Condition != "" {
		// Parse the schedule.Condition JSON string into a Condition struct
		var condition Condition
		err := json.Unmarshal([]byte(schedule.Condition), &condition)
		if err != nil {
			fmt.Printf("Error parsing condition for schedule %s: %v\n", schedule.TransactionID, err)
			return false
		}

		// Now that we have the condition data, we can pass it to checkCondition
		conditionMet, err := ts.checkCondition(condition.ContractID, condition.ConditionKey, condition.ExpectedValue)
		if err != nil {
			fmt.Printf("Error checking condition for schedule %s: %v\n", schedule.TransactionID, err)
			return false
		}

		// If the condition is met, return true
		if conditionMet {
			return true
		}
	}

	return false
}


// checkCondition evaluates a condition based on the actual state of the smart contract.
func (ts *TransactionScheduler) checkCondition(contractID, conditionKey string, expectedValue interface{}) (bool, error) {
	// Retrieve the smart contract from the ContractManager
	contract, exists := ts.ContractManager.Contracts[contractID]
	if !exists {
		return false, fmt.Errorf("smart contract with ID %s not found", contractID)
	}

	// Check if the conditionKey exists in the contract state
	stateValue, exists := contract.State[conditionKey]
	if !exists {
		return false, fmt.Errorf("condition %s not found in contract %s", conditionKey, contractID)
	}

	// Evaluate if the state value matches the expected value
	if stateValue == expectedValue {
		return true, nil
	}

	return false, nil
}



// executeTransaction decrypts the transaction and logs its execution to the ledger.
func (ts *TransactionScheduler) executeTransaction(schedule *TransactionSchedule) {
    // Create encryption instance for decrypting transactions
    encryptionInstance := &common.Encryption{}

    // Convert schedule.EncryptedTx from string to []byte
    encryptedTxBytes, err := hex.DecodeString(schedule.EncryptedTx)
    if err != nil {
        fmt.Printf("Error decoding encrypted transaction %s: %v\n", schedule.TransactionID, err)
        return
    }

    // Decrypt the transaction
    decryptedTx, err := encryptionInstance.DecryptData(encryptedTxBytes, common.EncryptionKey)
    if err != nil {
        fmt.Printf("Error decrypting transaction %s: %v\n", schedule.TransactionID, err)
        return
    }

    // Convert decryptedTx (which is []byte) to string
    decryptedTxStr := string(decryptedTx)

    // Assuming RecordTransactionExecution only needs the TransactionID
    err = ts.Ledger.RecordTransactionExecution(schedule.TransactionID)
    if err != nil {
        fmt.Printf("Error logging transaction execution for %s: %v\n", schedule.TransactionID, err)
        return
    }

    // If you need to log the decrypted transaction separately, you can do it here
    fmt.Printf("Decrypted transaction for %s: %s\n", schedule.TransactionID, decryptedTxStr)

    fmt.Printf("Transaction %s executed successfully at %v.\n", schedule.TransactionID, time.Now())
}


// CancelScheduledTransaction cancels a previously scheduled transaction.
func (ts *TransactionScheduler) CancelScheduledTransaction(txID string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	if _, exists := ts.ScheduleList[txID]; !exists {
		return errors.New("transaction not found")
	}

	delete(ts.ScheduleList, txID)
	fmt.Printf("Transaction %s has been successfully cancelled.\n", txID)
	return nil
}

// GetScheduledTransactions returns the list of scheduled transactions.
func (ts *TransactionScheduler) GetScheduledTransactions() map[string]*TransactionSchedule {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	return ts.ScheduleList
}

// StartScheduler is an example of how the scheduler could run periodically.
func (ts *TransactionScheduler) StartScheduler(interval time.Duration, currentBlockHeight *uint64) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ts.CheckAndExecuteScheduledTransactions(time.Now(), *currentBlockHeight)
		}
	}
}
