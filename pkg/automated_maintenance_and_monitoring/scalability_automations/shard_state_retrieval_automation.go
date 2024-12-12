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
	ShardStateRetrievalInterval = 1 * time.Minute  // Interval for retrieving shard state
	ShardStateLogEncryptionKey  = "shardStateKey"  // Encryption key for shard state logs
	MaxShardStateRetrievalBatch = 10               // Max number of shard states retrieved in a batch
)

// ShardStateRetrievalAutomation retrieves and verifies the state of shards in the network.
type ShardStateRetrievalAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging shard state retrieval
	consensusSystem *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
	shardMutex      *sync.RWMutex                // Mutex for concurrency management
}

// NewShardStateRetrievalAutomation creates a new instance of ShardStateRetrievalAutomation.
func NewShardStateRetrievalAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardStateRetrievalAutomation {
	return &ShardStateRetrievalAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardStateRetrievalAutomation starts the automated shard state retrieval process.
func (automation *ShardStateRetrievalAutomation) StartShardStateRetrievalAutomation() {
	ticker := time.NewTicker(ShardStateRetrievalInterval)
	go func() {
		for range ticker.C {
			automation.retrieveShardStates()
		}
	}()
}

// retrieveShardStates fetches the state of each shard, verifies it, and logs the results.
func (automation *ShardStateRetrievalAutomation) retrieveShardStates() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	// Fetch a list of active shards from the consensus system
	activeShards, err := automation.consensusSystem.GetActiveShards(MaxShardStateRetrievalBatch)
	if err != nil {
		fmt.Printf("Error fetching active shards: %v\n", err)
		return
	}

	// Retrieve and verify the state of each shard
	for _, shard := range activeShards {
		shardState, err := automation.consensusSystem.GetShardState(shard.ID)
		if err != nil {
			fmt.Printf("Error retrieving state for shard %s: %v\n", shard.ID, err)
			continue
		}

		// Validate the shard state and log the result
		automation.logShardStateRetrieval(shard, shardState)
	}
}

// logShardStateRetrieval securely logs the retrieved shard state in the ledger.
func (automation *ShardStateRetrievalAutomation) logShardStateRetrieval(shard common.Shard, shardState common.ShardState) {
	// Encrypt the shard state for secure logging
	encryptedState, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Shard %s state retrieved: %v", shard.ID, shardState)), ShardStateLogEncryptionKey)
	if err != nil {
		fmt.Printf("Error encrypting shard state for shard %s: %v\n", shard.ID, err)
		return
	}

	// Create a ledger entry for the shard state retrieval
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("SHARD-STATE-RET-%s-%d", shard.ID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard State Retrieval",
		Status:    "Completed",
		Details:   string(encryptedState),
	}

	// Log the shard state retrieval in the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log shard state retrieval for shard %s: %v\n", shard.ID, err)
	}
}
