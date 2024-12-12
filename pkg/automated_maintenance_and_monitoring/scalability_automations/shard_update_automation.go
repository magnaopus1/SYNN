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
	ShardUpdateInterval        = 5 * time.Minute  // Interval for checking and updating shard states
	MaxShardUpdateBatch        = 10               // Max number of shard updates in a batch
	ShardUpdateLogEncryptionKey = "shardUpdateKey" // Encryption key for shard update logs
)

// ShardUpdateAutomation handles the continuous update of shard states within the network.
type ShardUpdateAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging shard updates
	consensusSystem *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
	shardMutex      *sync.RWMutex                // Mutex for concurrency management
}

// NewShardUpdateAutomation creates a new instance of ShardUpdateAutomation.
func NewShardUpdateAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardUpdateAutomation {
	return &ShardUpdateAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardUpdateAutomation starts the automated shard update process in a continuous loop.
func (automation *ShardUpdateAutomation) StartShardUpdateAutomation() {
	ticker := time.NewTicker(ShardUpdateInterval)
	go func() {
		for range ticker.C {
			automation.updateShards()
		}
	}()
}

// updateShards retrieves shard updates, validates them, and applies the updates.
func (automation *ShardUpdateAutomation) updateShards() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	// Fetch a list of shards needing updates from the consensus system
	shardsToUpdate, err := automation.consensusSystem.GetShardsNeedingUpdate(MaxShardUpdateBatch)
	if err != nil {
		fmt.Printf("Error fetching shards for update: %v\n", err)
		return
	}

	// Apply updates to each shard and log the results
	for _, shard := range shardsToUpdate {
		err := automation.applyShardUpdate(shard)
		if err != nil {
			fmt.Printf("Error updating shard %s: %v\n", shard.ID, err)
			continue
		}
		automation.logShardUpdate(shard)
	}
}

// applyShardUpdate applies the necessary updates to the shard and integrates with the ledger.
func (automation *ShardUpdateAutomation) applyShardUpdate(shard common.Shard) error {
	// Fetch the latest state of the shard
	shardState, err := automation.consensusSystem.GetShardState(shard.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch shard state: %w", err)
	}

	// Validate the shard state before applying the update
	valid, validationErr := automation.consensusSystem.ValidateShardState(shard.ID, shardState)
	if !valid || validationErr != nil {
		return fmt.Errorf("shard state validation failed for shard %s: %v", shard.ID, validationErr)
	}

	// Apply the shard update in the consensus system
	updateErr := automation.consensusSystem.UpdateShard(shard.ID, shardState)
	if updateErr != nil {
		return fmt.Errorf("failed to update shard %s: %v", shard.ID, updateErr)
	}

	return nil
}

// logShardUpdate securely logs the shard update in the ledger.
func (automation *ShardUpdateAutomation) logShardUpdate(shard common.Shard) {
	// Encrypt the shard update details for secure logging
	encryptedUpdateDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Shard %s updated successfully", shard.ID)), ShardUpdateLogEncryptionKey)
	if err != nil {
		fmt.Printf("Error encrypting shard update log for shard %s: %v\n", shard.ID, err)
		return
	}

	// Create a ledger entry for the shard update
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("SHARD-UPDATE-%s-%d", shard.ID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Update",
		Status:    "Success",
		Details:   string(encryptedUpdateDetails),
	}

	// Log the shard update in the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log shard update for shard %s: %v\n", shard.ID, err)
	}
}
