package transactions

import (
	"crypto/sha256"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
)

// TransactionReceipt represents a receipt for a processed transaction.
type TransactionReceipt struct {
	TransactionID      string    // Unique ID of the transaction
	BlockID            string    // ID of the block the transaction is part of
	SubBlockID         string    // ID of the sub-block the transaction is part of
	Status             string    // Transaction status: "SUCCESS", "FAILURE"
	Timestamp          time.Time // Timestamp when the transaction was processed
	GasUsed            uint64    // Amount of gas used by the transaction
	TransactionOutput  string    // Output data from the transaction execution
	ValidatorSignature string    // Validator signature confirming the transaction inclusion
	EncryptionHash     string    // Hash of the encrypted transaction data
}

// ReceiptManager manages the creation and handling of transaction receipts.
type ReceiptManager struct {
	receipts          map[string]*TransactionReceipt // Stores transaction receipts by transaction ID
	encryptionService common.Encryption              // Encryption service for handling encryption
}

// NewReceiptManager initializes a new ReceiptManager.
func NewReceiptManager(encryptionService common.Encryption) *ReceiptManager {
	return &ReceiptManager{
		receipts:          make(map[string]*TransactionReceipt),
		encryptionService: encryptionService,
	}
}

// CreateReceipt creates a receipt for a processed transaction.
func (rm *ReceiptManager) CreateReceipt(tx *common.Transaction, blockID, subBlockID string, gasUsed uint64, status, output string, validatorSignature string) (*TransactionReceipt, error) {
	// Generate timestamp
	timestamp := time.Now()

	// Encrypt the transaction data for integrity
	encryptedTx, err := rm.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%+v", tx)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt transaction data for receipt: %v", err)
	}

	// Create a hash of the encrypted transaction
	hash := sha256.Sum256(encryptedTx)
	hashStr := fmt.Sprintf("%x", hash)

	// Create a new transaction receipt
	receipt := &TransactionReceipt{
		TransactionID:      tx.TransactionID, // Use the correct field from Transaction struct
		BlockID:            blockID,
		SubBlockID:         subBlockID,
		Status:             status,
		Timestamp:          timestamp,
		GasUsed:            gasUsed,
		TransactionOutput:  output,
		ValidatorSignature: validatorSignature,
		EncryptionHash:     hashStr,
	}

	// Store the receipt
	rm.receipts[tx.TransactionID] = receipt

	return receipt, nil
}


// GetReceipt retrieves the receipt for a given transaction ID.
func (rm *ReceiptManager) GetReceipt(txID string) (*TransactionReceipt, error) {
	receipt, exists := rm.receipts[txID]
	if !exists {
		return nil, fmt.Errorf("transaction receipt not found")
	}
	return receipt, nil
}

// VerifyReceipt verifies the integrity and authenticity of a transaction receipt.
func (rm *ReceiptManager) VerifyReceipt(receipt *TransactionReceipt, tx *common.Transaction) (bool, error) {
	// Recompute the hash of the encrypted transaction
	encryptedTx, err := rm.encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%+v", tx)), common.EncryptionKey)
	if err != nil {
		return false, fmt.Errorf("failed to re-encrypt transaction data: %v", err)
	}

	// Compute the hash of the encrypted transaction
	hash := sha256.Sum256(encryptedTx)
	hashStr := fmt.Sprintf("%x", hash)

	// Compare the computed hash with the one stored in the receipt
	if receipt.EncryptionHash != hashStr {
		return false, fmt.Errorf("receipt verification failed: transaction hash mismatch")
	}


	return true, nil
}


// ListReceipts returns all stored transaction receipts.
func (rm *ReceiptManager) ListReceipts() []*TransactionReceipt {
	receipts := []*TransactionReceipt{}
	for _, receipt := range rm.receipts {
		receipts = append(receipts, receipt)
	}
	return receipts
}

// RemoveReceipt removes a transaction receipt from the system.
func (rm *ReceiptManager) RemoveReceipt(txID string) error {
	if _, exists := rm.receipts[txID]; !exists {
		return fmt.Errorf("receipt not found")
	}
	delete(rm.receipts, txID)
	return nil
}
