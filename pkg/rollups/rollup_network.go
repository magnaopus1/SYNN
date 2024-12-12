package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// NewRollupNetwork initializes the RollupNetwork for managing rollups and nodes
func NewRollupNetwork(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.RollupNetwork {
	return &common.RollupNetwork{
		Nodes:          make(map[string]*common.Node),
		Rollups:        make(map[string]*common.Rollup),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// AddNode adds a new node to the rollup network
func (rn *common.RollupNetwork) AddNode(nodeID, ipAddress string, nodeType common.NodeType) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if _, exists := rn.Nodes[nodeID]; exists {
		return errors.New("node already exists")
	}

	// Create the new node
	newNode := &common.Node{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		NodeType:  nodeType,
	}

	rn.Nodes[nodeID] = newNode

	// Log the node addition in the ledger
	err := rn.Ledger.RecordNodeAddition(nodeID, ipAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node addition: %v", err)
	}

	fmt.Printf("Node %s added to the rollup network\n", nodeID)
	return nil
}

// RemoveNode removes a node from the rollup network
func (rn *common.RollupNetwork) RemoveNode(nodeID string) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if _, exists := rn.Nodes[nodeID]; !exists {
		return errors.New("node does not exist")
	}

	delete(rn.Nodes, nodeID)

	// Log the node removal in the ledger
	err := rn.Ledger.RecordNodeRemoval(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node removal: %v", err)
	}

	fmt.Printf("Node %s removed from the rollup network\n", nodeID)
	return nil
}

// CreateRollup creates a new rollup within the network
func (rn *common.RollupNetwork) CreateRollup(rollupID string) (*common.Rollup, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if _, exists := rn.Rollups[rollupID]; exists {
		return nil, errors.New("rollup already exists")
	}

	// Create the rollup
	rollup := NewRollup(rollupID, rn.Ledger, rn.Encryption, rn.NetworkManager)
	rn.Rollups[rollupID] = rollup

	// Log the rollup creation in the ledger
	err := rn.Ledger.RecordRollupCreation(rollupID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log rollup creation: %v", err)
	}

	fmt.Printf("Rollup %s created in the rollup network\n", rollupID)
	return rollup, nil
}

// SyncRollup synchronizes a rollup between nodes on the rollup network
func (rn *common.RollupNetwork) SyncRollup(rollupID string, sourceNodeID string, destinationNodeID string) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	rollup, exists := rn.Rollups[rollupID]
	if !exists {
		return fmt.Errorf("rollup %s not found", rollupID)
	}

	sourceNode, exists := rn.Nodes[sourceNodeID]
	if !exists {
		return fmt.Errorf("source node %s not found", sourceNodeID)
	}

	destNode, exists := rn.Nodes[destinationNodeID]
	if !exists {
		return fmt.Errorf("destination node %s not found", destinationNodeID)
	}

	// Encrypt rollup data before transferring
	encryptedRollup, err := rn.Encryption.EncryptData([]byte(rollupID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt rollup data: %v", err)
	}

	// Perform the rollup sync via the network manager
	err = rn.NetworkManager.TransferData(sourceNode, destNode, encryptedRollup)
	if err != nil {
		return fmt.Errorf("failed to sync rollup %s: %v", rollupID, err)
	}

	// Log the rollup sync operation
	err = rn.Ledger.RecordRollupSync(rollupID, sourceNodeID, destinationNodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup sync: %v", err)
	}

	fmt.Printf("Rollup %s synced from node %s to node %s\n", rollupID, sourceNodeID, destinationNodeID)
	return nil
}

// RetrieveRollup retrieves a rollup by its ID from the rollup network
func (rn *common.RollupNetwork) RetrieveRollup(rollupID string) (*common.Rollup, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	rollup, exists := rn.Rollups[rollupID]
	if !exists {
		return nil, fmt.Errorf("rollup %s not found", rollupID)
	}

	fmt.Printf("Retrieved rollup %s\n", rollupID)
	return rollup, nil
}

// BroadcastRollup broadcasts a finalized rollup to all nodes in the network
func (rn *common.RollupNetwork) BroadcastRollup(rollupID string) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	rollup, exists := rn.Rollups[rollupID]
	if !exists {
		return fmt.Errorf("rollup %s not found", rollupID)
	}

	if !rollup.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Encrypt the rollup state root before broadcasting
	encryptedStateRoot, err := rn.Encryption.EncryptData([]byte(rollup.StateRoot), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt rollup state root: %v", err)
	}

	// Broadcast the rollup to all nodes
	for _, node := range rn.Nodes {
		err := rn.NetworkManager.BroadcastData(node, encryptedStateRoot)
		if err != nil {
			return fmt.Errorf("failed to broadcast rollup %s to node %s: %v", rollupID, node.NodeID, err)
		}
	}

	// Log the rollup broadcast event
	err = rn.Ledger.RecordRollupBroadcast(rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup broadcast: %v", err)
	}

	fmt.Printf("Rollup %s broadcast to all nodes\n", rollupID)
	return nil
}

// ValidateRollup ensures a rollup is valid by verifying its state root and transactions
func (rn *common.RollupNetwork) ValidateRollup(rollupID string) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	rollup, exists := rn.Rollups[rollupID]
	if !exists {
		return fmt.Errorf("rollup %s not found", rollupID)
	}

	// Validate the rollup state root and its transactions
	err := common.ValidateMerkleRoot(rollup.Transactions, rollup.StateRoot)
	if err != nil {
		return fmt.Errorf("rollup validation failed: %v", err)
	}

	// Log the rollup validation in the ledger
	err = rn.Ledger.RecordRollupValidation(rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup validation: %v", err)
	}

	fmt.Printf("Rollup %s validated successfully\n", rollupID)
	return nil
}

