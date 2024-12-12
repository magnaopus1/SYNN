package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/common"
)


// NewAffordabilityManager initializes a new manager for handling affordability checks.
func NewAffordabilityManager(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *encryption.Encryption) *common.AffordabilityManager {
	return &common.AffordabilityManager{
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		EncryptionService: encryptionService,
		Submissions:       make(map[string]*common.AffordabilityCheck),
		ApprovalQueue:     []*common.AffordabilityCheck{},
	}
}

// SubmitAffordabilityCheck allows applicants to submit their financial details for affordability checks.
func (am *common.AffordabilityManager) SubmitAffordabilityCheck(loanID, applicantWallet string, income, expenses, dependentCosts, otherDebts float64, dependents int, workingStatus string, workProof []byte) (*common.AffordabilityCheck, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Validate input
	if loanID == "" || applicantWallet == "" || income <= 0 || expenses <= 0 || dependents < 0 || len(workProof) == 0 {
		return nil, errors.New("invalid affordability check details")
	}

	// Create a new affordability check
	affordabilityCheck := &common.AffordabilityCheck{
		LoanID:          loanID,
		ApplicantWallet: applicantWallet,
		Income:          income,
		Expenses:        expenses,
		Dependents:      dependents,
		DependentCosts:  dependentCosts,
		WorkingStatus:   workingStatus,
		OtherDebts:      otherDebts,
		WorkProof:       workProof,
		SubmissionTime:  time.Now(),
		ApprovalStatus:  "Pending",
	}

	// Encrypt the work proof before storing it
	encryptedWorkProof, err := am.EncryptionService.EncryptData(workProof, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt work proof: %v", err)
	}
	affordabilityCheck.WorkProof = encryptedWorkProof

	// Store the affordability check submission in the ledger
	err = am.Ledger.RecordAffordabilityCheck(loanID, affordabilityCheck)
	if err != nil {
		return nil, fmt.Errorf("failed to record affordability check in the ledger: %v", err)
	}

	// Add to the approval queue
	am.Submissions[loanID] = affordabilityCheck
	am.ApprovalQueue = append(am.ApprovalQueue, affordabilityCheck)

	fmt.Printf("Affordability check for loan %s successfully submitted and recorded.\n", loanID)
	return affordabilityCheck, nil
}

// ApproveAffordabilityCheck processes and approves an affordability check through Synnergy Consensus.
func (am *common.AffordabilityManager) ApproveAffordabilityCheck(loanID, approvingNode string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Retrieve the affordability check submission
	affordabilityCheck, exists := am.Submissions[loanID]
	if !exists {
		return errors.New("affordability check not found for this loan ID")
	}

	// Validate the check using Synnergy Consensus
	approved, err := am.Consensus.ValidateAffordability(loanID, affordabilityCheck.Income, affordabilityCheck.Expenses, affordabilityCheck.Dependents, affordabilityCheck.DependentCosts, affordabilityCheck.OtherDebts)
	if err != nil || !approved {
		return fmt.Errorf("failed to approve affordability check for loan %s: %v", loanID, err)
	}

	// Mark the check as approved and record the approval time
	affordabilityCheck.ApprovalStatus = "Approved"
	affordabilityCheck.ApprovedBy = approvingNode
	affordabilityCheck.ApprovalTime = time.Now()

	// Update the affordability check in the ledger
	err = am.Ledger.UpdateAffordabilityApproval(loanID, affordabilityCheck)
	if err != nil {
		return fmt.Errorf("failed to update affordability approval in the ledger: %v", err)
	}

	fmt.Printf("Affordability check for loan %s approved by node %s.\n", loanID, approvingNode)
	return nil
}

// RejectAffordabilityCheck handles the rejection of an affordability check.
func (am *common.AffordabilityManager) RejectAffordabilityCheck(loanID, rejectingNode string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Retrieve the affordability check submission
	affordabilityCheck, exists := am.Submissions[loanID]
	if !exists {
		return errors.New("affordability check not found for this loan ID")
	}

	// Mark the check as rejected and record the rejection time
	affordabilityCheck.ApprovalStatus = "Rejected"
	affordabilityCheck.ApprovedBy = rejectingNode
	affordabilityCheck.ApprovalTime = time.Now()

	// Update the affordability check in the ledger
	err := am.Ledger.UpdateAffordabilityRejection(loanID, affordabilityCheck)
	if err != nil {
		return fmt.Errorf("failed to update affordability rejection in the ledger: %v", err)
	}

	fmt.Printf("Affordability check for loan %s rejected by node %s.\n", loanID, rejectingNode)
	return nil
}

// ViewAffordabilityCheck allows users to view their affordability check and status.
func (am *common.AffordabilityManager) ViewAffordabilityCheck(loanID string) (*common.AffordabilityCheck, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	affordabilityCheck, exists := am.Submissions[loanID]
	if !exists {
		return nil, errors.New("affordability check not found for this loan ID")
	}

	// Decrypt the work proof before returning it
	decryptedWorkProof, err := am.EncryptionService.DecryptData(affordabilityCheck.WorkProof, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt work proof: %v", err)
	}
	affordabilityCheck.WorkProof = decryptedWorkProof

	return affordabilityCheck, nil
}

// ProcessApprovalQueue processes affordability checks that are pending approval.
func (am *common.AffordabilityManager) ProcessApprovalQueue() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	newQueue := []*common.AffordabilityCheck{}

	for _, submission := range am.ApprovalQueue {
		// Process pending submissions based on the business logic for validation and approval
		if submission.ApprovalStatus == "Pending" {
			newQueue = append(newQueue, submission)
		}
	}

	// Update the approval queue
	am.ApprovalQueue = newQueue
}
