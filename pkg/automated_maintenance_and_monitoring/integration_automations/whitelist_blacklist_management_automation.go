package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    ListManagementInterval = 3000 * time.Millisecond // Interval for checking and enforcing whitelist/blacklist rules
)

// WhitelistBlacklistManagementAutomation manages and enforces whitelist/blacklist rules for transactions and addresses
type WhitelistBlacklistManagementAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store whitelist/blacklist actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    listManagementCheckCount int                       // Counter for whitelist/blacklist check cycles
    whitelist             map[string]bool              // Mapping of whitelisted addresses/contracts
    blacklist             map[string]bool              // Mapping of blacklisted addresses/contracts
}

// NewWhitelistBlacklistManagementAutomation initializes the automation for whitelist/blacklist management
func NewWhitelistBlacklistManagementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *WhitelistBlacklistManagementAutomation {
    return &WhitelistBlacklistManagementAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        whitelist:             make(map[string]bool),
        blacklist:             make(map[string]bool),
        listManagementCheckCount: 0,
    }
}

// StartListManagementCheck starts the continuous loop for monitoring and enforcing whitelist/blacklist rules
func (automation *WhitelistBlacklistManagementAutomation) StartListManagementCheck() {
    ticker := time.NewTicker(ListManagementInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceListRules()
        }
    }()
}

// monitorAndEnforceListRules checks incoming contracts/addresses and enforces whitelist/blacklist rules
func (automation *WhitelistBlacklistManagementAutomation) monitorAndEnforceListRules() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch incoming transactions or contract addresses
    newContracts := automation.consensusSystem.GetPendingTransactions()

    for _, contract := range newContracts {
        contractAddress := contract.Address

        if automation.isBlacklisted(contractAddress) {
            fmt.Printf("Blacklisted contract/address detected: %s. Blocking transaction.\n", contractAddress)
            automation.blockTransaction(contract)
        } else if automation.isWhitelisted(contractAddress) {
            fmt.Printf("Whitelisted contract/address detected: %s. Allowing transaction.\n", contractAddress)
            automation.allowTransaction(contract)
        } else {
            fmt.Printf("Contract/address %s is neither whitelisted nor blacklisted. Requires action.\n", contractAddress)
        }
    }

    automation.listManagementCheckCount++
    fmt.Printf("Whitelist/Blacklist check cycle #%d executed.\n", automation.listManagementCheckCount)

    if automation.listManagementCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeListManagementCycle()
    }
}

// isWhitelisted checks if a contract/address is whitelisted
func (automation *WhitelistBlacklistManagementAutomation) isWhitelisted(address string) bool {
    return automation.whitelist[address]
}

// isBlacklisted checks if a contract/address is blacklisted
func (automation *WhitelistBlacklistManagementAutomation) isBlacklisted(address string) bool {
    return automation.blacklist[address]
}

// blockTransaction prevents the execution of transactions associated with blacklisted contracts/addresses
func (automation *WhitelistBlacklistManagementAutomation) blockTransaction(contract common.Transaction) {
    encryptedData := encryption.EncryptData(contract)

    success := automation.consensusSystem.BlockTransaction(encryptedData)
    if success {
        fmt.Printf("Transaction for blacklisted contract %s successfully blocked.\n", contract.Address)
        automation.logListManagementEvent(contract.Address, "Blocked")
    } else {
        fmt.Printf("Error blocking transaction for blacklisted contract %s.\n", contract.Address)
    }
}

// allowTransaction allows the execution of transactions associated with whitelisted contracts/addresses
func (automation *WhitelistBlacklistManagementAutomation) allowTransaction(contract common.Transaction) {
    encryptedData := encryption.EncryptData(contract)

    success := automation.consensusSystem.AllowTransaction(encryptedData)
    if success {
        fmt.Printf("Transaction for whitelisted contract %s successfully allowed.\n", contract.Address)
        automation.logListManagementEvent(contract.Address, "Allowed")
    } else {
        fmt.Printf("Error allowing transaction for whitelisted contract %s.\n", contract.Address)
    }
}

// logListManagementEvent logs the result of the whitelist/blacklist action into the ledger for traceability
func (automation *WhitelistBlacklistManagementAutomation) logListManagementEvent(address string, action string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("list-management-%s-%d", address, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Whitelist/Blacklist Management",
        Status:    action,
        Details:   fmt.Sprintf("Address: %s, Action: %s", address, action),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with whitelist/blacklist event for address %s.\n", address)
}

// finalizeListManagementCycle finalizes the whitelist/blacklist check cycle and logs the result in the ledger
func (automation *WhitelistBlacklistManagementAutomation) finalizeListManagementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeWhitelistBlacklistCycle()
    if success {
        fmt.Println("Whitelist/Blacklist check cycle finalized successfully.")
        automation.logListManagementCycleFinalization()
    } else {
        fmt.Println("Error finalizing whitelist/blacklist check cycle.")
    }
}

// logListManagementCycleFinalization logs the finalization of a whitelist/blacklist check cycle into the ledger
func (automation *WhitelistBlacklistManagementAutomation) logListManagementCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("list-management-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Whitelist/Blacklist Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with whitelist/blacklist cycle finalization.")
}

// AddToWhitelist adds a contract or address to the whitelist
func (automation *WhitelistBlacklistManagementAutomation) AddToWhitelist(address string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    automation.whitelist[address] = true
    fmt.Printf("Address %s added to whitelist.\n", address)
}

// AddToBlacklist adds a contract or address to the blacklist
func (automation *WhitelistBlacklistManagementAutomation) AddToBlacklist(address string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    automation.blacklist[address] = true
    fmt.Printf("Address %s added to blacklist.\n", address)
}
