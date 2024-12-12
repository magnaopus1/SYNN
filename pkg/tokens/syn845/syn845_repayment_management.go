package syn845

import (
    "fmt"
    "sync"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
)

// SYN845Token defines the structure for debt repayment management in SYN845 tokens.
type SYN845Token struct {
    mutex                   sync.Mutex
    DebtID                  string
    LoanerID                string
    LoaneeID                string
    PrincipalAmount         float64
    InterestRate            float64
    RepaymentSchedule       []RepaymentEntry
    PaymentHistory          []PaymentEntry
    GracePeriodEnabled      bool
    PenaltyRate             float64
    EarlyRepaymentPenalty   float64
    AssetMetadata           AssetMetadata
    Ledger                  *ledger.Ledger
    Consensus               *consensus.SynnergyConsensus
    EncryptionService       *encryption.Encryption
}

// CREATE_REPAYMENT_SCHEDULE creates a structured repayment schedule for the debt instrument.
func (token *SYN845Token) CREATE_REPAYMENT_SCHEDULE(schedule []RepaymentEntry) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.RepaymentSchedule = schedule
    return token.Ledger.RecordLog("RepaymentScheduleCreated", fmt.Sprintf("Repayment schedule created for debt %s", token.DebtID))
}

// LOG_REPAYMENT_ENTRY logs an individual repayment made toward the debt.
func (token *SYN845Token) LOG_REPAYMENT_ENTRY(payment PaymentEntry) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.PaymentHistory = append(token.PaymentHistory, payment)
    return token.Ledger.RecordLog("RepaymentEntryLogged", fmt.Sprintf("Repayment entry logged for debt %s", token.DebtID))
}

// SET_PENALTY_RATE sets the penalty rate for missed or late payments.
func (token *SYN845Token) SET_PENALTY_RATE(rate float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.PenaltyRate = rate
    return token.Ledger.RecordLog("PenaltyRateSet", fmt.Sprintf("Penalty rate set to %.2f%% for debt %s", rate, token.DebtID))
}

// CHECK_DEFAULT_STATUS checks if the debt has entered default status based on missed payments.
func (token *SYN845Token) CHECK_DEFAULT_STATUS() (bool, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    for _, entry := range token.RepaymentSchedule {
        if !entry.Paid && time.Now().After(entry.DueDate) {
            return true, nil
        }
    }
    return false, nil
}

// DEFAULT_DEBT_INSTRUMENT flags the debt as defaulted and logs the event.
func (token *SYN845Token) DEFAULT_DEBT_INSTRUMENT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Status = Defaulted
    return token.Ledger.RecordLog("DebtDefaulted", fmt.Sprintf("Debt %s has defaulted", token.DebtID))
}

// SET_EARLY_REPAYMENT_PENALTY sets the penalty for early repayment of the debt.
func (token *SYN845Token) SET_EARLY_REPAYMENT_PENALTY(penalty float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.EarlyRepaymentPenalty = penalty
    return token.Ledger.RecordLog("EarlyRepaymentPenaltySet", fmt.Sprintf("Early repayment penalty set for debt %s", token.DebtID))
}

// ENABLE_GRACE_PERIOD enables a grace period for debt repayments.
func (token *SYN845Token) ENABLE_GRACE_PERIOD() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.GracePeriodEnabled = true
    return token.Ledger.RecordLog("GracePeriodEnabled", fmt.Sprintf("Grace period enabled for debt %s", token.DebtID))
}

// DISABLE_GRACE_PERIOD disables the grace period for debt repayments.
func (token *SYN845Token) DISABLE_GRACE_PERIOD() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.GracePeriodEnabled = false
    return token.Ledger.RecordLog("GracePeriodDisabled", fmt.Sprintf("Grace period disabled for debt %s", token.DebtID))
}

// NOTIFY_BORROWER_OF_PAYMENT_DUE sends a notification to the borrower about an upcoming payment.
func (token *SYN845Token) NOTIFY_BORROWER_OF_PAYMENT_DUE(dueDate time.Time) error {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    notification := fmt.Sprintf("Payment due on %v for debt %s", dueDate, token.DebtID)
    return token.Ledger.RecordNotification("PaymentDueNotification", notification)
}

// INITIATE_DEFAULT_RECOVERY_PROCEDURE initiates the recovery process for a defaulted debt.
func (token *SYN845Token) INITIATE_DEFAULT_RECOVERY_PROCEDURE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.RecordLog("DefaultRecoveryInitiated", fmt.Sprintf("Default recovery procedure initiated for debt %s", token.DebtID))
}

// LOG_SETTLEMENT_ENTRY logs a settlement event for the debt.
func (token *SYN845Token) LOG_SETTLEMENT_ENTRY(settlement SettlementEntry) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.SettlementHistory = append(token.SettlementHistory, settlement)
    return token.Ledger.RecordLog("SettlementEntryLogged", fmt.Sprintf("Settlement entry logged for debt %s", token.DebtID))
}

// GET_SETTLEMENT_HISTORY retrieves the settlement history of the debt.
func (token *SYN845Token) GET_SETTLEMENT_HISTORY() ([]SettlementEntry, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.SettlementHistory, nil
}

// CHECK_ASSET_VALUATION checks the current valuation of the collateral asset associated with the debt.
func (token *SYN845Token) CHECK_ASSET_VALUATION() (float64, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.AssetValuation.ValuationAmount, nil
}

// GET_ASSET_METADATA retrieves metadata about the collateral asset.
func (token *SYN845Token) GET_ASSET_METADATA() (AssetMetadata, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.AssetMetadata, nil
}

// LOG_ASSET_VALUATION logs a new valuation of the asset collateral.
func (token *SYN845Token) LOG_ASSET_VALUATION(valuation AssetValuation) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AssetValuation = valuation
    return token.Ledger.RecordLog("AssetValuationLogged", fmt.Sprintf("New asset valuation logged for debt %s", token.DebtID))
}

// Supporting Structures and Enums

type DebtStatus string

const (
    Active    DebtStatus = "active"
    Defaulted DebtStatus = "defaulted"
    Repaid    DebtStatus = "repaid"
)

// RepaymentEntry represents a scheduled payment for the debt.
type RepaymentEntry struct {
    DueDate  time.Time `json:"due_date"`
    Amount   float64   `json:"amount"`
    Paid     bool      `json:"paid"`
}

// PaymentEntry represents an actual payment made towards the debt.
type PaymentEntry struct {
    PaymentDate time.Time `json:"payment_date"`
    Amount      float64   `json:"amount"`
    Interest    float64   `json:"interest"`
    Principal   float64   `json:"principal"`
    Balance     float64   `json:"balance"`
}

// SettlementEntry represents a settlement event for the debt.
type SettlementEntry struct {
    SettlementDate time.Time `json:"settlement_date"`
    SettledBy      string    `json:"settled_by"` // LoanerID or LoaneeID
    SettledAmount  float64   `json:"settled_amount"`
}

// AssetMetadata represents metadata related to the collateral or other assets linked to a debt instrument.
type AssetMetadata struct {
    AssetID         string    `json:"asset_id"`
    OwnerID         string    `json:"owner_id"`
    Description     string    `json:"description"`
    CreationDate    time.Time `json:"creation_date"`
    LastUpdatedDate time.Time `json:"last_updated_date"`
    Value           float64   `json:"value"`
}

// AssetValuation represents the valuation information for the asset linked to the debt instrument.
type AssetValuation struct {
    ValuationID     string    `json:"valuation_id"`
    AssetID         string    `json:"asset_id"`
    ValuationAmount float64   `json:"valuation_amount"`
    ValuationDate   time.Time `json:"valuation_date"`
}
