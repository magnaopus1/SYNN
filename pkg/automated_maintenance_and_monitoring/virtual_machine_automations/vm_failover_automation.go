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
    VMFailoverCheckInterval = 5000 * time.Millisecond // Interval for checking VM failover readiness
    SubBlocksPerBlock       = 1000                    // Number of sub-blocks in a block
)

// VMFailoverAutomation automates the failover process for VMs in the event of failure
type VMFailoverAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance   *ledger.Ledger               // Ledger to store failover actions
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    failoverCheckCount int                        // Counter for failover check cycles
}

// NewVMFailoverAutomation initializes the automation for VM failover handling
func NewVMFailoverAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMFailoverAutomation {
    return &VMFailoverAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        failoverCheckCount: 0,
    }
}

// StartFailoverCheck starts the continuous loop for monitoring and enforcing VM failover
func (automation *VMFailoverAutomation) StartFailoverCheck() {
    ticker := time.NewTicker(VMFailoverCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndTriggerFailover()
        }
    }()
}

// monitorAndTriggerFailover checks the VM health and triggers failover if necessary
func (automation *VMFailoverAutomation) monitorAndTriggerFailover() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of active VMs and their statuses
    vmList := automation.consensusSystem.GetActiveVMs()

    for _, vm := range vmList {
        fmt.Printf("Checking failover readiness for VM %s (Status: %s).\n", vm.ID, vm.Status)
        if automation.isVMFailureDetected(vm) {
            fmt.Printf("VM %s has failed. Initiating failover.\n", vm.ID)
            automation.triggerFailoverForVM(vm)
        }
    }

    automation.failoverCheckCount++
    fmt.Printf("VM failover check cycle #%d executed.\n", automation.failoverCheckCount)

    if automation.failoverCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeFailoverCycle()
    }
}

// isVMFailureDetected checks if the VM has encountered a failure
func (automation *VMFailoverAutomation) isVMFailureDetected(vm common.VirtualMachine) bool {
    // VM is considered failed if its status indicates failure
    return vm.Status == "failed"
}

// triggerFailoverForVM initiates the failover process for a failed VM
func (automation *VMFailoverAutomation) triggerFailoverForVM(vm common.VirtualMachine) {
    // Encrypt VM failover data before triggering failover
    encryptedVMData := automation.encryptVMFailoverData(vm)

    // Trigger failover for the VM through the Synnergy Consensus
    failoverSuccess := automation.consensusSystem.TriggerVMFailover(encryptedVMData)

    if failoverSuccess {
        fmt.Printf("Failover successfully triggered for VM %s.\n", vm.ID)
        automation.logFailoverEvent(vm)
    } else {
        fmt.Printf("Error triggering failover for VM %s.\n", vm.ID)
    }
}

// finalizeFailoverCycle finalizes the failover check cycle and logs the result in the ledger
func (automation *VMFailoverAutomation) finalizeFailoverCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeFailoverCycle()
    if success {
        fmt.Println("Failover check cycle finalized successfully.")
        automation.logFailoverCycleFinalization()
    } else {
        fmt.Println("Error finalizing failover check cycle.")
    }
}

// logFailoverEvent logs the failover action for a specific VM into the ledger for traceability
func (automation *VMFailoverAutomation) logFailoverEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-failover-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Failover",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Failover successfully triggered for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with VM failover event for VM %s.\n", vm.ID)
}

// logFailoverCycleFinalization logs the finalization of a failover check cycle into the ledger
func (automation *VMFailoverAutomation) logFailoverCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("failover-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Failover Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with failover cycle finalization.")
}

// encryptVMFailoverData encrypts the VM data before triggering failover
func (automation *VMFailoverAutomation) encryptVMFailoverData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting VM failover data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("VM failover data successfully encrypted.")
    return vm
}

// ensureFailoverIntegrity ensures the failover process integrity and re-triggers checks if needed
func (automation *VMFailoverAutomation) ensureFailoverIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateFailoverIntegrity()
    if !integrityValid {
        fmt.Println("Failover integrity breach detected. Re-triggering failover checks.")
        automation.monitorAndTriggerFailover()
    } else {
        fmt.Println("Failover integrity is valid.")
    }
}
