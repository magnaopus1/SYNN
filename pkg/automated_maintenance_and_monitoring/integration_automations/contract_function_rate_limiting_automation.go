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
    RateLimitingCheckInterval = 3000 * time.Millisecond // Interval for checking contract function rate limiting
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// ContractFunctionRateLimitingAutomation automates the process of enforcing rate limits on contract functions
type ContractFunctionRateLimitingAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store rate limiting logs
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    rateLimitingCheckCount int                     // Counter for rate limiting check cycles
}

// NewContractFunctionRateLimitingAutomation initializes the automation for contract function rate limiting
func NewContractFunctionRateLimitingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractFunctionRateLimitingAutomation {
    return &ContractFunctionRateLimitingAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        rateLimitingCheckCount: 0,
    }
}

// StartRateLimitingCheck starts the continuous loop for checking rate limiting of contract functions
func (automation *ContractFunctionRateLimitingAutomation) StartRateLimitingCheck() {
    ticker := time.NewTicker(RateLimitingCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkContractFunctionRateLimiting()
        }
    }()
}

// checkContractFunctionRateLimiting checks if contract functions are exceeding predefined rate limits
func (automation *ContractFunctionRateLimitingAutomation) checkContractFunctionRateLimiting() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    contractFunctions := automation.consensusSystem.GetRunningContractFunctions() // Fetch currently running contract functions

    for _, function := range contractFunctions {
        fmt.Printf("Checking rate limit for contract function: %s\n", function.Name)
        isWithinLimit := automation.enforceRateLimitingForFunction(function)

        if isWithinLimit {
            fmt.Printf("Contract function %s is within the rate limit.\n", function.Name)
            automation.logRateLimitingResult(function.Name, "Within Limit")
        } else {
            fmt.Printf("Contract function %s exceeded the rate limit.\n", function.Name)
            automation.logRateLimitingResult(function.Name, "Exceeded Limit")
        }
    }

    automation.rateLimitingCheckCount++
    fmt.Printf("Contract function rate limiting check cycle #%d completed.\n", automation.rateLimitingCheckCount)

    if automation.rateLimitingCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRateLimitingCycle()
    }
}

// enforceRateLimitingForFunction checks and enforces the rate limiting for a specific contract function
func (automation *ContractFunctionRateLimitingAutomation) enforceRateLimitingForFunction(function common.ContractFunction) bool {
    // Encrypt function execution data before checking rate limits
    fmt.Printf("Encrypting contract function data for: %s\n", function.Name)

    encryptedFunctionData, err := encryption.EncryptData(function)
    if err != nil {
        fmt.Printf("Error encrypting contract function data for %s: %s\n", function.Name, err.Error())
        return false
    }

    function.EncryptedData = encryptedFunctionData
    fmt.Printf("Contract function data for %s encrypted successfully.\n", function.Name)

    // Verify and enforce rate limits through consensus
    return automation.consensusSystem.VerifyContractFunctionRateLimit(function)
}

// logRateLimitingResult logs the result of the rate limiting check into the ledger
func (automation *ContractFunctionRateLimitingAutomation) logRateLimitingResult(functionName string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rate-limiting-check-%s", functionName),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Function Rate Limiting",
        Status:    result,
        Details:   fmt.Sprintf("Rate limiting check result for contract function %s: %s", functionName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with rate limiting check result for function %s: %s.\n", functionName, result)
}

// finalizeRateLimitingCycle finalizes the rate limiting check cycle and logs the result in the ledger
func (automation *ContractFunctionRateLimitingAutomation) finalizeRateLimitingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeRateLimitingCycle()
    if success {
        fmt.Println("Contract function rate limiting check cycle finalized successfully.")
        automation.logRateLimitingCycleFinalization()
    } else {
        fmt.Println("Error finalizing contract function rate limiting check cycle.")
    }
}

// logRateLimitingCycleFinalization logs the finalization of a rate limiting check cycle into the ledger
func (automation *ContractFunctionRateLimitingAutomation) logRateLimitingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("rate-limiting-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Rate Limiting Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with contract function rate limiting cycle finalization.")
}
