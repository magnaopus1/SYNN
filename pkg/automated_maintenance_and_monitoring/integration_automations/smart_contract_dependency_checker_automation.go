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
    DependencyCheckInterval = 10000 * time.Millisecond // Interval for checking contract dependencies
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks per block
)

// SmartContractDependencyCheckerAutomation automates the process of checking smart contract dependencies
type SmartContractDependencyCheckerAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store contract dependency results
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    dependencyCheckCount int                         // Counter for dependency check cycles
}

// NewSmartContractDependencyCheckerAutomation initializes the automation for contract dependency checking
func NewSmartContractDependencyCheckerAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartContractDependencyCheckerAutomation {
    return &SmartContractDependencyCheckerAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        dependencyCheckCount: 0,
    }
}

// StartDependencyCheck starts the continuous loop for checking and enforcing contract dependencies
func (automation *SmartContractDependencyCheckerAutomation) StartDependencyCheck() {
    ticker := time.NewTicker(DependencyCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndEnforceDependencies()
        }
    }()
}

// checkAndEnforceDependencies checks all new smart contracts for unresolved dependencies and enforces validation
func (automation *SmartContractDependencyCheckerAutomation) checkAndEnforceDependencies() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Step 1: Fetch all new smart contracts for dependency checking
    contractsToCheck, err := automation.consensusSystem.GetContractsPendingDependencyCheck()
    if err != nil {
        fmt.Printf("Error fetching contracts for dependency check: %v\n", err)
        return
    }

    // Step 2: Process each contract, check for dependencies, and validate them
    for _, contract := range contractsToCheck {
        fmt.Printf("Checking dependencies for contract: %s\n", contract.ID)

        // Step 3: Encrypt the contract data before dependency checking
        encryptedContract, err := automation.encryptContractData(contract)
        if err != nil {
            fmt.Printf("Error encrypting contract: %s - %v\n", contract.ID, err)
            automation.logDependencyResult(contract, "Encryption Failed")
            continue
        }

        // Step 4: Check for unresolved dependencies
        dependenciesResolved := automation.validateDependencies(encryptedContract)
        if dependenciesResolved {
            fmt.Printf("Dependencies for contract %s are resolved.\n", contract.ID)
            automation.logDependencyResult(contract, "Dependencies Resolved")
        } else {
            fmt.Printf("Unresolved dependencies found for contract %s.\n", contract.ID)
            automation.logDependencyResult(contract, "Unresolved Dependencies")
        }
    }

    automation.dependencyCheckCount++
    fmt.Printf("Smart contract dependency check cycle #%d completed.\n", automation.dependencyCheckCount)

    if automation.dependencyCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeDependencyCheckCycle()
    }
}

// encryptContractData encrypts the smart contract data before checking dependencies
func (automation *SmartContractDependencyCheckerAutomation) encryptContractData(contract common.SmartContract) (common.SmartContract, error) {
    fmt.Println("Encrypting smart contract data.")

    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        return contract, fmt.Errorf("failed to encrypt contract data: %v", err)
    }

    contract.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract, nil
}

// validateDependencies checks if all required dependencies for the smart contract are resolved
func (automation *SmartContractDependencyCheckerAutomation) validateDependencies(contract common.SmartContract) bool {
    fmt.Printf("Validating dependencies for contract %s.\n", contract.ID)

    dependenciesResolved := automation.consensusSystem.CheckContractDependencies(contract)
    if !dependenciesResolved {
        fmt.Printf("Unresolved dependencies for contract %s.\n", contract.ID)
        return false
    }

    fmt.Printf("Dependencies for contract %s successfully validated.\n", contract.ID)
    return true
}

// logDependencyResult logs the result of a smart contract dependency check in the ledger
func (automation *SmartContractDependencyCheckerAutomation) logDependencyResult(contract common.SmartContract, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-dependency-%s", contract.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Dependency Check",
        Status:    result,
        Details:   fmt.Sprintf("Dependency check result for contract %s: %s", contract.ID, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with contract dependency check result for contract %s: %s\n", contract.ID, result)
}

// finalizeDependencyCheckCycle finalizes the contract dependency check cycle and logs the results in the ledger
func (automation *SmartContractDependencyCheckerAutomation) finalizeDependencyCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDependencyCheckCycle()
    if success {
        fmt.Println("Smart contract dependency check cycle finalized successfully.")
        automation.logDependencyCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing smart contract dependency check cycle.")
    }
}

// logDependencyCheckCycleFinalization logs the finalization of a smart contract dependency check cycle in the ledger
func (automation *SmartContractDependencyCheckerAutomation) logDependencyCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dependency-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Dependency Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with smart contract dependency check cycle finalization.")
}
