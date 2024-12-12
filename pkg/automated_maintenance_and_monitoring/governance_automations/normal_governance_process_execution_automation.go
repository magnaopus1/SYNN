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
    GovernanceCheckInterval  = 3000 * time.Millisecond // Interval for checking and processing governance proposals
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks in a block
)

// NormalGovernanceProcessExecutionAutomation automates the standard governance process execution for governance proposals
type NormalGovernanceProcessExecutionAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store governance proposal executions
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    governanceExecutionCount int                        // Counter for governance proposal executions
}

// NewNormalGovernanceProcessExecutionAutomation initializes the automation for standard governance proposal execution
func NewNormalGovernanceProcessExecutionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NormalGovernanceProcessExecutionAutomation {
    return &NormalGovernanceProcessExecutionAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        governanceExecutionCount: 0,
    }
}

// StartGovernanceExecution starts the continuous loop for monitoring and executing governance proposals
func (automation *NormalGovernanceProcessExecutionAutomation) StartGovernanceExecution() {
    ticker := time.NewTicker(GovernanceCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndExecuteGovernanceProposals()
        }
    }()
}

// monitorAndExecuteGovernanceProposals checks for standard governance proposals and triggers execution if voting is complete
func (automation *NormalGovernanceProcessExecutionAutomation) monitorAndExecuteGovernanceProposals() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch active governance proposals from the Synnergy Consensus
    activeProposals := automation.consensusSystem.GetActiveStandardGovernanceProposals()

    for _, proposal := range activeProposals {
        fmt.Printf("Processing governance proposal %s.\n", proposal.ID)
        automation.executeProposalIfComplete(proposal)
    }

    automation.governanceExecutionCount++
    fmt.Printf("Governance proposal execution cycle #%d executed.\n", automation.governanceExecutionCount)

    if automation.governanceExecutionCount%SubBlocksPerBlock == 0 {
        automation.finalizeExecutionCycle()
    }
}

// executeProposalIfComplete checks if voting is complete and executes the standard governance proposal
func (automation *NormalGovernanceProcessExecutionAutomation) executeProposalIfComplete(proposal common.GovernanceProposal) {
    if proposal.VotingComplete {
        // Encrypt the proposal data before execution
        encryptedProposal := automation.AddEncryptionToProposalData(proposal)

        // Execute the proposal based on the final vote
        executionSuccess := automation.consensusSystem.ExecuteStandardGovernanceProposal(encryptedProposal)

        if executionSuccess {
            fmt.Printf("Proposal %s executed successfully based on standard governance.\n", proposal.ID)
            automation.logProposalExecution(proposal)
        } else {
            fmt.Printf("Error executing proposal %s based on standard governance.\n", proposal.ID)
        }
    } else {
        fmt.Printf("Proposal %s voting is not complete.\n", proposal.ID)
    }
}

// finalizeExecutionCycle finalizes the standard governance execution cycle and logs the result in the ledger
func (automation *NormalGovernanceProcessExecutionAutomation) finalizeExecutionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeGovernanceExecutionCycle()
    if success {
        fmt.Println("Standard governance execution cycle finalized successfully.")
        automation.logExecutionCycleFinalization()
    } else {
        fmt.Println("Error finalizing standard governance execution cycle.")
    }
}

// logProposalExecution logs each executed proposal into the ledger for traceability
func (automation *NormalGovernanceProcessExecutionAutomation) logProposalExecution(proposal common.GovernanceProposal) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("proposal-execution-%s", proposal.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Proposal Execution",
        Status:    "Executed",
        Details:   fmt.Sprintf("Proposal %s executed based on standard governance results.", proposal.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with execution of proposal %s.\n", proposal.ID)
}

// logExecutionCycleFinalization logs the finalization of a standard governance execution cycle into the ledger
func (automation *NormalGovernanceProcessExecutionAutomation) logExecutionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("governance-execution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Governance Execution Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with standard governance execution cycle finalization.")
}

// AddEncryptionToProposalData encrypts the proposal data before execution
func (automation *NormalGovernanceProcessExecutionAutomation) AddEncryptionToProposalData(proposal common.GovernanceProposal) common.GovernanceProposal {
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
func (automation *NormalGovernanceProcessExecutionAutomation) ensureProposalDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateProposalDataIntegrity()
    if !integrityValid {
        fmt.Println("Proposal data integrity breach detected. Re-triggering governance proposal execution.")
        automation.monitorAndExecuteGovernanceProposals()
    } else {
        fmt.Println("Proposal data integrity is valid.")
    }
}
