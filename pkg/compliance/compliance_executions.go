package compliance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewComplianceExecution initializes a new compliance execution
func NewComplianceExecution(actionID, executor string, rules []string, ledgerInstance *ledger.Ledger) *ComplianceExecution {
    return &ComplianceExecution{
        ExecutionID:    generateExecutionID(actionID, executor),
        ActionID:       actionID,
        Executor:       executor,
        RulesApplied:   rules,
        Timestamp:      time.Now(),
        LedgerInstance: ledgerInstance,
    }
}

// ExecuteCompliance runs the compliance checks and records the result in the ledger
func (ce *ComplianceExecution) ExecuteCompliance(actionData string) (*ComplianceExecutionResult, error) {
    ce.mutex.Lock()
    defer ce.mutex.Unlock()

    fmt.Printf("Starting compliance execution for Action ID %s by %s\n", ce.ActionID, ce.Executor)

    // Run the compliance checks
    result := ce.runComplianceChecks(actionData)
    executionResult := &ComplianceExecutionResult{
        ExecutionID:  ce.ExecutionID,
        ActionID:     ce.ActionID,
        IsValid:      result.IsValid,
        Reason:       result.Reason,
        Timestamp:    result.Timestamp,
    }

    // Ensure you have an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Encrypt the execution result
    encryptedResult, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", executionResult)), common.EncryptionKey)
    if err != nil {
        return executionResult, fmt.Errorf("failed to encrypt compliance execution result: %v", err)
    }

    // Convert the encrypted result from []byte to string
    encryptedResultString := string(encryptedResult)

    // RecordComplianceExecution expects only one argument, the encrypted result as a string
    recordResult, err := ce.LedgerInstance.ComplianceLedger.RecordComplianceExecution(encryptedResultString)
    if err != nil {
        return executionResult, fmt.Errorf("failed to record compliance execution in ledger: %v", err)
    }

    fmt.Printf("Compliance execution result for Action ID %s recorded successfully. Result: %s\n", ce.ActionID, recordResult)
    return executionResult, nil
}



// runComplianceChecks applies the compliance rules and validates the action data
func (ce *ComplianceExecution) runComplianceChecks(actionData string) ComplianceExecutionResult {
    for _, rule := range ce.RulesApplied {
        if !applyRule(rule, actionData) {
            return ComplianceExecutionResult{
                IsValid:   false,
                Reason:    fmt.Sprintf("Failed rule: %s", rule),
                Timestamp: time.Now(),
            }
        }
    }

    return ComplianceExecutionResult{
        IsValid:   true,
        Reason:    "All rules passed",
        Timestamp: time.Now(),
    }
}

// applyRule checks a specific rule against the action data (placeholder for real logic)
func applyRule(rule, actionData string) bool {
    return true // Simulate rule passing (implement actual rule logic here)
}

// generateExecutionID creates a unique identifier for each compliance execution
func generateExecutionID(actionID, executor string) string {
    input := fmt.Sprintf("%s-%s-%d", actionID, executor, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}

// RetrieveExecutionResult retrieves and decrypts a compliance execution result from the ledger
func (ce *ComplianceExecution) RetrieveExecutionResult(executionID string) (*ComplianceExecutionResult, error) {
    // Get the encrypted compliance record from the ledger
    encryptedRecord, err := ce.LedgerInstance.ComplianceLedger.GetComplianceExecutionRecord(executionID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve compliance execution result: %v", err)
    }

    // Assuming there's an EncryptedData field in ComplianceRecord that holds the encrypted result
    encryptedData := []byte(encryptedRecord.EncryptedData)

    // Decrypt the result using the encryption package
    encryptionInstance := &common.Encryption{}
    decryptedResult, err := encryptionInstance.DecryptData(encryptedData, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt compliance execution result: %v", err)
    }

    // Parse the decrypted data into a ComplianceExecutionResult struct
    var result ComplianceExecutionResult
    if err := json.Unmarshal(decryptedResult, &result); err != nil {
        return nil, fmt.Errorf("failed to parse compliance execution result: %v", err)
    }

    return &result, nil
}

