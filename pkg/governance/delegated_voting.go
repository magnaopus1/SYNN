package governance

import (
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// DelegatedVoting represents the structure for delegated voting
type DelegatedVoting struct {
	Votes          map[string]map[string]bool // Map[proposalID]map[delegatorID]bool
	Delegations    map[string]string          // Map[delegatorID]delegateeID (who the vote was delegated to)
	LedgerInstance *ledger.Ledger             // Ledger to store delegation and voting records
	mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// NewDelegatedVoting initializes a new DelegatedVoting system
func NewDelegatedVoting(ledgerInstance *ledger.Ledger) *DelegatedVoting {
    return &DelegatedVoting{
        Votes:          make(map[string]map[string]bool),
        Delegations:    make(map[string]string),
        LedgerInstance: ledgerInstance,
    }
}

// DelegateVote allows a user to delegate their vote to another user using their Syn-900 ID token
func (dv *DelegatedVoting) DelegateVote(proposalID, delegatorID, delegateeID string, syn900Token *common.SYN900Token, burnManager *common.Syn900BurnManager, ledgerInstance *ledger.Ledger) error {
    dv.mutex.Lock()
    defer dv.mutex.Unlock()

    // Check if the Syn-900 token is valid and belongs to the delegator
    isValid, err := dv.verifySyn900Token(syn900Token, delegatorID, ledgerInstance) // Pass the token by reference and add ledgerInstance
    if err != nil || !isValid {
        return fmt.Errorf("invalid Syn-900 token for delegator %s: %v", delegatorID, err)
    }

    // Check if the delegator has already delegated or voted on this proposal
    if dv.hasVotedOrDelegated(proposalID, delegatorID) {
        return fmt.Errorf("delegator %s has already voted or delegated for proposal %s", delegatorID, proposalID)
    }

    // Record the delegation
    dv.Delegations[delegatorID] = delegateeID
    dv.logDelegationToLedger(proposalID, delegatorID, delegateeID)

    // Destroy the Syn-900 ID token after successful delegation
    err = dv.destroySyn900Token(syn900Token) // Only pass the token, no burn manager
    if err != nil {
        return fmt.Errorf("failed to destroy Syn-900 token: %v", err)
    }

    fmt.Printf("Delegator %s delegated their vote to %s for proposal %s.\n", delegatorID, delegateeID, proposalID)
    return nil
}






// verifySyn900Token verifies if the SYN900 token is valid and belongs to the delegator
func (dv *DelegatedVoting) verifySyn900Token(token *common.SYN900Token, delegatorID string, ledgerInstance *ledger.Ledger) (bool, error) {
    // Check if the token is valid using IsTokenValid
    isValid, err := IsTokenValid(*token, ledgerInstance) // Dereference token to pass it by value
    if err != nil || !isValid {
        return false, fmt.Errorf("invalid token: %v", err)
    }

    // Check if the token belongs to the delegator
    voterID, err := GetVoterIDFromToken(*token) // Dereference token to get its value
    if err != nil {
        return false, fmt.Errorf("failed to retrieve voter ID from token: %v", err)
    }

    if voterID != delegatorID {
        return false, fmt.Errorf("token does not belong to the delegator %s", delegatorID)
    }

    return true, nil
}



// CastDelegatedVote allows the delegatee to cast the vote for a proposal using the delegated authority
func (dv *DelegatedVoting) CastDelegatedVote(proposalID, delegateeID string) error {
    dv.mutex.Lock()
    defer dv.mutex.Unlock()

    // Check if the delegatee has received a delegated vote
    for delegatorID, assignedDelegateeID := range dv.Delegations {
        if assignedDelegateeID == delegateeID {
            // Check if the delegatee has already voted for this proposal
            if dv.hasVoted(proposalID, delegateeID) {
                return fmt.Errorf("delegatee %s has already voted on behalf of delegator %s", delegateeID, delegatorID)
            }

            // Record the vote in place of the delegator
            if dv.Votes[proposalID] == nil {
                dv.Votes[proposalID] = make(map[string]bool)
            }
            dv.Votes[proposalID][delegatorID] = true

            // Log the vote in the ledger
            dv.logVoteToLedger(proposalID, delegatorID, delegateeID)
            fmt.Printf("Delegatee %s voted on behalf of delegator %s for proposal %s.\n", delegateeID, delegatorID, proposalID)

            return nil
        }
    }

    return fmt.Errorf("no delegation found for delegatee %s on proposal %s", delegateeID, proposalID)
}

// CastVoteWithToken allows the user to cast a vote using their SYN900 token
func (dv *DelegatedVoting) CastVoteWithToken(token *common.SYN900Token, proposalID string, ledgerInstance *ledger.Ledger) error {
    dv.mutex.Lock()
    defer dv.mutex.Unlock()

    // Check if the token is valid for voting
    isValid, err := IsTokenValid(*token, ledgerInstance) // Pass token by reference and ledgerInstance
    if err != nil {
        return fmt.Errorf("error validating token: %v", err)
    }

    if !isValid {
        return fmt.Errorf("invalid token, vote cannot be cast")
    }

    // Retrieve the VoterID associated with the token
    voterID, err := GetVoterIDFromToken(*token) // Dereference token pointer
    if err != nil {
        return fmt.Errorf("failed to retrieve voter ID: %v", err)
    }

    // Check if the voter has already voted for this proposal
    if dv.hasVoted(proposalID, voterID) {
        return fmt.Errorf("voter %s has already voted on proposal %s", voterID, proposalID)
    }

    // Record the vote in the Votes map
    if dv.Votes[proposalID] == nil {
        dv.Votes[proposalID] = make(map[string]bool)
    }
    dv.Votes[proposalID][voterID] = true

    // Log the vote in the ledger
    dv.logVoteToLedger(proposalID, voterID, "")

    fmt.Printf("Voter %s cast a vote on proposal %s using their token.\n", voterID, proposalID)

    return nil
}


// hasVotedOrDelegated checks if a delegator has already voted or delegated their vote on a proposal
func (dv *DelegatedVoting) hasVotedOrDelegated(proposalID, delegatorID string) bool {
    if dv.Votes[proposalID] == nil {
        return false
    }
    if dv.Votes[proposalID][delegatorID] || dv.Delegations[delegatorID] != "" {
        return true
    }
    return false
}

// logDelegationToLedger logs the delegation to the ledger for transparency
func (dv *DelegatedVoting) logDelegationToLedger(proposalID, delegatorID, delegateeID string) error {
    // Create an Encryption instance
    encryptionInstance := &common.Encryption{}

    // Prepare the delegation record
    delegationRecord := fmt.Sprintf("Delegator %s delegated vote to %s for proposal %s", delegatorID, delegateeID, proposalID)

    // Encrypt the record with AES algorithm
    encryptedRecord, err := encryptionInstance.EncryptData("AES", []byte(delegationRecord), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt delegation record: %v", err)
    }

    // Log the delegation to the ledger (pass only proposalID and the encrypted record)
    err = dv.LedgerInstance.RecordDelegation(proposalID, string(encryptedRecord))
    if err != nil {
        return fmt.Errorf("failed to log delegation to ledger: %v", err)
    }

    fmt.Printf("Delegation from %s to %s on proposal %s logged to ledger.\n", delegatorID, delegateeID, proposalID)
    return nil
}



// logVoteToLedger logs the delegated vote to the ledger for transparency
func (dv *DelegatedVoting) logVoteToLedger(proposalID, delegatorID, delegateeID string) error {
    // Create an Encryption instance
    encryptionInstance := &common.Encryption{}

    // Prepare the vote record
    voteRecord := fmt.Sprintf("Delegatee %s voted on behalf of delegator %s for proposal %s", delegateeID, delegatorID, proposalID)

    // Encrypt the record with AES algorithm
    encryptedRecord, err := encryptionInstance.EncryptData("AES", []byte(voteRecord), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt vote record: %v", err)
    }

    // Log the vote to the ledger
    err = dv.LedgerInstance.RecordVote(proposalID, delegatorID, string(encryptedRecord))
    if err != nil {
        return fmt.Errorf("failed to log vote to ledger: %v", err)
    }

    fmt.Printf("Vote by %s on behalf of %s for proposal %s logged to ledger.\n", delegateeID, delegatorID, proposalID)
    return nil
}

// destroySyn900Token ensures that the Syn-900 ID token is disposed of and cannot be reused after delegation
func (dv *DelegatedVoting) destroySyn900Token(token *common.SYN900Token) error {
    // Ensure dv.LedgerInstance is not nil
    if dv.LedgerInstance == nil {
        return fmt.Errorf("ledger instance is nil")
    }

    // Cast LedgerInstance to the appropriate ledger type
    ledgerInstance, ok := interface{}(dv.LedgerInstance).(*ledger.Ledger)
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

    fmt.Printf("Syn-900 ID token for delegator %s has been successfully burned after delegation.\n", voterID)
    return nil
}




// CountDelegatedVotes returns the total number of delegated votes for a proposal
func (dv *DelegatedVoting) CountDelegatedVotes(proposalID string) (int, error) {
    dv.mutex.Lock()
    defer dv.mutex.Unlock()

    if dv.Votes[proposalID] == nil {
        return 0, fmt.Errorf("no votes found for proposal %s", proposalID)
    }

    voteCount := len(dv.Votes[proposalID])
    return voteCount, nil
}

// GetDelegators returns a list of delegator IDs for a proposal
func (dv *DelegatedVoting) GetDelegators(proposalID string) ([]string, error) {
    dv.mutex.Lock()
    defer dv.mutex.Unlock()

    if dv.Votes[proposalID] == nil {
        return nil, fmt.Errorf("no votes found for proposal %s", proposalID)
    }

    delegators := make([]string, 0, len(dv.Votes[proposalID]))
    for delegatorID := range dv.Votes[proposalID] {
        delegators = append(delegators, delegatorID)
    }

    return delegators, nil
}

// hasVoted checks if the delegatee has already voted on behalf of a delegator for a specific proposal
func (dv *DelegatedVoting) hasVoted(proposalID, delegateeID string) bool {
    // Check if a vote exists for the given proposalID
    if votesForProposal, exists := dv.Votes[proposalID]; exists {
        // Check if any delegator has already voted through this delegateeID
        for delegatorID, voted := range votesForProposal {
            if voted && dv.Delegations[delegatorID] == delegateeID {
                return true
            }
        }
    }
    return false
}

// IsTokenValid checks if the token is valid for voting by interacting with the token's contract
func IsTokenValid(token common.SYN900Token, ledgerInstance *ledger.Ledger) (bool, error) {
    // Step 1: Check the token balance
    balance, err := token.CheckBalanceOnChain(ledgerInstance) // Pass ledgerInstance to CheckBalanceOnChain
    if err != nil {
        return false, fmt.Errorf("failed to check token balance: %v", err)
    }

    if balance <= 0 {
        return false, errors.New("token balance is zero or less, not valid for voting")
    }

    // Step 2: Check if the token is frozen
    frozen, err := token.IsFrozenOnChain(ledgerInstance) // Pass ledgerInstance to IsFrozenOnChain
    if err != nil {
        return false, fmt.Errorf("failed to check if token is frozen: %v", err)
    }

    if frozen {
        return false, errors.New("token is frozen, not valid for voting")
    }

    // Step 3: Check if the token is burned
    burned, err := token.IsBurnedOnChain(ledgerInstance) // Pass ledgerInstance to IsBurnedOnChain
    if err != nil {
        return false, fmt.Errorf("failed to check if token is burned: %v", err)
    }

    if burned {
        return false, errors.New("token is burned, not valid for voting")
    }

    // Step 4: Check if the token is expired
    expired, err := token.IsExpiredOnChain(ledgerInstance) // Pass ledgerInstance to IsExpiredOnChain
    if err != nil {
        return false, fmt.Errorf("failed to check if token is expired: %v", err)
    }

    if expired {
        return false, errors.New("token is expired, not valid for voting")
    }

    // If all conditions pass, the token is valid for voting
    return true, nil
}


// GetVoterIDFromToken retrieves the VoterID from the token's metadata stored in the struct
func GetVoterIDFromToken(token common.SYN900Token) (string, error) {
    // Step 1: Check if metadata exists
    if token.Metadata == nil {
        return "", errors.New("token metadata not found")
    }

    // Step 2: Check if the voterID (in this case, token.Owner or token.Metadata.Owner) exists
    voterID := token.Metadata.Owner
    if voterID == "" {
        return "", errors.New("VoterID not found in token metadata")
    }

    // Step 3: Return the voterID
    return voterID, nil
}
