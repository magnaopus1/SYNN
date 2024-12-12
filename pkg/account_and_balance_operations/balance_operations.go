package account_and_balance_operations

import (
	"errors"
	"fmt"
	"log"
	"math"
	"synnergy_network/pkg/ledger"
	"time"
)

// BalanceGet retrieves the current balance of an account.
func BalanceGet(l *ledger.Ledger, accountID string) (float64, error) {
	logAction("BalanceGet - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceGet - Validation Failed: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceGet - Account Retrieval Failed: %v", err)
		return 0, fmt.Errorf("account not found: %w", err)
	}

	log.Printf("BalanceGet - Success: Account %s, Balance %.2f", accountID, account.Balance)
	return account.Balance, nil
}


// BalanceSet sets an account’s balance to a specific amount.
func BalanceSet(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceSet - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceSet - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount < 0 {
		log.Printf("BalanceSet - Validation Failed: Negative Balance")
		return errors.New("balance amount cannot be negative")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceSet - Account Retrieval Failed: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	// Set the new balance
	account.Balance = amount

	// Update the ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceSet - Ledger Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceSet - Success: Account %s, New Balance %.2f", accountID, amount)
	return nil
}


// BalanceIncrease adds a specific amount to the account’s balance.
func BalanceIncrease(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceIncrease - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceIncrease - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceIncrease - Validation Failed: Non-Positive Amount")
		return errors.New("increase amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceIncrease - Account Retrieval Failed: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	// Increase the balance
	account.Balance += amount

	// Update the ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceIncrease - Ledger Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceIncrease - Success: Account %s, Increased By %.2f, New Balance %.2f", accountID, amount, account.Balance)
	return nil
}


// BalanceDecrease subtracts a specific amount from the account’s balance.
func BalanceDecrease(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceDecrease - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceDecrease - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceDecrease - Validation Failed: Non-Positive Amount")
		return errors.New("decrease amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceDecrease - Account Retrieval Failed: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		log.Printf("BalanceDecrease - Insufficient Funds: Account %s, Available %.2f, Required %.2f", accountID, account.Balance, amount)
		return errors.New("insufficient funds")
	}

	// Decrease the balance
	account.Balance -= amount

	// Update the ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceDecrease - Ledger Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceDecrease - Success: Account %s, Decreased By %.2f, New Balance %.2f", accountID, amount, account.Balance)
	return nil
}


// AccountBalanceTransfer moves funds from one account to another.
func AccountBalanceTransfer(l *ledger.Ledger, fromID, toID string, amount float64) error {
	logAction("AccountBalanceTransfer - Start", fromID)

	// Validate inputs
	if err := validateAccountID(fromID); err != nil {
		log.Printf("AccountBalanceTransfer - Validation Failed: %v", err)
		return fmt.Errorf("invalid source account ID: %w", err)
	}
	if err := validateAccountID(toID); err != nil {
		log.Printf("AccountBalanceTransfer - Validation Failed: %v", err)
		return fmt.Errorf("invalid destination account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("AccountBalanceTransfer - Validation Failed: Non-Positive Amount")
		return errors.New("transfer amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve accounts
	fromAccount, err := l.AccountsWalletLedger.GetAccount(fromID)
	if err != nil {
		log.Printf("AccountBalanceTransfer - Source Account Retrieval Failed: %v", err)
		return fmt.Errorf("source account not found: %w", err)
	}
	if fromAccount.Balance < amount {
		log.Printf("AccountBalanceTransfer - Insufficient Funds: %s, Balance: %.2f, Required: %.2f", fromID, fromAccount.Balance, amount)
		return errors.New("insufficient funds for transfer")
	}

	toAccount, err := l.AccountsWalletLedger.GetAccount(toID)
	if err != nil {
		log.Printf("AccountBalanceTransfer - Destination Account Retrieval Failed: %v", err)
		return fmt.Errorf("destination account not found: %w", err)
	}

	// Perform transfer
	fromAccount.Balance -= amount
	toAccount.Balance += amount

	// Update ledger
	if err := l.AccountsWalletLedger.UpdateAccount(fromID, *fromAccount); err != nil {
		log.Printf("AccountBalanceTransfer - Source Account Update Failed: %v", err)
		return fmt.Errorf("failed to update source account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(toID, *toAccount); err != nil {
		log.Printf("AccountBalanceTransfer - Destination Account Update Failed: %v", err)
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	log.Printf("AccountBalanceTransfer - Success: From %s To %s Amount %.2f", fromID, toID, amount)
	return nil
}


// BalanceTransferFrom authorizes and executes a transfer from a specific source account.
func BalanceTransferFrom(l *ledger.Ledger, approvedAccountID, fromID, toID string, amount float64) error {
	logAction("BalanceTransferFrom - Start", fromID)

	// Validate inputs
	if approvedAccountID != fromID {
		log.Printf("BalanceTransferFrom - Unauthorized Transfer Attempt: ApprovedAccountID %s, FromID %s", approvedAccountID, fromID)
		return errors.New("unauthorized transfer")
	}

	// Delegate to AccountBalanceTransfer
	err := AccountBalanceTransfer(l, fromID, toID, amount)
	if err != nil {
		log.Printf("BalanceTransferFrom - Transfer Failed: %v", err)
		return err
	}

	log.Printf("BalanceTransferFrom - Success: Authorized by %s From %s To %s Amount %.2f", approvedAccountID, fromID, toID, amount)
	return nil
}


// BalanceApprove sets authorization for a specific transfer or operation.
func BalanceApprove(l *ledger.Ledger, accountID, approverID, transactionID string, amount float64, expiresAt time.Time, remarks string) error {
	logAction("BalanceApprove - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceApprove - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceApprove - Validation Failed: Non-Positive Amount")
		return errors.New("approval amount must be positive")
	}
	if time.Now().After(expiresAt) {
		log.Printf("BalanceApprove - Validation Failed: Expired Approval Time")
		return errors.New("expiration time must be in the future")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceApprove - Account Retrieval Failed: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	// Create a new approval
	approval := ledger.Approval{
		ID:            generateUniqueID(),
		AccountID:     accountID,
		ApproverID:    approverID,
		TransactionID: transactionID,
		Amount:        amount,
		ApprovedAt:    time.Now(),
		ExpiresAt:     expiresAt,
		Status:        "Approved",
		Remarks:       remarks,
	}

	// Append approval to account
	account.Approvals = append(account.Approvals, approval)

	// Update ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceApprove - Ledger Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceApprove - Success: Account %s Approved By %s For Transaction %s Amount %.2f ExpiresAt %s Remarks %s",
		accountID, approverID, transactionID, amount, expiresAt.Format(time.RFC3339), remarks)
	return nil
}


// BalanceRevoke removes a previously set authorization from an account.
func BalanceRevoke(l *ledger.Ledger, accountID, approvedID string) error {
	logAction("BalanceRevoke - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceRevoke - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceRevoke - Account Retrieval Failed: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	// Search for the approval and remove it
	removed := false
	for i, approval := range account.Approvals {
		if approval.ID == approvedID {
			account.Approvals = append(account.Approvals[:i], account.Approvals[i+1:]...)
			removed = true
			break
		}
	}
	if !removed {
		log.Printf("BalanceRevoke - Approval Not Found: %s", approvedID)
		return fmt.Errorf("approval ID %s not found", approvedID)
	}

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceRevoke - Ledger Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceRevoke - Success: Approval %s Removed for Account %s", approvedID, accountID)
	return nil
}


// BalanceUnlock unlocks a previously locked amount in the account balance.
func BalanceUnlock(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceUnlock - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceUnlock - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceUnlock - Validation Failed: Non-Positive Amount")
		return errors.New("unlock amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceUnlock - Account Retrieval Failed: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}
	if account.LockedBalance < amount {
		log.Printf("BalanceUnlock - Insufficient Locked Balance: Locked %.2f, Requested %.2f", account.LockedBalance, amount)
		return errors.New("insufficient locked balance")
	}

	// Adjust balances
	account.LockedBalance -= amount
	account.Balance += amount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceUnlock - Ledger Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceUnlock - Success: Unlocked %.2f for Account %s", amount, accountID)
	return nil
}


// BalanceQuery returns an account’s complete balance information.
func BalanceQuery(l *ledger.Ledger, accountID string) (BalanceInfo, error) {
	logAction("BalanceQuery - Start", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceQuery - Validation Failed: %v", err)
		return BalanceInfo{}, fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceQuery - Account Retrieval Failed: %v", err)
		return BalanceInfo{}, fmt.Errorf("account not found: %w", err)
	}

	info := BalanceInfo{
		AccountID:         accountID,
		AvailableBalance:  account.Balance,
		LockedBalance:     account.LockedBalance,
		HeldBalance:       account.ReservedBalance,
		LastTransactionID: account.LastTransactionID,
		TotalDeposited:    account.TotalDeposited,
		TotalWithdrawn:    account.TotalWithdrawn,
		Status:            account.BalanceStatus,
		LastUpdated:       account.LastUpdated,
	}

	log.Printf("BalanceQuery - Success: AccountID %s, BalanceInfo %+v", accountID, info)
	return info, nil
}



// BalanceSum calculates the sum of two account balances.
func BalanceSum(l *ledger.Ledger, accountID1, accountID2 string) (float64, error) {
	logAction("BalanceSum - Start", fmt.Sprintf("%s + %s", accountID1, accountID2))

	// Validate inputs
	if err := validateAccountID(accountID1); err != nil {
		log.Printf("BalanceSum - Validation Failed for AccountID1: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}
	if err := validateAccountID(accountID2); err != nil {
		log.Printf("BalanceSum - Validation Failed for AccountID2: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}

	// Retrieve balances
	balance1, err := BalanceGet(l, accountID1)
	if err != nil {
		log.Printf("BalanceSum - Failed to Retrieve Balance for AccountID1: %v", err)
		return 0, err
	}
	balance2, err := BalanceGet(l, accountID2)
	if err != nil {
		log.Printf("BalanceSum - Failed to Retrieve Balance for AccountID2: %v", err)
		return 0, err
	}

	sum := balance1 + balance2
	log.Printf("BalanceSum - Success: AccountID1 %s Balance %.2f + AccountID2 %s Balance %.2f = Total %.2f", accountID1, balance1, accountID2, balance2, sum)
	return sum, nil
}


// BalanceSubtract calculates the difference between two account balances.
func BalanceSubtract(l *ledger.Ledger, accountID1, accountID2 string) (float64, error) {
	logAction("BalanceSubtract - Start", fmt.Sprintf("Subtract: %s - %s", accountID1, accountID2))

	// Validate inputs
	if err := validateAccountID(accountID1); err != nil {
		log.Printf("BalanceSubtract - Validation Failed for AccountID1: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}
	if err := validateAccountID(accountID2); err != nil {
		log.Printf("BalanceSubtract - Validation Failed for AccountID2: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}

	// Get balances
	balance1, err := BalanceGet(l, accountID1)
	if err != nil {
		log.Printf("BalanceSubtract - Balance Retrieval Failed for AccountID1: %v", err)
		return 0, err
	}
	balance2, err := BalanceGet(l, accountID2)
	if err != nil {
		log.Printf("BalanceSubtract - Balance Retrieval Failed for AccountID2: %v", err)
		return 0, err
	}

	result := balance1 - balance2
	log.Printf("BalanceSubtract - Success: %s Balance %.2f - %s Balance %.2f = %.2f", accountID1, balance1, accountID2, balance2, result)
	return result, nil
}


// BalanceMultiply multiplies an account balance by a specified factor.
func BalanceMultiply(l *ledger.Ledger, accountID string, factor float64) (float64, error) {
	logAction("BalanceMultiply - Start", fmt.Sprintf("Account: %s, Factor: %.2f", accountID, factor))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceMultiply - Validation Failed: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}
	if factor == 0 {
		log.Printf("BalanceMultiply - Validation Failed: Factor Cannot Be Zero")
		return 0, errors.New("multiplication factor cannot be zero")
	}

	// Get balance
	balance, err := BalanceGet(l, accountID)
	if err != nil {
		log.Printf("BalanceMultiply - Balance Retrieval Failed: %v", err)
		return 0, err
	}

	result := balance * factor
	log.Printf("BalanceMultiply - Success: %s Balance %.2f * %.2f = %.2f", accountID, balance, factor, result)
	return result, nil
}


// BalanceDivide divides an account balance by a specified divisor.
func BalanceDivide(l *ledger.Ledger, accountID string, divisor float64) (float64, error) {
	logAction("BalanceDivide - Start", fmt.Sprintf("Account: %s, Divisor: %.2f", accountID, divisor))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceDivide - Validation Failed: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}
	if divisor == 0 {
		log.Printf("BalanceDivide - Validation Failed: Division by Zero")
		return 0, errors.New("division by zero is not allowed")
	}

	// Get balance
	balance, err := BalanceGet(l, accountID)
	if err != nil {
		log.Printf("BalanceDivide - Balance Retrieval Failed: %v", err)
		return 0, err
	}

	result := balance / divisor
	log.Printf("BalanceDivide - Success: %s Balance %.2f / %.2f = %.2f", accountID, balance, divisor, result)
	return result, nil
}


// BalanceModulo computes the remainder of an account balance divided by a specified divisor.
func BalanceModulo(l *ledger.Ledger, accountID string, divisor float64) (float64, error) {
	logAction("BalanceModulo - Start", fmt.Sprintf("Account: %s, Divisor: %.2f", accountID, divisor))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceModulo - Validation Failed: %v", err)
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}
	if divisor == 0 {
		log.Printf("BalanceModulo - Validation Failed: Division by Zero")
		return 0, errors.New("division by zero is not allowed")
	}

	// Get balance
	balance, err := BalanceGet(l, accountID)
	if err != nil {
		log.Printf("BalanceModulo - Balance Retrieval Failed: %v", err)
		return 0, err
	}

	result := math.Mod(balance, divisor)
	log.Printf("BalanceModulo - Success: %s Balance %.2f %% %.2f = %.2f", accountID, balance, divisor, result)
	return result, nil
}
