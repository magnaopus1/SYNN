package syn1301

import (
	"errors"
	"time"

)

// SYN1301TokenManager manages the lifecycle of SYN1301 tokens, including creation, validation, and management.
type SYN1301TokenManager struct {
	Ledger            *ledger.Ledger                // Ledger system for managing tokens and assets
	EncryptionService *encryption.EncryptionService // Service to encrypt/decrypt data
	SecurityService   *security.SecurityService     // Service for managing security policies and access control
	Consensus         *synnergy_consensus.Consensus // Synnergy Consensus for validating transactions and sub-blocks
}

// CreateSYN1301Token creates and stores a new SYN1301 token after validating with Synnergy Consensus.
func (tm *SYN1301TokenManager) CreateSYN1301Token(token SYN1301Token) (SYN1301Token, error) {
	// Perform initial security checks
	err := tm.SecurityService.PerformSecurityChecks(token)
	if err != nil {
		return SYN1301Token{}, errors.New("security check failed: " + err.Error())
	}

	// Encrypt token metadata before storing
	metadata := map[string]string{
		"asset_id":    token.AssetID,
		"description": token.Description,
		"location":    token.Location,
		"status":      token.Status,
		"owner":       token.Owner,
	}
	encryptedMetadata, err := tm.EncryptionService.Encrypt(metadata)
	if err != nil {
		return SYN1301Token{}, errors.New("encryption failed: " + err.Error())
	}
	token.EncryptedMetadata = encryptedMetadata

	// Generate token ID
	token.TokenID = common.GenerateUniqueID()

	// Validate the token creation transaction via Synnergy Consensus
	subBlock := tm.Consensus.ValidateTransactionIntoSubBlock(token.TokenID, token)
	block, err := tm.Consensus.ValidateSubBlockIntoBlock(subBlock)
	if err != nil {
		return SYN1301Token{}, errors.New("consensus validation failed: " + err.Error())
	}

	// Store token in the ledger
	err = tm.Ledger.StoreToken(token.TokenID, token)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to store token in ledger: " + err.Error())
	}

	// Record the block ID after successful transaction validation
	token.BlockID = block.BlockID

	return token, nil
}

// GetSYN1301Token retrieves a SYN1301 token from the ledger and decrypts its metadata.
func (tm *SYN1301TokenManager) GetSYN1301Token(tokenID string) (SYN1301Token, error) {
	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, errors.New("token retrieval failed: " + err.Error())
	}

	// Decrypt the metadata
	decryptedMetadata, err := tm.EncryptionService.Decrypt(token.EncryptedMetadata)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to decrypt token metadata: " + err.Error())
	}

	// Populate decrypted fields into the token
	token.AssetID = decryptedMetadata["asset_id"]
	token.Description = decryptedMetadata["description"]
	token.Location = decryptedMetadata["location"]
	token.Status = decryptedMetadata["status"]
	token.Owner = decryptedMetadata["owner"]

	return token, nil
}

// UpdateSYN1301Token updates the metadata of an existing SYN1301 token and stores it back in the ledger.
func (tm *SYN1301TokenManager) UpdateSYN1301Token(tokenID string, updatedMetadata map[string]string) (SYN1301Token, error) {
	// Retrieve the token
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, errors.New("token not found: " + err.Error())
	}

	// Update the token fields with new metadata
	if description, ok := updatedMetadata["description"]; ok {
		token.Description = description
	}
	if location, ok := updatedMetadata["location"]; ok {
		token.Location = location
	}
	if status, ok := updatedMetadata["status"]; ok {
		token.Status = status
	}
	if owner, ok := updatedMetadata["owner"]; ok {
		token.Owner = owner
	}

	// Encrypt the updated metadata
	newEncryptedMetadata, err := tm.EncryptionService.Encrypt(updatedMetadata)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to encrypt updated metadata: " + err.Error())
	}
	token.EncryptedMetadata = newEncryptedMetadata

	// Validate the update through Synnergy Consensus
	subBlock := tm.Consensus.ValidateTransactionIntoSubBlock(tokenID, token)
	block, err := tm.Consensus.ValidateSubBlockIntoBlock(subBlock)
	if err != nil {
		return SYN1301Token{}, errors.New("consensus validation for update failed: " + err.Error())
	}

	// Update the token in the ledger
	err = tm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to update token in ledger: " + err.Error())
	}

	// Record the block ID for the update
	token.BlockID = block.BlockID

	return token, nil
}

// TransferSYN1301Token transfers ownership of a SYN1301 token from one owner to another.
func (tm *SYN1301TokenManager) TransferSYN1301Token(tokenID, newOwner string) (SYN1301Token, error) {
	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return SYN1301Token{}, errors.New("token not found: " + err.Error())
	}

	// Perform security checks for the transfer
	err = tm.SecurityService.ValidateOwnershipTransfer(token.Owner, newOwner)
	if err != nil {
		return SYN1301Token{}, errors.New("ownership transfer validation failed: " + err.Error())
	}

	// Update the token ownership
	token.Owner = newOwner

	// Encrypt the updated metadata
	updatedMetadata := map[string]string{
		"owner": newOwner,
	}
	encryptedMetadata, err := tm.EncryptionService.Encrypt(updatedMetadata)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to encrypt updated metadata: " + err.Error())
	}
	token.EncryptedMetadata = encryptedMetadata

	// Validate the transfer through Synnergy Consensus
	subBlock := tm.Consensus.ValidateTransactionIntoSubBlock(tokenID, token)
	block, err := tm.Consensus.ValidateSubBlockIntoBlock(subBlock)
	if err != nil {
		return SYN1301Token{}, errors.New("consensus validation for transfer failed: " + err.Error())
	}

	// Update the token in the ledger
	err = tm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return SYN1301Token{}, errors.New("failed to update token in ledger: " + err.Error())
	}

	// Record the block ID for the transfer
	token.BlockID = block.BlockID

	return token, nil
}

// DeleteSYN1301Token removes a token from the ledger after validation.
func (tm *SYN1301TokenManager) DeleteSYN1301Token(tokenID string) error {
	// Retrieve the token to ensure it exists
	_, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found: " + err.Error())
	}

	// Validate deletion rights through security service
	err = tm.SecurityService.ValidateDeletionRights(tokenID)
	if err != nil {
		return errors.New("deletion rights validation failed: " + err.Error())
	}

	// Validate the deletion through Synnergy Consensus
	subBlock := tm.Consensus.ValidateTransactionIntoSubBlock(tokenID, nil)
	block, err := tm.Consensus.ValidateSubBlockIntoBlock(subBlock)
	if err != nil {
		return errors.New("consensus validation for deletion failed: " + err.Error())
	}

	// Delete the token from the ledger
	err = tm.Ledger.DeleteToken(tokenID)
	if err != nil {
		return errors.New("failed to delete token from ledger: " + err.Error())
	}

	// Record the block ID for the deletion
	return nil
}
