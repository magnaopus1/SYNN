package syn1401

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// SYN1401Management represents the management layer for SYN1401 investment tokens.
type SYN1401Management struct {
	Ledger common.LedgerInterface // Ledger interface for interacting with the blockchain ledger.
	Keys   map[string][]byte      // Encryption keys mapped by the token owner.
}

// CreateSYN1401Token creates and registers a new SYN1401Token with encrypted metadata, and stores it in the ledger.
func (m *SYN1401Management) CreateSYN1401Token(
	owner string, principalAmount float64, interestRate float64, startDate time.Time, 
	maturityDate time.Time, compoundInterest bool, investmentPurpose string) (*common.SYN1401Token, error) {

	tokenID := generateUniqueID()
	token := &common.SYN1401Token{
		TokenID:                tokenID,
		Owner:                  owner,
		PrincipalAmount:         principalAmount,
		InterestRate:            interestRate,
		StartDate:               startDate,
		MaturityDate:            maturityDate,
		AccruedInterest:         0.0, // Start with no accrued interest.
		InterestPaymentSchedule: "Annual", // Default to annual payments.
		CompoundInterest:        compoundInterest,
		RedemptionStatus:        "Active",
		InvestmentPurpose:       investmentPurpose,
		IsTransferable:          true,
		ComplianceStatus:        "Pending Audit",
	}

	// Encrypt sensitive metadata
	encryptionKey, exists := m.Keys[owner]
	if !exists {
		return nil, errors.New("encryption key for owner not found")
	}
	encryptedMetadata, err := m.encryptMetadata(encryptionKey, token)
	if err != nil {
		return nil, err
	}
	token.EncryptedMetadata = encryptedMetadata

	// Create initial event log
	token.EventLogs = append(token.EventLogs, common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Creation",
		Description: fmt.Sprintf("SYN1401Token created with principal of %.2f for %s", principalAmount, owner),
		EventDate:   time.Now(),
		PerformedBy: owner,
	})

	// Save token in the ledger
	err = m.Ledger.SaveToken(tokenID, token)
	if err != nil {
		return nil, fmt.Errorf("error saving token to ledger: %w", err)
	}

	return token, nil
}

// RedeemToken processes the redemption of a SYN1401Token, including early redemptions with penalties if applicable.
func (m *SYN1401Management) RedeemToken(
	tokenID string, principalPaid float64, interestPaid float64, isEarly bool, penalty float64, performedBy string) error {

	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	if token.RedemptionStatus == "Redeemed" {
		return errors.New("token already redeemed")
	}

	redemptionType := "Matured"
	if isEarly {
		redemptionType = "Early"
	}

	redemptionLog := common.RedemptionLog{
		RedemptionID:   generateUniqueID(),
		RedemptionType: redemptionType,
		PrincipalPaid:  principalPaid,
		InterestPaid:   interestPaid,
		RedemptionDate: time.Now(),
		PerformedBy:    performedBy,
		PenaltyApplied: isEarly,
		PenaltyDetails: fmt.Sprintf("Penalty of %.2f applied", penalty),
	}

	token.RedemptionLogs = append(token.RedemptionLogs, redemptionLog)
	token.RedemptionStatus = "Redeemed"

	// Update token in the ledger
	err = m.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return fmt.Errorf("error updating token in ledger: %w", err)
	}

	// Create event log for redemption
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Redemption",
		Description: fmt.Sprintf("Token redeemed by %s. Principal: %.2f, Interest: %.2f", performedBy, principalPaid, interestPaid),
		EventDate:   time.Now(),
		PerformedBy: performedBy,
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	return nil
}

// LogAccrual logs accrued interest on a SYN1401Token and updates its accrual history in the ledger.
func (m *SYN1401Management) LogAccrual(tokenID string, accruedAmount float64, accrualMethod string) error {
	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	token.AccruedInterest += accruedAmount

	accrualLog := common.AccrualLog{
		AccrualID:     generateUniqueID(),
		AccruedAmount: accruedAmount,
		AccrualDate:   time.Now(),
		AccrualMethod: accrualMethod,
		InterestRate:  token.InterestRate,
	}

	token.AccrualHistory = append(token.AccrualHistory, accrualLog)

	// Update token in the ledger
	err = m.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return fmt.Errorf("error updating token accrual in ledger: %w", err)
	}

	// Log event for accrual
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Accrual",
		Description: fmt.Sprintf("Accrued %.2f interest for token %s", accruedAmount, tokenID),
		EventDate:   time.Now(),
		PerformedBy: "System", // System typically accrues interest automatically
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	return nil
}

// TransferToken allows for the transfer of ownership of a SYN1401Token to a new owner.
func (m *SYN1401Management) TransferToken(tokenID string, newOwner string, performedBy string) error {
	token, err := m.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token from ledger: %w", err)
	}

	if !token.IsTransferable {
		return errors.New("token is not transferable")
	}

	token.Owner = newOwner

	// Update token in the ledger
	err = m.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return fmt.Errorf("error updating token in ledger: %w", err)
	}

	// Log transfer event
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Transfer",
		Description: fmt.Sprintf("Token transferred to %s by %s", newOwner, performedBy),
		EventDate:   time.Now(),
		PerformedBy: performedBy,
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	return nil
}

// encryptMetadata encrypts the sensitive metadata for a SYN1401Token.
func (m *SYN1401Management) encryptMetadata(key []byte, token *common.SYN1401Token) ([]byte, error) {
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

// generateUniqueID generates a unique identifier for events, tokens, etc.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
