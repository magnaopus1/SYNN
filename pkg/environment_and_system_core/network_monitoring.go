package environment_and_system_core

import (
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// CheckNetworkStatus retrieves the current network status, encrypts it, and logs it in the ledger.
func CheckNetworkStatus(ledger *ledger.Ledger) (string, error) {
	// Retrieve network status
	status, err := ledger.EnvironmentSystemCoreLedger.NetworkManager.GetStatus()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve network status: %w", err)
	}

	// Encrypt the network status
	encryptedStatus, err := common.EncryptData(status)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt network status: %w", err)
	}

	// Log the encrypted status in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogNetworkStatus(encryptedStatus); err != nil {
		return "", fmt.Errorf("failed to log network status in ledger: %w", err)
	}

	// Log success
	log.Println("Network status checked, encrypted, and logged in ledger.")
	return status, nil
}


// GetBlockHeight retrieves the current block height from the ledger.
func GetBlockHeight(ledger *ledger.Ledger) (int, error) {
	// Fetch block height
	blockHeight, err := ledger.BlockchainConsensusCoinLedger.GetCurrentBlockHeight()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve current block height: %w", err)
	}

	// Log success
	log.Printf("Current block height retrieved: %d", blockHeight)
	return blockHeight, nil
}

// GetCurrentSubBlock retrieves information about the latest sub-block from the ledger.
func GetCurrentSubBlock(ledger *ledger.Ledger) (string, error) {
	// Fetch latest sub-block
	subBlockInfo, err := ledger.BlockchainConsensusCoinLedger.GetLatestSubBlock()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve latest sub-block: %w", err)
	}

	// Log success
	log.Printf("Latest sub-block retrieved: %s", subBlockInfo)
	return subBlockInfo, nil
}


// ValidateNodeHealth validates the health status of a specific node, encrypts the data, and logs it in the ledger.
func ValidateNodeHealth(ledger *ledger.Ledger, nodeID string) (string, error) {
	// Validate input
	if nodeID == "" {
		return "", fmt.Errorf("nodeID cannot be empty")
	}

	// Retrieve node health status
	nodeStatus, err := ledger.EnvironmentSystemCoreLedger.NodeManager.GetHealthStatus(nodeID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve health status for node %s: %w", nodeID, err)
	}

	// Encrypt the node health status
	encryptedStatus, err := common.EncryptData(nodeStatus)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt node health status for node %s: %w", nodeID, err)
	}

	// Record the health status in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordNodeHealth(nodeID, encryptedStatus); err != nil {
		return "", fmt.Errorf("failed to record node health status for node %s in ledger: %w", nodeID, err)
	}

	// Log success
	log.Printf("Health status for node %s validated, encrypted, and recorded in ledger.", nodeID)
	return nodeStatus, nil
}

// SetNodeMaintenanceMode enables or disables maintenance mode for a specific node.
func SetNodeMaintenanceMode(nodeID string, enable bool, ledger *ledger.Ledger) error {
	// Validate input
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}

	// Determine status
	status := "enabled"
	if !enable {
		status = "disabled"
	}

	// Update maintenance mode in the system
	if err := ledger.EnvironmentSystemCoreLedger.NodeManager.SetMaintenanceMode(nodeID, enable); err != nil {
		return fmt.Errorf("failed to set maintenance mode for node %s: %w", nodeID, err)
	}

	// Record the maintenance mode change in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.RecordMaintenanceEvent(nodeID, status); err != nil {
		return fmt.Errorf("failed to record maintenance mode change for node %s in ledger: %w", nodeID, err)
	}

	// Log success
	log.Printf("Maintenance mode %s for node %s and recorded in ledger.", status, nodeID)
	return nil
}


// CheckNodeSyncStatus checks if a node is synced with the blockchain.
func CheckNodeSyncStatus(nodeID string, ledger *ledger.Ledger) (bool, error) {
	// Validate input
	if nodeID == "" {
		return false, fmt.Errorf("nodeID cannot be empty")
	}

	// Check node sync status from the ledger
	isSynced, err := ledger.EnvironmentSystemCoreLedger.GetNodeSyncStatus(nodeID)
	if err != nil {
		return false, fmt.Errorf("failed to check sync status for node %s: %w", nodeID, err)
	}

	// Log success
	log.Printf("Node %s sync status: %v", nodeID, isSynced)
	return isSynced, nil
}


// GetEnvironmentVariable retrieves the value of a specified environment variable from the ledger.
func GetEnvironmentVariable(ledger *ledger.Ledger, variableName string) (string, error) {
	// Validate input
	if variableName == "" {
		return "", fmt.Errorf("variableName cannot be empty")
	}

	// Fetch the variable value
	value, exists := ledger.EnvironmentSystemCoreLedger.GetEnvironmentVariable(variableName)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", variableName)
	}

	// Log success
	log.Printf("Environment variable %s retrieved: %s", variableName, value)
	return value, nil
}


// SetEnvironmentVariable sets the value of a specified environment variable in the ledger.
func SetEnvironmentVariable(ledger *ledger.Ledger, variableName, value string) error {
	// Validate input
	if variableName == "" || value == "" {
		return fmt.Errorf("variableName and value cannot be empty")
	}

	// Set the variable in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.SetEnvironmentVariable(variableName, value); err != nil {
		return fmt.Errorf("failed to set environment variable %s: %w", variableName, err)
	}

	// Log success
	log.Printf("Environment variable %s set to %s.", variableName, value)
	return nil
}


// RecordBlockEvent records a block-related event in the ledger with encrypted details.
func RecordBlockEvent(ledger *ledger.Ledger, eventType, blockID, message string) error {
	// Validate input
	if eventType == "" || blockID == "" || message == "" {
		return fmt.Errorf("eventType, blockID, and message cannot be empty")
	}

	// Encrypt the message
	encryptedMessage, err := common.EncryptData(message)
	if err != nil {
		return fmt.Errorf("failed to encrypt block event message: %w", err)
	}

	// Record the block event in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.LogBlockEvent(eventType, blockID, encryptedMessage); err != nil {
		return fmt.Errorf("failed to log block event %s for block %s: %w", eventType, blockID, err)
	}

	// Log success
	log.Printf("Block event logged: %s for block %s.", eventType, blockID)
	return nil
}


// QueryBlockchainParameter retrieves the value of a specified blockchain parameter from the ledger.
func QueryBlockchainParameter(ledger *ledger.Ledger, paramName string) (string, error) {
	// Validate input
	if paramName == "" {
		return "", fmt.Errorf("parameter name cannot be empty")
	}

	// Fetch the parameter value from the ledger
	paramValue, err := ledger.EnvironmentSystemCoreLedger.GetBlockchainParameter(paramName)
	if err != nil {
		return "", fmt.Errorf("failed to query blockchain parameter %s: %w", paramName, err)
	}

	// Log success
	log.Printf("Blockchain parameter %s queried successfully: %s", paramName, paramValue)
	return paramValue, nil
}

// SetBlockchainParameter updates the value of a specified blockchain parameter in the ledger.
func SetBlockchainParameter(ledger *ledger.Ledger, paramName string, paramValue string) error {
	// Validate input
	if paramName == "" || paramValue == "" {
		return fmt.Errorf("parameter name and value cannot be empty")
	}

	// Set the parameter value in the ledger
	if err := ledger.EnvironmentSystemCoreLedger.SetBlockchainParameter(paramName, paramValue); err != nil {
		return fmt.Errorf("failed to set blockchain parameter %s: %w", paramName, err)
	}

	// Log success
	log.Printf("Blockchain parameter %s set to %s.", paramName, paramValue)
	return nil
}


// CheckBlockValidity verifies the integrity and validity of a block in the blockchain.
func CheckBlockValidity(ledger *ledger.Ledger, blockID string) (bool, error) {
	// Validate input
	if blockID == "" {
		return false, fmt.Errorf("blockID cannot be empty")
	}

	// Validate the block integrity
	isValid, err := ledger.EnvironmentSystemCoreLedger.ValidateBlockIntegrity(blockID)
	if err != nil {
		return false, fmt.Errorf("failed to validate block integrity for block %s: %w", blockID, err)
	}

	// Log success
	log.Printf("Block %s validity checked: %v", blockID, isValid)
	return isValid, nil
}



// ResetNodeState resets a node to its default state.
func ResetNodeState(nodeID string, ledger *ledger.Ledger) error {
	// Validate input
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}

	// Reset the node state
	if err := ledger.EnvironmentSystemCoreLedger.ResetNodeToDefault(nodeID); err != nil {
		return fmt.Errorf("failed to reset node state for node %s: %w", nodeID, err)
	}

	// Log success
	log.Printf("Node %s state reset to default successfully.", nodeID)
	return nil
}

