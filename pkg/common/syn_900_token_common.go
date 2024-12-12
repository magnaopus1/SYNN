package common

import (
	"fmt"
	"math/big"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// SYN900Token represents the structure of an identity token (SYN900).
type SYN900Token struct {
	mutex          sync.Mutex                 // For thread safety
	TokenID        string                     // Unique identifier for the identity token
	Owner          string                     // Owner of the identity token (typically a wallet address)
	Metadata       *IdentityMetadata          // Metadata associated with the identity
	EncryptedData  string                     // Encrypted token data for security
	Status         string                     // Status of the token (e.g., active, revoked, pending verification)
	VerificationLog []VerificationRecord      // Record of verifications
	AuditTrail     []AuditRecord              // Audit trail of the token's transactions
	ComplianceLog  []ComplianceRecord         // Compliance logs related to this token
	Ledger         *ledger.Ledger             // Reference to the ledger for recording transactions
	Encryption     *Encryption     // Encryption service for securing data
}

// SYN900TokenManager handles all operations for SYN900 tokens.
type SYN900TokenManager struct {
	mutex        sync.Mutex                 // For thread safety
	Tokens       map[string]*SYN900Token    // In-memory store for SYN900 tokens
	Ledger       *ledger.Ledger             // Ledger for recording transactions
	Encryption   *Encryption     // Encryption service for securing data
}

// Verifier represents the structure of a verifier for SYN900 tokens.
type Syn900Verifier struct {
	VerifierID        string                 // Unique identifier for the verifier
	RegisteredTokens  map[string]*SYN900Token // Tokens registered to the verifier
	LedgerInstance    *ledger.Ledger         // Ledger instance for recording verification events
	EncryptionService *Encryption // Encryption service for securing verification data
	mutex             sync.Mutex             // Mutex for thread-safe operations
}

// Syn900Validator handles the validation of wallets and ensures only one wallet is connected at a time.
type Syn900Validator struct {
	mutex           sync.Mutex               // For thread safety
	ValidatorID     string                   // Unique identifier for the validator
	ConnectedWallet *Wallet                  // The currently connected wallet (only one can be connected at a time)
	Ledger          *ledger.Ledger           // Reference to the ledger for logging validations
	Encryption      *Encryption   // Encryption service for securing wallet validation data
}

// IDVerificationNode represents a node responsible for verifying SYN900 token deployments.
type IDVerificationNode struct {
	NodeID       string // Unique identifier for the node
	Verified     bool   // Whether the node has already verified the token
	Confirmation bool   // Whether the node confirmed or rejected the token
	NodeType     string // Type of node (e.g., Exchange, Bank, Government, etc.)
}

// NodeTypeList defines the valid node types for SYN900 verification.
var NodeTypeList = []string{
	"Exchange",   // Exchange node
	"Bank",       // Bank node
	"Government", // Government node
	"CentralBank",// Central Bank node
	"Regulator",  // Regulatory node
	"Creditor",   // Creditor node
}

// TokenVerificationProcess manages the process of authorizing SYN900 tokens through ID verification nodes.
type TokenVerificationProcess struct {
	mutex                sync.Mutex                             // For thread safety
	TokenID              string                                 // Unique token ID for the SYN900 token
	OwnerID              string                                 // ID of the token owner
	PendingNodes         []*IDVerificationNode                  // List of ID verification nodes selected to verify the token
	Confirmations        int                                    // Number of confirmations received
	Rejections           int                                    // Number of rejections received
	Status               string                                 // Current status of the token authorization process (Pending, Confirmed, Rejected)
	Ledger               *ledger.Ledger                         // Ledger reference for logging
	Encryption           *Encryption                 // Encryption service for security
	MaxNodes             int                                    // Maximum number of nodes for verification (5)
	InitialNodes         int                                    // Initial nodes for verification (3)
	RequiredConfirmations int                                    // Required confirmations to pass (3)
	RequiredRejections   int                                    // Rejections threshold for failure (3)
}

// SYN900Transaction represents a transaction for the SYN900 identity token.
type SYN900Transaction struct {
	TransactionID  string             // Unique identifier for the transaction
	Sender         string             // Sender's address
	Receiver       string             // Receiver's address (can be a wallet or identity)
	Amount         *big.Int           // Amount to be transferred (if applicable)
	TokenID        string             // SYN900 Token ID (for identity verification)
	Metadata       *IdentityMetadata  // Identity metadata associated with the transaction
	EncryptedData  string             // Encrypted transaction data
	TransactionLog []TransactionLog   // Transaction log for audit trail
	Status         string             // Status of the transaction (e.g., pending, confirmed)
}

// TransactionLog holds the log entries for the transaction.
type TransactionLog struct {
	Timestamp string // Timestamp of the transaction
	Action    string // Action performed (e.g., transfer, mint, burn)
	Details   string // Additional details about the action
}

// SYN900TransactionManager manages all transactions related to SYN900 tokens.
type SYN900TransactionManager struct {
	mutex       sync.Mutex                    // For thread safety
	Ledger      *ledger.Ledger                // Reference to the ledger for transaction logging
	Encryption  *Encryption        // Encryption service for secure data storage
	Transactions map[string]*SYN900Transaction // In-memory store for transactions
}

// SYN900Storage handles the storage of SYN900 identity token data.
type SYN900Storage struct {
	mutex      sync.Mutex                 // For thread safety
	Identities map[string]*IdentityMetadata // Map of TokenID to IdentityMetadata
	Ledger     *ledger.Ledger             // Ledger for permanent storage and transaction logging
	Encryption *Encryption     // Encryption service for secure storage
}

// IdentityMetadata represents detailed personal information for an identity token.
type IdentityMetadata struct {
	FullName            string               `json:"full_name"`
	DateOfBirth         string               `json:"date_of_birth"`
	Nationality         string               `json:"nationality"`
	PhotographHash      string               `json:"photograph_hash"`
	PhysicalAddress     string               `json:"physical_address"`
	DrivingLicense      string               `json:"driving_license"`
	EncryptedPassNum    string               `json:"encrypted_pass_num"`
	TokenID             string               `json:"token_id"`
	Owner               string               `json:"owner"`
	VerificationLog     []VerificationRecord `json:"verification_log"`
	AuditTrail          []AuditRecord        `json:"audit_trail"`
	ComplianceRecords   []ComplianceRecord   `json:"compliance_records"`
	RegisteredWallets   []WalletAddress      `json:"registered_wallets"`
	IsKYCCompliant  bool      // Indicates if the token meets KYC compliance
    IsAMLCompliant  bool      // Indicates if the token meets AML compliance
	ExpirationDate  time.Time // Expiration date of the token
}

// VerificationRecord logs the results of the verification process.
type VerificationRecord struct {
	NodeID    string `json:"node_id"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"` // Confirmed or Rejected
}

// AuditRecord logs the history of actions taken on the identity metadata.
type AuditRecord struct {
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
	Actor     string `json:"actor"` // Who performed the action
}

// WalletAddress represents wallets registered to the identity.
type WalletAddress struct {
	WalletID string `json:"wallet_id"`
}

// OwnershipRecord represents the structure for tracking ownership of SYN900 tokens.
type OwnershipRecord struct {
	TokenID       string                     // Unique identifier for the SYN900 token
	Owner         string                     // Current owner of the token (identity holder)
	PreviousOwner string                     // This will remain empty since ownership cannot be transferred
	EncryptedData string                     // Encrypted ownership data for security
}

// OwnershipManager manages SYN900 token ownership and ensures non-transferability.
type OwnershipManager struct {
	mutex        sync.Mutex                        // For thread safety
	OwnershipLog map[string]*OwnershipRecord       // In-memory store for ownership logs
	Ledger       *ledger.Ledger                    // Reference to the ledger for logging ownership changes
	Encryption   *Encryption            // Encryption service for securing ownership data
}

// MintingRecord stores details about a minted SYN900 token.
type MintingRecord struct {
	TokenID       string    // Unique ID for the SYN900 token
	Owner         string    // The owner/identity holder of the token
	MintedAt      time.Time // Timestamp of when the token was minted
	EncryptedData string    // Encrypted minting data
}

// MintingManager handles the minting of SYN900 tokens, ensuring validation and recording.
type MintingManager struct {
	mutex      sync.Mutex                        // For thread-safe operations
	Ledger     *ledger.Ledger                    // Ledger for recording all minting activities
	Encryption *Encryption            // Encryption service for securing minting data
	Tokens     map[string]*MintingRecord         // Map of all minted tokens
}

// SYN900Factory handles the creation and authorization of SYN900 tokens.
type SYN900Factory struct {
	mutex         sync.Mutex                 // For thread-safe operations
	Ledger        *ledger.Ledger             // Ledger instance for recording all token operations
	Encryption    *Encryption     // Encryption service for security
	TokenSupply   map[string]*SYN900Token    // Map to store all tokens by TokenID
	MintingMgr    *MintingManager            // Manager for minting SYN900 tokens
	AuthProcesses map[string]*AuthProcess    // Tracks authorization processes for tokens
}

// AuthProcess represents the authorization process for a SYN900 token.
type AuthProcess struct {
	TokenID        string            // Token ID undergoing authorization
	Owner          string            // Owner of the token
	VerificationNodes []*VerificationNode // Selected ID verification nodes
	Confirmations  int               // Number of confirmations
	Rejections     int               // Number of rejections
	IsCompleted    bool              // Whether the authorization is complete
	Status         string            // Status of the authorization process
}

// VerificationNode represents an ID verification node
type VerificationNode struct {
	NodeID   string
	NodeType string // Bank, Government, Central Bank, Regulator, Creditor, etc.
	Approved bool
}

// EventRecord represents a single event in the blockchain system.
type EventRecord struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`   // Event type (Minting, Transfer, Burn, etc.)
	Details        string    `json:"details"`      // Detailed description of the event
	Timestamp      time.Time `json:"timestamp"`    // Time of the event
	AssociatedToken string   `json:"associated_token"` // Token ID associated with this event
	Initiator      string    `json:"initiator"`    // Address of the event initiator
	EncryptedData  string    `json:"encrypted_data"`// Encrypted event data
}

// EventLogger manages the logging of events in the blockchain system.
type EventLogger struct {
	mutex      sync.Mutex
	EventList  []*EventRecord
	Ledger     *ledger.Ledger             // Reference to the ledger for logging
	Encryption *Encryption     // Encryption service
}

// DeploymentRequest represents the request to deploy a Syn900 token.
type DeploymentRequest struct {
	TokenID          string `json:"token_id"`
	Owner            string `json:"owner"`
	TokenDetails     string `json:"token_details"`     // Basic details of the token
	EncryptedDetails string `json:"encrypted_details"` // Encrypted details for security
	Status           string `json:"status"`            // Status of the deployment (Pending, Confirmed, Rejected)
	ValidatorNodes   []string `json:"validator_nodes"` // Nodes assigned to validate
	Confirmations    int    `json:"confirmations"`     // Number of confirmations received
	Rejections       int    `json:"rejections"`        // Number of rejections received
}

// TokenDeploymentManager manages the deployment of Syn900 tokens.
type TokenDeploymentManager struct {
	mutex         sync.Mutex
	Deployments   map[string]*DeploymentRequest // Map of TokenID to deployment requests
	Ledger        *ledger.Ledger                // Reference to the ledger for logging
	Encryption    *Encryption        // Encryption service
	ValidatorPool []string                      // List of available validator nodes
}

// Syn900BurnManager handles the burning of Syn900 tokens.
type Syn900BurnManager struct {
	mutex      sync.Mutex
	Ledger     *ledger.Ledger             // Reference to the ledger for recording transactions
	Encryption *Encryption     // Encryption service
}

// Syn900BatchTransferManager manages batch transfers of Syn900 tokens.
type Syn900BatchTransferManager struct {
	mutex      sync.Mutex
	Ledger     *ledger.Ledger             // Reference to the ledger for recording transactions
	Encryption *Encryption     // Encryption service
}


// ComplianceRecord stores the compliance check data for a specific action or transaction
type ComplianceRecord struct {
    ActionID      string          // Unique identifier for the action or transaction
    Status        ComplianceStatus // Status of the compliance check
    CheckedBy     string          // Compliance officer or module responsible for the check
    EncryptedData string          // Field to hold encrypted data
}

// ComplianceStatus represents the result of a compliance check
type ComplianceStatus struct {
	IsValid   bool      // Whether the action/transaction complies with rules
	Reason    string    // Reason for failure (if applicable)
	Timestamp time.Time // Timestamp of the compliance check
}

// ConvertBigIntToFloat64 safely converts a *big.Int to float64
func ConvertBigIntToFloat64(balance *big.Int) float64 {
    floatBalance, _ := new(big.Float).SetInt(balance).Float64() // Convert *big.Int to *big.Float, then to float64
    return floatBalance
}

// CheckBalanceOnChain checks the balance of the token in the ledger
func (t *SYN900Token) CheckBalanceOnChain(ledger *ledger.Ledger) (float64, error) {
    balance, err := ledger.GetTokenBalance(t.TokenID)
    if err != nil {
        return 0, fmt.Errorf("failed to get balance for token %s: %v", t.TokenID, err)
    }
    
    // Convert *big.Int to float64
    return ConvertBigIntToFloat64(balance), nil
}
// IsFrozenOnChain checks if the token is frozen in the ledger
func (t *SYN900Token) IsFrozenOnChain(ledger *ledger.Ledger) (bool, error) {
    isFrozen, err := ledger.IsTokenFrozen(t.TokenID)
    if err != nil {
        return false, fmt.Errorf("failed to check if token %s is frozen: %v", t.TokenID, err)
    }
    return isFrozen, nil
}

// IsBurnedOnChain checks if the token has been burned in the ledger
func (t *SYN900Token) IsBurnedOnChain(ledger *ledger.Ledger) (bool, error) {
    isBurned, err := ledger.IsTokenBurned(t.TokenID)
    if err != nil {
        return false, fmt.Errorf("failed to check if token %s is burned: %v", t.TokenID, err)
    }
    return isBurned, nil
}

// IsExpiredOnChain checks if the token is expired in the ledger
func (t *SYN900Token) IsExpiredOnChain(ledger *ledger.Ledger) (bool, error) {
    isExpired, err := ledger.IsTokenExpired(t.TokenID)
    if err != nil {
        return false, fmt.Errorf("failed to check if token %s is expired: %v", t.TokenID, err)
    }
    return isExpired, nil
}
