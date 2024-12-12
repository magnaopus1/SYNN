package syn300

import (
    "sync"
    "fmt"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
    "math/big"
    "time"
)

// SYN300Token represents a governance token with full token management and voting functionality.
type SYN300Token struct {
    ID                      string
    BalanceOf               map[string]*big.Int // UserID -> token balance
    Metadata                SYN300TokenMetadata
    VotingPowerDelegations  map[string]string   // UserID -> DelegateUserID
    TransparentVoting       bool
    Proposals               map[string]GovernanceProposal
    mutex                   sync.RWMutex
    Ledger                  *ledger.Ledger
}

// SYN300TokenMetadata contains essential metadata for the governance token.
type SYN300TokenMetadata struct {
    Name        string
    Symbol      string
    Decimals    uint8
    TotalSupply *big.Int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// TRANSFER_GOVERNANCE_TOKEN transfers governance tokens from one user to another.
func (token *SYN300Token) TRANSFER_GOVERNANCE_TOKEN(from, to string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.BalanceOf[from].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance for user %s", from)
    }
    token.BalanceOf[from].Sub(token.BalanceOf[from], amount)
    if _, exists := token.BalanceOf[to]; !exists {
        token.BalanceOf[to] = big.NewInt(0)
    }
    token.BalanceOf[to].Add(token.BalanceOf[to], amount)
    return token.Ledger.RecordLog("TransferGovernanceToken", fmt.Sprintf("Transferred %s tokens from %s to %s", amount.String(), from, to))
}

// CHECK_GOVERNANCE_TOKEN_BALANCE retrieves the balance of governance tokens for a specific user.
func (token *SYN300Token) CHECK_GOVERNANCE_TOKEN_BALANCE(userID string) *big.Int {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    if balance, exists := token.BalanceOf[userID]; exists {
        return balance
    }
    return big.NewInt(0)
}

// GET_GOVERNANCE_TOKEN_METADATA retrieves metadata for the governance token.
func (token *SYN300Token) GET_GOVERNANCE_TOKEN_METADATA() SYN300TokenMetadata {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Metadata
}

// UPDATE_GOVERNANCE_TOKEN_METADATA updates the metadata of the governance token.
func (token *SYN300Token) UPDATE_GOVERNANCE_TOKEN_METADATA(metadata SYN300TokenMetadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata = metadata
    token.Metadata.UpdatedAt = time.Now()
    return token.Ledger.RecordLog("GovernanceTokenMetadataUpdated", "Governance token metadata updated successfully")
}

// DELEGATE_VOTING_POWER delegates voting power from one user to another.
func (token *SYN300Token) DELEGATE_VOTING_POWER(from, to string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.VotingPowerDelegations[from] = to
    return token.Ledger.RecordLog("VotingPowerDelegated", fmt.Sprintf("User %s delegated voting power to %s", from, to))
}

// REVOKE_DELEGATION revokes voting power delegation for a user.
func (token *SYN300Token) REVOKE_DELEGATION(from string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if _, exists := token.VotingPowerDelegations[from]; !exists {
        return fmt.Errorf("no delegation found for user %s", from)
    }
    delete(token.VotingPowerDelegations, from)
    return token.Ledger.RecordLog("DelegationRevoked", fmt.Sprintf("Delegation revoked for user %s", from))
}

// CAST_VOTE allows a user to cast a vote on a proposal.
func (token *SYN300Token) CAST_VOTE(userID, proposalID, vote string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    proposal, exists := token.Proposals[proposalID]
    if !exists {
        return fmt.Errorf("proposal %s not found", proposalID)
    }
    proposal.Votes[userID] = vote
    token.Proposals[proposalID] = proposal
    return token.Ledger.RecordLog("VoteCast", fmt.Sprintf("User %s cast vote on proposal %s", userID, proposalID))
}

// GET_VOTE_WEIGHT retrieves the vote weight for a user based on their token balance.
func (token *SYN300Token) GET_VOTE_WEIGHT(userID string) *big.Int {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    if balance, exists := token.BalanceOf[userID]; exists {
        return balance
    }
    return big.NewInt(0)
}

// SUBMIT_PROPOSAL submits a new proposal for governance voting.
func (token *SYN300Token) SUBMIT_PROPOSAL(proposalID, title, description string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    proposal := GovernanceProposal{
        ProposalID:  proposalID,
        Title:       title,
        Description: description,
        Status:      "Active",
        CreatedAt:   time.Now(),
        Votes:       make(map[string]string),
    }
    token.Proposals[proposalID] = proposal
    return token.Ledger.RecordLog("ProposalSubmitted", fmt.Sprintf("Proposal %s submitted with title '%s'", proposalID, title))
}

// TRACK_PROPOSAL_STATUS retrieves the current status of a proposal.
func (token *SYN300Token) TRACK_PROPOSAL_STATUS(proposalID string) (string, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    proposal, exists := token.Proposals[proposalID]
    if !exists {
        return "", fmt.Errorf("proposal %s not found", proposalID)
    }
    return proposal.Status, nil
}

// ENABLE_TRANSPARENT_VOTING_RECORDS enables transparent voting records, making votes visible to all participants.
func (token *SYN300Token) ENABLE_TRANSPARENT_VOTING_RECORDS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TransparentVoting = true
    return token.Ledger.RecordLog("TransparentVotingEnabled", "Transparent voting records enabled")
}

// DISABLE_TRANSPARENT_VOTING_RECORDS disables transparent voting records, hiding individual votes.
func (token *SYN300Token) DISABLE_TRANSPARENT_VOTING_RECORDS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TransparentVoting = false
    return token.Ledger.RecordLog("TransparentVotingDisabled", "Transparent voting records disabled")
}
