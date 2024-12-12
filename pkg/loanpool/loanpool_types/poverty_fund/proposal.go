package loanpool

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewProposalManager initializes a new manager for handling poverty fund proposals.
func NewProposalManager(ledgerInstance *ledger.Ledger, syn900Validator *common.Syn900Validator, encryptionService *common.Encryption) *PovertyFundProposalManager {
	return &PovertyFundProposalManager{
		Ledger:          ledgerInstance,
		Proposals:       make(map[string]*PovertyFundProposal),
		Syn900Validator: syn900Validator,
		Encryption:      encryptionService,
	}
}

// SubmitPovertyFundProposal allows an individual to submit a poverty fund proposal with evidence attachments.
func (pm *PovertyFundProposalManager) SubmitPovertyFundProposal(applicantName, contact, incomeDetails, bankBalanceDetails, statementOfReason, benefitStatus, walletAddress string, incomeEvidence, bankBalanceEvidence []byte, amount float64) (*PovertyFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if a proposal already exists for this applicant
	if _, exists := pm.Proposals[applicantName]; exists {
		return nil, errors.New("a proposal for this applicant has already been submitted")
	}

	// Verify if the wallet is verified by syn900
	verified, err := pm.Syn900Validator.IsWalletVerified(walletAddress)
	if err != nil || !verified {
		return nil, errors.New("wallet is not verified by syn900")
	}

	// Encrypt evidence attachments
	encryptedIncomeEvidence, err := pm.Encryption.EncryptData(incomeEvidence, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting income evidence: %v", err)
	}

	encryptedBankBalanceEvidence, err := pm.Encryption.EncryptData(bankBalanceEvidence, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting bank balance evidence: %v", err)
	}

	// Create the poverty fund proposal
	proposal := &PovertyFundProposal{
		ApplicantName:       applicantName,
		ApplicantContact:    contact,
		IncomeDetails:       incomeDetails,
		BankBalanceDetails:  bankBalanceDetails,
		IncomeEvidence:      encryptedIncomeEvidence,
		BankBalanceEvidence: encryptedBankBalanceEvidence,
		StatementOfReason:   statementOfReason,
		BenefitStatus:       benefitStatus,
		WalletAddress:       walletAddress,
		AmountAppliedFor:    amount,
		SubmissionTimestamp: time.Now(),
		VerifiedBySyn900:    true,
		Status:              "Pending",
		LastUpdated:         time.Now(),
		Comments:            []common.ProposalComment{},
	}

	// Encrypt proposal data before storing in the ledger
	encryptedProposalData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting proposal data: %v", err)
	}

	// Store the proposal in the map and record it in the ledger
	pm.Proposals[applicantName] = proposal
	err = pm.Ledger.RecordProposal(encryptedProposalData)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal in the ledger: %v", err)
	}

	fmt.Printf("Poverty fund proposal for applicant %s successfully submitted and recorded.\n", applicantName)
	return proposal, nil
}

// ViewProposal allows users to view a proposal by the applicant's name.
func (pm *PovertyFundProposalManager) ViewProposal(applicantName string) (*PovertyFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found for this applicant")
	}

	return proposal, nil
}

// UpdateProposal allows editing or updating an existing proposal.
func (pm *PovertyFundProposalManager) UpdateProposal(applicantName, contact, incomeDetails, bankBalanceDetails, statementOfReason, benefitStatus string, incomeEvidence, bankBalanceEvidence []byte, amount float64) (*PovertyFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	// Encrypt updated evidence attachments
	encryptedIncomeEvidence, err := pm.Encryption.EncryptData(incomeEvidence, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated income evidence: %v", err)
	}

	encryptedBankBalanceEvidence, err := pm.Encryption.EncryptData(bankBalanceEvidence, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated bank balance evidence: %v", err)
	}

	// Update the proposal details
	proposal.ApplicantContact = contact
	proposal.IncomeDetails = incomeDetails
	proposal.BankBalanceDetails = bankBalanceDetails
	proposal.IncomeEvidence = encryptedIncomeEvidence
	proposal.BankBalanceEvidence = encryptedBankBalanceEvidence
	proposal.StatementOfReason = statementOfReason
	proposal.BenefitStatus = benefitStatus
	proposal.AmountAppliedFor = amount
	proposal.LastUpdated = time.Now()

	// Encrypt updated proposal data
	encryptedProposalData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated proposal data: %v", err)
	}

	// Update the proposal in the ledger
	err = pm.Ledger.UpdateProposal(applicantName, encryptedProposalData)
	if err != nil {
		return nil, fmt.Errorf("failed to update proposal in the ledger: %v", err)
	}

	fmt.Printf("Proposal for applicant %s successfully updated.\n", applicantName)
	return proposal, nil
}

// AddComment allows users to add comments to a proposal.
func (pm *PovertyFundProposalManager) AddComment(applicantName, commenter, comment string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return errors.New("proposal not found")
	}

	// Create the new comment
	newComment := ProposalComment{
		CommentID: GenerateUniqueID(),
		Commenter: commenter,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// Add the comment to the proposal
	proposal.Comments = append(proposal.Comments, newComment)
	proposal.LastUpdated = time.Now()

	// Encrypt updated proposal data
	encryptedProposalData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting proposal comment data: %v", err)
	}

	// Update the proposal in the ledger with the new comment
	err = pm.Ledger.UpdateProposal(applicantName, encryptedProposalData)
	if err != nil {
		return fmt.Errorf("failed to update proposal with comment in the ledger: %v", err)
	}

	fmt.Printf("Comment added to proposal for applicant %s successfully.\n", applicantName)
	return nil
}

// GetProposalComments retrieves all comments associated with a proposal.
func (pm *PovertyFundProposalManager) GetProposalComments(applicantName string) ([]ProposalComment, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	return proposal.Comments, nil
}
