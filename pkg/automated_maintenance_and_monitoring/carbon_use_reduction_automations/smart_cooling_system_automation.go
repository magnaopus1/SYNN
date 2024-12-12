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
    TemperatureCheckInterval       = 1500 * time.Millisecond // Interval for checking node temperatures
    MaxAllowedTemperature          = 75                      // Maximum allowed temperature in degrees Celsius
    SubBlocksPerBlock              = 1000                    // Number of sub-blocks in a block
    CoolingThresholdTemperature    = 70                      // Threshold temperature to trigger cooling system
)

// SmartCoolingSystemAutomation automates the cooling system for nodes to maintain optimal temperatures
type SmartCoolingSystemAutomation struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger to store cooling system actions
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    temperatureViolationCount int                        // Counter for temperature violations
}

// NewSmartCoolingSystemAutomation initializes the automation for smart cooling of nodes and validators
func NewSmartCoolingSystemAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartCoolingSystemAutomation {
    return &SmartCoolingSystemAutomation{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        temperatureViolationCount: 0,
    }
}

// StartSmartCoolingSystemAutomation starts the continuous loop for monitoring node temperatures
func (automation *SmartCoolingSystemAutomation) StartSmartCoolingSystemAutomation() {
    ticker := time.NewTicker(TemperatureCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndControlTemperature()
        }
    }()
}

// monitorAndControlTemperature checks the temperatures of nodes and triggers the cooling system if required
func (automation *SmartCoolingSystemAutomation) monitorAndControlTemperature() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch temperature data for all nodes and validators from the Synnergy Consensus
    temperatureData := automation.consensusSystem.GetNodeTemperatureData()

    for _, node := range temperatureData {
        if node.Temperature > MaxAllowedTemperature {
            fmt.Printf("Node %s exceeds temperature limit with %d°C. Triggering emergency cooling action.\n", node.Address, node.Temperature)
            automation.triggerEmergencyCooling(node)
        } else if node.Temperature > CoolingThresholdTemperature {
            fmt.Printf("Node %s exceeds threshold temperature with %d°C. Activating smart cooling system.\n", node.Address, node.Temperature)
            automation.activateSmartCooling(node)
        } else {
            fmt.Printf("Node %s is operating at a safe temperature (%d°C).\n", node.Address, node.Temperature)
        }
    }

    automation.temperatureViolationCount++
    fmt.Printf("Temperature monitoring cycle #%d executed.\n", automation.temperatureViolationCount)

    if automation.temperatureViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeCoolingCycle()
    }
}

// activateSmartCooling activates the smart cooling system for nodes that exceed the threshold temperature
func (automation *SmartCoolingSystemAutomation) activateSmartCooling(node common.NodeTemperatureData) {
    // Encrypt the node's temperature data before triggering the cooling system
    encryptedData := automation.AddEncryptionToTemperatureData(node)

    // Trigger the smart cooling system via the Synnergy Consensus
    coolingSuccess := automation.consensusSystem.ActivateSmartCooling(encryptedData)
    if coolingSuccess {
        fmt.Printf("Smart cooling system activated for node %s.\n", node.Address)
        automation.logCoolingAction(node, "Smart Cooling Activated")
    } else {
        fmt.Printf("Error activating smart cooling system for node %s.\n", node.Address)
    }
}

// triggerEmergencyCooling triggers emergency cooling actions for nodes that exceed the maximum allowed temperature
func (automation *SmartCoolingSystemAutomation) triggerEmergencyCooling(node common.NodeTemperatureData) {
    // Encrypt the node's temperature data before triggering emergency cooling actions
    encryptedData := automation.AddEncryptionToTemperatureData(node)

    // Trigger emergency cooling via the Synnergy Consensus
    coolingSuccess := automation.consensusSystem.TriggerEmergencyCooling(encryptedData)
    if coolingSuccess {
        fmt.Printf("Emergency cooling system activated for node %s.\n", node.Address)
        automation.logCoolingAction(node, "Emergency Cooling Activated")
    } else {
        fmt.Printf("Error activating emergency cooling system for node %s.\n", node.Address)
    }
}

// finalizeCoolingCycle finalizes the smart cooling cycle and logs the result in the ledger
func (automation *SmartCoolingSystemAutomation) finalizeCoolingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCoolingCycle()
    if success {
        fmt.Println("Smart cooling cycle finalized successfully.")
        automation.logCoolingCycleFinalization()
    } else {
        fmt.Println("Error finalizing smart cooling cycle.")
    }
}

// logCoolingAction logs each cooling action into the ledger for traceability
func (automation *SmartCoolingSystemAutomation) logCoolingAction(node common.NodeTemperatureData, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cooling-action-%s-%s", node.Address, action),
        Timestamp: time.Now().Unix(),
        Type:      "Cooling System Action",
        Status:    action,
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with cooling system action for NodeID %s: %s.\n", node.Address, action)
}

// logCoolingCycleFinalization logs the finalization of a smart cooling cycle into the ledger
func (automation *SmartCoolingSystemAutomation) logCoolingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("cooling-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cooling Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with smart cooling system cycle finalization.")
}

// AddEncryptionToTemperatureData encrypts the node's temperature data before triggering any cooling actions
func (automation *SmartCoolingSystemAutomation) AddEncryptionToTemperatureData(node common.NodeTemperatureData) common.NodeTemperatureData {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node temperature data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node temperature data successfully encrypted.")
    return node
}

// ensureTemperatureDataIntegrity checks the integrity of temperature data and triggers cooling actions if necessary
func (automation *SmartCoolingSystemAutomation) ensureTemperatureDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateTemperatureDataIntegrity()
    if !integrityValid {
        fmt.Println("Temperature data integrity breach detected. Triggering smart cooling system.")
        automation.monitorAndControlTemperature()
    } else {
        fmt.Println("Temperature data integrity is valid.")
    }
}
