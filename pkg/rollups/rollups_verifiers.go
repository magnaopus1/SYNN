package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRollupVerifier initializes a new rollup verifier
func NewRollupVerifier(verifierID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager, consensus *common.SynnergyConsensus) *common.RollupVerifier {
	return &common.RollupVerifier{
		VerifierID:        verifierID,
		Ledger:            ledgerInstance,
		Encryption:        encryptionService,
		NetworkManager:    networkManager,
		SynnergyConsensus: consensus,
	}
}

// VerifyRollup validates the integrity of the rollup's state root, transactions, and ensures no fraudulent activity
func (rv *common.RollupVerifier) VerifyRollup(rollup *common.Rollup) error {
	rv.mu.Lock()
	defer rv.mu.Unlock()

	// Verify if rollup is finalized before attempting verification
	if !rollup.IsFinalized {
		return errors.New("cannot verify an unfinalized rollup")
	}

	// Validate the rollup's state root using the Synnergy Consensus mechanism
	valid, err := rv.SynnergyConsensus.ValidateStateRoot(rollup.StateRoot, rollup.Transactions)
	if err != nil {
		return fmt.Errorf("failed to validate rollup state root: %v", err)
	}

	if !valid {
		return errors.New("invalid state root detected in rollup")
	}

	// Log the rollup verification in the ledger
	err = rv.Ledger.RecordRollupVerification(rollup.RollupID, rv.VerifierID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup verification: %v", err)
	}

	fmt.Printf("Rollup %s verified successfully by verifier %s\n", rollup.RollupID, rv.VerifierID)
	return nil
}

// VerifyTransaction validates an individual transaction within the rollup
func (rv *common.RollupVerifier) VerifyTransaction(rollup *common.Rollup, txID string) error {
	rv.mu.Lock()
	defer rv.mu.Unlock()

	// Find the transaction within the rollup
	tx, err := rollup.RetrieveTransaction(txID)
	if err != nil {
		return fmt.Errorf("transaction %s not found in rollup: %v", txID, err)
	}

	// Validate the transaction using Synnergy Consensus
	valid, err := rv.SynnergyConsensus.ValidateTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to validate transaction %s: %v", txID, err)
	}

	if !valid {
		return fmt.Errorf("transaction %s is invalid", txID)
	}

	// Log the transaction verification in the ledger
	err = rv.Ledger.RecordTransactionVerification(txID, rollup.RollupID, rv.VerifierID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction verification: %v", err)
	}

	fmt.Printf("Transaction %s in rollup %s verified successfully by verifier %s\n", txID, rollup.RollupID, rv.VerifierID)
	return nil
}

// VerifyFraudProof checks for fraud proofs submitted by network participants
func (rv *common.RollupVerifier) VerifyFraudProof(rollup *common.Rollup, proof *common.FraudProof) error {
	rv.mu.Lock()
	defer rv.mu.Unlock()

	// Ensure that the rollup is finalized before fraud proof verification
	if !rollup.IsFinalized {
		return errors.New("cannot verify fraud proof for an unfinalized rollup")
	}

	// Validate the fraud proof using Synnergy Consensus
	valid, err := rv.SynnergyConsensus.VerifyFraudProof(proof)
	if err != nil {
		return fmt.Errorf("failed to verify fraud proof: %v", err)
	}

	if !valid {
		return fmt.Errorf("fraud proof is invalid for rollup %s", rollup.RollupID)
	}

	// Log the fraud proof verification in the ledger
	err = rv.Ledger.RecordFraudProofVerification(proof.ProofID, rollup.RollupID, rv.VerifierID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log fraud proof verification: %v", err)
	}

	fmt.Printf("Fraud proof %s for rollup %s verified successfully by verifier %s\n", proof.ProofID, rollup.RollupID, rv.VerifierID)
	return nil
}

// BroadcastVerification broadcasts the verification result to the network
func (rv *common.RollupVerifier) BroadcastVerification(rollup *common.Rollup) error {
	rv.mu.Lock()
	defer rv.mu.Unlock()

	if !rollup.IsFinalized {
		return errors.New("cannot broadcast verification for an unfinalized rollup")
	}

	// Prepare verification data for broadcasting
	data := fmt.Sprintf("RollupID: %s, VerifierID: %s, StateRoot: %s, VerificationTime: %s",
		rollup.RollupID, rv.VerifierID, rollup.StateRoot, time.Now().String())

	// Encrypt the verification data
	encryptedData, err := rv.Encryption.EncryptData([]byte(data), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt verification data: %v", err)
	}

	// Broadcast the verification result to the network
	err = rv.NetworkManager.BroadcastData(rv.VerifierID, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to broadcast verification: %v", err)
	}

	fmt.Printf("Verification result for rollup %s broadcasted by verifier %s\n", rollup.RollupID, rv.VerifierID)
	return nil
}
