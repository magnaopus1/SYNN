package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// SYN130Token represents a token with comprehensive attributes for tangible assets.
type SYN130Token struct {
    ID                    string
    Name                  string
    Owner                 string
    Value                 *big.Int
    Metadata              SYN130Metadata
    LeaseAgreements       []LeaseAgreement
    TransactionHistory    []TransactionRecord
    mutex                 sync.Mutex
}

// TRANSFER_TANGIBLE_ASSET transfers ownership of the tangible asset to a new owner.
func (token *SYN130Token) TRANSFER_TANGIBLE_ASSET(newOwner string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    previousOwner := token.Owner
    token.Owner = newOwner
    transferRecord := TransactionRecord{
        Type:        "OwnershipTransfer",
        Description: fmt.Sprintf("Asset transferred from %s to %s", previousOwner, newOwner),
        Timestamp:   time.Now(),
    }
    token.TransactionHistory = append(token.TransactionHistory, transferRecord)

    return token.Ledger.RecordTransaction("AssetTransfer", previousOwner, newOwner, token.Value)
}

// CHECK_ASSET_BALANCE retrieves the current value or balance of the tangible asset.
func (token *SYN130Token) CHECK_ASSET_BALANCE() *big.Int {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Value
}

// GET_ASSET_METADATA returns the metadata associated with the tangible asset.
func (token *SYN130Token) GET_ASSET_METADATA() SYN130Metadata {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Metadata
}

// UPDATE_ASSET_METADATA updates the metadata for the asset, requiring encryption.
func (token *SYN130Token) UPDATE_ASSET_METADATA(newMetadata SYN130Metadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    encryptedData, err := token.Encryption.Encrypt(newMetadata)
    if err != nil {
        return fmt.Errorf("failed to encrypt metadata: %v", err)
    }

    token.Metadata = newMetadata
    token.Metadata.EncryptedData = encryptedData
    return token.Ledger.RecordLog("MetadataUpdated", fmt.Sprintf("Metadata updated for asset %s", token.ID))
}

// SET_ASSET_VALUE sets a new value for the tangible asset.
func (token *SYN130Token) SET_ASSET_VALUE(newValue *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Value = newValue
    return token.Ledger.RecordLog("AssetValueSet", fmt.Sprintf("Asset value set to %s for asset %s", newValue.String(), token.ID))
}

// FETCH_ASSET_VALUE retrieves the current market or appraised value of the asset.
func (token *SYN130Token) FETCH_ASSET_VALUE() *big.Int {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Value
}

// LOCK_TANGIBLE_ASSET prevents any changes or transfers for the asset.
func (token *SYN130Token) LOCK_TANGIBLE_ASSET() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.Locked = true
    return token.Ledger.RecordLog("AssetLocked", fmt.Sprintf("Asset %s is now locked", token.ID))
}

// UNLOCK_TANGIBLE_ASSET allows changes or transfers for the asset to resume.
func (token *SYN130Token) UNLOCK_TANGIBLE_ASSET() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.Locked = false
    return token.Ledger.RecordLog("AssetUnlocked", fmt.Sprintf("Asset %s is now unlocked", token.ID))
}

// CREATE_LEASE_AGREEMENT establishes a lease agreement for the asset.
func (token *SYN130Token) CREATE_LEASE_AGREEMENT(lessee string, terms LeaseTerms) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    leaseAgreement := LeaseAgreement{
        Lessee: lessee,
        Terms:  terms,
    }
    token.LeaseAgreements = append(token.LeaseAgreements, leaseAgreement)
    
    return token.Ledger.RecordLog("LeaseAgreementCreated", fmt.Sprintf("Lease agreement created for asset %s with lessee %s", token.ID, lessee))
}

// DELETE_LEASE_AGREEMENT removes an existing lease agreement for the asset.
func (token *SYN130Token) DELETE_LEASE_AGREEMENT(lessee string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for i, agreement := range token.LeaseAgreements {
        if agreement.Lessee == lessee {
            token.LeaseAgreements = append(token.LeaseAgreements[:i], token.LeaseAgreements[i+1:]...)
            return token.Ledger.RecordLog("LeaseAgreementDeleted", fmt.Sprintf("Lease agreement for asset %s with lessee %s deleted", token.ID, lessee))
        }
    }
    return fmt.Errorf("lease agreement with lessee %s not found", lessee)
}

// QUERY_LEASE_AGREEMENT_STATUS checks the status of a lease agreement for the asset.
func (token *SYN130Token) QUERY_LEASE_AGREEMENT_STATUS(lessee string) (LeaseAgreement, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for _, agreement := range token.LeaseAgreements {
        if agreement.Lessee == lessee {
            return agreement, nil
        }
    }
    return LeaseAgreement{}, fmt.Errorf("lease agreement with lessee %s not found", lessee)
}

// INITIATE_LEASE_PAYMENT processes a lease payment for an active lease agreement.
func (token *SYN130Token) INITIATE_LEASE_PAYMENT(lessee string, paymentAmount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for i, agreement := range token.LeaseAgreements {
        if agreement.Lessee == lessee {
            paymentRecord := LeasePayment{
                Lessee:      lessee,
                Amount:      paymentAmount,
                PaymentDate: time.Now(),
            }
            token.LeaseAgreements[i].Payments = append(token.LeaseAgreements[i].Payments, paymentRecord)
            return token.Ledger.RecordTransaction("LeasePayment", lessee, token.Owner, paymentAmount)
        }
    }
    return fmt.Errorf("active lease agreement with lessee %s not found", lessee)
}
