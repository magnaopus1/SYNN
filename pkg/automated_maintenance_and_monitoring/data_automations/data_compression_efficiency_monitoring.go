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
    CompressionEfficiencyCheckInterval = 900 * time.Millisecond // Interval for checking compression efficiency
    MinCompressionEfficiencyThreshold  = 85                     // Minimum percentage efficiency for data compression
    MaxCompressionViolationLimit       = 5                      // Max number of compression violations allowed before triggering enforcement
)

// DataCompressionEfficiencyMonitoringAutomation monitors the data compression efficiency and triggers corrective actions if the efficiency drops below a certain threshold
type DataCompressionEfficiencyMonitoringAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store compression-related data
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    compressionViolationCount int                        // Counter for compression violations
}

// NewDataCompressionEfficiencyMonitoringAutomation initializes the automation for data compression efficiency monitoring
func NewDataCompressionEfficiencyMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataCompressionEfficiencyMonitoringAutomation {
    return &DataCompressionEfficiencyMonitoringAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        compressionViolationCount: 0,
    }
}

// StartCompressionEfficiencyMonitoring starts the continuous loop for monitoring data compression efficiency
func (automation *DataCompressionEfficiencyMonitoringAutomation) StartCompressionEfficiencyMonitoring() {
    ticker := time.NewTicker(CompressionEfficiencyCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkCompressionEfficiency()
        }
    }()
}

// checkCompressionEfficiency checks the compression efficiency and triggers enforcement if necessary
func (automation *DataCompressionEfficiencyMonitoringAutomation) checkCompressionEfficiency() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current compression efficiency from the consensus system
    compressionEfficiency := automation.consensusSystem.GetCompressionEfficiency()

    if compressionEfficiency < MinCompressionEfficiencyThreshold {
        fmt.Printf("Compression efficiency below threshold (%d%%). Triggering enforcement actions.\n", compressionEfficiency)
        automation.triggerCompressionEnforcement(compressionEfficiency)
    } else {
        fmt.Printf("Compression efficiency is acceptable (%d%%).\n", compressionEfficiency)
    }

    automation.compressionViolationCount++
    fmt.Printf("Compression efficiency monitoring cycle #%d executed.\n", automation.compressionViolationCount)

    if automation.compressionViolationCount%100 == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerCompressionEnforcement takes action to enforce better data compression efficiency when violations are detected
func (automation *DataCompressionEfficiencyMonitoringAutomation) triggerCompressionEnforcement(compressionEfficiency int) {
    validator := automation.consensusSystem.PoS.SelectValidator()
    if validator == nil {
        fmt.Println("Error selecting validator for compression enforcement.")
        return
    }

    // Encrypt data before enforcement
    encryptedEfficiency := automation.AddEncryptionToEfficiencyData(compressionEfficiency)

    fmt.Printf("Validator %s selected for enforcing compression efficiency improvements.\n", validator.Address)

    // Enforce improved compression via the selected validator
    enforcementSuccess := automation.consensusSystem.EnforceCompressionEfficiency(validator, encryptedEfficiency)
    if enforcementSuccess {
        fmt.Println("Compression efficiency successfully enforced.")
    } else {
        fmt.Println("Error enforcing compression efficiency.")
    }

    // Log the enforcement action into the ledger
    automation.logCompressionEnforcement(compressionEfficiency)
}

// finalizeEnforcementCycle finalizes the current cycle of compression monitoring and logs the event into the ledger
func (automation *DataCompressionEfficiencyMonitoringAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCompressionCycle()
    if success {
        fmt.Println("Compression efficiency enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing compression efficiency enforcement cycle.")
    }
}

// logCompressionEnforcement logs every compression enforcement action into the ledger
func (automation *DataCompressionEfficiencyMonitoringAutomation) logCompressionEnforcement(compressionEfficiency int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-enforcement-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compression Enforcement",
        Status:    "Enforced",
        Details:   fmt.Sprintf("Compression efficiency was %d%%", compressionEfficiency),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression enforcement action (Efficiency: %d%%).\n", compressionEfficiency)
}

// logEnforcementCycleFinalization logs the finalization of a compression enforcement cycle into the ledger
func (automation *DataCompressionEfficiencyMonitoringAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression enforcement cycle finalization.\n")
}

// AddEncryptionToEfficiencyData encrypts compression efficiency data before enforcement
func (automation *DataCompressionEfficiencyMonitoringAutomation) AddEncryptionToEfficiencyData(compressionEfficiency int) []byte {
    encryptedData, err := encryption.EncryptData([]byte(fmt.Sprintf("%d", compressionEfficiency)))
    if err != nil {
        fmt.Println("Error encrypting compression efficiency data:", err)
        return nil
    }
    fmt.Println("Compression efficiency data successfully encrypted.")
    return encryptedData
}

// ensureCompressionIntegrity monitors for breaches in compression efficiency and triggers enforcement when necessary
func (automation *DataCompressionEfficiencyMonitoringAutomation) ensureCompressionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCompressionIntegrity()
    if !integrityValid {
        fmt.Println("Compression integrity breach detected. Triggering enforcement.")
        automation.triggerCompressionEnforcement(automation.consensusSystem.GetCompressionEfficiency())
    } else {
        fmt.Println("Compression efficiency integrity is valid.")
    }
}
