package syn1000

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

// SYN1000StorageManager handles storage and retrieval of SYN1000 tokens and their associated data
type SYN1000StorageManager struct {
	TokenStorage        map[string]string // TokenID -> encrypted token data
	TransactionHistory  map[string][]string // TokenID -> encrypted transaction history
	MetadataStore       map[string]string // TokenID -> encrypted metadata
	ReserveStore        map[string]string // TokenID -> encrypted reserve data
	mutex               sync.RWMutex
	Ledger              *ledger.Ledger
	ConsensusEngine     *consensus.SynnergyConsensus
	EncryptionService   *encryption.EncryptionService
}

// NewSYN1000StorageManager creates a new instance of the storage manager
func NewSYN1000StorageManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN1000StorageManager {
	return &SYN1000StorageManager{
		TokenStorage:       make(map[string]string),
		TransactionHistory: make(map[string][]string),
		MetadataStore:      make(map[string]string),
		ReserveStore:       make(map[string]string),
		Ledger:             ledger,
		ConsensusEngine:    consensusEngine,
		EncryptionService:  encryptionService,
	}
}

// StoreToken stores SYN1000 token data securely, integrating encryption
func (sm *SYN1000StorageManager) StoreToken(tokenID string, tokenData interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize and encrypt the token data
	dataBytes, err := json.Marshal(tokenData)
	if err != nil {
		return err
	}

	encryptedData, encryptionKey, err := sm.EncryptData(dataBytes)
	if err != nil {
		return err
	}

	// Store the encrypted data in the storage
	sm.TokenStorage[tokenID] = encryptedData

	// Store in the ledger for immutability
	return sm.Ledger.StoreToken(tokenID, encryptedData, encryptionKey)
}

// RetrieveToken retrieves and decrypts SYN1000 token data
func (sm *SYN1000StorageManager) RetrieveToken(tokenID string) (interface{}, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve encrypted token data from storage
	encryptedData, exists := sm.TokenStorage[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Retrieve from ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from ledger")
	}

	// Decrypt the token data
	decryptedData, err := sm.DecryptData(encryptedData, encryptionKey)
	if err != nil {
		return nil, err
	}

	var token interface{}
	if err := json.Unmarshal(decryptedData, &token); err != nil {
		return nil, err
	}

	return token, nil
}

// StoreTransaction records encrypted transaction history for a SYN1000 token
func (sm *SYN1000StorageManager) StoreTransaction(tokenID string, transactionData interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize and encrypt transaction data
	dataBytes, err := json.Marshal(transactionData)
	if err != nil {
		return err
	}

	encryptedData, encryptionKey, err := sm.EncryptData(dataBytes)
	if err != nil {
		return err
	}

	// Append encrypted transaction data to the history
	sm.TransactionHistory[tokenID] = append(sm.TransactionHistory[tokenID], encryptedData)

	// Store transaction in ledger
	return sm.Ledger.StoreTransaction(tokenID, encryptedData, encryptionKey)
}

// RetrieveTransactions retrieves all transaction history for a specific token, decrypting data
func (sm *SYN1000StorageManager) RetrieveTransactions(tokenID string) ([]interface{}, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve encrypted transaction history from storage
	encryptedTransactions, exists := sm.TransactionHistory[tokenID]
	if !exists {
		return nil, errors.New("no transaction history found for the token")
	}

	// Retrieve encryption key from ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from ledger")
	}

	var transactions []interface{}
	for _, encryptedTransaction := range encryptedTransactions {
		// Decrypt transaction data
		decryptedData, err := sm.DecryptData(encryptedTransaction, encryptionKey)
		if err != nil {
			return nil, err
		}

		var transaction interface{}
		if err := json.Unmarshal(decryptedData, &transaction); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// StoreMetadata securely stores metadata for a token
func (sm *SYN1000StorageManager) StoreMetadata(tokenID string, metadata interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize and encrypt metadata
	dataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	encryptedData, encryptionKey, err := sm.EncryptData(dataBytes)
	if err != nil {
		return err
	}

	// Store encrypted metadata
	sm.MetadataStore[tokenID] = encryptedData

	// Store in the ledger
	return sm.Ledger.StoreMetadata(tokenID, encryptedData, encryptionKey)
}

// RetrieveMetadata retrieves and decrypts metadata for a token
func (sm *SYN1000StorageManager) RetrieveMetadata(tokenID string) (interface{}, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve encrypted metadata from storage
	encryptedData, exists := sm.MetadataStore[tokenID]
	if !exists {
		return nil, errors.New("metadata not found")
	}

	// Retrieve encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt metadata
	decryptedData, err := sm.DecryptData(encryptedData, encryptionKey)
	if err != nil {
		return nil, err
	}

	var metadata interface{}
	if err := json.Unmarshal(decryptedData, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

// StoreReserve securely stores reserve data for a token
func (sm *SYN1000StorageManager) StoreReserve(tokenID string, reserveData interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Serialize and encrypt reserve data
	dataBytes, err := json.Marshal(reserveData)
	if err != nil {
		return err
	}

	encryptedData, encryptionKey, err := sm.EncryptData(dataBytes)
	if err != nil {
		return err
	}

	// Store encrypted reserve data
	sm.ReserveStore[tokenID] = encryptedData

	// Store in the ledger
	return sm.Ledger.StoreReserve(tokenID, encryptedData, encryptionKey)
}

// RetrieveReserve retrieves and decrypts reserve data for a token
func (sm *SYN1000StorageManager) RetrieveReserve(tokenID string) (interface{}, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Retrieve encrypted reserve data from storage
	encryptedData, exists := sm.ReserveStore[tokenID]
	if !exists {
		return nil, errors.New("reserve data not found")
	}

	// Retrieve encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt reserve data
	decryptedData, err := sm.DecryptData(encryptedData, encryptionKey)
	if err != nil {
		return nil, err
	}

	var reserveData interface{}
	if err := json.Unmarshal(decryptedData, &reserveData); err != nil {
		return nil, err
	}

	return reserveData, nil
}

// EncryptData encrypts data for storage
func (sm *SYN1000StorageManager) EncryptData(data []byte) (string, string, error) {
	block, err := aes.NewCipher(sm.EncryptionService.Key[:])
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", "", err
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encryptedData), base64.StdEncoding.EncodeToString(sm.EncryptionService.Key[:]), nil
}

// DecryptData decrypts stored data
func (sm *SYN1000StorageManager) DecryptData(encryptedData string, encryptionKey string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(decodedData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := decodedData[:nonceSize], decodedData[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
