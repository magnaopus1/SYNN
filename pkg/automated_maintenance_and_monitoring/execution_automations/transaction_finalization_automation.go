package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/blocks"
	"synnergy_network_demo/transactions"
)

const (
	SubBlockThreshold        = 1000              // Number of sub-blocks to validate into a full block
	TransactionFinalizationInterval = 1 * time.Minute // Interval for checking pending finalization
	SubBlockFinalizationEntryType   = "SubBlockFinalization"
	BlockFinalizationEntryType      = "BlockFinalization"
)

// TransactionFinalizationAutomation handles the finalization of sub-blocks into full blocks
type TransactionFinalizationAutomation struct {
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	blockchain       *blocks.Blockchain
	finalizationMutex *sync.RWMutex
}

// NewTransactionFinalizationAutomation initializes the finalization automation with consensus, ledger, and blockchain integration
func NewTransactionFinalizationAutomation(consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, blockchain *blocks.Blockchain, finalizationMutex *sync.RWMutex) *TransactionFinalizationAutomation {
	return &TransactionFinalizationAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		blockchain:       blockchain,
		finalizationMutex: finalizationMutex,
	}
}

// StartFinalizationMonitor starts continuous monitoring and finalization of sub-blocks into blocks
func (automation *TransactionFinalizationAutomation) StartFinalizationMonitor() {
	ticker := time.NewTicker(TransactionFinalizationInterval)

	go func() {
		for range ticker.C {
			automation.finalizePendingSubBlocks()
		}
	}()
}

// finalizePendingSubBlocks finalizes the pending sub-blocks into full blocks
func (automation *TransactionFinalizationAutomation) finalizePendingSubBlocks() {
	automation.finalizationMutex.Lock()
	defer automation.finalizationMutex.Unlock()

	// Get all pending sub-blocks that need finalization
	pendingSubBlocks := automation.blockchain.GetPendingSubBlocks()

	// If there are 1000 sub-blocks or more, finalize into a full block
	if len(pendingSubBlocks) >= SubBlockThreshold {
		fullBlock, err := automation.createFullBlock(pendingSubBlocks)
		if err != nil {
			fmt.Printf("Failed to finalize sub-blocks into a full block: %v\n", err)
			return
		}

		// Log the full block finalization into the ledger
		automation.logBlockFinalization(fullBlock)

		// Clear finalized sub-blocks from the blockchain
		automation.blockchain.ClearSubBlocks()
	}
}

// createFullBlock takes the pending sub-blocks and validates them into a full block
func (automation *TransactionFinalizationAutomation) createFullBlock(subBlocks []*blocks.SubBlock) (*blocks.Block, error) {
	// Validate the sub-blocks using Synnergy Consensus
	valid, err := automation.consensusEngine.ValidateSubBlocks(subBlocks)
	if err != nil || !valid {
		return nil, fmt.Errorf("consensus validation failed for sub-blocks: %v", err)
	}

	// Create the full block from the validated sub-blocks
	fullBlock, err := automation.blockchain.CreateBlockFromSubBlocks(subBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create full block from sub-blocks: %v", err)
	}

	// Broadcast the newly created block to the network
	err = automation.consensusEngine.BroadcastBlock(fullBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast block: %v", err)
	}

	fmt.Printf("Full block %s successfully created and broadcasted.\n", fullBlock.ID)
	return fullBlock, nil
}

// logBlockFinalization logs the successful finalization of a block into the ledger
func (automation *TransactionFinalizationAutomation) logBlockFinalization(block *blocks.Block) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("block-finalization-%s-%d", block.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      BlockFinalizationEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("Block %s successfully finalized.", block.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log block finalization %s in the ledger: %v\n", block.ID, err)
	} else {
		fmt.Println("Block finalization successfully logged in the ledger.")
	}
}

// logSubBlockFinalization logs the successful finalization of a sub-block into the ledger
func (automation *TransactionFinalizationAutomation) logSubBlockFinalization(subBlock *blocks.SubBlock) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("subblock-finalization-%s-%d", subBlock.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      SubBlockFinalizationEntryType,
		Status:    "Success",
		Details:   fmt.Sprintf("Sub-block %s successfully finalized.", subBlock.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log sub-block finalization %s in the ledger: %v\n", subBlock.ID, err)
	} else {
		fmt.Println("Sub-block finalization successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TransactionFinalizationAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualFinalization allows administrators to manually trigger the finalization of sub-blocks into a full block
func (automation *TransactionFinalizationAutomation) TriggerManualFinalization() {
	fmt.Println("Manually triggering finalization of pending sub-blocks...")

	automation.finalizePendingSubBlocks()
}
