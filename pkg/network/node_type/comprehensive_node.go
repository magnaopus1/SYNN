package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/network"
	"synnergy_network_demo/common"
	"synnergy_network_demo/synnergy_vm"
)

// SuperUltraNode represents a node that can run all standard node types' functionalities in one super ultra configuration.
type SynnergyNode struct {
	NodeID            string
	EnabledNodes      []string // List of enabled node types
	SubNodes          map[string]NodeInterface
	ConsensusEngine   *synnergy_consensus.Engine
	EncryptionService *encryption.Encryption
	NetworkManager    *network.NetworkManager
	mutex             sync.Mutex
	SNVM              *synnergy_vm.VirtualMachine
	SyncInterval      time.Duration
}


// NewSuperUltraNode initializes a Super Ultra Node that runs all the other node types.
func NewSuperUltraNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, snvm *synnergy_vm.VirtualMachine, syncInterval time.Duration) *SuperUltraNode {
	return &SuperUltraNode{
		NodeID:            nodeID,
		SubNodes:          make(map[string]NodeInterface),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              snvm,
		SyncInterval:      syncInterval,
	}
}

// EnableNode enables a specific node type within the SuperUltraNode.
func (sun *SuperUltraNode) EnableNode(nodeType string, nodeInterface NodeInterface) {
	sun.mutex.Lock()
	defer sun.mutex.Unlock()

	if _, exists := sun.SubNodes[nodeType]; exists {
		fmt.Printf("Node type %s is already enabled.\n", nodeType)
		return
	}
	sun.SubNodes[nodeType] = nodeInterface
	fmt.Printf("Node type %s enabled in Super Ultra Node.\n", nodeType)
}

// StartNode starts all enabled sub-nodes in the SuperUltraNode.
func (sun *SuperUltraNode) StartNode() error {
	sun.mutex.Lock()
	defer sun.mutex.Unlock()

	// Start each sub-node based on the enabled node types
	for nodeType, subNode := range sun.SubNodes {
		fmt.Printf("Starting node type: %s\n", nodeType)
		err := subNode.StartNode()
		if err != nil {
			return fmt.Errorf("failed to start node type %s: %v", nodeType, err)
		}
	}
	fmt.Printf("Super Ultra Node %s started successfully with enabled nodes: %v\n", sun.NodeID, sun.EnabledNodes)
	return nil
}

// StopNode stops all running sub-nodes in the SuperUltraNode.
func (sun *SuperUltraNode) StopNode() error {
	sun.mutex.Lock()
	defer sun.mutex.Unlock()

	// Stop each sub-node
	for nodeType, subNode := range sun.SubNodes {
		fmt.Printf("Stopping node type: %s\n", nodeType)
		err := subNode.StopNode()
		if err != nil {
			return fmt.Errorf("failed to stop node type %s: %v", nodeType, err)
		}
	}
	fmt.Printf("Super Ultra Node %s stopped successfully.\n", sun.NodeID)
	return nil
}

// EncryptData encrypts sensitive data before sending it to other nodes, using a common encryption service.
func (sun *SuperUltraNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := sun.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted data from other nodes, using a common encryption service.
func (sun *SuperUltraNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := sun.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// AddSubNode allows dynamically adding new sub-node functionality into the SuperUltraNode.
func (sun *SuperUltraNode) AddSubNode(nodeType string, nodeInterface NodeInterface) error {
	sun.mutex.Lock()
	defer sun.mutex.Unlock()

	if _, exists := sun.SubNodes[nodeType]; exists {
		return fmt.Errorf("node type %s is already added", nodeType)
	}

	sun.SubNodes[nodeType] = nodeInterface
	sun.EnabledNodes = append(sun.EnabledNodes, nodeType)
	fmt.Printf("Sub-node of type %s added to Super Ultra Node %s.\n", nodeType, sun.NodeID)
	return nil
}

// SyncWithNetwork synchronizes the Super Ultra Node with the network.
func (sun *SuperUltraNode) SyncWithNetwork() {
	ticker := time.NewTicker(sun.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		sun.mutex.Lock()
		otherNodes := sun.NetworkManager.DiscoverOtherNodes(sun.NodeID)
		for _, node := range otherNodes {
			// Sync with discovered nodes.
			fmt.Printf("Syncing Super Ultra Node with %s...\n", node)
			sun.syncWithNode(node)
		}
		sun.mutex.Unlock()
	}
}

// syncWithNode synchronizes the state with a peer node.
func (sun *SuperUltraNode) syncWithNode(peerNode string) {
	// Sync logic across multiple node types
	fmt.Printf("Syncing with node %s across all sub-nodes...\n", peerNode)

	for nodeType, subNode := range sun.SubNodes {
		// Each node type should have its own syncing mechanism.
		fmt.Printf("Syncing node type %s with peer node %s...\n", nodeType, peerNode)
		// Assume sub-node has a specific sync function
	}
}

