package syn11

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

  	"synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// Syn11Transaction represents a token transaction for SYN11.
type Syn11Transaction struct {
	TokenID        string                      // Token ID for the transaction
	FromAddress    string                      // Sender address
	ToAddress      string                      // Receiver address
	Amount         uint64                      // Amount being transferred
	Timestamp      time.Time                   // Timestamp of the transaction
	VerificationID string                      // Verification ID for the transaction
}

// Syn11TransactionManager handles transaction processing for SYN11 tokens.
type Syn11TransactionManager struct {
	mutex          sync.Mutex
	Ledger         *ledger.Ledger                // Ledger for recording transactions
	Consensus      *consensus.SynnergyConsensus  // Consensus engine for transaction validation
	Encryption     *encryption.EncryptionService // Encryption service for securing transaction data
	Compliance     *compliance.KYCAmlService     // KYC/AML Compliance service for regulatory checks
}

// NewSyn11TransactionManager initializes a new Syn11TransactionManager.
func NewSyn11TransactionManager(ledger *ledger.Ledger, consensus *consensus.SynnergyConsensus, encryption *encryption.EncryptionService, compliance *compliance.KYCAmlService) *Syn11TransactionManager {
	return &Syn11TransactionManager{
		Ledger:     ledger,
		Consensus:  consensus,
		Encryption: encryption,
		Compliance: compliance,
	}
}

// ProcessTransaction processes a token transaction.
func (tm *Syn11TransactionManager) ProcessTransaction(tx Syn11Transaction) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// KYC/AML Compliance Check for the sender and receiver
	if err := tm.Compliance.VerifyUser(tx.FromAddress); err != nil {
		return fmt.Errorf("KYC/AML verification failed for sender: %v", err)
	}
	if err := tm.Compliance.VerifyUser(tx.ToAddress); err != nil {
		return fmt.Errorf("KYC/AML verification failed for receiver: %v", err)
	}

	// Consensus validation before proceeding with the transaction
	if err := tm.Consensus.ValidateTransaction(tx.TokenID, tx.FromAddress, tx.ToAddress, tx.Amount); err != nil {
		return fmt.Errorf("consensus validation failed: %v", err)
	}

	// Encrypt transaction details before storing
	encryptedTxID, err := tm.Encryption.Encrypt([]byte(tx.VerificationID))
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction ID: %v", err)
	}

	// Store the transaction in the ledger
	err = tm.Ledger.RecordTransaction(tx.TokenID, tx.FromAddress, tx.ToAddress, tx.Amount, tx.Timestamp, string(encryptedTxID))
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}

	log.Printf("Successfully processed transaction: %v from %s to %s", tx.TokenID, tx.FromAddress, tx.ToAddress)
	return nil
}

// ValidateTransaction performs validation on a transaction.
func (tm *Syn11TransactionManager) ValidateTransaction(tx Syn11Transaction) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Check if the sender has sufficient balance
	balance, err := tm.Ledger.GetBalance(tx.FromAddress)
	if err != nil {
		return fmt.Errorf("failed to retrieve balance: %v", err)
	}
	if balance < tx.Amount {
		return errors.New("insufficient balance for the transaction")
	}

	// Consensus validation
	if err := tm.Consensus.ValidateTransaction(tx.TokenID, tx.FromAddress, tx.ToAddress, tx.Amount); err != nil {
		return fmt.Errorf("consensus validation failed: %v", err)
	}

	return nil
}

// ReverseTransaction handles the reversal of a transaction in case of a dispute or error.
func (tm *Syn11TransactionManager) ReverseTransaction(tx Syn11Transaction) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Consensus authorization for reversing the transaction
	if err := tm.Consensus.AuthorizeTransactionReversal(tx.TokenID, tx.FromAddress, tx.ToAddress, tx.Amount); err != nil {
		return fmt.Errorf("transaction reversal not authorized: %v", err)
	}

	// Reverse the transaction in the ledger
	err := tm.Ledger.ReverseTransaction(tx.TokenID, tx.FromAddress, tx.ToAddress, tx.Amount)
	if err != nil {
		return fmt.Errorf("failed to reverse transaction: %v", err)
	}

	log.Printf("Successfully reversed transaction: %v from %s to %s", tx.TokenID, tx.FromAddress, tx.ToAddress)
	return nil
}

// RetrieveTransaction retrieves transaction details using the transaction ID.
func (tm *Syn11TransactionManager) RetrieveTransaction(txID string) (*Syn11Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Decrypt the transaction ID
	decryptedTxID, err := tm.Encryption.Decrypt([]byte(txID))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transaction ID: %v", err)
	}

	// Retrieve the transaction from the ledger
	txDetails, err := tm.Ledger.GetTransaction(string(decryptedTxID))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	// Return the transaction details
	return &Syn11Transaction{
		TokenID:        txDetails.TokenID,
		FromAddress:    txDetails.FromAddress,
		ToAddress:      txDetails.ToAddress,
		Amount:         txDetails.Amount,
		Timestamp:      txDetails.Timestamp,
		VerificationID: txDetails.VerificationID,
	}, nil
}

// ListTransactions lists all transactions involving a particular token.
func (tm *Syn11TransactionManager) ListTransactions(tokenID string) ([]Syn11Transaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the transactions from the ledger
	transactions, err := tm.Ledger.GetTransactionsByToken(tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions: %v", err)
	}

	// Format the transactions into a slice of Syn11Transaction
	var txList []Syn11Transaction
	for _, tx := range transactions {
		txList = append(txList, Syn11Transaction{
			TokenID:        tx.TokenID,
			FromAddress:    tx.FromAddress,
			ToAddress:      tx.ToAddress,
			Amount:         tx.Amount,
			Timestamp:      tx.Timestamp,
			VerificationID: tx.VerificationID,
		})
	}

	return txList, nil
}
