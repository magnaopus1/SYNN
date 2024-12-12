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
    ComplianceCheckInterval          = 1000 * time.Millisecond // Interval for checking data protection compliance
    MaxComplianceViolationLimit      = 10                      // Maximum allowed compliance violations before triggering enforcement
)

// DataProtectionComplianceMonitoring automates the monitoring and enforcement of data protection compliance within the blockchain
type DataProtectionComplianceMonitoring struct {
    consensusSystem           *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance            *ledger.Ledger               // Ledger to store compliance-related logs
    stateMutex                *sync.RWMutex                // Mutex for thread-safe access
    complianceViolationCount  int                          // Counter for compliance violations
}

// NewDataProtectionComplianceMonitoring initializes the automation for data protection compliance monitoring
func NewDataProtectionComplianceMonitoring(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataProtectionComplianceMonitoring {
    return &DataProtectionComplianceMonitoring{
        consensusSystem:          consensusSystem,
        ledgerInstance:           ledgerInstance,
        stateMutex:               stateMutex,
        complianceViolationCount: 0,
    }
}

// StartComplianceMonitoringAutomation starts the continuous loop for monitoring data protection compliance
func (automation *DataProtectionComplianceMonitoring) StartComplianceMonitoringAutomation() {
    ticker := time.NewTicker(ComplianceCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceCompliance()
        }
    }()
}

// monitorAndEnforceCompliance checks for data protection compliance violations and triggers enforcement if necessary
func (automation *DataProtectionComplianceMonitoring) monitorAndEnforceCompliance() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Check for compliance violations from the consensus system
    complianceViolations := automation.consensusSystem.CheckComplianceViolations()

    if len(complianceViolations) >= MaxComplianceViolationLimit {
        fmt.Printf("Compliance violations exceed limit (%d). Triggering enforcement.\n", len(complianceViolations))
        automation.triggerComplianceEnforcement(complianceViolations)
    } else {
        fmt.Printf("Compliance violations are within acceptable range (%d).\n", len(complianceViolations))
    }

    automation.complianceViolationCount++
    fmt.Printf("Compliance monitoring cycle #%d executed.\n", automation.complianceViolationCount)

    if automation.complianceViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerComplianceEnforcement triggers enforcement of data protection policies based on violations
func (automation *DataProtectionComplianceMonitoring) triggerComplianceEnforcement(violations []common.ComplianceViolation) {
    for _, violation := range violations {
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting validator for compliance enforcement.")
            continue
        }

        // Encrypt compliance violation data before enforcement
        encryptedViolation := automation.AddEncryptionToViolationData(violation)

        fmt.Printf("Validator %s selected for enforcing compliance.\n", validator.Address)

        // Enforce compliance via consensus using the selected validator
        enforcementSuccess := automation.consensusSystem.EnforceCompliancePolicy(validator, encryptedViolation)
        if enforcementSuccess {
            fmt.Println("Compliance policy successfully enforced.")
        } else {
            fmt.Println("Error enforcing compliance policy.")
        }

        // Log the enforcement action into the ledger
        automation.logComplianceEnforcement(violation)
    }
}

// finalizeEnforcementCycle finalizes the compliance enforcement cycle and logs the result in the ledger
func (automation *DataProtectionComplianceMonitoring) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeComplianceCycle()
    if success {
        fmt.Println("Compliance enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing compliance enforcement cycle.")
    }
}

// logComplianceEnforcement logs every compliance enforcement action into the ledger
func (automation *DataProtectionComplianceMonitoring) logComplianceEnforcement(violation common.ComplianceViolation) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compliance-enforcement-%s", violation.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Compliance Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compliance enforcement action for ViolationID %s.\n", violation.ID)
}

// logEnforcementCycleFinalization logs the finalization of a compliance enforcement cycle into the ledger
func (automation *DataProtectionComplianceMonitoring) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compliance-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compliance Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with compliance enforcement cycle finalization.")
}

// AddEncryptionToViolationData encrypts compliance violation data before enforcement
func (automation *DataProtectionComplianceMonitoring) AddEncryptionToViolationData(violation common.ComplianceViolation) common.ComplianceViolation {
    encryptedData, err := encryption.EncryptData(violation.Data)
    if err != nil {
        fmt.Println("Error encrypting compliance violation data:", err)
        return violation
    }
    violation.Data = encryptedData
    fmt.Println("Compliance violation data successfully encrypted.")
    return violation
}

// ensureComplianceIntegrity checks the integrity of data protection compliance and triggers enforcement if necessary
func (automation *DataProtectionComplianceMonitoring) ensureComplianceIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateComplianceIntegrity()
    if !integrityValid {
        fmt.Println("Compliance integrity breach detected. Triggering enforcement.")
        automation.triggerComplianceEnforcement(automation.consensusSystem.CheckComplianceViolations())
    } else {
        fmt.Println("Compliance integrity is valid.")
    }
}
