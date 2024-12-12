package syn200

import (
	"errors"
	"sync"
	"time"
)

// EncryptSensitiveData encrypts all sensitive data in a SYN200 token for secure storage.
func EncryptSensitiveData(token *common.SYN200Token) error {
    encryptedMetadata, err := encryption.EncryptMetadata(token.CreditMetadata)
    if err != nil {
        return fmt.Errorf("failed to encrypt token metadata: %v", err)
    }
    token.EncryptedMetadata = encryptedMetadata

    // Record encryption event in the ledger
    encryptionEvent := common.TokenEvent{
        EventID:        generateEventID(),
        TokenID:        token.TokenID,
        EventType:      "Encryption",
        EventTimestamp: time.Now(),
        Description:    "Sensitive metadata encrypted for security",
    }

    encryptedEvent, err := encryption.EncryptMetadata(encryptionEvent)
    if err != nil {
        return fmt.Errorf("failed to encrypt encryption event: %v", err)
    }

    if err := ledger.RecordEvent(token, encryptedEvent); err != nil {
        return fmt.Errorf("failed to log encryption event in ledger: %v", err)
    }

    return nil
}

// ValidateAccessControl checks if a user has the necessary permissions to perform actions on a SYN200 token.
func ValidateAccessControl(tokenID, userID string, role string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    if role != "Owner" && role != "Verifier" && role != "Issuer" {
        return fmt.Errorf("invalid role: %s", role)
    }

    // Check for role-based access control
    hasAccess := false
    switch role {
    case "Owner":
        for _, ownerRecord := range token.OwnershipHistory {
            if ownerRecord.OwnerID == userID {
                hasAccess = true
                break
            }
        }
    case "Verifier":
        for _, verifier := range token.VerificationLogs {
            if verifier.VerifierName == userID {
                hasAccess = true
                break
            }
        }
    case "Issuer":
        if token.Issuer.IssuerID == userID {
            hasAccess = true
        }
    }

    if !hasAccess {
        return fmt.Errorf("user %s does not have %s access to token %s", userID, role, tokenID)
    }

    return nil
}

// DetectFraudulentActivity scans for patterns of fraudulent activity on a SYN200 token.
func DetectFraudulentActivity(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    fraudDetected, fraudDetails := detectFraud(token)
    if fraudDetected {
        fraudEvent := common.TokenEvent{
            EventID:        generateEventID(),
            TokenID:        tokenID,
            EventType:      "Fraud Detection",
            EventTimestamp: time.Now(),
            Description:    "Fraudulent activity detected and logged",
        }

        // Encrypt and record the fraud detection event
        encryptedEvent, err := encryption.EncryptMetadata(fraudEvent)
        if err != nil {
            return fmt.Errorf("failed to encrypt fraud event: %v", err)
        }

        if err := ledger.RecordFraudAttempt(token, encryptedEvent); err != nil {
            return fmt.Errorf("failed to record fraud attempt in ledger: %v", err)
        }

        return fmt.Errorf("fraud detected in token %s", tokenID)
    }

    return nil
}

// RunComplianceCheck performs a compliance check to ensure SYN200 tokens meet required regulatory standards.
func RunComplianceCheck(tokenID string) error {
    token, err := ledger.GetTokenByID(tokenID)
    if err != nil {
        return fmt.Errorf("token not found: %v", err)
    }

    complianceRecord := common.ComplianceRecord{
        RecordID:         generateComplianceID(),
        ComplianceDate:   time.Now(),
        RegulationBody:   "Environmental Compliance Agency",
        ComplianceStatus: "Compliant",
    }

    // Encrypt and log the compliance check in the ledger
    encryptedRecord, err := encryption.EncryptMetadata(complianceRecord)
    if err != nil {
        return fmt.Errorf("failed to encrypt compliance record: %v", err)
    }

    if err := ledger.RecordCompliance(token, encryptedRecord); err != nil {
        return fmt.Errorf("failed to log compliance record in ledger: %v", err)
    }

    return nil
}

// generateEventID creates a unique identifier for events.
func generateEventID() string {
    b := make([]byte, 12)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

// detectFraud is an internal function that examines the token's history for fraud patterns.
func detectFraud(token *common.SYN200Token) (bool, common.FraudDetails) {
    for _, record := range token.OwnershipHistory {
        // Detect anomalies (e.g., rapid transfers, invalid ownership claims)
        if record.TransferMethod == "Suspicious Transfer" {
            return true, common.FraudDetails{
                TokenID:          token.TokenID,
                DetectedAt:       time.Now(),
                Description:      "Suspicious transfer pattern detected",
            }
        }
    }
    return false, common.FraudDetails{}
}

// generateComplianceID creates a unique ID for compliance records.
func generateComplianceID() string {
    b := make([]byte, 12)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
