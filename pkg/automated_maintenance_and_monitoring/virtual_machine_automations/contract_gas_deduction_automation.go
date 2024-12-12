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
    GasDeductionCheckInterval = 3000 * time.Millisecond // Interval for checking contract gas usage
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// ContractGasDeductionAutomation automates the process of deducting gas from contract executions
type ContractGasDeductionAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance      *ledger.Ledger               // Ledger for logging gas deductions
    stateMutex          *sync.RWMutex                // Mutex for thread-safe state access
    gasDeductionCheckCount int                       // Counter for gas deduction checks
}

// NewContractGasDeductionAutomation initializes the automation for contract gas deduction
func NewContractGasDeductionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractGasDeductionAutomation {
    return &ContractGasDeductionAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        gasDeductionCheckCount: 0,
    }
}

// StartGasDeductionCheck starts the continuous loop for checking and deducting gas from contract executions
func (automation *ContractGasDeductionAutomation) StartGasDeductionCheck() {
    ticker := time.NewTicker(GasDeductionCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndDeductGas()
        }
    }()
}

// monitorAndDeductGas checks the gas consumption of contract executions and deducts gas accordingly
func (automation *ContractGasDeductionAutomation) monitorAndDeductGas() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the gas usage data for all active contracts
    contractGasUsage := automation.consensusSystem.GetContractGasUsage()

    for _, usage := range contractGasUsage {
        if usage.GasUsed > 0 {
            fmt.Printf("Deducting gas for contract %s. Gas used: %d\n", usage.ContractID, usage.GasUsed)
            automation.deductGasForContract(usage)
        }
    }

    automation.gasDeductionCheckCount++
    if automation.gasDeductionCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeGasDeductionCycle()
    }
}

// deductGasForContract deducts gas from the contract and logs the deduction
func (automation *ContractGasDeductionAutomation) deductGasForContract(usage common.ContractGasUsage) {
    // Encrypt contract data before deducting gas
    encryptedContractData := automation.encryptContractData(usage)

    // Trigger gas deduction through Synnergy Consensus
    gasDeductionSuccess := automation.consensusSystem.DeductGas(encryptedContractData)

    if gasDeductionSuccess {
        fmt.Printf("Gas deduction for contract %s successfully triggered.\n", usage.ContractID)
        automation.logGasDeductionEvent(usage)
    } else {
        fmt.Printf("Error deducting gas for contract %s.\n", usage.ContractID)
    }
}

// finalizeGasDeductionCycle finalizes the gas deduction cycle and logs it in the ledger
func (automation *ContractGasDeductionAutomation) finalizeGasDeductionCycle() {
    success := automation.consensusSystem.FinalizeGasDeductionCycle()
    if success {
        fmt.Println("Gas deduction cycle finalized successfully.")
        automation.logGasDeductionCycleFinalization()
    } else {
        fmt.Println("Error finalizing gas deduction cycle.")
    }
}

// logGasDeductionEvent logs the gas deduction event for a specific contract in the ledger
func (automation *ContractGasDeductionAutomation) logGasDeductionEvent(usage common.ContractGasUsage) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-deduction-%s", usage.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Deduction",
        Status:    "Deducted",
        Details:   fmt.Sprintf("Gas deduction successfully processed for contract %s. Gas used: %d", usage.ContractID, usage.GasUsed),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with gas deduction event for contract %s.\n", usage.ContractID)
}

// logGasDeductionCycleFinalization logs the finalization of the gas deduction cycle into the ledger
func (automation *ContractGasDeductionAutomation) logGasDeductionCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gas-deduction-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Gas Deduction Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with gas deduction cycle finalization.")
}

// encryptContractData encrypts the contract gas usage data before processing
func (automation *ContractGasDeductionAutomation) encryptContractData(usage common.ContractGasUsage) common.ContractGasUsage {
    encryptedData, err := encryption.EncryptData(usage)
    if err != nil {
        fmt.Println("Error encrypting contract data:", err)
        return usage
    }
    usage.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return usage
}

