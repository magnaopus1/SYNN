package syn2600

import (

	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

)

// SYN2600Token represents a financial asset in the form of an investor token.
type SYN2600Token struct {
	TokenID         string            // Unique identifier for the token
	AssetDetails    string            // Details of the underlying financial asset (loan, bond, security)
	Owner           string            // Current owner of the token
	Shares          float64           // Number of shares or fractional ownership represented by the token
	IssuedDate      time.Time         // Date when the token was issued
	ExpiryDate      time.Time         // Date when the token expires
	ActiveStatus    bool              // Indicates if the token is active
	TransactionLogs []TransactionLog  // Log of all transactions for this token
	ReturnLogs      []ReturnLog       // Log of returns (dividends, interest) for this token
	ComplianceHash  string            // Hash for ensuring compliance integrity
}

// TransactionLog keeps track of transactions related to the investor token.
type TransactionLog struct {
	Timestamp        time.Time // Time of the transaction
	FromAddress      string    // Address of the sender
	ToAddress        string    // Address of the recipient
	SharesTransferred float64  // Amount of shares transferred
	TransactionHash  string    // Hash of the transaction for verification
}

// ReturnLog keeps track of returns (e.g., dividends) distributed to token holders.
type ReturnLog struct {
	DistributionDate time.Time // Time when returns were distributed
	ReturnAmount     float64   // The amount distributed to the token holder
	DistributionHash string    // Hash of the return transaction for verification
}

// CreateSYN2600Token creates a new investor token representing a financial asset.
func CreateSYN2600Token(assetDetails, owner string, shares float64, issuedDate, expiryDate time.Time) (*SYN2600Token, error) {
	if shares <= 0 {
		return nil, errors.New("invalid shares: must be greater than zero")
	}

	// Create a new token
	token := &SYN2600Token{
		TokenID:        generateUniqueTokenID(assetDetails, owner),
		AssetDetails:   assetDetails,
		Owner:          owner,
		Shares:         shares,
		IssuedDate:     issuedDate,
		ExpiryDate:     expiryDate,
		ActiveStatus:   true,
		TransactionLogs: []TransactionLog{},
		ReturnLogs:     []ReturnLog{},
	}

	// Generate compliance hash for the token
	token.ComplianceHash = generateComplianceHash(token)

	// Store the token in the ledger
	err := ledger.StoreInvestorToken(token)
	if err != nil {
		return nil, errors.New("failed to store token in the ledger")
	}

	return token, nil
}

// TransferOwnership transfers ownership of investor token shares between parties.
func (t *SYN2600Token) TransferOwnership(to string, shares float64) error {
	if !t.ActiveStatus {
		return errors.New("token is inactive")
	}
	if shares <= 0 || shares > t.Shares {
		return errors.New("invalid share amount for transfer")
	}

	// Create transaction log
	txLog := TransactionLog{
		Timestamp:        time.Now(),
		FromAddress:      t.Owner,
		ToAddress:        to,
		SharesTransferred: shares,
		TransactionHash:  generateTransactionHash(t.TokenID, t.Owner, to, shares),
	}
	t.TransactionLogs = append(t.TransactionLogs, txLog)

	// Update ownership and shares
	t.Shares -= shares
	if t.Shares == 0 {
		t.Owner = to // Transfer complete ownership
	} else {
		// Otherwise, partial transfer occurs, and ownership is shared.
	}

	// Store updated token in the ledger
	err := ledger.UpdateInvestorToken(t)
	if err != nil {
		return errors.New("failed to update token ownership in the ledger")
	}

	return nil
}

// DistributeReturns distributes returns (dividends, interest) to the token holder.
func (t *SYN2600Token) DistributeReturns(returnAmount float64) error {
	if !t.ActiveStatus {
		return errors.New("token is inactive")
	}

	// Create return log
	returnLog := ReturnLog{
		DistributionDate: time.Now(),
		ReturnAmount:     returnAmount,
		DistributionHash: generateReturnHash(t.TokenID, returnAmount),
	}
	t.ReturnLogs = append(t.ReturnLogs, returnLog)

	// Automate return distribution using smart contracts (placeholder for logic)
	err := ledger.UpdateInvestorToken(t)
	if err != nil {
		return errors.New("failed to update token returns in the ledger")
	}

	return nil
}

// CalculateReturn calculates potential returns based on historical performance.
func (t *SYN2600Token) CalculateReturn() (float64, error) {
	// Aggregate all return amounts in the logs
	totalReturns := 0.0
	for _, retLog := range t.ReturnLogs {
		totalReturns += retLog.ReturnAmount
	}

	// Return the total calculated returns
	return totalReturns, nil
}

// Generate compliance hash to ensure token's integrity and compliance tracking
func generateComplianceHash(t *SYN2600Token) string {
	hashInput := t.TokenID + t.AssetDetails + t.Owner + string(t.Shares) + t.IssuedDate.String() + t.ExpiryDate.String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateTransactionHash generates a hash for each transaction to ensure integrity.
func generateTransactionHash(tokenID, from, to string, shares float64) string {
	hashInput := tokenID + from + to + string(shares) + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateReturnHash generates a hash for return (dividend/interest) distribution.
func generateReturnHash(tokenID string, returnAmount float64) string {
	hashInput := tokenID + string(returnAmount) + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateUniqueTokenID generates a unique token ID based on asset and owner details.
func generateUniqueTokenID(assetDetails, owner string) string {
	hashInput := assetDetails + owner + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// ValidateTransaction validates the investor token's transaction through Synnergy Consensus.
func (t *SYN2600Token) ValidateTransaction(txLog TransactionLog) error {
	// Verify transaction hash
	validHash := generateTransactionHash(t.TokenID, txLog.FromAddress, txLog.ToAddress, txLog.SharesTransferred)
	if txLog.TransactionHash != validHash {
		return errors.New("invalid transaction hash")
	}

	// Perform consensus validation via Synnergy Consensus
	err := synconsensus.ValidateSubBlock(txLog.TransactionHash)
	if err != nil {
		return errors.New("transaction validation failed via Synnergy Consensus")
	}

	return nil
}

// CreateSYN2600Token creates a new investor token representing a financial asset in the system.
func CreateSYN2600Token(assetDetails, owner string, shares float64, issuedDate, expiryDate time.Time) (*SYN2600Token, error) {
	if shares <= 0 {
		return nil, errors.New("invalid shares: must be greater than zero")
	}

	// Create the new investor token with essential data
	token := &SYN2600Token{
		TokenID:        generateUniqueTokenID(assetDetails, owner),
		AssetDetails:   assetDetails,
		Owner:          owner,
		Shares:         shares,
		IssuedDate:     issuedDate,
		ExpiryDate:     expiryDate,
		ActiveStatus:   true,
		TransactionLogs: []TransactionLog{},
		ReturnLogs:     []ReturnLog{},
	}

	// Generate compliance hash for integrity and audit purposes
	token.ComplianceHash = generateComplianceHash(token)

	// Encrypt token data before storing
	encryptedToken, err := encryption.EncryptTokenData(token)
	if err != nil {
		return nil, errors.New("encryption failed for token data")
	}

	// Store the encrypted token in the ledger
	err = ledger.StoreInvestorToken(encryptedToken)
	if err != nil {
		return nil, errors.New("failed to store token in the ledger")
	}

	// Record the creation event on Synnergy Consensus by validating it into a sub-block
	err = synconsensus.ValidateSubBlock(token.TokenID)
	if err != nil {
		return nil, errors.New("validation of token creation failed in Synnergy Consensus")
	}

	return token, nil
}

// TransferSYN2600Token handles the transfer of investor tokens between two parties.
func TransferSYN2600Token(tokenID, from, to string, shares float64) error {
	// Fetch the token from the ledger
	token, err := ledger.FetchInvestorToken(tokenID)
	if err != nil {
		return errors.New("failed to fetch the token from the ledger")
	}

	// Decrypt the token data before processing
	decryptedToken, err := encryption.DecryptTokenData(token)
	if err != nil {
		return errors.New("failed to decrypt token data")
	}

	// Ensure the sender owns the token and has enough shares to transfer
	if decryptedToken.Owner != from {
		return errors.New("ownership mismatch: the sender is not the current owner")
	}
	if shares > decryptedToken.Shares {
		return errors.New("insufficient shares for transfer")
	}

	// Create a transaction log for the transfer
	txLog := TransactionLog{
		Timestamp:        time.Now(),
		FromAddress:      from,
		ToAddress:        to,
		SharesTransferred: shares,
		TransactionHash:  generateTransactionHash(tokenID, from, to, shares),
	}
	decryptedToken.TransactionLogs = append(decryptedToken.TransactionLogs, txLog)

	// Update ownership and shares
	decryptedToken.Shares -= shares
	if decryptedToken.Shares == 0 {
		decryptedToken.Owner = to
	}

	// Encrypt the token after updating
	encryptedToken, err := encryption.EncryptTokenData(decryptedToken)
	if err != nil {
		return errors.New("encryption failed after token transfer")
	}

	// Store the updated token back in the ledger
	err = ledger.UpdateInvestorToken(encryptedToken)
	if err != nil {
		return errors.New("failed to update token in the ledger")
	}

	// Validate the transaction through Synnergy Consensus
	err = synconsensus.ValidateSubBlock(txLog.TransactionHash)
	if err != nil {
		return errors.New("transaction validation failed in Synnergy Consensus")
	}

	return nil
}

// generateUniqueTokenID generates a unique token ID for the investor token.
func generateUniqueTokenID(assetDetails, owner string) string {
	hashInput := assetDetails + owner + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateTransactionHash generates a unique hash for each token transaction.
func generateTransactionHash(tokenID, from, to string, shares float64) string {
	hashInput := tokenID + from + to + string(shares) + time.Now().String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateComplianceHash generates a hash to ensure integrity for investor token compliance.
func generateComplianceHash(token *SYN2600Token) string {
	hashInput := token.TokenID + token.AssetDetails + token.Owner + string(token.Shares) + token.IssuedDate.String() + token.ExpiryDate.String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}
