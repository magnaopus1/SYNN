package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/smartcontracts"
)

const (
    VersionCheckInterval = 8000 * time.Millisecond // Interval for checking smart contract versioning
    SubBlocksPerBlock    = 1000                    // Number of sub-blocks per block
)

// SmartContractVersioningUpgradeAutomation automates the process of upgrading and managing smart contract versions
type SmartContractVersioningUpgradeAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger to store contract upgrade records
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    versionCheckCount    int                          // Counter for version upgrade check cycles
}

// NewSmartContractVersioningUpgradeAutomation initializes the automation for smart contract versioning and upgrade handling
func NewSmartContractVersioningUpgradeAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartContractVersioningUpgradeAutomation {
    return &SmartContractVersioningUpgradeAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        versionCheckCount: 0,
    }
}

// StartVersionUpgradeCheck starts the continuous loop for checking and enforcing contract version upgrades
func (automation *SmartContractVersioningUpgradeAutomation) StartVersionUpgradeCheck() {
    ticker := time.NewTicker(VersionCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndApplyUpgrades()
        }
    }()
}

// checkAndApplyUpgrades checks all contracts for new versions and applies upgrades if necessary
func (automation *SmartContractVersioningUpgradeAutomation) checkAndApplyUpgrades() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Step 1: Fetch all contracts pending version check
    contractsPendingUpgrade, err := automation.consensusSystem.GetContractsPendingVersionCheck()
    if err != nil {
        fmt.Printf("Error fetching contracts for version upgrade check: %v\n", err)
        return
    }

    // Step 2: Process each contract and upgrade if necessary
    for _, contract := range contractsPendingUpgrade {
        fmt.Printf("Checking version for contract: %s\n", contract.ID)

        // Step 3: Encrypt contract data before upgrade check
        encryptedContract, err := automation.encryptContractData(contract)
        if err != nil {
            fmt.Printf("Error encrypting contract %s: %v\n", contract.ID, err)
            automation.logUpgradeResult(contract, "Encryption Failed")
            continue
        }

        // Step 4: Check for version mismatch and apply upgrade
        versionUpgraded := automation.applyVersionUpgradeIfNecessary(encryptedContract)
        if versionUpgraded {
            fmt.Printf("Contract %s upgraded successfully.\n", contract.ID)
            automation.logUpgradeResult(contract, "Upgrade Successful")
        } else {
            fmt.Printf("No version upgrade needed for contract %s.\n", contract.ID)
            automation.logUpgradeResult(contract, "No Upgrade Needed")
        }
    }

    automation.versionCheckCount++
    fmt.Printf("Smart contract version upgrade check cycle #%d completed.\n", automation.versionCheckCount)

    if automation.versionCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeUpgradeCheckCycle()
    }
}

// encryptContractData encrypts the smart contract data before performing version checks and upgrades
func (automation *SmartContractVersioningUpgradeAutomation) encryptContractData(contract common.SmartContract) (common.SmartContract, error) {
    fmt.Println("Encrypting smart contract data.")

    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        return contract, fmt.Errorf("failed to encrypt contract data: %v", err)
    }

    contract.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract, nil
}

// applyVersionUpgradeIfNecessary checks if a new version exists and applies the upgrade if necessary
func (automation *SmartContractVersioningUpgradeAutomation) applyVersionUpgradeIfNecessary(contract common.SmartContract) bool {
    fmt.Printf("Checking if contract %s requires a version upgrade.\n", contract.ID)

    upgradeAvailable := automation.consensusSystem.IsNewVersionAvailable(contract)
    if upgradeAvailable {
        fmt.Printf("New version available for contract %s. Upgrading...\n", contract.ID)
        success := automation.consensusSystem.UpgradeContractVersion(contract)
        if success {
            fmt.Printf("Contract %s upgraded successfully to the new version.\n", contract.ID)
            return true
        } else {
            fmt.Printf("Error upgrading contract %s.\n", contract.ID)
            return false
        }
    }

    return false
}

// logUpgradeResult logs the result of a smart contract version upgrade in the ledger
func (automation *SmartContractVersioningUpgradeAutomation) logUpgradeResult(contract common.SmartContract, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-upgrade-%s", contract.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Version Upgrade",
        Status:    result,
        Details:   fmt.Sprintf("Version upgrade result for contract %s: %s", contract.ID, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with contract version upgrade result for contract %s: %s\n", contract.ID, result)
}

// finalizeUpgradeCheckCycle finalizes the smart contract version upgrade check cycle and logs the results
func (automation *SmartContractVersioningUpgradeAutomation) finalizeUpgradeCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeUpgradeCheckCycle()
    if success {
        fmt.Println("Smart contract version upgrade check cycle finalized successfully.")
        automation.logUpgradeCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing smart contract version upgrade check cycle.")
    }
}

// logUpgradeCheckCycleFinalization logs the finalization of the contract version upgrade check cycle in the ledger
func (automation *SmartContractVersioningUpgradeAutomation) logUpgradeCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("upgrade-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Version Upgrade Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with smart contract version upgrade check cycle finalization.")
}
