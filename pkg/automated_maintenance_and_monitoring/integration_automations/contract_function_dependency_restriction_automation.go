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
    DependencyRestrictionCheckInterval = 3000 * time.Millisecond // Interval for checking contract function dependency restrictions
    SubBlocksPerBlock                  = 1000                    // Number of sub-blocks in a block
)

// ContractFunctionDependencyRestrictionAutomation automates the restriction enforcement for contract function dependencies
type ContractFunctionDependencyRestrictionAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store restriction enforcement logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    restrictionCheckCount  int                          // Counter for restriction check cycles
}

// NewContractFunctionDependencyRestrictionAutomation initializes the automation for contract function dependency restrictions
func NewContractFunctionDependencyRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractFunctionDependencyRestrictionAutomation {
    return &ContractFunctionDependencyRestrictionAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        restrictionCheckCount: 0,
    }
}

// StartDependencyRestrictionCheck starts the continuous loop for checking contract function dependency restrictions
func (automation *ContractFunctionDependencyRestrictionAutomation) StartDependencyRestrictionCheck() {
    ticker := time.NewTicker(DependencyRestrictionCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkContractFunctionDependencies()
        }
    }()
}

// checkContractFunctionDependencies checks for contract functions violating dependency restrictions
func (automation *ContractFunctionDependencyRestrictionAutomation) checkContractFunctionDependencies() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newContractFunctions := automation.consensusSystem.GetNewContractFunctions() // Fetch newly added or updated contract functions

    for _, function := range newContractFunctions {
        fmt.Printf("Checking dependency restrictions for contract function: %s\n", function.Name)
        isValid := automation.checkDependenciesForFunction(function)

        if isValid {
            fmt.Printf("Contract function %s passed dependency restriction check.\n", function.Name)
            automation.logDependencyCheckResult(function.Name, "Passed")
        } else {
            fmt.Printf("Contract function %s failed dependency restriction check.\n", function.Name)
            automation.logDependencyCheckResult(function.Name, "Failed")
        }
    }

    automation.restrictionCheckCount++
    fmt.Printf("Contract function dependency restriction check cycle #%d completed.\n", automation.restrictionCheckCount)

    if automation.restrictionCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRestrictionCheckCycle()
    }
}

// checkDependenciesForFunction runs the dependency check for a specific contract function
func (automation *ContractFunctionDependencyRestrictionAutomation) checkDependenciesForFunction(function common.ContractFunction) bool {
    // Encrypt function data before checking dependencies
    fmt.Printf("Encrypting contract function data for: %s\n", function.Name)

    encryptedFunctionData, err := encryption.EncryptData(function)
    if err != nil {
        fmt.Printf("Error encrypting contract function data for %s: %s\n", function.Name, err.Error())
        return false
    }

    function.EncryptedData = encryptedFunctionData
    fmt.Printf("Contract function data for %s encrypted successfully.\n", function.Name)

    // Verify dependency compliance through consensus
    return automation.consensusSystem.VerifyContractFunctionDependencies(function)
}

// logDependencyCheckResult logs the result of the contract function dependency restriction check into the ledger
func (automation *ContractFunctionDependencyRestrictionAutomation) logDependencyCheckResult(functionName string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dependency-restriction-check-%s", functionName),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Function Dependency Check",
        Status:    result,
        Details:   fmt.Sprintf("Dependency check result for contract function %s: %s", functionName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with dependency check result for function %s: %s.\n", functionName, result)
}

// finalizeRestrictionCheckCycle finalizes the restriction check cycle and logs the result in the ledger
func (automation *ContractFunctionDependencyRestrictionAutomation) finalizeRestrictionCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDependencyRestrictionCheckCycle()
    if success {
        fmt.Println("Contract function dependency restriction check cycle finalized successfully.")
        automation.logRestrictionCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing contract function dependency restriction check cycle.")
    }
}

// logRestrictionCheckCycleFinalization logs the finalization of a restriction check cycle into the ledger
func (automation *ContractFunctionDependencyRestrictionAutomation) logRestrictionCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("dependency-restriction-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Dependency Restriction Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with contract function dependency restriction check cycle finalization.")
}
