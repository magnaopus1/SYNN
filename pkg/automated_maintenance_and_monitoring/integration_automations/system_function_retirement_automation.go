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
    FunctionRetirementInterval = 5000 * time.Millisecond // Interval for checking and retiring system functions
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// SystemFunctionRetirementAutomation automates the retirement process for deprecated system functions
type SystemFunctionRetirementAutomation struct {
    consensusSystem         *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance          *ledger.Ledger               // Ledger to store retirement logs
    stateMutex              *sync.RWMutex                // Mutex for thread-safe access
    retirementCheckCount    int                          // Counter for retirement check cycles
}

// NewSystemFunctionRetirementAutomation initializes the automation for retiring system functions
func NewSystemFunctionRetirementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SystemFunctionRetirementAutomation {
    return &SystemFunctionRetirementAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        retirementCheckCount: 0,
    }
}

// StartFunctionRetirementMonitoring starts the continuous loop for monitoring and retiring system functions
func (automation *SystemFunctionRetirementAutomation) StartFunctionRetirementMonitoring() {
    ticker := time.NewTicker(FunctionRetirementInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndRetireFunctions()
        }
    }()
}

// monitorAndRetireFunctions checks for deprecated system functions and retires them from the system
func (automation *SystemFunctionRetirementAutomation) monitorAndRetireFunctions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of system functions marked for retirement
    functions, err := automation.consensusSystem.GetDeprecatedFunctions()
    if err != nil {
        fmt.Printf("Error fetching deprecated functions: %v\n", err)
        return
    }

    // Process each function for retirement
    for _, function := range functions {
        fmt.Printf("Retiring system function: %s\n", function.FunctionID)

        // Encrypt function data before retirement
        encryptedFunction, err := automation.encryptFunctionData(function)
        if err != nil {
            fmt.Printf("Error encrypting function data for %s: %v\n", function.FunctionID, err)
            automation.logFunctionRetirement(function, "Encryption Failed")
            continue
        }

        // Retire the function using consensus
        automation.retireFunction(encryptedFunction)
    }

    automation.retirementCheckCount++
    fmt.Printf("Function retirement check cycle #%d completed.\n", automation.retirementCheckCount)

    if automation.retirementCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRetirementCycle()
    }
}

// encryptFunctionData encrypts the function data before retirement
func (automation *SystemFunctionRetirementAutomation) encryptFunctionData(function common.SystemFunction) (common.SystemFunction, error) {
    fmt.Println("Encrypting system function data for retirement.")

    encryptedData, err := encryption.EncryptData(function)
    if err != nil {
        return function, fmt.Errorf("failed to encrypt function data: %v", err)
    }

    function.EncryptedData = encryptedData
    fmt.Println("System function data successfully encrypted for retirement.")
    return function, nil
}

// retireFunction retires the system function from the system through Synnergy Consensus
func (automation *SystemFunctionRetirementAutomation) retireFunction(function common.SystemFunction) {
    success := automation.consensusSystem.RetireSystemFunction(function)
    if success {
        fmt.Printf("Function %s retired successfully.\n", function.FunctionID)
        automation.logFunctionRetirement(function, "Retired Successfully")
    } else {
        fmt.Printf("Error retiring function %s.\n", function.FunctionID)
        automation.logFunctionRetirement(function, "Retirement Failed")
    }
}

// logFunctionRetirement logs the function retirement into the ledger for auditability
func (automation *SystemFunctionRetirementAutomation) logFunctionRetirement(function common.SystemFunction, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-retirement-%s", function.FunctionID),
        Timestamp: time.Now().Unix(),
        Type:      "System Function Retirement",
        Status:    status,
        Details:   fmt.Sprintf("Function %s: %s", function.FunctionID, status),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with retirement event for function %s: %s\n", function.FunctionID, status)
}

// finalizeRetirementCycle finalizes the function retirement cycle and logs the result in the ledger
func (automation *SystemFunctionRetirementAutomation) finalizeRetirementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeFunctionRetirementCycle()
    if success {
        fmt.Println("Function retirement cycle finalized successfully.")
        automation.logRetirementCycleFinalization()
    } else {
        fmt.Println("Error finalizing function retirement cycle.")
    }
}

// logRetirementCycleFinalization logs the finalization of the function retirement cycle in the ledger
func (automation *SystemFunctionRetirementAutomation) logRetirementCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-retirement-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Function Retirement Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with system function retirement cycle finalization.")
}
