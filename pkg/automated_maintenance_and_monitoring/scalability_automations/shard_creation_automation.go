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
	ShardCreationInterval       = 45 * time.Second // Interval for shard creation check
	ShardCreationLogEncryptionKey = "creationKey"  // Encryption key for shard creation logs
	MaxShardLoadThreshold       = 90               // Maximum load percentage before creating a new shard
)

// ShardCreationAutomation automates the process of creating new shards when load thresholds are exceeded.
type ShardCreationAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
	consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
	shardMutex      *sync.RWMutex                // Mutex for concurrency control
}

// NewShardCreationAutomation creates a new instance of ShardCreationAutomation.
func NewShardCreationAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardCreationAutomation {
	return &ShardCreationAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardCreation initiates the continuous monitoring for shard creation based on network load.
func (automation *ShardCreationAutomation) StartShardCreation() {
	ticker := time.NewTicker(ShardCreationInterval)
	go func() {
		for range ticker.C {
			automation.checkAndCreateShard()
		}
	}()
}

// checkAndCreateShard monitors the current network load and triggers shard creation if necessary.
func (automation *ShardCreationAutomation) checkAndCreateShard() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	networkLoad, err := automation.consensusSystem.GetNetworkLoad()
	if err != nil {
		fmt.Printf("Error fetching network load: %v\n", err)
		return
	}

	if automation.shouldCreateNewShard(networkLoad) {
		newShardID, err := automation.createShard()
		if err != nil {
			fmt.Printf("Error creating shard: %v\n", err)
			return
		}
		automation.logShardCreation(newShardID, networkLoad)
	}
}

// shouldCreateNewShard checks if a new shard should be created based on network load.
func (automation *ShardCreationAutomation) shouldCreateNewShard(loadPercentage int) bool {
	return loadPercentage >= MaxShardLoadThreshold
}

// createShard creates a new shard in the network.
func (automation *ShardCreationAutomation) createShard() (string, error) {
	newShardID, err := automation.consensusSystem.CreateNewShard()
	if err != nil {
		return "", fmt.Errorf("failed to create new shard: %v", err)
	}
	fmt.Printf("New shard created with ID: %s\n", newShardID)
	return newShardID, nil
}

// logShardCreation logs the creation of a new shard to the ledger.
func (automation *ShardCreationAutomation) logShardCreation(shardID string, loadPercentage int) {
	// Encrypt log before storing it
	logData := fmt.Sprintf("ShardID: %s, LoadPercentage: %d%%, Time: %s", shardID, loadPercentage, time.Now().String())
	encryptedLog, err := encryption.EncryptDataWithKey([]byte(logData), ShardCreationLogEncryptionKey)
	if err != nil {
		fmt.Printf("Error encrypting shard creation log: %v\n", err)
		return
	}

	// Log entry into ledger
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("SHARD-CREATION-%s-%d", shardID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Creation",
		Status:    "Completed",
		Details:   fmt.Sprintf("Encrypted Log: %s", string(encryptedLog)),
	}

	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log shard creation for shard %s: %v\n", shardID, err)
	}
}
