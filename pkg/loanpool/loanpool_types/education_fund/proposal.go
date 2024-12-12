package loanpool

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)




// NewProposalManager initializes a new manager for handling education fund proposals.
func NewProposalManager(ledgerInstance *ledger.Ledger, syn900Validator *syn900.Validator, encryptionService *encryption.Encryption) *ProposalManager {
	return &ProposalManager{
		Ledger:          ledgerInstance,
		Proposals:       make(map[string]*EducationFundProposal),
		Syn900Validator: syn900Validator,
		Encryption:      encryptionService,
	}
}

// SubmitEducationFundProposal allows an individual to submit an education fund proposal.
func (pm *ProposalManager) SubmitEducationFundProposal(applicantName, contact, walletAddress, institutionName, courseName, courseLevel, applicationEvidence, personalStatement, sponsorName, sponsorContact string, amount float64) (*EducationFundProposal, error) {
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

	// Create the education fund proposal
	proposal := &EducationFundProposal{
		ApplicantName:       applicantName,
		ApplicantContact:    contact,
		WalletAddress:       walletAddress,
		InstitutionName:     institutionName,
		CourseName:          courseName,
		CourseLevel:         courseLevel,
		ApplicationEvidence: applicationEvidence,
		PersonalStatement:   personalStatement,
		SponsorName:         sponsorName,
		SponsorContactInfo:  sponsorContact,
		AmountAppliedFor:    amount,
		SubmissionTimestamp: time.Now(),
		VerifiedBySyn900:    true,
		Status:              "Pending",
		LastUpdated:         time.Now(),
		Comments:            []ProposalComment{},
	}

	// Encrypt proposal data
	encryptedData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), "encryption-key")
	if err != nil {
		return nil, fmt.Errorf("error encrypting proposal data: %v", err)
	}

	// Store the proposal in the map and record it in the ledger
	pm.Proposals[applicantName] = proposal
	err = pm.Ledger.RecordEducationFundProposal(proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal in the ledger: %v", err)
	}

	fmt.Printf("Education fund proposal for applicant %s successfully submitted and recorded.\n", applicantName)
	return proposal, nil
}

// ViewProposal allows users to view a proposal by the applicant's name.
func (pm *ProposalManager) ViewProposal(applicantName string) (*EducationFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found for this applicant")
	}

	return proposal, nil
}

// UpdateProposal allows editing or updating an existing proposal.
func (pm *ProposalManager) UpdateProposal(applicantName, contact, institutionName, courseName, courseLevel, applicationEvidence, personalStatement, sponsorName, sponsorContact string, amount float64) (*EducationFundProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	// Update the proposal details
	proposal.ApplicantContact = contact
	proposal.InstitutionName = institutionName
	proposal.CourseName = courseName
	proposal.CourseLevel = courseLevel
	proposal.ApplicationEvidence = applicationEvidence
	proposal.PersonalStatement = personalStatement
	proposal.SponsorName = sponsorName
	proposal.SponsorContactInfo = sponsorContact
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
func (pm *ProposalManager) AddComment(applicantName, commenter, comment string) error {
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
func (pm *ProposalManager) GetProposalComments(applicantName string) ([]ProposalComment, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[applicantName]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	return proposal.Comments, nil
}
