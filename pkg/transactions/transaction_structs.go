package transactions

import (
	"log"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/smart_contract"
	"time"
)

// TransactionRecord keeps track of all transactions in the ledger
type TransactionRecord struct {
    From        string  // Sender's address
    To          string  // Recipient's address
    Amount      float64 // Amount transferred
    Fee         float64 // Transaction fee
    Hash        string  // Transaction hash (unique ID)
    Status      string  // Status of the transaction (e.g., "pending", "confirmed")
    BlockIndex  int     // Block in which the transaction was confirmed
}



// TransactionPool manages the unconfirmed transactions waiting to be validated.
type TransactionPool struct {
    transactions      map[string]*common.Transaction  // Map of transaction ID to Transaction
    pendingSubBlocks  map[string][]*common.Transaction // Map of sub-block IDs to transactions
    mu                sync.Mutex
    maxPoolSize       int
    ledger            *ledger.Ledger
    encryptionService *common.Encryption              // Pointer to Encryption service
}

// EscrowTransaction holds the details of the escrow agreement.
type EscrowTransaction struct {
    EscrowID      string
    SenderID      string
    ReceiverID    string
    Amount        float64
    Status        EscrowStatus
    CreationTime  time.Time
    ReleaseTime   time.Time
    Condition     string // Optional: Condition to release funds (e.g., service completion)
}

// EscrowManager manages escrow transactions and ensures the proper release or cancellation of funds.
type EscrowManager struct {
    ledgerInstance *ledger.Ledger
    mutex          sync.Mutex
    escrows        map[string]*EscrowTransaction
}

// SmartLegalContractManager manages multiple smart legal contracts.
type SmartLegalContractManager struct {
	Contracts      map[string]*smart_contract.SmartLegalContract // All deployed legal contracts
	LedgerInstance *ledger.Ledger                 // Ledger for recording contract deployments and executions
	mutex          sync.Mutex                     // Mutex for safe concurrency
}

// TransactionCancellationManager handles transaction cancellation requests and processing.
type TransactionCancellationManager struct {
	mutex              sync.Mutex
	Consensus          *common.SynnergyConsensus
	TimeoutPeriod      time.Duration              // Time allowed for requesting a cancellation
	ResponseTimeout    time.Duration              // Time allowed for node response
	NotificationService common.Notification       // Notifies involved parties
	Ledger             *ledger.Ledger
	Encryption         *common.Encryption
    Logger             *log.Logger             // Standard Go logger

}

type TransactionReversalManager struct {
	mutex              sync.Mutex
	Consensus          *common.SynnergyConsensus
	Ledger             *ledger.Ledger
	ReversalTimeLimit  time.Duration            // Time allowed for requesting reversal (within 28 days)
	NotificationService common.Notification     // Notifies involved parties
	Encryption         *common.Encryption
	Logger             *log.Logger              // Logger for logging events (standard Go log package)
}





// TransactionMetricsManager manages the collection of transaction metrics in the blockchain.
type TransactionMetricsManager struct {
    totalTransactions     int
    totalSubBlocks        int
    totalBlocks           int
    totalGasConsumed      int
    totalFeesCollected    float64
    transactionThroughput float64 // transactions per second (TPS)
    gasEfficiency         float64 // ratio of gas used vs gas limit

    metricsLock sync.Mutex
    ledger      *ledger.Ledger
}

// TransactionReceiptManager manages the creation, storage, and validation of transaction receipts.
type TransactionReceiptManager struct {
    receipts           map[string]*TransactionReceipt // Map of transaction IDs to their receipts
    encryptionService  *common.Encryption         // Encryption service for receipt integrity
}

// EscrowStatus represents the status of an escrow transaction.
type EscrowStatus string

const (
    EscrowStatusPending   EscrowStatus = "Pending"
    EscrowStatusReleased  EscrowStatus = "Released"
    EscrowStatusCancelled EscrowStatus = "Cancelled"
)

// IsFinalized checks if the escrow status is finalized (either Released or Cancelled)
func (status EscrowStatus) IsFinalized() bool {
    return status == EscrowStatusReleased || status == EscrowStatusCancelled
}
