package network

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// BootstrapNode represents the node responsible for bootstrapping the network
type BootstrapNode struct {
	ActiveNodes    map[string]NodeInfo    // Map of active nodes in the network
	mutex          sync.Mutex             // Mutex for safe concurrent access
	ListenAddress  string                 // The IP address and port to listen on
	LedgerInstance *ledger.Ledger         // Pointer to the ledger for recording network activities
	NodeType       string                 // Type of node (e.g., Bootstrap, Validator)
}

// HealthCheckMessage is the message sent to check the node's health
const HealthCheckMessage = "ping"

// HealthCheckAcknowledgment is the expected acknowledgment message from the node
const HealthCheckAcknowledgment = "pong"

// NewBootstrapNode initializes a new bootstrap node
func NewBootstrapNode(listenAddress string, ledgerInstance *ledger.Ledger) *BootstrapNode {
    return &BootstrapNode{
        ActiveNodes:   make(map[string]NodeInfo),
        ListenAddress: listenAddress,
        LedgerInstance: ledgerInstance,
    }
}

// Start begins the operation of the bootstrap node, accepting connections and managing active nodes
func (bn *BootstrapNode) Start() error {
    ln, err := net.Listen("tcp", bn.ListenAddress)
    if err != nil {
        return fmt.Errorf("failed to start the bootstrap node: %v", err)
    }
    defer ln.Close()

    fmt.Printf("Bootstrap Node started on %s...\n", bn.ListenAddress)
    
    go bn.monitorNodeHealth() // Continuously monitor the health of active nodes

    // Listen for new node connections
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Failed to accept connection:", err)
            continue
        }

        // Handle the new connection in a separate goroutine
        go bn.handleNewNode(conn)
    }
}

// handleNewNode handles the connection from a new node, registers it, and sends the node list
func (bn *BootstrapNode) handleNewNode(conn net.Conn) {
    defer conn.Close()

    var newNode NodeInfo
    // Read node information from the connection, expecting a NodeInfo struct
    err := ReadFromConnection(conn, &newNode) // Ensure ReadFromConnection handles NodeInfo
    if err != nil {
        fmt.Println("Error reading node info:", err)
        return
    }

    // Register the new node
    bn.registerNode(newNode)

    // Get the list of active nodes and serialize it to JSON
    activeNodeList := bn.getActiveNodeList()
    nodeListJSON, err := json.Marshal(activeNodeList)
    if err != nil {
        fmt.Println("Error serializing node list:", err)
        return
    }

    // Send the list of active nodes back to the new node
    err = WriteToConnection(conn, string(nodeListJSON)) // Send the node list as a JSON string
    if err != nil {
        fmt.Println("Error sending node list:", err)
        return
    }

    fmt.Printf("New node registered: %s at %s\n", newNode.NodeID, newNode.Address)
}

// registerNode adds a new node to the active node list
func (bn *BootstrapNode) registerNode(newNode NodeInfo) {
    bn.mutex.Lock()
    defer bn.mutex.Unlock()

    // Register the new node in the active nodes map
    bn.ActiveNodes[newNode.NodeID] = newNode

    // Convert newNode (of type NodeInfo) to ledger.NodeInfo
    ledgerNode := ledger.NodeInfo{
        NodeID:         newNode.NodeID,
        Address:        newNode.Address,
        IPAddress:      newNode.IPAddress,
        Port:           newNode.Port,
        NodeType:       ledger.NodeType(newNode.NodeType),  
        GeoLocation:    ledger.GeoLocation{Latitude: newNode.GeoLocation.Latitude, Longitude: newNode.GeoLocation.Longitude},
        LastActiveTime: newNode.LastActiveTime,
        IsOnline:       newNode.IsOnline,
    }

    // Log the registration event to the ledger
    bn.LedgerInstance.LogNodeJoin(ledgerNode)
}


// getActiveNodeList retrieves a list of all active nodes in the network
func (bn *BootstrapNode) getActiveNodeList() []NodeInfo {
    bn.mutex.Lock()
    defer bn.mutex.Unlock()

    nodeList := make([]NodeInfo, 0, len(bn.ActiveNodes))
    for _, nodeInfo := range bn.ActiveNodes {
        nodeList = append(nodeList, nodeInfo)
    }

    return nodeList
}

// monitorNodeHealth periodically checks the health of all active nodes and removes any unresponsive nodes
func (bn *BootstrapNode) monitorNodeHealth() {
    for {
        time.Sleep(1 * time.Minute) // Perform health checks every minute
        fmt.Println("Monitoring node health...")

        bn.mutex.Lock()
        for nodeID, nodeInfo := range bn.ActiveNodes {
            if !bn.checkNodeHealth(nodeInfo) {
                fmt.Printf("Node %s is unresponsive, removing from active nodes.\n", nodeID)
                delete(bn.ActiveNodes, nodeID)
                
                // Log the node removal event in the ledger, assuming it expects the node's address
                bn.LedgerInstance.LogNodeRemoval(nodeInfo.Address) 
            }
        }
        bn.mutex.Unlock()
    }
}




// checkNodeHealth sends a ping to the node to verify its health
func (bn *BootstrapNode) checkNodeHealth(nodeInfo NodeInfo) bool {
    // Establish a connection to the node's address
    conn, err := net.DialTimeout("tcp", nodeInfo.Address, 5*time.Second)
    if err != nil {
        return false // Node is unresponsive
    }
    defer conn.Close()

    // Send a health check message
    err = WriteToConnection(conn, HealthCheckMessage)
    if err != nil {
        return false // Failed to send health check
    }

    // Expect an acknowledgment response
    var response string
    err = ReadFromConnection(conn, &response)
    if err != nil || response != HealthCheckAcknowledgment {
        return false // No valid response
    }

    return true // Node is healthy
}


// RemoveNode removes a node from the active list manually
func (bn *BootstrapNode) RemoveNode(nodeID string) {
    bn.mutex.Lock()
    defer bn.mutex.Unlock()

    if nodeInfo, exists := bn.ActiveNodes[nodeID]; exists {
        delete(bn.ActiveNodes, nodeID)
        fmt.Printf("Node %s has been manually removed from the network.\n", nodeID)
        
        // Log the manual removal event to the ledger
        bn.LedgerInstance.LogNodeRemoval(nodeInfo.Address) // Use node's address for logging
    }
}


// PrintNodeList prints the list of currently active nodes
func (bn *BootstrapNode) PrintNodeList() {
    bn.mutex.Lock()
    defer bn.mutex.Unlock()

    fmt.Println("Active nodes in the network:")
    for _, node := range bn.ActiveNodes {
        fmt.Printf("Node ID: %s, Address: %s\n", node.NodeID, node.Address)
    }
}


// WriteToConnection writes a message to the connection
func WriteToConnection(conn net.Conn, message string) error {
    // Encode the message and send it over the connection
    encoder := json.NewEncoder(conn)
    return encoder.Encode(message)
}

// ReadFromConnection reads a message from the connection and deserializes it into the given struct.
func ReadFromConnection(conn net.Conn, v interface{}) error {
    decoder := json.NewDecoder(conn)
    return decoder.Decode(v) // v can be a NodeInfo or any other struct
}


