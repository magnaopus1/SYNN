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
    BackupCheckInterval  = 4000 * time.Millisecond // Interval for checking contract state backups
    SubBlocksPerBlock    = 1000                    // Number of sub-blocks in a block
)

// ContractStateBackupAutomation automates the process of backing up and recovering contract states
type ContractStateBackupAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to Synnergy Consensus
    ledgerInstance   *ledger.Ledger               // Ledger for logging backup and recovery events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe state access
    backupCheckCount int                          // Counter for backup and recovery checks
}

// NewContractStateBackupAutomation initializes the automation for contract state backup and recovery
func NewContractStateBackupAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ContractStateBackupAutomation {
    return &ContractStateBackupAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        backupCheckCount: 0,
    }
}

// StartBackupAndRecoveryCheck starts the continuous loop for checking contract state backups and recovery
func (automation *ContractStateBackupAutomation) StartBackupAndRecoveryCheck() {
    ticker := time.NewTicker(BackupCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndBackupState()
        }
    }()
}

// monitorAndBackupState checks contract states and creates backups when needed
func (automation *ContractStateBackupAutomation) monitorAndBackupState() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the state data for all active contracts
    contractStates := automation.consensusSystem.GetContractStates()

    for _, state := range contractStates {
        fmt.Printf("Backing up state for contract %s.\n", state.ContractID)
        automation.backupContractState(state)
    }

    automation.backupCheckCount++
    if automation.backupCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeBackupCycle()
    }
}

// backupContractState backs up the state of the contract and logs the event in the ledger
func (automation *ContractStateBackupAutomation) backupContractState(state common.ContractState) {
    // Encrypt the contract state before backing up
    encryptedState := automation.encryptContractState(state)

    // Trigger state backup through Synnergy Consensus
    backupSuccess := automation.consensusSystem.BackupContractState(encryptedState)

    if backupSuccess {
        fmt.Printf("State backup for contract %s successfully triggered.\n", state.ContractID)
        automation.logStateBackupEvent(state)
    } else {
        fmt.Printf("Error backing up state for contract %s.\n", state.ContractID)
    }
}

// finalizeBackupCycle finalizes the backup cycle and logs the finalization event
func (automation *ContractStateBackupAutomation) finalizeBackupCycle() {
    success := automation.consensusSystem.FinalizeBackupCycle()
    if success {
        fmt.Println("Backup cycle finalized successfully.")
        automation.logBackupCycleFinalization()
    } else {
        fmt.Println("Error finalizing backup cycle.")
    }
}

// logStateBackupEvent logs the backup event for a specific contract state in the ledger
func (automation *ContractStateBackupAutomation) logStateBackupEvent(state common.ContractState) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("state-backup-%s", state.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "State Backup",
        Status:    "Backed Up",
        Details:   fmt.Sprintf("State backup successfully completed for contract %s.", state.ContractID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with state backup event for contract %s.\n", state.ContractID)
}

// logBackupCycleFinalization logs the finalization of the backup cycle into the ledger
func (automation *ContractStateBackupAutomation) logBackupCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("backup-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Backup Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with backup cycle finalization.")
}

// encryptContractState encrypts the contract state data before backup
func (automation *ContractStateBackupAutomation) encryptContractState(state common.ContractState) common.ContractState {
    encryptedData, err := encryption.EncryptData(state)
    if err != nil {
        fmt.Println("Error encrypting contract state data:", err)
        return state
    }
    state.EncryptedData = encryptedData
    fmt.Println("Contract state successfully encrypted.")
    return state
}

// recoverContractState recovers the state of a contract in case of failure and logs the recovery event
func (automation *ContractStateBackupAutomation) recoverContractState(contractID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Trigger contract state recovery through Synnergy Consensus
    recoveredState, recoverySuccess := automation.consensusSystem.RecoverContractState(contractID)
    if recoverySuccess {
        fmt.Printf("Contract state for %s recovered successfully.\n", contractID)
        automation.logStateRecoveryEvent(recoveredState)
    } else {
        fmt.Printf("Error recovering state for contract %s.\n", contractID)
    }
}

// logStateRecoveryEvent logs the recovery event for a specific contract state in the ledger
func (automation *ContractStateBackupAutomation) logStateRecoveryEvent(state common.ContractState) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("state-recovery-%s", state.ContractID),
        Timestamp: time.Now().Unix(),
        Type:      "State Recovery",
        Status:    "Recovered",
        Details:   fmt.Sprintf("State recovery successfully completed for contract %s.", state.ContractID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with state recovery event for contract %s.\n", state.ContractID)
}
