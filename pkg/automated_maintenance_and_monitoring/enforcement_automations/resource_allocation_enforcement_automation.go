package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
)

// Configuration for resource allocation enforcement automation
const (
	ResourceCheckInterval       = 15 * time.Second // Interval to check resource allocation
	MaxCPUUsage                 = 85              // Max CPU usage threshold
	MaxMemoryUsage              = 85              // Max memory usage threshold
	MaxBandwidthUsage           = 90              // Max bandwidth usage threshold
	ResourceViolationThreshold  = 3               // Allowed violations before adjusting resources
)

// ResourceAllocationEnforcementAutomation monitors and enforces resource allocation requirements
type ResourceAllocationEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	resourceViolationMap map[string]int // Tracks violations for each node based on resource usage
	nodeResourceMap      map[string]map[string]int // Tracks resource usage stats for each node (CPU, Memory, Bandwidth)
}

// NewResourceAllocationEnforcementAutomation initializes the resource allocation enforcement automation
func NewResourceAllocationEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ResourceAllocationEnforcementAutomation {
	return &ResourceAllocationEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		resourceViolationMap: make(map[string]int),
		nodeResourceMap:      make(map[string]map[string]int),
	}
}

// StartResourceAllocationEnforcement begins continuous monitoring and enforcement of resource allocation
func (automation *ResourceAllocationEnforcementAutomation) StartResourceAllocationEnforcement() {
	ticker := time.NewTicker(ResourceCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkResourceCompliance()
		}
	}()
}

// checkResourceCompliance monitors each node's resource usage and adjusts allocation as necessary
func (automation *ResourceAllocationEnforcementAutomation) checkResourceCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateResourceUsage()
	automation.adjustResourceAllocation()
}

// evaluateResourceUsage checks each node’s CPU, memory, and bandwidth usage and flags non-compliance
func (automation *ResourceAllocationEnforcementAutomation) evaluateResourceUsage() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		cpuUsage := automation.networkManager.GetNodeCPUUsage(nodeID)
		memoryUsage := automation.networkManager.GetNodeMemoryUsage(nodeID)
		bandwidthUsage := automation.networkManager.GetNodeBandwidthUsage(nodeID)

		automation.nodeResourceMap[nodeID] = map[string]int{
			"CPU":       cpuUsage,
			"Memory":    memoryUsage,
			"Bandwidth": bandwidthUsage,
		}

		// Track violations if resource usage exceeds thresholds
		if cpuUsage > MaxCPUUsage || memoryUsage > MaxMemoryUsage || bandwidthUsage > MaxBandwidthUsage {
			automation.resourceViolationMap[nodeID]++
		} else {
			automation.resourceViolationMap[nodeID] = 0 // Reset if compliant
		}
	}
}

// adjustResourceAllocation redistributes resources for nodes that exceed allowable usage thresholds
func (automation *ResourceAllocationEnforcementAutomation) adjustResourceAllocation() {
	for nodeID, violations := range automation.resourceViolationMap {
		if violations >= ResourceViolationThreshold {
			fmt.Printf("Resource allocation enforcement triggered for node %s due to excessive resource usage.\n", nodeID)
			automation.applyResourceAdjustment(nodeID)
		}
	}
}

// applyResourceAdjustment adjusts resource allocation for nodes that exceed usage thresholds
func (automation *ResourceAllocationEnforcementAutomation) applyResourceAdjustment(nodeID string) {
	err := automation.networkManager.AdjustNodeResources(nodeID)
	if err != nil {
		fmt.Printf("Failed to adjust resources for node %s: %v\n", nodeID, err)
		automation.logResourceAllocationAction(nodeID, "Adjustment Failed", fmt.Sprintf("Resources Exceeded: CPU=%d, Memory=%d, Bandwidth=%d",
			automation.nodeResourceMap[nodeID]["CPU"],
			automation.nodeResourceMap[nodeID]["Memory"],
			automation.nodeResourceMap[nodeID]["Bandwidth"]))
	} else {
		fmt.Printf("Resources successfully adjusted for node %s.\n", nodeID)
		automation.logResourceAllocationAction(nodeID, "Resources Adjusted", fmt.Sprintf("Adjusted Resources: CPU=%d, Memory=%d, Bandwidth=%d",
			automation.nodeResourceMap[nodeID]["CPU"],
			automation.nodeResourceMap[nodeID]["Memory"],
			automation.nodeResourceMap[nodeID]["Bandwidth"]))
	}
}

// logResourceAllocationAction securely logs actions related to resource allocation enforcement
func (automation *ResourceAllocationEnforcementAutomation) logResourceAllocationAction(nodeID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Details: %s", action, nodeID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("resource-allocation-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Resource Allocation Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log resource allocation enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Resource allocation enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ResourceAllocationEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualResourceAdjustment allows administrators to manually adjust resources for nodes with high usage
func (automation *ResourceAllocationEnforcementAutomation) TriggerManualResourceAdjustment(nodeID string) {
	fmt.Printf("Manually adjusting resources for node: %s\n", nodeID)

	if automation.nodeResourceMap[nodeID]["CPU"] > MaxCPUUsage || automation.nodeResourceMap[nodeID]["Memory"] > MaxMemoryUsage || automation.nodeResourceMap[nodeID]["Bandwidth"] > MaxBandwidthUsage {
		automation.applyResourceAdjustment(nodeID)
	} else {
		fmt.Printf("Node %s is within acceptable resource usage limits.\n", nodeID)
		automation.logResourceAllocationAction(nodeID, "Manual Adjustment Skipped", "Resources Within Limits")
	}
}
