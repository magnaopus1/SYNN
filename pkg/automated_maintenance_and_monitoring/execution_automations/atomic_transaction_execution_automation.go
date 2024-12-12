package execution_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    TransactionMonitorInterval    = 5 * time.Second // Interval for checking pending atomic transactions
    AtomicTransactionTimeout      = 10 * time.Minute // Timeout for an atomic transaction to be processed
)

// AtomicTransactionExecutionAutomation manages atomic transaction execution across the blockchain
type AtomicTransactionExecutionAutomation struct {
    consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    ledgerInstance    *ledger.Ledger                         // Ledger instance to track transaction history
    stateMutex        *sync.RWMutex                          // Mutex for thread-safe operations
}

// NewAtomicTransactionExecutionAutomation initializes the atomic transaction automation
func NewAtomicTransactionExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AtomicTransactionExecutionAutomation {
    return &AtomicTransactionExecutionAutomation{
        consensusEngine:  consensusEngine,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
    }
}

// StartTransactionMonitor begins the continuous monitoring of atomic transactions
func (automation *AtomicTransactionExecutionAutomation) StartTransactionMonitor() {
    ticker := time.NewTicker(TransactionMonitorInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Checking for pending atomic transactions...")
            automation.monitorPendingTransactions()
        }
    }()
}

// monitorPendingTransactions checks for any atomic transactions that need to be executed
func (automation *AtomicTransactionExecutionAutomation) monitorPendingTransactions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    pendingTransactions := automation.consensusEngine.GetPendingAtomicTransactions()

    for _, tx := range pendingTransactions {
        if time.Since(time.Unix(tx.Timestamp, 0)) > AtomicTransactionTimeout {
            fmt.Printf("Atomic transaction %s timed out and will be canceled.\n", tx.ID)
            automation.cancelAtomicTransaction(tx)
        } else {
            fmt.Printf("Executing atomic transaction %s.\n", tx.ID)
            automation.executeAtomicTransaction(tx)
        }
    }
}

// executeAtomicTransaction performs the atomic transaction and logs it to the ledger
func (automation *AtomicTransactionExecutionAutomation) executeAtomicTransaction(tx common.Transaction) {
    // Execute the transaction using the consensus engine
    success := automation.consensusEngine.ProcessTransaction(tx)

    // Log the result to the ledger
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("atomic-tx-execution-%s", tx.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Atomic Transaction Execution",
        Status:    "Success",
        Details:   fmt.Sprintf("Atomic transaction %s successfully executed.", tx.ID),
    }

    if !success {
        entry.Status = "Failure"
        entry.Details = fmt.Sprintf("Failed to execute atomic transaction %s.", tx.ID)
    }

    // Encrypt and log the ledger entry
    encryptedEntry := automation.encryptData(entry.Details)
    entry.Details = encryptedEntry

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Error logging atomic transaction execution: %v\n", err)
    } else {
        fmt.Printf("Atomic transaction %s successfully logged in the ledger.\n", tx.ID)
    }
}

// cancelAtomicTransaction cancels a timed-out atomic transaction and logs it to the ledger
func (automation *AtomicTransactionExecutionAutomation) cancelAtomicTransaction(tx common.Transaction) {
    // Cancel the transaction using the consensus engine
    success := automation.consensusEngine.CancelTransaction(tx.ID)

    // Log the cancellation to the ledger
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("atomic-tx-cancel-%s", tx.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Atomic Transaction Cancellation",
        Status:    "Success",
        Details:   fmt.Sprintf("Atomic transaction %s was canceled due to timeout.", tx.ID),
    }

    if !success {
        entry.Status = "Failure"
        entry.Details = fmt.Sprintf("Failed to cancel atomic transaction %s.", tx.ID)
    }

    // Encrypt and log the ledger entry
    encryptedEntry := automation.encryptData(entry.Details)
    entry.Details = encryptedEntry

    if err := automation.ledgerInstance.AddEntry(entry); err != nil {
        fmt.Printf("Error logging atomic transaction cancellation: %v\n", err)
    } else {
        fmt.Printf("Atomic transaction %s cancellation successfully logged in the ledger.\n", tx.ID)
    }
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AtomicTransactionExecutionAutomation) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}
