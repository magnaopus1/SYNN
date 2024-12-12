package syn1600

import (
	"errors"
	"time"
)

// StorageManager handles the storage, retrieval, and validation of SYN1600 tokens.
type StorageManager struct {
	Ledger ledger.Ledger  // The ledger to interact with for storing tokens
}

// StoreSYN1600Token stores a new SYN1600Token on the blockchain ledger.
func (sm *StorageManager) StoreSYN1600Token(token *common.SYN1600Token, encryptionKey []byte) error {
	// Step 1: Encrypt sensitive data before storing the token
	err := sm.encryptTokenData(token, encryptionKey)
	if err != nil {
		return err
	}

	// Step 2: Validate the token using Synnergy Consensus
	err = synnergy.ValidateToken(token.TokenID)
	if err != nil {
		return err
	}

	// Step 3: Store the token in the ledger
	err = sm.Ledger.StoreToken(token.TokenID, token)
	if err != nil {
		return err
	}

	// Step 4: Log the event of storing the token
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "TokenStored",
		Description: "SYN1600 token stored successfully.",
		EventDate:   time.Now(),
		PerformedBy: "StorageSystem",
	}
	token.AuditTrail = append(token.AuditTrail, eventLog)

	// Step 5: Update the ledger with the audit trail
	return sm.Ledger.UpdateToken(token.TokenID, token)
}

// RetrieveSYN1600Token retrieves a SYN1600Token from the blockchain ledger.
func (sm *StorageManager) RetrieveSYN1600Token(tokenID string, decryptionKey []byte) (*common.SYN1600Token, error) {
	// Step 1: Retrieve the token from the ledger
	token, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, err
	}

	// Step 2: Decrypt sensitive data within the token
	err = sm.decryptTokenData(token.(*common.SYN1600Token), decryptionKey)
	if err != nil {
		return nil, err
	}

	// Step 3: Log the event of token retrieval
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "TokenRetrieved",
		Description: "SYN1600 token retrieved successfully.",
		EventDate:   time.Now(),
		PerformedBy: "StorageSystem",
	}
	token.(*common.SYN1600Token).AuditTrail = append(token.(*common.SYN1600Token).AuditTrail, eventLog)

	// Step 4: Update the ledger with the audit trail
	err = sm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return nil, err
	}

	return token.(*common.SYN1600Token), nil
}

// DeleteSYN1600Token deletes a SYN1600Token from the blockchain ledger.
func (sm *StorageManager) DeleteSYN1600Token(tokenID string) error {
	// Step 1: Validate that the token exists
	_, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("token not found")
	}

	// Step 2: Remove the token from the ledger
	err = sm.Ledger.DeleteToken(tokenID)
	if err != nil {
		return err
	}

	// Step 3: Log the event of token deletion
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "TokenDeleted",
		Description: "SYN1600 token deleted successfully.",
		EventDate:   time.Now(),
		PerformedBy: "StorageSystem",
	}

	// Step 4: Since the token is deleted, add the event log to the ledger separately
	return sm.Ledger.StoreEventLog(tokenID, eventLog)
}

// encryptTokenData encrypts the sensitive fields of the SYN1600Token before storage.
func (sm *StorageManager) encryptTokenData(token *common.SYN1600Token, encryptionKey []byte) error {
	// Encrypt ownership records
	for i, ownership := range token.OwnershipRights {
		encryptedOwnerID, err := encryption.Encrypt([]byte(ownership.OwnerID), encryptionKey)
		if err != nil {
			return err
		}
		token.OwnershipRights[i].OwnerID = string(encryptedOwnerID)
	}

	// Encrypt royalty distribution logs
	for i, log := range token.RevenueDistribution {
		encryptedRecipientID, err := encryption.Encrypt([]byte(log.RecipientID), encryptionKey)
		if err != nil {
			return err
		}
		token.RevenueDistribution[i].RecipientID = string(encryptedRecipientID)
	}

	// Encrypt sensitive metadata
	token.EncryptedMetadata = encryption.EncryptBytes(token.EncryptedMetadata, encryptionKey)

	return nil
}

// decryptTokenData decrypts the sensitive fields of the SYN1600Token after retrieval.
func (sm *StorageManager) decryptTokenData(token *common.SYN1600Token, decryptionKey []byte) error {
	// Decrypt ownership records
	for i, ownership := range token.OwnershipRights {
		decryptedOwnerID, err := encryption.Decrypt([]byte(ownership.OwnerID), decryptionKey)
		if err != nil {
			return err
		}
		token.OwnershipRights[i].OwnerID = string(decryptedOwnerID)
	}

	// Decrypt royalty distribution logs
	for i, log := range token.RevenueDistribution {
		decryptedRecipientID, err := encryption.Decrypt([]byte(log.RecipientID), decryptionKey)
		if err != nil {
			return err
		}
		token.RevenueDistribution[i].RecipientID = string(decryptedRecipientID)
	}

	// Decrypt sensitive metadata
	token.EncryptedMetadata = encryption.DecryptBytes(token.EncryptedMetadata, decryptionKey)

	return nil
}

// generateUniqueID generates a unique ID for events or logs.
func generateUniqueID() string {
	return "ID_" + time.Now().Format("20060102150405")
}
