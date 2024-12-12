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
)

const (
    NodeConnectionInterval      = 5 * time.Second    // Interval to check for node connection updates
    MaxConnectionsPerCycle      = 15                 // Maximum node connections to handle per cycle
    NodeConnectionTimeout       = 30 * time.Second   // Timeout for each node connection process
    MaxConnectionRetries        = 3                  // Maximum retries for failed connections
)

// NodeConnectionAutomation manages the connection of nodes to the network
type NodeConnectionAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus   // Reference to the consensus system
    ledgerInstance     *ledger.Ledger                 // Reference to the ledger system
    stateMutex         *sync.RWMutex                  // Mutex for thread-safe access
    activeConnections  map[string]bool                // Map of active node connections
    connectionAttempts map[string]int                 // Map of connection retry attempts
    connectionCycle    int                            // Counter for connection monitoring cycles
}

// NewNodeConnectionAutomation initializes the node connection automation system
func NewNodeConnectionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NodeConnectionAutomation {
    return &NodeConnectionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        activeConnections:  make(map[string]bool),
        connectionAttempts: make(map[string]int),
        connectionCycle:    0,
    }
}

// StartNodeConnectionMonitoring starts the continuous loop for monitoring node connections
func (automation *NodeConnectionAutomation) StartNodeConnectionMonitoring() {
    ticker := time.NewTicker(NodeConnectionInterval)

    go func() {
        for range ticker.C {
            automation.checkAndUpdateConnections()
        }
    }()
}

// checkAndUpdateConnections verifies and updates the status of node connections
func (automation *NodeConnectionAutomation) checkAndUpdateConnections() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    nodes := automation.consensusSystem.GetConnectedNodes()
    
    for _, node := range nodes {
        if _, connected := automation.activeConnections[node.ID]; !connected {
            fmt.Printf("Node %s is not connected. Initiating connection...\n", node.ID)
            automation.connectNode(node)
        }
    }

    automation.connectionCycle++
    fmt.Printf("Node connection cycle #%d completed.\n", automation.connectionCycle)
}

// connectNode attempts to establish a connection with the given node
func (automation *NodeConnectionAutomation) connectNode(node common.Node) {
    encryptedNodeData := automation.encryptNodeData(node)
    
    connectionSuccess := automation.consensusSystem.EstablishNodeConnection(node, encryptedNodeData)
    if connectionSuccess {
        fmt.Printf("Node %s successfully connected.\n", node.ID)
        automation.logConnectionEvent(node, "Connected")
        automation.activeConnections[node.ID] = true
        automation.resetConnectionRetry(node.ID)
    } else {
        fmt.Printf("Failed to connect node %s. Retrying...\n", node.ID)
        automation.retryNodeConnection(node)
    }
}

// retryNodeConnection retries a failed node connection a limited number of times
func (automation *NodeConnectionAutomation) retryNodeConnection(node common.Node) {
    automation.connectionAttempts[node.ID]++
    
    if automation.connectionAttempts[node.ID] < MaxConnectionRetries {
        automation.connectNode(node)
    } else {
        fmt.Printf("Max retries reached for node %s. Connection failed.\n", node.ID)
        automation.logConnectionFailure(node)
    }
}

// resetConnectionRetry resets the retry count for a node connection
func (automation *NodeConnectionAutomation) resetConnectionRetry(nodeID string) {
    automation.connectionAttempts[nodeID] = 0
}

// logConnectionEvent logs a successful node connection event to the ledger
func (automation *NodeConnectionAutomation) logConnectionEvent(node common.Node, eventType string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("connection-%s-%s", node.ID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Node Connection Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Node %s successfully %s.", node.ID, eventType),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with connection event for node %s.\n", node.ID)
}

// logConnectionFailure logs the failure of a node connection to the ledger
func (automation *NodeConnectionAutomation) logConnectionFailure(node common.Node) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("connection-failure-%s", node.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Node Connection Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to connect node %s after maximum retries.", node.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with connection failure for node %s.\n", node.ID)
}

// encryptNodeData encrypts node data before connection attempts
func (automation *NodeConnectionAutomation) encryptNodeData(node common.Node) common.Node {
    encryptedData, err := encryption.EncryptData(node.Data)
    if err != nil {
        fmt.Println("Error encrypting node data for connection:", err)
        return node
    }

    node.EncryptedData = encryptedData
    fmt.Println("Node data successfully encrypted for connection.")
    return node
}

// emergencyNodeDisconnection triggers emergency disconnection of a node in case of critical network events
func (automation *NodeConnectionAutomation) emergencyNodeDisconnection(node common.Node) {
    fmt.Printf("Emergency disconnection triggered for node %s.\n", node.ID)
    success := automation.consensusSystem.TriggerEmergencyDisconnection(node)

    if success {
        automation.logConnectionEvent(node, "Disconnected")
        fmt.Println("Emergency disconnection executed successfully.")
        delete(automation.activeConnections, node.ID)
    } else {
        fmt.Println("Emergency disconnection failed.")
    }
}
