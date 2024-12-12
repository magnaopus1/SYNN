package syn1967

import (
	"fmt"
	"time"
	"sync"
)

// TokenManager handles operations and management of SYN1967 tokens.
type TokenManager struct {
	mu sync.Mutex // Ensures thread-safe token management
}

// IssueToken issues a new SYN1967 commodity token, records it on the ledger, and returns the issued token.
func (tm *TokenManager) IssueToken(commodityName string, amount float64, unitOfMeasure string, owner string, certification string, traceability string, origin string, expiryDate time.Time) (*common.SYN1967Token, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Validate input data
	if commodityName == "" || amount <= 0 || unitOfMeasure == "" || owner == "" {
		return nil, errors.New("invalid input parameters for issuing token")
	}

	// Generate a new token ID and create the token
	tokenID := generateUniqueTokenID()
	newToken := &common.SYN1967Token{
		TokenID:          tokenID,
		CommodityName:    commodityName,
		Amount:           amount,
		UnitOfMeasure:    unitOfMeasure,
		Owner:            owner,
		Certification:    certification,
		Traceability:     traceability,
		IssuedDate:       time.Now(),
		Origin:           origin,
		ExpiryDate:       expiryDate,
		CollateralStatus: "Uncollateralized", // Default
		Fractionalized:   false,
	}

	// Encrypt sensitive data before storing the token in the ledger
	encryptedToken, err := encryption.Encrypt(newToken)
	if err != nil {
		return nil, fmt.Errorf("error encrypting token data: %v", err)
	}

	// Store token in the ledger
	err = ledger.StoreToken(tokenID, encryptedToken)
	if err != nil {
		return nil, fmt.Errorf("error storing token in ledger: %v", err)
	}

	return newToken, nil
}

// TransferToken transfers ownership of a SYN1967 token from one owner to another.
func (tm *TokenManager) TransferToken(tokenID string, newOwner string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the token from the ledger
	token, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token from ledger: %v", err)
	}

	// Check if the token is fractionalized or restricted from transfer
	if token.Fractionalized || token.RestrictedTransfers {
		return errors.New("token is restricted from transfer or fractionalized")
	}

	// Update the owner of the token
	token.Owner = newOwner

	// Encrypt updated token
	encryptedToken, err := encryption.Encrypt(token)
	if err != nil {
		return fmt.Errorf("error encrypting updated token data: %v", err)
	}

	// Store updated token in the ledger
	err = ledger.StoreToken(tokenID, encryptedToken)
	if err != nil {
		return fmt.Errorf("error updating token ownership in ledger: %v", err)
	}

	return nil
}

// FractionalizeToken allows splitting a SYN1967 token into fractional ownership.
func (tm *TokenManager) FractionalizeToken(tokenID string, fractions []float64, owners []string) ([]*common.SYN1967Token, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the token from the ledger
	token, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving token for fractionalization: %v", err)
	}

	// Ensure the token is eligible for fractionalization
	if token.Fractionalized || len(fractions) != len(owners) || len(fractions) == 0 {
		return nil, errors.New("invalid token or fractionalization parameters")
	}

	// Fractionalize the token into multiple tokens
	var fractionalTokens []*common.SYN1967Token
	for i, fraction := range fractions {
		if fraction <= 0 || fraction > token.Amount {
			return nil, fmt.Errorf("invalid fraction value: %v", fraction)
		}

		// Create fractional tokens
		fractionalToken := &common.SYN1967Token{
			TokenID:          generateUniqueTokenID(),
			CommodityName:    token.CommodityName,
			Amount:           fraction,
			UnitOfMeasure:    token.UnitOfMeasure,
			Owner:            owners[i],
			Certification:    token.Certification,
			Traceability:     token.Traceability,
			IssuedDate:       token.IssuedDate,
			Origin:           token.Origin,
			ExpiryDate:       token.ExpiryDate,
			CollateralStatus: token.CollateralStatus,
			Fractionalized:   true,
		}

		// Encrypt each fractional token
		encryptedFractionalToken, err := encryption.Encrypt(fractionalToken)
		if err != nil {
			return nil, fmt.Errorf("error encrypting fractional token: %v", err)
		}

		// Store the fractional token in the ledger
		err = ledger.StoreToken(fractionalToken.TokenID, encryptedFractionalToken)
		if err != nil {
			return nil, fmt.Errorf("error storing fractional token in ledger: %v", err)
		}

		fractionalTokens = append(fractionalTokens, fractionalToken)
	}

	return fractionalTokens, nil
}

// UpdateCertification updates the certification status of a SYN1967 token.
func (tm *TokenManager) UpdateCertification(tokenID string, newCertification string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve the token from the ledger
	token, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return fmt.Errorf("error retrieving token for certification update: %v", err)
	}

	// Update the certification
	token.Certification = newCertification

	// Encrypt updated token
	encryptedToken, err := encryption.Encrypt(token)
	if err != nil {
		return fmt.Errorf("error encrypting updated token data: %v", err)
	}

	// Store the updated token in the ledger
	err = ledger.StoreToken(tokenID, encryptedToken)
	if err != nil {
		return fmt.Errorf("error updating token certification in ledger: %v", err)
	}

	return nil
}

// AuditToken generates an audit report of a SYN1967 token, showing the full history of ownership, transfers, and certifications.
func (tm *TokenManager) AuditToken(tokenID string) (*common.AuditReport, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Retrieve audit records from the ledger
	auditRecords, err := ledger.RetrieveAuditRecords(tokenID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving audit records for token: %v", err)
	}

	// Compile the audit report
	auditReport := &common.AuditReport{
		TokenID:      tokenID,
		Timestamp:    time.Now(),
		AuditRecords: auditRecords,
	}

	return auditReport, nil
}

// generateUniqueTokenID generates a unique token ID.
func generateUniqueTokenID() string {
	uniqueID, _ := rand.Int(rand.Reader, big.NewInt(1e12))
	return fmt.Sprintf("SYN1967-%d", uniqueID)
}
