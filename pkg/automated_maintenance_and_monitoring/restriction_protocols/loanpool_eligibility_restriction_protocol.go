package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	EligibilityCheckInterval      = 15 * time.Second // Interval for checking loan pool eligibility
	MaxLoanDefaultsAllowed        = 2               // Maximum number of defaults allowed for user eligibility
	MinCollateralRatio            = 1.5             // Minimum collateral-to-loan ratio for eligibility
)

// LoanPoolEligibilityRestrictionAutomation monitors and restricts loan pool eligibility across the network
type LoanPoolEligibilityRestrictionAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	userLoanDefaults      map[string]int   // Tracks the number of loan defaults per user
	userCollateralRatios  map[string]float64 // Tracks the user's collateral-to-loan ratio
}

// NewLoanPoolEligibilityRestrictionAutomation initializes and returns an instance of LoanPoolEligibilityRestrictionAutomation
func NewLoanPoolEligibilityRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *LoanPoolEligibilityRestrictionAutomation {
	return &LoanPoolEligibilityRestrictionAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		userLoanDefaults:     make(map[string]int),
		userCollateralRatios: make(map[string]float64),
	}
}

// StartEligibilityMonitoring starts continuous monitoring of loan pool eligibility
func (automation *LoanPoolEligibilityRestrictionAutomation) StartEligibilityMonitoring() {
	ticker := time.NewTicker(EligibilityCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorEligibility()
		}
	}()
}

// monitorEligibility checks each user's loan eligibility based on collateral ratio and default history
func (automation *LoanPoolEligibilityRestrictionAutomation) monitorEligibility() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch user eligibility data from Synnergy Consensus
	userData := automation.consensusSystem.GetLoanUserData()

	for userID, data := range userData {
		// Check if the user has exceeded the allowed number of loan defaults
		if automation.userLoanDefaults[userID] > MaxLoanDefaultsAllowed {
			automation.flagEligibilityViolation(userID, "Exceeded maximum number of loan defaults")
			continue
		}

		// Check if the user's collateral-to-loan ratio meets the minimum requirement
		if automation.userCollateralRatios[userID] < MinCollateralRatio {
			automation.flagEligibilityViolation(userID, "Collateral-to-loan ratio below minimum threshold")
		}
	}
}

// flagEligibilityViolation flags a violation of loan eligibility rules and logs it in the ledger
func (automation *LoanPoolEligibilityRestrictionAutomation) flagEligibilityViolation(userID string, reason string) {
	fmt.Printf("Loan pool eligibility violation: User ID %s, Reason: %s\n", userID, reason)

	// Log the violation in the ledger
	automation.logEligibilityViolation(userID, reason)
}

// logEligibilityViolation logs the flagged loan eligibility violation into the ledger with full details
func (automation *LoanPoolEligibilityRestrictionAutomation) logEligibilityViolation(userID string, violationReason string) {
	// Create a ledger entry for loan pool eligibility violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("loan-eligibility-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Loan Pool Eligibility Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated loan pool eligibility rules. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptViolationData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log loan pool eligibility violation:", err)
	} else {
		fmt.Println("Loan pool eligibility violation logged.")
	}
}

// encryptViolationData encrypts the violation data before logging for security
func (automation *LoanPoolEligibilityRestrictionAutomation) encryptViolationData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting violation data:", err)
		return data
	}
	return string(encryptedData)
}
