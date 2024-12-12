package common

import (
	"fmt"
	"strings"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// CompilationDebugger handles the compilation and debugging of smart contract bytecode.
type CompilationDebugger struct {
	LedgerInstance *ledger.Ledger // Ledger instance for logging compilation and debugging
	mutex          sync.Mutex     // Mutex for thread-safe compilation and debugging
}

// NewCompilationDebugger initializes a new CompilationDebugger.
func NewCompilationDebugger(ledgerInstance *ledger.Ledger) *CompilationDebugger {
    return &CompilationDebugger{
        LedgerInstance: ledgerInstance,
    }
}

// CompileBytecode attempts to compile the given smart contract bytecode.
func (cd *CompilationDebugger) CompileBytecode(contractID, bytecode string, encryptionInstance *Encryption, encryptionKey []byte) (string, error) {
    cd.mutex.Lock()
    defer cd.mutex.Unlock()

    // Step 1: Decrypt the bytecode for analysis
    decryptedBytecodeBytes, err := encryptionInstance.DecryptData([]byte(bytecode), encryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
    }

    // Convert decrypted bytecode from []byte to string for further processing
    decryptedBytecode := string(decryptedBytecodeBytes)

    // Step 2: Perform compilation checks on the bytecode
    if strings.TrimSpace(decryptedBytecode) == "" {
        return "", fmt.Errorf("empty bytecode provided for contract %s", contractID)
    }

    compilationErrors := cd.performCompilationChecks(decryptedBytecode)
    if len(compilationErrors) > 0 {
        return "", fmt.Errorf("compilation errors: %v", strings.Join(compilationErrors, ", "))
    }

    // Step 3: Log successful compilation in the ledger (pass encryption instance and key)
    err = cd.logCompilation(contractID, decryptedBytecode, encryptionInstance, encryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to log successful compilation to the ledger: %v", err)
    }

    return "Compilation successful", nil
}



// DebugBytecode attempts to debug the given smart contract bytecode and returns debugging information.
func (cd *CompilationDebugger) DebugBytecode(contractID, bytecode string, encryptionInstance *Encryption, encryptionKey []byte) (string, error) {
    cd.mutex.Lock()
    defer cd.mutex.Unlock()

    // Step 1: Decrypt the bytecode for analysis
    decryptedBytecodeBytes, err := encryptionInstance.DecryptData([]byte(bytecode), encryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt bytecode: %v", err)
    }

    // Convert decrypted bytecode from []byte to string
    decryptedBytecode := string(decryptedBytecodeBytes)

    // Step 2: Perform debugging checks
    debugLogs := cd.performDebuggingChecks(decryptedBytecode)

    // Step 3: Log debugging results in the ledger (pass encryption instance and key)
    err = cd.logDebugging(contractID, debugLogs, encryptionInstance, encryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to log debugging information to the ledger: %v", err)
    }

    return debugLogs, nil
}


// performCompilationChecks validates the bytecode to ensure it is suitable for execution.
func (cd *CompilationDebugger) performCompilationChecks(bytecode string) []string {
    var errors []string

    // Check for forbidden instructions, syntax errors, or malformed bytecode
    if strings.Contains(bytecode, "FORBIDDEN_OP") {
        errors = append(errors, "Found forbidden operation in bytecode")
    }

    if len(strings.Fields(bytecode)) == 0 {
        errors = append(errors, "Bytecode contains no executable instructions")
    }

    // Additional checks like control flow, loop validation, memory overflow checks, etc., can be added here.

    return errors
}

// performDebuggingChecks analyzes the bytecode to identify potential issues during execution.
func (cd *CompilationDebugger) performDebuggingChecks(bytecode string) string {
    // Step 1: Log potential issues like improper memory access, invalid operations, etc.
    var debugLog []string

    if strings.Contains(bytecode, "INVALID_MEM_ACCESS") {
        debugLog = append(debugLog, "Warning: Invalid memory access detected in bytecode")
    }

    if strings.Contains(bytecode, "INFINITE_LOOP") {
        debugLog = append(debugLog, "Warning: Infinite loop detected in bytecode")
    }

    if strings.Contains(bytecode, "DEPRECATED_OP") {
        debugLog = append(debugLog, "Warning: Deprecated operation detected in bytecode")
    }

    // Add more debugging checks based on typical execution patterns and security risks.

    // Step 2: Return the debug logs as a single string
    if len(debugLog) == 0 {
        debugLog = append(debugLog, "No issues detected during debugging")
    }

    return strings.Join(debugLog, "\n")
}


// logCompilation records the successful bytecode compilation in the ledger.
func (cd *CompilationDebugger) logCompilation(contractID, bytecode string, encryptionInstance *Encryption, encryptionKey []byte) error {
    logEntry := fmt.Sprintf("Compilation successful for contract %s at %s. Bytecode: %s", contractID, time.Now().String(), bytecode)
    
    // Encrypt the log entry using the encryption instance
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logEntry), encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt compilation log: %v", err)
    }

    // Create a map to store the encrypted log
    encryptedLogMap := map[string]interface{}{
        "encryptedLog": string(encryptedLog),
    }

    // Record the compilation log in the ledger
    return cd.LedgerInstance.VirtualMachineLedger.RecordContractExecution(contractID, encryptedLogMap)
}

// logDebugging records the debugging information in the ledger.
func (cd *CompilationDebugger) logDebugging(contractID, debugLog string, encryptionInstance *Encryption, encryptionKey []byte) error {
    logEntry := fmt.Sprintf("Debugging information for contract %s at %s. Debug Logs: %s", contractID, time.Now().String(), debugLog)
    
    // Encrypt the log entry using the encryption instance
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logEntry), encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt debugging log: %v", err)
    }

    // Create a map to store the encrypted log
    encryptedLogMap := map[string]interface{}{
        "encryptedLog": string(encryptedLog),
    }

    // Record the debugging log in the ledger
    return cd.LedgerInstance.VirtualMachineLedger.RecordContractExecution(contractID, encryptedLogMap)
}


// SimulateCompilationDebugging simulates the compilation and debugging of multiple contracts.
func (cd *CompilationDebugger) SimulateCompilationDebugging(contracts []*SmartContract, encryptionInstance *Encryption, encryptionKey []byte) {
    for _, contract := range contracts {
        // Print the start of compilation and debugging
        fmt.Printf("Compiling and debugging contract %s...\n", contract.ID)

        // Compile the contract
        _, err := cd.CompileBytecode(contract.ID, contract.Bytecode, encryptionInstance, encryptionKey)
        if err != nil {
            fmt.Printf("Compilation failed for contract %s: %v\n", contract.ID, err)
            continue // Skip debugging if compilation fails
        }

        // Debug the contract
        _, err = cd.DebugBytecode(contract.ID, contract.Bytecode, encryptionInstance, encryptionKey)
        if err != nil {
            fmt.Printf("Debugging failed for contract %s: %v\n", contract.ID, err)
        } else {
            fmt.Printf("Debugging completed successfully for contract %s\n", contract.ID)
        }
    }
}



