package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewProposalManager initializes a new manager for handling healthcare support fund proposals.
func NewProposalManager(ledgerInstance *ledger.Ledger, syn900Validator *common.Syn900Validator, encryptionService *common.Encryption) *HealthcareSupportFundProposalManager {
	return &HealthcareSupportFundProposalManager{
		Ledger:          ledgerInstance,
		Proposals:       make(map[string]*HealthcareSupportFundProposal),
		Syn900Validator: syn900Validator,
		Encryption:      encryptionService,
	}
}

// SubmitHealthcareSupportFundProposal allows an individual to submit a healthcare support fund proposal.
func (pm *HealthcareSupportFundProposalManager) SubmitHealthcareSupportFundProposal(applicantName, applicantContact, medicalProfessionalName, medicalProfessionalContact, walletAddress, hospitalName, medicalProcedure, costBreakdownEvidence, hospitalAddress, hospitalContactInfo string, amount float64) (*HealthcareSupportFundProposal, error) {
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

	// Create the healthcare support fund proposal
	proposal := &HealthcareSupportFundProposal{
		ApplicantName:              applicantName,
		ApplicantContact:           applicantContact,
		MedicalProfessionalName:    medicalProfessionalName,
		MedicalProfessionalContact: medicalProfessionalContact,
		WalletAddress:              walletAddress,
		HospitalName:               hospitalName,
		MedicalProcedure:           medicalProcedure,
		CostBreakdownEvidence:      costBreakdownEvidence,
		HospitalAddress:            hospitalAddress,
		HospitalContactInfo:        hospitalContactInfo,
		AmountAppliedFor:           amount,
		SubmissionTimestamp:        time.Now(),
		VerifiedBySyn900:           true,
		Status:                     "Pending",
		LastUpdated:                time.Now(),
		Comments:                   []ProposalComment{},
	}

	// Encrypt proposal data
	encryptedData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), "encryption-key")
	if err != nil {
		return nil, fmt.Errorf("error encrypting proposal data: %v", err)
	}

	// Store the proposal in the map and record it in the ledger
	pm.Proposals[applicantName] = proposal
	err = pm.Ledger.RecordProposal(proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal in the ledger: %v", err)
	}

	fmt.Printf("Healthcare support fund proposal for applicant %s successfully submitted and recorded.\n", applicantName)
	return proposal, nil
}

// ViewProposal allows users to view a proposal by the applicant's name.
func (pm *HealthcareSupportFundProposalManager) ViewProposal(applicantName string) (*HealthcareSupportFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found for this applicant")
	}

	return proposal, nil
}

// UpdateProposal allows editing or updating an existing proposal.
func (pm *HealthcareSupportFundProposalManager) UpdateProposal(applicantName, applicantContact, medicalProfessionalName, medicalProfessionalContact, hospitalName, medicalProcedure, costBreakdownEvidence, hospitalAddress, hospitalContactInfo string, amount float64) (*HealthcareSupportFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	// Update the proposal details
	proposal.ApplicantContact = applicantContact
	proposal.MedicalProfessionalName = medicalProfessionalName
	proposal.MedicalProfessionalContact = medicalProfessionalContact
	proposal.HospitalName = hospitalName
	proposal.MedicalProcedure = medicalProcedure
	proposal.CostBreakdownEvidence = costBreakdownEvidence
	proposal.HospitalAddress = hospitalAddress
	proposal.HospitalContactInfo = hospitalContactInfo
	proposal.AmountAppliedFor = amount
	proposal.LastUpdated = time.Now()

	// Encrypt updated proposal data
	encryptedData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), "encryption-key")
	if err != nil {
		return nil, fmt.Errorf("error encrypting updated proposal data: %v", err)
	}

	// Update the proposal in the ledger
	err = pm.Ledger.UpdateProposal(applicantName, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to update proposal in ledger: %v", err)
	}

	fmt.Printf("Proposal for applicant %s successfully updated.\n", applicantName)
	return proposal, nil
}

// AddComment allows users to add comments to a proposal.
func (pm *HealthcareSupportFundProposalManager) AddComment(applicantName, commenter, comment string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return errors.New("proposal not found")
	}

	// Create the new comment
	newComment := ProposalComment{
		CommentID: common.GenerateUniqueID(),
		Commenter: commenter,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	// Add the comment to the proposal
	proposal.Comments = append(proposal.Comments, newComment)
	proposal.LastUpdated = time.Now()

	// Encrypt updated proposal data
	encryptedData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), "encryption-key")
	if err != nil {
		return fmt.Errorf("error encrypting proposal comment data: %v", err)
	}

	// Update the proposal in the ledger with the new comment
	err = pm.Ledger.UpdateProposal(applicantName, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to update proposal with comment in ledger: %v", err)
	}

	fmt.Printf("Comment added to proposal for applicant %s successfully.\n", applicantName)
	return nil
}

// GetProposalComments retrieves all comments associated with a proposal.
func (pm *HealthcareSupportFundProposalManager) GetProposalComments(applicantName string) ([]ProposalComment, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	return proposal.Comments, nil
}
