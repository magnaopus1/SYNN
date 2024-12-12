package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    ContractCompilationInterval = 5000 * time.Millisecond // Interval for checking and compiling contracts
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
)

// MultiLanguageContractCompilationAutomation automates the process of compiling contracts written in multiple languages
type MultiLanguageContractCompilationAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to Synnergy Consensus for contract compilation validation
    ledgerInstance       *ledger.Ledger               // Ledger to store contract compilation actions
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    compilationCycleCount int                         // Counter for contract compilation cycles
}

// NewMultiLanguageContractCompilationAutomation initializes the automation for multi-language contract compilation
func NewMultiLanguageContractCompilationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *MultiLanguageContractCompilationAutomation {
    return &MultiLanguageContractCompilationAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        compilationCycleCount: 0,
    }
}

// StartContractCompilationCheck starts the continuous loop for monitoring and compiling smart contracts across languages
func (automation *MultiLanguageContractCompilationAutomation) StartContractCompilationCheck() {
    ticker := time.NewTicker(ContractCompilationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndCompileContracts()
        }
    }()
}

// monitorAndCompileContracts checks for pending contracts and triggers compilation, validation, and logging
func (automation *MultiLanguageContractCompilationAutomation) monitorAndCompileContracts() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch queued smart contracts pending compilation across different languages
    contractQueue := automation.consensusSystem.GetPendingContractQueue()

    for _, contract := range contractQueue {
        fmt.Printf("Compiling contract %s in language %s.\n", contract.ContractID, contract.Language)
        automation.compileAndValidateContract(contract)
    }

    automation.compilationCycleCount++
    if automation.compilationCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeCompilationCycle()
    }
}

// compileAndValidateContract compiles and validates the smart contract, ensuring security and integrity
func (automation *MultiLanguageContractCompilationAutomation) compileAndValidateContract(contract common.Contract) {
    // Encrypt contract code before validation
    encryptedContract := automation.encryptContract(contract)

    // Compile contract through Synnergy Consensus system
    compilationSuccess := automation.consensusSystem.CompileContract(encryptedContract)

    if compilationSuccess {
        fmt.Printf("Contract compilation successful for contract %s in language %s.\n", contract.ContractID, contract.Language)
        automation.validateCompiledContract(encryptedContract)
        automation.logContractCompilationEvent(contract)
    } else {
        fmt.Printf("Error compiling contract %s in language %s.\n", contract.ContractID, contract.Language)
    }
}

// validateCompiledContract validates the compiled contract through Synnergy Consensus
func (automation *MultiLanguageContractCompilationAutomation) validateCompiledContract(contract common.Contract) {
    validationSuccess := automation.consensusSystem.ValidateCompiledContract(contract)

    if validationSuccess {
        fmt.Printf("Contract validation successful for contract %s in language %s.\n", contract.ContractID, contract.Language)
    } else {
        fmt.Printf("Error validating compiled contract %s in language %s.\n", contract.ContractID, contract.Language)
    }
}

// finalizeCompilationCycle finalizes the contract compilation cycle and logs the result in the ledger
func (automation *MultiLanguageContractCompilationAutomation) finalizeCompilationCycle() {
    success := automation.consensusSystem.FinalizeCompilationCycle()
    if success {
        fmt.Println("Contract compilation cycle finalized successfully.")
        automation.logCompilationCycleFinalization()
    } else {
        fmt.Println("Error finalizing contract compilation cycle.")
    }
}

// logContractCompilationEvent logs the contract compilation event into the ledger for traceability
func (automation *MultiLanguageContractCompilationAutomation) logContractCompilationEvent(contract common.Contract) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-compilation-%s", contract.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Compilation",
        Status:    "Compiled and Validated",
        Details:   fmt.Sprintf("Contract %s in language %s was successfully compiled and validated.", contract.ContractID, contract.Language),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with contract compilation event for contract %s in language %s.\n", contract.ContractID, contract.Language)
}

// logCompilationCycleFinalization logs the finalization of a contract compilation cycle into the ledger
func (automation *MultiLanguageContractCompilationAutomation) logCompilationCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-compilation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Compilation Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with contract compilation cycle finalization.")
}

// encryptContract encrypts the contract code before compilation and validation
func (automation *MultiLanguageContractCompilationAutomation) encryptContract(contract common.Contract) common.Contract {
    encryptedData, err := encryption.EncryptData(contract.Code)
    if err != nil {
        fmt.Println("Error encrypting contract data:", err)
        return contract
    }
    contract.EncryptedCode = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract
}

// ensureCompilationIntegrity checks the integrity of the contract compilation process and triggers recompilation if necessary
func (automation *MultiLanguageContractCompilationAutomation) ensureCompilationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateCompilationIntegrity()
    if !integrityValid {
        fmt.Println("Compilation integrity breach detected. Re-triggering contract compilation checks.")
        automation.monitorAndCompileContracts()
    } else {
        fmt.Println("Compilation integrity is valid.")
    }
}
