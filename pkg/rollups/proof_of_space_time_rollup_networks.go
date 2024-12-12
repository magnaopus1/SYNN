package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/space_time_proofs" 
)

// NewPoSTRollup initializes a new Proof-of-Space-Time rollup.
func NewPoSTRollup(rollupID string, spaceTimeProof *space_time_proofs.SpaceTimeProof, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.PoSTRollup {
	return &common.PoSTRollup{
		RollupID:         rollupID,
		Transactions:     []*common.Transaction{},
		IsFinalized:      false,
		SpaceTimeProof:   spaceTimeProof,
		Ledger:           ledgerInstance,
		Encryption:       encryptionService,
		NetworkManager:   networkManager,
		SynnergyConsensus: consensus,
	}
}

// AddTransaction adds a new transaction to the rollup.
func (r *common.PoSTRollup) AddTransaction(tx *common.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized, no new transactions can be added")
	}

	// Encrypt the transaction data.
	encryptedTx, err := r.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the rollup.
	r.Transactions = append(r.Transactions, tx)

	// Log the transaction addition in the ledger.
	err = r.Ledger.RecordTransactionAddition(r.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to PoST rollup %s\n", tx.TxID, r.RollupID)
	return nil
}

// GenerateSpaceTimeProof generates the proof-of-space-time for the rollup's data.
func (r *common.PoSTRollup) GenerateSpaceTimeProof() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized, no PoST can be generated")
	}

	// Generate the PoST for the rollup's transactions.
	err := r.SpaceTimeProof.Generate(r.Transactions)
	if err != nil {
		return fmt.Errorf("failed to generate space-time proof: %v", err)
	}

	// Log the space-time proof generation in the ledger.
	err = r.Ledger.RecordSpaceTimeProofGeneration(r.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log space-time proof generation: %v", err)
	}

	fmt.Printf("Proof-of-space-time generated for PoST rollup %s\n", r.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup by computing the final state root, generating the PoST, and closing it.
func (r *common.PoSTRollup) FinalizeRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Compute the final state root based on all transactions.
	r.StateRoot = common.GenerateMerkleRoot(r.Transactions)

	// Generate the space-time proof before finalization.
	err := r.GenerateSpaceTimeProof()
	if err != nil {
		return fmt.Errorf("failed to generate space-time proof: %v", err)
	}

	r.IsFinalized = true

	// Log the rollup finalization in the ledger.
	err = r.Ledger.RecordRollupFinalization(r.RollupID, r.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("PoST rollup %s finalized with state root %s\n", r.RollupID, r.StateRoot)
	return nil
}

// VerifySpaceTimeProof verifies the proof-of-space-time for the rollup.
func (r *common.PoSTRollup) VerifySpaceTimeProof() (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verify the space-time proof using Synnergy Consensus.
	valid, err := r.SynnergyConsensus.VerifySpaceTimeProof(r.SpaceTimeProof)
	if err != nil {
		return false, fmt.Errorf("failed to verify space-time proof: %v", err)
	}

	// Log the space-time proof verification in the ledger.
	err = r.Ledger.RecordSpaceTimeProofVerification(r.RollupID, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to log space-time proof verification: %v", err)
	}

	fmt.Printf("Proof-of-space-time for rollup %s verified successfully\n", r.RollupID)
	return valid, nil
}

// RetrieveTransaction retrieves a transaction by its ID from the rollup.
func (r *common.PoSTRollup) RetrieveTransaction(txID string) (*common.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, tx := range r.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from PoST rollup %s\n", txID, r.RollupID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found in PoST rollup %s", txID, r.RollupID)
}

// BroadcastRollup sends the finalized rollup data and the proof-of-space-time to the network.
func (r *common.PoSTRollup) BroadcastRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Encrypt the space-time proof before broadcasting.
	encryptedProof, err := r.Encryption.EncryptData([]byte(r.SpaceTimeProof.String()), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt space-time proof: %v", err)
	}

	// Broadcast the rollup data and the space-time proof to the network.
	err = r.NetworkManager.BroadcastData(r.RollupID, encryptedProof)
	if err != nil {
		return fmt.Errorf("failed to broadcast PoST rollup: %v", err)
	}

	fmt.Printf("PoST rollup %s and its proof broadcasted to the network\n", r.RollupID)
	return nil
}
