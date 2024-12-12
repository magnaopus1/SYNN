package compliance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewComplianceRestriction initializes new compliance restrictions with a set of rules
func NewComplianceRestriction(restrictionID, enforcer string, rules []string, ledgerInstance *ledger.Ledger) *ComplianceRestrictions {
    return &ComplianceRestrictions{
        RestrictionID:    generateRestrictionID(restrictionID, enforcer),
        RestrictionRules: rules,
        CreatedAt:        time.Now(),
        EnforcedBy:       enforcer,
        LedgerInstance:   ledgerInstance,
    }
}

// ApplyRestrictions applies the compliance restrictions to an action and records the result
func (cr *ComplianceRestrictions) ApplyRestrictions(actionID string, actionData string) (*RestrictionResult, error) {
    cr.mutex.Lock()
    defer cr.mutex.Unlock()

    fmt.Printf("Applying compliance restrictions for Action ID %s by %s\n", actionID, cr.EnforcedBy)

    // Run restriction checks
    result := cr.runRestrictionChecks(actionData)
    restrictionResult := &RestrictionResult{
        RestrictionID: cr.RestrictionID,
        ActionID:      actionID,
        IsRestricted:  result.IsRestricted,
        Reason:        result.Reason,
        Timestamp:     result.Timestamp,
    }

    // Ensure you have an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Encrypt and store the restriction result in the ledger
    encryptedResult, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", restrictionResult)), common.EncryptionKey)
    if err != nil {
        return restrictionResult, fmt.Errorf("failed to encrypt restriction result: %v", err)
    }

    // Convert the encrypted result to a string
    encryptedResultString := string(encryptedResult)

    // Record the restriction in the ledger with four arguments
    recordResult, err := cr.LedgerInstance.ComplianceLedger.RecordRestriction(
        cr.RestrictionID,                // Restriction ID
        encryptedResultString,           // Encrypted result as string
        restrictionResult.Reason,        // Reason for the restriction
        restrictionResult.Timestamp.Unix(), // Timestamp as int64 (Unix time)
    )
    if err != nil {
        return restrictionResult, fmt.Errorf("failed to record restriction in ledger: %v", err)
    }

    fmt.Printf("Compliance restriction result for Action ID %s recorded successfully. Result: %s\n", actionID, recordResult)
    return restrictionResult, nil
}



// runRestrictionChecks applies the compliance restriction rules and validates the action data
func (cr *ComplianceRestrictions) runRestrictionChecks(actionData string) RestrictionResult {
    for _, rule := range cr.RestrictionRules {
        if !applyRestrictionRule(rule, actionData) {
            return RestrictionResult{
                IsRestricted: true,
                Reason:       fmt.Sprintf("Restricted by rule: %s", rule),
                Timestamp:    time.Now(),
            }
        }
    }

    return RestrictionResult{
        IsRestricted: false,
        Reason:       "No restrictions applied",
        Timestamp:    time.Now(),
    }
}

// applyRestrictionRule checks a specific rule against the action data (placeholder for real logic)
func applyRestrictionRule(rule, actionData string) bool {
    // Add real rule-checking logic here, return false if the rule is violated
    return true
}

// generateRestrictionID creates a unique identifier for each compliance restriction
func generateRestrictionID(restrictionID, enforcer string) string {
    input := fmt.Sprintf("%s-%s-%d", restrictionID, enforcer, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}

// RetrieveRestrictionResult retrieves and decrypts a restriction result from the ledger
func (cr *ComplianceRestrictions) RetrieveRestrictionResult(restrictionID string) (*RestrictionResult, error) {
    // Get the encrypted restriction record from the ledger
    encryptedRecord, err := cr.LedgerInstance.ComplianceLedger.GetRestrictionRecord(restrictionID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve restriction result: %v", err)
    }

    // Assuming the ledger.Restriction struct has a field EncryptedData that stores the encrypted result
    encryptedData := []byte(encryptedRecord.EncryptedData) // Convert the EncryptedData string to a byte slice for decryption

    // Decrypt the result using the encryption package
    encryptionInstance := &common.Encryption{}
    decryptedResult, err := encryptionInstance.DecryptData(encryptedData, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt restriction result: %v", err)
    }

    // Parse the decrypted data into a RestrictionResult struct
    var result RestrictionResult
    if err := json.Unmarshal(decryptedResult, &result); err != nil {
        return nil, fmt.Errorf("failed to parse restriction result: %v", err)
    }

    return &result, nil
}



// ValidateRestrictions checks if a particular action is compliant with set restrictions
func (cr *ComplianceRestrictions) ValidateRestrictions(actionID string, actionData string) error {
    result, err := cr.ApplyRestrictions(actionID, actionData)
    if err != nil {
        return err
    }

    if result.IsRestricted {
        return errors.New(result.Reason)
    }

    return nil
}
