package layer2_consensus

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewElasticConsensusManager initializes the Elastic Consensus Manager
func NewElasticConsensusManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *ElasticConsensusManager {
	return &ElasticConsensusManager{
		ConsensusLayers:   make(map[string]*ConsensusLayer),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AddConsensusLayer adds a new elastic consensus layer to the system
func (ecm *ElasticConsensusManager) AddConsensusLayer(layerID, layerType string, maxLoad float64) (*ConsensusLayer, error) {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	// Encrypt layer data
	layerData := fmt.Sprintf("LayerID: %s, Type: %s, MaxLoad: %f", layerID, layerType, maxLoad)
	encryptedData, err := ecm.EncryptionService.EncryptData(layerID, []byte(layerData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt layer data: %v", err)
	}

	// Use the encrypted data (e.g., log it or store it)
	fmt.Printf("Encrypted layer data: %x\n", encryptedData)

	// Create the consensus layer
	layer := &ConsensusLayer{
		LayerID:        layerID,
		LayerType:      layerType,
		CurrentLoad:    0,
		MaxLoad:        maxLoad,
		TransitionTime: time.Now(),
		Active:         false,
		TransitionCount: 0,
	}

	// Add the layer to the manager
	ecm.ConsensusLayers[layerID] = layer

	// Log the addition of the new layer in the ledger
	ecm.Ledger.BlockchainConsensusCoinLedger.RecordConsensusLayerAddition(layerID, layerType) // Removed maxLoad

	fmt.Printf("Consensus layer %s of type %s added\n", layerID, layerType)
	return layer, nil
}




// ActivateConsensusLayer activates a specific consensus layer and transitions from the current one
func (ecm *ElasticConsensusManager) ActivateConsensusLayer(layerID string) error {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	// Retrieve the consensus layer
	layer, exists := ecm.ConsensusLayers[layerID]
	if !exists {
		return fmt.Errorf("consensus layer %s not found", layerID)
	}

	// Deactivate the current layer if one is active
	if ecm.ActiveLayer != nil {
		ecm.ActiveLayer.Active = false
	}

	// Activate the new layer
	layer.Active = true
	layer.TransitionCount++
	layer.TransitionTime = time.Now()
	ecm.ActiveLayer = layer

	// Log the transition in the ledger (no need to assign to err)
	ecm.Ledger.BlockchainConsensusCoinLedger.RecordConsensusLayerTransition(ecm.ActiveLayer.LayerID, layer.LayerType)

	fmt.Printf("Consensus layer %s is now active\n", layerID)
	return nil
}


// MonitorLayerLoad monitors the load on a specific consensus layer and triggers a transition if necessary
func (ecm *ElasticConsensusManager) MonitorLayerLoad(layerID string, currentLoad float64) error {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	// Retrieve the consensus layer
	layer, exists := ecm.ConsensusLayers[layerID]
	if !exists {
		return fmt.Errorf("consensus layer %s not found", layerID)
	}

	// Update the layer's current load
	layer.CurrentLoad = currentLoad

	// Trigger a transition if the load exceeds the maximum allowed
	if currentLoad > layer.MaxLoad {
		fmt.Printf("Load on %s exceeds the maximum (%f > %f), transitioning to another layer...\n", layer.LayerType, currentLoad, layer.MaxLoad)
		for id, nextLayer := range ecm.ConsensusLayers {
			if nextLayer.CurrentLoad < nextLayer.MaxLoad && nextLayer.LayerID != layerID {
				return ecm.ActivateConsensusLayer(id)
			}
		}
	}

	return nil
}

// GetActiveConsensusLayer returns the currently active consensus layer
func (ecm *ElasticConsensusManager) GetActiveConsensusLayer() (*ConsensusLayer, error) {
	ecm.mu.Lock()
	defer ecm.mu.Unlock()

	if ecm.ActiveLayer == nil {
		return nil, errors.New("no active consensus layer")
	}

	return ecm.ActiveLayer, nil
}

