package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/common"
)

const (
    StateSyncCheckInterval = 3000 * time.Millisecond // Interval for checking contract state synchronization
    SubBlocksPerBlock      = 1000                    // Number of sub-blocks in a block
)

// ContractStateSynchronizationAutomation automates the process of synchronizing contract states across nodes
type ContractStateSynchronizationAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to Synnergy Consensus for validation and state synchronization
    ledgerInstance   *ledger.Ledger               // Ledger for logging synchronization events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe state access
    syncCheckCount   int                          // Counter for synchronization check cycles
}

// NewContractStateSynchronizationAutomation initializes the automation for contract state synchronization
func NewContractStateSynchronizationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractStateSynchronizationAutomation {
    return &ContractStateSynchronizationAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        syncCheckCount:   0,
    }
}

// StartStateSyncCheck starts the continuous loop for monitoring and synchronizing contract states across nodes
func (automation *ContractStateSynchronizationAutomation) StartStateSyncCheck() {
    ticker := time.NewTicker(StateSyncCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndSyncState()
        }
    }()
}

// monitorAndSyncState checks the state across nodes and ensures synchronization
func (automation *ContractStateSynchronizationAutomation) monitorAndSyncState() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the current contract states from all nodes
    contractStates := automation.consensusSystem.GetContractStates()

    for _, state := range contractStates {
        fmt.Printf("Checking state synchronization for contract %s.\n", state.ContractID)
        automation.syncContractState(state)
    }

    automation.syncCheckCount++
    if automation.syncCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeSyncCycle()
    }
}

// syncContractState ensures the contract state is synchronized across nodes
func (automation *ContractStateSynchronizationAutomation) syncContractState(state common.ContractState) {
    // Encrypt the contract state before synchronization
    encryptedState := automation.encryptContractState(state)

    // Trigger synchronization through Synnergy Consensus
    syncSuccess := automation.consensusSystem.SynchronizeContractState(encryptedState)

    if syncSuccess {
        fmt.Printf("State synchronization for contract %s successfully triggered.\n", state.ContractID)
        automation.logStateSyncEvent(state)
    } else {
        fmt.Printf("Error synchronizing state for contract %s.\n", state.ContractID)
    }
}

// finalizeSyncCycle finalizes the synchronization cycle and logs the result in the ledger
func (automation *ContractStateSynchronizationAutomation) finalizeSyncCycle() {
    success := automation.consensusSystem.FinalizeStateSyncCycle()
    if success {
        fmt.Println("State synchronization cycle finalized successfully.")
        automation.logSyncCycleFinalization()
    } else {
        fmt.Println("Error finalizing state synchronization cycle.")
    }
}

// logStateSyncEvent logs the synchronization event for a specific contract in the ledger
func (automation *ContractStateSynchronizationAutomation) logStateSyncEvent(state common.ContractState) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("state-sync-%s", state.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "State Synchronization",
        Status:    "Synchronized",
        Details:   fmt.Sprintf("State synchronization successfully completed for contract %s.", state.ContractID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with state synchronization event for contract %s.\n", state.ContractID)
}

// logSyncCycleFinalization logs the finalization of a state synchronization cycle into the ledger
func (automation *ContractStateSynchronizationAutomation) logSyncCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("state-sync-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "State Sync Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with state synchronization cycle finalization.")
}

// encryptContractState encrypts the contract state data before synchronization
func (automation *ContractStateSynchronizationAutomation) encryptContractState(state common.ContractState) common.ContractState {
    encryptedData, err := encryption.EncryptData(state)
    if err != nil {
        fmt.Println("Error encrypting contract state data:", err)
        return state
    }
    state.EncryptedData = encryptedData
    fmt.Println("Contract state successfully encrypted.")
    return state
}

// ensureStateIntegrity checks the integrity of synchronized contract states and re-triggers sync if necessary
func (automation *ContractStateSynchronizationAutomation) ensureStateIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateStateIntegrity()
    if !integrityValid {
        fmt.Println("State integrity breach detected. Re-triggering synchronization checks.")
        automation.monitorAndSyncState()
    } else {
        fmt.Println("State integrity is valid across nodes.")
    }
}
