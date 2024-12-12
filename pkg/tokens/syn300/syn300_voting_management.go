package syn300

import (
    "sync"
    "fmt"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
)

// SYN300Token represents a governance token in the Synthron ecosystem with advanced voting management functionalities.
type SYN300Token struct {
    ID                  string
    DelegationHistory   map[string][]DelegationRecord // UserID -> list of DelegationRecords
    VotingPower         map[string]float64            // UserID -> Voting Power
    VotingFeedback      map[string][]FeedbackRecord   // ProposalID -> list of FeedbackRecords
    ActiveProposals     map[string]GovernanceProposal // ProposalID -> GovernanceProposal
    ArchivedProposals   map[string]GovernanceProposal // Archived proposals for historical reference
    CommentsEnabled     bool
    VotingPeriodStatus  VotingPeriodStatus
    mutex               sync.RWMutex
    Ledger              *ledger.Ledger
}

// DelegationRecord represents a record of delegation actions.
type DelegationRecord struct {
    Delegate    string
    Delegator   string
    VotingPower float64
    Timestamp   time.Time
}

// FeedbackRecord stores feedback submitted by voters on proposals.
type FeedbackRecord struct {
    UserID    string
    Feedback  string
    Timestamp time.Time
}

// VotingPeriodStatus stores the status of the current voting period.
type VotingPeriodStatus struct {
    Active      bool
    StartTime   time.Time
    EndTime     time.Time
}

// FETCH_DELEGATION_HISTORY retrieves the delegation history for a specific user.
func (token *SYN300Token) FETCH_DELEGATION_HISTORY(userID string) ([]DelegationRecord, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    history, exists := token.DelegationHistory[userID]
    if !exists {
        return nil, fmt.Errorf("no delegation history found for user %s", userID)
    }
    return history, nil
}

// GET_VOTING_POWER_DETAILS retrieves the voting power for a specific user.
func (token *SYN300Token) GET_VOTING_POWER_DETAILS(userID string) (float64, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    power, exists := token.VotingPower[userID]
    if !exists {
        return 0, fmt.Errorf("voting power details not found for user %s", userID)
    }
    return power, nil
}

// UPDATE_VOTING_POWER updates the voting power for a specific user.
func (token *SYN300Token) UPDATE_VOTING_POWER(userID string, newPower float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.VotingPower[userID] = newPower
    return token.Ledger.RecordLog("VotingPowerUpdated", fmt.Sprintf("Voting power for user %s updated to %.2f", userID, newPower))
}

// LOG_VOTE_CAST logs a vote cast by a user on a proposal.
func (token *SYN300Token) LOG_VOTE_CAST(userID, proposalID, vote string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.RecordLog("VoteCast", fmt.Sprintf("User %s cast vote %s on proposal %s", userID, vote, proposalID))
}

// SUBMIT_FEEDBACK_ON_PROPOSAL allows a user to submit feedback on a specific proposal.
func (token *SYN300Token) SUBMIT_FEEDBACK_ON_PROPOSAL(userID, proposalID, feedback string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    record := FeedbackRecord{
        UserID:    userID,
        Feedback:  feedback,
        Timestamp: time.Now(),
    }
    token.VotingFeedback[proposalID] = append(token.VotingFeedback[proposalID], record)
    return token.Ledger.RecordLog("FeedbackSubmitted", fmt.Sprintf("Feedback submitted by user %s on proposal %s", userID, proposalID))
}

// GET_FEEDBACK_FOR_PROPOSAL retrieves all feedback for a given proposal.
func (token *SYN300Token) GET_FEEDBACK_FOR_PROPOSAL(proposalID string) ([]FeedbackRecord, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    feedback, exists := token.VotingFeedback[proposalID]
    if !exists {
        return nil, fmt.Errorf("no feedback found for proposal %s", proposalID)
    }
    return feedback, nil
}

// CHECK_ACTIVE_PROPOSALS checks and returns all active proposals.
func (token *SYN300Token) CHECK_ACTIVE_PROPOSALS() []GovernanceProposal {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    var activeProposals []GovernanceProposal
    for _, proposal := range token.ActiveProposals {
        activeProposals = append(activeProposals, proposal)
    }
    return activeProposals
}

// ARCHIVE_INACTIVE_PROPOSALS moves proposals that have ended into the archive.
func (token *SYN300Token) ARCHIVE_INACTIVE_PROPOSALS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for id, proposal := range token.ActiveProposals {
        if time.Now().After(proposal.ExpiresAt) {
            token.ArchivedProposals[id] = proposal
            delete(token.ActiveProposals, id)
        }
    }
    return token.Ledger.RecordLog("ProposalsArchived", "Inactive proposals archived successfully")
}

// VIEW_ARCHIVED_PROPOSALS retrieves all archived proposals.
func (token *SYN300Token) VIEW_ARCHIVED_PROPOSALS() []GovernanceProposal {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    var archivedProposals []GovernanceProposal
    for _, proposal := range token.ArchivedProposals {
        archivedProposals = append(archivedProposals, proposal)
    }
    return archivedProposals
}

// ENABLE_COMMENTS_ON_PROPOSALS enables comments on proposals for voter interaction.
func (token *SYN300Token) ENABLE_COMMENTS_ON_PROPOSALS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.CommentsEnabled = true
    return token.Ledger.RecordLog("CommentsEnabled", "Comments on proposals enabled")
}

// DISABLE_COMMENTS_ON_PROPOSALS disables comments on proposals.
func (token *SYN300Token) DISABLE_COMMENTS_ON_PROPOSALS() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.CommentsEnabled = false
    return token.Ledger.RecordLog("CommentsDisabled", "Comments on proposals disabled")
}

// INITIATE_VOTING_PERIOD initiates a voting period for active proposals.
func (token *SYN300Token) INITIATE_VOTING_PERIOD(startTime, endTime time.Time) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.VotingPeriodStatus = VotingPeriodStatus{
        Active:    true,
        StartTime: startTime,
        EndTime:   endTime,
    }
    return token.Ledger.RecordLog("VotingPeriodInitiated", fmt.Sprintf("Voting period initiated from %s to %s", startTime, endTime))
}

// END_VOTING_PERIOD ends the current voting period for all proposals.
func (token *SYN300Token) END_VOTING_PERIOD() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.VotingPeriodStatus.Active = false
    return token.Ledger.RecordLog("VotingPeriodEnded", "Current voting period ended")
}

// GET_VOTING_PERIOD_STATUS retrieves the status of the current voting period.
func (token *SYN300Token) GET_VOTING_PERIOD_STATUS() VotingPeriodStatus {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.VotingPeriodStatus
}
