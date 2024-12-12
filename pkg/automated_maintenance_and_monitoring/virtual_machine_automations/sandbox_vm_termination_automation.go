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
    VMSandboxTerminationCheckInterval = 3000 * time.Millisecond // Interval for checking VM sandbox termination conditions
    SubBlocksPerBlock                 = 1000                    // Number of sub-blocks in a block
)

// VMSandboxTerminationAutomation automates the process of terminating virtual machine sandboxes
type VMSandboxTerminationAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
    ledgerInstance      *ledger.Ledger               // Ledger to store sandbox termination actions
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    terminationCycleCount int                        // Counter for termination check cycles
}

// NewVMSandboxTerminationAutomation initializes the automation for terminating VM sandboxes
func NewVMSandboxTerminationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMSandboxTerminationAutomation {
    return &VMSandboxTerminationAutomation{
        consensusSystem:     consensusSystem,
        ledgerInstance:      ledgerInstance,
        stateMutex:          stateMutex,
        terminationCycleCount: 0,
    }
}

// StartSandboxTerminationCheck starts the continuous loop for monitoring and terminating virtual machine sandboxes
func (automation *VMSandboxTerminationAutomation) StartSandboxTerminationCheck() {
    ticker := time.NewTicker(VMSandboxTerminationCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndTerminateSandbox()
        }
    }()
}

// monitorAndTerminateSandbox checks for VM sandboxes that need to be terminated and performs the termination
func (automation *VMSandboxTerminationAutomation) monitorAndTerminateSandbox() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the sandbox states of all VMs
    sandboxStates := automation.consensusSystem.GetVMSandboxStates()

    for _, sandbox := range sandboxStates {
        if sandbox.ShouldTerminate {
            fmt.Printf("Terminating sandbox for VM %s.\n", sandbox.VMID)
            automation.terminateVMSandbox(sandbox)
        }
    }

    automation.terminationCycleCount++
    if automation.terminationCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeTerminationCycle()
    }
}

// terminateVMSandbox handles the termination of the specified VM sandbox and logs the action
func (automation *VMSandboxTerminationAutomation) terminateVMSandbox(sandbox common.VMSandbox) {
    // Encrypt sandbox data before termination
    encryptedSandbox := automation.encryptSandboxData(sandbox)

    // Terminate the sandbox through Synnergy Consensus system
    terminationSuccess := automation.consensusSystem.TerminateVMSandbox(encryptedSandbox)

    if terminationSuccess {
        fmt.Printf("Sandbox for VM %s terminated successfully.\n", sandbox.VMID)
        automation.logSandboxTerminationEvent(sandbox)
    } else {
        fmt.Printf("Error terminating sandbox for VM %s.\n", sandbox.VMID)
    }
}

// finalizeTerminationCycle finalizes the termination cycle and logs the result in the ledger
func (automation *VMSandboxTerminationAutomation) finalizeTerminationCycle() {
    success := automation.consensusSystem.FinalizeTerminationCycle()
    if success {
        fmt.Println("Sandbox termination cycle finalized successfully.")
        automation.logTerminationCycleFinalization()
    } else {
        fmt.Println("Error finalizing sandbox termination cycle.")
    }
}

// logSandboxTerminationEvent logs the sandbox termination event into the ledger for traceability
func (automation *VMSandboxTerminationAutomation) logSandboxTerminationEvent(sandbox common.VMSandbox) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-termination-%s", sandbox.VMID),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Termination",
        Status:    "Terminated",
        Details:   fmt.Sprintf("Sandbox for VM %s was successfully terminated.", sandbox.VMID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sandbox termination event for VM %s.\n", sandbox.VMID)
}

// logTerminationCycleFinalization logs the finalization of a sandbox termination cycle into the ledger
func (automation *VMSandboxTerminationAutomation) logTerminationCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-termination-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Termination Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with sandbox termination cycle finalization.")
}

// encryptSandboxData encrypts the sandbox data before performing termination operations
func (automation *VMSandboxTerminationAutomation) encryptSandboxData(sandbox common.VMSandbox) common.VMSandbox {
    encryptedData, err := encryption.EncryptData(sandbox)
    if err != nil {
        fmt.Println("Error encrypting sandbox data:", err)
        return sandbox
    }
    sandbox.EncryptedData = encryptedData
    fmt.Println("Sandbox data successfully encrypted.")
    return sandbox
}

// ensureSandboxTerminationIntegrity checks the integrity of sandbox data and triggers termination revalidation if necessary
func (automation *VMSandboxTerminationAutomation) ensureSandboxTerminationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateSandboxTerminationIntegrity()
    if !integrityValid {
        fmt.Println("Sandbox termination integrity breach detected. Re-triggering termination checks.")
        automation.monitorAndTerminateSandbox()
    } else {
        fmt.Println("Sandbox termination data integrity is valid.")
    }
}
