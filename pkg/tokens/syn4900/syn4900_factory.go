package syn4900

import (
	"errors"
	"sync"
	"time"
)

// Syn4900Token represents an agricultural asset in the form of a blockchain token.
type Syn4900Token struct {
	TokenID         string          `json:"token_id"`
	Metadata        Syn4900Metadata `json:"metadata"`
	mutex           sync.Mutex      // Mutex for thread-safe operations
	ledgerService   *ledger.Ledger  // Ledger for logging token operations
	encryptionService *encryption.Encryptor // Encryption service for securing token data
	consensusService *consensus.SynnergyConsensus // Consensus service for validating token creation and updates
}

// Syn4900Metadata holds the core metadata for the agricultural token.
type Syn4900Metadata struct {
	Name             string    `json:"name"`			//e.g 1kg Corn
	Symbol           string    `json:"symbol"`		
	Value            float64   `json:"value"`          // Can be updated based on market conditions
	Location         string    `json:"location"`       // Where the asset is stored or produced
	AssetType        string    `json:"asset_type"`     // Type of asset, e.g., grain, livestock, etc.
	Quantity         float64   `json:"quantity"`       // Quantity of the asset
	Owner            string    `json:"owner"`          // Owner's ID (e.g., wallet address)
	Origin           string    `json:"origin"`         // Where the asset was harvested or produced
	HarvestDate      time.Time `json:"harvest_date"`   // Date of harvest or production
	ExpiryDate       time.Time `json:"expiry_date"`    // Date after which the asset is no longer viable
	Status           string    `json:"status"`         // Status of the asset, e.g., "Available", "Expired", etc.
	Certification    string    `json:"certification"`  // Certification details of the asset, such as organic or fair trade certifications
}

// AssetLink represents the linkage between an agricultural token and a real-world agricultural asset.
type AssetLink struct {
	TokenID          string    `json:"token_id"`
	AssetID          string    `json:"asset_id"`
	AssetDetails     string    `json:"asset_details"`
	LinkDate         time.Time `json:"link_date"`
	VerificationStatus bool    `json:"verification_status"`
}

// TokenManager manages the lifecycle of SYN4900 tokens.
type TokenManager struct {
	tokens           map[string]*Syn4900Token    // In-memory storage for tokens
	mutex            sync.Mutex                  // Mutex for thread-safe operations
	ledgerService    *ledger.Ledger              // Ledger integration for tracking token events
	encryptionService *encryption.Encryptor       // Encryption service for securing token data
	consensusService *consensus.SynnergyConsensus // Consensus for validating token operations
}

// NewTokenManager creates a new instance of TokenManager.
func NewTokenManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TokenManager {
	return &TokenManager{
		tokens:           make(map[string]*Syn4900Token),
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// CreateToken creates a new agricultural token and stores it in the blockchain.
func (tm *TokenManager) CreateToken(metadata Syn4900Metadata) (*Syn4900Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a unique token ID
	tokenID := generateUniqueTokenID()

	// Create the token
	token := &Syn4900Token{
		TokenID: tokenID,
		Metadata: metadata,
		ledgerService: tm.ledgerService,
		encryptionService: tm.encryptionService,
		consensusService: tm.consensusService,
	}

	// Encrypt the token data
	encryptedToken, err := tm.encryptionService.EncryptData(token)
	if err != nil {
		return nil, err
	}

	// Store the token in memory
	tm.tokens[tokenID] = encryptedToken.(*Syn4900Token)

	// Log the token creation in the ledger
	if err := tm.ledgerService.LogEvent("TokenCreated", time.Now(), tokenID); err != nil {
		return nil, err
	}

	// Validate the token creation using Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, err
	}

	return token, nil
}

// GetToken retrieves a token by its ID.
func (tm *TokenManager) GetToken(tokenID string) (*Syn4900Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from in-memory storage
	encryptedToken, exists := tm.tokens[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Decrypt the token data
	decryptedToken, err := tm.encryptionService.DecryptData(encryptedToken)
	if err != nil {
		return nil, err
	}

	return decryptedToken.(*Syn4900Token), nil
}

// UpdateTokenValue updates the value of an existing agricultural token.
func (tm *TokenManager) UpdateTokenValue(tokenID string, newValue float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token to update
	token, err := tm.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Update the value of the token
	token.Metadata.Value = newValue

	// Encrypt the updated token
	encryptedToken, err := tm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Store the updated token
	tm.tokens[tokenID] = encryptedToken.(*Syn4900Token)

	// Log the value update in the ledger
	if err := tm.ledgerService.LogEvent("TokenValueUpdated", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the update using Synnergy Consensus
	return tm.consensusService.ValidateSubBlock(tokenID)
}

// LinkRealAsset links a real-world agricultural asset to the token.
func (tm *TokenManager) LinkRealAsset(tokenID, assetID, assetDetails string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Create the asset link
	link := &AssetLink{
		TokenID:          tokenID,
		AssetID:          assetID,
		AssetDetails:     assetDetails,
		LinkDate:         time.Now(),
		VerificationStatus: false, // Default to false until verified
	}

	// Log the asset link in the ledger
	if err := tm.ledgerService.LogEvent("AssetLinked", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the asset link using consensus
	return tm.consensusService.ValidateSubBlock(tokenID)
}
package syn4900

import (
	"errors"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/consensus"
)

// Syn4900Token represents an agricultural asset in the form of a blockchain token.
type Syn4900Token struct {
	TokenID         string          `json:"token_id"`
	Metadata        Syn4900Metadata `json:"metadata"`
	mutex           sync.Mutex      // Mutex for thread-safe operations
	ledgerService   *ledger.Ledger  // Ledger for logging token operations
	encryptionService *encryption.Encryptor // Encryption service for securing token data
	consensusService *consensus.SynnergyConsensus // Consensus service for validating token creation and updates
}

// Syn4900Metadata holds the core metadata for the agricultural token.
type Syn4900Metadata struct {
	Name             string    `json:"name"`
	Symbol           string    `json:"symbol"`
	Value            float64   `json:"value"`          // Can be updated based on market conditions
	Location         string    `json:"location"`       // Where the asset is stored or produced
	AssetType        string    `json:"asset_type"`     // Type of asset, e.g., grain, livestock, etc.
	Quantity         float64   `json:"quantity"`       // Quantity of the asset
	Owner            string    `json:"owner"`          // Owner's ID (e.g., wallet address)
	Origin           string    `json:"origin"`         // Where the asset was harvested or produced
	HarvestDate      time.Time `json:"harvest_date"`   // Date of harvest or production
	ExpiryDate       time.Time `json:"expiry_date"`    // Date after which the asset is no longer viable
	Status           string    `json:"status"`         // Status of the asset, e.g., "Available", "Expired", etc.
	Certification    string    `json:"certification"`  // Certification details of the asset, such as organic or fair trade certifications
}

// AssetLink represents the linkage between an agricultural token and a real-world agricultural asset.
type AssetLink struct {
	TokenID          string    `json:"token_id"`
	AssetID          string    `json:"asset_id"`
	AssetDetails     string    `json:"asset_details"`
	LinkDate         time.Time `json:"link_date"`
	VerificationStatus bool    `json:"verification_status"`
}

// TokenManager manages the lifecycle of SYN4900 tokens.
type TokenManager struct {
	tokens           map[string]*Syn4900Token    // In-memory storage for tokens
	assetLinks       map[string]*AssetLink       // Asset links for real-world asset tracking
	mutex            sync.Mutex                  // Mutex for thread-safe operations
	ledgerService    *ledger.Ledger              // Ledger integration for tracking token events
	encryptionService *encryption.Encryptor       // Encryption service for securing token data
	consensusService *consensus.SynnergyConsensus // Consensus for validating token operations
}

// NewTokenManager creates a new instance of TokenManager.
func NewTokenManager(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TokenManager {
	return &TokenManager{
		tokens:           make(map[string]*Syn4900Token),
		assetLinks:       make(map[string]*AssetLink),
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// CreateToken creates a new agricultural token and stores it in the blockchain.
func (tm *TokenManager) CreateToken(metadata Syn4900Metadata) (*Syn4900Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Generate a unique token ID
	tokenID := generateUniqueTokenID()

	// Create the token
	token := &Syn4900Token{
		TokenID: tokenID,
		Metadata: metadata,
		ledgerService: tm.ledgerService,
		encryptionService: tm.encryptionService,
		consensusService: tm.consensusService,
	}

	// Encrypt the token data
	encryptedToken, err := tm.encryptionService.EncryptData(token)
	if err != nil {
		return nil, err
	}

	// Store the token in memory
	tm.tokens[tokenID] = encryptedToken.(*Syn4900Token)

	// Log the token creation in the ledger
	if err := tm.ledgerService.LogEvent("TokenCreated", time.Now(), tokenID); err != nil {
		return nil, err
	}

	// Validate the token creation using Synnergy Consensus
	if err := tm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, err
	}

	return token, nil
}

// GetToken retrieves a token by its ID.
func (tm *TokenManager) GetToken(tokenID string) (*Syn4900Token, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token from in-memory storage
	encryptedToken, exists := tm.tokens[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Decrypt the token data
	decryptedToken, err := tm.encryptionService.DecryptData(encryptedToken)
	if err != nil {
		return nil, err
	}

	return decryptedToken.(*Syn4900Token), nil
}

// UpdateTokenValue updates the value of an existing agricultural token.
func (tm *TokenManager) UpdateTokenValue(tokenID string, newValue float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the token to update
	token, err := tm.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Update the value of the token
	token.Metadata.Value = newValue

	// Encrypt the updated token
	encryptedToken, err := tm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Store the updated token
	tm.tokens[tokenID] = encryptedToken.(*Syn4900Token)

	// Log the value update in the ledger
	if err := tm.ledgerService.LogEvent("TokenValueUpdated", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the update using Synnergy Consensus
	return tm.consensusService.ValidateSubBlock(tokenID)
}

// LinkRealAsset links a real-world agricultural asset to the token.
func (tm *TokenManager) LinkRealAsset(tokenID, assetID, assetDetails string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Create the asset link
	link := &AssetLink{
		TokenID:          tokenID,
		AssetID:          assetID,
		AssetDetails:     assetDetails,
		LinkDate:         time.Now(),
		VerificationStatus: false, // Default to false until verified
	}

	// Store the asset link
	tm.assetLinks[tokenID] = link

	// Log the asset link in the ledger
	if err := tm.ledgerService.LogEvent("AssetLinked", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the asset link using consensus
	return tm.consensusService.ValidateSubBlock(tokenID)
}

// VerifyAssetLink verifies the link between the token and the real-world asset.
func (tm *TokenManager) VerifyAssetLink(tokenID, assetID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Retrieve the asset link
	link, exists := tm.assetLinks[tokenID]
	if !exists || link.AssetID != assetID {
		return errors.New("asset link not found or mismatch")
	}

	// Verify the asset link
	link.VerificationStatus = true

	// Encrypt the updated asset link (optional depending on how it's stored)
	encryptedLink, err := tm.encryptionService.EncryptData(link)
	if err != nil {
		return err
	}
	tm.assetLinks[tokenID] = encryptedLink.(*AssetLink)

	// Log the asset verification
	if err := tm.ledgerService.LogEvent("AssetVerified", time.Now(), tokenID); err != nil {
		return err
	}

	// Validate the asset verification using consensus
	return tm.consensusService.ValidateSubBlock(tokenID)
}

// generateUniqueTokenID generates a unique identifier for a new token.
func generateUniqueTokenID() string {
	// Generate a timestamp-based unique ID, or use UUID
	return time.Now().Format("20060102150405") // e.g., "20230918104523"
}
