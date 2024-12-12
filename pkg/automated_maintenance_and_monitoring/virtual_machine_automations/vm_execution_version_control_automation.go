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
    VMVersionCheckInterval = 4000 * time.Millisecond // Interval for checking VM version control
    SubBlocksPerBlock      = 1000                    // Number of sub-blocks in a block
)

// VMExecutionVersionControlAutomation automates the process of managing version control for VM executions
type VMExecutionVersionControlAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance    *ledger.Ledger               // Ledger for recording version control actions
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    versionCheckCount int                          // Counter for version control check cycles
}

// NewVMExecutionVersionControlAutomation initializes the automation for managing VM version control
func NewVMExecutionVersionControlAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMExecutionVersionControlAutomation {
    return &VMExecutionVersionControlAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        versionCheckCount: 0,
    }
}

// StartVersionControlCheck starts the continuous loop for monitoring version control of VM executions
func (automation *VMExecutionVersionControlAutomation) StartVersionControlCheck() {
    ticker := time.NewTicker(VMVersionCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceVersionControl()
        }
    }()
}

// monitorAndEnforceVersionControl checks the version control across VMs and enforces upgrades or downgrades if necessary
func (automation *VMExecutionVersionControlAutomation) monitorAndEnforceVersionControl() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of active VMs and their versions
    vmList := automation.consensusSystem.GetActiveVMs()

    for _, vm := range vmList {
        fmt.Printf("Checking version control for VM %s (Current Version: %s).\n", vm.ID, vm.Version)
        automation.checkAndUpdateVMVersion(vm)
    }

    automation.versionCheckCount++
    fmt.Printf("VM version control check cycle #%d executed.\n", automation.versionCheckCount)

    if automation.versionCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeVersionControlCycle()
    }
}

// checkAndUpdateVMVersion verifies if the VM version is compliant and triggers updates if necessary
func (automation *VMExecutionVersionControlAutomation) checkAndUpdateVMVersion(vm common.VirtualMachine) {
    // Encrypt VM version data before processing
    encryptedVMData := automation.encryptVMVersionData(vm)

    // Verify if the VM version is up-to-date
    isVersionCompliant := automation.consensusSystem.VerifyVMVersionCompliance(encryptedVMData)

    if !isVersionCompliant {
        fmt.Printf("VM %s is not compliant with the latest version. Triggering update.\n", vm.ID)
        automation.triggerVMVersionUpdate(vm)
    } else {
        fmt.Printf("VM %s is running the latest compliant version.\n", vm.ID)
        automation.logVersionCompliance(vm)
    }
}

// triggerVMVersionUpdate initiates the update process for non-compliant VMs
func (automation *VMExecutionVersionControlAutomation) triggerVMVersionUpdate(vm common.VirtualMachine) {
    // Initiate the version upgrade or downgrade based on the latest compliant version
    updateSuccess := automation.consensusSystem.TriggerVMVersionUpdate(vm.ID)

    if updateSuccess {
        fmt.Printf("Version update successfully triggered for VM %s.\n", vm.ID)
        automation.logVersionUpdate(vm)
    } else {
        fmt.Printf("Failed to trigger version update for VM %s.\n", vm.ID)
    }
}

// finalizeVersionControlCycle finalizes the version control check cycle and logs the result in the ledger
func (automation *VMExecutionVersionControlAutomation) finalizeVersionControlCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVersionControlCycle()
    if success {
        fmt.Println("Version control check cycle finalized successfully.")
        automation.logVersionControlCycleFinalization()
    } else {
        fmt.Println("Error finalizing version control check cycle.")
    }
}

// logVersionCompliance logs the VM's compliance with the current version into the ledger
func (automation *VMExecutionVersionControlAutomation) logVersionCompliance(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("version-compliance-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Version Compliance",
        Status:    "Compliant",
        Details:   fmt.Sprintf("VM %s is running the latest compliant version.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with version compliance for VM %s.\n", vm.ID)
}

// logVersionUpdate logs the version update action for a specific VM into the ledger
func (automation *VMExecutionVersionControlAutomation) logVersionUpdate(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("version-update-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Version Update",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Version update successfully triggered for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with version update event for VM %s.\n", vm.ID)
}

// logVersionControlCycleFinalization logs the finalization of a version control check cycle into the ledger
func (automation *VMExecutionVersionControlAutomation) logVersionControlCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("version-control-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Version Control Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with version control cycle finalization.")
}

// encryptVMVersionData encrypts the VM version data before processing for version control
func (automation *VMExecutionVersionControlAutomation) encryptVMVersionData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting VM version data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("VM version data successfully encrypted.")
    return vm
}

// ensureVersionControlIntegrity ensures the version control process is running correctly and re-triggers checks if needed
func (automation *VMExecutionVersionControlAutomation) ensureVersionControlIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVersionControlIntegrity()
    if !integrityValid {
        fmt.Println("Version control integrity breach detected. Re-triggering version control checks.")
        automation.monitorAndEnforceVersionControl()
    } else {
        fmt.Println("Version control integrity is valid.")
    }
}
