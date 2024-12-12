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

// NewSidechainNetwork initializes the SidechainNetwork
func NewSidechainNetwork(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.Manager, consensus *common.SynnergyConsensus) *common.SidechainNetwork {
	return &common.SidechainNetwork{
		Nodes:          make(map[string]*common.Node),
		Sidechains:     make(map[string]*common.Sidechain),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AddSidechainNode adds a new node to a sidechain in the network
func (sn *common.SidechainNetwork) AddSidechainNode(sidechainID, nodeID, ipAddress string, nodeType common.NodeType) error {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return errors.New("sidechain not found")
	}

	if _, exists := sn.Nodes[nodeID]; exists {
		return errors.New("node already exists in the network")
	}

	// Create the new node
	newNode := &network.Node{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		NodeType:  nodeType,
	}

	sn.Nodes[nodeID] = newNode
	sidechain.Consensus.AddNode(newNode) // Add the node to consensus validation

	// Log the new node creation in the ledger
	err := sn.Ledger.RecordNodeAddition(nodeID, ipAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node addition: %v", err)
	}

	fmt.Printf("Node %s added to sidechain %s\n", nodeID, sidechainID)
	return nil
}

// RemoveSidechainNode removes a node from the sidechain
func (sn *common.SidechainNetwork) RemoveSidechainNode(sidechainID, nodeID string) error {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return errors.New("sidechain not found")
	}

	if _, exists := sn.Nodes[nodeID]; !exists {
		return errors.New("node does not exist in the network")
	}

	delete(sn.Nodes, nodeID)
	sidechain.Consensus.RemoveNode(nodeID) // Remove the node from consensus validation

	// Log the node removal in the ledger
	err := sn.Ledger.RecordNodeRemoval(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node removal: %v", err)
	}

	fmt.Printf("Node %s removed from sidechain %s\n", nodeID, sidechainID)
	return nil
}

// SyncSidechainBlock synchronizes a block between nodes on a sidechain
func (sn *common.SidechainNetwork) SyncSidechainBlock(sidechainID, blockID, sourceNodeID, destinationNodeID string) error {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return errors.New("sidechain not found")
	}

	block, err := sidechain.RetrieveBlock(blockID)
	if err != nil {
		return fmt.Errorf("block not found in sidechain %s: %v", sidechainID, err)
	}

	sourceNode, exists := sn.Nodes[sourceNodeID]
	if !exists {
		return fmt.Errorf("source node %s not found", sourceNodeID)
	}

	destNode, exists := sn.Nodes[destinationNodeID]
	if !exists {
		return fmt.Errorf("destination node %s not found", destinationNodeID)
	}

	// Encrypt block data before transferring
	encryptedBlock, err := sn.Encryption.EncryptData([]byte(blockID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt block data: %v", err)
	}

	// Perform the block sync via the network manager
	err = sn.NetworkManager.TransferData(sourceNode, destNode, encryptedBlock)
	if err != nil {
		return fmt.Errorf("failed to sync block %s: %v", blockID, err)
	}

	// Log the block sync operation
	err = sn.Ledger.RecordBlockSync(blockID, sourceNodeID, destinationNodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block sync: %v", err)
	}

	fmt.Printf("Block %s synced from node %s to node %s on sidechain %s\n", blockID, sourceNodeID, destinationNodeID, sidechainID)
	return nil
}

// BroadcastSidechainSubBlock broadcasts a sub-block to all nodes on a sidechain
func (sn *common.SidechainNetwork) BroadcastSidechainSubBlock(sidechainID, subBlockID string) error {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", sidechainID)
	}

	subBlock, err := sidechain.RetrieveSubBlock(subBlockID)
	if err != nil {
		return fmt.Errorf("sub-block %s not found on sidechain %s: %v", subBlockID, sidechainID, err)
	}

	// Encrypt the sub-block before broadcasting
	encryptedSubBlock, err := sn.Encryption.EncryptData([]byte(subBlockID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block data: %v", err)
	}

	// Broadcast the sub-block to all nodes in the sidechain
	for _, node := range sn.Nodes {
		err := sn.NetworkManager.BroadcastData(node, encryptedSubBlock)
		if err != nil {
			return fmt.Errorf("failed to broadcast sub-block %s to node %s: %v", subBlockID, node.NodeID, err)
		}
	}

	// Log the sub-block broadcast event
	err = sn.Ledger.RecordSubBlockBroadcast(subBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block broadcast: %v", err)
	}

	fmt.Printf("Sub-block %s broadcast to all nodes on sidechain %s\n", subBlockID, sidechainID)
	return nil
}

// ValidateSidechainNode validates a node’s integrity within a sidechain
func (sn *common.SidechainNetwork) ValidateSidechainNode(sidechainID, nodeID string) error {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", sidechainID)
	}

	node, exists := sn.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found in sidechain %s", nodeID, sidechainID)
	}

	// Validate the node's status and check if it's operational
	if !sn.NetworkManager.CheckNodeStatus(node) {
		return fmt.Errorf("node %s is not operational", nodeID)
	}

	// Log the node validation event
	err := sn.Ledger.RecordNodeValidation(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node validation: %v", err)
	}

	fmt.Printf("Node %s validated successfully on sidechain %s\n", nodeID, sidechainID)
	return nil
}

// HandleSidechainNetworkReconfiguration handles dynamic reconfiguration of a sidechain's network
func (sn *common.SidechainNetwork) HandleSidechainNetworkReconfiguration(sidechainID string, newTopology map[string]*common.Node) error {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", sidechainID)
	}

	// Log the network reconfiguration event
	err := sn.Ledger.RecordNetworkReconfiguration(time.Now())
	if err != nil {
		return fmt.Errorf("failed to log network reconfiguration: %v", err)
	}

	// Reconfigure the node topology within the sidechain
	sidechain.Consensus.ReconfigureNodes(newTopology)
	fmt.Printf("Network reconfiguration complete for sidechain %s\n", sidechainID)
	return nil
}

// RetrieveSidechainNode retrieves a node by its ID from a sidechain
func (sn *common.SidechainNetwork) RetrieveSidechainNode(nodeID string) (*common.Node, error) {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	node, exists := sn.Nodes[nodeID]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	fmt.Printf("Retrieved node %s\n", nodeID)
	return node, nil
}

// RetrieveSidechainBlock retrieves a block by its ID from a sidechain
func (sn *common.SidechainNetwork) RetrieveSidechainBlock(sidechainID, blockID string) (*common.SideBlock, error) {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return nil, fmt.Errorf("sidechain %s not found", sidechainID)
	}

	block, err := sidechain.RetrieveBlock(blockID)
	if err != nil {
		return nil, fmt.Errorf("block %s not found in sidechain %s: %v", blockID, sidechainID, err)
	}

	fmt.Printf("Retrieved block %s from sidechain %s\n", blockID, sidechainID)
	return block, nil
}

// RetrieveSidechainSubBlock retrieves a sub-block by its ID from a sidechain
func (sn *common.SidechainNetwork) RetrieveSidechainSubBlock(sidechainID, subBlockID string) (*common.SubBlock, error) {
	sn.mu.Lock()
	defer sn.mu.Unlock()

	sidechain, exists := sn.Sidechains[sidechainID]
	if !exists {
		return nil, fmt.Errorf("sidechain %s not found", sidechainID)
	}

	subBlock, err := sidechain.RetrieveSubBlock(subBlockID)
	if err != nil {
		return nil, fmt.Errorf("sub-block %s not found in sidechain %s: %v", subBlockID, sidechainID, err)
	}

	fmt.Printf("Retrieved sub-block %s from sidechain %s\n", subBlockID, sidechainID)
	return subBlock, nil
}
