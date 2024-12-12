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
    ScriptReplicationCheckInterval = 2000 * time.Millisecond // Interval for checking script replication across nodes
    SubBlocksPerBlock              = 1000                    // Number of sub-blocks in a block
)

// ScriptRunningReplicationAutomation automates the replication of scripts across all nodes, except authority nodes, to ensure high availability
type ScriptRunningReplicationAutomation struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger to store script replication events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    replicationCheckCount    int                          // Counter for replication check cycles
}

// NewScriptRunningReplicationAutomation initializes the automation for script replication across nodes
func NewScriptRunningReplicationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ScriptRunningReplicationAutomation {
    return &ScriptRunningReplicationAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        replicationCheckCount: 0,
    }
}

// StartScriptReplicationCheck starts the continuous loop for monitoring and enforcing script replication across all non-authority nodes
func (automation *ScriptRunningReplicationAutomation) StartScriptReplicationCheck() {
    ticker := time.NewTicker(ScriptReplicationCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndReplicateScripts()
        }
    }()
}

// monitorAndReplicateScripts checks all running scripts throughout the system and ensures they are replicated across all non-authority nodes
func (automation *ScriptRunningReplicationAutomation) monitorAndReplicateScripts() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch running scripts across the system from Synnergy Consensus
    runningScripts := automation.consensusSystem.GetRunningScriptsAcrossNodes()

    for _, script := range runningScripts {
        if !script.IsAuthorityNodeSpecific { // Ensure we don't replicate authority node-specific scripts
            if automation.checkScriptReplication(script) {
                fmt.Printf("Script %s is successfully running on all non-authority nodes.\n", script.ID)
            } else {
                fmt.Printf("Script %s is not running on all nodes. Replicating to remaining nodes.\n", script.ID)
                automation.replicateScriptToNodes(script)
            }
        } else {
            fmt.Printf("Skipping replication for authority node-specific script %s.\n", script.ID)
        }
    }

    automation.replicationCheckCount++
    fmt.Printf("Script replication check cycle #%d executed.\n", automation.replicationCheckCount)

    if automation.replicationCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeReplicationCycle()
    }
}

// checkScriptReplication checks if the script is running on all non-authority nodes in the network
func (automation *ScriptRunningReplicationAutomation) checkScriptReplication(script common.Script) bool {
    return automation.consensusSystem.IsScriptRunningOnAllNonAuthorityNodes(script.ID)
}

// replicateScriptToNodes replicates a script to any non-authority node that is not currently running the script
func (automation *ScriptRunningReplicationAutomation) replicateScriptToNodes(script common.Script) {
    // Encrypt the script data before replicating
    encryptedScript := automation.AddEncryptionToScriptData(script)

    // Replicate the script to all non-authority nodes via the Synnergy Consensus
    replicationSuccess := automation.consensusSystem.ReplicateScriptAcrossNonAuthorityNodes(encryptedScript)

    if replicationSuccess {
        fmt.Printf("Script %s successfully replicated across all non-authority nodes.\n", script.ID)
        automation.logScriptReplication(script)
    } else {
        fmt.Printf("Error replicating script %s to all non-authority nodes.\n", script.ID)
    }
}

// finalizeReplicationCycle finalizes the script replication check cycle and logs the result in the ledger
func (automation *ScriptRunningReplicationAutomation) finalizeReplicationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeScriptReplicationCycle()
    if success {
        fmt.Println("Script replication check cycle finalized successfully.")
        automation.logReplicationCycleFinalization()
    } else {
        fmt.Println("Error finalizing script replication check cycle.")
    }
}

// logScriptReplication logs each script replication event into the ledger for traceability
func (automation *ScriptRunningReplicationAutomation) logScriptReplication(script common.Script) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("script-replication-%s", script.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Script Replication",
        Status:    "Replicated",
        Details:   fmt.Sprintf("Script %s successfully replicated across all non-authority nodes.", script.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with script replication event for script %s.\n", script.ID)
}

// logReplicationCycleFinalization logs the finalization of a script replication check cycle into the ledger
func (automation *ScriptRunningReplicationAutomation) logReplicationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("script-replication-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Script Replication Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with script replication cycle finalization.")
}

// AddEncryptionToScriptData encrypts the script data before replication
func (automation *ScriptRunningReplicationAutomation) AddEncryptionToScriptData(script common.Script) common.Script {
    encryptedData, err := encryption.EncryptData(script)
    if err != nil {
        fmt.Println("Error encrypting script data:", err)
        return script
    }
    script.EncryptedData = encryptedData
    fmt.Println("Script data successfully encrypted.")
    return script
}

// ensureScriptReplicationIntegrity checks the integrity of script replication data and triggers replication if necessary
func (automation *ScriptRunningReplicationAutomation) ensureScriptReplicationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateScriptReplicationIntegrity()
    if !integrityValid {
        fmt.Println("Script replication data integrity breach detected. Re-triggering script replication.")
        automation.monitorAndReplicateScripts()
    } else {
        fmt.Println("Script replication data integrity is valid.")
    }
}
