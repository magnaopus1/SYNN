package scalability

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewMetaLayerOrchestrator initializes the Meta-Layer Orchestration system
func NewMetaLayerOrchestrator(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.MetaLayerOrchestrator {
	return &common.MetaLayerOrchestrator{
		Shards:           make(map[string]*common.MetaLayerShard),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// RegisterShard registers a shard for orchestration
func (mlo *common.MetaLayerOrchestrator) RegisterShard(shardID, shardType string) (*common.MetaLayerShard, error) {
	mlo.mu.Lock()
	defer mlo.mu.Unlock()

	if _, exists := mlo.Shards[shardID]; exists {
		return nil, fmt.Errorf("shard %s is already registered", shardID)
	}

	// Create the shard
	shard := &common.MetaLayerShard{
		ShardID:      shardID,
		ShardType:    shardType,
		LastActivity: time.Now(),
	}

	mlo.Shards[shardID] = shard

	// Log shard registration in the ledger
	err := mlo.Ledger.RecordShardRegistration(shardID, shardType, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard registration: %v", err)
	}

	fmt.Printf("Shard %s of type %s registered for meta-layer orchestration\n", shardID, shardType)
	return shard, nil
}

// OrchestrateTransaction orchestrates a transaction between shards, handling cross-layer coordination
func (mlo *common.MetaLayerOrchestrator) OrchestrateTransaction(txID string, involvedShards []string, txData []byte) error {
	mlo.mu.Lock()
	defer mlo.mu.Unlock()

	if len(involvedShards) == 0 {
		return errors.New("no shards provided for orchestration")
	}

	// Encrypt the transaction data
	encryptedData, err := mlo.EncryptionService.EncryptData(txData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	// Process the transaction across the involved shards
	for _, shardID := range involvedShards {
		shard, exists := mlo.Shards[shardID]
		if !exists {
			return fmt.Errorf("shard %s is not registered for orchestration", shardID)
		}

		// Log orchestration in the ledger
		err = mlo.Ledger.RecordOrchestratedTransaction(txID, shardID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log orchestration for shard %s: %v", shardID, err)
		}

		shard.LastActivity = time.Now()
		fmt.Printf("Transaction %s orchestrated through shard %s\n", txID, shardID)
	}

	fmt.Printf("Transaction %s successfully orchestrated across shards: %v\n", txID, involvedShards)
	return nil
}

// AdjustOrchestration dynamically adjusts the orchestration paths between shards for performance optimization
func (mlo *common.MetaLayerOrchestrator) AdjustOrchestration(shardID string, newInvolvedShards []string) error {
	mlo.mu.Lock()
	defer mlo.mu.Unlock()

	shard, exists := mlo.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Adjust orchestration paths and log in the ledger
	err := mlo.Ledger.RecordOrchestrationAdjustment(shardID, newInvolvedShards, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log orchestration adjustment for shard %s: %v", shardID, err)
	}

	shard.LastActivity = time.Now()
	fmt.Printf("Orchestration paths adjusted for shard %s with new involved shards: %v\n", shardID, newInvolvedShards)
	return nil
}

// RetrieveOrchestrationLog retrieves the orchestration log for a specific shard
func (mlo *common.MetaLayerOrchestrator) RetrieveOrchestrationLog(shardID string) ([]ledger.OrchestrationLog, error) {
	mlo.mu.Lock()
	defer mlo.mu.Unlock()

	// Retrieve the orchestration logs from the ledger
	logs, err := mlo.Ledger.GetOrchestrationLogs(shardID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orchestration logs for shard %s: %v", shardID, err)
	}

	fmt.Printf("Retrieved orchestration logs for shard %s\n", shardID)
	return logs, nil
}
