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
    PropagationSpeedCheckInterval = 10 * time.Minute // Interval for checking block propagation speed
    MaxPropagationSpeedRetries    = 3                // Maximum retries for enforcing block propagation speed
    SubBlocksPerBlock             = 1000             // Number of sub-blocks in a block
    PropagationSpeedThreshold     = 200              // Block propagation speed threshold in milliseconds
)

// BlockPropagationSpeedEnforcement manages the enforcement of block propagation speed across the network
type BlockPropagationSpeedEnforcement struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger for logging propagation speed enforcement events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    propagationRetryCount map[string]int               // Counter for retrying propagation enforcement
    enforcementCycle      int                          // Counter for propagation enforcement cycles
}

// NewBlockPropagationSpeedEnforcement initializes the automation for block propagation speed enforcement
func NewBlockPropagationSpeedEnforcement(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *BlockPropagationSpeedEnforcement {
    return &BlockPropagationSpeedEnforcement{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        propagationRetryCount: make(map[string]int),
        enforcementCycle:      0,
    }
}

// StartBlockPropagationSpeedEnforcement starts the continuous loop for enforcing block propagation speed
func (automation *BlockPropagationSpeedEnforcement) StartBlockPropagationSpeedEnforcement() {
    ticker := time.NewTicker(PropagationSpeedCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforcePropagationSpeed()
        }
    }()
}

// monitorAndEnforcePropagationSpeed checks and enforces propagation speed limits across the network
func (automation *BlockPropagationSpeedEnforcement) monitorAndEnforcePropagationSpeed() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch reports on block propagation speed
    propagationReports := automation.consensusSystem.FetchPropagationReports()

    for _, report := range propagationReports {
        if report.PropagationSpeed > PropagationSpeedThreshold {
            fmt.Printf("Block %s exceeds propagation speed threshold with %d ms. Enforcing speed limit.\n", report.BlockID, report.PropagationSpeed)
            automation.enforcePropagationSpeed(report)
        } else {
            fmt.Printf("Block %s is within speed limits: %d ms.\n", report.BlockID, report.PropagationSpeed)
        }
    }

    automation.enforcementCycle++
    fmt.Printf("Block propagation speed enforcement cycle #%d completed.\n", automation.enforcementCycle)

    if automation.enforcementCycle%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// enforcePropagationSpeed handles the enforcement of block propagation speed
func (automation *BlockPropagationSpeedEnforcement) enforcePropagationSpeed(report common.PropagationReport) {
    encryptedData := automation.encryptPropagationData(report)

    // Attempt to enforce the propagation speed limit
    success := automation.consensusSystem.EnforcePropagationSpeed(report.BlockID, encryptedData)
    if success {
        fmt.Printf("Propagation speed enforcement successful for block %s.\n", report.BlockID)
        automation.logEnforcementEvent(report, "Enforced")
        automation.resetPropagationRetry(report.BlockID)
    } else {
        fmt.Printf("Propagation speed enforcement failed for block %s. Retrying...\n", report.BlockID)
        automation.retryPropagationEnforcement(report)
    }
}

// retryPropagationEnforcement retries the enforcement of block propagation speed if it fails
func (automation *BlockPropagationSpeedEnforcement) retryPropagationEnforcement(report common.PropagationReport) {
    automation.propagationRetryCount[report.BlockID]++
    if automation.propagationRetryCount[report.BlockID] < MaxPropagationSpeedRetries {
        automation.enforcePropagationSpeed(report)
    } else {
        fmt.Printf("Max retries reached for block %s. Propagation speed enforcement failed.\n", report.BlockID)
        automation.logEnforcementFailure(report)
    }
}

// resetPropagationRetry resets the retry count for propagation enforcement
func (automation *BlockPropagationSpeedEnforcement) resetPropagationRetry(blockID string) {
    automation.propagationRetryCount[blockID] = 0
}

// encryptPropagationData encrypts the block data before enforcing propagation speed
func (automation *BlockPropagationSpeedEnforcement) encryptPropagationData(report common.PropagationReport) []byte {
    encryptedData, err := encryption.EncryptData(report.Data)
    if err != nil {
        fmt.Printf("Error encrypting data for block %s: %v\n", report.BlockID, err)
        return report.Data
    }

    fmt.Printf("Block data successfully encrypted for block %s.\n", report.BlockID)
    return encryptedData
}

// logEnforcementEvent logs a successful propagation speed enforcement event into the ledger
func (automation *BlockPropagationSpeedEnforcement) logEnforcementEvent(report common.PropagationReport, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("propagation-enforcement-%s-%s", report.BlockID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Propagation Speed Enforcement",
        Status:    eventType,
        Details:   fmt.Sprintf("Block %s %s successfully.", report.BlockID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with propagation enforcement event for block %s.\n", report.BlockID)
}

// logEnforcementFailure logs the failure to enforce propagation speed for a specific block into the ledger
func (automation *BlockPropagationSpeedEnforcement) logEnforcementFailure(report common.PropagationReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("propagation-enforcement-failure-%s", report.BlockID),
        Timestamp: time.Now().Unix(),
        Type:      "Propagation Speed Enforcement Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Propagation speed enforcement failed for block %s after maximum retries.", report.BlockID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with propagation enforcement failure for block %s.\n", report.BlockID)
}

// finalizeEnforcementCycle finalizes the propagation enforcement cycle and logs the result in the ledger
func (automation *BlockPropagationSpeedEnforcement) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEnforcementCycle()
    if success {
        fmt.Println("Block propagation speed enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing block propagation enforcement cycle.")
    }
}

// logEnforcementCycleFinalization logs the finalization of an enforcement cycle into the ledger
func (automation *BlockPropagationSpeedEnforcement) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Propagation Speed Enforcement Cycle",
        Status:    "Finalized",
        Details:   "Block propagation speed enforcement cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with propagation enforcement cycle finalization.")
}
