package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// COMPLETE_LEASE_PAYMENT records a lease payment for the asset.
func (token *Syn131Token) COMPLETE_LEASE_PAYMENT(lessee string, amount *big.Int, paymentDate time.Time) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    paymentRecord := LeasePayment{
        Lessee:      lessee,
        Amount:      amount,
        PaymentDate: paymentDate,
    }
    token.LeaseTerms[len(token.LeaseTerms)-1].Payments = append(token.LeaseTerms[len(token.LeaseTerms)-1].Payments, paymentRecord)
    
    return token.Ledger.RecordTransaction("LeasePaymentCompleted", lessee, token.Owner, amount)
}

// APPROVE_LEASE_TERMS approves lease terms for the asset, making them active.
func (token *Syn131Token) APPROVE_LEASE_TERMS(terms LeaseTerms) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.LeaseTerms = append(token.LeaseTerms, terms)
    return token.Ledger.RecordLog("LeaseTermsApproved", fmt.Sprintf("Lease terms approved for asset %s", token.ID))
}

// CHECK_LEASE_ALLOWANCE checks if a lessee has sufficient allowance for lease payments.
func (token *Syn131Token) CHECK_LEASE_ALLOWANCE(lessee string, requiredAmount *big.Int) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    allowance := token.Ledger.GetLeaseAllowance(lessee)
    if allowance.Cmp(requiredAmount) < 0 {
        return false, fmt.Errorf("insufficient lease allowance for %s", lessee)
    }
    return true, nil
}

// SET_CO_OWNERSHIP sets co-ownership terms for the asset.
func (token *Syn131Token) SET_CO_OWNERSHIP(coOwners []string, ownershipPercents []float64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    var agreements []CoOwnershipAgreement
    for i, owner := range coOwners {
        agreements = append(agreements, CoOwnershipAgreement{
            Owner:   owner,
            Percent: ownershipPercents[i],
        })
    }
    token.CoOwnershipAgreements = agreements
    return token.Ledger.RecordLog("CoOwnershipSet", fmt.Sprintf("Co-ownership agreements set for asset %s", token.ID))
}

// GET_CO_OWNERSHIP retrieves the current co-ownership details.
func (token *Syn131Token) GET_CO_OWNERSHIP() []CoOwnershipAgreement {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.CoOwnershipAgreements
}

// ENABLE_ASSET_TRACKING enables tracking for the asset’s location or status.
func (token *Syn131Token) ENABLE_ASSET_TRACKING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AssetTrackingEnabled = true
    return token.Ledger.RecordLog("AssetTrackingEnabled", fmt.Sprintf("Asset tracking enabled for asset %s", token.ID))
}

// DISABLE_ASSET_TRACKING disables tracking for the asset.
func (token *Syn131Token) DISABLE_ASSET_TRACKING() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.AssetTrackingEnabled = false
    return token.Ledger.RecordLog("AssetTrackingDisabled", fmt.Sprintf("Asset tracking disabled for asset %s", token.ID))
}

// VALIDATE_ASSET_OWNERSHIP checks that the provided owner is a valid co-owner or primary owner.
func (token *Syn131Token) VALIDATE_ASSET_OWNERSHIP(owner string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if owner == token.Owner {
        return true, nil
    }
    for _, coOwner := range token.CoOwnershipAgreements {
        if coOwner.Owner == owner {
            return true, nil
        }
    }
    return false, fmt.Errorf("ownership validation failed for %s", owner)
}

// SET_ASSET_CLASSIFICATION sets the classification for the asset (e.g., "Intellectual Property", "Trademark").
func (token *Syn131Token) SET_ASSET_CLASSIFICATION(classification string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Classification = classification
    return token.Ledger.RecordLog("AssetClassificationSet", fmt.Sprintf("Classification set to %s for asset %s", classification, token.ID))
}

// GET_ASSET_CLASSIFICATION retrieves the classification of the asset.
func (token *Syn131Token) GET_ASSET_CLASSIFICATION() string {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Classification
}

// RECORD_SALE logs the sale of the asset, updating ownership and sale history.
func (token *Syn131Token) RECORD_SALE(buyer string, salePrice *big.Int, saleDate time.Time) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    saleRecord := SaleRecord{
        Buyer: buyer,
        Price: salePrice,
        Date:  saleDate,
    }
    token.SaleHistory = append(token.SaleHistory, saleRecord)
    token.Owner = buyer
    
    return token.Ledger.RecordTransaction("AssetSale", token.Owner, buyer, salePrice)
}

// LOG_ASSET_TRANSACTION logs a generic transaction related to the asset.
func (token *Syn131Token) LOG_ASSET_TRANSACTION(transactionType, description string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    transactionRecord := TransactionRecord{
        Type:        transactionType,
        Description: description,
        Timestamp:   time.Now(),
    }
    token.TransactionHistory = append(token.TransactionHistory, transactionRecord)
    
    return token.Ledger.RecordLog("AssetTransaction", description)
}

// VIEW_ASSET_TRANSACTION_HISTORY retrieves the complete transaction history for the asset.
func (token *Syn131Token) VIEW_ASSET_TRANSACTION_HISTORY() []TransactionRecord {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.TransactionHistory
}

// SET_ASSET_PROVENANCE adds a provenance record for tracking asset history.
func (token *Syn131Token) SET_ASSET_PROVENANCE(owner string, action string, timestamp time.Time) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    provenanceRecord := ProvenanceRecord{
        Owner:     owner,
        Action:    action,
        Timestamp: timestamp,
    }
    token.Provenance = append(token.Provenance, provenanceRecord)
    
    return token.Ledger.RecordLog("AssetProvenance", fmt.Sprintf("Provenance record added: %s by %s on %s", action, owner, timestamp))
}

// FETCH_ASSET_PROVENANCE retrieves the provenance history for the asset.
func (token *Syn131Token) FETCH_ASSET_PROVENANCE() []ProvenanceRecord {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Provenance
}