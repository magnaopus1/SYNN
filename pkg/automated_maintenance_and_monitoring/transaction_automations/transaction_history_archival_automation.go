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

// TransactionHistoryArchivalAutomation automates the archival of transaction history from finalized blocks.
type TransactionHistoryArchivalAutomation struct {
	ledgerInstance *ledger.Ledger
	archiveManager *transactions.ArchiveManager
	mutex          sync.Mutex
	stopChan       chan bool
}

// NewTransactionHistoryArchivalAutomation initializes a new TransactionHistoryArchivalAutomation.
func NewTransactionHistoryArchivalAutomation(ledgerInstance *ledger.Ledger, archiveManager *transactions.ArchiveManager) *TransactionHistoryArchivalAutomation {
	return &TransactionHistoryArchivalAutomation{
		ledgerInstance: ledgerInstance,
		archiveManager: archiveManager,
		stopChan:       make(chan bool),
	}
}

// Start begins the process of periodically archiving transaction history.
func (t *TransactionHistoryArchivalAutomation) Start() {
	go t.runArchivalLoop()
	log.Println("Transaction History Archival Automation started.")
}

// Stop halts the archival automation.
func (t *TransactionHistoryArchivalAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction History Archival Automation stopped.")
}

// runArchivalLoop continuously checks for finalized blocks and archives their transaction history every hour.
func (t *TransactionHistoryArchivalAutomation) runArchivalLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.archiveFinalizedBlockHistory()
		case <-t.stopChan:
			return
		}
	}
}

// archiveFinalizedBlockHistory archives the transaction history for finalized blocks.
func (t *TransactionHistoryArchivalAutomation) archiveFinalizedBlockHistory() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Retrieve finalized blocks from the ledger that are not archived yet.
	finalizedBlocks, err := t.ledgerInstance.GetFinalizedBlocks()
	if err != nil {
		log.Printf("Failed to retrieve finalized blocks for archival: %v", err)
		return
	}

	for _, block := range finalizedBlocks {
		err := t.archiveBlockHistory(block)
		if err != nil {
			log.Printf("Failed to archive transaction history for block %s: %v", block.ID, err)
		}
	}
}

// archiveBlockHistory archives the transaction history of a specific finalized block.
func (t *TransactionHistoryArchivalAutomation) archiveBlockHistory(block common.Block) error {
	// Retrieve the transaction history for the block.
	transactions, err := t.ledgerInstance.GetBlockTransactions(block.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transactions for block %s: %v", block.ID, err)
	}

	// Encrypt the transaction data before archiving.
	encryptedHistory, err := encryption.EncryptData(fmt.Sprintf("Block:%s Transactions:%v", block.ID, transactions), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction history for block %s: %v", block.ID, err)
	}

	// Archive the transaction history.
	err = t.archiveManager.ArchiveTransactionHistory(block.ID, encryptedHistory)
	if err != nil {
		return fmt.Errorf("failed to archive transaction history for block %s: %v", block.ID, err)
	}

	// Mark the block as archived in the ledger.
	err = t.ledgerInstance.MarkBlockAsArchived(block.ID)
	if err != nil {
		return fmt.Errorf("failed to mark block %s as archived: %v", block.ID, err)
	}

	// Log the archival action.
	log.Printf("Transaction history for block %s archived successfully.", block.ID)

	return nil
}

// RetrieveArchivedBlockHistory retrieves the archived transaction history for a specific block.
func (t *TransactionHistoryArchivalAutomation) RetrieveArchivedBlockHistory(blockID string) (string, error) {
	// Retrieve archived transaction history from the archive manager.
	archivedHistory, err := t.archiveManager.RetrieveArchivedHistory(blockID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve archived history for block %s: %v", blockID, err)
	}

	// Decrypt the archived history for viewing.
	decryptedHistory, err := encryption.DecryptData(archivedHistory, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt archived history for block %s: %v", blockID, err)
	}

	log.Printf("Archived transaction history for block %s retrieved successfully.", blockID)
	return decryptedHistory, nil
}

// ManualArchive allows for manual archival of a specific block's transaction history.
func (t *TransactionHistoryArchivalAutomation) ManualArchive(blockID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Retrieve the finalized block from the ledger.
	block, err := t.ledgerInstance.GetBlockByID(blockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve block %s for manual archival: %v", blockID, err)
	}

	// Archive the block's transaction history.
	err = t.archiveBlockHistory(block)
	if err != nil {
		return fmt.Errorf("failed to manually archive transaction history for block %s: %v", blockID, err)
	}

	log.Printf("Manually archived transaction history for block %s.", blockID)
	return nil
}
