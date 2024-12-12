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
    DistributionAdjustmentInterval = 10 * time.Minute // Interval for checking network load and adjusting distribution
    MaxDistributionRetries         = 3                // Maximum retry attempts for failed distribution adjustments
    SubBlocksPerBlock              = 1000             // Number of sub-blocks in a block
    MaxNetworkLoad                 = 85               // Max acceptable load percentage before adjustment
)

// AdaptiveDistributionAutomation adjusts the distribution of network resources based on load
type AdaptiveDistributionAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger for logging distribution adjustments
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    distributionRetryMap map[string]int               // Counter for retrying failed distribution operations
    distributionCycle    int                          // Counter for distribution adjustment cycles
}

// NewAdaptiveDistributionAutomation initializes the adaptive distribution automation
func NewAdaptiveDistributionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AdaptiveDistributionAutomation {
    return &AdaptiveDistributionAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        distributionRetryMap: make(map[string]int),
        distributionCycle:    0,
    }
}

// StartDistributionAdjustment starts the continuous loop for adjusting network distribution
func (automation *AdaptiveDistributionAutomation) StartDistributionAdjustment() {
    ticker := time.NewTicker(DistributionAdjustmentInterval)

    go func() {
        for range ticker.C {
            automation.checkAndAdjustDistribution()
        }
    }()
}

// checkAndAdjustDistribution checks the network load and adjusts distribution if necessary
func (automation *AdaptiveDistributionAutomation) checkAndAdjustDistribution() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    currentNetworkLoad := automation.consensusSystem.GetNetworkLoad()

    if currentNetworkLoad >= MaxNetworkLoad {
        fmt.Printf("Network load exceeds acceptable threshold: %d%%. Adjusting distribution.\n", currentNetworkLoad)
        automation.adjustDistribution(currentNetworkLoad)
    } else {
        fmt.Printf("Network load within acceptable limits: %d%%.\n", currentNetworkLoad)
    }

    automation.distributionCycle++
    fmt.Printf("Adaptive distribution adjustment cycle #%d completed.\n", automation.distributionCycle)

    if automation.distributionCycle%SubBlocksPerBlock == 0 {
        automation.finalizeDistributionCycle()
    }
}

// adjustDistribution handles the redistribution of resources across the network based on the load
func (automation *AdaptiveDistributionAutomation) adjustDistribution(currentNetworkLoad int) {
    err := automation.triggerDistributionAdjustment(currentNetworkLoad)
    if err != nil {
        fmt.Printf("Distribution adjustment failed due to network load: %v. Retrying...\n", err)
        automation.retryDistributionAdjustment(currentNetworkLoad)
        return
    }

    fmt.Printf("Distribution successfully adjusted for network load: %d%%.\n", currentNetworkLoad)
    automation.logDistributionEvent(currentNetworkLoad, "Adjusted", "High Load")
    automation.resetDistributionRetry()
}

// triggerDistributionAdjustment triggers the actual distribution adjustment based on the load
func (automation *AdaptiveDistributionAutomation) triggerDistributionAdjustment(currentNetworkLoad int) error {
    fmt.Println("Triggering distribution adjustment...")

    success := automation.consensusSystem.AdjustNetworkDistribution(currentNetworkLoad)
    if success {
        fmt.Println("Distribution adjustment successful.")
        return nil
    }

    return fmt.Errorf("distribution adjustment failed for network load: %d%%", currentNetworkLoad)
}

// retryDistributionAdjustment retries the distribution adjustment process in case of failure
func (automation *AdaptiveDistributionAutomation) retryDistributionAdjustment(currentNetworkLoad int) {
    automation.distributionRetryMap["distribution"]++
    if automation.distributionRetryMap["distribution"] < MaxDistributionRetries {
        automation.adjustDistribution(currentNetworkLoad)
    } else {
        fmt.Printf("Max retries reached for distribution adjustment. Adjustment failed for network load: %d%%.\n", currentNetworkLoad)
        automation.logDistributionFailure(currentNetworkLoad)
    }
}

// resetDistributionRetry resets the retry count for distribution adjustment operations
func (automation *AdaptiveDistributionAutomation) resetDistributionRetry() {
    automation.distributionRetryMap["distribution"] = 0
}

// logDistributionEvent logs a distribution-related event into the ledger
func (automation *AdaptiveDistributionAutomation) logDistributionEvent(networkLoad int, eventType, reason string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("distribution-event-%s-%d", eventType, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Network Distribution Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Network distribution %s due to %s. Network load: %d%%.", eventType, reason, networkLoad),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with distribution event: %s for network load: %d%%.\n", eventType, networkLoad)
}

// logDistributionFailure logs the failure to adjust distribution into the ledger
func (automation *AdaptiveDistributionAutomation) logDistributionFailure(networkLoad int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("distribution-failure-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Network Distribution Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to adjust distribution for network load: %d%% after max retries.", networkLoad),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with distribution adjustment failure for network load: %d%%.\n", networkLoad)
}

// finalizeDistributionCycle finalizes the distribution adjustment cycle and logs the result in the ledger
func (automation *AdaptiveDistributionAutomation) finalizeDistributionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDistributionCycle()
    if success {
        fmt.Println("Distribution adjustment cycle finalized successfully.")
        automation.logDistributionCycleFinalization()
    } else {
        fmt.Println("Error finalizing distribution adjustment cycle.")
    }
}

// logDistributionCycleFinalization logs the finalization of a distribution cycle into the ledger
func (automation *AdaptiveDistributionAutomation) logDistributionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("distribution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Distribution Cycle Finalization",
        Status:    "Finalized",
        Details:   "Network distribution cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with distribution cycle finalization.")
}
