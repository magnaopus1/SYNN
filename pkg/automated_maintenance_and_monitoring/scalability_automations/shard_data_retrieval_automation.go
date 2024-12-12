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
	DataRetrievalInterval       = 30 * time.Second // Interval for retrieving shard data
	ShardRetrievalLogEncryptionKey = "retrievalKey" // Encryption key for shard data retrieval logs
)

// ShardDataRetrievalAutomation handles automated retrieval of data from various shards.
type ShardDataRetrievalAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
	consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
	shardMutex      *sync.RWMutex                // Mutex for concurrency control
}

// NewShardDataRetrievalAutomation creates a new instance of ShardDataRetrievalAutomation.
func NewShardDataRetrievalAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardDataRetrievalAutomation {
	return &ShardDataRetrievalAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardDataRetrieval begins continuous monitoring and data retrieval from shards at regular intervals.
func (automation *ShardDataRetrievalAutomation) StartShardDataRetrieval() {
	ticker := time.NewTicker(DataRetrievalInterval)
	go func() {
		for range ticker.C {
			automation.retrieveDataFromShards()
		}
	}()
}

// retrieveDataFromShards fetches the necessary data from each shard and integrates it into the ledger.
func (automation *ShardDataRetrievalAutomation) retrieveDataFromShards() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	shardList, err := automation.consensusSystem.GetAllShards()
	if err != nil {
		fmt.Printf("Error fetching shards: %v\n", err)
		return
	}

	for _, shard := range shardList {
		data, err := automation.consensusSystem.GetShardData(shard.ID)
		if err != nil {
			fmt.Printf("Error retrieving data from shard %s: %v\n", shard.ID, err)
			continue
		}

		// Log the data retrieval process into the ledger
		automation.logShardDataRetrieval(shard.ID, data)
	}
}

// logShardDataRetrieval logs the data retrieval operation for a specific shard to the ledger.
func (automation *ShardDataRetrievalAutomation) logShardDataRetrieval(shardID string, data []byte) {
	// Encrypt the shard data before logging
	encryptedData, err := encryption.EncryptDataWithKey(data, ShardRetrievalLogEncryptionKey)
	if err != nil {
		fmt.Printf("Error encrypting shard data for shard %s: %v\n", shardID, err)
		return
	}

	// Create a ledger entry for the retrieval
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("SHARD-RETRIEVAL-%s-%d", shardID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Data Retrieval",
		Status:    "Completed",
		Details:   fmt.Sprintf("Encrypted Shard Data: %s", string(encryptedData)),
	}

	// Add the entry to the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log data retrieval for shard %s: %v\n", shardID, err)
	}
}
