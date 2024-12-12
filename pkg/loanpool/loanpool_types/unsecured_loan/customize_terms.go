package loanpool

import (
	"errors"
	"fmt"

	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/common"
)


// NewUnsecuredLoanTermManager initializes a new manager for handling unsecured loan terms.
func NewUnsecuredLoanTermManager(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *encryption.Encryption) *common.UnsecuredLoanTermManager {
	return &common.UnsecuredLoanTermManager{
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		LoanTermRecords:   make(map[string]*common.LoanTerms),
		EncryptionService: encryptionService,
	}
}

// CustomizeLoanTerms allows users to customize their loan terms, including repayment length, amount, and selecting Islamic finance.
func (ltm *common.UnsecuredLoanTermManager) CustomizeLoanTerms(loanID string, repaymentLength int, amountBorrowed, interestRate, feeOnTop float64, useIslamicFinance bool) (*common.LoanTerms, error) {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	// Validate inputs
	if repaymentLength <= 0 || amountBorrowed <= 0 {
		return nil, errors.New("invalid repayment length or loan amount")
	}

	// Create customized loan terms
	loanTerms := &common.LoanTerms{
		RepaymentLength: repaymentLength,
		AmountBorrowed:  amountBorrowed,
		InterestRate:    interestRate,
		IslamicFinance:  useIslamicFinance,
		FeeOnTop:        feeOnTop,
	}

	// Calculate the total repayment amount based on the loan terms
	if useIslamicFinance {
		loanTerms.TotalRepaymentAmount = loanTerms.AmountBorrowed + feeOnTop
	} else {
		interest := loanTerms.AmountBorrowed * (loanTerms.InterestRate / 100)
		loanTerms.TotalRepaymentAmount = loanTerms.AmountBorrowed + interest
	}

	// Encrypt the loan terms before storing them
	encryptedLoanTerms, err := ltm.EncryptionService.EncryptData(fmt.Sprintf("%v", loanTerms), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt loan terms: %v", err)
	}

	// Store the encrypted loan terms in the ledger
	err = ltm.Ledger.RecordLoanTerms(loanID, encryptedLoanTerms)
	if err != nil {
		return nil, fmt.Errorf("failed to record loan terms in ledger: %v", err)
	}

	// Store the loan terms in the local record
	ltm.LoanTermRecords[loanID] = loanTerms

	fmt.Printf("Customized loan terms for loan %s successfully created.\n", loanID)
	return loanTerms, nil
}

// SwitchToIslamicFinance allows users to switch their loan to Islamic finance terms (if not already).
func (ltm *common.UnsecuredLoanTermManager) SwitchToIslamicFinance(loanID string, feeOnTop float64) (*common.LoanTerms, error) {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	// Retrieve existing loan terms
	loanTerms, exists := ltm.LoanTermRecords[loanID]
	if !exists {
		return nil, errors.New("loan terms not found for this loan ID")
	}

	// If already Islamic finance, return an error
	if loanTerms.IslamicFinance {
		return nil, errors.New("loan is already under Islamic finance terms")
	}

	// Switch to Islamic finance terms: remove interest, apply a flat fee on top
	loanTerms.IslamicFinance = true
	loanTerms.FeeOnTop = feeOnTop
	loanTerms.TotalRepaymentAmount = loanTerms.AmountBorrowed + feeOnTop

	// Encrypt the updated loan terms
	encryptedLoanTerms, err := ltm.EncryptionService.EncryptData(fmt.Sprintf("%v", loanTerms), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt updated loan terms: %v", err)
	}

	// Update the loan terms in the ledger
	err = ltm.Ledger.UpdateLoanTerms(loanID, encryptedLoanTerms)
	if err != nil {
		return nil, fmt.Errorf("failed to update loan terms in ledger: %v", err)
	}

	fmt.Printf("Loan %s successfully switched to Islamic finance terms.\n", loanID)
	return loanTerms, nil
}

// ValidateLoanTerms validates the customized loan terms using the Synnergy Consensus mechanism.
func (ltm *common.UnsecuredLoanTermManager) ValidateLoanTerms(loanID string) (bool, error) {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	// Retrieve the loan terms
	loanTerms, exists := ltm.LoanTermRecords[loanID]
	if !exists {
		return false, errors.New("loan terms not found for this loan ID")
	}

	// Use Synnergy Consensus to validate the loan terms
	validated, err := ltm.Consensus.ValidateLoanTerms(loanID, loanTerms.TotalRepaymentAmount)
	if err != nil || !validated {
		return false, fmt.Errorf("failed to validate loan terms for loan %s: %v", loanID, err)
	}

	fmt.Printf("Loan terms for loan %s have been validated.\n", loanID)
	return true, nil
}

// ViewLoanTerms allows users to view their customized loan terms.
func (ltm *common.UnsecuredLoanTermManager) ViewLoanTerms(loanID string) (*common.LoanTerms, error) {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	loanTerms, exists := ltm.LoanTermRecords[loanID]
	if !exists {
		return nil, errors.New("loan terms not found for this loan ID")
	}

	return loanTerms, nil
}
