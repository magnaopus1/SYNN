package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
)

const (
    SubBlockValidationInterval = 600 * time.Millisecond // Validate sub-block every 0.6 seconds
    SubBlocksPerBlock          = 1000                   // Number of sub-blocks per block
)

// ConsensusMechanismExecutionAutomation automates the execution of Synnergy Consensus using PoH, PoS, and PoW
type ConsensusMechanismExecutionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store consensus-related data
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    subBlockCount     int                          // Counter for sub-block validation
}

// NewConsensusMechanismExecutionAutomation initializes the automation for Synnergy Consensus execution
func NewConsensusMechanismExecutionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusMechanismExecutionAutomation {
    return &ConsensusMechanismExecutionAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        subBlockCount:   0,
    }
}

// StartConsensusExecution starts the execution of sub-block validation every 0.6 seconds and block finalization after 1000 sub-blocks
func (automation *ConsensusMechanismExecutionAutomation) StartConsensusExecution() {
    ticker := time.NewTicker(SubBlockValidationInterval)

    go func() {
        for range ticker.C {
            automation.validateSubBlock()
        }
    }()
}

// validateSubBlock alternates between PoH and PoS for sub-block validation every 0.6 seconds
func (automation *ConsensusMechanismExecutionAutomation) validateSubBlock() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if automation.subBlockCount%2 == 0 {
        // Step 1: Use PoH to validate the current sub-block
        success := automation.consensusSystem.PoH.GeneratePoHProof()
        if success {
            fmt.Println("PoH proof generated and sub-block validated by PoH.")
        } else {
            fmt.Println("Error generating PoH proof.")
            return
        }
    } else {
        // Step 2: Use PoS to validate the next sub-block
        validator := automation.consensusSystem.PoS.SelectValidator()
        if validator == nil {
            fmt.Println("Error selecting PoS validator.")
            return
        }
        fmt.Printf("PoS validator %s selected.\n", validator.Address)

        // Validate sub-block using PoS
        success := automation.consensusSystem.PoS.ValidateSubBlock(validator)
        if success {
            fmt.Println("Sub-block validated by PoS.")
        } else {
            fmt.Println("Error validating sub-block via PoS.")
            return
        }
    }

    // Increment the sub-block counter
    automation.subBlockCount++
    fmt.Printf("Sub-block %d validated.\n", automation.subBlockCount)

    // Finalize the block after 1000 sub-blocks
    if automation.subBlockCount >= SubBlocksPerBlock {
        automation.finalizeBlock()
        automation.subBlockCount = 0
    }
}

// finalizeBlock finalizes 1000 sub-blocks into a full block using PoW
func (automation *ConsensusMechanismExecutionAutomation) finalizeBlock() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.PoW.MineBlock()
    if success {
        fmt.Println("Block successfully finalized with PoW.")
        automation.logBlockFinalization()
    } else {
        fmt.Println("Error finalizing block with PoW.")
    }
}

// logBlockFinalization logs the block finalization event into the ledger
func (automation *ConsensusMechanismExecutionAutomation) logBlockFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("block-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Block Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with block finalization event.\n")
}

// ensureConsensusConsistency checks and validates the consensus chain to ensure it's operating correctly
func (automation *ConsensusMechanismExecutionAutomation) ensureConsensusConsistency() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    valid := automation.consensusSystem.ValidateConsensusChain()
    if !valid {
        fmt.Println("Consensus validation failed.")
        return
    }

    fmt.Println("Consensus chain is valid and consistent.")
}
