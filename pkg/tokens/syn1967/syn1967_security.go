package syn1967

import (
	"fmt"
	"time"
	"sync"
)

// TokenSecurityManager handles the security aspects of SYN1967 token management.
type TokenSecurityManager struct {
	privateKey *rsa.PrivateKey // Private key for digital signing
	publicKey  *rsa.PublicKey  // Public key for verification
	mu         sync.Mutex      // To ensure thread-safe security management
}

// NewTokenSecurityManager creates a new TokenSecurityManager and generates RSA keys.
func NewTokenSecurityManager() (*TokenSecurityManager, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	return &TokenSecurityManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// SignToken signs the data of a SYN1967 token, ensuring data integrity and verification.
func (tsm *TokenSecurityManager) SignToken(token *common.SYN1967Token) ([]byte, error) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	tokenData := fmt.Sprintf("%s|%s|%f|%s|%s", token.TokenID, token.CommodityName, token.Amount, token.Owner, token.ExpiryDate.String())
	hashed := sha256.Sum256([]byte(tokenData))
	signature, err := rsa.SignPSS(rand.Reader, tsm.privateKey, sha256.New(), hashed[:], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}
	return signature, nil
}

// VerifyTokenSignature verifies the digital signature of a SYN1967 token.
func (tsm *TokenSecurityManager) VerifyTokenSignature(token *common.SYN1967Token, signature []byte) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	tokenData := fmt.Sprintf("%s|%s|%f|%s|%s", token.TokenID, token.CommodityName, token.Amount, token.Owner, token.ExpiryDate.String())
	hashed := sha256.Sum256([]byte(tokenData))

	err := rsa.VerifyPSS(tsm.publicKey, sha256.New(), hashed[:], signature, nil)
	if err != nil {
		return errors.New("failed to verify token signature")
	}
	return nil
}

// EncryptTokenData encrypts the sensitive data of a SYN1967 token.
func (tsm *TokenSecurityManager) EncryptTokenData(token *common.SYN1967Token) ([]byte, error) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	encryptedData, err := encryption.Encrypt(token)
	if err != nil {
		return nil, fmt.Errorf("error encrypting token data: %v", err)
	}
	return encryptedData, nil
}

// DecryptTokenData decrypts the sensitive data of a SYN1967 token.
func (tsm *TokenSecurityManager) DecryptTokenData(encryptedData []byte) (*common.SYN1967Token, error) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	token, err := encryption.Decrypt(encryptedData, &common.SYN1967Token{})
	if err != nil {
		return nil, fmt.Errorf("error decrypting token data: %v", err)
	}
	return token, nil
}

// ValidateSubBlock ensures that the token's transactions are validated within a sub-block.
func (tsm *TokenSecurityManager) ValidateSubBlock(tokenID string, subBlockID string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Validate sub-block transaction
	valid, err := subblock.ValidateSubBlockTransaction(subBlockID, tokenID)
	if err != nil {
		return fmt.Errorf("error validating sub-block transaction: %v", err)
	}
	if !valid {
		return errors.New("sub-block validation failed")
	}
	return nil
}

// RecordTransactionInLedger records secure token transactions into the ledger after encryption.
func (tsm *TokenSecurityManager) RecordTransactionInLedger(token *common.SYN1967Token, transactionType string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Encrypt token data before ledger entry
	encryptedToken, err := tsm.EncryptTokenData(token)
	if err != nil {
		return fmt.Errorf("error encrypting token data before ledger entry: %v", err)
	}

	// Record transaction in the ledger
	err = ledger.RecordTransaction(token.TokenID, encryptedToken, transactionType)
	if err != nil {
		return fmt.Errorf("error recording transaction in ledger: %v", err)
	}

	return nil
}

// AuditTrail generates an audit trail for all token actions including encryption, validation, and transactions.
func (tsm *TokenSecurityManager) AuditTrail(tokenID string) (*common.AuditReport, error) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Retrieve audit records from the ledger
	auditRecords, err := ledger.RetrieveAuditRecords(tokenID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving audit trail for token: %v", err)
	}

	// Compile the audit report
	auditReport := &common.AuditReport{
		TokenID:      tokenID,
		Timestamp:    time.Now(),
		AuditRecords: auditRecords,
	}

	return auditReport, nil
}

// RevokeToken revokes the ownership and permissions of a SYN1967 token, useful for emergency cases.
func (tsm *TokenSecurityManager) RevokeToken(tokenID string, reason string) error {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()

	// Retrieve token from the ledger
	token, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token for revocation: %v", err)
	}

	// Mark the token as revoked and add the reason
	token.RevocationStatus = true
	token.RevocationReason = reason

	// Encrypt updated token data
	encryptedToken, err := tsm.EncryptTokenData(token)
	if err != nil {
		return fmt.Errorf("error encrypting revoked token data: %v", err)
	}

	// Update the ledger with the revoked token
	err = ledger.StoreToken(tokenID, encryptedToken)
	if err != nil {
		return fmt.Errorf("error storing revoked token in ledger: %v", err)
	}

	return nil
}

// generateUniqueID generates a unique identifier for token-related actions.
func generateUniqueID() string {
	id, _ := rand.Int(rand.Reader, big.NewInt(1e12))
	return fmt.Sprintf("SYN-SEC-%d", id)
}
