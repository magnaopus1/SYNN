package account_and_balance_operations

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"

	"github.com/google/uuid"
)

// BalanceWithdraw withdraws a specified amount from an account.
func BalanceWithdraw(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceWithdraw - Start", fmt.Sprintf("AccountID: %s, Amount: %.2f", accountID, amount))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceWithdraw - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceWithdraw - Validation Failed: Withdrawal Amount Must Be Positive")
		return errors.New("withdrawal amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceWithdraw - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		log.Printf("BalanceWithdraw - Insufficient Funds: Account Balance %.2f, Requested %.2f", account.Balance, amount)
		return errors.New("insufficient funds for withdrawal")
	}

	account.Balance -= amount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceWithdraw - Update Account Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceWithdraw - Success: %.2f Withdrawn from AccountID %s", amount, accountID)
	return nil
}


// BalanceDeposit deposits a specified amount into an account.
func BalanceDeposit(l *ledger.Ledger, accountID string, amount float64) error {
	logAction("BalanceDeposit - Start", fmt.Sprintf("AccountID: %s, Amount: %.2f", accountID, amount))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceDeposit - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceDeposit - Validation Failed: Deposit Amount Must Be Positive")
		return errors.New("deposit amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceDeposit - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	account.Balance += amount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceDeposit - Update Account Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceDeposit - Success: %.2f Deposited to AccountID %s", amount, accountID)
	return nil
}

// BalanceRequestRefund requests a refund for a specific transaction.
func BalanceRequestRefund(l *ledger.Ledger, accountID, transactionID string) error {
	logAction("BalanceRequestRefund - Start", fmt.Sprintf("AccountID: %s, TransactionID: %s", accountID, transactionID))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceRequestRefund - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if transactionID == "" {
		log.Printf("BalanceRequestRefund - Validation Failed: Transaction ID Cannot Be Empty")
		return errors.New("transaction ID cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	transaction, err := l.BlockchainConsensusCoinLedger.GetTransaction(transactionID)
	if err != nil {
		log.Printf("BalanceRequestRefund - Transaction Not Found: %v", err)
		return fmt.Errorf("transaction not found: %w", err)
	}
	if transaction.Status != "Completed" {
		log.Printf("BalanceRequestRefund - Refund Not Allowed: Transaction Status %s", transaction.Status)
		return errors.New("only completed transactions can be refunded")
	}

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceRequestRefund - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	account.Balance += transaction.Amount
	transaction.Status = "Refunded"

	if err := l.BlockchainConsensusCoinLedger.UpdateTransaction(transactionID, *transaction); err != nil {
		log.Printf("BalanceRequestRefund - Update Transaction Failed: %v", err)
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceRequestRefund - Update Account Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceRequestRefund - Success: Refunded %.2f to AccountID %s for TransactionID %s", transaction.Amount, accountID, transactionID)
	return nil
}


// BalanceAllocate allocates a specified amount for a project or purpose.
func BalanceAllocate(l *ledger.Ledger, accountID string, amount float64, purpose string) error {
	logAction("BalanceAllocate - Start", fmt.Sprintf("AccountID: %s, Amount: %.2f, Purpose: %s", accountID, amount, purpose))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceAllocate - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if amount <= 0 {
		log.Printf("BalanceAllocate - Validation Failed: Allocation Amount Must Be Positive")
		return errors.New("allocation amount must be positive")
	}
	if purpose == "" {
		log.Printf("BalanceAllocate - Validation Failed: Purpose Cannot Be Empty")
		return errors.New("purpose cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceAllocate - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < amount {
		log.Printf("BalanceAllocate - Insufficient Funds: Account Balance %.2f, Requested %.2f", account.Balance, amount)
		return errors.New("insufficient funds for allocation")
	}

	allocation := ledger.Allocation{
		ID:           GenerateUniqueID(),
		AccountID:    accountID,
		ResourceType: "funds",
		AllocatedAt:  time.Now(),
		Amount:       amount,
		Status:       "active",
		Remarks:      purpose,
		Allocated:    true,
	}

	account.Balance -= amount
	account.Allocations = append(account.Allocations, allocation)

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceAllocate - Update Account Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceAllocate - Success: Allocated %.2f from AccountID %s for Purpose: %s", amount, accountID, purpose)
	return nil
}

// BalanceRedeem redeems an allocation and releases funds back to the balance.
func BalanceRedeem(l *ledger.Ledger, accountID, allocationID string) error {
	logAction("BalanceRedeem - Start", fmt.Sprintf("AccountID: %s, AllocationID: %s", accountID, allocationID))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceRedeem - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if allocationID == "" {
		log.Printf("BalanceRedeem - Validation Failed: Allocation ID Cannot Be Empty")
		return errors.New("allocation ID cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceRedeem - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	var redeemedAllocation ledger.Allocation
	found := false
	for i, allocation := range account.Allocations {
		if allocation.ID == allocationID && allocation.Allocated {
			redeemedAllocation = allocation
			account.Allocations[i].Allocated = false
			found = true
			break
		}
	}

	if !found {
		log.Printf("BalanceRedeem - Allocation Not Found or Already Redeemed: AllocationID %s", allocationID)
		return errors.New("allocation not found or already redeemed")
	}

	account.Balance += redeemedAllocation.Amount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceRedeem - Update Account Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceRedeem - Success: Redeemed %.2f to AccountID %s from AllocationID %s", redeemedAllocation.Amount, accountID, allocationID)
	return nil
}


// BalanceVerifyChecksum verifies the integrity of an account’s balance.
func BalanceVerifyChecksum(l *ledger.Ledger, accountID, expectedChecksum string) error {
	logAction("BalanceVerifyChecksum - Start", fmt.Sprintf("AccountID: %s", accountID))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceVerifyChecksum - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if expectedChecksum == "" {
		log.Printf("BalanceVerifyChecksum - Validation Failed: Expected Checksum Cannot Be Empty")
		return errors.New("expected checksum cannot be empty")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceVerifyChecksum - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	checksum := GenerateChecksum(fmt.Sprintf("%f", account.Balance))
	if checksum != expectedChecksum {
		log.Printf("BalanceVerifyChecksum - Checksum Mismatch: Expected %s, Got %s", expectedChecksum, checksum)
		return errors.New("checksum mismatch, data integrity compromised")
	}

	log.Printf("BalanceVerifyChecksum - Success: Checksum Verified for AccountID %s", accountID)
	return nil
}


// GenerateChecksum creates a SHA-256 checksum for a given string.
func GenerateChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}


// GenerateUniqueID creates a unique identifier as a string.
func GenerateUniqueID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String())
}

// BalanceValidate ensures an account meets balance requirements.
func BalanceValidate(l *ledger.Ledger, accountID string, minimumRequired float64) (bool, error) {
	logAction("BalanceValidate - Start", fmt.Sprintf("AccountID: %s, MinimumRequired: %.2f", accountID, minimumRequired))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceValidate - Validation Failed: %v", err)
		return false, fmt.Errorf("invalid account ID: %w", err)
	}

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceValidate - Account Not Found: %v", err)
		return false, fmt.Errorf("account not found: %w", err)
	}

	isValid := account.Balance >= minimumRequired
	log.Printf("BalanceValidate - Result: AccountID %s has Valid Balance: %t", accountID, isValid)
	return isValid, nil
}


// BalanceAdjustForFee deducts a transaction fee from an account balance.
func BalanceAdjustForFee(l *ledger.Ledger, accountID string, feeAmount float64) error {
	logAction("BalanceAdjustForFee - Start", fmt.Sprintf("AccountID: %s, FeeAmount: %.2f", accountID, feeAmount))

	// Validate inputs
	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceAdjustForFee - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if feeAmount <= 0 {
		log.Printf("BalanceAdjustForFee - Validation Failed: Fee Amount Must Be Positive")
		return errors.New("fee amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceAdjustForFee - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}
	if account.Balance < feeAmount {
		log.Printf("BalanceAdjustForFee - Insufficient Funds: Account Balance %.2f, Fee Amount %.2f", account.Balance, feeAmount)
		return errors.New("insufficient balance to cover the fee")
	}

	account.Balance -= feeAmount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceAdjustForFee - Update Account Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceAdjustForFee - Success: Deducted Fee %.2f from AccountID %s", feeAmount, accountID)
	return nil
}


// BalanceAdjustForReward adds a reward to an account balance.
func BalanceAdjustForReward(l *ledger.Ledger, accountID string, rewardAmount float64) error {
	logAction("BalanceAdjustForReward - Start", fmt.Sprintf("AccountID: %s, RewardAmount: %.2f", accountID, rewardAmount))

	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceAdjustForReward - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}
	if rewardAmount <= 0 {
		log.Printf("BalanceAdjustForReward - Validation Failed: Reward amount must be positive")
		return errors.New("reward amount must be positive")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceAdjustForReward - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	account.Balance += rewardAmount

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceAdjustForReward - Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceAdjustForReward - Success: Reward of %.2f added to AccountID %s", rewardAmount, accountID)
	return nil
}


// BalanceDistributeProfit distributes profit among multiple accounts.
func BalanceDistributeProfit(l *ledger.Ledger, accountIDs []string, totalProfit float64) error {
	logAction("BalanceDistributeProfit - Start", fmt.Sprintf("TotalProfit: %.2f, Accounts: %v", totalProfit, accountIDs))

	if len(accountIDs) == 0 {
		log.Printf("BalanceDistributeProfit - Validation Failed: No accounts specified")
		return errors.New("no accounts specified for profit distribution")
	}
	if totalProfit <= 0 {
		log.Printf("BalanceDistributeProfit - Validation Failed: Total profit must be positive")
		return errors.New("total profit must be positive")
	}

	profitPerAccount := totalProfit / float64(len(accountIDs))
	for _, accountID := range accountIDs {
		if err := validateAccountID(accountID); err != nil {
			log.Printf("BalanceDistributeProfit - Validation Failed for AccountID %s: %v", accountID, err)
			return fmt.Errorf("invalid account ID: %w", err)
		}

		accountMutex.Lock()
		account, err := l.AccountsWalletLedger.GetAccount(accountID)
		if err != nil {
			accountMutex.Unlock()
			log.Printf("BalanceDistributeProfit - Account Not Found: %s", accountID)
			return fmt.Errorf("account not found: %w", err)
		}

		account.Balance += profitPerAccount
		if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
			accountMutex.Unlock()
			log.Printf("BalanceDistributeProfit - Update Failed for AccountID %s: %v", accountID, err)
			return err
		}
		accountMutex.Unlock()
	}

	log.Printf("BalanceDistributeProfit - Success: Distributed %.2f profit to %d accounts", totalProfit, len(accountIDs))
	return nil
}


// BalanceFinalize finalizes an account transaction, marking it as complete.
func BalanceFinalize(l *ledger.Ledger, transactionID string) error {
	logAction("BalanceFinalize - Start", fmt.Sprintf("TransactionID: %s", transactionID))

	transaction, err := l.BlockchainConsensusCoinLedger.GetTransaction(transactionID)
	if err != nil {
		log.Printf("BalanceFinalize - Transaction Not Found: %v", err)
		return fmt.Errorf("transaction not found: %w", err)
	}
	if transaction.Status == "Completed" {
		log.Printf("BalanceFinalize - Transaction Already Finalized: %s", transactionID)
		return errors.New("transaction already finalized")
	}

	transaction.Status = "Completed"
	if err := l.BlockchainConsensusCoinLedger.UpdateTransaction(transactionID, *transaction); err != nil {
		log.Printf("BalanceFinalize - Update Failed: %v", err)
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	log.Printf("BalanceFinalize - Success: Transaction %s marked as Completed", transactionID)
	return nil
}


// BalanceRollback reverts a transaction and adjusts balances accordingly.
func BalanceRollback(l *ledger.Ledger, transactionID string) error {
	logAction("BalanceRollback - Start", fmt.Sprintf("TransactionID: %s", transactionID))

	transaction, err := l.BlockchainConsensusCoinLedger.GetTransaction(transactionID)
	if err != nil {
		log.Printf("BalanceRollback - Transaction Not Found: %v", err)
		return fmt.Errorf("transaction not found: %w", err)
	}
	if transaction.Status != "Completed" {
		log.Printf("BalanceRollback - Invalid Status for Rollback: %s", transaction.Status)
		return errors.New("only completed transactions can be rolled back")
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	fromAccount, err := l.AccountsWalletLedger.GetAccount(transaction.FromAddress)
	if err != nil {
		log.Printf("BalanceRollback - Source Account Not Found: %v", err)
		return fmt.Errorf("source account not found: %w", err)
	}
	toAccount, err := l.AccountsWalletLedger.GetAccount(transaction.ToAddress)
	if err != nil {
		log.Printf("BalanceRollback - Destination Account Not Found: %v", err)
		return fmt.Errorf("destination account not found: %w", err)
	}

	// Revert the transaction
	fromAccount.Balance += transaction.Amount
	toAccount.Balance -= transaction.Amount
	transaction.Status = "Rolled Back"
	transaction.RefundAmount = transaction.Amount

	if err := l.BlockchainConsensusCoinLedger.UpdateTransaction(transactionID, *transaction); err != nil {
		log.Printf("BalanceRollback - Transaction Update Failed: %v", err)
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(transaction.FromAddress, *fromAccount); err != nil {
		log.Printf("BalanceRollback - Source Account Update Failed: %v", err)
		return fmt.Errorf("failed to update source account: %w", err)
	}
	if err := l.AccountsWalletLedger.UpdateAccount(transaction.ToAddress, *toAccount); err != nil {
		log.Printf("BalanceRollback - Destination Account Update Failed: %v", err)
		return fmt.Errorf("failed to update destination account: %w", err)
	}

	log.Printf("BalanceRollback - Success: Transaction %s rolled back", transactionID)
	return nil
}


// BalanceReview marks an account for balance review, possibly flagging it.
func BalanceReview(l *ledger.Ledger, accountID string) error {
	logAction("BalanceReview - Start", fmt.Sprintf("AccountID: %s", accountID))

	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceReview - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceReview - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}

	account.RequiresReview = true

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceReview - Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceReview - Success: AccountID %s marked for review", accountID)
	return nil
}


// BalanceConfirmation confirms the validity of a balance after a review.
func BalanceConfirmation(l *ledger.Ledger, accountID string) error {
	logAction("BalanceConfirmation - Start", fmt.Sprintf("AccountID: %s", accountID))

	if err := validateAccountID(accountID); err != nil {
		log.Printf("BalanceConfirmation - Validation Failed: %v", err)
		return fmt.Errorf("invalid account ID: %w", err)
	}

	accountMutex.Lock()
	defer accountMutex.Unlock()

	account, err := l.AccountsWalletLedger.GetAccount(accountID)
	if err != nil {
		log.Printf("BalanceConfirmation - Account Not Found: %v", err)
		return fmt.Errorf("account not found: %w", err)
	}
	if !account.RequiresReview {
		log.Printf("BalanceConfirmation - No Review Required for AccountID: %s", accountID)
		return errors.New("no review required for this account")
	}

	account.RequiresReview = false

	if err := l.AccountsWalletLedger.UpdateAccount(accountID, *account); err != nil {
		log.Printf("BalanceConfirmation - Update Failed: %v", err)
		return fmt.Errorf("failed to update account: %w", err)
	}

	log.Printf("BalanceConfirmation - Success: AccountID %s review confirmed", accountID)
	return nil
}
