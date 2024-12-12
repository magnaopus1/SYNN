package syn1301

import (
	"errors"
	"time"

)

// SYN1301Token defines the structure for supply chain tokens under the SYN1301 standard.
type SYN1301Token struct {
	TokenID           string    `json:"token_id"`            // Unique token identifier for the supply chain asset
	AssetID           string    `json:"asset_id"`            // Unique asset identifier associated with this token
	Description       string    `json:"description"`         // Description of the asset being tracked
	Location          string    `json:"location"`            // Current physical location of the asset
	Timestamp         time.Time `json:"timestamp"`           // Timestamp of the last status update
	Status            string    `json:"status"`              // Status of the asset (e.g., in-transit, delivered, etc.)
	Owner             string    `json:"owner"`               // Current owner of the asset
	EncryptedMetadata string    `json:"encrypted_metadata"`  // Encrypted metadata of the asset for secure storage
	AssetValue        float64   `json:"asset_value"`         // Original purchase price of the asset
	LastSalePrice     float64   `json:"last_sale_price"`     // Price at which the asset was last sold
	Quantity          int       `json:"quantity"`            // Quantity of the asset being tracked in the supply chain
	BatchID           string    `json:"batch_id"`            // Batch or lot identifier for the asset
	ExpiryDate        time.Time `json:"expiry_date"`         // Expiry date for perishable goods (if applicable)
}

// SYN1301Factory manages the creation and operation of SYN1301Tokens.
type SYN1301Factory struct {
	Ledger            *ledger.Ledger                // Ledger system for token storage and management
	EncryptionService *encryption.EncryptionService // Encryption service to securely store metadata
}

// NewSYN1301Token creates a new supply chain token, encrypts its metadata, and stores it in the ledger.
func (f *SYN1301Factory) NewSYN1301Token(assetID, description, location, owner string, status string, assetValue, lastSalePrice float64, quantity int, batchID string, expiryDate time.Time) (SYN1301Token, error) {
	// Validate inputs
	if assetID == "" || description == "" || location == "" || owner == "" || status == "" || quantity <= 0 || assetValue < 0 || lastSalePrice < 0 {
		return SYN1301Token{}, errors.New("missing or invalid asset information")
	}

	// Generate a unique token ID
	tokenID := common.GenerateUUID()

	// Create the token with the current timestamp
	token := SYN1301Token{
		TokenID:       tokenID,
		AssetID:       assetID,
		Description:   description,
		Location:      location,
		Timestamp:     time.Now(),
		Status:        status,
		Owner:         owner,
		AssetValue:    assetValue,
		LastSalePrice: lastSalePrice,
		Quantity:      quantity,
		BatchID:       batchID,
		ExpiryDate:    expiryDate,
	}

	// Encrypt asset metadata
	metadata := map[string]interface{}{
		"description":    description,
		"location":       location,
		"status":         status,
		"asset_value":    assetValue,
		"last_sale_price": lastSalePrice,
		"quantity":       quantity,
		"batch_id":       batchID,
		"expiry_date":    expiryDate,
	}
	encryptedMetadata, err := f.EncryptionService.Encrypt(metadata)
	if err != nil {
		return SYN1301Token{}, err
	}

	// Assign encrypted metadata
	token.EncryptedMetadata = encryptedMetadata

	// Store the token in the ledger
	err = f.Ledger.StoreToken(token)
	if err != nil {
		return SYN1301Token{}, err
	}

	// Return the newly created token
	return token, nil
}

// UpdateTokenStatus updates the status, location, and potentially last sale price of an existing SYN1301Token and records it in the ledger.
func (f *SYN1301Factory) UpdateTokenStatus(tokenID, newStatus, newLocation string, newLastSalePrice float64) (SYN1301Token, error) {
	// Fetch the token from the ledger
	token, err := f.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, err
	}

	// Update status, location, and last sale price if applicable
	token.Status = newStatus
	token.Location = newLocation
	token.Timestamp = time.Now()

	if newLastSalePrice > 0 {
		token.LastSalePrice = newLastSalePrice
	}

	// Encrypt updated metadata
	updatedMetadata := map[string]interface{}{
		"description":    token.Description,
		"location":       newLocation,
		"status":         newStatus,
		"asset_value":    token.AssetValue,
		"last_sale_price": token.LastSalePrice,
		"quantity":       token.Quantity,
		"batch_id":       token.BatchID,
		"expiry_date":    token.ExpiryDate,
	}
	encryptedMetadata, err := f.EncryptionService.Encrypt(updatedMetadata)
	if err != nil {
		return SYN1301Token{}, err
	}
	token.EncryptedMetadata = encryptedMetadata

	// Store the updated token back into the ledger
	err = f.Ledger.UpdateToken(token)
	if err != nil {
		return SYN1301Token{}, err
	}

	// Return the updated token
	return token, nil
}

// GetTokenByID retrieves the details of a SYN1301Token by its TokenID from the ledger.
func (f *SYN1301Factory) GetTokenByID(tokenID string) (SYN1301Token, error) {
	// Retrieve the token from the ledger
	token, err := f.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, err
	}

	// Decrypt metadata for viewing
	metadata, err := f.EncryptionService.Decrypt(token.EncryptedMetadata)
	if err != nil {
		return SYN1301Token{}, err
	}

	// Populate decrypted metadata into the token object
	// Note: It's assumed here that metadata includes the same fields as the struct
	token.Description = metadata["description"].(string)
	token.Location = metadata["location"].(string)
	token.Status = metadata["status"].(string)
	token.AssetValue = metadata["asset_value"].(float64)
	token.LastSalePrice = metadata["last_sale_price"].(float64)
	token.Quantity = metadata["quantity"].(int)
	token.BatchID = metadata["batch_id"].(string)
	token.ExpiryDate = metadata["expiry_date"].(time.Time)

	return token, nil
}
