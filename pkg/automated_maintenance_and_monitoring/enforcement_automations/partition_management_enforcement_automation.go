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

// Configuration for partition management enforcement automation
const (
	PartitionCheckInterval        = 15 * time.Second // Interval to check partition status
	MaxLoadPerPartition           = 100000           // Maximum allowed load per partition before redistributing
	MinRequiredPartitions         = 3                // Minimum required partitions for optimal performance
	OverloadedPartitionThreshold  = 80000            // Threshold to flag partitions as overloaded
)

// PartitionManagementEnforcementAutomation monitors and enforces efficient partition management
type PartitionManagementEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	partitionLoadMap     map[string]int // Tracks load for each partition
	overloadedPartitions map[string]int // Tracks overloaded partitions and their load
}

// NewPartitionManagementEnforcementAutomation initializes the partition management enforcement automation
func NewPartitionManagementEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *PartitionManagementEnforcementAutomation {
	return &PartitionManagementEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		partitionLoadMap:     make(map[string]int),
		overloadedPartitions: make(map[string]int),
	}
}

// StartPartitionManagementEnforcement begins continuous monitoring and enforcement of partition management
func (automation *PartitionManagementEnforcementAutomation) StartPartitionManagementEnforcement() {
	ticker := time.NewTicker(PartitionCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkPartitionManagement()
		}
	}()
}

// checkPartitionManagement monitors each partition’s load and enforces load balancing when necessary
func (automation *PartitionManagementEnforcementAutomation) checkPartitionManagement() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluatePartitionLoad()
	automation.balanceOverloadedPartitions()
	automation.ensurePartitionCount()
}

// evaluatePartitionLoad updates the load status of each partition and flags overloaded partitions
func (automation *PartitionManagementEnforcementAutomation) evaluatePartitionLoad() {
	for _, partitionID := range automation.networkManager.GetAllPartitions() {
		load := automation.networkManager.GetPartitionLoad(partitionID)
		automation.partitionLoadMap[partitionID] = load

		if load > OverloadedPartitionThreshold {
			automation.overloadedPartitions[partitionID] = load
		} else {
			delete(automation.overloadedPartitions, partitionID) // Remove if no longer overloaded
		}
	}
}

// balanceOverloadedPartitions redistributes load from overloaded partitions to maintain performance
func (automation *PartitionManagementEnforcementAutomation) balanceOverloadedPartitions() {
	for partitionID, load := range automation.overloadedPartitions {
		if load > MaxLoadPerPartition {
			fmt.Printf("Partition management enforcement triggered for partition %s due to excessive load.\n", partitionID)
			automation.redistributePartitionLoad(partitionID, load)
		}
	}
}

// redistributePartitionLoad redistributes data from overloaded partitions to improve performance
func (automation *PartitionManagementEnforcementAutomation) redistributePartitionLoad(partitionID string, load int) {
	err := automation.networkManager.RedistributePartitionLoad(partitionID)
	if err != nil {
		fmt.Printf("Failed to redistribute load for partition %s: %v\n", partitionID, err)
		automation.logPartitionAction(partitionID, "Redistribution Failed", fmt.Sprintf("Load: %d", load))
	} else {
		fmt.Printf("Load successfully redistributed from partition %s.\n", partitionID)
		automation.logPartitionAction(partitionID, "Redistributed", fmt.Sprintf("Load: %d", load))
	}
}

// ensurePartitionCount checks if the network has a minimum number of partitions and activates additional ones if necessary
func (automation *PartitionManagementEnforcementAutomation) ensurePartitionCount() {
	currentPartitionCount := len(automation.networkManager.GetAllPartitions())
	if currentPartitionCount < MinRequiredPartitions {
		fmt.Println("Partition management enforcement triggered due to insufficient partitions.")
		automation.activateAdditionalPartitions(MinRequiredPartitions - currentPartitionCount)
	}
}

// activateAdditionalPartitions adds more partitions to maintain network stability and performance
func (automation *PartitionManagementEnforcementAutomation) activateAdditionalPartitions(count int) {
	err := automation.networkManager.ActivateAdditionalPartitions(count)
	if err != nil {
		fmt.Println("Failed to activate additional partitions:", err)
		automation.logPartitionAction("Network", "Failed to Activate Additional Partitions", "Partition Increase Required")
	} else {
		fmt.Printf("%d additional partitions activated to maintain network stability.\n", count)
		automation.logPartitionAction("Network", "Additional Partitions Activated", fmt.Sprintf("Activated Count: %d", count))
	}
}

// logPartitionAction securely logs actions related to partition management enforcement
func (automation *PartitionManagementEnforcementAutomation) logPartitionAction(partitionID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Partition ID: %s, Details: %s", action, partitionID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("partition-management-enforcement-%s-%d", partitionID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Partition Management Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log partition management enforcement action for partition %s: %v\n", partitionID, err)
	} else {
		fmt.Println("Partition management enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *PartitionManagementEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualPartitionAdjustment allows administrators to manually trigger partition adjustments
func (automation *PartitionManagementEnforcementAutomation) TriggerManualPartitionAdjustment() {
	fmt.Println("Manually triggering partition load redistribution for performance improvement.")

	for partitionID, load := range automation.partitionLoadMap {
		if load > OverloadedPartitionThreshold {
			automation.redistributePartitionLoad(partitionID, load)
		}
	}
	automation.logPartitionAction("Network", "Manual Partition Redistribution Triggered", "Administrator Action")
}
