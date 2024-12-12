package automations

import (
    "fmt"
    "time"
    "sync"
    "errors"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/peer_to_peer"
)

const (
    NodeAdditionInterval       = 10 * time.Second  // Interval to check for new node requests
    MaxNodesPerCycle           = 10                // Maximum number of nodes that can be added per cycle
    NodeAdditionTimeout        = 60 * time.Second  // Timeout for node addition process
)

// NodeAdditionAutomation handles automation for node addition in the network
type NodeAdditionAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus
    ledgerInstance   *ledger.Ledger
    peerNetwork      *peer_to_peer.PeerNetwork
    stateMutex       *sync.RWMutex
    additionCycle    int
    activeNodes      int
}

// NewNodeAdditionAutomation initializes the node addition automation
func NewNodeAdditionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, peerNetwork *peer_to_peer.PeerNetwork, stateMutex *sync.RWMutex) *NodeAdditionAutomation {
    return &NodeAdditionAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        peerNetwork:     peerNetwork,
        stateMutex:      stateMutex,
        additionCycle:   0,
        activeNodes:     0,
    }
}

// StartNodeAdditionMonitoring starts the automation to monitor and add nodes continuously
func (automation *NodeAdditionAutomation) StartNodeAdditionMonitoring() {
    ticker := time.NewTicker(NodeAdditionInterval)

    go func() {
        for range ticker.C {
            automation.processNodeAdditionRequests()
        }
    }()
}

// processNodeAdditionRequests handles the logic for node addition requests
func (automation *NodeAdditionAutomation) processNodeAdditionRequests() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newNodes := automation.fetchNewNodeRequests()
    if len(newNodes) == 0 {
        return
    }

    // Ensure we are not adding more than the maximum allowed nodes per cycle
    nodesToAdd := newNodes
    if len(newNodes) > MaxNodesPerCycle {
        nodesToAdd = newNodes[:MaxNodesPerCycle]
    }

    for _, node := range nodesToAdd {
        err := automation.addNodeToNetwork(node)
        if err != nil {
            fmt.Printf("Failed to add node: %v\n", err)
            automation.logNodeAdditionFailure(node, err)
        } else {
            automation.logNodeAdditionSuccess(node)
        }
    }

    automation.additionCycle++
    fmt.Printf("Node addition cycle #%d completed.\n", automation.additionCycle)

    if automation.additionCycle % 50 == 0 {
        automation.finalizeNodeAdditionCycle()
    }
}

// fetchNewNodeRequests fetches new node requests from the peer network
func (automation *NodeAdditionAutomation) fetchNewNodeRequests() []common.Node {
    return automation.peerNetwork.GetPendingNodeRequests()
}

// addNodeToNetwork adds a node to the network with encryption and consensus validation
func (automation *NodeAdditionAutomation) addNodeToNetwork(node common.Node) error {
    fmt.Printf("Adding node to network: %s\n", node.ID)

    // Encrypt the node ID and other sensitive data
    encryptedNodeID, err := encryption.EncryptData([]byte(node.ID))
    if err != nil {
        return fmt.Errorf("failed to encrypt node ID: %v", err)
    }

    // Consensus validation
    isValid := automation.consensusSystem.ValidateNodeAddition(node)
    if !isValid {
        return errors.New("node addition failed validation")
    }

    // Add the node to the peer network
    success := automation.peerNetwork.AddNodeToNetwork(encryptedNodeID)
    if !success {
        return errors.New("failed to add node to the network")
    }

    // Increment the number of active nodes
    automation.activeNodes++
    return nil
}

// finalizeNodeAdditionCycle finalizes the node addition cycle in the consensus system
func (automation *NodeAdditionAutomation) finalizeNodeAdditionCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeNodeAdditionCycle()
    if success {
        fmt.Println("Node addition cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing node addition cycle.")
    }
}

// logNodeAdditionSuccess logs a successful node addition to the ledger
func (automation *NodeAdditionAutomation) logNodeAdditionSuccess(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-addition-success-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Addition",
        Status:    "Success",
        Details:   fmt.Sprintf("Node added successfully: %s", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with successful node addition for node %s.\n", node.ID)
}

// logNodeAdditionFailure logs a failed node addition attempt to the ledger
func (automation *NodeAdditionAutomation) logNodeAdditionFailure(node common.Node, err error) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-addition-failure-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Addition",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to add node %s: %v", node.ID, err),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with failed node addition for node %s.\n", node.ID)
}

// logCycleFinalization logs the finalization of the node addition cycle to the ledger
func (automation *NodeAdditionAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("node-addition-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Node Addition Cycle",
        Status:    "Finalized",
        Details:   "Node addition cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with node addition cycle finalization.")
}
