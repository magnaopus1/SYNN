package common


import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"synnergy_network/pkg/ledger"
)

// PreProcessor is responsible for preparing and validating bytecode before execution.
type PreProcessor struct {
	EncryptionModule  *Encryption // Encryption module for securing sensitive data
	LedgerInstance    *ledger.Ledger               // Ledger instance for logging pre-processing events
	OptimizationLevel int                           // Level of optimization to apply (e.g., 1 to 3)
}

// NewPreProcessor initializes a new PreProcessor instance with encryption and ledger access.
func NewPreProcessor(encryptionModule *Encryption, ledgerInstance *ledger.Ledger, optimizationLevel int) *PreProcessor {
	return &PreProcessor{
		EncryptionModule:  encryptionModule,
		LedgerInstance:    ledgerInstance,
		OptimizationLevel: optimizationLevel,
	}
}

// ValidateBytecodeHash verifies the integrity of the bytecode by checking its hash.
func (pp *PreProcessor) ValidateBytecodeHash(bytecode []byte, expectedHash string) error {
	hash := sha256.Sum256(bytecode)
	calculatedHash := hex.EncodeToString(hash[:])
	if calculatedHash != expectedHash {
		return fmt.Errorf("bytecode validation failed: hash mismatch (expected %s, got %s)", expectedHash, calculatedHash)
	}
	return nil
}

// EncryptSensitiveSections encrypts sensitive sections of the bytecode using the encryption module.
func (pp *PreProcessor) EncryptSensitiveSections(bytecode []byte, sensitiveSections []int, key []byte) ([]byte, error) {
	// Encrypt specified sections of the bytecode
	for _, start := range sensitiveSections {
		sectionEnd := start + 32 // Assume 32-byte blocks for simplicity
		if sectionEnd > len(bytecode) {
			return nil, errors.New("specified section exceeds bytecode bounds")
		}
		// Call EncryptData with the algorithm, data to encrypt, and key
		encryptedSection, err := pp.EncryptionModule.EncryptData("AES", bytecode[start:sectionEnd], key)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt section starting at %d: %v", start, err)
		}
		copy(bytecode[start:sectionEnd], encryptedSection)
	}
	return bytecode, nil
}


// OptimizeBytecode performs bytecode optimization based on the optimization level.
func (pp *PreProcessor) OptimizeBytecode(bytecode []byte) ([]byte, error) {
	optimizedBytecode := bytecode
	switch pp.OptimizationLevel {
	case 1:
		optimizedBytecode = pp.basicOptimization(bytecode)
	case 2:
		optimizedBytecode = pp.intermediateOptimization(bytecode)
	case 3:
		optimizedBytecode = pp.advancedOptimization(bytecode)
	default:
		return nil, errors.New("invalid optimization level")
	}
	return optimizedBytecode, nil
}

// basicOptimization applies a basic level of optimization, removing redundant instructions.
func (pp *PreProcessor) basicOptimization(bytecode []byte) []byte {
	return RemoveRedundantInstructions(bytecode)
}

// intermediateOptimization applies further optimizations, such as loop unrolling.
func (pp *PreProcessor) intermediateOptimization(bytecode []byte) []byte {
	bytecode = pp.basicOptimization(bytecode)
	return UnrollLoops(bytecode)
}

// advancedOptimization performs in-depth optimizations like dead code elimination and reordering.
func (pp *PreProcessor) advancedOptimization(bytecode []byte) []byte {
	bytecode = pp.intermediateOptimization(bytecode)
	return RemoveDeadCode(bytecode)
}

// RemoveRedundantInstructions removes any redundant instructions in the bytecode by analyzing consecutive repetitions.
func RemoveRedundantInstructions(bytecode []byte) []byte {
	optimizedBytecode := []byte{}
	lastInstruction := byte(0)  // Placeholder for the last unique instruction

	for _, instruction := range bytecode {
		// Only append the instruction if it’s not a repeat of the last one
		if instruction != lastInstruction || !isRedundant(instruction) {
			optimizedBytecode = append(optimizedBytecode, instruction)
			lastInstruction = instruction
		}
	}

	return optimizedBytecode
}

// isRedundant checks if an instruction is commonly redundant (e.g., repeated no-ops or repeated assignments).
func isRedundant(instruction byte) bool {
	// Define redundancy for common cases like no-ops or consecutive arithmetic with no effect.
	return instruction == 0x00 || instruction == 0x01 // Example: 0x00 for "no-op", 0x01 for redundant assignment
}

// UnrollLoops optimizes loop execution by unrolling small, predictable loops in the bytecode.
func UnrollLoops(bytecode []byte) []byte {
	optimizedBytecode := []byte{}
	loopStartIndex := -1
	loopCounter := 0

	for i := 0; i < len(bytecode); i++ {
		if isLoopStart(bytecode[i]) {
			loopStartIndex = i
			loopCounter = getLoopCounter(bytecode[i+1:])
			i += loopCounter - 1  // Move to the loop end based on detected counter
		} else if isLoopEnd(bytecode[i]) && loopStartIndex != -1 {
			loopBody := bytecode[loopStartIndex:i]
			for j := 0; j < loopCounter; j++ {
				optimizedBytecode = append(optimizedBytecode, loopBody...)
			}
			loopStartIndex = -1
		} else if loopStartIndex == -1 {
			optimizedBytecode = append(optimizedBytecode, bytecode[i])
		}
	}

	return optimizedBytecode
}

// isLoopStart checks if an instruction indicates the start of a loop.
func isLoopStart(instruction byte) bool {
	return instruction == 0x10 // Example opcode for loop start
}

// getLoopCounter retrieves the loop repetition count from subsequent bytes.
func getLoopCounter(operands []byte) int {
	if len(operands) == 0 {
		return 1 // Default to 1 if no operand provided
	}
	return int(operands[0]) // Use the first operand as loop count
}

// isLoopEnd checks if an instruction indicates the end of a loop.
func isLoopEnd(instruction byte) bool {
	return instruction == 0x11 // Example opcode for loop end
}

// RemoveDeadCode removes unreachable or unnecessary instructions from the bytecode.
func RemoveDeadCode(bytecode []byte) []byte {
	optimizedBytecode := []byte{}
	visited := make(map[int]bool)

	for i := 0; i < len(bytecode); i++ {
		if isBranchInstruction(bytecode[i]) {
			targetIndex := getBranchTarget(bytecode, i)
			if targetIndex >= 0 && targetIndex < len(bytecode) {
				visited[targetIndex] = true
			}
		}

		if !isDeadInstruction(bytecode[i], visited, i) {
			optimizedBytecode = append(optimizedBytecode, bytecode[i])
		}
	}

	return optimizedBytecode
}

// isBranchInstruction checks if an instruction is a branch or jump operation.
func isBranchInstruction(instruction byte) bool {
	return instruction == 0x20 || instruction == 0x21 // Example opcodes for branches/jumps
}

// getBranchTarget retrieves the target index for a branch/jump instruction.
func getBranchTarget(bytecode []byte, index int) int {
	if index+1 < len(bytecode) {
		return int(bytecode[index+1]) // Next byte gives target offset or address
	}
	return -1 // Invalid target if out of range
}

// isDeadInstruction checks if an instruction is dead code by analyzing its reachability.
func isDeadInstruction(instruction byte, visited map[int]bool, index int) bool {
	return instruction == 0x00 || !visited[index] // No-op or unreachable instructions are considered dead
}


// LogPreProcessEvent records a pre-processing event in the ledger.
func (pp *PreProcessor) LogPreProcessEvent(eventType, details string) {
	pp.LedgerInstance.RecordLedgerEvent(eventType, details)
}


// PrepareBytecode handles the entire pre-processing pipeline: validation, optimization, and encryption.
func (pp *PreProcessor) PrepareBytecode(bytecode []byte, expectedHash string, sensitiveSections []int, encryptionKey []byte) ([]byte, error) {
	// Validate the bytecode hash
	if err := pp.ValidateBytecodeHash(bytecode, expectedHash); err != nil {
		pp.LogPreProcessEvent("Bytecode Validation Failed", err.Error())
		return nil, err
	}
	pp.LogPreProcessEvent("Bytecode Validation Passed", "Hash matched successfully")

	// Encrypt sensitive sections of the bytecode
	encryptedBytecode, err := pp.EncryptSensitiveSections(bytecode, sensitiveSections, encryptionKey)
	if err != nil {
		pp.LogPreProcessEvent("Bytecode Encryption Failed", err.Error())
		return nil, err
	}
	pp.LogPreProcessEvent("Bytecode Encryption Completed", "Sensitive sections encrypted")

	// Optimize the bytecode based on the selected level
	optimizedBytecode, err := pp.OptimizeBytecode(encryptedBytecode)
	if err != nil {
		pp.LogPreProcessEvent("Bytecode Optimization Failed", err.Error())
		return nil, err
	}
	pp.LogPreProcessEvent("Bytecode Optimization Completed", "Optimization level applied successfully")

	return optimizedBytecode, nil
}

