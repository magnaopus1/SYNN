package common

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// SolidityCompiler manages the compilation and deployment of Solidity smart contracts.
type SolidityCompiler struct {
	LedgerInstance *ledger.Ledger          // Ledger instance to log contract deployment
	CompiledCode   map[string]string       // Stores the compiled bytecode for each contract
	mutex          sync.Mutex              // Mutex for thread-safety during compilation and deployment
}

// NewSolidityCompiler initializes a new SolidityCompiler instance.
func NewSolidityCompiler(ledgerInstance *ledger.Ledger) *SolidityCompiler {
	return &SolidityCompiler{
		LedgerInstance: ledgerInstance,
		CompiledCode:   make(map[string]string),
	}
}

// CompileSolidityContract compiles the Solidity contract and stores the resulting bytecode.
func (sc *SolidityCompiler) CompileSolidityContract(contractID, contractSourcePath string) error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Step 1: Execute the Solidity compiler (solc) to compile the contract.
	compiledBytecode, err := sc.runSolidityCompiler(contractSourcePath)
	if err != nil {
		return fmt.Errorf("compilation failed: %v", err)
	}

	// Step 2: Create an encryption instance and define an encryption key and IV.
	encryptionInstance := &Encryption{}                // Assuming Encryption struct exists.
	encryptionKey := string(EncryptionKey)             // Convert EncryptionKey []byte to string if necessary.
	iv := []byte("random-iv-16bytes")                  // Define a 16-byte initialization vector (IV).

	// Step 3: Encrypt the compiled bytecode using the key and IV.
	encryptedBytecode, err := encryptionInstance.EncryptData(encryptionKey, []byte(compiledBytecode), iv) // Pass key, bytecode, and IV.
	if err != nil {
		return fmt.Errorf("failed to encrypt compiled bytecode: %v", err)
	}

	// Step 4: Store the encrypted bytecode.
	sc.CompiledCode[contractID] = string(encryptedBytecode)

	// Step 5: Log the successful compilation in the ledger.
	logEntry := fmt.Sprintf("Contract %s compiled successfully at %s", contractID, time.Now().String())
	err = sc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Corrected to pass both logEntry and contractID.
	if err != nil {
		return fmt.Errorf("failed to log contract compilation: %v", err)
	}

	fmt.Printf("Contract %s successfully compiled and encrypted.\n", contractID)
	return nil
}



// DeployContract deploys the compiled contract bytecode into the blockchain.
func (sc *SolidityCompiler) DeployContract(contractID string, parameters map[string]interface{}) (string, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Step 1: Retrieve the compiled bytecode.
	encryptedBytecode, exists := sc.CompiledCode[contractID]
	if !exists {
		return "", fmt.Errorf("compiled bytecode for contract %s not found", contractID)
	}

	// Step 2: Create an encryption instance and decrypt the bytecode.
	encryptionInstance := &Encryption{} // Assuming you have an Encryption struct or package.
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(encryptedBytecode), EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 3: Validate and execute the bytecode (part of the consensus mechanism).
	executionResult, err := sc.executeBytecode(contractID, string(decryptedBytecode), parameters)
	if err != nil {
		return "", fmt.Errorf("contract execution failed: %v", err)
	}

	// Step 4: Log the contract deployment.
	logEntry := fmt.Sprintf("Contract %s deployed successfully at %s", contractID, time.Now().String())
	err = sc.LedgerInstance.VirtualMachineLedger.LogEntry(logEntry, contractID) // Corrected to pass both logEntry and contractID.
	if err != nil {
		return "", fmt.Errorf("failed to log contract deployment: %v", err)
	}

	fmt.Printf("Contract %s successfully deployed.\n", contractID)
	return executionResult, nil
}


// runSolidityCompiler runs the Solidity compiler (solc) and returns the compiled bytecode.
func (sc *SolidityCompiler) runSolidityCompiler(sourcePath string) (string, error) {
	// Execute the Solidity compiler (solc) to generate the bytecode
	cmd := exec.Command("solc", "--bin", sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("solc compilation error: %v\nOutput: %s", err, string(output))
	}

	compiledBytecode := strings.TrimSpace(string(output))
	if compiledBytecode == "" {
		return "", fmt.Errorf("no bytecode generated")
	}

	return compiledBytecode, nil
}

// executeBytecode simulates the execution of a smart contract's bytecode.
func (sc *SolidityCompiler) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
	
	// Step 1: Create an encryption instance and define the encryption key.
	encryptionInstance := &Encryption{}                 // Assuming Encryption struct is defined elsewhere.
	encryptionKey := []byte("your-encryption-key")       // Example encryption key.

	// Step 2: Execute the bytecode using the VM.
	vm := NewBytecodeInterpreter(sc.LedgerInstance)
	executionResult, err := vm.ExecuteBytecode(contractID, bytecode, parameters, encryptionInstance, encryptionKey) // Pass the correct number of arguments.
	if err != nil {
		return "", fmt.Errorf("bytecode execution failed: %v", err)
	}

	// Step 3: Return the execution result.
	return fmt.Sprintf("Contract %s executed successfully: %v", contractID, executionResult), nil
}
