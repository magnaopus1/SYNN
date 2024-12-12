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
    FunctionCompatibilityCheckInterval = 2500 * time.Millisecond // Interval for checking contract function compatibility
    SubBlocksPerBlock                  = 1000                    // Number of sub-blocks in a block
)

// ContractFunctionCompatibilityCheckerAutomation automates the process of checking the compatibility of newly added contract functions
type ContractFunctionCompatibilityCheckerAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store compatibility check logs
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    compatibilityCheckCount int                    // Counter for compatibility check cycles
}

// NewContractFunctionCompatibilityCheckerAutomation initializes the automation for contract function compatibility checks
func NewContractFunctionCompatibilityCheckerAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractFunctionCompatibilityCheckerAutomation {
    return &ContractFunctionCompatibilityCheckerAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        compatibilityCheckCount: 0,
    }
}

// StartFunctionCompatibilityCheck starts the continuous loop for checking the compatibility of contract functions
func (automation *ContractFunctionCompatibilityCheckerAutomation) StartFunctionCompatibilityCheck() {
    ticker := time.NewTicker(FunctionCompatibilityCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkContractFunctionCompatibility()
        }
    }()
}

// checkContractFunctionCompatibility checks if new or updated contract functions are compatible with the existing system
func (automation *ContractFunctionCompatibilityCheckerAutomation) checkContractFunctionCompatibility() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newContractFunctions := automation.consensusSystem.GetNewContractFunctions() // Fetch newly added or updated contract functions

    for _, function := range newContractFunctions {
        fmt.Printf("Checking compatibility for contract function: %s\n", function.Name)
        isCompatible := automation.checkCompatibilityForFunction(function)

        if isCompatible {
            fmt.Printf("Contract function %s is compatible.\n", function.Name)
            automation.logCompatibilityCheckResult(function.Name, "Compatible")
        } else {
            fmt.Printf("Contract function %s is incompatible.\n", function.Name)
            automation.logCompatibilityCheckResult(function.Name, "Incompatible")
        }
    }

    automation.compatibilityCheckCount++
    fmt.Printf("Contract function compatibility check cycle #%d completed.\n", automation.compatibilityCheckCount)

    if automation.compatibilityCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeCompatibilityCheckCycle()
    }
}

// checkCompatibilityForFunction runs the compatibility check for a specific contract function
func (automation *ContractFunctionCompatibilityCheckerAutomation) checkCompatibilityForFunction(function common.ContractFunction) bool {
    // Simulate a compatibility check
    fmt.Printf("Running compatibility check on function: %s\n", function.Name)
    
    // Encrypt function data before verifying compatibility
    encryptedFunctionData, err := encryption.EncryptData(function)
    if err != nil {
        fmt.Printf("Error encrypting contract function data for %s: %s\n", function.Name, err.Error())
        return false
    }

    function.EncryptedData = encryptedFunctionData
    fmt.Printf("Contract function data for %s encrypted successfully.\n", function.Name)

    // Verify compatibility through consensus (example logic here)
    return automation.consensusSystem.VerifyContractFunctionCompatibility(function)
}

// logCompatibilityCheckResult logs the result of the contract function compatibility check into the ledger
func (automation *ContractFunctionCompatibilityCheckerAutomation) logCompatibilityCheckResult(functionName string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-compatibility-%s", functionName),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Function Compatibility Check",
        Status:    result,
        Details:   fmt.Sprintf("Compatibility check result for contract function %s: %s", functionName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with compatibility check result for function %s: %s.\n", functionName, result)
}

// finalizeCompatibilityCheckCycle finalizes the compatibility check cycle and logs the result in the ledger
func (automation *ContractFunctionCompatibilityCheckerAutomation) finalizeCompatibilityCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeCompatibilityCheckCycle()
    if success {
        fmt.Println("Contract function compatibility check cycle finalized successfully.")
        automation.logCompatibilityCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing contract function compatibility check cycle.")
    }
}

// logCompatibilityCheckCycleFinalization logs the finalization of a compatibility check cycle into the ledger
func (automation *ContractFunctionCompatibilityCheckerAutomation) logCompatibilityCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("compatibility-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Compatibility Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with contract function compatibility check cycle finalization.")
}
