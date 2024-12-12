package automations

import (
	"fmt"
	"sync"
	"time"
	"math"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	VotingCheckInterval         = 15 * time.Second // Interval for checking quadratic voting compliance
	MaxAllowedVotingViolations  = 3                // Maximum number of quadratic voting violations allowed per user
)

// QuadraticVotingRestrictionAutomation monitors and enforces quadratic voting rules across the network
type QuadraticVotingRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	votingViolationCount   map[string]int // Tracks voting violation counts per user
	voteCredits            map[string]float64 // Tracks available vote credits for each user
}

// NewQuadraticVotingRestrictionAutomation initializes and returns an instance of QuadraticVotingRestrictionAutomation
func NewQuadraticVotingRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *QuadraticVotingRestrictionAutomation {
	return &QuadraticVotingRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		votingViolationCount: make(map[string]int),
		voteCredits:          make(map[string]float64),
	}
}

// StartVotingMonitoring starts continuous monitoring of quadratic voting compliance
func (automation *QuadraticVotingRestrictionAutomation) StartVotingMonitoring() {
	ticker := time.NewTicker(VotingCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorVotingCompliance()
		}
	}()
}

// monitorVotingCompliance checks for quadratic voting violations and enforces restrictions if necessary
func (automation *QuadraticVotingRestrictionAutomation) monitorVotingCompliance() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch voting data from Synnergy Consensus
	votingData := automation.consensusSystem.GetVotingData()

	for userID, votesCast := range votingData {
		// Calculate the number of credits used based on quadratic voting
		creditsUsed := math.Pow(float64(votesCast), 2)

		// Check if the user has enough vote credits to cast the votes
		if automation.voteCredits[userID] < creditsUsed {
			automation.flagVotingViolation(userID, votesCast, "Insufficient vote credits for quadratic voting")
		} else {
			// Deduct the credits from the user's available vote credits
			automation.voteCredits[userID] -= creditsUsed
		}
	}
}

// flagVotingViolation flags a quadratic voting violation and logs it in the ledger
func (automation *QuadraticVotingRestrictionAutomation) flagVotingViolation(userID string, votesCast int, reason string) {
	fmt.Printf("Quadratic voting violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.votingViolationCount[userID]++

	// Log the violation in the ledger
	automation.logVotingViolation(userID, votesCast, reason)
}

// logVotingViolation logs the flagged quadratic voting violation into the ledger with full details
func (automation *QuadraticVotingRestrictionAutomation) logVotingViolation(userID string, votesCast int, violationReason string) {
	// Create a ledger entry for voting violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("quadratic-voting-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Quadratic Voting Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated quadratic voting rules. Votes Cast: %d. Reason: %s", userID, votesCast, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptVotingData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log quadratic voting violation:", err)
	} else {
		fmt.Println("Quadratic voting violation logged.")
	}
}

// encryptVotingData encrypts the voting data before logging for security
func (automation *QuadraticVotingRestrictionAutomation) encryptVotingData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting voting data:", err)
		return data
	}
	return string(encryptedData)
}

// resetVoteCredits resets the user's vote credits for a new voting cycle
func (automation *QuadraticVotingRestrictionAutomation) resetVoteCredits(userID string, credits float64) {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	automation.voteCredits[userID] = credits
	fmt.Printf("Vote credits for user %s have been reset to %.2f\n", userID, credits)
}
