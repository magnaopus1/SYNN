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
    ReplicationEnforcementInterval    = 800 * time.Millisecond // Interval for checking data replication issues
    MaxReplicationViolationLimit      = 10                     // Maximum number of replication violations before enforcement
)

// DataReplicationEnforcementAutomation automates the enforcement of data replication across the blockchain network
type DataReplicationEnforcementAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store replication-related logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    replicationViolationCount int                        // Counter for replication violations
}

// NewDataReplicationEnforcementAutomation initializes the automation for data replication enforcement
func NewDataReplicationEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataReplicationEnforcementAutomation {
    return &DataReplicationEnforcementAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        replicationViolationCount: 0,
    }
}

// StartReplicationEnforcementAutomation starts the continuous loop for monitoring and enforcing data replication
func (automation *DataReplicationEnforcementAutomation) StartReplicationEnforcementAutomation() {
    ticker := time.NewTicker(ReplicationEnforcementInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceReplication()
        }
    }()
}

// monitorAndEnforceReplication checks for replication violations and triggers enforcement if necessary
func (automation *DataReplicationEnforcementAutomation) monitorAndEnforceReplication() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch replication violations from the consensus system
    replicationViolations := automation.consensusSystem.CheckReplicationViolations()

    if len(replicationViolations) >= MaxReplicationViolationLimit {
        fmt.Printf("Replication violations exceed limit (%d). Triggering enforcement.\n", len(replicationViolations))
        automation.triggerReplicationEnforcement(replicationViolations)
    } else {
        fmt.Printf("Replication violations are within acceptable range (%d).\n", len(replicationViolations))
    }

    automation.replicationViolationCount++
    fmt.Printf("Replication enforcement cycle #%d executed.\n", automation.replicationViolationCount)

    // Finalize the replication enforcement cycle after the set number of sub-blocks
    if automation.replicationViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerReplicationEnforcement enforces data replication policies based on violations
func (automation *DataReplicationEnforcementAutomation) triggerReplicationEnforcement(violations []common.ReplicationViolation) {
    for _, violation := range violations {
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting validator for replication enforcement.")
            continue
        }

        // Encrypt replication violation data before enforcement
        encryptedViolation := automation.AddEncryptionToViolationData(violation)

        fmt.Printf("Validator %s selected for enforcing data replication.\n", validator.Address)

        // Enforce data replication using the selected validator
        enforcementSuccess := automation.consensusSystem.EnforceReplicationPolicy(validator, encryptedViolation)
        if enforcementSuccess {
            fmt.Println("Data replication policy successfully enforced.")
        } else {
            fmt.Println("Error enforcing data replication policy.")
        }

        // Log the enforcement action into the ledger
        automation.logReplicationEnforcement(violation)
    }
}

// finalizeEnforcementCycle finalizes the replication enforcement cycle and logs the result into the ledger
func (automation *DataReplicationEnforcementAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeReplicationCycle()
    if success {
        fmt.Println("Replication enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing replication enforcement cycle.")
    }
}

// logReplicationEnforcement logs each replication enforcement action into the ledger
func (automation *DataReplicationEnforcementAutomation) logReplicationEnforcement(violation common.ReplicationViolation) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("replication-enforcement-%s", violation.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Replication Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with replication enforcement action for ViolationID %s.\n", violation.ID)
}

// logEnforcementCycleFinalization logs the finalization of a replication enforcement cycle into the ledger
func (automation *DataReplicationEnforcementAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("replication-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Replication Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with replication enforcement cycle finalization.")
}

// AddEncryptionToViolationData encrypts replication violation data before enforcement
func (automation *DataReplicationEnforcementAutomation) AddEncryptionToViolationData(violation common.ReplicationViolation) common.ReplicationViolation {
    encryptedData, err := encryption.EncryptData(violation.Data)
    if err != nil {
        fmt.Println("Error encrypting replication violation data:", err)
        return violation
    }
    violation.Data = encryptedData
    fmt.Println("Replication violation data successfully encrypted.")
    return violation
}

// ensureReplicationIntegrity checks the integrity of data replication and triggers enforcement if necessary
func (automation *DataReplicationEnforcementAutomation) ensureReplicationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateReplicationIntegrity()
    if !integrityValid {
        fmt.Println("Replication integrity breach detected. Triggering enforcement.")
        automation.triggerReplicationEnforcement(automation.consensusSystem.CheckReplicationViolations())
    } else {
        fmt.Println("Replication integrity is valid.")
    }
}
