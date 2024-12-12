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

// SYN722Token struct defines the main structure for the Syn722 token, supporting dual-mode functionality.
type SYN722Token struct {
    mutex               sync.Mutex
    ID                  string
    Name                string
    Owner               string
    Mode                string // "fungible" or "non-fungible"
    Quantity            uint64
    Metadata            SYN722Metadata
    TokenURI            string
    ApprovedAddress     string
    Ledger              *ledger.Ledger
    Consensus           *consensus.SynnergyConsensus
    EncryptionService   *encryption.Encryption
    ModeChangeHistory   []common.ModeChangeLog
    TransferHistory     []common.TransferRecord
    CreatedAt           time.Time
    UpdatedAt           time.Time
}

// TRANSFER_SYN722_TOKEN transfers ownership of a token.
func (token *SYN722Token) TRANSFER_SYN722_TOKEN(newOwner string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Mode == "non-fungible" && newOwner == token.Owner {
        return fmt.Errorf("transfer failed: token %s is already owned by %s", token.ID, newOwner)
    }

    err := token.Ledger.Transfer(token.Owner, newOwner, 1)
    if err != nil {
        return err
    }

    token.Owner = newOwner
    token.UpdatedAt = time.Now()
    return token.Ledger.RecordLog("Transfer", fmt.Sprintf("Token %s transferred to %s", token.ID, newOwner))
}

// APPROVE_SYN722_TOKEN_TRANSFER grants transfer rights to an address.
func (token *SYN722Token) APPROVE_SYN722_TOKEN_TRANSFER(address string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ApprovedAddress = address
    return token.Ledger.RecordLog("Approval", fmt.Sprintf("Approved %s for transferring token %s", address, token.ID))
}

// CHECK_SYN722_TOKEN_BALANCE returns the balance of the token.
func (token *SYN722Token) CHECK_SYN722_TOKEN_BALANCE() uint64 {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Quantity
}

// GET_SYN722_TOKEN_METADATA retrieves the token's metadata.
func (token *SYN722Token) GET_SYN722_TOKEN_METADATA() SYN722Metadata {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Metadata
}

// UPDATE_SYN722_TOKEN_METADATA updates the token's metadata.
func (token *SYN722Token) UPDATE_SYN722_TOKEN_METADATA(newMetadata SYN722Metadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata = newMetadata
    token.UpdatedAt = time.Now()
    return token.Ledger.RecordLog("MetadataUpdate", fmt.Sprintf("Metadata updated for token %s", token.ID))
}

// SET_SYN722_TOKEN_URI updates the URI associated with the token's metadata.
func (token *SYN722Token) SET_SYN722_TOKEN_URI(uri string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TokenURI = uri
    return token.Ledger.RecordLog("TokenURIUpdate", fmt.Sprintf("Token URI set to %s for token %s", uri, token.ID))
}

// GET_SYN722_TOKEN_URI retrieves the token's URI.
func (token *SYN722Token) GET_SYN722_TOKEN_URI() string {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.TokenURI
}

// GET_SYN722_TOKEN_OWNER returns the owner of the token.
func (token *SYN722Token) GET_SYN722_TOKEN_OWNER() string {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Owner
}

// ENABLE_FUNGIBLE_MODE enables fungible mode for the token.
func (token *SYN722Token) ENABLE_FUNGIBLE_MODE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Mode = "fungible"
    return token.Ledger.RecordLog("ModeChange", fmt.Sprintf("Token %s switched to fungible mode", token.ID))
}

// ENABLE_NON_FUNGIBLE_MODE enables non-fungible mode for the token.
func (token *SYN722Token) ENABLE_NON_FUNGIBLE_MODE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Mode = "non-fungible"
    return token.Ledger.RecordLog("ModeChange", fmt.Sprintf("Token %s switched to non-fungible mode", token.ID))
}

// SWITCH_SYN722_TOKEN_MODE switches the token mode between fungible and non-fungible.
func (token *SYN722Token) SWITCH_SYN722_TOKEN_MODE(targetMode string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if targetMode != "fungible" && targetMode != "non-fungible" {
        return fmt.Errorf("invalid target mode: %s", targetMode)
    }

    if token.Mode == targetMode {
        return fmt.Errorf("token %s is already in %s mode", token.ID, targetMode)
    }

    token.Mode = targetMode
    token.UpdatedAt = time.Now()
    modeChangeLog := common.ModeChangeLog{
        TokenID:    token.ID,
        FromMode:   token.Mode,
        ToMode:     targetMode,
        ChangeDate: time.Now(),
    }
    token.ModeChangeHistory = append(token.ModeChangeHistory, modeChangeLog)
    return token.Ledger.RecordLog("ModeSwitch", fmt.Sprintf("Token %s switched to %s mode", token.ID, targetMode))
}

// LOG_MODE_CHANGE logs the change of mode in the token's history.
func (token *SYN722Token) LOG_MODE_CHANGE(fromMode, toMode string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    modeChangeLog := common.ModeChangeLog{
        TokenID:    token.ID,
        FromMode:   fromMode,
        ToMode:     toMode,
        ChangeDate: time.Now(),
    }
    token.ModeChangeHistory = append(token.ModeChangeHistory, modeChangeLog)
    return token.Ledger.RecordLog("ModeChange", fmt.Sprintf("Mode changed from %s to %s for token %s", fromMode, toMode, token.ID))
}

// CHECK_MODE_CHANGE_HISTORY returns the mode change history of the token.
func (token *SYN722Token) CHECK_MODE_CHANGE_HISTORY() []common.ModeChangeLog {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.ModeChangeHistory
}
