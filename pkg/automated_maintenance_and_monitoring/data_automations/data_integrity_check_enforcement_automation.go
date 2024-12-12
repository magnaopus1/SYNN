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
    IntegrityCheckInterval        = 800 * time.Millisecond // Interval for performing integrity checks
    MaxIntegrityViolationLimit    = 10                     // Max number of integrity violations allowed before triggering enforcement
)

// DataIntegrityCheckEnforcementAutomation monitors and enforces data integrity within the blockchain network
type DataIntegrityCheckEnforcementAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store integrity-related data
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    integrityViolationCount int                        // Counter for data integrity violations
}

// NewDataIntegrityCheckEnforcementAutomation initializes the automation for data integrity checks
func NewDataIntegrityCheckEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataIntegrityCheckEnforcementAutomation {
    return &DataIntegrityCheckEnforcementAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        integrityViolationCount: 0,
    }
}

// StartIntegrityCheckAutomation starts the continuous loop for monitoring and enforcing data integrity
func (automation *DataIntegrityCheckEnforcementAutomation) StartIntegrityCheckAutomation() {
    ticker := time.NewTicker(IntegrityCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkDataIntegrity()
        }
    }()
}

// checkDataIntegrity performs data integrity checks and triggers enforcement if necessary
func (automation *DataIntegrityCheckEnforcementAutomation) checkDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current integrity status from the consensus system
    integrityStatus := automation.consensusSystem.ValidateDataIntegrity()

    if !integrityStatus {
        fmt.Println("Data integrity check failed. Triggering enforcement.")
        automation.triggerIntegrityEnforcement()
    } else {
        fmt.Println("Data integrity check passed.")
    }

    automation.integrityViolationCount++
    fmt.Printf("Data integrity check cycle #%d executed.\n", automation.integrityViolationCount)

    // After a set number of cycles, finalize the integrity enforcement process
    if automation.integrityViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerIntegrityEnforcement takes action to resolve any detected integrity violations
func (automation *DataIntegrityCheckEnforcementAutomation) triggerIntegrityEnforcement() {
    validator := automation.consensusSystem.PoS.SelectValidator()
    if validator == nil {
        fmt.Println("Error selecting validator for data integrity enforcement.")
        return
    }

    // Perform encryption before initiating the enforcement action
    encryptedData := automation.AddEncryptionToIntegrityData()

    fmt.Printf("Validator %s selected for enforcing data integrity.\n", validator.Address)

    // Enforce data integrity using the selected validator
    enforcementSuccess := automation.consensusSystem.EnforceDataIntegrity(validator, encryptedData)
    if enforcementSuccess {
        fmt.Println("Data integrity successfully enforced.")
    } else {
        fmt.Println("Error enforcing data integrity.")
    }

    // Log the enforcement action into the ledger
    automation.logIntegrityEnforcement()
}

// finalizeEnforcementCycle finalizes an integrity enforcement cycle and logs the result into the ledger
func (automation *DataIntegrityCheckEnforcementAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeIntegrityCycle()
    if success {
        fmt.Println("Data integrity enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing data integrity enforcement cycle.")
    }
}

// logIntegrityEnforcement logs every integrity enforcement action into the ledger
func (automation *DataIntegrityCheckEnforcementAutomation) logIntegrityEnforcement() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("integrity-enforcement-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Data Integrity Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with data integrity enforcement action.")
}

// logEnforcementCycleFinalization logs the finalization of an integrity enforcement cycle into the ledger
func (automation *DataIntegrityCheckEnforcementAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("integrity-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Integrity Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with integrity enforcement cycle finalization.")
}

// AddEncryptionToIntegrityData encrypts any sensitive data related to integrity checks
func (automation *DataIntegrityCheckEnforcementAutomation) AddEncryptionToIntegrityData() []byte {
    dataToEncrypt := []byte("Data integrity details or sensitive information") // Example data to encrypt
    encryptedData, err := encryption.EncryptData(dataToEncrypt)
    if err != nil {
        fmt.Println("Error encrypting integrity data:", err)
        return nil
    }
    fmt.Println("Integrity data successfully encrypted.")
    return encryptedData
}

// ensureDataIntegrity ensures that the integrity of the data is valid and triggers enforcement if necessary
func (automation *DataIntegrityCheckEnforcementAutomation) ensureDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateDataIntegrity()
    if !integrityValid {
        fmt.Println("Data integrity breach detected. Triggering enforcement.")
        automation.triggerIntegrityEnforcement()
    } else {
        fmt.Println("Data integrity is valid.")
    }
}
