package loanpool

import (
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// HealthcareSupportFundApprovalProcess manages the approval process for healthcare fund proposals.
type HealthcareSupportFundApprovalProcess struct {
	mutex               sync.Mutex
	Ledger              *ledger.Ledger                     // Ledger to store proposal and approval status
	Nodes               []*common.AuthorityNodeTypes                   // List of all authority nodes in the network
	ActiveProposals     map[string]*HealthcareSupportFundActiveProposal         // Map of active proposals being reviewed
	EncryptionService   *common.Encryption             // Encryption service for secure transmission
	RequeueDuration     time.Duration                      // Duration before a proposal is requeued (7 days)
	MaxConfirmations    int                                // Required confirmations for proposal approval
	MaxRejections       int                                // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, and assigned nodes for a proposal.
type HealthcareSupportFundActiveProposal struct {
	ProposalID        string                          // Unique proposal ID
	ProposalData      *HealthcareSupportFundProposal // The healthcare support fund proposal details
	ConfirmedNodes    map[string]bool                 // Nodes that confirmed the proposal
	RejectedNodes     map[string]bool                 // Nodes that rejected the proposal
	AssignedNodes     map[string]*common.AuthorityNodeTypes       // Nodes currently assigned for review
	Status            string                          // Status of the proposal (Pending, Approved, Rejected)
	LastDistribution  time.Time                       // Timestamp of last node distribution
	ProposalDeadline  time.Time                       // Deadline for the proposal before requeuing
}

// HealthcareSupportFundDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Healthcare Support Fund.
type HealthcareSupportFundDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// HealthcareSupportFundDisbursementManager manages the disbursement of funds for confirmed proposals in the Healthcare Support Fund.
type HealthcareSupportFundDisbursementManager struct {
	mutex                  sync.Mutex                                     // Mutex for thread safety
	Ledger                 *ledger.Ledger                                 // Reference to the ledger
	Consensus              *common.SynnergyConsensus                    // Synnergy Consensus engine
	FundBalance            float64                                        // Current balance of the Healthcare Support Fund
	DisbursementQueue      []*HealthcareSupportFundDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime           time.Duration                                  // Maximum time a proposal can wait in the queue (7 days)
	EncryptionService      *common.Encryption                         // Encryption service for secure data
}

// HealthcareSupportFund manages the details of the fund, including balance and distributed healthcare grants.
type HealthcareSupportFund struct {
	mutex             sync.Mutex     // Mutex for thread safety
	TotalBalance      *big.Int       // Total balance available in the fund
	GrantsDistributed *big.Int       // Total amount of healthcare grants distributed
	Ledger            *ledger.Ledger // Reference to the ledger for storing transactions
	Consensus         *common.SynnergyConsensus // Synnergy Consensus engine for transaction validation
	EncryptionService *common.Encryption // Encryption service for securing sensitive data
}

// HealthcareSupportFundProposal represents the structure of the healthcare support fund proposal.
type HealthcareSupportFundProposal struct {
	ApplicantName              string            // Name of the applicant (person requiring medical treatment)
	ApplicantContact           string            // Applicant's contact information
	MedicalProfessionalName    string            // Name of the supporting medical professional
	MedicalProfessionalContact string            // Contact information for the medical professional
	WalletAddress              string            // The wallet address of the applicant
	HospitalName               string            // Name of the hospital, medical practice, or provider
	MedicalProcedure           string            // Details of the medical procedure required
	CostBreakdownEvidence      string            // Evidence of the cost breakdown of the medical treatment
	HospitalAddress            string            // Full address of the hospital or medical provider
	HospitalContactInfo        string            // Contact information for the hospital or provider
	AmountAppliedFor           float64           // Amount of funds being applied for
	SubmissionTimestamp        time.Time         // Timestamp of proposal submission
	VerifiedBySyn900           bool              // Whether the proposal has been verified with syn900
	Status                     string            // Proposal status (e.g., Pending, Approved, Rejected)
	Comments                   []ProposalComment // Comments made on the proposal
	LastUpdated                time.Time         // Last update timestamp for the proposal
}

// ProposalManager manages the submission and verification of healthcare fund proposals.
type HealthcareSupportFundProposalManager struct {
	mutex           sync.Mutex                           // Mutex for thread safety
	Ledger          *ledger.Ledger                       // Reference to the ledger for storing proposals
	Proposals       map[string]*HealthcareSupportFundProposal // Map of proposals by applicant name
	Syn900Validator *common.Syn900Validator                   // Reference to syn900 validator for wallet verification
	Encryption      *common.Encryption               // Encryption service for secure proposal data
}
