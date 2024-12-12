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
	LoanCheckInterval       = 15 * time.Second  // Interval for checking DeFi loan activities
	MaxLoansPerUser         = 5                 // Maximum number of active loans allowed per user
	MaxLoanAmount           = 50000.0           // Maximum loan amount allowed per loan
	LoanDurationWindow      = 24 * 7 * time.Hour // Time window for loan validation
)

// DefiLoanRestrictionAutomation monitors and restricts decentralized finance loan activities across the network
type DefiLoanRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	userLoanCount          map[string]int    // Tracks active loans per user
}

// NewDefiLoanRestrictionAutomation initializes and returns an instance of DefiLoanRestrictionAutomation
func NewDefiLoanRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DefiLoanRestrictionAutomation {
	return &DefiLoanRestrictionAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		userLoanCount:       make(map[string]int),
	}
}

// StartLoanMonitoring starts continuous monitoring of DeFi loans
func (automation *DefiLoanRestrictionAutomation) StartLoanMonitoring() {
	ticker := time.NewTicker(LoanCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorLoans()
		}
	}()
}

// monitorLoans checks recent loan activities and enforces loan restrictions
func (automation *DefiLoanRestrictionAutomation) monitorLoans() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent loan activities from Synnergy Consensus
	recentLoans := automation.consensusSystem.GetRecentLoans()

	for _, loan := range recentLoans {
		// Validate loan limits
		if !automation.validateLoanLimit(loan) {
			automation.flagLoanViolation(loan, "Exceeded maximum number of active loans for this user")
		} else if !automation.validateLoanAmount(loan) {
			automation.flagLoanViolation(loan, "Loan amount exceeds the maximum allowed limit")
		}
	}
}

// validateLoanLimit checks if a user has exceeded the active loan limit within the time window
func (automation *DefiLoanRestrictionAutomation) validateLoanLimit(loan common.Loan) bool {
	currentLoanCount := automation.userLoanCount[loan.UserID]
	if currentLoanCount+1 > MaxLoansPerUser {
		return false
	}

	// Update the loan count for the user
	automation.userLoanCount[loan.UserID]++
	return true
}

// validateLoanAmount checks if the loan amount exceeds the maximum allowed amount
func (automation *DefiLoanRestrictionAutomation) validateLoanAmount(loan common.Loan) bool {
	return loan.Amount <= MaxLoanAmount
}

// flagLoanViolation flags a loan activity that violates system rules and logs it in the ledger
func (automation *DefiLoanRestrictionAutomation) flagLoanViolation(loan common.Loan, reason string) {
	fmt.Printf("DeFi loan violation: User %s, Reason: %s\n", loan.UserID, reason)

	// Log the violation into the ledger
	automation.logLoanViolation(loan, reason)
}

// logLoanViolation logs the flagged DeFi loan violation into the ledger with full details
func (automation *DefiLoanRestrictionAutomation) logLoanViolation(loan common.Loan, violationReason string) {
	// Encrypt the loan data before logging
	encryptedData := automation.encryptLoanData(loan)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("loan-violation-%s-%d", loan.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DeFi Loan Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for loan violation. Reason: %s. Encrypted Data: %s", loan.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log loan violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Loan violation logged for user: %s\n", loan.UserID)
	}
}

// encryptLoanData encrypts loan data before logging for security
func (automation *DefiLoanRestrictionAutomation) encryptLoanData(loan common.Loan) string {
	data := fmt.Sprintf("User ID: %s, Loan Amount: %.2f, Timestamp: %d", loan.UserID, loan.Amount, loan.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting loan data:", err)
		return data
	}
	return string(encryptedData)
}
