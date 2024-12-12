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
    CarbonEmissionCheckInterval   = 1500 * time.Millisecond // Interval for checking carbon emission data
    MaxAllowedEmissionViolations  = 10                      // Maximum emission violations before triggering offset
    SubBlocksPerBlock             = 1000                    // Number of sub-blocks in a block
    EmissionOffsetThreshold       = 10000                   // Threshold for carbon emissions that require offset (in tons)
)

// CarbonEmissionOffsetAutomation automates the offset of carbon emissions on the blockchain
type CarbonEmissionOffsetAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store carbon offset-related logs
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    emissionViolationCount int                         // Counter for carbon emission violations
}

// NewCarbonEmissionOffsetAutomation initializes the automation for carbon emission offsetting
func NewCarbonEmissionOffsetAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CarbonEmissionOffsetAutomation {
    return &CarbonEmissionOffsetAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        emissionViolationCount: 0,
    }
}

// StartCarbonOffsetAutomation starts the continuous loop for monitoring and enforcing carbon emission offsets
func (automation *CarbonEmissionOffsetAutomation) StartCarbonOffsetAutomation() {
    ticker := time.NewTicker(CarbonEmissionCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndOffsetEmissions()
        }
    }()
}

// monitorAndOffsetEmissions checks carbon emissions data and triggers offsets when necessary
func (automation *CarbonEmissionOffsetAutomation) monitorAndOffsetEmissions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch carbon emission violations from the consensus system
    emissionViolations := automation.consensusSystem.CheckCarbonEmissions()

    if len(emissionViolations) >= MaxAllowedEmissionViolations {
        fmt.Printf("Carbon emissions exceed limit (%d violations). Triggering offset process.\n", len(emissionViolations))
        automation.triggerEmissionOffset(emissionViolations)
    } else {
        fmt.Printf("Carbon emissions within acceptable range (%d violations).\n", len(emissionViolations))
    }

    automation.emissionViolationCount++
    fmt.Printf("Carbon emission offset cycle #%d executed.\n", automation.emissionViolationCount)

    if automation.emissionViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeOffsetCycle()
    }
}

// triggerEmissionOffset triggers the carbon emission offset process when violations occur
func (automation *CarbonEmissionOffsetAutomation) triggerEmissionOffset(violations []common.EmissionViolation) {
    for _, violation := range violations {
        validator := automation.consensusSystem.SelectValidatorForOffset()
        if validator == nil {
            fmt.Println("Error selecting validator for emission offset.")
            continue
        }

        // Encrypt the carbon emission data before enforcing offset
        encryptedViolation := automation.AddEncryptionToEmissionData(violation)

        fmt.Printf("Validator %s selected for enforcing carbon emission offset using Synnergy Consensus.\n", validator.Address)

        // Enforce the carbon emission offset using the selected validator
        offsetSuccess := automation.consensusSystem.EnforceEmissionOffset(validator, encryptedViolation)
        if offsetSuccess {
            fmt.Println("Carbon emission offset successfully enforced.")
        } else {
            fmt.Println("Error enforcing carbon emission offset.")
        }

        // Log the emission offset action into the ledger
        automation.logEmissionOffset(violation)
    }
}

// finalizeOffsetCycle finalizes the offset cycle and logs it in the ledger
func (automation *CarbonEmissionOffsetAutomation) finalizeOffsetCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeOffsetCycle()
    if success {
        fmt.Println("Carbon emission offset cycle finalized successfully.")
        automation.logOffsetCycleFinalization()
    } else {
        fmt.Println("Error finalizing carbon emission offset cycle.")
    }
}

// logEmissionOffset logs each emission offset action into the ledger for traceability
func (automation *CarbonEmissionOffsetAutomation) logEmissionOffset(violation common.EmissionViolation) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("carbon-offset-%s", violation.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Carbon Emission Offset",
        Status:    "Offset",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with carbon emission offset action for ViolationID %s.\n", violation.ID)
}

// logOffsetCycleFinalization logs the finalization of a carbon emission offset cycle into the ledger
func (automation *CarbonEmissionOffsetAutomation) logOffsetCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("carbon-offset-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Offset Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with carbon emission offset cycle finalization.")
}

// AddEncryptionToEmissionData encrypts carbon emission violation data before enforcing offset
func (automation *CarbonEmissionOffsetAutomation) AddEncryptionToEmissionData(violation common.EmissionViolation) common.EmissionViolation {
    encryptedData, err := encryption.EncryptData(violation.Data)
    if err != nil {
        fmt.Println("Error encrypting carbon emission data:", err)
        return violation
    }
    violation.Data = encryptedData
    fmt.Println("Carbon emission data successfully encrypted.")
    return violation
}

// ensureEmissionIntegrity checks the integrity of carbon emission data and triggers offset when necessary
func (automation *CarbonEmissionOffsetAutomation) ensureEmissionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateEmissionIntegrity()
    if !integrityValid {
        fmt.Println("Carbon emission integrity breach detected. Triggering offset process.")
        automation.triggerEmissionOffset(automation.consensusSystem.CheckCarbonEmissions())
    } else {
        fmt.Println("Carbon emission data integrity is valid.")
    }
}
