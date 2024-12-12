package common

import (
    "fmt"
    "synnergy_network/pkg/ledger"
    "sync"
)

// Coin represents the structure of the Synnergy Network Coin
type SynthronCoin struct {
	Name              string  // Name of the coin
	Symbol            string  // Coin ticker symbol
	TotalSupply       float64 // Total supply of the coin
	CirculatingSupply float64 // Circulating supply (optional)
	Decimals          int     // Number of decimal places the coin supports
}


// InitializeSynthronCoin creates and initializes the SynthronCoin
func InitializeSynthronCoin() SynthronCoin {
    // Initialize the SynthronCoin inside a function
    synthronCoin := SynthronCoin{
        Name:              "Synthron Coin",
        Symbol:            "SYNN",
        TotalSupply:       500000000, // 500 million SYNN
        CirculatingSupply: 0,
        Decimals:          18,
    }
    return synthronCoin
}

// Mutex for thread-safe operations
var mutex = &sync.Mutex{} // Use pointer to sync.Mutex for better performance

// Initialize the SynthronCoin once for the package
var synthronCoin = InitializeSynthronCoin()

// InitSupply initializes the total supply of SYNN within the system and stores it in the ledger.
func InitSupply(ledgerInstance *ledger.Ledger) {
    mutex.Lock()
    defer mutex.Unlock()

    // Initialize the system with the total supply of the coin
    ledgerInstance.BlockchainConsensusCoinLedger.SetSystemBalance(synthronCoin.TotalSupply)
    fmt.Printf("System initialized with total supply of %f %s.\n", synthronCoin.TotalSupply, synthronCoin.Symbol)
}

// GetSystemBalance retrieves the total balance of SYNN in the system.
func GetSystemBalance(ledgerInstance *ledger.Ledger) float64 {
    mutex.Lock()
    defer mutex.Unlock()

    return ledgerInstance.BlockchainConsensusCoinLedger.GetSystemBalance()
}

// TransferCoins transfers coins between users in the system and records the transaction in the ledger.
func TransferCoins(from string, to string, amount float64, ledgerInstance *ledger.Ledger) error {
    mutex.Lock()
    defer mutex.Unlock()

    // Check if the sender has enough balance
    hasBalance, err := ledgerInstance.AccountsWalletLedger.HasSufficientBalance(from, amount) // Capture both bool and error
    if err != nil {
        return fmt.Errorf("error checking balance for user %s: %v", from, err)
    }
    if !hasBalance {
        return fmt.Errorf("insufficient balance for user %s", from)
    }

    // Perform the transfer and record the transaction in the ledger
    err = ledgerInstance.AccountsWalletLedger.DebitBalance(from, amount) // Use DebitBalance with error handling
    if err != nil {
        return fmt.Errorf("failed to debit balance for user %s: %v", from, err)
    }

    err = ledgerInstance.AccountsWalletLedger.CreditBalance(to, amount) // Use CreditBalance with error handling
    if err != nil {
        return fmt.Errorf("failed to credit balance for user %s: %v", to, err)
    }

    // Log the transaction in the ledger (passing correct arguments)
    err = ledgerInstance.BlockchainConsensusCoinLedger.AddTransaction(from, to, amount)
    if err != nil {
        return fmt.Errorf("failed to log transaction: %v", err)
    }

    fmt.Printf("Transferred %f %s from %s to %s.\n", amount, synthronCoin.Symbol, from, to)
    return nil
}

