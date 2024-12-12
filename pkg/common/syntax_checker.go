package common

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)

// SyntaxChecker is responsible for verifying the syntax of smart contracts in multiple languages
// (Solidity, JavaScript, Rust, Golang, Yul) and ensuring transaction compliance.
type SyntaxChecker struct {
	LedgerInstance *ledger.Ledger   // Ledger for logging results and transaction histories
	mutex          sync.Mutex       // Mutex for thread-safety
}

// NewSyntaxChecker initializes a new SyntaxChecker instance.
func NewSyntaxChecker(ledgerInstance *ledger.Ledger) *SyntaxChecker {
	return &SyntaxChecker{
		LedgerInstance: ledgerInstance,
	}
}

// ValidateContractSyntax performs a syntax check on a smart contract before deployment, supporting multiple languages.
func (sc *SyntaxChecker) ValidateContractSyntax(contractID string, bytecode string, language string) (bool, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Step 1: Create an encryption instance and decrypt the bytecode.
	encryptionInstance := &Encryption{} // Create an instance of Encryption
	decryptedBytecode, err := encryptionInstance.DecryptData([]byte(bytecode), EncryptionKey) // Ensure EncryptionKey is defined globally
	if err != nil {
		return false, fmt.Errorf("failed to decrypt bytecode: %v", err)
	}

	// Step 2: Check syntax based on the contract's programming language.
	isValid := false
	switch language {
	case "solidity":
		isValid, err = sc.checkSoliditySyntax(string(decryptedBytecode))
	case "javascript":
		isValid, err = sc.checkJavaScriptSyntax(string(decryptedBytecode))
	case "golang":
		isValid, err = sc.checkGolangSyntax(string(decryptedBytecode))
	case "rust":
		isValid, err = sc.checkRustSyntax(string(decryptedBytecode))
	case "yul":
		isValid, err = sc.checkYulSyntax(string(decryptedBytecode))
	default:
		return false, fmt.Errorf("unsupported contract language: %s", language)
	}

	if err != nil {
		return false, err
	}

	// Step 3: Log the syntax validation result into the ledger.
	logMessage := fmt.Sprintf("Contract %s (language: %s) syntax validation result: %v at %s", contractID, language, isValid, time.Now())
	err = sc.LedgerInstance.VirtualMachineLedger.LogEntry(logMessage, contractID) // Corrected to pass both the message and the contract ID
	if err != nil {
		return false, fmt.Errorf("failed to log syntax validation: %v", err)
	}

	return isValid, nil
}


// ValidateTransactionSyntax checks that the transaction follows the Synnergy Network standards before inclusion in a sub-block.
func (sc *SyntaxChecker) ValidateTransactionSyntax(tx *Transaction) (bool, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Step 1: Check for required fields in the transaction.
	if tx.TransactionID == "" || tx.FromAddress == "" || tx.ToAddress == "" || tx.Amount < 0 {
		return false, fmt.Errorf("transaction syntax invalid: missing required fields or invalid amount")
	}

	// Step 2: Log the transaction validation result into the ledger.
	err := sc.LedgerInstance.VirtualMachineLedger.LogEntry(fmt.Sprintf("Transaction %s syntax validation result: Passed at %s", tx.TransactionID, time.Now()), tx.TransactionID)
	if err != nil {
		return false, fmt.Errorf("failed to log transaction syntax validation: %v", err)
	}

	return true, nil
}


// ValidateSubBlockSyntax validates the syntax of a sub-block before aggregation.
func (sc *SyntaxChecker) ValidateSubBlockSyntax(subBlock *SubBlock) (bool, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Step 1: Perform a basic syntax check on the sub-block's structure.
	if subBlock == nil {
		return false, fmt.Errorf("sub-block is nil")
	}
	if len(subBlock.Transactions) == 0 {
		return false, fmt.Errorf("sub-block has no transactions")
	}

	// Step 2: Log the validation result into the ledger.
	err := sc.LedgerInstance.VirtualMachineLedger.LogEntry(fmt.Sprintf("Sub-block %d syntax validation passed at %s", subBlock.Index, time.Now()), fmt.Sprintf("%d", subBlock.Index))
	if err != nil {
		return false, fmt.Errorf("failed to log sub-block syntax validation: %v", err)
	}

	return true, nil
}


// checkSoliditySyntax checks for common syntax issues in Solidity code.
func (sc *SyntaxChecker) checkSoliditySyntax(bytecode string) (bool, error) {
	if bytecode == "" {
		return false, errors.New("bytecode is empty")
	}
	if !strings.Contains(bytecode, "contract") || !strings.Contains(bytecode, "function") {
		return false, errors.New("syntax error: missing 'contract' or 'function' keyword in Solidity code")
	}
	return true, nil
}

// checkJavaScriptSyntax checks for common syntax issues in JavaScript contracts.
func (sc *SyntaxChecker) checkJavaScriptSyntax(bytecode string) (bool, error) {
	if bytecode == "" {
		return false, errors.New("bytecode is empty")
	}
	if !strings.Contains(bytecode, "function") {
		return false, errors.New("syntax error: missing 'function' keyword in JavaScript contract")
	}
	// Add more JS-specific checks if necessary
	return true, nil
}

// checkGolangSyntax checks for syntax issues in Golang contracts.
func (sc *SyntaxChecker) checkGolangSyntax(bytecode string) (bool, error) {
	if bytecode == "" {
		return false, errors.New("bytecode is empty")
	}
	if !strings.Contains(bytecode, "func") || !strings.Contains(bytecode, "package") {
		return false, errors.New("syntax error: missing 'func' or 'package' keyword in Golang contract")
	}
	return true, nil
}

// checkRustSyntax checks for common syntax issues in Rust contracts.
func (sc *SyntaxChecker) checkRustSyntax(bytecode string) (bool, error) {
	if bytecode == "" {
		return false, errors.New("bytecode is empty")
	}
	if !strings.Contains(bytecode, "fn") || !strings.Contains(bytecode, "struct") {
		return false, errors.New("syntax error: missing 'fn' or 'struct' keyword in Rust contract")
	}
	return true, nil
}

// checkYulSyntax performs syntax checks for Yul contracts.
func (sc *SyntaxChecker) checkYulSyntax(bytecode string) (bool, error) {
	if bytecode == "" {
		return false, errors.New("bytecode is empty")
	}
	if !strings.Contains(bytecode, "function") {
		return false, errors.New("syntax error: missing 'function' keyword in Yul code")
	}
	return true, nil
}
