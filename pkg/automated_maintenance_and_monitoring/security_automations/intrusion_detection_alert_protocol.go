package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    TransactionValidationInterval = 15 * time.Second // Interval for transaction validation checks
    MaxValidationRetries          = 5                // Maximum retries for transaction validation
    SubBlocksPerBlock             = 1000             // Number of sub-blocks in a block
)

// TransactionValidationAutomation handles automated validation of transactions into sub-blocks and blocks
type TransactionValidationAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging transaction validation events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    validationRetryCount map[string]int               // Counter for retrying transaction validation
    validationCycleCount int                          // Counter for validation cycles
}

// NewTransactionValidationAutomation initializes the automation for transaction validation
func NewTransactionValidationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *TransactionValidationAutomation {
    return &TransactionValidationAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        validationRetryCount: make(map[string]int),
        validationCycleCount: 0,
    }
}

// StartTransactionValidation starts the continuous loop for validating transactions into sub-blocks and blocks
func (automation *TransactionValidationAutomation) StartTransactionValidation() {
    ticker := time.NewTicker(TransactionValidationInterval)

    go func() {
        for range ticker.C {
            automation.validateAndProcessTransactions()
        }
    }()
}

// validateAndProcessTransactions validates pending transactions and organizes them into sub-blocks and blocks
func (automation *TransactionValidationAutomation) validateAndProcessTransactions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch pending transactions from the consensus system
    pendingTransactions := automation.consensusSystem.FetchPendingTransactions()

    if len(pendingTransactions) > 0 {
        for _, tx := range pendingTransactions {
            fmt.Printf("Validating transaction ID: %s\n", tx.ID)
            automation.validateTransaction(tx)
        }
    } else {
        fmt.Println("No pending transactions detected this cycle.")
    }

    automation.validationCycleCount++
    fmt.Printf("Transaction validation cycle #%d completed.\n", automation.validationCycleCount)

    if automation.validationCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeBlockCycle()
    }
}

// validateTransaction attempts to validate a transaction and add it to a sub-block
func (automation *TransactionValidationAutomation) validateTransaction(tx common.Transaction) {
    encryptedTx := automation.encryptTransactionData(tx)

    // Attempt to validate and add the transaction to a sub-block through the consensus system
    validationSuccess := automation.consensusSystem.ValidateTransaction(tx, encryptedTx)

    if validationSuccess {
        fmt.Printf("Transaction ID: %s validated successfully.\n", tx.ID)
        automation.logValidationEvent(tx, "Validated")
        automation.resetValidationRetry(tx.ID)
    } else {
        fmt.Printf("Error validating transaction ID: %s. Retrying...\n", tx.ID)
        automation.retryTransactionValidation(tx)
    }
}

// retryTransactionValidation retries the transaction validation if the first attempt fails
func (automation *TransactionValidationAutomation) retryTransactionValidation(tx common.Transaction) {
    automation.validationRetryCount[tx.ID]++
    if automation.validationRetryCount[tx.ID] < MaxValidationRetries {
        automation.validateTransaction(tx)
    } else {
        fmt.Printf("Max retries reached for validating transaction ID: %s. Validation failed.\n", tx.ID)
        automation.logValidationFailure(tx)
    }
}

// resetValidationRetry resets the retry count for a transaction validation
func (automation *TransactionValidationAutomation) resetValidationRetry(txID string) {
    automation.validationRetryCount[txID] = 0
}

// finalizeBlockCycle finalizes the validation of sub-blocks into a block and logs the result in the ledger
func (automation *TransactionValidationAutomation) finalizeBlockCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeBlockCycle()
    if success {
        fmt.Println("Block finalized successfully.")
        automation.logBlockFinalization()
    } else {
        fmt.Println("Error finalizing block cycle.")
    }
}

// logValidationEvent logs a successful transaction validation into the ledger
func (automation *TransactionValidationAutomation) logValidationEvent(tx common.Transaction, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("transaction-validation-%s-%s", tx.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Validation Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Transaction ID: %s was %s.", tx.ID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with validation event for transaction ID: %s.\n", tx.ID)
}

// logValidationFailure logs the failure of a transaction validation into the ledger
func (automation *TransactionValidationAutomation) logValidationFailure(tx common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("transaction-validation-failure-%s", tx.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Transaction Validation Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Transaction ID: %s failed validation after maximum retries.", tx.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with validation failure for transaction ID: %s.\n", tx.ID)
}

// logBlockFinalization logs the finalization of a block into the ledger
func (automation *TransactionValidationAutomation) logBlockFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("block-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Block Finalization",
        Status:    "Finalized",
        Details:   "Block successfully finalized and sub-blocks merged.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with block finalization.")
}

// encryptTransactionData encrypts transaction data before processing
func (automation *TransactionValidationAutomation) encryptTransactionData(tx common.Transaction) common.Transaction {
    encryptedData, err := encryption.EncryptData(tx.Data)
    if err != nil {
        fmt.Println("Error encrypting transaction data:", err)
        return tx
    }

    tx.EncryptedData = encryptedData
    fmt.Println("Transaction data successfully encrypted.")
    return tx
}

// triggerEmergencyValidationSuspension suspends transaction validation in case of emergency situations
func (automation *TransactionValidationAutomation) triggerEmergencyValidationSuspension(reason string) {
    fmt.Println("Emergency validation suspension triggered. Reason:", reason)
    success := automation.consensusSystem.TriggerEmergencySuspension(reason)

    if success {
        fmt.Println("Emergency validation suspension executed successfully.")
        automation.logEmergencySuspension(reason)
    } else {
        fmt.Println("Emergency validation suspension failed.")
    }
}

// logEmergencySuspension logs an emergency suspension into the ledger
func (automation *TransactionValidationAutomation) logEmergencySuspension(reason string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("emergency-suspension-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Emergency Suspension",
        Status:    "Suspended",
        Details:   fmt.Sprintf("Emergency suspension triggered due to: %s.", reason),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with emergency suspension event.")
}
