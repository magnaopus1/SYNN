package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordCompression compresses transaction records and stores them in the cache
func (l *ScalabilityLedger) RecordCompression(t *BlockchainConsensusCoinLedger, record TransactionRecord) error {
	l.Lock()
	defer l.Unlock()

	compressedRecord := compressTransactionRecord(record) // Placeholder for compression logic
	if compressedRecord == nil {
		return errors.New("failed to compress transaction record")
	}
	t.TransactionCache[record.ID] = *compressedRecord // Save compressed transaction to cache
	return nil
}

// compressTransactionRecord converts a TransactionRecord to a Transaction.
func compressTransactionRecord(record TransactionRecord) *Transaction {
    // Mapping fields from TransactionRecord to Transaction
    return &Transaction{
        TransactionID:    record.ID,           // Mapping ID to TransactionID
        FromAddress:      record.From,         // Mapping From to FromAddress
        ToAddress:        record.To,           // Mapping To to ToAddress
        Amount:           record.Amount,       // Mapping Amount
        Fee:              record.Fee,          // Mapping Fee
        Status:           record.Status,       // Mapping Status
        Timestamp:        record.Timestamp,    // Mapping Timestamp
        ValidatorID:      record.ValidatorID,   // Mapping ValidatorID
        BlockID:          "",                  // BlockID not available in TransactionRecord, so leave it empty
        SubBlockID:       "",                  // SubBlockID not available, so leave it empty
        Signature:        "",                  // Signature not present in TransactionRecord, so leave it empty
        TokenStandard:    "",                  // TokenStandard not available, set to empty
        TokenID:          "",                  // TokenID not available, set to empty
        EncryptedData:    "",                  // No encrypted data in TransactionRecord, leave it empty
        DecryptedData:    "",                  // No decrypted data in TransactionRecord, leave it empty
        ExecutionResult:  "",                  // Execution result not applicable
        FrozenAmount:     0,                   // No FrozenAmount in TransactionRecord, set to 0
        RefundAmount:     0,                   // No RefundAmount in TransactionRecord, set to 0
        ReversalRequested: false,              // No ReversalRequested flag in TransactionRecord, set to false
    }
}


// RecordShardCreation logs the creation of a new shard
func (l *ScalabilityLedger) RecordShardCreation(T *BlockchainConsensusCoinLedger, shardID string) error {
	l.Lock()
	defer l.Unlock()

	newShard := Shard{
		ShardID:   shardID,
		CreatedAt: time.Now(),
		Status:    "active",
	}
	l.Shards[shardID] = newShard // Add shard to ledger
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{ID: shardID, Action: "ShardCreated", Timestamp: time.Now()})
	return nil
}

// RecordShardUpdate updates the state of a shard
func (l *ScalabilityLedger) RecordShardUpdate(T *BlockchainConsensusCoinLedger,shardID string, newState string) error {
	l.Lock()
	defer l.Unlock()

	if shard, exists := l.Shards[shardID]; exists {
		shard.Status = newState
		l.Shards[shardID] = shard // Update shard state
		T.BlockchainConsensusCoinState.TransactionHistory = append( T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{ID: shardID, Action: "ShardUpdated: " + newState, Timestamp: time.Now()})
		return nil
	}
	return errors.New("shard not found")
}

// RecordShardDelegation logs the delegation of a shard
func (l *ScalabilityLedger) RecordShardDelegation(T *BlockchainConsensusCoinLedger,shardID, delegator string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.Shards[shardID]; !exists {
		return errors.New("shard not found")
	}
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{ID: shardID, Action: "ShardDelegated", Delegator: delegator, Timestamp: time.Now()})
	return nil
}

// RecordShardActivity logs activity on a shard.
func (l *ScalabilityLedger) RecordShardActivity(T *BlockchainConsensusCoinLedger, shardID string) error {
    l.Lock()
    defer l.Unlock()

    // Ensure the Shards map is initialized
    if l.Shards == nil {
        l.Shards = make(map[string]ShardRecord)
    }

    // Check if the shard exists
    if _, exists := l.Shards[shardID]; !exists {
        return errors.New("shard not found")
    }

    // Log the shard activity
    T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
        ID:        shardID,
        Action:    "ShardActivity",
        Details:   "Shard activity recorded",
        Timestamp: time.Now(),
    })

    return nil
}

// RecordShardReallocation reallocates a shard to a new node.
func (l *ScalabilityLedger) RecordShardReallocation(T *BlockchainConsensusCoinLedger, shardID, newNodeID string) error {
    l.Lock()
    defer l.Unlock()

    // Ensure the Shards map is initialized
    if l.Shards == nil {
        l.Shards = make(map[string]ShardRecord)
    }

    // Check if the shard exists
    if shard, exists := l.Shards[shardID]; exists {
        // Update the shard's NodeID
        shard.NodeID = newNodeID
        shard.Timestamp = time.Now()
        l.Shards[shardID] = shard

        // Log the reallocation
        T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
            ID:        shardID,
            Action:    "ShardReallocated",
            Details:   fmt.Sprintf("Reallocated to node: %s", newNodeID),
            Timestamp: time.Now(),
        })
        return nil
    }

    return errors.New("shard not found")
}

// RecordShardMerge logs the merge of two shards.
func (l *ScalabilityLedger) RecordShardMerge(T *BlockchainConsensusCoinLedger,shardID1, shardID2 string) error {
    l.Lock()
    defer l.Unlock()

    // Ensure the Shards map is initialized
    if l.Shards == nil {
        l.Shards = make(map[string]ShardRecord)
    }

    // Log the merge of two shards
    mergeID := shardID1 + "_" + shardID2
    T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
        ID:        mergeID,
        Action:    "ShardMerged",
        Details:   fmt.Sprintf("Merged shards: %s and %s", shardID1, shardID2),
        Timestamp: time.Now(),
    })
    return nil
}

// RecordShardSplit logs the split of a shard into a new shard.
func (l *ScalabilityLedger) RecordShardSplit(T *BlockchainConsensusCoinLedger,shardID, newShardID string) error {
    l.Lock()
    defer l.Unlock()

    // Ensure the Shards map is initialized
    if l.Shards == nil {
        l.Shards = make(map[string]ShardRecord)
    }

    // Log the split of a shard
    T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
        ID:        shardID,
        Action:    "ShardSplit",
        Details:   fmt.Sprintf("Shard split into new shard: %s", newShardID),
        Timestamp: time.Now(),
    })
    return nil
}

// RecordShardAvailability logs the availability of a shard.
func (l *ScalabilityLedger) RecordShardAvailability(T *BlockchainConsensusCoinLedger,shardID string, available bool) error {
    l.Lock()
    defer l.Unlock()

    // Ensure the Shards map is initialized
    if l.Shards == nil {
        l.Shards = make(map[string]ShardRecord)
    }

    // Determine availability status
    status := "Unavailable"
    if available {
        status = "Available"
    }

    // Log shard availability
    T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
        ID:        shardID,
        Action:    "ShardAvailability",
        Details:   status,
        Timestamp: time.Now(),
    })
    return nil
}

// RecordShardAvailabilityChange logs a change in the availability of a shard.
func (l *ScalabilityLedger) RecordShardAvailabilityChange(T *BlockchainConsensusCoinLedger,shardID string, available bool) error {
    l.Lock()
    defer l.Unlock()

    // Ensure the Shards map is initialized
    if l.Shards == nil {
        l.Shards = make(map[string]ShardRecord)
    }

    // Determine availability status
    status := "Unavailable"
    if available {
        status = "Available"
    }

    // Log shard availability change
    T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
        ID:        shardID,
        Action:    "ShardAvailabilityChange",
        Details:   fmt.Sprintf("Availability changed to: %s", status),
        Timestamp: time.Now(),
    })
    return nil
}



// RecordOrchestratedTransaction logs an orchestrated transaction in the ledger.
func (l *ScalabilityLedger) RecordOrchestratedTransaction(T *BlockchainConsensusCoinLedger,txID, orchestrator string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure the Orchestrations map is initialized
	if l.Orchestrations == nil {
		l.Orchestrations = make(map[string]OrchestrationRecord)
	}

	// Record the orchestrated transaction
	l.Orchestrations[txID] = OrchestrationRecord{
		ID:           txID,
		Orchestrator: orchestrator,
		Status:       "completed",
		Timestamp:    time.Now(),
	}

	// Add to history
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        txID,
		Action:    "OrchestratedTransaction",
		Details:   fmt.Sprintf("Orchestrator: %s", orchestrator),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordOrchestrationAdjustment logs an orchestration adjustment.
func (l *ScalabilityLedger) RecordOrchestrationAdjustment(T *BlockchainConsensusCoinLedger,orchestrationID, newConfig string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure the Orchestrations map is initialized
	if l.Orchestrations == nil {
		l.Orchestrations = make(map[string]OrchestrationRecord)
	}

	// Check if the orchestration exists
	orchestration, exists := l.Orchestrations[orchestrationID]
	if !exists {
		return fmt.Errorf("orchestration with ID %s does not exist", orchestrationID)
	}

	// Update the orchestration configuration
	orchestration.Configuration = newConfig
	orchestration.Timestamp = time.Now()
	l.Orchestrations[orchestrationID] = orchestration

	// Add to history
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        orchestrationID,
		Action:    "OrchestrationAdjusted",
		Details:   fmt.Sprintf("New Configuration: %s", newConfig),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordPartitionCreation logs the creation of a partition.
func (l *ScalabilityLedger) RecordPartitionCreation(T *BlockchainConsensusCoinLedger,partitionID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure the Partitions map is initialized
	if l.Partitions == nil {
		l.Partitions = make(map[string]PartitionRecord)
	}

	// Record the new partition
	l.Partitions[partitionID] = PartitionRecord{
		ID:        partitionID,
		Status:    "created",
		Timestamp: time.Now(),
	}

	// Add to history
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        partitionID,
		Action:    "PartitionCreated",
		Details:   "New partition created",
		Timestamp: time.Now(),
	})
	return nil
}

// RecordPartitionRebalance logs a rebalance of a partition.
func (l *ScalabilityLedger) RecordPartitionRebalance(T *BlockchainConsensusCoinLedger,partitionID string, newConfig string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure the Partitions map is initialized
	if l.Partitions == nil {
		l.Partitions = make(map[string]PartitionRecord)
	}

	// Check if the partition exists
	partition, exists := l.Partitions[partitionID]
	if !exists {
		return fmt.Errorf("partition with ID %s does not exist", partitionID)
	}

	// Update the partition configuration
	partition.Configuration = newConfig
	partition.Timestamp = time.Now()
	l.Partitions[partitionID] = partition

	// Add to history
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        partitionID,
		Action:    "PartitionRebalanced",
		Details:   fmt.Sprintf("Rebalanced Configuration: %s", newConfig),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordHorizontalPartitioning logs horizontal partitioning in the ledger.
func (l *ScalabilityLedger) RecordHorizontalPartitioning(T *BlockchainConsensusCoinLedger,partitionID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Partitions map is initialized
	if l.Partitions == nil {
		l.Partitions = make(map[string]PartitionRecord)
	}

	// Log horizontal partitioning
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        partitionID,
		Action:    "HorizontalPartitioning",
		Details:   "Horizontal partitioning applied",
		Timestamp: time.Now(),
	})
	return nil
}

// RecordVerticalPartitioning logs vertical partitioning in the ledger.
func (l *ScalabilityLedger) RecordVerticalPartitioning(T *BlockchainConsensusCoinLedger,partitionID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Partitions map is initialized
	if l.Partitions == nil {
		l.Partitions = make(map[string]PartitionRecord)
	}

	// Log vertical partitioning
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        partitionID,
		Action:    "VerticalPartitioning",
		Details:   "Vertical partitioning applied",
		Timestamp: time.Now(),
	})
	return nil
}



// RecordShardAdjustment logs a dynamic adjustment for a shard in the ledger.
func (l *ScalabilityLedger) RecordShardAdjustment(T *BlockchainConsensusCoinLedger,shardID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	// Check if the shard exists
	if _, exists := l.Shards[shardID]; !exists {
		return errors.New("shard not found")
	}

	// Log shard adjustment
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        shardID,
		Action:    "ShardAdjusted",
		Details:   "Shard dynamically adjusted",
		Timestamp: time.Now(),
	})
	return nil
}

// RecordCrossShardCommunication logs communication between two shards.
func (l *ScalabilityLedger) RecordCrossShardCommunication(T *BlockchainConsensusCoinLedger,shardID1, shardID2, data string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	// Check if the shards exist
	if _, exists1 := l.Shards[shardID1]; !exists1 {
		return errors.New("source shard not found")
	}
	if _, exists2 := l.Shards[shardID2]; !exists2 {
		return errors.New("target shard not found")
	}

	// Log cross-shard communication
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        shardID1 + "_" + shardID2,
		Action:    "CrossShardCommunication",
		Details:   fmt.Sprintf("Data: %s", data),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordPartitionAdjustment logs a partition adjustment.
func (l *ScalabilityLedger) RecordPartitionAdjustment(T *BlockchainConsensusCoinLedger,partitionID, adjustment string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Partitions map is initialized
	if l.Partitions == nil {
		l.Partitions = make(map[string]PartitionRecord)
	}

	// Log the partition adjustment
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        partitionID,
		Action:    "PartitionAdjusted",
		Details:   adjustment,
		Timestamp: time.Now(),
	})
	return nil
}

// RecordShardTransaction logs a transaction within a shard.
func (l *ScalabilityLedger) RecordShardTransaction(shardID string, tx TransactionRecord) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	shard, exists := l.Shards[shardID]
	if !exists {
		return errors.New("shard not found")
	}

	// Add transaction to the shard
	if shard.Transactions == nil {
		shard.Transactions = []TransactionRecord{}
	}
	shard.Transactions = append(shard.Transactions, tx)
	l.Shards[shardID] = shard
	return nil
}

// RecordShardClosure logs the closure of a shard in the system.
func (l *ScalabilityLedger) RecordShardClosure(shardID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	shard, exists := l.Shards[shardID]
	if !exists {
		return errors.New("shard not found")
	}

	// Mark the shard as closed
	shard.Status = "closed"
	shard.ClosedAt = time.Now()

	// Update the shard in the ledger
	l.Shards[shardID] = shard
	return nil
}

// RecordShardSync logs the synchronization of a shard with the main ledger.
func (l *ScalabilityLedger) RecordShardSync(shardID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	shard, exists := l.Shards[shardID]
	if !exists {
		return errors.New("shard not found")
	}

	// Update the shard's sync timestamp
	shard.SyncedAt = time.Now()

	// Update the shard in the ledger
	l.Shards[shardID] = shard
	return nil
}

// RecordShardRebalance logs the rebalancing of resources between shards.
func (l *ScalabilityLedger) RecordShardRebalance(shardID string, newResources int) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	shard, exists := l.Shards[shardID]
	if !exists {
		return errors.New("shard not found")
	}

	// Update the shard's resource allocation
	shard.Resources = newResources

	// Update the shard in the ledger
	l.Shards[shardID] = shard
	return nil
}

// RecordShardValidation logs the validation of a shard's state or proof.
func (l *ScalabilityLedger) RecordShardValidation(shardID, validator string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	shard, exists := l.Shards[shardID]
	if !exists {
		return errors.New("shard not found")
	}

	// Log the validation
	shard.ValidatedBy = validator
	shard.ValidatedAt = time.Now()

	// Update the shard in the ledger
	l.Shards[shardID] = shard
	return nil
}

// RecordStateChannelShardCreation logs the creation of a shard within the state channel system.
func (l *ScalabilityLedger) RecordStateChannelShardCreation(shardID, channelID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Shards map is initialized
	if l.Shards == nil {
		l.Shards = make(map[string]ShardRecord)
	}

	if _, exists := l.Shards[shardID]; exists {
		return errors.New("shard already exists")
	}

	// Create a new shard
	shard := ShardRecord{
		ID:           shardID,
		StateChannel: channelID,
		CreatedAt:    time.Now(),
		Status:       "active",
		IsAvailable:  true,
	}

	// Store the shard in the ledger
	l.Shards[shardID] = shard
	return nil
}
