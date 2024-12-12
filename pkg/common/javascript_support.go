package common

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
	"github.com/dop251/goja"
)

// JavaScriptSupport manages JavaScript contract execution and validation within the Synnergy Network.
type JavaScriptSupport struct {
	LedgerInstance *ledger.Ledger             // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock         // List of sub-blocks pending full block aggregation
	mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// NewJavaScriptSupport initializes a new JavaScriptSupport instance.
func NewJavaScriptSupport(ledgerInstance *ledger.Ledger) *JavaScriptSupport {
	return &JavaScriptSupport{
		LedgerInstance: ledgerInstance,
		SubBlocks:      make([]*SubBlock, 0),
	}
}

// DeployAndExecuteContract deploys and executes a JavaScript smart contract on the blockchain.
func (js *JavaScriptSupport) DeployAndExecuteContract(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	// Step 1: Ensure encryption instance is correctly defined.
	encryptionInstance := &Encryption{} // Initialize or use an existing encryption instance

	// Step 2: Decrypt the contract bytecode.
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Execute the JavaScript bytecode.
	result, err := js.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Add the transaction result to sub-blocks.
	err = js.addTransactionToSubBlock(contractID, result)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction to sub-block: %v", err)
	}

	// Step 5: Log the contract deployment and execution to the ledger.
	// Provide both a message and the contract ID to LogEntry.
	err = js.LedgerInstance.VirtualMachineLedger.LogEntry(fmt.Sprintf("Contract %s deployed and executed at %s", contractID, time.Now().String()), contractID)
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	return result, nil
}


// ValidateSubBlock validates a sub-block and adds it to the list of sub-blocks waiting for full block aggregation.
func (js *JavaScriptSupport) ValidateSubBlock(subBlock *SubBlock) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	// Step 1: Check if the sub-block has already been validated.
	for _, sb := range js.SubBlocks {
		if sb.SubBlockID == subBlock.SubBlockID {
			return fmt.Errorf("sub-block %s is already validated", subBlock.SubBlockID)
		}
	}

	// Step 2: Add the validated sub-block to the pending sub-blocks list.
	js.SubBlocks = append(js.SubBlocks, subBlock)
	fmt.Printf("Sub-block %s validated and added for full block aggregation.\n", subBlock.SubBlockID)

	// Step 3: If 1000 sub-blocks have been validated, create a new block.
	if len(js.SubBlocks) == 1000 {
		err := js.aggregateSubBlocksIntoBlock()
		if err != nil {
			return fmt.Errorf("failed to aggregate sub-blocks into block: %v", err)
		}
	}

	return nil
}

// executeBytecode simulates executing a JavaScript smart contract bytecode on the Synnergy Network.
func (js *JavaScriptSupport) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Utilize a virtual machine or JavaScript engine to execute the bytecode.
	vm := NewJavaScriptInterpreter(js.LedgerInstance)
	result, err := vm.ExecuteBytecode(contractID, bytecode, parameters)
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}
	return result, nil
}

// addTransactionToSubBlock adds a contract execution result to a new or existing sub-block.
func (js *JavaScriptSupport) addTransactionToSubBlock(contractID, result string) error {
	// Create a new transaction.
	transaction := Transaction{
		TransactionID:   fmt.Sprintf("tx-%d", time.Now().UnixNano()), // Correct field name
		FromAddress:     "contract",  // Assuming FromAddress is the correct field name
		ToAddress:       "network",   // Assuming ToAddress is the correct field name
		Amount:          0,           // No monetary value for contract execution
		Timestamp:       time.Now(),
		ExecutionResult: result,      // Store the result of the contract execution
		SubBlockID:      "",          // Initially empty, will be populated later
	}

	// Create a new sub-block for the transaction.
	subBlock := SubBlock{
		SubBlockID:   fmt.Sprintf("sb-%d", time.Now().UnixNano()),  // Sub-block ID
		Transactions: []Transaction{transaction},                  // Transactions in this sub-block
		Timestamp:    time.Now(),                                  // Timestamp of sub-block creation
	}

	// Set the SubBlockID in the transaction after creating the sub-block.
	transaction.SubBlockID = subBlock.SubBlockID

	// Validate the sub-block.
	return js.ValidateSubBlock(&subBlock)
}

// aggregateSubBlocksIntoBlock aggregates 1000 validated sub-blocks into a new block.
func (js *JavaScriptSupport) aggregateSubBlocksIntoBlock() error {
	if len(js.SubBlocks) < 1000 {
		return fmt.Errorf("insufficient sub-blocks for block creation")
	}

	// Step 1: Convert []*SubBlock to []ledger.SubBlock (dereference pointers and convert transactions).
	var subBlocks []ledger.SubBlock
	for _, sb := range js.SubBlocks {
		// Convert each Transaction to ledger.Transaction
		var ledgerTransactions []ledger.Transaction
		for _, tx := range sb.Transactions {
			ledgerTransactions = append(ledgerTransactions, ledger.Transaction{
				TransactionID:     tx.TransactionID,
				FromAddress:       tx.FromAddress,
				ToAddress:         tx.ToAddress,
				Amount:            tx.Amount,
				Fee:               tx.Fee,
				TokenStandard:     tx.TokenStandard,
				TokenID:           tx.TokenID,
				Timestamp:         tx.Timestamp,
				SubBlockID:        tx.SubBlockID,
				BlockID:           tx.BlockID,
				ValidatorID:       tx.ValidatorID,
				Signature:         tx.Signature,
				Status:            tx.Status,
				EncryptedData:     tx.EncryptedData,
				DecryptedData:     tx.DecryptedData,
				ExecutionResult:   tx.ExecutionResult,
				FrozenAmount:      tx.FrozenAmount,
				RefundAmount:      tx.RefundAmount,
				ReversalRequested: tx.ReversalRequested,
			})
		}

		// Append to the ledger.SubBlock slice
		subBlocks = append(subBlocks, ledger.SubBlock{
			SubBlockID:   sb.SubBlockID,
			Transactions: ledgerTransactions, // Converted transactions
			Timestamp:    sb.Timestamp,
			PrevHash:     sb.PrevHash,
			Hash:         sb.Hash,
			Status:       sb.Status,
		})
	}

	// Step 2: Create the new block.
	newBlock := ledger.Block{
		BlockID:   fmt.Sprintf("block-%d", time.Now().UnixNano()),
		SubBlocks: subBlocks, // Use the converted []ledger.SubBlock
		Timestamp: time.Now(),
	}

	// Step 3: Store the block in the ledger.
	err := js.LedgerInstance.BlockchainConsensusCoinLedger.AddBlock(newBlock) // Directly store the Block in the ledger without encryption
	if err != nil {
		return fmt.Errorf("failed to record block in ledger: %v", err)
	}

	// Step 4: Clear the sub-block list after successful block creation.
	js.SubBlocks = make([]*SubBlock, 0) // Clear the sub-block list
	fmt.Printf("Block %s created and stored in the ledger.\n", newBlock.BlockID)

	return nil
}

// JavaScriptInterpreter represents a virtual machine that can interpret and execute JavaScript smart contract bytecode.
type JavaScriptInterpreter struct {
	LedgerInstance *ledger.Ledger // Reference to the ledger for storing transactions and logs
}

// NewJavaScriptInterpreter creates and returns a new JavaScript bytecode interpreter.
func NewJavaScriptInterpreter(ledger *ledger.Ledger) *JavaScriptInterpreter {
	return &JavaScriptInterpreter{
		LedgerInstance: ledger,
	}
}

// ExecuteBytecode executes JavaScript bytecode with provided parameters.
func (vm *JavaScriptInterpreter) ExecuteBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Step 1: Initialize the JavaScript engine.
	runtime := goja.New()

	// Step 2: Define contract parameters in the JavaScript runtime.
	for key, value := range parameters {
		runtime.Set(key, value)
	}

	// Step 3: Try to run the JavaScript bytecode.
	var executionResult string
	_, err := runtime.RunString(bytecode) // Capture both the result and the error
	if err != nil {
		return "", fmt.Errorf("execution failed for contract %s: %v", contractID, err)
	}

	// Step 4: Assume the execution result is stored in a variable named 'result' in the JavaScript context.
	// You can extract the 'result' variable from the runtime after the bytecode has been executed.
	if val := runtime.Get("result"); val != nil {
		executionResult = val.String()
	} else {
		executionResult = "No result returned from contract"
	}

	// Step 5: Log the execution details (optional).
	fmt.Printf("Executed JavaScript contract with ID: %s\n", contractID)
	fmt.Printf("Bytecode executed: %s\n", bytecode)
	fmt.Printf("Parameters: %+v\n", parameters)
	fmt.Printf("Execution Result: %s\n", executionResult)

	// Return the execution result with no error.
	return executionResult, nil
}

