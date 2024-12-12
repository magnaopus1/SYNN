package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// NewInteroperableRollupLayer initializes a new Interoperable Rollup Layer (IRL)
func NewInteroperableRollupLayer(layerID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.InteroperableRollupLayer {
	return &common.InteroperableRollupLayer{
		LayerID:        layerID,
		Rollups:        make(map[string]*common.Rollup),
		SharedState:    make(map[string]interface{}),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AddRollup adds a new rollup to the interoperable rollup layer
func (irl *common.InteroperableRollupLayer) AddRollup(rollup *common.Rollup) error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	if _, exists := irl.Rollups[rollup.RollupID]; exists {
		return errors.New("rollup already exists in the layer")
	}

	irl.Rollups[rollup.RollupID] = rollup

	// Log the addition of the rollup in the ledger
	err := irl.Ledger.RecordRollupAddition(irl.LayerID, rollup.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup addition: %v", err)
	}

	fmt.Printf("Rollup %s added to interoperable layer %s\n", rollup.RollupID, irl.LayerID)
	return nil
}

// RemoveRollup removes a rollup from the interoperable rollup layer
func (irl *common.InteroperableRollupLayer) RemoveRollup(rollupID string) error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	if _, exists := irl.Rollups[rollupID]; !exists {
		return errors.New("rollup not found in the layer")
	}

	delete(irl.Rollups, rollupID)

	// Log the removal of the rollup in the ledger
	err := irl.Ledger.RecordRollupRemoval(irl.LayerID, rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup removal: %v", err)
	}

	fmt.Printf("Rollup %s removed from interoperable layer %s\n", rollupID, irl.LayerID)
	return nil
}

// UpdateSharedState securely updates the shared state across rollups in the layer
func (irl *common.InteroperableRollupLayer) UpdateSharedState(key string, value interface{}) error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	// Update the shared state
	irl.SharedState[key] = value

	// Encrypt the state update before recording it
	encryptedKey, err := irl.Encryption.EncryptData([]byte(key), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt shared state key: %v", err)
	}

	// Log the state update in the ledger
	err = irl.Ledger.RecordStateUpdate(irl.LayerID, string(encryptedKey), value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shared state update: %v", err)
	}

	fmt.Printf("Shared state updated in layer %s: %s = %v\n", irl.LayerID, key, value)
	return nil
}

// RetrieveSharedState retrieves the current shared state across rollups in the layer
func (irl *common.InteroperableRollupLayer) RetrieveSharedState(key string) (interface{}, error) {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	value, exists := irl.SharedState[key]
	if !exists {
		return nil, fmt.Errorf("shared state key %s not found in layer %s", key, irl.LayerID)
	}

	fmt.Printf("Retrieved shared state from layer %s: %s = %v\n", irl.LayerID, key, value)
	return value, nil
}

// CrossRollupTransaction facilitates a transaction between two rollups within the interoperable layer
func (irl *common.InteroperableRollupLayer) CrossRollupTransaction(senderRollupID string, receiverRollupID string, tx *common.Transaction) error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	senderRollup, exists := irl.Rollups[senderRollupID]
	if !exists {
		return fmt.Errorf("sender rollup %s not found in layer %s", senderRollupID, irl.LayerID)
	}

	receiverRollup, exists := irl.Rollups[receiverRollupID]
	if !exists {
		return fmt.Errorf("receiver rollup %s not found in layer %s", receiverRollupID, irl.LayerID)
	}

	// Encrypt the transaction before processing
	encryptedTx, err := irl.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Process the transaction on both rollups
	err = senderRollup.AddTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to process transaction on sender rollup: %v", err)
	}

	err = receiverRollup.AddTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to process transaction on receiver rollup: %v", err)
	}

	// Log the cross-rollup transaction in the ledger
	err = irl.Ledger.RecordCrossRollupTransaction(irl.LayerID, senderRollupID, receiverRollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log cross-rollup transaction: %v", err)
	}

	fmt.Printf("Cross-rollup transaction %s processed between %s and %s in layer %s\n", tx.TxID, senderRollupID, receiverRollupID, irl.LayerID)
	return nil
}

// BroadcastSharedState broadcasts the shared state to all rollups in the layer
func (irl *common.InteroperableRollupLayer) BroadcastSharedState() error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	// Serialize the shared state for broadcast
	stateData := fmt.Sprintf("%v", irl.SharedState)

	// Broadcast the shared state to all rollups in the layer
	err := irl.NetworkManager.BroadcastData(irl.LayerID, []byte(stateData))
	if err != nil {
		return fmt.Errorf("failed to broadcast shared state: %v", err)
	}

	fmt.Printf("Shared state of layer %s broadcasted to all rollups\n", irl.LayerID)
	return nil
}

// FinalizeLayer finalizes all rollups within the interoperable rollup layer
func (irl *common.InteroperableRollupLayer) FinalizeLayer() error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	for _, rollup := range irl.Rollups {
		if !rollup.IsFinalized {
			// Finalize each rollup in the layer
			err := rollup.FinalizeRollup()
			if err != nil {
				return fmt.Errorf("failed to finalize rollup %s: %v", rollup.RollupID, err)
			}
		}
	}

	// Log the finalization of the layer in the ledger
	err := irl.Ledger.RecordLayerFinalization(irl.LayerID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log layer finalization: %v", err)
	}

	fmt.Printf("Interoperable rollup layer %s finalized\n", irl.LayerID)
	return nil
}

// ValidateSharedState uses Synnergy Consensus to validate the shared state across all rollups in the layer
func (irl *common.InteroperableRollupLayer) ValidateSharedState() error {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	// Generate a hash of the shared state for validation
	stateRoot := common.GenerateMerkleRootFromMap(irl.SharedState)

	// Validate the shared state using the Synnergy Consensus
	err := irl.Consensus.ValidateStateRoot(stateRoot)
	if err != nil {
		return fmt.Errorf("shared state validation failed for layer %s: %v", irl.LayerID, err)
	}

	// Log the validation in the ledger
	err = irl.Ledger.RecordSharedStateValidation(irl.LayerID, stateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log shared state validation: %v", err)
	}

	fmt.Printf("Shared state validated for layer %s\n", irl.LayerID)
	return nil
}
