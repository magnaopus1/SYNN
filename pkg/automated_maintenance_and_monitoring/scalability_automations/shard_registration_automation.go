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
	ShardRegistrationInterval      = 30 * time.Second  // Interval to check for shard registration events
	ShardRegistrationLogEncryption = "shardRegKey"     // Encryption key for shard registration logs
	MaxShardsPerRegistrationBatch  = 5                 // Max number of shards registered in a batch
)

// ShardRegistrationAutomation manages the registration and initialization of new shards in the network.
type ShardRegistrationAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging shard registration events
	consensusSystem *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
	shardMutex      *sync.RWMutex                // Mutex for managing concurrency when registering shards
}

// NewShardRegistrationAutomation creates a new instance of ShardRegistrationAutomation.
func NewShardRegistrationAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardRegistrationAutomation {
	return &ShardRegistrationAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardRegistrationAutomation starts a loop that monitors and handles new shard registrations.
func (automation *ShardRegistrationAutomation) StartShardRegistrationAutomation() {
	ticker := time.NewTicker(ShardRegistrationInterval)
	go func() {
		for range ticker.C {
			automation.registerNewShards()
		}
	}()
}

// registerNewShards handles the logic for registering new shards when necessary.
func (automation *ShardRegistrationAutomation) registerNewShards() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	// Fetch any pending shards that need to be registered
	pendingShards, err := automation.consensusSystem.GetPendingShards(MaxShardsPerRegistrationBatch)
	if err != nil {
		fmt.Printf("Error fetching pending shards: %v\n", err)
		return
	}

	for _, shard := range pendingShards {
		// Initialize shard in the consensus system
		if err := automation.consensusSystem.InitializeShard(shard); err != nil {
			fmt.Printf("Error initializing shard %s: %v\n", shard.ID, err)
			continue
		}

		// Register the shard in the ledger
		automation.logShardRegistration(shard)
	}
}

// logShardRegistration securely logs the registration of a new shard in the ledger.
func (automation *ShardRegistrationAutomation) logShardRegistration(shard common.Shard) {
	// Encrypt the shard registration details for secure logging
	encryptedDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Shard %s successfully registered at %s", shard.ID, time.Now().String())), ShardRegistrationLogEncryption)
	if err != nil {
		fmt.Printf("Error encrypting shard registration details for shard %s: %v\n", shard.ID, err)
		return
	}

	// Create a ledger entry for the shard registration
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("SHARD-REG-%s-%d", shard.ID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Registration",
		Status:    "Completed",
		Details:   string(encryptedDetails),
	}

	// Log the shard registration in the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log shard registration for shard %s: %v\n", shard.ID, err)
	}
}
