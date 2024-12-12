package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/zk_proofs"
)

// NewRecursiveProofAggregation initializes a new RPA system for zk-SNARK proof aggregation.
func NewRecursiveProofAggregation(aggregationID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.RecursiveProofAggregation {
	return &common.RecursiveProofAggregation{
		AggregationID:  aggregationID,
		Rollups:        make(map[string]*common.Rollup),
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AddRollupProof adds a rollup's zk-SNARK proof to the recursive aggregation process.
func (rpa *common.RecursiveProofAggregation) AddRollupProof(rollupID string, zkProof *zk_proofs.ZkSnarkProof) error {
	rpa.mu.Lock()
	defer rpa.mu.Unlock()

	if rpa.IsFinalized {
		return errors.New("aggregation process is already finalized")
	}

	// Ensure rollup exists
	rollup, exists := rpa.Rollups[rollupID]
	if !exists {
		return fmt.Errorf("rollup %s not found", rollupID)
	}

	// Aggregate the proof using recursive zk-SNARK proof aggregation
	if rpa.AggregatedProof == nil {
		rpa.AggregatedProof = zkProof
	} else {
		// Recursively combine zk-SNARK proofs
		aggregatedProof, err := proof.AggregateZkSnarkProofs(rpa.AggregatedProof, zkProof)
		if err != nil {
			return fmt.Errorf("failed to aggregate proofs: %v", err)
		}
		rpa.AggregatedProof = aggregatedProof
	}

	// Log the proof addition in the ledger
	err := rpa.Ledger.RecordProofAddition(rpa.AggregationID, rollupID, zkProof.ProofID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log proof addition: %v", err)
	}

	fmt.Printf("Proof from rollup %s added to recursive proof aggregation %s\n", rollupID, rpa.AggregationID)
	return nil
}

// FinalizeAggregation finalizes the recursive proof aggregation process and computes the final proof.
func (rpa *common.RecursiveProofAggregation) FinalizeAggregation() error {
	rpa.mu.Lock()
	defer rpa.mu.Unlock()

	if rpa.IsFinalized {
		return errors.New("aggregation is already finalized")
	}

	if rpa.AggregatedProof == nil {
		return errors.New("no proofs available for aggregation")
	}

	rpa.IsFinalized = true

	// Log the finalization in the ledger
	err := rpa.Ledger.RecordProofAggregationFinalization(rpa.AggregationID, rpa.AggregatedProof.ProofID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log aggregation finalization: %v", err)
	}

	fmt.Printf("Recursive proof aggregation %s finalized with proof ID %s\n", rpa.AggregationID, rpa.AggregatedProof.ProofID)
	return nil
}

// ValidateAggregatedProof validates the aggregated zk-SNARK proof using Synnergy Consensus.
func (rpa *common.RecursiveProofAggregation) ValidateAggregatedProof() (bool, error) {
	rpa.mu.Lock()
	defer rpa.mu.Unlock()

	if !rpa.IsFinalized {
		return false, errors.New("aggregation is not finalized, cannot validate proof")
	}

	// Validate the aggregated proof using consensus
	valid, err := rpa.Consensus.ValidateZkSnarkProof(rpa.AggregatedProof)
	if err != nil {
		return false, fmt.Errorf("failed to validate zk-SNARK proof: %v", err)
	}

	// Log the proof validation event in the ledger
	err = rpa.Ledger.RecordProofValidation(rpa.AggregationID, rpa.AggregatedProof.ProofID, valid, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to log proof validation: %v", err)
	}

	fmt.Printf("Aggregated proof %s validated with status: %v\n", rpa.AggregatedProof.ProofID, valid)
	return valid, nil
}

// BroadcastAggregatedProof broadcasts the aggregated zk-SNARK proof across the network.
func (rpa *common.RecursiveProofAggregation) BroadcastAggregatedProof() error {
	rpa.mu.Lock()
	defer rpa.mu.Unlock()

	if !rpa.IsFinalized {
		return errors.New("aggregation is not finalized, cannot broadcast proof")
	}

	// Encrypt the aggregated proof before broadcasting
	encryptedProof, err := rpa.Encryption.EncryptData([]byte(rpa.AggregatedProof.ProofID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt proof data: %v", err)
	}

	// Broadcast the proof across the network
	err = rpa.NetworkManager.BroadcastData(rpa.AggregationID, encryptedProof)
	if err != nil {
		return fmt.Errorf("failed to broadcast aggregated proof: %v", err)
	}

	// Log the broadcast event in the ledger
	err = rpa.Ledger.RecordProofBroadcast(rpa.AggregationID, rpa.AggregatedProof.ProofID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log proof broadcast: %v", err)
	}

	fmt.Printf("Aggregated proof %s broadcasted for aggregation %s\n", rpa.AggregatedProof.ProofID, rpa.AggregationID)
	return nil
}

// RetrieveRollupProof retrieves a zk-SNARK proof from a rollup by ID.
func (rpa *common.RecursiveProofAggregation) RetrieveRollupProof(rollupID string, proofID string) (*zk_proofs.ZkSnarkProof, error) {
	rpa.mu.Lock()
	defer rpa.mu.Unlock()

	rollup, exists := rpa.Rollups[rollupID]
	if !exists {
		return nil, fmt.Errorf("rollup %s not found in aggregation %s", rollupID, rpa.AggregationID)
	}

	// Find and retrieve the proof within the rollup
	for _, tx := range rollup.Transactions {
		if tx.ProofID == proofID {
			fmt.Printf("Proof %s retrieved from rollup %s in aggregation %s\n", proofID, rollupID, rpa.AggregationID)
			return tx.ZkProof, nil
		}
	}

	return nil, fmt.Errorf("proof %s not found in rollup %s", proofID, rollupID)
}
