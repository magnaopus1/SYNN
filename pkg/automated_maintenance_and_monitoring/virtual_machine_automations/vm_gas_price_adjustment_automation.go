package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/transaction"
)

const (
    GasPriceAdjustmentInterval = 5000 * time.Millisecond // Interval for adjusting gas prices
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// VMGasPriceAdjustmentAutomation automates the process of adjusting gas prices on the virtual machine
type VMGasPriceAdjustmentAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance   *ledger.Ledger               // Ledger to store gas price adjustment actions
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    gasPriceAdjustmentCount int                   // Counter for gas price adjustment cycles
    currentGasPrice  int                          // Current gas price in the system
}

// NewVMGasPriceAdjustmentAutomation initializes the automation for gas price adjustment on the VM
func NewVMGasPriceAdjustmentAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMGasPriceAdjustmentAutomation {
    return &VMGasPriceAdjustmentAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        gasPriceAdjustmentCount: 0,
        currentGasPrice:        transaction.InitialGasPrice, // Starting gas price from the transaction package
    }
}

// StartGasPriceAdjustmentCheck starts the continuous loop for monitoring and adjusting gas prices
func (automation *VMGasPriceAdjustmentAutomation) StartGasPriceAdjustmentCheck() {
    ticker := time.NewTicker(GasPriceAdjustmentInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndAdjustGasPrice()
        }
    }()
}

// monitorAndAdjustGasPrice monitors network load and adjusts gas prices dynamically
func (automation *VMGasPriceAdjustmentAutomation) monitorAndAdjustGasPrice() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the network load and congestion statistics
    networkLoad := automation.consensusSystem.GetNetworkLoad()

    // Adjust gas price based on network load
    newGasPrice := automation.calculateAdjustedGasPrice(networkLoad)

    // Enforce fee ceiling and fee floor
    if newGasPrice > transaction.FeeCeiling {
        newGasPrice = transaction.FeeCeiling
        fmt.Println("Gas price adjusted to the fee ceiling.")
    } else if newGasPrice < transaction.FeeFloor {
        newGasPrice = transaction.FeeFloor
        fmt.Println("Gas price adjusted to the fee floor.")
    }

    if newGasPrice != automation.currentGasPrice {
        fmt.Printf("Gas price adjusted from %d to %d.\n", automation.currentGasPrice, newGasPrice)
        automation.currentGasPrice = newGasPrice
        automation.logGasPriceAdjustment(newGasPrice)
    }

    automation.gasPriceAdjustmentCount++
    fmt.Printf("Gas price adjustment cycle #%d executed.\n", automation.gasPriceAdjustmentCount)

    if automation.gasPriceAdjustmentCount%SubBlocksPerBlock == 0 {
        automation.finalizeGasPriceAdjustmentCycle()
    }
}

// calculateAdjustedGasPrice calculates the new gas price based on network load
func (automation *VMGasPriceAdjustmentAutomation) calculateAdjustedGasPrice(networkLoad int) int {
    // Basic adjustment based on network load
    adjustmentFactor := 1.0
    if networkLoad > common.HighNetworkLoadThreshold {
        adjustmentFactor = 1.2
    } else if networkLoad < common.LowNetworkLoadThreshold {
        adjustmentFactor = 0.8
    }

    return int(float64(automation.currentGasPrice) * adjustmentFactor)
}

// logGasPriceAdjustment logs the gas price adjustment into the ledger for traceability
func (automation *VMGasPriceAdjustmentAutomation) logGasPriceAdjustment(newGasPrice int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-price-adjustment-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Price Adjustment",
        Status:    "Adjusted",
        Details:   fmt.Sprintf("Gas price adjusted to %d units.", newGasPrice),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with new gas price: %d units.\n", newGasPrice)
}

// finalizeGasPriceAdjustmentCycle finalizes the gas price adjustment cycle and logs the result in the ledger
func (automation *VMGasPriceAdjustmentAutomation) finalizeGasPriceAdjustmentCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeGasPriceAdjustmentCycle()
    if success {
        fmt.Println("Gas price adjustment cycle finalized successfully.")
        automation.logGasPriceAdjustmentCycleFinalization()
    } else {
        fmt.Println("Error finalizing gas price adjustment cycle.")
    }
}

// logGasPriceAdjustmentCycleFinalization logs the finalization of a gas price adjustment cycle into the ledger
func (automation *VMGasPriceAdjustmentAutomation) logGasPriceAdjustmentCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-price-adjustment-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Price Adjustment Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with gas price adjustment cycle finalization.")
}

// ensureGasPriceIntegrity ensures the integrity of gas prices and triggers recalculation if necessary
func (automation *VMGasPriceAdjustmentAutomation) ensureGasPriceIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateGasPriceIntegrity()
    if !integrityValid {
        fmt.Println("Gas price integrity breach detected. Re-triggering gas price adjustments.")
        automation.monitorAndAdjustGasPrice()
    } else {
        fmt.Println("Gas price integrity is valid.")
    }
}
