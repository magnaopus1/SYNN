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
    SandboxExecutionInterval = 5000 * time.Millisecond // Interval for checking VM sandbox execution
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks in a block
    MaxSandboxExecutionTime  = 300                     // Max execution time for sandboxed contracts
)

// VMSandboxExecutionAutomation manages the sandbox execution of VMs
type VMSandboxExecutionAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging sandbox actions
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    sandboxCheckCount  int                          // Counter for sandbox execution check cycles
}

// NewVMSandboxExecutionAutomation initializes the automation for VM sandbox execution
func NewVMSandboxExecutionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMSandboxExecutionAutomation {
    return &VMSandboxExecutionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        sandboxCheckCount:  0,
    }
}

// StartSandboxExecutionCheck starts the continuous loop for monitoring and executing VMs in the sandbox
func (automation *VMSandboxExecutionAutomation) StartSandboxExecutionCheck() {
    ticker := time.NewTicker(SandboxExecutionInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndExecuteInSandbox()
        }
    }()
}

// monitorAndExecuteInSandbox checks the VM contracts queued for sandbox execution and runs them
func (automation *VMSandboxExecutionAutomation) monitorAndExecuteInSandbox() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the pending VM contracts for sandbox execution
    pendingContracts := automation.consensusSystem.GetPendingSandboxContracts()

    if len(pendingContracts) > 0 {
        for _, contract := range pendingContracts {
            fmt.Printf("Executing contract %s in sandbox.\n", contract.ID)
            automation.executeContractInSandbox(contract)
        }
    } else {
        fmt.Println("No contracts pending for sandbox execution.")
    }

    automation.sandboxCheckCount++
    fmt.Printf("Sandbox execution check cycle #%d executed.\n", automation.sandboxCheckCount)

    if automation.sandboxCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeSandboxExecutionCycle()
    }
}

// executeContractInSandbox executes a smart contract in the sandbox environment
func (automation *VMSandboxExecutionAutomation) executeContractInSandbox(contract common.SmartContract) {
    // Encrypt contract data before execution
    encryptedContract := automation.encryptContractData(contract)

    // Trigger sandbox execution through the Synnergy Consensus system
    executionSuccess := automation.consensusSystem.TriggerSandboxExecution(encryptedContract)

    if executionSuccess {
        fmt.Printf("Contract %s executed successfully in the sandbox.\n", contract.ID)
        automation.logSandboxExecution(contract)
    } else {
        fmt.Printf("Error executing contract %s in the sandbox.\n", contract.ID)
    }
}

// finalizeSandboxExecutionCycle finalizes the sandbox execution cycle and logs the result in the ledger
func (automation *VMSandboxExecutionAutomation) finalizeSandboxExecutionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeSandboxExecutionCycle()
    if success {
        fmt.Println("Sandbox execution check cycle finalized successfully.")
        automation.logSandboxExecutionCycleFinalization()
    } else {
        fmt.Println("Error finalizing sandbox execution check cycle.")
    }
}

// logSandboxExecution logs the sandbox execution action for a specific contract into the ledger
func (automation *VMSandboxExecutionAutomation) logSandboxExecution(contract common.SmartContract) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-execution-%s", contract.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Execution",
        Status:    "Executed",
        Details:   fmt.Sprintf("Contract %s successfully executed in the sandbox.", contract.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sandbox execution event for contract %s.\n", contract.ID)
}

// logSandboxExecutionCycleFinalization logs the finalization of a sandbox execution check cycle into the ledger
func (automation *VMSandboxExecutionAutomation) logSandboxExecutionCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-execution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Execution Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with sandbox execution cycle finalization.")
}

// encryptContractData encrypts the smart contract data before executing it in the sandbox
func (automation *VMSandboxExecutionAutomation) encryptContractData(contract common.SmartContract) common.SmartContract {
    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        fmt.Println("Error encrypting contract data:", err)
        return contract
    }

    contract.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract
}

// ensureSandboxExecutionIntegrity checks the integrity of sandbox execution data and re-executes if necessary
func (automation *VMSandboxExecutionAutomation) ensureSandboxExecutionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateSandboxExecutionIntegrity()
    if !integrityValid {
        fmt.Println("Sandbox execution integrity breach detected. Re-triggering sandbox execution.")
        automation.monitorAndExecuteInSandbox()
    } else {
        fmt.Println("Sandbox execution integrity is valid.")
    }
}
