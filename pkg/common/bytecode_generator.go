package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
)

// BytecodeGenerator generates bytecode for smart contracts and virtual machine execution.
type BytecodeGenerator struct {
	LedgerInstance *ledger.Ledger    // Ledger instance for logging bytecode deployments
	mutex          sync.Mutex        // Mutex for thread-safe bytecode generation
}

// NewBytecodeGenerator initializes a new bytecode generator.
func NewBytecodeGenerator(ledgerInstance *ledger.Ledger) *BytecodeGenerator {
    return &BytecodeGenerator{
        LedgerInstance: ledgerInstance,
    }
}

// GenerateBytecode generates a unique bytecode string from contract source code and parameters.
func (bg *BytecodeGenerator) GenerateBytecode(contractID string, sourceCode string, parameters map[string]interface{}, encryption *Encryption) ([]byte, error) {
    bg.mutex.Lock()
    defer bg.mutex.Unlock()

    // Combine contract ID, source code, and parameters to generate the bytecode hash
    bytecodeInput := fmt.Sprintf("%s%s%v", contractID, sourceCode, parameters)
    bytecode := sha256.New()
    bytecode.Write([]byte(bytecodeInput))
    bytecodeHash := hex.EncodeToString(bytecode.Sum(nil))

    // Define the encryption key
    encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Ensure this is properly sized for AES

    // Encrypt the bytecode for security using the provided encryption instance
    encryptedBytecode, err := encryption.EncryptData("AES", []byte(bytecodeHash), encryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt bytecode: %v", err)
    }

    // Log the bytecode in the ledger (no need to convert to string, pass as []byte)
    err = bg.LedgerInstance.VirtualMachineLedger.RecordBytecodeDeployment(contractID, encryptedBytecode)
    if err != nil {
        return nil, fmt.Errorf("failed to log bytecode in ledger: %v", err)
    }

    fmt.Printf("Bytecode generated for contract %s and logged to ledger.\n", contractID)
    return encryptedBytecode, nil
}



// ValidateBytecode validates a given bytecode by comparing it with the expected hash.
func (bg *BytecodeGenerator) ValidateBytecode(contractID, providedBytecode string, sourceCode string, parameters map[string]interface{}) (bool, error) {
    bg.mutex.Lock()
    defer bg.mutex.Unlock()

    // Recompute the bytecode hash from the provided source code and parameters
    bytecodeInput := fmt.Sprintf("%s%s%v", contractID, sourceCode, parameters)
    bytecode := sha256.New()
    bytecode.Write([]byte(bytecodeInput))
    expectedBytecodeHash := hex.EncodeToString(bytecode.Sum(nil))

    // Compare the provided bytecode with the expected one
    if providedBytecode == expectedBytecodeHash {
        fmt.Printf("Bytecode validation successful for contract %s.\n", contractID)
        return true, nil
    }

    fmt.Printf("Bytecode validation failed for contract %s.\n", contractID)
    return false, nil
}



// GetBytecode retrieves the bytecode for a given contract from the ledger.
func (bg *BytecodeGenerator) GetBytecode(contractID string) (string, error) {
    bg.mutex.Lock()
    defer bg.mutex.Unlock()

    // Retrieve the bytecode from the ledger
    bytecode, err := bg.LedgerInstance.VirtualMachineLedger.GetBytecode(contractID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve bytecode from ledger: %v", err)
    }

    // Assuming bytecode is stored as []byte in the ledger, convert to string
    bytecodeStr := string(bytecode)

    fmt.Printf("Bytecode retrieved for contract %s from ledger.\n", contractID)
    return bytecodeStr, nil
}




