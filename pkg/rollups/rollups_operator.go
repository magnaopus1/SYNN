package rollups

import (
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewRollupOperator initializes a new RollupOperator
func NewRollupOperator(operatorID, nodeID, ipAddress string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.RollupOperator {
	return &common.RollupOperator{
		OperatorID:     operatorID,
		NodeID:         nodeID,
		IPAddress:      ipAddress,
		ManagedBatches: make(map[string]*common.RollupBatch),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		Consensus:      consensus,
		NetworkManager: networkManager,
	}
}

// AggregateAndSubmitBatch aggregates transactions into a rollup batch and submits it to the network
func (ro *common.RollupOperator) AggregateAndSubmitBatch(transactions []*common.Transaction) (*common.RollupBatch, error) {
	ro.mu.Lock()
	defer ro.mu.Unlock()

	// Create a new rollup batch
	batch := &common.RollupBatch{
		BatchID:      common.GenerateUniqueID(),
		Transactions: transactions,
		MerkleRoot:   common.GenerateMerkleRoot(transactions), // Assuming a helper function for Merkle root calculation
		Timestamp:    time.Now(),
	}

	// Encrypt the transactions before submitting
	for _, tx := range transactions {
		encryptedTx, err := ro.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
		}
		tx.TxID = string(encryptedTx)
	}

	// Add the batch to the managed batches
	ro.ManagedBatches[batch.BatchID] = batch

	// Log the batch creation in the ledger
	err := ro.Ledger.RecordBatchCreation(batch.BatchID, batch.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to log batch creation: %v", err)
	}

	// Submit the batch to the rollup network
	err = ro.SubmitBatch(batch)
	if err != nil {
		return nil, fmt.Errorf("failed to submit batch: %v", err)
	}

	fmt.Printf("Batch %s aggregated and submitted by operator %s\n", batch.BatchID, ro.OperatorID)
	return batch, nil
}

// SubmitBatch submits an existing batch to the rollup network for processing
func (ro *common.RollupOperator) SubmitBatch(batch *common.RollupBatch) error {
	ro.mu.Lock()
	defer ro.mu.Unlock()

	// Encrypt the batch data before submitting
	encryptedBatch, err := ro.Encryption.EncryptData([]byte(batch.BatchID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt batch data: %v", err)
	}

	// Broadcast the batch to connected nodes
	for _, node := range ro.NetworkManager.GetConnectedNodes() {
		err := ro.NetworkManager.BroadcastData(node, encryptedBatch)
		if err != nil {
			return fmt.Errorf("failed to broadcast batch %s to node %s: %v", batch.BatchID, node.NodeID, err)
		}
		fmt.Printf("Batch %s submitted to node %s\n", batch.BatchID, node.NodeID)
	}

	// Log the batch submission in the ledger
	err = ro.Ledger.RecordBatchSubmission(batch.BatchID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log batch submission: %v", err)
	}

	return nil
}

// ValidateBatch uses Synnergy Consensus to validate the aggregated transactions in a batch
func (ro *common.RollupOperator) ValidateBatch(batch *common.RollupBatch) error {
	ro.mu.Lock()
	defer ro.mu.Unlock()

	// Validate the batch using consensus
	err := ro.Consensus.ValidateBatch(batch.BatchID, batch.Transactions)
	if err != nil {
		return fmt.Errorf("batch validation failed: %v", err)
	}

	// Log the validation in the ledger
	err = ro.Ledger.RecordBatchValidation(batch.BatchID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log batch validation: %v", err)
	}

	fmt.Printf("Batch %s validated by operator %s\n", batch.BatchID, ro.OperatorID)
	return nil
}

// RetrieveBatch retrieves a batch by its ID from the managed batches
func (ro *common.RollupOperator) RetrieveBatch(batchID string) (*common.RollupBatch, error) {
	ro.mu.Lock()
	defer ro.mu.Unlock()

	batch, exists := ro.ManagedBatches[batchID]
	if !exists {
		return nil, fmt.Errorf("batch %s not found", batchID)
	}

	fmt.Printf("Retrieved batch %s\n", batchID)
	return batch, nil
}

// RemoveBatch removes a batch from the operator's managed batches
func (ro *common.RollupOperator) RemoveBatch(batchID string) error {
	ro.mu.Lock()
	defer ro.mu.Unlock()

	if _, exists := ro.ManagedBatches[batchID]; !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	delete(ro.ManagedBatches, batchID)

	// Log the batch removal in the ledger
	err := ro.Ledger.RecordBatchRemoval(batchID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log batch removal: %v", err)
	}

	fmt.Printf("Batch %s removed by operator %s\n", batchID, ro.OperatorID)
	return nil
}

// MonitorBatchStatus monitors the status of the batch and resubmits if necessary
func (ro *common.RollupOperator) MonitorBatchStatus(batchID string, interval time.Duration) error {
	for {
		time.Sleep(interval)

		batch, err := ro.RetrieveBatch(batchID)
		if err != nil {
			return fmt.Errorf("failed to retrieve batch for monitoring: %v", err)
		}

		// If batch validation fails, try resubmitting the batch
		err = ro.ValidateBatch(batch)
		if err != nil {
			fmt.Printf("Batch %s validation failed: %v. Resubmitting...\n", batchID, err)
			err = ro.SubmitBatch(batch)
			if err != nil {
				return fmt.Errorf("failed to resubmit batch: %v", err)
			}
		}
	}
}
