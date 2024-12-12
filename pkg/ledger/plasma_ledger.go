package ledger

import (
    "fmt"
    "time"
)

// RecordTransaction logs a new transaction in the Plasma ledger.
func (ledger *PlasmaLedger) RecordTransaction(tx TransactionRecord) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.transactionCache == nil {
		ledger.transactionCache = make(map[string]TransactionRecord)
	}

	ledger.pendingTransactions = append(ledger.pendingTransactions, &tx)  // Pass the pointer to tx
	ledger.transactionCache[tx.Hash] = tx
	fmt.Printf("Transaction %s recorded in the Plasma ledger.\n", tx.Hash)
	return nil
}



// RecordTransactionSubmission logs the submission of a transaction to the Plasma chain.
func (ledger *PlasmaLedger) RecordTransactionSubmission(tx TransactionRecord) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Transaction %s submitted to Plasma ledger for processing.\n", tx.Hash)
	return nil
}

// RecordSubBlockValidation records the validation of a sub-block in the Plasma chain.
func (ledger *PlasmaLedger) RecordSubBlockValidation(subBlock SubBlock) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Sub-block #%d validated.\n", subBlock.Index)
	return nil
}

// RecordSubBlockCreation records the creation of a new sub-block in the Plasma chain.
func (ledger *PlasmaLedger) RecordSubBlockCreation(subBlock *PlasmaSubBlock) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Add the new sub-block to the map of sub-blocks, using its ID as the key
	ledger.subBlocks[subBlock.SubBlockID] = subBlock
	ledger.subBlockIndex++

	fmt.Printf("Sub-block #%s created with hash %s.\n", subBlock.SubBlockID, subBlock.ParentBlockID)
	return nil
}

// RecordSubBlockAddition adds a validated sub-block to the Plasma chain ledger.
func (ledger *PlasmaLedger) RecordSubBlockAddition(subBlock *PlasmaSubBlock) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Move the sub-block from the current map to finalized sub-blocks
	ledger.finalizedSubBlocks[subBlock.SubBlockID] = subBlock
	ledger.State.SubBlockHeight++
	ledger.State.LastSubBlockHash = subBlock.ParentBlockID

	fmt.Printf("Sub-block #%s added to finalized sub-blocks.\n", subBlock.SubBlockID)
	return nil
}


// RecordBlockCreation handles the creation of a block from a set of sub-blocks.
func (ledger *PlasmaLedger) RecordBlockCreation() error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if len(ledger.subBlocks) >= 1000 {
        newBlock := PlasmaBlock{
            BlockID:       generateBlockID(),               // Generate a unique block ID
            PreviousBlock: ledger.State.LastSubBlockHash,   // Use the hash of the last sub-block as the previous block
            SubBlocks:     ledger.subBlocks,               // Assign subBlocks as a map directly
            Timestamp:     time.Now(),
            ValidatorID:   "validator-id",                 // Placeholder for the validator's ID
        }
        // Move finalized sub-blocks and reset subBlocks
        ledger.finalizedSubBlocks = ledger.subBlocks
        ledger.subBlocks = make(map[string]*PlasmaSubBlock) // Reset sub-blocks after creating the block
        fmt.Printf("Plasma block created with ID %s.\n", newBlock.BlockID)
        return nil
    }
    fmt.Println("Not enough sub-blocks to create a Plasma block.")
    return nil
}

// FinalizeSubBlocksIntoPlasmaBlock finalizes and packages 1000 sub-blocks into a Plasma block.
func (ledger *PlasmaLedger) FinalizeSubBlocksIntoPlasmaBlock() error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if len(ledger.subBlocks) >= 1000 {
        // Create a new Plasma block
        newBlock := PlasmaBlock{
            BlockID:       generateBlockID(),
            PreviousBlock: ledger.State.LastSubBlockHash,
            SubBlocks:     ledger.subBlocks,               // No need for conversion; use the map directly
            Timestamp:     time.Now(),
            ValidatorID:   "validator-id",                 // Validator responsible for this block
        }
        // Finalize sub-blocks into the new block and reset subBlocks map
        ledger.finalizedSubBlocks = ledger.subBlocks
        ledger.subBlocks = make(map[string]*PlasmaSubBlock) // Reset after finalization
        fmt.Printf("Finalized Plasma block #%s.\n", newBlock.BlockID)
    } else {
        fmt.Println("Not enough sub-blocks to finalize into a Plasma block.")
    }
    return nil
}


// generateBlockID creates a unique block ID (this is a simple placeholder implementation).
func generateBlockID() string {
	return fmt.Sprintf("block-%d", time.Now().UnixNano())
}


// RecordBlockFinalization finalizes a Plasma block and syncs it with the root chain.
func (ledger *PlasmaLedger) RecordBlockFinalization(rootChainTx RootChainTxRecord) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	ledger.State.RootChainTx = rootChainTx
	fmt.Printf("Sub-block set finalized and submitted to root chain with Tx: %s.\n", rootChainTx.TxHash)
	return nil
}

// RecordTransactionSpent marks a transaction as spent within the Plasma ledger.
func (ledger *PlasmaLedger) RecordTransactionSpent(txHash string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if tx, exists := ledger.transactionCache[txHash]; exists {
		tx.Status = "spent"
		ledger.transactionCache[txHash] = tx
		fmt.Printf("Transaction %s marked as spent.\n", txHash)
		return nil
	}
	return fmt.Errorf("Transaction %s not found in cache.\n", txHash)
}

// RecordCrossChainTransfer logs a cross-chain transfer between Plasma and another chain.
func (ledger *PlasmaLedger) RecordCrossChainTransfer(tx CrossChainTransfer) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Cross-chain transfer initiated from Plasma to chain %s.\n", tx.ToChain)
	return nil
}

// RecordCrossChainTransactionFinalization finalizes a cross-chain transaction.
func (ledger *PlasmaLedger) RecordCrossChainTransactionFinalization(tx CrossChainTransfer) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Cross-chain transaction finalized from Plasma to chain %s.\n", tx.ToChain)
	return nil
}

// RecordNodeAddition logs the addition of a node to the Plasma network.
func (ledger *PlasmaLedger) RecordNodeAddition(node NodeInfo) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	// Use NodeID instead of ID
	ledger.State.Nodes[node.NodeID] = node
	fmt.Printf("Node %s added to Plasma network.\n", node.NodeID)
	return nil
}


// RecordNodeRemoval logs the removal of a node from the Plasma network.
func (ledger *PlasmaLedger) RecordNodeRemoval(nodeID string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if _, exists := ledger.State.Nodes[nodeID]; exists {
		delete(ledger.State.Nodes, nodeID)
		fmt.Printf("Node %s removed from Plasma network.\n", nodeID)
		return nil
	}
	return fmt.Errorf("Node with ID %s not found.\n", nodeID)
}

// RecordBlockSync logs synchronization between the Plasma chain and the root chain.
func (ledger *PlasmaLedger) RecordBlockSync(rootChainTx RootChainTxRecord) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	ledger.State.RootChainTx = rootChainTx
	fmt.Printf("Plasma chain synced with root chain via Tx: %s.\n", rootChainTx.TxHash)
	return nil
}

// RecordSubBlockBroadcast handles the broadcasting of sub-blocks to other nodes.
func (ledger *PlasmaLedger) RecordSubBlockBroadcast(subBlock SubBlock) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Sub-block #%d broadcasted to Plasma network.\n", subBlock.Index)
	return nil
}

// RecordNodeValidation logs the validation activity of a node within the Plasma network.
func (ledger *PlasmaLedger) RecordNodeValidation(nodeID string, subBlock SubBlock) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Node %s validated Sub-block #%d.\n", nodeID, subBlock.Index)
	return nil
}
// RecordNetworkReconfiguration logs a network reconfiguration event.
func (ledger *PlasmaLedger) RecordNetworkReconfiguration(reconfig NetworkReconfig) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Plasma network reconfiguration: %s.\n", reconfig.Details)
	return nil
}



// FinalizeAndSubmitToRootChain submits the Plasma block to the root chain for security.
func (ledger *PlasmaLedger) FinalizeAndSubmitToRootChain(rootChainTx RootChainTxRecord) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	ledger.State.RootChainTx = rootChainTx
	fmt.Printf("Plasma block submitted to root chain with transaction: %s.\n", rootChainTx.TxHash)
	return nil
}