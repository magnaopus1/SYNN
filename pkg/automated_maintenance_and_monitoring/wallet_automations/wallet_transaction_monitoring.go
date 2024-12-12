package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

const (
    TransactionMonitoringInterval = 10 * time.Second // Interval for monitoring wallet transactions
    MaxAllowedTransactionValue    = 1000000          // Threshold for high-value transactions that require extra monitoring
    IrregularTransactionRate      = 0.25             // Percentage deviation to detect irregular transaction patterns
)

// WalletTransactionMonitoringAutomation continuously monitors transactions to detect irregularities
type WalletTransactionMonitoringAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging transaction events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    flaggedTransactions  map[string]bool              // Store flagged transactions
}

// NewWalletTransactionMonitoringAutomation initializes the automation for monitoring wallet transactions
func NewWalletTransactionMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *WalletTransactionMonitoringAutomation {
    return &WalletTransactionMonitoringAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        flaggedTransactions: make(map[string]bool),
    }
}

// StartTransactionMonitoring starts the continuous loop for monitoring wallet transactions
func (automation *WalletTransactionMonitoringAutomation) StartTransactionMonitoring() {
    ticker := time.NewTicker(TransactionMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorWalletTransactions()
        }
    }()
}

// monitorWalletTransactions checks the transactions for any irregular or suspicious behavior
func (automation *WalletTransactionMonitoringAutomation) monitorWalletTransactions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of recent wallet transactions from the consensus system
    recentTransactions := automation.consensusSystem.GetRecentTransactions()

    for _, transaction := range recentTransactions {
        if automation.isIrregularTransaction(transaction) {
            automation.handleIrregularTransaction(transaction)
        } else {
            fmt.Printf("Transaction %s is within normal parameters.\n", transaction.TransactionID)
        }
    }
}

// isIrregularTransaction checks if the transaction has any irregularities that need attention
func (automation *WalletTransactionMonitoringAutomation) isIrregularTransaction(transaction common.Transaction) bool {
    // Check if the transaction exceeds the allowed transaction value
    if transaction.Amount > MaxAllowedTransactionValue {
        fmt.Printf("Transaction %s flagged for high value: %f\n", transaction.TransactionID, transaction.Amount)
        return true
    }

    // Check if the transaction rate is significantly different from historical patterns
    if automation.consensusSystem.DetectIrregularTransactionRate(transaction, IrregularTransactionRate) {
        fmt.Printf("Transaction %s flagged for irregular rate.\n", transaction.TransactionID)
        return true
    }

    return false
}

// handleIrregularTransaction processes flagged transactions
func (automation *WalletTransactionMonitoringAutomation) handleIrregularTransaction(transaction common.Transaction) {
    fmt.Printf("Handling irregular transaction: %s\n", transaction.TransactionID)

    // Encrypt and log the irregular transaction
    encryptedTransaction := automation.encryptTransactionData(transaction)
    automation.logIrregularTransaction(encryptedTransaction)

    // Take action depending on the severity of the irregularity
    if transaction.Amount > MaxAllowedTransactionValue {
        automation.enforceTransactionLimit(transaction)
    }

    if automation.consensusSystem.IsFraudulentTransaction(transaction) {
        err := automation.consensusSystem.LockWallet(transaction.FromWalletID)
        if err != nil {
            fmt.Printf("Error locking wallet %s: %v\n", transaction.FromWalletID, err)
        } else {
            fmt.Printf("Wallet %s locked due to fraudulent transaction.\n", transaction.FromWalletID)
        }
    }
}

// logIrregularTransaction logs the irregular transaction into the ledger
func (automation *WalletTransactionMonitoringAutomation) logIrregularTransaction(transaction common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("irregular-transaction-%s", transaction.TransactionID),
        Timestamp: time.Now().Unix(),
        Type:      "Irregular Transaction",
        Status:    "Flagged",
        Details:   fmt.Sprintf("Transaction %s flagged as irregular. Action taken.", transaction.TransactionID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with irregular transaction %s.\n", transaction.TransactionID)
}

// encryptTransactionData encrypts transaction data before processing or logging
func (automation *WalletTransactionMonitoringAutomation) encryptTransactionData(transaction common.Transaction) common.Transaction {
    encryptedData, err := encryption.EncryptData(transaction)
    if err != nil {
        fmt.Println("Error encrypting transaction data:", err)
        return transaction
    }

    transaction.EncryptedData = encryptedData
    fmt.Println("Transaction data successfully encrypted.")
    return transaction
}

// enforceTransactionLimit blocks or reviews transactions that exceed a certain threshold
func (automation *WalletTransactionMonitoringAutomation) enforceTransactionLimit(transaction common.Transaction) {
    fmt.Printf("Enforcing transaction limit for transaction %s.\n", transaction.TransactionID)

    // Depending on business logic, either block, hold, or notify admins
    if err := automation.consensusSystem.BlockTransaction(transaction.TransactionID); err != nil {
        fmt.Printf("Error blocking transaction %s: %v\n", transaction.TransactionID, err)
    } else {
        fmt.Printf("Transaction %s blocked due to exceeding value limits.\n", transaction.TransactionID)
    }

    // Optionally notify admins
    automation.notifyAdmin(transaction)
}

// notifyAdmin sends a notification for admin review of suspicious transactions
func (automation *WalletTransactionMonitoringAutomation) notifyAdmin(transaction common.Transaction) {
    fmt.Printf("Admin notified of suspicious transaction: %s\n", transaction.TransactionID)
    // Implementation for sending a notification to admins
}

// ensureTransactionIntegrity checks the integrity of transaction data and the monitoring system
func (automation *WalletTransactionMonitoringAutomation) ensureTransactionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateTransactionIntegrity()
    if !integrityValid {
        fmt.Println("Transaction integrity breach detected. Re-running transaction monitoring.")
        automation.monitorWalletTransactions()
    } else {
        fmt.Println("Transaction integrity is valid.")
    }
}
