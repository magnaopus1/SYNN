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
    SnapshotCreationInterval  = 3000 * time.Millisecond // Interval for checking if VM snapshots should be created
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// VMSnapshotCreationAutomation automates the process of creating VM snapshots and logging them into the ledger
type VMSnapshotCreationAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging snapshot creations
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    snapshotCheckCount int                         // Counter for snapshot creation cycles
}

// NewVMSnapshotCreationAutomation initializes the automation for VM snapshot creation
func NewVMSnapshotCreationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMSnapshotCreationAutomation {
    return &VMSnapshotCreationAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        snapshotCheckCount: 0,
    }
}

// StartSnapshotCreationCheck starts the continuous loop for monitoring and creating VM snapshots
func (automation *VMSnapshotCreationAutomation) StartSnapshotCreationCheck() {
    ticker := time.NewTicker(SnapshotCreationInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndCreateSnapshots()
        }
    }()
}

// monitorAndCreateSnapshots checks if snapshots should be created and creates them if necessary
func (automation *VMSnapshotCreationAutomation) monitorAndCreateSnapshots() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of VMs that need snapshot creation
    vmList := automation.consensusSystem.GetVMsForSnapshot()

    if len(vmList) > 0 {
        for _, vm := range vmList {
            fmt.Printf("Creating snapshot for VM %s.\n", vm.ID)
            automation.createSnapshotForVM(vm)
        }
    } else {
        fmt.Println("No VMs require snapshot creation at this time.")
    }

    automation.snapshotCheckCount++
    fmt.Printf("Snapshot creation check cycle #%d executed.\n", automation.snapshotCheckCount)

    if automation.snapshotCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeSnapshotCreationCycle()
    }
}

// createSnapshotForVM creates a snapshot of the specified VM
func (automation *VMSnapshotCreationAutomation) createSnapshotForVM(vm common.VirtualMachine) {
    // Encrypt snapshot data before creating the snapshot
    encryptedSnapshotData := automation.encryptSnapshotData(vm)

    // Trigger snapshot creation through the Synnergy Consensus system
    snapshotSuccess := automation.consensusSystem.CreateSnapshot(vm, encryptedSnapshotData)

    if snapshotSuccess {
        fmt.Printf("Snapshot successfully created for VM %s.\n", vm.ID)
        automation.logSnapshotEvent(vm)
    } else {
        fmt.Printf("Error creating snapshot for VM %s.\n", vm.ID)
    }
}

// finalizeSnapshotCreationCycle finalizes the snapshot creation check cycle and logs the result in the ledger
func (automation *VMSnapshotCreationAutomation) finalizeSnapshotCreationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeSnapshotCreationCycle()
    if success {
        fmt.Println("Snapshot creation check cycle finalized successfully.")
        automation.logSnapshotCycleFinalization()
    } else {
        fmt.Println("Error finalizing snapshot creation check cycle.")
    }
}

// logSnapshotEvent logs the snapshot creation event for a specific VM into the ledger
func (automation *VMSnapshotCreationAutomation) logSnapshotEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("snapshot-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Snapshot Creation",
        Status:    "Completed",
        Details:   fmt.Sprintf("Snapshot successfully created for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with snapshot creation event for VM %s.\n", vm.ID)
}

// logSnapshotCycleFinalization logs the finalization of a snapshot creation check cycle into the ledger
func (automation *VMSnapshotCreationAutomation) logSnapshotCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("snapshot-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Snapshot Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with snapshot creation cycle finalization.")
}

// encryptSnapshotData encrypts the snapshot data before creating the snapshot
func (automation *VMSnapshotCreationAutomation) encryptSnapshotData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting snapshot data:", err)
        return vm
    }

    vm.EncryptedData = encryptedData
    fmt.Println("Snapshot data successfully encrypted.")
    return vm
}

// ensureSnapshotIntegrity checks the integrity of the snapshot data and triggers re-creation if necessary
func (automation *VMSnapshotCreationAutomation) ensureSnapshotIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateSnapshotIntegrity()
    if !integrityValid {
        fmt.Println("Snapshot data integrity breach detected. Re-triggering snapshot creation.")
        automation.monitorAndCreateSnapshots()
    } else {
        fmt.Println("Snapshot data integrity is valid.")
    }
}
