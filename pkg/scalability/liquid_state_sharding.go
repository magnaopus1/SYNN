package scalability

import (
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewLiquidStateSharding initializes the Liquid-State Sharding system
func NewLiquidStateSharding(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.LiquidStateSharding {
	return &common.LiquidStateSharding{
		LiquidShards:     make(map[string]*common.LiquidShard),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateShard creates a new shard and allocates it to specific chains
func (lss *common.LiquidStateSharding) CreateShard(shardID, shardType string, allocatedChains []string, shardState []byte) (*common.LiquidShard, error) {
	lss.mu.Lock()
	defer lss.mu.Unlock()

	if _, exists := lss.LiquidShards[shardID]; exists {
		return nil, fmt.Errorf("shard %s already exists", shardID)
	}

	// Encrypt the shard state before storing it
	encryptedState, err := lss.EncryptionService.EncryptData(shardState, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt shard state: %v", err)
	}

	// Create the shard
	shard := &common.LiquidShard{
		ShardID:         shardID,
		ShardType:       shardType,
		AllocatedChains: allocatedChains,
		ShardState:      encryptedState,
		LastAdjustment:  time.Now(),
		ReallocationTime: time.Now(),
	}

	lss.LiquidShards[shardID] = shard

	// Log the shard creation in the ledger
	err = lss.Ledger.RecordShardCreation(shardID, shardType, allocatedChains, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard creation: %v", err)
	}

	fmt.Printf("Shard %s of type %s created and allocated to chains: %v\n", shardID, shardType, allocatedChains)
	return shard, nil
}

// AdjustShard dynamically adjusts the state of a shard based on cross-chain operations
func (lss *common.LiquidStateSharding) AdjustShard(shardID string, newState []byte) error {
	lss.mu.Lock()
	defer lss.mu.Unlock()

	shard, exists := lss.LiquidShards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Encrypt the new shard state before updating
	encryptedState, err := lss.EncryptionService.EncryptData(newState, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt new shard state: %v", err)
	}

	// Update shard state and log adjustment
	shard.ShardState = encryptedState
	shard.LastAdjustment = time.Now()

	err = lss.Ledger.RecordShardAdjustment(shardID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard adjustment: %v", err)
	}

	fmt.Printf("Shard %s state adjusted with new data\n", shardID)
	return nil
}

// ReallocateShard reassigns a shard to different chains based on network load or cross-chain demands
func (lss *common.LiquidStateSharding) ReallocateShard(shardID string, newChains []string) error {
	lss.mu.Lock()
	defer lss.mu.Unlock()

	shard, exists := lss.LiquidShards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Update allocated chains and log reallocation
	shard.AllocatedChains = newChains
	shard.ReallocationTime = time.Now()

	err := lss.Ledger.RecordShardReallocation(shardID, newChains, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard reallocation: %v", err)
	}

	fmt.Printf("Shard %s reallocated to chains: %v\n", shardID, newChains)
	return nil
}

// RetrieveShardState retrieves the decrypted state of a shard
func (lss *common.LiquidStateSharding) RetrieveShardState(shardID string) ([]byte, error) {
	lss.mu.Lock()
	defer lss.mu.Unlock()

	shard, exists := lss.LiquidShards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %s not found", shardID)
	}

	// Decrypt the shard state before returning
	decryptedState, err := lss.EncryptionService.DecryptData(shard.ShardState, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt shard state: %v", err)
	}

	fmt.Printf("Shard %s state retrieved\n", shardID)
	return decryptedState, nil
}

// LogShardActivity logs any activity related to a shard for transparency and auditability
func (lss *common.LiquidStateSharding) LogShardActivity(shardID, activityType string) error {
	lss.mu.Lock()
	defer lss.mu.Unlock()

	// Log the activity in the ledger
	err := lss.Ledger.RecordShardActivity(shardID, activityType, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard activity: %v", err)
	}

	fmt.Printf("Activity '%s' logged for shard %s\n", activityType, shardID)
	return nil
}
