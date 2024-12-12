package scalability

import (

	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewDelegatedAdaptiveLayerSharding initializes the sharding system
func NewDelegatedAdaptiveLayerSharding(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.DelegatedAdaptiveLayerSharding {
	return &common.DelegatedAdaptiveLayerSharding{
		Shards:            make(map[string]*common.Shard),
		NodeShards:        make(map[string][]string),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// CreateShard creates a new shard and assigns nodes for management
func (dals *common.DelegatedAdaptiveLayerSharding) CreateShard(shardID, layerID string, assignedNodes []string, stateData []byte) (*common.Shard, error) {
	dals.mu.Lock()
	defer dals.mu.Unlock()

	if _, exists := dals.Shards[shardID]; exists {
		return nil, fmt.Errorf("shard %s already exists", shardID)
	}

	// Encrypt the state data before storing it
	encryptedState, err := dals.EncryptionService.EncryptData(stateData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt state data: %v", err)
	}

	// Create the shard
	shard := &common.Shard{
		ShardID:       shardID,
		LayerID:       layerID,
		AssignedNodes: assignedNodes,
		StateData:     encryptedState,
		LastUpdate:    time.Now(),
		ShardSize:     len(stateData),
	}
	dals.Shards[shardID] = shard

	// Assign shard to nodes
	for _, node := range assignedNodes {
		dals.NodeShards[node] = append(dals.NodeShards[node], shardID)
	}

	// Log the shard creation in the ledger
	err = dals.Ledger.RecordShardCreation(shardID, layerID, assignedNodes, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log shard creation: %v", err)
	}

	fmt.Printf("Shard %s created for layer %s, assigned to nodes: %v\n", shardID, layerID, assignedNodes)
	return shard, nil
}

// UpdateShard updates the state data of an existing shard
func (dals *common.DelegatedAdaptiveLayerSharding) UpdateShard(shardID string, newStateData []byte) error {
	dals.mu.Lock()
	defer dals.mu.Unlock()

	shard, exists := dals.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Encrypt the new state data before updating
	encryptedState, err := dals.EncryptionService.EncryptData(newStateData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt new state data: %v", err)
	}

	// Update the shard state
	shard.StateData = encryptedState
	shard.LastUpdate = time.Now()
	shard.ShardSize = len(newStateData)

	// Log the shard update in the ledger
	err = dals.Ledger.RecordShardUpdate(shardID, shard.ShardSize, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard update: %v", err)
	}

	fmt.Printf("Shard %s updated with new state data (size: %d bytes)\n", shardID, len(newStateData))
	return nil
}

// DelegateShard dynamically assigns new nodes to an existing shard based on load
func (dals *common.DelegatedAdaptiveLayerSharding) DelegateShard(shardID string, newNodes []string) error {
	dals.mu.Lock()
	defer dals.mu.Unlock()

	shard, exists := dals.Shards[shardID]
	if !exists {
		return fmt.Errorf("shard %s not found", shardID)
	}

	// Assign new nodes to the shard
	shard.AssignedNodes = append(shard.AssignedNodes, newNodes...)
	for _, node := range newNodes {
		dals.NodeShards[node] = append(dals.NodeShards[node], shardID)
	}

	// Log the delegation in the ledger
	err := dals.Ledger.RecordShardDelegation(shardID, newNodes, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard delegation: %v", err)
	}

	fmt.Printf("Shard %s delegated to new nodes: %v\n", shardID, newNodes)
	return nil
}

// RetrieveShardData retrieves the state data of a shard (decrypted)
func (dals *common.DelegatedAdaptiveLayerSharding) RetrieveShardData(shardID string) ([]byte, error) {
	dals.mu.Lock()
	defer dals.mu.Unlock()

	shard, exists := dals.Shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %s not found", shardID)
	}

	// Decrypt the state data before returning
	decryptedData, err := dals.EncryptionService.DecryptData(shard.StateData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt shard data: %v", err)
	}

	fmt.Printf("Shard %s data retrieved (size: %d bytes)\n", shardID, len(decryptedData))
	return decryptedData, nil
}

// LogShardActivity logs any activity related to a shard for transparency and auditability
func (dals *common.DelegatedAdaptiveLayerSharding) LogShardActivity(shardID string, activityType string) error {
	dals.mu.Lock()
	defer dals.mu.Unlock()

	// Log the activity in the ledger
	err := dals.Ledger.RecordShardActivity(shardID, activityType, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shard activity: %v", err)
	}

	fmt.Printf("Activity '%s' logged for shard %s\n", activityType, shardID)
	return nil
}
