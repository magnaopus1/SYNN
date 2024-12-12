package governance

import (
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)





// ExecutionRecord keeps track of the governance actions to be executed after voting
type ExecutionRecord struct {
	ProposalID string    // ID of the proposal to be executed
	Executed   bool      // Whether the execution was successful
	Timestamp  time.Time // Time when the execution happened
}

// GovernanceExecution manages the execution of approved proposals
type GovernanceExecution struct {
	ExecutionQueue []ExecutionRecord // Queue of proposals to be executed
	mutex          sync.Mutex        // Mutex for thread-safe operations
	LedgerInstance *ledger.Ledger    // Ledger instance for tracking executed proposals
}


// GovernanceTimelock defines the timelock contract for governance
type GovernanceTimelock struct {
	PendingExecutions map[string]*TimelockExecution // Map of proposal ID to pending executions
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	LedgerInstance    *ledger.Ledger                // Ledger instance for logging execution
}

// TimelockExecution represents a proposal pending execution after timelock
type TimelockExecution struct {
	ProposalID        string    // ID of the proposal
	ExecutionTime     time.Time // Time when the proposal can be executed
	EncryptedProposal string    // Encrypted details of the proposal
	Creator           string    // Address of the proposer
}


// GovernanceTracking manages the tracking of governance proposals and generates reports
type GovernanceTracking struct {
	ProposalHistory map[string]*GovernanceProposalStatus // Map of proposal ID to status
	LedgerInstance  *ledger.Ledger                      // Instance of the ledger
	mutex           sync.Mutex                          // Mutex for thread-safe operations
}

// GovernanceProposalStatus tracks the status of a single proposal
type GovernanceProposalStatus struct {
	ProposalID      string    // ID of the proposal
	Status          string    // Current status (e.g., "Pending", "Approved", "Executed", "Rejected")
	Timestamps      []time.Time // Important timestamps (e.g., creation, approval, execution)
	EncryptedDetails string    // Encrypted details of the proposal
}


// ReputationVoting represents the structure for reputation-based voting
type ReputationVoting struct {
	Votes          map[string]map[string]float64 // Map[proposalID]map[voterID]reputationScore
	LedgerInstance *ledger.Ledger                // Ledger to store voting records
	mutex          sync.Mutex                    // Mutex for thread-safe operations
}
