package community_engagement

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// Poll struct to store information about polls
type Poll struct {
	ID        string            // Unique identifier for the poll
	CreatorID string            // ID of the user who created the poll
	Question  string            // Poll question
	Options   []string          // Options for voting
	Expiry    time.Time         // Expiration time of the poll
	Open      bool              // Status indicating if the poll is open for voting
	Votes     map[string]int    // Vote counts for each option
	VoterList map[string]string // Records the option each user voted for to prevent double-voting
	CreatedAt time.Time         // Timestamp when the poll was created
}

// createPoll initializes a new poll with options and records it in the ledger
func CreatePoll(creatorID, question string, options []string, expiry time.Time) (string, error) {
	// Ensure at least two options are provided
	if len(options) < 2 {
		return "", errors.New("at least two options are required to create a poll")
	}

	// Generate a unique ID for the poll
	pollID := generateUniqueID(creatorID + question)

	// Create the poll struct
	poll := ledger.Poll{
		ID:        pollID,
		CreatorID: creatorID,
		Question:  question,
		Options:   options,
		Expiry:    expiry,
		Open:      true,
		Votes:     make(map[string]int),
		VoterList: make(map[string]string),
		CreatedAt: time.Now(),
	}

	// Initialize vote counts for each option
	for _, option := range options {
		poll.Votes[option] = 0
	}

	// Record the poll in the ledger
	l := &ledger.Ledger{}
	if err := l.CommunityEngagementLedger.RecordPoll(poll); err != nil {
		return "", fmt.Errorf("failed to record poll: %v", err)
	}
	return pollID, nil
}

// voteInPoll allows a user to cast a vote in an open poll
func VoteInPoll(pollID, userID, option string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the poll
	poll, err := l.CommunityEngagementLedger.FetchPoll(pollID)
	if err != nil {
		return fmt.Errorf("poll not found: %v", err)
	}

	// Check if the poll is open and not expired
	if !poll.Open || time.Now().After(poll.Expiry) {
		return errors.New("poll is closed or expired")
	}

	// Check if the option exists in the poll
	if _, exists := poll.Votes[option]; !exists {
		return errors.New("invalid option selected")
	}

	// Check if the user has already voted
	if l.CommunityEngagementLedger.HasUserVoted(pollID, userID) {
		return errors.New("user has already voted in this poll")
	}

	// Record the vote in the ledger
	if err := l.CommunityEngagementLedger.RecordVote(pollID, userID, option); err != nil {
		return fmt.Errorf("failed to record vote: %v", err)
	}

	return nil
}

// closePoll closes a poll, preventing any further votes
func ClosePoll(adminID, pollID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the poll
	poll, err := l.CommunityEngagementLedger.FetchPoll(pollID)
	if err != nil {
		return fmt.Errorf("poll not found: %v", err)
	}

	// Check if the poll is already closed
	if !poll.Open {
		return errors.New("poll is already closed")
	}

	// Close the poll
	poll.Open = false
	if err := l.CommunityEngagementLedger.UpdatePollStatus(poll); err != nil {
		return fmt.Errorf("failed to close poll: %v", err)
	}

	return nil
}

// listAllPolls retrieves a list of all polls from the ledger
func ListAllPolls() ([]ledger.Poll, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve all polls
	polls, err := l.CommunityEngagementLedger.ListAllPolls()
	if err != nil {
		return nil, fmt.Errorf("failed to list all polls: %v", err)
	}
	return polls, nil
}

// viewPollResults retrieves the current results of a specified poll
func ViewPollResults(pollID string) (map[string]int, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the specified poll
	poll, err := l.CommunityEngagementLedger.FetchPoll(pollID)
	if err != nil {
		return nil, fmt.Errorf("poll not found: %v", err)
	}

	return poll.Votes, nil
}
