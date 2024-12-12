package syn845

import (
    "fmt"
    "sync"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
    "math"
)

// SYN845Token defines the debt management structure for SYN845 tokens.
type SYN845Token struct {
    mutex               sync.Mutex
    DebtID              string
    LoanerID            string
    LoaneeID            string
    PrincipalAmount     float64
    InterestRate        float64
    CollateralID        string
    Status              DebtStatus
    AccruedInterest     float64
    PaymentHistory      []PaymentEntry
    Ledger              *ledger.Ledger
    Consensus           *consensus.SynnergyConsensus
    EncryptionService   *encryption.Encryption
}

// TRANSFER_SYN845_TOKEN transfers the debt token to a new owner.
func (token *SYN845Token) TRANSFER_SYN845_TOKEN(newOwnerID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.LoaneeID = newOwnerID
    return token.Ledger.RecordLog("TokenTransfer", fmt.Sprintf("Debt token %s transferred to %s", token.DebtID, newOwnerID))
}

// APPROVE_SYN845_TOKEN_TRANSFER approves a token transfer to a specified address.
func (token *SYN845Token) APPROVE_SYN845_TOKEN_TRANSFER(approverID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.RecordLog("TransferApproved", fmt.Sprintf("Transfer of debt token %s approved by %s", token.DebtID, approverID))
}

// CHECK_SYN845_TOKEN_BALANCE returns the balance of the debt associated with the token.
func (token *SYN845Token) CHECK_SYN845_TOKEN_BALANCE() (float64, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    outstandingBalance := token.PrincipalAmount + token.AccruedInterest
    return outstandingBalance, nil
}

// GET_SYN845_TOKEN_METADATA retrieves metadata for the debt token.
func (token *SYN845Token) GET_SYN845_TOKEN_METADATA() (map[string]interface{}, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    metadata := map[string]interface{}{
        "DebtID":          token.DebtID,
        "LoanerID":        token.LoanerID,
        "LoaneeID":        token.LoaneeID,
        "PrincipalAmount": token.PrincipalAmount,
        "InterestRate":    token.InterestRate,
        "CollateralID":    token.CollateralID,
        "Status":          token.Status,
        "AccruedInterest": token.AccruedInterest,
    }
    return metadata, nil
}

// UPDATE_SYN845_TOKEN_METADATA updates metadata for the debt token.
func (token *SYN845Token) UPDATE_SYN845_TOKEN_METADATA(metadata map[string]interface{}) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if principal, ok := metadata["PrincipalAmount"].(float64); ok {
        token.PrincipalAmount = principal
    }
    if interestRate, ok := metadata["InterestRate"].(float64); ok {
        token.InterestRate = interestRate
    }
    if status, ok := metadata["Status"].(DebtStatus); ok {
        token.Status = status
    }
    return token.Ledger.RecordLog("MetadataUpdated", fmt.Sprintf("Metadata updated for debt token %s", token.DebtID))
}

// SET_SYN845_COLLATERAL_ID sets the collateral ID associated with the debt token.
func (token *SYN845Token) SET_SYN845_COLLATERAL_ID(collateralID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.CollateralID = collateralID
    return token.Ledger.RecordLog("CollateralIDSet", fmt.Sprintf("Collateral ID %s set for debt %s", collateralID, token.DebtID))
}

// GET_SYN845_COLLATERAL_ID retrieves the collateral ID associated with the debt token.
func (token *SYN845Token) GET_SYN845_COLLATERAL_ID() (string, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.CollateralID, nil
}

// GET_SYN845_DEBT_STATUS retrieves the current status of the debt (e.g., Active, Defaulted).
func (token *SYN845Token) GET_SYN845_DEBT_STATUS() (DebtStatus, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Status, nil
}

// LOG_DEBT_STATUS_CHANGE logs changes to the debt status.
func (token *SYN845Token) LOG_DEBT_STATUS_CHANGE(newStatus DebtStatus) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Status = newStatus
    return token.Ledger.RecordLog("StatusChanged", fmt.Sprintf("Debt %s status changed to %s", token.DebtID, newStatus))
}

// CALCULATE_ACCRUED_INTEREST calculates the accrued interest on the debt token based on its interest rate and time.
func (token *SYN845Token) CALCULATE_ACCRUED_INTEREST(durationInDays int) (float64, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    dailyRate := token.InterestRate / 365
    accruedInterest := token.PrincipalAmount * dailyRate * float64(durationInDays) / 100
    token.AccruedInterest += accruedInterest
    return accruedInterest, nil
}

// ADD_PAYMENT_ENTRY adds a new payment entry to the debt token's payment history.
func (token *SYN845Token) ADD_PAYMENT_ENTRY(payment PaymentEntry) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.PaymentHistory = append(token.PaymentHistory, payment)
    return token.Ledger.RecordLog("PaymentEntryAdded", fmt.Sprintf("Payment of %.2f added for debt %s", payment.Amount, token.DebtID))
}

// GET_PAYMENT_HISTORY retrieves the payment history for the debt token.
func (token *SYN845Token) GET_PAYMENT_HISTORY() ([]PaymentEntry, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.PaymentHistory, nil
}

// Supporting Structures and Enums

type DebtStatus string

const (
    Active    DebtStatus = "active"
    Defaulted DebtStatus = "defaulted"
    Repaid    DebtStatus = "repaid"
)

// PaymentEntry represents an individual payment toward the debt.
type PaymentEntry struct {
    PaymentDate time.Time `json:"payment_date"`
    Amount      float64   `json:"amount"`
    Interest    float64   `json:"interest"`
    Principal   float64   `json:"principal"`
    Balance     float64   `json:"balance"`
}
