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

// NewComplianceContract initializes a new compliance contract
func NewComplianceContract(creator string, rules []string, ledgerInstance *ledger.Ledger) *ComplianceContract {
    return &ComplianceContract{
        ContractID:      generateContractID(creator),
        Creator:         creator,
        ComplianceRules: rules,
        LedgerInstance:  ledgerInstance,
    }
}

// EnforceCompliance checks compliance on an action using predefined rules and records the result
func (cc *ComplianceContract) EnforceCompliance(actionID, actionData, initiator string) (*ComplianceResult, error) {
    cc.mutex.Lock()
    defer cc.mutex.Unlock()

    fmt.Printf("Compliance check initiated for Action ID %s by %s\n", actionID, initiator)

    // Perform the compliance checks
    result := cc.runComplianceChecks(actionData)
    complianceResult := &ComplianceResult{
        ActionID:  actionID,
        IsValid:   result.IsValid,
        Reason:    result.Reason,
        Timestamp: result.Timestamp,
    }

    if !result.IsValid {
        return complianceResult, fmt.Errorf("compliance check failed: %s", result.Reason)
    }

    // Ensure you have an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Encrypt and store the compliance result in the ledger
    // Pass the algorithm "AES" as the first argument
    encryptedResult, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", complianceResult)), common.EncryptionKey)
    if err != nil {
        return complianceResult, fmt.Errorf("failed to encrypt compliance result: %v", err)
    }

    // RecordCompliance expects two string arguments, so use the result directly
    recordResult, err := cc.LedgerInstance.ComplianceLedger.RecordCompliance(actionID, string(encryptedResult))
    if err != nil {
        return complianceResult, fmt.Errorf("failed to record compliance result in ledger: %v", err)
    }

    fmt.Printf("Compliance result for Action ID %s recorded successfully. Result: %s\n", actionID, recordResult)
    return complianceResult, nil
}



// runComplianceChecks performs compliance checks on the action based on the rules
func (cc *ComplianceContract) runComplianceChecks(actionData string) ComplianceResult {
    for _, rule := range cc.ComplianceRules {
        if !checkRule(rule, actionData) {
            return ComplianceResult{
                IsValid:   false,
                Reason:    fmt.Sprintf("Failed rule: %s", rule),
                Timestamp: time.Now(),
            }
        }
    }

    return ComplianceResult{
        IsValid:   true,
        Reason:    "All rules passed",
        Timestamp: time.Now(),
    }
}

// checkRule simulates a compliance rule check (placeholder for real logic)
func checkRule(rule, actionData string) bool {
    return true // Assume all rules are passed for now
}

// generateContractID creates a unique identifier for the compliance contract based on the creator's details
func generateContractID(creator string) string {
    hashInput := fmt.Sprintf("%s-%d", creator, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// RetrieveComplianceResult retrieves and decrypts a compliance result from the ledger
func (cc *ComplianceContract) RetrieveComplianceResult(actionID string) (*ComplianceResult, error) {
    // Get the encrypted compliance result from the ledger
    encryptedRecord, err := cc.LedgerInstance.ComplianceLedger.GetComplianceRecord(actionID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve compliance result: %v", err)
    }

    // Assuming there's a field EncryptedData in the ComplianceRecord struct that holds the encrypted result
    encryptedData := []byte(encryptedRecord.EncryptedData) // Convert encrypted data to []byte if it's a string

    // Decrypt the result using the encryption package
    encryptionInstance := &common.Encryption{}
    decryptedResult, err := encryptionInstance.DecryptData(encryptedData, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt compliance result: %v", err)
    }

    // Parse the decrypted data into a ComplianceResult struct
    var result ComplianceResult
    if err := json.Unmarshal(decryptedResult, &result); err != nil {
        return nil, fmt.Errorf("failed to parse compliance result: %v", err)
    }

    return &result, nil
}



// SmartContractInvocation invokes the compliance contract as a smart contract
func (cc *ComplianceContract) SmartContractInvocation(actionID, actionData, initiator string) error {
    fmt.Printf("Invoking compliance contract for Action ID %s by %s\n", actionID, initiator)

    // Create the contract invocation
    contractInvocation := common.ContractInvocation{
        ContractAddress: cc.ContractID, // Assuming cc.ContractID holds the contract address
        Method:          "invokeComplianceContract", // This could be a contract method name
        Params: map[string]string{
            "actionID":  actionID,
            "actionData": actionData,
            "initiator":  initiator,
        },
        CallerAddress: initiator, // Assuming the initiator is the caller
        GasLimit:      21000, // Set a gas limit as appropriate
        GasPrice:      100, // Set a gas price as appropriate
        Timestamp:     time.Now(),
    }

    // Add the invocation to the ledger as a smart contract action
    result, err := cc.LedgerInstance.ComplianceLedger.RecordContractInvocation(contractInvocation.ContractAddress, fmt.Sprintf("%+v", contractInvocation))
    if err != nil {
        return fmt.Errorf("failed to record contract invocation: %v", err)
    }

    fmt.Printf("Smart contract for Action ID %s invoked successfully. Result: %s\n", actionID, result)
    return nil
}

