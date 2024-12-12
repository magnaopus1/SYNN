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

// Configuration for shard merging enforcement automation
const (
	ShardMergingCheckInterval      = 30 * time.Second // Interval to check for shard merging requirements
	MinTransactionThreshold        = 25000            // Minimum transaction volume to avoid merging
	ShardMergingViolationThreshold = 3                // Allowed violations before triggering merging
)

// ShardMergingEnforcementAutomation monitors and enforces shard merging for underutilized shards
type ShardMergingEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	shardTransactionMap  map[string]int // Tracks transaction volume for each shard
	shardMergeViolationMap map[string]int // Tracks violation count for each shard eligible for merging
}

// NewShardMergingEnforcementAutomation initializes the shard merging enforcement automation
func NewShardMergingEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ShardMergingEnforcementAutomation {
	return &ShardMergingEnforcementAutomation{
		networkManager:         networkManager,
		consensusEngine:        consensusEngine,
		ledgerInstance:         ledgerInstance,
		enforcementMutex:       enforcementMutex,
		shardTransactionMap:    make(map[string]int),
		shardMergeViolationMap: make(map[string]int),
	}
}

// StartShardMergingEnforcement begins continuous monitoring and enforcement of shard merging
func (automation *ShardMergingEnforcementAutomation) StartShardMergingEnforcement() {
	ticker := time.NewTicker(ShardMergingCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkShardMergingRequirements()
		}
	}()
}

// checkShardMergingRequirements monitors each shard’s transaction volume and enforces merging where needed
func (automation *ShardMergingEnforcementAutomation) checkShardMergingRequirements() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.evaluateShardTransactionVolumes()
	automation.performShardMerging()
}

// evaluateShardTransactionVolumes checks each shard's transaction volume and flags underutilized shards
func (automation *ShardMergingEnforcementAutomation) evaluateShardTransactionVolumes() {
	for _, shardID := range automation.networkManager.GetAllShards() {
		transactionCount := automation.networkManager.GetShardTransactionVolume(shardID)
		automation.shardTransactionMap[shardID] = transactionCount

		// Track violations for shards under the minimum transaction threshold
		if transactionCount < MinTransactionThreshold {
			automation.shardMergeViolationMap[shardID]++
		} else {
			automation.shardMergeViolationMap[shardID] = 0 // Reset if within limits
		}
	}
}

// performShardMerging merges underutilized shards based on violation threshold
func (automation *ShardMergingEnforcementAutomation) performShardMerging() {
	for shardID, violations := range automation.shardMergeViolationMap {
		if violations >= ShardMergingViolationThreshold {
			fmt.Printf("Shard merging enforcement triggered for shard %s due to low transaction volume.\n", shardID)
			automation.mergeShard(shardID)
		}
	}
}

// mergeShard consolidates a low-volume shard to optimize network resources
func (automation *ShardMergingEnforcementAutomation) mergeShard(lowVolumeShardID string) {
	err := automation.networkManager.MergeShard(lowVolumeShardID)
	if err != nil {
		fmt.Printf("Failed to merge shard %s: %v\n", lowVolumeShardID, err)
		automation.logShardMergingAction(lowVolumeShardID, "Shard Merge Failed", fmt.Sprintf("Transaction Count: %d", automation.shardTransactionMap[lowVolumeShardID]))
	} else {
		fmt.Printf("Shard %s successfully merged to optimize resources.\n", lowVolumeShardID)
		automation.logShardMergingAction(lowVolumeShardID, "Shard Merged", fmt.Sprintf("Transaction Count: %d", automation.shardTransactionMap[lowVolumeShardID]))
	}
}

// logShardMergingAction securely logs actions related to shard merging enforcement
func (automation *ShardMergingEnforcementAutomation) logShardMergingAction(shardID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Shard ID: %s, Details: %s", action, shardID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-merging-enforcement-%s-%d", shardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Merging Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log shard merging enforcement action for shard %s: %v\n", shardID, err)
	} else {
		fmt.Println("Shard merging enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ShardMergingEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualShardMerge allows administrators to manually merge a shard if needed
func (automation *ShardMergingEnforcementAutomation) TriggerManualShardMerge(shardID string) {
	fmt.Printf("Manually triggering shard merge for shard: %s\n", shardID)

	transactionCount := automation.shardTransactionMap[shardID]
	if transactionCount < MinTransactionThreshold {
		automation.mergeShard(shardID)
	} else {
		fmt.Printf("Shard %s meets transaction volume requirements, merge not required.\n", shardID)
		automation.logShardMergingAction(shardID, "Manual Merge Skipped", "Sufficient Transaction Volume")
	}
}
