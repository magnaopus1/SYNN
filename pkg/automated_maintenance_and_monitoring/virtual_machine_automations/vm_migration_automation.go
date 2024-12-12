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
    VMMigrationCheckInterval = 5000 * time.Millisecond // Interval for checking VM migration requirements
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks in a block
    MigrationThreshold       = 0.80                    // Load threshold for triggering VM migration
)

// VMMigrationAutomation automates virtual machine migration between nodes in case of high load or failover requirements
type VMMigrationAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store migration actions
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    migrationCheckCount int                          // Counter for migration check cycles
}

// NewVMMigrationAutomation initializes the VM migration automation
func NewVMMigrationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMMigrationAutomation {
    return &VMMigrationAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        migrationCheckCount: 0,
    }
}

// StartVMMigrationCheck starts the continuous loop for monitoring and triggering VM migrations
func (automation *VMMigrationAutomation) StartVMMigrationCheck() {
    ticker := time.NewTicker(VMMigrationCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndMigrateVMs()
        }
    }()
}

// monitorAndMigrateVMs checks the load distribution and triggers VM migrations if necessary
func (automation *VMMigrationAutomation) monitorAndMigrateVMs() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the load statistics of all virtual machines
    vmLoadStats := automation.consensusSystem.GetVMLoadStatistics()

    // Identify VMs that need to be migrated based on the load threshold
    vmsToMigrate := automation.findVMsToMigrate(vmLoadStats)

    if len(vmsToMigrate) > 0 {
        for _, vm := range vmsToMigrate {
            fmt.Printf("VM %s exceeds migration threshold. Initiating migration.\n", vm.ID)
            automation.migrateVM(vm)
        }
    } else {
        fmt.Println("No VMs exceed migration threshold.")
    }

    automation.migrationCheckCount++
    fmt.Printf("VM migration check cycle #%d executed.\n", automation.migrationCheckCount)

    if automation.migrationCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeVMMigrationCycle()
    }
}

// findVMsToMigrate identifies VMs that exceed the migration load threshold
func (automation *VMMigrationAutomation) findVMsToMigrate(loadStats []common.VMLoad) []common.VirtualMachine {
    var vmsToMigrate []common.VirtualMachine

    for _, stat := range loadStats {
        if stat.Load > MigrationThreshold {
            vmsToMigrate = append(vmsToMigrate, stat.VM)
        }
    }

    return vmsToMigrate
}

// migrateVM handles the migration process for a virtual machine to a new node
func (automation *VMMigrationAutomation) migrateVM(vm common.VirtualMachine) {
    // Encrypt VM data for secure migration
    encryptedVMData := automation.encryptVMData(vm)

    // Trigger the migration process through the consensus system
    migrationSuccess := automation.consensusSystem.TriggerVMMigration(encryptedVMData)

    if migrationSuccess {
        fmt.Printf("VM %s successfully migrated to a new node.\n", vm.ID)
        automation.logVMMigrationEvent(vm)
    } else {
        fmt.Printf("Error migrating VM %s.\n", vm.ID)
    }
}

// finalizeVMMigrationCycle finalizes the VM migration check cycle and logs the result in the ledger
func (automation *VMMigrationAutomation) finalizeVMMigrationCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeVMMigrationCycle()
    if success {
        fmt.Println("VM migration check cycle finalized successfully.")
        automation.logVMMigrationCycleFinalization()
    } else {
        fmt.Println("Error finalizing VM migration check cycle.")
    }
}

// logVMMigrationEvent logs the migration action for a specific VM into the ledger
func (automation *VMMigrationAutomation) logVMMigrationEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-migration-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "VM Migration",
        Status:    "Migrated",
        Details:   fmt.Sprintf("VM %s successfully migrated.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with migration event for VM %s.\n", vm.ID)
}

// logVMMigrationCycleFinalization logs the finalization of a VM migration check cycle into the ledger
func (automation *VMMigrationAutomation) logVMMigrationCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("vm-migration-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "VM Migration Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with VM migration cycle finalization.")
}

// encryptVMData encrypts the VM data before initiating the migration process
func (automation *VMMigrationAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachineData {
    vmData := common.VirtualMachineData{
        VM: vm,
    }

    encryptedData, err := encryption.EncryptData(vmData)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return vmData
    }

    fmt.Println("VM migration data successfully encrypted.")
    return encryptedData
}

// ensureVMMigrationIntegrity checks the integrity of migration data and re-triggers migrations if necessary
func (automation *VMMigrationAutomation) ensureVMMigrationIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateVMMigrationIntegrity()
    if !integrityValid {
        fmt.Println("VM migration integrity breach detected. Re-triggering migration checks.")
        automation.monitorAndMigrateVMs()
    } else {
        fmt.Println("VM migration integrity is valid.")
    }
}
