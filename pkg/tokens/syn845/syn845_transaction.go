package syn845

import (
	"errors"
	"sync"
	"time"

)

// TransactionType defines the type of transaction
type TransactionType string

const (
	Issuance     TransactionType = "Issuance"
	Repayment    TransactionType = "Repayment"
	Refinancing  TransactionType = "Refinancing"
	Transfer     TransactionType = "Transfer"
)

// SYN845Transaction represents a transaction related to a debt instrument (SYN845)
type SYN845Transaction struct {
	ID            string                 `json:"id"`
	Type          TransactionType        `json:"type"`
	InstrumentID  string                 `json:"instrument_id"`
	From          string                 `json:"from"`
	To            string                 `json:"to"`
	Amount        float64                `json:"amount"`
	InterestRate  float64                `json:"interest_rate"`
	Date          time.Time              `json:"date"`
	ExtraData     map[string]interface{} `json:"extra_data"`
}

// TransactionManager handles the creation and management of SYN845 transactions
type TransactionManager struct {
	transactions   map[string]*SYN845Transaction
	mu             sync.RWMutex
	Ledger         *ledger.Ledger               // Ledger for recording transactions
	ConsensusEngine *consensus.SynnergyConsensus // Synnergy Consensus for validating transactions
}

// NewTransactionManager creates a new TransactionManager instance
func NewTransactionManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus) *TransactionManager {
	return &TransactionManager{
		transactions:    make(map[string]*SYN845Transaction),
		Ledger:          ledger,
		ConsensusEngine: consensusEngine,
	}
}

// CreateTransaction creates a new SYN845 transaction and validates it through consensus
func (tm *TransactionManager) CreateTransaction(tType TransactionType, instrumentID, from, to string, amount, interestRate float64, extraData map[string]interface{}) (*SYN845Transaction, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Generate a unique ID for the transaction
	transactionID := generateUniqueID()

	transaction := &SYN845Transaction{
		ID:           transactionID,
		Type:         tType,
		InstrumentID: instrumentID,
		From:         from,
		To:           to,
		Amount:       amount,
		InterestRate: interestRate,
		Date:         time.Now(),
		ExtraData:    extraData,
	}

	// Validate the transaction using Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateTransaction(transaction); err != nil {
		return nil, errors.New("transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordTransaction(transaction.ID, transaction); err != nil {
		return nil, errors.New("failed to record transaction in the ledger")
	}

	tm.transactions[transactionID] = transaction
	return transaction, nil
}

// GetTransaction retrieves a SYN845 transaction by its ID
func (tm *TransactionManager) GetTransaction(transactionID string) (*SYN845Transaction, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	transaction, exists := tm.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	return transaction, nil
}

// GetTransactionsByInstrument retrieves all transactions for a specific debt instrument (SYN845)
func (tm *TransactionManager) GetTransactionsByInstrument(instrumentID string) ([]*SYN845Transaction, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var transactions []*SYN845Transaction
	for _, transaction := range tm.transactions {
		if transaction.InstrumentID == instrumentID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

// GetAllTransactions retrieves all SYN845 transactions
func (tm *TransactionManager) GetAllTransactions() []*SYN845Transaction {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var transactions []*SYN845Transaction
	for _, transaction := range tm.transactions {
		transactions = append(transactions, transaction)
	}

	return transactions
}

// generateUniqueID generates a unique ID for the transaction
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
