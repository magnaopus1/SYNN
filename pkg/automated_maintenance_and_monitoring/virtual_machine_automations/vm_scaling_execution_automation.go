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
    ScalingCheckInterval  = 4000 * time.Millisecond // Interval for checking VM scaling requirements
    SubBlocksPerBlock     = 1000                    // Number of sub-blocks in a block
    MaxResourceThreshold  = 80                      // Threshold for resource utilization before scaling
    MinResourceThreshold  = 30                      // Threshold for scaling down resources
)

// VMScalingExecutionAutomation automates scaling of VM resources based on load
type VMScalingExecutionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging scaling actions
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    scalingCheckCount int                          // Counter for scaling check cycles
}

// NewVMScalingExecutionAutomation initializes the automation for VM scaling execution
func NewVMScalingExecutionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMScalingExecutionAutomation {
    return &VMScalingExecutionAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        scalingCheckCount: 0,
    }
}

// StartScalingCheck starts the continuous loop for monitoring and scaling VM resources
func (automation *VMScalingExecutionAutomation) StartScalingCheck() {
    ticker := time.NewTicker(ScalingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndScaleResources()
        }
    }()
}

// monitorAndScaleResources checks the resource utilization and triggers scaling actions if required
func (automation *VMScalingExecutionAutomation) monitorAndScaleResources() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current resource statistics for VMs
    resourceStats := automation.consensusSystem.GetVMResourceStats()

    // Check if resources need scaling up or down
    scalingRequired := false
    for _, stat := range resourceStats {
        if stat.Utilization > MaxResourceThreshold {
            fmt.Printf("Scaling up VM %s due to high resource usage.\n", stat.VM.ID)
            automation.scaleUpVMResources(stat.VM)
            scalingRequired = true
        } else if stat.Utilization < MinResourceThreshold {
            fmt.Printf("Scaling down VM %s due to low resource usage.\n", stat.VM.ID)
            automation.scaleDownVMResources(stat.VM)
            scalingRequired = true
        }
    }

    if !scalingRequired {
        fmt.Println("No VMs require scaling. Resource utilization is optimal.")
    }

    automation.scalingCheckCount++
    fmt.Printf("VM scaling check cycle #%d executed.\n", automation.scalingCheckCount)

    if automation.scalingCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeScalingCycle()
    }
}

// scaleUpVMResources triggers scaling up of resources for a VM
func (automation *VMScalingExecutionAutomation) scaleUpVMResources(vm common.VirtualMachine) {
    // Encrypt VM data before scaling
    encryptedVM := automation.encryptVMData(vm)

    // Trigger resource scaling through the Synnergy Consensus system
    scalingSuccess := automation.consensusSystem.ScaleUpVMResources(encryptedVM)

    if scalingSuccess {
        fmt.Printf("VM %s resources successfully scaled up.\n", vm.ID)
        automation.logScalingEvent(vm, "Scale Up")
    } else {
        fmt.Printf("Error scaling up resources for VM %s.\n", vm.ID)
    }
}

// scaleDownVMResources triggers scaling down of resources for a VM
func (automation *VMScalingExecutionAutomation) scaleDownVMResources(vm common.VirtualMachine) {
    // Encrypt VM data before scaling
    encryptedVM := automation.encryptVMData(vm)

    // Trigger resource scaling down through the Synnergy Consensus system
    scalingSuccess := automation.consensusSystem.ScaleDownVMResources(encryptedVM)

    if scalingSuccess {
        fmt.Printf("VM %s resources successfully scaled down.\n", vm.ID)
        automation.logScalingEvent(vm, "Scale Down")
    } else {
        fmt.Printf("Error scaling down resources for VM %s.\n", vm.ID)
    }
}

// finalizeScalingCycle finalizes the scaling check cycle and logs the result in the ledger
func (automation *VMScalingExecutionAutomation) finalizeScalingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeScalingCycle()
    if success {
        fmt.Println("VM scaling check cycle finalized successfully.")
        automation.logScalingCycleFinalization()
    } else {
        fmt.Println("Error finalizing VM scaling check cycle.")
    }
}

// logScalingEvent logs the scaling action for a specific VM into the ledger
func (automation *VMScalingExecutionAutomation) logScalingEvent(vm common.VirtualMachine, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-scaling-%s-%s", vm.ID, action),
        Timestamp: time.Now().Unix(),
        Type:      "VM Scaling",
        Status:    action,
        Details:   fmt.Sprintf("VM %s successfully performed scaling action: %s.", vm.ID, action),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM scaling event for VM %s.\n", vm.ID)
}

// logScalingCycleFinalization logs the finalization of a scaling check cycle into the ledger
func (automation *VMScalingExecutionAutomation) logScalingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-scaling-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Scaling Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with VM scaling cycle finalization.")
}

// encryptVMData encrypts the VM data before performing scaling actions
func (automation *VMScalingExecutionAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return vm
    }

    vm.EncryptedData = encryptedData
    fmt.Println("VM data successfully encrypted.")
    return vm
}

// ensureScalingIntegrity checks the integrity of scaling data and re-performs scaling if necessary
func (automation *VMScalingExecutionAutomation) ensureScalingIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateScalingIntegrity()
    if !integrityValid {
        fmt.Println("VM scaling integrity breach detected. Re-triggering scaling actions.")
        automation.monitorAndScaleResources()
    } else {
        fmt.Println("VM scaling integrity is valid.")
    }
}
