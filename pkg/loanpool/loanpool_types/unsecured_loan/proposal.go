package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/tokens/syn900"
	"synnergy_network/pkg/common"
)

// NewProposalManager initializes a new manager for handling unsecured loan proposals.
func NewProposalManager(
	ledgerInstance *ledger.Ledger,
	consensusEngine *common.SynnergyConsensus,
	syn900Validator *syn900.Validator,
	encryptionService *encryption.Encryption,
	creditChecker *common.CreditCheckManager,
	affordabilityMgr *common.AffordabilityManager,
	termsManager *common.TermsCustomizationManager,
) *common.ProposalManager {
	return &common.ProposalManager{
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		Syn900Validator:   syn900Validator,
		EncryptionService: encryptionService,
		CreditChecker:     creditChecker,
		AffordabilityMgr:  affordabilityMgr,
		TermsManager:      termsManager,
		Proposals:         make(map[string]*common.UnsecuredLoanProposal),
	}
}

// StartProposalProcess begins the loan application process with the applicant details.
func (pm *common.ProposalManager) StartProposalProcess(applicantName, applicantID, walletAddress string) (*common.UnsecuredLoanProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Validate applicant ID using Syn900
	verified, err := pm.Syn900Validator.ValidateID(applicantID)
	if err != nil || !verified {
		return nil, errors.New("applicant ID validation failed with syn900")
	}

	// Create a new proposal
	loanID := common.GenerateUniqueID()
	proposal := &common.UnsecuredLoanProposal{
		LoanID:              loanID,
		ApplicantName:       applicantName,
		ApplicantID:         applicantID,
		WalletAddress:       walletAddress,
		SubmissionTimestamp: time.Now(),
		ProposalStatus:      "Pending",
		ApprovalStage:       "Application",
		LastUpdated:         time.Now(),
	}

	// Store the proposal in the ledger
	err = pm.Ledger.RecordProposal(proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal in the ledger: %v", err)
	}

	// Add proposal to in-memory map for further processing
	pm.Proposals[loanID] = proposal

	fmt.Printf("Unsecured loan application process started for applicant %s (Loan ID: %s).\n", applicantName, loanID)
	return proposal, nil
}

// RunCreditCheck performs the decentralized credit check for the applicant.
func (pm *common.ProposalManager) RunCreditCheck(loanID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve the proposal
	proposal, exists := pm.Proposals[loanID]
	if !exists {
		return errors.New("loan proposal not found")
	}

	// Perform credit check
	creditScore, err := pm.CreditChecker.RunCreditCheck(proposal.ApplicantID, proposal.WalletAddress)
	if err != nil {
		return fmt.Errorf("credit check failed: %v", err)
	}

	// Update proposal with the credit score
	proposal.CreditScore = creditScore
	proposal.ApprovalStage = "CreditCheck"
	proposal.LastUpdated = time.Now()

	// Store the updated proposal in the ledger
	err = pm.Ledger.UpdateProposal(proposal)
	if err != nil {
		return fmt.Errorf("failed to update proposal after credit check: %v", err)
	}

	fmt.Printf("Credit check completed for Loan ID: %s. Credit Score: %.2f.\n", loanID, creditScore)
	return nil
}

// RunAffordabilityCheck processes the affordability check for the applicant.
func (pm *common.ProposalManager) RunAffordabilityCheck(loanID string, income, expenses, dependentCosts, otherDebts float64, dependents int, workingStatus string, workProof []byte) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve the proposal
	proposal, exists := pm.Proposals[loanID]
	if !exists {
		return errors.New("loan proposal not found")
	}

	// Run affordability check
	affordabilityCheck, err := pm.AffordabilityMgr.SubmitAffordabilityCheck(loanID, proposal.WalletAddress, income, expenses, dependentCosts, otherDebts, dependents, workingStatus, workProof)
	if err != nil {
		return fmt.Errorf("affordability check failed: %v", err)
	}

	// Update proposal with the affordability status
	proposal.AffordabilityStatus = affordabilityCheck.ApprovalStatus
	proposal.ApprovalStage = "Affordability"
	proposal.LastUpdated = time.Now()

	// Store the updated proposal in the ledger
	err = pm.Ledger.UpdateProposal(proposal)
	if err != nil {
		return fmt.Errorf("failed to update proposal after affordability check: %v", err)
	}

	fmt.Printf("Affordability check completed for Loan ID: %s. Status: %s.\n", loanID, affordabilityCheck.ApprovalStatus)
	return nil
}

// CustomizeTerms allows the applicant to customize their loan terms.
func (pm *common.ProposalManager) CustomizeTerms(loanID string, repaymentLength int, loanAmount float64, interestRate float64, islamicOption bool) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve the proposal
	proposal, exists := pm.Proposals[loanID]
	if !exists {
		return errors.New("loan proposal not found")
	}

	// Customize terms
	err := pm.TermsManager.CustomizeLoanTerms(loanID, repaymentLength, loanAmount, interestRate, islamicOption)
	if err != nil {
		return fmt.Errorf("customization of terms failed: %v", err)
	}

	// Update proposal with terms customization
	proposal.TermsCustomization = true
	proposal.ApprovalStage = "Terms"
	proposal.LastUpdated = time.Now()

	// Store the updated proposal in the ledger
	err = pm.Ledger.UpdateProposal(proposal)
	if err != nil {
		return fmt.Errorf("failed to update proposal after terms customization: %v", err)
	}

	fmt.Printf("Loan terms customized for Loan ID: %s.\n", loanID)
	return nil
}

// FinalizeProposal finalizes the proposal and marks it as complete.
func (pm *common.ProposalManager) FinalizeProposal(loanID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Retrieve the proposal
	proposal, exists := pm.Proposals[loanID]
	if !exists {
		return errors.New("loan proposal not found")
	}

	// Ensure all stages are complete
	if proposal.CreditScore == 0 || proposal.AffordabilityStatus == "" || !proposal.TermsCustomization {
		return errors.New("proposal process is incomplete")
	}

	// Mark proposal as approved and finalize
	proposal.ProposalStatus = "Approved"
	proposal.LastUpdated = time.Now()

	// Store the finalized proposal in the ledger
	err := pm.Ledger.FinalizeProposal(proposal)
	if err != nil {
		return fmt.Errorf("failed to finalize proposal: %v", err)
	}

	fmt.Printf("Loan proposal (Loan ID: %s) finalized and approved.\n", loanID)
	return nil
}
