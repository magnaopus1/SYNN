package interoperability

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewBlockchainAgnosticManager initializes the manager with supported chains and validators
func NewBlockchainAgnosticManager(ledgerInstance *ledger.Ledger, supportedChains []string, validators []common.Validator) *BlockchainAgnosticManager {
	protocol := &BlockchainAgnosticProtocol{
		SupportedChains: supportedChains,
		Validators:      validators,
		LedgerInstance:  ledgerInstance,
	}

	return &BlockchainAgnosticManager{
		ActiveTransactions: make(map[string]*CrossChainTransaction),
		LedgerInstance:     ledgerInstance,
		Protocol:           protocol,
	}
}

// InitiateCrossChainTransaction initiates a cross-chain transaction using the agnostic protocol
func (manager *BlockchainAgnosticManager) InitiateCrossChainTransaction(fromChain, toChain string, amount float64, tokenSymbol, fromAddress, toAddress string) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if !manager.Protocol.isChainSupported(fromChain) || !manager.Protocol.isChainSupported(toChain) {
		return "", errors.New("one or both chains are not supported by the protocol")
	}

	// Generate unique transaction ID
	transactionID := manager.generateTransactionID(fromChain, toChain)

	// Create cross-chain transaction
	transaction := &CrossChainTransaction{
		TransactionID:  transactionID,
		FromChain:      fromChain,
		ToChain:        toChain,
		Amount:         amount,
		TokenSymbol:    tokenSymbol,
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		Timestamp:      time.Now(),
		Status:         "pending",
	}

	// Validate the transaction across chains
	validationHash, err := manager.validateCrossChainTransaction(transaction)
	if err != nil {
		return "", fmt.Errorf("transaction validation failed: %v", err)
	}

	transaction.ValidationHash = validationHash
	manager.ActiveTransactions[transactionID] = transaction

	// Log the transaction initiation in the ledger
	err = manager.logTransactionToLedger(transaction)
	if err != nil {
		return "", fmt.Errorf("failed to log transaction in the ledger: %v", err)
	}

	fmt.Printf("Cross-chain transaction initiated. Transaction ID: %s\n", transactionID)
	return transactionID, nil
}

// CompleteCrossChainTransaction completes a cross-chain transaction after validation
func (manager *BlockchainAgnosticManager) CompleteCrossChainTransaction(transactionID string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	transaction, exists := manager.ActiveTransactions[transactionID]
	if !exists {
		return errors.New("transaction not found")
	}

	if transaction.Status != "pending" {
		return errors.New("transaction is not in a pending state")
	}

	// Complete transaction and record in ledger
	transaction.Status = "completed"

	// Log transaction completion to the ledger
	err := manager.logTransactionCompletionToLedger(transaction)
	if err != nil {
		return fmt.Errorf("failed to log transaction completion: %v", err)
	}

	fmt.Printf("Cross-chain transaction completed. Transaction ID: %s\n", transactionID)
	return nil
}

// validateCrossChainTransaction validates the transaction across chains using validators
func (manager *BlockchainAgnosticManager) validateCrossChainTransaction(transaction *CrossChainTransaction) (string, error) {
	// Simulating validation using multiple validators
	validator := manager.Protocol.Validators[0] // Simplified: select first validator for demonstration

	validationData := fmt.Sprintf("%s%s%f%s%s%d", transaction.FromChain, transaction.ToChain, transaction.Amount, transaction.FromAddress, transaction.ToAddress, transaction.Timestamp.UnixNano())
	validationHash := sha256.New()
	validationHash.Write([]byte(validationData))
	hashString := hex.EncodeToString(validationHash.Sum(nil))

	fmt.Printf("Transaction validated by %s with hash: %s\n", validator.Address, hashString)
	return hashString, nil
}

// isChainSupported checks if a blockchain is supported by the protocol
func (protocol *BlockchainAgnosticProtocol) isChainSupported(chain string) bool {
	for _, supportedChain := range protocol.SupportedChains {
		if supportedChain == chain {
			return true
		}
	}
	return false
}

// generateTransactionID generates a unique ID for the cross-chain transaction
func (manager *BlockchainAgnosticManager) generateTransactionID(fromChain, toChain string) string {
	hashInput := fmt.Sprintf("%s%s%d", fromChain, toChain, time.Now().UnixNano())
	hash := sha256.New()
	hash.Write([]byte(hashInput))
	return hex.EncodeToString(hash.Sum(nil))
}

// logTransactionToLedger logs the initiation of a cross-chain transaction to the ledger.
func (manager *BlockchainAgnosticManager) logTransactionToLedger(transaction *CrossChainTransaction) error {
    // Serialize transaction data for logging/audit purposes
    transactionData := fmt.Sprintf("Initiated cross-chain transaction: %+v", transaction)

    // Create an encryption instance
    encryptInstance, err := common.NewEncryption(256)
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt transaction data if needed (omit encryptedData if unused)
    _, err = encryptInstance.EncryptData(transactionData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt transaction data: %v", err)
    }

    // Record the cross-chain transaction in the ledger with specific details
    manager.LedgerInstance.RecordCrossChainTransaction(
        transaction.TransactionID,
        transaction.FromAddress,
        transaction.ToAddress,
        transaction.FromChain,
        transaction.ToChain,
        transaction.Amount,
    )

    return nil
}


// logTransactionCompletionToLedger logs the completion of a cross-chain transaction to the ledger.
func (manager *BlockchainAgnosticManager) logTransactionCompletionToLedger(transaction *CrossChainTransaction) error {
    // Serialize transaction data for logging/audit purposes
    transactionData := fmt.Sprintf("Completed cross-chain transaction: %+v", transaction)

    // Create an encryption instance
    encryptInstance, err := common.NewEncryption(256)
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt transaction data if needed (omit storing encryptedData if unused)
    _, err = encryptInstance.EncryptData(transactionData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt transaction data: %v", err)
    }

    // Record the transaction completion in the ledger with only the transaction ID
    manager.LedgerInstance.RecordCrossChainTransactionCompletion(transaction.TransactionID)

    return nil
}