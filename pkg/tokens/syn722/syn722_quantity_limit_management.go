package syn722

import (
    "sync"
    "fmt"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
    "path/to/common"
)

// SYN722Token defines the structure of a SYN722 token with quantity and metadata limit management.
type SYN722Token struct {
    mutex                 sync.Mutex
    ID                    string
    Name                  string
    Owner                 string
    Mode                  string // "fungible" or "non-fungible"
    Quantity              uint64 // Used only in fungible mode
    Metadata              SYN722Metadata
    MaxQuantity           uint64
    QuantityLimitEnabled  bool
    MetadataImmutable     bool
    TransferHistory       []common.TransferRecord
    Ledger                *ledger.Ledger
    Consensus             *consensus.SynnergyConsensus
    EncryptionService     *encryption.Encryption
    CreatedAt             time.Time
    UpdatedAt             time.Time
}

// ENABLE_METADATA_IMMUTABILITY locks the token metadata, making it immutable.
func (token *SYN722Token) ENABLE_METADATA_IMMUTABILITY() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.MetadataImmutable = true
    return token.Ledger.RecordLog("MetadataImmutabilityEnabled", fmt.Sprintf("Metadata immutability enabled for token %s", token.ID))
}

// DISABLE_METADATA_IMMUTABILITY unlocks the token metadata, allowing updates.
func (token *SYN722Token) DISABLE_METADATA_IMMUTABILITY() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.MetadataImmutable {
        return fmt.Errorf("metadata is already mutable for token %s", token.ID)
    }
    token.MetadataImmutable = false
    return token.Ledger.RecordLog("MetadataImmutabilityDisabled", fmt.Sprintf("Metadata immutability disabled for token %s", token.ID))
}

// GET_METADATA_IMMUTABILITY_STATUS checks if metadata immutability is enabled.
func (token *SYN722Token) GET_METADATA_IMMUTABILITY_STATUS() bool {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.MetadataImmutable
}

// ENABLE_QUANTITY_LIMIT enables a maximum quantity limit on the token.
func (token *SYN722Token) ENABLE_QUANTITY_LIMIT(maxQuantity uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if maxQuantity == 0 {
        return fmt.Errorf("maximum quantity must be greater than zero")
    }
    token.QuantityLimitEnabled = true
    token.MaxQuantity = maxQuantity
    return token.Ledger.RecordLog("QuantityLimitEnabled", fmt.Sprintf("Quantity limit enabled for token %s with maximum of %d", token.ID, maxQuantity))
}

// DISABLE_QUANTITY_LIMIT disables the maximum quantity limit on the token.
func (token *SYN722Token) DISABLE_QUANTITY_LIMIT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.QuantityLimitEnabled {
        return fmt.Errorf("quantity limit is already disabled for token %s", token.ID)
    }
    token.QuantityLimitEnabled = false
    return token.Ledger.RecordLog("QuantityLimitDisabled", fmt.Sprintf("Quantity limit disabled for token %s", token.ID))
}

// SET_MAXIMUM_QUANTITY sets the maximum allowable quantity for the token.
func (token *SYN722Token) SET_MAXIMUM_QUANTITY(maxQuantity uint64) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if !token.QuantityLimitEnabled {
        return fmt.Errorf("quantity limit is not enabled for token %s", token.ID)
    }
    if maxQuantity == 0 {
        return fmt.Errorf("maximum quantity must be greater than zero")
    }
    token.MaxQuantity = maxQuantity
    return token.Ledger.RecordLog("MaxQuantitySet", fmt.Sprintf("Maximum quantity set to %d for token %s", maxQuantity, token.ID))
}

// GET_MAXIMUM_QUANTITY retrieves the maximum quantity limit for the token.
func (token *SYN722Token) GET_MAXIMUM_QUANTITY() (uint64, error) {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    if !token.QuantityLimitEnabled {
        return 0, fmt.Errorf("quantity limit is not enabled for token %s", token.ID)
    }
    return token.MaxQuantity, nil
}

// LOG_QUANTITY_CHANGE logs any changes in the quantity of the token.
func (token *SYN722Token) LOG_QUANTITY_CHANGE(changeDetails string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    logEntry := common.QuantityChangeLog{
        TokenID:   token.ID,
        Details:   changeDetails,
        Timestamp: time.Now(),
    }
    token.TransferHistory = append(token.TransferHistory, common.TransferRecord{
        TokenID:   token.ID,
        Details:   changeDetails,
        Timestamp: time.Now(),
    })
    return token.Ledger.RecordLog("QuantityChange", fmt.Sprintf("Quantity change logged for token %s: %s", token.ID, changeDetails))
}
