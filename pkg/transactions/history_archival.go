package transactions

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// ArchiveManager handles archiving and managing old blockchain data.
type ArchiveManager struct {
	ledgerInstance *ledger.Ledger // Reference to the ledger instance
	archivePath    string         // Path to store archived data
	mutex          sync.Mutex     // Mutex for thread-safe operations
	memoryStorage  map[string][]byte // Store encrypted archives in memory
}

// NewArchiveManager initializes a new ArchiveManager with the provided ledger instance and archive path.
func NewArchiveManager(ledgerInstance *ledger.Ledger, archivePath string) *ArchiveManager {
	return &ArchiveManager{
		ledgerInstance: ledgerInstance,
		archivePath:    archivePath,
		memoryStorage:  make(map[string][]byte),
	}
}

// Convert ledger transactions to common transactions.
func convertToCommonTransactions(ledgerTxns []*ledger.Transaction) []common.Transaction {
    commonTxns := make([]common.Transaction, len(ledgerTxns))
    for i, txn := range ledgerTxns {
        // Convert each ledger.Transaction to a common.Transaction by mapping fields.
        commonTxns[i] = common.Transaction{
            TransactionID:     txn.TransactionID,
            FromAddress:       txn.FromAddress,
            ToAddress:         txn.ToAddress,
            Amount:            txn.Amount,
            Fee:               txn.Fee,
            TokenStandard:     txn.TokenStandard,
            TokenID:           txn.TokenID,
            Timestamp:         txn.Timestamp,
            SubBlockID:        txn.SubBlockID,
            BlockID:           txn.BlockID,
            ValidatorID:       txn.ValidatorID,
            Signature:         txn.Signature,
            Status:            txn.Status,
            EncryptedData:     txn.EncryptedData,
            DecryptedData:     txn.DecryptedData,
            ExecutionResult:   txn.ExecutionResult,
            FrozenAmount:      txn.FrozenAmount,
            RefundAmount:      txn.RefundAmount,
            ReversalRequested: txn.ReversalRequested,
        }
    }
    return commonTxns
}


// ArchiveTransactionHistory archives all transaction history for the current block.
func (am *ArchiveManager) ArchiveTransactionHistory(blockID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Retrieve all transactions for the block from the ledger.
	transactions, err := am.ledgerInstance.GetTransactionsByBlockID(blockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transactions for block %s: %v", blockID, err)
	}

	// Convert ledger transactions to common transactions.
	commonTransactions := convertToCommonTransactions(transactions)

	// Generate archive data.
	archiveData := fmt.Sprintf("Block: %s\nTimestamp: %s\n\nTransactions:\n%s", blockID, time.Now().String(), formatTransactions(commonTransactions))

	// Create an encryption instance
	encryptionInstance := &common.Encryption{}

	// Encrypt the archive data for security.
	encryptedArchiveData, err := encryptionInstance.EncryptData("AES", []byte(archiveData), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt archive data for block %s: %v", blockID, err)
	}

	// Store the encrypted archive data in memory
	am.memoryStorage[blockID] = encryptedArchiveData

	// Generate archive file path.
	archiveFilePath := filepath.Join(am.archivePath, fmt.Sprintf("block_%s_archive.dat", blockID))

	// Write the encrypted archive data to a file.
	err = os.WriteFile(archiveFilePath, encryptedArchiveData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write archive file for block %s: %v", blockID, err)
	}

	fmt.Printf("Transaction history for block %s archived successfully at %s\n", blockID, archiveFilePath)
	return nil
}



// RetrieveArchivedHistory retrieves the archived history for a given block, decrypts, and returns the transaction data.
func (am *ArchiveManager) RetrieveArchivedHistory(blockID string) (string, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Generate archive file path.
	archiveFilePath := filepath.Join(am.archivePath, fmt.Sprintf("block_%s_archive.dat", blockID))

	// Check if the archive file exists.
	if _, err := os.Stat(archiveFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("no archive found for block %s", blockID)
	}

	// Read the encrypted archive data from the file.
	encryptedArchiveData, err := os.ReadFile(archiveFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read archive file for block %s: %v", blockID, err)
	}

	// Create an encryption instance for decrypting data
	encryptionInstance := &common.Encryption{}

	// Decrypt the archive data using the encryption instance
	decryptedArchiveData, err := encryptionInstance.DecryptData(encryptedArchiveData, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt archive data for block %s: %v", blockID, err)
	}

	fmt.Printf("Transaction history for block %s retrieved successfully.\n", blockID)
	return string(decryptedArchiveData), nil
}

// DeleteArchivedHistory permanently deletes the archived history for a given block.
func (am *ArchiveManager) DeleteArchivedHistory(blockID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Generate archive file path.
	archiveFilePath := filepath.Join(am.archivePath, fmt.Sprintf("block_%s_archive.dat", blockID))

	// Check if the archive file exists.
	if _, err := os.Stat(archiveFilePath); os.IsNotExist(err) {
		return fmt.Errorf("no archive found for block %s", blockID)
	}

	// Delete the archive file.
	err := os.Remove(archiveFilePath)
	if err != nil {
		return fmt.Errorf("failed to delete archive file for block %s: %v", blockID, err)
	}

	fmt.Printf("Archived transaction history for block %s deleted successfully.\n", blockID)
	return nil
}

// formatTransactions formats the list of transactions into a human-readable string.
func formatTransactions(transactions []common.Transaction) string {
	var formatted string
	for _, tx := range transactions {
		formatted += fmt.Sprintf("Transaction ID: %s, Amount: %f, Timestamp: %s\n", tx.TransactionID, tx.Amount, tx.Timestamp.String())
	}
	return formatted
}

