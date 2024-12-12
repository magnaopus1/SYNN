package common

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// RustSupport manages Rust-based contract execution and validation within the Synnergy Network.
type RustSupport struct {
	LedgerInstance *ledger.Ledger             // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock         // List of sub-blocks pending full block aggregation
	mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// NewRustSupport initializes a new RustSupport instance.
func NewRustSupport(ledgerInstance *ledger.Ledger) *RustSupport {
	return &RustSupport{
		LedgerInstance: ledgerInstance,
		SubBlocks:      make([]*SubBlock, 0),
	}
}

// DeployAndExecuteContract deploys and executes a Rust smart contract on the blockchain.
func (rs *RustSupport) DeployAndExecuteContract(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// Step 1: Create an encryption instance (assuming Encryption is defined elsewhere in your project).
	encryptionInstance := &Encryption{} // Replace with the actual initialization of your encryption system.

	// Step 2: Decrypt the contract bytecode using the encryption instance.
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Execute the Rust bytecode.
	result, err := rs.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Add the transaction result to sub-blocks.
	err = rs.addTransactionToSubBlock(contractID, result)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction to sub-block: %v", err)
	}

	// Step 5: Log the contract deployment and execution to the ledger.
	err = rs.LedgerInstance.VirtualMachineLedger.LogEntry(fmt.Sprintf("Contract %s deployed and executed at %s", contractID, time.Now().String()), contractID)
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	return result, nil
}


// ValidateSubBlock validates a sub-block and adds it to the list of sub-blocks waiting for full block aggregation.
func (rs *RustSupport) ValidateSubBlock(subBlock *SubBlock) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// Step 1: Check if the sub-block has already been validated.
	for _, sb := range rs.SubBlocks {
		if sb.SubBlockID == subBlock.SubBlockID {
			return fmt.Errorf("sub-block %s is already validated", subBlock.SubBlockID)
		}
	}

	// Step 2: Add the validated sub-block to the pending sub-blocks list.
	rs.SubBlocks = append(rs.SubBlocks, subBlock)
	fmt.Printf("Sub-block %s validated and added for full block aggregation.\n", subBlock.SubBlockID)

	// Step 3: If 1000 sub-blocks have been validated, create a new block.
	if len(rs.SubBlocks) == 1000 {
		err := rs.aggregateSubBlocksIntoBlock()
		if err != nil {
			return fmt.Errorf("failed to aggregate sub-blocks into block: %v", err)
		}
	}

	return nil
}

// executeBytecode simulates executing a Rust smart contract bytecode on the Synnergy Network.
func (rs *RustSupport) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Utilize the virtual machine to execute the bytecode.
	vm := NewRustCompiler(rs.LedgerInstance)
	result, err := vm.executeBytecode(contractID, bytecode, parameters) // Change to lowercase "executeBytecode"
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}
	return result, nil
}


// addTransactionToSubBlock adds a contract execution result to a new or existing sub-block.
func (rs *RustSupport) addTransactionToSubBlock(contractID, result string) error {
	// Create a new transaction.
	transaction := Transaction{
		TransactionID:   fmt.Sprintf("tx-%d", time.Now().UnixNano()), // Fixed field name
		FromAddress:     "contract",
		ToAddress:       "network",
		Amount:          0, // No monetary value for contract execution.
		Timestamp:       time.Now(),
		ExecutionResult: result, // Use the result from the contract execution
	}

	// Create a new sub-block for the transaction.
	subBlock := SubBlock{
		SubBlockID:   fmt.Sprintf("sb-%d", time.Now().UnixNano()),
		Transactions: []Transaction{transaction},
		Timestamp:    time.Now(),
	}

	// Validate the sub-block.
	return rs.ValidateSubBlock(&subBlock)
}

// aggregateSubBlocksIntoBlock aggregates 1000 validated sub-blocks into a new block.
func (rs *RustSupport) aggregateSubBlocksIntoBlock() error {
	if len(rs.SubBlocks) < 1000 {
		return fmt.Errorf("insufficient sub-blocks for block creation")
	}

	// Step 1: Convert []*SubBlock to []ledger.SubBlock (dereference pointers and convert type).
	var subBlocks []ledger.SubBlock
	for _, sb := range rs.SubBlocks {
		// Convert each SubBlock to ledger.SubBlock
		ledgerSubBlock := ledger.SubBlock{
			SubBlockID:   sb.SubBlockID,
			Transactions: convertTransactions(sb.Transactions), // Convert transactions as needed
			Timestamp:    sb.Timestamp,
			Validator:    sb.Validator,
			PoHProof:     convertPoHProof(sb.PoHProof), // Convert PoHProof as needed
			PrevHash:     sb.PrevHash,
			Hash:         sb.Hash,
			Status:       sb.Status,
		}
		subBlocks = append(subBlocks, ledgerSubBlock)
	}

	// Step 2: Create the new block.
	newBlock := ledger.Block{
		BlockID:   fmt.Sprintf("block-%d", time.Now().UnixNano()),
		SubBlocks: subBlocks, // Now using converted []ledger.SubBlock
		Timestamp: time.Now(),
	}

	// Step 3: Add the block to the ledger (without using encryptedBlock).
	err := rs.LedgerInstance.BlockchainConsensusCoinLedger.AddBlock(newBlock) // Add the actual ledger.Block structure
	if err != nil {
		return fmt.Errorf("failed to record block in ledger: %v", err)
	}

	// Step 4: Clear the sub-block list after successful block creation.
	rs.SubBlocks = make([]*SubBlock, 0)
	fmt.Printf("Block %s created and stored in the ledger.\n", newBlock.BlockID)

	return nil
}


// Helper function to convert transactions to ledger.Transaction
func convertTransactions(txs []Transaction) []ledger.Transaction {
	var ledgerTxs []ledger.Transaction
	for _, tx := range txs {
		ledgerTx := ledger.Transaction{
			TransactionID:   tx.TransactionID,
			FromAddress:     tx.FromAddress,
			ToAddress:       tx.ToAddress,
			Amount:          tx.Amount,
			Fee:             tx.Fee,
			TokenStandard:   tx.TokenStandard,
			TokenID:         tx.TokenID,
			Timestamp:       tx.Timestamp,
			SubBlockID:      tx.SubBlockID,
			BlockID:         tx.BlockID,
			ValidatorID:     tx.ValidatorID,
			Signature:       tx.Signature,
			Status:          tx.Status,
			EncryptedData:   tx.EncryptedData,
			DecryptedData:   tx.DecryptedData,
			ExecutionResult: tx.ExecutionResult,
			FrozenAmount:    tx.FrozenAmount,
			RefundAmount:    tx.RefundAmount,
			ReversalRequested: tx.ReversalRequested,
		}
		ledgerTxs = append(ledgerTxs, ledgerTx)
	}
	return ledgerTxs
}
// Helper function to convert PoHProof to ledger.PoHProof
func convertPoHProof(poh PoHProof) ledger.PoHProof {
	return ledger.PoHProof{
		Hash:     poh.Hash,        // Use existing field `Hash` from PoHProof
		Sequence: poh.Sequence,    // Use existing field `Sequence` from PoHProof
		Timestamp:     poh.Timestamp,   // Use existing field `Timestamp` from PoHProof as `Time`
	}
}
