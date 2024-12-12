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
    PrivacyEnforcementInterval       = 700 * time.Millisecond // Interval for checking privacy violations
    MaxPrivacyViolationLimit         = 5                      // Max allowed privacy violations before triggering enforcement
)

// DataPrivacyEnforcementAutomation automates the enforcement of data privacy rules in the blockchain network
type DataPrivacyEnforcementAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store privacy-related data
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    privacyViolationCount  int                          // Counter for privacy violations
}

// NewDataPrivacyEnforcementAutomation initializes the automation for data privacy enforcement
func NewDataPrivacyEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DataPrivacyEnforcementAutomation {
    return &DataPrivacyEnforcementAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        privacyViolationCount:  0,
    }
}

// StartPrivacyEnforcementAutomation starts the continuous loop for monitoring and enforcing data privacy
func (automation *DataPrivacyEnforcementAutomation) StartPrivacyEnforcementAutomation() {
    ticker := time.NewTicker(PrivacyEnforcementInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceDataPrivacy()
        }
    }()
}

// monitorAndEnforceDataPrivacy checks for data privacy violations and triggers enforcement actions if necessary
func (automation *DataPrivacyEnforcementAutomation) monitorAndEnforceDataPrivacy() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Check if there are any privacy violations within the consensus system
    privacyViolations := automation.consensusSystem.CheckPrivacyViolations()

    if len(privacyViolations) >= MaxPrivacyViolationLimit {
        fmt.Printf("Privacy violations exceed limit (%d). Triggering enforcement.\n", len(privacyViolations))
        automation.triggerPrivacyEnforcement(privacyViolations)
    } else {
        fmt.Printf("Privacy violations are within acceptable range (%d).\n", len(privacyViolations))
    }

    automation.privacyViolationCount++
    fmt.Printf("Privacy enforcement cycle #%d executed.\n", automation.privacyViolationCount)

    if automation.privacyViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerPrivacyEnforcement enforces data privacy actions based on detected violations
func (automation *DataPrivacyEnforcementAutomation) triggerPrivacyEnforcement(violations []common.PrivacyViolation) {
    for _, violation := range violations {
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting validator for privacy enforcement.")
            continue
        }

        // Encrypt privacy violation data before enforcement
        encryptedViolation := automation.AddEncryptionToViolationData(violation)

        fmt.Printf("Validator %s selected for enforcing privacy.\n", validator.Address)

        // Enforce privacy rules using the selected validator
        enforcementSuccess := automation.consensusSystem.EnforcePrivacy(validator, encryptedViolation)
        if enforcementSuccess {
            fmt.Println("Privacy successfully enforced for violation.")
        } else {
            fmt.Println("Error enforcing privacy for violation.")
        }

        // Log the enforcement action into the ledger
        automation.logPrivacyEnforcement(violation)
    }
}

// finalizeEnforcementCycle finalizes the privacy enforcement cycle and logs the result into the ledger
func (automation *DataPrivacyEnforcementAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizePrivacyCycle()
    if success {
        fmt.Println("Privacy enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing privacy enforcement cycle.")
    }
}

// logPrivacyEnforcement logs every privacy enforcement action into the ledger for audit purposes
func (automation *DataPrivacyEnforcementAutomation) logPrivacyEnforcement(violation common.PrivacyViolation) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("privacy-enforcement-%s", violation.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Privacy Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with privacy enforcement action for ViolationID %s.\n", violation.ID)
}

// logEnforcementCycleFinalization logs the finalization of a privacy enforcement cycle into the ledger
func (automation *DataPrivacyEnforcementAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("privacy-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Privacy Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with privacy enforcement cycle finalization.\n")
}

// AddEncryptionToViolationData encrypts privacy violation data before enforcement
func (automation *DataPrivacyEnforcementAutomation) AddEncryptionToViolationData(violation common.PrivacyViolation) common.PrivacyViolation {
    encryptedData, err := encryption.EncryptData(violation.Data)
    if err != nil {
        fmt.Println("Error encrypting privacy violation data:", err)
        return violation
    }
    violation.Data = encryptedData
    fmt.Println("Privacy violation data successfully encrypted.")
    return violation
}

// ensurePrivacyIntegrity monitors privacy integrity and triggers enforcement if any breach occurs
func (automation *DataPrivacyEnforcementAutomation) ensurePrivacyIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidatePrivacyIntegrity()
    if !integrityValid {
        fmt.Println("Privacy integrity breach detected. Triggering enforcement.")
        automation.triggerPrivacyEnforcement(automation.consensusSystem.CheckPrivacyViolations())
    } else {
        fmt.Println("Privacy integrity is valid.")
    }
}
