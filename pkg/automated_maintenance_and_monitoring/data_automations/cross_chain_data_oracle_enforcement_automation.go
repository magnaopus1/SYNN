package data_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
)

const (
    OracleEnforcementInterval = 500 * time.Millisecond // Run enforcement checks every 0.5 seconds
    DataConsistencyThreshold  = 95                    // Required consistency threshold for cross-chain data
)

// CrossChainDataOracleEnforcementAutomation automates the enforcement of cross-chain data consistency using Synnergy Consensus and integrates with the ledger
type CrossChainDataOracleEnforcementAutomation struct {
    consensusSystem   *common.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store consensus and oracle enforcement data
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    enforcementCount  int                          // Counter for enforcement actions
}

// NewCrossChainDataOracleEnforcementAutomation initializes the automation for cross-chain data enforcement
func NewCrossChainDataOracleEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *CrossChainDataOracleEnforcementAutomation {
    return &CrossChainDataOracleEnforcementAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        enforcementCount: 0,
    }
}

// StartEnforcementAutomation starts the continuous automation loop for cross-chain data enforcement
func (automation *CrossChainDataOracleEnforcementAutomation) StartEnforcementAutomation() {
    ticker := time.NewTicker(OracleEnforcementInterval)

    go func() {
        for range ticker.C {
            automation.enforceDataConsistency()
        }
    }()
}

// enforceDataConsistency checks the consistency of cross-chain data and triggers necessary enforcement actions
func (automation *CrossChainDataOracleEnforcementAutomation) enforceDataConsistency() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    consistency := automation.consensusSystem.CrossChainDataConsistency()

    if consistency < DataConsistencyThreshold {
        fmt.Printf("Data consistency below threshold (%d%%). Triggering enforcement.\n", consistency)
        automation.triggerEnforcementActions()
    } else {
        fmt.Printf("Cross-chain data consistency is within acceptable range (%d%%).\n", consistency)
    }

    automation.enforcementCount++
    fmt.Printf("Enforcement action #%d executed.\n", automation.enforcementCount)

    if automation.enforcementCount%100 == 0 {
        automation.finalizeEnforcementCycle()
    }
}

// triggerEnforcementActions triggers actions to enforce data consistency across chains
func (automation *CrossChainDataOracleEnforcementAutomation) triggerEnforcementActions() {
    validator := automation.consensusSystem.PoS.SelectValidator()
    if validator == nil {
        fmt.Println("Error selecting validator for enforcement.")
        return
    }

    fmt.Printf("Validator %s selected for enforcement.\n", validator.Address)

    enforcementSuccess := automation.consensusSystem.EnforceCrossChainConsistency(validator)
    if enforcementSuccess {
        fmt.Println("Cross-chain data consistency successfully enforced.")
    } else {
        fmt.Println("Error enforcing cross-chain data consistency.")
    }

    automation.logEnforcementAction()
}

// finalizeEnforcementCycle finalizes the enforcement cycle and integrates with the ledger for accountability
func (automation *CrossChainDataOracleEnforcementAutomation) finalizeEnforcementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeEnforcementCycle()
    if success {
        fmt.Println("Enforcement cycle finalized successfully.")
        automation.logFinalization()
    } else {
        fmt.Println("Error finalizing enforcement cycle.")
    }
}

// logEnforcementAction logs every enforcement action into the ledger for audit purposes
func (automation *CrossChainDataOracleEnforcementAutomation) logEnforcementAction() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("enforcement-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Cross-Chain Data Enforcement",
        Status:    "Executed",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with enforcement action.\n")
}

// logFinalization logs the finalization of an enforcement cycle into the ledger
func (automation *CrossChainDataOracleEnforcementAutomation) logFinalization() {
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

// AddEncryption adds encryption to sensitive data within the enforcement process
func (automation *CrossChainDataOracleEnforcementAutomation) AddEncryption(data []byte) []byte {
    encryptedData, err := encryption.EncryptData(data)
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return nil
    }
    fmt.Println("Data successfully encrypted.")
    return encryptedData
}

// ensureCrossChainIntegrity triggers only when integrity breaches are detected
func (automation *CrossChainDataOracleEnforcementAutomation) ensureCrossChainIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCrossChainIntegrity()
    if !integrityValid {
        fmt.Println("Cross-chain integrity breach detected. Triggering enforcement.")
        automation.triggerEnforcementActions()
    }
}
