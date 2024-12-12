package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn900"
)

// NewSmallBusinessProposalManager initializes a new manager for handling small business grant proposals.
func NewSmallBusinessProposalManager(ledgerInstance *ledger.Ledger, syn900Validator *syn900.Validator, encryptionService *encryption.Encryption) *common.SmallBusinessProposalManager {
	return &common.SmallBusinessProposalManager{
		Ledger:          ledgerInstance,
		Proposals:       make(map[string]*common.SmallBusinessGrantProposal),
		Syn900Validator: syn900Validator,
		Encryption:      encryptionService,
	}
}

// SubmitSmallBusinessGrantProposal allows a small business to submit a grant proposal.
func (pm *common.SmallBusinessProposalManager) SubmitSmallBusinessGrantProposal(businessName, businessAddress, regNumber, country, website, activities, applicantName, walletAddress, usageDesc, financialPosition string, amount float64, startup bool, employeeCount int) (*common.SmallBusinessGrantProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Check if the proposal for this business already exists
	if _, exists := pm.Proposals[businessName]; exists {
		return nil, errors.New("a proposal for this business has already been submitted")
	}

	// Validate the business criteria (must be either a startup or have <= 50 employees)
	if !startup && employeeCount > 50 {
		return nil, errors.New("the business must have 50 or fewer employees or be a startup")
	}

	// Verify if the wallet is verified by syn900
	verified, err := pm.Syn900Validator.IsWalletVerified(walletAddress)
	if err != nil || !verified {
		return nil, errors.New("wallet is not verified by syn900")
	}

	// Create the grant proposal
	proposal := &common.SmallBusinessGrantProposal{
		BusinessName:        businessName,
		BusinessAddress:     businessAddress,
		RegistrationNumber:  regNumber,
		Country:             country,
		Website:             website,
		BusinessActivities:  activities,
		ApplicantName:       applicantName,
		WalletAddress:       walletAddress,
		AmountAppliedFor:    amount,
		UsageDescription:    usageDesc,
		FinancialPosition:   financialPosition,
		SubmissionTimestamp: time.Now(),
		VerifiedBySyn900:    true,
		Status:              "Pending",
		LastUpdated:         time.Now(),
		Startup:             startup,
		EmployeeCount:       employeeCount,
		Comments:            []common.ProposalComment{},
	}

	// Encrypt proposal data
	encryptedData, err := pm.Encryption.EncryptData(fmt.Sprintf("%v", proposal), "encryption-key")
	if err != nil {
		return nil, fmt.Errorf("error encrypting proposal data: %v", err)
	}

	// Store the proposal in the map and record it in the ledger
	pm.Proposals[businessName] = proposal
	err = pm.Ledger.RecordGrantProposal(proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to record proposal in the ledger: %v", err)
	}

	fmt.Printf("Grant proposal for small business %s successfully submitted and recorded.\n", businessName)
	return proposal, nil
}

// ViewProposal allows users to view a proposal by the business name.
func (pm *common.SmallBusinessProposalManager) ViewProposal(businessName string) (*common.SmallBusinessGrantProposal, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[businessName]
	if !exists {
		return nil, errors.New("proposal not found for this business")
	}

	return proposal, nil
}

// AddComment allows users to add comments to a proposal.
func (pm *common.SmallBusinessProposalManager) AddComment(proposalID, commenter, comment string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[proposalID]
	if !exists {
		return errors.New("proposal not found")
	}

	// Create the new comment
	newComment := common.ProposalComment{
		CommentID:   common.GenerateUniqueID(),
		Commenter:   commenter,
		Comment:     comment,
		CreatedAt:   time.Now(),
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
	err = pm.Ledger.UpdateProposal(proposalID, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to update proposal with comment in ledger: %v", err)
	}

	fmt.Printf("Comment added to proposal %s successfully.\n", proposalID)
	return nil
}

// GetProposalComments retrieves all comments associated with a proposal.
func (pm *common.SmallBusinessProposalManager) GetProposalComments(proposalID string) ([]common.ProposalComment, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	proposal, exists := pm.Proposals[proposalID]
	if !exists {
		return nil, errors.New("proposal not found")
	}

	return proposal.Comments, nil
}
