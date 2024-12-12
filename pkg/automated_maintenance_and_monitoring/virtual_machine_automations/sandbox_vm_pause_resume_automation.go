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
    VMSandboxCheckInterval = 3000 * time.Millisecond // Interval for checking VM sandbox state
    SubBlocksPerBlock      = 1000                    // Number of sub-blocks in a block
)

// VMSandboxPauseResumeAutomation automates the process of pausing and resuming virtual machine sandboxes
type VMSandboxPauseResumeAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
    ledgerInstance   *ledger.Ledger               // Ledger to store sandbox pause and resume actions
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    sandboxCycleCount int                         // Counter for sandbox check cycles
}

// NewVMSandboxPauseResumeAutomation initializes the automation for sandbox VM pause and resume functionality
func NewVMSandboxPauseResumeAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMSandboxPauseResumeAutomation {
    return &VMSandboxPauseResumeAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        sandboxCycleCount: 0,
    }
}

// StartSandboxCheck starts the continuous loop for monitoring and pausing/resuming virtual machine sandboxes
func (automation *VMSandboxPauseResumeAutomation) StartSandboxCheck() {
    ticker := time.NewTicker(VMSandboxCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndControlSandbox()
        }
    }()
}

// monitorAndControlSandbox checks the state of VM sandboxes and triggers pause or resume actions if necessary
func (automation *VMSandboxPauseResumeAutomation) monitorAndControlSandbox() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the sandbox states of all VMs
    sandboxStates := automation.consensusSystem.GetVMSandboxStates()

    for _, sandbox := range sandboxStates {
        if sandbox.ShouldPause {
            fmt.Printf("Pausing sandbox for VM %s.\n", sandbox.VMID)
            automation.pauseVMSandbox(sandbox)
        } else if sandbox.ShouldResume {
            fmt.Printf("Resuming sandbox for VM %s.\n", sandbox.VMID)
            automation.resumeVMSandbox(sandbox)
        }
    }

    automation.sandboxCycleCount++
    if automation.sandboxCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeSandboxCycle()
    }
}

// pauseVMSandbox pauses the specified VM sandbox and logs the action
func (automation *VMSandboxPauseResumeAutomation) pauseVMSandbox(sandbox common.VMSandbox) {
    // Encrypt sandbox data before pausing
    encryptedSandbox := automation.encryptSandboxData(sandbox)

    // Pause the sandbox through Synnergy Consensus system
    pauseSuccess := automation.consensusSystem.PauseVMSandbox(encryptedSandbox)

    if pauseSuccess {
        fmt.Printf("Sandbox for VM %s paused successfully.\n", sandbox.VMID)
        automation.logSandboxPauseEvent(sandbox)
    } else {
        fmt.Printf("Error pausing sandbox for VM %s.\n", sandbox.VMID)
    }
}

// resumeVMSandbox resumes the specified VM sandbox and logs the action
func (automation *VMSandboxPauseResumeAutomation) resumeVMSandbox(sandbox common.VMSandbox) {
    // Encrypt sandbox data before resuming
    encryptedSandbox := automation.encryptSandboxData(sandbox)

    // Resume the sandbox through Synnergy Consensus system
    resumeSuccess := automation.consensusSystem.ResumeVMSandbox(encryptedSandbox)

    if resumeSuccess {
        fmt.Printf("Sandbox for VM %s resumed successfully.\n", sandbox.VMID)
        automation.logSandboxResumeEvent(sandbox)
    } else {
        fmt.Printf("Error resuming sandbox for VM %s.\n", sandbox.VMID)
    }
}

// finalizeSandboxCycle finalizes the sandbox cycle and logs the result in the ledger
func (automation *VMSandboxPauseResumeAutomation) finalizeSandboxCycle() {
    success := automation.consensusSystem.FinalizeSandboxCycle()
    if success {
        fmt.Println("Sandbox cycle finalized successfully.")
        automation.logSandboxCycleFinalization()
    } else {
        fmt.Println("Error finalizing sandbox cycle.")
    }
}

// logSandboxPauseEvent logs the sandbox pause event into the ledger for traceability
func (automation *VMSandboxPauseResumeAutomation) logSandboxPauseEvent(sandbox common.VMSandbox) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-pause-%s", sandbox.VMID),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Pause",
        Status:    "Paused",
        Details:   fmt.Sprintf("Sandbox for VM %s was successfully paused.", sandbox.VMID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sandbox pause event for VM %s.\n", sandbox.VMID)
}

// logSandboxResumeEvent logs the sandbox resume event into the ledger for traceability
func (automation *VMSandboxPauseResumeAutomation) logSandboxResumeEvent(sandbox common.VMSandbox) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-resume-%s", sandbox.VMID),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Resume",
        Status:    "Resumed",
        Details:   fmt.Sprintf("Sandbox for VM %s was successfully resumed.", sandbox.VMID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with sandbox resume event for VM %s.\n", sandbox.VMID)
}

// logSandboxCycleFinalization logs the finalization of a sandbox cycle into the ledger
func (automation *VMSandboxPauseResumeAutomation) logSandboxCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("sandbox-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Sandbox Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with sandbox cycle finalization.")
}

// encryptSandboxData encrypts the sandbox data before performing pause or resume operations
func (automation *VMSandboxPauseResumeAutomation) encryptSandboxData(sandbox common.VMSandbox) common.VMSandbox {
    encryptedData, err := encryption.EncryptData(sandbox)
    if err != nil {
        fmt.Println("Error encrypting sandbox data:", err)
        return sandbox
    }
    sandbox.EncryptedData = encryptedData
    fmt.Println("Sandbox data successfully encrypted.")
    return sandbox
}

// ensureSandboxIntegrity checks the integrity of sandbox data and triggers revalidation if necessary
func (automation *VMSandboxPauseResumeAutomation) ensureSandboxIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateSandboxIntegrity()
    if !integrityValid {
        fmt.Println("Sandbox integrity breach detected. Re-triggering sandbox checks.")
        automation.monitorAndControlSandbox()
    } else {
        fmt.Println("Sandbox data integrity is valid.")
    }
}
