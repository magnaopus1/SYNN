package syn4700

import (
	"errors"
	"sync"
	"time"

)

// Syn4700Token represents a legal document in the form of a blockchain token.
type Syn4700Token struct {
	TokenID           string            `json:"token_id"`           // Unique token ID
	Metadata          Syn4700Metadata   `json:"metadata"`           // Detailed metadata related to the legal token
	mutex             sync.Mutex        // Mutex for thread-safe operations
}

// Syn4700Metadata captures detailed legal document metadata.
type Syn4700Metadata struct {
	ContractTitle    string             `json:"contract_title"`      // Title of the contract or legal document
	DocumentType     string             `json:"document_type"`       // Type of legal document (e.g., contract, agreement)
	PartiesInvolved  []string           `json:"parties_involved"`    // List of parties involved in the legal document
	ContentHash      string             `json:"content_hash"`        // Hash of the document content for integrity
	CreationDate     time.Time          `json:"creation_date"`       // Date when the document was created
	ExpiryDate       *time.Time         `json:"expiry_date"`         // Expiry date of the document (optional)
	Status           string             `json:"status"`              // Status of the token (e.g., active, expired, terminated)
	Signatures       map[string]string  `json:"signatures"`          // Map of signatures (party ID -> signature hash)
	ContractVersion  string             `json:"contract_version"`    // Version of the contract
	DocumentCopy     []byte             `json:"document_copy"`       // Copy of the legal document as a byte array (encrypted)
}

// LegalTokenManager handles operations for managing SYN4700 legal tokens.
type LegalTokenManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
}

// NewLegalTokenManager initializes a new LegalTokenManager.
func NewLegalTokenManager(ledgerService *ledger.LedgerService, encryptionService *encryption.Encryptor, consensusService *consensus.SynnergyConsensus) *LegalTokenManager {
	return &LegalTokenManager{
		ledgerService:     ledgerService,
		encryptionService: encryptionService,
		consensusService:  consensusService,
	}
}

// CreateLegalToken creates a new SYN4700 legal token and stores it in the ledger.
func (ltm *LegalTokenManager) CreateLegalToken(title, documentType string, partiesInvolved []string, contentHash, contractVersion string, expiryDate *time.Time, documentCopy []byte) (*Syn4700Token, error) {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	// Generate a unique token ID.
	tokenID := generateUniqueTokenID()

	// Create the metadata for the token.
	metadata := Syn4700Metadata{
		ContractTitle:   title,
		DocumentType:    documentType,
		PartiesInvolved: partiesInvolved,
		ContentHash:     contentHash,
		CreationDate:    time.Now(),
		ExpiryDate:      expiryDate,
		Status:          "active", // Default status is active
		Signatures:      make(map[string]string),
		ContractVersion: contractVersion,
		DocumentCopy:    documentCopy,
	}

	// Create the token.
	token := &Syn4700Token{
		TokenID:  tokenID,
		Metadata: metadata,
	}

	// Encrypt the token before storing.
	encryptedToken, err := ltm.encryptionService.EncryptData(token)
	if err != nil {
		return nil, err
	}

	// Log the token creation in the ledger.
	if err := ltm.ledgerService.LogEvent("LegalTokenCreated", time.Now(), tokenID); err != nil {
		return nil, err
	}

	// Store the token in the ledger.
	if err := ltm.ledgerService.StoreToken(tokenID, encryptedToken); err != nil {
		return nil, err
	}

	// Validate the transaction using the consensus mechanism.
	if err := ltm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, err
	}

	return token, nil
}

// AddSignature adds a signature to the SYN4700 token.
func (ltm *LegalTokenManager) AddSignature(tokenID, partyID, signatureHash string) error {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	// Retrieve the token from the ledger.
	token, err := ltm.retrieveToken(tokenID)
	if err != nil {
		return err
	}

	// Add the signature to the token.
	token.Metadata.Signatures[partyID] = signatureHash

	// Encrypt the updated token.
	encryptedToken, err := ltm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Log the signature addition in the ledger.
	if err := ltm.ledgerService.LogEvent("SignatureAdded", time.Now(), tokenID); err != nil {
		return err
	}

	// Store the updated token in the ledger.
	if err := ltm.ledgerService.StoreToken(tokenID, encryptedToken); err != nil {
		return err
	}

	// Validate the update with the consensus mechanism.
	return ltm.consensusService.ValidateSubBlock(tokenID)
}

// Retrieve the token from the ledger and decrypt it.
func (ltm *LegalTokenManager) retrieveToken(tokenID string) (*Syn4700Token, error) {
	// Retrieve the encrypted token from the ledger.
	encryptedData, err := ltm.ledgerService.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}

	// Decrypt the token data.
	decryptedToken, err := ltm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedToken.(*Syn4700Token), nil
}

// ExpireLegalToken updates the status of the legal token to expired.
func (ltm *LegalTokenManager) ExpireLegalToken(tokenID string) error {
	ltm.mutex.Lock()
	defer ltm.mutex.Unlock()

	// Retrieve the token from the ledger.
	token, err := ltm.retrieveToken(tokenID)
	if err != nil {
		return err
	}

	// Update the status to expired.
	token.Metadata.Status = "expired"

	// Encrypt the updated token.
	encryptedToken, err := ltm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}

	// Log the token expiration in the ledger.
	if err := ltm.ledgerService.LogEvent("LegalTokenExpired", time.Now(), tokenID); err != nil {
		return err
	}

	// Store the updated token in the ledger.
	if err := ltm.ledgerService.StoreToken(tokenID, encryptedToken); err != nil {
		return err
	}

	// Validate the update using the consensus mechanism.
	return ltm.consensusService.ValidateSubBlock(tokenID)
}

// generateUniqueTokenID generates a unique identifier for a new token.
func generateUniqueTokenID() string {
	return "syn4700-token-" + time.Now().Format("20060102150405")
}
