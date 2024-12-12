package syn2500

import (
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
)

// DAOTokenTransaction represents a DAO token transaction.
type DAOTokenTransaction struct {
	TransactionID  string              // Unique identifier for the transaction
	TokenID        string              // DAO Token involved in the transaction
	Sender         string              // Address of the sender
	Recipient      string              // Address of the recipient
	Timestamp      time.Time           // Timestamp of the transaction
	TransactionHash string             // Hash of the transaction for verification
	Signature      string              // Cryptographic signature of the transaction
}

// NewDAOTokenTransaction creates a new DAO token transaction.
func NewDAOTokenTransaction(tokenID, sender, recipient string) (*DAOTokenTransaction, error) {
	// Check if the token exists and is valid
	token, err := ledger.GetDAOToken(tokenID)
	if err != nil {
		return nil, errors.New("DAO token not found in the ledger")
	}

	// Ensure the sender owns the token
	if token.Owner != sender {
		return nil, errors.New("sender is not the owner of the DAO token")
	}

	// Create a new transaction
	tx := &DAOTokenTransaction{
		TransactionID:  generateUniqueTransactionID(),
		TokenID:        tokenID,
		Sender:         sender,
		Recipient:      recipient,
		Timestamp:      time.Now(),
	}

	// Generate transaction hash
	tx.TransactionHash = generateTransactionHash(tx)

	// Sign the transaction (signature encryption)
	tx.Signature, err = encryption.SignTransaction(tx.TransactionHash, sender)
	if err != nil {
		return nil, errors.New("failed to sign the transaction")
	}

	return tx, nil
}

// ValidateTransaction validates the DAO token transaction using Synnergy Consensus.
func (tx *DAOTokenTransaction) ValidateTransaction() error {
	// Verify the transaction's signature
	valid, err := encryption.VerifySignature(tx.TransactionHash, tx.Signature, tx.Sender)
	if err != nil || !valid {
		return errors.New("invalid transaction signature")
	}

	// Validate the transaction through Synnergy Consensus
	err = synconsensus.ValidateSubBlock(tx.TransactionHash)
	if err != nil {
		return errors.New("transaction validation failed through Synnergy Consensus")
	}

	return nil
}

// ExecuteTransaction executes the DAO token transaction and updates the ledger.
func (tx *DAOTokenTransaction) ExecuteTransaction() error {
	// Ensure the transaction is valid
	err := tx.ValidateTransaction()
	if err != nil {
		return err
	}

	// Retrieve the token and update its owner
	token, err := ledger.GetDAOToken(tx.TokenID)
	if err != nil {
		return errors.New("DAO token not found in the ledger")
	}
	token.Owner = tx.Recipient
	token.Timestamp = time.Now()

	// Store the updated token in the ledger
	err = ledger.StoreDAOToken(token)
	if err != nil {
		return errors.New("failed to update DAO token in the ledger")
	}

	// Log the transaction in the ledger
	err = ledger.StoreDAOTokenTransaction(tx)
	if err != nil {
		return errors.New("failed to store DAO token transaction in the ledger")
	}

	return nil
}

// RevokeTransaction allows the revocation of a DAO token transaction if conditions allow.
func (tx *DAOTokenTransaction) RevokeTransaction(reason string) error {
	// Check if the transaction can be revoked based on time or other conditions
	if time.Since(tx.Timestamp).Hours() > 24 {
		return errors.New("transaction cannot be revoked after 24 hours")
	}

	// Log the revocation in the ledger
	err := ledger.RevokeDAOTokenTransaction(tx.TransactionID, reason)
	if err != nil {
		return errors.New("failed to revoke the DAO token transaction in the ledger")
	}

	return nil
}

// generateTransactionHash generates a unique hash for the transaction.
func generateTransactionHash(tx *DAOTokenTransaction) string {
	hashInput := tx.TokenID + tx.Sender + tx.Recipient + tx.Timestamp.String()
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

// generateUniqueTransactionID generates a unique transaction ID using a combination of timestamp and hash.
func generateUniqueTransactionID() string {
	currentTime := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(string(currentTime)))
	return hex.EncodeToString(hash[:])
}
