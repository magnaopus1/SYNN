package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/encryption"
)

const (
    ConcurrencyCheckInterval = 2000 * time.Millisecond // Interval for checking concurrency across VMs
    SubBlocksPerBlock        = 1000                    // Number of sub-blocks in a block
)

// VMConcurrencyManagementAutomation automates the management of concurrent operations across VMs
type VMConcurrencyManagementAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance     *ledger.Ledger               // Ledger for recording concurrency management actions
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    concurrencyCycleCount int                       // Counter for concurrency check cycles
}

// NewVMConcurrencyManagementAutomation initializes the automation for managing concurrency across VMs
func NewVMConcurrencyManagementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMConcurrencyManagementAutomation {
    return &VMConcurrencyManagementAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        concurrencyCycleCount: 0,
    }
}

// StartConcurrencyManagement starts the continuous loop to monitor and manage concurrency in VMs
func (automation *VMConcurrencyManagementAutomation) StartConcurrencyManagement() {
    ticker := time.NewTicker(ConcurrencyCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndManageConcurrency()
        }
    }()
}

// monitorAndManageConcurrency checks and manages concurrent operations across VMs
func (automation *VMConcurrencyManagementAutomation) monitorAndManageConcurrency() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of active VMs
    vmList := automation.consensusSystem.GetActiveVMs()

    for _, vm := range vmList {
        fmt.Printf("Managing concurrency for VM %s.\n", vm.ID)
        automation.manageVMConcurrency(vm)
    }

    automation.concurrencyCycleCount++
    fmt.Printf("Concurrency management cycle #%d executed.\n", automation.concurrencyCycleCount)

    if automation.concurrencyCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeConcurrencyManagementCycle()
    }
}

// manageVMConcurrency ensures that each VM is properly handling concurrent tasks without overload
func (automation *VMConcurrencyManagementAutomation) manageVMConcurrency(vm common.VirtualMachine) {
    // Encrypt VM data before managing concurrency
    encryptedVMData := automation.encryptVMData(vm)

    // Manage concurrency for the VM through the Synnergy Consensus system
    concurrencyBalanced := automation.consensusSystem.BalanceVMConcurrency(encryptedVMData)

    if concurrencyBalanced {
        fmt.Printf("Concurrency for VM %s successfully balanced.\n", vm.ID)
        automation.logConcurrencyManagementEvent(vm)
    } else {
        fmt.Printf("Concurrency management for VM %s failed. Overload detected.\n", vm.ID)
        automation.triggerConcurrencyCorrection(vm)
    }
}

// triggerConcurrencyCorrection triggers actions to correct concurrency overload for a VM
func (automation *VMConcurrencyManagementAutomation) triggerConcurrencyCorrection(vm common.VirtualMachine) {
    // Initiating automatic concurrency correction
    success := automation.consensusSystem.InitiateConcurrencyCorrection(vm.ID)

    if success {
        fmt.Printf("Concurrency correction successfully initiated for VM %s.\n", vm.ID)
        automation.logConcurrencyCorrectionEvent(vm)
    } else {
        fmt.Printf("Failed to initiate concurrency correction for VM %s.\n", vm.ID)
    }
}

// finalizeConcurrencyManagementCycle finalizes the concurrency management cycle and logs the result in the ledger
func (automation *VMConcurrencyManagementAutomation) finalizeConcurrencyManagementCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeConcurrencyManagementCycle()
    if success {
        fmt.Println("Concurrency management cycle finalized successfully.")
        automation.logConcurrencyManagementCycleFinalization()
    } else {
        fmt.Println("Error finalizing concurrency management cycle.")
    }
}

// logConcurrencyManagementEvent logs the concurrency management action for a specific VM into the ledger for traceability
func (automation *VMConcurrencyManagementAutomation) logConcurrencyManagementEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("concurrency-management-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Concurrency Management",
        Status:    "Balanced",
        Details:   fmt.Sprintf("Concurrency successfully managed for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with concurrency management event for VM %s.\n", vm.ID)
}

// logConcurrencyCorrectionEvent logs the concurrency correction action into the ledger
func (automation *VMConcurrencyManagementAutomation) logConcurrencyCorrectionEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("concurrency-correction-vm-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Concurrency Correction",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Concurrency correction successfully triggered for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with concurrency correction event for VM %s.\n", vm.ID)
}

// logConcurrencyManagementCycleFinalization logs the finalization of a concurrency management check cycle into the ledger
func (automation *VMConcurrencyManagementAutomation) logConcurrencyManagementCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("concurrency-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Concurrency Management Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with concurrency management cycle finalization.")
}

// encryptVMData encrypts the VM data before managing concurrency
func (automation *VMConcurrencyManagementAutomation) encryptVMData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm)
    if err != nil {
        fmt.Println("Error encrypting VM data:", err)
        return vm
    }
    vm.EncryptedData = encryptedData
    fmt.Println("VM data successfully encrypted.")
    return vm
}

// ensureConcurrencyIntegrity checks the integrity of the concurrency management process and re-triggers if necessary
func (automation *VMConcurrencyManagementAutomation) ensureConcurrencyIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateConcurrencyIntegrity()
    if !integrityValid {
        fmt.Println("Concurrency management integrity breach detected. Re-triggering checks.")
        automation.monitorAndManageConcurrency()
    } else {
        fmt.Println("Concurrency management integrity is valid.")
    }
}
