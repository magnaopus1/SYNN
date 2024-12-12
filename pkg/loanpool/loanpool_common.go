package loanpool

import (
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"time"
)

// LoanPool represents the structure of the main loan pool.
type LoanPool struct {
	mutex                 sync.Mutex                 // For thread safety
	MainFund              *big.Int                   // Main loanpool fund
	PersonalGrantFund     *big.Int                   // 25%
	EcosystemGrantFund    *big.Int                   // 25%
	EducationFund         *big.Int                   // 5%
	HealthcareSupportFund *big.Int                   // 5%
	PovertyFund           *big.Int                   // 5%
	SecuredFund           *big.Int                   // 15%
	BusinessGrantFund     *big.Int                   // 25%
	UnsecuredLoanFund     *big.Int                   // 15%
	EnvironmentalFund     *big.Int                   // 5%
	Ledger                *ledger.Ledger             // Reference to the ledger for transaction logging
	Consensus             *common.SynnergyConsensus // Synnergy Consensus engine for validating fund transfers
	Encryption            *common.Encryption     // Encryption service for secure data handling
}

// LoanPoolManager provides functions to view the balances of the loan pool and its associated sub-funds.
type LoanPoolManager struct {
	mutex    sync.Mutex     // For thread-safe operations
	LoanPool *LoanPool      // Reference to the LoanPool structure
	Ledger   *ledger.Ledger // Reference to the ledger for transaction logging
}



// AuthorityNode represents an authority node in the network responsible for approving/rejecting proposals.
type AuthorityNode struct {
	NodeID     string
	NodeStatus string // E.g., Online, Busy, Offline
}




// SecuredLoanApprovalProcess manages the approval process for secured loan proposals.
type SecuredLoanApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *ledger.Ledger              // Ledger to store proposal and approval status
	Nodes             []*AuthorityNode            // List of all authority nodes in the network
	ActiveProposals   map[string]*ActiveProposal  // Map of active proposals being reviewed
	EncryptionService *common.Encryption      // Encryption service for secure transmission
	NetworkManager    *network.NetworkManager     // Network manager to handle proposal transmission
	RequeueDuration   time.Duration               // Duration before a proposal is requeued (48 hours)
	MaxConfirmations  int                         // Required confirmations for proposal approval
	MaxRejections     int                         // Required rejections for proposal rejection
}

// ActiveProposal keeps track of confirmations, rejections, interest rates, and assigned nodes for a proposal.
type SecuredLoanActiveProposal struct {
	ProposalID        string                     // Unique proposal ID
	ProposalData      *SecuredLoanProposal       // The loan proposal details
	ConfirmedNodes    map[string]bool            // Nodes that confirmed the proposal
	RejectedNodes     map[string]bool            // Nodes that rejected the proposal
	AssignedNodes     map[string]*common.AuthorityNodeTypes  // Nodes currently assigned for review
	Status            string                     // Status of the proposal (Pending, Approved, Rejected)
	InterestRates     []float64                  // List of interest rates submitted by authority nodes
	LastDistribution  time.Time                  // Timestamp of last node distribution
	ProposalDeadline  time.Time                  // Deadline for the proposal before requeuing
	AverageInterest   float64                    // Running average of the interest rates
	AllDocsOpened     bool                       // Whether all nodes have opened the attached documents
}

// CollateralSubmission represents the structure for collateral proof submission.
type CollateralSubmission struct {
	LoanID           string    // Unique loan ID linked to the collateral
	ProposerWallet   string    // Wallet address of the proposer (borrower)
	CollateralType   string    // Type of collateral (e.g., property, car, assets)
	CollateralValue  float64   // Value of the collateral being offered
	CollateralProof  []byte    // Digital document providing proof of collateral (e.g., title deed)
	IOULegalDocument []byte    // Digital IOU or legal agreement document
	SubmissionTime   time.Time // Time of collateral submission
	ApprovalStatus   string    // Status of the collateral (Pending, Approved, Rejected)
	ApprovedBy       string    // Authority node that approved/rejected the collateral
	ApprovalTime     time.Time // Time of approval or rejection
}

// CollateralManager manages the submission and approval process for collateral in a secured loan.
type CollateralManager struct {
	mutex             sync.Mutex                       // Mutex for thread safety
	Ledger            *ledger.Ledger                   // Ledger to record collateral submissions and approvals
	Consensus         *common.SynnergyConsensus      // Synnergy Consensus engine for validating collateral
	EncryptionService *common.Encryption           // Encryption service for securing sensitive collateral data
	Submissions       map[string]*CollateralSubmission // Map of collateral submissions by loan ID
	ApprovalQueue     []*CollateralSubmission          // Queue of submissions pending approval
}

// LoanTerms represents the structure of customized loan terms.
type LoanTerms struct {
	RepaymentLength      int     // Number of months to repay the loan
	AmountBorrowed       float64 // Total loan amount
	InterestRate         float64 // Interest rate applied (unless Islamic terms)
	IslamicFinance       bool    // If true, switches to Islamic terms (no interest, fee applied)
	FeeOnTop             float64 // Flat fee applied if Islamic finance is selected
	TotalRepaymentAmount float64 // Total amount to be repaid (calculated)
}

// SecuredLoanTermManager manages customization of secured loan terms.
type SecuredLoanTermManager struct {
	mutex             sync.Mutex                 // Mutex for thread safety
	Ledger            *ledger.Ledger             // Reference to the ledger for storing loan term records
	Consensus         *common.SynnergyConsensus // Synnergy Consensus engine for validating terms
	LoanTermRecords   map[string]*LoanTerms      // Records of loan terms by loan ID
	EncryptionService *common.Encryption     // Encryption service for securing sensitive data
}


// SecuredLoanDisbursementQueueEntry represents a loan waiting for disbursement in the secured loan pool.
type SecuredLoanDisbursementQueueEntry struct {
	ProposalID        string    // The loan proposal ID
	ProposerWallet    string    // Wallet address of the borrower
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
	AverageInterest   float64   // Average interest rate to be applied to the loan
}

// SecuredLoanDisbursementManager manages the disbursement of approved loans in the secured loan pool.
type SecuredLoanDisbursementManager struct {
	mutex               sync.Mutex                           // Mutex for thread safety
	Ledger              *ledger.Ledger                       // Reference to the ledger for recording disbursements
	Consensus           *common.SynnergyConsensus           // Synnergy Consensus engine
	FundBalance         float64                              // Current balance of the loan pool
	DisbursementQueue   []*SecuredLoanDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime        time.Duration                        // Maximum time a proposal can wait in the queue
	EncryptionService   *common.Encryption               // Encryption service for securing sensitive data
	IssuerFeePercentage float64                              // Issuer fee (0.5%)
}

// SecuredLoanProposal represents the full loan application process.
type SecuredLoanProposal struct {
	LoanID               string    // Unique loan ID
	ApplicantName        string    // Name of the applicant
	ApplicantID          string    // Unique applicant ID (validated by syn900)
	WalletAddress        string    // Wallet address of the applicant
	SubmissionTimestamp  time.Time // Time when the proposal was submitted
	ProposalStatus       string    // Status of the proposal (e.g., Pending, Approved, Rejected)
	LastUpdated          time.Time // Timestamp of last update
	ApprovalStage        string    // Current stage of the approval process (e.g., Application, CreditCheck, Affordability, Collateral, Terms)
	CreditScore          float64   // Applicant's credit score (from decentralized credit check)
	AffordabilityStatus  string    // Result of affordability check (e.g., Approved, Rejected)
	CollateralStatus     string    // Result of collateral submission
	TermsCustomization   bool      // Whether terms customization has been completed
}

// ProposalManager handles the overall proposal process for secured loans.
type SecuredLoanProposalManager struct {
	mutex             sync.Mutex                          // Mutex for thread safety
	Ledger            *ledger.Ledger                      // Ledger for storing proposal data
	Consensus         *common.SynnergyConsensus        // Consensus engine for validation
	Syn900Validator   *common.Syn900Validator                   // Syn900 validator for ID validation
	EncryptionService *common.Encryption              // Encryption service for secure proposal data
	CreditChecker     *CreditCheckManager                 // Decentralized Credit Check manager
	AffordabilityMgr  *AffordabilityManager               // Affordability Check manager
	CollateralMgr     *CollateralManager                  // Collateral submission manager
	TermsManager      *TermsCustomizationManager          // Customization of loan terms
	Proposals         map[string]*SecuredLoanProposal     // Proposals mapped by LoanID
}

// LoanRepaymentDetails stores repayment information for a loan.
type SecuredLoanRepaymentDetails struct {
	LoanID            string            // Unique loan ID
	ProposerWallet    string            // Wallet address of the borrower
	TotalAmount       float64           // Total amount to be repaid
	RemainingAmount   float64           // Remaining amount to be repaid
	InterestRate      float64           // Interest rate applied to the loan
	RepaymentDates    []time.Time       // Scheduled payment dates
	NextPaymentDue    time.Time         // Next payment due date
	Status            string            // Loan status (Active, Defaulted, Satisfied)
	DefaultedAt       *time.Time        // If loan defaulted, record the default date
	CollateralContact string            // Contact email for collateral request
	AuthorityWallets  []string          // Wallet addresses of authority nodes
}

// SecuredLoanRepaymentManager manages the repayment and settlement process for secured loans.
type SecuredLoanRepaymentManager struct {
	mutex               sync.Mutex                          // Mutex for thread safety
	Ledger              *ledger.Ledger                      // Reference to the ledger
	EncryptionService   *common.Encryption              // Encryption service for secure data
	Syn900Registry      *syn900Registry                    // Reference to Syn900 for record keeping
	LoanRepayments      map[string]*SecuredLoanRepaymentDetails    // Map of loan repayments by loan ID
	DefaultThreshold    time.Duration                       // Time duration for default (e.g., 6 months)
}


// SecuredLoanManagement handles management tasks such as authority node updates, borrower detail changes, and term change requests.
type SecuredLoanManagement struct {
	mutex             sync.Mutex
	Ledger            *ledger.Ledger                          // Ledger reference for recording updates
	EncryptionService *common.Encryption                  // Encryption service for data security
	ConsensusEngine   *common.SynnergyConsensus             // Consensus engine for approval processes
	NetworkManager    *network.NetworkManager                 // Network manager for sending requests
	AuthorityNodes    map[string]*common.AuthorityNodeTypes        // Map of authority nodes that can manage loans
	LoanBorrowerInfo  map[string]*BorrowerDetails             // Stores borrower details by loan ID
	TermChangeRequests map[string]*BorrowerTermChangeRequest  // Stores term change requests by loan ID
}

// BorrowerTermChangeRequest represents a request for changing loan terms.
type BorrowerTermChangeRequest struct {
	LoanID              string
	RequestedTerms      string            // New terms requested by the borrower
	ApprovalStatus      string            // Approval status (Pending, Accepted, Rejected)
	ConfirmedNodes      map[string]bool   // Nodes that confirmed the request
	RejectedNodes       map[string]bool   // Nodes that rejected the request
	AssignedNodes       map[string]*common.AuthorityNodeTypes // Nodes currently assigned for review
	LastDistribution    time.Time         // Last node distribution time
	RequeueDeadline     time.Time         // Requeue if not processed within this time
}


// SecuredLoanPool manages the details of the fund, including balance, loan disbursements, repayments, and defaults.
type SecuredLoanPool struct {
	mutex               sync.Mutex     // Mutex for thread safety
	TotalBalance        *big.Int       // Total balance available in the loan pool
	LoansDistributed    *big.Int       // Total amount of loans distributed
	LoansRepaid         *big.Int       // Total amount of loans repaid
	LoansDefaulted      *big.Int       // Total amount of loans that have defaulted
	Ledger              *ledger.Ledger // Reference to the ledger for storing loan transactions
	Consensus           *common.SynnergyConsensus // Synnergy Consensus engine for transaction validation
	EncryptionService   *common.Encryption // Encryption service for securing sensitive loan data
	LoanRecords         map[string]*LoanRecord  // Map of loan records by applicant wallet
}



// SmallBusinessGrantApprovalProcess handles the two-stage approval process of a small business grant proposal.
type SmallBusinessGrantApprovalProcess struct {
	mutex             sync.Mutex                                 // Mutex for thread safety
	Ledger            *ledger.Ledger                             // Reference to the ledger
	Consensus         *common.SynnergyConsensus                 // Synnergy Consensus engine
	Proposals         map[string]*SmallBusinessGrantProposalApproval // Map to hold grant proposals by proposal ID
	AuthorityNodes    []common.AuthorityNodeTypes                          // List of valid authority node types (bank, government, central bank, etc.)
	PublicVotePeriod  time.Duration                              // Time allowed for public voting
	AuthorityVoteTime time.Duration                              // Time window for authority nodes to vote
}

// SmallBusinessGrantProposalApproval represents a grant proposal along with its voting data.
type SmallBusinessGrantProposalApproval struct {
	Proposal          *SmallBusinessGrantProposal                // Reference to the grant proposal
	PublicVotes       map[string]bool                            // Map of public votes (address -> vote)
	Stage             ApprovalStage                             // Current approval stage
	AuthorityVotes    map[string]bool                            // Authority node votes
	VoteStartTime     time.Time                                 // Time when voting starts
	ConfirmationCount int                                       // Count of authority confirmations
	RejectionCount    int                                       // Count of authority rejections
}

// SmallBusinessGrantDisbursementQueueEntry represents a proposal waiting for funds to be disbursed for the Small Business Grant Fund.
type SmallBusinessGrantDisbursementQueueEntry struct {
	ProposalID        string    // The proposal ID
	ProposerWallet    string    // Wallet address of the proposer
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
}

// SmallBusinessGrantDisbursementManager manages the disbursement of funds for confirmed proposals in the Small Business Grant Fund.
type SmallBusinessGrantDisbursementManager struct {
	mutex                  sync.Mutex                                 // Mutex for thread safety
	Ledger                 *ledger.Ledger                             // Reference to the ledger
	Consensus              *common.SynnergyConsensus               // Synnergy Consensus engine
	FundBalance            float64                                    // Current balance of the Small Business Grant Fund
	DisbursementQueue      []*SmallBusinessGrantDisbursementQueueEntry // Queue for proposals waiting for disbursement
	QueueMaxTime           time.Duration                              // Maximum time a proposal can wait in the queue (30 days)
}

// SmallBusinessGrantProposal represents the structure of the small business grant application.
type SmallBusinessGrantProposal struct {
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
	Startup             bool                // Is the business a startup?
	EmployeeCount       int                 // Number of employees (required if not a startup)
}

// SmallBusinessProposalManager manages the submission and verification of small business grant proposals.
type SmallBusinessGrantProposalManager struct {
	mutex           sync.Mutex                         // Mutex for thread safety
	Ledger          *ledger.Ledger                     // Reference to the ledger for storing proposals
	Proposals       map[string]*SmallBusinessGrantProposal // Map of proposals by business name
	Syn900Validator *common.Syn900Validator                 // Reference to syn900 validator for wallet verification
	Encryption      *common.Encryption             // Encryption service for secure proposal data
}

// SmallBusinessGrantFund holds the details of the fund such as balance and distributed grants.
type SmallBusinessGrantFund struct {
	mutex             sync.Mutex     // Mutex for thread safety
	TotalBalance      *big.Int       // Total balance available in the fund
	GrantsDistributed *big.Int       // Total amount of grants distributed
	Ledger            *ledger.Ledger // Reference to the ledger for storing transactions
}

// AffordabilityCheck represents the structure for an affordability assessment.
type AffordabilityCheck struct {
	LoanID             string    // Unique loan ID
	ApplicantWallet    string    // Wallet address of the applicant
	Income             float64   // Applicant's monthly income
	Expenses           float64   // Applicant's monthly expenses
	Dependents         int       // Number of dependents the applicant has
	DependentCosts     float64   // Monthly cost for dependents
	WorkingStatus      string    // Employment status (e.g., employed, self-employed, unemployed)
	OtherDebts         float64   // Total amount of other debts the applicant owes
	WorkProof          []byte    // Proof of employment (e.g., employment contract or income statement)
	SubmissionTime     time.Time // Timestamp of affordability check submission
	ApprovalStatus     string    // Status of the affordability check (Pending, Approved, Rejected)
	ApprovedBy         string    // Authority node that approved or rejected the check
	ApprovalTime       time.Time // Time of approval or rejection
}

// AffordabilityManager handles the submission and approval process for affordability checks.
type AffordabilityManager struct {
	mutex             sync.Mutex                          // Mutex for thread safety
	Ledger            *ledger.Ledger                      // Ledger for recording affordability checks
	Consensus         *common.SynnergyConsensus         // Synnergy Consensus engine for validation
	EncryptionService *common.Encryption              // Encryption service for securing sensitive financial data
	Submissions       map[string]*AffordabilityCheck      // Map of affordability submissions by loan ID
	ApprovalQueue     []*AffordabilityCheck               // Queue of submissions pending approval
}



// UnsecuredLoanApprovalProcess manages the approval process for unsecured loan proposals.
type UnsecuredLoanApprovalProcess struct {
	mutex             sync.Mutex
	Ledger            *ledger.Ledger              // Ledger to store proposal and approval status
	Nodes             []*common.AuthorityNodeTypes            // List of all authority nodes in the network
	ActiveProposals   map[string]*ActiveProposal  // Map of active proposals being reviewed
	EncryptionService *common.Encryption      // Encryption service for secure transmission
	NetworkManager    *network.NetworkManager     // Network manager to handle proposal transmission
	RequeueDuration   time.Duration               // Duration before a proposal is requeued (48 hours)
	MaxConfirmations  int                         // Required confirmations for proposal approval
	MaxRejections     int                         // Required rejections for proposal rejection
}

// ActiveProposal tracks confirmations, rejections, interest rates, and assigned nodes for a proposal.
type ActiveProposal struct {
	ProposalID        string                     // Unique proposal ID
	ProposalData      *UnsecuredLoanProposal     // The loan proposal details
	ConfirmedNodes    map[string]bool            // Nodes that confirmed the proposal
	RejectedNodes     map[string]bool            // Nodes that rejected the proposal
	AssignedNodes     map[string]*common.AuthorityNodeTypes  // Nodes currently assigned for review
	Status            string                     // Status of the proposal (Pending, Approved, Rejected)
	InterestRates     []float64                  // List of interest rates submitted by authority nodes
	LastDistribution  time.Time                  // Timestamp of last node distribution
	ProposalDeadline  time.Time                  // Deadline for the proposal before requeuing
	AverageInterest   float64                    // Running average of the interest rates
	AllDocsOpened     bool                       // Whether all nodes have opened the attached documents
}

// UnsecuredLoanTermManager manages customization of unsecured loan terms.
type UnsecuredLoanTermManager struct {
	mutex              sync.Mutex                    // Mutex for thread safety
	Ledger             *ledger.Ledger                // Reference to the ledger for storing loan term records
	Consensus          *common.SynnergyConsensus   // Synnergy Consensus engine for validating terms
	LoanTermRecords    map[string]*LoanTerms         // Records of loan terms by loan ID
	EncryptionService  *common.Encryption        // Encryption service for securing sensitive data
}

// DecentralizedCreditCheck tracks spending for wallets and stores credit score documents.
type DecentralizedCreditCheck struct {
	mutex                sync.Mutex                    // Mutex for thread safety
	Ledger               *ledger.Ledger                // Ledger for storing credit check and transaction data
	Consensus            *common.SynnergyConsensus   // Consensus engine for validation
	WalletSpendingRecords map[string]*SpendingRecord   // Stores spending records by wallet address
	CreditScoreDocuments  map[string][]byte            // Stores encrypted credit score documents by wallet address
	EncryptionService     *common.Encryption       // Encryption service for securing sensitive data
}

// SpendingRecord represents a record of wallet spending and associated transactions.
type SpendingRecord struct {
	WalletAddress   string                      // Address of the wallet being tracked
	TotalSpent      float64                     // Total amount spent from the wallet
	Transactions    []*common.Transaction       // List of transactions made from the wallet
	LastUpdated     time.Time                   // Timestamp of the last update
}

// UnsecuredLoanDisbursementQueueEntry represents a loan waiting for disbursement in the unsecured loan pool.
type UnsecuredLoanDisbursementQueueEntry struct {
	ProposalID        string    // The loan proposal ID
	ProposerWallet    string    // Wallet address of the borrower
	RequestedAmount   float64   // Amount requested for disbursement
	DisbursementStart time.Time // The time the proposal was added to the queue
	AverageInterest   float64   // Average interest rate to be applied to the loan
}

// UnsecuredLoanDisbursementManager manages the disbursement of approved loans in the unsecured loan pool.
type UnsecuredLoanDisbursementManager struct {
	mutex               sync.Mutex                               // Mutex for thread safety
	Ledger              *ledger.Ledger                           // Reference to the ledger for recording disbursements
	Consensus           *common.SynnergyConsensus               // Synnergy Consensus engine
	FundBalance         float64                                  // Current balance of the loan pool
	DisbursementQueue   []*UnsecuredLoanDisbursementQueueEntry   // Queue for proposals waiting for disbursement
	QueueMaxTime        time.Duration                            // Maximum time a proposal can wait in the queue
	EncryptionService   *common.Encryption                   // Encryption service for securing sensitive data
	IssuerFeePercentage float64                                  // Issuer fee (0.5%)
}

// UnsecuredLoanProposal represents the full loan application process.
type UnsecuredLoanProposal struct {
	LoanID               string    // Unique loan ID
	ApplicantName        string    // Name of the applicant
	ApplicantID          string    // Unique applicant ID (validated by syn900)
	WalletAddress        string    // Wallet address of the applicant
	SubmissionTimestamp  time.Time // Time when the proposal was submitted
	ProposalStatus       string    // Status of the proposal (e.g., Pending, Approved, Rejected)
	LastUpdated          time.Time // Timestamp of last update
	ApprovalStage        string    // Current stage of the approval process (e.g., Application, CreditCheck, Affordability, Terms)
	CreditScore          float64   // Applicant's credit score (from decentralized credit check)
	AffordabilityStatus  string    // Result of affordability check (e.g., Approved, Rejected)
	TermsCustomization   bool      // Whether terms customization has been completed
}

// ProposalManager handles the overall proposal process for unsecured loans.
type UnsecuredLoanProposalManager struct {
	mutex             sync.Mutex                           // Mutex for thread safety
	Ledger            *ledger.Ledger                       // Ledger for storing proposal data
	Consensus         *common.SynnergyConsensus          // Consensus engine for validation
	Syn900Validator   *common.Syn900Validator                    // Syn900 validator for ID validation
	EncryptionService *common.Encryption               // Encryption service for secure proposal data
	CreditChecker     *CreditCheckManager                  // Decentralized Credit Check manager
	AffordabilityMgr  *AffordabilityManager                // Affordability Check manager
	TermsManager      *TermsCustomizationManager           // Customization of loan terms
	Proposals         map[string]*UnsecuredLoanProposal    // Proposals mapped by LoanID
}

// LoanRepaymentDetails stores repayment information for a loan.
type UnsecuredLoanRepaymentDetails struct {
	LoanID            string            // Unique loan ID
	ProposerWallet    string            // Wallet address of the borrower
	TotalAmount       float64           // Total amount to be repaid
	RemainingAmount   float64           // Remaining amount to be repaid
	InterestRate      float64           // Interest rate applied to the loan
	RepaymentDates    []time.Time       // Scheduled payment dates
	NextPaymentDue    time.Time         // Next payment due date
	Status            string            // Loan status (Active, Defaulted, Satisfied)
	DefaultedAt       *time.Time        // If loan defaulted, record the default date
	DefaultContact    string            // Contact email for default notification
	AuthorityWallets  []string          // Wallet addresses of authority nodes
}

// SecuredLoanRepaymentManager manages the repayment and settlement process for secured loans.
type UnsecuredLoanRepaymentManager struct {
	mutex               sync.Mutex                          // Mutex for thread safety
	Ledger              *ledger.Ledger                      // Reference to the ledger
	EncryptionService   *common.Encryption              // Encryption service for secure data
	Syn900Registry      *syn900Registry                    // Reference to Syn900 for record keeping
	LoanRepayments      map[string]*UnsecuredLoanRepaymentDetails    // Map of loan repayments by loan ID
	DefaultThreshold    time.Duration                       // Time duration for default (e.g., 6 months)
}

// BorrowerDetails represents the borrower information in a loan.
type BorrowerDetails struct {
	LoanID          string
	BorrowerName    string
	BorrowerEmail   string
	BorrowerContact string
	WalletAddress   string // Borrower's wallet address
}

// UnsecuredLoanManagement handles management tasks such as authority node updates, borrower detail changes, and term change requests.
type UnsecuredLoanManagement struct {
	mutex             sync.Mutex
	Ledger            *ledger.Ledger                          // Ledger reference for recording updates
	EncryptionService *common.Encryption                  // Encryption service for data security
	ConsensusEngine   *common.SynnergyConsensus            // Consensus engine for approval processes
	NetworkManager    *network.NetworkManager                 // Network manager for sending requests
	AuthorityNodes    map[string]*common.AuthorityNodeTypes        // Map of authority nodes that can manage loans
	LoanBorrowerInfo  map[string]*BorrowerDetails             // Stores borrower details by loan ID
	TermChangeRequests map[string]*BorrowerTermChangeRequest  // Stores term change requests by loan ID
}

// Syn900Registry represents the registry for associating a loan with an SYN900 token.
type Syn900Registry struct {
	mutex        sync.Mutex         // For thread safety
	LoanID       string             // Unique identifier for the loan
	SYN900TokenID string            // The SYN900 token ID associated with the borrower
	Ledger       *ledger.Ledger     // Ledger reference to track loan transactions
}


// UnsecuredLoanPool manages the details of the unsecured loan pool, including balance, loan disbursements, repayments, and defaults.
type UnsecuredLoanPool struct {
	mutex             sync.Mutex                 // Mutex for thread safety
	TotalBalance      *big.Int                   // Total balance available in the loan pool
	LoansDistributed  *big.Int                   // Total amount of loans distributed
	LoansRepaid       *big.Int                   // Total amount of loans repaid
	LoansDefaulted    *big.Int                   // Total amount of loans that have defaulted
	Ledger            *ledger.Ledger             // Reference to the ledger for storing loan transactions
	Consensus         *common.SynnergyConsensus // Synnergy Consensus engine for transaction validation
	EncryptionService *common.Encryption     // Encryption service for securing sensitive loan data
	LoanRecords       map[string]*LoanRecord     // Map of loan records by applicant wallet
}

// LoanRecord tracks details for each unsecured loan.
type LoanRecord struct {
	ProposalID       string    // Unique ID for the loan proposal
	ApplicantWallet  string    // Wallet address of the loan applicant
	LoanAmount       *big.Int  // Amount of the loan distributed
	LoanRepaid       *big.Int  // Amount repaid so far
	LoanStatus       string    // Status of the loan: Active, Repaid, Defaulted
	RepaymentDueDate time.Time // Due date for full repayment
}

// ApprovalStage represents a stage in a multi-step approval process.
type ApprovalStage struct {
	StageID        string    // Unique identifier for the approval stage
	Description    string    // Description of the approval stage (e.g., "KYC Verification", "AML Check")
	IsApproved     bool      // Whether this stage has been approved
	ApproverID     string    // The identifier of the entity or individual who approved this stage
	ApprovalTime   time.Time // The timestamp when the stage was approved
	RejectionReason string   // Reason for rejection, if applicable
	IsRejected     bool      // Whether this stage has been rejected
}
