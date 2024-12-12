package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    SmartContractExecutionInterval = 1000 * time.Millisecond // Interval for checking and executing smart contracts on the VM
    SubBlocksPerBlock              = 1000                    // Number of sub-blocks in a block
)

// SmartContractExecutionAutomation automates the execution of smart contracts on the VM
type SmartContractExecutionAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to Synnergy Consensus system
    ledgerInstance      *ledger.Ledger               // Ledger for logging contract executions
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    executionCycleCount int                          // Counter for execution check cycles
}

// NewSmartContractExecutionAutomation initializes the automation for smart contract execution
func NewSmartContractExecutionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartContractExecutionAutomation {
    return &SmartContractExecutionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        executionCycleCount: 0,
    }
}

// StartExecutionLoop starts the continuous loop for monitoring and executing smart contracts on the VM
func (automation *SmartContractExecutionAutomation) StartExecutionLoop() {
    ticker := time.NewTicker(SmartContractExecutionInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndExecuteSmartContracts()
        }
    }()
}

// monitorAndExecuteSmartContracts checks for pending smart contracts and executes them on the virtual machine
func (automation *SmartContractExecutionAutomation) monitorAndExecuteSmartContracts() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the pending smart contracts from the consensus system
    pendingContracts := automation.consensusSystem.GetPendingSmartContracts()

    for _, contract := range pendingContracts {
        fmt.Printf("Executing smart contract %s.\n", contract.ContractAddress)
        automation.executeSmartContract(contract)
    }

    automation.executionCycleCount++
    if automation.executionCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeExecutionCycle()
    }
}

// executeSmartContract handles the actual execution of the smart contract on the VM
func (automation *SmartContractExecutionAutomation) executeSmartContract(contract common.SmartContract) {
    // Encrypt contract data before executing
    encryptedContract := automation.encryptContractData(contract)

    // Execute the contract on the VM
    success := automation.consensusSystem.ExecuteSmartContract(encryptedContract)

    if success {
        fmt.Printf("Smart contract %s executed successfully.\n", contract.ContractAddress)
        automation.logContractExecution(contract)
    } else {
        fmt.Printf("Error executing smart contract %s.\n", contract.ContractAddress)
    }
}

// logContractExecution logs the execution of the smart contract into the ledger for traceability
func (automation *SmartContractExecutionAutomation) logContractExecution(contract common.SmartContract) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("smart-contract-execution-%s", contract.ContractAddress),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Execution",
        Status:    "Executed",
        Details:   fmt.Sprintf("Smart contract %s executed successfully.", contract.ContractAddress),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with execution of smart contract %s.\n", contract.ContractAddress)
}

// finalizeExecutionCycle finalizes the smart contract execution cycle and logs the result in the ledger
func (automation *SmartContractExecutionAutomation) finalizeExecutionCycle() {
    success := automation.consensusSystem.FinalizeExecutionCycle()
    if success {
        fmt.Println("Smart contract execution cycle finalized successfully.")
        automation.logExecutionCycleFinalization()
    } else {
        fmt.Println("Error finalizing smart contract execution cycle.")
    }
}

// logExecutionCycleFinalization logs the finalization of an execution cycle into the ledger
func (automation *SmartContractExecutionAutomation) logExecutionCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("execution-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Execution Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with smart contract execution cycle finalization.")
}

// encryptContractData encrypts the smart contract data before execution
func (automation *SmartContractExecutionAutomation) encryptContractData(contract common.SmartContract) common.SmartContract {
    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        fmt.Println("Error encrypting contract data:", err)
        return contract
    }
    contract.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract
}

// ensureExecutionIntegrity checks the integrity of executed smart contracts and triggers re-execution if necessary
func (automation *SmartContractExecutionAutomation) ensureExecutionIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateExecutionIntegrity()
    if !integrityValid {
        fmt.Println("Execution data integrity breach detected. Re-triggering execution checks.")
        automation.monitorAndExecuteSmartContracts()
    } else {
        fmt.Println("Smart contract execution data integrity is valid.")
    }
}
