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
    AtomicSwapCheckInterval = 2 * time.Minute // Interval for checking atomic swaps
    MaxRetryCount           = 3               // Maximum number of retries for atomic swaps
    SubBlocksPerBlock       = 1000            // Number of sub-blocks in a block
)

// AtomicSwapSecurityProtocolAutomation automates the security process for atomic swaps between chains
type AtomicSwapSecurityProtocolAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging atomic swap events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    swapRetryCount    map[string]int               // Retry counter for failed atomic swaps
    swapCycleCount    int                          // Counter for atomic swap monitoring cycles
}

// NewAtomicSwapSecurityProtocolAutomation initializes the automation for atomic swap security
func NewAtomicSwapSecurityProtocolAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AtomicSwapSecurityProtocolAutomation {
    return &AtomicSwapSecurityProtocolAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        swapRetryCount:   make(map[string]int),
        swapCycleCount:   0,
    }
}

// StartAtomicSwapMonitoring starts the continuous loop for monitoring and enforcing atomic swap security
func (automation *AtomicSwapSecurityProtocolAutomation) StartAtomicSwapMonitoring() {
    ticker := time.NewTicker(AtomicSwapCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceAtomicSwaps()
        }
    }()
}

// monitorAndEnforceAtomicSwaps checks for atomic swap actions and enforces security protocols
func (automation *AtomicSwapSecurityProtocolAutomation) monitorAndEnforceAtomicSwaps() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of ongoing atomic swaps from the consensus system
    ongoingSwaps := automation.consensusSystem.GetOngoingAtomicSwaps()

    if len(ongoingSwaps) > 0 {
        for _, swap := range ongoingSwaps {
            fmt.Printf("Monitoring atomic swap %s.\n", swap.ID)
            if swap.Status == "Pending" {
                automation.enforceAtomicSwap(swap)
            } else if swap.Status == "Failed" {
                automation.retryAtomicSwap(swap)
            }
        }
    } else {
        fmt.Println("No ongoing atomic swaps to monitor at this time.")
    }

    automation.swapCycleCount++
    fmt.Printf("Atomic swap monitoring cycle #%d executed.\n", automation.swapCycleCount)

    if automation.swapCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeSwapCycle()
    }
}

// enforceAtomicSwap enforces the security protocols for atomic swaps
func (automation *AtomicSwapSecurityProtocolAutomation) enforceAtomicSwap(swap common.AtomicSwap) {
    encryptedSwapData := automation.encryptSwapData(swap)

    swapSuccess := automation.consensusSystem.EnforceAtomicSwap(swap, encryptedSwapData)

    if swapSuccess {
        fmt.Printf("Atomic swap %s enforced successfully.\n", swap.ID)
        automation.logSwapEvent(swap, "Enforced")
    } else {
        fmt.Printf("Error enforcing atomic swap %s. Retrying...\n", swap.ID)
        automation.retryAtomicSwap(swap)
    }
}

// retryAtomicSwap retries a failed atomic swap up to a maximum number of attempts
func (automation *AtomicSwapSecurityProtocolAutomation) retryAtomicSwap(swap common.AtomicSwap) {
    automation.swapRetryCount[swap.ID]++
    if automation.swapRetryCount[swap.ID] < MaxRetryCount {
        fmt.Printf("Retrying atomic swap %s (attempt %d).\n", swap.ID, automation.swapRetryCount[swap.ID])
        automation.enforceAtomicSwap(swap)
    } else {
        fmt.Printf("Max retries reached for atomic swap %s. Swap failed.\n", swap.ID)
        automation.logSwapFailure(swap)
    }
}

// finalizeSwapCycle finalizes the atomic swap monitoring cycle and logs the result in the ledger
func (automation *AtomicSwapSecurityProtocolAutomation) finalizeSwapCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeSwapCycle()
    if success {
        fmt.Println("Atomic swap monitoring cycle finalized successfully.")
        automation.logSwapCycleFinalization()
    } else {
        fmt.Println("Error finalizing atomic swap monitoring cycle.")
    }
}

// logSwapEvent logs atomic swap events into the ledger
func (automation *AtomicSwapSecurityProtocolAutomation) logSwapEvent(swap common.AtomicSwap, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("atomic-swap-%s-%s", swap.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Atomic Swap",
        Status:    eventType,
        Details:   fmt.Sprintf("Atomic swap %s %s successfully.", swap.ID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with atomic swap %s event for swap %s.\n", eventType, swap.ID)
}

// logSwapFailure logs the failure of an atomic swap into the ledger
func (automation *AtomicSwapSecurityProtocolAutomation) logSwapFailure(swap common.AtomicSwap) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("atomic-swap-failure-%s", swap.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Atomic Swap Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Atomic swap %s failed after maximum retries.", swap.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with atomic swap failure for swap %s.\n", swap.ID)
}

// logSwapCycleFinalization logs the finalization of an atomic swap monitoring cycle into the ledger
func (automation *AtomicSwapSecurityProtocolAutomation) logSwapCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("atomic-swap-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Atomic Swap Monitoring Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with atomic swap monitoring cycle finalization.")
}

// encryptSwapData encrypts the atomic swap data for secure enforcement
func (automation *AtomicSwapSecurityProtocolAutomation) encryptSwapData(swap common.AtomicSwap) common.AtomicSwap {
    encryptedData, err := encryption.EncryptData(swap.Data)
    if err != nil {
        fmt.Println("Error encrypting atomic swap data:", err)
        return swap
    }

    swap.EncryptedData = encryptedData
    fmt.Println("Atomic swap data successfully encrypted.")
    return swap
}

// emergencyManualIntervention allows for emergency manual intervention in case of atomic swap failure
func (automation *AtomicSwapSecurityProtocolAutomation) emergencyManualIntervention(swap common.AtomicSwap, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if action == "enforce" {
        fmt.Printf("Manually enforcing atomic swap %s.\n", swap.ID)
        automation.enforceAtomicSwap(swap)
    } else if action == "cancel" {
        fmt.Printf("Manually cancelling atomic swap %s.\n", swap.ID)
        automation.cancelAtomicSwap(swap)
    } else {
        fmt.Println("Invalid action for manual intervention.")
    }
}

// cancelAtomicSwap cancels an atomic swap due to security risks or manual intervention
func (automation *AtomicSwapSecurityProtocolAutomation) cancelAtomicSwap(swap common.AtomicSwap) {
    success := automation.consensusSystem.CancelAtomicSwap(swap)

    if success {
        fmt.Printf("Atomic swap %s cancelled successfully.\n", swap.ID)
        automation.logSwapEvent(swap, "Cancelled")
    } else {
        fmt.Printf("Error cancelling atomic swap %s.\n", swap.ID)
    }
}
