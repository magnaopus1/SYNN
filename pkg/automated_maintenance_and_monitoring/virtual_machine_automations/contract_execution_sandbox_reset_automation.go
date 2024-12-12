package automations

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

const (
    SandboxResetCheckInterval = 4000 * time.Millisecond // Interval for checking sandbox execution state
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// ContractExecutionSandboxResetAutomation automates the process of resetting the execution sandbox when necessary
type ContractExecutionSandboxResetAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance       *ledger.Ledger               // Ledger for logging resets
    stateMutex           *sync.RWMutex                // Mutex for thread-safe state access
    sandboxResetCheckCount int                        // Counter for sandbox reset checks
}

// NewContractExecutionSandboxResetAutomation initializes the automation for contract execution sandbox reset
func NewContractExecutionSandboxResetAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractExecutionSandboxResetAutomation {
    return &ContractExecutionSandboxResetAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        sandboxResetCheckCount: 0,
    }
}

// StartSandboxResetCheck starts the continuous loop for checking and resetting the sandbox environment
func (automation *ContractExecutionSandboxResetAutomation) StartSandboxResetCheck() {
    ticker := time.NewTicker(SandboxResetCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndResetSandbox()
        }
    }()
}

// monitorAndResetSandbox checks the sandbox state and triggers a reset if needed
func (automation *ContractExecutionSandboxResetAutomation) monitorAndResetSandbox() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current state of all sandbox environments
    sandboxStates := automation.consensusSystem.GetSandboxExecutionStates()

    for _, sandbox := range sandboxStates {
        if sandbox.NeedsReset {
            fmt.Printf("Sandbox for contract %s requires reset.\n", sandbox.ContractID)
            automation.resetSandboxEnvironment(sandbox)
        }
    }

    automation.sandboxResetCheckCount++
    if automation.sandboxResetCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeSandboxResetCycle()
    }
}

// resetSandboxEnvironment triggers a reset of the sandbox environment for a specific contract
func (automation *ContractExecutionSandboxResetAutomation) resetSandboxEnvironment(sandbox common.SandboxState) {
    // Encrypt sandbox data before triggering reset
    encryptedSandboxData := automation.encryptSandboxData(sandbox)

    // Trigger sandbox reset through the Synnergy Consensus
    resetSuccess := automation.consensusSystem.ResetSandbox(encryptedSandboxData)

    if resetSuccess {
        fmt.Printf("Sandbox reset for contract %s successfully triggered.\n", sandbox.ContractID)
        automation.logSandboxResetEvent(sandbox)
    } else {
        fmt.Printf("Error triggering sandbox reset for contract %s.\n", sandbox.ContractID)
    }
}

// finalizeSandboxResetCycle finalizes the sandbox reset check cycle and logs it in the ledger
func (automation *ContractExecutionSandboxResetAutomation) finalizeSandboxResetCycle() {
    success := automation.consensusSystem.FinalizeSandboxResetCycle()
    if success {
        fmt.Println("Sandbox reset check cycle finalized successfully.")
        automation.logSandboxResetCycleFinalization()
    } else {
        fmt.Println("Error finalizing sandbox reset check cycle.")
    }
}

// logSandboxResetEvent logs the sandbox reset event in the ledger for traceability
func (automation *ContractExecutionSandboxResetAutomation) logSandboxResetEvent(sandbox common.SandboxState) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-reset-%s", sandbox.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Reset",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Sandbox reset successfully triggered for contract %s.", sandbox.ContractID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sandbox reset event for contract %s.\n", sandbox.ContractID)
}

// logSandboxResetCycleFinalization logs the finalization of the sandbox reset cycle into the ledger
func (automation *ContractExecutionSandboxResetAutomation) logSandboxResetCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-reset-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Reset Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with sandbox reset cycle finalization.")
}

// encryptSandboxData encrypts the sandbox state data before triggering the reset
func (automation *ContractExecutionSandboxResetAutomation) encryptSandboxData(sandbox common.SandboxState) common.SandboxState {
    encryptedData, err := encryption.EncryptData(sandbox)
    if err != nil {
        fmt.Println("Error encrypting sandbox data:", err)
        return sandbox
    }
    sandbox.EncryptedData = encryptedData
    fmt.Println("Sandbox data successfully encrypted.")
    return sandbox
}

