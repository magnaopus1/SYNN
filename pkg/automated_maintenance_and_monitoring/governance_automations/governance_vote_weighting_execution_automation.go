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
    VoteWeightingCheckInterval  = 2000 * time.Millisecond // Interval for checking voting process and execution
    SubBlocksPerBlock           = 1000                    // Number of sub-blocks in a block
)

// GovernanceVoteWeightingExecutionAutomation automates vote weighting and execution when this voting style is used for governance proposals
type GovernanceVoteWeightingExecutionAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger to store weighted vote execution
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    proposalExecutionCount int                        // Counter for proposals executed
}

// NewGovernanceVoteWeightingExecutionAutomation initializes the automation for executing weighted voting proposals
func NewGovernanceVoteWeightingExecutionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GovernanceVoteWeightingExecutionAutomation {
    return &GovernanceVoteWeightingExecutionAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        proposalExecutionCount: 0,
    }
}

// StartVoteWeightingExecution starts the continuous loop for monitoring and executing weighted voting proposals
func (automation *GovernanceVoteWeightingExecutionAutomation) StartVoteWeightingExecution() {
    ticker := time.NewTicker(VoteWeightingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndExecuteWeightedVotes()
        }
    }()
}

// monitorAndExecuteWeightedVotes checks for active proposals using weighted voting and triggers execution based on the vote tally
func (automation *GovernanceVoteWeightingExecutionAutomation) monitorAndExecuteWeightedVotes() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch active proposals using weighted voting from the Synnergy Consensus
    activeProposals := automation.consensusSystem.GetActiveWeightedVotingProposals()

    for _, proposal := range activeProposals {
        fmt.Printf("Processing weighted voting proposal %s.\n", proposal.ID)
        automation.executeProposalIfComplete(proposal)
    }

    automation.proposalExecutionCount++
    fmt.Printf("Vote weighting execution cycle #%d executed.\n", automation.proposalExecutionCount)

    if automation.proposalExecutionCount%SubBlocksPerBlock == 0 {
        automation.finalizeExecutionCycle()
    }
}

// executeProposalIfComplete checks if the voting is complete and executes the proposal based on weighted votes
func (automation *GovernanceVoteWeightingExecutionAutomation) executeProposalIfComplete(proposal common.GovernanceProposal) {
    if proposal.VotingComplete {
        // Encrypt the voting data before execution
        encryptedProposal := automation.AddEncryptionToProposalData(proposal)

        // Calculate the final vote based on weighting and execute the proposal
        executionSuccess := automation.consensusSystem.ExecuteWeightedVoteProposal(encryptedProposal)

        if executionSuccess {
            fmt.Printf("Proposal %s executed successfully based on weighted voting.\n", proposal.ID)
            automation.logProposalExecution(proposal)
        } else {
            fmt.Printf("Error executing proposal %s based on weighted voting.\n", proposal.ID)
        }
    } else {
        fmt.Printf("Proposal %s voting is not complete.\n", proposal.ID)
    }
}

// finalizeExecutionCycle finalizes the weighted vote execution cycle and logs the result in the ledger
func (automation *GovernanceVoteWeightingExecutionAutomation) finalizeExecutionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVoteExecutionCycle()
    if success {
        fmt.Println("Weighted vote execution cycle finalized successfully.")
        automation.logExecutionCycleFinalization()
    } else {
        fmt.Println("Error finalizing weighted vote execution cycle.")
    }
}

// logProposalExecution logs each executed proposal into the ledger for traceability
func (automation *GovernanceVoteWeightingExecutionAutomation) logProposalExecution(proposal common.GovernanceProposal) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("proposal-execution-%s", proposal.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Proposal Execution",
        Status:    "Executed",
        Details:   fmt.Sprintf("Proposal %s executed based on weighted voting results.", proposal.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with execution of proposal %s.\n", proposal.ID)
}

// logExecutionCycleFinalization logs the finalization of a weighted vote execution cycle into the ledger
func (automation *GovernanceVoteWeightingExecutionAutomation) logExecutionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("weighted-vote-execution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Weighted Vote Execution Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with weighted vote execution cycle finalization.")
}

// AddEncryptionToProposalData encrypts the proposal data before execution
func (automation *GovernanceVoteWeightingExecutionAutomation) AddEncryptionToProposalData(proposal common.GovernanceProposal) common.GovernanceProposal {
    encryptedData, err := encryption.EncryptData(proposal)
    if err != nil {
        fmt.Println("Error encrypting proposal data:", err)
        return proposal
    }
    proposal.EncryptedData = encryptedData
    fmt.Println("Proposal data successfully encrypted.")
    return proposal
}

// ensureVoteDataIntegrity checks the integrity of voting data and triggers re-execution if necessary
func (automation *GovernanceVoteWeightingExecutionAutomation) ensureVoteDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVoteDataIntegrity()
    if !integrityValid {
        fmt.Println("Vote data integrity breach detected. Re-triggering weighted vote execution.")
        automation.monitorAndExecuteWeightedVotes()
    } else {
        fmt.Println("Vote data integrity is valid.")
    }
}
