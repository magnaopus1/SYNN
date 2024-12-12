package syn4300

import (
	"errors"
	"sync"
	"time"
)

// Syn4300Token represents an energy asset, REC, or carbon credit token.
type SYN4300Token struct {
	TokenID           string                `json:"token_id"`             // Unique token identifier
	Metadata          SYN4300Metadata       `json:"metadata"`             // Detailed metadata for the energy token
	TransactionHistory []TokenTransaction   `json:"transaction_history"`  // Historical transactions of the token
	CreationDate      time.Time             `json:"creation_date"`        // Date when the token was created
	LastModified      time.Time             `json:"last_modified"`        // Last modification date
	mutex             sync.Mutex            // Mutex for thread-safe operations
	ledgerService     *ledger.LedgerService // Integration with the ledger
	encryptionService *encryption.Encryptor // Integration with encryption
	consensusService  *consensus.SynnergyConsensus // Integration with Synnergy Consensus
}

// Syn4300Metadata defines the structure for energy tokens, RECs, or carbon credits.
type SYN4300Metadata struct {
	Name            string          `json:"name"`             // Name of the energy asset
	Symbol          string          `json:"symbol"`           // Symbol representing the token
	AssetType       string          `json:"asset_type"`       // Type of energy asset (e.g., REC, carbon credit, etc.)
	Owner           string          `json:"owner"`            // Current owner of the token
	IssuanceDate    time.Time       `json:"issuance_date"`    // Date of issuance
	Quantity        float64         `json:"quantity"`         // Quantity of energy represented by the token
	ValidUntil      *time.Time      `json:"valid_until"`      // Validity period (optional)
	Status          string          `json:"status"`           // Status of the token (e.g., active, expired, revoked)
	Location        string          `json:"location"`         // Location of the energy asset or carbon credit generation
	Certification   []Certification `json:"certification"`    // List of certifications for the renewable energy or carbon credits
	EnergyDetails   EnergyDetails   `json:"energy_details"`   // Additional details specific to the energy asset
	AssetLink       AssetLink       `json:"asset_link"`       // Link to the specific energy asset
}

// Certification holds detailed certification information for the energy asset.
type Certification struct {
	CertifyingBody   string    `json:"certifying_body"`   // Entity providing the certification
	CertificateID    string    `json:"certificate_id"`    // Unique ID of the certificate
	IssuedDate       time.Time `json:"issued_date"`       // Date when the certificate was issued
	ValidUntil       *time.Time `json:"valid_until"`      // Certificate validity period (optional)
	CertificationURL string    `json:"certification_url"` // URL to view or verify the certificate
}

// EnergyDetails captures technical and operational details specific to the energy token.
type EnergyDetails struct {
	EnergyType   string  `json:"energy_type"`   // Type of energy (e.g., solar, wind, hydro, etc.)
	Production   float64 `json:"production"`    // Total energy produced (e.g., in MWh)
	Unit         string  `json:"unit"`          // Unit of measurement (e.g., MWh, kWh)
	CarbonOffset float64 `json:"carbon_offset"` // Amount of carbon offset by the energy
}

// AssetLink links the energy token to a specific asset.
type AssetLink struct {
	AssetID           string    `json:"asset_id"`           // ID of the associated energy asset
	AssetDetails      string    `json:"asset_details"`      // Detailed description of the energy asset
	LinkDate          time.Time `json:"link_date"`          // Date when the token was linked to the asset
	VerificationStatus bool     `json:"verification_status"` // Status of asset verification (verified or pending)
	Location          string    `json:"location"`           // Physical or virtual location of the asset
}

// TokenTransaction represents a historical transaction associated with the token.
type TokenTransaction struct {
	TransactionID string    `json:"transaction_id"` // Unique ID of the transaction
	Timestamp     time.Time `json:"timestamp"`      // Time of the transaction
	Action        string    `json:"action"`         // Action taken (e.g., transferred, updated, etc.)
	From          string    `json:"from"`           // Previous owner or system entity
	To            string    `json:"to"`             // New owner or system entity
	Details       string    `json:"details"`        // Any additional details or notes for the transaction
}

// NewSyn4300Token creates a new Syn4300 token with the provided metadata.
func NewSyn4300Token(metadata Syn4300Metadata, ledger *ledger.LedgerService, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) (*Syn4300Token, error) {
	// Generate a unique token ID
	tokenID := generateUniqueTokenID()

	// Create a new Syn4300Token object
	token := &Syn4300Token{
		TokenID:           tokenID,
		Metadata:          metadata,
		TransactionHistory: []TokenTransaction{},
		CreationDate:      time.Now(),
		LastModified:      time.Now(),
		ledgerService:     ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}

	// Log the token creation event
	err := token.ledgerService.LogEvent("TokenCreated", time.Now(), tokenID)
	if err != nil {
		return nil, err
	}

	// Validate the new token using Synnergy Consensus
	err = token.consensusService.ValidateSubBlock(tokenID)
	if err != nil {
		return nil, err
	}

	// Store the token in the ledger
	err = token.storeTokenInLedger()
	if err != nil {
		return nil, err
	}

	return token, nil
}

// storeTokenInLedger stores the token details in the ledger.
func (t *Syn4300Token) storeTokenInLedger() error {
	// Encrypt the token data
	encryptedToken, err := t.encryptionService.EncryptData(t)
	if err != nil {
		return err
	}

	// Store the encrypted token in the ledger
	err = t.ledgerService.StoreToken(t.TokenID, encryptedToken)
	if err != nil {
		return err
	}

	return nil
}

// TransferToken transfers ownership of the token to a new owner.
func (t *Syn4300Token) TransferToken(newOwner string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Update the token's owner
	t.Metadata.Owner = newOwner
	t.LastModified = time.Now()

	// Record the transfer transaction
	transfer := TokenTransaction{
		TransactionID: generateUniqueTransactionID(),
		Timestamp:     time.Now(),
		Action:        "Transfer",
		From:          t.Metadata.Owner,
		To:            newOwner,
		Details:       "Ownership transfer",
	}
	t.TransactionHistory = append(t.TransactionHistory, transfer)

	// Log the transfer event in the ledger
	err := t.ledgerService.LogEvent("TokenTransferred", time.Now(), t.TokenID)
	if err != nil {
		return err
	}

	// Validate the transfer using Synnergy Consensus
	err = t.consensusService.ValidateSubBlock(t.TokenID)
	if err != nil {
		return err
	}

	// Update the token in the ledger
	err = t.storeTokenInLedger()
	if err != nil {
		return err
	}

	return nil
}


