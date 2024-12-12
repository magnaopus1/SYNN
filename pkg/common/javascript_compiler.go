package common

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// JavaScriptCompiler manages the compilation, execution, and deployment of JavaScript smart contracts.
type JavaScriptCompiler struct {
	LedgerInstance *ledger.Ledger    // Ledger instance for logging contract activities
	CompiledCode   map[string]string // Stores compiled bytecode for each contract
	mutex          sync.Mutex        // Mutex for thread-safe operations
}

// NewJavaScriptCompiler initializes a new JavaScriptCompiler instance.
func NewJavaScriptCompiler(ledgerInstance *ledger.Ledger) *JavaScriptCompiler {
	return &JavaScriptCompiler{
		LedgerInstance: ledgerInstance,
		CompiledCode:   make(map[string]string),
	}
}

// CompileJavaScriptContract compiles the JavaScript contract and stores the resulting bytecode.
func (jc *JavaScriptCompiler) CompileJavaScriptContract(contractID, contractSourcePath string) error {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	// Step 1: Execute the JavaScript compiler (using Node.js) to compile the contract
	compiledBytecode, err := jc.runJavaScriptCompiler(contractSourcePath)
	if err != nil {
		return fmt.Errorf("compilation failed: %v", err)
	}

	// Step 2: Ensure encryption instance is available and encrypt the compiled bytecode
	encryptionInstance := &Encryption{} // Assuming Encryption struct is available
	iv := []byte("random-iv-16bytes")   // Example IV (should be 16 bytes for encryption algorithms)

	encryptedBytecode, err := encryptionInstance.EncryptData("encryption-key", []byte(compiledBytecode), iv)
	if err != nil {
		return fmt.Errorf("failed to encrypt compiled bytecode: %v", err)
	}

	// Step 3: Store the encrypted bytecode
	jc.CompiledCode[contractID] = string(encryptedBytecode)

	// Step 4: Log the successful compilation in the ledger
	logEntry := fmt.Sprintf("Contract %s compiled successfully at %s", contractID, time.Now().String())
	err = jc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Pass both the log entry and contract ID
	if err != nil {
		return fmt.Errorf("failed to log contract compilation: %v", err)
	}

	fmt.Printf("Contract %s successfully compiled and encrypted.\n", contractID)
	return nil
}

// DeployContract deploys the compiled contract bytecode into the blockchain.
func (jc *JavaScriptCompiler) DeployContract(contractID string, parameters map[string]interface{}) (string, error) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	// Step 1: Retrieve the compiled bytecode
	encryptedBytecode, exists := jc.CompiledCode[contractID]
	if !exists {
		return "", fmt.Errorf("compiled bytecode for contract %s not found", contractID)
	}

	// Step 2: Decrypt the bytecode using the encryption instance.
	encryptionInstance := &Encryption{} // Assuming Encryption struct is available
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(encryptedBytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Validate and execute the bytecode (Synnergy Consensus)
	executionResult, err := jc.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Log the contract deployment
	logEntry := fmt.Sprintf("Contract %s deployed successfully at %s", contractID, time.Now().String())
	err = jc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Corrected to pass both arguments
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	fmt.Printf("Contract %s successfully deployed.\n", contractID)
	return executionResult, nil
}

// runJavaScriptCompiler runs the JavaScript compiler (Node.js) and returns the compiled bytecode.
func (jc *JavaScriptCompiler) runJavaScriptCompiler(sourcePath string) (string, error) {
	// Execute Node.js to generate the bytecode
	cmd := exec.Command("node", sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("JavaScript compilation error: %v\nOutput: %s", err, string(output))
	}

	compiledBytecode := strings.TrimSpace(string(output))
	if compiledBytecode == "" {
		return "", fmt.Errorf("no bytecode generated")
	}

	return compiledBytecode, nil
}

// executeBytecode simulates the execution of a JavaScript smart contract's bytecode.
func (jc *JavaScriptCompiler) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Use Synnergy Consensus for bytecode execution
	vm := NewJavaScriptInterpreter(jc.LedgerInstance)
	executionResult, err := vm.ExecuteBytecode(contractID, bytecode, parameters)
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}

	return fmt.Sprintf("Contract %s executed successfully: %v", contractID, executionResult), nil
}

// ValidateSubBlock validates a sub-block and adds it to the list of validated sub-blocks.
func (jc *JavaScriptCompiler) ValidateSubBlock(subBlock *SubBlock) error {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	// Step 1: Check if sub-block is already validated
	for _, sb := range jc.LedgerInstance.BlockchainConsensusCoinLedger.GetSubBlocks() {
		if sb.SubBlockID == subBlock.SubBlockID {
			return fmt.Errorf("sub-block %s is already validated", subBlock.SubBlockID)
		}
	}

	// Step 2: Convert []Transaction to []ledger.Transaction
	var ledgerTransactions []ledger.Transaction
	for _, tx := range subBlock.Transactions {
		ledgerTx := ledger.Transaction{
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
		}
		ledgerTransactions = append(ledgerTransactions, ledgerTx)
	}

	// Step 3: Convert *SubBlock to ledger.SubBlock before adding it to the ledger.
	ledgerSubBlock := ledger.SubBlock{
		SubBlockID:   subBlock.SubBlockID,
		Transactions: ledgerTransactions, // Use the converted transactions
		Timestamp:    subBlock.Timestamp,
		PrevHash:     subBlock.PrevHash,
		Hash:         subBlock.Hash,
		Status:       subBlock.Status,
		Validator:    subBlock.Validator,
	}

	// Step 4: Add sub-block to ledger
	err := jc.LedgerInstance.BlockchainConsensusCoinLedger.AddSubBlock(ledgerSubBlock) // Pass the converted sub-block
	if err != nil {
		return fmt.Errorf("failed to record sub-block: %v", err)
	}

	fmt.Printf("Sub-block %s validated.\n", subBlock.SubBlockID)

	// Step 5: If 1000 sub-blocks validated, aggregate them into a block.
	if len(jc.LedgerInstance.BlockchainConsensusCoinLedger.GetSubBlocks()) == 1000 {
		err := jc.aggregateSubBlocksIntoBlock()
		if err != nil {
			return fmt.Errorf("failed to aggregate sub-blocks into block: %v", err)
		}
	}

	return nil
}

// aggregateSubBlocksIntoBlock aggregates 1000 validated sub-blocks into a new block.
func (jc *JavaScriptCompiler) aggregateSubBlocksIntoBlock() error {
	subBlocks := jc.LedgerInstance.BlockchainConsensusCoinLedger.GetSubBlocks()
	if len(subBlocks) < 1000 {
		return fmt.Errorf("insufficient sub-blocks for block creation")
	}

	// Step 1: Create the new block.
	newBlock := ledger.Block{ // Assuming Block is part of the ledger package
		BlockID:   fmt.Sprintf("block-%d", time.Now().UnixNano()),
		SubBlocks: subBlocks, // Use the sub-blocks from the ledger
		Timestamp: time.Now(),
	}

	// Step 2: Ensure encryption instance is available.
	encryptionInstance := &Encryption{} // Initialize or use an existing encryption instance
	iv := []byte("random-iv-16bytes")   // Example IV (should be 16 bytes for encryption algorithms)

	// Step 3: Encrypt the block (optional, if encryption is required).
	_, err := encryptionInstance.EncryptData("encryption-key", []byte(fmt.Sprintf("%+v", newBlock)), iv)
	if err != nil {
		return fmt.Errorf("failed to encrypt block: %v", err)
	}

	// Step 4: Store the block in the ledger.
	err = jc.LedgerInstance.BlockchainConsensusCoinLedger.AddBlock(newBlock) // Pass the full block object, not BlockID and encryptedBlock
	if err != nil {
		return fmt.Errorf("failed to record block: %v", err)
	}

	// Step 5: Print success message.
	fmt.Printf("Block %s created and stored in the ledger.\n", newBlock.BlockID)

	return nil
}
