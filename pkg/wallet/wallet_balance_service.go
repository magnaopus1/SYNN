package wallet

import (
	"fmt"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewWalletBalanceService initializes the WalletBalanceService.
func NewWalletBalanceService(walletID string, ledgerInstance *ledger.Ledger) *WalletBalanceService {
	return &WalletBalanceService{
		walletID:       walletID,
		ledgerInstance: ledgerInstance,
		balances:       make(map[string]float64),
	}
}

// GetBalance retrieves the current balance of a specific currency or token.
func (wbs *WalletBalanceService) GetBalance(currency string) (float64, error) {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Fetch the balance as a float64 directly from the ledger, no need for encryption.
	encryptedBalance, err := wbs.ledgerInstance.GetBalance(wbs.walletID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve balance for wallet %s: %v", wbs.walletID, err)
	}

	// Since encryptedBalance is a float64, you need to handle it directly without decryption
	return encryptedBalance, nil
}



// UpdateBalance updates the wallet's balance for a specific currency or token.
func (wbs *WalletBalanceService) UpdateBalance(currency string, amount float64) error {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Update the balance in memory
	wbs.balances[currency] += amount

	// Convert the updated balance to uint64 for storage
	updatedBalance := uint64(wbs.balances[currency] * 1000000000) // Example: convert float64 to a fixed-point uint64

	// Store the updated balance in the ledger using only walletID and updatedBalance
	err := wbs.ledgerInstance.UpdateBalance(wbs.walletID, updatedBalance) // Pass only walletID and updatedBalance
	if err != nil {
		return fmt.Errorf("failed to update balance for currency %s in ledger: %v", currency, err)
	}

	return nil
}



// FetchAllBalances retrieves and decrypts all balances for the wallet from the ledger.
func (wbs *WalletBalanceService) FetchAllBalances() (map[string]float64, error) {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Fetch encrypted balances from the ledger (assuming they are stored as strings)
	encryptedBalances, err := wbs.ledgerInstance.GetAllBalances()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve balances for wallet: %v", err)
	}

	encryptionInstance := &common.Encryption{} // Create an instance of Encryption

	// Decrypt each balance and store it in the in-memory balances map
	for currency, encryptedBalance := range encryptedBalances {
		// Convert the float64 balance to string before decryption
		encryptedBalanceStr := fmt.Sprintf("%f", encryptedBalance)

		// Convert encryptedBalanceStr to []byte (assuming it's stored as a string)
		encryptedBalanceBytes := []byte(encryptedBalanceStr)

		// Decrypt the balance
		decryptedBalanceBytes, err := encryptionInstance.DecryptData(encryptedBalanceBytes, common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt balance for currency %s: %v", currency, err)
		}

		// Convert decrypted data from string to float64
		var balance float64
		_, err = fmt.Sscanf(string(decryptedBalanceBytes), "%f", &balance)
		if err != nil {
			return nil, fmt.Errorf("failed to parse decrypted balance for currency %s: %v", currency, err)
		}

		wbs.balances[currency] = balance
	}

	return wbs.balances, nil
}




// TransferFunds transfers funds from the wallet to another wallet, adjusting balances accordingly.
func (wbs *WalletBalanceService) TransferFunds(destinationWalletID string, currency string, amount float64) error {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Check if sufficient balance is available in the sender's wallet
	currentBalance, err := wbs.GetBalance(currency)
	if err != nil {
		return fmt.Errorf("failed to retrieve balance for currency %s: %v", currency, err)
	}

	if currentBalance < amount {
		return fmt.Errorf("insufficient funds in %s to transfer %.2f", currency, amount)
	}

	// Deduct the balance from the sender's wallet
	wbs.balances[currency] -= amount
	if err := wbs.UpdateBalance(currency, wbs.balances[currency]); err != nil {
		return fmt.Errorf("failed to update balance after deduction for sender: %v", err)
	}

	// Assume we have a method to get the destination wallet's ledger instance
	destinationLedgerInstance, err := ledger.GetInstanceForWallet(destinationWalletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve ledger instance for destination wallet %s: %v", destinationWalletID, err)
	}

	// Load the destination wallet's balance service
	destinationWalletService := NewWalletBalanceService(destinationWalletID, destinationLedgerInstance)

	// Add the balance to the destination wallet
	err = destinationWalletService.UpdateBalance(currency, amount)
	if err != nil {
		return fmt.Errorf("failed to update balance for destination wallet %s: %v", destinationWalletID, err)
	}

	fmt.Printf("Successfully transferred %.2f %s to wallet %s\n", amount, currency, destinationWalletID)
	return nil
}


// AdjustBalance adjusts the balance by either increasing or decreasing based on the transaction type.
func (wbs *WalletBalanceService) AdjustBalance(amount float64, isCredit bool) error {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Adjust balance based on credit or debit
	if isCredit {
		wbs.balances["default"] += amount // Assuming "default" is the currency type or key in balances map
	} else {
		wbs.balances["default"] -= amount
	}

	// Convert the adjusted balance to uint64 to match the expected type in the ledger.
	adjustedBalance := uint64(wbs.balances["default"])

	// Update ledger (only pass walletID and adjusted balance)
	err := wbs.ledgerInstance.UpdateBalance(wbs.walletID, adjustedBalance)
	if err != nil {
		return fmt.Errorf("failed to update adjusted balance in ledger: %v", err)
	}

	return nil
}


