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
    VMLoadBalancingCheckInterval = 3000 * time.Millisecond // Interval for checking VM load balancing
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
    MaxLoadThreshold             = 0.75                    // Max load threshold before triggering balancing
    MinLoadThreshold             = 0.25                    // Min load threshold for reallocation
)

// VMLoadBalancingAutomation automates load balancing across virtual machines (VMs) in the blockchain
type VMLoadBalancingAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store load balancing actions
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    loadCheckCount      int                          // Counter for load balancing check cycles
}

// NewVMLoadBalancingAutomation initializes the load balancing automation across VMs
func NewVMLoadBalancingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMLoadBalancingAutomation {
    return &VMLoadBalancingAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        loadCheckCount:  0,
    }
}

// StartLoadBalancingCheck starts the continuous loop for monitoring and enforcing load balancing across VMs
func (automation *VMLoadBalancingAutomation) StartLoadBalancingCheck() {
    ticker := time.NewTicker(VMLoadBalancingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndBalanceVMs()
        }
    }()
}

// monitorAndBalanceVMs checks the load distribution across virtual machines and enforces load balancing if necessary
func (automation *VMLoadBalancingAutomation) monitorAndBalanceVMs() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the load statistics of all virtual machines
    vmLoadStats := automation.consensusSystem.GetVMLoadStatistics()

    overloadedVMs := automation.findOverloadedVMs(vmLoadStats)
    underloadedVMs := automation.findUnderloadedVMs(vmLoadStats)

    if len(overloadedVMs) > 0 || len(underloadedVMs) > 0 {
        automation.reallocateLoad(overloadedVMs, underloadedVMs)
    } else {
        fmt.Println("All virtual machines are balanced. No reallocation needed.")
    }

    automation.loadCheckCount++
    fmt.Printf("VM load balancing check cycle #%d executed.\n", automation.loadCheckCount)

    if automation.loadCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeVMLoadBalancingCycle()
    }
}

// findOverloadedVMs identifies VMs that are overloaded and require load balancing
func (automation *VMLoadBalancingAutomation) findOverloadedVMs(loadStats []common.VMLoad) []common.VirtualMachine {
    var overloadedVMs []common.VirtualMachine

    for _, stat := range loadStats {
        if stat.Load > MaxLoadThreshold {
            overloadedVMs = append(overloadedVMs, stat.VM)
        }
    }

    return overloadedVMs
}

// findUnderloadedVMs identifies VMs that are underloaded and can receive load
func (automation *VMLoadBalancingAutomation) findUnderloadedVMs(loadStats []common.VMLoad) []common.VirtualMachine {
    var underloadedVMs []common.VirtualMachine

    for _, stat := range loadStats {
        if stat.Load < MinLoadThreshold {
            underloadedVMs = append(underloadedVMs, stat.VM)
        }
    }

    return underloadedVMs
}

// reallocateLoad redistributes the load from overloaded VMs to underloaded VMs
func (automation *VMLoadBalancingAutomation) reallocateLoad(overloadedVMs, underloadedVMs []common.VirtualMachine) {
    for _, overloadedVM := range overloadedVMs {
        for _, underloadedVM := range underloadedVMs {
            fmt.Printf("Reallocating load from overloaded VM %s to underloaded VM %s.\n", overloadedVM.ID, underloadedVM.ID)
            encryptedData := automation.encryptVMData(overloadedVM, underloadedVM)
            automation.balanceLoadBetweenVMs(encryptedData)
        }
    }
}

// balanceLoadBetweenVMs triggers the load balancing action via the consensus system
func (automation *VMLoadBalancingAutomation) balanceLoadBetweenVMs(encryptedData common.VirtualMachineData) {
    success := automation.consensusSystem.TriggerVMLoadBalancing(encryptedData)

    if success {
        fmt.Printf("Load successfully reallocated between VMs.\n")
        automation.logVMLoadBalancingEvent(encryptedData)
    } else {
        fmt.Println("Error in reallocating load between VMs.")
    }
}

// finalizeVMLoadBalancingCycle finalizes the VM load balancing check cycle and logs the result in the ledger
func (automation *VMLoadBalancingAutomation) finalizeVMLoadBalancingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVMLoadBalancingCycle()
    if success {
        fmt.Println("VM load balancing cycle finalized successfully.")
        automation.logVMLoadBalancingCycleFinalization()
    } else {
        fmt.Println("Error finalizing VM load balancing cycle.")
    }
}

// logVMLoadBalancingEvent logs the VM load balancing action into the ledger for traceability
func (automation *VMLoadBalancingAutomation) logVMLoadBalancingEvent(data common.VirtualMachineData) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-load-balancing-%s-to-%s", data.OverloadedVM.ID, data.UnderloadedVM.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Load Balancing",
        Status:    "Completed",
        Details:   fmt.Sprintf("Load reallocated from VM %s to VM %s.", data.OverloadedVM.ID, data.UnderloadedVM.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with load balancing event between VM %s and VM %s.\n", data.OverloadedVM.ID, data.UnderloadedVM.ID)
}

// logVMLoadBalancingCycleFinalization logs the finalization of a VM load balancing cycle into the ledger
func (automation *VMLoadBalancingAutomation) logVMLoadBalancingCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-load-balancing-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Load Balancing Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with VM load balancing cycle finalization.")
}

// encryptVMData encrypts the VM load balancing data before triggering the load reallocation
func (automation *VMLoadBalancingAutomation) encryptVMData(overloadedVM, underloadedVM common.VirtualMachine) common.VirtualMachineData {
    data := common.VirtualMachineData{
        OverloadedVM:  overloadedVM,
        UnderloadedVM: underloadedVM,
    }

    encryptedData, err := encryption.EncryptData(data)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return data
    }

    fmt.Println("VM load balancing data successfully encrypted.")
    return encryptedData
}

// ensureVMLoadBalancingIntegrity checks the integrity of load balancing data and triggers balancing if necessary
func (automation *VMLoadBalancingAutomation) ensureVMLoadBalancingIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVMLoadBalancingIntegrity()
    if !integrityValid {
        fmt.Println("VM load balancing integrity breach detected. Re-triggering balancing checks.")
        automation.monitorAndBalanceVMs()
    } else {
        fmt.Println("VM load balancing integrity is valid.")
    }
}
