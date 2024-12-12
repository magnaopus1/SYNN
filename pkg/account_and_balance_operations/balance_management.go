package account_and_balance_operations

import (
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)

// BalanceRemoveFromExternalAccount removes funds from an external account.
func BalanceRemoveFromExternalAccount(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceRemoveFromExternalAccount", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		return errors.New("insufficient balance in external account")
	}

	// Deduct the amount
	account.Balance -= amount

	// Update the account in the ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Funds successfully removed from external account. AccountID: %s, Amount: %.2f, Remaining Balance: %.2f",
		accountID, amount, account.Balance)
	return nil
}


// BalanceRollup consolidates smaller balances into a main balance.
func BalanceRollup(l *ledger.Ledger, accountID string, subAccounts []string) error {
	logAction("BalanceRollup", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid main account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the main account
	mainAccount, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("main account not found: %w", err)
	}

	// Consolidate balances
	for _, subAccountID := range subAccounts {
		if err := validateAccountID(subAccountID); err != nil {
			return fmt.Errorf("invalid sub-account ID: %w", err)
		}
		subAccount, err := l.AccountsWalletLedger.GetAccount(subAccountID)
		if err != nil {
			return fmt.Errorf("sub-account not found: %w", err)
		}

		mainAccount.Balance += subAccount.Balance
		subAccount.Balance = 0

		// Update the sub-account
		if err := l.AccountsWalletLedger.UpdateAccount(subAccountID, *subAccount); err != nil {
			return fmt.Errorf("failed to update sub-account: %w", err)
		}

		log.Printf("Sub-account %s rolled up into main account %s. Transferred Balance: %.2f",
			subAccountID, accountID, subAccount.Balance)
	}

	// Update the main account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *mainAccount); err != nil {
		return fmt.Errorf("failed to update main account: %w", err)
	}

	log.Printf("Balance rollup completed successfully. Main AccountID: %s, Total Balance: %.2f",
		accountID, mainAccount.Balance)
	return nil
}


// BalanceSplit divides a balance into smaller amounts across multiple accounts.
func BalanceSplit(l *ledger.Ledger, accountID string, targetAccounts []string, amounts []float64) error {
	logAction("BalanceSplit", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid source account ID: %w", err)
	}
	if len(targetAccounts) != len(amounts) {
		return errors.New("number of target accounts must match the number of amounts")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the source account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("source account not found: %w", err)
	}

	// Validate and calculate total amount
	total := 0.0
	for _, amount := range amounts {
		if amount <= 0 {
			return errors.New("split amounts must be positive")
		}
		total += amount
	}
	if account.Balance < total {
		return errors.New("insufficient funds to split")
	}

	// Distribute balances
	for i, targetID := range targetAccounts {
		if err := validateAccountID(targetID); err != nil {
			return fmt.Errorf("invalid target account ID: %w", err)
		}
		targetAccount, err := l.AccountsWalletLedger.GetAccount(targetID)
		if err != nil {
			return fmt.Errorf("target account not found: %w", err)
		}

		// Adjust balances
		targetAccount.Balance += amounts[i]
		account.Balance -= amounts[i]

		// Update target account
		if err := l.AccountsWalletLedger.UpdateAccount(targetID, *targetAccount); err != nil {
			return fmt.Errorf("failed to update target account: %w", err)
		}

		log.Printf("Funds split to target account %s. Amount: %.2f", targetID, amounts[i])
	}

	// Update the source account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update source account: %w", err)
	}

	log.Printf("Balance split completed successfully. Source AccountID: %s, Remaining Balance: %.2f",
		accountID, account.Balance)
	return nil
}


// BalanceMerge combines balances from multiple accounts into one.
func BalanceMerge(l *ledger.Ledger, targetID string, sourceAccounts []string) error {
	logAction("BalanceMerge", targetID)

	// Validate inputs
	if err := validateAccountID(targetID); err != nil {
		return fmt.Errorf("invalid target account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve target account
	targetAccount, err := l.AccountsWalletLedger.GetAccount(targetID)
	if err != nil {
		return fmt.Errorf("target account not found: %w", err)
	}

	// Merge balances
	for _, sourceID := range sourceAccounts {
		if err := validateAccountID(sourceID); err != nil {
			return fmt.Errorf("invalid source account ID: %w", err)
		}
		sourceAccount, err := l.AccountsWalletLedger.GetAccount(sourceID)
		if err != nil {
			return fmt.Errorf("source account not found: %w", err)
		}

		// Transfer balance
		targetAccount.Balance += sourceAccount.Balance
		sourceAccount.Balance = 0

		// Update source account
		if err := l.AccountsWalletLedger.UpdateAccount(sourceID, *sourceAccount); err != nil {
			return fmt.Errorf("failed to update source account: %w", err)
		}

		log.Printf("Merged balance from source account %s into target account %s. Transferred Amount: %.2f",
			sourceID, targetID, sourceAccount.Balance)
	}

	// Update target account
	if err := l.AccountsWalletLedger.UpdateAccount(targetID, *targetAccount); err != nil {
		return fmt.Errorf("failed to update target account: %w", err)
	}

	log.Printf("Balance merge completed successfully. Target AccountID: %s, Total Balance: %.2f",
		targetID, targetAccount.Balance)
	return nil
}


// AccountUpdateBalanceStatus updates an account’s balance status.
func AccountUpdateBalanceStatus(l *ledger.Ledger, accountID string, status ledger.BalanceStatus) error {
	logAction("AccountUpdateBalanceStatus", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Update balance status
	account.BalanceStatus = status

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account balance status updated. AccountID: %s, Status: %+v", accountID, status)
	return nil
}


// AccountRevertBalanceStatus reverts an account’s balance status to ACTIVE.
func AccountRevertBalanceStatus(l *ledger.Ledger, accountID string) error {
	logAction("AccountRevertBalanceStatus", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Define active status
	activeStatus := ledger.BalanceStatus{
		IsActive:     true,
		IsFrozen:     false,
		IsOnHold:     false,
		FreezeReason: "",
		HoldReason:   "",
		UpdatedAt:    time.Now(),
	}

	// Reuse AccountUpdateBalanceStatus
	if err := AccountUpdateBalanceStatus(l, accountID, activeStatus); err != nil {
		return fmt.Errorf("failed to revert balance status: %w", err)
	}

	log.Printf("Account balance status reverted to ACTIVE. AccountID: %s", accountID)
	return nil
}

// BalanceReserve reserves a specific amount in an account.
func BalanceReserve(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceReserve", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("reservation amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		return errors.New("insufficient funds to reserve")
	}

	// Reserve balance
	account.Balance -= amount
	account.ReservedBalance += amount

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Balance reserved successfully. AccountID: %s, Reserved Amount: %.2f", accountID, amount)
	return nil
}


// BalanceReleaseReserve releases reserved funds back to the main balance.
func BalanceReleaseReserve(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceReleaseReserve", accountID)

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
	if account.ReservedBalance < amount {
		return errors.New("insufficient reserved balance for release")
	}

	// Release reserved balance
	account.ReservedBalance -= amount
	account.Balance += amount

	// Update account
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Reserved balance released successfully. AccountID: %s, Released Amount: %.2f", accountID, amount)
	return nil
}


// BalanceTempFreeze temporarily freezes an account’s balance.
func BalanceTempFreeze(l *ledger.Ledger, accountID string, duration time.Duration) error {
	logAction("BalanceTempFreeze", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if duration <= 0 {
		return errors.New("freeze duration must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Apply the freeze
	account.FreezeUntil = time.Now().Add(duration)
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account %s temporarily frozen until %s", accountID, account.FreezeUntil.Format(time.RFC3339))
	return nil
}


// BalanceTempUnfreeze unfreezes an account’s balance if the freeze period has expired.
func BalanceTempUnfreeze(l *ledger.Ledger, accountID string) error {
	logAction("BalanceTempUnfreeze", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Check freeze expiration
	if time.Now().Before(account.FreezeUntil) {
		return errors.New("account still in freeze period")
	}

	// Remove freeze
	account.FreezeUntil = time.Time{}
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account %s unfrozen successfully", accountID)
	return nil
}


// BalanceFlagAsSuspicious flags an account balance as suspicious.
func BalanceFlagAsSuspicious(l *ledger.Ledger, accountID string) error {
	logAction("BalanceFlagAsSuspicious", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Flag as suspicious
	account.IsSuspicious = true
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account %s flagged as suspicious", accountID)
	return nil
}

// BalanceRemoveSuspiciousFlag removes the suspicious flag from an account balance.
func BalanceRemoveSuspiciousFlag(l *ledger.Ledger, accountID string) error {
	logAction("BalanceRemoveSuspiciousFlag", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Remove suspicious flag
	account.IsSuspicious = false
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Suspicious flag removed from account %s", accountID)
	return nil
}


// BalanceValidateThreshold checks if an account’s balance exceeds a given threshold.
func BalanceValidateThreshold(l *ledger.Ledger, accountID string, threshold float64) (bool, error) {
	logAction("BalanceValidateThreshold", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return false, fmt.Errorf("invalid account ID: %w", err)
	}
	if threshold < 0 {
		return false, errors.New("threshold must be non-negative")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return false, fmt.Errorf("account not found: %w", err)
	}

	isAboveThreshold := account.Balance >= threshold
	log.Printf("Account %s balance check: %.2f >= %.2f? %v", accountID, account.Balance, threshold, isAboveThreshold)
	return isAboveThreshold, nil
}

// BalanceTransferInBatch processes multiple balance transfers in a single batch.
func BalanceTransferInBatch(l *ledger.Ledger, transfers []ledger.BalanceTransfer) error {
	logAction("BalanceTransferInBatch", "BatchProcessing")

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Process each transfer in the batch
	for _, transfer := range transfers {
		// Validate inputs
		if err := validateAccountID(transfer.FromID); err != nil {
			return fmt.Errorf("invalid source account ID: %w", err)
		}
		if err := validateAccountID(transfer.ToID); err != nil {
			return fmt.Errorf("invalid destination account ID: %w", err)
		}
		if transfer.Amount <= 0 {
			return fmt.Errorf("transfer amount must be positive for transfer from %s to %s", transfer.FromID, transfer.ToID)
		}

		// Retrieve accounts
		fromAccount, err := l.AccountsWalletLedger.GetAccount(transfer.FromID)
		if err != nil {
			return fmt.Errorf("source account not found: %w", err)
		}
		toAccount, err := l.AccountsWalletLedger.GetAccount(transfer.ToID)
		if err != nil {
			return fmt.Errorf("destination account not found: %w", err)
		}

		// Check balance and execute transfer
		if fromAccount.Balance < transfer.Amount {
			return fmt.Errorf("insufficient funds in account: %s", transfer.FromID)
		}

		fromAccount.Balance -= transfer.Amount
		toAccount.Balance += transfer.Amount

		// Update accounts
		if err := l.AccountsWalletLedger.UpdateAccount(transfer.FromID, *fromAccount); err != nil {
			return fmt.Errorf("failed to update source account: %w", err)
		}
		if err := l.AccountsWalletLedger.UpdateAccount(transfer.ToID, *toAccount); err != nil {
			return fmt.Errorf("failed to update destination account: %w", err)
		}

		log.Printf("Transferred %.2f from %s to %s", transfer.Amount, transfer.FromID, transfer.ToID)
	}

	log.Printf("Batch transfer completed successfully with %d transfers", len(transfers))
	return nil
}
