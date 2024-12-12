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
    FunctionRegistrationInterval = 4000 * time.Millisecond // Interval for checking and registering new system functions
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks per block
)

// SystemFunctionRegistrationAutomation automates the process of registering new system functions
type SystemFunctionRegistrationAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store function registration logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    registrationCheckCount int                          // Counter for registration check cycles
}

// NewSystemFunctionRegistrationAutomation initializes the automation for registering new system functions
func NewSystemFunctionRegistrationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SystemFunctionRegistrationAutomation {
    return &SystemFunctionRegistrationAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        registrationCheckCount: 0,
    }
}

// StartFunctionRegistrationMonitoring starts the continuous loop for monitoring and registering new system functions
func (automation *SystemFunctionRegistrationAutomation) StartFunctionRegistrationMonitoring() {
    ticker := time.NewTicker(FunctionRegistrationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndRegisterFunctions()
        }
    }()
}

// monitorAndRegisterFunctions checks for new system functions and registers them into the system
func (automation *SystemFunctionRegistrationAutomation) monitorAndRegisterFunctions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of new system functions awaiting registration
    functions, err := automation.consensusSystem.GetPendingFunctionRegistrations()
    if err != nil {
        fmt.Printf("Error fetching new function registrations: %v\n", err)
        return
    }

    // Process each function for registration
    for _, function := range functions {
        fmt.Printf("Registering new system function: %s\n", function.FunctionID)

        // Encrypt function data before registration
        encryptedFunction, err := automation.encryptFunctionData(function)
        if err != nil {
            fmt.Printf("Error encrypting function data for %s: %v\n", function.FunctionID, err)
            automation.logFunctionRegistration(function, "Encryption Failed")
            continue
        }

        // Register the function using consensus
        automation.registerFunction(encryptedFunction)
    }

    automation.registrationCheckCount++
    fmt.Printf("Function registration check cycle #%d completed.\n", automation.registrationCheckCount)

    if automation.registrationCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRegistrationCycle()
    }
}

// encryptFunctionData encrypts the function data before registering it to the system
func (automation *SystemFunctionRegistrationAutomation) encryptFunctionData(function common.SystemFunction) (common.SystemFunction, error) {
    fmt.Println("Encrypting system function data.")

    encryptedData, err := encryption.EncryptData(function)
    if err != nil {
        return function, fmt.Errorf("failed to encrypt function data: %v", err)
    }

    function.EncryptedData = encryptedData
    fmt.Println("System function data successfully encrypted.")
    return function, nil
}

// registerFunction registers the new system function into the system through Synnergy Consensus
func (automation *SystemFunctionRegistrationAutomation) registerFunction(function common.SystemFunction) {
    success := automation.consensusSystem.RegisterSystemFunction(function)
    if success {
        fmt.Printf("Function %s registered successfully.\n", function.FunctionID)
        automation.logFunctionRegistration(function, "Registered Successfully")
    } else {
        fmt.Printf("Error registering function %s.\n", function.FunctionID)
        automation.logFunctionRegistration(function, "Registration Failed")
    }
}

// logFunctionRegistration logs the function registration into the ledger for auditability
func (automation *SystemFunctionRegistrationAutomation) logFunctionRegistration(function common.SystemFunction, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-registration-%s", function.FunctionID),
        Timestamp: time.Now().Unix(),
        Type:      "System Function Registration",
        Status:    status,
        Details:   fmt.Sprintf("Function %s: %s", function.FunctionID, status),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with registration event for function %s: %s\n", function.FunctionID, status)
}

// finalizeRegistrationCycle finalizes the function registration cycle and logs the result in the ledger
func (automation *SystemFunctionRegistrationAutomation) finalizeRegistrationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeFunctionRegistrationCycle()
    if success {
        fmt.Println("Function registration cycle finalized successfully.")
        automation.logRegistrationCycleFinalization()
    } else {
        fmt.Println("Error finalizing function registration cycle.")
    }
}

// logRegistrationCycleFinalization logs the finalization of the function registration cycle in the ledger
func (automation *SystemFunctionRegistrationAutomation) logRegistrationCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-registration-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Function Registration Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with system function registration cycle finalization.")
}
