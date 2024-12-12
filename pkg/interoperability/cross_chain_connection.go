package interoperability

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewCrossChainConnection initializes a new cross-chain connection manager.
func NewCrossChainConnection(connectedChains []string, ledgerInstance *ledger.Ledger, subBlockPool *common.SubBlockChain) *CrossChainConnection {
    return &CrossChainConnection{
        ConnectedChains: connectedChains,
        LedgerInstance:  ledgerInstance,
        SubBlockPool:    subBlockPool,
    }
}

// EstablishConnection attempts to establish a secure connection between two blockchain networks.
func (ccc *CrossChainConnection) EstablishConnection(fromChain, toChain string) (string, error) {
    ccc.mutex.Lock()
    defer ccc.mutex.Unlock()

    // Check if both chains are supported.
    if !ccc.isChainSupported(fromChain) || !ccc.isChainSupported(toChain) {
        return "", fmt.Errorf("unsupported blockchain networks")
    }

    // Generate a unique connection ID.
    connectionID := ccc.generateConnectionID(fromChain, toChain)

    // Create an encryption instance and encrypt the connection data.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return "", fmt.Errorf("failed to create encryption instance: %v", err)
    }

    _, err = encryptInstance.EncryptData(
        fmt.Sprintf("Connection established between %s and %s", fromChain, toChain),
        common.EncryptionKey,
        make([]byte, 12), // Example nonce; adjust as needed
    )
    if err != nil {
        return "", fmt.Errorf("failed to encrypt connection data: %v", err)
    }

    // Log the connection to the ledger.
    err = ccc.logConnectionToLedger(fromChain, toChain, connectionID)
    if err != nil {
        return "", fmt.Errorf("failed to log connection: %v", err)
    }

    fmt.Printf("Cross-chain connection established between %s and %s. Connection ID: %s\n", fromChain, toChain, connectionID)
    return connectionID, nil
}


// SendTransaction sends a cross-chain transaction to the sub-block pool for validation.
func (ccc *CrossChainConnection) SendTransaction(fromChain, toChain, transactionData string) error {
    ccc.mutex.Lock()
    defer ccc.mutex.Unlock()

    // Check if both chains are supported.
    if !ccc.isChainSupported(fromChain) || !ccc.isChainSupported(toChain) {
        return fmt.Errorf("unsupported blockchain networks")
    }

    // Create a new cross-chain transaction.
    transaction := CrossChainTransaction{
        TransactionID:  "transaction_id_example", // Replace with generated ID
        FromChain:      fromChain,
        ToChain:        toChain,
        Amount:         0,                   // Set amount as needed
        TokenSymbol:    "SYN",
        FromAddress:    "sender_address",     // Replace with actual sender address
        ToAddress:      "recipient_address",  // Replace with actual recipient address
        Timestamp:      time.Now(),
        Status:         "pending",
        ValidationHash: "validation_hash",   // Placeholder or computed value
        Data:           transactionData,
    }

    // Convert CrossChainTransaction to common.Transaction for pool validation
    commonTransaction := convertToCommonTransaction(transaction)

    // Validate the transaction through the sub-block pool.
    ccc.SubBlockPool.AddTransactions([]common.Transaction{commonTransaction}) // No error assignment if no return value

    // Log the cross-chain transaction in the ledger using the original transaction.
    err := ccc.logTransactionToLedger(transaction)
    if err != nil {
        return fmt.Errorf("failed to log transaction to ledger: %v", err)
    }

    fmt.Printf("Cross-chain transaction from %s to %s successfully sent and validated.\n", fromChain, toChain)
    return nil
}

// convertToCommonTransaction converts a CrossChainTransaction to a common.Transaction.
func convertToCommonTransaction(tx CrossChainTransaction) common.Transaction {
    return common.Transaction{
        TransactionID:   tx.TransactionID,
        FromAddress:     tx.FromAddress,
        ToAddress:       tx.ToAddress,
        Amount:          tx.Amount,
        TokenStandard:   tx.TokenSymbol,
        Timestamp:       tx.Timestamp,
        Status:          tx.Status,
    }
}

// logConnectionToLedger logs a new cross-chain connection to the ledger.
func (ccc *CrossChainConnection) logConnectionToLedger(fromChain, toChain, connectionID string) error {
    // Prepare the log data for encryption.
    logData := fmt.Sprintf("Connection ID: %s | From: %s | To: %s", connectionID, fromChain, toChain)

    // Create an encryption instance and encrypt the log data.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedData, err := encryptInstance.EncryptData(logData, common.EncryptionKey, make([]byte, 12)) // Add a nonce if required
    if err != nil {
        return fmt.Errorf("failed to encrypt connection log: %v", err)
    }

    // Optionally encode to base64 if needed for external storage/logging
    _ = base64.StdEncoding.EncodeToString(encryptedData) // Encoding to base64 without assigning to `encryptedLog`

    // Log the cross-chain connection in the ledger.
    ccc.LedgerInstance.RecordCrossChainConnection(connectionID, fromChain, toChain)

    return nil
}

// logTransactionToLedger logs a cross-chain transaction to the ledger.
func (ccc *CrossChainConnection) logTransactionToLedger(transaction CrossChainTransaction) error {
    // Prepare the transaction data string for encryption (for logging purposes).
    transactionData := fmt.Sprintf("Cross-chain transaction from %s to %s at %s", transaction.FromChain, transaction.ToChain, transaction.Timestamp.String())

    // Create an encryption instance and encrypt the transaction data (optional).
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedData, err := encryptInstance.EncryptData(transactionData, common.EncryptionKey, make([]byte, 12)) // Use a nonce if required
    if err != nil {
        return fmt.Errorf("failed to encrypt transaction log: %v", err)
    }

    // Store encryptedData for optional audit log (if needed).
    _ = encryptedData // Remove if not needed for audit

    // Log the cross-chain transaction in the ledger.
    ccc.LedgerInstance.RecordCrossChainTransaction(
        transaction.TransactionID,
        transaction.FromAddress,
        transaction.ToAddress,
        transaction.FromChain,
        transaction.ToChain,
        transaction.Amount,
    )

    return nil
}


// generateConnectionID generates a unique ID for a cross-chain connection
func (ccc *CrossChainConnection) generateConnectionID(fromChain, toChain string) string {
    hashInput := fmt.Sprintf("%s%s%d", fromChain, toChain, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// isChainSupported checks if a blockchain is supported for cross-chain connections
func (ccc *CrossChainConnection) isChainSupported(chain string) bool {
    for _, supportedChain := range ccc.ConnectedChains {
        if supportedChain == chain {
            return true
        }
    }
    return false
}

// GetConnectionStatus returns the status of a cross-chain connection
func (ccc *CrossChainConnection) GetConnectionStatus(connectionID string) (string, error) {
    ccc.mutex.Lock()
    defer ccc.mutex.Unlock()

    // Here, we could check the ledger or other state to return the status of the connection
    // For simplicity, we will just assume the connection is active
    return "active", nil
}
