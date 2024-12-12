package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/scalability"

)


// NewMultiDimensionalCompressionRollup initializes a new Multi-Dimensional Compression Rollup (MDCR) system
func NewMultiDimensionalCompressionRollup(rollupID string, compressionAlgo *common.CompressionAlgorithm, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus, networkManager *common.NetworkManager) *common.MultiDimensionalCompressionRollup {
	return &common.MultiDimensionalCompressionRollup{
		RollupID:        rollupID,
		Transactions:    []*common.Transaction{},
		CompressionAlgo: compressionAlgo,
		IsCompressed:    false,
		Ledger:          ledgerInstance,
		Consensus:       consensus,
		Encryption:      encryptionService,
		NetworkManager:  networkManager,
	}
}

// AddTransaction adds a new transaction to the rollup
func (mdcr *common.MultiDimensionalCompressionRollup) AddTransaction(tx *common.Transaction) error {
	mdcr.mu.Lock()
	defer mdcr.mu.Unlock()

	if mdcr.IsCompressed {
		return errors.New("rollup data is already compressed, no new transactions can be added")
	}

	// Encrypt transaction data before adding
	encryptedTx, err := mdcr.Encryption.EncryptData([]byte(tx.TxID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}
	tx.TxID = string(encryptedTx)

	// Add transaction to the rollup
	mdcr.Transactions = append(mdcr.Transactions, tx)

	// Log transaction addition to the ledger
	err = mdcr.Ledger.RecordTransactionAddition(mdcr.RollupID, tx.TxID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log transaction addition: %v", err)
	}

	fmt.Printf("Transaction %s added to rollup %s\n", tx.TxID, mdcr.RollupID)
	return nil
}

// ApplyCompression applies multi-dimensional compression to the rollup data
func (mdcr *common.MultiDimensionalCompressionRollup) ApplyCompression() error {
	mdcr.mu.Lock()
	defer mdcr.mu.Unlock()

	if mdcr.IsCompressed {
		return errors.New("rollup data is already compressed")
	}

	// Apply the compression using the scalability package
	compressedData, err := scalability.CompressData(mdcr.Transactions, mdcr.CompressionAlgo)
	if err != nil {
		return fmt.Errorf("failed to compress rollup data: %v", err)
	}

	// Encrypt the compressed data
	encryptedData, err := mdcr.Encryption.EncryptData(compressedData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt compressed data: %v", err)
	}
	mdcr.CompressedData = encryptedData
	mdcr.IsCompressed = true

	// Log the compression event in the ledger
	err = mdcr.Ledger.RecordCompressionEvent(mdcr.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log compression event: %v", err)
	}

	fmt.Printf("Rollup %s compressed successfully\n", mdcr.RollupID)
	return nil
}

// FinalizeRollup finalizes the rollup by marking it as compressed and validating the state
func (mdcr *common.MultiDimensionalCompressionRollup) FinalizeRollup() error {
	mdcr.mu.Lock()
	defer mdcr.mu.Unlock()

	if !mdcr.IsCompressed {
		return errors.New("rollup must be compressed before finalizing")
	}

	// Use Synnergy Consensus to validate the rollup state
	err := mdcr.Consensus.ValidateCompression(mdcr.RollupID, mdcr.CompressedData)
	if err != nil {
		return fmt.Errorf("failed to validate compressed rollup: %v", err)
	}

	// Log finalization in the ledger
	err = mdcr.Ledger.RecordRollupFinalization(mdcr.RollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup finalization: %v", err)
	}

	fmt.Printf("Rollup %s finalized successfully\n", mdcr.RollupID)
	return nil
}

// BroadcastRollup sends the compressed rollup data to the network
func (mdcr *common.MultiDimensionalCompressionRollup) BroadcastRollup() error {
	mdcr.mu.Lock()
	defer mdcr.mu.Unlock()

	if !mdcr.IsCompressed {
		return errors.New("rollup is not compressed, cannot broadcast")
	}

	// Broadcast the compressed rollup data to the network
	err := mdcr.NetworkManager.BroadcastData(mdcr.RollupID, mdcr.CompressedData)
	if err != nil {
		return fmt.Errorf("failed to broadcast rollup: %v", err)
	}

	fmt.Printf("Compressed rollup %s broadcasted to the network\n", mdcr.RollupID)
	return nil
}
