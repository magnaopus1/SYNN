package scalability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewGossipSystem initializes the gossip system
func NewGossipSystem(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.GossipSystem {
	return &common.GossipSystem{
		Nodes:             []*common.GossipNode{},
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AddNode adds a new node to the gossip system
func (gs *common.GossipSystem) AddNode(nodeID, nodeType string) *common.GossipNode {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	node := &common.GossipNode{
		NodeID:       nodeID,
		NodeType:     nodeType,
		LastSyncTime: time.Now(),
		Neighbors:    []*common.GossipNode{},
	}

	gs.Nodes = append(gs.Nodes, node)

	// Log the addition of the new node
	err := gs.Ledger.RecordNodeAddition(nodeID, nodeType, time.Now())
	if err != nil {
		fmt.Printf("Failed to log node addition: %v\n", err)
	}

	fmt.Printf("Node %s of type %s added to the gossip system\n", nodeID, nodeType)
	return node
}

// ConnectNodes establishes a gossip connection between two nodes (neighbors)
func (gs *common.GossipSystem) ConnectNodes(node1, node2 *common.GossipNode) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	node1.Neighbors = append(node1.Neighbors, node2)
	node2.Neighbors = append(node2.Neighbors, node1)

	// Log the connection between nodes
	err := gs.Ledger.RecordNodeConnection(node1.NodeID, node2.NodeID, time.Now())
	if err != nil {
		fmt.Printf("Failed to log node connection: %v\n", err)
	}

	fmt.Printf("Nodes %s and %s are now connected\n", node1.NodeID, node2.NodeID)
}

// GossipMessage sends a message to all neighbors of a node
func (gs *common.GossipSystem) GossipMessage(originNode *common.GossipNode, messageData []byte) (*common.GossipMessage, error) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// Encrypt the message
	encryptedMessage, err := gs.EncryptionService.EncryptData(messageData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt gossip message: %v", err)
	}

	// Create a new gossip message
	message := &common.GossipMessage{
		MessageID:  common.GenerateUniqueID(),
		Data:       encryptedMessage,
		Timestamp:  time.Now(),
		OriginNode: originNode.NodeID,
	}

	// Log the message gossip
	err = gs.Ledger.RecordGossipMessage(message.MessageID, originNode.NodeID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log gossip message: %v", err)
	}

	// Send the message to all neighbors
	for _, neighbor := range originNode.Neighbors {
		fmt.Printf("Gossiping message %s to node %s\n", message.MessageID, neighbor.NodeID)
		gs.receiveMessage(neighbor, message)
	}

	return message, nil
}

// receiveMessage processes a message received by a node
func (gs *common.GossipSystem) receiveMessage(node *common.GossipNode, message *common.GossipMessage) {
	// Decrypt the message for processing
	decryptedMessage, err := gs.EncryptionService.DecryptData(message.Data, common.EncryptionKey)
	if err != nil {
		fmt.Printf("Failed to decrypt message at node %s: %v\n", node.NodeID, err)
		return
	}

	fmt.Printf("Node %s received gossip message: %s\n", node.NodeID, string(decryptedMessage))

	// Forward the message to the node's neighbors
	for _, neighbor := range node.Neighbors {
		if neighbor.NodeID != message.OriginNode { // Avoid sending back to origin
			fmt.Printf("Forwarding message %s to neighbor %s\n", message.MessageID, neighbor.NodeID)
			gs.receiveMessage(neighbor, message)
		}
	}
}

// RedundancyProtocol ensures that messages are redundantly distributed across the network
func (gs *common.GossipSystem) RedundancyProtocol(node *common.GossipNode, message *common.GossipMessage) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// Send the message to all neighbors twice for redundancy
	for i := 0; i < 2; i++ {
		for _, neighbor := range node.Neighbors {
			fmt.Printf("Redundantly gossiping message %s to node %s (attempt %d)\n", message.MessageID, neighbor.NodeID, i+1)
			gs.receiveMessage(neighbor, message)
		}
	}

	// Log the redundancy action
	err := gs.Ledger.RecordRedundancyAction(message.MessageID, node.NodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log redundancy action: %v", err)
	}

	fmt.Printf("Redundancy protocol applied for message %s from node %s\n", message.MessageID, node.NodeID)
	return nil
}

// SyncProtocol synchronizes messages between two nodes to ensure consistency
func (gs *common.GossipSystem) SyncProtocol(node1, node2 *common.GossipNode) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// Simulate message synchronization between two nodes
	fmt.Printf("Synchronizing messages between node %s and node %s\n", node1.NodeID, node2.NodeID)

	// Log the synchronization action in the ledger
	err := gs.Ledger.RecordSyncAction(node1.NodeID, node2.NodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sync action: %v", err)
	}

	node1.LastSyncTime = time.Now()
	node2.LastSyncTime = time.Now()

	fmt.Printf("Sync protocol completed between node %s and node %s\n", node1.NodeID, node2.NodeID)
	return nil
}
