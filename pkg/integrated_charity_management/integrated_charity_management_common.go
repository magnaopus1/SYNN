package integrated_charity_management

import (
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// CharityPoolManagement handles the distribution of transaction fees into internal and external charity pools
type CharityPoolManagement struct {
	mutex               sync.Mutex
	InternalPoolBalance float64        // Balance for the internal charity pool
	ExternalPoolBalance float64        // Balance for the external charity pool
	LedgerInstance      *ledger.Ledger // Ledger instance for tracking pool activities
}

// CharityPool represents the external and internal charity pools and manages their balances
type CharityPool struct {
	externalPool   float64        // External charity pool balance
	internalPool   float64        // Internal charity pool balance
	totalBalance   float64        // Total balance to be distributed between both pools
	LedgerInstance *ledger.Ledger // Ledger instance for tracking pool activities
	mutex          sync.Mutex     // Mutex for thread-safe operations
}

// CharityProposal represents a charity that enters into the external charity pool
type CharityProposal struct {
	CharityID     string    // Unique ID for the charity
	Name          string    // Charity name
	CharityNumber string    // Registered charity number
	Description   string    // Charity description
	Website       string    // Charity website
	Addresses     []string  // Encrypted addresses to receive funds
	CreatedAt     time.Time // Timestamp of when the charity entered the pool
	VoteCount     int       // Number of votes received
	IsValid       bool      // Is the charity valid
}

// ExternalCharityPoolManager manages the external charity proposal process and fund distribution
type ExternalCharityPoolManager struct {
	mutex               sync.Mutex
	CurrentCycle        []*CharityProposal       // Charities in the current 90-day cycle
	CharityEntries      map[string]*CharityProposal // Entries for the current round
	ProposalStart       time.Time               // Start of the proposal submission
	VotingEnd           time.Time               // End of the voting period
	LedgerInstance      *ledger.Ledger          // Ledger instance for tracking charity activity
	ExternalPoolBalance float64                 // Balance of the external charity pool
}

// InternalCharityPool manages the internal charity pool, distributing funds every 24 hours
type InternalCharityPool struct {
	mutex            sync.Mutex
	PoolBalance      float64                // Current balance of the internal charity pool
	WalletAddresses  map[string]float64     // Map of charity wallet addresses to their respective balances
	OwnerAddress     string                 // Blockchain owner's address (for access control)
	LedgerInstance   *ledger.Ledger         // Ledger instance for tracking all transactions and activities
	stopChan         chan bool              // Channel to stop the 24-hour distribution
}
