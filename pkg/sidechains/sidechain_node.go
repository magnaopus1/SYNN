package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/consensus"
)


// NewSidechainNodeManager initializes the SidechainNodeManager
func NewSidechainNodeManager(ledgerInstance *ledger.Ledger, consensus *common.SynnergyConsensus, encryptionService *encryption.Encryption) *common.SidechainNodeManager {
	return &common.SidechainNodeManager{
		Nodes:     make(map[string]*common.SidechainNode),
		Consensus: consensus,
		Ledger:    ledgerInstance,
		Encryption: encryptionService,
	}
}

// AddNode adds a new node to the sidechain network
func (nm *common.SidechainNodeManager) AddNode(nodeID, ipAddress string, nodeType common.NodeType) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if _, exists := nm.Nodes[nodeID]; exists {
		return errors.New("node already exists in the sidechain")
	}

	newNode := &common.SidechainNode{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		NodeType:  nodeType,
		Status:    Active,
		Ledger:    nm.Ledger,
		Encryption: nm.Encryption,
	}

	nm.Nodes[nodeID] = newNode

	// Log the addition of the node in the ledger
	err := nm.Ledger.RecordNodeAddition(nodeID, ipAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node addition: %v", err)
	}

	fmt.Printf("Node %s added to sidechain network\n", nodeID)
	return nil
}

// RemoveNode removes a node from the sidechain network
func (nm *common.SidechainNodeManager) RemoveNode(nodeID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.Nodes[nodeID]
	if !exists {
		return errors.New("node does not exist in the sidechain")
	}

	delete(nm.Nodes, nodeID)

	// Log the removal of the node in the ledger
	err := nm.Ledger.RecordNodeRemoval(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node removal: %v", err)
	}

	fmt.Printf("Node %s removed from sidechain network\n", nodeID)
	return nil
}

// UpdateNodeStatus updates the status of a node in the sidechain network
func (nm *common.SidechainNodeManager) UpdateNodeStatus(nodeID string, status common.NodeStatus) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.Nodes[nodeID]
	if !exists {
		return errors.New("node not found in sidechain")
	}

	node.Status = status

	// Log the status update in the ledger
	err := nm.Ledger.RecordNodeStatusUpdate(nodeID, string(status), time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node status update: %v", err)
	}

	fmt.Printf("Node %s status updated to %s\n", nodeID, status)
	return nil
}

// ValidateNode validates the status of a node in the sidechain network using Synnergy Consensus
func (nm *common.SidechainNodeManager) ValidateNode(nodeID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found in sidechain", nodeID)
	}

	// Validate the node through Synnergy Consensus
	err := nm.Consensus.ValidateNode(node.NodeID)
	if err != nil {
		return fmt.Errorf("validation failed for node %s: %v", nodeID, err)
	}

	// Log the validation in the ledger
	err = nm.Ledger.RecordNodeValidation(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node validation: %v", err)
	}

	fmt.Printf("Node %s validated successfully\n", nodeID)
	return nil
}

// BroadcastDataToNode broadcasts encrypted data to a specific node in the sidechain network
func (nm *common.SidechainNodeManager) BroadcastDataToNode(nodeID string, data []byte) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found in sidechain", nodeID)
	}

	// Encrypt the data before sending
	encryptedData, err := nm.Encryption.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Send the encrypted data to the node
	err = nm.Consensus.NetworkManager.SendData(node, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to send data to node %s: %v", nodeID, err)
	}

	fmt.Printf("Encrypted data sent to node %s\n", nodeID)
	return nil
}

// RetrieveNode retrieves a node's information by nodeID
func (nm *common.SidechainNodeManager) RetrieveNode(nodeID string) (*common.SidechainNode, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	node, exists := nm.Nodes[nodeID]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	fmt.Printf("Node %s retrieved from sidechain\n", nodeID)
	return node, nil
}
