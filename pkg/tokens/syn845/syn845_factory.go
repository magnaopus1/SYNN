package syn845

import (
	"errors"
	"sync"
	"time"

)

// DebtStatus represents the status of the debt instrument
type DebtStatus string

const (
	Active    DebtStatus = "active"
	Defaulted DebtStatus = "defaulted"
	Repaid    DebtStatus = "repaid"
)

// SYN845Token represents the SYN845 debt instrument token standard
type SYN845Token struct {
	DebtID               string           `json:"debt_id"`
	LoanerID             string           `json:"loaner_id"`
	LoaneeID             string           `json:"loanee_id"`
	PrincipalAmount      float64          `json:"principal_amount"`
	InterestRate         float64          `json:"interest_rate"`
	RepaymentPeriod      int              `json:"repayment_period"` // in months
	PenaltyRate          float64          `json:"penalty_rate"`
	CollateralID         string           `json:"collateral_id"`
	Status               DebtStatus       `json:"status"`
	CreationDate         time.Time        `json:"creation_date"`
	LastUpdatedDate      time.Time        `json:"last_updated_date"`
	AccruedInterest      float64          `json:"accrued_interest"`
	RepaymentSchedule    []RepaymentEntry `json:"repayment_schedule"`
	PaymentHistory       []PaymentEntry   `json:"payment_history"`
	EarlyRepaymentPenalty float64         `json:"early_repayment_penalty"`
	SettlementHistory    []SettlementEntry `json:"settlement_history"`
	AssetMetadata        AssetMetadata    `json:"asset_metadata"`   // Linked metadata
	AssetValuation       AssetValuation   `json:"asset_valuation"`  // Linked valuation
}

// RepaymentEntry represents an entry in the repayment schedule
type RepaymentEntry struct {
	DueDate  time.Time `json:"due_date"`
	Amount   float64   `json:"amount"`
	Paid     bool      `json:"paid"`
}

// PaymentEntry represents a payment made towards the debt instrument
type PaymentEntry struct {
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	Interest    float64   `json:"interest"`
	Principal   float64   `json:"principal"`
	Balance     float64   `json:"balance"`
}

// SettlementEntry represents an entry in the settlement history
type SettlementEntry struct {
	SettlementDate time.Time `json:"settlement_date"`
	SettledBy      string    `json:"settled_by"` // LoanerID or LoaneeID
	SettledAmount  float64   `json:"settled_amount"`
}

// AssetMetadata represents metadata related to the collateral or other assets linked to a debt instrument
type AssetMetadata struct {
	AssetID         string    `json:"asset_id"`
	OwnerID         string    `json:"owner_id"`
	Description     string    `json:"description"`
	CreationDate    time.Time `json:"creation_date"`
	LastUpdatedDate time.Time `json:"last_updated_date"`
	Value           float64   `json:"value"`
}

// AssetValuation represents the valuation information for the asset linked to the debt instrument
type AssetValuation struct {
	ValuationID     string    `json:"valuation_id"`
	AssetID         string    `json:"asset_id"`
	ValuationAmount float64   `json:"valuation_amount"`
	ValuationDate   time.Time `json:"valuation_date"`
}

var (
	syn845Store         = make(map[string]SYN845Token)
	mutex               = &sync.Mutex{}
)

// CreateSYN845Token creates a new SYN845Token debt instrument and validates it with consensus
func CreateSYN845Token(loanerID, loaneeID string, principalAmount, interestRate, penaltyRate, earlyRepaymentPenalty float64, repaymentPeriod int, collateralID string, metadata AssetMetadata, valuation AssetValuation) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	debtID := generateDebtID()
	creationDate := time.Now()
	repaymentSchedule := generateRepaymentSchedule(principalAmount, interestRate, repaymentPeriod)

	syn845Token := SYN845Token{
		DebtID:               debtID,
		LoanerID:             loanerID,
		LoaneeID:             loaneeID,
		PrincipalAmount:      principalAmount,
		InterestRate:         interestRate,
		RepaymentPeriod:      repaymentPeriod,
		PenaltyRate:          penaltyRate,
		CollateralID:         collateralID,
		Status:               Active,
		CreationDate:         creationDate,
		LastUpdatedDate:      creationDate,
		AccruedInterest:      0,
		RepaymentSchedule:    repaymentSchedule,
		PaymentHistory:       []PaymentEntry{},
		EarlyRepaymentPenalty: earlyRepaymentPenalty,
		SettlementHistory:    []SettlementEntry{},
		AssetMetadata:        metadata,
		AssetValuation:       valuation,
	}

	// Validate debt creation with Synnergy Consensus
	if err := consensus.ValidateDebtCreation(syn845Token); err != nil {
		return "", errors.New("consensus validation failed for debt creation")
	}

	// Record the transaction in the ledger
	if err := ledger.RecordEntry(debtID, "debt_creation", principalAmount, 0, 0, 0, ""); err != nil {
		return "", err
	}

	syn845Store[debtID] = syn845Token

	// Store the encrypted data
	if err := saveSYN845TokenToStorage(syn845Token); err != nil {
		return "", err
	}

	return debtID, nil
}

// UpdateSYN845Token updates an existing SYN845Token debt instrument and validates with consensus
func UpdateSYN845Token(debtID string, principalAmount, interestRate, penaltyRate, earlyRepaymentPenalty float64, repaymentPeriod int, collateralID string, metadata AssetMetadata, valuation AssetValuation, status DebtStatus) error {
	mutex.Lock()
	defer mutex.Unlock()

	syn845Token, exists := syn845Store[debtID]
	if !exists {
		return errors.New("debt instrument not found")
	}

	syn845Token.PrincipalAmount = principalAmount
	syn845Token.InterestRate = interestRate
	syn845Token.RepaymentPeriod = repaymentPeriod
	syn845Token.PenaltyRate = penaltyRate
	syn845Token.CollateralID = collateralID
	syn845Token.AssetMetadata = metadata
	syn845Token.AssetValuation = valuation
	syn845Token.Status = status
	syn845Token.LastUpdatedDate = time.Now()

	// Validate the update with Synnergy Consensus
	if err := consensus.ValidateDebtUpdate(syn845Token); err != nil {
		return errors.New("consensus validation failed for debt update")
	}

	// Record the update in the ledger
	if err := ledger.RecordEntry(debtID, "debt_update", principalAmount, 0, 0, 0, ""); err != nil {
		return err
	}

	syn845Store[debtID] = syn845Token
	return saveSYN845TokenToStorage(syn845Token)
}

// AddPayment records a payment for a SYN845Token debt instrument and validates the payment with consensus
func AddPayment(debtID string, amount, interest, principal float64) error {
	mutex.Lock()
	defer mutex.Unlock()

	syn845Token, exists := syn845Store[debtID]
	if !exists {
		return errors.New("debt instrument not found")
	}

	paymentDate := time.Now()
	balance := syn845Token.PrincipalAmount - principal
	syn845Token.AccruedInterest += interest

	paymentEntry := PaymentEntry{
		PaymentDate: paymentDate,
		Amount:      amount,
		Interest:    interest,
		Principal:   principal,
		Balance:     balance,
	}

	syn845Token.PaymentHistory = append(syn845Token.PaymentHistory, paymentEntry)
	syn845Token.LastUpdatedDate = paymentDate

	// Validate payment with Synnergy Consensus
	if err := consensus.ValidateDebtPayment(syn845Token); err != nil {
		return errors.New("consensus validation failed for payment")
	}

	// Record the payment in the ledger
	if err := ledger.RecordEntry(debtID, "payment", amount, balance, interest, principal, ""); err != nil {
		return err
	}

	for i, entry := range syn845Token.RepaymentSchedule {
		if !entry.Paid && entry.DueDate.Before(paymentDate) {
			syn845Token.RepaymentSchedule[i].Paid = true
		}
	}

	syn845Store[debtID] = syn845Token
	return saveSYN845TokenToStorage(syn845Token)
}

// SettleDebt settles a SYN845Token debt instrument, validates with consensus, and updates the ledger
func SettleDebt(debtID, settledBy string, settledAmount float64) error {
	mutex.Lock()
	defer mutex.Unlock()

	syn845Token, exists := syn845Store[debtID]
	if !exists {
		return errors.New("debt instrument not found")
	}

	settlementDate := time.Now()

	settlementEntry := SettlementEntry{
		SettlementDate: settlementDate,
		SettledBy:      settledBy,
		SettledAmount:  settledAmount,
	}

	syn845Token.SettlementHistory = append(syn845Token.SettlementHistory, settlementEntry)
	syn845Token.LastUpdatedDate = settlementDate

	// Check if the debt is fully settled
	if settledAmount >= syn845Token.PrincipalAmount+syn845Token.AccruedInterest {
		syn845Token.Status = Repaid
	}

	// Validate settlement with Synnergy Consensus
	if err := consensus.ValidateDebtSettlement(syn845Token); err != nil {
		return errors.New("consensus validation failed for settlement")
	}

	// Record the settlement in the ledger
	if err := ledger.RecordEntry(debtID, "debt_settlement", settledAmount, 0, 0, 0, ""); err != nil {
		return err
	}

	syn845Store[debtID] = syn845Token
	return saveSYN845TokenToStorage(syn845Token)
}

// Ledger and Storage Integration for SYN845 Token
func saveSYN845TokenToStorage(syn845Token SYN845Token) error {
	data, err := json.Marshal(syn845Token)
	if err != nil {
		return err
	}

	encryptedData, err := encryption.Encrypt(data)
	if err != nil {
		return err
	}

	return storage.Save("syn845", syn845Token.DebtID, encryptedData)
}

func deleteSYN845TokenFromStorage(debtID string) error {
	return storage.Delete("syn845", debtID)
}

// Utility functions
func generateDebtID() string {
	return "debt-" + uuid.New().String()
}

func generateRepaymentSchedule(principalAmount, interestRate float64, repaymentPeriod int) []RepaymentEntry {
	var schedule []RepaymentEntry
	for i := 0; i < repaymentPeriod; i++ {
		schedule = append(schedule, RepaymentEntry{
			DueDate: time.Now().AddDate(0, i+1, 0),
			Amount:  principalAmount / float64(repaymentPeriod),
			Paid:    false,
		})
	}
	return schedule
}
