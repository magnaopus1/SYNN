package account_and_balance_operations

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// BalanceClear resets an account's balance to zero.
func BalanceClear(l *ledger.Ledger, accountID string) error {
	logAction("BalanceClear", accountID)

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

	// Clear balance
	account.Balance = 0
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account %s balance cleared successfully", accountID)
	return nil
}


// BalanceTransferToTrust moves funds to a trust account.
func BalanceTransferToTrust(l *ledger.Ledger, accountID, trustID string, amount float64) error {
	logAction("BalanceTransferToTrust", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if err := validateAccountID(trustID); err != nil {
		return fmt.Errorf("invalid trust ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve accounts
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		return errors.New("insufficient funds to transfer to trust")
	}

	trustAccount, err := l.AccountsWalletLedger.GetTrustAccount(trustID)
	if err != nil {
		return fmt.Errorf("trust account not found: %w", err)
	}

	// Perform transfer
	account.Balance -= amount
	amountFloat := new(big.Float).SetFloat64(amount)
	trustAccount.Balance = new(big.Float).Add(trustAccount.Balance, amountFloat)

	// Update accounts in ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateTrustAccount(trustID, *trustAccount); err != nil {
		return fmt.Errorf("failed to update trust account: %w", err)
	}

	log.Printf("Transferred %.2f from account %s to trust %s", amount, accountID, trustID)
	return nil
}


// BalanceTransferFromTrust retrieves funds from a trust account.
func BalanceTransferFromTrust(l *ledger.Ledger, trustID, accountID string, amount float64) error {
	logAction("BalanceTransferFromTrust", trustID)

	// Validate inputs
	if err := validateAccountID(trustID); err != nil {
		return fmt.Errorf("invalid trust ID: %w", err)
	}
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve accounts
	trustAccount, err := l.AccountsWalletLedger.GetTrustAccount(trustID)
	if err != nil {
		return fmt.Errorf("trust account not found: %w", err)
	}
	amountFloat := new(big.Float).SetFloat64(amount)
	if trustAccount.Balance.Cmp(amountFloat) < 0 {
		return errors.New("insufficient funds in trust account")
	}

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("destination account not found: %w", err)
	}

	// Perform transfer
	trustAccount.Balance = new(big.Float).Sub(trustAccount.Balance, amountFloat)
	account.Balance += amount

	// Update accounts in ledger
	if err := l.AccountsWalletLedger.UpdateTrustAccount(trustID, *trustAccount); err != nil {
		return fmt.Errorf("failed to update trust account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Transferred %.2f from trust %s to account %s", amount, trustID, accountID)
	return nil
}


// BalanceFreeze fully freezes an account’s balance.
func BalanceFreeze(l *ledger.Ledger, accountID string) error {
	logAction("BalanceFreeze", accountID)

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

	// Apply freeze
	account.IsFrozen = true
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account %s successfully frozen", accountID)
	return nil
}


// BalanceUnfreeze unfreezes an account’s balance.
func BalanceUnfreeze(l *ledger.Ledger, accountID string) error {
	logAction("BalanceUnfreeze", accountID)

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

	// Update freeze status
	if !account.IsFrozen {
		log.Printf("Account %s is not currently frozen", accountID)
		return nil
	}
	account.IsFrozen = false

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Account %s successfully unfrozen", accountID)
	return nil
}


// BalanceEncrypt encrypts an account balance using AES encryption.
func BalanceEncrypt(l *ledger.Ledger, accountID string) error {
	logAction("BalanceEncrypt", accountID)

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

	// Encrypt balance
	encryption := &common.Encryption{}
	balanceString := fmt.Sprintf("%.2f", account.Balance)
	encryptedBalance, err := encryption.EncryptData("AES", []byte(balanceString), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	account.EncryptedBalance = string(encryptedBalance)

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account with encrypted balance: %w", err)
	}

	log.Printf("Balance for account %s successfully encrypted", accountID)
	return nil
}

// BalanceDecrypt decrypts an account balance.
func BalanceDecrypt(l *ledger.Ledger, accountID string) (float64, error) {
	logAction("BalanceDecrypt", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return 0, fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	// Retrieve account
	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return 0, fmt.Errorf("account not found: %w", err)
	}

	// Decrypt balance
	encryption := &common.Encryption{}
	decryptedData, err := encryption.DecryptData([]byte(account.EncryptedBalance), common.EncryptionKey)
	if err != nil {
		return 0, fmt.Errorf("decryption failed: %w", err)
	}

	// Parse balance
	var balance float64
	_, err = fmt.Sscanf(string(decryptedData), "%f", &balance)
	if err != nil {
		return 0, fmt.Errorf("balance parsing failed: %w", err)
	}

	log.Printf("Balance for account %s successfully decrypted", accountID)
	return balance, nil
}


// BalanceAuthorizationSet assigns an authorization to an account.
func BalanceAuthorizationSet(l *ledger.Ledger, accountID string, authorization ledger.Authorization) error {
	logAction("BalanceAuthorizationSet", accountID)

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

	// Check for duplicate authorization
	for _, auth := range account.Authorizations {
		if auth == authorization {
			log.Printf("Authorization already exists for account %s", accountID)
			return nil
		}
	}

	// Add authorization
	account.Authorizations = append(account.Authorizations, authorization)

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account with new authorization: %w", err)
	}

	log.Printf("Authorization added to account %s", accountID)
	return nil
}


// BalanceAuthorizationRemove removes an authorization from an account.
func BalanceAuthorizationRemove(l *ledger.Ledger, accountID, authID string) error {
	logAction("BalanceAuthorizationRemove", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if authID == "" {
		return errors.New("authorization ID cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Find and remove authorization
	authRemoved := false
	for i, auth := range account.Authorizations {
		if auth.ID == authID {
			account.Authorizations = append(account.Authorizations[:i], account.Authorizations[i+1:]...)
			authRemoved = true
			break
		}
	}

	if !authRemoved {
		return fmt.Errorf("authorization ID %s not found for account %s", authID, accountID)
	}

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Authorization ID %s removed from account %s", authID, accountID)
	return nil
}


// AccountQueryBalanceStatus checks an account’s balance status.
func AccountQueryBalanceStatus(l *ledger.Ledger, accountID string) (ledger.BalanceStatus, error) {
	logAction("AccountQueryBalanceStatus", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return ledger.BalanceStatus{}, fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return ledger.BalanceStatus{}, fmt.Errorf("account not found: %w", err)
	}

	log.Printf("Queried balance status for account %s: %+v", accountID, account.BalanceStatus)
	return account.BalanceStatus, nil
}


// BalanceTransfer transfers a specified amount from one account to another.
func BalanceTransfer(l *ledger.Ledger, fromID, toID string, amount float64) error {
	logAction("BalanceTransfer", fromID)

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

	accountMutex.Lock()
	defer accountMutex.Unlock()

	fromAccount, err := l.AccountsWalletLedger.GetAccount(fromID)
	if err != nil {
		return fmt.Errorf("source account not found: %w", err)
	}

	toAccount, err := l.AccountsWalletLedger.GetAccount(toID)
	if err != nil {
		return fmt.Errorf("destination account not found: %w", err)
	}

	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient funds in source account %s", fromID)
	}

	// Perform the transfer
	fromAccount.Balance -= amount
	toAccount.Balance += amount

	if err := l.AccountsWalletLedger.UpdateAccount(fromID, *fromAccount); err != nil {
		return fmt.Errorf("failed to update source account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(toID, *toAccount); err != nil {
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	log.Printf("Successfully transferred %.2f from account %s to account %s", amount, fromID, toID)
	return nil
}


// BalanceTransferBatch processes a batch of transfers.
func BalanceTransferBatch(l *ledger.Ledger, transfers []ledger.BalanceTransfer) error {
	logAction("BalanceTransferBatch", "BatchProcessing")

	accountMutex.Lock()
	defer accountMutex.Unlock()

	for _, transfer := range transfers {
		if err := validateAccountID(transfer.FromID); err != nil {
			return fmt.Errorf("invalid source account ID %s: %w", transfer.FromID, err)
		}
		if err := validateAccountID(transfer.ToID); err != nil {
			return fmt.Errorf("invalid destination account ID %s: %w", transfer.ToID, err)
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

		// Validate balance
		if fromAccount.Balance < transfer.Amount {
			return fmt.Errorf("insufficient funds in source account %s", transfer.FromID)
		}

		// Perform transfer
		fromAccount.Balance -= transfer.Amount
		toAccount.Balance += transfer.Amount

		if err := l.AccountsWalletLedger.UpdateAccount(transfer.FromID, *fromAccount); err != nil {
			return fmt.Errorf("failed to update source account: %w", err)
		}
		if err := l.AccountsWalletLedger.UpdateAccount(transfer.ToID, *toAccount); err != nil {
			return fmt.Errorf("failed to update destination account: %w", err)
		}

		log.Printf("Batch transfer: %.2f from %s to %s", transfer.Amount, transfer.FromID, transfer.ToID)
	}

	log.Printf("Batch transfer processing complete for %d transfers", len(transfers))
	return nil
}


// BalanceTransferAll transfers the entire balance from one account to another.
func BalanceTransferAll(l *ledger.Ledger, fromID, toID string) error {
	logAction("BalanceTransferAll", fromID)

	// Validate inputs
	if err := validateAccountID(fromID); err != nil {
		return fmt.Errorf("invalid source account ID: %w", err)
	}
	if err := validateAccountID(toID); err != nil {
		return fmt.Errorf("invalid destination account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	fromAccount, err := l.AccountsWalletLedger.GetAccount(fromID)
	if err != nil {
		return fmt.Errorf("source account not found: %w", err)
	}

	toAccount, err := l.AccountsWalletLedger.GetAccount(toID)
	if err != nil {
		return fmt.Errorf("destination account not found: %w", err)
	}

	if fromAccount.Balance <= 0 {
		return errors.New("source account has no funds to transfer")
	}

	// Perform the transfer
	toAccount.Balance += fromAccount.Balance
	fromAccount.Balance = 0

	// Update accounts in ledger
	if err := l.AccountsWalletLedger.UpdateAccount(fromID, *fromAccount); err != nil {
		return fmt.Errorf("failed to update source account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(toID, *toAccount); err != nil {
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	log.Printf("Successfully transferred all balance from %s to %s", fromID, toID)
	return nil
}


// BalanceHoldForVerification places a hold on funds for verification.
func BalanceHoldForVerification(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceHoldForVerification", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("hold amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		return errors.New("insufficient balance to hold for verification")
	}

	// Place hold on funds
	account.Balance -= amount
	account.VerificationHold += amount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Successfully held %.2f for verification on account %s", amount, accountID)
	return nil
}


// BalanceReleaseAfterVerification releases funds after verification is complete.
func BalanceReleaseAfterVerification(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceReleaseAfterVerification", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		return errors.New("release amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.VerificationHold < amount {
		return errors.New("insufficient held funds for release")
	}

	// Release funds
	account.VerificationHold -= amount
	account.Balance += amount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Successfully released %.2f after verification on account %s", amount, accountID)
	return nil
}


// BalanceAdjustForLoss adjusts an account's balance for any reported loss.
func BalanceAdjustForLoss(l *ledger.Ledger, accountID string, lossAmount float64) error {
	logAction("BalanceAdjustForLoss", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if lossAmount <= 0 {
		return errors.New("loss amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < lossAmount {
		return errors.New("insufficient balance to cover loss adjustment")
	}

	// Adjust balance for loss
	account.Balance -= lossAmount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Successfully adjusted balance for loss of %.2f on account %s", lossAmount, accountID)
	return nil
}


// BalanceReflectOnLedger syncs the current balance state with the ledger.
func BalanceReflectOnLedger(l *ledger.Ledger, accountID string) error {
	logAction("BalanceReflectOnLedger", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Sync the account state with the ledger
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account in ledger: %w", err)
	}

	log.Printf("Successfully reflected balance state on ledger for account %s", accountID)
	return nil
}


// BalanceMapToExternalAccount links the balance with an external account.
func BalanceMapToExternalAccount(l *ledger.Ledger, accountID, externalID string) error {
	logAction("BalanceMapToExternalAccount", accountID)

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if externalID == "" {
		return errors.New("external account ID cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		return fmt.Errorf("internal account not found: %w", err)
	}

	// Link account to external ID
	account.ExternalAccountID = externalID

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("Successfully linked account %s to external account %s", accountID, externalID)
	return nil
}
