package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/lending"
)

const (
	LendingCheckInterval          = 10 * time.Minute // Interval to monitor lending protocol health
	DefaultLiquidationThreshold   = 0.8              // Collateral-to-loan value ratio threshold for liquidation
	MaxLoanDuration               = 365 * 24 * time.Hour // Maximum loan duration (1 year)
	LendingLedgerEntryType        = "Lending Protocol Event"
	LiquidationLedgerEntryType    = "Loan Liquidation"
)

// LendingProtocolExecutionAutomation handles the monitoring and execution of lending contracts
type LendingProtocolExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
	ledgerInstance    *ledger.Ledger                        // Ledger instance for logging events
	lendingManager    *lending.Manager                      // Lending protocol manager
	executionMutex    *sync.RWMutex                         // Mutex for thread-safe lending execution
}

// NewLendingProtocolExecutionAutomation initializes lending protocol automation
func NewLendingProtocolExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, lendingManager *lending.Manager, executionMutex *sync.RWMutex) *LendingProtocolExecutionAutomation {
	return &LendingProtocolExecutionAutomation{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		lendingManager:    lendingManager,
		executionMutex:    executionMutex,
	}
}

// StartLendingProtocolMonitor starts the continuous monitoring of lending contracts
func (automation *LendingProtocolExecutionAutomation) StartLendingProtocolMonitor() {
	ticker := time.NewTicker(LendingCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAndExecuteLendingActions()
		}
	}()
}

// checkAndExecuteLendingActions monitors the lending contracts and executes required actions
func (automation *LendingProtocolExecutionAutomation) checkAndExecuteLendingActions() {
	automation.executionMutex.Lock()
	defer automation.executionMutex.Unlock()

	// Fetch all active loans from the lending manager
	activeLoans := automation.lendingManager.GetActiveLoans()

	for _, loan := range activeLoans {
		automation.evaluateLoanHealth(loan)
	}
}

// evaluateLoanHealth checks if any actions are required based on the loan's collateralization ratio or duration
func (automation *LendingProtocolExecutionAutomation) evaluateLoanHealth(loan *lending.Loan) {
	collateralRatio := loan.GetCollateralToLoanRatio()

	if collateralRatio < DefaultLiquidationThreshold {
		automation.triggerLiquidation(loan)
	}

	// Check if the loan exceeds its maximum duration
	if time.Since(loan.StartTime) > MaxLoanDuration {
		automation.triggerLoanClosure(loan)
	}
}

// triggerLiquidation executes a loan liquidation when collateral falls below the required threshold
func (automation *LendingProtocolExecutionAutomation) triggerLiquidation(loan *lending.Loan) {
	fmt.Printf("Triggering liquidation for loan %s (Collateral ratio: %.2f)\n", loan.ID, loan.GetCollateralToLoanRatio())

	err := automation.lendingManager.LiquidateLoan(loan)
	if err != nil {
		fmt.Printf("Failed to liquidate loan %s: %v\n", loan.ID, err)
		return
	}

	// Log the liquidation into the ledger
	automation.logLiquidationInLedger(loan)

	// Notify the consensus engine about the liquidation
	automation.consensusEngine.NotifyLoanLiquidation(loan.ID)
}

// triggerLoanClosure handles loan closure if the loan duration exceeds the maximum allowed duration
func (automation *LendingProtocolExecutionAutomation) triggerLoanClosure(loan *lending.Loan) {
	fmt.Printf("Closing loan %s due to exceeded duration.\n", loan.ID)

	err := automation.lendingManager.CloseLoan(loan)
	if err != nil {
		fmt.Printf("Failed to close loan %s: %v\n", loan.ID, err)
		return
	}

	// Log the loan closure in the ledger
	automation.logLoanClosureInLedger(loan)

	// Notify the consensus engine about the loan closure
	automation.consensusEngine.NotifyLoanClosure(loan.ID)
}

// logLiquidationInLedger securely logs a loan liquidation event in the ledger
func (automation *LendingProtocolExecutionAutomation) logLiquidationInLedger(loan *lending.Loan) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("loan-liquidation-%s-%d", loan.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      LiquidationLedgerEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("Liquidated loan %s due to collateral ratio %.2f", loan.ID, loan.GetCollateralToLoanRatio()),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log loan liquidation for loan %s: %v\n", loan.ID, err)
	} else {
		fmt.Printf("Loan liquidation logged successfully for loan %s.\n", loan.ID)
	}
}

// logLoanClosureInLedger securely logs a loan closure event in the ledger
func (automation *LendingProtocolExecutionAutomation) logLoanClosureInLedger(loan *lending.Loan) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("loan-closure-%s-%d", loan.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      LendingLedgerEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("Closed loan %s due to exceeded duration.", loan.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log loan closure for loan %s: %v\n", loan.ID, err)
	} else {
		fmt.Printf("Loan closure logged successfully for loan %s.\n", loan.ID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *LendingProtocolExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLiquidation allows an administrator to manually trigger liquidation for a specific loan
func (automation *LendingProtocolExecutionAutomation) TriggerManualLiquidation(loanID string) {
	fmt.Printf("Manually triggering liquidation for loan %s...\n", loanID)

	loan := automation.lendingManager.GetLoanByID(loanID)
	if loan != nil {
		automation.triggerLiquidation(loan)
	} else {
		fmt.Printf("Loan %s not found.\n", loanID)
	}
}
