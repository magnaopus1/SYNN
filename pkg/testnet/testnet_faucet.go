package testnet

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewTestnetFaucet initializes a new Testnet faucet with unlimited Synn available
func NewTestnetFaucet(ledgerInstance *ledger.Ledger) *common.TestnetFaucet {
    return &common.TestnetFaucet{
        Balance:        -1, // Testnet faucet has "unlimited" balance, no strict balance tracking
        Claims:         make(map[string]time.Time),
        LedgerInstance: ledgerInstance,
    }
}

// Claim allows a user to claim Testnet Synn from the faucet
func (f *common.TestnetFaucet) Claim(walletAddress string) (float64, error) {
    f.mutex.Lock()
    defer f.mutex.Unlock()

    // Enforce cooldown period (if needed)
    now := time.Now()
    if lastClaim, exists := f.Claims[walletAddress]; exists {
        if now.Sub(lastClaim) < ClaimCooldown {
            return 0, fmt.Errorf("claim cooldown period not yet over for wallet %s", walletAddress)
        }
    }

    // Update claim timestamp
    f.Claims[walletAddress] = now

    // Record the claim in the ledger
    err := f.LedgerInstance.RecordTestnetFaucetClaim(walletAddress, TestnetFaucetAmount)
    if err != nil {
        return 0, fmt.Errorf("failed to record testnet faucet claim in ledger: %v", err)
    }

    fmt.Printf("Wallet %s claimed %.2f Testnet SYNN from the faucet.\n", walletAddress, TestnetFaucetAmount)
    return TestnetFaucetAmount, nil
}

// Ledger integration

// RecordTestnetFaucetClaim logs a claim from the testnet faucet in the ledger
func (l *ledger.Ledger) RecordTestnetFaucetClaim(walletAddress string, amount float64) error {
    encryptedData, err := encryption.EncryptData(fmt.Sprintf("Claim %.2f Testnet SYNN by %s from testnet faucet", amount, walletAddress), common.EncryptionKey)
    if err != nil {
        return err
    }
    return l.RecordTransaction("TestnetFaucet", walletAddress, amount, encryptedData)
}
