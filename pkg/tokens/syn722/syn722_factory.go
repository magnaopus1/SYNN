package syn722

import (
	"errors"
	"sync"
	"time"

)

// SYN722Metadata defines the metadata for the SYN722 token in both fungible and non-fungible states
type SYN722Metadata struct {
    ID              string                 `json:"id"`
    Name            string                 `json:"name"`
    Description     string                 `json:"description"`
    ImageURL        string                 `json:"image_url,omitempty"` // Used for non-fungible mode
    DocumentHash    string                 `json:"document_hash,omitempty"` // Used for non-fungible mode
    Properties      map[string]interface{} `json:"properties,omitempty"` // Additional metadata for non-fungible tokens
}

// SYN722Token represents a token that can switch between fungible and non-fungible states
type SYN722Token struct {
    ID                   string                   `json:"id"`
    Name                 string                   `json:"name"`
    Owner                string                   `json:"owner"`
    Mode                 string                   `json:"mode"` // "fungible" or "non-fungible"
    Quantity             uint64                   `json:"quantity,omitempty"` // Used only in fungible mode
    Metadata             SYN722Metadata            `json:"metadata,omitempty"` // Used in non-fungible mode
    RoyaltyInfo          common.RoyaltyInfo       `json:"royalty_info"`
    TransferHistory      []common.TransferRecord  `json:"transfer_history"`
    ModeChangeHistory    []common.ModeChangeLog   `json:"mode_change_history"`
    EncryptedData        string                   `json:"encrypted_data"`
    EncryptionKey        string                   `json:"encryption_key"`
    CreatedAt            time.Time                `json:"created_at"`
    UpdatedAt            time.Time                `json:"updated_at"`
}

// SYN722Factory manages the creation, issuance, and storage of SYN722 tokens
type SYN722Factory struct {
    Ledger            *ledger.Ledger                // Ledger for recording token transactions
    ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for validating token operations
    EncryptionService *encryption.EncryptionService // Encryption service for securing data
    mutex             sync.Mutex                    // Mutex for safe concurrent access
}

// NewSYN722Factory initializes a new SYN722Factory instance
func NewSYN722Factory(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN722Factory {
    return &SYN722Factory{
        Ledger:            ledger,
        ConsensusEngine:   consensusEngine,
        EncryptionService: encryptionService,
    }
}

// CreateToken creates a new SYN722 token in either fungible or non-fungible mode
func (sf *SYN722Factory) CreateToken(name, owner string, metadata SYN722Metadata, mode string, quantity uint64) (*SYN722Token, error) {
    sf.mutex.Lock()
    defer sf.mutex.Unlock()

    // Validate input based on the mode
    if mode == "fungible" && quantity == 0 {
        return nil, errors.New("quantity must be greater than zero in fungible mode")
    }
    if mode == "non-fungible" && (metadata.ID == "" || metadata.Name == "") {
        return nil, errors.New("metadata must contain a valid ID and Name in non-fungible mode")
    }

    // Generate a new token ID
    tokenID := generateTokenID(name, owner)

    // Encrypt metadata for security
    encryptedData, encryptionKey, err := sf.EncryptionService.EncryptData([]byte(metadata.Description))
    if err != nil {
        return nil, errors.New("failed to encrypt token data")
    }

    // Initialize the SYN722 token
    token := &SYN722Token{
        ID:            tokenID,
        Name:          name,
        Owner:         owner,
        Mode:          mode,
        Quantity:      quantity,
        Metadata:      metadata,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
        EncryptedData: encryptedData,
        EncryptionKey: encryptionKey,
    }

    // Validate token creation with Synnergy Consensus
    if err := sf.ConsensusEngine.ValidateTokenCreation(token); err != nil {
        return nil, errors.New("token creation validation failed via Synnergy Consensus")
    }

    // Record the token in the ledger
    if err := sf.Ledger.RecordTokenCreation(tokenID, token); err != nil {
        return nil, errors.New("failed to record token creation in the ledger")
    }

    return token, nil
}

// SwitchMode allows the token to switch between fungible and non-fungible modes
func (sf *SYN722Factory) SwitchMode(tokenID, newMode string) error {
    sf.mutex.Lock()
    defer sf.mutex.Unlock()

    // Retrieve the token from the ledger
    token, err := sf.Ledger.GetToken(tokenID)
    if err != nil {
        return errors.New("failed to retrieve token from ledger")
    }

    // Validate mode switch logic
    if token.Mode == newMode {
        return errors.New("token is already in the specified mode")
    }

    // Update the token's mode
    token.Mode = newMode
    token.ModeChangeHistory = append(token.ModeChangeHistory, common.ModeChangeLog{
        PreviousMode: token.Mode,
        NewMode:      newMode,
        Timestamp:    time.Now(),
    })
    token.UpdatedAt = time.Now()

    // Validate the mode switch through Synnergy Consensus
    if err := sf.ConsensusEngine.ValidateTokenUpdate(token); err != nil {
        return errors.New("token mode switch validation failed via Synnergy Consensus")
    }

    // Update the token in the ledger
    if err := sf.Ledger.UpdateToken(tokenID, token); err != nil {
        return errors.New("failed to update token in ledger")
    }

    return nil
}

// TransferTokenOwnership transfers the ownership of a SYN722 token to a new owner
func (sf *SYN722Factory) TransferTokenOwnership(tokenID, newOwner string) error {
    sf.mutex.Lock()
    defer sf.mutex.Unlock()

    // Retrieve the token from the ledger
    token, err := sf.Ledger.GetToken(tokenID)
    if err != nil {
        return errors.New("failed to retrieve token from ledger")
    }

    // Update the token's owner
    token.Owner = newOwner
    token.TransferHistory = append(token.TransferHistory, common.TransferRecord{
        From:      token.Owner,
        To:        newOwner,
        Timestamp: time.Now(),
    })
    token.UpdatedAt = time.Now()

    // Validate the transfer through Synnergy Consensus
    if err := sf.ConsensusEngine.ValidateTokenTransfer(token); err != nil {
        return errors.New("token transfer validation failed via Synnergy Consensus")
    }

    // Update the token in the ledger
    if err := sf.Ledger.UpdateToken(tokenID, token); err != nil {
        return errors.New("failed to update token in ledger")
    }

    return nil
}

// generateTokenID creates a unique token ID based on name and owner
func generateTokenID(name, owner string) string {
    return name + "_" + owner + "_" + time.Now().Format("20060102150405")
}
