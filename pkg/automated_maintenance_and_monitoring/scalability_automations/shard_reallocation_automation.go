package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	ReallocationInterval          = 60 * time.Second  // Interval for checking shard reallocation
	ShardReallocationLogKey       = "shardReallocKey" // Encryption key for shard reallocation logs
	ReallocationThreshold         = 70               // Percentage threshold for shard load before reallocation
	MaxValidatorReallocationBatch = 20               // Maximum validators to reallocate per interval
)

// ShardReallocationAutomation handles the dynamic reallocation of validators between shards based on load.
type ShardReallocationAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
	consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
	shardMutex      *sync.RWMutex                // Mutex for concurrency control
}

// NewShardReallocationAutomation creates a new instance of ShardReallocationAutomation.
func NewShardReallocationAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardReallocationAutomation {
	return &ShardReallocationAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardReallocationAutomation starts the continuous process of monitoring shard load and reallocation.
func (automation *ShardReallocationAutomation) StartShardReallocationAutomation() {
	ticker := time.NewTicker(ReallocationInterval)
	go func() {
		for range ticker.C {
			automation.checkAndReallocateShards()
		}
	}()
}

// checkAndReallocateShards monitors the shard loads and reallocates validators when certain thresholds are reached.
func (automation *ShardReallocationAutomation) checkAndReallocateShards() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	shardList, err := automation.consensusSystem.GetAllShards()
	if err != nil {
		fmt.Printf("Error fetching shards: %v\n", err)
		return
	}

	for _, shard := range shardList {
		shardLoad, err := automation.consensusSystem.GetShardLoad(shard.ID)
		if err != nil {
			fmt.Printf("Error fetching load for shard %s: %v\n", shard.ID, err)
			continue
		}

		if shardLoad >= ReallocationThreshold {
			automation.reallocateValidatorsFromShard(shard.ID, shardLoad)
		}
	}
}

// reallocateValidatorsFromShard reallocates validators from an overloaded shard to underloaded shards.
func (automation *ShardReallocationAutomation) reallocateValidatorsFromShard(shardID string, shardLoad int) {
	validators, err := automation.consensusSystem.GetValidatorsForShard(shardID, MaxValidatorReallocationBatch)
	if err != nil {
		fmt.Printf("Error fetching validators from shard %s: %v\n", shardID, err)
		return
	}

	// Find underloaded shards for reallocation
	underloadedShards, err := automation.consensusSystem.GetUnderloadedShards()
	if err != nil {
		fmt.Printf("Error fetching underloaded shards: %v\n", err)
		return
	}

	if len(underloadedShards) == 0 {
		fmt.Printf("No underloaded shards available for reallocation\n")
		return
	}

	// Distribute validators to underloaded shards
	for i, validator := range validators {
		targetShard := underloadedShards[i%len(underloadedShards)]
		if err := automation.consensusSystem.AssignValidatorToShard(targetShard.ID, validator); err != nil {
			fmt.Printf("Error reallocating validator %s to shard %s: %v\n", validator, targetShard.ID, err)
			continue
		}

		// Log the reallocation
		automation.logValidatorReallocation(shardID, targetShard.ID, validator)
	}
}

// logValidatorReallocation logs the reallocation of validators from one shard to another in the ledger.
func (automation *ShardReallocationAutomation) logValidatorReallocation(fromShardID string, toShardID string, validatorID string) {
	// Encrypt the reallocation details
	encryptedDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Validator %s reallocated from Shard %s to Shard %s", validatorID, fromShardID, toShardID)), ShardReallocationLogKey)
	if err != nil {
		fmt.Printf("Error encrypting reallocation details for shard %s: %v\n", fromShardID, err)
		return
	}

	// Create a ledger entry for the validator reallocation
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("REALLOCATION-%s-%s-%s-%d", fromShardID, toShardID, validatorID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Reallocation",
		Status:    "Completed",
		Details:   string(encryptedDetails),
	}

	// Add the entry to the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log validator reallocation for shard %s: %v\n", fromShardID, err)
	}
}
