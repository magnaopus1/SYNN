package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordChannelCreation creates a new state channel and records it in the ledger.
func (l *StateChannelLedger ) RecordChannelCreation(creator string, participants []string) (string, error) {
    l.Lock()
    defer l.Unlock()

    if l.Channels == nil {
        l.Channels = make(map[string]StateChannel) // Ensure it's a map of StateChannel
    }

    id := l.GenerateUniqueID(creator)

    // Create a new StateChannel based on the provided participants and creator.
    stateChannel := StateChannel{
        ChannelID:    id,
        Participants: participants,
        Status:       "open", // Initial status of the channel
        CreatedAt:    time.Now(),
        State:        "initial", // Assuming an initial state
        Transactions: []TransactionRecord{}, // Empty transactions list
        DataTransfers: make(map[string]*DataBlock), // Empty map of data transfers
        Flexibility:  1.0, // Default flexibility value, you can change this
        Collateral:   make(map[string]*CollateralRecord), // Empty collateral map
        IsOpen:       true, // Channel is initially open
    }

    // Store the new state channel in the ledger
    l.Channels[id] = stateChannel
    return id, nil
}

// RecordChannelClosure closes the specified channel and records it in the ledger.
func (l *StateChannelLedger ) RecordChannelClosure(channelID string) error {
    l.Lock()
    defer l.Unlock()

    channel, exists := l.Channels[channelID]
    if !exists {
        return errors.New("channel not found")
    }

    now := time.Now()
    channel.ClosedAt = now // Directly assign now (type time.Time) to ClosedAt
    channel.Status = "closed"
    l.Channels[channelID] = channel

    return nil
}

// RecordChannelOpening logs the opening of a state channel.
func (l *StateChannelLedger ) RecordChannelOpening(channelID string, participants []string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.Channels[channelID]; exists {
		return errors.New("channel already exists")
	}

	channel := StateChannel{
		ChannelID:    channelID,
		Participants: participants,
		Status:       "open",
		CreatedAt:    time.Now(),
		State:        "active",
	}

	// Record the channel in the ledger
	l.Channels[channelID] = channel
	return nil
}

// RecordLoadMetrics logs the load metrics for a state channel.
func (l *StateChannelLedger) RecordLoadMetrics(channelID string, load int) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	if _, exists := l.Channels[channelID]; !exists {
		return errors.New("channel not found")
	}

	// Create the load metric and log it
	loadMetric := LoadMetric{
		ChannelID: channelID,
		Load:      load,
		Timestamp: time.Now(),
	}

	l.LoadMetrics[channelID] = loadMetric
	return nil
}

// RecordResourceReallocation logs resource reallocation for a state channel.
func (l *StateChannelLedger) RecordResourceReallocation(channelID string, newResources int) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	if _, exists := l.Channels[channelID]; !exists {
		return errors.New("channel not found")
	}

	// Log the resource reallocation
	allocation := ResourceAllocation{
		ChannelID:    channelID,
		Resources:    newResources,
		ReallocatedAt: time.Now(),
	}

	l.ResourceAllocations[channelID] = allocation
	return nil
}

// RecordScalingEvent logs a scaling event for a state channel.
func (l *StateChannelLedger) RecordScalingEvent(eventID, channelID string, oldResources, newResources int) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	if _, exists := l.Channels[channelID]; !exists {
		return errors.New("channel not found")
	}

	// Log the scaling event
	event := ScalingEvent{
		EventID:      eventID,
		ChannelID:    channelID,
		OldResources: oldResources,
		NewResources: newResources,
		Timestamp:    time.Now(),
	}

	l.ScalingEvents[eventID] = event
	return nil
}

// RecordStateFragmentation logs the fragmentation of a state for a state channel.
func (l *StateChannelLedger) RecordStateFragmentation(fragmentID, channelID, fragmentData string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	if _, exists := l.Channels[channelID]; !exists {
		return errors.New("channel not found")
	}

	// Create the fragmented state record
	fragment := FragmentedState{
		FragmentID:   fragmentID,
		ChannelID:    channelID,
		FragmentData: fragmentData,
		Timestamp:    time.Now(),
	}

	l.FragmentedStates[fragmentID] = fragment
	return nil
}

// RecordStateReassembly logs the reassembly of a fragmented state.
func (l *StateChannelLedger) RecordStateReassembly(channelID string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	if _, exists := l.Channels[channelID]; !exists {
		return errors.New("channel not found")
	}

	// Remove all fragments associated with the channel
	for fragmentID, fragment := range l.FragmentedStates {
		if fragment.ChannelID == channelID {
			delete(l.FragmentedStates, fragmentID)
		}
	}

	return nil
}

// RecordCollateralAdjustment logs adjustments to the collateral within a state channel.
func (l *StateChannelLedger) RecordCollateralAdjustment(channelID string, participant string, newCollateral float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	channel, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Check if the participant exists in the channel's collateral map
	collateralRecord, participantExists := channel.Collateral[participant]
	if !participantExists {
		return errors.New("participant not found in the channel")
	}

	// Update the collateral amount for the participant
	collateralRecord.Amount = newCollateral
	collateralRecord.LastUpdated = time.Now()

	// Save the updated channel in the ledger
	l.Channels[channelID] = channel
	return nil
}


// RecordStateChannelStateUpdate logs a state update for a state channel or shard.
func (l *StateChannelLedger) RecordStateChannelStateUpdate(entityID, newState string) error {
	l.Lock()
	defer l.Unlock()

	// Check if it's a state channel or shard and update accordingly
	if channel, exists := l.Channels[entityID]; exists {
		// Update the state for a state channel
		channel.State = newState
		l.Channels[entityID] = channel
	} else if shard, exists := l.Shards[entityID]; exists {
		// Assuming that newState represents availability or another property of the shard
		if newState == "available" {
			shard.IsAvailable = true
		} else if newState == "unavailable" {
			shard.IsAvailable = false
		}
		// Optionally, update other fields such as LastUpdate or StateData if necessary
		shard.LastUpdate = time.Now()
		l.Shards[entityID] = shard
	} else {
		return errors.New("entity not found")
	}

	return nil
}




// RecordStateChannelClosure logs the closure of a state channel.
func (l *StateChannelLedger) RecordStateChannelClosure(channelID string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	channel, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Close the channel
	channel.Status = "closed"
	channel.ClosedAt = time.Now()

	// Update the channel in the ledger
	l.Channels[channelID] = channel
	return nil
}



// RecordTransaction logs a transaction within a shard or channel.
func (l *StateChannelLedger) RecordStateChannelShardTransaction(entityID string, tx TransactionRecord) error {
	l.Lock()
	defer l.Unlock()

	if channel, exists := l.Channels[entityID]; exists {
		// Add transaction to the channel
		channel.Transactions = append(channel.Transactions, tx)
		l.Channels[entityID] = channel
	} else if shard, exists := l.shards[entityID]; exists {
		// Add transaction to the shard
		shard.Transactions = append(shard.Transactions, tx)
		l.shards[entityID] = shard
	} else {
		return errors.New("entity not found")
	}

	return nil
}

// RecordFinality logs the finality of a transaction within a shard or channel.
func (l *StateChannelLedger) RecordFinality(entityID, txID string) error {
	l.Lock()
	defer l.Unlock()

	if channel, exists := l.Channels[entityID]; exists {
		// Update the transaction status in the channel
		for i, tx := range channel.Transactions {
			if tx.Hash == txID {
				channel.Transactions[i].Status = "finalized"
				l.Channels[entityID] = channel
				return nil
			}
		}
	} else if shard, exists := l.shards[entityID]; exists {
		// Update the transaction status in the shard
		for i, tx := range shard.Transactions {
			if tx.Hash == txID {
				shard.Transactions[i].Status = "finalized"
				l.shards[entityID] = shard
				return nil
			}
		}
	} else {
		return errors.New("entity not found")
	}

	return errors.New("transaction not found")
}



// RecordDataTransfer logs a data transfer event between participants in a state channel.
func (l *StateChannelLedger) RecordDataTransfer(channelID, data string) error {
	l.Lock()
	defer l.Unlock()

	channel, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Create the data transfer record
	dataTransfer := DataBlock{
		BlockID:    fmt.Sprintf("block-%d", len(channel.DataTransfers)+1),
		Data:       data,
		Timestamp:  time.Now(),
		MerkleRoot: "", // Assume Merkle root generation logic is handled elsewhere
	}

	// Add the data transfer to the channel
	channel.DataTransfers[dataTransfer.BlockID] = &dataTransfer
	l.Channels[channelID] = channel
	return nil
}




// RecordCollateralValidation logs the validation of collateral in a state channel.
func (l *StateChannelLedger) RecordCollateralValidation(channelID, participant, validator string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the channel exists
	_, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Initialize the collateral map for the channel if it's not already present
	if l.Collateral == nil {
		l.Collateral = make(map[string]map[string]*CollateralRecord)
	}
	if _, exists := l.Collateral[channelID]; !exists {
		return errors.New("collateral not found for channel")
	}

	// Check if the collateral exists for the participant
	record, exists := l.Collateral[channelID][participant]
	if !exists {
		return errors.New("collateral not found for participant")
	}

	// Update the collateral record with validation information
	record.ValidationStatus = "validated"
	record.Validator = validator
	record.ValidatedAt = time.Now()

	// Store the updated collateral record
	l.Collateral[channelID][participant] = record
	return nil
}



// RecordFlexibilityAdjustment logs flexibility adjustments in a state channel.
func (l *StateChannelLedger) RecordFlexibilityAdjustment(channelID string, newFlexibility float64) error {
	l.Lock()
	defer l.Unlock()

	channel, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Adjust the flexibility setting
	channel.Flexibility = newFlexibility
	l.Channels[channelID] = channel
	return nil
}

// RecordStateValidation logs the validation of a state channel update.
func (l *StateChannelLedger) RecordStateValidation(channelID string, validator string) error {
	l.Lock()
	defer l.Unlock()

	channel, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Log the validation event
	channel.LastValidatedBy = validator
	channel.LastValidatedAt = time.Now()

	l.Channels[channelID] = channel
	return nil
}

// RecordStateSync logs the synchronization of the state for a state channel.
func (l *StateChannelLedger) RecordStateSync(channelID string) error {
	l.Lock()
	defer l.Unlock()

	// Sync the state (this could involve complex off-chain synchronization logic)
	channel, exists := l.Channels[channelID]
	if !exists {
		return errors.New("channel not found")
	}

	// Update sync status
	channel.LastSyncedAt = time.Now()
	l.Channels[channelID] = channel

	return nil
}









