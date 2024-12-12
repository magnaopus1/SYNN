package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/resource_manager"
)

const (
	ResourceCheckInterval        = 5 * time.Minute // Interval for checking resource allocation across nodes
	ResourceOverloadThreshold    = 0.85            // Threshold for resource overload (e.g., 85% CPU/memory utilization)
	ResourceAllocationLedgerType = "Resource Allocation Event"
)

// ResourceAllocationExecutionAutomation monitors and allocates resources dynamically
type ResourceAllocationExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validating changes
	ledgerInstance    *ledger.Ledger                        // Ledger instance for logging resource events
	resourceManager   *resource_manager.Manager             // Manager for handling resource allocation
	allocationMutex   *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewResourceAllocationExecutionAutomation initializes resource allocation automation
func NewResourceAllocationExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, resourceManager *resource_manager.Manager, allocationMutex *sync.RWMutex) *ResourceAllocationExecutionAutomation {
	return &ResourceAllocationExecutionAutomation{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		resourceManager:   resourceManager,
		allocationMutex:   allocationMutex,
	}
}

// StartResourceAllocationMonitor begins the monitoring process for resource allocation
func (automation *ResourceAllocationExecutionAutomation) StartResourceAllocationMonitor() {
	ticker := time.NewTicker(ResourceCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAndAllocateResources()
		}
	}()
}

// checkAndAllocateResources checks resource usage and reallocates based on current demand and availability
func (automation *ResourceAllocationExecutionAutomation) checkAndAllocateResources() {
	automation.allocationMutex.Lock()
	defer automation.allocationMutex.Unlock()

	// Fetch resource usage data for all nodes
	nodeResourceData := automation.resourceManager.GetAllNodeResourceUsage()

	for _, node := range nodeResourceData {
		automation.evaluateAndAllocate(node)
	}
}

// evaluateAndAllocate evaluates if resource reallocation is necessary based on node load
func (automation *ResourceAllocationExecutionAutomation) evaluateAndAllocate(node *resource_manager.NodeResourceUsage) {
	if node.CPUUsage >= ResourceOverloadThreshold || node.MemoryUsage >= ResourceOverloadThreshold {
		fmt.Printf("Node %s is overloaded (CPU: %.2f%%, Memory: %.2f%%). Reallocating resources.\n", node.ID, node.CPUUsage*100, node.MemoryUsage*100)
		automation.reallocateResources(node)
	} else {
		automation.logNormalResourceUsage(node)
	}
}

// reallocateResources dynamically reallocates resources from low-demand to high-demand nodes
func (automation *ResourceAllocationExecutionAutomation) reallocateResources(overloadedNode *resource_manager.NodeResourceUsage) {
	// Fetch nodes with underutilized resources for potential reallocation
	underloadedNodes := automation.resourceManager.GetUnderloadedNodes()

	for _, underloadedNode := range underloadedNodes {
		// Reallocate resources between nodes
		err := automation.resourceManager.ReallocateResources(underloadedNode.ID, overloadedNode.ID)
		if err != nil {
			fmt.Printf("Failed to reallocate resources from node %s to node %s: %v\n", underloadedNode.ID, overloadedNode.ID, err)
			continue
		}
		fmt.Printf("Successfully reallocated resources from node %s to node %s.\n", underloadedNode.ID, overloadedNode.ID)
		automation.logResourceReallocation(underloadedNode, overloadedNode)
	}
}

// logResourceReallocation logs resource reallocation events in the ledger securely
func (automation *ResourceAllocationExecutionAutomation) logResourceReallocation(fromNode, toNode *resource_manager.NodeResourceUsage) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("resource-reallocation-%s-%s-%d", fromNode.ID, toNode.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      ResourceAllocationLedgerType,
		Status:    "Success",
		Details:   fmt.Sprintf("Reallocated resources from node %s to node %s (CPU: %.2f%%, Memory: %.2f%%).", fromNode.ID, toNode.ID, toNode.CPUUsage*100, toNode.MemoryUsage*100),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log resource reallocation from node %s to node %s: %v\n", fromNode.ID, toNode.ID, err)
	} else {
		fmt.Println("Resource reallocation logged successfully.")
	}
}

// logNormalResourceUsage logs normal resource usage events in the ledger
func (automation *ResourceAllocationExecutionAutomation) logNormalResourceUsage(node *resource_manager.NodeResourceUsage) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("resource-usage-%s-%d", node.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      ResourceAllocationLedgerType,
		Status:    "Normal",
		Details:   fmt.Sprintf("Node %s is operating within normal resource limits (CPU: %.2f%%, Memory: %.2f%%).", node.ID, node.CPUUsage*100, node.MemoryUsage*100),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log normal resource usage for node %s: %v\n", node.ID, err)
	} else {
		fmt.Println("Normal resource usage logged successfully.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ResourceAllocationExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualResourceCheck allows administrators to manually trigger a resource allocation check for a specific node
func (automation *ResourceAllocationExecutionAutomation) TriggerManualResourceCheck(nodeID string) {
	fmt.Printf("Manually triggering resource allocation check for node %s...\n", nodeID)

	node := automation.resourceManager.GetNodeByID(nodeID)
	if node != nil {
		automation.evaluateAndAllocate(node)
	} else {
		fmt.Printf("Node %s not found.\n", nodeID)
	}
}
