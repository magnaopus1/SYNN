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
    DebuggingCheckInterval = 3000 * time.Millisecond // Interval for checking and logging errors and debugging
    SubBlocksPerBlock      = 1000                    // Number of sub-blocks in a block
)

// VMDebuggingAndErrorLoggingAutomation automates the debugging and error logging process for VMs
type VMDebuggingAndErrorLoggingAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance      *ledger.Ledger               // Ledger for recording debugging and error logs
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    debuggingCycleCount int                          // Counter for debugging and error check cycles
}

// NewVMDebuggingAndErrorLoggingAutomation initializes the automation for debugging and error logging across VMs
func NewVMDebuggingAndErrorLoggingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMDebuggingAndErrorLoggingAutomation {
    return &VMDebuggingAndErrorLoggingAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        debuggingCycleCount: 0,
    }
}

// StartDebuggingCheck starts the continuous loop to monitor debugging and error logs in VMs
func (automation *VMDebuggingAndErrorLoggingAutomation) StartDebuggingCheck() {
    ticker := time.NewTicker(DebuggingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndLogDebuggingAndErrors()
        }
    }()
}

// monitorAndLogDebuggingAndErrors checks for debugging issues and errors across VMs and logs them
func (automation *VMDebuggingAndErrorLoggingAutomation) monitorAndLogDebuggingAndErrors() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of active VMs and their statuses
    vmList := automation.consensusSystem.GetActiveVMs()

    for _, vm := range vmList {
        fmt.Printf("Checking for errors and debugging issues for VM %s.\n", vm.ID)
        automation.handleVMErrorsAndDebugging(vm)
    }

    automation.debuggingCycleCount++
    fmt.Printf("Debugging check cycle #%d executed.\n", automation.debuggingCycleCount)

    if automation.debuggingCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeDebuggingCycle()
    }
}

// handleVMErrorsAndDebugging manages errors and debugging issues for each VM and logs them
func (automation *VMDebuggingAndErrorLoggingAutomation) handleVMErrorsAndDebugging(vm common.VirtualMachine) {
    // Encrypt VM data before logging
    encryptedVMData := automation.encryptVMData(vm)

    // Check for errors and debugging flags in the VM
    debuggingSuccess := automation.consensusSystem.CheckAndLogDebugging(encryptedVMData)

    if debuggingSuccess {
        fmt.Printf("No critical errors found in VM %s. Debugging successful.\n", vm.ID)
        automation.logDebuggingEvent(vm)
    } else {
        fmt.Printf("Errors found in VM %s. Triggering error handling.\n", vm.ID)
        automation.triggerErrorHandling(vm)
    }
}

// triggerErrorHandling triggers error handling for a VM with detected issues
func (automation *VMDebuggingAndErrorLoggingAutomation) triggerErrorHandling(vm common.VirtualMachine) {
    // Initiating automatic error handling for the VM
    success := automation.consensusSystem.InitiateErrorHandling(vm.ID)

    if success {
        fmt.Printf("Error handling successfully initiated for VM %s.\n", vm.ID)
        automation.logErrorHandlingEvent(vm)
    } else {
        fmt.Printf("Failed to initiate error handling for VM %s.\n", vm.ID)
    }
}

// finalizeDebuggingCycle finalizes the debugging and error check cycle and logs the result in the ledger
func (automation *VMDebuggingAndErrorLoggingAutomation) finalizeDebuggingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeDebuggingCycle()
    if success {
        fmt.Println("Debugging check cycle finalized successfully.")
        automation.logDebuggingCycleFinalization()
    } else {
        fmt.Println("Error finalizing debugging check cycle.")
    }
}

// logDebuggingEvent logs the debugging event for a specific VM into the ledger
func (automation *VMDebuggingAndErrorLoggingAutomation) logDebuggingEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("debugging-event-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Debugging",
        Status:    "Successful",
        Details:   fmt.Sprintf("Debugging and error check successful for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with debugging event for VM %s.\n", vm.ID)
}

// logErrorHandlingEvent logs the error handling event for a VM into the ledger
func (automation *VMDebuggingAndErrorLoggingAutomation) logErrorHandlingEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("error-handling-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Error Handling",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Error handling successfully triggered for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with error handling event for VM %s.\n", vm.ID)
}

// logDebuggingCycleFinalization logs the finalization of the debugging check cycle into the ledger
func (automation *VMDebuggingAndErrorLoggingAutomation) logDebuggingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("debugging-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Debugging Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with debugging cycle finalization.")
}

// encryptVMData encrypts the VM data before processing debugging and error logging
func (automation *VMDebuggingAndErrorLoggingAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("VM data successfully encrypted.")
    return vm
}

// ensureDebuggingIntegrity ensures that the debugging and error logging processes are running correctly
func (automation *VMDebuggingAndErrorLoggingAutomation) ensureDebuggingIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateDebuggingIntegrity()
    if !integrityValid {
        fmt.Println("Debugging integrity breach detected. Re-triggering debugging checks.")
        automation.monitorAndLogDebuggingAndErrors()
    } else {
        fmt.Println("Debugging integrity is valid.")
    }
}
