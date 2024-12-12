package automations

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy_network_demo/ledger"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/transactions"
)

// TransactionFinalizationOptimizationAutomation automates the finalization of validated transactions into sub-blocks and blocks.
type TransactionFinalizationOptimizationAutomation struct {
	ledgerInstance     *ledger.Ledger
	subBlockValidator  *transactions.SubBlockValidator
	blockValidator     *transactions.BlockValidator
	mutex              sync.Mutex
	stopChan           chan bool
}

// NewTransactionFinalizationOptimizationAutomation initializes a new TransactionFinalizationOptimizationAutomation.
func NewTransactionFinalizationOptimizationAutomation(ledgerInstance *ledger.Ledger, subBlockValidator *transactions.SubBlockValidator, blockValidator *transactions.BlockValidator) *TransactionFinalizationOptimizationAutomation {
	return &TransactionFinalizationOptimizationAutomation{
		ledgerInstance:    ledgerInstance,
		subBlockValidator: subBlockValidator,
		blockValidator:    blockValidator,
		stopChan:          make(chan bool),
	}
}

// Start begins the process of monitoring sub-blocks and blocks for finalization.
func (t *TransactionFinalizationOptimizationAutomation) Start() {
	go t.runFinalizationLoop()
	log.Println("Transaction Finalization Optimization Automation started.")
}

// Stop halts the finalization optimization automation.
func (t *TransactionFinalizationOptimizationAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Finalization Optimization Automation stopped.")
}

// runFinalizationLoop continuously monitors the ledger for validated sub-blocks and blocks, finalizing them when ready.
func (t *TransactionFinalizationOptimizationAutomation) runFinalizationLoop() {
	ticker := time.NewTicker(50 * time.Millisecond) // Check every 50ms for sub-block and block validation
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.finalizeValidatedSubBlocks()
			t.finalizeValidatedBlocks()
		case <-t.stopChan:
			return
		}
	}
}

// finalizeValidatedSubBlocks checks the ledger for validated sub-blocks and finalizes them.
func (t *TransactionFinalizationOptimizationAutomation) finalizeValidatedSubBlocks() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Get all sub-blocks that are validated but not yet finalized
	validatedSubBlocks, err := t.ledgerInstance.GetValidatedSubBlocks()
	if err != nil {
		log.Printf("Failed to retrieve validated sub-blocks: %v", err)
		return
	}

	for _, subBlock := range validatedSubBlocks {
		err := t.finalizeSubBlock(subBlock)
		if err != nil {
			log.Printf("Failed to finalize sub-block %s: %v", subBlock.ID, err)
		}
	}
}

// finalizeSubBlock finalizes a validated sub-block by updating the ledger and notifying relevant parties.
func (t *TransactionFinalizationOptimizationAutomation) finalizeSubBlock(subBlock common.SubBlock) error {
	// Encrypt the sub-block ID for secure finalization logging
	encryptedSubBlockID, err := encryption.EncryptData(fmt.Sprintf("SubBlock:%s", subBlock.ID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block ID: %v", err)
	}

	// Record the finalized sub-block in the ledger
	err = t.ledgerInstance.FinalizeSubBlock(subBlock, encryptedSubBlockID)
	if err != nil {
		return fmt.Errorf("failed to finalize sub-block in ledger: %v", err)
	}

	log.Printf("Sub-block %s finalized successfully", subBlock.ID)
	return nil
}

// finalizeValidatedBlocks checks the ledger for validated blocks and finalizes them.
func (t *TransactionFinalizationOptimizationAutomation) finalizeValidatedBlocks() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Get all blocks that are validated but not yet finalized
	validatedBlocks, err := t.ledgerInstance.GetValidatedBlocks()
	if err != nil {
		log.Printf("Failed to retrieve validated blocks: %v", err)
		return
	}

	for _, block := range validatedBlocks {
		err := t.finalizeBlock(block)
		if err != nil {
			log.Printf("Failed to finalize block %s: %v", block.ID, err)
		}
	}
}

// finalizeBlock finalizes a validated block by updating the ledger and notifying relevant parties.
func (t *TransactionFinalizationOptimizationAutomation) finalizeBlock(block common.Block) error {
	// Encrypt the block ID for secure finalization logging
	encryptedBlockID, err := encryption.EncryptData(fmt.Sprintf("Block:%s", block.ID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt block ID: %v", err)
	}

	// Record the finalized block in the ledger
	err = t.ledgerInstance.FinalizeBlock(block, encryptedBlockID)
	if err != nil {
		return fmt.Errorf("failed to finalize block in ledger: %v", err)
	}

	log.Printf("Block %s finalized successfully", block.ID)
	return nil
}
