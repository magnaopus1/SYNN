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
    VMMMonitoringInterval   = 3000 * time.Millisecond // Monitoring interval for the virtual machines
    SubBlocksPerBlock       = 1000                    // Number of sub-blocks in a block
)

// VirtualMachineMonitoringAutomation automates the monitoring and health checking of virtual machines
type VirtualMachineMonitoringAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
    ledgerInstance        *ledger.Ledger               // Ledger for logging monitoring events
    stateMutex            *sync.RWMutex                // Mutex for thread-safe state access
    monitoringCycleCount  int                          // Counter for monitoring check cycles
}

// NewVirtualMachineMonitoringAutomation initializes the automation for monitoring virtual machines
func NewVirtualMachineMonitoringAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VirtualMachineMonitoringAutomation {
    return &VirtualMachineMonitoringAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        monitoringCycleCount: 0,
    }
}

// StartMonitoring starts the continuous loop for monitoring virtual machines' health and performance
func (automation *VirtualMachineMonitoringAutomation) StartMonitoring() {
    ticker := time.NewTicker(VMMMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorVMsHealthAndPerformance()
        }
    }()
}

// monitorVMsHealthAndPerformance checks the health and performance of virtual machines
func (automation *VirtualMachineMonitoringAutomation) monitorVMsHealthAndPerformance() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the status of virtual machines
    vmHealthStatus := automation.consensusSystem.GetVMHealthStatus()

    unhealthyVMs := automation.findUnhealthyVMs(vmHealthStatus)

    if len(unhealthyVMs) > 0 {
        for _, vm := range unhealthyVMs {
            fmt.Printf("Virtual machine %s is unhealthy. Investigating and fixing.\n", vm.ID)
            automation.fixUnhealthyVM(vm)
        }
    } else {
        fmt.Println("All virtual machines are healthy.")
    }

    automation.monitoringCycleCount++
    fmt.Printf("VM Monitoring check cycle #%d executed.\n", automation.monitoringCycleCount)

    if automation.monitoringCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeMonitoringCycle()
    }
}

// findUnhealthyVMs identifies unhealthy virtual machines that need to be fixed
func (automation *VirtualMachineMonitoringAutomation) findUnhealthyVMs(vmHealthStatus []common.VMHealthStatus) []common.VirtualMachine {
    var unhealthyVMs []common.VirtualMachine

    for _, status := range vmHealthStatus {
        if !status.IsHealthy {
            unhealthyVMs = append(unhealthyVMs, status.VirtualMachine)
        }
    }

    return unhealthyVMs
}

// fixUnhealthyVM triggers actions to resolve issues with unhealthy virtual machines
func (automation *VirtualMachineMonitoringAutomation) fixUnhealthyVM(vm common.VirtualMachine) {
    // Encrypt health data before triggering a fix process
    encryptedVMHealthData := automation.encryptVMHealthData(vm)

    // Trigger resolution process through Synnergy Consensus
    resolutionSuccess := automation.consensusSystem.TriggerVMResolution(encryptedVMHealthData)

    if resolutionSuccess {
        fmt.Printf("Resolution for virtual machine %s successfully triggered.\n", vm.ID)
        automation.logVMResolutionEvent(vm)
    } else {
        fmt.Printf("Error triggering resolution for virtual machine %s.\n", vm.ID)
    }
}

// finalizeMonitoringCycle finalizes the monitoring cycle and logs the result in the ledger
func (automation *VirtualMachineMonitoringAutomation) finalizeMonitoringCycle() {
    success := automation.consensusSystem.FinalizeMonitoringCycle()
    if success {
        fmt.Println("Monitoring cycle finalized successfully.")
        automation.logMonitoringCycleFinalization()
    } else {
        fmt.Println("Error finalizing monitoring cycle.")
    }
}

// logVMResolutionEvent logs the resolution event for a virtual machine into the ledger
func (automation *VirtualMachineMonitoringAutomation) logVMResolutionEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-resolution-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Resolution",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Resolution successfully triggered for virtual machine %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM resolution event for virtual machine %s.\n", vm.ID)
}

// logMonitoringCycleFinalization logs the finalization of a monitoring cycle into the ledger
func (automation *VirtualMachineMonitoringAutomation) logMonitoringCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-monitoring-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Monitoring Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with virtual machine monitoring cycle finalization.")
}

// encryptVMHealthData encrypts the health data of the virtual machine before resolution
func (automation *VirtualMachineMonitoringAutomation) encryptVMHealthData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting virtual machine health data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("Virtual machine health data successfully encrypted.")
    return vm
}

// ensureMonitoringIntegrity checks the integrity of the monitoring process and triggers a recheck if necessary
func (automation *VirtualMachineMonitoringAutomation) ensureMonitoringIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateMonitoringIntegrity()
    if !integrityValid {
        fmt.Println("Monitoring integrity breach detected. Re-triggering monitoring checks.")
        automation.monitorVMsHealthAndPerformance()
    } else {
        fmt.Println("Monitoring integrity is valid.")
    }
}
