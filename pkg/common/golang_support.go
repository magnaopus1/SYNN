package common

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/ledger"
)

// GoSupport manages Go contract execution and validation within the Synnergy Network's virtual machine.
type GoSupport struct {
	LedgerInstance *ledger.Ledger // Ledger for logging transactions and contract activities
	SubBlocks      []*SubBlock    // List of sub-blocks pending full block aggregation
	mutex          sync.Mutex     // Mutex for thread-safe operations
}

// NewGoSupport initializes a new GoSupport instance.
func NewGoSupport(ledgerInstance *ledger.Ledger) *GoSupport {
	return &GoSupport{
		LedgerInstance: ledgerInstance,
		SubBlocks:      make([]*SubBlock, 0),
	}
}

// DeployAndExecuteContract deploys and executes a Go smart contract on the blockchain.
func (gs *GoSupport) DeployAndExecuteContract(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// Step 1: Create an encryption instance and decrypt the contract bytecode.
	encryptionInstance := &Encryption{}
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 2: Execute the bytecode.
	result, err := gs.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 3: Convert result string to map[string]interface{}.
	var resultMap map[string]interface{}
	err = json.Unmarshal([]byte(result), &resultMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal execution result: %v", err)
	}

	// Step 4: Add transaction to sub-blocks.
	err = gs.addTransactionToSubBlock(contractID, resultMap)
	if err != nil {
		return "", fmt.Errorf("failed to add transaction to sub-block: %v", err)
	}

	// Step 5: Log the contract deployment and execution to the ledger.
	err = gs.LedgerInstance.VirtualMachineLedger.LogEntry(fmt.Sprintf("Contract %s deployed and executed at %s", contractID, time.Now().String()), contractID)
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	// Step 6: Return the result string as is.
	return result, nil
}

// ValidateSubBlock validates a sub-block and adds it to the list of sub-blocks waiting for full block validation.
func (gs *GoSupport) ValidateSubBlock(subBlock *SubBlock) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	// Step 1: Check if the sub-block has already been validated.
	for _, sb := range gs.SubBlocks {
		if sb.SubBlockID == subBlock.SubBlockID {
			return fmt.Errorf("sub-block %s is already validated", subBlock.SubBlockID)
		}
	}

	// Step 2: Add the validated sub-block to the pending sub-blocks list.
	gs.SubBlocks = append(gs.SubBlocks, subBlock)
	fmt.Printf("Sub-block %s validated and added for full block aggregation.\n", subBlock.SubBlockID)

	// Step 3: If 1000 sub-blocks have been validated, create a new block.
	if len(gs.SubBlocks) == 1000 {
		err := gs.aggregateSubBlocksIntoBlock()
		if err != nil {
			return fmt.Errorf("failed to aggregate sub-blocks into block: %v", err)
		}
	}

	return nil
}

// ExecuteBytecode simulates executing a Go smart contract bytecode on the Synnergy Network.
func (gs *GoSupport) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Utilize the virtual machine to execute the bytecode.
	vm := NewBytecodeInterpreter(gs.LedgerInstance)

	// Create an encryption instance and pass the encryption key.
	encryptionInstance := &Encryption{}
	encryptionKey := []byte("your-encryption-key")

	// Execute the bytecode with all necessary arguments.
	result, err := vm.ExecuteBytecode(contractID, bytecode, parameters, encryptionInstance, encryptionKey)
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}

	// Convert the result (map[string]interface{}) to a string using json.Marshal.
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to convert result to string: %v", err)
	}

	return string(resultBytes), nil
}

// addTransactionToSubBlock adds a contract execution result to a new or existing sub-block.
func (gs *GoSupport) addTransactionToSubBlock(contractID string, result map[string]interface{}) error {
	// Convert the result map to a string (e.g., JSON format).
	resultString, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to convert result to string: %v", err)
	}

	// Create a new transaction.
	transaction := Transaction{
		TransactionID:   fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		FromAddress:     "contract",
		ToAddress:       "network",
		Amount:          0, // No monetary value for contract execution.
		Timestamp:       time.Now(),
		ExecutionResult: string(resultString), // Store the result in ExecutionResult.
	}

	// Create a new sub-block for the transaction.
	subBlock := SubBlock{
		SubBlockID:   fmt.Sprintf("sb-%d", time.Now().UnixNano()),
		Transactions: []Transaction{transaction},
		Timestamp:    time.Now(),
	}

	// Validate the sub-block.
	return gs.ValidateSubBlock(&subBlock)
}

// aggregateSubBlocksIntoBlock aggregates 1000 validated sub-blocks into a new block.
func (gs *GoSupport) aggregateSubBlocksIntoBlock() error {
	if len(gs.SubBlocks) < 1000 {
		return fmt.Errorf("insufficient sub-blocks for block creation")
	}

	// Step 1: Convert []*SubBlock to []SubBlock (slice of pointers to slice of values).
	var subBlocks []SubBlock
	for _, subBlock := range gs.SubBlocks {
		subBlocks = append(subBlocks, *subBlock) // Dereference pointer to get value
	}

	// Step 2: Create the new block.
	newBlock := Block{
		BlockID:    fmt.Sprintf("block-%d", time.Now().UnixNano()),
		SubBlocks:  subBlocks, // Pass the converted []SubBlock
		Timestamp:  time.Now(),
		Validators: []string{"validator-1"}, // Use Synnergy Consensus to determine the validator.
	}

	// Step 3: Encrypt the new block (Optional).
	encryptionInstance := &Encryption{} // Create an instance of Encryption
	iv := []byte("random-iv-16bytes")   // Ensure this is 16 bytes or adjust based on encryption needs
	_, err := encryptionInstance.EncryptData(string(EncryptionKey), []byte(fmt.Sprintf("%+v", newBlock)), iv)
	if err != nil {
		return fmt.Errorf("failed to encrypt block: %v", err)
	}

	// Step 4: Convert newBlock to ledger.Block and store it in the ledger.
	ledgerBlock := convertBlockToLedgerBlock(newBlock) // Convert the custom block to ledger.Block
	err = gs.LedgerInstance.BlockchainConsensusCoinLedger.AddBlock(ledgerBlock)      // Pass the ledgerBlock
	if err != nil {
		return fmt.Errorf("failed to record block in ledger: %v", err)
	}

	// Step 5: Clear the sub-block list after successful block creation.
	gs.SubBlocks = make([]*SubBlock, 0)
	fmt.Printf("Block %s created and stored in the ledger.\n", newBlock.BlockID)

	return nil
}
