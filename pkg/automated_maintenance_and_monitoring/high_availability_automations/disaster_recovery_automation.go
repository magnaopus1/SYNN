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
    DisasterRecoveryCheckInterval = 5000 * time.Millisecond // Interval for checking node health and disaster recovery actions
    SubBlocksPerBlock             = 1000                    // Number of sub-blocks in a block
)

// DisasterRecoveryAutomation automates the process of ensuring high availability and disaster recovery across the network
type DisasterRecoveryAutomation struct {
    consensusSystem       *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance        *ledger.Ledger               // Ledger to store disaster recovery actions
    stateMutex            *sync.RWMutex                // Mutex for thread-safe access
    recoveryCheckCount    int                          // Counter for recovery check cycles
}

// NewDisasterRecoveryAutomation initializes the automation for disaster recovery across the blockchain network
func NewDisasterRecoveryAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DisasterRecoveryAutomation {
    return &DisasterRecoveryAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        recoveryCheckCount: 0,
    }
}

// StartDisasterRecoveryMonitoring starts the continuous loop for monitoring and triggering disaster recovery across nodes
func (automation *DisasterRecoveryAutomation) StartDisasterRecoveryMonitoring() {
    ticker := time.NewTicker(DisasterRecoveryCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndTriggerRecovery()
        }
    }()
}

// monitorAndTriggerRecovery checks for any failures across nodes and triggers disaster recovery as necessary
func (automation *DisasterRecoveryAutomation) monitorAndTriggerRecovery() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the health status of all nodes
    failedNodes := automation.consensusSystem.CheckFailedNodes()

    if len(failedNodes) > 0 {
        for _, node := range failedNodes {
            fmt.Printf("Node %s has failed. Triggering recovery.\n", node.ID)
            automation.triggerRecoveryForNode(node)
        }
    } else {
        fmt.Println("No node failures detected. Network is healthy.")
    }

    automation.recoveryCheckCount++
    fmt.Printf("Disaster recovery check cycle #%d executed.\n", automation.recoveryCheckCount)

    if automation.recoveryCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeRecoveryCycle()
    }
}

// triggerRecoveryForNode triggers the recovery process for a failed node
func (automation *DisasterRecoveryAutomation) triggerRecoveryForNode(node common.Node) {
    // Encrypt node recovery data before triggering
    encryptedNodeData := automation.AddEncryptionToNodeData(node)

    // Trigger recovery for the node through the Synnergy Consensus
    recoverySuccess := automation.consensusSystem.TriggerNodeRecovery(encryptedNodeData)

    if recoverySuccess {
        fmt.Printf("Node %s recovery successfully triggered.\n", node.ID)
        automation.logRecoveryEvent(node)
    } else {
        fmt.Printf("Error triggering recovery for node %s.\n", node.ID)
    }
}

// finalizeRecoveryCycle finalizes the disaster recovery check cycle and logs the result in the ledger
func (automation *DisasterRecoveryAutomation) finalizeRecoveryCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeRecoveryCycle()
    if success {
        fmt.Println("Disaster recovery check cycle finalized successfully.")
        automation.logRecoveryCycleFinalization()
    } else {
        fmt.Println("Error finalizing disaster recovery check cycle.")
    }
}

// logRecoveryEvent logs the recovery action for a specific node into the ledger for traceability
func (automation *DisasterRecoveryAutomation) logRecoveryEvent(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-recovery-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Recovery",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Recovery successfully triggered for node %s.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with node recovery event for node %s.\n", node.ID)
}

// logRecoveryCycleFinalization logs the finalization of a disaster recovery check cycle into the ledger
func (automation *DisasterRecoveryAutomation) logRecoveryCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("recovery-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Recovery Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with recovery cycle finalization.")
}

// AddEncryptionToNodeData encrypts the node data before triggering recovery
func (automation *DisasterRecoveryAutomation) AddEncryptionToNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureRecoveryIntegrity checks the integrity of disaster recovery data and triggers actions if necessary
func (automation *DisasterRecoveryAutomation) ensureRecoveryIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateRecoveryDataIntegrity()
    if !integrityValid {
        fmt.Println("Recovery data integrity breach detected. Re-triggering disaster recovery checks.")
        automation.monitorAndTriggerRecovery()
    } else {
        fmt.Println("Recovery data integrity is valid.")
    }
}
