package syn3900

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"sync"

)

// Syn3900Token represents a benefit token on the SYN3900 token standard.
type Syn3900Token struct {
	TokenID           string               `json:"token_id"`           // Unique token identifier
	Metadata          Syn3900Metadata      `json:"metadata"`           // Detailed metadata for the benefit token
	AllocationHistory []BenefitAllocation  `json:"allocation_history"` // History of all benefit allocations
	TransactionHistory []BenefitTransaction `json:"transaction_history"` // History of benefit-related transactions
	OwnershipHistory  []OwnershipChange    `json:"ownership_history"`  // History of ownership changes
	mutex             sync.Mutex           // Mutex for thread-safe operations
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
}

// Syn3900Metadata defines the metadata for benefit tokens, capturing comprehensive benefit details.
type Syn3900Metadata struct {
	BenefitName        string    `json:"benefit_name"`        // Name of the benefit (e.g., Pension, Welfare)
	BenefitType        string    `json:"benefit_type"`        // Type of benefit (e.g., Recurring, One-time)
	Recipient          string    `json:"recipient"`           // The recipient of the benefit
	Amount             float64   `json:"amount"`              // The current benefit amount
	ValidFrom          time.Time `json:"valid_from"`          // Start date for the benefit
	ValidUntil         *time.Time `json:"valid_until"`        // Expiry date of the benefit, if applicable
	IssuedDate         time.Time `json:"issued_date"`         // Date the benefit was issued
	Conditions         string    `json:"conditions"`          // Conditions required for the benefit to be issued
	NextPaymentDate    time.Time `json:"next_payment_date"`   // Next benefit payment date
	NextPaymentAmount  float64   `json:"next_payment_amount"` // Amount for the next payment
	TokenOrCurrency    string    `json:"token_or_currency"`   // The currency or token in which the benefit is issued
	BenefitIssuer      string    `json:"benefit_issuer"`      // Issuer of the benefit (e.g., Government agency)
	Country            string    `json:"country"`             // Country of the benefit
	City               string    `json:"city"`                // City of the benefit
	Jurisdiction       string    `json:"jurisdiction"`        // Legal jurisdiction governing the benefit
	Status             string    `json:"status"`              // Current status of the benefit (e.g., Active, Expired)
}

// BenefitAllocation records the details of benefit allocations.
type BenefitAllocation struct {
	AllocationID       string    `json:"allocation_id"`       // Unique identifier for the allocation
	AllocationDate     time.Time `json:"allocation_date"`     // Date of allocation
	Recipient          string    `json:"recipient"`           // Benefit recipient
	AmountAllocated    float64   `json:"amount_allocated"`    // Amount allocated during this allocation
	UsageBreakdown     string    `json:"usage_breakdown"`     // Description of how the benefit is used
}

// BenefitTransaction represents immutable records of benefit-related transactions.
type BenefitTransaction struct {
	TransactionID      string    `json:"transaction_id"`      // Unique identifier for the transaction
	Timestamp          time.Time `json:"timestamp"`           // Time when the transaction was made
	TransactionType    string    `json:"transaction_type"`    // Type of transaction (e.g., Payment, Adjustment)
	Amount             float64   `json:"amount"`              // Transaction amount
	Recipient          string    `json:"recipient"`           // Benefit recipient involved in the transaction
}

// OwnershipChange tracks the ownership history of benefit tokens.
type OwnershipChange struct {
	OwnershipID        string    `json:"ownership_id"`        // Unique identifier for the ownership change
	OldOwner           string    `json:"old_owner"`           // Previous owner of the benefit token
	NewOwner           string    `json:"new_owner"`           // New owner of the benefit token
	ChangeDate         time.Time `json:"change_date"`         // Date when the ownership change occurred
}

// TokenFactory manages the creation and issuance of SYN3900 benefit tokens.
type TokenFactory struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
}

// NewTokenFactory creates a new TokenFactory.
func NewTokenFactory(ledgerService *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TokenFactory {
	return &TokenFactory{
		ledgerService:     ledgerService,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// CreateToken creates a new SYN3900 benefit token with the given metadata and records it in the ledger.
func (tf *TokenFactory) CreateToken(metadata Syn3900Metadata) (*Syn3900Token, error) {
	tf.ledgerService.Mutex.Lock()
	defer tf.ledgerService.Mutex.Unlock()

	// Generate a unique token ID
	tokenID := generateUniqueTokenID()

	// Create the token object
	token := &Syn3900Token{
		TokenID:           tokenID,
		Metadata:          metadata,
		AllocationHistory: []BenefitAllocation{},
		TransactionHistory: []BenefitTransaction{},
		OwnershipHistory:  []OwnershipChange{},
	}

	// Encrypt the token metadata before storing it
	encryptedToken, err := tf.encryptionService.EncryptData(token)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Store the encrypted token in the ledger
	if err := tf.ledgerService.StoreData(tokenID, encryptedToken); err != nil {
		return nil, fmt.Errorf("failed to store token in ledger: %w", err)
	}

	// Validate the token creation using Synnergy Consensus
	if err := tf.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, fmt.Errorf("failed to validate token in consensus: %w", err)
	}

	// Log the token creation event in the ledger
	tf.logTokenEvent(tokenID, "TokenCreated")

	return token, nil
}

// UpdateNextPayment updates the next payment amount and date for a given SYN3900 token.
func (tf *TokenFactory) UpdateNextPayment(tokenID string, nextPaymentAmount float64, nextPaymentDate time.Time) error {
	tf.ledgerService.Mutex.Lock()
	defer tf.ledgerService.Mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := tf.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %w", err)
	}

	// Update the next payment details
	token.Metadata.NextPaymentAmount = nextPaymentAmount
	token.Metadata.NextPaymentDate = nextPaymentDate

	// Encrypt and store the updated token
	encryptedToken, err := tf.encryptionService.EncryptData(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated token: %w", err)
	}

	if err := tf.ledgerService.StoreData(tokenID, encryptedToken); err != nil {
		return fmt.Errorf("failed to store updated token in ledger: %w", err)
	}

	// Log the payment update event in the ledger
	tf.logTokenEvent(tokenID, "NextPaymentUpdated")

	return nil
}

// RetrieveToken retrieves a SYN3900 token from the ledger and decrypts it.
func (tf *TokenFactory) RetrieveToken(tokenID string) (*Syn3900Token, error) {
	// Retrieve encrypted token data from the ledger
	encryptedData, err := tf.ledgerService.RetrieveData(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token: %w", err)
	}

	// Decrypt the token data
	decryptedToken, err := tf.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}

	return decryptedToken.(*Syn3900Token), nil
}

// generateUniqueTokenID generates a unique identifier for the benefit token.
func generateUniqueTokenID() string {
	// Implement unique ID generation logic (e.g., UUID or timestamp-based)
	return fmt.Sprintf("syn3900-%d", time.Now().UnixNano())
}

// logTokenEvent logs token-related events in the ledger.
func (tf *TokenFactory) logTokenEvent(tokenID, eventType string) {
	_ = tf.ledgerService.LogEvent(eventType, time.Now(), tokenID)
}
