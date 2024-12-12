package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

const (
    WalletCreationCheckInterval = 10 * time.Second // Interval for checking wallet creation requests
    MaxWalletCreationAttempts   = 5               // Maximum number of wallet creation attempts before restriction
)

// WalletCreationRestrictionAutomation enforces restrictions and policies for wallet creation
type WalletCreationRestrictionAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger for logging wallet creation events
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    creationAttemptTracker map[string]int               // Track wallet creation attempts per IP or account
}

// NewWalletCreationRestrictionAutomation initializes the automation for wallet creation restriction and monitoring
func NewWalletCreationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *WalletCreationRestrictionAutomation {
    return &WalletCreationRestrictionAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        creationAttemptTracker: make(map[string]int),
    }
}

// StartWalletCreationCheck starts the continuous loop for monitoring and restricting wallet creation
func (automation *WalletCreationRestrictionAutomation) StartWalletCreationCheck() {
    ticker := time.NewTicker(WalletCreationCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorWalletCreationRequests()
        }
    }()
}

// monitorWalletCreationRequests fetches new wallet creation requests and checks if they comply with the restriction protocols
func (automation *WalletCreationRestrictionAutomation) monitorWalletCreationRequests() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of pending wallet creation requests
    walletCreationRequests := automation.consensusSystem.GetPendingWalletCreationRequests()

    for _, request := range walletCreationRequests {
        fmt.Printf("Processing wallet creation request from account %s or IP %s.\n", request.AccountID, request.IP)
        if automation.isWalletCreationAllowed(request) {
            automation.processWalletCreation(request)
        } else {
            automation.blockWalletCreation(request)
        }
    }
}

// isWalletCreationAllowed checks if a wallet creation request is allowed based on restrictions
func (automation *WalletCreationRestrictionAutomation) isWalletCreationAllowed(request common.WalletCreationRequest) bool {
    // Check if the request has exceeded the maximum allowed creation attempts
    if automation.creationAttemptTracker[request.IP] >= MaxWalletCreationAttempts {
        fmt.Printf("Wallet creation request from IP %s has been blocked due to too many attempts.\n", request.IP)
        return false
    }

    // Additional checks could include blacklists, duplicate wallet detection, etc.
    if automation.consensusSystem.IsAccountBlacklisted(request.AccountID) {
        fmt.Printf("Wallet creation request from account %s has been blocked (blacklisted).\n", request.AccountID)
        return false
    }

    return true
}

// processWalletCreation processes a valid wallet creation request
func (automation *WalletCreationRestrictionAutomation) processWalletCreation(request common.WalletCreationRequest) {
    // Encrypt wallet data before processing
    encryptedWalletData := automation.encryptWalletData(request)

    // Create the wallet and update the ledger
    success := automation.consensusSystem.CreateWallet(request.AccountID, encryptedWalletData)

    if success {
        fmt.Printf("Wallet for account %s created successfully.\n", request.AccountID)
        automation.logWalletCreationEvent(request)
        automation.resetCreationAttempts(request.IP)
    } else {
        fmt.Printf("Error creating wallet for account %s.\n", request.AccountID)
        automation.creationAttemptTracker[request.IP]++
    }
}

// blockWalletCreation blocks a wallet creation request and logs the event
func (automation *WalletCreationRestrictionAutomation) blockWalletCreation(request common.WalletCreationRequest) {
    // Log the failed creation attempt
    automation.logFailedCreationAttempt(request)

    // Send a block event to Synnergy Consensus
    err := automation.consensusSystem.BlockWalletCreation(request.AccountID)
    if err != nil {
        fmt.Printf("Failed to block wallet creation for account %s: %v\n", request.AccountID, err)
    } else {
        fmt.Printf("Wallet creation blocked for account %s.\n", request.AccountID)
    }
}

// resetCreationAttempts resets the creation attempt count for a specific IP
func (automation *WalletCreationRestrictionAutomation) resetCreationAttempts(IP string) {
    automation.creationAttemptTracker[IP] = 0
}

// logWalletCreationEvent logs the wallet creation event into the ledger
func (automation *WalletCreationRestrictionAutomation) logWalletCreationEvent(request common.WalletCreationRequest) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-creation-%s", request.AccountID),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Creation",
        Status:    "Success",
        Details:   fmt.Sprintf("Wallet created successfully for account %s.", request.AccountID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with wallet creation event for account %s.\n", request.AccountID)
}

// logFailedCreationAttempt logs a failed wallet creation attempt into the ledger
func (automation *WalletCreationRestrictionAutomation) logFailedCreationAttempt(request common.WalletCreationRequest) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-creation-failed-%s", request.AccountID),
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Creation",
        Status:    "Failed",
        Details:   fmt.Sprintf("Wallet creation failed for account %s due to restrictions or maximum attempts.", request.AccountID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with failed wallet creation event for account %s.\n", request.AccountID)
}

// encryptWalletData encrypts the wallet data before processing the creation
func (automation *WalletCreationRestrictionAutomation) encryptWalletData(request common.WalletCreationRequest) common.WalletData {
    encryptedData, err := encryption.EncryptData(request.Data)
    if err != nil {
        fmt.Println("Error encrypting wallet data:", err)
        return request.Data
    }

    request.EncryptedData = encryptedData
    fmt.Println("Wallet data successfully encrypted.")
    return request.Data
}

// ensureWalletCreationIntegrity ensures that wallet creation processes maintain their integrity
func (automation *WalletCreationRestrictionAutomation) ensureWalletCreationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateWalletCreationIntegrity()
    if !integrityValid {
        fmt.Println("Wallet creation integrity breach detected. Re-validating all wallet creation requests.")
        automation.monitorWalletCreationRequests()
    } else {
        fmt.Println("Wallet creation integrity is valid.")
    }
}
