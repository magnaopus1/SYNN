package wallet

import (
	"encoding/json"
	"fmt"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewWalletDisplayService initializes the WalletDisplayService.
func NewWalletDisplayService(walletID string, ledgerInstance *ledger.Ledger) *WalletDisplayService {
	return &WalletDisplayService{
		walletID:       walletID,
		ledgerInstance: ledgerInstance,
	}
}

// DisplayBalances shows the balances of all supported currencies and tokens.
func (wds *WalletDisplayService) DisplayBalances(balances map[string]float64) {
	wds.mutex.Lock()
	defer wds.mutex.Unlock()

	fmt.Printf("Balances for Wallet ID: %s\n", wds.walletID)
	for symbol, balance := range balances {
		fmt.Printf("Currency/Token: %s, Balance: %f\n", symbol, balance)
	}
}

// DisplayTransactionHistory shows the transaction history of the wallet.
func (wds *WalletDisplayService) DisplayTransactionHistory() error {
    wds.mutex.Lock()
    defer wds.mutex.Unlock()

    fmt.Printf("Transaction History for Wallet ID: %s\n", wds.walletID)

    // Fetch the transaction history from the ledger (adjust based on actual return type)
    history, err := wds.ledgerInstance.GetTransactionHistory() // Removed walletID from arguments

    // Handle any errors that occur during fetching
    if err != nil {
        return fmt.Errorf("failed to get transaction history: %v", err)
    }

    // Iterate over the transaction history and print it
    for _, tx := range history {
        fmt.Printf("Transaction ID: %s, Amount: %.2f, From: %s, To: %s, Timestamp: %s\n", 
            tx.ID, tx.Amount, tx.From, tx.To, tx.Timestamp.String())
    }

    return nil
}







// DisplaySupportedCurrencies shows all supported currencies and tokens in the wallet.
func (wds *WalletDisplayService) DisplaySupportedCurrencies(balances map[string]float64) {
	wds.mutex.Lock()
	defer wds.mutex.Unlock()

	fmt.Printf("Supported Currencies and Tokens for Wallet ID: %s\n", wds.walletID)
	for symbol := range balances {
		fmt.Printf("Currency/Token: %s\n", symbol)
	}
}

// DisplayContractExecution displays the details of executed smart contracts associated with the wallet.
func (wds *WalletDisplayService) DisplayContractExecution() error {
	wds.mutex.Lock()
	defer wds.mutex.Unlock()

	fmt.Printf("Smart Contract Execution for Wallet ID: %s\n", wds.walletID)

	// Fetch contract execution history from the ledger (assuming this returns a struct or other serializable data)
	contractHistory, err := wds.ledgerInstance.GetContractExecutionHistory(wds.walletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve contract execution history: %v", err)
	}

	// If contractHistory is not []byte, serialize it. Assuming it's a struct or map.
	contractHistoryBytes, err := json.Marshal(contractHistory)
	if err != nil {
		return fmt.Errorf("failed to serialize contract execution history: %v", err)
	}

	// Create an encryption instance
	encryptionInstance := &common.Encryption{}

	// Decrypt contract execution history before display
	decryptedHistoryBytes, err := encryptionInstance.DecryptData(contractHistoryBytes, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt contract execution history: %v", err)
	}

	// Convert the decrypted bytes to a readable format (assuming it's a string)
	decryptedHistory := string(decryptedHistoryBytes)

	// Display the decrypted contract execution history
	fmt.Printf("Decrypted Contract Execution History: %s\n", decryptedHistory)

	return nil
}




// DisplayTokenMints shows the token minting history of the wallet.
func (wds *WalletDisplayService) DisplayTokenMints() error {
	wds.mutex.Lock()
	defer wds.mutex.Unlock()

	fmt.Printf("Token Minting History for Wallet ID: %s\n", wds.walletID)

	// Fetch token minting history from the ledger
	tokenMints, err := wds.ledgerInstance.GetTokenMintHistory(wds.walletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token minting history: %v", err)
	}

	// If tokenMints is not []byte, serialize it. Assuming it's a struct, map, or similar.
	tokenMintsBytes, err := json.Marshal(tokenMints)
	if err != nil {
		return fmt.Errorf("failed to serialize token minting history: %v", err)
	}

	// Create an encryption instance
	encryptionInstance := &common.Encryption{}

	// Decrypt token minting history before display
	decryptedMintsBytes, err := encryptionInstance.DecryptData(tokenMintsBytes, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt token minting history: %v", err)
	}

	// Convert the decrypted bytes to a readable format (assuming it's a string or JSON)
	decryptedMints := string(decryptedMintsBytes)

	// Display the decrypted token minting history
	fmt.Printf("Decrypted Token Minting History: %s\n", decryptedMints)

	return nil
}



// DisplayTokenBurns shows the token burning history of the wallet.
func (wds *WalletDisplayService) DisplayTokenBurns() error {
	wds.mutex.Lock()
	defer wds.mutex.Unlock()

	fmt.Printf("Token Burning History for Wallet ID: %s\n", wds.walletID)

	// Fetch token burning history from the ledger
	tokenBurns, err := wds.ledgerInstance.GetTokenBurnHistory(wds.walletID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token burning history: %v", err)
	}

	// Serialize tokenBurns if it's not already a []byte type
	tokenBurnsBytes, err := json.Marshal(tokenBurns)
	if err != nil {
		return fmt.Errorf("failed to serialize token burning history: %v", err)
	}

	// Create an encryption instance
	encryptionInstance := &common.Encryption{}

	// Decrypt token burning history before display
	decryptedBurnsBytes, err := encryptionInstance.DecryptData(tokenBurnsBytes, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt token burning history: %v", err)
	}

	// Convert the decrypted bytes to a readable format (assuming it's a string or JSON)
	decryptedBurns := string(decryptedBurnsBytes)

	// Display the decrypted token burning history
	fmt.Printf("Decrypted Token Burning History: %s\n", decryptedBurns)

	return nil
}



