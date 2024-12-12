package security_automations

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
    GasLimitMonitoringInterval    = 15 * time.Second // Interval for monitoring gas limits
    MaxGasUsageThreshold          = 1000000          // Maximum gas limit per transaction
    EmergencyGasLimitAdjustment   = 500000           // Emergency reduction in gas limit
    GasLimitAdjustmentThreshold   = 10               // Number of transactions exceeding gas limit to trigger adjustment
)

// GasLimitSecurityAutomation ensures that gas usage is within acceptable limits and enforces gas security protocols
type GasLimitSecurityAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging gas limit adjustments and events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    gasUsageTracker        map[string]int               // Tracks gas usage per wallet or contract
    gasLimitExceedCount    int                          // Counter for the number of transactions exceeding the gas limit
}

// NewGasLimitSecurityAutomation initializes the gas limit security automation
func NewGasLimitSecurityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GasLimitSecurityAutomation {
    return &GasLimitSecurityAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        gasUsageTracker:     make(map[string]int),
        gasLimitExceedCount: 0,
    }
}

// StartGasLimitMonitoring starts the continuous loop for monitoring gas limit usage
func (automation *GasLimitSecurityAutomation) StartGasLimitMonitoring() {
    ticker := time.NewTicker(GasLimitMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorGasUsage()
        }
    }()
}

// monitorGasUsage monitors all transactions for gas limit violations
func (automation *GasLimitSecurityAutomation) monitorGasUsage() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    transactionList := automation.consensusSystem.GetPendingTransactions()

    if len(transactionList) > 0 {
        for _, transaction := range transactionList {
            automation.enforceGasLimit(transaction)
        }
    } else {
        fmt.Println("No pending transactions at this time.")
    }

    if automation.gasLimitExceedCount >= GasLimitAdjustmentThreshold {
        automation.triggerGasLimitAdjustment()
    }
}

// enforceGasLimit checks the gas usage of a transaction and applies limits if exceeded
func (automation *GasLimitSecurityAutomation) enforceGasLimit(transaction common.Transaction) {
    if transaction.GasUsed > MaxGasUsageThreshold {
        fmt.Printf("Gas limit exceeded for transaction %s. Gas used: %d\n", transaction.ID, transaction.GasUsed)
        automation.blockTransaction(transaction)
        automation.logGasLimitViolation(transaction)
        automation.gasLimitExceedCount++
    } else {
        fmt.Printf("Transaction %s is within gas limits.\n", transaction.ID)
        automation.resetGasLimitExceedCount()
    }
}

// blockTransaction prevents the execution of a transaction that exceeds gas limits
func (automation *GasLimitSecurityAutomation) blockTransaction(transaction common.Transaction) {
    automation.consensusSystem.BlockTransaction(transaction)
    fmt.Printf("Transaction %s blocked due to gas limit violation.\n", transaction.ID)
}

// logGasLimitViolation logs a gas limit violation event into the ledger
func (automation *GasLimitSecurityAutomation) logGasLimitViolation(transaction common.Transaction) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-limit-violation-%s", transaction.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Limit Violation",
        Status:    "Blocked",
        Details:   fmt.Sprintf("Transaction %s exceeded gas limit with %d gas used.", transaction.ID, transaction.GasUsed),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with gas limit violation for transaction %s.\n", transaction.ID)
}

// triggerGasLimitAdjustment adjusts the gas limit temporarily if violations exceed threshold
func (automation *GasLimitSecurityAutomation) triggerGasLimitAdjustment() {
    fmt.Println("Gas limit adjustment triggered due to repeated violations.")
    automation.consensusSystem.AdjustGasLimit(EmergencyGasLimitAdjustment)
    automation.logGasLimitAdjustment()
}

// logGasLimitAdjustment logs the gas limit adjustment event into the ledger
func (automation *GasLimitSecurityAutomation) logGasLimitAdjustment() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-limit-adjustment-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Limit Adjustment",
        Status:    "Adjusted",
        Details:   fmt.Sprintf("Gas limit adjusted to %d due to repeated violations.", EmergencyGasLimitAdjustment),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with gas limit adjustment.")
}

// resetGasLimitExceedCount resets the gas limit exceed counter
func (automation *GasLimitSecurityAutomation) resetGasLimitExceedCount() {
    automation.gasLimitExceedCount = 0
}
