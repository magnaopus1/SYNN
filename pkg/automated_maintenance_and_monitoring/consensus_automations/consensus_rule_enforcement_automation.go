package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
)

const (
    RuleCheckInterval  = 15 * time.Second // Interval for checking rule enforcement
    RuleEnforcementKey = "rule_enforcement_key" // Encryption key for rule enforcement logs
)

// ConsensusRuleEnforcementAutomation automates rule enforcement across Synnergy Consensus
type ConsensusRuleEnforcementAutomation struct {
    consensusSystem *consensus.SynnergyConsensus // SynnergyConsensus struct
    ledgerInstance  *ledger.Ledger               // Ledger for storing rule enforcement data
    stateMutex      *sync.RWMutex                // Mutex for thread-safe access
}

// NewConsensusRuleEnforcementAutomation initializes rule enforcement automation
func NewConsensusRuleEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusRuleEnforcementAutomation {
    return &ConsensusRuleEnforcementAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
    }
}

// StartRuleEnforcementMonitoring initiates continuous monitoring for rule enforcement in Synnergy Consensus
func (automation *ConsensusRuleEnforcementAutomation) StartRuleEnforcementMonitoring() {
    ticker := time.NewTicker(RuleCheckInterval)

    go func() {
        for range ticker.C {
            automation.enforcePoHCompliance()
            automation.enforcePoSCompliance()
            automation.enforcePoWCompliance()
            automation.logRuleEnforcementMetrics()
        }
    }()
}

// enforcePoHCompliance ensures that PoH adheres to the rules of Synnergy Consensus
func (automation *ConsensusRuleEnforcementAutomation) enforcePoHCompliance() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    // Check recent PoH proofs and validate them
    proofs := automation.consensusSystem.PoH.GenerateMultipleProofs(10)
    for _, proof := range proofs {
        if !automation.consensusSystem.PoH.ValidatePoHProof(proof, "enforce") {
            fmt.Printf("PoH Proof %s flagged as non-compliant.\n", proof.Hash)
            automation.flagNonCompliantPoH(proof.Hash)
        }
    }
}

// flagNonCompliantPoH flags non-compliant PoH proofs and halts further processing
func (automation *ConsensusRuleEnforcementAutomation) flagNonCompliantPoH(proofHash string) {
    fmt.Printf("PoH Proof %s flagged as non-compliant. Halting further processing.\n", proofHash)
    automation.haltTransactionProcessing("PoH", proofHash)
}

// enforcePoSCompliance ensures that PoS validator selection follows Synnergy Consensus rules
func (automation *ConsensusRuleEnforcementAutomation) enforcePoSCompliance() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    activeValidators := automation.consensusSystem.PoS.State.Validators

    // Check for compliance in validator selection or any suspicious activity
    for _, validator := range activeValidators {
        if validator.Stake < automation.consensusSystem.PoS.State.TotalStake*0.01 { // Example: too low a stake
            fmt.Printf("PoS Validator %s flagged as non-compliant.\n", validator.Address)
            automation.flagNonCompliantValidator(validator.Address)
        }
    }
}

// flagNonCompliantValidator flags non-compliant PoS validators and halts their activities
func (automation *ConsensusRuleEnforcementAutomation) flagNonCompliantValidator(validatorAddress string) {
    fmt.Printf("Validator %s flagged as non-compliant. Halting further validation.\n", validatorAddress)
    automation.haltTransactionProcessing("PoS", validatorAddress)
}

// enforcePoWCompliance ensures that PoW block validation adheres to consensus rules
func (automation *ConsensusRuleEnforcementAutomation) enforcePoWCompliance() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    lastBlock := automation.consensusSystem.LedgerInstance.GetLastBlock()

    // Validate the last mined block for compliance
    if !automation.consensusSystem.PoW.ValidateBlock(&lastBlock) {
        fmt.Printf("PoW Block %s flagged as non-compliant.\n", lastBlock.Hash)
        automation.flagNonCompliantBlock(lastBlock.Hash)
    }
}

// flagNonCompliantBlock flags non-compliant PoW blocks and halts further block finalization
func (automation *ConsensusRuleEnforcementAutomation) flagNonCompliantBlock(blockHash string) {
    fmt.Printf("Block %s flagged as non-compliant. Halting further processing.\n", blockHash)
    automation.haltTransactionProcessing("PoW", blockHash)
}

// haltTransactionProcessing halts any ongoing transaction, sub-block, or block processing
func (automation *ConsensusRuleEnforcementAutomation) haltTransactionProcessing(consensusStage string, identifier string) {
    fmt.Printf("Halting transaction processing in %s stage for identifier %s.\n", consensusStage, identifier)
    
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    switch consensusStage {
    case "PoH":
        automation.haltPoHProcessing(identifier)
    case "PoS":
        automation.haltPoSProcessing(identifier)
    case "PoW":
        automation.haltPoWProcessing(identifier)
    default:
        fmt.Printf("Unknown consensus stage: %s\n", consensusStage)
    }
}

// haltPoHProcessing halts PoH-related processing for the given proof identifier
func (automation *ConsensusRuleEnforcementAutomation) haltPoHProcessing(proofHash string) {
    fmt.Printf("Halting PoH processing for proof: %s\n", proofHash)
    proof := automation.consensusSystem.PoH.GetProofByHash(proofHash)
    if proof == nil {
        fmt.Printf("PoH proof %s not found.\n", proofHash)
        return
    }

    // Mark the proof as invalid and halt timestamp generation
    automation.consensusSystem.PoH.MarkProofAsInvalid(proofHash)
    automation.consensusSystem.PoH.HaltTimestampGeneration(proofHash)
}

// haltPoSProcessing halts PoS-related processing for the given validator identifier
func (automation *ConsensusRuleEnforcementAutomation) haltPoSProcessing(validatorAddress string) {
    fmt.Printf("Halting PoS processing for validator: %s\n", validatorAddress)
    validator := automation.consensusSystem.PoS.GetValidatorByAddress(validatorAddress)
    if validator == nil {
        fmt.Printf("Validator %s not found.\n", validatorAddress)
        return
    }

    // Disqualify the validator and freeze their stake
    automation.consensusSystem.PoS.DisqualifyValidator(validatorAddress)
    automation.consensusSystem.PoS.FreezeValidatorStake(validatorAddress)
}

// haltPoWProcessing halts PoW-related processing for the given block identifier
func (automation *ConsensusRuleEnforcementAutomation) haltPoWProcessing(blockHash string) {
    fmt.Printf("Halting PoW processing for block: %s\n", blockHash)
    block := automation.consensusSystem.PoW.GetBlockByHash(blockHash)
    if block == nil {
        fmt.Printf("PoW block %s not found.\n", blockHash)
        return
    }

    // Mark the block as invalid and stop mining operations
    automation.consensusSystem.PoW.MarkBlockAsInvalid(blockHash)
    automation.consensusSystem.PoW.HaltBlockMining(blockHash)
}

// logRuleEnforcementMetrics logs rule enforcement data and encrypts it for security
func (automation *ConsensusRuleEnforcementAutomation) logRuleEnforcementMetrics() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rule-enforcement-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Rule Enforcement",
        Status:    "Logged",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(RuleEnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting rule enforcement log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Rule enforcement metrics logged and stored in the ledger.")
}

