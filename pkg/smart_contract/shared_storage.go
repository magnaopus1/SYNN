// shared_storage.go

package smart_contract

import (
    "fmt"
    "errors"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/network"
)

// ResolveContractAlias resolves an alias to the full contract ID, helping manage contract identifiers across the network.
func ResolveContractAlias(alias string, ledgerInstance *ledger.Ledger) (string, error) {
    contractID, err := ledgerInstance.GetContractIDByAlias(alias)
    if err != nil {
        return "", fmt.Errorf("failed to resolve alias %s: %v", alias, err)
    }
    fmt.Printf("Alias %s resolved to contract ID %s\n", alias, contractID)
    return contractID, nil
}

// RegisterSharedStorageAccess registers a contract's access to a shared storage.
func RegisterSharedStorageAccess(contractID string, storageKey string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.GrantStorageAccess(contractID, storageKey); err != nil {
        return fmt.Errorf("failed to register shared storage access for contract %s: %v", contractID, err)
    }
    fmt.Printf("Shared storage access granted to contract %s for storage key %s\n", contractID, storageKey)
    return nil
}

// RemoveSharedStorageAccess revokes access to shared storage from a contract.
func RemoveSharedStorageAccess(contractID string, storageKey string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevokeStorageAccess(contractID, storageKey); err != nil {
        return fmt.Errorf("failed to remove shared storage access for contract %s: %v", contractID, err)
    }
    fmt.Printf("Shared storage access revoked for contract %s on storage key %s\n", contractID, storageKey)
    return nil
}

// LogContractMessage records a message sent or received by a contract.
func LogContractMessage(contractID string, message common.Message, ledgerInstance *ledger.Ledger) error {
    encryptedMessage, err := encryption.EncryptMessage(message)
    if err != nil {
        return fmt.Errorf("failed to encrypt message for contract %s: %v", contractID, err)
    }
    if err := ledgerInstance.StoreContractMessage(contractID, encryptedMessage); err != nil {
        return fmt.Errorf("failed to log message for contract %s: %v", contractID, err)
    }
    fmt.Printf("Message logged for contract %s\n", contractID)
    return nil
}

// ValidateContractResponse checks if a response from a contract matches expected criteria.
func ValidateContractResponse(contractID string, response common.Response, ledgerInstance *ledger.Ledger) (bool, error) {
    isValid := ledgerInstance.VerifyResponse(contractID, response)
    if !isValid {
        return false, fmt.Errorf("invalid response received from contract %s", contractID)
    }
    fmt.Printf("Response from contract %s validated\n", contractID)
    return true, nil
}

// RetryFailedCommunication attempts to resend a message that failed previously.
func RetryFailedCommunication(contractID string, messageID string, ledgerInstance *ledger.Ledger) error {
    message, err := ledgerInstance.GetFailedMessage(contractID, messageID)
    if err != nil {
        return fmt.Errorf("failed to retrieve message %s for retry: %v", messageID, err)
    }
    if err := network.ResendMessage(contractID, message); err != nil {
        return fmt.Errorf("failed to retry message %s for contract %s: %v", messageID, contractID, err)
    }
    fmt.Printf("Retry successful for message %s on contract %s\n", messageID, contractID)
    return nil
}

// EnableContractBroadcastMode enables a mode where a contract can broadcast messages to multiple listeners.
func EnableContractBroadcastMode(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetBroadcastMode(contractID, true); err != nil {
        return fmt.Errorf("failed to enable broadcast mode for contract %s: %v", contractID, err)
    }
    fmt.Printf("Broadcast mode enabled for contract %s\n", contractID)
    return nil
}

// DisableContractBroadcastMode disables broadcast mode, limiting communication to direct messages only.
func DisableContractBroadcastMode(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetBroadcastMode(contractID, false); err != nil {
        return fmt.Errorf("failed to disable broadcast mode for contract %s: %v", contractID, err)
    }
    fmt.Printf("Broadcast mode disabled for contract %s\n", contractID)
    return nil
}

// GetActiveContractListeners retrieves a list of active listeners for a given contract.
func GetActiveContractListeners(contractID string, ledgerInstance *ledger.Ledger) ([]string, error) {
    listeners, err := ledgerInstance.FetchActiveListeners(contractID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch active listeners for contract %s: %v", contractID, err)
    }
    fmt.Printf("Active listeners for contract %s: %v\n", contractID, listeners)
    return listeners, nil
}

// SynchronizeSharedData syncs shared data across contracts that have access permissions.
func SynchronizeSharedData(contractID string, data common.SharedData, ledgerInstance *ledger.Ledger) error {
    encryptedData, err := encryption.EncryptData(data)
    if err != nil {
        return fmt.Errorf("failed to encrypt data for synchronization with contract %s: %v", contractID, err)
    }
    if err := ledgerInstance.SyncSharedData(contractID, encryptedData); err != nil {
        return fmt.Errorf("failed to synchronize shared data for contract %s: %v", contractID, err)
    }
    fmt.Printf("Shared data synchronized for contract %s\n", contractID)
    return nil
}
