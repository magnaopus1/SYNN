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

// NewBridge initializes the bridge with supported chains, validators, and an initial balance
func NewBridge(supportedChains []string, validators []common.Validator, ledgerInstance *ledger.Ledger) *Bridge {
    return &Bridge{
        SupportedChains: supportedChains,
        Validators:      validators,
        LedgerInstance:  ledgerInstance,
        BridgeBalance:   make(map[string]float64),
    }
}

// InitiateTransfer initializes a cross-chain transfer using the bridge
func (b *Bridge) InitiateTransfer(fromChain, toChain string, amount float64, tokenSymbol, fromAddress, toAddress string) (string, error) {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    // Check if both chains are supported by the bridge
    if !b.isChainSupported(fromChain) || !b.isChainSupported(toChain) {
        return "", errors.New("unsupported blockchain networks")
    }

    // Check if the bridge has sufficient balance for the token
    if b.BridgeBalance[tokenSymbol] < amount {
        return "", fmt.Errorf("insufficient bridge balance for token %s", tokenSymbol)
    }

    // Generate a unique transfer ID
    transferID := b.generateTransferID(fromChain, toChain)

    // Create a cross-chain transfer
    transfer := &CrossChainTransfer{
        TransferID:     transferID,
        FromChain:      fromChain,
        ToChain:        toChain,
        Amount:         amount,
        TokenSymbol:    tokenSymbol,
        FromAddress:    fromAddress,
        ToAddress:      toAddress,
        Timestamp:      time.Now(),
        Status:         "pending",
    }

    // Validate the transfer
    validationHash, err := b.validateTransfer(transfer)
    if err != nil {
        return "", fmt.Errorf("transfer validation failed: %v", err)
    }

    transfer.ValidationHash = validationHash

    // Deduct balance from the bridge for the token
    b.BridgeBalance[tokenSymbol] -= amount

    // Log the transfer initiation in the ledger
    err = b.logTransferToLedger(transfer)
    if err != nil {
        return "", fmt.Errorf("failed to log transfer in the ledger: %v", err)
    }

    fmt.Printf("Cross-chain transfer initiated. Transfer ID: %s\n", transferID)
    return transferID, nil
}

// CompleteTransfer completes the cross-chain transfer after validation
func (b *Bridge) CompleteTransfer(transferID string) error {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    transfer, err := b.getTransferByID(transferID)
    if err != nil {
        return err
    }

    if transfer.Status != "pending" {
        return fmt.Errorf("transfer is not in a pending state")
    }

    // Mark the transfer as completed
    transfer.Status = "completed"

    // Log transfer completion to the ledger
    err = b.logTransferCompletionToLedger(transfer)
    if err != nil {
        return fmt.Errorf("failed to log transfer completion: %v", err)
    }

    fmt.Printf("Cross-chain transfer completed. Transfer ID: %s\n", transferID)
    return nil
}

// AddFundsToBridge allows adding tokens to the bridge balance
func (b *Bridge) AddFundsToBridge(tokenSymbol string, amount float64) {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    b.BridgeBalance[tokenSymbol] += amount
    fmt.Printf("Added %.2f %s to bridge balance. New balance: %.2f %s\n", amount, tokenSymbol, b.BridgeBalance[tokenSymbol], tokenSymbol)
}

// validateTransfer validates the transfer across validators
func (b *Bridge) validateTransfer(transfer *CrossChainTransfer) (string, error) {
    // Select first validator (simplified for demonstration)
    validator := b.Validators[0]

    // Generate validation hash
    validationData := fmt.Sprintf("%s%s%f%s%s%d", transfer.FromChain, transfer.ToChain, transfer.Amount, transfer.FromAddress, transfer.ToAddress, transfer.Timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(validationData))
    validationHash := hex.EncodeToString(hash.Sum(nil))

    fmt.Printf("Transfer validated by %s with hash: %s\n", validator.Address, validationHash)
    return validationHash, nil
}

// generateTransferID creates a unique ID for the cross-chain transfer
func (b *Bridge) generateTransferID(fromChain, toChain string) string {
    hashInput := fmt.Sprintf("%s%s%d", fromChain, toChain, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// logTransferToLedger logs the initiation of a transfer to the ledger.
func (b *Bridge) logTransferToLedger(transfer *CrossChainTransfer) error {
    // Serialize transfer data for logging/audit purposes
    transferData := fmt.Sprintf("Initiated transfer: %+v", transfer)

    // Create an encryption instance
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt transfer data without using the result if not needed for auditing
    _, err = encryptInstance.EncryptData(transferData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt transfer data: %v", err)
    }

    // Record the cross-chain transfer in the ledger with specific details
    b.LedgerInstance.RecordCrossChainTransfer(
        transfer.TransferID,
        transfer.TokenSymbol,   // Map TokenSymbol to asset
        transfer.FromChain,     // Map FromChain to sourceChainID
        transfer.ToChain,       // Map ToChain to targetChainID
        transfer.Amount,
    )

    return nil
}

// logTransferCompletionToLedger logs the completion of a transfer to the ledger.
func (b *Bridge) logTransferCompletionToLedger(transfer *CrossChainTransfer) error {
    // Serialize transfer data for logging/audit purposes
    transferData := fmt.Sprintf("Completed transfer: %+v", transfer)

    // Create an encryption instance
    encryptInstance, err := common.NewEncryption(256)
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt transfer data without using the result if it’s not needed
    _, err = encryptInstance.EncryptData(transferData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt transfer data: %v", err)
    }

    // Record the completion of the cross-chain transfer using only the transfer ID
    b.LedgerInstance.RecordCrossChainTransferCompletion(transfer.TransferID)

    return nil
}

// isChainSupported checks if a blockchain is supported by the bridge
func (b *Bridge) isChainSupported(chain string) bool {
    for _, supportedChain := range b.SupportedChains {
        if supportedChain == chain {
            return true
        }
    }
    return false
}

// getTransferByID retrieves a transfer by its ID
func (b *Bridge) getTransferByID(transferID string) (*CrossChainTransfer, error) {
    // In practice, this would look up the transfer in a data store.
    // Placeholder function for demo purposes.
    return nil, fmt.Errorf("transfer ID %s not found", transferID)
}
