package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
	"synnergy_network_demo/scalability/compression"
)

const (
	TransactionCompressionInterval    = 1 * time.Minute  // Time interval to check for transaction compression
	TransactionCompressionBatchSize   = 100              // Number of transactions to process in each batch
	TransactionCompressionLogKey      = "transactionCompressionKey" // Encryption key for log entries
	TransactionCompressionThreshold   = 5000             // Compression threshold for transactions in bytes
)

// TransactionCompressionAutomation handles compressing transaction data to optimize storage and speed.
type TransactionCompressionAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging compression events
	consensusSystem *consensus.SynnergyConsensus // Reference to the consensus system for transaction handling
	compressionTool *compression.CompressionTool // Reference to the compression tool
	syncMutex       *sync.RWMutex                // Mutex to handle concurrent compression
}

// NewTransactionCompressionAutomation creates a new instance of the automation process.
func NewTransactionCompressionAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, syncMutex *sync.RWMutex) *TransactionCompressionAutomation {
	return &TransactionCompressionAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		compressionTool: compression.NewCompressionTool(), // Use raw compression, without any external package
		syncMutex:       syncMutex,
	}
}

// StartTransactionCompressionAutomation starts the continuous loop for transaction compression.
func (automation *TransactionCompressionAutomation) StartTransactionCompressionAutomation() {
	ticker := time.NewTicker(TransactionCompressionInterval)
	go func() {
		for range ticker.C {
			automation.compressTransactions()
		}
	}()
}

// compressTransactions fetches uncompressed transactions and compresses them in batches.
func (automation *TransactionCompressionAutomation) compressTransactions() {
	automation.syncMutex.Lock()
	defer automation.syncMutex.Unlock()

	// Fetch uncompressed transactions from consensus
	transactions, err := automation.consensusSystem.GetUncompressedTransactions(TransactionCompressionBatchSize)
	if err != nil {
		fmt.Printf("Error fetching transactions for compression: %v\n", err)
		return
	}

	// Compress each transaction in the batch
	for _, tx := range transactions {
		err := automation.compressTransaction(tx)
		if err != nil {
			fmt.Printf("Error compressing transaction %s: %v\n", tx.ID, err)
			continue
		}
		automation.logCompressionEvent(tx)
	}
}

// compressTransaction compresses the individual transaction data using raw compression logic.
func (automation *TransactionCompressionAutomation) compressTransaction(tx common.Transaction) error {
	// Check if the transaction size meets the compression threshold
	if len(tx.Data) < TransactionCompressionThreshold {
		return nil // No need to compress if below threshold
	}

	// Compress the transaction data using the raw compression tool
	compressedData, err := automation.compressionTool.Compress(tx.Data)
	if err != nil {
		return fmt.Errorf("failed to compress transaction %s: %w", tx.ID, err)
	}

	// Update the transaction with compressed data
	tx.CompressedData = compressedData
	tx.IsCompressed = true

	// Save the compressed transaction state in the consensus system
	err = automation.consensusSystem.SaveCompressedTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to save compressed transaction %s: %w", tx.ID, err)
	}

	return nil
}

// logCompressionEvent logs the compression event to the ledger with encrypted details.
func (automation *TransactionCompressionAutomation) logCompressionEvent(tx common.Transaction) {
	// Encrypt the log entry for security
	encryptedLogDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Transaction %s compressed successfully", tx.ID)), TransactionCompressionLogKey)
	if err != nil {
		fmt.Printf("Error encrypting compression log for transaction %s: %v\n", tx.ID, err)
		return
	}

	// Create a ledger entry for the transaction compression event
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("TX-COMPRESSION-%s-%d", tx.ID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Compression",
		Status:    "Success",
		Details:   string(encryptedLogDetails),
	}

	// Log the compression event in the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log compression event for transaction %s: %v\n", tx.ID, err)
	}
}
