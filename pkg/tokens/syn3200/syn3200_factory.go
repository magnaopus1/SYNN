package syn3200

import (
	"time"
	"errors"
	"sync"

)

// SYN3200Token represents a bill token in the SYN3200 standard.
type SYN3200Token struct {
	TokenID          string        `json:"token_id"`         // Unique token identifier
	BillMetadata     BillMetadata  `json:"metadata"`         // Metadata related to the bill
	mutex            sync.Mutex                              // Mutex for thread-safe operations
	ledgerService    *ledger.Ledger                         // Ledger service for tracking token events
	encryptionService *encryption.Encryptor                 // Encryption service for securing the token data
	consensusService *consensus.SynnergyConsensus           // Consensus service for transaction validation
}

// NewSYN3200Token creates a new instance of SYN3200Token.
func NewSYN3200Token(tokenID string, billMetadata BillMetadata, ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *SYN3200Token {
	return &SYN3200Token{
		TokenID:          tokenID,
		BillMetadata:     billMetadata,
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// TokenManager manages SYN3200 tokens, including creation, validation, and transaction handling.
type TokenManager struct {
	ledgerManager      *ledger.Ledger           // Ledger manager for tracking token ownership and transactions
	transactionManager *TransactionManager      // Manager for handling token-related transactions
	mutex              sync.Mutex               // Mutex for thread-safe operations
}

// NewTokenManager creates a new TokenManager instance.
func NewTokenManager(ledgerManager *ledger.Ledger, transactionManager *TransactionManager) *TokenManager {
	return &TokenManager{
		ledgerManager:      ledgerManager,
		transactionManager: transactionManager,
	}
}

// BillMetadata represents the detailed metadata associated with a bill token.
type BillMetadata struct {
	BillID           string    `json:"bill_id"`            // Unique identifier for the bill
	Issuer           string    `json:"issuer"`             // The entity issuing the bill
	Payer            string    `json:"payer"`              // The entity responsible for paying the bill
	OriginalAmount   float64   `json:"original_amount"`    // The original bill amount
	RemainingAmount  float64   `json:"remaining_amount"`   // The remaining balance of the bill
	DueDate          time.Time `json:"due_date"`           // The date by which the bill must be paid
	PaidStatus       bool      `json:"paid_status"`        // Whether the bill is paid or not
	TermsConditions  string    `json:"terms_conditions"`   // Any terms or conditions attached to the bill
	Timestamp        time.Time `json:"timestamp"`          // The time when the bill was created or updated
}

// CreateBill creates a new SYN3200Token with the specified BillMetadata.
func (tm *TokenManager) CreateBill(tokenID string, metadata BillMetadata) (*SYN3200Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate if a token with the same ID already exists.
	if tm.ledgerManager.TokenExists(tokenID) {
		return nil, errors.New("token already exists")
	}

	// Create a new SYN3200Token.
	token := NewSYN3200Token(tokenID, metadata, tm.ledgerManager, tm.transactionManager.encryptionService, tm.transactionManager.consensusService)

	// Encrypt the token metadata for security.
	encryptedMetadata, err := token.encryptionService.EncryptData(&metadata)
	if err != nil {
		return nil, err
	}

	// Log the token creation in the ledger.
	err = tm.ledgerManager.LogEvent("BillCreated", time.Now(), tokenID)
	if err != nil {
		return nil, err
	}

	// Store the token data in the ledger.
	err = tm.ledgerManager.StoreToken(tokenID, encryptedMetadata)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ValidateBill validates the authenticity and integrity of a bill token.
func (tm *TokenManager) ValidateBill(tokenID string) (bool, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from the ledger.
	encryptedMetadata, err := tm.ledgerManager.RetrieveToken(tokenID)
	if err != nil {
		return false, err
	}

	// Decrypt the token metadata.
	metadata, err := tm.transactionManager.encryptionService.DecryptData(encryptedMetadata)
	if err != nil {
		return false, err
	}

	// Validate the token using the Synnergy Consensus mechanism.
	err = tm.transactionManager.consensusService.ValidateSubBlock(tokenID)
	if err != nil {
		return false, err
	}

	// Verify the bill status.
	bill := metadata.(*BillMetadata)
	if bill.PaidStatus {
		return true, nil
	}

	return false, nil
}

// PayBill marks a bill token as paid and updates the ledger.
func (tm *TokenManager) PayBill(tokenID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from the ledger.
	encryptedMetadata, err := tm.ledgerManager.RetrieveToken(tokenID)
	if err != nil {
		return err
	}

	// Decrypt the token metadata.
	metadata, err := tm.transactionManager.encryptionService.DecryptData(encryptedMetadata)
	if err != nil {
		return err
	}

	// Update the bill status to paid.
	bill := metadata.(*BillMetadata)
	if bill.PaidStatus {
		return errors.New("bill is already paid")
	}
	bill.PaidStatus = true

	// Encrypt the updated metadata.
	updatedEncryptedMetadata, err := tm.transactionManager.encryptionService.EncryptData(bill)
	if err != nil {
		return err
	}

	// Update the ledger with the paid status.
	err = tm.ledgerManager.LogEvent("BillPaid", time.Now(), tokenID)
	if err != nil {
		return err
	}

	// Store the updated metadata back into the ledger.
	err = tm.ledgerManager.StoreToken(tokenID, updatedEncryptedMetadata)
	if err != nil {
		return err
	}

	return nil
}
