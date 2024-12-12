// contract_interactions.go

package smart_contract

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
    "time"
)

// RequestEventFeedback sends a feedback request for an event triggered by a contract.
func RequestEventFeedback(contractID, eventID string, ledgerInstance *ledger.Ledger) (common.EventFeedback, error) {
    feedback, err := ledgerInstance.FetchEventFeedback(contractID, eventID)
    if err != nil {
        return common.EventFeedback{}, fmt.Errorf("failed to request feedback for event %s in contract %s: %v", eventID, contractID, err)
    }
    fmt.Printf("Feedback for event %s in contract %s retrieved.\n", eventID, contractID)
    return feedback, nil
}

// ConfigureCommunicationRetryPolicy sets the retry policy for contract communications.
func ConfigureCommunicationRetryPolicy(contractID string, retryPolicy common.RetryPolicy, ledgerInstance *ledger.Ledger) error {
    encryptedPolicy := encryption.EncryptRetryPolicy(retryPolicy)
    if err := ledgerInstance.StoreRetryPolicy(contractID, encryptedPolicy); err != nil {
        return fmt.Errorf("failed to configure retry policy for contract %s: %v", contractID, err)
    }
    fmt.Printf("Retry policy configured for contract %s.\n", contractID)
    return nil
}

// ValidateSharedResourceState verifies the current state of a shared resource among contracts.
func ValidateSharedResourceState(resourceID string, ledgerInstance *ledger.Ledger) (common.ResourceState, error) {
    state, err := ledgerInstance.FetchSharedResourceState(resourceID)
    if err != nil {
        return common.ResourceState{}, fmt.Errorf("failed to validate state for shared resource %s: %v", resourceID, err)
    }
    fmt.Printf("State for shared resource %s validated.\n", resourceID)
    return state, nil
}

// InitiateCrossContractCall initializes a call to another contract, handling all encryption and state management.
func InitiateCrossContractCall(originContractID, targetContractID, function string, params common.ExecutionParams, ledgerInstance *ledger.Ledger) (common.ExecutionResult, error) {
    encryptedParams := encryption.EncryptExecutionParams(params)
    result, err := ledgerInstance.ExecuteCrossContractFunction(originContractID, targetContractID, function, encryptedParams)
    if err != nil {
        return common.ExecutionResult{}, fmt.Errorf("cross-contract call from %s to %s failed for function %s: %v", originContractID, targetContractID, function, err)
    }
    fmt.Printf("Cross-contract call from %s to %s executed for function %s.\n", originContractID, targetContractID, function)
    return result, nil
}

// RevokeSharedResourceAccess revokes access to a shared resource for a specific contract.
func RevokeSharedResourceAccess(contractID, resourceID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevokeResourceAccess(contractID, resourceID); err != nil {
        return fmt.Errorf("failed to revoke access for resource %s in contract %s: %v", resourceID, contractID, err)
    }
    fmt.Printf("Access to resource %s revoked for contract %s.\n", resourceID, contractID)
    return nil
}

// QueryContractConnectionStatus retrieves the current connection status of a contract.
func QueryContractConnectionStatus(contractID string, ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.FetchContractConnectionStatus(contractID)
    if err != nil {
        return "", fmt.Errorf("failed to fetch connection status for contract %s: %v", contractID, err)
    }
    fmt.Printf("Connection status for contract %s: %s.\n", contractID, status)
    return status, nil
}

// LogSharedResourceAccess logs access to a shared resource by a contract.
func LogSharedResourceAccess(contractID, resourceID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LogResourceAccess(contractID, resourceID); err != nil {
        return fmt.Errorf("failed to log access for resource %s in contract %s: %v", resourceID, contractID, err)
    }
    fmt.Printf("Access to resource %s logged for contract %s.\n", resourceID, contractID)
    return nil
}

// EnableEncryptedCommunication enables encryption for communications involving a contract.
func EnableEncryptedCommunication(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateCommunicationEncryption(contractID, true); err != nil {
        return fmt.Errorf("failed to enable encrypted communication for contract %s: %v", contractID, err)
    }
    fmt.Printf("Encrypted communication enabled for contract %s.\n", contractID)
    return nil
}

// DisableEncryptedCommunication disables encryption for communications involving a contract.
func DisableEncryptedCommunication(contractID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateCommunicationEncryption(contractID, false); err != nil {
        return fmt.Errorf("failed to disable encrypted communication for contract %s: %v", contractID, err)
    }
    fmt.Printf("Encrypted communication disabled for contract %s.\n", contractID)
    return nil
}

// SetContractCommTimeout sets a communication timeout for a contract.
func SetContractCommTimeout(contractID string, timeout time.Duration, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateCommunicationTimeout(contractID, timeout); err != nil {
        return fmt.Errorf("failed to set communication timeout for contract %s: %v", contractID, err)
    }
    fmt.Printf("Communication timeout of %s set for contract %s.\n", timeout, contractID)
    return nil
}
