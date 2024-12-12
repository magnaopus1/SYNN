package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/transactions"
	"time"
)

// WalletTransactionService handles creating, signing, broadcasting, and validating transactions for wallets.
type WalletTransactionService struct {
	ledgerInstance  *ledger.Ledger
	transactionPool *transactions.TransactionPool // Use actual TransactionPool struct
	mutex           sync.Mutex
	networkManager  *network.NetworkManager
}

// NewWalletTransactionService initializes the WalletTransactionService with the ledger and network.
func NewWalletTransactionService(ledgerInstance *ledger.Ledger, networkManager *network.NetworkManager, encryptionService *common.Encryption) *WalletTransactionService {
	maxPoolSize := 1000000 // Example pool size, adjust this value as necessary

	return &WalletTransactionService{
		ledgerInstance:  ledgerInstance,
		networkManager:  networkManager,
		transactionPool: transactions.NewTransactionPool(maxPoolSize, ledgerInstance, encryptionService), // Pass encryptionService as a pointer
	}
}


// CreateTransaction creates a new transaction from the wallet.
func (wts *WalletTransactionService) CreateTransaction(senderAddress, recipientAddress string, amount float64, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// Step 1: Validate inputs
	if amount <= 0 {
		return nil, errors.New("invalid transaction amount")
	}
	if senderAddress == "" || recipientAddress == "" {
		return nil, errors.New("sender or recipient address cannot be empty")
	}

	// Step 2: Create the transaction structure
	transactionData := fmt.Sprintf("%s:%s:%f:%d", senderAddress, recipientAddress, amount, time.Now().UnixNano())
	transactionHash := sha256.Sum256([]byte(transactionData))

	// Step 3: Sign the transaction with the wallet's private key
	r, s, err := ecdsa.Sign(nil, privateKey, transactionHash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Step 4: Assemble the signed transaction
	signedTransaction := fmt.Sprintf("%s:%x:%x", transactionData, r, s)

	// Step 5: Create the encryption instance and encrypt the transaction
	encryption := &common.Encryption{} // Encryption instance created inside the function
	encryptedTransaction, err := encryption.EncryptData("AES", []byte(signedTransaction), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Step 6: Store the encrypted transaction in the transaction pool
	transactionID := hex.EncodeToString(transactionHash[:])
	newTransaction := &common.Transaction{
		TransactionID: transactionID,
		FromAddress:   senderAddress,
		ToAddress:     recipientAddress,
		Amount:        amount,
		EncryptedData: hex.EncodeToString(encryptedTransaction),
		Timestamp:     time.Now(),
		Status:        "pending",
	}

	// Add the transaction to the pool using the AddTransaction method of the TransactionPool
	err = wts.transactionPool.AddTransaction(newTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to add transaction to the pool: %v", err)
	}

	// Step 7: Log the transaction creation event in the ledger (fix: pass 3 arguments)
	err = wts.ledgerInstance.RecordTransaction(senderAddress, recipientAddress, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to record transaction in ledger: %v", err)
	}

	return encryptedTransaction, nil
}



// BroadcastTransaction broadcasts a signed and encrypted transaction to the blockchain network.
func (wts *WalletTransactionService) BroadcastTransaction(transactionID string) error {
	// Step 1: Retrieve the transaction from the transaction pool
	transaction, err := wts.transactionPool.GetTransaction(transactionID)
	if err != nil {
		return fmt.Errorf("transaction not found in the pool: %v", err)
	}

	// Step 2: Broadcast the transaction via the network manager
	if err := wts.networkManager.SendEncryptedMessage(transactionID, string(transaction.EncryptedData)); err != nil {
		return fmt.Errorf("failed to broadcast transaction: %v", err)
	}

	// Step 3: Log the broadcast event in the ledger
	if err := wts.ledgerInstance.RecordTransactionBroadcast(transactionID); err != nil {
		return fmt.Errorf("failed to log transaction broadcast: %v", err)
	}

	// Step 4: Remove the transaction from the pool after broadcasting (no value returned)
	wts.transactionPool.RemoveTransaction(transactionID)

	return nil
}




// ValidateTransaction validates the integrity of a transaction, ensuring its correctness and validity.
func (wts *WalletTransactionService) ValidateTransaction(transaction []byte, publicKey *ecdsa.PublicKey) (bool, error) {
	// Step 1: Decrypt the transaction using the encryption package
	encryptionService := &common.Encryption{} // Create an instance of the encryption service
	decryptedTransaction, err := encryptionService.DecryptData(transaction, common.EncryptionKey)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt transaction: %v", err)
	}

	// Step 2: Parse the transaction (expected format: "sender:recipient:amount:timestamp:r:s")
	var sender, recipient string
	var amount float64
	var timestamp int64
	var r, s string
	_, err = fmt.Sscanf(string(decryptedTransaction), "%s:%s:%f:%d:%x:%x", &sender, &recipient, &amount, &timestamp, &r, &s)
	if err != nil {
		return false, fmt.Errorf("failed to parse transaction: %v", err)
	}

	// Step 3: Recreate the original transaction data (before signing)
	transactionData := fmt.Sprintf("%s:%s:%f:%d", sender, recipient, amount, timestamp)
	transactionHash := sha256.Sum256([]byte(transactionData))

	// Step 4: Verify the signature using the public key
	var rInt, sInt big.Int
	rInt.SetString(r, 16)
	sInt.SetString(s, 16)

	if ecdsa.Verify(publicKey, transactionHash[:], &rInt, &sInt) {
		return true, nil
	}
	return false, errors.New("transaction signature verification failed")
}

// GetTransactionStatus checks the status of a given transaction (e.g., pending, confirmed).
func (wts *WalletTransactionService) GetTransactionStatus(transactionID string) (string, error) {
	// Check if the transaction exists in the ledger
	status, err := wts.ledgerInstance.GetTransactionStatus(transactionID)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction status: %v", err)
	}

	return status, nil
}

// ConfirmSubBlock validates and confirms transactions within a sub-block.
func (wts *WalletTransactionService) ConfirmSubBlock(subBlock []byte) error {
	// Convert subBlock from []byte to string
	subBlockStr := string(subBlock)

	// Validate sub-block and integrate with ledger
	if err := wts.ledgerInstance.ValidateSubBlock(subBlockStr); err != nil {
		return fmt.Errorf("sub-block validation failed: %v", err)
	}

	return nil
}

// ConfirmBlock validates a full block containing 1000 sub-blocks.
func (wts *WalletTransactionService) ConfirmBlock(block [][]byte) error {
	if len(block) != 1000 {
		return errors.New("invalid block size: block must contain 1000 sub-blocks")
	}

	// Loop through each sub-block for validation
	for _, subBlock := range block {
		if err := wts.ConfirmSubBlock(subBlock); err != nil {
			return fmt.Errorf("block validation failed at sub-block: %v", err)
		}
	}

	// Concatenate the entire block into one string
	var blockStr string
	for _, subBlock := range block {
		blockStr += string(subBlock)
	}

	// After all sub-blocks are validated, log the confirmed block in the ledger
	if err := wts.ledgerInstance.RecordConfirmedBlock(blockStr); err != nil {
		return fmt.Errorf("failed to record confirmed block in ledger: %v", err)
	}

	return nil
}

