package account_and_balance_operations

import (
	"errors"
	"fmt"
	"log"
	"math"
	"synnergy_network/pkg/ledger"
	"time"
)

// BalanceIsZero checks if an account's balance is zero.
func BalanceIsZero(l *ledger.Ledger, accountID string) (bool, error) {
	logAction("CheckBalanceIsZero", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return false, fmt.Errorf("invalid account ID: %w", err)
	}

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return false, fmt.Errorf("account not found: %w", err)
	}

	// Check if the balance is zero
	isZero := account.Balance == 0
	log.Printf("Balance check completed for account: %s. Is Zero: %t", accountID, isZero)
	return isZero, nil
}


// BalanceIsNegative checks if an account's balance is negative.
func BalanceIsNegative(l *ledger.Ledger, accountID string) (bool, error) {
	logAction("CheckBalanceIsNegative", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return false, fmt.Errorf("invalid account ID: %w", err)
	}

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return false, fmt.Errorf("account not found: %w", err)
	}

	// Check if the balance is negative
	isNegative := account.Balance < 0
	log.Printf("Balance check completed for account: %s. Is Negative: %t", accountID, isNegative)
	return isNegative, nil
}


// BalanceIsPositive checks if an account's balance is positive.
func BalanceIsPositive(l *ledger.Ledger, accountID string) (bool, error) {
	logAction("CheckBalanceIsPositive", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return false, fmt.Errorf("invalid account ID: %w", err)
	}

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return false, fmt.Errorf("account not found: %w", err)
	}

	// Check if the balance is positive
	isPositive := account.Balance > 0
	log.Printf("Balance check completed for account: %s. Is Positive: %t", accountID, isPositive)
	return isPositive, nil
}


// BalanceCompareEqual compares two account balances for equality.
func BalanceCompareEqual(l *ledger.Ledger, accountID1, accountID2 string) (bool, error) {
	logAction("CompareBalancesEqual", accountID1)

	// Validate inputs
	if err := validateAccountID(accountID1); err != nil {
		return false, fmt.Errorf("invalid account ID1: %w", err)
	}
	if err := validateAccountID(accountID2); err != nil {
		return false, fmt.Errorf("invalid account ID2: %w", err)
	}

	// Retrieve the accounts
	account1, err := l.AccountsWalletLedger.GetAccount(accountID1)
	if err != nil {
		return false, fmt.Errorf("account1 not found: %w", err)
	}

	account2, err := l.AccountsWalletLedger.GetAccount(accountID2)
	if err != nil {
		return false, fmt.Errorf("account2 not found: %w", err)
	}

	// Compare the balances
	isEqual := account1.Balance == account2.Balance
	log.Printf("Balance comparison completed. Account1: %s, Balance1: %.2f, Account2: %s, Balance2: %.2f, Equal: %t",
		accountID1, account1.Balance, accountID2, account2.Balance, isEqual)
	return isEqual, nil
}


// BalanceCompareGreater checks if accountID1 has a greater balance than accountID2.
func BalanceCompareGreater(l *ledger.Ledger, accountID1, accountID2 string) (bool, error) {
	logAction("CompareBalancesGreater", accountID1)

	// Validate inputs
	if err := validateAccountID(accountID1); err != nil {
		return false, fmt.Errorf("invalid account ID1: %w", err)
	}
	if err := validateAccountID(accountID2); err != nil {
		return false, fmt.Errorf("invalid account ID2: %w", err)
	}

	// Retrieve accounts
	account1, err := l.AccountsWalletLedger.GetAccount(accountID1)
	if err != nil {
		return false, fmt.Errorf("account1 not found: %w", err)
	}
	account2, err := l.AccountsWalletLedger.GetAccount(accountID2)
	if err != nil {
		return false, fmt.Errorf("account2 not found: %w", err)
	}

	// Compare balances
	isGreater := account1.Balance > account2.Balance
	log.Printf("Balance comparison completed. Account1: %s (%.2f), Account2: %s (%.2f), Result: %t",
		accountID1, account1.Balance, accountID2, account2.Balance, isGreater)
	return isGreater, nil
}


// BalanceCompareLess checks if accountID1 has a lesser balance than accountID2.
func BalanceCompareLess(l *ledger.Ledger, accountID1, accountID2 string) (bool, error) {
	logAction("CompareBalancesLess", accountID1)

	// Validate inputs
	if err := validateAccountID(accountID1); err != nil {
		return false, fmt.Errorf("invalid account ID1: %w", err)
	}
	if err := validateAccountID(accountID2); err != nil {
		return false, fmt.Errorf("invalid account ID2: %w", err)
	}

	// Retrieve accounts
	account1, err := l.AccountsWalletLedger.GetAccount(accountID1)
	if err != nil {
		return false, fmt.Errorf("account1 not found: %w", err)
	}
	account2, err := l.AccountsWalletLedger.GetAccount(accountID2)
	if err != nil {
		return false, fmt.Errorf("account2 not found: %w", err)
	}

	// Compare balances
	isLesser := account1.Balance < account2.Balance
	log.Printf("Balance comparison completed. Account1: %s (%.2f), Account2: %s (%.2f), Result: %t",
		accountID1, account1.Balance, accountID2, account2.Balance, isLesser)
	return isLesser, nil
}


// BalanceTransferWithLock transfers funds with a lock period.
func BalanceTransferWithLock(l *ledger.Ledger, fromID, toID string, amount float64, lockUntil time.Time) error {
	logAction("TransferWithLock", fromID)

	// Validate inputs
	if err := validateAccountID(fromID); err != nil {
		return fmt.Errorf("invalid source account ID: %w", err)
	}
	if err := validateAccountID(toID); err != nil {
		return fmt.Errorf("invalid destination account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}
	if time.Now().After(lockUntil) {
		return errors.New("lockUntil time must be in the future")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve accounts
	fromAccount, err := l.AccountsWalletLedger.GetAccount(fromID)
	if err != nil {
		return fmt.Errorf("source account not found: %w", err)
	}
	if fromAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	toAccount, err := l.AccountsWalletLedger.GetAccount(toID)
	if err != nil {
		return fmt.Errorf("destination account not found: %w", err)
	}

	// Perform transfer
	fromAccount.Balance -= amount
	toAccount.Balance += amount

	// Create and apply lock
	lock := ledger.BalanceLock{
		ID:        generateUniqueID(),
		AccountID: toID,
		Amount:    amount,
		UnlockAt:  lockUntil,
	}
	toAccount.LockedBalances = append(toAccount.LockedBalances, lock)

	// Update accounts in ledger
	if err := l.AccountsWalletLedger.UpdateAccount(fromID, *fromAccount); err != nil {
		return fmt.Errorf("failed to update source account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(toID, *toAccount); err != nil {
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	log.Printf("Transfer with lock successful. From: %s, To: %s, Amount: %.2f, LockUntil: %s",
		fromID, toID, amount, lockUntil.Format(time.RFC3339))
	return nil
}



// generateUniqueID creates a unique identifier for locks.
func generateUniqueID() string {
	return fmt.Sprintf("LOCK-%d", time.Now().UnixNano())
}


// BalanceTransferWithUnlock unlocks a locked balance and transfers it.
func BalanceTransferWithUnlock(l *ledger.Ledger, accountID, toID, lockID string) error {
	logAction("TransferWithUnlock", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if err := validateAccountID(toID); err != nil {
		return fmt.Errorf("invalid destination account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve source account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("source account not found: %w", err)
	}

	// Locate and validate lock
	var lock ledger.BalanceLock
	lockIndex := -1
	for i, lockedBalance := range account.LockedBalances {
		if lockedBalance.ID == lockID && lockedBalance.UnlockAt.Before(time.Now()) {
			lock = lockedBalance
			lockIndex = i
			break
		}
	}
	if lockIndex == -1 {
		return errors.New("no valid lock found for transfer")
	}

	// Retrieve destination account
	toAccount, err := l.AccountsWalletLedger.GetAccount(toID)
	if err != nil {
		return fmt.Errorf("destination account not found: %w", err)
	}

	// Perform transfer and update balances
	toAccount.Balance += lock.Amount
	account.LockedBalances = append(account.LockedBalances[:lockIndex], account.LockedBalances[lockIndex+1:]...)

	// Update accounts in ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update source account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(toID, *toAccount); err != nil {
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	log.Printf("Transfer with unlock successful. From: %s, To: %s, Amount: %.2f", accountID, toID, lock.Amount)
	return nil
}


// BalanceHold places a hold on a portion of the balance.
func BalanceHold(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceHold", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("hold amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		return errors.New("insufficient funds for hold")
	}

	// Apply hold
	account.Balance -= amount
	account.HeldBalance += amount

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Hold placed successfully. AccountID: %s, Amount: %.2f", accountID, amount)
	return nil
}


// BalanceRelease releases a held amount back to the balance.
func BalanceRelease(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceRelease", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("release amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.HeldBalance < amount {
		return errors.New("insufficient held balance for release")
	}

	// Release hold
	account.HeldBalance -= amount
	account.Balance += amount

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Held amount released successfully. AccountID: %s, Amount: %.2f", accountID, amount)
	return nil
}


// BalanceMint creates new balance in the specified account.
func BalanceMint(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceMint", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("mint amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Mint balance
	account.Balance += amount

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Balance minted successfully. AccountID: %s, Amount: %.2f", accountID, amount)
	return nil
}


// BalanceBurn destroys balance in the specified account.
func BalanceBurn(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceBurn", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("burn amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		return errors.New("insufficient funds for burn")
	}

	// Burn balance
	account.Balance -= amount

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Balance burned successfully. AccountID: %s, Amount: %.2f", accountID, amount)
	return nil
}


// BalanceCalculateInterest applies interest to an account's balance over a specified period.
func BalanceCalculateInterest(l *ledger.Ledger, accountID string, rate float64, duration time.Duration) error {
	logAction("CalculateInterest", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if rate <= 0 {
		return errors.New("interest rate must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Calculate interest based on annual rate and duration
	interest := account.Balance * rate * (float64(duration.Hours()) / (24 * 365))
	account.Balance += math.Round(interest*100) / 100 // Round to two decimal places

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Interest applied successfully. AccountID: %s, Rate: %.2f%%, Duration: %s, Interest: %.2f",
		accountID, rate*100, duration.String(), interest)
	return nil
}
