package rollups

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewHierarchicalProofRollup initializes a new rollup with hierarchical proofing
func NewHierarchicalProofRollup(rollupID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.HierarchicalProofRollup {
	return &common.HierarchicalProofRollup{
		RollupID:       rollupID,
		ProofHierarchy: make(map[string]*common.RollupProofLayer),
		Transactions:   []*common.Transaction{},
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AddTransaction adds a new transaction to the rollup
func (hpr *common.HierarchicalProofRollup) AddTransaction(tx *common.Transaction) error {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	if hpr.IsFinalized {
		return errors.New("rollup is already finalized, cannot add new transactions")
	}

	// Encrypt the transaction data
	encryptedTx, err := hpr.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the rollup
	hpr.Transactions = append(hpr.Transactions, tx)

	// Log the transaction addition in the ledger
	err = hpr.Ledger.RecordTransactionAddition(hpr.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to rollup %s\n", tx.TxID, hpr.RollupID)
	return nil
}

// CreateProofLayer generates a new proof for a layer in the rollup's proof hierarchy
func (hpr *common.HierarchicalProofRollup) CreateProofLayer(layerID, parentLayer string) (*common.RollupProofLayer, error) {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	// Check if the layer already exists
	if _, exists := hpr.ProofHierarchy[layerID]; exists {
		return nil, errors.New("proof layer already exists")
	}

	// Generate a proof for the current transactions (Merkle root)
	layerProof := common.GenerateMerkleRoot(hpr.Transactions)

	// Create the proof layer and add it to the hierarchy
	newLayer := &common.RollupProofLayer{
		LayerID:    layerID,
		Proof:      layerProof,
		ParentLayer: parentLayer,
		IsVerified: false,
	}

	hpr.ProofHierarchy[layerID] = newLayer

	// Log the creation of the proof layer in the ledger
	err := hpr.Ledger.RecordProofLayerCreation(hpr.RollupID, layerID, layerProof, parentLayer, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log proof layer creation: %v", err)
	}

	fmt.Printf("Proof layer %s created for rollup %s\n", layerID, hpr.RollupID)
	return newLayer, nil
}

// VerifyProofLayer verifies a specific proof layer using Synnergy Consensus
func (hpr *common.HierarchicalProofRollup) VerifyProofLayer(layerID string) error {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	layer, exists := hpr.ProofHierarchy[layerID]
	if !exists {
		return fmt.Errorf("proof layer %s not found in rollup %s", layerID, hpr.RollupID)
	}

	// Verify the proof using Synnergy Consensus
	err := hpr.Consensus.ValidateLayerProof(layer.Proof)
	if err != nil {
		return fmt.Errorf("proof verification failed for layer %s: %v", layerID, err)
	}

	layer.IsVerified = true

	// Log the proof verification in the ledger
	err = hpr.Ledger.RecordProofLayerVerification(hpr.RollupID, layerID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log proof verification: %v", err)
	}

	fmt.Printf("Proof layer %s verified for rollup %s\n", layerID, hpr.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup, ensuring all proof layers are verified and the final state is computed
func (hpr *common.HierarchicalProofRollup) FinalizeRollup() error {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	if hpr.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Ensure all proof layers are verified before finalizing
	for layerID, layer := range hpr.ProofHierarchy {
		if !layer.IsVerified {
			return fmt.Errorf("proof layer %s is not verified", layerID)
		}
	}

	// Compute the final state root based on the transactions
	hpr.StateRoot = common.GenerateMerkleRoot(hpr.Transactions)
	hpr.IsFinalized = true

	// Log the rollup finalization in the ledger
	err := hpr.Ledger.RecordRollupFinalization(hpr.RollupID, hpr.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("Rollup %s finalized with state root %s\n", hpr.RollupID, hpr.StateRoot)
	return nil
}

// BroadcastRollupProofs broadcasts the finalized rollup's proof hierarchy to the network
func (hpr *common.HierarchicalProofRollup) BroadcastRollupProofs() error {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	if !hpr.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast proofs")
	}

	// Encrypt and broadcast each proof layer in the hierarchy
	for layerID, layer := range hpr.ProofHierarchy {
		encryptedProof, err := hpr.Encryption.EncryptData([]byte(layer.Proof), common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt proof for layer %s: %v", layerID, err)
		}

		err = hpr.NetworkManager.BroadcastData(layerID, encryptedProof)
		if err != nil {
			return fmt.Errorf("failed to broadcast proof for layer %s: %v", layerID, err)
		}
	}

	fmt.Printf("Rollup %s proofs broadcasted to the network\n", hpr.RollupID)
	return nil
}

// RetrieveProofLayer retrieves the proof layer by its ID from the hierarchical rollup
func (hpr *common.HierarchicalProofRollup) RetrieveProofLayer(layerID string) (*common.RollupProofLayer, error) {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	layer, exists := hpr.ProofHierarchy[layerID]
	if !exists {
		return nil, fmt.Errorf("proof layer %s not found in rollup %s", layerID, hpr.RollupID)
	}

	fmt.Printf("Retrieved proof layer %s from rollup %s\n", layerID, hpr.RollupID)
	return layer, nil
}

// CrossLayerProofVerification verifies proofs between two rollup layers in the hierarchy
func (hpr *common.HierarchicalProofRollup) CrossLayerProofVerification(sourceLayerID, targetLayerID string) error {
	hpr.mu.Lock()
	defer hpr.mu.Unlock()

	sourceLayer, exists := hpr.ProofHierarchy[sourceLayerID]
	if !exists {
		return fmt.Errorf("source proof layer %s not found in rollup %s", sourceLayerID, hpr.RollupID)
	}

	targetLayer, exists := hpr.ProofHierarchy[targetLayerID]
	if !exists {
		return fmt.Errorf("target proof layer %s not found in rollup %s", targetLayerID, hpr.RollupID)
	}

	// Perform cross-layer proof verification (this would involve more advanced verification logic)
	verified := sourceLayer.Proof == targetLayer.ParentLayer // Simplified verification

	if verified {
		fmt.Printf("Cross-layer proof verification successful between layer %s and layer %s in rollup %s\n", sourceLayerID, targetLayerID, hpr.RollupID)
	} else {
		return fmt.Errorf("cross-layer proof verification failed between layer %s and layer %s", sourceLayerID, targetLayerID)
	}

	return nil
}
