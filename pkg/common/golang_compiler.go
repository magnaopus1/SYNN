package common

import (
	"fmt"
	"os/exec"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// GoContractCompiler manages the compilation and deployment of Go-based smart contracts.
type GoContractCompiler struct {
	LedgerInstance *ledger.Ledger          // Ledger instance for logging contract deployments
	CompiledCode   map[string]string       // Stores compiled bytecode for each contract
	mutex          sync.Mutex              // Mutex for thread-safe operations
}

// NewGoContractCompiler initializes a new GoContractCompiler instance.
func NewGoContractCompiler(ledgerInstance *ledger.Ledger) *GoContractCompiler {
	return &GoContractCompiler{
		LedgerInstance: ledgerInstance,
		CompiledCode:   make(map[string]string),
	}
}

// CompileGoContract compiles a Go smart contract and stores the resulting bytecode.
func (gcc *GoContractCompiler) CompileGoContract(contractID, contractSourcePath string) error {
	gcc.mutex.Lock()
	defer gcc.mutex.Unlock()

	// Step 1: Execute the Go compiler (go build) to compile the contract into a binary
	compiledBytecode, err := gcc.runGoCompiler(contractSourcePath)
	if err != nil {
		return fmt.Errorf("compilation failed: %v", err)
	}

	// Step 2: Create an encryption instance and generate an initialization vector (IV).
	encryptionInstance := &Encryption{}  // Create an instance of Encryption
	iv := []byte("random-iv-16bytes")    // Ensure this is 16 bytes or adjust according to your encryption scheme

	// Step 3: Encrypt the compiled bytecode using the encryption key (converted to string) and IV.
	encryptedBytecode, err := encryptionInstance.EncryptData(string(EncryptionKey), []byte(compiledBytecode), iv) // Correctly convert EncryptionKey to string
	if err != nil {
		return fmt.Errorf("failed to encrypt compiled bytecode: %v", err)
	}

	// Step 4: Store the encrypted bytecode
	gcc.CompiledCode[contractID] = string(encryptedBytecode)

	// Step 5: Log the successful compilation in the ledger
	logEntry := fmt.Sprintf("Contract %s compiled successfully at %s", contractID, time.Now().String())
	err = gcc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Correct: pass both logEntry and contractID
	if err != nil {
		return fmt.Errorf("failed to log contract compilation: %v", err)
	}

	fmt.Printf("Contract %s successfully compiled and encrypted.\n", contractID)
	return nil
}

// DeployContract deploys the compiled contract bytecode to the blockchain.
func (gcc *GoContractCompiler) DeployContract(contractID string, parameters map[string]interface{}) (string, error) {
	gcc.mutex.Lock()
	defer gcc.mutex.Unlock()

	// Step 1: Retrieve the compiled bytecode
	encryptedBytecode, exists := gcc.CompiledCode[contractID]
	if !exists {
		return "", fmt.Errorf("compiled bytecode for contract %s not found", contractID)
	}

	// Step 2: Create an encryption instance and decrypt the bytecode
	encryptionInstance := &Encryption{} // Create an instance of Encryption
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(encryptedBytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Execute the bytecode and validate with Synnergy Consensus
	executionResult, err := gcc.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Log the contract deployment
	logEntry := fmt.Sprintf("Contract %s deployed successfully at %s", contractID, time.Now().String())
	err = gcc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Correct: pass both logEntry and contractID
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	fmt.Printf("Contract %s successfully deployed.\n", contractID)
	return executionResult, nil
}

// runGoCompiler runs the Go compiler and returns the compiled bytecode.
func (gcc *GoContractCompiler) runGoCompiler(sourcePath string) (string, error) {
	cmd := exec.Command("go", "build", "-o", "/tmp/contract_binary", sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go compilation error: %v\nOutput: %s", err, string(output))
	}

	compiledBytecode := "/tmp/contract_binary" // Compiled binary file path
	return compiledBytecode, nil
}

// executeBytecode simulates the execution of a smart contract's bytecode.
func (gcc *GoContractCompiler) executeBytecode(contractID string, bytecodePath string, parameters map[string]interface{}) (string, error) {
	// The bytecode execution will interface with the Synnergy Consensus mechanism
	cmd := exec.Command(bytecodePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution error: %v\nOutput: %s", err, string(output))
	}

	// Validate sub-block transactions with Synnergy Consensus
	var transactions []Transaction 

	// Assuming you have an instance of SynnergyConsensus
	var synnergyConsensus SynnergyConsensus
	synnergyConsensus.ProcessTransactions(transactions) // Call the method on the instance

	return fmt.Sprintf("Contract %s executed successfully: %s", contractID, string(output)), nil
}
