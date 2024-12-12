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
    CompressionMaintenanceInterval = 1000 * time.Millisecond // Interval for running compression maintenance checks
    MaxCompressionIssueLimit       = 5                       // Max number of issues allowed before triggering enforcement
    CompressionThreshold           = 90                      // Minimum percentage for acceptable compression efficiency
)

// DataCompressionMaintenanceAutomation automates the maintenance of data compression efficiency within the blockchain
type DataCompressionMaintenanceAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store compression-related data
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    compressionIssueCount int                          // Counter for compression issues detected
}

// NewDataCompressionMaintenanceAutomation initializes the automation for data compression maintenance
func NewDataCompressionMaintenanceAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataCompressionMaintenanceAutomation {
    return &DataCompressionMaintenanceAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        compressionIssueCount: 0,
    }
}

// StartCompressionMaintenanceAutomation starts the continuous loop for maintaining data compression
func (automation *DataCompressionMaintenanceAutomation) StartCompressionMaintenanceAutomation() {
    ticker := time.NewTicker(CompressionMaintenanceInterval)

    go func() {
        for range ticker.C {
            automation.runCompressionMaintenance()
        }
    }()
}

// runCompressionMaintenance checks compression efficiency and resolves issues when necessary
func (automation *DataCompressionMaintenanceAutomation) runCompressionMaintenance() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current compression efficiency from the consensus system
    compressionEfficiency := automation.consensusSystem.GetCompressionEfficiency()

    if compressionEfficiency < CompressionThreshold {
        fmt.Printf("Compression efficiency below acceptable level (%d%%). Triggering maintenance.\n", compressionEfficiency)
        automation.triggerCompressionMaintenance(compressionEfficiency)
    } else {
        fmt.Printf("Compression efficiency is acceptable (%d%%).\n", compressionEfficiency)
    }

    automation.compressionIssueCount++
    fmt.Printf("Compression maintenance cycle #%d executed.\n", automation.compressionIssueCount)

    if automation.compressionIssueCount%100 == 0 {
        automation.finalizeMaintenanceCycle()
    }
}

// triggerCompressionMaintenance performs the necessary actions to restore acceptable compression efficiency
func (automation *DataCompressionMaintenanceAutomation) triggerCompressionMaintenance(compressionEfficiency int) {
    validator := automation.consensusSystem.PoS.SelectValidator()
    if validator == nil {
        fmt.Println("Error selecting validator for compression maintenance.")
        return
    }

    // Encrypt the efficiency data before any actions
    encryptedEfficiency := automation.AddEncryptionToEfficiencyData(compressionEfficiency)

    fmt.Printf("Validator %s selected for compression maintenance.\n", validator.Address)

    // Execute the maintenance action to improve compression efficiency
    maintenanceSuccess := automation.consensusSystem.MaintainCompressionEfficiency(validator, encryptedEfficiency)
    if maintenanceSuccess {
        fmt.Println("Compression efficiency successfully maintained.")
    } else {
        fmt.Println("Error maintaining compression efficiency.")
    }

    // Log the maintenance action into the ledger
    automation.logCompressionMaintenance(compressionEfficiency)
}

// finalizeMaintenanceCycle finalizes and logs the completion of a maintenance cycle
func (automation *DataCompressionMaintenanceAutomation) finalizeMaintenanceCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCompressionCycle()
    if success {
        fmt.Println("Compression maintenance cycle finalized successfully.")
        automation.logMaintenanceCycleFinalization()
    } else {
        fmt.Println("Error finalizing compression maintenance cycle.")
    }
}

// logCompressionMaintenance logs every maintenance action into the ledger for accountability
func (automation *DataCompressionMaintenanceAutomation) logCompressionMaintenance(compressionEfficiency int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-maintenance-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compression Maintenance",
        Status:    "Maintained",
        Details:   fmt.Sprintf("Compression efficiency was %d%%", compressionEfficiency),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression maintenance action (Efficiency: %d%%).\n", compressionEfficiency)
}

// logMaintenanceCycleFinalization logs the finalization of a maintenance cycle into the ledger
func (automation *DataCompressionMaintenanceAutomation) logMaintenanceCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compression-maintenance-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Maintenance Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compression maintenance cycle finalization.\n")
}

// AddEncryptionToEfficiencyData encrypts compression efficiency data before maintenance
func (automation *DataCompressionMaintenanceAutomation) AddEncryptionToEfficiencyData(compressionEfficiency int) []byte {
    encryptedData, err := encryption.EncryptData([]byte(fmt.Sprintf("%d", compressionEfficiency)))
    if err != nil {
        fmt.Println("Error encrypting compression efficiency data:", err)
        return nil
    }
    fmt.Println("Compression efficiency data successfully encrypted.")
    return encryptedData
}

// ensureCompressionIntegrity checks for any breaches in compression efficiency and triggers maintenance if necessary
func (automation *DataCompressionMaintenanceAutomation) ensureCompressionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCompressionIntegrity()
    if !integrityValid {
        fmt.Println("Compression integrity breach detected. Triggering maintenance.")
        automation.triggerCompressionMaintenance(automation.consensusSystem.GetCompressionEfficiency())
    } else {
        fmt.Println("Compression efficiency integrity is valid.")
    }
}
