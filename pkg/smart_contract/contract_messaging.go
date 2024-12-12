// contract_messaging.go

package smart_contract

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// SendMessage securely sends a message to another contract.
func SendMessage(senderID, receiverID, message string, ledgerInstance *ledger.Ledger) error {
    encryptedMessage := encryption.EncryptMessage(message)
    if err := ledgerInstance.StoreMessage(senderID, receiverID, encryptedMessage); err != nil {
        return fmt.Errorf("failed to send message from %s to %s: %v", senderID, receiverID, err)
    }
    fmt.Printf("Message sent from %s to %s.\n", senderID, receiverID)
    return nil
}

// InvokeFunction calls a function on another contract, passing encrypted parameters.
func InvokeFunction(contractID, function string, params common.ExecutionParams, ledgerInstance *ledger.Ledger) (common.ExecutionResult, error) {
    encryptedParams := encryption.EncryptExecutionParams(params)
    result, err := ledgerInstance.ExecuteContractFunction(contractID, function, encryptedParams)
    if err != nil {
        return common.ExecutionResult{}, fmt.Errorf("failed to invoke function %s on contract %s: %v", function, contractID, err)
    }
    fmt.Printf("Function %s invoked on contract %s.\n", function, contractID)
    return result, nil
}

// TransferData moves encrypted data between contracts.
func TransferData(senderID, receiverID string, data common.ContractData, ledgerInstance *ledger.Ledger) error {
    encryptedData := encryption.EncryptContractData(data)
    if err := ledgerInstance.TransferContractData(senderID, receiverID, encryptedData); err != nil {
        return fmt.Errorf("failed to transfer data from %s to %s: %v", senderID, receiverID, err)
    }
    fmt.Printf("Data transferred from %s to %s.\n", senderID, receiverID)
    return nil
}

// RequestContractData fetches data from another contract securely.
func RequestContractData(contractID string, dataID string, ledgerInstance *ledger.Ledger) (common.ContractData, error) {
    encryptedData, err := ledgerInstance.FetchContractData(contractID, dataID)
    if err != nil {
        return common.ContractData{}, fmt.Errorf("failed to request data %s from contract %s: %v", dataID, contractID, err)
    }
    data := encryption.DecryptContractData(encryptedData)
    fmt.Printf("Data %s retrieved from contract %s.\n", dataID, contractID)
    return data, nil
}

// EventSubscribe subscribes to a specific event on a contract.
func EventSubscribe(contractID, eventName string, subscriberID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AddEventSubscriber(contractID, eventName, subscriberID); err != nil {
        return fmt.Errorf("failed to subscribe %s to event %s on contract %s: %v", subscriberID, eventName, contractID, err)
    }
    fmt.Printf("Subscriber %s added to event %s on contract %s.\n", subscriberID, eventName, contractID)
    return nil
}

// EventUnsubscribe removes a subscription to a specific event on a contract.
func EventUnsubscribe(contractID, eventName string, subscriberID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RemoveEventSubscriber(contractID, eventName, subscriberID); err != nil {
        return fmt.Errorf("failed to unsubscribe %s from event %s on contract %s: %v", subscriberID, eventName, contractID, err)
    }
    fmt.Printf("Subscriber %s removed from event %s on contract %s.\n", subscriberID, eventName, contractID)
    return nil
}

// BroadcastEvent sends a notification of an event to all subscribers.
func BroadcastEvent(contractID, eventName string, eventData common.EventData, ledgerInstance *ledger.Ledger) error {
    encryptedEventData := encryption.EncryptEventData(eventData)
    if err := ledgerInstance.NotifyEventSubscribers(contractID, eventName, encryptedEventData); err != nil {
        return fmt.Errorf("failed to broadcast event %s on contract %s: %v", eventName, contractID, err)
    }
    fmt.Printf("Event %s broadcasted on contract %s.\n", eventName, contractID)
    return nil
}

// PollContractStatus checks and returns the current status of a contract.
func PollContractStatus(contractID string, ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.FetchContractStatus(contractID)
    if err != nil {
        return "", fmt.Errorf("failed to poll status for contract %s: %v", contractID, err)
    }
    fmt.Printf("Status of contract %s: %s.\n", contractID, status)
    return status, nil
}

// SetSharedStorage configures a shared storage value for a contract.
func SetSharedStorage(contractID string, storageKey string, storageValue common.SharedStorage, ledgerInstance *ledger.Ledger) error {
    encryptedValue := encryption.EncryptSharedStorage(storageValue)
    if err := ledgerInstance.SetSharedStorage(contractID, storageKey, encryptedValue); err != nil {
        return fmt.Errorf("failed to set shared storage for %s on contract %s: %v", storageKey, contractID, err)
    }
    fmt.Printf("Shared storage %s set on contract %s.\n", storageKey, contractID)
    return nil
}

// GetSharedStorage retrieves a shared storage value for a contract.
func GetSharedStorage(contractID string, storageKey string, ledgerInstance *ledger.Ledger) (common.SharedStorage, error) {
    encryptedValue, err := ledgerInstance.FetchSharedStorage(contractID, storageKey)
    if err != nil {
        return common.SharedStorage{}, fmt.Errorf("failed to retrieve shared storage %s from contract %s: %v", storageKey, contractID, err)
    }
    storageValue := encryption.DecryptSharedStorage(encryptedValue)
    fmt.Printf("Shared storage %s retrieved from contract %s.\n", storageKey, contractID)
    return storageValue, nil
}
