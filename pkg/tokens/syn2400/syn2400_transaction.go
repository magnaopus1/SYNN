package syn2400

import (
	"errors"
	"time"

)

// SYN2400Transaction handles the transactions and validation of SYN2400 tokens
type SYN2400Transaction struct {
	Ledger      ledger.LedgerInterface          // Interface for interacting with the blockchain ledger
	Encrypt     encryption.EncryptionInterface   // Interface for encryption
	Compress    compression.CompressionInterface // Interface for compression
	Consensus   consensus.ConsensusInterface     // Interface for consensus validation (Synnergy Consensus)
	Audit       audit.AuditInterface             // Interface for auditing activities
}

// NewSYN2400Transaction initializes a new instance of SYN2400Transaction
func NewSYN2400Transaction(ledger ledger.LedgerInterface, encrypt encryption.EncryptionInterface, compress compression.CompressionInterface, consensus consensus.ConsensusInterface, audit audit.AuditInterface) *SYN2400Transaction {
	return &SYN2400Transaction{
		Ledger:   ledger,
		Encrypt:  encrypt,
		Compress: compress,
		Consensus: consensus,
		Audit:    audit,
	}
}

// InitiateTransaction starts a new transaction for SYN2400 token
func (tx *SYN2400Transaction) InitiateTransaction(sender string, receiver string, tokenID string, price float64, compress bool) (common.SYN2400Token, error) {
	// Step 1: Retrieve the token from the ledger
	token, err := tx.Ledger.GetToken(tokenID)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Step 2: Decrypt the token
	decryptedToken, err := tx.Encrypt.DecryptTokenData(token)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to decrypt token: " + err.Error())
	}

	// Step 3: Verify ownership and initiate transaction
	if decryptedToken.Owner != sender {
		return common.SYN2400Token{}, errors.New("unauthorized transaction attempt, sender is not the owner")
	}

	// Step 4: Update the ownership of the token
	decryptedToken.Owner = receiver
	decryptedToken.Price = price
	decryptedToken.UpdateDate = time.Now()

	// Step 5: Re-encrypt the updated token
	updatedEncryptedToken, err := tx.Encrypt.EncryptTokenData(decryptedToken)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to encrypt updated token: " + err.Error())
	}

	// Step 6: Optionally compress the encrypted token
	if compress {
		updatedEncryptedToken, err = tx.Compress.CompressData(updatedEncryptedToken)
		if err != nil {
			return common.SYN2400Token{}, errors.New("failed to compress encrypted token: " + err.Error())
		}
	}

	// Step 7: Store the updated token in the ledger
	err = tx.Ledger.UpdateToken(tokenID, updatedEncryptedToken)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to store updated token in ledger: " + err.Error())
	}

	// Step 8: Validate the transaction using Synnergy Consensus
	subBlock, err := tx.Consensus.ValidateSubBlockTransaction(tokenID, sender, receiver)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to validate sub-block transaction: " + err.Error())
	}

	// Step 9: Finalize the transaction into a block
	err = tx.Consensus.FinalizeBlock(subBlock)
	if err != nil {
		return common.SYN2400Token{}, errors.New("failed to finalize block: " + err.Error())
	}

	// Step 10: Audit the transaction
	tx.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "InitiateTransaction",
		PerformedBy: sender,
		Timestamp:   time.Now(),
		Details:     "Transferred token " + tokenID + " from " + sender + " to " + receiver,
	})

	return decryptedToken, nil
}

// ValidateTransaction verifies the transaction before finalizing it
func (tx *SYN2400Transaction) ValidateTransaction(tokenID string, sender string, receiver string) error {
	// Step 1: Retrieve the token from the ledger
	token, err := tx.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Step 2: Decrypt the token
	decryptedToken, err := tx.Encrypt.DecryptTokenData(token)
	if err != nil {
		return errors.New("failed to decrypt token: " + err.Error())
	}

	// Step 3: Validate ownership and transaction details
	if decryptedToken.Owner != sender {
		return errors.New("ownership validation failed, unauthorized transaction attempt by sender")
	}

	// Step 4: Validate transaction using Synnergy Consensus
	_, err = tx.Consensus.ValidateSubBlockTransaction(tokenID, sender, receiver)
	if err != nil {
		return errors.New("failed to validate transaction using consensus: " + err.Error())
	}

	// Step 5: Audit the validation event
	tx.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "ValidateTransaction",
		PerformedBy: sender,
		Timestamp:   time.Now(),
		Details:     "Validated transaction for token " + tokenID,
	})

	return nil
}

// CompleteTransaction finalizes the transaction and ensures the token is fully transferred
func (tx *SYN2400Transaction) CompleteTransaction(tokenID string, receiver string) error {
	// Step 1: Retrieve the token from the ledger
	token, err := tx.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Step 2: Decrypt the token
	decryptedToken, err := tx.Encrypt.DecryptTokenData(token)
	if err != nil {
		return errors.New("failed to decrypt token: " + err.Error())
	}

	// Step 3: Verify that the receiver is the current owner
	if decryptedToken.Owner != receiver {
		return errors.New("transaction completion failed, receiver is not the current owner")
	}

	// Step 4: Finalize the block using Synnergy Consensus
	err = tx.Consensus.FinalizeBlock(decryptedToken.TokenID)
	if err != nil {
		return errors.New("failed to finalize the block: " + err.Error())
	}

	// Step 5: Audit the transaction completion event
	tx.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "CompleteTransaction",
		PerformedBy: receiver,
		Timestamp:   time.Now(),
		Details:     "Completed transaction for token " + tokenID + " to receiver " + receiver,
	})

	return nil
}

// RevertTransaction allows rolling back a transaction in case of failure
func (tx *SYN2400Transaction) RevertTransaction(tokenID string, previousOwner string) error {
	// Retrieve the token and decrypt
	token, err := tx.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	decryptedToken, err := tx.Encrypt.DecryptTokenData(token)
	if err != nil {
		return errors.New("failed to decrypt token: " + err.Error())
	}

	// Revert the ownership
	decryptedToken.Owner = previousOwner
	decryptedToken.UpdateDate = time.Now()

	// Re-encrypt the updated token
	encryptedToken, err := tx.Encrypt.EncryptTokenData(decryptedToken)
	if err != nil {
		return errors.New("failed to encrypt reverted token: " + err.Error())
	}

	// Update the token in the ledger
	err = tx.Ledger.UpdateToken(tokenID, encryptedToken)
	if err != nil {
		return errors.New("failed to update reverted token in ledger: " + err.Error())
	}

	// Log the audit for the revert event
	tx.Audit.LogAuditEvent(audit.AuditRecord{
		Action:      "RevertTransaction",
		PerformedBy: "System",
		Timestamp:   time.Now(),
		Details:     "Reverted transaction for token " + tokenID + " to previous owner " + previousOwner,
	})

	return nil
}

