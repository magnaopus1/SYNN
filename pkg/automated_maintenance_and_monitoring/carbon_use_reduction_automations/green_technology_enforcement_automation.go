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
    GreenTechnologyCheckInterval   = 1800 * time.Millisecond // Interval for checking green technology compliance
    MaxAllowedNonGreenViolations   = 5                       // Maximum number of non-compliance issues before triggering enforcement
    SubBlocksPerBlock              = 1000                    // Number of sub-blocks in a block
)

// GreenTechnologyEnforcementAutomation automates the enforcement of green technology compliance across nodes and validators
type GreenTechnologyEnforcementAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store compliance-related logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    nonGreenViolationCount int                          // Counter for non-compliance violations
}

// NewGreenTechnologyEnforcementAutomation initializes the automation for enforcing green technology standards
func NewGreenTechnologyEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GreenTechnologyEnforcementAutomation {
    return &GreenTechnologyEnforcementAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        nonGreenViolationCount: 0,
    }
}

// StartGreenTechnologyEnforcement starts the continuous loop for monitoring green technology compliance across nodes and validators
func (automation *GreenTechnologyEnforcementAutomation) StartGreenTechnologyEnforcement() {
    ticker := time.NewTicker(GreenTechnologyCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceGreenTechnology()
        }
    }()
}

// monitorAndEnforceGreenTechnology checks nodes and validators for green technology compliance and triggers enforcement actions if required
func (automation *GreenTechnologyEnforcementAutomation) monitorAndEnforceGreenTechnology() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch compliance data for all nodes and validators from the Synnergy Consensus
    complianceData := automation.consensusSystem.GetGreenTechnologyComplianceData()

    for _, node := range complianceData {
        if !node.IsGreenTechCompliant {
            fmt.Printf("Node %s is non-compliant with green technology standards. Triggering enforcement action.\n", node.Address)
            automation.triggerGreenTechnologyEnforcement(node)
        } else {
            fmt.Printf("Node %s is compliant with green technology standards.\n", node.Address)
        }
    }

    automation.nonGreenViolationCount++
    fmt.Printf("Green technology enforcement cycle #%d executed.\n", automation.nonGreenViolationCount)

    if automation.nonGreenViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerGreenTechnologyEnforcement triggers enforcement actions for nodes or validators that are non-compliant with green technology standards
func (automation *GreenTechnologyEnforcementAutomation) triggerGreenTechnologyEnforcement(node common.NodeComplianceData) {
    // Encrypt the node's compliance data before triggering enforcement actions
    encryptedData := automation.AddEncryptionToComplianceData(node)

    // Trigger enforcement using the Synnergy Consensus to bring the node into green technology compliance
    enforcementSuccess := automation.consensusSystem.EnforceGreenTechnologyCompliance(encryptedData)
    if enforcementSuccess {
        fmt.Printf("Green technology compliance enforced for node %s.\n", node.Address)
        automation.logGreenTechnologyEnforcement(node)
    } else {
        fmt.Printf("Error enforcing green technology compliance for node %s.\n", node.Address)
    }
}

// finalizeEnforcementCycle finalizes the green technology enforcement cycle and logs the result in the ledger
func (automation *GreenTechnologyEnforcementAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeGreenTechnologyEnforcementCycle()
    if success {
        fmt.Println("Green technology enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing green technology enforcement cycle.")
    }
}

// logGreenTechnologyEnforcement logs each green technology enforcement action into the ledger for traceability
func (automation *GreenTechnologyEnforcementAutomation) logGreenTechnologyEnforcement(node common.NodeComplianceData) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("green-tech-enforcement-%s", node.Address),
        Timestamp: time.Now().Unix(),
        Type:      "Green Technology Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with green technology enforcement action for NodeID %s.\n", node.Address)
}

// logEnforcementCycleFinalization logs the finalization of a green technology enforcement cycle into the ledger
func (automation *GreenTechnologyEnforcementAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("green-tech-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Green Technology Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with green technology enforcement cycle finalization.")
}

// AddEncryptionToComplianceData encrypts the node's compliance data before triggering enforcement actions
func (automation *GreenTechnologyEnforcementAutomation) AddEncryptionToComplianceData(node common.NodeComplianceData) common.NodeComplianceData {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node compliance data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node compliance data successfully encrypted.")
    return node
}

// ensureComplianceIntegrity checks the integrity of compliance data and ensures proper actions are triggered
func (automation *GreenTechnologyEnforcementAutomation) ensureComplianceIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateComplianceDataIntegrity()
    if !integrityValid {
        fmt.Println("Compliance data integrity breach detected. Triggering green technology enforcement.")
        automation.monitorAndEnforceGreenTechnology()
    } else {
        fmt.Println("Compliance data integrity is valid.")
    }
}
