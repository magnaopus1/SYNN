package layer2_consensus

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewProofOfCollaborationManager initializes the Proof-of-Collaboration manager
func NewProofOfCollaborationManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *ProofOfCollaborationManager {
	return &ProofOfCollaborationManager{
		Nodes:             make(map[string]*CollaborationNode),
		ActiveTasks:       make(map[string]*CollaborationTask),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AddCollaborationNode adds a new node to participate in Proof-of-Collaboration
func (poc *ProofOfCollaborationManager) AddCollaborationNode(nodeID string) (*CollaborationNode, error) {
	poc.mu.Lock()
	defer poc.mu.Unlock()

	// Encrypt node data
	nodeData := fmt.Sprintf("NodeID: %s", nodeID)
	encryptedData, err := poc.EncryptionService.EncryptData(nodeID, []byte(nodeData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt node data: %v", err)
	}

	// Use the encrypted data (e.g., log it or store it)
	fmt.Printf("Encrypted node data: %x\n", encryptedData)

	// Create the node
	node := &CollaborationNode{
		NodeID:     nodeID,
		Reputation: 100, // Default starting reputation
		Active:     true,
	}

	// Add the node to the PoC manager
	poc.Nodes[nodeID] = node

	// Log the node addition in the ledger
	poc.Ledger.BlockchainConsensusCoinLedger.RecordCollaborationNodeAddition(nodeID, nodeData) // Pass nodeData instead of time.Now()

	fmt.Printf("Collaboration node %s added to the PoC network\n", nodeID)
	return node, nil
}


// AssignCollaborationTask assigns a collaborative task to a group of nodes
func (poc *ProofOfCollaborationManager) AssignCollaborationTask(taskID string, nodes []string, computationData string) (*CollaborationTask, error) {
	poc.mu.Lock()
	defer poc.mu.Unlock()

	// Encrypt task data
	taskData := fmt.Sprintf("TaskID: %s, Nodes: %v, ComputationData: %s", taskID, nodes, computationData)
	encryptedData, err := poc.EncryptionService.EncryptData(taskID, []byte(taskData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt task data: %v", err)
	}

	// Convert encrypted data to a string representation (e.g., base64)
	encryptedDataStr := base64.StdEncoding.EncodeToString(encryptedData)

	// Create the collaboration task
	task := &CollaborationTask{
		TaskID:           taskID,
		AssignedNodes:    nodes,
		ComputationResult: "",
		CompletionStatus: "Pending",
		AssignedTime:     time.Now(),
		EncryptedData:    encryptedDataStr, // Use string representation here
	}

	// Add the task to active tasks
	poc.ActiveTasks[taskID] = task

	// Log the task assignment in the ledger (assuming you want to log the first node only)
	if len(nodes) > 0 {
		poc.Ledger.BlockchainConsensusCoinLedger.RecordCollaborationTaskAssignment(taskID, nodes[0]) // Log assignment for the first node
	}

	fmt.Printf("Collaboration task %s assigned to nodes %v\n", taskID, nodes)
	return task, nil
}

// CompleteCollaborationTask completes a collaborative task and distributes rewards
func (poc *ProofOfCollaborationManager) CompleteCollaborationTask(taskID, result string) error {
	poc.mu.Lock()
	defer poc.mu.Unlock()

	// Retrieve the collaboration task
	task, exists := poc.ActiveTasks[taskID]
	if !exists {
		return fmt.Errorf("collaboration task %s not found", taskID)
	}

	// Mark the task as completed and store the result
	task.ComputationResult = result
	task.CompletionStatus = "Completed"
	task.CompletedTime = time.Now()

	// Update node reputations for completing the task
	for _, nodeID := range task.AssignedNodes {
		node, exists := poc.Nodes[nodeID]
		if exists {
			node.Reputation += 10 // Increase reputation for successful collaboration
			node.LastCollabTime = time.Now()
		}
	}

	// Log the task completion in the ledger (removed err assignment)
	poc.Ledger.BlockchainConsensusCoinLedger.RecordCollaborationTaskCompletion(taskID, result) // Removed time.Now()

	fmt.Printf("Collaboration task %s completed with result: %s\n", taskID, result)
	return nil
}


// GetActiveTaskDetails retrieves the details of an active collaboration task
func (poc *ProofOfCollaborationManager) GetActiveTaskDetails(taskID string) (*CollaborationTask, error) {
	poc.mu.Lock()
	defer poc.mu.Unlock()

	// Retrieve the task
	task, exists := poc.ActiveTasks[taskID]
	if !exists {
		return nil, fmt.Errorf("collaboration task %s not found", taskID)
	}

	return task, nil
}

// GetNodeDetails retrieves the details of a collaboration node
func (poc *ProofOfCollaborationManager) GetNodeDetails(nodeID string) (*CollaborationNode, error) {
	poc.mu.Lock()
	defer poc.mu.Unlock()

	// Retrieve the node
	node, exists := poc.Nodes[nodeID]
	if !exists {
		return nil, fmt.Errorf("collaboration node %s not found", nodeID)
	}

	return node, nil
}

// generateUniqueID creates a cryptographically secure unique ID
func generateUniqueID() string {
	id := make([]byte, 16)
	rand.Read(id)
	return hex.EncodeToString(id)
}
