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
    QuadraticVotingCheckInterval  = 2000 * time.Millisecond // Interval for checking quadratic voting proposals
    SubBlocksPerBlock             = 1000                    // Number of sub-blocks in a block
)

// QuadraticVotingAutomation automates the quadratic voting process for governance proposals
type QuadraticVotingAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store quadratic voting executions
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    proposalExecutionCount int                          // Counter for proposals executed
}

// NewQuadraticVotingAutomation initializes the automation for handling governance proposals using quadratic voting
func NewQuadraticVotingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *QuadraticVotingAutomation {
    return &QuadraticVotingAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        proposalExecutionCount: 0,
    }
}

// StartQuadraticVotingExecution starts the continuous loop for monitoring and executing quadratic voting proposals
func (automation *QuadraticVotingAutomation) StartQuadraticVotingExecution() {
    ticker := time.NewTicker(QuadraticVotingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndExecuteQuadraticVotingProposals()
        }
    }()
}

// monitorAndExecuteQuadraticVotingProposals checks for active proposals using quadratic voting and triggers execution if voting is complete
func (automation *QuadraticVotingAutomation) monitorAndExecuteQuadraticVotingProposals() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch active quadratic voting proposals from the Synnergy Consensus
    activeProposals := automation.consensusSystem.GetActiveQuadraticVotingProposals()

    for _, proposal := range activeProposals {
        fmt.Printf("Processing quadratic voting proposal %s.\n", proposal.ID)
        automation.executeProposalIfComplete(proposal)
    }

    automation.proposalExecutionCount++
    fmt.Printf("Quadratic voting execution cycle #%d executed.\n", automation.proposalExecutionCount)

    if automation.proposalExecutionCount%SubBlocksPerBlock == 0 {
        automation.finalizeExecutionCycle()
    }
}

// executeProposalIfComplete checks if the voting process is complete and executes the proposal based on quadratic voting results
func (automation *QuadraticVotingAutomation) executeProposalIfComplete(proposal common.GovernanceProposal) {
    if proposal.VotingComplete {
        // Encrypt the proposal data before execution
        encryptedProposal := automation.AddEncryptionToProposalData(proposal)

        // Execute the proposal using quadratic voting results
        executionSuccess := automation.consensusSystem.ExecuteQuadraticVotingProposal(encryptedProposal)

        if executionSuccess {
            fmt.Printf("Proposal %s executed successfully based on quadratic voting.\n", proposal.ID)
            automation.logProposalExecution(proposal)
        } else {
            fmt.Printf("Error executing proposal %s based on quadratic voting.\n", proposal.ID)
        }
    } else {
        fmt.Printf("Proposal %s voting is not complete.\n", proposal.ID)
    }
}

// finalizeExecutionCycle finalizes the quadratic voting execution cycle and logs the result in the ledger
func (automation *QuadraticVotingAutomation) finalizeExecutionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeQuadraticVotingCycle()
    if success {
        fmt.Println("Quadratic voting execution cycle finalized successfully.")
        automation.logExecutionCycleFinalization()
    } else {
        fmt.Println("Error finalizing quadratic voting execution cycle.")
    }
}

// logProposalExecution logs each executed quadratic voting proposal into the ledger for traceability
func (automation *QuadraticVotingAutomation) logProposalExecution(proposal common.GovernanceProposal) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("quadratic-vote-execution-%s", proposal.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Quadratic Voting Proposal Execution",
        Status:    "Executed",
        Details:   fmt.Sprintf("Proposal %s executed based on quadratic voting results.", proposal.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with execution of proposal %s.\n", proposal.ID)
}

// logExecutionCycleFinalization logs the finalization of a quadratic voting execution cycle into the ledger
func (automation *QuadraticVotingAutomation) logExecutionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("quadratic-vote-execution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Quadratic Vote Execution Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with quadratic vote execution cycle finalization.")
}

// AddEncryptionToProposalData encrypts the proposal data before execution
func (automation *QuadraticVotingAutomation) AddEncryptionToProposalData(proposal common.GovernanceProposal) common.GovernanceProposal {
    encryptedData, err := encryption.EncryptData(proposal)
    if err != nil {
        fmt.Println("Error encrypting proposal data:", err)
        return proposal
    }
    proposal.EncryptedData = encryptedData
    fmt.Println("Proposal data successfully encrypted.")
    return proposal
}

// ensureProposalDataIntegrity checks the integrity of governance proposal data and triggers re-execution if necessary
func (automation *QuadraticVotingAutomation) ensureProposalDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateProposalDataIntegrity()
    if !integrityValid {
        fmt.Println("Proposal data integrity breach detected. Re-triggering quadratic voting proposal execution.")
        automation.monitorAndExecuteQuadraticVotingProposals()
    } else {
        fmt.Println("Proposal data integrity is valid.")
    }
}
