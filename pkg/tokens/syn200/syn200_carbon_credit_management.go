package syn200

import (
    "sync"
    "fmt"
    "time"
    "math/big"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
)

// SYN200Token represents a carbon credit token under the SYN200 standard.
type SYN200Token struct {
    TokenID                string
    CreditMetadata         CarbonCreditMetadata
    Issuer                 IssuerRecord
    OwnershipHistory       []OwnershipRecord
    ValidityStatus         string
    ExpirationDate         *time.Time
    Transferable           bool
    ProjectLinkage         EmissionProjectRecord
    EncryptedMetadata      []byte
    RealTimeUpdatesEnabled bool
    mutex                  sync.Mutex
}

// TRANSFER_CARBON_CREDIT transfers ownership of the carbon credit to a new owner.
func (token *SYN200Token) TRANSFER_CARBON_CREDIT(newOwner string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.Transferable {
        return fmt.Errorf("carbon credit %s is not transferable", token.TokenID)
    }
    transferRecord := OwnershipRecord{
        Owner:     newOwner,
        Timestamp: time.Now(),
        Amount:    amount,
    }
    token.OwnershipHistory = append(token.OwnershipHistory, transferRecord)

    return token.Ledger.RecordTransaction("CarbonCreditTransfer", token.Issuer.IssuerName, newOwner, amount)
}

// CHECK_CARBON_CREDIT_BALANCE retrieves the remaining balance or amount of carbon offset available.
func (token *SYN200Token) CHECK_CARBON_CREDIT_BALANCE() *big.Int {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Assuming `CreditMetadata.CO2OffsetAmount` is the remaining balance in metric tons.
    return big.NewInt(int64(token.CreditMetadata.CO2OffsetAmount))
}

// GET_CARBON_CREDIT_METADATA retrieves the metadata associated with the carbon credit.
func (token *SYN200Token) GET_CARBON_CREDIT_METADATA() CarbonCreditMetadata {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.CreditMetadata
}

// UPDATE_CARBON_CREDIT_METADATA updates the metadata for the carbon credit.
func (token *SYN200Token) UPDATE_CARBON_CREDIT_METADATA(newMetadata CarbonCreditMetadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedData, err := token.Encryption.Encrypt(newMetadata)
    if err != nil {
        return fmt.Errorf("failed to encrypt metadata: %v", err)
    }

    token.CreditMetadata = newMetadata
    token.EncryptedMetadata = encryptedData

    return token.Ledger.RecordLog("MetadataUpdated", fmt.Sprintf("Metadata updated for carbon credit %s", token.TokenID))
}

// SET_CARBON_CREDIT_VALIDITY_STATUS sets the validity status for the carbon credit.
func (token *SYN200Token) SET_CARBON_CREDIT_VALIDITY_STATUS(status string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ValidityStatus = status
    return token.Ledger.RecordLog("ValidityStatusUpdated", fmt.Sprintf("Validity status set to %s for carbon credit %s", status, token.TokenID))
}

// FETCH_CARBON_CREDIT_VALIDITY_STATUS retrieves the current validity status of the carbon credit.
func (token *SYN200Token) FETCH_CARBON_CREDIT_VALIDITY_STATUS() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.ValidityStatus
}

// LOCK_CARBON_CREDIT locks the carbon credit, preventing any changes or transfers.
func (token *SYN200Token) LOCK_CARBON_CREDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Transferable = false
    return token.Ledger.RecordLog("CarbonCreditLocked", fmt.Sprintf("Carbon credit %s is now locked", token.TokenID))
}

// UNLOCK_CARBON_CREDIT unlocks the carbon credit, allowing transfers to resume.
func (token *SYN200Token) UNLOCK_CARBON_CREDIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Transferable = true
    return token.Ledger.RecordLog("CarbonCreditUnlocked", fmt.Sprintf("Carbon credit %s is now unlocked", token.TokenID))
}

// CREATE_CARBON_CREDIT_PROJECT links a carbon credit to a new emission reduction project.
func (token *SYN200Token) CREATE_CARBON_CREDIT_PROJECT(project EmissionProjectRecord) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ProjectLinkage = project
    return token.Ledger.RecordLog("ProjectLinked", fmt.Sprintf("Emission reduction project %s linked to carbon credit %s", project.ProjectID, token.TokenID))
}

// DELETE_CARBON_CREDIT_PROJECT removes the linkage of the carbon credit to an emission reduction project.
func (token *SYN200Token) DELETE_CARBON_CREDIT_PROJECT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    projectID := token.ProjectLinkage.ProjectID
    token.ProjectLinkage = EmissionProjectRecord{}
    return token.Ledger.RecordLog("ProjectUnlinked", fmt.Sprintf("Emission reduction project %s unlinked from carbon credit %s", projectID, token.TokenID))
}

// QUERY_CARBON_CREDIT_PROJECT_STATUS retrieves the current status of the emission reduction project linked to the carbon credit.
func (token *SYN200Token) QUERY_CARBON_CREDIT_PROJECT_STATUS() (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.ProjectLinkage.ProjectID == "" {
        return "", fmt.Errorf("no emission reduction project linked to carbon credit %s", token.TokenID)
    }
    return token.ProjectLinkage.Status, nil
}
