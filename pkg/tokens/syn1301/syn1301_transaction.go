package syn1301

import (
	"errors"
	"time"

)

// SYN1301TransactionManager handles the processing of transactions for the SYN1301 Token Standard.
type SYN1301TransactionManager struct {
	Ledger            *ledger.Ledger                // Ledger instance for token data and transaction management
	EncryptionService *encryption.EncryptionService // Service for encryption and decryption of transaction data
	Consensus         *synnergy_consensus.Consensus // Synnergy Consensus system for transaction validation
}

// ProcessTransaction handles the full life cycle of a token transaction, including validation, encryption, and ledger integration.
func (tm *SYN1301TransactionManager) ProcessTransaction(tokenID string, senderID string, receiverID string, metadata map[string]string) error {
	// Step 1: Validate the transaction through Synnergy Consensus
	subBlock, err := tm.Consensus.ValidateTransactionIntoSubBlock(tokenID, metadata)
	if err != nil {
		return errors.New("failed to validate transaction into sub-block: " + err.Error())
	}

	// Step 2: Validate sub-block into a full block
	block, err := tm.Consensus.ValidateSubBlockIntoBlock(subBlock)
	if err != nil {
		return errors.New("failed to validate sub-block into block: " + err.Error())
	}

	// Step 3: Encrypt the transaction metadata before storing it in the ledger
	encryptedMetadata, err := tm.EncryptionService.Encrypt(metadata)
	if err != nil {
		return errors.New("encryption of transaction metadata failed: " + err.Error())
	}

	// Step 4: Record the transaction in the ledger
	transactionRecord := ledger.Transaction{
		TokenID:           tokenID,
		SenderID:          senderID,
		ReceiverID:        receiverID,
		EncryptedMetadata: encryptedMetadata,
		BlockID:           block.BlockID,
		Timestamp:         time.Now(),
	}
	err = tm.Ledger.RecordTransaction(transactionRecord)
	if err != nil {
		return errors.New("failed to record transaction in the ledger: " + err.Error())
	}

	// Transaction successfully processed and secured with encryption
	return nil
}

// ValidateAndPrepareToken validates and prepares a token for a new transaction, ensuring it meets all conditions for transfer.
func (tm *SYN1301TransactionManager) ValidateAndPrepareToken(tokenID string, senderID string, receiverID string) (*common.Token, error) {
	// Step 1: Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("token not found in ledger: " + err.Error())
	}

	// Step 2: Validate token ownership
	if token.OwnerID != senderID {
		return nil, errors.New("token ownership validation failed: sender is not the token owner")
	}

	// Step 3: Ensure the token is not locked or restricted for transfer
	if token.IsLocked || token.Restricted {
		return nil, errors.New("token cannot be transferred: it is either locked or restricted")
	}

	// Token is validated and ready for the transaction
	return token, nil
}

// RecordTokenTransfer updates the token's ownership and logs the transfer in the ledger after a successful transaction.
func (tm *SYN1301TransactionManager) RecordTokenTransfer(tokenID string, newOwnerID string) error {
	// Step 1: Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Step 2: Update the token's ownership
	token.OwnerID = newOwnerID

	// Step 3: Encrypt the updated token metadata
	updatedMetadata := map[string]string{
		"owner": newOwnerID,
		"status": "transferred",
	}
	encryptedMetadata, err := tm.EncryptionService.Encrypt(updatedMetadata)
	if err != nil {
		return errors.New("failed to encrypt updated token metadata: " + err.Error())
	}
	token.EncryptedMetadata = encryptedMetadata

	// Step 4: Update the token in the ledger
	err = tm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return errors.New("failed to update token ownership in ledger: " + err.Error())
	}

	// Step 5: Log the token transfer in the ledger
	err = tm.Ledger.LogEvent(ledger.EventLog{
		EventType:   "TRANSFER",
		TokenID:     tokenID,
		UserID:      newOwnerID,
		Description: "Ownership transferred",
		Timestamp:   time.Now(),
	})
	if err != nil {
		return errors.New("failed to log token transfer event: " + err.Error())
	}

	// Ownership transfer successfully recorded and encrypted
	return nil
}

// SecureTransactionData encrypts and securely stores transaction metadata.
func (tm *SYN1301TransactionManager) SecureTransactionData(transactionData map[string]string) (string, error) {
	encryptedData, err := tm.EncryptionService.Encrypt(transactionData)
	if err != nil {
		return "", errors.New("failed to encrypt transaction data: " + err.Error())
	}
	return encryptedData, nil
}

// DecryptTransactionData decrypts the encrypted transaction metadata for further use.
func (tm *SYN1301TransactionManager) DecryptTransactionData(encryptedData string) (map[string]string, error) {
	decryptedData, err := tm.EncryptionService.Decrypt(encryptedData)
	if err != nil {
		return nil, errors.New("failed to decrypt transaction data: " + err.Error())
	}
	return decryptedData, nil
}

// ReverseTransaction securely reverses a token transaction, reverting the token ownership and updating the ledger accordingly.
func (tm *SYN1301TransactionManager) ReverseTransaction(tokenID, senderID, receiverID string) error {
	// Step 1: Validate and retrieve the token
	token, err := tm.ValidateAndPrepareToken(tokenID, receiverID, senderID)
	if err != nil {
		return errors.New("failed to validate token for reversal: " + err.Error())
	}

	// Step 2: Record the reversal in the ledger and update ownership
	token.OwnerID = senderID
	updatedMetadata := map[string]string{
		"owner": senderID,
		"status": "reversed",
	}
	encryptedMetadata, err := tm.EncryptTransactionData(updatedMetadata)
	if err != nil {
		return errors.New("failed to encrypt reversed token metadata: " + err.Error())
	}
	token.EncryptedMetadata = encryptedMetadata

	// Step 3: Update the token and ledger
	err = tm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return errors.New("failed to update token for reversal in ledger: " + err.Error())
	}

	// Step 4: Log the reversal in the ledger
	err = tm.Ledger.LogEvent(ledger.EventLog{
		EventType:   "REVERSE",
		TokenID:     tokenID,
		UserID:      senderID,
		Description: "Ownership reversal",
		Timestamp:   time.Now(),
	})
	if err != nil {
		return errors.New("failed to log transaction reversal event: " + err.Error())
	}

	// Reversal successfully completed
	return nil
}
