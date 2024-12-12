package account_and_balance_operations

import (
	"math/big"
	"synnergy_network/pkg/ledger"
	"time"
)

// Supplementary Structures
type BalanceSnapshot struct {
	AccountID string
	Balance   float64
	Timestamp time.Time
}

// BalanceLock represents a lock on a specified amount within an account until a specific time.
type BalanceLock struct {
	ID        string    // Unique identifier for the lock
	AccountID string    // The account to which this lock applies
	Amount    float64   // Amount locked
	UnlockAt  time.Time // Time when the lock can be released
}


// LiquidityPool represents a liquidity pool in DeFi operations.
type LiquidityPool struct {
	PoolID             string    // Unique ID for the liquidity pool
	TotalLiquidity     *big.Float // Total liquidity in the pool
	AvailableLiquidity *big.Float // Available liquidity for operations
	RewardRate         float64   // Reward rate for liquidity providers
	CreatedAt          time.Time // Creation timestamp
	Status             string    // Status of the pool ("Active", "Paused", etc.)
}


// Account represents a user's balance and other account-related information.
type Account struct {
	Balance               float64
	HeldBalance         float64            // Add this field to hold a locked amount
	CreatedAt             time.Time
	Address               string
	Nonce                 uint64
	Verified              bool
	EncryptedKey          string
	ConnectionEvents      []ConnectionEvent
	ContractExecutionLogs []ContractExecutionLog
	MintRecords           []MintRecord
	BurnRecords           []BurnRecord
	CurrencyExchanges     []CurrencyExchange
	CustomName            string
	PrivateKey            string
	PublicKey             string
	Stake                 float64
	IsAdmin               bool
	Permissions           []string
	LockedBalances        []BalanceLock // Holds locked balance entries
	ReservedBalance       float64       // Holds reserved balance
	FreezeUntil           time.Time     // Indicates until when the account balance is frozen
	IsSuspicious          bool          // Flag for suspicious accounts
	BalanceStatus         BalanceStatus // Holds the balance status of the account
	IsFrozen         bool    // Indicates if the account is frozen
    EncryptedBalance string  // Encrypted representation of the balance (if applicable)
	ExternalAccountID string // Link to an external account
    Authorizations  []Authorization // List of authorizations for this account
    VerificationHold  float64 // Holds funds for verification purposes
	LockedBalance     float64          // Balance that is locked
	Approvals         []Approval       // List of approved transactions
	LastTransactionID string           // ID of the last transaction affecting this account
	TotalDeposited    float64          // Total deposited into the account
	TotalWithdrawn    float64          // Total withdrawn from the account
	LastUpdated       time.Time        // Last time the account was updated
	RequiresReview  bool          // Indicates if the account is flagged for review
	Allocations     []Allocation  // List of allocations for specific purposes
}




// ConnectionEvent represents an event related to wallet connections.
type ConnectionEvent struct {
	EventID        string    // Unique identifier for the event
	ConnectionID   string
	WalletID       string    // Wallet associated with the event
	EventType      string    // Type of event (e.g., "connection", "disconnection")
	EventTime      time.Time // Timestamp of the event
	Details        string    // Additional details related to the event
}

// ContractExecutionLog represents a log of a smart contract execution.
type ContractExecutionLog struct {
	LogID          string    // Unique identifier for the log
	ContractID     string    // ID of the contract being executed
	ExecutedBy     string    // Address of the entity executing the contract
	ExecutionTime  time.Time // Timestamp of the execution
	InputData      string    // Input data for the contract execution
	OutputData     string    // Output data from the contract execution
	Status         string    // Status of the execution (e.g., "success", "failure")
}

// MintRecord represents a record of token minting.
type MintRecord struct {
	RecordID       string    // Unique identifier for the mint record
	TokenID        string    // ID of the token being minted
	Amount         *big.Int  // Amount of tokens minted
	MintedBy       string    // Address of the entity that minted the tokens
	Timestamp      time.Time // Timestamp of the minting
}

// BurnRecord represents a record of token burning.
type BurnRecord struct {
	RecordID       string    // Unique identifier for the burn record
	TokenID        string    // ID of the token being burned
	Amount         *big.Int  // Amount of tokens burned
	BurnedBy       string    // Address of the entity that burned the tokens
	Timestamp      time.Time // Timestamp of the burning
}

// CurrencyExchange represents a currency exchange transaction.
type CurrencyExchange struct {
    ExchangeID     string    // Unique identifier for the exchange
    FromCurrency   string    // Currency being exchanged
    ToCurrency     string    // Currency being received
    Amount         *big.Int  // Amount being exchanged
    ExchangedAmount *big.Int // Amount received after the exchange
    ExchangeRate   float64   // Exchange rate applied
    ExecutedAt     time.Time // Timestamp of the exchange
}


// BalanceStatus represents the current status of an account's balance.
type BalanceStatus struct {
	AccountID   string    // Unique identifier for the account
	IsActive    bool      // Indicates if the balance is active
	IsFrozen    bool      // Indicates if the balance is frozen
	IsOnHold    bool      // Indicates if the balance is on hold
	FreezeReason string   // Reason for balance freeze, if applicable
	HoldReason   string   // Reason for balance hold, if applicable
	UpdatedAt   time.Time // Last time the status was updated
}


// Authorization represents permissions or access rights assigned to an account.
type Authorization struct {
    ID          string // Unique ID for the authorization
    Description string // Description of the authorization
    Permissions []string // List of permissions granted
}

// BalanceInfo represents detailed information about an account’s balance.
type BalanceInfo struct {
	AccountID          string               // Unique identifier for the account
	AvailableBalance   float64              // Available balance for transactions
	LockedBalance      float64              // Balance that is locked and cannot be used
	HeldBalance        float64              // Balance that is held (e.g., for pending transactions)
	LastTransactionID  string               // ID of the last transaction affecting the balance
	TotalDeposited     float64              // Total deposited into the account
	TotalWithdrawn     float64              // Total withdrawn from the account
	Status             ledger.BalanceStatus // Use ledger.BalanceStatus directly
	LastUpdated        time.Time            // Last time the balance was updated
}

// TrustAccount represents a trust account within the ledger.
type TrustAccount struct {
    ID      string
    Balance *big.Float // Use *big.Float for high precision balance
}

// Approval represents an approval granted for specific transactions or operations.
type Approval struct {
    ID             string    // Unique identifier for the approval
    AccountID      string    // ID of the account granting the approval
    ApproverID     string    // ID of the approving authority or user
    TransactionID  string    // ID of the related transaction (if applicable)
    ApprovedAt     time.Time // Timestamp of approval
    ExpiresAt      time.Time // Expiration timestamp for time-limited approvals
    Status         string    // Status of the approval (e.g., "Pending", "Approved", "Rejected")
    Remarks        string    // Optional remarks or comments related to the approval
}

// Allocation represents an allocation of resources or funds.
type Allocation struct {
	ID            string    // Unique identifier for the allocation
	AccountID     string    // ID of the account to which the allocation is associated
	ResourceType  string    // Type of resource being allocated (e.g., "funds", "CPU", "memory")
	AllocatedAt   time.Time // Timestamp when the allocation was made
	ExpiresAt     time.Time // Expiration time for the allocation (if applicable)
	Amount        float64   // Amount of resource allocated
	Status        string    // Current status of the allocation (e.g., "active", "expired", "pending")
	Remarks       string    // Optional remarks or comments related to the allocation
	ApprovalID    string    // ID of any approval associated with the allocation
	Allocated bool      // Indicates if the allocation is active (not yet redeemed)

}
