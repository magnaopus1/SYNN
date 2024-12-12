// contract_state_management.go

package smart_contract

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// SyncContractState synchronizes the current state of a contract with the ledger.
func SyncContractState(contractID string, state common.ContractState, ledgerInstance *ledger.Ledger) error {
    encryptedState := encryption.EncryptContractState(state)
    if err := ledgerInstance.UpdateContractState(contractID, encryptedState); err != nil {
        return fmt.Errorf("failed to sync state for contract %s: %v", contractID, err)
    }
    fmt.Printf("Contract state synced for %s.\n", contractID)
    return nil
}

// QueryContractVersion retrieves the current version of a contract.
func QueryContractVersion(contractID string, ledgerInstance *ledger.Ledger) (string, error) {
    version, err := ledgerInstance.GetContractVersion(contractID)
    if err != nil {
        return "", fmt.Errorf("failed to query version for contract %s: %v", contractID, err)
    }
    fmt.Printf("Version %s retrieved for contract %s.\n", version, contractID)
    return version, nil
}

// LockSharedResource locks a shared resource for exclusive access by a contract.
func LockSharedResource(contractID, resourceID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LockResource(contractID, resourceID); err != nil {
        return fmt.Errorf("failed to lock resource %s for contract %s: %v", resourceID, contractID, err)
    }
    fmt.Printf("Resource %s locked for contract %s.\n", resourceID, contractID)
    return nil
}

// UnlockSharedResource unlocks a previously locked shared resource.
func UnlockSharedResource(contractID, resourceID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockResource(contractID, resourceID); err != nil {
        return fmt.Errorf("failed to unlock resource %s for contract %s: %v", resourceID, contractID, err)
    }
    fmt.Printf("Resource %s unlocked for contract %s.\n", resourceID, contractID)
    return nil
}

// RegisterInterContractCallback sets up a callback function for inter-contract communication.
func RegisterInterContractCallback(contractID string, callback common.CallbackFunction, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AddContractCallback(contractID, callback); err != nil {
        return fmt.Errorf("failed to register callback for contract %s: %v", contractID, err)
    }
    fmt.Printf("Callback registered for contract %s.\n", contractID)
    return nil
}

// DeregisterInterContractCallback removes a registered callback function for a contract.
func DeregisterInterContractCallback(contractID string, callbackID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RemoveContractCallback(contractID, callbackID); err != nil {
        return fmt.Errorf("failed to deregister callback for contract %s: %v", contractID, err)
    }
    fmt.Printf("Callback deregistered for contract %s.\n", contractID)
    return nil
}

// InitiateAsyncMessage initiates an asynchronous message to another contract.
func InitiateAsyncMessage(senderID, receiverID, message string, ledgerInstance *ledger.Ledger) error {
    encryptedMessage := encryption.EncryptMessage(message)
    if err := ledgerInstance.StoreAsyncMessage(senderID, receiverID, encryptedMessage); err != nil {
        return fmt.Errorf("failed to initiate async message from %s to %s: %v", senderID, receiverID, err)
    }
    fmt.Printf("Async message initiated from %s to %s.\n", senderID, receiverID)
    return nil
}

// AcknowledgeReceipt confirms the receipt of a message.
func AcknowledgeReceipt(contractID, messageID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmMessageReceipt(contractID, messageID); err != nil {
        return fmt.Errorf("failed to acknowledge receipt for message %s on contract %s: %v", messageID, contractID, err)
    }
    fmt.Printf("Receipt acknowledged for message %s on contract %s.\n", messageID, contractID)
    return nil
}

// VerifyMessageIntegrity checks the integrity of a received message using encryption.
func VerifyMessageIntegrity(contractID, messageID string, ledgerInstance *ledger.Ledger) (bool, error) {
    message, err := ledgerInstance.FetchMessage(contractID, messageID)
    if err != nil {
        return false, fmt.Errorf("failed to fetch message %s for integrity check on contract %s: %v", messageID, contractID, err)
    }
    isValid := encryption.VerifyIntegrity(message)
    fmt.Printf("Message integrity for %s on contract %s: %v.\n", messageID, contractID, isValid)
    return isValid, nil
}

// RequestDataConsistencyCheck verifies that the contract's data is consistent with the ledger.
func RequestDataConsistencyCheck(contractID string, ledgerInstance *ledger.Ledger) (bool, error) {
    consistencyStatus, err := ledgerInstance.CheckDataConsistency(contractID)
    if err != nil {
        return false, fmt.Errorf("failed to perform data consistency check for contract %s: %v", contractID, err)
    }
    fmt.Printf("Data consistency for contract %s: %v.\n", contractID, consistencyStatus)
    return consistencyStatus, nil
}
