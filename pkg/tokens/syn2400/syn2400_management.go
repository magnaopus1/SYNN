package syn2400

import (
	"errors"
	"time"

)

// SYN2400Management provides functions for managing SYN2400 tokens on the blockchain
type SYN2400Management struct {
	Ledger   ledger.LedgerInterface        // Interface for interacting with the blockchain ledger
	Encrypt  encryption.EncryptionInterface // Interface for encryption
}

// NewSYN2400Management initializes a new instance of SYN2400Management
func NewSYN2400Management(ledger ledger.LedgerInterface, encrypt encryption.EncryptionInterface) *SYN2400Management {
	return &SYN2400Management{
		Ledger:  ledger,
		Encrypt: encrypt,
	}
}

// CreateNewDataToken handles the creation and management of a new SYN2400 token representing data
func (manager *SYN2400Management) CreateNewDataToken(
	owner string,
	dataHash string,
	description string,
	accessRights string,
	price float64) (common.SYN2400Token, error) {

	// Generate unique token ID
	tokenID := generateUniqueID()

	// Create the SYN2400Token structure
	dataToken := common.SYN2400Token{
		TokenID:      tokenID,
		Owner:        owner,
		DataHash:     dataHash,
		Description:  description,
		AccessRights: accessRights,
		CreationDate: time.Now(),
		UpdateDate:   time.Now(),
		Price:        price,
		Status:       "Available",
	}

	// Log the creation event
	dataToken.AuditTrail = append(dataToken.AuditTrail, common.AuditRecord{
		Action:      "Token Created",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Data token created for " + description,
	})

	// Encrypt the token data before saving it
	encryptedToken, err := manager.Encrypt.EncryptTokenData(dataToken)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the encrypted token in the ledger
	if err := manager.Ledger.CreateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return dataToken, nil
}

// UpdateDataTokenMetadata handles updates to metadata and price for an existing SYN2400 token
func (manager *SYN2400Management) UpdateDataTokenMetadata(
	tokenID string,
	updatedDescription string,
	updatedAccessRights string,
	updatedPrice float64,
	owner string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := manager.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Verify ownership
	if token.Owner != owner {
		return common.SYN2400Token{}, errors.New("only the owner can update this token")
	}

	// Update the token's metadata
	token.Description = updatedDescription
	token.AccessRights = updatedAccessRights
	token.Price = updatedPrice
	token.UpdateDate = time.Now()

	// Log the update event
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:      "Metadata Updated",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Updated description and price",
	})

	// Encrypt the updated token before saving it
	encryptedToken, err := manager.Encrypt.EncryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := manager.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// TransferOwnership transfers the ownership of a SYN2400 token to a new user
func (manager *SYN2400Management) TransferOwnership(
	tokenID string,
	newOwner string,
	currentOwner string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := manager.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Verify current ownership
	if token.Owner != currentOwner {
		return common.SYN2400Token{}, errors.New("only the current owner can transfer ownership")
	}

	// Transfer ownership
	token.Owner = newOwner
	token.UpdateDate = time.Now()

	// Log the transfer event
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:      "Ownership Transferred",
		PerformedBy: currentOwner,
		Timestamp:   time.Now(),
		Details:     "Ownership transferred to " + newOwner,
	})

	// Encrypt the updated token before saving it
	encryptedToken, err := manager.Encrypt.EncryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := manager.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// RetireDataToken marks a SYN2400 token as retired and unavailable for future transactions
func (manager *SYN2400Management) RetireDataToken(tokenID string, owner string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := manager.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Verify ownership
	if token.Owner != owner {
		return common.SYN2400Token{}, errors.New("only the owner can retire this token")
	}

	// Mark the token as retired
	token.Status = "Retired"
	token.UpdateDate = time.Now()

	// Log the retirement event
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:      "Token Retired",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Data token retired and removed from active circulation",
	})

	// Encrypt the updated token before saving it
	encryptedToken, err := manager.Encrypt.EncryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := manager.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// QueryToken retrieves a SYN2400 token by its ID and decrypts it for viewing
func (manager *SYN2400Management) QueryToken(tokenID string) (common.SYN2400Token, error) {

	// Retrieve the encrypted token from the ledger
	encryptedToken, err := manager.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Decrypt the token data
	decryptedToken, err := manager.Encrypt.DecryptTokenData(encryptedToken)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	return decryptedToken, nil
}

// generateUniqueID generates a unique token ID using current time and random string
func generateUniqueID() string {
	return "SYN2400-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random alphanumeric string for unique identifiers
func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
