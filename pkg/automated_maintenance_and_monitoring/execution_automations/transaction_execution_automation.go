package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/transactions"
)

const (
	TransactionCheckInterval   = 30 * time.Second // Interval for transaction monitoring
	TransactionLedgerEntryType = "Transaction"
)

// TransactionExecutionAutomation handles automatic validation, execution, and logging of transactions
type TransactionExecutionAutomation struct {
	consensusEngine *consensus.SynnergyConsensus // Synnergy Consensus engine
	ledgerInstance  *ledger.Ledger               // Ledger for transaction logging
	transactionPool *transactions.TransactionPool // Transaction pool to validate and execute
	executionMutex  *sync.RWMutex                 // Mutex for thread-safe execution
}

// NewTransactionExecutionAutomation initializes the transaction automation with consensus and ledger integration
func NewTransactionExecutionAutomation(consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionPool *transactions.TransactionPool, executionMutex *sync.RWMutex) *TransactionExecutionAutomation {
	return &TransactionExecutionAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		transactionPool: transactionPool,
		executionMutex:  executionMutex,
	}
}

// StartTransactionMonitor starts continuous monitoring and execution of transactions
func (automation *TransactionExecutionAutomation) StartTransactionMonitor() {
	ticker := time.NewTicker(TransactionCheckInterval)

	go func() {
		for range ticker.C {
			automation.processPendingTransactions()
		}
	}()
}

// processPendingTransactions processes the pending transactions in the pool by validating and executing them
func (automation *TransactionExecutionAutomation) processPendingTransactions() {
	automation.executionMutex.Lock()
	defer automation.executionMutex.Unlock()

	// Get all pending transactions from the pool
	pendingTransactions := automation.transactionPool.GetPendingTransactions()

	for _, tx := range pendingTransactions {
		automation.validateAndExecuteTransaction(tx)
	}
}

// validateAndExecuteTransaction validates the transaction with consensus and executes it if valid
func (automation *TransactionExecutionAutomation) validateAndExecuteTransaction(tx *transactions.Transaction) {
	// Validate the transaction using Synnergy Consensus
	valid, err := automation.consensusEngine.ValidateTransaction(tx)
	if err != nil {
		fmt.Printf("Failed to validate transaction %s: %v\n", tx.ID, err)
		automation.logTransactionFailure(tx, "Consensus validation failed")
		return
	}

	if !valid {
		fmt.Printf("Transaction %s failed consensus validation.\n", tx.ID)
		automation.logTransactionFailure(tx, "Consensus validation failed")
		return
	}

	// Execute the transaction
	err = automation.transactionPool.ExecuteTransaction(tx)
	if err != nil {
		fmt.Printf("Failed to execute transaction %s: %v\n", tx.ID, err)
		automation.logTransactionFailure(tx, "Transaction execution failed")
		return
	}

	// Log the successful transaction execution in the ledger
	automation.logTransactionSuccess(tx)
}

// logTransactionSuccess logs the successful execution of a transaction into the ledger
func (automation *TransactionExecutionAutomation) logTransactionSuccess(tx *transactions.Transaction) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("transaction-success-%s-%d", tx.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      TransactionLedgerEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("Transaction %s successfully executed.", tx.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log successful transaction %s in the ledger: %v\n", tx.ID, err)
	} else {
		fmt.Println("Transaction successfully logged in the ledger.")
	}
}

// logTransactionFailure logs a failure event when a transaction fails validation or execution
func (automation *TransactionExecutionAutomation) logTransactionFailure(tx *transactions.Transaction, reason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("transaction-failure-%s-%d", tx.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      TransactionLedgerEntryType,
		Status:    "Failure",
		Details:   fmt.Sprintf("Transaction %s failed: %s", tx.ID, reason),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log transaction failure %s in the ledger: %v\n", tx.ID, err)
	} else {
		fmt.Println("Transaction failure successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TransactionExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualTransactionExecution allows administrators to manually trigger the execution of a specific transaction
func (automation *TransactionExecutionAutomation) TriggerManualTransactionExecution(transactionID string) {
	fmt.Printf("Manually triggering execution of transaction %s...\n", transactionID)

	tx := automation.transactionPool.GetTransactionByID(transactionID)
	if tx != nil {
		automation.validateAndExecuteTransaction(tx)
	} else {
		fmt.Printf("Transaction %s not found in the pool.\n", transactionID)
	}
}
