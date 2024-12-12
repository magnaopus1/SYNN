package plasma

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/network"
)

// NewPlasmaNode initializes a new PlasmaNode with the provided parameters
func NewPlasmaNode(nodeID string, nodeType common.PlasmaNodeType, ipAddress string, encryptionService *encryption.Encryption) *common.PlasmaNode {
	return &common.PlasmaNode{
		NodeID:         nodeID,
		NodeType:       nodeType,
		IPAddress:      ipAddress,
		Encryption:     encryptionService,
		LastActiveTime: time.Now(),
		NodeHealth:     network.Healthy,         // Default health status
		NodeState:      network.Active,          // Default state of the node
	}
}

// UpdateHealthStatus updates the health status of the Plasma node
func (pn *common.PlasmaNode) UpdateHealthStatus(newStatus common.NodeHealthStatus) {
	pn.NodeHealth = newStatus
	pn.LastActiveTime = time.Now()
	fmt.Printf("Node %s health status updated to %v\n", pn.NodeID, newStatus)
}

// DeactivateNode deactivates the Plasma node from active participation
func (pn *common.PlasmaNode) DeactivateNode() {
	pn.NodeState = network.Inactive
	pn.LastActiveTime = time.Now()
	fmt.Printf("Node %s is now inactive\n", pn.NodeID)
}

// ReactivateNode reactivates the Plasma node for participation in the Plasma network
func (pn *common.PlasmaNode) ReactivateNode() {
	pn.NodeState = network.Active
	pn.LastActiveTime = time.Now()
	fmt.Printf("Node %s is now active again\n", pn.NodeID)
}

// ValidateNode ensures the node is operational and in compliance with network rules
func (pn *common.PlasmaNode) ValidateNode() error {
	if pn.NodeHealth != network.Healthy {
		return fmt.Errorf("node %s is not healthy, current status: %v", pn.NodeID, pn.NodeHealth)
	}
	if pn.NodeState != network.Active {
		return fmt.Errorf("node %s is inactive", pn.NodeID)
	}

	// Log the node validation operation
	fmt.Printf("Node %s validated successfully\n", pn.NodeID)
	return nil
}

// EncryptData encrypts a given data block using the node's encryption service
func (pn *common.PlasmaNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := pn.Encryption.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data on node %s: %v", pn.NodeID, err)
	}
	return encryptedData, nil
}

// DecryptData decrypts a given data block using the node's encryption service
func (pn *common.PlasmaNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := pn.Encryption.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data on node %s: %v", pn.NodeID, err)
	}
	return decryptedData, nil
}

// HandleSubBlock handles the validation and processing of a sub-block on the Plasma node
func (pn *common.PlasmaNode) HandleSubBlock(subBlock *common.PlasmaSubBlock) error {
	if pn.NodeState != network.Active {
		return fmt.Errorf("node %s is inactive, cannot handle sub-blocks", pn.NodeID)
	}

	// Validate the sub-block (example, could include consensus-related logic)
	fmt.Printf("Node %s is handling sub-block %s\n", pn.NodeID, subBlock.SubBlockID)
	// You can add additional processing/validation logic for the sub-block here

	return nil
}

// FetchNodeStatus retrieves the current operational status of the Plasma node
func (pn *common.PlasmaNode) FetchNodeStatus() common.NodeState {
	return pn.NodeState
}

// FetchHealthStatus retrieves the health status of the Plasma node
func (pn *common.PlasmaNode) FetchHealthStatus() common.NodeHealthStatus {
	return pn.NodeHealth
}

// IsHealthy checks if the node is in a healthy state
func (pn *common.PlasmaNode) IsHealthy() bool {
	return pn.NodeHealth == network.Healthy
}

// IsActive checks if the node is active and can participate in the network
func (pn *common.PlasmaNode) IsActive() bool {
	return pn.NodeState == network.Active
}

