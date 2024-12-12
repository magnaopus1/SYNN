package automations

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/transactions"
)

// TransactionReceiptManagementAutomation automates the creation and verification of transaction receipts.
type TransactionReceiptManagementAutomation struct {
	ledgerInstance     *ledger.Ledger
	receiptManager     *transactions.ReceiptManager
	subBlockValidator  *transactions.SubBlockValidator
	mutex              sync.Mutex
	stopChan           chan bool
}

// NewTransactionReceiptManagementAutomation initializes the automation.
func NewTransactionReceiptManagementAutomation(ledgerInstance *ledger.Ledger, receiptManager *transactions.ReceiptManager, subBlockValidator *transactions.SubBlockValidator) *TransactionReceiptManagementAutomation {
	return &TransactionReceiptManagementAutomation{
		ledgerInstance:     ledgerInstance,
		receiptManager:     receiptManager,
		subBlockValidator:  subBlockValidator,
		stopChan:           make(chan bool),
	}
}

// Start begins the continuous transaction receipt creation and verification process.
func (t *TransactionReceiptManagementAutomation) Start() {
	go t.runReceiptManagementLoop()
	log.Println("Transaction Receipt Management Automation started.")
}

// Stop stops the transaction receipt management process.
func (t *TransactionReceiptManagementAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Receipt Management Automation stopped.")
}

// runReceiptManagementLoop continuously monitors confirmed transactions and generates receipts every 100 milliseconds.
func (t *TransactionReceiptManagementAutomation) runReceiptManagementLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.generateReceiptsForConfirmedTransactions()
		case <-t.stopChan:
			return
		}
	}
}

// generateReceiptsForConfirmedTransactions creates receipts for all confirmed transactions within a sub-block.
func (t *TransactionReceiptManagementAutomation) generateReceiptsForConfirmedTransactions() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch validated sub-blocks from the ledger
	subBlocks, err := t.ledgerInstance.GetValidatedSubBlocks()
	if err != nil {
		log.Printf("Failed to retrieve validated sub-blocks: %v", err)
		return
	}

	// Loop through each validated sub-block and create receipts for confirmed transactions
	for _, subBlock := range subBlocks {
		for _, transaction := range subBlock.Transactions {
			if !transaction.ReceiptGenerated {
				err := t.createTransactionReceipt(transaction, subBlock)
				if err != nil {
					log.Printf("Failed to generate receipt for transaction %s: %v", transaction.ID, err)
				} else {
					log.Printf("Receipt generated for transaction %s", transaction.ID)
				}
			}
		}
	}
}

// createTransactionReceipt generates a receipt for a confirmed transaction.
func (t *TransactionReceiptManagementAutomation) createTransactionReceipt(transaction common.Transaction, subBlock common.SubBlock) error {
	// Encrypt transaction data for security
	encryptedTxData, err := encryption.EncryptData(fmt.Sprintf("TransactionID:%s, Amount:%f", transaction.ID, transaction.Amount), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	// Create a receipt and record it in the ledger
	receipt, err := t.receiptManager.CreateReceipt(transaction, subBlock.ID, subBlock.BlockID, encryptedTxData)
	if err != nil {
		return fmt.Errorf("failed to create receipt: %v", err)
	}

	// Update the transaction as having generated a receipt
	err = t.ledgerInstance.MarkTransactionAsReceiptGenerated(transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to update transaction receipt status in ledger: %v", err)
	}

	log.Printf("Transaction receipt created: %v", receipt)
	return nil
}

// Trigger: Verifies receipts automatically if requested by the user or other triggers.
func (t *TransactionReceiptManagementAutomation) VerifyReceiptOnRequest(transactionID string, providedReceipt common.TransactionReceipt) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch the transaction from the ledger
	transaction, err := t.ledgerInstance.GetTransactionByID(transactionID)
	if err != nil {
		return fmt.Errorf("transaction not found: %v", err)
	}

	// Verify the receipt against the transaction
	valid, err := t.receiptManager.VerifyReceipt(&providedReceipt, transaction)
	if err != nil {
		return fmt.Errorf("failed to verify receipt: %v", err)
	}

	if !valid {
		return fmt.Errorf("receipt verification failed for transaction %s", transaction.ID)
	}

	log.Printf("Receipt verified successfully for transaction %s", transaction.ID)
	return nil
}

// generateReceiptOnBlockFinalization generates receipts when a block is finalized, providing a final confirmation.
func (t *TransactionReceiptManagementAutomation) generateReceiptOnBlockFinalization(block common.Block) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, subBlock := range block.SubBlocks {
		for _, transaction := range subBlock.Transactions {
			if !transaction.ReceiptGenerated {
				err := t.createTransactionReceipt(transaction, subBlock)
				if err != nil {
					log.Printf("Failed to generate receipt for transaction %s in block %s: %v", transaction.ID, block.ID, err)
				} else {
					log.Printf("Receipt generated for transaction %s in block %s", transaction.ID, block.ID)
				}
			}
		}
	}
}
