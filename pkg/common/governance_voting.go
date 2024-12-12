package common

import (
    "fmt"
    "sync"
    "synnergy_network/pkg/ledger"
)


// GovernanceVoting handles the voting mechanism for governance proposals
type GovernanceVoting struct {
    Votes          map[string]map[string]bool // Map[proposalID]map[voterID]bool
    LedgerInstance *ledger.Ledger             // Ledger to store voting records
    mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// NewGovernanceVoting initializes a new GovernanceVoting system
func NewGovernanceVoting(ledgerInstance *ledger.Ledger) *GovernanceVoting {
    return &GovernanceVoting{
        Votes:          make(map[string]map[string]bool),
        LedgerInstance: ledgerInstance,
    }
}

func (gv *GovernanceVoting) CastVote(proposalID, voterID string, syn900Token *SYN900Token) error {
    gv.mutex.Lock()
    defer gv.mutex.Unlock()

    // Check if the Syn-900 token is valid and matches the voter ID
    if !gv.verifySyn900Token(syn900Token, voterID) {
        return fmt.Errorf("invalid Syn-900 token for voter %s", voterID)
    }

    // Check if the user has already voted on this proposal
    if gv.hasVoted(proposalID, voterID) {
        return fmt.Errorf("voter %s has already voted on proposal %s", voterID, proposalID)
    }

    // Record the vote
    if gv.Votes[proposalID] == nil {
        gv.Votes[proposalID] = make(map[string]bool)
    }
    gv.Votes[proposalID][voterID] = true

    // Log the vote in the ledger
    gv.logVoteToLedger(proposalID, voterID)

    // Destroy the Syn-900 ID token after a successful vote
    err := gv.destroySyn900Token(syn900Token)
    if err != nil {
        return fmt.Errorf("failed to destroy Syn-900 token: %v", err)
    }

    fmt.Printf("Vote cast by voter %s on proposal %s.\n", voterID, proposalID)
    return nil
}


// Verify Syn-900 token to ensure it's valid
func (gv *GovernanceVoting) verifySyn900Token(syn900Token *SYN900Token, voterID string) bool {
    return syn900Token.Owner == voterID && syn900Token.Status == "active"
}

// hasVoted checks if a voter has already voted on a proposal
func (gv *GovernanceVoting) hasVoted(proposalID, voterID string) bool {
    if gv.Votes[proposalID] == nil {
        return false
    }
    return gv.Votes[proposalID][voterID]
}

// logVoteToLedger logs the vote to the ledger for transparency
func (gv *GovernanceVoting) logVoteToLedger(proposalID, voterID string) error {
    // Create an instance of Encryption
    encryptionInstance := &Encryption{}

    // Create the vote record
    voteRecord := fmt.Sprintf("Voter %s voted on proposal %s", voterID, proposalID)

    // Encrypt the vote record using AES
    encryptedRecord, err := encryptionInstance.EncryptData("AES", []byte(voteRecord), EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt vote record: %v", err)
    }

    // Convert the encrypted record to a string and record it in the ledger
    err = gv.LedgerInstance.GovernanceLedger.RecordVote(proposalID, voterID, string(encryptedRecord))
    if err != nil {
        return fmt.Errorf("failed to log vote to ledger: %v", err)
    }

    fmt.Printf("Vote by %s on proposal %s logged to ledger.\n", voterID, proposalID)
    return nil
}


// Destroy the Syn-900 token after use
func (gv *GovernanceVoting) destroySyn900Token(syn900Token *SYN900Token) error {
    syn900Token.Status = "destroyed"
    fmt.Printf("Syn-900 token for voter %s destroyed.\n", syn900Token.Owner)
    return nil
}

// CountVotes returns the total number of votes for a proposal
func (gv *GovernanceVoting) CountVotes(proposalID string) (int, error) {
    gv.mutex.Lock()
    defer gv.mutex.Unlock()

    if gv.Votes[proposalID] == nil {
        return 0, fmt.Errorf("no votes found for proposal %s", proposalID)
    }

    voteCount := len(gv.Votes[proposalID])
    return voteCount, nil
}

// GetVoters returns a list of voter IDs who voted on a proposal
func (gv *GovernanceVoting) GetVoters(proposalID string) ([]string, error) {
    gv.mutex.Lock()
    defer gv.mutex.Unlock()

    if gv.Votes[proposalID] == nil {
        return nil, fmt.Errorf("no votes found for proposal %s", proposalID)
    }

    voters := make([]string, 0, len(gv.Votes[proposalID]))
    for voterID := range gv.Votes[proposalID] {
        voters = append(voters, voterID)
    }

    return voters, nil
}
