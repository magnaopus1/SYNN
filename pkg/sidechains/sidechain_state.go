package sidechains

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
)

// NewSidechainState initializes the state for a sidechain
func NewSidechainState(chainID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.SidechainState {
	return &common.SidechainState{
		ChainID:        chainID,
		StateData:      make(map[string]*common.StateObject),
		BlockStates:    make(map[string]*common.BlockState),
		SubBlockStates: make(map[string]*common.SubBlockState),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
	}
}

// UpdateState updates the overall state of the sidechain
func (ss *common.SidechainState) UpdateState(stateID string, stateObject *common.StateObject) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Encrypt the state data before updating
	encryptedStateData, err := ss.Encryption.EncryptData([]byte(stateObject.Data), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state data: %v", err)
	}
	stateObject.Data = string(encryptedStateData)

	// Update the state data
	ss.StateData[stateID] = stateObject

	// Log the state update in the ledger
	err = ss.Ledger.RecordStateUpdate(ss.ChainID, stateID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State for %s updated in sidechain %s\n", stateID, ss.ChainID)
	return nil
}

// GetState retrieves the current state of the sidechain
func (ss *common.SidechainState) GetState(stateID string) (*common.StateObject, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	state, exists := ss.StateData[stateID]
	if !exists {
		return nil, fmt.Errorf("state %s not found in sidechain %s", stateID, ss.ChainID)
	}

	// Decrypt the state data before returning
	decryptedStateData, err := ss.Encryption.DecryptData([]byte(state.Data), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state data: %v", err)
	}
	state.Data = string(decryptedStateData)

	fmt.Printf("State for %s retrieved from sidechain %s\n", stateID, ss.ChainID)
	return state, nil
}

// UpdateBlockState updates the state of a specific block in the sidechain
func (ss *common.SidechainState) UpdateBlockState(blockID string, stateData map[string]*common.StateObject) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Encrypt each state object before updating
	for stateID, stateObject := range stateData {
		encryptedStateData, err := ss.Encryption.EncryptData([]byte(stateObject.Data), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt state data for block %s: %v", blockID, err)
		}
		stateObject.Data = string(encryptedStateData)
		stateData[stateID] = stateObject
	}

	// Update the block state
	blockState := &common.BlockState{
		BlockID:   blockID,
		StateData: stateData,
		Timestamp: time.Now(),
	}
	ss.BlockStates[blockID] = blockState

	// Log the block state update in the ledger
	err := ss.Ledger.RecordBlockStateUpdate(ss.ChainID, blockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block state update: %v", err)
	}

	fmt.Printf("State for block %s updated in sidechain %s\n", blockID, ss.ChainID)
	return nil
}

// GetBlockState retrieves the state of a specific block in the sidechain
func (ss *common.SidechainState) GetBlockState(blockID string) (*common.BlockState, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	blockState, exists := ss.BlockStates[blockID]
	if !exists {
		return nil, fmt.Errorf("block state %s not found in sidechain %s", blockID, ss.ChainID)
	}

	// Decrypt each state object in the block state before returning
	for stateID, stateObject := range blockState.StateData {
		decryptedStateData, err := ss.Encryption.DecryptData([]byte(stateObject.Data), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt state data for block %s: %v", blockID, err)
		}
		stateObject.Data = string(decryptedStateData)
		blockState.StateData[stateID] = stateObject
	}

	fmt.Printf("State for block %s retrieved from sidechain %s\n", blockID, ss.ChainID)
	return blockState, nil
}

// UpdateSubBlockState updates the state of a specific sub-block in the sidechain
func (ss *common.SidechainState) UpdateSubBlockState(subBlockID string, stateData map[string]*common.StateObject) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Encrypt each state object before updating
	for stateID, stateObject := range stateData {
		encryptedStateData, err := ss.Encryption.EncryptData([]byte(stateObject.Data), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt state data for sub-block %s: %v", subBlockID, err)
		}
		stateObject.Data = string(encryptedStateData)
		stateData[stateID] = stateObject
	}

	// Update the sub-block state
	subBlockState := &common.SubBlockState{
		SubBlockID: subBlockID,
		StateData:  stateData,
		Timestamp:  time.Now(),
	}
	ss.SubBlockStates[subBlockID] = subBlockState

	// Log the sub-block state update in the ledger
	err := ss.Ledger.RecordSubBlockStateUpdate(ss.ChainID, subBlockID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sub-block state update: %v", err)
	}

	fmt.Printf("State for sub-block %s updated in sidechain %s\n", subBlockID, ss.ChainID)
	return nil
}

// GetSubBlockState retrieves the state of a specific sub-block in the sidechain
func (ss *common.SidechainState) GetSubBlockState(subBlockID string) (*common.SubBlockState, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	subBlockState, exists := ss.SubBlockStates[subBlockID]
	if !exists {
		return nil, fmt.Errorf("sub-block state %s not found in sidechain %s", subBlockID, ss.ChainID)
	}

	// Decrypt each state object in the sub-block state before returning
	for stateID, stateObject := range subBlockState.StateData {
		decryptedStateData, err := ss.Encryption.DecryptData([]byte(stateObject.Data), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt state data for sub-block %s: %v", subBlockID, err)
		}
		stateObject.Data = string(decryptedStateData)
		subBlockState.StateData[stateID] = stateObject
	}

	fmt.Printf("State for sub-block %s retrieved from sidechain %s\n", subBlockID, ss.ChainID)
	return subBlockState, nil
}

// SyncStateAcrossChains synchronizes the state between the sidechain and the mainchain or another sidechain
func (ss *common.SidechainState) SyncStateAcrossChains(destinationChainID, stateID string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Encrypt the state data before syncing
	state, exists := ss.StateData[stateID]
	if !exists {
		return fmt.Errorf("state %s not found in sidechain %s", stateID, ss.ChainID)
	}

	encryptedStateData, err := ss.Encryption.EncryptData([]byte(state.Data), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state data: %v", err)
	}
	state.Data = string(encryptedStateData)

	// Sync the state across chains (implementation of the network sync depends on the network package)
	err = ss.Consensus.SyncState(ss.ChainID, destinationChainID, stateID, encryptedStateData)
	if err != nil {
		return fmt.Errorf("failed to sync state across chains: %v", err)
	}

	// Log the state sync in the ledger
	err = ss.Ledger.RecordStateSync(ss.ChainID, destinationChainID, stateID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state sync: %v", err)
	}

	fmt.Printf("State %s synced from sidechain %s to chain %s\n", stateID, ss.ChainID, destinationChainID)
	return nil
}

// ValidateState ensures the integrity of the sidechain's state using Synnergy Consensus
func (ss *common.SidechainState) ValidateState(stateID string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	state, exists := ss.StateData[stateID]
	if !exists {
		return fmt.Errorf("state %s not found in sidechain %s", stateID, ss.ChainID)
	}

	// Validate the state using Synnergy Consensus
	err := ss.Consensus.ValidateState(ss.ChainID, stateID, state)
	if err != nil {
		return fmt.Errorf("state validation failed: %v", err)
	}

	// Log the state validation in the ledger
	err = ss.Ledger.RecordStateValidation(ss.ChainID, stateID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state validation: %v", err)
	}

	fmt.Printf("State %s validated in sidechain %s\n", stateID, ss.ChainID)
	return nil
}
