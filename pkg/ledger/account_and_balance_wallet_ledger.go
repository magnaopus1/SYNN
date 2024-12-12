package ledger

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"synnergy_network/pkg/ledger/token_ledgers"
	"time"
)

// RecordAccount adds a new account to the ledger.
func (l *AccountsWalletLedger) RecordAccount(accountID string, account Account) error {
    l.Lock()
    defer l.Unlock()

    // Validate input
    if accountID == "" {
        return fmt.Errorf("accountID cannot be empty")
    }

    // Check for duplicate account
    if _, exists := l.AccountsWalletLedgerState.Accounts[accountID]; exists {
        return fmt.Errorf("account %s already exists", accountID)
    }

    // Add account to ledger state
    l.AccountsWalletLedgerState.Accounts[accountID] = account

    // Log account creation in the ledger
    log.Printf("[INFO] Account %s added successfully", accountID)
    return nil
}


// DebitBalance subtracts a specified amount from an account balance.
func (l *AccountsWalletLedger) DebitBalance(accountID string, amount float64) error {
    l.Lock()
    defer l.Unlock()

    // Validate inputs
    if accountID == "" {
        return fmt.Errorf("accountID cannot be empty")
    }
    if amount <= 0 {
        return fmt.Errorf("debit amount must be greater than zero")
    }

    // Check if account exists
    account, exists := l.Balances[accountID]
    if !exists {
        return fmt.Errorf("account %s does not exist", accountID)
    }

    // Verify sufficient balance
    if account.Balance < amount {
        return fmt.Errorf("insufficient funds in account %s. Available: %.2f, Requested: %.2f", accountID, account.Balance, amount)
    }

    // Deduct the amount
    account.Balance -= amount
    l.Balances[accountID] = account

    // Log the transaction
    log.Printf("[INFO] Debited %.2f from account %s. Remaining balance: %.2f", amount, accountID, account.Balance)
    return nil
}


// GetVerifiedWallets retrieves wallets that hold the SYN900 token and are verified.
func (l *AccountsWalletLedger) GetVerifiedWallets(tokenID string) ([]WalletData, error) {
    // Validate input
    if tokenID == "" {
        return nil, fmt.Errorf("tokenID cannot be empty")
    }

    // Retrieve verified wallets
    var verifiedWallets []WalletData
    for _, wallet := range l.Wallets {
        if wallet.TokenID == tokenID && wallet.Verified {
            verifiedWallets = append(verifiedWallets, wallet)
        }
    }

    // Handle case where no wallets are found
    if len(verifiedWallets) == 0 {
        return nil, fmt.Errorf("no verified wallets found for token %s", tokenID)
    }

    // Log success
    log.Printf("[INFO] Found %d verified wallets for token %s", len(verifiedWallets), tokenID)
    return verifiedWallets, nil
}


// AdjustBalance modifies the balance of the given account by the specified amount.
func (l *AccountsWalletLedger) AdjustBalance(accountID string, amount float64) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return fmt.Errorf("accountID cannot be empty")
    }

    // Retrieve the account
    account, exists := l.Balances[accountID]
    if !exists {
        return fmt.Errorf("account %s not found", accountID)
    }

    // Adjust the account balance
    newBalance := account.Balance + amount
    if newBalance < 0 {
        return fmt.Errorf("insufficient funds in account %s. Current balance: %.2f, Attempted adjustment: %.2f", accountID, account.Balance, amount)
    }

    // Update the balance
    account.Balance = newBalance
    l.Balances[accountID] = account

    // Log the adjustment
    log.Printf("[INFO] Account %s balance adjusted by %.2f. New balance: %.2f", accountID, amount, newBalance)
    return nil
}


// GetTokenByWalletID retrieves the SYN900 token associated with the given walletID.
func (l *AccountsWalletLedger) GetTokenByWalletID(walletID string) (*tokenledgers.SYN900Token, error) {
    // Input validation
    if walletID == "" {
        return nil, fmt.Errorf("walletID cannot be empty")
    }

    // Check if the token exists
    token, exists := l.SYN900tokens[walletID]
    if !exists {
        return nil, fmt.Errorf("no SYN900 token found for wallet ID %s", walletID)
    }

    // Log retrieval success
    log.Printf("[INFO] Retrieved SYN900 token for wallet ID %s", walletID)
    return token, nil
}


// TransferFunds facilitates the transfer of funds between two accounts (float64 version).
func (l *AccountsWalletLedger) TransferFundsFloat(fromAccountID, toAccountID string, amount float64) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if fromAccountID == "" || toAccountID == "" {
        return fmt.Errorf("both fromAccountID and toAccountID must be provided")
    }
    if amount <= 0 {
        return fmt.Errorf("transfer amount must be greater than zero")
    }

    // Retrieve source account
    fromAccount, exists := l.AccountsWalletLedgerState.Accounts[fromAccountID]
    if !exists {
        return fmt.Errorf("source account %s not found", fromAccountID)
    }

    // Retrieve destination account
    toAccount, exists := l.AccountsWalletLedgerState.Accounts[toAccountID]
    if !exists {
        return fmt.Errorf("destination account %s not found", toAccountID)
    }

    // Check for sufficient funds
    if fromAccount.Balance < amount {
        return fmt.Errorf("insufficient funds in source account %s. Available: %.2f, Requested: %.2f", fromAccountID, fromAccount.Balance, amount)
    }

    // Perform the transfer
    fromAccount.Balance -= amount
    toAccount.Balance += amount
    l.AccountsWalletLedgerState.Accounts[fromAccountID] = fromAccount
    l.AccountsWalletLedgerState.Accounts[toAccountID] = toAccount

    // Log the transfer
    log.Printf("[INFO] Transferred %.2f from account %s to account %s. Source balance: %.2f, Destination balance: %.2f", amount, fromAccountID, toAccountID, fromAccount.Balance, toAccount.Balance)
    return nil
}


// RecordValidatorStake updates the validator's stake in the ledger.
func (l *AccountsWalletLedger) RecordValidatorStake(validatorID string, amount float64) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if validatorID == "" {
        return fmt.Errorf("validatorID cannot be empty")
    }
    if amount <= 0 {
        return fmt.Errorf("stake amount must be greater than zero")
    }

    // Check if the validator account exists
    account, exists := l.AccountsWalletLedgerState.Accounts[validatorID]
    if !exists {
        return fmt.Errorf("validator account %s does not exist", validatorID)
    }

    // Update the stake
    account.Stake += amount
    l.AccountsWalletLedgerState.Accounts[validatorID] = account

    // Log the update
    log.Printf("[INFO] Validator %s stake updated. New stake: %.2f", validatorID, account.Stake)
    return nil
}


// CreditBalance adds a specified amount to an account balance.
func (l *AccountsWalletLedger) CreditBalance(accountID string, amount float64) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return fmt.Errorf("accountID cannot be empty")
    }
    if amount <= 0 {
        return fmt.Errorf("credit amount must be greater than zero")
    }

    // Retrieve the account
    account, exists := l.Balances[accountID]
    if !exists {
        return fmt.Errorf("account %s does not exist", accountID)
    }

    // Update the balance
    account.Balance += amount
    l.Balances[accountID] = account

    // Log the credit operation
    log.Printf("[INFO] Account %s credited with %.2f. New balance: %.2f", accountID, amount, account.Balance)
    return nil
}


// HasSufficientBalance checks whether an account has sufficient balance for a transaction.
func (l *AccountsWalletLedger) HasSufficientBalance(accountID string, amount float64) (bool, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return false, fmt.Errorf("accountID cannot be empty")
    }
    if amount <= 0 {
        return false, fmt.Errorf("amount to check must be greater than zero")
    }

    // Retrieve the account
    account, exists := l.Balances[accountID]
    if !exists {
        return false, fmt.Errorf("account %s does not exist", accountID)
    }

    // Check balance sufficiency
    if account.Balance < amount {
        return false, fmt.Errorf("account %s has insufficient balance. Available: %.2f, Required: %.2f", accountID, account.Balance, amount)
    }

    log.Printf("[INFO] Account %s has sufficient balance for the transaction. Available: %.2f, Required: %.2f", accountID, account.Balance, amount)
    return true, nil
}


// RecordMultiSigWallet creates a new multi-signature wallet record
func (l *AccountsWalletLedger) RecordMultiSigWallet(walletID string, owners []string, requiredSigs int) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if len(owners) == 0 {
        return fmt.Errorf("owners list cannot be empty")
    }
    if requiredSigs <= 0 || requiredSigs > len(owners) {
        return fmt.Errorf("requiredSigs must be between 1 and the number of owners")
    }

    // Check if the wallet already exists
    if _, exists := l.MultiSigWallets[walletID]; exists {
        return fmt.Errorf("multi-signature wallet %s already exists", walletID)
    }

    // Create the multi-signature wallet
    l.MultiSigWallets[walletID] = MultiSigWallet{
        WalletID:    walletID,
        Owners:      owners,
        RequiredSigs: requiredSigs,
        CreatedAt:   time.Now(),
    }

    // Log the creation
    log.Printf("[INFO] MultiSig wallet %s created with %d required signatures.", walletID, requiredSigs)
    return nil
}


// FetchMultiSigWallet retrieves a multi-signature wallet's details
func (l *AccountsWalletLedger) FetchMultiSigWallet(walletID string) (MultiSigWallet, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return MultiSigWallet{}, fmt.Errorf("walletID cannot be empty")
    }

    // Fetch the wallet details
    wallet, exists := l.MultiSigWallets[walletID]
    if !exists {
        return MultiSigWallet{}, fmt.Errorf("multi-signature wallet %s not found", walletID)
    }

    // Log the retrieval
    log.Printf("[INFO] MultiSig wallet %s retrieved successfully.", walletID)
    return wallet, nil
}


// RevokeSignature removes a signature from a multi-signature wallet
func (l *AccountsWalletLedger) RevokeSignature(walletID, signerID string) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if signerID == "" {
        return fmt.Errorf("signerID cannot be empty")
    }

    // Check if the wallet exists
    wallet, exists := l.MultiSigWallets[walletID]
    if !exists {
        return fmt.Errorf("multi-signature wallet %s not found", walletID)
    }

    // Verify signer exists
    signerExists := false
    for _, owner := range wallet.Owners {
        if owner == signerID {
            signerExists = true
            break
        }
    }
    if !signerExists {
        return fmt.Errorf("signer %s is not an owner of wallet %s", signerID, walletID)
    }

    // Remove the signer from the owners list
    for i, owner := range wallet.Owners {
        if owner == signerID {
            wallet.Owners = append(wallet.Owners[:i], wallet.Owners[i+1:]...)
            l.MultiSigWallets[walletID] = wallet
            log.Printf("[INFO] Signature revoked for signer %s on wallet %s.", signerID, walletID)
            return nil
        }
    }

    return fmt.Errorf("failed to revoke signer %s for wallet %s", signerID, walletID)
}


// UpdateMultiSigWallet updates a multi-signature wallet’s required signatures or owners
func (l *AccountsWalletLedger) UpdateMultiSigWallet(walletID string, owners []string, requiredSigs int) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if len(owners) == 0 {
        return fmt.Errorf("owners list cannot be empty")
    }
    if requiredSigs <= 0 || requiredSigs > len(owners) {
        return fmt.Errorf("requiredSigs must be between 1 and the number of owners")
    }

    // Retrieve and update the wallet
    wallet, exists := l.MultiSigWallets[walletID]
    if !exists {
        return fmt.Errorf("multi-signature wallet %s not found", walletID)
    }

    wallet.Owners = owners
    wallet.RequiredSigs = requiredSigs
    l.MultiSigWallets[walletID] = wallet

    log.Printf("[INFO] MultiSig wallet %s updated with %d required signatures and new owners.", walletID, requiredSigs)
    return nil
}


// GetAccount retrieves an account by its ID from the ledger.
func (l *AccountsWalletLedger) GetAccount(accountID string) (*Account, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return nil, fmt.Errorf("accountID cannot be empty")
    }

    // Retrieve the account
    account, exists := l.AccountsWalletLedgerState.Accounts[accountID]
    if !exists {
        return nil, fmt.Errorf("account %s not found", accountID)
    }

    log.Printf("[INFO] Account %s retrieved successfully.", accountID)
    return &account, nil
}


// UpdateAccount modifies an existing account in the ledger.
func (l *AccountsWalletLedger) UpdateAccount(accountID string, updatedAccount Account) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return fmt.Errorf("accountID cannot be empty")
    }

    // Check if the account exists
    if _, exists := l.AccountsWalletLedgerState.Accounts[accountID]; !exists {
        return fmt.Errorf("account %s not found", accountID)
    }

    // Update the account
    l.AccountsWalletLedgerState.Accounts[accountID] = updatedAccount
    log.Printf("[INFO] Account %s updated successfully.", accountID)
    return nil
}


// SaveBalanceSnapshot creates a snapshot of the current balance of an account.
func (l *AccountsWalletLedger) SaveBalanceSnapshot(accountID string) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return fmt.Errorf("accountID cannot be empty")
    }

    // Retrieve the account
    account, exists := l.AccountsWalletLedgerState.Accounts[accountID]
    if !exists {
        return fmt.Errorf("account %s not found", accountID)
    }

    // Create a balance snapshot
    balanceSnapshot := BalanceSnapshot{
        AccountID: accountID,
        Balance:   account.Balance,
        Timestamp: time.Now(),
    }

    // Initialize the BalanceSnapshots map for the account if not already present
    if _, exists := l.AccountsWalletLedgerState.BalanceSnapshots[accountID]; !exists {
        l.AccountsWalletLedgerState.BalanceSnapshots[accountID] = []BalanceSnapshot{}
    }

    // Append the snapshot to the list
    l.AccountsWalletLedgerState.BalanceSnapshots[accountID] = append(l.AccountsWalletLedgerState.BalanceSnapshots[accountID], balanceSnapshot)

    log.Printf("[INFO] Balance snapshot saved for account %s with balance %.2f.", accountID, account.Balance)
    return nil
}


// GetBalanceHistory retrieves the balance history of a given account.
func (l *AccountsWalletLedger) GetBalanceHistory(accountID string) ([]BalanceSnapshot, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return nil, fmt.Errorf("accountID cannot be empty")
    }

    // Check if history exists for the account
    history, exists := l.AccountsWalletLedgerState.BalanceSnapshots[accountID]
    if !exists {
        return nil, fmt.Errorf("no balance history found for account %s", accountID)
    }

    log.Printf("[INFO] Retrieved balance history for account %s. Number of snapshots: %d", accountID, len(history))
    return history, nil
}


// GetBalanceAt retrieves the balance of an account at a specific timestamp.
func (l *AccountsWalletLedger) GetBalanceAt(accountID string, timestamp time.Time) (*big.Int, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if accountID == "" {
        return nil, fmt.Errorf("accountID cannot be empty")
    }

    // Retrieve balance history
    history, exists := l.AccountsWalletLedgerState.BalanceSnapshots[accountID]
    if !exists {
        return nil, fmt.Errorf("no balance history found for account %s", accountID)
    }

    // Find the balance at or before the specified timestamp
    var closestSnapshot *BalanceSnapshot
    for _, snapshot := range history {
        if snapshot.Timestamp.Before(timestamp) || snapshot.Timestamp.Equal(timestamp) {
            if closestSnapshot == nil || snapshot.Timestamp.After(closestSnapshot.Timestamp) {
                closestSnapshot = &snapshot
            }
        }
    }

    if closestSnapshot == nil {
        return nil, fmt.Errorf("no balance record available at or before timestamp %s for account %s", timestamp, accountID)
    }

    // Convert balance to big.Int
    balanceInt := new(big.Int)
    balanceInt.SetInt64(int64(closestSnapshot.Balance))

    log.Printf("[INFO] Balance for account %s at %s: %s", accountID, timestamp, balanceInt.String())
    return balanceInt, nil
}




// GetTrustAccount retrieves a trust account from the ledger by ID.
func (l *AccountsWalletLedger) GetTrustAccount(trustID string) (*TrustAccount, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if trustID == "" {
        return nil, fmt.Errorf("trustID cannot be empty")
    }

    // Retrieve the trust account
    trustAccount, exists := l.AccountsWalletLedgerState.TrustAccounts[trustID]
    if !exists {
        return nil, fmt.Errorf("trust account %s not found", trustID)
    }

    log.Printf("[INFO] Trust account %s retrieved successfully.", trustID)
    return &trustAccount, nil
}


// UpdateTrustAccount updates an existing trust account in the ledger.
func (l *AccountsWalletLedger) UpdateTrustAccount(trustID string, updatedTrustAccount TrustAccount) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if trustID == "" {
        return fmt.Errorf("trustID cannot be empty")
    }

    // Check if the trust account exists
    if _, exists := l.AccountsWalletLedgerState.TrustAccounts[trustID]; !exists {
        return fmt.Errorf("trust account %s not found", trustID)
    }

    // Update the trust account
    l.AccountsWalletLedgerState.TrustAccounts[trustID] = updatedTrustAccount
    log.Printf("[INFO] Trust account %s updated successfully.", trustID)
    return nil
}



var (
	ledgerInstances = make(map[string]*Ledger) // A map to store ledger instances per walletID
	mutex           sync.Mutex                 // Mutex to ensure thread-safe access
)


// LogRecoverySetup logs the setup of wallet recovery mechanisms.
func (l *AccountsWalletLedger) LogRecoverySetup(walletID string, recoveryDetails string) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if recoveryDetails == "" {
        return fmt.Errorf("recoveryDetails cannot be empty")
    }

    // Check if wallet exists
    identity, exists := l.Identities[walletID]
    if !exists {
        return fmt.Errorf("wallet with ID %s not found", walletID)
    }

    // Log the recovery details
    identity.RecoverySetup = recoveryDetails
    l.Identities[walletID] = identity

    log.Printf("[INFO] Recovery setup updated for wallet %s. Details: %s", walletID, recoveryDetails)
    return nil
}


// StoreIdentity stores identity information for a wallet.
func (l *AccountsWalletLedger) StoreIdentity(walletID string, identity Identity) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }

    // Store identity
    l.Identities[walletID] = identity

    log.Printf("[INFO] Identity stored for wallet %s.", walletID)
    return nil
}




// UpdateWalletVerificationStatus updates the verification status of a wallet.
func (l *AccountsWalletLedger) UpdateWalletVerificationStatus(walletID string, status bool) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }

    // Check if the wallet exists
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet with ID %s not found", walletID)
    }

    // Update verification status
    account.Verified = status
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Verification status updated for wallet %s to %v.", walletID, status)
    return nil
}



// StoreWalletKey securely stores the key for a wallet.
func (l *AccountsWalletLedger) StoreWalletKey(walletID, encryptedKey string) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if encryptedKey == "" {
        return fmt.Errorf("encryptedKey cannot be empty")
    }

    // Check if the wallet exists
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet with ID %s not found", walletID)
    }

    // Store the encrypted key
    account.EncryptedKey = encryptedKey
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Encrypted key stored for wallet %s.", walletID)
    return nil
}


// RecordConnectionEvent logs a connection event for a wallet.
func (l *AccountsWalletLedger) RecordConnectionEvent(walletID string, event ConnectionEvent) error {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if event.Timestamp.IsZero() {
        return fmt.Errorf("event timestamp must be valid")
    }

    // Retrieve the account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet with ID %s not found", walletID)
    }

    // Append the event
    account.ConnectionEvents = append(account.ConnectionEvents, event)
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Connection event recorded for wallet %s: %v", walletID, event)
    return nil
}



// GetContractExecutionHistory retrieves the contract execution history for a wallet.
func (l *AccountsWalletLedger) GetContractExecutionHistory(walletID string) ([]ContractExecutionLog, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return nil, fmt.Errorf("walletID cannot be empty")
    }

    // Retrieve the account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return nil, fmt.Errorf("wallet with ID %s not found", walletID)
    }

    if len(account.ContractExecutionLogs) == 0 {
        return nil, fmt.Errorf("no contract execution history found for wallet %s", walletID)
    }

    log.Printf("[INFO] Retrieved %d contract execution logs for wallet %s.", len(account.ContractExecutionLogs), walletID)
    return account.ContractExecutionLogs, nil
}



// GetTokenMintHistory retrieves the minting history of a wallet.
func (l *AccountsWalletLedger) GetTokenMintHistory(walletID string) ([]MintRecord, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return nil, fmt.Errorf("walletID cannot be empty")
    }

    // Retrieve the account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return nil, fmt.Errorf("wallet with ID %s not found", walletID)
    }

    if len(account.MintRecords) == 0 {
        return nil, fmt.Errorf("no mint records found for wallet %s", walletID)
    }

    log.Printf("[INFO] Retrieved %d mint records for wallet %s.", len(account.MintRecords), walletID)
    return account.MintRecords, nil
}



// GetTokenBurnHistory retrieves the burn history of a wallet.
func (l *AccountsWalletLedger) GetTokenBurnHistory(walletID string) ([]BurnRecord, error) {
    l.Lock()
    defer l.Unlock()

    // Input validation
    if walletID == "" {
        return nil, fmt.Errorf("walletID cannot be empty")
    }

    // Retrieve the account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return nil, fmt.Errorf("wallet with ID %s not found", walletID)
    }

    if len(account.BurnRecords) == 0 {
        return nil, fmt.Errorf("no burn records found for wallet %s", walletID)
    }

    log.Printf("[INFO] Retrieved %d burn records for wallet %s.", len(account.BurnRecords), walletID)
    return account.BurnRecords, nil
}



// RecordTokenMint records a token minting event for a wallet.
func (l *AccountsWalletLedger) RecordTokenMint(walletID string, mintRecord MintRecord) error {
    l.Lock()
    defer l.Unlock()

    // Validate inputs
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if mintRecord.Timestamp.IsZero() {
        return fmt.Errorf("mintRecord timestamp must be valid")
    }

    // Retrieve the wallet account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet %s not found", walletID)
    }

    // Append the mint record and update the ledger
    account.MintRecords = append(account.MintRecords, mintRecord)
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Token mint recorded for wallet %s: %+v", walletID, mintRecord)
    return nil
}


// RecordTokenBurn records a token burn event for a wallet.
func (l *AccountsWalletLedger) RecordTokenMint(walletID string, mintRecord MintRecord) error {
    l.Lock()
    defer l.Unlock()

    // Validate inputs
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if mintRecord.Timestamp.IsZero() {
        return fmt.Errorf("mintRecord timestamp must be valid")
    }

    // Retrieve the wallet account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet %s not found", walletID)
    }

    // Append the mint record and update the ledger
    account.MintRecords = append(account.MintRecords, mintRecord)
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Token mint recorded for wallet %s: %+v", walletID, mintRecord)
    return nil
}



// RecordCurrencyExchange logs a currency exchange event between tokens.
func (l *AccountsWalletLedger) RecordCurrencyExchange(walletID string, exchange CurrencyExchange) error {
    l.Lock()
    defer l.Unlock()

    // Validate inputs
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if exchange.Timestamp.IsZero() || exchange.FromToken == "" || exchange.ToToken == "" || exchange.Amount <= 0 {
        return fmt.Errorf("invalid currency exchange data")
    }

    // Retrieve the wallet account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet %s not found", walletID)
    }

    // Append the currency exchange record
    account.CurrencyExchanges = append(account.CurrencyExchanges, exchange)
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Currency exchange logged for wallet %s: %+v", walletID, exchange)
    return nil
}


// RecordWalletNaming logs a custom name for a wallet.
func (l *AccountsWalletLedger) RecordWalletNaming(walletID, customName string) error {
    l.Lock()
    defer l.Unlock()

    // Validate inputs
    if walletID == "" {
        return fmt.Errorf("walletID cannot be empty")
    }
    if customName == "" {
        return fmt.Errorf("customName cannot be empty")
    }

    // Retrieve the wallet account
    account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
    if !exists {
        return fmt.Errorf("wallet %s not found", walletID)
    }

    // Update the custom name
    account.CustomName = customName
    l.AccountsWalletLedgerState.Accounts[walletID] = account

    log.Printf("[INFO] Custom name assigned to wallet %s: %s", walletID, customName)
    return nil
}



// GetNextNonce retrieves the next nonce for a wallet's transaction.
func (l *AccountsWalletLedger) GetNextNonce(walletID string) (uint64, error) {
	l.Lock()
	defer l.Unlock()

	// Validate wallet ID
	if walletID == "" {
		return 0, fmt.Errorf("walletID cannot be empty")
	}

	// Retrieve account details
	account, exists := l.AccountsWalletLedgerState.Accounts[walletID]
	if !exists {
		return 0, fmt.Errorf("wallet %s not found", walletID)
	}

	// Increment and return the nonce
	nextNonce := account.Nonce + 1
	log.Printf("[INFO] Retrieved next nonce for wallet %s: %d", walletID, nextNonce)
	return nextNonce, nil
}


// RetrieveMnemonic retrieves the mnemonic phrase for a given wallet ID from the ledger.
func (l *AccountsWalletLedger) RetrieveMnemonic(walletID string) (string, error) {
	l.Lock()
	defer l.Unlock()

	// Validate wallet ID
	if walletID == "" {
		return "", fmt.Errorf("walletID cannot be empty")
	}

	// Check if the mnemonic exists
	mnemonicList, exists := l.AccountsWalletLedgerState.Mnemonic[walletID]
	if !exists || len(mnemonicList) == 0 {
		return "", fmt.Errorf("no mnemonic stored for wallet ID: %s", walletID)
	}

	// Combine mnemonics into a single string
	stringList := make([]string, len(mnemonicList))
	for i, mnemonic := range mnemonicList {
		stringList[i] = mnemonic.Phrase
	}
	mnemonicPhrase := strings.Join(stringList, " ")

	log.Printf("[INFO] Retrieved mnemonic for wallet %s", walletID)
	return mnemonicPhrase, nil
}

// GetAllBalances returns all wallet balances stored in the ledger.
func (l *AccountsWalletLedger) GetAllBalances() (map[string]float64, error) {
	l.Lock()
	defer l.Unlock()

	// Check if accounts exist in the ledger
	if len(l.AccountsWalletLedgerState.Accounts) == 0 {
		return nil, fmt.Errorf("no balances found in the ledger")
	}

	// Create a copy of balances to avoid data exposure
	balancesCopy := make(map[string]float64)
	for walletID, account := range l.AccountsWalletLedgerState.Accounts {
		balancesCopy[walletID] = account.Balance
	}

	log.Printf("[INFO] Retrieved all wallet balances. Total wallets: %d", len(balancesCopy))
	return balancesCopy, nil
}


// GetWalletKeys retrieves the private and public keys for a wallet by walletID.
func (l *AccountsWalletLedger) GetWalletKeys(walletID string) (string, string, error) {
	l.Lock()
	defer l.Unlock()

	// Validate wallet ID
	if walletID == "" {
		return "", "", fmt.Errorf("walletID cannot be empty")
	}

	// Retrieve wallet keys
	walletData, exists := l.AccountsWalletLedgerState.Accounts[walletID]
	if !exists {
		return "", "", fmt.Errorf("walletID %s not found in the ledger", walletID)
	}

	// Ensure keys are present
	if walletData.PrivateKey == "" || walletData.PublicKey == "" {
		return "", "", fmt.Errorf("keys not found for walletID %s", walletID)
	}

	log.Printf("[INFO] Retrieved keys for wallet %s", walletID)
	return walletData.PrivateKey, walletData.PublicKey, nil
}


// GetInstanceForWallet returns the ledger instance for a specific wallet ID.
func GetInstanceForWallet(walletID string) (*Ledger, error) {
	// Input validation
	if walletID == "" {
		return nil, fmt.Errorf("walletID cannot be empty")
	}

	// Lock for thread-safe access
	mutex.Lock()
	defer mutex.Unlock()

	// Check for the existence of the ledger instance
	ledgerInstance, exists := ledgerInstances[walletID]
	if !exists {
		err := fmt.Errorf("ledger instance not found for wallet %s", walletID)
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	log.Printf("[INFO] Retrieved ledger instance for wallet %s", walletID)
	return ledgerInstance, nil
}


// RegisterLedgerInstance registers a ledger instance for a specific wallet ID.
func RegisterLedgerInstance(walletID string, ledger *Ledger) error {
	// Input validation
	if walletID == "" {
		return fmt.Errorf("walletID cannot be empty")
	}
	if ledger == nil {
		return fmt.Errorf("ledger instance cannot be nil")
	}

	// Lock for thread-safe access
	mutex.Lock()
	defer mutex.Unlock()

	// Register the ledger instance
	ledgerInstances[walletID] = ledger
	log.Printf("[INFO] Ledger instance registered for wallet %s", walletID)
	return nil
}


