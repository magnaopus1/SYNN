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
    StatePruningCheckInterval = 3000 * time.Millisecond // Interval for checking the need for state pruning
    SubBlocksPerBlock         = 1000                    // Number of sub-blocks in a block
)

// VMStatePruningAutomation automates the process of pruning outdated or unnecessary VM states and logs the pruning events
type VMStatePruningAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger for logging state pruning actions
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    pruningCheckCount int                         // Counter for state pruning check cycles
}

// NewVMStatePruningAutomation initializes the automation for VM state pruning
func NewVMStatePruningAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *VMStatePruningAutomation {
    return &VMStatePruningAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        pruningCheckCount: 0,
    }
}

// StartStatePruningCheck starts the continuous loop for monitoring and pruning outdated VM states
func (automation *VMStatePruningAutomation) StartStatePruningCheck() {
    ticker := time.NewTicker(StatePruningCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndPruneStates()
        }
    }()
}

// monitorAndPruneStates checks if any VM has outdated or unnecessary states that require pruning
func (automation *VMStatePruningAutomation) monitorAndPruneStates() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of VMs that require state pruning
    vmList := automation.consensusSystem.GetVMsForStatePruning()

    if len(vmList) > 0 {
        for _, vm := range vmList {
            fmt.Printf("Pruning state for VM %s.\n", vm.ID)
            automation.pruneStateForVM(vm)
        }
    } else {
        fmt.Println("No VMs require state pruning at this time.")
    }

    automation.pruningCheckCount++
    fmt.Printf("State pruning check cycle #%d executed.\n", automation.pruningCheckCount)

    if automation.pruningCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeStatePruningCycle()
    }
}

// pruneStateForVM prunes the outdated or unnecessary state for the specified VM
func (automation *VMStatePruningAutomation) pruneStateForVM(vm common.VirtualMachine) {
    // Encrypt VM state before pruning
    encryptedStateData := automation.encryptStateData(vm)

    // Trigger state pruning through the Synnergy Consensus system
    pruningSuccess := automation.consensusSystem.PruneVMState(vm, encryptedStateData)

    if pruningSuccess {
        fmt.Printf("State successfully pruned for VM %s.\n", vm.ID)
        automation.logStatePruningEvent(vm)
    } else {
        fmt.Printf("Error pruning state for VM %s.\n", vm.ID)
    }
}

// finalizeStatePruningCycle finalizes the state pruning check cycle and logs the result in the ledger
func (automation *VMStatePruningAutomation) finalizeStatePruningCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeStatePruningCycle()
    if success {
        fmt.Println("State pruning check cycle finalized successfully.")
        automation.logStatePruningCycleFinalization()
    } else {
        fmt.Println("Error finalizing state pruning check cycle.")
    }
}

// logStatePruningEvent logs the state pruning event for a specific VM into the ledger
func (automation *VMStatePruningAutomation) logStatePruningEvent(vm common.VirtualMachine) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("state-pruning-%s", vm.ID),
        Timestamp: time.Now().Unix(),
        Type:      "State Pruning",
        Status:    "Completed",
        Details:   fmt.Sprintf("State successfully pruned for VM %s.", vm.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with state pruning event for VM %s.\n", vm.ID)
}

// logStatePruningCycleFinalization logs the finalization of a state pruning check cycle into the ledger
func (automation *VMStatePruningAutomation) logStatePruningCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("state-pruning-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "State Pruning Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with state pruning cycle finalization.")
}

// encryptStateData encrypts the VM state data before pruning
func (automation *VMStatePruningAutomation) encryptStateData(vm common.VirtualMachine) common.VirtualMachine {
    encryptedData, err := encryption.EncryptData(vm.StateData)
    if err != nil {
        fmt.Println("Error encrypting state data:", err)
        return vm
    }

    vm.EncryptedStateData = encryptedData
    fmt.Println("VM state data successfully encrypted.")
    return vm
}

// ensureStatePruningIntegrity checks the integrity of the state pruning process
func (automation *VMStatePruningAutomation) ensureStatePruningIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateStatePruningIntegrity()
    if !integrityValid {
        fmt.Println("State pruning data integrity breach detected. Re-triggering pruning.")
        automation.monitorAndPruneStates()
    } else {
        fmt.Println("State pruning data integrity is valid.")
    }
}
