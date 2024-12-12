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
    ShardTransactionEnforcementInterval = 700 * time.Millisecond // Run enforcement checks every 0.7 seconds
    MaxPendingShardTransactions         = 150                    // Max pending transactions across shards before enforcement
)

// CrossShardTransactionEnforcementAutomation automates the enforcement of cross-shard transactions using Synnergy Consensus and integrates with the ledger
type CrossShardTransactionEnforcementAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store consensus and transaction enforcement data
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    pendingShardTxCount int                          // Counter for pending cross-shard transactions
}

// NewCrossShardTransactionEnforcementAutomation initializes the automation for cross-shard transaction enforcement
func NewCrossShardTransactionEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossShardTransactionEnforcementAutomation {
    return &CrossShardTransactionEnforcementAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        pendingShardTxCount: 0,
    }
}

// StartTransactionEnforcementAutomation starts the continuous automation loop for cross-shard transaction enforcement
func (automation *CrossShardTransactionEnforcementAutomation) StartTransactionEnforcementAutomation() {
    ticker := time.NewTicker(ShardTransactionEnforcementInterval)

    go func() {
        for range ticker.C {
            automation.enforceShardTransactionProcessing()
        }
    }()
}

// enforceShardTransactionProcessing checks pending cross-shard transactions and triggers enforcement when thresholds are reached
func (automation *CrossShardTransactionEnforcementAutomation) enforceShardTransactionProcessing() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch pending cross-shard transactions from the consensus system
    pendingShardTxs := automation.consensusSystem.GetPendingCrossShardTransactions()

    if len(pendingShardTxs) >= MaxPendingShardTransactions {
        fmt.Printf("Pending cross-shard transactions exceed limit (%d). Triggering transaction enforcement.\n", len(pendingShardTxs))
        automation.triggerShardTransactionEnforcement(pendingShardTxs)
    } else {
        fmt.Printf("Pending cross-shard transactions are within acceptable range (%d).\n", len(pendingShardTxs))
    }

    automation.pendingShardTxCount++
    fmt.Printf("Shard transaction enforcement cycle #%d executed.\n", automation.pendingShardTxCount)

    // Finalize sub-blocks and blocks when necessary
    if automation.pendingShardTxCount%SubBlocksPerBlock == 0 {
        automation.finalizeBlock()
    }
}

// triggerShardTransactionEnforcement enforces cross-shard transaction validation when pending transactions exceed the threshold
func (automation *CrossShardTransactionEnforcementAutomation) triggerShardTransactionEnforcement(pendingShardTxs []common.Transaction) {
    for _, tx := range pendingShardTxs {
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting validator for cross-shard transaction enforcement.")
            continue
        }

        // Encrypt transaction data before enforcement
        encryptedTx := automation.AddEncryption(tx)

        fmt.Printf("Validator %s selected for enforcing cross-shard transaction.\n", validator.Address)

        // Validate cross-shard transaction via consensus
        validationSuccess := automation.consensusSystem.EnforceCrossShardTransaction(validator, encryptedTx)
        if validationSuccess {
            fmt.Println("Cross-shard transaction successfully enforced.")
        } else {
            fmt.Println("Error enforcing cross-shard transaction.")
        }

        // Log the enforcement action into the ledger
        automation.logShardTransactionEnforcement(tx)
    }
}

// finalizeBlock finalizes 1000 sub-blocks into a full block using consensus PoW
func (automation *CrossShardTransactionEnforcementAutomation) finalizeBlock() {
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

// logShardTransactionEnforcement logs every cross-shard transaction enforcement into the ledger
func (automation *CrossShardTransactionEnforcementAutomation) logShardTransactionEnforcement(tx common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cross-shard-transaction-enforcement-%s", tx.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Cross-Shard Transaction Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cross-shard transaction enforcement event for TxID %s.\n", tx.ID)
}

// logBlockFinalization logs the block finalization event into the ledger
func (automation *CrossShardTransactionEnforcementAutomation) logBlockFinalization() {
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
func (automation *CrossShardTransactionEnforcementAutomation) AddEncryption(tx common.Transaction) common.Transaction {
    encryptedData, err := encryption.EncryptData(tx.Data)
    if err != nil {
        fmt.Println("Error encrypting transaction data:", err)
        return tx
    }
    tx.Data = encryptedData
    fmt.Println("Transaction data successfully encrypted.")
    return tx
}

// ensureCrossShardIntegrity checks the integrity of cross-shard transactions and triggers enforcement when integrity breaches are detected
func (automation *CrossShardTransactionEnforcementAutomation) ensureCrossShardIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCrossShardIntegrity()
    if !integrityValid {
        fmt.Println("Cross-shard integrity breach detected. Triggering transaction enforcement.")
        automation.triggerShardTransactionEnforcement(automation.consensusSystem.GetPendingCrossShardTransactions())
    } else {
        fmt.Println("Cross-shard transaction integrity is valid.")
    }
}
