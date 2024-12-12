package token

import (
    "fmt"
    "sync"
    "synnergy_network_demo/ledger"
)

// TokenManager manages multiple token standards using the TokenInterface
type TokenManager struct {
    Tokens          map[string]TokenInterface // Token ID -> Token Interface
    LedgerInstance  *ledger.Ledger            // Ledger instance to store transactions
    mutex           sync.Mutex                // Mutex for thread-safe operations
}

// NewTokenManager initializes a new TokenManager
func NewTokenManager(ledgerInstance *ledger.Ledger) *TokenManager {
    return &TokenManager{
        Tokens:         make(map[string]TokenInterface),
        LedgerInstance: ledgerInstance,
    }
}

// DeployToken deploys a new token on the blockchain using a universal interface
func (tm *TokenManager) DeployToken(tokenID string, token TokenInterface) error {
    tm.mutex.Lock()
    defer tm.mutex.Unlock()

    tm.Tokens[tokenID] = token
    fmt.Printf("Token %s deployed.\n", tokenID)
    return nil
}

// Transfer tokens between addresses, based on the token standard
func (tm *TokenManager) Transfer(tokenID, from, to string, amount float64) error {
    tm.mutex.Lock()
    defer tm.mutex.Unlock()

    token, exists := tm.Tokens[tokenID]
    if !exists {
        return fmt.Errorf("token with ID %s not found", tokenID)
    }

    return token.Transfer(from, to, amount)
}

// Get balance of a token for a specific address
func (tm *TokenManager) BalanceOf(tokenID, address string) (float64, error) {
    tm.mutex.Lock()
    defer tm.mutex.Unlock()

    token, exists := tm.Tokens[tokenID]
    if !exists {
        return 0, fmt.Errorf("token with ID %s not found", tokenID)
    }

    return token.BalanceOf(address)
}

// Validate ensures the integrity of all token balances in the system
func (tm *TokenManager) Validate(tokenID string) error {
    tm.mutex.Lock()
    defer tm.mutex.Unlock()

    token, exists := tm.Tokens[tokenID]
    if !exists {
        return fmt.Errorf("token with ID %s not found", tokenID)
    }

    return token.Validate()
}
