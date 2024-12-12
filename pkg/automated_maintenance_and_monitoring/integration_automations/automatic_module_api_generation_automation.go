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
    APIGenerationInterval   = 2000 * time.Millisecond // Interval for checking and generating new APIs
    SubBlocksPerBlock       = 1000                    // Number of sub-blocks in a block
)

// AutomaticModuleAPIGenerationAutomation automates the API generation for newly added modules
type AutomaticModuleAPIGenerationAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger to store API generation details
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    apiGenerationCount int                        // Counter for API generation cycles
}

// NewAutomaticModuleAPIGenerationAutomation initializes the automation for API generation
func NewAutomaticModuleAPIGenerationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AutomaticModuleAPIGenerationAutomation {
    return &AutomaticModuleAPIGenerationAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        apiGenerationCount: 0,
    }
}

// StartAPIGeneration starts the continuous loop for checking and generating APIs for new modules
func (automation *AutomaticModuleAPIGenerationAutomation) StartAPIGeneration() {
    ticker := time.NewTicker(APIGenerationInterval)

    go func() {
        for range ticker.C {
            automation.generateAPIsForNewModules()
        }
    }()
}

// generateAPIsForNewModules checks for new modules and generates corresponding APIs
func (automation *AutomaticModuleAPIGenerationAutomation) generateAPIsForNewModules() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch newly added modules from the consensus system
    newModules := automation.consensusSystem.GetNewModules()

    for _, module := range newModules {
        fmt.Printf("Generating API for module: %s\n", module.Name)
        apiGenerated := automation.generateAPIForModule(module)

        if apiGenerated {
            fmt.Printf("API for module %s generated successfully.\n", module.Name)
            automation.logAPIGeneration(module.Name, "Success")
        } else {
            fmt.Printf("Error generating API for module %s.\n", module.Name)
            automation.logAPIGeneration(module.Name, "Failed")
        }
    }

    automation.apiGenerationCount++
    fmt.Printf("API generation cycle #%d completed.\n", automation.apiGenerationCount)

    if automation.apiGenerationCount%SubBlocksPerBlock == 0 {
        automation.finalizeAPIGenerationCycle()
    }
}

// generateAPIForModule generates an API for the provided module and encrypts the API data
func (automation *AutomaticModuleAPIGenerationAutomation) generateAPIForModule(module common.Module) bool {
    // Simulate API generation for the module
    fmt.Printf("Starting encryption for module API: %s\n", module.Name)
    
    // Encrypt API data before finalizing the API generation
    encryptedAPIdata, err := encryption.EncryptData(module)
    if err != nil {
        fmt.Printf("Error encrypting API data for module %s: %s\n", module.Name, err.Error())
        return false
    }

    module.EncryptedAPIdata = encryptedAPIdata
    fmt.Printf("API data for module %s successfully encrypted.\n", module.Name)

    // Add API to the system (this could involve adding routes/endpoints based on your implementation)
    return automation.consensusSystem.RegisterAPIForModule(module)
}

// logAPIGeneration logs the API generation result into the ledger for traceability
func (automation *AutomaticModuleAPIGenerationAutomation) logAPIGeneration(moduleName string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("api-generation-%s", moduleName),
        Timestamp: time.Now().Unix(),
        Type:      "API Generation",
        Status:    result,
        Details:   fmt.Sprintf("API generation for module %s: %s", moduleName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with API generation result for module %s: %s.\n", moduleName, result)
}

// finalizeAPIGenerationCycle finalizes the API generation cycle and logs the result in the ledger
func (automation *AutomaticModuleAPIGenerationAutomation) finalizeAPIGenerationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeAPIGenerationCycle()
    if success {
        fmt.Println("API generation cycle finalized successfully.")
        automation.logAPIGenerationCycleFinalization()
    } else {
        fmt.Println("Error finalizing API generation cycle.")
    }
}

// logAPIGenerationCycleFinalization logs the finalization of an API generation cycle into the ledger
func (automation *AutomaticModuleAPIGenerationAutomation) logAPIGenerationCycleFinalization() {
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
