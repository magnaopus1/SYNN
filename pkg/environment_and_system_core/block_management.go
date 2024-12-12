package environment_and_system_core

import (
	"fmt"
	"log"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Mutex for node operations and block management
var blockManagementLock sync.Mutex

// logChainActivity records blockchain-related activity in the ledger.
func logChainActivity(ledgerInstance *ledger.Ledger, activityType, description string) error {
    if ledgerInstance == nil {
        return fmt.Errorf("logChainActivity: ledger instance cannot be nil")
    }
    if activityType == "" {
        return fmt.Errorf("logChainActivity: activityType cannot be empty")
    }
    if description == "" {
        return fmt.Errorf("logChainActivity: description cannot be empty")
    }

    entry := ledger.ChainActivityLog{
        Timestamp:    time.Now().UTC(),
        ActivityType: activityType,
        Description:  description,
    }

    if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordChainActivity(entry); err != nil {
        return fmt.Errorf("logChainActivity: failed to record chain activity: %w", err)
    }

    log.Printf("[INFO] Chain activity logged: Type=%s, Description=%s", activityType, description)
    return nil
}



// checkBlockFinality verifies if a block has reached finality.
func checkBlockFinality(ledgerInstance *ledger.Ledger, blockID string) (bool, error) {
    // Validate inputs
    if ledgerInstance == nil {
        return false, fmt.Errorf("checkBlockFinality: ledger instance cannot be nil")
    }
    if blockID == "" {
        return false, fmt.Errorf("checkBlockFinality: blockID cannot be empty")
    }

    // Retrieve block finality status
    finalityStatus, err := ledgerInstance.EnvironmentSystemCoreLedger.GetBlockFinalityStatus(blockID)
    if err != nil {
        return false, fmt.Errorf("checkBlockFinality: failed to check finality for block %s: %w", blockID, err)
    }

    // Log status
    log.Printf("Block finality for block %s: %t", blockID, finalityStatus)
    return finalityStatus, nil
}


// suspendNodeOperations temporarily halts operations on a node.
func suspendNodeOperations(ledgerInstance *ledger.Ledger, nodeID string) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("suspendNodeOperations: ledger instance cannot be nil")
    }
    if nodeID == "" {
        return fmt.Errorf("suspendNodeOperations: nodeID cannot be empty")
    }

    // Update node status in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateNodeStatus(nodeID, "suspended"); err != nil {
        return fmt.Errorf("suspendNodeOperations: failed to suspend operations for node %s: %w", nodeID, err)
    }

    // Log success
    log.Printf("Node %s operations suspended.", nodeID)
    return nil
}


// resumeNodeOperations reactivates a suspended node.
func resumeNodeOperations(ledgerInstance *ledger.Ledger, nodeID string) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("resumeNodeOperations: ledger instance cannot be nil")
    }
    if nodeID == "" {
        return fmt.Errorf("resumeNodeOperations: nodeID cannot be empty")
    }

    // Update node status in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateNodeStatus(nodeID, "active"); err != nil {
        return fmt.Errorf("resumeNodeOperations: failed to resume operations for node %s: %w", nodeID, err)
    }

    // Log success
    log.Printf("Node %s operations resumed.", nodeID)
    return nil
}


// queryNodeRole retrieves the role assigned to a specific node.
func queryNodeRole(ledgerInstance *ledger.Ledger, nodeID string) (string, error) {
    // Validate inputs
    if ledgerInstance == nil {
        return "", fmt.Errorf("queryNodeRole: ledger instance cannot be nil")
    }
    if nodeID == "" {
        return "", fmt.Errorf("queryNodeRole: nodeID cannot be empty")
    }

    // Fetch node role from the ledger
    role, err := ledgerInstance.EnvironmentSystemCoreLedger.GetNodeRole(nodeID)
    if err != nil {
        return "", fmt.Errorf("queryNodeRole: failed to retrieve role for node %s: %w", nodeID, err)
    }

    // Log success
    log.Printf("Role retrieved for node %s: %s", nodeID, role)
    return role, nil
}


// assignNodeRole assigns a role to a specific node.
func assignNodeRole(ledgerInstance *ledger.Ledger, nodeID, role string) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("assignNodeRole: ledger instance cannot be nil")
    }
    if nodeID == "" || role == "" {
        return fmt.Errorf("assignNodeRole: nodeID and role cannot be empty")
    }

    // Set node role in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.SetNodeRole(nodeID, role); err != nil {
        return fmt.Errorf("assignNodeRole: failed to assign role %s to node %s: %w", role, nodeID, err)
    }

    // Log success
    log.Printf("Role %s assigned to node %s.", role, nodeID)
    return nil
}


// fetchBlockData retrieves data for a specified block.
func fetchBlockData(ledgerInstance *ledger.Ledger, blockID string) (ledger.BlockData, error) {
    // Validate inputs
    if ledgerInstance == nil {
        return ledger.BlockData{}, fmt.Errorf("fetchBlockData: ledger instance cannot be nil")
    }
    if blockID == "" {
        return ledger.BlockData{}, fmt.Errorf("fetchBlockData: blockID cannot be empty")
    }

    // Fetch block data from the ledger
    data, err := ledgerInstance.EnvironmentSystemCoreLedger.GetBlockData(blockID)
    if err != nil {
        return ledger.BlockData{}, fmt.Errorf("fetchBlockData: failed to fetch data for block %s: %w", blockID, err)
    }

    // Log success
    log.Printf("Data fetched for block %s.", blockID)
    return data, nil
}


// setBlockVerificationThreshold updates the verification threshold for blocks.
func setBlockVerificationThreshold(ledgerInstance *ledger.Ledger, threshold int) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("setBlockVerificationThreshold: ledger instance cannot be nil")
    }
    if threshold <= 0 {
        return fmt.Errorf("setBlockVerificationThreshold: threshold must be greater than zero")
    }

    // Update threshold in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateVerificationThreshold(threshold); err != nil {
        return fmt.Errorf("setBlockVerificationThreshold: failed to set verification threshold: %w", err)
    }

    // Log success
    log.Printf("Block verification threshold set to %d.", threshold)
    return nil
}


// retrieveNetworkMetrics gathers key performance metrics for the network.
func retrieveNetworkMetrics(ledgerInstance *ledger.Ledger) (map[string]interface{}, error) {
    // Validate inputs
    if ledgerInstance == nil {
        return nil, fmt.Errorf("retrieveNetworkMetrics: ledger instance cannot be nil")
    }

    // Retrieve metrics from the ledger
    metrics, err := ledgerInstance.EnvironmentSystemCoreLedger.GetNetworkMetrics()
    if err != nil {
        return nil, fmt.Errorf("retrieveNetworkMetrics: failed to retrieve metrics: %w", err)
    }

    // Format metrics into a map
    metricsMap := map[string]interface{}{
        "NodeCount":       metrics.NodeCount,
        "TransactionRate": metrics.TransactionRate,
        "BlockLatency":    metrics.BlockLatency,
    }

    // Log success
    log.Printf("Network metrics retrieved.")
    return metricsMap, nil
}


// logNodeHealthCheck records the health status of a node.
func logNodeHealthCheck(ledgerInstance *ledger.Ledger, nodeID string, healthScore int) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("logNodeHealthCheck: ledger instance cannot be nil")
    }
    if nodeID == "" {
        return fmt.Errorf("logNodeHealthCheck: nodeID cannot be empty")
    }
    if healthScore < 0 || healthScore > 100 {
        return fmt.Errorf("logNodeHealthCheck: healthScore must be between 0 and 100")
    }

    // Create health log entry
    entry := ledger.NodeHealthLog{
        NodeID:      nodeID,
        HealthScore: healthScore,
        Timestamp:   time.Now().UTC(),
    }

    // Record health status in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.RecordNodeHealth(entry); err != nil {
        return fmt.Errorf("logNodeHealthCheck: failed to log health for node %s: %w", nodeID, err)
    }

    // Log success
    log.Printf("Health check logged for node %s with score %d.", nodeID, healthScore)
    return nil
}



// reconfigureNode updates the configuration of a specific node.
func reconfigureNode(ledgerInstance *ledger.Ledger, nodeID string, config map[string]interface{}) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("reconfigureNode: ledger instance cannot be nil")
    }
    if nodeID == "" {
        return fmt.Errorf("reconfigureNode: nodeID cannot be empty")
    }
    if len(config) == 0 {
        return fmt.Errorf("reconfigureNode: configuration cannot be empty")
    }

    // Update node configuration in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateNodeConfig(nodeID, config); err != nil {
        return fmt.Errorf("reconfigureNode: failed to reconfigure node %s: %w", nodeID, err)
    }

    // Log success
    log.Printf("Node %s successfully reconfigured with new configuration.", nodeID)
    return nil
}


// activateEmergencyProtocol sets the system to emergency mode.
func activateEmergencyProtocol(ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("activateEmergencyProtocol: ledger instance cannot be nil")
    }

    // Activate emergency status
    if err := ledgerInstance.EnvironmentSystemCoreLedger.SetEmergencyStatus(true); err != nil {
        return fmt.Errorf("activateEmergencyProtocol: failed to activate emergency protocol: %w", err)
    }

    // Log success
    log.Printf("Emergency protocol successfully activated.")
    return nil
}


// deactivateEmergencyProtocol disables the system's emergency mode.
func deactivateEmergencyProtocol(ledgerInstance *ledger.Ledger) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("deactivateEmergencyProtocol: ledger instance cannot be nil")
    }

    // Deactivate emergency status
    if err := ledgerInstance.EnvironmentSystemCoreLedger.SetEmergencyStatus(false); err != nil {
        return fmt.Errorf("deactivateEmergencyProtocol: failed to deactivate emergency protocol: %w", err)
    }

    // Log success
    log.Printf("Emergency protocol successfully deactivated.")
    return nil
}



// verifyEnvironmentStatus checks the current status of the environment.
func verifyEnvironmentStatus(ledgerInstance *ledger.Ledger) (bool, error) {
    // Validate inputs
    if ledgerInstance == nil {
        return false, fmt.Errorf("verifyEnvironmentStatus: ledger instance cannot be nil")
    }

    // Check environment status
    status, err := ledgerInstance.EnvironmentSystemCoreLedger.CheckEnvironmentStatus()
    if err != nil {
        return false, fmt.Errorf("verifyEnvironmentStatus: failed to verify environment status: %w", err)
    }

    // Log status
    log.Printf("Environment status verified: %t", status)
    return status, nil
}


// retrieveNodeHealthScore fetches the health score of a specific node.
func retrieveNodeHealthScore(ledgerInstance *ledger.Ledger, nodeID string) (int, error) {
    // Validate inputs
    if ledgerInstance == nil {
        return 0, fmt.Errorf("retrieveNodeHealthScore: ledger instance cannot be nil")
    }
    if nodeID == "" {
        return 0, fmt.Errorf("retrieveNodeHealthScore: nodeID cannot be empty")
    }

    // Fetch health score from the ledger
    score, err := ledgerInstance.EnvironmentSystemCoreLedger.GetNodeHealthScore(nodeID)
    if err != nil {
        return 0, fmt.Errorf("retrieveNodeHealthScore: failed to retrieve health score for node %s: %w", nodeID, err)
    }

    // Log success
    log.Printf("Node %s health score retrieved: %d", nodeID, score)
    return score, nil
}


// setBlockMiningDifficulty updates the mining difficulty for the blockchain.
func setBlockMiningDifficulty(ledgerInstance *ledger.Ledger, difficultyLevel int) error {
    // Validate inputs
    if ledgerInstance == nil {
        return fmt.Errorf("setBlockMiningDifficulty: ledger instance cannot be nil")
    }
    if difficultyLevel <= 0 {
        return fmt.Errorf("setBlockMiningDifficulty: difficulty level must be greater than zero")
    }

    // Update mining difficulty in the ledger
    if err := ledgerInstance.EnvironmentSystemCoreLedger.UpdateMiningDifficulty(difficultyLevel); err != nil {
        return fmt.Errorf("setBlockMiningDifficulty: failed to set mining difficulty to %d: %w", difficultyLevel, err)
    }

    // Log success
    log.Printf("Mining difficulty successfully set to %d.", difficultyLevel)
    return nil
}

