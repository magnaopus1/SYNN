package governance

import (
	"encoding/base64"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewReputationVoting initializes a new ReputationVoting system
func NewReputationVoting(ledgerInstance *ledger.Ledger) *ReputationVoting {
    return &ReputationVoting{
        Votes:          make(map[string]map[string]float64),
        LedgerInstance: ledgerInstance,
    }
}

// CastReputationVote allows a user to cast their reputation-based vote using their Syn-900 ID token
func (rv *ReputationVoting) CastReputationVote(proposalID, voterID string, reputationScore float64, syn900Token *common.SYN900Token, burnManager *common.Syn900BurnManager, ledgerInstance *ledger.Ledger) error {
    rv.mutex.Lock()
    defer rv.mutex.Unlock()

    // Check if the Syn-900 token is valid and matches the voter's ID
    isValid, err := rv.verifySyn900Token(syn900Token, voterID, ledgerInstance) // Pass token by reference and ledgerInstance
    if err != nil {
        return fmt.Errorf("error verifying Syn-900 token: %v", err)
    }
    if !isValid {
        return fmt.Errorf("invalid Syn-900 token for voter %s", voterID)
    }

    // Check if the voter has already voted on this proposal
    if rv.hasVoted(proposalID, voterID) {
        return fmt.Errorf("voter %s has already cast their vote for proposal %s", voterID, proposalID)
    }

    // Record the vote
    if rv.Votes[proposalID] == nil {
        rv.Votes[proposalID] = make(map[string]float64)
    }
    rv.Votes[proposalID][voterID] = reputationScore

    // Log the vote in the ledger
    rv.logVoteToLedger(proposalID, voterID, reputationScore)

    // Destroy the Syn-900 ID token after successful vote casting (no need for burnManager)
    err = rv.destroySyn900Token(syn900Token) // Only pass the token by reference
    if err != nil {
        return fmt.Errorf("failed to destroy Syn-900 token: %v", err)
    }

    fmt.Printf("Voter %s cast their reputation-based vote for proposal %s with a reputation score of %.2f.\n", voterID, proposalID, reputationScore)
    return nil
}



// destroySyn900Token ensures that the Syn-900 ID token is disposed of and logs the action after voting ends
func (rv *ReputationVoting) destroySyn900Token(token *common.SYN900Token) error {
    // Ensure rv.LedgerInstance is not nil
    if rv.LedgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    // Cast LedgerInstance to the appropriate ledger type
    ledgerInstance, ok := interface{}(rv.LedgerInstance).(*ledger.Ledger)
    if !ok {
        return fmt.Errorf("failed to cast LedgerInstance to *ledger.Ledger")
    }

    // Step 1: Check if the token is already burned by accessing TokenRecords, not TokenBalances
    tokenRecord, exists := ledgerInstance.TokenRecords[token.TokenID]
    if !exists {
        return fmt.Errorf("token %s not found in ledger", token.TokenID)
    }

    if tokenRecord.IsBurned {
        return fmt.Errorf("token %s has already been burned", token.TokenID)
    }

    // Step 2: Mark the token as burned
    tokenRecord.IsBurned = true // Mark the token as burned

    // Step 3: Optionally log the burn action
    voterID, err := GetVoterIDFromToken(*token)
    if err != nil {
        return fmt.Errorf("failed to retrieve VoterID from token: %v", err)
    }

    fmt.Printf("Syn-900 ID token for voter %s has been successfully burned and disposed after voting.\n", voterID)
    return nil
}




// verifySyn900Token ensures the Syn-900 ID token is valid for the voter
func (rv *ReputationVoting) verifySyn900Token(token *common.SYN900Token, voterID string, ledgerInstance *ledger.Ledger) (bool, error) {
    // Check if the token is valid using IsTokenValid
    isValid, err := IsTokenValid(*token, ledgerInstance) // Dereference token and pass ledgerInstance
    if err != nil || !isValid {
        return false, fmt.Errorf("invalid token: %v", err)
    }

    // Check if the token belongs to the voter
    associatedVoterID, err := GetVoterIDFromToken(*token) // Dereference token
    if err != nil {
        return false, fmt.Errorf("failed to retrieve voter ID from token: %v", err)
    }

    if associatedVoterID != voterID {
        return false, fmt.Errorf("token does not belong to the voter %s", voterID)
    }

    return true, nil
}



// hasVoted checks if a voter has already cast their vote on a proposal
func (rv *ReputationVoting) hasVoted(proposalID, voterID string) bool {
    if rv.Votes[proposalID] == nil {
        return false
    }
    if _, voted := rv.Votes[proposalID][voterID]; voted {
        return true
    }
    return false
}

// logVoteToLedger logs the vote to the ledger for transparency
func (rv *ReputationVoting) logVoteToLedger(proposalID, voterID string, reputationScore float64) error {
    // Create an instance of the Encryption struct
    encryption := &common.Encryption{}

    // Prepare the vote record to be encrypted
    voteRecord := fmt.Sprintf("Voter %s cast a reputation-based vote with a score of %.2f for proposal %s", voterID, reputationScore, proposalID)
    voteRecordBytes := []byte(voteRecord) // Convert vote record to byte slice

    // Encrypt the vote record using AES encryption
    encryptedRecord, err := encryption.EncryptData("AES", voteRecordBytes, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt vote record: %v", err)
    }

    // Convert the encryptedRecord (byte slice) to a string using base64 encoding
    encryptedRecordStr := base64.StdEncoding.EncodeToString(encryptedRecord)

    // Log the encrypted vote to the ledger with three string arguments
    err = rv.LedgerInstance.RecordVote(proposalID, voterID, encryptedRecordStr)
    if err != nil {
        return fmt.Errorf("failed to log vote to ledger: %v", err)
    }

    fmt.Printf("Reputation vote by %s for proposal %s logged to ledger.\n", voterID, proposalID)
    return nil
}








// CountReputationVotes returns the total reputation score for a proposal
func (rv *ReputationVoting) CountReputationVotes(proposalID string) (float64, error) {
    rv.mutex.Lock()
    defer rv.mutex.Unlock()

    if rv.Votes[proposalID] == nil {
        return 0, fmt.Errorf("no votes found for proposal %s", proposalID)
    }

    totalReputationScore := 0.0
    for _, score := range rv.Votes[proposalID] {
        totalReputationScore += score
    }

    return totalReputationScore, nil
}

// GetVoters returns a list of voter IDs for a proposal
func (rv *ReputationVoting) GetVoters(proposalID string) ([]string, error) {
    rv.mutex.Lock()
    defer rv.mutex.Unlock()

    if rv.Votes[proposalID] == nil {
        return nil, fmt.Errorf("no votes found for proposal %s", proposalID)
    }

    voters := make([]string, 0, len(rv.Votes[proposalID]))
    for voterID := range rv.Votes[proposalID] {
        voters = append(voters, voterID)
    }

    return voters, nil
}
