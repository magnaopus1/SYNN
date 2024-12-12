package syn722

import (
	"errors"
	"sync"
	"time"

)

// SYN722Transaction represents a transaction for SYN722 tokens
type SYN722Transaction struct {
	ID            string                 `json:"id"`
	TokenID       string                 `json:"token_id"`
	From          string                 `json:"from"`
	To            string                 `json:"to"`
	Quantity      uint64                 `json:"quantity"`
	Mode          string                 `json:"mode"` // "fungible" or "non-fungible"
	Status        string                 `json:"status"` // "pending", "completed", "failed"
	Timestamp     time.Time              `json:"timestamp"`
	Details       map[string]interface{} `json:"details"` // Any additional transaction details
	EncryptedData string                 `json:"encrypted_data,omitempty"`
	EncryptionKey string                 `json:"encryption_key,omitempty"`
}

// SYN722TransactionManager handles all transaction-related operations for SYN722 tokens
type SYN722TransactionManager struct {
	Ledger            *ledger.Ledger                // Ledger for recording transactions
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for validating transactions
	EncryptionService *encryption.EncryptionService // Encryption service for securing transaction data
	mutex             sync.Mutex                    // Mutex for safe concurrent access
}

// NewSYN722TransactionManager initializes a new instance of SYN722TransactionManager
func NewSYN722TransactionManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN722TransactionManager {
	return &SYN722TransactionManager{
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
	}
}

// CreateTransaction initializes and validates a new transaction, handling both fungible and non-fungible transfers
func (tm *SYN722TransactionManager) CreateTransaction(tokenID, from, to string, quantity uint64, mode string, details map[string]interface{}, encrypt bool) (*SYN722Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a unique transaction ID
	txID := generateTransactionID(tokenID, from, to)

	// Encrypt transaction details if requested
	var encryptedData, encryptionKey string
	var err error
	if encrypt {
		encryptedData, encryptionKey, err = tm.EncryptionService.EncryptData([]byte(common.MapToString(details)))
		if err != nil {
			return nil, errors.New("failed to encrypt transaction details")
		}
	}

	// Create the transaction object
	tx := &SYN722Transaction{
		ID:            txID,
		TokenID:       tokenID,
		From:          from,
		To:            to,
		Quantity:      quantity,
		Mode:          mode,
		Status:        "pending",
		Timestamp:     time.Now(),
		Details:       details,
		EncryptedData: encryptedData,
		EncryptionKey: encryptionKey,
	}

	// Validate the transaction through Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateTransaction(tx); err != nil {
		tx.Status = "failed"
		return nil, errors.New("transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordTransaction(tx.ID, tx); err != nil {
		return nil, errors.New("failed to record transaction in the ledger")
	}

	tx.Status = "completed"
	return tx, nil
}

// TransferToken processes a transfer of SYN722 tokens between accounts, including both fungible and non-fungible transfers
func (tm *SYN722TransactionManager) TransferToken(tokenID, from, to string, quantity uint64) (*SYN722Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve token from ledger")
	}

	// Ensure the token has sufficient quantity for the transfer
	if token.Mode == "fungible" && token.Quantity < quantity {
		return nil, errors.New("insufficient token quantity for transfer")
	} else if token.Mode == "non-fungible" && quantity > 1 {
		return nil, errors.New("non-fungible tokens cannot be split or partially transferred")
	}

	// Create the transaction for transferring tokens
	details := map[string]interface{}{
		"action":  "transfer",
		"from":    from,
		"to":      to,
		"quantity": quantity,
	}
	tx, err := tm.CreateTransaction(tokenID, from, to, quantity, token.Mode, details, false)
	if err != nil {
		return nil, err
	}

	// Update the token's owner (for non-fungible) or reduce the quantity (for fungible)
	if token.Mode == "fungible" {
		token.Quantity -= quantity
	} else {
		token.Owner = to
	}

	// Update the token in the ledger
	if err := tm.Ledger.UpdateToken(tokenID, token); err != nil {
		return nil, errors.New("failed to update token in the ledger")
	}

	return tx, nil
}

// RetrieveTransaction fetches a specific transaction by its ID and decrypts it if necessary
func (tm *SYN722TransactionManager) RetrieveTransaction(txID string) (*SYN722Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	tx, err := tm.Ledger.GetTransaction(txID)
	if err != nil {
		return nil, errors.New("failed to retrieve transaction from ledger")
	}

	// Decrypt the transaction details if it was encrypted
	if tx.EncryptedData != "" {
		decryptedData, err := tm.EncryptionService.DecryptData([]byte(tx.EncryptedData), tx.EncryptionKey)
		if err != nil {
			return nil, errors.New("failed to decrypt transaction details")
		}
		tx.Details = common.StringToMap(string(decryptedData))
	}

	return tx, nil
}

// ListTransactionsByToken retrieves all transactions related to a specific token
func (tm *SYN722TransactionManager) ListTransactionsByToken(tokenID string) ([]*SYN722Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve all transactions for the given tokenID from the ledger
	transactions, err := tm.Ledger.GetTransactionsByToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve transactions for token from ledger")
	}

	// Decrypt any encrypted transactions
	for _, tx := range transactions {
		if tx.EncryptedData != "" {
			decryptedData, err := tm.EncryptionService.DecryptData([]byte(tx.EncryptedData), tx.EncryptionKey)
			if err != nil {
				return nil, errors.New("failed to decrypt transaction details for transaction ID: " + tx.ID)
			}
			tx.Details = common.StringToMap(string(decryptedData))
		}
	}

	return transactions, nil
}

// generateTransactionID creates a unique transaction ID based on token ID, sender, and receiver
func generateTransactionID(tokenID, from, to string) string {
	return tokenID + "_" + from + "_" + to + "_" + time.Now().Format("20060102150405")
}
