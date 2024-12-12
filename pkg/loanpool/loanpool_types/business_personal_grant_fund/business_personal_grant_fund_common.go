package loanpool

import (
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// GrantApprovalProcess handles the two-stage approval process of a business personal grant proposal.
type BusinessPersonalGrantApprovalProcess struct {
	mutex             sync.Mutex                          // Mutex for thread safety
	Ledger            *ledger.Ledger                      // Reference to the ledger
	Consensus         *common.SynnergyConsensus       // Synnergy Consensus engine
	Proposals         map[string]*BusinessPersonalGrantProposalApproval   // Map to hold grant proposals by proposal ID
	AuthorityNodes    []common.AuthorityNodeTypes                   // List of valid authority node types (bank, government, central bank, etc.)
	PublicVotePeriod  time.Duration                       // Time allowed for public voting
	AuthorityVoteTime time.Duration                       // Time window for authority nodes to vote
}

// GrantProposalApproval represents a grant proposal along with its voting data
type BusinessPersonalGrantProposalApproval struct {
	Proposal          *BusinessPersonalGrantProposal                      // Reference to the grant proposal
	PublicVotes       map[string]bool                     // Map of public votes (address -> vote)
	Stage             ApprovalStage                      // Current approval stage
	AuthorityVotes    map[string]bool                     // Authority node votes
	VoteStartTime     time.Time                          // Time when voting starts
	ConfirmationCount int                                // Count of authority confirmations
	RejectionCount    int                                // Count of authority rejections
}

// BusinessPersonalGrantFund holds the details of the fund such as balance and distributed grants.
type BusinessPersonalGrantFund struct {
	mutex             sync.Mutex     // Mutex for thread safety
	TotalBalance      *big.Int       // Total balance available in the fund
	GrantsDistributed *big.Int       // Total amount of grants distributed
	Ledger            *ledger.Ledger // Reference to the ledger for storing transactions
}

// BusinessPersonalGrantDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Business Personal Grant Fund.
type BusinessPersonalGrantDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// BusinessPersonalGrantDisbursementManager manages the disbursement of funds for confirmed proposals in the Business Personal Grant Fund.
type BusinessPersonalGrantDisbursementManager struct {
	mutex                  sync.Mutex                                 // Mutex for thread safety
	Ledger                 *ledger.Ledger                             // Reference to the ledger
	Consensus              *common.SynnergyConsensus                // Synnergy Consensus engine
	FundBalance            float64                                    // Current balance of the Business Personal Grant Fund
	DisbursementQueue      []*BusinessPersonalGrantDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime           time.Duration                              // Maximum time a proposal can wait in the queue (30 days)
}

// GrantProposal represents the structure of the business grant application.
type BusinessPersonalGrantProposal struct {
	BusinessName        string              // Name of the business
	BusinessAddress     string              // Address of the business
	RegistrationNumber  string              // Business registration number
	Country             string              // Country of registration
	Website             string              // Business website (optional)
	BusinessActivities  string              // Description of business activities
	ApplicantName       string              // Name of the acting member applying for the funds
	WalletAddress       string              // The wallet address of the applicant
	AmountAppliedFor    float64             // Amount of grant funds being applied for
	UsageDescription    string              // Full description of how the funds will be used
	FinancialPosition   string              // Financial position or last submitted accounts (or state if it's a startup)
	SubmissionTimestamp time.Time           // Timestamp of proposal submission
	VerifiedBySyn900    bool                // Whether the proposal has been verified with syn900
	Status              string              // Proposal status (e.g., Pending, Approved, Rejected)
	Comments            []ProposalComment   // Comments made on the proposal
	LastUpdated         time.Time           // Last update timestamp for the proposal
}



// ProposalManager manages the submission and verification of grant proposals.
type BusinessPersonalGrantProposalManager struct {
	mutex           sync.Mutex                         // Mutex for thread safety
	Ledger          *ledger.Ledger                     // Reference to the ledger for storing proposals
	Proposals       map[string]*BusinessPersonalGrantProposal          // Map of proposals by business name
	Syn900Validator *common.Syn900Validator                  // Reference to syn900 validator for wallet verification
	Encryption      *common.Encryption             // Encryption service for secure proposal data
}

// ProposalComment represents a comment added to a proposal.
type ProposalComment struct {
	CommentID   string    // Unique ID for the comment
	Commenter   string    // Name or wallet address of the commenter
	Comment     string    // The content of the comment
	CreatedAt   time.Time // Timestamp of when the comment was added
}
