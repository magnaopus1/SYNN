package transactions

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// TransactionSearchService handles searching transactions based on different criteria.
type TransactionSearchService struct {
	Ledger       *ledger.Ledger                     // Reference to the ledger to search transactions
	Cache        map[string]*common.Transaction     // Cache to store frequently accessed transactions
	Index        map[string][]*common.Transaction   // Index to improve search performance
	cacheMutex   sync.RWMutex                       // Mutex for thread-safe cache access
	mutex        sync.RWMutex                       // Mutex for safe read/write operations
	cacheEnabled bool                               // Toggle cache usage
}

// TransactionSearchCriteria defines the criteria for searching transactions.
type TransactionSearchCriteria struct {
	TransactionID string    // Filter by Transaction ID
	SenderID      string    // Filter by Sender ID
	RecipientID   string    // Filter by Recipient ID
	DateFrom      time.Time // Filter by start date (inclusive)
	DateTo        time.Time // Filter by end date (inclusive)
	MinAmount     float64   // Filter by minimum transaction amount
	MaxAmount     float64   // Filter by maximum transaction amount
	Status        string    // Filter by transaction status (e.g., "Pending", "Confirmed")
}

// NewTransactionSearchService initializes a new TransactionSearchService.
func NewTransactionSearchService(ledger *ledger.Ledger, cacheEnabled bool) *TransactionSearchService {
	return &TransactionSearchService{
		Ledger:       ledger,
		Cache:        make(map[string]*common.Transaction),
		Index:        make(map[string][]*common.Transaction),
		cacheEnabled: cacheEnabled,
	}
}

// SearchByTransactionID searches for a transaction by its unique transaction ID.
func (tss *TransactionSearchService) SearchByTransactionID(txID string) (*common.Transaction, error) {
    // Acquire read lock for thread-safe access
    tss.mutex.RLock()
    defer tss.mutex.RUnlock()

    // Lookup the transaction in the ledger using only the transaction ID
    txRecord, err := tss.Ledger.GetTransactionByID(txID)
    if err != nil {
        return nil, fmt.Errorf("transaction with ID %s not found: %v", txID, err)
    }

    // Map ledger.TransactionRecord to common.Transaction
    commonTx := &common.Transaction{
        TransactionID:   txRecord.ID,           // Assuming ID is the transaction ID in TransactionRecord
        FromAddress:     txRecord.From,         // Mapping From in TransactionRecord to FromAddress in Transaction
        ToAddress:       txRecord.To,           // Mapping To in TransactionRecord to ToAddress in Transaction
        Amount:          txRecord.Amount,
        Fee:             txRecord.Fee,
        Timestamp:       txRecord.Timestamp,
        BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight), // Assuming BlockHeight corresponds to BlockID in some way
        ValidatorID:     txRecord.ValidatorID,
        Status:          txRecord.Status,
        ExecutionResult: txRecord.Action,       // Assuming Action corresponds to execution result
    }

    return commonTx, nil
}

// SearchByBlockHeight searches for all transactions executed at a specific block height.
func (tss *TransactionSearchService) SearchByBlockHeight(blockHeight uint64) ([]*common.Transaction, error) {
    tss.mutex.RLock()
    defer tss.mutex.RUnlock()

    // Convert blockHeight from uint64 to int to match the expected type
    txRecords, err := tss.Ledger.GetTransactionsByBlockHeight(int(blockHeight))
    if err != nil {
        return nil, fmt.Errorf("no transactions found at block height %d: %v", blockHeight, err)
    }

    // Create an array of common.Transaction to return
    commonTxns := make([]*common.Transaction, len(txRecords))

    // Map each ledger.TransactionRecord to common.Transaction
    for i, txRecord := range txRecords {
        commonTxns[i] = &common.Transaction{
            TransactionID:   txRecord.ID,           // Assuming ID is the transaction ID
            FromAddress:     txRecord.From,         // Map From in TransactionRecord to FromAddress
            ToAddress:       txRecord.To,           // Map To in TransactionRecord to ToAddress
            Amount:          txRecord.Amount,
            Fee:             txRecord.Fee,
            Timestamp:       txRecord.Timestamp,
            BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight), // Assuming BlockHeight corresponds to BlockID
            ValidatorID:     txRecord.ValidatorID,
            Status:          txRecord.Status,
            ExecutionResult: txRecord.Action,       // Assuming Action corresponds to execution result
        }
    }

    return commonTxns, nil
}


// SearchByValidatorID searches for all transactions that were validated by a specific validator.
func (tss *TransactionSearchService) SearchByValidatorID(validatorID string) ([]*common.Transaction, error) {
	tss.mutex.RLock()
	defer tss.mutex.RUnlock()

	// Lookup transactions from the ledger by validator ID
	txRecords, err := tss.Ledger.GetTransactionsByValidator(validatorID)
	if err != nil {
		return nil, fmt.Errorf("no transactions found for validator %s: %v", validatorID, err)
	}

	// Create encryption instance for decrypting transactions
	encryptionInstance := &common.Encryption{}

	// Convert and decrypt all transaction data
	commonTxns := make([]*common.Transaction, len(txRecords))
	for i, txRecord := range txRecords {
		// Decrypt the transaction data (assuming transaction details are stored in the Details field)
		decryptedTx, err := encryptionInstance.DecryptData([]byte(txRecord.Details), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting transaction %s: %v", txRecord.ID, err)
		}

		// Map TransactionRecord to common.Transaction
		commonTxns[i] = &common.Transaction{
			TransactionID:   txRecord.ID,
			FromAddress:     txRecord.From,
			ToAddress:       txRecord.To,
			Amount:          txRecord.Amount,
			Fee:             txRecord.Fee,
			Timestamp:       txRecord.Timestamp,
			BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight),
			ValidatorID:     txRecord.ValidatorID,
			Status:          txRecord.Status,
			ExecutionResult: txRecord.Action,
			DecryptedData:   string(decryptedTx), // Convert []byte to string
		}
	}

	return commonTxns, nil
}



// SearchByTimeRange searches for transactions within a specified time range.
func (tss *TransactionSearchService) SearchByTimeRange(startTime, endTime time.Time) ([]*common.Transaction, error) {
	tss.mutex.RLock()
	defer tss.mutex.RUnlock()

	// Lookup transactions from the ledger by time range
	txRecords, err := tss.Ledger.GetTransactionsByTimeRange(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("no transactions found between %v and %v: %v", startTime, endTime, err)
	}

	// Create encryption instance for decrypting transactions
	encryptionInstance := &common.Encryption{}

	// Convert and decrypt all transaction data
	commonTxns := make([]*common.Transaction, len(txRecords))
	for i, txRecord := range txRecords {
		// Decrypt the transaction data (assuming transaction details are stored in the Details field)
		decryptedTx, err := encryptionInstance.DecryptData([]byte(txRecord.Details), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting transaction %s: %v", txRecord.ID, err)
		}

		// Map TransactionRecord to common.Transaction
		commonTxns[i] = &common.Transaction{
			TransactionID:   txRecord.ID,
			FromAddress:     txRecord.From,
			ToAddress:       txRecord.To,
			Amount:          txRecord.Amount,
			Fee:             txRecord.Fee,
			Timestamp:       txRecord.Timestamp,
			BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight),
			ValidatorID:     txRecord.ValidatorID,
			Status:          txRecord.Status,
			ExecutionResult: txRecord.Action,
			DecryptedData:   string(decryptedTx), // Convert []byte to string
		}
	}

	return commonTxns, nil
}



// SearchByCondition searches for transactions based on a condition (e.g., smart contract state).
func (tss *TransactionSearchService) SearchByCondition(condition string) ([]*common.Transaction, error) {
	tss.mutex.RLock()
	defer tss.mutex.RUnlock()

	// Define a function that checks if a transaction record matches the condition
	matchCondition := func(txRecord ledger.TransactionRecord) bool {
		// Example: Check if the condition is a substring of the transaction details
		return strings.Contains(txRecord.Details, condition)
	}

	// Lookup transactions from the ledger by condition function
	txRecords, err := tss.Ledger.GetTransactionsByCondition(matchCondition)
	if err != nil {
		return nil, fmt.Errorf("no transactions found matching condition %s: %v", condition, err)
	}

	// Create encryption instance for decrypting transactions
	encryptionInstance := &common.Encryption{}

	// Convert and decrypt all transaction data
	commonTxns := make([]*common.Transaction, len(txRecords))
	for i, txRecord := range txRecords {
		// Decrypt the transaction data (assuming transaction details are stored in the Details field)
		decryptedTx, err := encryptionInstance.DecryptData([]byte(txRecord.Details), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting transaction %s: %v", txRecord.ID, err)
		}

		// Map TransactionRecord to common.Transaction
		commonTxns[i] = &common.Transaction{
			TransactionID:   txRecord.ID,
			FromAddress:     txRecord.From,
			ToAddress:       txRecord.To,
			Amount:          txRecord.Amount,
			Fee:             txRecord.Fee,
			Timestamp:       txRecord.Timestamp,
			BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight),
			ValidatorID:     txRecord.ValidatorID,
			Status:          txRecord.Status,
			ExecutionResult: txRecord.Action,
			DecryptedData:   string(decryptedTx), // Convert []byte to string
		}
	}

	return commonTxns, nil
}



// SearchByStatus searches for transactions based on their execution status (e.g., pending, executed).
func (tss *TransactionSearchService) SearchByStatus(status string) ([]*common.Transaction, error) {
	tss.mutex.RLock()
	defer tss.mutex.RUnlock()

	// Lookup transactions from the ledger by status
	txRecords, err := tss.Ledger.GetTransactionsByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("no transactions found with status %s: %v", status, err)
	}

	// Create encryption instance for decrypting transactions
	encryptionInstance := &common.Encryption{}

	// Convert and decrypt all transaction data
	commonTxns := make([]*common.Transaction, len(txRecords))
	for i, txRecord := range txRecords {
		// Decrypt the transaction data (assuming transaction details are stored in the Details field)
		decryptedTx, err := encryptionInstance.DecryptData([]byte(txRecord.Details), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting transaction %s: %v", txRecord.ID, err)
		}

		// Map TransactionRecord to common.Transaction
		commonTxns[i] = &common.Transaction{
			TransactionID:   txRecord.ID,
			FromAddress:     txRecord.From,
			ToAddress:       txRecord.To,
			Amount:          txRecord.Amount,
			Fee:             txRecord.Fee,
			Timestamp:       txRecord.Timestamp,
			BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight),
			ValidatorID:     txRecord.ValidatorID,
			Status:          txRecord.Status,
			ExecutionResult: txRecord.Action,
			DecryptedData:   string(decryptedTx), // Convert []byte to string
		}
	}

	return commonTxns, nil
}


// SearchByGasFeeRange searches for transactions within a specified gas fee range.
func (tss *TransactionSearchService) SearchByGasFeeRange(minGasFee, maxGasFee float64) ([]*common.Transaction, error) {
	tss.mutex.RLock()
	defer tss.mutex.RUnlock()

	// Cast the minGasFee and maxGasFee to uint64 as required by GetTransactionsByGasFeeRange
	minGasFeeUint := uint64(minGasFee)
	maxGasFeeUint := uint64(maxGasFee)

	// Lookup transactions from the ledger by gas fee range
	txRecords, err := tss.Ledger.GetTransactionsByGasFeeRange(minGasFeeUint, maxGasFeeUint)
	if err != nil {
		return nil, fmt.Errorf("no transactions found with gas fee between %.10f and %.10f: %v", minGasFee, maxGasFee, err)
	}

	// Create encryption instance for decrypting transactions
	encryptionInstance := &common.Encryption{}

	// Prepare the list of common transactions
	commonTxns := make([]*common.Transaction, len(txRecords))

	// Decrypt and map all transaction data
	for i, txRecord := range txRecords {
		// Decrypt the transaction data (assuming it's stored in Details)
		decryptedTx, err := encryptionInstance.DecryptData([]byte(txRecord.Details), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("error decrypting transaction %s: %v", txRecord.ID, err)
		}

		// Map TransactionRecord to common.Transaction
		commonTxns[i] = &common.Transaction{
			TransactionID:   txRecord.ID,
			FromAddress:     txRecord.From,
			ToAddress:       txRecord.To,
			Amount:          txRecord.Amount,
			Fee:             txRecord.Fee,
			Timestamp:       txRecord.Timestamp,
			BlockID:         fmt.Sprintf("%d", txRecord.BlockHeight),
			ValidatorID:     txRecord.ValidatorID,
			Status:          txRecord.Status,
			ExecutionResult: txRecord.Action,
			DecryptedData:   string(decryptedTx), // Convert []byte to string
		}
	}

	return commonTxns, nil
}
