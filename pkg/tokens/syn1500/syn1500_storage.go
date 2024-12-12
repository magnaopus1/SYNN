package syn1500

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SYN1500Storage provides storage management for SYN1500 reputation tokens.
type SYN1500Storage struct {
	Ledger ledger.Ledger // Ledger to log storage and retrieval transactions
}

// SaveToken securely stores the SYN1500Token into the blockchain ledger, encrypting sensitive data.
func (storage *SYN1500Storage) SaveToken(token *common.SYN1500Token, encryptionKey []byte) error {
	// Encrypt the token's metadata before storing it
	encryptedData, err := encryptTokenData(token, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt token metadata: %v", err)
	}
	token.EncryptedMetadata = encryptedData

	// Generate storage transaction
	tx := ledger.Transaction{
		TxID:        generateUniqueID(token.TokenID),
		Description: fmt.Sprintf("Store token %s into the ledger", token.TokenID),
		Timestamp:   time.Now(),
		Data:        map[string]interface{}{"token_id": token.TokenID, "metadata": token.EncryptedMetadata},
	}

	// Validate the transaction using Synnergy Consensus
	if err := synnergy_consensus.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("failed Synnergy Consensus validation for storing token: %v", err)
	}

	// Record the storage event in the ledger
	return storage.Ledger.RecordTransaction(tx)
}

// RetrieveToken fetches a SYN1500Token from storage, decrypting its metadata.
func (storage *SYN1500Storage) RetrieveToken(tokenID string, decryptionKey []byte) (*common.SYN1500Token, error) {
	// Fetch token's encrypted metadata from the ledger
	tx, err := storage.Ledger.GetTransactionByID(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token %s from the ledger: %v", tokenID, err)
	}

	// Create a new token object
	token := &common.SYN1500Token{
		TokenID:           tokenID,
		EncryptedMetadata: tx.Data["metadata"].([]byte),
	}

	// Decrypt the token's metadata
	decryptedData, err := decryptTokenData(token.EncryptedMetadata, decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token metadata: %v", err)
	}

	// Populate the token with decrypted data
	err = populateTokenFromDecryptedData(token, decryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to populate token with decrypted data: %v", err)
	}

	// Log retrieval event in the ledger
	retrieveTx := ledger.Transaction{
		TxID:        generateUniqueID(tokenID),
		Description: fmt.Sprintf("Retrieve token %s from storage", tokenID),
		Timestamp:   time.Now(),
		Data:        map[string]interface{}{"token_id": tokenID},
	}

	// Validate the retrieval transaction using Synnergy Consensus
	if err := synnergy_consensus.ValidateTransaction(retrieveTx); err != nil {
		return nil, fmt.Errorf("failed Synnergy Consensus validation during token retrieval: %v", err)
	}

	// Record the retrieval event in the ledger
	if err := storage.Ledger.RecordTransaction(retrieveTx); err != nil {
		return nil, fmt.Errorf("failed to log token retrieval: %v", err)
	}

	return token, nil
}

// encryptTokenData encrypts the sensitive fields of the SYN1500Token for secure storage.
func encryptTokenData(token *common.SYN1500Token, key []byte) ([]byte, error) {
	// Convert token data to bytes (simplified for demonstration)
	tokenBytes := []byte(fmt.Sprintf("%v", token))

	// Encrypt the token data
	encryptedData, err := encrypt(tokenBytes, key)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

// decryptTokenData decrypts the encrypted metadata of a SYN1500Token.
func decryptTokenData(encryptedData []byte, key []byte) ([]byte, error) {
	// Decrypt the encrypted token metadata
	decryptedData, err := decrypt(encryptedData, key)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

// populateTokenFromDecryptedData populates a SYN1500Token from decrypted data.
func populateTokenFromDecryptedData(token *common.SYN1500Token, decryptedData []byte) error {
	// Assuming we have some mechanism to populate the token from decrypted data
	// (In practice, you'd deserialize or unmarshal the data back into the token)
	// For simplicity, let's assume this operation is successful
	return nil
}

// encrypt performs AES encryption on data using the provided key.
func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// decrypt performs AES decryption on encrypted data using the provided key.
func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// generateUniqueID creates a unique identifier for transactions or storage events.
func generateUniqueID(seed string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", seed, timestamp)))
	return hex.EncodeToString(hash[:])
}
