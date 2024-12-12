package syn721

import (
    "sync"
    "fmt"
    "time"
    "path/to/ledger"
    "path/to/consensus"
    "path/to/encryption"
)

// SYN721Token struct defines the core attributes and operations of a SYN721 NFT.
type SYN721Token struct {
    mutex            sync.Mutex
    TokenID          string
    TokenURI         string
    Owner            string
    ApprovedAddress  string
    Ledger           *ledger.Ledger
    Consensus        *consensus.SynnergyConsensus
    Encryption       *encryption.Encryption
    Metadata         SYN721Metadata
    EscrowEnabled    bool
}

// TRANSFER_SYN721_TOKEN transfers ownership of the token to a new address.
func (token *SYN721Token) TRANSFER_SYN721_TOKEN(newOwner string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Owner != token.ApprovedAddress && token.Owner != newOwner {
        return fmt.Errorf("transfer not approved for new owner %s", newOwner)
    }

    oldOwner := token.Owner
    token.Owner = newOwner
    token.ApprovedAddress = "" // Clear any previous approval after transfer
    return token.Ledger.RecordLog("Transfer", fmt.Sprintf("Token %s transferred from %s to %s", token.TokenID, oldOwner, newOwner))
}

// APPROVE_SYN721_TOKEN_TRANSFER approves an address for transferring the token.
func (token *SYN721Token) APPROVE_SYN721_TOKEN_TRANSFER(address string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ApprovedAddress = address
    return token.Ledger.RecordLog("Approval", fmt.Sprintf("Token %s approved for transfer by %s", token.TokenID, address))
}

// CHECK_SYN721_TOKEN_BALANCE checks if the token has a defined owner.
func (token *SYN721Token) CHECK_SYN721_TOKEN_BALANCE() bool {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Owner != ""
}

// GET_SYN721_TOKEN_METADATA retrieves the metadata associated with the token.
func (token *SYN721Token) GET_SYN721_TOKEN_METADATA() SYN721Metadata {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Metadata
}

// UPDATE_SYN721_TOKEN_METADATA updates the token's metadata.
func (token *SYN721Token) UPDATE_SYN721_TOKEN_METADATA(newMetadata SYN721Metadata) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata = newMetadata
    return token.Ledger.RecordLog("MetadataUpdate", fmt.Sprintf("Metadata updated for token %s", token.TokenID))
}

// SET_SYN721_TOKEN_URI sets a new URI for the token's metadata.
func (token *SYN721Token) SET_SYN721_TOKEN_URI(newURI string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.TokenURI = newURI
    return token.Ledger.RecordLog("SetTokenURI", fmt.Sprintf("Token URI set for token %s", token.TokenID))
}

// GET_SYN721_TOKEN_URI retrieves the token's URI.
func (token *SYN721Token) GET_SYN721_TOKEN_URI() string {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.TokenURI
}

// GET_SYN721_TOKEN_OWNER retrieves the current owner of the token.
func (token *SYN721Token) GET_SYN721_TOKEN_OWNER() string {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.Owner
}

// SET_APPROVED_ADDRESS sets the approved address for token transfers.
func (token *SYN721Token) SET_APPROVED_ADDRESS(address string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.ApprovedAddress = address
    return token.Ledger.RecordLog("SetApprovedAddress", fmt.Sprintf("Approved address set to %s for token %s", address, token.TokenID))
}

// GET_APPROVED_ADDRESS retrieves the approved address for token transfers.
func (token *SYN721Token) GET_APPROVED_ADDRESS() string {
    token.mutex.RLock()
    defer token.mutex.RUnlock()

    return token.ApprovedAddress
}

// ENABLE_ESCROW_SERVICES enables escrow services for the token.
func (token *SYN721Token) ENABLE_ESCROW_SERVICES() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.EscrowEnabled = true
    return token.Ledger.RecordLog("EnableEscrow", fmt.Sprintf("Escrow enabled for token %s", token.TokenID))
}

// DISABLE_ESCROW_SERVICES disables escrow services for the token.
func (token *SYN721Token) DISABLE_ESCROW_SERVICES() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.EscrowEnabled = false
    return token.Ledger.RecordLog("DisableEscrow", fmt.Sprintf("Escrow disabled for token %s", token.TokenID))
}
