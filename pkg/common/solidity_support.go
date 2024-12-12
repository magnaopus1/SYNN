package common

import (
	"encoding/json"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// SoliditySupport manages Solidity contract execution and validation within the virtual machine.
type SoliditySupport struct {
	LedgerInstance *ledger.Ledger // Ledger for logging transactions and contract activities
	mutex          sync.Mutex     // Mutex for thread-safe operations
	SubBlocks      []*SubBlock     // List of sub-blocks pending full block aggregation

}

// NewSoliditySupport initializes a new SoliditySupport instance.
func NewSoliditySupport(ledgerInstance *ledger.Ledger) *SoliditySupport {
	return &SoliditySupport{
		LedgerInstance: ledgerInstance,
	}
}

// DeployAndExecuteContract deploys and executes a Solidity smart contract on the blockchain.
func (ss *SoliditySupport) DeployAndExecuteContract(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Step 1: Ensure encryption is correctly defined or imported.
	encryptionInstance := &Encryption{} // Assume Encryption struct or package exists and is imported.

	// Step 2: Decrypt the contract bytecode using the encryption instance.
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), []byte("iv-placeholder")) // Corrected to only pass two arguments
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Execute the bytecode.
	result, err := ss.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Add transaction to sub-blocks.
	err = ss.addTransactionToSubBlock(contractID, result)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction to sub-block: %v", err)
	}

	// Step 5: Log the contract deployment and execution to the ledger.
	err = ss.LedgerInstance.VirtualMachineLedger.LogEntry(fmt.Sprintf("Contract %s deployed and executed at %s", contractID, time.Now().String()), contractID)
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	return result, nil
}




// ValidateSubBlock validates a sub-block and adds it to the list of sub-blocks waiting for full block validation.
func (ss *SoliditySupport) ValidateSubBlock(subBlock *SubBlock) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Step 1: Check if the sub-block has already been validated.
	for _, sb := range ss.SubBlocks {
		if sb.SubBlockID == subBlock.SubBlockID {
			return fmt.Errorf("sub-block %s is already validated", subBlock.SubBlockID)
		}
	}

	// Step 2: Add the validated sub-block to the pending sub-blocks list.
	ss.SubBlocks = append(ss.SubBlocks, subBlock)
	fmt.Printf("Sub-block %s validated and added for full block aggregation.\n", subBlock.SubBlockID)

	// Step 3: If 1000 sub-blocks have been validated, create a new block.
	if len(ss.SubBlocks) == 1000 {
		err := ss.aggregateSubBlocksIntoBlock()
		if err != nil {
			return fmt.Errorf("failed to aggregate sub-blocks into block: %v", err)
		}
	}

	return nil
}

// executeBytecode simulates executing a smart contract bytecode on the Synnergy Network.
func (ss *SoliditySupport) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Utilize the virtual machine to execute the bytecode.
	vm := NewBytecodeInterpreter(ss.LedgerInstance)

	// Initialize the encryption instance and initialization vector.
	encryptionInstance := &Encryption{}
	initializationVector := []byte("randomIV") // Replace with actual IV as needed.

	// Execute the bytecode with all necessary arguments.
	result, err := vm.ExecuteBytecode(contractID, bytecode, parameters, encryptionInstance, initializationVector)
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}

	// Convert the result (map) to a JSON string before returning.
	resultString, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to convert result to string: %v", err)
	}

	return string(resultString), nil
}


// addTransactionToSubBlock adds a contract execution result to a new or existing sub-block.
func (ss *SoliditySupport) addTransactionToSubBlock(contractID, result string) error {
	// Create a new transaction.
	transaction := Transaction{
		TransactionID:   fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		FromAddress:     "contract",
		ToAddress:       "network",
		Amount:          0, // No monetary value for contract execution.
		Timestamp:       time.Now(),
		ExecutionResult: result, // Store the contract execution result.
		Status:          "pending", // Set initial status as pending.
	}

	// Create a new sub-block for the transaction.
	subBlock := SubBlock{
		SubBlockID:    fmt.Sprintf("sb-%d", time.Now().UnixNano()),
		Transactions:  []Transaction{transaction}, // Add the transaction to the sub-block.
		Timestamp:     time.Now(),
	}

	// Validate the sub-block.
	return ss.ValidateSubBlock(&subBlock)
}

// aggregateSubBlocksIntoBlock aggregates 1000 validated sub-blocks into a new block.
func (ss *SoliditySupport) aggregateSubBlocksIntoBlock() error {
	if len(ss.SubBlocks) < 1000 {
		return fmt.Errorf("insufficient sub-blocks for block creation")
	}

	// Step 1: Convert []*SubBlock to []ledger.SubBlock (dereference pointers and convert type).
	var subBlocks []ledger.SubBlock
	for _, sb := range ss.SubBlocks {
		// Convert []Transaction to []ledger.Transaction
		var ledgerTransactions []ledger.Transaction
		for _, tx := range sb.Transactions {
			ledgerTransaction := ledger.Transaction{
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
			ledgerTransactions = append(ledgerTransactions, ledgerTransaction)
		}

		// Convert PoHProof to ledger.PoHProof
		ledgerPoHProof := ledger.PoHProof{
			Sequence:  sb.PoHProof.Sequence,  // Ensure this matches the field types
			Timestamp: sb.PoHProof.Timestamp, // Mapping the time field correctly
			Hash:      sb.PoHProof.Hash,      // Mapping the hash field
		}

		// Assuming ledger.SubBlock and SubBlock structures are different, adjust accordingly
		ledgerSubBlock := ledger.SubBlock{
			SubBlockID:   sb.SubBlockID,
			Transactions: ledgerTransactions, // Converted transactions
			Timestamp:    sb.Timestamp,
			Validator:    sb.Validator,
			PoHProof:     ledgerPoHProof,     // Converted PoHProof
			PrevHash:     sb.PrevHash,
			Hash:         sb.Hash,
			Status:       sb.Status,
		}
		subBlocks = append(subBlocks, ledgerSubBlock)
	}

	// Step 2: Create the new block.
	ledgerBlock := ledger.Block{
		BlockID:   fmt.Sprintf("block-%d", time.Now().UnixNano()),
		SubBlocks: subBlocks, // Now using converted []ledger.SubBlock
		Timestamp: time.Now(),
	}

	// Step 3: Store the block in the ledger.
	err := ss.LedgerInstance.BlockchainConsensusCoinLedger.AddBlock(ledgerBlock) // Ensure ledger.Block is the expected type.
	if err != nil {
		return fmt.Errorf("failed to record block in ledger: %v", err)
	}

	// Step 4: Clear the sub-block list after successful block creation.
	ss.SubBlocks = make([]*SubBlock, 0)
	fmt.Printf("Block %s created and stored in the ledger.\n", ledgerBlock.BlockID)

	return nil
}
