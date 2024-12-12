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
    CrossChainCheckInterval = 4000 * time.Millisecond // Interval for checking cross-chain compatibility
    SubBlocksPerBlock       = 1000                    // Number of sub-blocks in a block
)

// CrossChainCompatibilityIntegrationAutomation automates the process of ensuring compatibility between Synnergy Network and external blockchains
type CrossChainCompatibilityIntegrationAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger to store cross-chain integration logs
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    crossChainCheckCount int                        // Counter for cross-chain compatibility check cycles
}

// NewCrossChainCompatibilityIntegrationAutomation initializes the automation for cross-chain compatibility checks
func NewCrossChainCompatibilityIntegrationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossChainCompatibilityIntegrationAutomation {
    return &CrossChainCompatibilityIntegrationAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        crossChainCheckCount: 0,
    }
}

// StartCrossChainCheck starts the continuous loop for checking cross-chain compatibility
func (automation *CrossChainCompatibilityIntegrationAutomation) StartCrossChainCheck() {
    ticker := time.NewTicker(CrossChainCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkCrossChainCompatibility()
        }
    }()
}

// checkCrossChainCompatibility validates if cross-chain transactions and data comply with Synnergy Network's consensus and ledger requirements
func (automation *CrossChainCompatibilityIntegrationAutomation) checkCrossChainCompatibility() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    crossChainTransactions := automation.consensusSystem.GetPendingCrossChainTransactions() // Fetch cross-chain transactions

    for _, tx := range crossChainTransactions {
        fmt.Printf("Validating cross-chain transaction: %s\n", tx.ID)
        isValid := automation.validateCrossChainTransaction(tx)

        if isValid {
            fmt.Printf("Cross-chain transaction %s is valid.\n", tx.ID)
            automation.logCrossChainTransactionResult(tx.ID, "Valid")
        } else {
            fmt.Printf("Cross-chain transaction %s is invalid.\n", tx.ID)
            automation.logCrossChainTransactionResult(tx.ID, "Invalid")
        }
    }

    automation.crossChainCheckCount++
    fmt.Printf("Cross-chain compatibility check cycle #%d completed.\n", automation.crossChainCheckCount)

    if automation.crossChainCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeCrossChainCheckCycle()
    }
}

// validateCrossChainTransaction validates the cross-chain transaction against Synnergy Network's consensus and rules
func (automation *CrossChainCompatibilityIntegrationAutomation) validateCrossChainTransaction(tx common.CrossChainTransaction) bool {
    // Encrypt transaction data before validation
    fmt.Printf("Encrypting cross-chain transaction data for: %s\n", tx.ID)

    encryptedTxData, err := encryption.EncryptData(tx)
    if err != nil {
        fmt.Printf("Error encrypting cross-chain transaction data for %s: %s\n", tx.ID, err.Error())
        return false
    }

    tx.EncryptedData = encryptedTxData
    fmt.Printf("Cross-chain transaction data for %s encrypted successfully.\n", tx.ID)

    // Validate transaction via Synnergy Consensus
    return automation.consensusSystem.ValidateCrossChainTransaction(tx)
}

// logCrossChainTransactionResult logs the result of the cross-chain transaction validation into the ledger
func (automation *CrossChainCompatibilityIntegrationAutomation) logCrossChainTransactionResult(txID string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cross-chain-tx-validation-%s", txID),
        Timestamp: time.Now().Unix(),
        Type:      "Cross-Chain Transaction Validation",
        Status:    result,
        Details:   fmt.Sprintf("Validation result for cross-chain transaction %s: %s", txID, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cross-chain transaction validation result for transaction %s: %s.\n", txID, result)
}

// finalizeCrossChainCheckCycle finalizes the cross-chain compatibility check cycle and logs the result in the ledger
func (automation *CrossChainCompatibilityIntegrationAutomation) finalizeCrossChainCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCrossChainCheckCycle()
    if success {
        fmt.Println("Cross-chain compatibility check cycle finalized successfully.")
        automation.logCrossChainCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing cross-chain compatibility check cycle.")
    }
}

// logCrossChainCheckCycleFinalization logs the finalization of a cross-chain check cycle into the ledger
func (automation *CrossChainCompatibilityIntegrationAutomation) logCrossChainCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cross-chain-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cross-Chain Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with cross-chain compatibility check cycle finalization.")
}
