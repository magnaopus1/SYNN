package common

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// RustCompiler manages the compilation, execution, and deployment of Rust-based smart contracts.
type RustCompiler struct {
	LedgerInstance *ledger.Ledger          // Ledger instance for logging contract activities
	CompiledCode   map[string]string       // Stores compiled bytecode for each contract
	mutex          sync.Mutex              // Mutex for thread-safe operations
}

// NewRustCompiler initializes a new RustCompiler instance.
func NewRustCompiler(ledgerInstance *ledger.Ledger) *RustCompiler {
	return &RustCompiler{
		LedgerInstance: ledgerInstance,
		CompiledCode:   make(map[string]string),
	}
}

// CompileRustContract compiles the Rust contract and stores the resulting bytecode.
func (rc *RustCompiler) CompileRustContract(contractID, contractSourcePath string) error {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Step 1: Execute the Rust compiler (cargo) to compile the contract.
	compiledBytecode, err := rc.runRustCompiler(contractSourcePath)
	if err != nil {
		return fmt.Errorf("compilation failed: %v", err)
	}

	// Step 2: Create an encryption instance and generate an initialization vector (IV).
	encryptionInstance := &Encryption{}  // Create an instance of Encryption
	iv := []byte("random-iv-16bytes")    // Ensure this is 16 bytes or adjust according to your encryption scheme

	// Step 3: Encrypt the compiled bytecode using the encryption key and IV.
	encryptedBytecode, err := encryptionInstance.EncryptData(string(EncryptionKey), []byte(compiledBytecode), iv) // Convert EncryptionKey to string
	if err != nil {
		return fmt.Errorf("failed to encrypt compiled bytecode: %v", err)
	}

	// Step 4: Store the encrypted bytecode.
	rc.CompiledCode[contractID] = string(encryptedBytecode)

	// Step 5: Log the successful compilation in the ledger.
	logEntry := fmt.Sprintf("Contract %s compiled successfully at %s", contractID, time.Now().String())
	err = rc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Pass both the log message and the contract ID
	if err != nil {
		return fmt.Errorf("failed to log contract compilation: %v", err)
	}

	fmt.Printf("Contract %s successfully compiled and encrypted.\n", contractID)
	return nil
}

// DeployContract deploys the compiled contract bytecode into the blockchain.
func (rc *RustCompiler) DeployContract(contractID string, parameters map[string]interface{}) (string, error) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Step 1: Retrieve the compiled bytecode.
	encryptedBytecode, exists := rc.CompiledCode[contractID]
	if !exists {
		return "", fmt.Errorf("compiled bytecode for contract %s not found", contractID)
	}

	// Step 2: Create an encryption instance and decrypt the bytecode.
	encryptionInstance := &Encryption{} // Create an instance of Encryption
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(encryptedBytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Validate and execute the bytecode (part of the Synnergy Consensus mechanism).
	executionResult, err := rc.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Log the contract deployment in the ledger.
	logEntry := fmt.Sprintf("Contract %s deployed successfully at %s", contractID, time.Now().String())
	err = rc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Corrected: Pass both the log message and contractID
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	fmt.Printf("Contract %s successfully deployed.\n", contractID)
	return executionResult, nil
}

// runRustCompiler compiles the Rust smart contract using cargo.
func (rc *RustCompiler) runRustCompiler(sourcePath string) (string, error) {
	// Execute the Rust compiler (cargo) to generate the bytecode.
	cmd := exec.Command("cargo", "build", "--release", "--manifest-path", sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Rust compilation error: %v\nOutput: %s", err, string(output))
	}

	compiledBytecode := strings.TrimSpace(string(output))
	if compiledBytecode == "" {
		return "", fmt.Errorf("no bytecode generated")
	}

	return compiledBytecode, nil
}

// executeBytecode simulates the execution of a smart contract's bytecode.
func (rc *RustCompiler) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	// Step 1: Create an encryption instance to pass to the VM.
	encryptionInstance := &Encryption{} // Create an instance of Encryption

	// Step 2: The bytecode execution will interface with the blockchain's virtual machine and Synnergy Consensus.
	vm := NewBytecodeInterpreter(rc.LedgerInstance)

	// Step 3: Execute the bytecode and pass all required arguments including the encryption instance and key.
	executionResult, err := vm.ExecuteBytecode(contractID, bytecode, parameters, encryptionInstance, EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}

	return fmt.Sprintf("Contract %s executed successfully: %v", contractID, executionResult), nil
}
