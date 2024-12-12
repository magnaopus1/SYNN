package automations

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy_network_demo/transactions"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
)

// TransactionBatchingOptimizationAutomation automates the batching of transactions for optimization.
type TransactionBatchingOptimizationAutomation struct {
	ledgerInstance     *ledger.Ledger
	transactionPool    *transactions.TransactionPool
	subBlockBatchSize  int
	mutex              sync.Mutex
	stopChan           chan bool
	maxTransactions    int
	networkLoadMonitor func() float64 // Function to monitor network load and adjust batching size
	subBlockValidator  func(common.SubBlock) error
}

// NewTransactionBatchingOptimizationAutomation initializes a new TransactionBatchingOptimizationAutomation.
func NewTransactionBatchingOptimizationAutomation(ledgerInstance *ledger.Ledger, transactionPool *transactions.TransactionPool, networkLoadMonitor func() float64, subBlockValidator func(common.SubBlock) error) *TransactionBatchingOptimizationAutomation {
	return &TransactionBatchingOptimizationAutomation{
		ledgerInstance:     ledgerInstance,
		transactionPool:    transactionPool,
		subBlockBatchSize:  100, // Default batch size
		maxTransactions:    1000,
		networkLoadMonitor: networkLoadMonitor,
		subBlockValidator:  subBlockValidator,
		stopChan:           make(chan bool),
	}
}

// Start begins the continuous batching and optimization process.
func (t *TransactionBatchingOptimizationAutomation) Start() {
	go t.runBatchingLoop()
	log.Println("Transaction Batching Optimization Automation started.")
}

// Stop stops the continuous batching process.
func (t *TransactionBatchingOptimizationAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Batching Optimization Automation stopped.")
}

// runBatchingLoop continuously batches transactions into sub-blocks for validation.
func (t *TransactionBatchingOptimizationAutomation) runBatchingLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.optimizeBatchAndProcess()
		case <-t.stopChan:
			return
		}
	}
}

// optimizeBatchAndProcess optimizes the batch size based on network load and processes the transactions.
func (t *TransactionBatchingOptimizationAutomation) optimizeBatchAndProcess() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Adjust the batch size based on network load
	networkLoad := t.networkLoadMonitor()
	t.subBlockBatchSize = t.adjustBatchSize(networkLoad)

	// Fetch transactions from the pool
	transactions, err := t.transactionPool.ListTransactions()
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return
	}

	// Group transactions into batches
	for len(transactions) > 0 {
		batch := t.batchTransactions(transactions)
		transactions = transactions[len(batch):]

		// Validate and process the sub-block
		err := t.processSubBlock(batch)
		if err != nil {
			log.Printf("Error processing sub-block: %v", err)
		}
	}
}

// batchTransactions batches transactions based on the current batch size.
func (t *TransactionBatchingOptimizationAutomation) batchTransactions(transactions []common.Transaction) []common.Transaction {
	batchSize := t.subBlockBatchSize
	if len(transactions) < batchSize {
		batchSize = len(transactions)
	}
	return transactions[:batchSize]
}

// adjustBatchSize dynamically adjusts the batch size based on network load.
func (t *TransactionBatchingOptimizationAutomation) adjustBatchSize(networkLoad float64) int {
	if networkLoad < 0.5 {
		// Low load, increase batch size
		return min(t.maxTransactions, t.subBlockBatchSize+100)
	} else if networkLoad > 0.8 {
		// High load, decrease batch size
		return max(10, t.subBlockBatchSize-100)
	}
	// Moderate load, keep the batch size constant
	return t.subBlockBatchSize
}

// processSubBlock validates and processes a batch of transactions into a sub-block.
func (t *TransactionBatchingOptimizationAutomation) processSubBlock(transactions []common.Transaction) error {
	subBlockID := generateSubBlockID()
	subBlock := common.SubBlock{
		ID:           subBlockID,
		Transactions: transactions,
		Validated:    false,
	}

	// Validate the sub-block using Synnergy Consensus
	err := t.subBlockValidator(subBlock)
	if err != nil {
		return fmt.Errorf("sub-block validation failed: %v", err)
	}

	// Encrypt and add the sub-block to the ledger
	err = t.addSubBlockToLedger(subBlock)
	if err != nil {
		return fmt.Errorf("failed to add sub-block to ledger: %v", err)
	}

	log.Printf("Sub-block %s processed with %d transactions", subBlock.ID, len(transactions))
	return nil
}

// addSubBlockToLedger records a validated sub-block in the ledger with encryption.
func (t *TransactionBatchingOptimizationAutomation) addSubBlockToLedger(subBlock common.SubBlock) error {
	// Encrypt the sub-block data
	encryptedData, err := encryption.EncryptData(fmt.Sprintf("SubBlock:%s", subBlock.ID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block data: %v", err)
	}
	// Record the sub-block in the ledger
	err = t.ledgerInstance.RecordSubBlock(subBlock, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record sub-block in ledger: %v", err)
	}
	return nil
}

// generateSubBlockID generates a new sub-block ID (UUID, cryptographically secure hash, or similar).
func generateSubBlockID() string {
	return common.GenerateID()
}

// Utility functions to ensure batch size is within bounds
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
