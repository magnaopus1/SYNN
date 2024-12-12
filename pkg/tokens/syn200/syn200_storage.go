package syn200

import (
	"errors"
	"sync"
	"time"
)

// StoreToken securely stores a SYN200 token, encrypting sensitive data and logging the storage event.
func StoreToken(token *common.SYN200Token) error {
    encryptedMetadata, err := encryption.EncryptMetadata(token.CreditMetadata)
    if err != nil {
        return fmt.Errorf("failed to encrypt token metadata: %v", err)
    }
    token.EncryptedMetadata = encryptedMetadata

    // Record storage event in the ledger
    storageEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        token.TokenID,
        EventType:      "Storage",
        EventTimestamp: time.Now(),
        Description:    "Token securely stored with encrypted metadata",
    }

    encryptedEvent, err := encryption.EncryptMetadata(storageEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt storage event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log storage event in ledger: %v", err)
    }

    // Add token to ledger storage
    if err := ledger.StoreToken(token); err != nil {
        return fmt.Errorf("failed to store token in ledger: %v", err)
    }

    return nil
}

// RetrieveToken retrieves a SYN200 token from storage, decrypting sensitive data for secure access.
func RetrieveToken(tokenID string) (*common.SYN200Token, error) {
    encryptedToken, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return nil, fmt.Errorf("token not found: %v", err)
    }

    // Decrypt metadata for secure access
    decryptedMetadata, err := encryption.DecryptMetadata(encryptedToken.EncryptedMetadata)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt token metadata: %v", err)
    }
    encryptedToken.CreditMetadata = decryptedMetadata.(common.CarbonCreditMetadata)

    // Log retrieval in the ledger
    retrievalEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Retrieval",
        EventTimestamp: time.Now(),
        Description:    "Token retrieved and decrypted",
    }

    encryptedEvent, err := encryption.EncryptMetadata(retrievalEvent)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt retrieval event: %v", err)
    }

    if err := ledger.RecordEvent(encryptedToken, encryptedEvent); err != nil {
        return nil, fmt.Errorf("failed to log retrieval event in ledger: %v", err)
    }

    return encryptedToken, nil
}

// ArchiveToken archives a SYN200 token, marking it as archived in the ledger for long-term storage.
func ArchiveToken(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    token.ValidityStatus = "Archived"

    archiveEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Archive",
        EventTimestamp: time.Now(),
        Description:    "Token archived for long-term storage",
    }

    // Encrypt and record archive event in ledger
    encryptedEvent, err := encryption.EncryptMetadata(archiveEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt archive event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log archive event in ledger: %v", err)
    }

    return nil
}

// DeleteToken removes a SYN200 token from storage, logging the deletion and encrypting the record for secure tracking.
func DeleteToken(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    deletionEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        tokenID,
        EventType:      "Deletion",
        EventTimestamp: time.Now(),
        Description:    fmt.Sprintf("Token %s permanently deleted", tokenID),
    }

    // Encrypt and record the deletion event in ledger
    encryptedEvent, err := encryption.EncryptMetadata(deletionEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt deletion event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log deletion event in ledger: %v", err)
    }

    // Remove token from ledger storage
    if err := ledger.DeleteToken(tokenID); err != nil {
        return fmt.Errorf("failed to delete token from ledger: %v", err)
    }

    return nil
}

// generateEventID creates a unique identifier for events.
func generateEventID() string {
    b := make([]byte, 12)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
