package automations

import (
    "fmt"
    "log"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "errors"
)

const (
    APIGenerationCheckInterval = 5000 * time.Millisecond // Interval for checking API generation for new system functions
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// DynamicAPIGenerationAutomation automates the generation of APIs for new system functions dynamically
type DynamicAPIGenerationAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger to store API generation logs
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    apiGenerationCount int                          // Counter for API generation check cycles
}

// NewDynamicAPIGenerationAutomation initializes the automation for dynamic API generation
func NewDynamicAPIGenerationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DynamicAPIGenerationAutomation {
    return &DynamicAPIGenerationAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        apiGenerationCount: 0,
    }
}

// StartAPIGenerationCheck starts the continuous loop for checking and generating APIs for new system functions
func (automation *DynamicAPIGenerationAutomation) StartAPIGenerationCheck() {
    ticker := time.NewTicker(APIGenerationCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkForNewSystemFunctions()
        }
    }()
}

// checkForNewSystemFunctions verifies if new system functions have been added and dynamically generates their APIs
func (automation *DynamicAPIGenerationAutomation) checkForNewSystemFunctions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newSystemFunctions, err := automation.consensusSystem.GetNewSystemFunctions()
    if err != nil {
        log.Printf("Error fetching new system functions: %v", err)
        return
    }

    for _, function := range newSystemFunctions {
        fmt.Printf("Generating API for system function: %s\n", function.Name)
        isGenerated, err := automation.generateAPIForSystemFunction(function)
        if err != nil {
            fmt.Printf("Failed to generate API for system function %s: %v\n", function.Name, err)
            automation.logAPIGenerationResult(function.Name, "Failed", err.Error())
        } else if isGenerated {
            fmt.Printf("API for system function %s generated successfully.\n", function.Name)
            automation.logAPIGenerationResult(function.Name, "Generated", "Success")
        }
    }

    automation.apiGenerationCount++
    fmt.Printf("API generation check cycle #%d completed.\n", automation.apiGenerationCount)

    if automation.apiGenerationCount%SubBlocksPerBlock == 0 {
        automation.finalizeAPIGenerationCycle()
    }
}

// generateAPIForSystemFunction generates the API for a specific system function
func (automation *DynamicAPIGenerationAutomation) generateAPIForSystemFunction(function common.SystemFunction) (bool, error) {
    // Validate function signature
    if !automation.validateFunctionSignature(function) {
        return false, errors.New("invalid function signature")
    }

    // Encrypt function data before generating the API
    fmt.Printf("Encrypting system function data for: %s\n", function.Name)

    encryptedFunctionData, err := encryption.EncryptData(function)
    if err != nil {
        return false, fmt.Errorf("error encrypting system function data for %s: %v", function.Name, err)
    }

    function.EncryptedData = encryptedFunctionData
    fmt.Printf("System function data for %s encrypted successfully.\n", function.Name)

    // Generate API endpoint through the Synnergy Consensus
    return automation.consensusSystem.GenerateAPIForSystemFunction(function)
}

// validateFunctionSignature ensures the function signature is valid before API generation
func (automation *DynamicAPIGenerationAutomation) validateFunctionSignature(function common.SystemFunction) bool {
    if function.Name == "" || len(function.Parameters) == 0 {
        fmt.Printf("Invalid function signature for %s. Missing name or parameters.\n", function.Name)
        return false
    }

    // Additional validation logic can be added here if needed
    return true
}

// logAPIGenerationResult logs the result of the API generation process into the ledger
func (automation *DynamicAPIGenerationAutomation) logAPIGenerationResult(functionName, result, details string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("api-generation-%s", functionName),
        Timestamp: time.Now().Unix(),
        Type:      "API Generation",
        Status:    result,
        Details:   fmt.Sprintf("API generation result for system function %s: %s. Details: %s", functionName, result, details),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with API generation result for system function %s: %s.\n", functionName, result)
}

// finalizeAPIGenerationCycle finalizes the API generation check cycle and logs the result in the ledger
func (automation *DynamicAPIGenerationAutomation) finalizeAPIGenerationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeAPIGenerationCycle()
    if success {
        fmt.Println("API generation check cycle finalized successfully.")
        automation.logAPIGenerationCycleFinalization()
    } else {
        fmt.Println("Error finalizing API generation check cycle.")
    }
}

// logAPIGenerationCycleFinalization logs the finalization of an API generation check cycle into the ledger
func (automation *DynamicAPIGenerationAutomation) logAPIGenerationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("api-generation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "API Generation Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with API generation cycle finalization.")
}
