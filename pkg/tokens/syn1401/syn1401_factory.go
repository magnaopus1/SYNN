package syn1401

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

// TokenFactory is responsible for creating and managing SYN1401 tokens.
type TokenFactory struct {
	Ledger common.LedgerInterface // Interface to interact with the blockchain ledger
	Keys   map[string][]byte      // Encryption keys for each token owner
}

// SYN1401Token represents an investment token under the SYN1401 standard.
type SYN1401Token struct {
	TokenID                 string            // Unique identifier for the token
	Owner                   string            // The owner of the token (issuer or investor)
	PrincipalAmount          float64           // The principal amount (investment) associated with the token
	InterestRate             float64           // The fixed interest rate (annual)
	StartDate                time.Time         // Date when the investment starts
	MaturityDate             time.Time         // Date when the investment matures (redemption date)
	AccruedInterest          float64           // Total accrued interest over the investment period
	InterestPaymentSchedule  string            // Schedule for interest payments (e.g., monthly, quarterly)
	CompoundInterest         bool              // Whether compound interest is applied
	LastInterestPayment      time.Time         // Timestamp of the last interest payment
	NextInterestPayment      time.Time         // Timestamp for the next interest payment due
	RedemptionStatus         string            // Current status of the token regarding redemption (e.g., "Active", "Matured", "Redeemed", "Early Redemption")
	CustomRedemptionConditions map[string]string // Custom conditions related to token redemption
	Collateralized           bool              // Whether the token is backed by collateral
	CollateralDetails        map[string]string // Details about collateral (if applicable)
	MarketValue              float64           // Current market value of the token (if applicable)
	InterestAccrualMode      string            // Interest accrual method (e.g., "Daily", "Monthly", etc.)
	YieldAtMaturity          float64           // The expected yield (interest + principal) at the maturity date
	InvestmentPurpose        string            // Description of the underlying purpose for the investment (e.g., "Bond", "Treasury Bill")
	IsTransferable           bool              // Whether the token can be transferred between owners
	ApprovalRequired         bool              // Whether transfers or certain actions require approval
	RestrictedTransfers      bool              // Whether there are restrictions on token transfers due to regulations
	ComplianceStatus         string            // Compliance status of the token (e.g., "Compliant", "Pending Audit", "Non-Compliant")
	AuditTrail               []AuditLog        // List of audit records related to the token (for compliance and security)
	RedemptionLogs           []RedemptionLog   // Logs for redemption events (early redemptions, maturity payouts)
	AccrualHistory           []AccrualLog      // Historical records of interest accruals
	EventLogs                []EventLog        // Logs of key events (creation, transfers, interest payments)
	EncryptedMetadata        []byte            // Encrypted metadata for sensitive information related to the token
}

// AuditLog represents an audit entry for SYN1401Token, capturing changes and reviews.
type AuditLog struct {
	AuditID      string    // Unique identifier for the audit entry
	PerformedBy  string    // ID of the auditor or system performing the audit
	Description  string    // Description of the audit or compliance check
	Timestamp    time.Time // When the audit was performed
}

// RedemptionLog represents a redemption event for a SYN1401Token, whether early or at maturity.
type RedemptionLog struct {
	RedemptionID    string    // Unique identifier for the redemption event
	RedemptionType  string    // Type of redemption ("Early", "Matured")
	PrincipalPaid   float64   // Amount of principal paid out
	InterestPaid    float64   // Amount of interest paid out
	RedemptionDate  time.Time // When the redemption occurred
	PerformedBy     string    // ID of the entity performing the redemption
	PenaltyApplied  bool      // Whether a penalty was applied for early redemption
	PenaltyDetails  string    // Details about the penalty applied (if any)
}

// AccrualLog represents a record of interest accrual for SYN1401Token.
type AccrualLog struct {
	AccrualID       string    // Unique identifier for the accrual log
	AccruedAmount   float64   // Amount of interest accrued during the period
	AccrualDate     time.Time // Date when the accrual was calculated
	AccrualMethod   string    // Method used for accrual ("Daily", "Compound", etc.)
	InterestRate    float64   // The interest rate applied for this accrual period
}

// EventLog represents a key event in the lifecycle of the SYN1401Token.
type EventLog struct {
	EventID       string    // Unique identifier for the event
	EventType     string    // Type of event ("Creation", "Transfer", "Interest Payment", etc.)
	Description   string    // Description of the event
	EventDate     time.Time // Timestamp when the event occurred
	PerformedBy   string    // ID of the entity or system performing the action
}

// CreateSYN1401Token creates a new SYN1401Token with the given details.
func (tf *TokenFactory) CreateSYN1401Token(owner string, principal float64, interestRate float64, startDate time.Time, maturityDate time.Time, isCompound bool, investmentPurpose string) (*SYN1401Token, error) {
	tokenID := generateUniqueID()
	initialAccruedInterest := 0.0

	token := &SYN1401Token{
		TokenID:                tokenID,
		Owner:                  owner,
		PrincipalAmount:         principal,
		InterestRate:            interestRate,
		StartDate:               startDate,
		MaturityDate:            maturityDate,
		AccruedInterest:         initialAccruedInterest,
		InterestPaymentSchedule: "Annual", // Default to annual
		CompoundInterest:        isCompound,
		RedemptionStatus:        "Active",
		InvestmentPurpose:       investmentPurpose,
		IsTransferable:          true, // Transferable by default
		ComplianceStatus:        "Pending Audit",
	}

	// Encrypt metadata if needed
	encryptionKey, exists := tf.Keys[owner]
	if !exists {
		return nil, errors.New("encryption key for owner not found")
	}
	token.EncryptedMetadata, _ = tf.encryptMetadata(encryptionKey, token)

	// Add creation event log
	token.EventLogs = append(token.EventLogs, EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Creation",
		Description: fmt.Sprintf("SYN1401Token created for %s with principal of %.2f", owner, principal),
		EventDate:   time.Now(),
		PerformedBy: owner,
	})

	// Save the token in the ledger
	tf.Ledger.SaveToken(tokenID, token)

	return token, nil
}

// RedeemToken handles redemption of the token at maturity or early redemption with applicable penalties.
func (tf *TokenFactory) RedeemToken(tokenID string, paidPrincipal float64, paidInterest float64, isEarly bool, penalty float64, performedBy string) error {
	token, err := tf.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	redemptionType := "Matured"
	if isEarly {
		redemptionType = "Early"
	}

	redemptionLog := RedemptionLog{
		RedemptionID:   generateUniqueID(),
		RedemptionType: redemptionType,
		PrincipalPaid:  paidPrincipal,
		InterestPaid:   paidInterest,
		RedemptionDate: time.Now(),
		PerformedBy:    performedBy,
		PenaltyApplied: isEarly,
		PenaltyDetails: fmt.Sprintf("Penalty of %.2f applied", penalty),
	}

	token.RedemptionLogs = append(token.RedemptionLogs, redemptionLog)
	token.RedemptionStatus = "Redeemed"

	// Update in the ledger
	tf.Ledger.UpdateToken(tokenID, token)

	return nil
}

// LogAccrual updates the token's accrued interest and adds an accrual log entry.
func (tf *TokenFactory) LogAccrual(tokenID string, accruedAmount float64, accrualMethod string) error {
	token, err := tf.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	token.AccruedInterest += accruedAmount
	accrualLog := AccrualLog{
		AccrualID:     generateUniqueID(),
		AccruedAmount: accruedAmount,
		AccrualDate:   time.Now(),
		AccrualMethod: accrualMethod,
		InterestRate:  token.InterestRate,
	}

	token.AccrualHistory = append(token.AccrualHistory, accrualLog)

	// Update in the ledger
	tf.Ledger.UpdateToken(tokenID, token)

	return nil
}

// encryptMetadata encrypts the sensitive metadata of a SYN1401Token.
func (tf *TokenFactory) encryptMetadata(key []byte, token *SYN1401Token) ([]byte, error) {
	plaintext := []byte(fmt.Sprintf("%s|%f|%s", token.TokenID, token.PrincipalAmount, token.Owner))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// generateUniqueID generates a unique ID for token and event identification.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
