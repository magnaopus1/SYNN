package security_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/hardware"
)

const (
    HardwareWalletCheckInterval = 30 * time.Second // Interval for checking hardware wallet connections
    MaxInvalidHardwareAttempts  = 3                // Max invalid attempts before locking out the wallet
)

// HardwareWalletSecurityAutomation monitors and integrates hardware wallets for secure usage
type HardwareWalletSecurityAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging hardware wallet events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    invalidAttempts   map[string]int               // Tracks invalid hardware wallet attempts by user or wallet ID
}

// NewHardwareWalletSecurityAutomation initializes the automation for hardware wallet security
func NewHardwareWalletSecurityAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *HardwareWalletSecurityAutomation {
    return &HardwareWalletSecurityAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        invalidAttempts:  make(map[string]int),
    }
}

// StartHardwareWalletMonitoring starts the continuous loop for hardware wallet monitoring and integration
func (automation *HardwareWalletSecurityAutomation) StartHardwareWalletMonitoring() {
    ticker := time.NewTicker(HardwareWalletCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorHardwareWallets()
        }
    }()
}

// monitorHardwareWallets checks for any unauthorized or invalid hardware wallet connections
func (automation *HardwareWalletSecurityAutomation) monitorHardwareWallets() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    connectedWallets := hardware.ListConnectedHardwareWallets()

    for _, wallet := range connectedWallets {
        if wallet.IsAuthorized {
            fmt.Printf("Hardware wallet %s is securely connected.\n", wallet.ID)
            automation.logWalletConnection(wallet)
        } else {
            fmt.Printf("Unauthorized hardware wallet %s detected.\n", wallet.ID)
            automation.handleUnauthorizedWallet(wallet)
        }
    }
}

// handleUnauthorizedWallet processes unauthorized hardware wallet connections
func (automation *HardwareWalletSecurityAutomation) handleUnauthorizedWallet(wallet hardware.HardwareWallet) {
    automation.invalidAttempts[wallet.ID]++

    if automation.invalidAttempts[wallet.ID] >= MaxInvalidHardwareAttempts {
        automation.lockWallet(wallet.ID)
    }

    automation.logUnauthorizedWallet(wallet)
}

// lockWallet locks a hardware wallet after repeated invalid attempts
func (automation *HardwareWalletSecurityAutomation) lockWallet(walletID string) {
    hardware.LockHardwareWallet(walletID)
    fmt.Printf("Hardware wallet %s has been locked due to repeated invalid connection attempts.\n", walletID)
    automation.logWalletLock(walletID)
}

// logWalletConnection logs a valid hardware wallet connection to the ledger
func (automation *HardwareWalletSecurityAutomation) logWalletConnection(wallet hardware.HardwareWallet) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("hardware-wallet-connection-%s", wallet.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Hardware Wallet Connection",
        Status:    "Connected",
        Details:   fmt.Sprintf("Hardware wallet %s connected successfully.", wallet.ID),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with hardware wallet connection for wallet %s.\n", wallet.ID)
}

// logUnauthorizedWallet logs unauthorized hardware wallet connection attempts to the ledger
func (automation *HardwareWalletSecurityAutomation) logUnauthorizedWallet(wallet hardware.HardwareWallet) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("unauthorized-wallet-%s", wallet.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Unauthorized Hardware Wallet",
        Status:    "Unauthorized",
        Details:   fmt.Sprintf("Unauthorized connection attempt detected for hardware wallet %s.", wallet.ID),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with unauthorized hardware wallet connection for wallet %s.\n", wallet.ID)
}

// logWalletLock logs the event when a hardware wallet is locked due to repeated invalid attempts
func (automation *HardwareWalletSecurityAutomation) logWalletLock(walletID string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("wallet-lock-%s", walletID),
        Timestamp: time.Now().Unix(),
        Type:      "Hardware Wallet Lock",
        Status:    "Locked",
        Details:   fmt.Sprintf("Hardware wallet %s has been locked due to repeated invalid attempts.", walletID),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with hardware wallet lock event for wallet %s.\n", walletID)
}

// ensureHardwareWalletIntegrity checks the integrity of the hardware wallet connections and ensures they are secure
func (automation *HardwareWalletSecurityAutomation) ensureHardwareWalletIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateHardwareWalletIntegrity()
    if !integrityValid {
        fmt.Println("Hardware wallet integrity breach detected. Re-checking connections.")
        automation.monitorHardwareWallets()
    } else {
        fmt.Println("Hardware wallet connections are secure.")
    }
}
