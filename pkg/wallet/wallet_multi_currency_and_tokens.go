package wallet

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// WalletMultiCurrencyService handles multiple currencies and tokens in the wallet.
type WalletMultiCurrencyService struct {
	walletID        string
	balances        map[string]float64 // Key is currency/token symbol, value is balance
	ledgerInstance  *ledger.Ledger
	networkManager  *network.NetworkManager // Added network manager to the struct
	mutex           sync.Mutex
}

type CurrencyExchange struct {
	FromSymbol      string  // Currency being exchanged from
	ToSymbol        string  // Currency being exchanged to
	Amount          float64 // Amount of the fromSymbol being exchanged
	ExchangedAmount float64 // Resulting amount in the toSymbol after exchange
}


// NewWalletMultiCurrencyService initializes the wallet service with multi-currency support.
func NewWalletMultiCurrencyService(walletID string, ledgerInstance *ledger.Ledger, networkManager *network.NetworkManager) *WalletMultiCurrencyService {
	return &WalletMultiCurrencyService{
		walletID:       walletID,
		balances:       make(map[string]float64),
		ledgerInstance: ledgerInstance,
		networkManager: networkManager, // Added network manager initialization
	}
}


// AddCurrency adds a new currency or token to the wallet.
func (wms *WalletMultiCurrencyService) AddCurrency(symbol string, initialBalance float64) error {
	wms.mutex.Lock()
	defer wms.mutex.Unlock()

	if _, exists := wms.balances[symbol]; exists {
		return fmt.Errorf("currency or token %s already exists in wallet", symbol)
	}

	wms.balances[symbol] = initialBalance
	return nil
}

// GetBalance returns the balance for a given currency or token.
func (wms *WalletMultiCurrencyService) GetBalance(symbol string) (float64, error) {
	wms.mutex.Lock()
	defer wms.mutex.Unlock()

	balance, exists := wms.balances[symbol]
	if !exists {
		return 0, fmt.Errorf("currency or token %s does not exist in the wallet", symbol)
	}
	return balance, nil
}

// Transfer transfers a certain amount of currency/token from one wallet to another.
func (wms *WalletMultiCurrencyService) Transfer(toWalletID, symbol string, amount float64) error {
	wms.mutex.Lock()
	defer wms.mutex.Unlock()

	// Step 1: Check if wallet contains the currency or token
	balance, exists := wms.balances[symbol]
	if !exists {
		return fmt.Errorf("currency or token %s not available in wallet", symbol)
	}

	// Step 2: Check if there is enough balance
	if balance < amount {
		return fmt.Errorf("insufficient balance in wallet for %s: current balance is %f", symbol, balance)
	}

	// Step 3: Deduct from current wallet and update the ledger
	newBalance := balance - amount
	wms.balances[symbol] = newBalance

	// Step 4: Record the transfer in the ledger
	err := wms.ledgerInstance.RecordTransaction(wms.walletID, toWalletID, amount) // Remove symbol argument
	if err != nil {
		return fmt.Errorf("failed to record transaction in ledger: %v", err)
	}

	// Step 5: Encrypt and send the transaction over the network
	encryptionInstance := &common.Encryption{}  // Create an encryption instance
	encryptedTransactionID, err := encryptionInstance.EncryptData("AES", []byte(wms.walletID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction ID: %v", err)
	}

	// Use the NetworkManager to send the encrypted transaction
	err = wms.networkManager.SendEncryptedMessage(toWalletID, string(encryptedTransactionID))
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	return nil
}


// MintTokens mints new tokens and adds them to the wallet.
func (wms *WalletMultiCurrencyService) MintTokens(symbol string, amount float64) error {
    wms.mutex.Lock()
    defer wms.mutex.Unlock()

    if _, exists := wms.balances[symbol]; !exists {
        return fmt.Errorf("token %s does not exist in wallet", symbol)
    }

    wms.balances[symbol] += amount

    // Step 1: Create a MintRecord
    mintRecord := ledger.MintRecord{
        RecordID:   fmt.Sprintf("mint-%d", time.Now().UnixNano()), // Unique ID for the record
        TokenID:    symbol, // Assuming the symbol is the token ID
        Amount:     big.NewInt(int64(amount)), // Convert float64 to *big.Int
        MintedBy:   wms.walletID, // Address of the wallet minting the tokens
        Timestamp:  time.Now(), // Timestamp of the minting
    }

    // Step 2: Record minting event in the ledger
    err := wms.ledgerInstance.RecordTokenMint(wms.walletID, mintRecord) // Pass wallet ID and mint record
    if err != nil {
        return fmt.Errorf("failed to record token minting in ledger: %v", err)
    }

    return nil
}

// BurnTokens burns a specified amount of tokens from the wallet.
func (wms *WalletMultiCurrencyService) BurnTokens(symbol string, amount float64) error {
    wms.mutex.Lock()
    defer wms.mutex.Unlock()

    balance, exists := wms.balances[symbol]
    if !exists {
        return fmt.Errorf("token %s does not exist in wallet", symbol)
    }

    if balance < amount {
        return fmt.Errorf("insufficient tokens to burn: current balance is %f", balance)
    }

    wms.balances[symbol] -= amount

    // Step 1: Create a BurnRecord
    burnRecord := ledger.BurnRecord{
        RecordID:   fmt.Sprintf("burn-%d", time.Now().UnixNano()), // Unique ID for the record
        TokenID:    symbol, // Assuming the symbol is the token ID
        Amount:     big.NewInt(int64(amount)), // Convert float64 to *big.Int
        BurnedBy:   wms.walletID, // Address of the wallet burning the tokens
        Timestamp:  time.Now(), // Timestamp of the burning
    }

    // Step 2: Record burning event in the ledger
    err := wms.ledgerInstance.RecordTokenBurn(wms.walletID, burnRecord) // Pass wallet ID and burn record
    if err != nil {
        return fmt.Errorf("failed to record token burning in ledger: %v", err)
    }

    return nil
}




// GetSupportedCurrencies returns the list of all currencies and tokens supported in the wallet.
func (wms *WalletMultiCurrencyService) GetSupportedCurrencies() []string {
	wms.mutex.Lock()
	defer wms.mutex.Unlock()

	var supportedCurrencies []string
	for symbol := range wms.balances {
		supportedCurrencies = append(supportedCurrencies, symbol)
	}
	return supportedCurrencies
}

// ExchangeCurrency exchanges one currency for another (e.g., USD to BTC).
func (wms *WalletMultiCurrencyService) ExchangeCurrency(fromSymbol, toSymbol string, amount float64, exchangeRate float64) error {
    wms.mutex.Lock()
    defer wms.mutex.Unlock()

    // Check if the wallet supports both currencies
    fromBalance, existsFrom := wms.balances[fromSymbol]
    if !existsFrom {
        return errors.New("source currency does not exist in wallet")
    }
    _, existsTo := wms.balances[toSymbol]
    if !existsTo {
        return errors.New("target currency does not exist in wallet")
    }

    // Check if there are enough funds to exchange
    if fromBalance < amount {
        return fmt.Errorf("insufficient funds in %s to exchange: balance %f, required %f", fromSymbol, fromBalance, amount)
    }

    // Calculate the equivalent amount in the target currency
    exchangedAmount := amount * exchangeRate

    // Update wallet balances
    wms.balances[fromSymbol] -= amount
    wms.balances[toSymbol] += exchangedAmount

    // Create the CurrencyExchange struct
    exchange := ledger.CurrencyExchange{
        FromCurrency:    fromSymbol,
        ToCurrency:      toSymbol,
        Amount:          big.NewInt(int64(amount * 1e18)), // Convert float64 to *big.Int with scaling
        ExchangedAmount: big.NewInt(int64(exchangedAmount * 1e18)), // Convert float64 to *big.Int with scaling
        ExecutedAt:      time.Now(),  // Add a timestamp for the exchange record
    }

    // Record the exchange in the ledger
    err := wms.ledgerInstance.RecordCurrencyExchange(wms.walletID, exchange) // Ensure the method matches expected arguments
    if err != nil {
        return fmt.Errorf("failed to record currency exchange in ledger: %v", err)
    }

    return nil
}


// GetLedgerInstance returns the wallet's ledger instance.
func (wms *WalletMultiCurrencyService) GetLedgerInstance() *ledger.Ledger {
	return wms.ledgerInstance
}

