package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

type BytecodeInterpreter struct {
    LedgerInstance *ledger.Ledger       // Ledger instance to log contract executions
    mutex          sync.Mutex           // Mutex for thread-safe bytecode execution
    State          map[string]interface{} // State of the contract, storing key-value pairs
}

// NewBytecodeInterpreter initializes a new BytecodeInterpreter with an empty state.
func NewBytecodeInterpreter(ledgerInstance *ledger.Ledger) *BytecodeInterpreter {
    return &BytecodeInterpreter{
        LedgerInstance: ledgerInstance,
        State:          make(map[string]interface{}), // Initialize an empty state map
    }
}

// ExecuteBytecode executes the provided bytecode for a smart contract.
func (bi *BytecodeInterpreter) ExecuteBytecode(contractID string, bytecode string, parameters map[string]interface{}, encryptionInstance *Encryption, encryptionKey []byte) (map[string]interface{}, error) {
    bi.mutex.Lock()
    defer bi.mutex.Unlock()

    // Step 1: Decrypt the bytecode
    decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), encryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt bytecode: %v", err)
    }

    // Convert decrypted bytecode from []byte to string
    decryptedBytecodeStr := string(decryptedBytecode)

    // Step 2: Parse the bytecode into individual instructions
    instructions, err := bi.parseBytecode(decryptedBytecodeStr)
    if err != nil {
        return nil, fmt.Errorf("failed to parse bytecode: %v", err)
    }

    // Convert []string to []interface{} for execution
    instructionsInterface := make([]interface{}, len(instructions))
    for i, instruction := range instructions {
        instructionsInterface[i] = instruction
    }

    // Step 3: Simulate execution of bytecode instructions
    executionResult, err := bi.executeInstructions(contractID, instructionsInterface, parameters)
    if err != nil {
        return nil, fmt.Errorf("failed to execute bytecode: %v", err)
    }

    // Step 4: Log the execution in the ledger
    err = bi.logExecution(contractID, executionResult)
    if err != nil {
        return nil, fmt.Errorf("failed to log bytecode execution in ledger: %v", err)
    }

    fmt.Printf("Bytecode for contract %s executed successfully.\n", contractID)
    return executionResult, nil
}




// parseBytecode breaks the bytecode string into individual instructions.
func (bi *BytecodeInterpreter) parseBytecode(bytecode string) ([]string, error) {
    if bytecode == "" {
        return nil, fmt.Errorf("empty bytecode")
    }

    instructions := splitBytecode(bytecode) // Simulates splitting by space, semicolon, etc.
    if len(instructions) == 0 {
        return nil, fmt.Errorf("no instructions found in bytecode")
    }

    return instructions, nil
}

// executeInstructions processes each instruction and determines whether to use opcode or standard instruction logic.
func (bi *BytecodeInterpreter) executeInstructions(contractID string, instructions []interface{}, parameters map[string]interface{}) (map[string]interface{}, error) {
    executionResult := make(map[string]interface{})
    executionResult["contractID"] = contractID
    executionResult["status"] = "executing"
    executionResult["executedInstructions"] = []string{}

    for _, instruction := range instructions {
        var result interface{}
        var err error

        // Determine if the instruction is opcode-based or string-based
        switch inst := instruction.(type) {
        case string:
            // Use executeInstruction for string-based instructions
            result, err = bi.executeInstruction(inst, parameters)
        case byte:
            // Use executeOpcode for byte-based (opcode) instructions
            result, err = bi.executeOpcode(inst, parameters)
        default:
            return nil, fmt.Errorf("unknown instruction type: %T", instruction)
        }

        if err != nil {
            return nil, fmt.Errorf("failed to execute instruction %v: %v", instruction, err)
        }

        // Type assertion to ensure result is a string
        if resultStr, ok := result.(string); ok {
            executionResult["executedInstructions"] = append(executionResult["executedInstructions"].([]string), resultStr)
        } else {
            return nil, fmt.Errorf("instruction %v returned a non-string result", instruction)
        }
    }

    executionResult["status"] = "executed"
    executionResult["timestamp"] = time.Now()

    return executionResult, nil
}



// logExecution records the results of bytecode execution into the ledger.
func (bi *BytecodeInterpreter) logExecution(contractID string, result map[string]interface{}) error {
    // Create an encryption instance
    encryptionInstance := &Encryption{}

    // Define the encryption key
    encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Adjust this key based on your system

    // Encrypt the result using the encryption instance
    encryptedResult, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", result)), encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt execution result: %v", err)
    }

    // Convert the encryptedResult back to the original map format (if needed)
    // Decrypt the result before passing it, using only two arguments (encrypted data and the key)
    decryptedResult, err := encryptionInstance.DecryptData(encryptedResult, encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt execution result: %v", err)
    }

    // Now you need to convert decryptedResult back into a map[string]interface{}
    var resultMap map[string]interface{}
    if err := json.Unmarshal(decryptedResult, &resultMap); err != nil {
        return fmt.Errorf("failed to unmarshal decrypted result: %v", err)
    }

    // Pass the decrypted result (map[string]interface{}) to the ledger's function
    return bi.LedgerInstance.VirtualMachineLedger.RecordContractExecution(contractID, resultMap)
} 

// splitBytecode simulates breaking bytecode into instructions based on spaces (or other delimiters).
func splitBytecode(bytecode string) []string {
    return strings.Fields(bytecode)
}

// executeInstruction interprets and executes a single bytecode instruction securely in the virtual machine.
func (bi *BytecodeInterpreter) executeInstruction(instruction string, parameters map[string]interface{}) (interface{}, error) {
    bi.mutex.Lock()
    defer bi.mutex.Unlock()

    // Log the execution attempt for auditing and tracing purposes
    fmt.Printf("Executing instruction: %s with parameters: %v\n", instruction, parameters)

    // Parse the instruction and execute based on its type
    switch instruction {
    case "ADD":
        return bi.executeAdd(parameters)

    case "SUB":
        return bi.executeSub(parameters)

    case "MUL":
        return bi.executeMul(parameters)

    case "DIV":
        return bi.executeDiv(parameters)

    case "STORE":
        err := bi.executeStore(parameters)
        if err != nil {
            return nil, fmt.Errorf("STORE operation failed: %v", err)
        }
        return "success", nil

    case "LOAD":
        return bi.executeLoad(parameters)

    case "CALL":
        return bi.executeCall(parameters)

    case "TRANSFER":
        return bi.executeTransfer(parameters)

    default:
        return nil, fmt.Errorf("unknown instruction: %s", instruction)
    }
}


// Helper functions to handle individual instructions
// executeAdd performs an addition operation.
func (bi *BytecodeInterpreter) executeAdd(parameters map[string]interface{}) (float64, error) {
    a, okA := parameters["a"].(float64)
    b, okB := parameters["b"].(float64)
    if !okA || !okB {
        return 0, fmt.Errorf("invalid parameters for ADD: %v", parameters)
    }
    return a + b, nil
}

// executeSub performs a subtraction operation.
func (bi *BytecodeInterpreter) executeSub(parameters map[string]interface{}) (float64, error) {
    a, okA := parameters["a"].(float64)
    b, okB := parameters["b"].(float64)
    if !okA || !okB {
        return 0, fmt.Errorf("invalid parameters for SUB: %v", parameters)
    }
    return a - b, nil
}

// executeMul performs a multiplication operation.
func (bi *BytecodeInterpreter) executeMul(parameters map[string]interface{}) (float64, error) {
    a, okA := parameters["a"].(float64)
    b, okB := parameters["b"].(float64)
    if !okA || !okB {
        return 0, fmt.Errorf("invalid parameters for MUL: %v", parameters)
    }
    return a * b, nil
}

// executeDiv performs a division operation with validation for division by zero.
func (bi *BytecodeInterpreter) executeDiv(parameters map[string]interface{}) (float64, error) {
    a, okA := parameters["a"].(float64)
    b, okB := parameters["b"].(float64)
    if !okA || !okB {
        return 0, fmt.Errorf("invalid parameters for DIV: %v", parameters)
    }
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// executeStore stores a key-value pair in the smart contract's state.
func (bi *BytecodeInterpreter) executeStore(parameters map[string]interface{}) error {
    key, okKey := parameters["key"].(string)
    value, okValue := parameters["value"].(interface{})
    if !okKey || !okValue {
        return fmt.Errorf("invalid parameters for STORE: %v", parameters)
    }

    // Ensure state storage is thread-safe
    bi.State[key] = value
    fmt.Printf("Stored value '%v' in key '%s'\n", value, key)
    return nil
}

func (bi *BytecodeInterpreter) executeLoad(parameters map[string]interface{}) (interface{}, error) {
    key, okKey := parameters["key"].(string)
    if !okKey {
        return nil, fmt.Errorf("invalid parameters for LOAD: %v", parameters)
    }

    // Retrieve the value from the state storage
    value, exists := bi.State[key]
    if !exists {
        return nil, fmt.Errorf("key '%s' not found in state", key)
    }

    fmt.Printf("Loaded value '%v' from key '%s'\n", value, key)
    return value, nil
}

// executeCall simulates calling another smart contract or external system.
func (bi *BytecodeInterpreter) executeCall(parameters map[string]interface{}) (interface{}, error) {
    contractID, okID := parameters["contractID"].(string)
    method, okMethod := parameters["method"].(string)
    args, okArgs := parameters["args"].(map[string]interface{})
    if !okID || !okMethod || !okArgs {
        return nil, fmt.Errorf("invalid parameters for CALL: %v", parameters)
    }

    // Simulate the call to another contract (this could be a cross-contract execution or external API call)
    fmt.Printf("Calling contract %s with method '%s' and arguments %v\n", contractID, method, args)

    // In a real-world implementation, you would retrieve the contract's bytecode, invoke the method, and pass the arguments.
    result := fmt.Sprintf("Simulated call result from contract %s, method '%s'", contractID, method)
    return result, nil
}

// executeTransfer handles transferring funds from the contract or an address.
func (bi *BytecodeInterpreter) executeTransfer(parameters map[string]interface{}) (interface{}, error) {
    fromAddress, okFrom := parameters["from"].(string)
    toAddress, okTo := parameters["to"].(string)
    amount, okAmount := parameters["amount"].(float64)
    if !okFrom || !okTo || !okAmount {
        return nil, fmt.Errorf("invalid parameters for TRANSFER: %v", parameters)
    }

    // Simulate the transfer of funds (this could be interacting with a balance ledger or sending tokens)
    fmt.Printf("Transferring %.2f units from %s to %s\n", amount, fromAddress, toAddress)

    // In a real-world scenario, you'd integrate with the ledger or balance mechanism here
    result := fmt.Sprintf("Transferred %.2f units from %s to %s", amount, fromAddress, toAddress)
    return result, nil
}


// executeOpcode interprets and executes a single opcode instruction in the virtual machine. Add rest from packages later
func (bi *BytecodeInterpreter) executeOpcode(opcode byte, parameters map[string]interface{}) (interface{}, error) {
    bi.mutex.Lock()
    defer bi.mutex.Unlock()

    // Log the execution attempt for auditing and tracing purposes
    fmt.Printf("Executing opcode: 0x%X with parameters: %v\n", opcode, parameters)

    switch opcode {
    // Arithmetic Operations (0x01 - 0x0F)
    case 0x01: // ADD
        return bi.executeAdd(parameters)
    case 0x02: // SUB
        return bi.executeSub(parameters)
    case 0x03: // MUL
        return bi.executeMul(parameters)
    case 0x04: // DIV
        return bi.executeDiv(parameters)
	
    default:
        if opcode >= 0xD0 && opcode <= 0xFF {
            return nil, fmt.Errorf("opcode 0x%X is not yet implemented", opcode)
        }
        return nil, fmt.Errorf("unknown opcode: 0x%X", opcode)
    }
}

