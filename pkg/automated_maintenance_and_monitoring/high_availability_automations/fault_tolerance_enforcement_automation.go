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
    FaultToleranceCheckInterval = 4000 * time.Millisecond // Interval for checking fault tolerance across nodes
    SubBlocksPerBlock           = 1000                    // Number of sub-blocks in a block
)

// FaultToleranceEnforcementAutomation automates the process of ensuring fault tolerance across the network
type FaultToleranceEnforcementAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store fault tolerance actions
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    faultToleranceCheckCount int                        // Counter for fault tolerance check cycles
}

// NewFaultToleranceEnforcementAutomation initializes the automation for fault tolerance across nodes
func NewFaultToleranceEnforcementAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *FaultToleranceEnforcementAutomation {
    return &FaultToleranceEnforcementAutomation{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        faultToleranceCheckCount: 0,
    }
}

// StartFaultToleranceCheck starts the continuous loop for monitoring and enforcing fault tolerance across nodes
func (automation *FaultToleranceEnforcementAutomation) StartFaultToleranceCheck() {
    ticker := time.NewTicker(FaultToleranceCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndEnforceFaultTolerance()
        }
    }()
}

// monitorAndEnforceFaultTolerance checks the fault tolerance status of nodes and enforces measures to ensure network stability
func (automation *FaultToleranceEnforcementAutomation) monitorAndEnforceFaultTolerance() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch node health status and replication states
    nodeHealth := automation.consensusSystem.GetNodeHealthStatus()

    for _, node := range nodeHealth {
        if automation.isNodeFaultTolerant(node) {
            fmt.Printf("Node %s is operating within fault tolerance limits.\n", node.ID)
        } else {
            fmt.Printf("Node %s is outside fault tolerance limits. Enforcing fault tolerance.\n", node.ID)
            automation.enforceFaultToleranceForNode(node)
        }
    }

    automation.faultToleranceCheckCount++
    fmt.Printf("Fault tolerance check cycle #%d executed.\n", automation.faultToleranceCheckCount)

    if automation.faultToleranceCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeFaultToleranceCycle()
    }
}

// isNodeFaultTolerant checks if a node is operating within fault tolerance limits
func (automation *FaultToleranceEnforcementAutomation) isNodeFaultTolerant(node common.Node) bool {
    return automation.consensusSystem.IsNodeWithinFaultToleranceLimits(node.ID)
}

// enforceFaultToleranceForNode triggers fault tolerance measures for a node that is operating outside safe limits
func (automation *FaultToleranceEnforcementAutomation) enforceFaultToleranceForNode(node common.Node) {
    // Encrypt node data before enforcing fault tolerance
    encryptedNodeData := automation.AddEncryptionToNodeData(node)

    // Trigger fault tolerance enforcement through Synnergy Consensus
    enforcementSuccess := automation.consensusSystem.EnforceFaultToleranceMeasures(encryptedNodeData)

    if enforcementSuccess {
        fmt.Printf("Fault tolerance successfully enforced for node %s.\n", node.ID)
        automation.logFaultToleranceEnforcement(node)
    } else {
        fmt.Printf("Error enforcing fault tolerance for node %s.\n", node.ID)
    }
}

// finalizeFaultToleranceCycle finalizes the fault tolerance check cycle and logs the result in the ledger
func (automation *FaultToleranceEnforcementAutomation) finalizeFaultToleranceCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeFaultToleranceCycle()
    if success {
        fmt.Println("Fault tolerance check cycle finalized successfully.")
        automation.logFaultToleranceCycleFinalization()
    } else {
        fmt.Println("Error finalizing fault tolerance check cycle.")
    }
}

// logFaultToleranceEnforcement logs the enforcement of fault tolerance measures for a specific node into the ledger for traceability
func (automation *FaultToleranceEnforcementAutomation) logFaultToleranceEnforcement(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("fault-tolerance-enforcement-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Fault Tolerance Enforcement",
        Status:    "Enforced",
        Details:   fmt.Sprintf("Fault tolerance successfully enforced for node %s.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with fault tolerance enforcement event for node %s.\n", node.ID)
}

// logFaultToleranceCycleFinalization logs the finalization of a fault tolerance check cycle into the ledger
func (automation *FaultToleranceEnforcementAutomation) logFaultToleranceCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("fault-tolerance-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Fault Tolerance Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with fault tolerance cycle finalization.")
}

// AddEncryptionToNodeData encrypts the node data before enforcing fault tolerance
func (automation *FaultToleranceEnforcementAutomation) AddEncryptionToNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureFaultToleranceIntegrity checks the integrity of fault tolerance measures and triggers actions if necessary
func (automation *FaultToleranceEnforcementAutomation) ensureFaultToleranceIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateFaultToleranceMeasures()
    if !integrityValid {
        fmt.Println("Fault tolerance integrity breach detected. Re-triggering fault tolerance measures.")
        automation.monitorAndEnforceFaultTolerance()
    } else {
        fmt.Println("Fault tolerance integrity is valid.")
    }
}
