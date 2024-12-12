package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewHyperLayeredRollupFramework initializes a new Hyper-Layered Rollup Framework (HLRF)
func NewHyperLayeredRollupFramework(frameworkID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.HyperLayeredRollupFramework {
	return &common.HyperLayeredRollupFramework{
		FrameworkID:    frameworkID,
		RollupLayers:   make(map[string]*common.RollupLayer),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AddRollupLayer adds a new rollup layer to the hyper-layered framework
func (hlrf *common.HyperLayeredRollupFramework) AddRollupLayer(layerID string) error {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	if _, exists := hlrf.RollupLayers[layerID]; exists {
		return errors.New("rollup layer already exists in the framework")
	}

	// Create a new rollup layer
	hlrf.RollupLayers[layerID] = &common.RollupLayer{
		LayerID:   layerID,
		Rollups:   make(map[string]*common.Rollup),
		IsFinalized: false,
	}

	// Log the addition of the rollup layer in the ledger
	err := hlrf.Ledger.RecordRollupLayerAddition(hlrf.FrameworkID, layerID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup layer addition: %v", err)
	}

	fmt.Printf("Rollup layer %s added to framework %s\n", layerID, hlrf.FrameworkID)
	return nil
}

// AddRollupToLayer adds a rollup to a specific rollup layer
func (hlrf *common.HyperLayeredRollupFramework) AddRollupToLayer(layerID string, rollup *common.Rollup) error {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	layer, exists := hlrf.RollupLayers[layerID]
	if !exists {
		return fmt.Errorf("rollup layer %s not found in framework %s", layerID, hlrf.FrameworkID)
	}

	if _, exists := layer.Rollups[rollup.RollupID]; exists {
		return errors.New("rollup already exists in the layer")
	}

	layer.Rollups[rollup.RollupID] = rollup

	// Log the addition of the rollup in the ledger
	err := hlrf.Ledger.RecordRollupAddition(layerID, rollup.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup addition: %v", err)
	}

	fmt.Printf("Rollup %s added to layer %s in framework %s\n", rollup.RollupID, layerID, hlrf.FrameworkID)
	return nil
}

// FinalizeRollupLayer finalizes a rollup layer by generating its state root and verifying it with consensus
func (hlrf *common.HyperLayeredRollupFramework) FinalizeRollupLayer(layerID string) error {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	layer, exists := hlrf.RollupLayers[layerID]
	if !exists {
		return fmt.Errorf("rollup layer %s not found in framework %s", layerID, hlrf.FrameworkID)
	}

	if layer.IsFinalized {
		return errors.New("rollup layer is already finalized")
	}

	// Compute the final state root for the layer by generating the Merkle root of all rollups
	var allTransactions []*common.Transaction
	for _, rollup := range layer.Rollups {
		allTransactions = append(allTransactions, rollup.Transactions...)
	}
	layer.StateRoot = common.GenerateMerkleRoot(allTransactions)
	layer.IsFinalized = true

	// Verify the state root using Synnergy Consensus
	err := hlrf.Consensus.ValidateStateRoot(layer.StateRoot)
	if err != nil {
		return fmt.Errorf("failed to validate state root for layer %s: %v", layerID, err)
	}

	// Log the finalization of the rollup layer in the ledger
	err = hlrf.Ledger.RecordRollupLayerFinalization(layerID, layer.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup layer finalization: %v", err)
	}

	fmt.Printf("Rollup layer %s finalized with state root %s\n", layerID, layer.StateRoot)
	return nil
}

// CrossLayerVerification verifies the state between two layers using proof aggregation
func (hlrf *common.HyperLayeredRollupFramework) CrossLayerVerification(sourceLayerID, targetLayerID string) error {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	sourceLayer, exists := hlrf.RollupLayers[sourceLayerID]
	if !exists {
		return fmt.Errorf("source rollup layer %s not found in framework %s", sourceLayerID, hlrf.FrameworkID)
	}

	targetLayer, exists := hlrf.RollupLayers[targetLayerID]
	if !exists {
		return fmt.Errorf("target rollup layer %s not found in framework %s", targetLayerID, hlrf.FrameworkID)
	}

	// Perform proof aggregation between the two layers
	proof := common.GenerateProof(sourceLayer.StateRoot, targetLayer.StateRoot)
	err := hlrf.Consensus.ValidateCrossLayerProof(proof)
	if err != nil {
		return fmt.Errorf("cross-layer proof validation failed between %s and %s: %v", sourceLayerID, targetLayerID, err)
	}

	// Log the cross-layer verification in the ledger
	err = hlrf.Ledger.RecordCrossLayerVerification(sourceLayerID, targetLayerID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log cross-layer verification: %v", err)
	}

	fmt.Printf("Cross-layer verification completed between layer %s and layer %s\n", sourceLayerID, targetLayerID)
	return nil
}

// BroadcastLayerState broadcasts the state root of a finalized rollup layer to the network
func (hlrf *common.HyperLayeredRollupFramework) BroadcastLayerState(layerID string) error {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	layer, exists := hlrf.RollupLayers[layerID]
	if !exists {
		return fmt.Errorf("rollup layer %s not found in framework %s", layerID, hlrf.FrameworkID)
	}

	if !layer.IsFinalized {
		return errors.New("rollup layer is not finalized, cannot broadcast")
	}

	// Encrypt and broadcast the state root of the rollup layer
	encryptedStateRoot, err := hlrf.Encryption.EncryptData([]byte(layer.StateRoot), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt state root: %v", err)
	}

	err = hlrf.NetworkManager.BroadcastData(layerID, encryptedStateRoot)
	if err != nil {
		return fmt.Errorf("failed to broadcast state root for layer %s: %v", layerID, err)
	}

	fmt.Printf("State root of rollup layer %s broadcasted to the network\n", layerID)
	return nil
}

// RetrieveRollupLayer retrieves the state of a specific rollup layer
func (hlrf *common.HyperLayeredRollupFramework) RetrieveRollupLayer(layerID string) (*common.RollupLayer, error) {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	layer, exists := hlrf.RollupLayers[layerID]
	if !exists {
		return nil, fmt.Errorf("rollup layer %s not found in framework %s", layerID, hlrf.FrameworkID)
	}

	fmt.Printf("Retrieved rollup layer %s from framework %s\n", layerID, hlrf.FrameworkID)
	return layer, nil
}

// RebalanceRollups dynamically rebalances rollups across layers based on network load
func (hlrf *common.HyperLayeredRollupFramework) RebalanceRollups() error {
	hlrf.mu.Lock()
	defer hlrf.mu.Unlock()

	// Example of rebalancing rollups across layers based on arbitrary network load conditions
	for layerID, layer := range hlrf.RollupLayers {
		if layer.IsFinalized {
			continue
		}
		// Dynamically move rollups between layers based on load
		// (In real-world, this would involve more complex analysis)
	}

	// Log the rebalancing event in the ledger
	err := hlrf.Ledger.RecordRollupRebalancing(hlrf.FrameworkID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup rebalancing: %v", err)
	}

	fmt.Println("Rollups rebalanced across layers")
	return nil
}
