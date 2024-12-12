package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRollupNode initializes a new RollupNode
func NewRollupNode(nodeID, ipAddress string, nodeType common.NodeType, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.RollupNode {
	return &common.RollupNode{
		NodeID:        nodeID,
		IPAddress:     ipAddress,
		NodeType:      nodeType,
		ConnectedNodes: make(map[string]*common.RollupNode),
		Ledger:        ledgerInstance,
		Encryption:    encryptionService,
		Consensus:     consensus,
	}
}

// AddConnectedNode connects another node to this rollup node
func (rn *common.RollupNode) AddConnectedNode(connectedNode *common.RollupNode) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if _, exists := rn.ConnectedNodes[connectedNode.NodeID]; exists {
		return errors.New("node is already connected")
	}

	rn.ConnectedNodes[connectedNode.NodeID] = connectedNode

	// Log the connection event in the ledger
	err := rn.Ledger.RecordNodeConnection(rn.NodeID, connectedNode.NodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node connection: %v", err)
	}

	fmt.Printf("Node %s connected to node %s\n", rn.NodeID, connectedNode.NodeID)
	return nil
}

// RemoveConnectedNode removes a connected node from this rollup node
func (rn *common.RollupNode) RemoveConnectedNode(nodeID string) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if _, exists := rn.ConnectedNodes[nodeID]; !exists {
		return errors.New("node is not connected")
	}

	delete(rn.ConnectedNodes, nodeID)

	// Log the disconnection event in the ledger
	err := rn.Ledger.RecordNodeDisconnection(rn.NodeID, nodeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log node disconnection: %v", err)
	}

	fmt.Printf("Node %s disconnected from node %s\n", rn.NodeID, nodeID)
	return nil
}

// AggregateTransactions creates a new RollupBatch by aggregating transactions
func (rn *common.RollupNode) AggregateTransactions(transactions []*common.Transaction) (*common.RollupBatch, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Create a new batch
	batch := &common.RollupBatch{
		BatchID:      common.GenerateUniqueID(), // Assuming a utility function to generate unique IDs
		Transactions: transactions,
		MerkleRoot:   common.GenerateMerkleRoot(transactions), // Assuming helper function to generate Merkle root
		Timestamp:    time.Now(),
	}

	// Log the batch creation in the ledger
	err := rn.Ledger.RecordBatchCreation(batch.BatchID, batch.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to log batch creation: %v", err)
	}

	fmt.Printf("Rollup batch %s created with %d transactions\n", batch.BatchID, len(transactions))
	return batch, nil
}

// ValidateBatch validates a rollup batch using Synnergy Consensus
func (rn *common.RollupNode) ValidateBatch(batch *common.RollupBatch) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Validate the batch using Synnergy Consensus
	err := rn.Consensus.ValidateBatch(batch.BatchID, batch.Transactions)
	if err != nil {
		return fmt.Errorf("batch validation failed: %v", err)
	}

	// Log the batch validation in the ledger
	err = rn.Ledger.RecordBatchValidation(batch.BatchID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log batch validation: %v", err)
	}

	fmt.Printf("Rollup batch %s validated successfully\n", batch.BatchID)
	return nil
}

// BroadcastBatch broadcasts the batch to all connected nodes
func (rn *common.RollupNode) BroadcastBatch(batch *common.RollupBatch) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Encrypt the batch before broadcasting
	encryptedBatch, err := rn.Encryption.EncryptData([]byte(batch.BatchID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt batch data: %v", err)
	}

	// Broadcast the batch to connected nodes
	for _, connectedNode := range rn.ConnectedNodes {
		err := rn.Ledger.RecordBatchBroadcast(batch.BatchID, connectedNode.NodeID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log batch broadcast: %v", err)
		}
		// Assuming the network manager will handle sending the encrypted batch
		fmt.Printf("Batch %s broadcast to node %s\n", batch.BatchID, connectedNode.NodeID)
	}

	return nil
}
