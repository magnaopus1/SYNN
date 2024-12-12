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
    FunctionLinkingCheckInterval = 3500 * time.Millisecond // Interval for checking contract-to-system function linking
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
)

// ContractToSystemFunctionLinkingAutomation automates the process of ensuring correct linking between contract and system functions
type ContractToSystemFunctionLinkingAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store function linking logs
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    linkingCheckCount   int                          // Counter for linking check cycles
}

// NewContractToSystemFunctionLinkingAutomation initializes the automation for contract-to-system function linking
func NewContractToSystemFunctionLinkingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractToSystemFunctionLinkingAutomation {
    return &ContractToSystemFunctionLinkingAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        linkingCheckCount: 0,
    }
}

// StartFunctionLinkingCheck starts the continuous loop for checking the linking of contract and system functions
func (automation *ContractToSystemFunctionLinkingAutomation) StartFunctionLinkingCheck() {
    ticker := time.NewTicker(FunctionLinkingCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkContractToSystemFunctionLinking()
        }
    }()
}

// checkContractToSystemFunctionLinking verifies if contract functions are properly linked to their respective system functions
func (automation *ContractToSystemFunctionLinkingAutomation) checkContractToSystemFunctionLinking() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newContractFunctions := automation.consensusSystem.GetNewContractFunctions() // Fetch newly added or updated contract functions

    for _, function := range newContractFunctions {
        fmt.Printf("Checking function linking for contract function: %s\n", function.Name)
        isLinked := automation.verifyLinkingForFunction(function)

        if isLinked {
            fmt.Printf("Contract function %s is correctly linked to system functions.\n", function.Name)
            automation.logFunctionLinkingResult(function.Name, "Linked")
        } else {
            fmt.Printf("Contract function %s is not correctly linked to system functions.\n", function.Name)
            automation.logFunctionLinkingResult(function.Name, "Not Linked")
        }
    }

    automation.linkingCheckCount++
    fmt.Printf("Contract-to-system function linking check cycle #%d completed.\n", automation.linkingCheckCount)

    if automation.linkingCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeLinkingCheckCycle()
    }
}

// verifyLinkingForFunction checks the linkage between a specific contract function and the corresponding system functions
func (automation *ContractToSystemFunctionLinkingAutomation) verifyLinkingForFunction(function common.ContractFunction) bool {
    // Encrypt function data before verifying linking
    fmt.Printf("Encrypting contract function data for: %s\n", function.Name)

    encryptedFunctionData, err := encryption.EncryptData(function)
    if err != nil {
        fmt.Printf("Error encrypting contract function data for %s: %s\n", function.Name, err.Error())
        return false
    }

    function.EncryptedData = encryptedFunctionData
    fmt.Printf("Contract function data for %s encrypted successfully.\n", function.Name)

    // Verify linking through consensus and check if system functions are correctly mapped
    return automation.consensusSystem.VerifyContractToSystemFunctionLinking(function)
}

// logFunctionLinkingResult logs the result of the contract-to-system function linking check into the ledger
func (automation *ContractToSystemFunctionLinkingAutomation) logFunctionLinkingResult(functionName string, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("function-linking-check-%s", functionName),
        Timestamp: time.Now().Unix(),
        Type:      "Contract-to-System Function Linking",
        Status:    result,
        Details:   fmt.Sprintf("Function linking check result for contract function %s: %s", functionName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with function linking check result for function %s: %s.\n", functionName, result)
}

// finalizeLinkingCheckCycle finalizes the linking check cycle and logs the result in the ledger
func (automation *ContractToSystemFunctionLinkingAutomation) finalizeLinkingCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeLinkingCheckCycle()
    if success {
        fmt.Println("Contract-to-system function linking check cycle finalized successfully.")
        automation.logLinkingCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing contract-to-system function linking check cycle.")
    }
}

// logLinkingCheckCycleFinalization logs the finalization of a linking check cycle into the ledger
func (automation *ContractToSystemFunctionLinkingAutomation) logLinkingCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("linking-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Function Linking Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with contract-to-system function linking cycle finalization.")
}
