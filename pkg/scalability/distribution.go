package scalability

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"
)

// NewDistributionSystem initializes the distribution system
func NewDistributionSystem(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.DistributionSystem {
	return &common.DistributionSystem{
		Nodes:             []*common.DistributionNode{},
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AddNode adds a new node to the distribution system
func (ds *common.DistributionSystem) AddNode(nodeID, nodeType string, weight int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	node := &common.DistributionNode{
		NodeID:   nodeID,
		NodeType: nodeType,
		Load:     0,
		Weight:   weight,
		LastTask: time.Now(),
	}

	ds.Nodes = append(ds.Nodes, node)

	// Log the addition of the new node
	err := ds.Ledger.RecordNodeAddition(nodeID, nodeType, time.Now())
	if err != nil {
		fmt.Printf("Failed to log node addition: %v\n", err)
	}

	fmt.Printf("Node %s of type %s added to the distribution system\n", nodeID, nodeType)
}

// AdaptiveDistribution dynamically selects a node based on the current load
func (ds *common.DistributionSystem) AdaptiveDistribution(taskID string) (*common.DistributionNode, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var selectedNode *common.DistributionNode
	for _, node := range ds.Nodes {
		if selectedNode == nil || node.Load < selectedNode.Load {
			selectedNode = node
		}
	}

	if selectedNode == nil {
		return nil, errors.New("no available nodes for adaptive distribution")
	}

	// Simulate assigning the task to the selected node
	selectedNode.Load++
	selectedNode.LastTask = time.Now()

	// Log the distribution in the ledger
	err := ds.Ledger.RecordTaskDistribution(taskID, selectedNode.NodeID, "adaptive", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log adaptive distribution: %v", err)
	}

	fmt.Printf("Task %s assigned to node %s using adaptive distribution\n", taskID, selectedNode.NodeID)
	return selectedNode, nil
}

// RoundRobinDistribution selects nodes in a round-robin manner for equal load balancing
func (ds *common.DistributionSystem) RoundRobinDistribution(taskID string, lastNodeIndex int) (*common.DistributionNode, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if len(ds.Nodes) == 0 {
		return nil, errors.New("no available nodes for round-robin distribution")
	}

	nextIndex := (lastNodeIndex + 1) % len(ds.Nodes)
	selectedNode := ds.Nodes[nextIndex]

	// Simulate assigning the task to the selected node
	selectedNode.Load++
	selectedNode.LastTask = time.Now()

	// Log the distribution in the ledger
	err := ds.Ledger.RecordTaskDistribution(taskID, selectedNode.NodeID, "round-robin", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log round-robin distribution: %v", err)
	}

	fmt.Printf("Task %s assigned to node %s using round-robin distribution\n", taskID, selectedNode.NodeID)
	return selectedNode, nil
}

// WeightedDistribution selects nodes based on their assigned weight
func (ds *common.DistributionSystem) WeightedDistribution(taskID string) (*common.DistributionNode, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if len(ds.Nodes) == 0 {
		return nil, errors.New("no available nodes for weighted distribution")
	}

	var totalWeight int
	for _, node := range ds.Nodes {
		totalWeight += node.Weight
	}

	if totalWeight == 0 {
		return nil, errors.New("no weighted nodes available for distribution")
	}

	// Select node based on weight (basic random weighted selection)
	var selectedNode *common.DistributionNode
	randomWeight := common.GenerateRandomWeight(totalWeight)
	cumulativeWeight := 0

	for _, node := range ds.Nodes {
		cumulativeWeight += node.Weight
		if cumulativeWeight >= randomWeight {
			selectedNode = node
			break
		}
	}

	// Simulate assigning the task to the selected node
	selectedNode.Load++
	selectedNode.LastTask = time.Now()

	// Log the distribution in the ledger
	err := ds.Ledger.RecordTaskDistribution(taskID, selectedNode.NodeID, "weighted", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log weighted distribution: %v", err)
	}

	fmt.Printf("Task %s assigned to node %s using weighted distribution\n", taskID, selectedNode.NodeID)
	return selectedNode, nil
}
