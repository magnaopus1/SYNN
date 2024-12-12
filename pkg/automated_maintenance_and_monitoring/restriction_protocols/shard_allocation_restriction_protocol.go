package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	ShardAllocationCheckInterval = 5 * time.Second // Interval to check shard allocation
	MaxShardCapacityThreshold    = 80              // Maximum shard capacity threshold in percentage
)

// ShardAllocationRestrictionAutomation manages shard allocation and enforces restrictions when shard over-utilization is detected
type ShardAllocationRestrictionAutomation struct {
	consensusSystem  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	stateMutex       *sync.RWMutex
	shardAllocations map[string]int // Tracks current shard allocation utilization
}

// NewShardAllocationRestrictionAutomation initializes the ShardAllocationRestrictionAutomation struct
func NewShardAllocationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ShardAllocationRestrictionAutomation {
	return &ShardAllocationRestrictionAutomation{
		consensusSystem:  consensusSystem,
		ledgerInstance:   ledgerInstance,
		stateMutex:       stateMutex,
		shardAllocations: make(map[string]int),
	}
}

// StartShardAllocationMonitoring starts the continuous process of monitoring shard allocation
func (automation *ShardAllocationRestrictionAutomation) StartShardAllocationMonitoring() {
	ticker := time.NewTicker(ShardAllocationCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateShardAllocations()
		}
	}()
}

// evaluateShardAllocations checks the shard allocation levels and triggers restrictions if necessary
func (automation *ShardAllocationRestrictionAutomation) evaluateShardAllocations() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch current shard allocations from the consensus system
	shardData := automation.consensusSystem.GetShardAllocationData()

	for shardID, usage := range shardData {
		automation.logShardAllocation(shardID, usage)

		// If the shard usage exceeds the capacity threshold, restrict further allocation
		if usage >= MaxShardCapacityThreshold {
			automation.restrictShardAllocation(shardID)
		}
	}
}

// logShardAllocation logs the shard allocation details in the ledger
func (automation *ShardAllocationRestrictionAutomation) logShardAllocation(shardID string, usage int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-allocation-%s-%d", shardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Allocation Monitoring",
		Status:    "In Progress",
		Details:   fmt.Sprintf("Shard %s is at %d%% capacity.", shardID, usage),
	}

	// Encrypt the shard allocation details before logging them in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	// Add the shard allocation entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log shard allocation:", err)
	} else {
		fmt.Println("Shard allocation logged successfully for shard:", shardID)
	}
}

// restrictShardAllocation restricts further allocation to a shard if it exceeds the maximum capacity threshold
func (automation *ShardAllocationRestrictionAutomation) restrictShardAllocation(shardID string) {
	fmt.Printf("Shard %s has exceeded the allocation threshold and will be restricted.\n", shardID)

	// Log the shard allocation restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-allocation-restriction-%s-%d", shardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Allocation Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Shard %s is restricted from further allocation due to exceeding %d%% capacity.", shardID, MaxShardCapacityThreshold),
	}

	// Encrypt the restriction details before adding them to the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log shard allocation restriction:", err)
	} else {
		fmt.Println("Shard allocation restriction applied successfully for shard:", shardID)
	}

	// Inform the consensus system to restrict further transactions or allocations on this shard
	automation.consensusSystem.RestrictShard(shardID)
}

// encryptData encrypts sensitive shard allocation data before storing it in the ledger
func (automation *ShardAllocationRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting shard allocation details:", err)
		return data
	}
	return string(encryptedData)
}
