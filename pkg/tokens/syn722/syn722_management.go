package syn722

import (
	"errors"
	"sync"
	"time"

)

// SYN722Manager handles all management operations for SYN722 tokens, including creation, updates, burning, and auditing
type SYN722Manager struct {
	Ledger            *ledger.Ledger                // Ledger for recording token management actions
	ConsensusEngine   *consensus.SynnergyConsensus  // Synnergy Consensus for validating management operations
	EncryptionService *encryption.EncryptionService // Encryption service for securing management data
	mutex             sync.Mutex                    // Mutex for safe concurrent access
}

// NewSYN722Manager initializes a new SYN722Manager instance
func NewSYN722Manager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) *SYN722Manager {
	return &SYN722Manager{
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
	}
}

// CreateToken creates a new SYN722 token with either fungible or non-fungible mode and records it in the ledger
func (tm *SYN722Manager) CreateToken(name, owner, mode string, quantity uint64, metadata SYN722Metadata) (*SYN722Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

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
	encryptedData, encryptionKey, err := tm.EncryptionService.EncryptData([]byte(metadata.Description))
	if err != nil {
		return nil, errors.New("failed to encrypt token metadata")
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
	if err := tm.ConsensusEngine.ValidateTokenCreation(token); err != nil {
		return nil, errors.New("token creation validation failed via Synnergy Consensus")
	}

	// Record the token in the ledger
	if err := tm.Ledger.RecordTokenCreation(tokenID, token); err != nil {
		return nil, errors.New("failed to record token creation in the ledger")
	}

	return token, nil
}

// UpdateToken updates an existing SYN722 token's metadata or state, ensuring compliance and recording the change in the ledger
func (tm *SYN722Manager) UpdateToken(tokenID, newOwner string, metadata SYN722Metadata, mode string, quantity uint64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger")
	}

	// Update token details based on the new inputs
	token.Owner = newOwner
	token.Mode = mode
	token.Quantity = quantity
	token.Metadata = metadata
	token.UpdatedAt = time.Now()

	// Validate the token update through Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateTokenUpdate(token); err != nil {
		return errors.New("token update validation failed via Synnergy Consensus")
	}

	// Record the token update in the ledger
	if err := tm.Ledger.UpdateToken(tokenID, token); err != nil {
		return errors.New("failed to update token in ledger")
	}

	return nil
}

// BurnToken destroys a specific quantity of a fungible SYN722 token or completely burns a non-fungible SYN722 token, removing it from circulation
func (tm *SYN722Manager) BurnToken(tokenID string, quantity uint64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger")
	}

	// Ensure proper logic for burning based on token mode
	if token.Mode == "fungible" {
		if token.Quantity < quantity {
			return errors.New("insufficient quantity to burn")
		}
		token.Quantity -= quantity
	} else if token.Mode == "non-fungible" {
		token.Quantity = 0 // Non-fungible tokens are fully burned
	}

	// Validate token burning with Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateTokenBurning(token); err != nil {
		return errors.New("token burning validation failed via Synnergy Consensus")
	}

	// Record the burning event in the ledger
	if err := tm.Ledger.BurnToken(tokenID, token); err != nil {
		return errors.New("failed to burn token in ledger")
	}

	return nil
}

// AuditToken retrieves a SYN722 token and its entire history for auditing purposes
func (tm *SYN722Manager) AuditToken(tokenID string) (*SYN722Token, []*common.TransferRecord, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, nil, errors.New("failed to retrieve token for auditing")
	}

	// Retrieve the token's transfer history from the ledger
	history, err := tm.Ledger.GetTokenTransferHistory(tokenID)
	if err != nil {
		return nil, nil, errors.New("failed to retrieve token transfer history")
	}

	return token, history, nil
}

// generateTokenID creates a unique token ID based on name and owner
func generateTokenID(name, owner string) string {
	return name + "_" + owner + "_" + time.Now().Format("20060102150405")
}
