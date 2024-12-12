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
    AccessControlEnforcementInterval = 800 * time.Millisecond // Interval for checking access control violations
    MaxAccessViolationLimit          = 10                     // Max access violations allowed before triggering enforcement
)

// DataAccessControlAutomation manages and enforces data access control policies in the blockchain network
type DataAccessControlAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store access control-related data
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    accessViolationCount int                         // Counter for access control violations
}

// NewDataAccessControlAutomation initializes the automation for data access control
func NewDataAccessControlAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataAccessControlAutomation {
    return &DataAccessControlAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        accessViolationCount: 0,
    }
}

// StartAccessControlAutomation starts the continuous loop for monitoring and enforcing access control policies
func (automation *DataAccessControlAutomation) StartAccessControlAutomation() {
    ticker := time.NewTicker(AccessControlEnforcementInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceAccessControl()
        }
    }()
}

// monitorAndEnforceAccessControl checks for access control violations and triggers enforcement if necessary
func (automation *DataAccessControlAutomation) monitorAndEnforceAccessControl() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch access violations from the consensus system
    accessViolations := automation.consensusSystem.GetAccessViolations()

    if len(accessViolations) >= MaxAccessViolationLimit {
        fmt.Printf("Access violations exceed limit (%d). Triggering enforcement actions.\n", len(accessViolations))
        automation.triggerAccessControlEnforcement(accessViolations)
    } else {
        fmt.Printf("Access violations are within acceptable range (%d).\n", len(accessViolations))
    }

    automation.accessViolationCount++
    fmt.Printf("Access control enforcement cycle #%d executed.\n", automation.accessViolationCount)

    if automation.accessViolationCount%100 == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerAccessControlEnforcement enforces access control rules based on detected violations
func (automation *DataAccessControlAutomation) triggerAccessControlEnforcement(violations []common.AccessViolation) {
    for _, violation := range violations {
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting validator for access control enforcement.")
            continue
        }

        // Encrypt sensitive data before enforcement
        encryptedViolation := automation.AddEncryptionToViolation(violation)

        fmt.Printf("Validator %s selected for enforcing access control.\n", validator.Address)

        // Enforce access control policy via consensus
        enforcementSuccess := automation.consensusSystem.EnforceAccessControlPolicy(validator, encryptedViolation)
        if enforcementSuccess {
            fmt.Println("Access control policy successfully enforced.")
        } else {
            fmt.Println("Error enforcing access control policy.")
        }

        // Log the enforcement action into the ledger
        automation.logAccessControlEnforcement(violation)
    }
}

// finalizeEnforcementCycle finalizes enforcement actions and updates the ledger
func (automation *DataAccessControlAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEnforcementCycle()
    if success {
        fmt.Println("Enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing enforcement cycle.")
    }
}

// logAccessControlEnforcement logs every access control enforcement action into the ledger
func (automation *DataAccessControlAutomation) logAccessControlEnforcement(violation common.AccessViolation) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("access-control-enforcement-%s", violation.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Access Control Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with access control enforcement event for ViolationID %s.\n", violation.ID)
}

// logEnforcementCycleFinalization logs the finalization of an enforcement cycle into the ledger
func (automation *DataAccessControlAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with enforcement cycle finalization.\n")
}

// AddEncryptionToViolation encrypts sensitive access violation data
func (automation *DataAccessControlAutomation) AddEncryptionToViolation(violation common.AccessViolation) common.AccessViolation {
    encryptedData, err := encryption.EncryptData(violation.Data)
    if err != nil {
        fmt.Println("Error encrypting access violation data:", err)
        return violation
    }
    violation.Data = encryptedData
    fmt.Println("Access violation data successfully encrypted.")
    return violation
}

// ensureAccessControlIntegrity monitors for breaches in access control and triggers enforcement when necessary
func (automation *DataAccessControlAutomation) ensureAccessControlIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateAccessControlIntegrity()
    if !integrityValid {
        fmt.Println("Access control integrity breach detected. Triggering enforcement.")
        automation.triggerAccessControlEnforcement(automation.consensusSystem.GetAccessViolations())
    } else {
        fmt.Println("Access control integrity is valid.")
    }
}
