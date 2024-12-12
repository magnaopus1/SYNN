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
    ConsensusSwitchInterval = 3 * time.Second  // Interval for automating the switch between PoH, PoS, and PoW
    ConsensusKey            = "consensus_switch_encryption_key" // Encryption key for consensus switch logs
)

// CrossConsensusSwitchAutomation automates the transition between PoH, PoS, and PoW stages of the Synnergy Consensus
type CrossConsensusSwitchAutomation struct {
    ledgerInstance   *ledger.Ledger                    // Blockchain ledger for tracking consensus switches
    consensusEngine  *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    stateMutex       *sync.RWMutex                     // Mutex for thread-safe ledger access
    subBlockCount    int                               // Counter for sub-blocks
}

// NewCrossConsensusSwitchAutomation initializes the automation for cross-consensus switching
func NewCrossConsensusSwitchAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *CrossConsensusSwitchAutomation {
    return &CrossConsensusSwitchAutomation{
        ledgerInstance:   ledgerInstance,
        consensusEngine:  consensusEngine,
        stateMutex:       stateMutex,
        subBlockCount:    0,
    }
}

// StartConsensusSwitching initiates the continuous transition between PoH, PoS, and PoW at regular intervals
func (automation *CrossConsensusSwitchAutomation) StartConsensusSwitching() {
    ticker := time.NewTicker(ConsensusSwitchInterval)

    go func() {
        for range ticker.C {
            automation.switchConsensusStages()
        }
    }()
}

// switchConsensusStages automates the transition between PoH, PoS, and PoW stages
func (automation *CrossConsensusSwitchAutomation) switchConsensusStages() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Alternate between PoH and PoS for sub-block validation
    if automation.subBlockCount%2 == 0 {
        // Step 1: Use PoH to validate the current sub-block
        success := automation.consensusEngine.PoH.GenerateProof()
        if !success {
            fmt.Println("Error generating PoH proof for sub-block.")
            return
        }
        fmt.Println("PoH proof generated successfully for sub-block.")
    } else {
        // Step 2: Use PoS to validate the next sub-block
        validatorSelected := automation.consensusEngine.PoS.SelectValidator()
        if !validatorSelected {
            fmt.Println("Error selecting PoS validator for sub-block.")
            return
        }
        success := automation.consensusEngine.PoS.ValidateSubBlock()
        if !success {
            fmt.Println("Error validating sub-block via PoS.")
            return
        }
        fmt.Println("Sub-block validated successfully via PoS.")
    }

    // Increment the sub-block counter
    automation.subBlockCount++
    fmt.Printf("Sub-block %d validated.\n", automation.subBlockCount)

    // After 1000 sub-blocks, finalize the block using PoW
    if automation.subBlockCount >= synnergy_consensus.SubBlocksPerBlock {
        automation.finalizeBlock()
        automation.subBlockCount = 0 // Reset for the next block
    }
}

// finalizeBlock finalizes 1000 sub-blocks into a full block using PoW
func (automation *CrossConsensusSwitchAutomation) finalizeBlock() {
    success := automation.consensusEngine.PoW.FinalizeBlock()
    if !success {
        fmt.Println("Error finalizing block with PoW.")
        return
    }
    fmt.Println("Block successfully finalized using PoW.")

    // Log the block finalization in the ledger
    automation.logConsensusSwitch()
}

// logConsensusSwitch logs the consensus switch and block finalization event into the ledger
func (automation *CrossConsensusSwitchAutomation) logConsensusSwitch() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("consensus-switch-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Consensus Switch and Block Finalization",
        Status:    "Finalized",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(ConsensusKey))
    if err != nil {
        fmt.Printf("Error encrypting consensus switch log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with consensus switch and block finalization event.\n")
}

// ensureConsensusConsistency checks that the overall consensus chain is consistent and correct
func (automation *CrossConsensusSwitchAutomation) ensureConsensusConsistency() {
    fmt.Println("Ensuring overall consensus consistency...")

    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Validate the entire chain for consistency
    err := automation.consensusEngine.ValidateChain()
    if err != nil {
        fmt.Printf("Consensus validation failed: %v\n", err)
    } else {
        fmt.Println("Consensus chain is valid and consistent.")
    }
}

// triggerSwitchOptimization adds functionality to trigger consensus stage optimizations based on specific conditions
func (automation *CrossConsensusSwitchAutomation) triggerSwitchOptimization() {
    fmt.Println("Triggering optimization for consensus stage switch...")

    // Reprocess transactions and rebalance stages if necessary
    success := automation.consensusEngine.OptimizeSwitching()
    if !success {
        fmt.Println("Error optimizing consensus switch process.")
    } else {
        fmt.Println("Consensus switch optimized successfully.")
    }
}

