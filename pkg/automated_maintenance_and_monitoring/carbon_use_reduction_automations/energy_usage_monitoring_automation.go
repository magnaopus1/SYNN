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
    EnergyUsageCheckInterval    = 1800 * time.Millisecond // Interval for checking energy usage
    MaxAllowedEnergyViolation   = 1000                    // Maximum energy usage in watts before triggering action
    SubBlocksPerBlock           = 1000                    // Number of sub-blocks in a block
)

// EnergyUsageMonitoringAutomation automates the monitoring of energy usage across all nodes and validators
type EnergyUsageMonitoringAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store energy usage monitoring-related logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    energyViolationCount   int                          // Counter for energy usage violations
}

// NewEnergyUsageMonitoringAutomation initializes the automation for energy usage monitoring
func NewEnergyUsageMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *EnergyUsageMonitoringAutomation {
    return &EnergyUsageMonitoringAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        energyViolationCount: 0,
    }
}

// StartEnergyUsageMonitoring starts the continuous loop for monitoring energy usage across nodes and validators
func (automation *EnergyUsageMonitoringAutomation) StartEnergyUsageMonitoring() {
    ticker := time.NewTicker(EnergyUsageCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndCheckEnergyUsage()
        }
    }()
}

// monitorAndCheckEnergyUsage monitors the energy usage of nodes and validators and triggers actions if violations occur
func (automation *EnergyUsageMonitoringAutomation) monitorAndCheckEnergyUsage() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the energy usage data for all nodes and validators from the consensus system
    energyUsageData := automation.consensusSystem.GetEnergyUsageData()

    for _, node := range energyUsageData {
        if node.EnergyUsage > MaxAllowedEnergyViolation {
            fmt.Printf("Node %s exceeds energy usage limit with %d watts. Triggering corrective action.\n", node.Address, node.EnergyUsage)
            automation.triggerEnergyUsageViolation(node)
        } else {
            fmt.Printf("Node %s is within the allowed energy usage limit (%d watts).\n", node.Address, node.EnergyUsage)
        }
    }

    automation.energyViolationCount++
    fmt.Printf("Energy usage monitoring cycle #%d executed.\n", automation.energyViolationCount)

    if automation.energyViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeMonitoringCycle()
    }
}

// triggerEnergyUsageViolation triggers corrective actions for nodes or validators that exceed energy usage limits
func (automation *EnergyUsageMonitoringAutomation) triggerEnergyUsageViolation(node common.NodeEnergyData) {
    // Encrypt the node's energy data before triggering corrective actions
    encryptedData := automation.AddEncryptionToEnergyData(node)

    // Trigger corrective action using Synnergy Consensus to adjust the node's operations
    adjustmentSuccess := automation.consensusSystem.AdjustNodeEnergyUsage(encryptedData)
    if adjustmentSuccess {
        fmt.Printf("Energy usage adjusted for node %s.\n", node.Address)
        automation.logEnergyUsageAction(node)
    } else {
        fmt.Printf("Error adjusting energy usage for node %s.\n", node.Address)
    }
}

// finalizeMonitoringCycle finalizes the monitoring cycle and logs the result in the ledger
func (automation *EnergyUsageMonitoringAutomation) finalizeMonitoringCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEnergyMonitoringCycle()
    if success {
        fmt.Println("Energy usage monitoring cycle finalized successfully.")
        automation.logMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing energy usage monitoring cycle.")
    }
}

// logEnergyUsageAction logs each energy usage corrective action into the ledger for traceability
func (automation *EnergyUsageMonitoringAutomation) logEnergyUsageAction(node common.NodeEnergyData) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("energy-usage-correction-%s", node.Address),
        Timestamp: time.Now().Unix(),
        Type:      "Energy Usage Correction",
        Status:    "Corrected",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with energy usage correction action for NodeID %s.\n", node.Address)
}

// logMonitoringCycleFinalization logs the finalization of an energy usage monitoring cycle into the ledger
func (automation *EnergyUsageMonitoringAutomation) logMonitoringCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("energy-usage-monitoring-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Energy Usage Monitoring Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with energy usage monitoring cycle finalization.")
}

// AddEncryptionToEnergyData encrypts the node's energy data before triggering corrective actions
func (automation *EnergyUsageMonitoringAutomation) AddEncryptionToEnergyData(node common.NodeEnergyData) common.NodeEnergyData {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node energy data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node energy data successfully encrypted.")
    return node
}

// ensureEnergyDataIntegrity checks the integrity of energy usage data and ensures actions are triggered as required
func (automation *EnergyUsageMonitoringAutomation) ensureEnergyDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateEnergyDataIntegrity()
    if !integrityValid {
        fmt.Println("Energy usage data integrity breach detected. Triggering monitoring actions.")
        automation.monitorAndCheckEnergyUsage()
    } else {
        fmt.Println("Energy usage data integrity is valid.")
    }
}
