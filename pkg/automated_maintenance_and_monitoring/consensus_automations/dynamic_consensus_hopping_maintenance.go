package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    ConsensusHoppingInterval = 5 * time.Second // Interval for dynamic consensus hopping checks
    OverrideKey              = "dynamic_hopping_override_key" // Encryption key for overriding consensus loop
)

// DynamicConsensusHoppingAutomation automates dynamic consensus hopping in Synnergy Consensus
type DynamicConsensusHoppingAutomation struct {
    ledgerInstance   *ledger.Ledger                   // Blockchain ledger for tracking consensus hopping actions
    consensusEngine  *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    stateMutex       *sync.RWMutex                    // Mutex for thread-safe ledger access
    dynamicHoppingOn bool                             // Flag to track if dynamic hopping override is active
}

// NewDynamicConsensusHoppingAutomation initializes the automation for dynamic consensus hopping
func NewDynamicConsensusHoppingAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *DynamicConsensusHoppingAutomation {
    return &DynamicConsensusHoppingAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        dynamicHoppingOn: false, // Default to off, i.e., main consensus loop runs by default
    }
}

// StartConsensusHopping initiates the continuous dynamic hopping between PoH and PoS based on system load and validator allocation
func (automation *DynamicConsensusHoppingAutomation) StartConsensusHopping() {
    ticker := time.NewTicker(ConsensusHoppingInterval)

    go func() {
        for range ticker.C {
            if automation.dynamicHoppingOn {
                automation.hopConsensus()
            } else {
                fmt.Println("Dynamic consensus hopping is off, using the main consensus loop.")
            }
        }
    }()
}

// hopConsensus automates the dynamic switch between PoH and PoS based on system load and validator allocation
func (automation *DynamicConsensusHoppingAutomation) hopConsensus() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    load := automation.consensusEngine.MeasureSystemLoad() // Measure system load
    validatorWeighting := automation.consensusEngine.CalculateValidatorWeighting() // Calculate PoH vs PoS validator weighting

    if load > 0.75 && validatorWeighting["PoH"] > validatorWeighting["PoS"] {
        fmt.Println("High load detected, favoring PoH for the next block.")
        automation.generatePoHProof()
    } else {
        fmt.Println("Moderate load detected, allocating more validators to PoS for the next block.")
        automation.validatePoSBlock()
    }

    // Finalize the block when the sub-block limit is reached
    if automation.consensusEngine.SubBlockCount >= synnergy_consensus.SubBlocksPerBlock {
        automation.finalizeBlock()
        automation.consensusEngine.SubBlockCount = 0 // Reset for the next block
    }
}

// generatePoHProof handles PoH proof generation and sub-block validation
func (automation *DynamicConsensusHoppingAutomation) generatePoHProof() {
    success := automation.consensusEngine.PoH.GenerateProof()
    if success {
        fmt.Println("PoH proof generated successfully.")
    } else {
        fmt.Println("Error generating PoH proof.")
    }
}

// validatePoSBlock handles PoS validation for sub-blocks
func (automation *DynamicConsensusHoppingAutomation) validatePoSBlock() {
    validatorSelected := automation.consensusEngine.PoS.SelectValidator()
    if !validatorSelected {
        fmt.Println("Error selecting PoS validator.")
        return
    }

    success := automation.consensusEngine.PoS.ValidateSubBlock()
    if success {
        fmt.Println("Sub-block validated successfully via PoS.")
    } else {
        fmt.Println("Error validating sub-block via PoS.")
    }
}

// finalizeBlock finalizes 1000 sub-blocks into a full block using PoW
func (automation *DynamicConsensusHoppingAutomation) finalizeBlock() {
    success := automation.consensusEngine.PoW.FinalizeBlock()
    if !success {
        fmt.Println("Error finalizing block using PoW.")
        return
    }

    fmt.Println("Block successfully finalized using PoW.")
    automation.logConsensusHopping()
}

// logConsensusHopping logs the consensus hopping event and block finalization into the ledger
func (automation *DynamicConsensusHoppingAutomation) logConsensusHopping() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("consensus-hopping-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Consensus Hopping and Block Finalization",
        Status:    "Finalized",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(ConsensusKey))
    if err != nil {
        fmt.Printf("Error encrypting consensus hopping log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with consensus hopping and block finalization event.\n")
}

// OverrideMainConsensusLoop allows this automation to override the main consensus loop and activate dynamic hopping
func (automation *DynamicConsensusHoppingAutomation) OverrideMainConsensusLoop() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    automation.dynamicHoppingOn = true

    // Log the override action in the ledger
    automation.logOverrideAction("Dynamic consensus hopping override activated.")
}

// StopOverride reverts the consensus back to the main loop without dynamic hopping
func (automation *DynamicConsensusHoppingAutomation) StopOverride() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    automation.dynamicHoppingOn = false

    // Log the stop override action in the ledger
    automation.logOverrideAction("Dynamic consensus hopping override deactivated.")
}

// logOverrideAction logs the override action in the ledger for auditing
func (automation *DynamicConsensusHoppingAutomation) logOverrideAction(action string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("override-action-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Override Action",
        Status:    action,
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(OverrideKey))
    if err != nil {
        fmt.Printf("Error encrypting override log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with override action: %s.\n", action)
}

// ensureHoppingConsistency checks that the consensus hopping and overall chain are consistent
func (automation *DynamicConsensusHoppingAutomation) ensureHoppingConsistency() {
    fmt.Println("Ensuring consistency in dynamic consensus hopping...")

    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Validate the chain to ensure proper transitions and consensus consistency
    err := automation.consensusEngine.ValidateChain()
    if err != nil {
        fmt.Printf("Consensus chain validation failed: %v\n", err)
    } else {
        fmt.Println("Consensus hopping and chain validation are consistent.")
    }
}

// triggerValidatorReallocation dynamically reallocates validators based on system load and consensus hopping needs
func (automation *DynamicConsensusHoppingAutomation) triggerValidatorReallocation() {
    fmt.Println("Reallocating validators based on system load and consensus hopping needs...")

    success := automation.consensusEngine.ReallocateValidators()
    if !success {
        fmt.Println("Error reallocating validators.")
    } else {
        fmt.Println("Validators reallocated successfully for consensus hopping optimization.")
    }
}
