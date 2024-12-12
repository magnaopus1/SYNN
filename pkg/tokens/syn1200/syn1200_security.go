package syn1200

import (
	"errors"
	"time"
)

// SYN1200SecurityManager handles security-related functionalities for SYN1200 interoperable tokens.
type SYN1200SecurityManager struct {
	Ledger            *ledger.Ledger                // Ledger integration for tracking and storing encrypted data
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy consensus validation mechanism
	EncryptionService *encryption.EncryptionService // Encryption service for secure token handling
	mutex             sync.Mutex                    // Mutex for thread-safe operations
}


// SecureTokenTransfer facilitates the secure cross-chain transfer of interoperable tokens.
func (sm *SYN1200SecurityManager) SecureTokenTransfer(tokenID, sourceChain, destinationChain string, amount int64) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := sm.GetInteroperableToken(tokenID)
	if err != nil {
		return "", errors.New("token retrieval failed")
	}

	// Check cross-chain compatibility
	if !common.Contains(token.LinkedBlockchains, destinationChain) {
		return "", errors.New("token not compatible with destination blockchain")
	}

	// Create a new transaction record
	transactionID := common.GenerateID()
	transaction := TransactionRecord{
		TransactionID:    transactionID,
		TokenID:          tokenID,
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		Amount:           amount,
		Status:           "pending",
		Timestamp:        time.Now(),
	}

	// Append the transaction record to the token history
	token.TransactionHistory = append(token.TransactionHistory, transaction)

	// Encrypt the token data for storage
	encryptedTokenData, err := sm.EncryptTokenData(tokenID, token)
	if err != nil {
		return "", err
	}

	// Store the encrypted token in the ledger
	if err := sm.Ledger.UpdateToken(tokenID, encryptedTokenData); err != nil {
		return "", errors.New("failed to update token in ledger")
	}

	// Validate the transaction through Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateTransaction(transactionID); err != nil {
		return "", errors.New("transaction validation failed")
	}

	// Mark the transaction as complete
	sm.UpdateTransactionStatus(tokenID, transactionID, "completed")

	return transactionID, nil
}

// EncryptTokenData encrypts the token data before storing it in the ledger.
func (sm *SYN1200SecurityManager) EncryptTokenData(tokenID string, token InteroperableToken) (string, error) {
	// Convert token struct to string
	tokenData := common.StructToString(token)

	// Generate encryption key for the token
	encryptionKey := sm.EncryptionService.GenerateKey()
	encryptedData, err := sm.EncryptionService.EncryptData([]byte(tokenData), encryptionKey)
	if err != nil {
		return "", errors.New("failed to encrypt token data")
	}

	// Store encryption key in the ledger
	if err := sm.Ledger.StoreEncryptionKey(tokenID, encryptionKey); err != nil {
		return "", errors.New("failed to store encryption key")
	}

	return string(encryptedData), nil
}

// DecryptTokenData decrypts token data retrieved from the ledger.
func (sm *SYN1200SecurityManager) DecryptTokenData(tokenID, encryptedData string) (*InteroperableToken, error) {
	// Retrieve the encryption key from the ledger
	encryptionKey, err := sm.Ledger.GetEncryptionKey(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key")
	}

	// Decrypt the data
	decryptedData, err := sm.EncryptionService.DecryptData([]byte(encryptedData), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt token data")
	}

	// Convert the decrypted string back to struct
	var token InteroperableToken
	if err := common.StringToStruct(string(decryptedData), &token); err != nil {
		return nil, errors.New("failed to deserialize token data")
	}

	return &token, nil
}

// UpdateTransactionStatus updates the status of a cross-chain transaction.
func (sm *SYN1200SecurityManager) UpdateTransactionStatus(tokenID, transactionID, status string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the token
	token, err := sm.GetInteroperableToken(tokenID)
	if err != nil {
		return err
	}

	// Update the transaction status
	for i, transaction := range token.TransactionHistory {
		if transaction.TransactionID == transactionID {
			token.TransactionHistory[i].Status = status
			token.TransactionHistory[i].Timestamp = time.Now()
			break
		}
	}

	// Encrypt and update token data in the ledger
	encryptedTokenData, err := sm.EncryptTokenData(tokenID, *token)
	if err != nil {
		return err
	}
	if err := sm.Ledger.UpdateToken(tokenID, encryptedTokenData); err != nil {
		return errors.New("failed to update token in ledger")
	}

	return nil
}

// GetInteroperableToken retrieves and decrypts an interoperable token from the ledger.
func (sm *SYN1200SecurityManager) GetInteroperableToken(tokenID string) (*InteroperableToken, error) {
	// Retrieve encrypted token from ledger
	encryptedToken, err := sm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve token from ledger")
	}

	// Decrypt token data
	token, err := sm.DecryptTokenData(tokenID, encryptedToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
