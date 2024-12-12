package syn1200

import (
	"errors"
	"time"
)




// StoreToken securely stores the interoperable token data in the ledger.
func (sm *SYN1200StorageManager) StoreToken(token InteroperableToken) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the token data before storing
	encryptedTokenData, err := sm.EncryptTokenData(token)
	if err != nil {
		return errors.New("failed to encrypt token data")
	}

	// Store the encrypted token in the ledger
	if err := sm.Ledger.StoreToken(token.TokenID, encryptedTokenData); err != nil {
		return errors.New("failed to store token in ledger")
	}

	return nil
}

// EncryptTokenData encrypts the interoperable token's data before storing it.
func (sm *SYN1200StorageManager) EncryptTokenData(token InteroperableToken) (string, error) {
	// Serialize token data
	tokenData := common.StructToString(token)

	// Generate encryption key for token
	encryptionKey := sm.EncryptionService.GenerateKey()
	encryptedData, err := sm.EncryptionService.EncryptData([]byte(tokenData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to encrypt token data")
	}

	// Store encryption key in the ledger for future decryption
	if err := sm.Ledger.StoreEncryptionKey(token.TokenID, encryptionKey); err != nil {
		return "", errors.New("failed to store encryption key")
	}

	return string(encryptedData), nil
}

// RetrieveToken retrieves and decrypts token data from the ledger.
func (sm *SYN1200StorageManager) RetrieveToken(tokenID string) (*InteroperableToken, error) {
	// Retrieve the encrypted token data from the ledger
	encryptedData, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve token from ledger")
	}

	// Decrypt the token data
	token, err := sm.DecryptTokenData(tokenID, encryptedData)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// DecryptTokenData decrypts the token data retrieved from the ledger.
func (sm *SYN1200StorageManager) DecryptTokenData(tokenID string, encryptedData string) (*InteroperableToken, error) {
	// Retrieve the encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key")
	}

	// Decrypt the token data
	decryptedData, err := sm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	// Deserialize the decrypted data back into the token struct
	var token InteroperableToken
	if err := common.StringToStruct(string(decryptedData), &token); err != nil {
		return nil, errors.New("failed to deserialize token data")
	}

	return &token, nil
}

// UpdateToken updates the token data in the ledger and encrypts the updates.
func (sm *SYN1200StorageManager) UpdateToken(token InteroperableToken) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Encrypt the updated token data
	encryptedTokenData, err := sm.EncryptTokenData(token)
	if err != nil {
		return errors.New("failed to encrypt updated token data")
	}

	// Update the encrypted token in the ledger
	if err := sm.Ledger.UpdateToken(token.TokenID, encryptedTokenData); err != nil {
		return errors.New("failed to update token in ledger")
	}

	return nil
}

// AddTransactionRecord adds a transaction record to the token's history and updates it in the ledger.
func (sm *SYN1200StorageManager) AddTransactionRecord(tokenID string, record TransactionRecord) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token
	token, err := sm.RetrieveToken(tokenID)
	if err != nil {
		return err
	}

	// Append the transaction record to the history
	token.TransactionHistory = append(token.TransactionHistory, record)

	// Encrypt and update the token in the ledger
	encryptedTokenData, err := sm.EncryptTokenData(*token)
	if err != nil {
		return err
	}

	if err := sm.Ledger.UpdateToken(token.TokenID, encryptedTokenData); err != nil {
		return errors.New("failed to update token with new transaction record")
	}

	return nil
}

// GetTransactionHistory retrieves the transaction history of a specific token.
func (sm *SYN1200StorageManager) GetTransactionHistory(tokenID string) ([]TransactionRecord, error) {
	// Retrieve the token from the ledger
	token, err := sm.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}

	// Return the transaction history
	return token.TransactionHistory, nil
}

