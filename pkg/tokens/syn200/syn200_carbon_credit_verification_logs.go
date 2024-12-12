package syn200

import (
    "sync"
    "fmt"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
)

// SYN200Token represents a carbon credit token under the SYN200 standard.
type SYN200Token struct {
    TokenID             string
    VerificationLogs    []VerificationLog
    ProjectLinkage      EmissionProjectRecord
    Issuer              IssuerRecord
    mutex               sync.Mutex
}

// VerificationLog contains information about verification activities related to a carbon credit.
type VerificationLog struct {
    Activity      string    // Description of the verification activity
    VerifiedBy    string    // Entity responsible for verification
    Timestamp     time.Time // Date and time of verification
    VerificationID string   // Unique identifier for the verification log entry
}

// EmissionProjectRecord holds data on the emission reduction project associated with the carbon credit.
type EmissionProjectRecord struct {
    ProjectID   string // Unique identifier for the emission project
    ProjectName string // Name of the emission reduction project
    Status      string // Current status of the project (e.g., "Active", "Completed", "Inactive")
}

// GET_CARBON_CREDIT_VERIFICATION_LOGS retrieves all verification logs for the carbon credit.
func (token *SYN200Token) GET_CARBON_CREDIT_VERIFICATION_LOGS() ([]VerificationLog, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if len(token.VerificationLogs) == 0 {
        return nil, fmt.Errorf("no verification logs found for carbon credit %s", token.TokenID)
    }

    return token.VerificationLogs, nil
}

// CHECK_CARBON_CREDIT_PROJECT_LINKAGE verifies if the carbon credit is currently linked to an emission reduction project.
func (token *SYN200Token) CHECK_CARBON_CREDIT_PROJECT_LINKAGE() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.ProjectLinkage.ProjectID == "" {
        return "", fmt.Errorf("no project linkage found for carbon credit %s", token.TokenID)
    }
    return fmt.Sprintf("Carbon credit %s is linked to project %s: %s", token.TokenID, token.ProjectLinkage.ProjectID, token.ProjectLinkage.Status), nil
}

// ADD_CARBON_CREDIT_VERIFICATION_LOG adds a new verification log entry for the carbon credit.
func (token *SYN200Token) ADD_CARBON_CREDIT_VERIFICATION_LOG(activity, verifiedBy, verificationID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    logEntry := VerificationLog{
        Activity:      activity,
        VerifiedBy:    verifiedBy,
        Timestamp:     time.Now(),
        VerificationID: verificationID,
    }

    token.VerificationLogs = append(token.VerificationLogs, logEntry)

    encryptedLog, err := token.Encryption.Encrypt(fmt.Sprintf("Activity: %s, Verified By: %s, Timestamp: %s", activity, verifiedBy, time.Now()))
    if err != nil {
        return fmt.Errorf("encryption failed for verification log: %v", err)
    }

    return token.Ledger.RecordLog("VerificationLogEntry", encryptedLog)
}

// GET_CARBON_CREDIT_LATEST_VERIFICATION retrieves the most recent verification log for the carbon credit.
func (token *SYN200Token) GET_CARBON_CREDIT_LATEST_VERIFICATION() (VerificationLog, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if len(token.VerificationLogs) == 0 {
        return VerificationLog{}, fmt.Errorf("no verification logs found for carbon credit %s", token.TokenID)
    }

    return token.VerificationLogs[len(token.VerificationLogs)-1], nil
}

// REMOVE_CARBON_CREDIT_VERIFICATION_LOG removes a verification log entry by its VerificationID.
func (token *SYN200Token) REMOVE_CARBON_CREDIT_VERIFICATION_LOG(verificationID string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for i, log := range token.VerificationLogs {
        if log.VerificationID == verificationID {
            token.VerificationLogs = append(token.VerificationLogs[:i], token.VerificationLogs[i+1:]...)
            return token.Ledger.RecordLog("VerificationLogRemoved", fmt.Sprintf("Verification log %s removed for carbon credit %s", verificationID, token.TokenID))
        }
    }
    return fmt.Errorf("verification log with ID %s not found for carbon credit %s", verificationID, token.TokenID)
}
