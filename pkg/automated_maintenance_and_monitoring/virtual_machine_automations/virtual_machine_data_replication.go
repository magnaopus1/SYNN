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
    DataReplicationInterval = 5000 * time.Millisecond // Interval for data replication between virtual machines
    SubBlocksPerBlock       = 1000                    // Number of sub-blocks in a block
)

// VirtualMachineDataReplicationAutomation automates the data replication across virtual machines
type VirtualMachineDataReplicationAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to Synnergy Consensus system
    ledgerInstance        *ledger.Ledger               // Ledger for logging replication actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    replicationCycleCount int                          // Counter for replication check cycles
}

// NewVirtualMachineDataReplicationAutomation initializes the automation for data replication
func NewVirtualMachineDataReplicationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VirtualMachineDataReplicationAutomation {
    return &VirtualMachineDataReplicationAutomation{
        consensusSystem:       consensusSystem,
        ledgerInstance:        ledgerInstance,
        stateMutex:            stateMutex,
        replicationCycleCount: 0,
    }
}

// StartDataReplication starts the continuous loop for monitoring and executing data replication across virtual machines
func (automation *VirtualMachineDataReplicationAutomation) StartDataReplication() {
    ticker := time.NewTicker(DataReplicationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndReplicateData()
        }
    }()
}

// monitorAndReplicateData checks the data synchronization status across virtual machines and triggers replication where necessary
func (automation *VirtualMachineDataReplicationAutomation) monitorAndReplicateData() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch replication status across virtual machines
    replicationStatus := automation.consensusSystem.GetVMReplicationStatus()

    outOfSyncVMs := automation.findOutOfSyncVMs(replicationStatus)

    if len(outOfSyncVMs) > 0 {
        for _, vm := range outOfSyncVMs {
            fmt.Printf("Virtual machine %s is out of sync. Replicating data.\n", vm.ID)
            automation.replicateDataForVM(vm)
        }
    } else {
        fmt.Println("All virtual machines are in sync. No replication needed.")
    }

    automation.replicationCycleCount++
    if automation.replicationCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeReplicationCycle()
    }
}

// findOutOfSyncVMs identifies virtual machines that are out of sync and require data replication
func (automation *VirtualMachineDataReplicationAutomation) findOutOfSyncVMs(replicationStatus []common.VMReplicationStatus) []common.VirtualMachine {
    var outOfSyncVMs []common.VirtualMachine

    for _, status := range replicationStatus {
        if !status.IsSynced {
            outOfSyncVMs = append(outOfSyncVMs, status.VirtualMachine)
        }
    }

    return outOfSyncVMs
}

// replicateDataForVM triggers the data replication process for a virtual machine that is out of sync
func (automation *VirtualMachineDataReplicationAutomation) replicateDataForVM(vm common.VirtualMachine) {
    // Encrypt replication data before triggering the process
    encryptedVMData := automation.encryptVMData(vm)

    // Trigger replication process through Synnergy Consensus
    replicationSuccess := automation.consensusSystem.TriggerDataReplication(encryptedVMData)

    if replicationSuccess {
        fmt.Printf("Data replication for virtual machine %s successfully completed.\n", vm.ID)
        automation.logReplicationEvent(vm)
    } else {
        fmt.Printf("Error replicating data for virtual machine %s.\n", vm.ID)
    }
}

// finalizeReplicationCycle finalizes the replication cycle and logs the result in the ledger
func (automation *VirtualMachineDataReplicationAutomation) finalizeReplicationCycle() {
    success := automation.consensusSystem.FinalizeReplicationCycle()
    if success {
        fmt.Println("Data replication cycle finalized successfully.")
        automation.logReplicationCycleFinalization()
    } else {
        fmt.Println("Error finalizing data replication cycle.")
    }
}

// logReplicationEvent logs the replication action for a virtual machine into the ledger
func (automation *VirtualMachineDataReplicationAutomation) logReplicationEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-data-replication-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Data Replication",
        Status:    "Replicated",
        Details:   fmt.Sprintf("Data replication successfully completed for virtual machine %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with data replication event for virtual machine %s.\n", vm.ID)
}

// logReplicationCycleFinalization logs the finalization of a data replication cycle into the ledger
func (automation *VirtualMachineDataReplicationAutomation) logReplicationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-replication-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Replication Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with virtual machine data replication cycle finalization.")
}

// encryptVMData encrypts the virtual machine data before replication
func (automation *VirtualMachineDataReplicationAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting virtual machine data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("Virtual machine data successfully encrypted.")
    return vm
}

// ensureReplicationIntegrity checks the integrity of replicated data and triggers re-replication if necessary
func (automation *VirtualMachineDataReplicationAutomation) ensureReplicationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateReplicationIntegrity()
    if !integrityValid {
        fmt.Println("Replication data integrity breach detected. Re-triggering replication checks.")
        automation.monitorAndReplicateData()
    } else {
        fmt.Println("Replication data integrity is valid.")
    }
}
