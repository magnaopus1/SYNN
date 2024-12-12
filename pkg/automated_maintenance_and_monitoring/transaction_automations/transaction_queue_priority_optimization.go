package automations

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/transactions"
)

// TransactionQueuePriorityOptimization automates the optimization of transaction prioritization.
type TransactionQueuePriorityOptimization struct {
	ledgerInstance    *ledger.Ledger
	transactionPool   *transactions.TransactionPool
	mutex             sync.Mutex
	stopChan          chan bool
	networkLoadMonitor func() float64 // To monitor network load for dynamic adjustments
}

// NewTransactionQueuePriorityOptimization initializes the automation.
func NewTransactionQueuePriorityOptimization(ledgerInstance *ledger.Ledger, transactionPool *transactions.TransactionPool, networkLoadMonitor func() float64) *TransactionQueuePriorityOptimization {
	return &TransactionQueuePriorityOptimization{
		ledgerInstance:    ledgerInstance,
		transactionPool:   transactionPool,
		stopChan:          make(chan bool),
		networkLoadMonitor: networkLoadMonitor,
	}
}

// Start begins the continuous transaction priority optimization process.
func (t *TransactionQueuePriorityOptimization) Start() {
	go t.runPriorityOptimizationLoop()
	log.Println("Transaction Queue Priority Optimization Automation started.")
}

// Stop stops the transaction priority optimization process.
func (t *TransactionQueuePriorityOptimization) Stop() {
	t.stopChan <- true
	log.Println("Transaction Queue Priority Optimization Automation stopped.")
}

// runPriorityOptimizationLoop continuously reorders the transaction queue every 50 milliseconds based on priority.
func (t *TransactionQueuePriorityOptimization) runPriorityOptimizationLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.optimizeTransactionQueue()
		case <-t.stopChan:
			return
		}
	}
}

// optimizeTransactionQueue reorders the transaction queue based on priority (fee, network load, transaction type).
func (t *TransactionQueuePriorityOptimization) optimizeTransactionQueue() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Get all transactions from the pool
	transactions, err := t.transactionPool.ListTransactions()
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return
	}

	// Check if there are any transactions to reorder
	if len(transactions) == 0 {
		log.Println("Transaction pool is empty, skipping optimization.")
		return
	}

	// Prioritize transactions by fee, network load, and type
	t.reorderTransactionsByPriority(transactions)

	log.Printf("Transaction pool reordered based on priority. Current pool size: %d", len(transactions))
}

// reorderTransactionsByPriority reorders transactions by priority based on fee, network load, and transaction type.
func (t *TransactionQueuePriorityOptimization) reorderTransactionsByPriority(transactions []common.Transaction) {
	networkLoad := t.networkLoadMonitor()

	// Sort transactions in-place based on fee, transaction type, and network load.
	sort.SliceStable(transactions, func(i, j int) bool {
		return t.calculatePriority(transactions[i], networkLoad) > t.calculatePriority(transactions[j], networkLoad)
	})

	// Update the transaction pool with the reordered list
	err := t.transactionPool.UpdateTransactionPool(transactions)
	if err != nil {
		log.Printf("Failed to update transaction pool with reordered transactions: %v", err)
	}
}

// calculatePriority calculates the priority of a transaction based on its fee, type, and network load.
func (t *TransactionQueuePriorityOptimization) calculatePriority(tx common.Transaction, networkLoad float64) float64 {
	basePriority := tx.Fee

	// Add weight based on transaction type (e.g., critical transactions get higher priority)
	switch tx.Type {
	case common.TransactionTypeCritical:
		basePriority += 10.0
	case common.TransactionTypeStandard:
		basePriority += 1.0
	}

	// Adjust priority based on network load
	if networkLoad > 0.8 {
		// During high load, give preference to higher-fee transactions
		basePriority *= 1.5
	}

	return basePriority
}

// Trigger: Add transactions based on priority directly into the validation process.
func (t *TransactionQueuePriorityOptimization) triggerValidationForHighPriorityTransactions() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch transactions ordered by priority
	transactions, err := t.transactionPool.ListTransactions()
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return
	}

	// Only trigger validation for the highest priority transactions (e.g., top 10)
	if len(transactions) > 10 {
		transactions = transactions[:10]
	}

	// Directly send high-priority transactions for validation
	for _, tx := range transactions {
		err := t.validateTransaction(tx)
		if err != nil {
			log.Printf("Failed to validate high-priority transaction %s: %v", tx.ID, err)
		}
	}
}

// validateTransaction sends a high-priority transaction for immediate validation.
func (t *TransactionQueuePriorityOptimization) validateTransaction(transaction common.Transaction) error {
	// Simulate validation process (real implementation would involve actual validation steps)
	encryptedData, err := encryption.EncryptData(fmt.Sprintf("Transaction:%s", transaction.ID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Record the transaction validation in the ledger
	err = t.ledgerInstance.RecordTransactionValidation(transaction, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record transaction validation: %v", err)
	}

	log.Printf("High-priority transaction %s validated.", transaction.ID)
	return nil
}
