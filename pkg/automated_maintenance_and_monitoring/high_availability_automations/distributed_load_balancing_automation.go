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
    LoadBalancingCheckInterval = 2000 * time.Millisecond // Interval for checking load distribution across nodes
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// DistributedLoadBalancingAutomation automates the process of monitoring and enforcing distributed load balancing across the network
type DistributedLoadBalancingAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger to store load balancing actions
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    loadBalancingCheckCount int                       // Counter for load balancing check cycles
}

// NewDistributedLoadBalancingAutomation initializes the automation for load balancing across nodes
func NewDistributedLoadBalancingAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DistributedLoadBalancingAutomation {
    return &DistributedLoadBalancingAutomation{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        loadBalancingCheckCount: 0,
    }
}

// StartLoadBalancingCheck starts the continuous loop for monitoring and enforcing load balancing across nodes
func (automation *DistributedLoadBalancingAutomation) StartLoadBalancingCheck() {
    ticker := time.NewTicker(LoadBalancingCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndBalanceLoad()
        }
    }()
}

// monitorAndBalanceLoad checks the load distribution across nodes and enforces load balancing if necessary
func (automation *DistributedLoadBalancingAutomation) monitorAndBalanceLoad() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the load statistics of all nodes
    loadStats := automation.consensusSystem.GetNodeLoadStatistics()

    overloadedNodes := automation.findOverloadedNodes(loadStats)

    if len(overloadedNodes) > 0 {
        for _, node := range overloadedNodes {
            fmt.Printf("Node %s is overloaded. Balancing load.\n", node.ID)
            automation.balanceLoadForNode(node)
        }
    } else {
        fmt.Println("No nodes are overloaded. Load distribution is balanced.")
    }

    automation.loadBalancingCheckCount++
    fmt.Printf("Load balancing check cycle #%d executed.\n", automation.loadBalancingCheckCount)

    if automation.loadBalancingCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeLoadBalancingCycle()
    }
}

// findOverloadedNodes identifies nodes that are overloaded and require load balancing
func (automation *DistributedLoadBalancingAutomation) findOverloadedNodes(loadStats []common.NodeLoad) []common.Node {
    var overloadedNodes []common.Node

    for _, stat := range loadStats {
        if stat.Load > common.MaxAllowedLoad {
            overloadedNodes = append(overloadedNodes, stat.Node)
        }
    }

    return overloadedNodes
}

// balanceLoadForNode triggers load balancing for the overloaded node
func (automation *DistributedLoadBalancingAutomation) balanceLoadForNode(node common.Node) {
    // Encrypt load balancing data before triggering load redistribution
    encryptedNodeData := automation.AddEncryptionToNodeData(node)

    // Trigger load balancing for the node through the Synnergy Consensus
    balancingSuccess := automation.consensusSystem.TriggerLoadBalancing(encryptedNodeData)

    if balancingSuccess {
        fmt.Printf("Load balancing for node %s successfully triggered.\n", node.ID)
        automation.logLoadBalancingEvent(node)
    } else {
        fmt.Printf("Error triggering load balancing for node %s.\n", node.ID)
    }
}

// finalizeLoadBalancingCycle finalizes the load balancing check cycle and logs the result in the ledger
func (automation *DistributedLoadBalancingAutomation) finalizeLoadBalancingCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeLoadBalancingCycle()
    if success {
        fmt.Println("Load balancing check cycle finalized successfully.")
        automation.logLoadBalancingCycleFinalization()
    } else {
        fmt.Println("Error finalizing load balancing check cycle.")
    }
}

// logLoadBalancingEvent logs the load balancing action for a specific node into the ledger for traceability
func (automation *DistributedLoadBalancingAutomation) logLoadBalancingEvent(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("load-balancing-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Load Balancing",
        Status:    "Triggered",
        Details:   fmt.Sprintf("Load balancing successfully triggered for node %s.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with load balancing event for node %s.\n", node.ID)
}

// logLoadBalancingCycleFinalization logs the finalization of a load balancing check cycle into the ledger
func (automation *DistributedLoadBalancingAutomation) logLoadBalancingCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("load-balancing-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Load Balancing Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with load balancing cycle finalization.")
}

// AddEncryptionToNodeData encrypts the node data before triggering load balancing
func (automation *DistributedLoadBalancingAutomation) AddEncryptionToNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node)
    if err != nil {
        fmt.Println("Error encrypting node data:", err)
        return node
    }
    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted.")
    return node
}

// ensureLoadBalancingIntegrity checks the integrity of load balancing data and triggers balancing if necessary
func (automation *DistributedLoadBalancingAutomation) ensureLoadBalancingIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateLoadBalancingIntegrity()
    if !integrityValid {
        fmt.Println("Load balancing data integrity breach detected. Re-triggering load balancing checks.")
        automation.monitorAndBalanceLoad()
    } else {
        fmt.Println("Load balancing data integrity is valid.")
    }
}
