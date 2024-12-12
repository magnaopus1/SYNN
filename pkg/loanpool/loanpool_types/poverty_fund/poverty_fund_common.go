package loanpool

import (
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// PovertyFundApprovalProcess manages the approval process for poverty fund proposals.
type PovertyFundApprovalProcess struct {
	mutex               sync.Mutex
	Ledger              *ledger.Ledger                     // Ledger to store proposal and approval status
	Nodes               []*common.AuthorityNodeTypes                   // List of all authority nodes in the network
	ActiveProposals     map[string]*PovertyFundActiveProposal         // Map of active proposals being reviewed
	EncryptionService   *common.Encryption             // Encryption service for secure transmission
	RequeueDuration     time.Duration                      // Duration before a proposal is requeued (7 days)
	MaxConfirmations    int                                // Required confirmations for proposal approval
	MaxRejections       int                                // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, and assigned nodes for a proposal.
type PovertyFundActiveProposal struct {
	ProposalID        string                        // Unique proposal ID
	ProposalData      *PovertyFundProposal          // The poverty fund proposal details
	ConfirmedNodes    map[string]bool               // Nodes that confirmed the proposal
	RejectedNodes     map[string]bool               // Nodes that rejected the proposal
	AssignedNodes     map[string]*common.AuthorityNodeTypes     // Nodes currently assigned for review
	Status            string                        // Status of the proposal (Pending, Approved, Rejected)
	LastDistribution  time.Time                     // Timestamp of last node distribution
	ProposalDeadline  time.Time                     // Deadline for the proposal before requeuing
}

// PovertyFundDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Poverty Fund.
type PovertyFundDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// PovertyFundDisbursementManager manages the disbursement of funds for confirmed proposals in the Poverty Fund.
type PovertyFundDisbursementManager struct {
	mutex                  sync.Mutex                                 // Mutex for thread safety
	Ledger                 *ledger.Ledger                             // Reference to the ledger
	Consensus              *common.SynnergyConsensus                 // Synnergy Consensus engine
	FundBalance            float64                                    // Current balance of the Poverty Fund
	DisbursementQueue      []*PovertyFundDisbursementQueueEntry       // Queue for proposals waiting for disbursement
	QueueMaxTime           time.Duration                              // Maximum time a proposal can wait in the queue (48 hours)
	EncryptionService      *common.Encryption                     // Encryption service for secure data
}

// PovertyFund manages the details of the fund, including balance and distributed grants.
type PovertyFund struct {
	mutex             sync.Mutex                    // Mutex for thread safety
	TotalBalance      *big.Int                      // Total balance available in the fund
	GrantsDistributed *big.Int                      // Total amount of grants distributed
	Ledger            *ledger.Ledger                // Reference to the ledger for storing transactions
	Consensus         *common.SynnergyConsensus // Synnergy Consensus engine for transaction validation
	EncryptionService *common.Encryption        // Encryption service for securing sensitive data
}

// PovertyFundProposal represents the structure of the poverty fund proposal.
type PovertyFundProposal struct {
	ApplicantName        string              // Name of the applicant
	ApplicantContact     string              // Applicant's contact information
	IncomeDetails        string              // Income details of the applicant
	BankBalanceDetails   string              // Current bank balance details of the applicant
	IncomeEvidence       []byte              // Attachment: Evidence of income (encrypted)
	BankBalanceEvidence  []byte              // Attachment: Evidence of bank balance (encrypted)
	StatementOfReason    string              // Statement of reason for the request
	BenefitStatus        string              // Current benefit status of the applicant
	WalletAddress        string              // The wallet address of the applicant
	AmountAppliedFor     float64             // Amount of funds being applied for
	SubmissionTimestamp  time.Time           // Timestamp of proposal submission
	VerifiedBySyn900     bool                // Whether the proposal has been verified with syn900
	Status               string              // Proposal status (e.g., Pending, Approved, Rejected)
	Comments             []ProposalComment   // Comments made on the proposal
	LastUpdated          time.Time           // Last update timestamp for the proposal
}

// ProposalComment represents a comment added to a proposal.
type ProposalComment struct {
	CommentID   string    // Unique ID for the comment
	Commenter   string    // Name or wallet address of the commenter
	Comment     string    // The content of the comment
	CreatedAt   time.Time // Timestamp of when the comment was added
}

// ProposalManager manages the submission and verification of poverty fund proposals.
type PovertyFundProposalManager struct {
	mutex           sync.Mutex                      // Mutex for thread safety
	Ledger          *ledger.Ledger                  // Reference to the ledger for storing proposals
	Proposals       map[string]*PovertyFundProposal // Map of proposals by applicant name
	Syn900Validator *common.Syn900Validator            // Reference to syn900 validator for wallet verification
	Encryption      *common.Encryption          // Encryption service for secure proposal data
}