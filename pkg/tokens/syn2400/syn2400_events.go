package syn2400

import (
	"errors"
	"time"

)

// SYN2400EventHandler handles all events related to SYN2400 tokens, such as creation, updates, and transfers.
type SYN2400EventHandler struct {
	Ledger   ledger.LedgerInterface        // Interface for interacting with the blockchain ledger
	Encrypt  encryption.EncryptionInterface // Interface for encryption
}

// NewSYN2400EventHandler initializes a new instance of SYN2400EventHandler
func NewSYN2400EventHandler(ledger ledger.LedgerInterface, encrypt encryption.EncryptionInterface) *SYN2400EventHandler {
	return &SYN2400EventHandler{
		Ledger:  ledger,
		Encrypt: encrypt,
	}
}

// CreateDataToken handles the creation of a new SYN2400 token representing a data set
func (handler *SYN2400EventHandler) CreateDataToken(
	owner string,
	dataHash string,
	description string,
	accessRights string,
	price float64) (common.SYN2400Token, error) {

	// Generate a new token ID
	tokenID := generateUniqueID()

	// Create a new SYN2400Token struct
	newToken := common.SYN2400Token{
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
	newToken.AuditTrail = append(newToken.AuditTrail, common.AuditRecord{
		Action:      "Data Token Created",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Token created with description: " + description,
	})

	// Encrypt the token data before saving it to the ledger
	encryptedToken, err := handler.Encrypt.EncryptTokenData(newToken)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the encrypted token in the ledger
	if err := handler.Ledger.CreateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return newToken, nil
}

// UpdateDataToken handles updates to an existing SYN2400 token's metadata or price
func (handler *SYN2400EventHandler) UpdateDataToken(
	tokenID string,
	updatedDescription string,
	updatedAccessRights string,
	updatedPrice float64,
	owner string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := handler.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Verify ownership
	if token.Owner != owner {
		return common.SYN2400Token{}, errors.New("only the token owner can update the data token")
	}

	// Update the token details
	token.Description = updatedDescription
	token.AccessRights = updatedAccessRights
	token.Price = updatedPrice
	token.UpdateDate = time.Now()

	// Log the update event
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:      "Data Token Updated",
		PerformedBy: owner,
		Timestamp:   time.Now(),
		Details:     "Updated description and price for the token",
	})

	// Encrypt the updated token before storing it back in the ledger
	encryptedToken, err := handler.Encrypt.EncryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := handler.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// TransferDataToken handles the secure transfer of ownership of a SYN2400 token
func (handler *SYN2400EventHandler) TransferDataToken(
	tokenID string,
	newOwner string,
	currentOwner string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := handler.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Verify current ownership
	if token.Owner != currentOwner {
		return common.SYN2400Token{}, errors.New("only the current owner can transfer the data token")
	}

	// Update ownership
	token.Owner = newOwner
	token.UpdateDate = time.Now()

	// Log the transfer event
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:      "Ownership Transferred",
		PerformedBy: currentOwner,
		Timestamp:   time.Now(),
		Details:     "Token transferred to new owner: " + newOwner,
	})

	// Encrypt the updated token before storing it back in the ledger
	encryptedToken, err := handler.Encrypt.EncryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := handler.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// RecordAccessEvent logs an access event whenever a user accesses data from a SYN2400 token
func (handler *SYN2400EventHandler) RecordAccessEvent(
	tokenID string,
	accessedBy string,
	accessDetails string) (common.SYN2400Token, error) {

	// Retrieve the token from the ledger
	token, err := handler.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Log the access event
	token.AuditTrail = append(token.AuditTrail, common.AuditRecord{
		Action:      "Data Accessed",
		PerformedBy: accessedBy,
		Timestamp:   time.Now(),
		Details:     accessDetails,
	})

	// Encrypt the updated token before storing it back in the ledger
	encryptedToken, err := handler.Encrypt.EncryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, err
	}

	// Store the updated token in the ledger
	if err := handler.Ledger.UpdateToken(tokenID, encryptedToken); err != nil {
		return common.SYN2400Token{}, err
	}

	return token, nil
}

// generateUniqueID is a helper function for generating unique token IDs and event logs
func generateUniqueID() string {
	return "SYN2400-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random alphanumeric string for creating unique identifiers
func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
