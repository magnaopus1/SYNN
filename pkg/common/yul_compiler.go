package common

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// YulCompiler manages the compilation, execution, and deployment of Yul-based smart contracts.
type YulCompiler struct {
	LedgerInstance *ledger.Ledger          // Ledger instance for logging contract activities
	CompiledCode   map[string]string       // Stores compiled bytecode for each contract
	mutex          sync.Mutex              // Mutex for thread-safe operations
}

// NewYulCompiler initializes a new YulCompiler instance.
func NewYulCompiler(ledgerInstance *ledger.Ledger) *YulCompiler {
	return &YulCompiler{
		LedgerInstance: ledgerInstance,
		CompiledCode:   make(map[string]string),
	}
}

// CompileYulContract compiles the Yul contract and stores the resulting bytecode.
func (yc *YulCompiler) CompileYulContract(contractID, contractSourcePath string) error {
    yc.mutex.Lock()
    defer yc.mutex.Unlock()

    // Step 1: Execute the Yul compiler to compile the contract.
    compiledBytecode, err := yc.runYulCompiler(contractSourcePath)
    if err != nil {
        return fmt.Errorf("compilation failed: %v", err)
    }

    // Step 2: Encrypt the compiled bytecode using the Encryption instance.
    encryptionInstance := &Encryption{}
    encryptedBytecode, err := encryptionInstance.EncryptData("AES", []byte(compiledBytecode), EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt compiled bytecode: %v", err)
    }

    // Step 3: Store the encrypted bytecode.
    yc.CompiledCode[contractID] = string(encryptedBytecode)

    // Step 4: Log the successful compilation in the ledger.
    logEntry := fmt.Sprintf("Contract %s compiled successfully at %s", contractID, time.Now().String())
    err = yc.LedgerInstance.VirtualMachineLedger.LogEntry(contractID, logEntry) // Pass contractID and log entry
    if err != nil {
        return fmt.Errorf("failed to log contract compilation: %v", err)
    }

    fmt.Printf("Contract %s successfully compiled and encrypted.\n", contractID)
    return nil
}

// DeployContract deploys the compiled contract bytecode into the blockchain.
func (yc *YulCompiler) DeployContract(contractID string, parameters map[string]interface{}) (string, error) {
    yc.mutex.Lock()
    defer yc.mutex.Unlock()

    // Step 1: Retrieve the compiled bytecode.
    encryptedBytecode, exists := yc.CompiledCode[contractID]
    if !exists {
        return "", fmt.Errorf("compiled bytecode for contract %s not found", contractID)
    }

    // Step 2: Decrypt the bytecode using the Encryption instance.
    encryptionInstance := &Encryption{}
    decryptedBytecode, err := encryptionInstance.DecryptData([]byte(encryptedBytecode), EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
    }

    // Step 3: Validate and execute the bytecode (part of the Synnergy Consensus mechanism).
    executionResult, err := yc.executeBytecode(contractID, string(decryptedBytecode), parameters)
    if err != nil {
        return "", fmt.Errorf("contract execution failed: %v", err)
    }

    // Step 4: Log the contract deployment in the ledger.
    logEntry := fmt.Sprintf("Contract %s deployed successfully at %s", contractID, time.Now().String())
    err = yc.LedgerInstance.VirtualMachineLedger.LogEntry(contractID, logEntry) // Pass contractID and log entry
    if err != nil {
        return "", fmt.Errorf("failed to log contract deployment: %v", err)
    }

    fmt.Printf("Contract %s successfully deployed.\n", contractID)
    return executionResult, nil
}



// runYulCompiler compiles the Yul smart contract using an external Yul compiler.
func (yc *YulCompiler) runYulCompiler(sourcePath string) (string, error) {
	// Execute the Yul compiler to generate the bytecode.
	cmd := exec.Command("yulc", "--bin", sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Yul compilation error: %v\nOutput: %s", err, string(output))
	}

	compiledBytecode := strings.TrimSpace(string(output))
	if compiledBytecode == "" {
		return "", fmt.Errorf("no bytecode generated")
	}

	return compiledBytecode, nil
}

// executeBytecode simulates the execution of a Yul smart contract's bytecode.
func (yc *YulCompiler) executeBytecode(contractID string, bytecode string, parameters map[string]interface{}) (string, error) {
    // The bytecode execution will interface with the blockchain's virtual machine and Synnergy Consensus.
    vm := NewBytecodeInterpreter(yc.LedgerInstance)
    
    // Create an encryption instance and key
    encryptionInstance := &Encryption{}
    encryptionKey := []byte("your-32-byte-key-for-aes-encryption")

    // Execute the bytecode with proper parameters and encryption
    executionResult, err := vm.ExecuteBytecode(contractID, bytecode, parameters, encryptionInstance, encryptionKey)
    if err != nil {
        return "", fmt.Errorf("bytecode execution failed: %v", err)
    }

    return fmt.Sprintf("Contract %s executed successfully: %v", contractID, executionResult), nil
}

