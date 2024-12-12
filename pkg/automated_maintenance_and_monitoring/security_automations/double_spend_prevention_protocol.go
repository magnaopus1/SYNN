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
    DoubleSpendCheckInterval = 5 * time.Second // Interval for checking double-spending attempts
    TransactionLogRetention  = 30 * time.Minute // Retention period for transaction logs in memory
)

// DoubleSpendPreventionAutomation automates the detection and prevention of double-spending in the blockchain
type DoubleSpendPreventionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus for sub-block validation
    ledgerInstance    *ledger.Ledger               // Ledger for logging double-spend attempts
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    recentTransactions map[string]time.Time        // In-memory cache of recent transactions for checking double spends
}

// NewDoubleSpendPreventionAutomation initializes the automation for double-spend prevention
func NewDoubleSpendPreventionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DoubleSpendPreventionAutomation {
    return &DoubleSpendPreventionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        recentTransactions: make(map[string]time.Time),
    }
}

// StartDoubleSpendPrevention starts the continuous loop for detecting and preventing double-spending attempts
func (automation *DoubleSpendPreventionAutomation) StartDoubleSpendPrevention() {
    ticker := time.NewTicker(DoubleSpendCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorDoubleSpending()
        }
    }()
}

// monitorDoubleSpending checks for any double-spending attempts in the network
func (automation *DoubleSpendPreventionAutomation) monitorDoubleSpending() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch list of recent transactions from the consensus system
    transactionList := automation.consensusSystem.GetRecentTransactions()

    if len(transactionList) > 0 {
        for _, tx := range transactionList {
            automation.checkForDoubleSpend(tx)
        }
    } else {
        fmt.Println("No transactions detected for double-spend monitoring at this time.")
    }

    automation.cleanupOldTransactions()
}

// checkForDoubleSpend checks if a transaction is an attempt to double spend
func (automation *DoubleSpendPreventionAutomation) checkForDoubleSpend(tx common.Transaction) {
    if _, exists := automation.recentTransactions[tx.ID]; exists {
        fmt.Printf("Double spend detected for transaction %s.\n", tx.ID)
        automation.handleDoubleSpendAttempt(tx)
    } else {
        automation.recentTransactions[tx.ID] = time.Now()
    }
}

// handleDoubleSpendAttempt handles the response to a detected double-spend attempt
func (automation *DoubleSpendPreventionAutomation) handleDoubleSpendAttempt(tx common.Transaction) {
    // Trigger double-spend prevention in the Synnergy Consensus system
    preventionSuccess := automation.consensusSystem.PreventDoubleSpend(tx)

    if preventionSuccess {
        fmt.Printf("Double-spend attempt for transaction %s prevented.\n", tx.ID)
        automation.logDoubleSpendAttempt(tx, "Prevented")
    } else {
        fmt.Printf("Failed to prevent double-spend attempt for transaction %s.\n", tx.ID)
        automation.logDoubleSpendAttempt(tx, "Failed")
    }
}

// logDoubleSpendAttempt logs the double-spend attempt in the ledger
func (automation *DoubleSpendPreventionAutomation) logDoubleSpendAttempt(tx common.Transaction, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("double-spend-%s", tx.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Double Spend",
        Status:    status,
        Details:   fmt.Sprintf("Double spend attempt detected for transaction %s. Status: %s", tx.ID, status),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with double-spend attempt for transaction %s.\n", tx.ID)
}

// cleanupOldTransactions removes old transactions from memory after the retention period
func (automation *DoubleSpendPreventionAutomation) cleanupOldTransactions() {
    for txID, timestamp := range automation.recentTransactions {
        if time.Since(timestamp) > TransactionLogRetention {
            delete(automation.recentTransactions, txID)
            fmt.Printf("Transaction %s removed from double-spend tracking after retention period.\n", txID)
        }
    }
}

// EncryptTransactionData encrypts transaction data before validation
func (automation *DoubleSpendPreventionAutomation) EncryptTransactionData(tx common.Transaction) common.Transaction {
    encryptedData, err := encryption.EncryptData(tx.Data)
    if err != nil {
        fmt.Println("Error encrypting transaction data:", err)
        return tx
    }

    tx.EncryptedData = encryptedData
    fmt.Println("Transaction data successfully encrypted.")    tx
}

// ValidateTransactionIntegrity checks the integrity of the transaction data
func (automation *DoubleSpendPreventionAutomation) ValidateTransactionIntegrity(tx common.Transaction) bool {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    isValid := automation.consensusSystem.ValidateTransactionIntegrity(tx)
    if isValid {
        fmt.Printf("Transaction %s passed integrity validation.\n", tx.ID)
    } else {
        fmt.Printf("Transaction %s failed integrity validation. Possible tampering detected.\n", tx.ID)
    }

    return isValid
}

// FinalizeDoubleSpendPreventionCycle finalizes the double-spend prevention cycle, clearing up state and logging completion
func (automation *DoubleSpendPreventionAutomation) FinalizeDoubleSpendPreventionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Log the cycle finalization in the ledger
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("double-spend-prevention-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Double Spend Prevention Cycle Finalization",
        Status:    "Finalized",
        Details:   "Double spend prevention cycle completed and finalized.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with double-spend prevention cycle finalization.")
}

   
