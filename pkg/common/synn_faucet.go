package common

import (
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

const (
    FaucetAmount       = 0.01       // Amount of Synn each user can claim from the faucet per request
    ClaimLimit         = 3          // Maximum claims allowed per wallet in 7 days
    CooldownPeriod     = 7 * 24 * time.Hour // Cooldown period for claim limit reset
)

// Faucet represents the Synn faucet
type Faucet struct {
	OwnerAddress   string              // Owner of the faucet (blockchain owner)
	Balance        float64             // Current balance of the faucet
	Claims         map[string][]time.Time // Claims by wallet address with timestamps
	mutex          sync.Mutex          // Mutex for thread-safe operations
	LedgerInstance *ledger.Ledger      // Ledger instance for tracking faucet transactions
}


// NewFaucet initializes the Synn faucet with the owner's wallet address
func NewFaucet(ownerAddress string, initialBalance float64, ledgerInstance *ledger.Ledger) *Faucet {
    return &Faucet{
        OwnerAddress:   ownerAddress,
        Balance:        initialBalance,
        Claims:         make(map[string][]time.Time),
        LedgerInstance: ledgerInstance,
    }
}

// Deposit allows the owner to deposit funds into the faucet
func (f *Faucet) Deposit(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("invalid deposit amount")
    }

    f.mutex.Lock()
    defer f.mutex.Unlock()

    f.Balance += amount

    // Record the deposit in the ledger, passing only the amount
    err := f.LedgerInstance.BlockchainConsensusCoinLedger.RecordFaucetDeposit(amount)
    if err != nil {
        return fmt.Errorf("failed to record deposit in ledger: %v", err)
    }

    fmt.Printf("Owner %s deposited %.2f SYNN to the faucet.\n", f.OwnerAddress, amount)
    return nil
}


// Claim allows a user to claim Synn from the faucet
func (f *Faucet) Claim(walletAddress string) (float64, error) {
    f.mutex.Lock()
    defer f.mutex.Unlock()

    if f.Balance < FaucetAmount {
        return 0, errors.New("faucet balance is too low to process the claim")
    }

    // Check claim limits for the wallet address
    now := time.Now()
    claims, exists := f.Claims[walletAddress]
    if exists {
        claims = filterOldClaims(claims, now)
        if len(claims) >= ClaimLimit {
            return 0, fmt.Errorf("claim limit reached for wallet %s, try again after the cooldown period", walletAddress)
        }
    } else {
        claims = []time.Time{}
    }

    // Update faucet balance and claim history
    f.Balance -= FaucetAmount
    claims = append(claims, now)
    f.Claims[walletAddress] = claims

    // Record the claim in the ledger
    err := f.LedgerInstance.BlockchainConsensusCoinLedger.RecordFaucetClaim(walletAddress, FaucetAmount)
    if err != nil {
        return 0, fmt.Errorf("failed to record faucet claim in ledger: %v", err)
    }

    fmt.Printf("Wallet %s claimed %.2f SYNN from the faucet.\n", walletAddress, FaucetAmount)
    return FaucetAmount, nil
}

// filterOldClaims filters out claims older than the cooldown period (7 days)
func filterOldClaims(claims []time.Time, now time.Time) []time.Time {
    filteredClaims := []time.Time{}
    for _, claimTime := range claims {
        if now.Sub(claimTime) <= CooldownPeriod {
            filteredClaims = append(filteredClaims, claimTime)
        }
    }
    return filteredClaims
}

// FaucetBalance returns the current balance of the faucet
func (f *Faucet) FaucetBalance() float64 {
    f.mutex.Lock()
    defer f.mutex.Unlock()
    return f.Balance
}





