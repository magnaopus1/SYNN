// monitoring_and_logging.go

package smart_contract

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// ClearContractMessageQueue clears all pending messages in a contract's queue.
func ClearContractMessageQueue(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ClearMessageQueue(contractID); err != nil {
        return fmt.Errorf("failed to clear message queue for contract %s: %v", contractID, err)
    }
    fmt.Printf("Message queue cleared for contract %s.\n", contractID)
    return nil
}

// MonitorContractInteractions logs and tracks all interactions for a contract.
func MonitorContractInteractions(contractID string, ledgerInstance *ledger.Ledger) error {
    interactions, err := ledgerInstance.FetchContractInteractions(contractID)
    if err != nil {
        return fmt.Errorf("failed to monitor interactions for contract %s: %v", contractID, err)
    }
    for _, interaction := range interactions {
        fmt.Printf("Interaction logged for contract %s: %v\n", contractID, interaction)
    }
    return nil
}

// InitiateContractHandshake initiates a handshake between two contracts.
func InitiateContractHandshake(senderID, receiverID string, ledgerInstance *ledger.Ledger) error {
    handshakeData := encryption.GenerateHandshakeData(senderID, receiverID)
    if err := ledgerInstance.RecordHandshake(senderID, receiverID, handshakeData); err != nil {
        return fmt.Errorf("failed to initiate handshake from %s to %s: %v", senderID, receiverID, err)
    }
    fmt.Printf("Handshake initiated from %s to %s.\n", senderID, receiverID)
    return nil
}

// ValidateContractHandshake verifies the authenticity of a handshake.
func ValidateContractHandshake(senderID, receiverID string, ledgerInstance *ledger.Ledger) (bool, error) {
    handshakeData, err := ledgerInstance.GetHandshakeData(senderID, receiverID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve handshake data from %s to %s: %v", senderID, receiverID, err)
    }
    isValid := encryption.ValidateHandshake(handshakeData)
    fmt.Printf("Handshake validation between %s and %s: %v\n", senderID, receiverID, isValid)
    return isValid, nil
}

// DefineCommunicationThreshold sets the threshold for communication retries or timeouts.
func DefineCommunicationThreshold(contractID string, threshold int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetCommunicationThreshold(contractID, threshold); err != nil {
        return fmt.Errorf("failed to set communication threshold for contract %s: %v", contractID, err)
    }
    fmt.Printf("Communication threshold set for contract %s.\n", contractID)
    return nil
}

// RegisterInteractionMonitor starts monitoring a contract's interactions.
func RegisterInteractionMonitor(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.EnableInteractionMonitor(contractID); err != nil {
        return fmt.Errorf("failed to register interaction monitor for contract %s: %v", contractID, err)
    }
    fmt.Printf("Interaction monitor registered for contract %s.\n", contractID)
    return nil
}

// UnregisterInteractionMonitor stops monitoring a contract's interactions.
func UnregisterInteractionMonitor(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.DisableInteractionMonitor(contractID); err != nil {
        return fmt.Errorf("failed to unregister interaction monitor for contract %s: %v", contractID, err)
    }
    fmt.Printf("Interaction monitor unregistered for contract %s.\n", contractID)
    return nil
}

// TrackSharedStorageUsage logs and tracks usage of shared storage by a contract.
func TrackSharedStorageUsage(contractID string, storageUsed int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateSharedStorageUsage(contractID, storageUsed); err != nil {
        return fmt.Errorf("failed to track shared storage usage for contract %s: %v", contractID, err)
    }
    fmt.Printf("Shared storage usage tracked for contract %s.\n", contractID)
    return nil
}

// GetLastMessageTimestamp retrieves the timestamp of the last message sent or received by a contract.
func GetLastMessageTimestamp(contractID string, ledgerInstance *ledger.Ledger) (time.Time, error) {
    timestamp, err := ledgerInstance.FetchLastMessageTimestamp(contractID)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to get last message timestamp for contract %s: %v", contractID, err)
    }
    fmt.Printf("Last message timestamp for contract %s: %v\n", contractID, timestamp)
    return timestamp, nil
}

// CheckContractDependencyStatus verifies the operational status of dependent contracts.
func CheckContractDependencyStatus(contractID string, dependencies []string, ledgerInstance *ledger.Ledger) (map[string]bool, error) {
    statusMap := make(map[string]bool)
    for _, depID := range dependencies {
        status, err := ledgerInstance.CheckContractStatus(depID)
        if err != nil {
            return nil, fmt.Errorf("failed to check status of dependency %s for contract %s: %v", depID, contractID, err)
        }
        statusMap[depID] = status
    }
    fmt.Printf("Dependency status for contract %s: %v\n", contractID, statusMap)
    return statusMap, nil
}

// RetrieveMessageHistory fetches the message history for a contract.
func RetrieveMessageHistory(contractID string, ledgerInstance *ledger.Ledger) ([]common.MessageLog, error) {
    messageHistory, err := ledgerInstance.FetchMessageHistory(contractID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve message history for contract %s: %v", contractID, err)
    }
    fmt.Printf("Message history retrieved for contract %s.\n", contractID)
    return messageHistory, nil
}

// ConfigureMessageRetryLimit sets the retry limit for messages sent by the contract.
func ConfigureMessageRetryLimit(contractID string, retryLimit int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetRetryLimit(contractID, retryLimit); err != nil {
        return fmt.Errorf("failed to set retry limit for contract %s: %v", contractID, err)
    }
    fmt.Printf("Retry limit set for contract %s.\n", contractID)
    return nil
}

// SetInterContractAuthorization grants permission for inter-contract communication.
func SetInterContractAuthorization(contractID, targetID string, permissionLevel string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.GrantAuthorization(contractID, targetID, permissionLevel); err != nil {
        return fmt.Errorf("failed to set inter-contract authorization from %s to %s: %v", contractID, targetID, err)
    }
    fmt.Printf("Inter-contract authorization set from %s to %s.\n", contractID, targetID)
    return nil
}

// RemoveInterContractAuthorization revokes authorization for inter-contract communication.
func RemoveInterContractAuthorization(contractID, targetID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevokeAuthorization(contractID, targetID); err != nil {
        return fmt.Errorf("failed to remove inter-contract authorization from %s to %s: %v", contractID, targetID, err)
    }
    fmt.Printf("Inter-contract authorization removed from %s to %s.\n", contractID, targetID)
    return nil
}
