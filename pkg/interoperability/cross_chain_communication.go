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

// NewCrossChainCommunication initializes the cross-chain communication system
func NewCrossChainCommunication(supportedChains []string, validators []common.Validator, ledgerInstance *ledger.Ledger) *CrossChainCommunication {
    return &CrossChainCommunication{
        SupportedChains: supportedChains,
        Validators:      validators,
        LedgerInstance:  ledgerInstance,
        MessagePool:     make(map[string]CrossChainMessage),
    }
}

// SendMessage sends a message from one blockchain to another.
func (cc *CrossChainCommunication) SendMessage(fromChain, toChain, payload string) (string, error) {
    cc.mutex.Lock()
    defer cc.mutex.Unlock()

    // Check if both chains are supported by the cross-chain communication system.
    if !cc.isChainSupported(fromChain) || !cc.isChainSupported(toChain) {
        return "", fmt.Errorf("unsupported blockchain networks")
    }

    // Create an encryption instance and encrypt the payload.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return "", fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedPayload, err := encryptInstance.EncryptData(payload, common.EncryptionKey, nil)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt message payload: %v", err)
    }

    // Convert encryptedPayload to a base64-encoded string.
    encodedPayload := base64.StdEncoding.EncodeToString(encryptedPayload)

    // Generate a unique message ID.
    messageID := cc.generateMessageID(fromChain, toChain)

    // Create the cross-chain message.
    message := CrossChainMessage{
        MessageID: messageID,
        FromChain: fromChain,
        ToChain:   toChain,
        Payload:   encodedPayload,
        Timestamp: time.Now(),
        Status:    "sent",
    }

    // Validate the message before sending.
    validationHash, err := cc.validateMessage(message)
    if err != nil {
        return "", fmt.Errorf("message validation failed: %v", err)
    }

    message.ValidationHash = validationHash

    // Store the message in the pool for future confirmation.
    cc.MessagePool[messageID] = message

    // Log the message sending to the ledger.
    err = cc.logMessageToLedger(message)
    if err != nil {
        return "", fmt.Errorf("failed to log message to ledger: %v", err)
    }

    fmt.Printf("Cross-chain message sent. Message ID: %s\n", messageID)
    return messageID, nil
}

// ConfirmMessage confirms that a message was received and processed by the destination chain
func (cc *CrossChainCommunication) ConfirmMessage(messageID string) error {
    cc.mutex.Lock()
    defer cc.mutex.Unlock()

    message, exists := cc.MessagePool[messageID]
    if !exists {
        return fmt.Errorf("message ID %s not found", messageID)
    }

    // Mark the message as confirmed
    message.Status = "confirmed"

    // Log the confirmation to the ledger
    err := cc.logMessageConfirmationToLedger(message)
    if err != nil {
        return fmt.Errorf("failed to log message confirmation: %v", err)
    }

    // Remove the message from the pool after confirmation
    delete(cc.MessagePool, messageID)

    fmt.Printf("Cross-chain message confirmed. Message ID: %s\n", messageID)
    return nil
}

// validateMessage validates the message across validators
func (cc *CrossChainCommunication) validateMessage(message CrossChainMessage) (string, error) {
    // Select the first validator for simplicity
    validator := cc.Validators[0]

    // Create the validation hash
    validationData := fmt.Sprintf("%s%s%s%d", message.FromChain, message.ToChain, message.Payload, message.Timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(validationData))
    validationHash := hex.EncodeToString(hash.Sum(nil))

    fmt.Printf("Message validated by %s with hash: %s\n", validator.Address, validationHash)
    return validationHash, nil
}



// logMessageToLedger logs the sending of a message to the ledger.
func (cc *CrossChainCommunication) logMessageToLedger(message CrossChainMessage) error {
    // Serialize message data for encryption (optional for audit purposes).
    messageData := fmt.Sprintf("Sent message: %+v", message)

    // Create an encryption instance and encrypt the message data.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    _, err = encryptInstance.EncryptData(messageData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt message data: %v", err)
    }

    // Record only the necessary details in the ledger.
    cc.LedgerInstance.RecordCrossChainMessage(
        message.MessageID,
        message.FromChain,
        message.ToChain,
        message.Payload,
    )

    return nil
}

// logMessageConfirmationToLedger logs the confirmation of a message to the ledger.
func (cc *CrossChainCommunication) logMessageConfirmationToLedger(message CrossChainMessage) error {
    // Serialize confirmation data for encryption (optional for audit purposes).
    confirmationData := fmt.Sprintf("Confirmed message: %+v", message)

    // Create an encryption instance and encrypt the confirmation data.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedData, err := encryptInstance.EncryptData(confirmationData, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt confirmation data: %v", err)
    }

    // Optionally log `encryptedData` for audit purposes if required.
    encodedData := base64.StdEncoding.EncodeToString(encryptedData)
    fmt.Printf("Audit Log Entry for Confirmation: %s\n", encodedData)

    // Record the confirmation of the cross-chain message in the ledger using only the message ID.
    cc.LedgerInstance.RecordCrossChainMessageConfirmation(message.MessageID)

    return nil
}

// generateMessageID generates a unique ID for a message
func (cc *CrossChainCommunication) generateMessageID(fromChain, toChain string) string {
    hashInput := fmt.Sprintf("%s%s%d", fromChain, toChain, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// isChainSupported checks if a blockchain is supported for cross-chain communication
func (cc *CrossChainCommunication) isChainSupported(chain string) bool {
    for _, supportedChain := range cc.SupportedChains {
        if supportedChain == chain {
            return true
        }
    }
    return false
}

// GetMessageStatus returns the status of a cross-chain message
func (cc *CrossChainCommunication) GetMessageStatus(messageID string) (string, error) {
    cc.mutex.Lock()
    defer cc.mutex.Unlock()

    message, exists := cc.MessagePool[messageID]
    if !exists {
        return "", fmt.Errorf("message ID %s not found", messageID)
    }

    return message.Status, nil
}
