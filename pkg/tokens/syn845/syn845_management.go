package syn845

import (
	"errors"
	"sync"
	"time"

)

// NodeType represents different types of nodes in the network
type NodeType string

const (
	GovernmentNode  NodeType = "government"
	CreditorNode    NodeType = "creditor"
	CentralBankNode NodeType = "central_bank"
	BankingNode     NodeType = "banking"
	OtherNode       NodeType = "other"
)

// DebtIssuanceManager manages the issuance and actions for SYN845 tokens
type DebtIssuanceManager struct {
	mu       sync.Mutex
	nodeType NodeType
}

// NewDebtIssuanceManager creates a new instance of DebtIssuanceManager
func NewDebtIssuanceManager(nodeType NodeType) *DebtIssuanceManager {
	return &DebtIssuanceManager{
		nodeType: nodeType,
	}
}

// IssueDebt creates a new SYN845Token debt instrument, validates it, and records it in the ledger
func (dim *DebtIssuanceManager) IssueDebt(ownerID string, principalAmount, interestRate, penaltyRate, earlyRepaymentPenalty float64, repaymentPeriod int, collateralID string, metadata syn845.AssetMetadata, valuation syn845.AssetValuation) (string, error) {
	dim.mu.Lock()
	defer dim.mu.Unlock()

	if !dim.isAuthorized() {
		return "", errors.New("unauthorized node type for issuing debt")
	}

	debtID, err := syn845.CreateSYN845Token(ownerID, principalAmount, interestRate, penaltyRate, earlyRepaymentPenalty, repaymentPeriod, collateralID, metadata, valuation)
	if err != nil {
		return "", err
	}

	if err := consensus.ValidateDebtIssuance(debtID); err != nil {
		return "", err
	}

	encryptedDebt, err := security.Encrypt(common.StructToJSON(debtID))
	if err != nil {
		return "", err
	}

	if err := ledger.RecordDebtIssuance(debtID, encryptedDebt); err != nil {
		return "", err
	}

	return debtID, nil
}

// EarlySettlement allows a debt instrument to be settled early with applicable penalties
func (dim *DebtIssuanceManager) EarlySettlement(debtID string, earlyPaymentAmount float64) (float64, error) {
	dim.mu.Lock()
	defer dim.mu.Unlock()

	debt, err := syn845.GetSYN845Token(debtID)
	if err != nil {
		return 0, errors.New("debt instrument not found")
	}

	if earlyPaymentAmount < debt.PrincipalAmount {
		return 0, errors.New("early payment is less than the remaining principal")
	}

	penalty := debt.EarlyRepaymentPenalty * earlyPaymentAmount / 100
	totalAmountDue := earlyPaymentAmount + penalty

	debt.Status = syn845.Repaid
	debt.LastUpdatedDate = time.Now()

	if err := consensus.ValidateEarlySettlement(debtID); err != nil {
		return 0, err
	}

	encryptedDebt, err := security.Encrypt(common.StructToJSON(debtID))
	if err != nil {
		return 0, err
	}

	if err := ledger.RecordEarlySettlement(debtID, encryptedDebt); err != nil {
		return 0, err
	}

	return totalAmountDue, nil
}

// ProcessRepayment handles repayments, including interest, and applies late penalties when necessary
func (dim *DebtIssuanceManager) ProcessRepayment(debtID string, paymentAmount, interest, principal float64) error {
	dim.mu.Lock()
	defer dim.mu.Unlock()

	err := syn845.AddPayment(debtID, paymentAmount, interest, principal)
	if err != nil {
		return err
	}

	if err := consensus.ValidateDebtPayment(debtID); err != nil {
		return err
	}

	encryptedDebt, err := security.Encrypt(common.StructToJSON(debtID))
	if err != nil {
		return err
	}

	return ledger.RecordDebtPayment(debtID, encryptedDebt)
}

// ApplyLatePaymentPenalty adds penalties for late payments
func (dim *DebtIssuanceManager) ApplyLatePaymentPenalty(debtID string, penaltyAmount float64) error {
	dim.mu.Lock()
	defer dim.mu.Unlock()

	debt, err := syn845.GetSYN845Token(debtID)
	if err != nil {
		return errors.New("debt instrument not found")
	}

	if debt.Status != syn845.Active {
		return errors.New("penalty can only be applied to active debts")
	}

	debt.PrincipalAmount += penaltyAmount
	debt.LastUpdatedDate = time.Now()

	if err := consensus.ValidateLatePenalty(debtID); err != nil {
		return err
	}

	encryptedDebt, err := security.Encrypt(common.StructToJSON(debtID))
	if err != nil {
		return err
	}

	return ledger.RecordLatePenalty(debtID, encryptedDebt)
}

// ManageCollateral ensures collateral management and updates valuation
func (dim *DebtIssuanceManager) ManageCollateral(debtID string, newCollateralID string, newValuation syn845.AssetValuation) error {
	dim.mu.Lock()
	defer dim.mu.Unlock()

	debt, err := syn845.GetSYN845Token(debtID)
	if err != nil {
		return errors.New("debt instrument not found")
	}

	debt.CollateralID = newCollateralID
	debt.AssetValuation = newValuation
	debt.LastUpdatedDate = time.Now()

	// Re-validate after collateral update
	if err := consensus.ValidateCollateralUpdate(debtID); err != nil {
		return err
	}

	encryptedDebt, err := security.Encrypt(common.StructToJSON(debtID))
	if err != nil {
		return err
	}

	return ledger.RecordCollateralUpdate(debtID, encryptedDebt)
}

// RequestDebtRefinancing allows a user to request refinancing for their debt instrument
func (dim *DebtIssuanceManager) RequestDebtRefinancing(debtID string, newTerms map[string]interface{}) (string, error) {
	dim.mu.Lock()
	defer dim.mu.Unlock()

	proposalID, err := dim.proposeChange(debtID, "Refinancing Request", newTerms)
	if err != nil {
		return "", err
	}

	return proposalID, nil
}

// proposeChange creates a proposal for changes such as refinancing
func (dim *DebtIssuanceManager) proposeChange(debtID, description string, newTerms map[string]interface{}) (string, error) {
	proposal := struct {
		ProposalID   string
		DebtID       string
		Description  string
		NewTerms     map[string]interface{}
		CreationDate time.Time
	}{
		ProposalID:   generateProposalID(),
		DebtID:       debtID,
		Description:  description,
		NewTerms:     newTerms,
		CreationDate: time.Now(),
	}

	encryptedProposal, err := security.Encrypt(common.StructToJSON(proposal))
	if err != nil {
		return "", err
	}

	if err := storage.Save("proposal", proposal.ProposalID, encryptedProposal); err != nil {
		return "", err
	}

	return proposal.ProposalID, nil
}

// isAuthorized checks if the current node type is authorized to perform debt issuance management actions
func (dim *DebtIssuanceManager) isAuthorized() bool {
	switch dim.nodeType {
	case GovernmentNode, CreditorNode, CentralBankNode, BankingNode:
		return true
	default:
		return false
	}
}

// generateProposalID generates a unique ID for proposals
func generateProposalID() string {
	return "proposal-" + time.Now().Format("20060102150405")
}
