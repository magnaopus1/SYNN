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
    RenewableEnergyCheckInterval    = 2000 * time.Millisecond // Interval for checking renewable energy compliance
    MaxAllowedNonRenewableViolations = 5                      // Maximum allowed non-compliance before enforcement
    SubBlocksPerBlock               = 1000                    // Number of sub-blocks in a block
)

// RenewableEnergyEnforcementAutomation automates the enforcement of renewable energy usage across all nodes and validators
type RenewableEnergyEnforcementAutomation struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger to store renewable energy enforcement logs
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    nonRenewableViolationCount int                       // Counter for non-compliance violations
}

// NewRenewableEnergyEnforcementAutomation initializes the automation for enforcing renewable energy standards
func NewRenewableEnergyEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RenewableEnergyEnforcementAutomation {
    return &RenewableEnergyEnforcementAutomation{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        nonRenewableViolationCount: 0,
    }
}

// StartRenewableEnergyEnforcement starts the continuous loop for monitoring renewable energy compliance across nodes and validators
func (automation *RenewableEnergyEnforcementAutomation) StartRenewableEnergyEnforcement() {
    ticker := time.NewTicker(RenewableEnergyCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceRenewableEnergy()
        }
    }()
}

// monitorAndEnforceRenewableEnergy checks renewable energy compliance across nodes and validators and triggers enforcement actions if necessary
func (automation *RenewableEnergyEnforcementAutomation) monitorAndEnforceRenewableEnergy() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch renewable energy compliance data for all nodes and validators from Synnergy Consensus
    complianceData := automation.consensusSystem.GetRenewableEnergyComplianceData()

    for _, node := range complianceData {
        if !node.IsRenewableEnergyCompliant {
            fmt.Printf("Node %s is non-compliant with renewable energy standards. Triggering enforcement action.\n", node.Address)
            automation.triggerRenewableEnergyEnforcement(node)
        } else {
            fmt.Printf("Node %s is compliant with renewable energy standards.\n", node.Address)
        }
    }

    automation.nonRenewableViolationCount++
    fmt.Printf("Renewable energy enforcement cycle #%d executed.\n", automation.nonRenewableViolationCount)

    if automation.nonRenewableViolationCount%SubBlocksPerBlock == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerRenewableEnergyEnforcement triggers enforcement actions for nodes or validators that are non-compliant with renewable energy usage
func (automation *RenewableEnergyEnforcementAutomation) triggerRenewableEnergyEnforcement(node common.NodeComplianceData) {
    // Encrypt the node's renewable energy compliance data before triggering enforcement actions
    encryptedData := automation.AddEncryptionToComplianceData(node)

    // Trigger enforcement using Synnergy Consensus to bring the node into renewable energy compliance
    enforcementSuccess := automation.consensusSystem.EnforceRenewableEnergyCompliance(encryptedData)
    if enforcementSuccess {
        fmt.Printf("Renewable energy compliance enforced for node %s.\n", node.Address)
        automation.logRenewableEnergyEnforcement(node)
    } else {
        fmt.Printf("Error enforcing renewable energy compliance for node %s.\n", node.Address)
    }
}

// finalizeEnforcementCycle finalizes the renewable energy enforcement cycle and logs the result in the ledger
func (automation *RenewableEnergyEnforcementAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeRenewableEnergyEnforcementCycle()
    if success {
        fmt.Println("Renewable energy enforcement cycle finalized successfully.")
        automation.logEnforcementCycleFinalization()
    } else {
        fmt.Println("Error finalizing renewable energy enforcement cycle.")
    }
}

// logRenewableEnergyEnforcement logs each renewable energy enforcement action into the ledger for traceability
func (automation *RenewableEnergyEnforcementAutomation) logRenewableEnergyEnforcement(node common.NodeComplianceData) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("renewable-energy-enforcement-%s", node.Address),
        Timestamp: time.Now().Unix(),
        Type:      "Renewable Energy Enforcement",
        Status:    "Enforced",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with renewable energy enforcement action for NodeID %s.\n", node.Address)
}

// logEnforcementCycleFinalization logs the finalization of a renewable energy enforcement cycle into the ledger
func (automation *RenewableEnergyEnforcementAutomation) logEnforcementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("renewable-energy-enforcement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Renewable Energy Enforcement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with renewable energy enforcement cycle finalization.")
}

// AddEncryptionToComplianceData encrypts the node's compliance data before triggering enforcement actions
func (automation *RenewableEnergyEnforcementAutomation) AddEncryptionToComplianceData(node common.NodeComplianceData) common.NodeComplianceData {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node compliance data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node compliance data successfully encrypted.")
    return node
}

// ensureComplianceIntegrity checks the integrity of compliance data and ensures enforcement actions are triggered if necessary
func (automation *RenewableEnergyEnforcementAutomation) ensureComplianceIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateComplianceDataIntegrity()
    if !integrityValid {
        fmt.Println("Compliance data integrity breach detected. Triggering renewable energy enforcement.")
        automation.monitorAndEnforceRenewableEnergy()
    } else {
        fmt.Println("Compliance data integrity is valid.")
    }
}
