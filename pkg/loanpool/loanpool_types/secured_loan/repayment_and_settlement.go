package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)


// NewSecuredLoanRepaymentManager initializes a new loan repayment manager.
func NewSecuredLoanRepaymentManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, syn900Registry *syn900.Registry) *common.SecuredLoanRepaymentManager {
	return &common.SecuredLoanRepaymentManager{
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		Syn900Registry:    syn900Registry,
		LoanRepayments:    make(map[string]*common.LoanRepaymentDetails),
		DefaultThreshold:  6 * 30 * 24 * time.Hour, // Default if payments missed for 6 months
	}
}

// SetupRepaymentPlan sets up the repayment schedule for a loan.
func (rm *common.SecuredLoanRepaymentManager) SetupRepaymentPlan(loanID string, totalAmount float64, interestRate float64, proposerWallet string, repaymentDates []time.Time, authorityWallets []string, collateralContact string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if the loan already has a repayment plan.
	if _, exists := rm.LoanRepayments[loanID]; exists {
		return errors.New("repayment plan already exists for this loan")
	}

	// Create a new repayment plan.
	repaymentDetails := &common.LoanRepaymentDetails{
		LoanID:           loanID,
		ProposerWallet:   proposerWallet,
		TotalAmount:      totalAmount,
		RemainingAmount:  totalAmount,
		InterestRate:     interestRate,
		RepaymentDates:   repaymentDates,
		NextPaymentDue:   repaymentDates[0],
		Status:           "Active",
		CollateralContact: collateralContact,
		AuthorityWallets: authorityWallets,
	}

	// Store the repayment plan.
	rm.LoanRepayments[loanID] = repaymentDetails

	// Log the repayment plan creation in the ledger.
	err := rm.Ledger.RecordRepaymentPlan(loanID, repaymentDetails)
	if err != nil {
		return fmt.Errorf("failed to record repayment plan in ledger: %v", err)
	}

	fmt.Printf("Repayment plan for loan %s successfully set up.\n", loanID)
	return nil
}

// ChangePaymentDate allows the borrower to change a scheduled payment date.
func (rm *common.SecuredLoanRepaymentManager) ChangePaymentDate(loanID string, newDate time.Time) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	repaymentDetails, exists := rm.LoanRepayments[loanID]
	if !exists {
		return errors.New("repayment plan not found for this loan")
	}

	// Change the next payment due date.
	repaymentDetails.NextPaymentDue = newDate
	repaymentDetails.RepaymentDates = append(repaymentDetails.RepaymentDates, newDate)

	// Update the ledger with the new payment date.
	err := rm.Ledger.UpdatePaymentDate(loanID, newDate)
	if err != nil {
		return fmt.Errorf("failed to update payment date in ledger: %v", err)
	}

	fmt.Printf("Payment date for loan %s changed to %s.\n", loanID, newDate.Format(time.RFC3339))
	return nil
}

// ProcessRepayment handles the repayment of a scheduled payment.
func (rm *common.SecuredLoanRepaymentManager) ProcessRepayment(loanID string, paymentAmount float64) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	repaymentDetails, exists := rm.LoanRepayments[loanID]
	if !exists {
		return errors.New("repayment plan not found for this loan")
	}

	// Deduct the payment from the remaining amount.
	repaymentDetails.RemainingAmount -= paymentAmount
	if repaymentDetails.RemainingAmount < 0 {
		repaymentDetails.RemainingAmount = 0
	}

	// Distribute the interest to authority node wallets.
	interest := (repaymentDetails.InterestRate / 100) * paymentAmount
	rm.distributeInterestToAuthorityNodes(loanID, interest)

	// Update the ledger with the repayment.
	err := rm.Ledger.RecordRepayment(loanID, paymentAmount, repaymentDetails.RemainingAmount)
	if err != nil {
		return fmt.Errorf("failed to record repayment in ledger: %v", err)
	}

	// Check if the loan has been fully repaid.
	if repaymentDetails.RemainingAmount == 0 {
		rm.markLoanAsSatisfied(loanID)
	}

	fmt.Printf("Payment of %.2f processed for loan %s. Remaining amount: %.2f.\n", paymentAmount, loanID, repaymentDetails.RemainingAmount)
	return nil
}

// HandleDefault checks if the loan has defaulted and sends a notification.
func (rm *common.SecuredLoanRepaymentManager) HandleDefault(loanID string) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	repaymentDetails, exists := rm.LoanRepayments[loanID]
	if !exists {
		return errors.New("repayment plan not found for this loan")
	}

	// Check if the loan has defaulted (missed payments for 6 months).
	if time.Since(repaymentDetails.NextPaymentDue) > rm.DefaultThreshold {
		defaultTime := time.Now()
		repaymentDetails.Status = "Defaulted"
		repaymentDetails.DefaultedAt = &defaultTime

		// Record the default in the Syn900 registry.
		err := rm.Syn900Registry.RecordDefault(loanID, repaymentDetails.ProposerWallet)
		if err != nil {
			return fmt.Errorf("failed to record default in Syn900: %v", err)
		}

		// Notify the authority nodes via email for collateral request.
		rm.notifyAuthorityNodesForCollateral(repaymentDetails.CollateralContact)

		// Update the loan status in the ledger.
		err = rm.Ledger.RecordLoanDefault(loanID, defaultTime)
		if err != nil {
			return fmt.Errorf("failed to record default in ledger: %v", err)
		}

		fmt.Printf("Loan %s has defaulted. Collateral request sent.\n", loanID)
		return nil
	}

	return nil
}

// distributeInterestToAuthorityNodes distributes interest payments to the wallets of authority nodes.
func (rm *common.SecuredLoanRepaymentManager) distributeInterestToAuthorityNodes(loanID string, interest float64) {
	repaymentDetails, exists := rm.LoanRepayments[loanID]
	if !exists {
		fmt.Printf("Repayment plan not found for loan %s.\n", loanID)
		return
	}

	// Split the interest between the authority nodes.
	interestPerNode := interest / float64(len(repaymentDetails.AuthorityWallets))
	for _, wallet := range repaymentDetails.AuthorityWallets {
		err := rm.Ledger.RecordInterestPayment(wallet, interestPerNode)
		if err != nil {
			fmt.Printf("Failed to record interest payment to wallet %s: %v\n", wallet, err)
		} else {
			fmt.Printf("Interest payment of %.2f distributed to wallet %s.\n", interestPerNode, wallet)
		}
	}
}

// markLoanAsSatisfied updates the loan status to "Satisfied" once fully repaid.
func (rm *common.SecuredLoanRepaymentManager) markLoanAsSatisfied(loanID string) {
	repaymentDetails, exists := rm.LoanRepayments[loanID]
	if !exists {
		fmt.Printf("Repayment plan not found for loan %s.\n", loanID)
		return
	}

	repaymentDetails.Status = "Satisfied"

	// Update the ledger and Syn900 registry.
	err := rm.Ledger.RecordLoanSatisfaction(loanID)
	if err != nil {
		fmt.Printf("Failed to record loan satisfaction for loan %s: %v\n", loanID, err)
	}

	err = rm.Syn900Registry.UpdateLoanStatusToSatisfied(loanID, repaymentDetails.ProposerWallet)
	if err != nil {
		fmt.Printf("Failed to update loan satisfaction in Syn900 for loan %s: %v\n", loanID, err)
	}

	fmt.Printf("Loan %s has been marked as satisfied.\n", loanID)
}

// notifyAuthorityNodesForCollateral sends a notification to authority nodes for collateral request on loan default.
func (rm *common.SecuredLoanRepaymentManager) notifyAuthorityNodesForCollateral(collateralContact string) {
	// Implement the notification mechanism here (e.g., email).
	fmt.Printf("Collateral request sent to %s due to loan default.\n", collateralContact)
}
