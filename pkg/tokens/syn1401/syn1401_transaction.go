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

// SYN1401TransactionManager manages all transactions related to SYN1401 tokens.
type SYN1401TransactionManager struct {
	Ledger common.LedgerInterface // Interface to interact with the ledger
}

// TransferToken transfers a SYN1401 token from one owner to another, handling ledger update and encryption.
func (tm *SYN1401TransactionManager) TransferToken(tokenID string, fromOwner string, toOwner string) error {
	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token: %w", err)
	}

	// Check if the token is transferable
	if !token.IsTransferable {
		return errors.New("token is not transferable")
	}

	// Ensure compliance with any transfer restrictions
	if token.RestrictedTransfers {
		return errors.New("token is restricted from transfers")
	}

	// Log transfer event
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "Transfer",
		Description: fmt.Sprintf("Token %s transferred from %s to %s", tokenID, fromOwner, toOwner),
		EventDate:   time.Now(),
		PerformedBy: "System",
	}
	token.EventLogs = append(token.EventLogs, eventLog)

	// Update the token's ownership and store the updated token in the ledger
	token.Owner = toOwner
	if err := tm.Ledger.UpdateToken(tokenID, token); err != nil {
		return fmt.Errorf("error updating token in ledger: %w", err)
	}

	return nil
}

// RedeemToken handles the redemption of a SYN1401 token, either at maturity or early redemption.
func (tm *SYN1401TransactionManager) RedeemToken(tokenID string, owner string) (*common.RedemptionLog, error) {
	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving token: %w", err)
	}

	// Check the redemption status
	if token.RedemptionStatus == "Redeemed" {
		return nil, errors.New("token has already been redeemed")
	}

	// If token is not matured, apply early redemption logic
	var principalPaid, interestPaid float64
	redemptionType := "Matured"
	if time.Now().Before(token.MaturityDate) {
		redemptionType = "Early"
		// Apply early redemption penalty (if applicable)
		if token.CustomRedemptionConditions["PenaltyOnEarlyRedemption"] == "true" {
			principalPaid = token.PrincipalAmount * 0.95 // 5% penalty
			interestPaid = token.AccruedInterest * 0.95
		} else {
			principalPaid = token.PrincipalAmount
			interestPaid = token.AccruedInterest
		}
	} else {
		// Matured redemption
		principalPaid = token.PrincipalAmount
		interestPaid = token.AccruedInterest
	}

	// Log redemption event
	redemptionLog := common.RedemptionLog{
		RedemptionID:    generateUniqueID(),
		RedemptionType:  redemptionType,
		PrincipalPaid:   principalPaid,
		InterestPaid:    interestPaid,
		RedemptionDate:  time.Now(),
		PerformedBy:     owner,
		PenaltyApplied:  redemptionType == "Early" && token.CustomRedemptionConditions["PenaltyOnEarlyRedemption"] == "true",
		PenaltyDetails:  "5% penalty applied for early redemption",
	}
	token.RedemptionLogs = append(token.RedemptionLogs, redemptionLog)

	// Update token as redeemed and store back in the ledger
	token.RedemptionStatus = "Redeemed"
	if err := tm.Ledger.UpdateToken(tokenID, token); err != nil {
		return nil, fmt.Errorf("error updating token after redemption: %w", err)
	}

	return &redemptionLog, nil
}

// AccrueInterest calculates the interest for a SYN1401 token based on its accrual method (daily, monthly, etc.).
func (tm *SYN1401TransactionManager) AccrueInterest(tokenID string) error {
	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token: %w", err)
	}

	// Ensure the token hasn't reached its maturity date
	if time.Now().After(token.MaturityDate) {
		return errors.New("token has matured and cannot accrue further interest")
	}

	// Calculate interest based on the accrual method
	var accruedInterest float64
	switch token.InterestAccrualMode {
	case "Daily":
		accruedInterest = token.PrincipalAmount * (token.InterestRate / 365)
	case "Monthly":
		accruedInterest = token.PrincipalAmount * (token.InterestRate / 12)
	default:
		return errors.New("unsupported interest accrual method")
	}

	// Update the accrued interest in the token
	token.AccruedInterest += accruedInterest

	// Log the accrual event
	accrualLog := common.AccrualLog{
		AccrualID:     generateUniqueID(),
		AccruedAmount: accruedInterest,
		AccrualDate:   time.Now(),
		AccrualMethod: token.InterestAccrualMode,
		InterestRate:  token.InterestRate,
	}
	token.AccrualHistory = append(token.AccrualHistory, accrualLog)

	// Update the token in the ledger
	if err := tm.Ledger.UpdateToken(tokenID, token); err != nil {
		return fmt.Errorf("error updating token with accrued interest: %w", err)
	}

	return nil
}

// EncryptSensitiveData encrypts the token's sensitive data before a transaction.
func (tm *SYN1401TransactionManager) EncryptSensitiveData(token *common.SYN1401Token, owner string) error {
	key, err := tm.getOwnerKey(owner)
	if err != nil {
		return err
	}

	plaintext := []byte(fmt.Sprintf("%s|%f|%s", token.TokenID, token.PrincipalAmount, token.Owner))
	encryptedData, err := encryptAES(key, plaintext)
	if err != nil {
		return err
	}

	token.EncryptedMetadata = encryptedData
	return nil
}

// DecryptSensitiveData decrypts the token's sensitive data for viewing.
func (tm *SYN1401TransactionManager) DecryptSensitiveData(token *common.SYN1401Token, owner string) error {
	key, err := tm.getOwnerKey(owner)
	if err != nil {
		return err
	}

	plaintext, err := decryptAES(key, token.EncryptedMetadata)
	if err != nil {
		return err
	}

	_, err = fmt.Sscanf(string(plaintext), "%s|%f|%s", &token.TokenID, &token.PrincipalAmount, &token.Owner)
	if err != nil {
		return fmt.Errorf("error parsing decrypted token data: %w", err)
	}

	return nil
}

// Helper functions for AES encryption and decryption
func encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// Helper function to retrieve owner's encryption key
func (tm *SYN1401TransactionManager) getOwnerKey(owner string) ([]byte, error) {
	ownerInfo, err := tm.Ledger.GetOwnerInfo(owner)
	if err != nil {
		return nil, fmt.Errorf("error retrieving encryption key: %w", err)
	}

	return hex.DecodeString(ownerInfo.EncryptionKey)
}

// Helper function to generate a unique ID
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
