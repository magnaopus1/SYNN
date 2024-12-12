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
    InsuranceMonitoringInterval   = 10 * time.Second  // Interval for monitoring cross-chain bridge transactions
    MaxInsuranceClaimRetries      = 3                 // Maximum number of retries for processing an insurance claim
    SubBlocksPerBlock             = 1000              // Number of sub-blocks in a block
)

// CrossChainBridgeInsuranceAutomation ensures protection of assets transferred over cross-chain bridges
type CrossChainBridgeInsuranceAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging insurance events
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    claimRetryCount      map[string]int               // Counter for retrying insurance claims on failure
    insuranceCycleCount  int                          // Counter for monitoring cycles
}

// NewCrossChainBridgeInsuranceAutomation initializes the automation for cross-chain bridge insurance
func NewCrossChainBridgeInsuranceAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossChainBridgeInsuranceAutomation {
    return &CrossChainBridgeInsuranceAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        claimRetryCount:     make(map[string]int),
        insuranceCycleCount: 0,
    }
}

// StartInsuranceMonitoring starts the continuous loop for monitoring cross-chain bridge transactions
func (automation *CrossChainBridgeInsuranceAutomation) StartInsuranceMonitoring() {
    ticker := time.NewTicker(InsuranceMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorCrossChainBridge()
        }
    }()
}

// monitorCrossChainBridge monitors cross-chain bridge transactions and ensures insurance coverage
func (automation *CrossChainBridgeInsuranceAutomation) monitorCrossChainBridge() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch a list of transactions needing insurance verification
    pendingBridgeTransactions := automation.consensusSystem.GetPendingBridgeTransactions()

    if len(pendingBridgeTransactions) > 0 {
        for _, transaction := range pendingBridgeTransactions {
            fmt.Printf("Processing insurance for cross-chain transaction %s.\n", transaction.ID)
            automation.processInsurance(transaction)
        }
    } else {
        fmt.Println("No cross-chain transactions need insurance this cycle.")
    }

    automation.insuranceCycleCount++
    fmt.Printf("Insurance monitoring cycle #%d executed.\n", automation.insuranceCycleCount)

    if automation.insuranceCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeInsuranceCycle()
    }
}

// processInsurance applies insurance coverage to cross-chain transactions
func (automation *CrossChainBridgeInsuranceAutomation) processInsurance(transaction common.Transaction) {
    // Encrypt transaction data before processing insurance
    encryptedTransactionData := automation.encryptTransactionData(transaction)

    // Attempt to apply insurance through the Synnergy Consensus system
    insuranceSuccess := automation.consensusSystem.ApplyInsurance(transaction, encryptedTransactionData)

    if insuranceSuccess {
        fmt.Printf("Insurance successfully applied for cross-chain transaction %s.\n", transaction.ID)
        automation.logInsuranceEvent(transaction, "Insured")
        automation.resetClaimRetry(transaction.ID)
    } else {
        fmt.Printf("Error processing insurance for transaction %s. Retrying...\n", transaction.ID)
        automation.retryInsuranceClaim(transaction)
    }
}

// retryInsuranceClaim attempts to retry a failed insurance claim a limited number of times
func (automation *CrossChainBridgeInsuranceAutomation) retryInsuranceClaim(transaction common.Transaction) {
    automation.claimRetryCount[transaction.ID]++
    if automation.claimRetryCount[transaction.ID] < MaxInsuranceClaimRetries {
        automation.processInsurance(transaction)
    } else {
        fmt.Printf("Max retries reached for insurance claim on transaction %s. Insurance claim failed.\n", transaction.ID)
        automation.logClaimFailure(transaction)
    }
}

// resetClaimRetry resets the retry count for insurance claims
func (automation *CrossChainBridgeInsuranceAutomation) resetClaimRetry(transactionID string) {
    automation.claimRetryCount[transactionID] = 0
}

// finalizeInsuranceCycle finalizes the cross-chain insurance cycle and logs the result in the ledger
func (automation *CrossChainBridgeInsuranceAutomation) finalizeInsuranceCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeInsuranceCycle()
    if success {
        fmt.Println("Cross-chain insurance cycle finalized successfully.")
        automation.logInsuranceCycleFinalization()
    } else {
        fmt.Println("Error finalizing cross-chain insurance cycle.")
    }
}

// logInsuranceEvent logs a successful insurance event into the ledger
func (automation *CrossChainBridgeInsuranceAutomation) logInsuranceEvent(transaction common.Transaction, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insurance-%s-%s", transaction.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Cross-Chain Insurance Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Transaction %s insured successfully.", transaction.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with insurance event for transaction %s.\n", transaction.ID)
}

// logClaimFailure logs the failure of an insurance claim into the ledger
func (automation *CrossChainBridgeInsuranceAutomation) logClaimFailure(transaction common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insurance-failure-%s", transaction.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Insurance Claim Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Insurance claim failed for transaction %s after maximum retries.", transaction.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with insurance claim failure for transaction %s.\n", transaction.ID)
}

// logInsuranceCycleFinalization logs the finalization of an insurance cycle into the ledger
func (automation *CrossChainBridgeInsuranceAutomation) logInsuranceCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("insurance-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Insurance Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with insurance cycle finalization.")
}

// encryptTransactionData encrypts transaction data before processing insurance
func (automation *CrossChainBridgeInsuranceAutomation) encryptTransactionData(transaction common.Transaction) common.Transaction {
    encryptedData, err := encryption.EncryptData(transaction.Data)
    if err != nil {
        fmt.Println("Error encrypting transaction data for insurance:", err)
        return transaction
    }

    transaction.EncryptedData = encryptedData
    fmt.Println("Transaction data successfully encrypted for insurance.")
    return transaction
}

// manualIntervention allows for manual intervention in the insurance process
func (automation *CrossChainBridgeInsuranceAutomation) manualIntervention(transaction common.Transaction, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if action == "insure" {
        fmt.Printf("Manually applying insurance to transaction %s.\n", transaction.ID)
        automation.processInsurance(transaction)
    } else if action == "ignore" {
        fmt.Printf("Manually ignoring insurance application for transaction %s.\n", transaction.ID)
    } else {
        fmt.Println("Invalid action for manual intervention.")
    }
}

// emergencyUninsurance triggers emergency uninsurance of a transaction in case of critical needs
func (automation *CrossChainBridgeInsuranceAutomation) emergencyUninsurance(transaction common.Transaction) {
    fmt.Printf("Emergency uninsurance triggered for transaction %s.\n", transaction.ID)
    success := automation.consensusSystem.TriggerEmergencyUninsurance(transaction)

    if success {
        automation.logInsuranceEvent(transaction, "Uninsured")
        fmt.Println("Emergency uninsurance executed successfully.")
    } else {
        fmt.Println("Emergency uninsurance failed.")
    }
}
