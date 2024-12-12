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
    PowerAdjustmentCheckInterval     = 2000 * time.Millisecond // Interval for checking node power usage
    MaxPowerAdjustmentViolationLimit = 5                       // Maximum power usage violations before triggering adjustments
    SubBlocksPerBlock                = 1000                    // Number of sub-blocks in a block
    MaxAllowedNodePowerUsage         = 1000                    // Maximum allowed node power usage in watts
)

// DynamicNodePowerAdjustmentAutomation automates dynamic power adjustments for nodes to meet global warming reduction goals
type DynamicNodePowerAdjustmentAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store power adjustment-related logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    powerViolationCount    int                          // Counter for power usage violations
}

// NewDynamicNodePowerAdjustmentAutomation initializes the automation for dynamic node power adjustments
func NewDynamicNodePowerAdjustmentAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DynamicNodePowerAdjustmentAutomation {
    return &DynamicNodePowerAdjustmentAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        powerViolationCount: 0,
    }
}

// StartPowerAdjustmentAutomation starts the continuous loop for monitoring and adjusting node power usage
func (automation *DynamicNodePowerAdjustmentAutomation) StartPowerAdjustmentAutomation() {
    ticker := time.NewTicker(PowerAdjustmentCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndAdjustNodePower()
        }
    }()
}

// monitorAndAdjustNodePower checks node power usage and triggers dynamic adjustments if usage exceeds allowed limits
func (automation *DynamicNodePowerAdjustmentAutomation) monitorAndAdjustNodePower() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch current node power usage from the consensus system
    powerViolations := automation.consensusSystem.CheckNodePowerUsage()

    if len(powerViolations) >= MaxPowerAdjustmentViolationLimit {
        fmt.Printf("Node power usage violations exceed limit (%d). Triggering power adjustment.\n", len(powerViolations))
        automation.triggerPowerAdjustment(powerViolations)
    } else {
        fmt.Printf("Node power usage within acceptable range (%d violations).\n", len(powerViolations))
    }

    automation.powerViolationCount++
    fmt.Printf("Node power adjustment cycle #%d executed.\n", automation.powerViolationCount)

    if automation.powerViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizePowerAdjustmentCycle()
    }
}

// triggerPowerAdjustment triggers the dynamic adjustment of node power usage when violations occur
func (automation *DynamicNodePowerAdjustmentAutomation) triggerPowerAdjustment(violations []common.PowerViolation) {
    for _, violation := range violations {
        validator := automation.consensusSystem.SelectValidatorForPowerAdjustment()
        if validator == nil {
            fmt.Println("Error selecting validator for power adjustment.")
            continue
        }

        // Encrypt the node power violation data before adjustment
        encryptedViolation := automation.AddEncryptionToPowerViolationData(violation)

        fmt.Printf("Validator %s selected for adjusting node power usage using Synnergy Consensus.\n", validator.Address)

        // Adjust the node's power usage dynamically using the selected validator
        adjustmentSuccess := automation.consensusSystem.AdjustNodePower(validator, encryptedViolation)
        if adjustmentSuccess {
            fmt.Println("Node power adjustment successfully enforced.")
        } else {
            fmt.Println("Error enforcing node power adjustment.")
        }

        // Log the power adjustment action into the ledger
        automation.logPowerAdjustment(violation)
    }
}

// finalizePowerAdjustmentCycle finalizes the power adjustment cycle and logs the result in the ledger
func (automation *DynamicNodePowerAdjustmentAutomation) finalizePowerAdjustmentCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizePowerAdjustmentCycle()
    if success {
        fmt.Println("Node power adjustment cycle finalized successfully.")
        automation.logPowerAdjustmentCycleFinalization()
    } else {
        fmt.Println("Error finalizing node power adjustment cycle.")
    }
}

// logPowerAdjustment logs each node power adjustment action into the ledger for traceability
func (automation *DynamicNodePowerAdjustmentAutomation) logPowerAdjustment(violation common.PowerViolation) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-power-adjustment-%s", violation.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Power Adjustment",
        Status:    "Adjusted",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with node power adjustment action for ViolationID %s.\n", violation.ID)
}

// logPowerAdjustmentCycleFinalization logs the finalization of a power adjustment cycle into the ledger
func (automation *DynamicNodePowerAdjustmentAutomation) logPowerAdjustmentCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("power-adjustment-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Power Adjustment Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with power adjustment cycle finalization.")
}

// AddEncryptionToPowerViolationData encrypts node power violation data before adjustment
func (automation *DynamicNodePowerAdjustmentAutomation) AddEncryptionToPowerViolationData(violation common.PowerViolation) common.PowerViolation {
    encryptedData, err := encryption.EncryptData(violation.Data)
    if err != nil {
        fmt.Println("Error encrypting node power violation data:", err)
        return violation
    }
    violation.Data = encryptedData
    fmt.Println("Node power violation data successfully encrypted.")
    return violation
}

// ensureNodePowerIntegrity checks the integrity of node power usage data and triggers adjustments if necessary
func (automation *DynamicNodePowerAdjustmentAutomation) ensureNodePowerIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidatePowerUsageIntegrity()
    if !integrityValid {
        fmt.Println("Node power integrity breach detected. Triggering power adjustment.")
        automation.triggerPowerAdjustment(automation.consensusSystem.CheckNodePowerUsage())
    } else {
        fmt.Println("Node power usage integrity is valid.")
    }
}
