package loanpool

import (
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// EducationFundApprovalProcess manages the approval process for education fund proposals.
type EducationFundApprovalProcess struct {
	mutex               sync.Mutex
	Ledger              *ledger.Ledger               // Ledger to store proposal and approval status
	Nodes               []*common.AuthorityNodeTypes             // List of all authority nodes in the network
	ActiveProposals     map[string]*EducationFundActiveProposal   // Map of active proposals being reviewed
	EncryptionService   *common.Encryption       // Encryption service for secure transmission
	RequeueDuration     time.Duration                // Duration before a proposal is requeued
	MaxConfirmations    int                          // Required confirmations for proposal approval
	MaxRejections       int                          // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, and assigned nodes for a proposal.
type EducationFundActiveProposal struct {
	ProposalID        string                     // Unique proposal ID
	ProposalData      *EducationFundProposal     // The education fund proposal details
	ConfirmedNodes    map[string]bool            // Nodes that confirmed the proposal
	RejectedNodes     map[string]bool            // Nodes that rejected the proposal
	AssignedNodes     map[string]*common.AuthorityNodeTypes  // Nodes currently assigned for review
	Status            string                     // Status of the proposal (Pending, Approved, Rejected)
	LastDistribution  time.Time                  // Timestamp of last node distribution
	ProposalDeadline  time.Time                  // Deadline for the proposal before requeuing
}


// EducationFund holds the details of the fund, including balance and distributed grants.
type EducationFund struct {
	mutex             sync.Mutex     // Mutex for thread safety
	TotalBalance      *big.Int       // Total balance available in the fund
	GrantsDistributed *big.Int       // Total amount of grants distributed
	Ledger            *ledger.Ledger // Reference to the ledger for storing transactions
	Consensus         *common.SynnergyConsensus // Synnergy Consensus engine for transaction validation
	EncryptionService *common.Encryption // Encryption service for securing sensitive data
}

// EducationFundDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Education Fund.
type EducationFundDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// EducationFundDisbursementManager manages the disbursement of funds for confirmed proposals in the Education Fund.
type EducationFundDisbursementManager struct {
	mutex             sync.Mutex                             // Mutex for thread safety
	Ledger            *ledger.Ledger                         // Reference to the ledger
	Consensus         *common.SynnergyConsensus             // Synnergy Consensus engine
	FundBalance       float64                                // Current balance of the Education Fund
	DisbursementQueue []*EducationFundDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime      time.Duration                          // Maximum time a proposal can wait in the queue (30 days)
	EncryptionService *common.Encryption                 // Encryption service for secure data
}

// EducationFundProposal represents the structure of the education fund proposal.
type EducationFundProposal struct {
	ApplicantName       string            // Name of the applicant
	ApplicantContact    string            // Applicant's contact information
	WalletAddress       string            // The wallet address of the applicant
	InstitutionName     string            // Name of the educational institution
	CourseName          string            // Name of the course
	CourseLevel         string            // Level of the course (e.g., Bachelor's, Master's)
	ApplicationEvidence string            // Evidence of course application or acceptance
	PersonalStatement   string            // Personal statement of the applicant
	AmountAppliedFor    float64           // Amount of funds being applied for
	SubmissionTimestamp time.Time         // Timestamp of proposal submission
	SponsorName         string            // Name of the sponsor (if applicable)
	SponsorContactInfo  string            // Contact information of the sponsor (if applicable)
	VerifiedBySyn900    bool              // Whether the proposal has been verified with syn900
	Status              string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments            []ProposalComment // Comments made on the proposal
	LastUpdated         time.Time         // Last update timestamp for the proposal
}


// ProposalManager manages the submission and verification of education fund proposals.
type EducationFundProposalManager struct {
	mutex           sync.Mutex                        // Mutex for thread safety
	Ledger          *ledger.Ledger                    // Reference to the ledger for storing proposals
	Proposals       map[string]*EducationFundProposal // Map of proposals by applicant name
	Syn900Validator *common.Syn900Validator                 // Reference to syn900 validator for wallet verification
	Encryption      *common.Encryption            // Encryption service for secure proposal data
}

// ProposalComment represents a comment added to a proposal.
type ProposalComment struct {
	CommentID   string    // Unique ID for the comment
	Commenter   string    // Name or wallet address of the commenter
	Comment     string    // The content of the comment
	CreatedAt   time.Time // Timestamp of when the comment was added
}
