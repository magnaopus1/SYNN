package plasma

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)


// NewPlasmaNetwork initializes the PlasmaNetwork for managing Plasma nodes and blockchain data
func NewPlasmaNetwork(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.Manager) *common.PlasmaNetwork {
	return &common.PlasmaNetwork{
		PlasmaNodes:    make(map[string]*common.PlasmaNode),
		Blocks:         make(map[string]*common.PlasmaBlock),
		SubBlocks:      make(map[string]*common.PlasmaSubBlock),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// AddPlasmaNode adds a new Plasma node to the Plasma network
func (pn *common.PlasmaNetwork) AddPlasmaNode(nodeID, ipAddress string, nodeType common.PlasmaNodeType) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	if _, exists := pn.PlasmaNodes[nodeID]; exists {
		return errors.New("Plasma node already exists")
	}

	// Create the new Plasma node
	newNode := &common.PlasmaNode{
		NodeID:    nodeID,
		IPAddress: ipAddress,
		NodeType:  nodeType,
	}

	pn.PlasmaNodes[nodeID] = newNode

	// Log the new Plasma node creation
	err := pn.Ledger.RecordNodeAddition(nodeID, ipAddress, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node addition: %v", err)
	}

	fmt.Printf("Plasma node %s added to the Plasma network\n", nodeID)
	return nil
}

// RemovePlasmaNode removes a Plasma node from the Plasma network
func (pn *common.PlasmaNetwork) RemovePlasmaNode(nodeID string) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	if _, exists := pn.PlasmaNodes[nodeID]; !exists {
		return errors.New("Plasma node does not exist")
	}

	delete(pn.PlasmaNodes, nodeID)

	// Log the Plasma node removal
	err := pn.Ledger.RecordNodeRemoval(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node removal: %v", err)
	}

	fmt.Printf("Plasma node %s removed from the Plasma network\n", nodeID)
	return nil
}

// SyncBlock synchronizes a block between Plasma nodes on the Plasma network
func (pn *common.PlasmaNetwork) SyncBlock(blockID string, sourceNodeID string, destinationNodeID string) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	block, exists := pn.Blocks[blockID]
	if !exists {
		return fmt.Errorf("block %s not found", blockID)
	}

	sourceNode, exists := pn.PlasmaNodes[sourceNodeID]
	if !exists {
		return fmt.Errorf("source Plasma node %s not found", sourceNodeID)
	}

	destNode, exists := pn.PlasmaNodes[destinationNodeID]
	if !exists {
		return fmt.Errorf("destination Plasma node %s not found", destinationNodeID)
	}

	// Encrypt block data before transferring
	encryptedBlock, err := pn.Encryption.EncryptData([]byte(blockID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt block data: %v", err)
	}

	// Perform the block sync via the network manager
	err = pn.NetworkManager.TransferData(sourceNode, destNode, encryptedBlock)
	if err != nil {
		return fmt.Errorf("failed to sync block %s: %v", blockID, err)
	}

	// Log the block sync operation
	err = pn.Ledger.RecordBlockSync(blockID, sourceNodeID, destinationNodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block sync: %v", err)
	}

	fmt.Printf("Block %s synced from Plasma node %s to Plasma node %s\n", blockID, sourceNodeID, destinationNodeID)
	return nil
}

// BroadcastSubBlock broadcasts a sub-block to all Plasma nodes in the Plasma network
func (pn *common.PlasmaNetwork) BroadcastSubBlock(subBlockID string) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	subBlock, exists := pn.SubBlocks[subBlockID]
	if !exists {
		return fmt.Errorf("sub-block %s not found", subBlockID)
	}

	// Encrypt the sub-block before broadcasting
	encryptedSubBlock, err := pn.Encryption.EncryptData([]byte(subBlockID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block data: %v", err)
	}

	// Broadcast the sub-block to all Plasma nodes
	for _, node := range pn.PlasmaNodes {
		err := pn.NetworkManager.BroadcastData(node, encryptedSubBlock)
		if err != nil {
			return fmt.Errorf("failed to broadcast sub-block %s to Plasma node %s: %v", subBlockID, node.NodeID, err)
		}
	}

	// Log the sub-block broadcast event
	err = pn.Ledger.RecordSubBlockBroadcast(subBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block broadcast: %v", err)
	}

	fmt.Printf("Sub-block %s broadcast to all Plasma nodes\n", subBlockID)
	return nil
}

// ValidatePlasmaNode validates a Plasma node’s integrity in the Plasma network based on consensus rules
func (pn *common.PlasmaNetwork) ValidatePlasmaNode(nodeID string) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	node, exists := pn.PlasmaNodes[nodeID]
	if !exists {
		return fmt.Errorf("Plasma node %s not found", nodeID)
	}

	// Validate the Plasma node's status and check if it's operational
	if !pn.NetworkManager.CheckNodeStatus(node) {
		return fmt.Errorf("Plasma node %s is not operational", nodeID)
	}

	// Log the node validation
	err := pn.Ledger.RecordNodeValidation(nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node validation: %v", err)
	}

	fmt.Printf("Plasma node %s validated successfully\n", nodeID)
	return nil
}

// HandleNetworkReconfiguration handles dynamic reconfiguration of the Plasma network
func (pn *common.PlasmaNetwork) HandleNetworkReconfiguration(newTopology map[string]*common.PlasmaNode) error {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	// Log the network reconfiguration event
	err := pn.Ledger.RecordNetworkReconfiguration(time.Now())
	if err != nil {
		return fmt.Errorf("failed to log network reconfiguration: %v", err)
	}

	// Reconfigure the Plasma node topology
	pn.PlasmaNodes = newTopology
	fmt.Println("Plasma network reconfiguration complete")
	return nil
}

// RetrievePlasmaNode retrieves a Plasma node by its ID from the Plasma network
func (pn *common.PlasmaNetwork) RetrievePlasmaNode(nodeID string) (*common.PlasmaNode, error) {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	node, exists := pn.PlasmaNodes[nodeID]
	if !exists {
		return nil, fmt.Errorf("Plasma node %s not found", nodeID)
	}

	fmt.Printf("Retrieved Plasma node %s\n", nodeID)
	return node, nil
}

// RetrieveBlock retrieves a block by its ID from the Plasma network
func (pn *common.PlasmaNetwork) RetrieveBlock(blockID string) (*common.PlasmaBlock, error) {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	block, exists := pn.Blocks[blockID]
	if !exists {
		return nil, fmt.Errorf("block %s not found", blockID)
	}

	fmt.Printf("Retrieved block %s\n", blockID)
	return block, nil
}

// RetrieveSubBlock retrieves a sub-block by its ID from the Plasma network
func (pn *common.PlasmaNetwork) RetrieveSubBlock(subBlockID string) (*common.PlasmaSubBlock, error) {
	pn.mu.Lock()
	defer pn.mu.Unlock()

	subBlock, exists := pn.SubBlocks[subBlockID]
	if !exists {
		return nil, fmt.Errorf("sub-block %s not found", subBlockID)
	}

	fmt.Printf("Retrieved sub-block %s\n", subBlockID)
	return subBlock, nil
}
