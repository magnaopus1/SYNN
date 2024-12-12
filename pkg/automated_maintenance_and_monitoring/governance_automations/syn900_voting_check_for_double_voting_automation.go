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
    DoubleVotingCheckInterval  = 1500 * time.Millisecond // Interval for checking double voting attempts
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// Syn900VotingCheckAutomation automates the process of checking for double voting attempts using SYN900 identity tokens
type Syn900VotingCheckAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store double voting actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    doubleVoteCheckCount  int                          // Counter for double voting checks
}

// NewSyn900VotingCheckAutomation initializes the automation for checking SYN900 tokens for double voting
func NewSyn900VotingCheckAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *Syn900VotingCheckAutomation {
    return &Syn900VotingCheckAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        doubleVoteCheckCount: 0,
    }
}

// StartDoubleVotingCheck starts the continuous loop for monitoring and preventing double voting using SYN900 identity tokens
func (automation *Syn900VotingCheckAutomation) StartDoubleVotingCheck() {
    ticker := time.NewTicker(DoubleVotingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndCheckDoubleVoting()
        }
    }()
}

// monitorAndCheckDoubleVoting checks for any double voting attempts using the same SYN900 identity token
func (automation *Syn900VotingCheckAutomation) monitorAndCheckDoubleVoting() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch ongoing voting processes from the Synnergy Consensus
    activeVotes := automation.consensusSystem.GetOngoingVotes()

    for _, vote := range activeVotes {
        if automation.checkForDoubleVoting(vote) {
            fmt.Printf("Double voting detected for SYN900 token %s in vote %s. Rejecting the vote.\n", vote.Syn900Token, vote.ID)
            automation.rejectDoubleVote(vote)
        } else {
            fmt.Printf("No double voting detected for SYN900 token %s in vote %s.\n", vote.Syn900Token, vote.ID)
        }
    }

    automation.doubleVoteCheckCount++
    fmt.Printf("Double voting check cycle #%d executed.\n", automation.doubleVoteCheckCount)

    if automation.doubleVoteCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeCheckCycle()
    }
}

// checkForDoubleVoting checks if the SYN900 token has been used multiple times in the same vote or across votes
func (automation *Syn900VotingCheckAutomation) checkForDoubleVoting(vote common.Vote) bool {
    // Check if the SYN900 token is being used in multiple votes simultaneously
    doubleVoteDetected := automation.consensusSystem.HasTokenBeenUsedTwice(vote.Syn900Token, vote.ID)
    return doubleVoteDetected
}

// rejectDoubleVote rejects the vote if double voting is detected
func (automation *Syn900VotingCheckAutomation) rejectDoubleVote(vote common.Vote) {
    // Encrypt the vote data before logging and rejection
    encryptedVote := automation.AddEncryptionToVoteData(vote)

    // Reject the vote and log the action
    rejectionSuccess := automation.consensusSystem.RejectVote(encryptedVote)

    if rejectionSuccess {
        fmt.Printf("Vote %s rejected due to double use of SYN900 token %s.\n", vote.ID, vote.Syn900Token)
        automation.logDoubleVoteRejection(vote)
    } else {
        fmt.Printf("Error rejecting vote %s for SYN900 token %s.\n", vote.ID, vote.Syn900Token)
    }
}

// finalizeCheckCycle finalizes the double voting check cycle and logs the result in the ledger
func (automation *Syn900VotingCheckAutomation) finalizeCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDoubleVoteCheckCycle()
    if success {
        fmt.Println("Double voting check cycle finalized successfully.")
        automation.logCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing double voting check cycle.")
    }
}

// logDoubleVoteRejection logs the rejection of a vote due to double use of a SYN900 token
func (automation *Syn900VotingCheckAutomation) logDoubleVoteRejection(vote common.Vote) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("double-vote-rejection-%s", vote.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Double Vote Rejection",
        Status:    "Rejected",
        Details:   fmt.Sprintf("Vote %s rejected due to double use of SYN900 token %s.", vote.ID, vote.Syn900Token),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with rejection of vote %s for SYN900 token %s.\n", vote.ID, vote.Syn900Token)
}

// logCheckCycleFinalization logs the finalization of a double voting check cycle into the ledger
func (automation *Syn900VotingCheckAutomation) logCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("double-vote-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Double Vote Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with double voting check cycle finalization.")
}

// AddEncryptionToVoteData encrypts the vote data before logging or rejecting the vote
func (automation *Syn900VotingCheckAutomation) AddEncryptionToVoteData(vote common.Vote) common.Vote {
    encryptedData, err := encryption.EncryptData(vote)
    if err != nil {
        fmt.Println("Error encrypting vote data:", err)
        return vote
    }
    vote.EncryptedData = encryptedData
    fmt.Println("Vote data successfully encrypted.")
    return vote
}

// ensureVoteDataIntegrity checks the integrity of vote data and triggers double voting checks if necessary
func (automation *Syn900VotingCheckAutomation) ensureVoteDataIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVoteDataIntegrity()
    if !integrityValid {
        fmt.Println("Vote data integrity breach detected. Re-triggering double voting check.")
        automation.monitorAndCheckDoubleVoting()
    } else {
        fmt.Println("Vote data integrity is valid.")
    }
}
