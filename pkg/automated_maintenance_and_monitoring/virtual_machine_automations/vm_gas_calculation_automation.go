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
    GasCalculationInterval = 3000 * time.Millisecond // Interval for gas calculation checks
    SubBlocksPerBlock      = 1000                    // Number of sub-blocks in a block
)

// VMGasCalculationAutomation automates the process of calculating gas for smart contract execution on the VM
type VMGasCalculationAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance   *ledger.Ledger               // Ledger to store gas calculation actions
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    gasCalculationCount int                       // Counter for gas calculation cycles
}

// NewVMGasCalculationAutomation initializes the automation for gas calculation on the VM
func NewVMGasCalculationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMGasCalculationAutomation {
    return &VMGasCalculationAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        gasCalculationCount: 0,
    }
}

// StartGasCalculationCheck starts the continuous loop for monitoring gas calculations
func (automation *VMGasCalculationAutomation) StartGasCalculationCheck() {
    ticker := time.NewTicker(GasCalculationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndCalculateGas()
        }
    }()
}

// monitorAndCalculateGas checks the contract executions and calculates gas used
func (automation *VMGasCalculationAutomation) monitorAndCalculateGas() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of active contract executions on the VM
    contractExecutions := automation.consensusSystem.GetActiveContractExecutions()

    for _, contract := range contractExecutions {
        fmt.Printf("Calculating gas for contract execution: %s.\n", contract.ID)
        automation.calculateGasForContract(contract)
    }

    automation.gasCalculationCount++
    fmt.Printf("Gas calculation cycle #%d executed.\n", automation.gasCalculationCount)

    if automation.gasCalculationCount%SubBlocksPerBlock == 0 {
        automation.finalizeGasCalculationCycle()
    }
}

// calculateGasForContract calculates the gas consumed by a specific contract execution
func (automation *VMGasCalculationAutomation) calculateGasForContract(contract common.SmartContractExecution) {
    // Encrypt contract data before calculating gas
    encryptedContractData := automation.encryptContractData(contract)

    // Calculate the gas used during the execution
    gasUsed := automation.consensusSystem.CalculateGasUsage(encryptedContractData)

    if gasUsed > 0 {
        fmt.Printf("Gas calculated for contract %s: %d units.\n", contract.ID, gasUsed)
        automation.logGasUsage(contract, gasUsed)
    } else {
        fmt.Printf("Error calculating gas for contract %s.\n", contract.ID)
    }
}

// finalizeGasCalculationCycle finalizes the gas calculation cycle and logs the result in the ledger
func (automation *VMGasCalculationAutomation) finalizeGasCalculationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeGasCalculationCycle()
    if success {
        fmt.Println("Gas calculation cycle finalized successfully.")
        automation.logGasCalculationCycleFinalization()
    } else {
        fmt.Println("Error finalizing gas calculation cycle.")
    }
}

// logGasUsage logs the gas usage for a specific contract execution into the ledger for traceability
func (automation *VMGasCalculationAutomation) logGasUsage(contract common.SmartContractExecution, gasUsed int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-usage-%s", contract.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Usage",
        Status:    "Calculated",
        Details:   fmt.Sprintf("Gas used for contract %s: %d units.", contract.ID, gasUsed),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with gas usage for contract %s.\n", contract.ID)
}

// logGasCalculationCycleFinalization logs the finalization of a gas calculation cycle into the ledger
func (automation *VMGasCalculationAutomation) logGasCalculationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-calculation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Calculation Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with gas calculation cycle finalization.")
}

// encryptContractData encrypts the contract data before processing gas calculations
func (automation *VMGasCalculationAutomation) encryptContractData(contract common.SmartContractExecution) common.SmartContractExecution {
    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        fmt.Println("Error encrypting contract data:", err)
        return contract
    }
    contract.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract
}

// ensureGasCalculationIntegrity ensures the integrity of gas calculations and triggers recalculation if necessary
func (automation *VMGasCalculationAutomation) ensureGasCalculationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateGasCalculationIntegrity()
    if !integrityValid {
        fmt.Println("Gas calculation integrity breach detected. Re-triggering gas calculations.")
        automation.monitorAndCalculateGas()
    } else {
        fmt.Println("Gas calculation integrity is valid.")
    }
}
