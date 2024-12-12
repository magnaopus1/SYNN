package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

const (
    WalletBackupInterval       = 24 * time.Hour  // Interval for performing wallet backups
    SubBlocksPerBlock          = 1000            // Number of sub-blocks in a block
    MaxBackupRetries           = 3               // Maximum number of retry attempts for a wallet backup
)

// WalletBackupMaintenanceAutomation automates the process of backing up wallets, ensuring encryption and ledger integration
type WalletBackupMaintenanceAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging wallet backup actions
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    backupRetryCount  map[string]int               // Counter for retrying backups on failure
    backupCycleCount  int                          // Counter for wallet backup cycles
}

// NewWalletBackupMaintenanceAutomation initializes the automation for wallet backups
func NewWalletBackupMaintenanceAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *WalletBackupMaintenanceAutomation {
    return &WalletBackupMaintenanceAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        backupRetryCount: make(map[string]int),
        backupCycleCount: 0,
    }
}

// StartWalletBackup starts the continuous loop for regularly backing up wallets
func (automation *WalletBackupMaintenanceAutomation) StartWalletBackup() {
    ticker := time.NewTicker(WalletBackupInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndBackupWallets()
        }
    }()
}

// monitorAndBackupWallets checks all wallets and performs backups if needed
func (automation *WalletBackupMaintenanceAutomation) monitorAndBackupWallets() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of all wallets that need to be backed up
    walletList := automation.consensusSystem.GetWalletsForBackup()

    if len(walletList) > 0 {
        for _, wallet := range walletList {
            fmt.Printf("Initiating backup for wallet %s.\n", wallet.ID)
            automation.backupWallet(wallet)
        }
    } else {
        fmt.Println("No wallets need backup at this time.")
    }

    automation.backupCycleCount++
    fmt.Printf("Wallet backup cycle #%d executed.\n", automation.backupCycleCount)

    if automation.backupCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeBackupCycle()
    }
}

// backupWallet encrypts the wallet data and triggers the backup process
func (automation *WalletBackupMaintenanceAutomation) backupWallet(wallet common.Wallet) {
    // Encrypt wallet data before backup
    encryptedWalletData := automation.encryptWalletData(wallet)

    // Trigger wallet backup through the Synnergy Consensus system
    backupSuccess := automation.consensusSystem.BackupWallet(wallet, encryptedWalletData)

    if backupSuccess {
        fmt.Printf("Wallet %s backed up successfully.\n", wallet.ID)
        automation.logWalletBackupEvent(wallet)
        automation.resetBackupRetry(wallet.ID)
    } else {
        fmt.Printf("Error backing up wallet %s. Retrying...\n", wallet.ID)
        automation.retryBackup(wallet)
    }
}

// retryBackup attempts to retry a failed backup a limited number of times
func (automation *WalletBackupMaintenanceAutomation) retryBackup(wallet common.Wallet) {
    automation.backupRetryCount[wallet.ID]++
    if automation.backupRetryCount[wallet.ID] < MaxBackupRetries {
        automation.backupWallet(wallet)
    } else {
        fmt.Printf("Max retries reached for wallet %s. Backup failed.\n", wallet.ID)
        automation.logBackupFailure(wallet)
    }
}

// resetBackupRetry resets the retry count for a wallet backup
func (automation *WalletBackupMaintenanceAutomation) resetBackupRetry(walletID string) {
    automation.backupRetryCount[walletID] = 0
}

// finalizeBackupCycle finalizes the wallet backup check cycle and logs the result in the ledger
func (automation *WalletBackupMaintenanceAutomation) finalizeBackupCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeBackupCycle()
    if success {
        fmt.Println("Wallet backup check cycle finalized successfully.")
        automation.logBackupCycleFinalization()
    } else {
        fmt.Println("Error finalizing wallet backup check cycle.")
    }
}

// logWalletBackupEvent logs the wallet backup event into the ledger
func (automation *WalletBackupMaintenanceAutomation) logWalletBackupEvent(wallet common.Wallet) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-backup-%s", wallet.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Backup",
        Status:    "Completed",
        Details:   fmt.Sprintf("Backup successfully completed for wallet %s.", wallet.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with wallet backup event for wallet %s.\n", wallet.ID)
}

// logBackupCycleFinalization logs the finalization of a wallet backup cycle into the ledger
func (automation *WalletBackupMaintenanceAutomation) logBackupCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-backup-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Backup Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with wallet backup cycle finalization.")
}

// logBackupFailure logs the failure of a wallet backup event into the ledger
func (automation *WalletBackupMaintenanceAutomation) logBackupFailure(wallet common.Wallet) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-backup-failure-%s", wallet.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Backup Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Backup failed for wallet %s after maximum retries.", wallet.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with wallet backup failure for wallet %s.\n", wallet.ID)
}

// encryptWalletData encrypts the wallet data before backup
func (automation *WalletBackupMaintenanceAutomation) encryptWalletData(wallet common.Wallet) common.Wallet {
    encryptedData, err := encryption.EncryptData(wallet.Data)
    if err != nil {
        fmt.Println("Error encrypting wallet data:", err)
        return wallet
    }

    wallet.EncryptedData = encryptedData
    fmt.Println("Wallet data successfully encrypted.")
    return wallet
}

// ensureBackupIntegrity checks the integrity of the wallet backup process
func (automation *WalletBackupMaintenanceAutomation) ensureBackupIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateBackupIntegrity()
    if !integrityValid {
        fmt.Println("Wallet backup data integrity breach detected. Re-triggering backup.")
        automation.monitorAndBackupWallets()
    } else {
        fmt.Println("Wallet backup data integrity is valid.")
    }
}
