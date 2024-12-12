package rollups

import (
	"errors"
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/zk_proofs" // Assuming a package for zero-knowledge proofs
)

// NewZKRollup initializes a new zero-knowledge proof-enabled rollup
func NewZKRollup(rollupID string, zkProof *zk_proofs.ZKProof, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.ZKRollup {
	return &common.ZKRollup{
		RollupID:       rollupID,
		Transactions:   []*common.Transaction{},
		ZKProof:        zkProof,
		IsFinalized:    false,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
		SynnergyConsensus: consensus,
	}
}

// AddTransaction adds a new transaction to the zero-knowledge proof-enabled rollup
func (r *common.ZKRollup) AddTransaction(tx *common.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized, no new transactions can be added")
	}

	// Encrypt the transaction data
	encryptedTx, err := r.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add the transaction to the rollup
	r.Transactions = append(r.Transactions, tx)

	// Log the transaction addition in the ledger
	err = r.Ledger.RecordTransactionAddition(r.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to ZK rollup %s\n", tx.TxID, r.RollupID)
	return nil
}

// GenerateZKProof generates the zero-knowledge proof for the rollup
func (r *common.ZKRollup) GenerateZKProof() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized, no ZK proof can be generated")
	}

	// Generate the ZK proof for the rollup's transactions
	err := r.ZKProof.Generate(r.Transactions)
	if err != nil {
		return fmt.Errorf("failed to generate ZK proof: %v", err)
	}

	// Log the ZK proof generation in the ledger
	err = r.Ledger.RecordZKProofGeneration(r.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log ZK proof generation: %v", err)
	}

	fmt.Printf("Zero-knowledge proof generated for ZK rollup %s\n", r.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup by computing the final state root, generating the ZK proof, and closing it
func (r *common.ZKRollup) FinalizeRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFinalized {
		return errors.New("rollup is already finalized")
	}

	// Compute the final state root based on all transactions
	r.StateRoot = common.GenerateMerkleRoot(r.Transactions)

	// Generate the ZK proof before finalization
	err := r.GenerateZKProof()
	if err != nil {
		return fmt.Errorf("failed to generate ZK proof: %v", err)
	}

	r.IsFinalized = true

	// Log the rollup finalization in the ledger
	err = r.Ledger.RecordRollupFinalization(r.RollupID, r.StateRoot, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("ZK rollup %s finalized with state root %s\n", r.RollupID, r.StateRoot)
	return nil
}

// VerifyZKProof verifies the zero-knowledge proof for the rollup
func (r *common.ZKRollup) VerifyZKProof() (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verify the ZK proof using Synnergy Consensus
	valid, err := r.SynnergyConsensus.VerifyZKProof(r.ZKProof)
	if err != nil {
		return false, fmt.Errorf("failed to verify ZK proof: %v", err)
	}

	// Log the ZK proof verification in the ledger
	err = r.Ledger.RecordZKProofVerification(r.RollupID, time.Now())
	if err != nil {
		return false, fmt.Errorf("failed to log ZK proof verification: %v", err)
	}

	fmt.Printf("ZK proof for rollup %s verified successfully\n", r.RollupID)
	return valid, nil
}

// RetrieveTransaction retrieves a transaction by its ID from the rollup
func (r *common.ZKRollup) RetrieveTransaction(txID string) (*common.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, tx := range r.Transactions {
		if tx.TxID == txID {
			fmt.Printf("Retrieved transaction %s from ZK rollup %s\n", txID, r.RollupID)
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found in ZK rollup %s", txID, r.RollupID)
}

// BroadcastRollup sends the finalized rollup data and the ZK proof to the network
func (r *common.ZKRollup) BroadcastRollup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.IsFinalized {
		return errors.New("rollup is not finalized, cannot broadcast")
	}

	// Encrypt the ZK proof before broadcasting
	encryptedProof, err := r.Encryption.EncryptData([]byte(r.ZKProof.String()), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt ZK proof: %v", err)
	}

	// Broadcast the rollup data and the ZK proof to the network
	err = r.NetworkManager.BroadcastData(r.RollupID, encryptedProof)
	if err != nil {
		return fmt.Errorf("failed to broadcast ZK rollup: %v", err)
	}

	fmt.Printf("ZK rollup %s and its proof broadcasted to the network\n", r.RollupID)
	return nil
}
