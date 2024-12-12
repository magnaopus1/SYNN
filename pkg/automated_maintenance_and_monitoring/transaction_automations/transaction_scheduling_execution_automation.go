package automations

import (
	"log"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/transactions"
)

// TransactionSchedulingExecutionAutomation automates the execution of scheduled transactions.
type TransactionSchedulingExecutionAutomation struct {
	ledgerInstance   *ledger.Ledger
	scheduler        *transactions.TransactionScheduler
	mutex            sync.Mutex
	stopChan         chan bool
	currentBlockHeight *uint64
}

// NewTransactionSchedulingExecutionAutomation initializes the automation.
func NewTransactionSchedulingExecutionAutomation(ledgerInstance *ledger.Ledger, scheduler *transactions.TransactionScheduler, blockHeight *uint64) *TransactionSchedulingExecutionAutomation {
	return &TransactionSchedulingExecutionAutomation{
		ledgerInstance:    ledgerInstance,
		scheduler:         scheduler,
		stopChan:          make(chan bool),
		currentBlockHeight: blockHeight,
	}
}

// Start begins the continuous loop for executing scheduled transactions every 200 milliseconds.
func (t *TransactionSchedulingExecutionAutomation) Start() {
	go t.runSchedulingExecutionLoop()
	log.Println("Transaction Scheduling Execution Automation started.")
}

// Stop halts the transaction scheduling execution process.
func (t *TransactionSchedulingExecutionAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Scheduling Execution Automation stopped.")
}

// runSchedulingExecutionLoop checks every 200 milliseconds for scheduled transactions to execute.
func (t *TransactionSchedulingExecutionAutomation) runSchedulingExecutionLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.executeScheduledTransactions()
		case <-t.stopChan:
			return
		}
	}
}

// executeScheduledTransactions looks for scheduled transactions that are ready to be executed and executes them.
func (t *TransactionSchedulingExecutionAutomation) executeScheduledTransactions() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch all scheduled transactions
	scheduledTxs := t.scheduler.GetScheduledTransactions()

	// Execute transactions that are due based on time, block height, or condition
	for _, scheduledTx := range scheduledTxs {
		if !scheduledTx.Executed && t.shouldExecuteTransaction(scheduledTx) {
			t.executeTransaction(scheduledTx)
			if !scheduledTx.Recurring {
				scheduledTx.Executed = true
			} else {
				// Update next execution for recurring transactions (example: daily execution)
				scheduledTx.NextExecution = scheduledTx.NextExecution.Add(24 * time.Hour)
			}
		}
	}
}

// shouldExecuteTransaction checks whether a scheduled transaction is due for execution.
func (t *TransactionSchedulingExecutionAutomation) shouldExecuteTransaction(schedule *transactions.TransactionSchedule) bool {
	currentTime := time.Now()

	// Check if the transaction's scheduled time has passed
	if !schedule.ScheduledTime.IsZero() && currentTime.After(schedule.ScheduledTime) {
		return true
	}

	// Check if the block height condition is met
	if schedule.BlockHeight > 0 && *t.currentBlockHeight >= schedule.BlockHeight {
		return true
	}

	// Check if a custom condition is satisfied (via smart contracts or other mechanisms)
	if schedule.Condition != "" && t.scheduler.CheckCondition(schedule.Condition) {
		return true
	}

	return false
}

// executeTransaction decrypts and processes the scheduled transaction.
func (t *TransactionSchedulingExecutionAutomation) executeTransaction(schedule *transactions.TransactionSchedule) {
	// Decrypt the transaction data
	decryptedTx, err := encryption.DecryptData(schedule.EncryptedTx, common.EncryptionKey)
	if err != nil {
		log.Printf("Failed to decrypt transaction %s: %v\n", schedule.TransactionID, err)
		return
	}

	// Log and execute the transaction on the ledger
	err = t.ledgerInstance.RecordTransactionExecution(schedule.TransactionID, decryptedTx)
	if err != nil {
		log.Printf("Failed to log execution of transaction %s: %v\n", schedule.TransactionID, err)
		return
	}

	log.Printf("Scheduled transaction %s executed successfully.\n", schedule.TransactionID)
}

// Trigger: Executes the transaction when the scheduled time or condition is met.
func (t *TransactionSchedulingExecutionAutomation) ExecuteTransactionOnCondition(scheduledTx *transactions.TransactionSchedule) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.shouldExecuteTransaction(scheduledTx) {
		t.executeTransaction(scheduledTx)
		if !scheduledTx.Recurring {
			scheduledTx.Executed = true
		}
	}
	return nil
}
