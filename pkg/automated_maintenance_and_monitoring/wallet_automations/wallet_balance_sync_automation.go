package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

const (
    WalletSyncInterval        = 10 * time.Minute  // Interval for performing wallet balance synchronization
    MaxSyncRetries            = 3                 // Maximum number of retry attempts for wallet synchronization
)

// WalletBalanceSyncAutomation automates the process of syncing wallet balances across the network
type WalletBalanceSyncAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger for logging wallet sync events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    syncRetryCount   map[string]int               // Counter for retrying balance sync on failure
    syncCycleCount   int                          // Counter for wallet sync cycles
}

// NewWalletBalanceSyncAutomation initializes the automation for wallet balance synchronization
func NewWalletBalanceSyncAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *WalletBalanceSyncAutomation {
    return &WalletBalanceSyncAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        syncRetryCount:  make(map[string]int),
        syncCycleCount:  0,
    }
}

// StartWalletBalanceSync starts the continuous loop for synchronizing wallet balances
func (automation *WalletBalanceSyncAutomation) StartWalletBalanceSync() {
    ticker := time.NewTicker(WalletSyncInterval)

    go func() {
        for range ticker.C {
            automation.syncWalletBalances()
        }
    }()
}

// syncWalletBalances fetches all wallet balances and performs synchronization
func (automation *WalletBalanceSyncAutomation) syncWalletBalances() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of all wallets that need balance synchronization
    walletList := automation.consensusSystem.GetWalletsForSync()

    if len(walletList) > 0 {
        for _, wallet := range walletList {
            fmt.Printf("Initiating balance sync for wallet %s.\n", wallet.ID)
            automation.syncWalletBalance(wallet)
        }
    } else {
        fmt.Println("No wallets need balance sync at this time.")
    }

    automation.syncCycleCount++
    fmt.Printf("Wallet sync cycle #%d executed.\n", automation.syncCycleCount)

    if automation.syncCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeSyncCycle()
    }
}

// syncWalletBalance triggers balance sync for a single wallet
func (automation *WalletBalanceSyncAutomation) syncWalletBalance(wallet common.Wallet) {
    // Encrypt wallet data before syncing
    encryptedWalletData := automation.encryptWalletData(wallet)

    // Trigger wallet balance sync through Synnergy Consensus
    syncSuccess := automation.consensusSystem.SyncWalletBalance(wallet, encryptedWalletData)

    if syncSuccess {
        fmt.Printf("Wallet balance for wallet %s synced successfully.\n", wallet.ID)
        automation.logWalletSyncEvent(wallet)
        automation.resetSyncRetry(wallet.ID)
    } else {
        fmt.Printf("Error syncing wallet balance for wallet %s. Retrying...\n", wallet.ID)
        automation.retrySync(wallet)
    }
}

// retrySync retries the wallet balance sync a limited number of times if it fails
func (automation *WalletBalanceSyncAutomation) retrySync(wallet common.Wallet) {
    automation.syncRetryCount[wallet.ID]++
    if automation.syncRetryCount[wallet.ID] < MaxSyncRetries {
        automation.syncWalletBalance(wallet)
    } else {
        fmt.Printf("Max retries reached for wallet %s. Balance sync failed.\n", wallet.ID)
        automation.logSyncFailure(wallet)
    }
}

// resetSyncRetry resets the retry count for a wallet balance sync
func (automation *WalletBalanceSyncAutomation) resetSyncRetry(walletID string) {
    automation.syncRetryCount[walletID] = 0
}

// finalizeSyncCycle finalizes the wallet balance sync cycle and logs the result in the ledger
func (automation *WalletBalanceSyncAutomation) finalizeSyncCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeSyncCycle()
    if success {
        fmt.Println("Wallet sync cycle finalized successfully.")
        automation.logSyncCycleFinalization()
    } else {
        fmt.Println("Error finalizing wallet sync cycle.")
    }
}

// logWalletSyncEvent logs the wallet sync event into the ledger
func (automation *WalletBalanceSyncAutomation) logWalletSyncEvent(wallet common.Wallet) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-sync-%s", wallet.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Sync",
        Status:    "Completed",
        Details:   fmt.Sprintf("Balance sync completed successfully for wallet %s.", wallet.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with wallet balance sync event for wallet %s.\n", wallet.ID)
}

// logSyncCycleFinalization logs the finalization of a wallet balance sync cycle into the ledger
func (automation *WalletBalanceSyncAutomation) logSyncCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-sync-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Sync Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with wallet sync cycle finalization.")
}

// logSyncFailure logs the failure of a wallet balance sync event into the ledger
func (automation *WalletBalanceSyncAutomation) logSyncFailure(wallet common.Wallet) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-sync-failure-%s", wallet.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Sync Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Balance sync failed for wallet %s after maximum retries.", wallet.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with wallet balance sync failure for wallet %s.\n", wallet.ID)
}

// encryptWalletData encrypts the wallet data before syncing
func (automation *WalletBalanceSyncAutomation) encryptWalletData(wallet common.Wallet) common.Wallet {
    encryptedData, err := encryption.EncryptData(wallet.Data)
    if err != nil {
        fmt.Println("Error encrypting wallet data:", err)
        return wallet
    }

    wallet.EncryptedData = encryptedData
    fmt.Println("Wallet data successfully encrypted.")
    return wallet
}

// ensureSyncIntegrity checks the integrity of the wallet balance sync process
func (automation *WalletBalanceSyncAutomation) ensureSyncIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateSyncIntegrity()
    if !integrityValid {
        fmt.Println("Wallet sync data integrity breach detected. Re-triggering sync.")
        automation.syncWalletBalances()
    } else {
        fmt.Println("Wallet sync data integrity is valid.")
    }
}
