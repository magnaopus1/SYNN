package network

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewTopologyManager initializes a new TopologyManager
func NewTopologyManager() *TopologyManager {
	return &TopologyManager{
		Nodes:        make(map[string]*common.Node),
		NetworkGraph: NewNetworkGraph(),
	}
}

// NetworkGraph represents the network topology, containing nodes and edges for connections between nodes.
type NetworkGraph struct {
	Nodes map[string]*GraphNode     // Nodes in the network, keyed by node ID
	Edges map[string][]*GraphEdge   // Connections between nodes, keyed by node ID (adjacency list)
}

// NewNetworkGraph initializes a new NetworkGraph
func NewNetworkGraph() *NetworkGraph {
	return &NetworkGraph{
		Nodes: make(map[string]*GraphNode),
		Edges: make(map[string][]*GraphEdge), // Correctly using []*GraphEdge as the value type
	}
}


// RegisterNode registers a new node in the network and adds it to the topology
func (tm *TopologyManager) RegisterNode(nodeID string, nodeAddress string) error {
    tm.lock.Lock()
    defer tm.lock.Unlock()

    if _, exists := tm.Nodes[nodeID]; exists {
        return fmt.Errorf("Node %s already exists in the network", nodeID)
    }

    // Manually create a new node using the correct fields, leaving NodeCategory and NodeType blank
    newNode := &common.Node{
        Address:    nodeAddress,
        Name:       nodeID,       // Assuming Name is used for nodeID
        NodeKey:    &common.NodeKey{}, // Placeholder for key; should be generated properly
        IsActive:   true,         // Set node as active initially
        // NodeCategory and NodeType are left blank here for flexibility
    }
    tm.Nodes[nodeID] = newNode

    // Add node to the network graph
    tm.NetworkGraph.Nodes[nodeID] = &GraphNode{
        NodeID:      nodeID,
        NodeInfo:    &NodeInfo{Address: nodeAddress}, // Initialize NodeInfo with address
        Connections: []*GraphEdge{},                  // Initialize with no connections
    }

    fmt.Printf("Node %s registered at %s\n", nodeID, nodeAddress)

    // Save the node's information in the ledger
    nodeInfo := ledger.NodeInfo{
        NodeID:   nodeID,
        Address:  nodeAddress,
        NodeType: "general", // Assign a type or change as needed
    }
    err := tm.LedgerInstance.AddNodeToLedger(nodeInfo)
    if err != nil {
        return fmt.Errorf("Failed to add node %s to the ledger: %v", nodeID, err)
    }

    return nil
}




// AddEdge adds a connection (edge) between two nodes in the network graph
func (ng *NetworkGraph) AddEdge(fromNodeID, toNodeID string, weight float64) {
    fromNode, exists := ng.Nodes[fromNodeID]
    if !exists {
        fmt.Printf("Node %s does not exist in the graph\n", fromNodeID)
        return
    }

    // Create the edge and add it to the node's connections
    edge := &GraphEdge{
        FromNodeID: fromNodeID,
        ToNodeID:   toNodeID,
        Weight:     weight,
    }
    fromNode.Connections = append(fromNode.Connections, edge)
    ng.Edges[fromNodeID] = append(ng.Edges[fromNodeID], edge)

    fmt.Printf("Edge added from %s to %s with weight %.2f\n", fromNodeID, toNodeID, weight)
}

// AddNode adds a node to the network graph
func (ng *NetworkGraph) AddNode(nodeID string, nodeInfo *NodeInfo) {
    // Check if the node already exists
    if _, exists := ng.Nodes[nodeID]; !exists {
        // Add the node to the graph with the provided node information
        ng.Nodes[nodeID] = &GraphNode{
            NodeID:      nodeID,
            NodeInfo:    nodeInfo,
            Connections: []*GraphEdge{}, // Initialize with an empty connection list
        }
        fmt.Printf("Node %s added to the network graph\n", nodeID)
    } else {
        fmt.Printf("Node %s already exists in the network graph\n", nodeID)
    }
}


// RecordNodeConnection records a connection event between two nodes in the ledger
func (tm *TopologyManager) ConnectNodes(nodeID1, nodeID2 string) error {
    tm.lock.Lock()
    defer tm.lock.Unlock()

    if _, exists1 := tm.Nodes[nodeID1]; !exists1 {
        return fmt.Errorf("Node %s not found in the network", nodeID1)
    }

    if _, exists2 := tm.Nodes[nodeID2]; !exists2 {
        return fmt.Errorf("Node %s not found in the network", nodeID2)
    }

    // Connect nodes in the network graph
    tm.NetworkGraph.ConnectNodes(nodeID1, nodeID2)

    // Generate a unique connection ID
    connectionID := fmt.Sprintf("%s-%s", nodeID1, nodeID2)

    // Record the connection in the ledger using a ConnectionEvent struct
    connectionEvent := ledger.ConnectionEvent{
        EventID:      fmt.Sprintf("conn-%s", connectionID),
        ConnectionID: connectionID,
        WalletID:     "", // You may want to associate it with a wallet or leave it empty if not applicable
        EventType:    "connection",
        EventTime:    time.Now(),
        Details:      fmt.Sprintf("Nodes %s and %s connected", nodeID1, nodeID2),
    }

    // Record the connection event in the ledger (pass nodeID1 as the string argument)
    err := tm.LedgerInstance.RecordNodeConnection(nodeID1, connectionEvent)
    if err != nil {
        return fmt.Errorf("Failed to record connection between %s and %s in the ledger: %v", nodeID1, nodeID2, err)
    }

    // Simulate secure communication setup
    fmt.Printf("Nodes %s and %s connected and secure connection established\n", nodeID1, nodeID2)

    return nil
}



// RemoveNode removes a node from the network topology and graph
func (tm *TopologyManager) RemoveNode(nodeID string) error {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	if _, exists := tm.Nodes[nodeID]; !exists {
		return fmt.Errorf("Node %s not found in the network", nodeID)
	}

	delete(tm.Nodes, nodeID)
	tm.NetworkGraph.RemoveNode(nodeID)

	fmt.Printf("Node %s removed from the network\n", nodeID)

	// Remove the node's information from the ledger
	err := tm.LedgerInstance.RemoveNodeFromLedger(nodeID)
	if err != nil {
		return fmt.Errorf("Failed to remove node %s from the ledger: %v", nodeID, err)
	}

	return nil
}

// BroadcastMessage securely broadcasts a message to all nodes in the network
func (tm *TopologyManager) BroadcastMessage(senderID string, message EncryptedMessage) error {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	if _, exists := tm.Nodes[senderID]; !exists {
		return fmt.Errorf("Sender node %s not found in the network", senderID)
	}

	commonMessage := common.EncryptedMessage{
		CipherText: message.CipherText,
		Hash:       message.Hash,

	}

	// Iterate over all nodes and send the converted encrypted message
	for nodeID, node := range tm.Nodes {
		if nodeID != senderID {
			err := node.ReceiveEncryptedMessage(commonMessage) // Now using the common.EncryptedMessage type
			if err != nil {
				fmt.Printf("Failed to send message to node %s: %v\n", nodeID, err)
			}
		}
	}

	return nil
}



// RemoveNode removes a node and its associated edges from the network graph
func (ng *NetworkGraph) RemoveNode(nodeID string) {
	if _, exists := ng.Nodes[nodeID]; exists {
		delete(ng.Nodes, nodeID) // Remove the node from the graph

		// Remove all edges connected to this node
		delete(ng.Edges, nodeID)

		// Iterate over the remaining nodes to remove edges to the deleted node
		for _, edges := range ng.Edges {
			for i := len(edges) - 1; i >= 0; i-- {
				if edges[i].ToNodeID == nodeID || edges[i].FromNodeID == nodeID {
					edges = append(edges[:i], edges[i+1:]...)
				}
			}
		}

		fmt.Printf("Node %s removed from the network graph\n", nodeID)
	}
}


// SecureMessageEncrypt encrypts a message using a shared encryption key
func (tm *TopologyManager) SecureMessageEncrypt(message []byte) (EncryptedMessage, error) {
	// Create an encryption instance
	encryption := &common.Encryption{}

	// Encrypt the message using the shared encryption key
	encryptedMessage, err := encryption.EncryptData("AES", message, common.EncryptionKey)
	if err != nil {
		return EncryptedMessage{}, fmt.Errorf("Failed to encrypt message: %v", err)
	}

	// Return the encrypted message with the required metadata
	return EncryptedMessage{
		CipherText: encryptedMessage,
		Hash:       generateMessageHash(message), // Store the hash of the original message
		CreatedAt:  time.Now(),                   // Record the timestamp
	}, nil
}

// SecureMessageDecrypt decrypts an encrypted message using a shared encryption key
func (tm *TopologyManager) SecureMessageDecrypt(encryptedMessage EncryptedMessage) ([]byte, error) {
    // Create an encryption instance
    encryption := &common.Encryption{}

    // Decrypt the message using the shared encryption key
    decryptedMessage, err := encryption.DecryptData(encryptedMessage.CipherText, common.EncryptionKey) // Removed the extra argument
    if err != nil {
        return nil, fmt.Errorf("Failed to decrypt message: %v", err)
    }

    // Validate the hash of the decrypted message
    if generateMessageHash(decryptedMessage) != encryptedMessage.Hash {
        return nil, fmt.Errorf("Message integrity check failed: hash mismatch")
    }

    return decryptedMessage, nil
}




// MonitorNetworkHealth monitors the health of the network by checking node connectivity
func (tm *TopologyManager) MonitorNetworkHealth() {
    for {
        time.Sleep(10 * time.Second) // Periodic monitoring

        for nodeID, node := range tm.Nodes {
            if !node.IsActive { // Use IsActive field instead of IsAlive method
                fmt.Printf("Node %s is not responding, removing from network\n", nodeID)
                tm.RemoveNode(nodeID)
            }
        }
    }
}

// ConnectNodes adds an edge between two nodes in the network graph
func (ng *NetworkGraph) ConnectNodes(nodeID1, nodeID2 string) {
	if _, exists := ng.Nodes[nodeID1]; !exists {
		ng.Nodes[nodeID1] = &GraphNode{NodeID: nodeID1}
	}
	if _, exists := ng.Nodes[nodeID2]; !exists {
		ng.Nodes[nodeID2] = &GraphNode{NodeID: nodeID2}
	}

	// Add edges between node1 and node2
	ng.Edges[nodeID1] = append(ng.Edges[nodeID1], &GraphEdge{FromNodeID: nodeID1, ToNodeID: nodeID2, Weight: 1.0})
	ng.Edges[nodeID2] = append(ng.Edges[nodeID2], &GraphEdge{FromNodeID: nodeID2, ToNodeID: nodeID1, Weight: 1.0})

	fmt.Printf("Nodes %s and %s connected in the graph\n", nodeID1, nodeID2)
}


