package account_and_balance_operations

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)
var accountMutex sync.Mutex

// HasPermission checks if the account has the specified permission.
func (a Account) HasPermission(permission string) bool {
	log.Printf("Checking permission: %s for account: %s", permission, a.ID)

	// Input validation
	if strings.TrimSpace(permission) == "" {
		log.Printf("Permission check failed: invalid permission string")
		return false
	}

	// Case-insensitive permission comparison
	for _, perm := range a.Permissions {
		if strings.EqualFold(perm, permission) {
			log.Printf("Permission granted: %s for account: %s", permission, a.ID)
			return true
		}
	}

	log.Printf("Permission denied: %s for account: %s", permission, a.ID)
	return false
}




// AccountCreate creates a new account with an initial setup.
func AccountCreate(l *ledger.Ledger, accountID string, initialBalance float64, creator Account) error {
	log.Printf("Starting account creation for ID: %s", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if initialBalance < 0 {
		return errors.New("initial balance cannot be negative")
	}

	// Validate creator’s admin privileges
	if !creator.HasPermission("admin:create_account") {
		return errors.New("account creation permission denied")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Initialize new account structure
	encryption := &common.Encryption{}
	encryptedData, err := encryption.EncryptData("AES", []byte(fmt.Sprintf("Account ID: %s", accountID)), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	account := ledger.Account{
		Address:       accountID,
		Balance:       initialBalance,
		CreatedAt:     time.Now(),
		EncryptedKey:  string(encryptedData),
	}

	// Record creation in ledger
	if err := l.AccountsWalletLedger.RecordAccount(accountID, account); err != nil {
		return fmt.Errorf("failed to record account: %w", err)
	}

	// Integrate with Synnergy Consensus
	if err := l.BlockchainConsensusCoinLedger.SyncWithConsensus(); err != nil {
		return fmt.Errorf("failed to sync with consensus: %w", err)
	}

	// Verify ledger consistency post-operation
	if _, err := l.AccountsWalletLedger.GetAccount(accountID); err != nil {
		return fmt.Errorf("post-verification failed: %w", err)
	}

	log.Printf("Account successfully created for ID: %s", accountID)
	return nil
}


// AccountDelete removes an account from the system after checks.
func AccountDelete(l *ledger.Ledger, accountID string, requester Account) error {
	log.Printf("Starting account deletion for ID: %s", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account from the ledger
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Validate requester's permission to delete account
	if !requester.HasPermission("admin:delete_account") {
		return errors.New("account deletion permission denied")
	}

	// Mark account as deleted and overwrite sensitive data
	encryption := &common.Encryption{}
	encryptedData, err := encryption.EncryptData("AES", []byte("Deleted Account"), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	account.EncryptedKey = string(encryptedData)
	account.Balance = 0.0 // Reset balance to 0

	// Update account status in ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account status: %w", err)
	}

	// Integrate with Synnergy Consensus
	if err := l.BlockchainConsensusCoinLedger.SyncWithConsensus(); err != nil {
		return fmt.Errorf("failed to sync with consensus: %w", err)
	}

	// Verify account removal from ledger
	if _, err := l.AccountsWalletLedger.GetAccount(accountID); err == nil {
		return fmt.Errorf("post-verification failed: account %s still exists", accountID)
	}

	log.Printf("Account successfully deleted for ID: %s", accountID)
	return nil
}


// AccountFreeze freezes an account, disabling transactions and setting its state to "frozen".
func AccountFreeze(l *ledger.Ledger, accountID string, requester Account) error {
	log.Printf("Starting account freeze for ID: %s", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account from the ledger
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Validate requester's permission
	if !requester.HasPermission("admin:freeze_account") {
		return errors.New("account freeze permission denied")
	}

	// Update account state to "frozen"
	account.IsFrozen = true

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	// Integrate with Synnergy Consensus
	if err := l.BlockchainConsensusCoinLedger.SyncWithConsensus(); err != nil {
		return fmt.Errorf("failed to sync with consensus: %w", err)
	}

	log.Printf("Account successfully frozen for ID: %s", accountID)
	return nil
}



// AccountUnfreeze reactivates a frozen account, restoring its operational state.
func AccountUnfreeze(l *ledger.Ledger, accountID string, requester Account) error {
	log.Printf("Starting account unfreeze for ID: %s", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account from the ledger
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Validate requester's permission
	if !requester.HasPermission("admin:unfreeze_account") {
		return errors.New("account unfreeze permission denied")
	}

	// Ensure account is frozen before attempting to unfreeze
	if !account.IsFrozen {
		return errors.New("account is not currently frozen")
	}

	// Update account state to "active"
	account.IsFrozen = false

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	// Integrate with Synnergy Consensus
	if err := l.BlockchainConsensusCoinLedger.SyncWithConsensus(); err != nil {
		return fmt.Errorf("failed to sync with consensus: %w", err)
	}

	log.Printf("Account successfully unfrozen for ID: %s", accountID)
	return nil
}


// AccountSetAdmin grants admin permissions to a user.
func AccountSetAdmin(l *ledger.Ledger, accountID string, requester Account) error {
	log.Printf("Starting admin rights assignment for ID: %s", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account from the ledger
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Validate requester's permission
	if !requester.HasPermission("admin:set_admin") {
		return errors.New("set admin permission denied")
	}

	// Grant admin status
	account.IsAdmin = true

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Admin rights successfully granted to account: %s", accountID)
	return nil
}


// AccountRemoveAdmin revokes admin permissions from a user.
func AccountRemoveAdmin(l *ledger.Ledger, accountID string, requester Account) error {
	log.Printf("Starting admin rights removal for ID: %s", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account from the ledger
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Validate requester's permission
	if !requester.HasPermission("admin:remove_admin") {
		return errors.New("remove admin permission denied")
	}

	// Revoke admin status
	account.IsAdmin = false

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Admin rights successfully removed from account: %s", accountID)
	return nil
}


// AccountAddPermission adds a specific permission to an account.
func AccountAddPermission(l *ledger.Ledger, accountID, permission string, requester Account) error {
	logAction("AddPermission", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if strings.TrimSpace(permission) == "" {
		return errors.New("permission cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Check requester permissions
	if !requester.HasPermission("admin:modify_permissions") {
		return errors.New("add permission denied")
	}

	// Check if the permission already exists
	for _, perm := range account.Permissions {
		if strings.EqualFold(perm, permission) {
			return fmt.Errorf("permission '%s' already exists for account %s", permission, accountID)
		}
	}

	// Add the permission
	account.Permissions = append(account.Permissions, permission)
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Permission '%s' successfully added to account: %s", permission, accountID)
	return nil
}


// validateAccountID ensures the account ID is valid.
func validateAccountID(accountID string) error {
	if accountID == "" {
		return errors.New("account ID cannot be empty")
	}
	if len(accountID) < 5 {
		return errors.New("account ID must be at least 5 characters long")
	}
	return nil
}

// logAction logs account-related actions.
func logAction(action, accountID string) {
	log.Printf("[%s] Action: %s, AccountID: %s", time.Now().Format(time.RFC3339), action, accountID)
}

// Enhanced Functions

// AccountRemovePermission removes a specific permission from an account.
func AccountRemovePermission(l *ledger.Ledger, accountID, permission string, requester Account) error {
	logAction("RemovePermission", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if strings.TrimSpace(permission) == "" {
		return errors.New("permission cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Check requester permissions
	if !requester.HasPermission("admin:modify_permissions") {
		return errors.New("remove permission denied")
	}

	// Find and remove the permission
	for i, perm := range account.Permissions {
		if strings.EqualFold(perm, permission) {
			account.Permissions = append(account.Permissions[:i], account.Permissions[i+1:]...)
			if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
				return fmt.Errorf("failed to update account: %w", err)
			}
			log.Printf("Permission '%s' successfully removed from account: %s", permission, accountID)
			return nil
		}
	}

	return fmt.Errorf("permission '%s' not found for account %s", permission, accountID)
}


// AccountGetBalance retrieves the current balance of an account.
func AccountGetBalance(l *ledger.Ledger, accountID string) (float64, error) {
	logAction("GetBalance", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return 0, fmt.Errorf("account not found: %w", err)
	}

	log.Printf("Retrieved balance for account %s: %.2f", accountID, account.Balance)
	return account.Balance, nil
}


// AccountCredit adds funds to an account.
func AccountCredit(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("CreditAccount", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Update the balance
	account.Balance += amount
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to credit account: %w", err)
	}

	log.Printf("Account credited successfully. AccountID: %s, Amount: %.2f, New Balance: %.2f", accountID, amount, account.Balance)
	return nil
}


// AccountDebit removes funds from an account if balance permits.
func AccountDebit(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("DebitAccount", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Check for sufficient balance
	if account.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Update the balance
	account.Balance -= amount
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to debit account: %w", err)
	}

	log.Printf("Account debited successfully. AccountID: %s, Amount: %.2f, New Balance: %.2f", accountID, amount, account.Balance)
	return nil
}


// AccountBalanceSnapshot creates a snapshot of the current balance.
func AccountBalanceSnapshot(l *ledger.Ledger, accountID string) (ledger.BalanceSnapshot, error) {
	logAction("BalanceSnapshot", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return ledger.BalanceSnapshot{}, fmt.Errorf("invalid account ID: %w", err)
	}

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return ledger.BalanceSnapshot{}, fmt.Errorf("account not found: %w", err)
	}

	// Create the balance snapshot
	snapshot := ledger.BalanceSnapshot{
		AccountID: accountID,
		Balance:   account.Balance,
		Timestamp: time.Now(),
	}

	// Save the snapshot to the ledger
	if err := l.AccountsWalletLedger.SaveBalanceSnapshot(snapshot); err != nil {
		return ledger.BalanceSnapshot{}, fmt.Errorf("failed to save balance snapshot: %w", err)
	}

	log.Printf("Balance snapshot created successfully for account: %s. Snapshot: %+v", accountID, snapshot)
	return snapshot, nil
}


// AccountBalanceRestore restores an account's balance to a previous snapshot.
func AccountBalanceRestore(l *ledger.Ledger, accountID string, snapshot ledger.BalanceSnapshot) error {
	logAction("BalanceRestore", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if snapshot.AccountID != accountID {
		return errors.New("snapshot account ID does not match")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Restore the balance
	account.Balance = snapshot.Balance
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to restore account balance: %w", err)
	}

	log.Printf("Balance successfully restored for account: %s. Restored Balance: %.2f, Snapshot Timestamp: %s", accountID, snapshot.Balance, snapshot.Timestamp)
	return nil
}


// AccountGetBalanceAt retrieves the balance of an account at a specific time.
func AccountGetBalanceAt(l *ledger.Ledger, accountID string, timestamp time.Time) (float64, error) {
	logAction("GetBalanceAt", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}
	if timestamp.After(time.Now()) {
		return 0, errors.New("timestamp cannot be in the future")
	}

	// Retrieve the balance at the specified timestamp
	balance, err := l.AccountsWalletLedger.GetBalanceAt(accountID, timestamp)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve balance at timestamp %s: %w", timestamp.Format(time.RFC3339), err)
	}

	// Convert big.Int to float64 for compatibility
	floatBalance, _ := new(big.Float).SetInt(&balance).Float64()

	log.Printf("Balance retrieved for account: %s at timestamp: %s, Balance: %.2f", accountID, timestamp.Format(time.RFC3339), floatBalance)
	return floatBalance, nil
}


// AccountBalanceHistory retrieves the balance history of an account.
func AccountBalanceHistory(l *ledger.Ledger, accountID string) ([]ledger.BalanceSnapshot, error) {
	logAction("BalanceHistory", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Retrieve balance history
	history, err := l.AccountsWalletLedger.GetBalanceHistory(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve balance history for account %s: %w", accountID, err)
	}

	log.Printf("Balance history retrieved for account: %s, Entries: %d", accountID, len(history))
	return history, nil
}


// AccountBalanceUpdate updates the balance of an account directly (restricted to admins).
func AccountBalanceUpdate(l *ledger.Ledger, accountID string, newBalance float64, requester Account) error {
	logAction("BalanceUpdate", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if newBalance < 0 {
		return errors.New("new balance cannot be negative")
	}
	if !requester.HasPermission("admin:update_balance") {
		return errors.New("balance update permission denied")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve the account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Update the account balance
	oldBalance := account.Balance
	account.Balance = newBalance

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	log.Printf("Account balance updated successfully. AccountID: %s, Old Balance: %.2f, New Balance: %.2f", accountID, oldBalance, newBalance)
	return nil
}
