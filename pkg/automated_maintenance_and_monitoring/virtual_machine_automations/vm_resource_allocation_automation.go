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
    VMResourceAllocationInterval = 4000 * time.Millisecond // Interval for checking resource allocation
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
    MaxResourceThreshold         = 0.85                    // Maximum resource threshold before reallocation
)

// VMResourceAllocationAutomation automates resource allocation for VMs across the network
type VMResourceAllocationAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store resource allocation actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    allocationCheckCount  int                          // Counter for resource allocation check cycles
}

// NewVMResourceAllocationAutomation initializes the VM resource allocation automation
func NewVMResourceAllocationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMResourceAllocationAutomation {
    return &VMResourceAllocationAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        allocationCheckCount:  0,
    }
}

// StartResourceAllocationCheck starts the continuous loop for monitoring and reallocating resources
func (automation *VMResourceAllocationAutomation) StartResourceAllocationCheck() {
    ticker := time.NewTicker(VMResourceAllocationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndReallocateResources()
        }
    }()
}

// monitorAndReallocateResources checks resource usage across VMs and triggers reallocation if necessary
func (automation *VMResourceAllocationAutomation) monitorAndReallocateResources() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the resource statistics of all virtual machines
    vmResourceStats := automation.consensusSystem.GetVMResourceStatistics()

    // Identify VMs that need additional resources based on the threshold
    vmsToReallocate := automation.findVMsForResourceReallocation(vmResourceStats)

    if len(vmsToReallocate) > 0 {
        for _, vm := range vmsToReallocate {
            fmt.Printf("VM %s exceeds resource threshold. Reallocating resources.\n", vm.ID)
            automation.reallocateResources(vm)
        }
    } else {
        fmt.Println("All VMs are within the resource thresholds.")
    }

    automation.allocationCheckCount++
    fmt.Printf("Resource allocation check cycle #%d executed.\n", automation.allocationCheckCount)

    if automation.allocationCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeResourceAllocationCycle()
    }
}

// findVMsForResourceReallocation identifies VMs that exceed the maximum resource threshold
func (automation *VMResourceAllocationAutomation) findVMsForResourceReallocation(resourceStats []common.VMResource) []common.VirtualMachine {
    var vmsToReallocate []common.VirtualMachine

    for _, stat := range resourceStats {
        if stat.ResourceUsage > MaxResourceThreshold {
            vmsToReallocate = append(vmsToReallocate, stat.VM)
        }
    }

    return vmsToReallocate
}

// reallocateResources handles the resource reallocation process for an overloaded VM
func (automation *VMResourceAllocationAutomation) reallocateResources(vm common.VirtualMachine) {
    // Encrypt VM data for secure resource reallocation
    encryptedVMData := automation.encryptVMData(vm)

    // Trigger the resource reallocation through the Synnergy Consensus system
    reallocationSuccess := automation.consensusSystem.TriggerVMResourceReallocation(encryptedVMData)

    if reallocationSuccess {
        fmt.Printf("Resources for VM %s successfully reallocated.\n", vm.ID)
        automation.logResourceReallocationEvent(vm)
    } else {
        fmt.Printf("Error reallocating resources for VM %s.\n", vm.ID)
    }
}

// finalizeResourceAllocationCycle finalizes the resource allocation check cycle and logs the result in the ledger
func (automation *VMResourceAllocationAutomation) finalizeResourceAllocationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVMResourceAllocationCycle()
    if success {
        fmt.Println("Resource allocation check cycle finalized successfully.")
        automation.logResourceAllocationCycleFinalization()
    } else {
        fmt.Println("Error finalizing resource allocation check cycle.")
    }
}

// logResourceReallocationEvent logs the resource reallocation action for a specific VM into the ledger
func (automation *VMResourceAllocationAutomation) logResourceReallocationEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("resource-reallocation-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Resource Reallocation",
        Status:    "Reallocated",
        Details:   fmt.Sprintf("Resources successfully reallocated for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with resource reallocation event for VM %s.\n", vm.ID)
}

// logResourceAllocationCycleFinalization logs the finalization of a resource allocation check cycle into the ledger
func (automation *VMResourceAllocationAutomation) logResourceAllocationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("resource-allocation-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Resource Allocation Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with resource allocation cycle finalization.")
}

// encryptVMData encrypts the VM data before initiating the resource reallocation process
func (automation *VMResourceAllocationAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachineData {
    vmData := common.VirtualMachineData{
        VM: vm,
    }

    encryptedData, err := encryption.EncryptData(vmData)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return vmData
    }

    fmt.Println("VM resource allocation data successfully encrypted.")
    return encryptedData
}

// ensureResourceAllocationIntegrity checks the integrity of resource allocation data and re-triggers allocations if necessary
func (automation *VMResourceAllocationAutomation) ensureResourceAllocationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVMResourceAllocationIntegrity()
    if !integrityValid {
        fmt.Println("Resource allocation integrity breach detected. Re-triggering resource allocation checks.")
        automation.monitorAndReallocateResources()
    } else {
        fmt.Println("Resource allocation integrity is valid.")
    }
}
