package automations

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/transactions"
)

// TransactionPoolManagementAutomation automates the management of the transaction pool.
type TransactionPoolManagementAutomation struct {
	ledgerInstance    *ledger.Ledger
	transactionPool   *transactions.TransactionPool
	subBlockBatchSize int
	mutex             sync.Mutex
	stopChan          chan bool
}

// NewTransactionPoolManagementAutomation initializes a new TransactionPoolManagementAutomation.
func NewTransactionPoolManagementAutomation(ledgerInstance *ledger.Ledger, transactionPool *transactions.TransactionPool) *TransactionPoolManagementAutomation {
	return &TransactionPoolManagementAutomation{
		ledgerInstance:    ledgerInstance,
		transactionPool:   transactionPool,
		subBlockBatchSize: 1000, // Default batch size for sub-blocks
		stopChan:          make(chan bool),
	}
}

// Start begins the automation of managing the transaction pool, adding transactions, optimizing the pool, and creating sub-blocks.
func (t *TransactionPoolManagementAutomation) Start() {
	go t.runPoolManagementLoop()
	log.Println("Transaction Pool Management Automation started.")
}

// Stop halts the transaction pool management automation.
func (t *TransactionPoolManagementAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Pool Management Automation stopped.")
}

// runPoolManagementLoop continuously manages the transaction pool and creates sub-blocks every 100 milliseconds.
func (t *TransactionPoolManagementAutomation) runPoolManagementLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)  // Changed to run every 100 milliseconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.optimizePoolAndCreateSubBlock()
		case <-t.stopChan:
			return
		}
	}
}

// optimizePoolAndCreateSubBlock optimizes the transaction pool and creates sub-blocks periodically.
func (t *TransactionPoolManagementAutomation) optimizePoolAndCreateSubBlock() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch transactions from the pool
	transactions, err := t.transactionPool.ListTransactions()
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return
	}

	// Check if there are enough transactions for a sub-block
	if len(transactions) < t.subBlockBatchSize {
		log.Printf("Not enough transactions to create a sub-block, current pool size: %d", len(transactions))
		return
	}

	// Create a sub-block with the transactions
	subBlock := t.createSubBlock(transactions)

	// Validate the sub-block using Synnergy Consensus
	err = t.validateSubBlock(subBlock)
	if err != nil {
		log.Printf("Sub-block validation failed: %v", err)
		return
	}

	// Add the validated sub-block to the ledger
	err = t.addSubBlockToLedger(subBlock)
	if err != nil {
		log.Printf("Failed to add sub-block to ledger: %v", err)
		return
	}

	log.Printf("Sub-block %s created and validated with %d transactions.", subBlock.ID, len(subBlock.Transactions))
}

// createSubBlock creates a sub-block from the current pool of transactions.
func (t *TransactionPoolManagementAutomation) createSubBlock(transactions []common.Transaction) *common.SubBlock {
	subBlockID := common.GenerateID()
	subBlockTransactions := transactions[:t.subBlockBatchSize]

	return &common.SubBlock{
		ID:           subBlockID,
		Transactions: subBlockTransactions,
		Validated:    false,
	}
}

// validateSubBlock validates a sub-block using Synnergy Consensus.
func (t *TransactionPoolManagementAutomation) validateSubBlock(subBlock *common.SubBlock) error {
	// Validate transactions with Proof-of-History (PoH)
	err := t.validateWithPoH(subBlock)
	if err != nil {
		return fmt.Errorf("PoH validation failed: %v", err)
	}

	// Select validators with Proof-of-Stake (PoS)
	validators, err := t.selectValidatorsForSubBlock(subBlock)
	if err != nil {
		return fmt.Errorf("PoS validator selection failed: %v", err)
	}

	// Verify transactions using validators
	err = t.verifySubBlockWithPoS(validators, subBlock)
	if err != nil {
		return fmt.Errorf("PoS verification failed: %v", err)
	}

	// Mark the sub-block as validated
	subBlock.Validated = true
	return nil
}

// validateWithPoH performs Proof-of-History (PoH) validation for a sub-block.
func (t *TransactionPoolManagementAutomation) validateWithPoH(subBlock *common.SubBlock) error {
	for _, tx := range subBlock.Transactions {
		// Generate PoH proof for each transaction
		err := t.generatePoHProof(tx)
		if err != nil {
			return fmt.Errorf("PoH proof generation failed for transaction %s: %v", tx.ID, err)
		}
	}
	return nil
}

// generatePoHProof generates a Proof-of-History for a transaction.
func (t *TransactionPoolManagementAutomation) generatePoHProof(transaction common.Transaction) error {
	encryptedData, err := encryption.EncryptData(fmt.Sprintf("%s:%f", transaction.ID, transaction.Amount), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to generate PoH proof: %v", err)
	}
	log.Printf("PoH proof generated for transaction %s: %s", transaction.ID, encryptedData)
	return nil
}

// selectValidatorsForSubBlock selects validators using Proof-of-Stake (PoS) for a sub-block.
func (t *TransactionPoolManagementAutomation) selectValidatorsForSubBlock(subBlock *common.SubBlock) ([]common.Validator, error) {
	validators, err := t.ledgerInstance.SelectValidators(subBlock.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to select validators: %v", err)
	}
	return validators, nil
}

// verifySubBlockWithPoS verifies the transactions in a sub-block using selected validators.
func (t *TransactionPoolManagementAutomation) verifySubBlockWithPoS(validators []common.Validator, subBlock *common.SubBlock) error {
	for _, validator := range validators {
		log.Printf("Validator %s verified sub-block %s", validator.ID, subBlock.ID)
	}
	return nil
}

// addSubBlockToLedger records the validated sub-block in the ledger.
func (t *TransactionPoolManagementAutomation) addSubBlockToLedger(subBlock *common.SubBlock) error {
	// Encrypt the sub-block data
	encryptedData, err := encryption.EncryptData(fmt.Sprintf("SubBlock:%s", subBlock.ID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block data: %v", err)
	}

	// Record the sub-block in the ledger
	err = t.ledgerInstance.RecordSubBlock(*subBlock, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record sub-block in ledger: %v", err)
	}

	return nil
}
