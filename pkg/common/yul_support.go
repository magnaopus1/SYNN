package common

import (
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/ledger"
)

// YulSupport manages Yul contract execution and validation within the Synnergy Network.
type YulSupport struct {
	LedgerInstance *ledger.Ledger // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock    // List of sub-blocks pending full block aggregation
	Blocks         []Block        // List of aggregated blocks
	mutex          sync.Mutex     // Mutex for thread-safe operations
}

// YulInterpreter represents a virtual machine that can interpret and execute Yul bytecode.
type YulInterpreter struct {
	LedgerInstance *ledger.Ledger // Update this to use *ledger.Ledger
}

// NewYulInterpreter creates and returns a new Yul bytecode interpreter.
func NewYulInterpreter(ledger *ledger.Ledger) *YulInterpreter { // Update the argument type
	return &YulInterpreter{
		LedgerInstance: ledger,
	}
}


// ExecuteBytecode executes the Yul bytecode within the YulInterpreter.
func (vm *YulInterpreter) ExecuteBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Simulate bytecode execution logic here.
	// In a real-world scenario, this would involve interpreting the Yul bytecode and executing the logic.
	fmt.Printf("Executing Yul bytecode for contractID: %s with parameters: %v\n", contractID, parameters)

	// For now, simulate a successful execution result.
	executionResult := "Yul contract execution successful"
	return executionResult, nil
}

// executeBytecode simulates executing a Yul smart contract bytecode on the Synnergy Network.
func (ys *YulSupport) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// No need to cast or convert types, as NewYulInterpreter now accepts *ledger.Ledger
	vm := NewYulInterpreter(ys.LedgerInstance)
	result, err := vm.ExecuteBytecode(contractID, bytecode, parameters)
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}
	return result, nil
}






// NewYulSupport initializes a new YulSupport instance.
func NewYulSupport(ledgerInstance *ledger.Ledger) *YulSupport {
	return &YulSupport{
		LedgerInstance: ledgerInstance,
		SubBlocks:      make([]*SubBlock, 0),
	}
}

// DeployAndExecuteContract deploys and executes a Yul smart contract on the blockchain.
func (ys *YulSupport) DeployAndExecuteContract(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	ys.mutex.Lock()
	defer ys.mutex.Unlock()

	// Step 1: Create an instance of the Encryption struct.
	encryptionInstance := &Encryption{}

	// Step 2: Decrypt the contract bytecode.
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Execute the Yul bytecode.
	result, err := ys.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Add the transaction result to sub-blocks.
	err = ys.addTransactionToSubBlock(contractID, result)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction to sub-block: %v", err)
	}

	// Step 5: Log the contract deployment and execution to the ledger.
	logEntry := fmt.Sprintf("Contract %s deployed and executed at %s", contractID, time.Now().String())
	err = ys.LedgerInstance.VirtualMachineLedger.LogEntry(contractID, logEntry) // Pass two arguments: contractID and log entry
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	return result, nil
}


// ValidateSubBlock validates a sub-block and adds it to the list of sub-blocks waiting for full block aggregation.
func (ys *YulSupport) ValidateSubBlock(subBlock *SubBlock) error {
	ys.mutex.Lock()
	defer ys.mutex.Unlock()

	// Step 1: Check if the sub-block has already been validated.
	for _, sb := range ys.SubBlocks {
		if sb.SubBlockID == subBlock.SubBlockID {
			return fmt.Errorf("sub-block %s is already validated", subBlock.SubBlockID)
		}
	}

	// Step 2: Add the validated sub-block to the pending sub-blocks list.
	ys.SubBlocks = append(ys.SubBlocks, subBlock)
	fmt.Printf("Sub-block %s validated and added for full block aggregation.\n", subBlock.SubBlockID)

	// Step 3: If 1000 sub-blocks have been validated, create a new block.
	if len(ys.SubBlocks) == 1000 {
		err := ys.aggregateSubBlocksIntoBlock()
		if err != nil {
			return fmt.Errorf("failed to aggregate sub-blocks into block: %v", err)
		}
	}

	return nil
}




// addTransactionToSubBlock adds a contract execution result to a new or existing sub-block.
func (ys *YulSupport) addTransactionToSubBlock(contractID, result string) error {
    // Create a new transaction.
    transaction := Transaction{
        TransactionID:   fmt.Sprintf("tx-%d", time.Now().UnixNano()),
        FromAddress:     "contract",
        ToAddress:       "network",
        Amount:          0, // No monetary value for contract execution.
        Timestamp:       time.Now(),
        ExecutionResult: result, // Add the execution result to the transaction.
        Status:          "pending",
    }

    // Create a new sub-block for the transaction.
    subBlock := SubBlock{
        Index:         len(ys.SubBlocks) + 1, // Increment sub-block index.
        Transactions:  []Transaction{transaction},
        Timestamp:     time.Now(),
    }

    // Validate the sub-block.
    return ys.ValidateSubBlock(&subBlock)
}

// aggregateSubBlocksIntoBlock aggregates 1000 validated sub-blocks into a new block.
func (ys *YulSupport) aggregateSubBlocksIntoBlock() error {
    if len(ys.SubBlocks) < 1000 {
        return fmt.Errorf("insufficient sub-blocks for block creation")
    }

    // Step 1: Convert []*SubBlock to []SubBlock.
    subBlocks := make([]SubBlock, len(ys.SubBlocks))
    for i, sb := range ys.SubBlocks {
        subBlocks[i] = *sb // Dereference the pointer to get the SubBlock value.
    }

    // Step 2: Gather validators.
    validatorMap := make(map[string]bool) // To avoid duplicates
    for _, sb := range subBlocks {
        validatorMap[sb.Validator] = true
    }

    // Convert validatorMap keys to a slice
    validators := make([]string, 0, len(validatorMap))
    for v := range validatorMap {
        validators = append(validators, v)
    }

    // Step 3: Create the new block (in your local format).
    newBlock := Block{
        BlockID:    fmt.Sprintf("block-%d", time.Now().UnixNano()), // Unique Block ID
        Index:      len(ys.Blocks) + 1,                             // Block index
        SubBlocks:  subBlocks,                                      // Use the converted sub-blocks
        Timestamp:  time.Now(),
        PrevHash:   "previous-block-hash",                          // Add logic to get the previous block hash
        Difficulty: 2,                                              // Set the difficulty for PoW
        Validators: validators,                                     // Set the validators
    }

    // Step 4: Convert newBlock to ledger.Block type before adding it to the ledger.
    ledgerBlock := ConvertToLedgerBlock(newBlock)

    // Step 5: Store the block in the ledger (no need for BlockID or encrypted data here).
    err := ys.LedgerInstance.BlockchainConsensusCoinLedger.AddBlock(ledgerBlock) // Pass the converted ledger block directly
    if err != nil {
        return fmt.Errorf("failed to record block in ledger: %v", err)
    }

    // Step 6: Clear the sub-block list after successful block creation.
    ys.SubBlocks = make([]*SubBlock, 0) // Reset the list
    fmt.Printf("Block %s created and stored in the ledger.\n", newBlock.BlockID)

    return nil
}
