package data_automations

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
    TransactionEnforcementInterval = 500 * time.Millisecond // Run enforcement checks every 0.5 seconds
    MaxPendingTransactions         = 100                    // Maximum number of pending transactions before enforcement
    SubBlocksPerBlock              = 1000                   // Number of sub-blocks in a block
)

// CrossChainTransactionEnforcementAutomation automates the enforcement of cross-chain transactions using Synnergy Consensus and integrates with the ledger
type CrossChainTransactionEnforcementAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store consensus and transaction enforcement data
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    pendingTxCount    int                          // Counter for pending cross-chain transactions
}

// NewCrossChainTransactionEnforcementAutomation initializes the automation for cross-chain transaction enforcement
func NewCrossChainTransactionEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossChainTransactionEnforcementAutomation {
    return &CrossChainTransactionEnforcementAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        pendingTxCount:  0,
    }
}

// StartTransactionEnforcementAutomation starts the continuous automation loop for cross-chain transaction enforcement
func (automation *CrossChainTransactionEnforcementAutomation) StartTransactionEnforcementAutomation() {
    ticker := time.NewTicker(TransactionEnforcementInterval)

    go func() {
        for range ticker.C {
            automation.enforceTransactionProcessing()
        }
    }()
}

// enforceTransactionProcessing checks pending cross-chain transactions and triggers enforcement when thresholds are reached
func (automation *CrossChainTransactionEnforcementAutomation) enforceTransactionProcessing() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch pending cross-chain transactions from the consensus system
    pendingTxs := automation.consensusSystem.GetPendingCrossChainTransactions()

    if len(pendingTxs) >= MaxPendingTransactions {
        fmt.Printf("Pending transactions exceed limit (%d). Triggering transaction enforcement.\n", len(pendingTxs))
        automation.triggerTransactionEnforcement(pendingTxs)
    } else {
        fmt.Printf("Pending transactions are within acceptable range (%d).\n", len(pendingTxs))
    }

    automation.pendingTxCount++
    fmt.Printf("Transaction enforcement cycle #%d executed.\n", automation.pendingTxCount)

    // Finalize transactions into sub-blocks and blocks periodically
    if automation.pendingTxCount%SubBlocksPerBlock == 0 {
        automation.finalizeBlock()
    }
}

// triggerTransactionEnforcement enforces cross-chain transaction validation when pending transactions exceed the threshold
func (automation *CrossChainTransactionEnforcementAutomation) triggerTransactionEnforcement(pendingTxs []common.Transaction) {
    for _, tx := range pendingTxs {
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting validator for transaction enforcement.")
            continue
        }

        // Encrypt transaction data before enforcement
        encryptedTx := automation.AddEncryption(tx)

        fmt.Printf("Validator %s selected for enforcing cross-chain transaction.\n", validator.Address)

        // Validate transaction via consensus
        validationSuccess := automation.consensusSystem.EnforceCrossChainTransaction(validator, encryptedTx)
        if validationSuccess {
            fmt.Println("Cross-chain transaction successfully enforced.")
        } else {
            fmt.Println("Error enforcing cross-chain transaction.")
        }

        // Log the enforcement action into the ledger
        automation.logTransactionEnforcement(tx)
    }
}

// finalizeBlock finalizes 1000 sub-blocks into a full block using consensus PoW
func (automation *CrossChainTransactionEnforcementAutomation) finalizeBlock() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.PoW.MineBlock()
    if success {
        fmt.Println("Block successfully finalized with PoW.")
        automation.logBlockFinalization()
    } else {
        fmt.Println("Error finalizing block with PoW.")
    }
}

// logTransactionEnforcement logs every cross-chain transaction enforcement into the ledger
func (automation *CrossChainTransactionEnforcementAutomation) logTransactionEnforcement(tx common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("transaction-enforcement-%s", tx.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Cross-Chain Transaction Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with transaction enforcement event for TxID %s.\n", tx.ID)
}

// logBlockFinalization logs the block finalization event into the ledger
func (automation *CrossChainTransactionEnforcementAutomation) logBlockFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("block-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Block Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with block finalization event.\n")
}

// AddEncryption adds encryption to sensitive transaction data during enforcement
func (automation *CrossChainTransactionEnforcementAutomation) AddEncryption(tx common.Transaction) common.Transaction {
    encryptedData, err := encryption.EncryptData(tx.Data)
    if err != nil {
        fmt.Println("Error encrypting transaction data:", err)
        return tx
    }
    tx.Data = encryptedData
    fmt.Println("Transaction data successfully encrypted.")
    return tx
}

// ensureCrossChainIntegrity checks the integrity of cross-chain transactions and triggers enforcement when integrity breaches are detected
func (automation *CrossChainTransactionEnforcementAutomation) ensureCrossChainIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCrossChainIntegrity()
    if !integrityValid {
        fmt.Println("Cross-chain integrity breach detected. Triggering transaction enforcement.")
        automation.triggerTransactionEnforcement(automation.consensusSystem.GetPendingCrossChainTransactions())
    } else {
        fmt.Println("Cross-chain transaction integrity is valid.")
    }
}
